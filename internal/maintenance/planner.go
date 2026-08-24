package maintenance

import (
	"fmt"
	"sort"

	"campgear/internal/domain"
)

type Priority string

const (
	PriorityRoutine Priority = "routine"
	PriorityUrgent  Priority = "urgent"
)

type PlanEntry struct {
	OrderID  string   `json:"order_id"`
	ItemID   string   `json:"item_id"`
	Priority Priority `json:"priority"`
	Reason   string   `json:"reason"`
}

func BuildPlan(orders []domain.MaintenanceOrder, items []domain.InventoryItem) ([]PlanEntry, error) {
	lookup := make(map[string]domain.InventoryItem, len(items))
	for _, item := range items {
		lookup[item.ID] = item
	}
	plan := make([]PlanEntry, 0, len(orders))
	for _, order := range orders {
		if order.Status != StatusOpen {
			continue
		}
		item, ok := lookup[order.ItemID]
		if !ok {
			return nil, fmt.Errorf("item %s missing", order.ItemID)
		}
		priority := PriorityRoutine
		if item.Category == domain.CategoryTent || item.Stock > 4 {
			priority = PriorityUrgent
		}
		plan = append(plan, PlanEntry{OrderID: order.ID, ItemID: order.ItemID, Priority: priority, Reason: order.Reason})
	}
	sort.SliceStable(plan, func(a, b int) bool {
		if plan[a].Priority == plan[b].Priority {
			return plan[a].OrderID < plan[b].OrderID
		}
		return plan[a].Priority == PriorityUrgent
	})
	return plan, nil
}
