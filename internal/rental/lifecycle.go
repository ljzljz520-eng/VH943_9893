package rental

import (
	"context"
	"database/sql"
	"fmt"

	"campgear/internal/domain"
)

func (s *Service) Activate(ctx context.Context, id string) (domain.RentalRecord, error) {
	return s.transition(ctx, id, StatusActive)
}

func (s *Service) Return(ctx context.Context, id, returnDate string) (domain.RentalRecord, error) {
	if returnDate == "" {
		return domain.RentalRecord{}, fmt.Errorf("return date is required")
	}
	record, err := s.Get(ctx, id)
	if err != nil {
		return domain.RentalRecord{}, err
	}
	if record.Status != StatusActive && record.Status != StatusReserved {
		return domain.RentalRecord{}, fmt.Errorf("rental cannot be returned from %s", record.Status)
	}
	err = s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		for _, line := range record.Lines {
			if err := s.repo.ChangeAvailability(ctx, tx, line.ItemID, line.Quantity); err != nil {
				return err
			}
		}
		if err := s.repo.UpdateRentalStatusTx(ctx, tx, id, StatusReturned); err != nil {
			return err
		}
		return s.repo.InsertAuditTx(ctx, tx, domain.AuditEvent{ID: "audit-return-" + id, EntityType: "RentalRecord", EntityID: id, Action: "return", Actor: "staff", OccurredAt: returnDate, Details: "gear returned"})
	})
	if err != nil {
		return domain.RentalRecord{}, err
	}
	record.Status = StatusReturned
	return record, nil
}

func (s *Service) Cancel(ctx context.Context, id, reason string) (domain.RentalRecord, error) {
	if reason == "" {
		return domain.RentalRecord{}, fmt.Errorf("cancel reason is required")
	}
	record, err := s.Get(ctx, id)
	if err != nil {
		return domain.RentalRecord{}, err
	}
	if record.Status == StatusReturned || record.Status == StatusCanceled {
		return domain.RentalRecord{}, fmt.Errorf("rental already closed")
	}
	err = s.repo.WithTx(ctx, func(tx *sql.Tx) error {
		for _, line := range record.Lines {
			if err := s.repo.ChangeAvailability(ctx, tx, line.ItemID, line.Quantity); err != nil {
				return err
			}
		}
		if err := s.repo.UpdateRentalStatusTx(ctx, tx, id, StatusCanceled); err != nil {
			return err
		}
		return s.repo.InsertAuditTx(ctx, tx, domain.AuditEvent{ID: "audit-cancel-" + id, EntityType: "RentalRecord", EntityID: id, Action: "cancel", Actor: "staff", OccurredAt: record.StartDate, Details: reason})
	})
	if err != nil {
		return domain.RentalRecord{}, err
	}
	record.Status = StatusCanceled
	return record, nil
}

func (s *Service) transition(ctx context.Context, id, status string) (domain.RentalRecord, error) {
	record, err := s.Get(ctx, id)
	if err != nil {
		return domain.RentalRecord{}, err
	}
	if status == StatusActive && record.Status != StatusReserved {
		return domain.RentalRecord{}, fmt.Errorf("only reserved rentals can activate")
	}
	if err := s.repo.UpdateRentalStatus(ctx, id, status); err != nil {
		return domain.RentalRecord{}, err
	}
	record.Status = status
	return record, nil
}
