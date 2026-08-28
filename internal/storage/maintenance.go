package storage

import (
	"context"
	"database/sql"

	"campgear/internal/domain"
)

func (r *Repository) InsertMaintenance(ctx context.Context, order domain.MaintenanceOrder) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO maintenance_orders (id, item_id, reason, opened_date, closed_date, status, technician) VALUES (?, ?, ?, ?, ?, ?, ?)`, order.ID, order.ItemID, order.Reason, order.OpenedDate, order.ClosedDate, order.Status, order.Technician)
	return err
}

func (r *Repository) InsertMaintenanceTx(ctx context.Context, tx *sql.Tx, order domain.MaintenanceOrder) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO maintenance_orders (id, item_id, reason, opened_date, closed_date, status, technician) VALUES (?, ?, ?, ?, ?, ?, ?)`, order.ID, order.ItemID, order.Reason, order.OpenedDate, order.ClosedDate, order.Status, order.Technician)
	return err
}

func scanMaintenance(scanner interface{ Scan(...any) error }) (domain.MaintenanceOrder, error) {
	var order domain.MaintenanceOrder
	err := scanner.Scan(&order.ID, &order.ItemID, &order.Reason, &order.OpenedDate, &order.ClosedDate, &order.Status, &order.Technician)
	if err != nil {
		return domain.MaintenanceOrder{}, wrapNotFound(err)
	}
	return order, nil
}

func (r *Repository) GetMaintenance(ctx context.Context, id string) (domain.MaintenanceOrder, error) {
	return scanMaintenance(r.db.QueryRowContext(ctx, `SELECT id, item_id, reason, opened_date, closed_date, status, technician FROM maintenance_orders WHERE id=?`, id))
}

func (r *Repository) ListMaintenance(ctx context.Context, status string) ([]domain.MaintenanceOrder, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, item_id, reason, opened_date, closed_date, status, technician FROM maintenance_orders WHERE (?='' OR status=?) ORDER BY opened_date, id`, status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.MaintenanceOrder, 0)
	for rows.Next() {
		order, err := scanMaintenance(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, order)
	}
	return result, rows.Err()
}

func (r *Repository) UpdateMaintenance(ctx context.Context, tx *sql.Tx, order domain.MaintenanceOrder) error {
	result, err := tx.ExecContext(ctx, `UPDATE maintenance_orders SET reason=?, opened_date=?, closed_date=?, status=?, technician=? WHERE id=?`, order.Reason, order.OpenedDate, order.ClosedDate, order.Status, order.Technician, order.ID)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return sql.ErrNoRows
	}
	return nil
}
