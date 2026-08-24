package storage

import (
	"context"
	"database/sql"

	"campgear/internal/domain"
)

func (r *Repository) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *Repository) InsertRental(ctx context.Context, tx *sql.Tx, record domain.RentalRecord) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO rental_records (id, customer, start_date, end_date, status, total, deposit_held) VALUES (?, ?, ?, ?, ?, ?, ?)`, record.ID, record.Customer, record.StartDate, record.EndDate, record.Status, record.Total, record.DepositHeld)
	return err
}

func (r *Repository) InsertRentalLine(ctx context.Context, tx *sql.Tx, rentalID string, line domain.RentalLine) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO rental_lines (rental_id, item_id, quantity, days, rate, deposit, subtotal) VALUES (?, ?, ?, ?, ?, ?, ?)`, rentalID, line.ItemID, line.Quantity, line.Days, line.Rate, line.Deposit, line.Subtotal)
	return err
}

func (r *Repository) GetRental(ctx context.Context, id string) (domain.RentalRecord, error) {
	var record domain.RentalRecord
	err := r.db.QueryRowContext(ctx, `SELECT id, customer, start_date, end_date, status, total, deposit_held FROM rental_records WHERE id=?`, id).Scan(&record.ID, &record.Customer, &record.StartDate, &record.EndDate, &record.Status, &record.Total, &record.DepositHeld)
	if err != nil {
		return domain.RentalRecord{}, wrapNotFound(err)
	}
	rows, err := r.db.QueryContext(ctx, `SELECT item_id, quantity, days, rate, deposit, subtotal FROM rental_lines WHERE rental_id=? ORDER BY item_id`, id)
	if err != nil {
		return domain.RentalRecord{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var line domain.RentalLine
		if err := rows.Scan(&line.ItemID, &line.Quantity, &line.Days, &line.Rate, &line.Deposit, &line.Subtotal); err != nil {
			return domain.RentalRecord{}, err
		}
		record.Lines = append(record.Lines, line)
	}
	return record, rows.Err()
}

func (r *Repository) ListRentals(ctx context.Context, status string) ([]domain.RentalRecord, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, customer, start_date, end_date, status, total, deposit_held FROM rental_records WHERE (?='' OR status=?) ORDER BY start_date, id`, status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.RentalRecord, 0)
	for rows.Next() {
		var record domain.RentalRecord
		if err := rows.Scan(&record.ID, &record.Customer, &record.StartDate, &record.EndDate, &record.Status, &record.Total, &record.DepositHeld); err != nil {
			return nil, err
		}
		result = append(result, record)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	for index := range result {
		lines, err := r.LoadRentalLines(ctx, result[index].ID)
		if err != nil {
			return nil, err
		}
		result[index].Lines = lines
	}
	return result, nil
}

func (r *Repository) UpdateRentalStatus(ctx context.Context, id, status string) error {
	result, err := r.db.ExecContext(ctx, `UPDATE rental_records SET status=? WHERE id=?`, status, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrNotFound
	}
	return nil
}

func (r *Repository) UpdateRentalStatusTx(ctx context.Context, tx *sql.Tx, id, status string) error {
	result, err := tx.ExecContext(ctx, `UPDATE rental_records SET status=? WHERE id=?`, status, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return ErrNotFound
	}
	return nil
}
