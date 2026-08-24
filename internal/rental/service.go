package rental

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
	StatusReserved = "reserved"
	StatusActive   = "active"
	StatusReturned = "returned"
	StatusCanceled = "canceled"
)

type Services struct {
	Repo        *storage.Repository
	Catalog     *catalog.Service
	Maintenance *MaintenanceGateway
}

type MaintenanceGateway struct{}

func NewServices(repo *storage.Repository) *Services {
	return &Services{Repo: repo, Catalog: catalog.NewService(repo), Maintenance: &MaintenanceGateway{}}
}

type Service struct {
	repo    *storage.Repository
	catalog *catalog.Service
}

func NewService(repo *storage.Repository, catalogService *catalog.Service) *Service {
	return &Service{repo: repo, catalog: catalogService}
}

func (s *Service) AddToCart(ctx context.Context, cart *Cart, itemID string, quantity, days int) error {
	if cart == nil {
		return fmt.Errorf("cart is required")
	}
	if err := ValidateDates(cart.StartDate, cart.EndDate); err != nil {
		return err
	}
	item, err := s.catalog.Get(ctx, itemID)
	if err != nil {
		return err
	}
	if item.Available < quantity {
		return fmt.Errorf("insufficient availability")
	}
	if err := cart.AddLine(item, quantity, days); err != nil {
		return err
	}
	return nil
}

func (s *Service) CreateRental(ctx context.Context, cart Cart, actor string) (domain.RentalRecord, error) {
	if cart.Empty() {
		return domain.RentalRecord{}, fmt.Errorf("rental cart is empty")
	}
	if strings.TrimSpace(cart.Customer) == "" {
		return domain.RentalRecord{}, fmt.Errorf("customer is required")
	}
	if err := ValidateDates(cart.StartDate, cart.EndDate); err != nil {
		return domain.RentalRecord{}, err
	}
	record := cart.ToRecord(StatusReserved)
	err := s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		if err := s.repo.InsertRental(ctx, tx, record); err != nil {
			return err
		}
		for _, line := range record.Lines {
			if err := s.repo.InsertRentalLine(ctx, tx, record.ID, line); err != nil {
				return err
			}
			if err := s.repo.ChangeAvailability(ctx, tx, line.ItemID, -line.Quantity); err != nil {
				return err
			}
		}
		return s.repo.InsertAuditTx(ctx, tx, domain.AuditEvent{ID: "audit-" + record.ID, EntityType: "RentalRecord", EntityID: record.ID, Action: "create", Actor: actor, OccurredAt: record.StartDate, Details: "rental reserved"})
	})
	if err != nil {
		return domain.RentalRecord{}, err
	}
	return record, nil
}

func ValidateDates(startDate, endDate string) error {
	if len(startDate) != 10 || len(endDate) != 10 {
		return fmt.Errorf("dates must use YYYY-MM-DD")
	}
	if startDate > endDate {
		return fmt.Errorf("end date must not precede start date")
	}
	return nil
}

func (s *Service) Get(ctx context.Context, id string) (domain.RentalRecord, error) {
	return s.repo.GetRental(ctx, id)
}

func (s *Service) List(ctx context.Context, status string) ([]domain.RentalRecord, error) {
	return s.repo.ListRentals(ctx, status)
}

func (s *Service) UpdateStatus(ctx context.Context, id, status string) error {
	if !validTransition(status) {
		return fmt.Errorf("unsupported rental status")
	}
	return s.repo.UpdateRentalStatus(ctx, id, status)
}

func validTransition(status string) bool {
	switch status {
	case StatusReserved, StatusActive, StatusReturned, StatusCanceled:
		return true
	default:
		return false
	}
}
