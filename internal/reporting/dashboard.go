package reporting

import (
	"context"

	"campgear/internal/storage"
)

type Dashboard struct {
	Inventory       InventorySummary `json:"inventory"`
	OpenMaintenance int              `json:"open_maintenance"`
	ReservedRentals int              `json:"reserved_rentals"`
	AuditEvents     int              `json:"audit_events"`
}

func BuildDashboard(ctx context.Context, repo *storage.Repository) (Dashboard, error) {
	service := NewService(repo)
	inventory, err := service.Inventory(ctx)
	if err != nil {
		return Dashboard{}, err
	}
	orders, err := repo.ListMaintenance(ctx, "open")
	if err != nil {
		return Dashboard{}, err
	}
	rentals, err := repo.ListRentals(ctx, "reserved")
	if err != nil {
		return Dashboard{}, err
	}
	audits, err := repo.ListAudit(ctx, "", "")
	if err != nil {
		return Dashboard{}, err
	}
	return Dashboard{Inventory: inventory, OpenMaintenance: len(orders), ReservedRentals: len(rentals), AuditEvents: len(audits)}, nil
}
