package catalog

import (
	"context"
	"fmt"
	"sort"

	"campgear/internal/domain"
)

type ManifestResult struct {
	Created int
	Skipped int
	IDs     []string
}

func (s *Service) ImportManifest(ctx context.Context, inputs []ItemInput) (ManifestResult, error) {
	if len(inputs) == 0 {
		return ManifestResult{}, fmt.Errorf("manifest is empty")
	}
	ordered := append([]ItemInput(nil), inputs...)
	sort.SliceStable(ordered, func(a, b int) bool { return ordered[a].ID < ordered[b].ID })
	result := ManifestResult{IDs: make([]string, 0, len(ordered))}
	for _, input := range ordered {
		if _, err := s.Get(ctx, input.ID); err == nil {
			result.Skipped++
			continue
		}
		if _, err := s.Create(ctx, input); err != nil {
			return result, err
		}
		result.Created++
		result.IDs = append(result.IDs, input.ID)
	}
	return result, nil
}

func ManifestByCategory(inputs []ItemInput) map[domain.Category]int {
	result := make(map[domain.Category]int)
	for _, input := range inputs {
		result[input.Category]++
	}
	return result
}
