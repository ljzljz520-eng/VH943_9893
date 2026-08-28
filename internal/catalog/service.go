package catalog

import (
	"context"
	"fmt"
	"strings"

	"campgear/internal/storage"
)

type Service struct {
	repo *storage.Repository
}

func NewService(repo *storage.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, input ItemInput) (InventoryItem, error) {
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		return InventoryItem{}, err
	}
	item := InventoryItem{ID: input.ID, SKU: input.SKU, Name: strings.TrimSpace(input.Name), Category: input.Category, DailyRate: input.DailyRate, Deposit: input.Deposit, Stock: input.Stock, Available: input.Stock, MaintenanceStatus: input.MaintenanceStatus, StorageBin: input.StorageBin, Listed: input.MaintenanceStatus == StatusReady, Version: 1}
	if err := s.repo.InsertItem(ctx, item); err != nil {
		if strings.Contains(err.Error(), "UNIQUE") {
			return InventoryItem{}, ErrDuplicateSKU
		}
		return InventoryItem{}, err
	}
	return item, nil
}

func (s *Service) Update(ctx context.Context, id string, input ItemInput) (InventoryItem, error) {
	input = input.Normalize()
	if input.ID == "" {
		input.ID = id
	}
	if input.ID != id {
		return InventoryItem{}, fmt.Errorf("item id cannot change")
	}
	if err := input.Validate(); err != nil {
		return InventoryItem{}, err
	}
	old, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return InventoryItem{}, err
	}
	available := old.Available
	delta := input.Stock - old.Stock
	if delta > 0 {
		available += delta
	}
	if available > input.Stock {
		available = input.Stock
	}
	listed := input.MaintenanceStatus == StatusReady && available > 0
	item := InventoryItem{ID: id, SKU: input.SKU, Name: strings.TrimSpace(input.Name), Category: input.Category, DailyRate: input.DailyRate, Deposit: input.Deposit, Stock: input.Stock, Available: available, MaintenanceStatus: input.MaintenanceStatus, StorageBin: input.StorageBin, Listed: listed, Version: old.Version + 1}
	if err := s.repo.UpdateItem(ctx, item, old.Version); err != nil {
		return InventoryItem{}, err
	}
	return item, nil
}

func (s *Service) Get(ctx context.Context, id string) (InventoryItem, error) {
	return s.repo.GetItem(ctx, id)
}

func (s *Service) List(ctx context.Context, category Category, onlyListed bool) ([]InventoryItem, error) {
	if category != "" && !ValidCategory(category) {
		return nil, ErrInvalidCategory
	}
	return s.repo.ListItems(ctx, category, onlyListed)
}

func (s *Service) SetMaintenance(ctx context.Context, id string, status MaintenanceStatus, reason string) (InventoryItem, error) {
	if !ValidMaintenanceStatus(status) {
		return InventoryItem{}, ErrInvalidState
	}
	item, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return InventoryItem{}, err
	}
	if status == StatusReady && strings.TrimSpace(reason) == "" {
		return InventoryItem{}, fmt.Errorf("release requires inspection reason")
	}
	item.MaintenanceStatus = status
	item.Listed = status == StatusReady && item.Available > 0
	item.Version++
	if err := s.repo.UpdateItem(ctx, item, item.Version-1); err != nil {
		return InventoryItem{}, err
	}
	return item, nil
}

func (s *Service) AdjustStock(ctx context.Context, id string, delta int) (InventoryItem, error) {
	if delta == 0 {
		return InventoryItem{}, fmt.Errorf("stock delta cannot be zero")
	}
	item, err := s.repo.GetItem(ctx, id)
	if err != nil {
		return InventoryItem{}, err
	}
	if item.Stock+delta < 0 || item.Available+delta < 0 {
		return InventoryItem{}, fmt.Errorf("stock would become negative")
	}
	item.Stock += delta
	item.Available += delta
	item.Listed = item.IsOperational() && item.Available > 0
	oldVersion := item.Version
	item.Version++
	if err := s.repo.UpdateItem(ctx, item, oldVersion); err != nil {
		return InventoryItem{}, err
	}
	return item, nil
}
