package storage

import (
	"context"
	"database/sql"

	"campgear/internal/domain"
)

func (r *Repository) InsertAudit(ctx context.Context, event domain.AuditEvent) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO audit_events (id, entity_type, entity_id, action, actor, occurred_at, details) VALUES (?, ?, ?, ?, ?, ?, ?)`, event.ID, event.EntityType, event.EntityID, event.Action, event.Actor, event.OccurredAt, event.Details)
	return err
}

func (r *Repository) InsertAuditTx(ctx context.Context, tx *sql.Tx, event domain.AuditEvent) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO audit_events (id, entity_type, entity_id, action, actor, occurred_at, details) VALUES (?, ?, ?, ?, ?, ?, ?)`, event.ID, event.EntityType, event.EntityID, event.Action, event.Actor, event.OccurredAt, event.Details)
	return err
}

func (r *Repository) ListAudit(ctx context.Context, entityType, entityID string) ([]domain.AuditEvent, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, entity_type, entity_id, action, actor, occurred_at, details FROM audit_events WHERE (?='' OR entity_type=?) AND (?='' OR entity_id=?) ORDER BY occurred_at, id`, entityType, entityType, entityID, entityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]domain.AuditEvent, 0)
	for rows.Next() {
		var event domain.AuditEvent
		if err := rows.Scan(&event.ID, &event.EntityType, &event.EntityID, &event.Action, &event.Actor, &event.OccurredAt, &event.Details); err != nil {
			return nil, err
		}
		result = append(result, event)
	}
	return result, rows.Err()
}

func (r *Repository) Count(ctx context.Context, table string) (int, error) {
	allowed := map[string]string{"inventory_items": "inventory_items", "rental_records": "rental_records", "maintenance_orders": "maintenance_orders", "audit_events": "audit_events"}
	name, ok := allowed[table]
	if !ok {
		return 0, ErrNotFound
	}
	var count int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM "+name).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}
