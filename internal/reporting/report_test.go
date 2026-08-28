package reporting

import (
	"context"
	"testing"

	"campgear/internal/catalog"
	"campgear/internal/storage"
)

func TestInventoryReport(t *testing.T) {
	repo, err := storage.Open("file:reporting-workflow?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	service := catalog.NewService(repo)
	if _, err := service.Create(context.Background(), catalog.ItemInput{ID: "r-tent", SKU: "RT-1", Name: "Report Tent", Category: catalog.CategoryTent, DailyRate: 1200, Deposit: 4000, Stock: 2, StorageBin: "X-1"}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.Create(context.Background(), catalog.ItemInput{ID: "r-light", SKU: "RL-1", Name: "Report Light", Category: catalog.CategoryLight, DailyRate: 300, Deposit: 1000, Stock: 3, StorageBin: "X-2", MaintenanceStatus: catalog.StatusMaintenance}); err != nil {
		t.Fatal(err)
	}
	reports := NewService(repo)
	summary, err := reports.Inventory(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if summary.TotalItems != 2 || summary.ListedItems != 1 || summary.MaintenanceItems != 1 || summary.AvailableUnits != 5 {
		t.Fatalf("summary: %+v", summary)
	}
	text, err := ExportInventory(summary)
	if err != nil || text == "" {
		t.Fatalf("export: %v %q", err, text)
	}
}
