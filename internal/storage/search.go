package storage

import (
	"context"
	"strings"

	"campgear/internal/domain"
)

func (r *Repository) SearchItems(ctx context.Context, query string, category domain.Category, listedOnly bool) ([]domain.InventoryItem, error) {
	query = strings.TrimSpace(query)
	items, err := r.ListItems(ctx, category, listedOnly)
	if err != nil {
		return nil, err
	}
	if query == "" {
		return items, nil
	}
	result := make([]domain.InventoryItem, 0, len(items))
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Name), strings.ToLower(query)) || strings.Contains(strings.ToLower(item.SKU), strings.ToLower(query)) {
			result = append(result, item)
		}
	}
	return result, nil
}

func (r *Repository) LoadRentalLines(ctx context.Context, rentalID string) ([]domain.RentalLine, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT item_id, quantity, days, rate, deposit, subtotal FROM rental_lines WHERE rental_id=? ORDER BY item_id`, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	lines := make([]domain.RentalLine, 0)
	for rows.Next() {
		var line domain.RentalLine
		if err := rows.Scan(&line.ItemID, &line.Quantity, &line.Days, &line.Rate, &line.Deposit, &line.Subtotal); err != nil {
			return nil, err
		}
		lines = append(lines, line)
	}
	return lines, rows.Err()
}
