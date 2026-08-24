package reporting

import (
	"context"
	"fmt"
	"sort"

	"campgear/internal/domain"
	"campgear/internal/storage"
)

type Service struct{ repo *storage.Repository }

func NewService(repo *storage.Repository) *Service { return &Service{repo: repo} }

type InventorySummary struct {
	TotalItems       int `json:"total_items"`
	ListedItems      int `json:"listed_items"`
	MaintenanceItems int `json:"maintenance_items"`
	AvailableUnits   int `json:"available_units"`
	StockUnits       int `json:"stock_units"`
}

func (s *Service) Inventory(ctx context.Context) (InventorySummary, error) {
	items, err := s.repo.ListItems(ctx, "", false)
	if err != nil {
		return InventorySummary{}, err
	}
	summary := InventorySummary{}
	for _, item := range items {
		summary.TotalItems++
		summary.StockUnits += item.Stock
		summary.AvailableUnits += item.Available
		if item.Listed {
			summary.ListedItems++
		}
		if item.MaintenanceStatus == domain.StatusMaintenance {
			summary.MaintenanceItems++
		}
	}
	return summary, nil
}

type RevenueReport struct {
	RentalCount int                       `json:"rental_count"`
	Gross       int64                     `json:"gross"`
	Deposits    int64                     `json:"deposits"`
	ByCategory  map[domain.Category]int64 `json:"by_category"`
}

func (s *Service) Revenue(ctx context.Context, status string, categoryOf func(string) (domain.Category, error)) (RevenueReport, error) {
	rentals, err := s.repo.ListRentals(ctx, status)
	if err != nil {
		return RevenueReport{}, err
	}
	report := RevenueReport{ByCategory: make(map[domain.Category]int64)}
	for _, rental := range rentals {
		report.RentalCount++
		report.Gross += rental.Total
		report.Deposits += rental.DepositHeld
		for _, line := range rental.Lines {
			category, err := categoryOf(line.ItemID)
			if err != nil {
				return RevenueReport{}, err
			}
			report.ByCategory[category] += line.Subtotal
		}
	}
	return report, nil
}

func (s *Service) Audit(ctx context.Context, entityType, entityID string) ([]domain.AuditEvent, error) {
	events, err := s.repo.ListAudit(ctx, entityType, entityID)
	if err != nil {
		return nil, err
	}
	sort.SliceStable(events, func(a, b int) bool {
		if events[a].OccurredAt == events[b].OccurredAt {
			return events[a].ID < events[b].ID
		}
		return events[a].OccurredAt < events[b].OccurredAt
	})
	return events, nil
}

func ExportInventory(summary InventorySummary) (string, error) {
	if summary.TotalItems < 0 || summary.StockUnits < summary.AvailableUnits {
		return "", fmt.Errorf("invalid inventory summary")
	}
	return fmt.Sprintf("items=%d listed=%d maintenance=%d available=%d stock=%d", summary.TotalItems, summary.ListedItems, summary.MaintenanceItems, summary.AvailableUnits, summary.StockUnits), nil
}
