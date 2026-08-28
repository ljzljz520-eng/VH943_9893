package catalog

import (
	"context"
	"testing"

	"campgear/internal/storage"
)

func TestCatalogWorkflow(t *testing.T) {
	repo, err := storage.Open("file:catalog-workflow?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	service := NewService(repo)
	item, err := service.Create(context.Background(), ItemInput{ID: "tent-1", SKU: "T-100", Name: "Ridge Tent", Category: CategoryTent, DailyRate: 1200, Deposit: 5000, Stock: 3, StorageBin: "A-01"})
	if err != nil {
		t.Fatal(err)
	}
	if !item.Listed || item.Available != 3 {
		t.Fatalf("unexpected item state: %+v", item)
	}
	items, err := service.List(context.Background(), CategoryTent, true)
	if err != nil || len(items) != 1 {
		t.Fatalf("listed items: %v %#v", err, items)
	}
	updated, err := service.Update(context.Background(), item.ID, ItemInput{ID: item.ID, SKU: item.SKU, Name: item.Name, Category: item.Category, DailyRate: item.DailyRate, Deposit: item.Deposit, Stock: 4, StorageBin: item.StorageBin})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Stock != 4 || updated.Available != 4 {
		t.Fatalf("unexpected update: %+v", updated)
	}
}

func TestFilterItems(t *testing.T) {
	items := []InventoryItem{{ID: "1", Name: "Alpine Stove", SKU: "ST-1", Category: CategoryStove, DailyRate: 900, Listed: true}, {ID: "2", Name: "Ridge Tent", SKU: "TN-1", Category: CategoryTent, DailyRate: 1200, Listed: false}}
	filtered := FilterItems(items, Filter{Query: "stove", OnlyListed: true})
	if len(filtered) != 1 || filtered[0].ID != "1" {
		t.Fatalf("unexpected filter: %#v", filtered)
	}
}
