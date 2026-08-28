package catalog

import "sort"

type Filter struct {
	Query        string
	Category     Category
	OnlyListed   bool
	MinDailyRate int64
	MaxDailyRate int64
}

func FilterItems(items []InventoryItem, filter Filter) []InventoryItem {
	result := make([]InventoryItem, 0, len(items))
	for _, item := range items {
		if filter.Category != "" && item.Category != filter.Category {
			continue
		}
		if filter.OnlyListed && !item.Listed {
			continue
		}
		if filter.MinDailyRate > 0 && item.DailyRate < filter.MinDailyRate {
			continue
		}
		if filter.MaxDailyRate > 0 && item.DailyRate > filter.MaxDailyRate {
			continue
		}
		if filter.Query != "" && !containsFold(item.Name, filter.Query) && !containsFold(item.SKU, filter.Query) {
			continue
		}
		result = append(result, item)
	}
	sort.Slice(result, func(a, b int) bool {
		if result[a].Category == result[b].Category {
			return result[a].Name < result[b].Name
		}
		return result[a].Category < result[b].Category
	})
	return result
}

func containsFold(value, query string) bool {
	return len(query) == 0 || len(value) >= len(query) && lower(value, query)
}

func lower(value, query string) bool {
	for i := 0; i+len(query) <= len(value); i++ {
		match := true
		for j := range query {
			v := value[i+j]
			q := query[j]
			if v >= 'A' && v <= 'Z' {
				v += 'a' - 'A'
			}
			if q >= 'A' && q <= 'Z' {
				q += 'a' - 'A'
			}
			if v != q {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
