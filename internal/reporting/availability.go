package reporting

import (
	"context"
	"sort"

	"campgear/internal/domain"
	"campgear/internal/storage"
)

type Availability struct {
	Category domain.Category `json:"category"`
	Units    int             `json:"units"`
	Listed   int             `json:"listed"`
}

func GroupAvailability(ctx context.Context, repo *storage.Repository) ([]Availability, error) {
	items, err := repo.ListItems(ctx, "", false)
	if err != nil {
		return nil, err
	}
	groups := make(map[domain.Category]*Availability)
	for _, item := range items {
		group := groups[item.Category]
		if group == nil {
			group = &Availability{Category: item.Category}
			groups[item.Category] = group
		}
		group.Units += item.Available
		if item.Listed {
			group.Listed++
		}
	}
	result := make([]Availability, 0, len(groups))
	for _, group := range groups {
		result = append(result, *group)
	}
	sort.Slice(result, func(a, b int) bool { return result[a].Category < result[b].Category })
	return result, nil
}
