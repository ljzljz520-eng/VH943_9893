package catalog

import (
	"fmt"
	"sort"

	"campgear/internal/domain"
)

type ListingDecision struct {
	Listed bool
	Reason string
}

func DecideListing(item domain.InventoryItem) ListingDecision {
	if item.MaintenanceStatus == domain.StatusRetired {
		return ListingDecision{Reason: "retired equipment"}
	}
	if item.MaintenanceStatus == domain.StatusMaintenance {
		return ListingDecision{Reason: "maintenance in progress"}
	}
	if item.Stock <= 0 || item.Available <= 0 {
		return ListingDecision{Reason: "no available units"}
	}
	return ListingDecision{Listed: true, Reason: "ready and available"}
}

func ValidateListing(item domain.InventoryItem) error {
	decision := DecideListing(item)
	if item.Listed != decision.Listed {
		return fmt.Errorf("listing mismatch: expected %t because %s", decision.Listed, decision.Reason)
	}
	if item.Available > item.Stock {
		return fmt.Errorf("available units exceed stock")
	}
	return nil
}

func Rebalance(items []domain.InventoryItem) []domain.InventoryItem {
	result := append([]domain.InventoryItem(nil), items...)
	for index := range result {
		if result[index].Available < 0 {
			result[index].Available = 0
		}
		if result[index].Available > result[index].Stock {
			result[index].Available = result[index].Stock
		}
		result[index].Listed = DecideListing(result[index]).Listed
	}
	sort.Slice(result, func(a, b int) bool { return result[a].StorageBin < result[b].StorageBin })
	return result
}

type CategoryTotals struct {
	Category  domain.Category
	Stock     int
	Available int
}

func TotalsByCategory(items []domain.InventoryItem) []CategoryTotals {
	groups := make(map[domain.Category]*CategoryTotals)
	for _, item := range items {
		group := groups[item.Category]
		if group == nil {
			group = &CategoryTotals{Category: item.Category}
			groups[item.Category] = group
		}
		group.Stock += item.Stock
		group.Available += item.Available
	}
	result := make([]CategoryTotals, 0, len(groups))
	for _, group := range groups {
		result = append(result, *group)
	}
	sort.Slice(result, func(a, b int) bool { return result[a].Category < result[b].Category })
	return result
}
