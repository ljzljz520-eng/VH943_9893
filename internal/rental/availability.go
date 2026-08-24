package rental

import (
	"context"
	"fmt"

	"campgear/internal/domain"
)

type AvailabilityRequest struct {
	ItemID   string
	Quantity int
	Days     int
}

type AvailabilityResult struct {
	ItemID    string `json:"item_id"`
	Requested int    `json:"requested"`
	Available int    `json:"available"`
	Rentable  bool   `json:"rentable"`
	Reason    string `json:"reason"`
}

func (s *Service) CheckAvailability(ctx context.Context, request AvailabilityRequest) (AvailabilityResult, error) {
	if request.Quantity <= 0 || request.Days <= 0 {
		return AvailabilityResult{}, fmt.Errorf("quantity and days must be positive")
	}
	item, err := s.catalog.Get(ctx, request.ItemID)
	if err != nil {
		return AvailabilityResult{}, err
	}
	result := AvailabilityResult{ItemID: request.ItemID, Requested: request.Quantity, Available: item.Available, Rentable: item.CanRent(request.Quantity)}
	if item.MaintenanceStatus != domain.StatusReady {
		result.Reason = "equipment is in maintenance"
	} else if !item.Listed {
		result.Reason = "equipment is not listed"
	} else if item.Available < request.Quantity {
		result.Reason = "insufficient stock"
	}
	return result, nil
}

func (s *Service) CheckBatch(ctx context.Context, requests []AvailabilityRequest) ([]AvailabilityResult, error) {
	result := make([]AvailabilityResult, 0, len(requests))
	for _, request := range requests {
		availability, err := s.CheckAvailability(ctx, request)
		if err != nil {
			return nil, err
		}
		result = append(result, availability)
	}
	return result, nil
}
