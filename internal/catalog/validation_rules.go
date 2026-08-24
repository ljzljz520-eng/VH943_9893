package catalog

import (
	"fmt"
	"sort"
	"strings"

	"campgear/internal/domain"
)

type ValidationIssue struct {
	ItemID string `json:"item_id"`
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func ValidateBatch(inputs []ItemInput) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	seenSKU := make(map[string]string)
	for _, original := range inputs {
		input := original.Normalize()
		if input.ID == "" {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "id", Reason: "id is required"})
		}
		if input.SKU == "" {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "sku", Reason: "sku is required"})
		} else if previous, exists := seenSKU[input.SKU]; exists {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "sku", Reason: "duplicates " + previous})
		} else {
			seenSKU[input.SKU] = input.ID
		}
		if strings.TrimSpace(input.Name) == "" {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "name", Reason: "name is required"})
		}
		if !domain.ValidCategory(input.Category) {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "category", Reason: "unsupported category"})
		}
		if input.DailyRate <= 0 {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "daily_rate", Reason: "must be positive"})
		}
		if input.Stock <= 0 {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "stock", Reason: "must be positive"})
		}
		if input.StorageBin == "" {
			issues = append(issues, ValidationIssue{ItemID: input.ID, Field: "storage_bin", Reason: "storage bin is required"})
		}
	}
	sort.SliceStable(issues, func(a, b int) bool {
		if issues[a].ItemID == issues[b].ItemID {
			return issues[a].Field < issues[b].Field
		}
		return issues[a].ItemID < issues[b].ItemID
	})
	return issues
}

type RestockRecommendation struct {
	ItemID       string `json:"item_id"`
	CurrentStock int    `json:"current_stock"`
	TargetStock  int    `json:"target_stock"`
	OrderUnits   int    `json:"order_units"`
}

func RecommendRestock(items []domain.InventoryItem, target int) ([]RestockRecommendation, error) {
	if target <= 0 {
		return nil, fmt.Errorf("target stock must be positive")
	}
	result := make([]RestockRecommendation, 0)
	for _, item := range items {
		if item.Stock >= target || item.MaintenanceStatus == domain.StatusRetired {
			continue
		}
		result = append(result, RestockRecommendation{ItemID: item.ID, CurrentStock: item.Stock, TargetStock: target, OrderUnits: target - item.Stock})
	}
	sort.Slice(result, func(a, b int) bool { return result[a].OrderUnits > result[b].OrderUnits })
	return result, nil
}

type PriceBand string

const (
	PriceBudget  PriceBand = "budget"
	PriceRegular PriceBand = "regular"
	PricePremium PriceBand = "premium"
)

func ClassifyPrice(dailyRate int64) (PriceBand, error) {
	if dailyRate <= 0 {
		return "", fmt.Errorf("daily rate must be positive")
	}
	if dailyRate < 500 {
		return PriceBudget, nil
	}
	if dailyRate < 1500 {
		return PriceRegular, nil
	}
	return PricePremium, nil
}
