package maintenance

import (
	"context"
	"testing"

	"campgear/internal/catalog"
	"campgear/internal/storage"
)

func TestMaintenanceWorkflow(t *testing.T) {
	repo, err := storage.Open("file:maintenance-workflow?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	catalogService := catalog.NewService(repo)
	item, err := catalogService.Create(context.Background(), catalog.ItemInput{ID: "tent-m1", SKU: "T-M1", Name: "Family Tent", Category: catalog.CategoryTent, DailyRate: 1500, Deposit: 5000, Stock: 1, StorageBin: "M-01"})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(repo, catalogService)
	order, err := service.Open(context.Background(), "wo-1", item.ID, "zipper repair", "2026-08-01", "tech-1")
	if err != nil {
		t.Fatal(err)
	}
	if order.Status != StatusOpen {
		t.Fatalf("order: %+v", order)
	}
	underRepair, err := catalogService.Get(context.Background(), item.ID)
	if err != nil || underRepair.MaintenanceStatus != catalog.StatusMaintenance || underRepair.Listed {
		t.Fatalf("maintenance state: %v %+v", err, underRepair)
	}
	closed, err := service.Close(context.Background(), order.ID, "2026-08-02", "tech-1", "zipper replaced")
	if err != nil || closed.Status != StatusClosed {
		t.Fatalf("close: %v %+v", err, closed)
	}
	ready, err := catalogService.Get(context.Background(), item.ID)
	if err != nil || ready.MaintenanceStatus != catalog.StatusReady || !ready.Listed {
		t.Fatalf("release state: %v %+v", err, ready)
	}
}
