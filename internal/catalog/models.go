package catalog

import "campgear/internal/domain"

type Category = domain.Category
type MaintenanceStatus = domain.MaintenanceStatus
type InventoryItem = domain.InventoryItem
type ItemInput = domain.ItemInput

const (
	CategoryTent        = domain.CategoryTent
	CategorySleepingBag = domain.CategorySleepingBag
	CategoryStove       = domain.CategoryStove
	CategoryLight       = domain.CategoryLight
	StatusReady         = domain.StatusReady
	StatusMaintenance   = domain.StatusMaintenance
	StatusRetired       = domain.StatusRetired
)

var (
	ErrInvalidItem     = domain.ErrInvalidItem
	ErrItemNotFound    = domain.ErrItemNotFound
	ErrDuplicateSKU    = domain.ErrDuplicateSKU
	ErrInvalidCategory = domain.ErrInvalidCategory
	ErrInvalidState    = domain.ErrInvalidState
)

func ValidCategory(category Category) bool { return domain.ValidCategory(category) }
func ValidMaintenanceStatus(status MaintenanceStatus) bool {
	return domain.ValidMaintenanceStatus(status)
}
