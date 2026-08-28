package maintenance

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"campgear/internal/catalog"
	"campgear/internal/domain"
	"campgear/internal/storage"
)

const (
	StatusOpen   = "open"
	StatusClosed = "closed"
)

type Service struct {
	repo    *storage.Repository
	catalog *catalog.Service
}

func NewService(repo *storage.Repository, catalogService *catalog.Service) *Service {
	return &Service{repo: repo, catalog: catalogService}
}

func (s *Service) Open(ctx context.Context, id, itemID, reason, openedDate, technician string) (domain.MaintenanceOrder, error) {
	if strings.TrimSpace(id) == "" || strings.TrimSpace(reason) == "" || openedDate == "" {
		return domain.MaintenanceOrder{}, fmt.Errorf("maintenance id, reason, and date are required")
	}
	item, err := s.catalog.Get(ctx, itemID)
	if err != nil {
		return domain.MaintenanceOrder{}, err
	}
	if item.MaintenanceStatus == catalog.StatusRetired {
		return domain.MaintenanceOrder{}, fmt.Errorf("retired item cannot enter maintenance")
	}
	order := domain.MaintenanceOrder{ID: id, ItemID: itemID, Reason: strings.TrimSpace(reason), OpenedDate: openedDate, Status: StatusOpen, Technician: technician}
	err = s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.InsertMaintenanceTx(ctx, tx, order); err != nil {
			return err
		}
		if err := s.repo.SetMaintenanceState(ctx, tx, itemID, catalog.StatusMaintenance, false); err != nil {
			return err
		}
		return s.repo.InsertAuditTx(ctx, tx, domain.AuditEvent{ID: "audit-maint-open-" + id, EntityType: "MaintenanceOrder", EntityID: id, Action: "open", Actor: technician, OccurredAt: openedDate, Details: reason})
	})
	if err != nil {
		return domain.MaintenanceOrder{}, err
	}
	return order, nil
}

func (s *Service) Close(ctx context.Context, id, closedDate, technician, inspection string) (domain.MaintenanceOrder, error) {
	if closedDate == "" || strings.TrimSpace(inspection) == "" {
		return domain.MaintenanceOrder{}, fmt.Errorf("close date and inspection are required")
	}
	order, err := s.repo.GetMaintenance(ctx, id)
	if err != nil {
		return domain.MaintenanceOrder{}, err
	}
	if order.Status != StatusOpen {
		return domain.MaintenanceOrder{}, fmt.Errorf("maintenance order already closed")
	}
	order.ClosedDate = closedDate
	order.Status = StatusClosed
	order.Technician = technician
	item, err := s.catalog.Get(ctx, order.ItemID)
	if err != nil {
		return domain.MaintenanceOrder{}, err
	}
	err = s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.UpdateMaintenance(ctx, tx, order); err != nil {
			return err
		}
		listed := item.Available > 0
		if err := s.repo.SetMaintenanceState(ctx, tx, order.ItemID, catalog.StatusReady, listed); err != nil {
			return err
		}
		return s.repo.InsertAuditTx(ctx, tx, domain.AuditEvent{ID: "audit-maint-close-" + id, EntityType: "MaintenanceOrder", EntityID: id, Action: "close", Actor: technician, OccurredAt: closedDate, Details: inspection})
	})
	if err != nil {
		return domain.MaintenanceOrder{}, err
	}
	return order, nil
}

func (s *Service) Get(ctx context.Context, id string) (domain.MaintenanceOrder, error) {
	return s.repo.GetMaintenance(ctx, id)
}

func (s *Service) List(ctx context.Context, status string) ([]domain.MaintenanceOrder, error) {
	if status != "" && status != StatusOpen && status != StatusClosed {
		return nil, fmt.Errorf("unsupported maintenance status")
	}
	return s.repo.ListMaintenance(ctx, status)
}
