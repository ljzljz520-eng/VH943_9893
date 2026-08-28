package storage

import (
	"context"
	"database/sql"

	"campgear/internal/domain"
)

func (r *Repository) InsertItem(ctx context.Context, item domain.InventoryItem) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO inventory_items (id, sku, name, category, daily_rate, deposit, stock, available, maintenance_status, storage_bin, listed, version) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`, item.ID, item.SKU, item.Name, item.Category, item.DailyRate, item.Deposit, item.Stock, item.Available, item.MaintenanceStatus, item.StorageBin, boolInt(item.Listed), item.Version)
	return err
}

func (r *Repository) UpdateItem(ctx context.Context, item domain.InventoryItem, expectedVersion int) error {
	result, err := r.db.ExecContext(ctx, `UPDATE inventory_items SET sku=?, name=?, category=?, daily_rate=?, deposit=?, stock=?, available=?, maintenance_status=?, storage_bin=?, listed=?, version=? WHERE id=? AND version=?`, item.SKU, item.Name, item.Category, item.DailyRate, item.Deposit, item.Stock, item.Available, item.MaintenanceStatus, item.StorageBin, boolInt(item.Listed), item.Version, item.ID, expectedVersion)
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

func scanItem(scanner interface{ Scan(...any) error }) (domain.InventoryItem, error) {
	var item domain.InventoryItem
	var listed int64
	if err := scanner.Scan(&item.ID, &item.SKU, &item.Name, &item.Category, &item.DailyRate, &item.Deposit, &item.Stock, &item.Available, &item.MaintenanceStatus, &item.StorageBin, &listed, &item.Version); err != nil {
		return domain.InventoryItem{}, wrapNotFound(err)
	}
	item.Listed = intBool(listed)
	return item, nil
}

func (r *Repository) GetItem(ctx context.Context, id string) (domain.InventoryItem, error) {
	return scanItem(r.db.QueryRowContext(ctx, `SELECT id, sku, name, category, daily_rate, deposit, stock, available, maintenance_status, storage_bin, listed, version FROM inventory_items WHERE id=?`, id))
}

func (r *Repository) ListItems(ctx context.Context, category domain.Category, onlyListed bool) ([]domain.InventoryItem, error) {
	query := `SELECT id, sku, name, category, daily_rate, deposit, stock, available, maintenance_status, storage_bin, listed, version FROM inventory_items WHERE (? = '' OR category=?) AND (? = 0 OR listed=1) ORDER BY category, name`
	rows, err := r.db.QueryContext(ctx, query, category, category, boolInt(onlyListed))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.InventoryItem, 0)
	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ChangeAvailability(ctx context.Context, tx *sql.Tx, id string, delta int) error {
	result, err := tx.ExecContext(ctx, `UPDATE inventory_items SET available=available+?, listed=CASE WHEN maintenance_status='ready' AND available+? > 0 THEN 1 ELSE 0 END, version=version+1 WHERE id=? AND available+? >= 0`, delta, delta, id, delta)
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

func (r *Repository) SetMaintenanceState(ctx context.Context, tx *sql.Tx, id string, status domain.MaintenanceStatus, listed bool) error {
	result, err := tx.ExecContext(ctx, `UPDATE inventory_items SET maintenance_status=?, listed=?, version=version+1 WHERE id=?`, status, boolInt(listed), id)
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
