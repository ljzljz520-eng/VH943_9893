package domain

import "errors"

type Category string

const (
	CategoryTent        Category = "tent"
	CategorySleepingBag Category = "sleeping_bag"
	CategoryStove       Category = "stove"
	CategoryLight       Category = "light"
)

type MaintenanceStatus string

const (
	StatusReady       MaintenanceStatus = "ready"
	StatusMaintenance MaintenanceStatus = "maintenance"
	StatusRetired     MaintenanceStatus = "retired"
)

type InventoryItem struct {
	ID                string            `json:"id"`
	SKU               string            `json:"sku"`
	Name              string            `json:"name"`
	Category          Category          `json:"category"`
	DailyRate         int64             `json:"daily_rate"`
	Deposit           int64             `json:"deposit"`
	Stock             int               `json:"stock"`
	Available         int               `json:"available"`
	MaintenanceStatus MaintenanceStatus `json:"maintenance_status"`
	StorageBin        string            `json:"storage_bin"`
	Listed            bool              `json:"listed"`
	Version           int               `json:"version"`
}

type ItemInput struct {
	ID                string            `json:"id"`
	SKU               string            `json:"sku"`
	Name              string            `json:"name"`
	Category          Category          `json:"category"`
	DailyRate         int64             `json:"daily_rate"`
	Deposit           int64             `json:"deposit"`
	Stock             int               `json:"stock"`
	StorageBin        string            `json:"storage_bin"`
	MaintenanceStatus MaintenanceStatus `json:"maintenance_status"`
}

type RentalLine struct {
	ItemID   string `json:"item_id"`
	Quantity int    `json:"quantity"`
	Days     int    `json:"days"`
	Rate     int64  `json:"rate"`
	Deposit  int64  `json:"deposit"`
	Subtotal int64  `json:"subtotal"`
}

type RentalRecord struct {
	ID          string       `json:"id"`
	Customer    string       `json:"customer"`
	StartDate   string       `json:"start_date"`
	EndDate     string       `json:"end_date"`
	Status      string       `json:"status"`
	Lines       []RentalLine `json:"lines"`
	Total       int64        `json:"total"`
	DepositHeld int64        `json:"deposit_held"`
}

type MaintenanceOrder struct {
	ID         string `json:"id"`
	ItemID     string `json:"item_id"`
	Reason     string `json:"reason"`
	OpenedDate string `json:"opened_date"`
	ClosedDate string `json:"closed_date"`
	Status     string `json:"status"`
	Technician string `json:"technician"`
}

type AuditEvent struct {
	ID         string `json:"id"`
	EntityType string `json:"entity_type"`
	EntityID   string `json:"entity_id"`
	Action     string `json:"action"`
	Actor      string `json:"actor"`
	OccurredAt string `json:"occurred_at"`
	Details    string `json:"details"`
}

type StaffMember struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

var (
	ErrInvalidItem     = errors.New("invalid inventory item")
	ErrItemNotFound    = errors.New("inventory item not found")
	ErrDuplicateSKU    = errors.New("sku already exists")
	ErrInvalidCategory = errors.New("unsupported category")
	ErrInvalidState    = errors.New("unsupported maintenance state")
)

func ValidCategory(category Category) bool {
	switch category {
	case CategoryTent, CategorySleepingBag, CategoryStove, CategoryLight:
		return true
	default:
		return false
	}
}

func ValidMaintenanceStatus(status MaintenanceStatus) bool {
	switch status {
	case StatusReady, StatusMaintenance, StatusRetired:
		return true
	default:
		return false
	}
}

func (i ItemInput) Normalize() ItemInput {
	if i.MaintenanceStatus == "" {
		i.MaintenanceStatus = StatusReady
	}
	if i.Stock < 0 {
		i.Stock = 0
	}
	if i.DailyRate < 0 {
		i.DailyRate = 0
	}
	if i.Deposit < 0 {
		i.Deposit = 0
	}
	return i
}

func (i ItemInput) Validate() error {
	if i.ID == "" || i.SKU == "" || i.Name == "" || i.StorageBin == "" {
		return ErrInvalidItem
	}
	if !ValidCategory(i.Category) {
		return ErrInvalidCategory
	}
	if !ValidMaintenanceStatus(i.MaintenanceStatus) {
		return ErrInvalidState
	}
	if i.DailyRate <= 0 || i.Deposit < 0 || i.Stock <= 0 {
		return ErrInvalidItem
	}
	return nil
}

func (i InventoryItem) CanRent(quantity int) bool {
	if quantity <= 0 {
		return false
	}
	if !i.Listed || i.MaintenanceStatus != StatusReady {
		return false
	}
	return i.Available >= quantity
}

func (i InventoryItem) IsOperational() bool {
	return i.MaintenanceStatus == StatusReady && i.Listed
}
