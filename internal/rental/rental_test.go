package rental

import (
	"context"
	"testing"

	"campgear/internal/catalog"
	"campgear/internal/storage"
)

func TestRentalWorkflow(t *testing.T) {
	repo, err := storage.Open("file:rental-workflow?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	catalogService := catalog.NewService(repo)
	item, err := catalogService.Create(context.Background(), catalog.ItemInput{ID: "stove-1", SKU: "S-1", Name: "Trail Stove", Category: catalog.CategoryStove, DailyRate: 800, Deposit: 2500, Stock: 2, StorageBin: "D-04"})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(repo, catalogService)
	cart := NewCart("rent-1", "Mina", "2026-08-10", "2026-08-12")
	if err := service.AddToCart(context.Background(), cart, item.ID, 1, 2); err != nil {
		t.Fatal(err)
	}
	quote, err := Quote(*cart)
	if err != nil || quote.Total != 4400 {
		t.Fatalf("quote: %v %#v", err, quote)
	}
	record, err := service.CreateRental(context.Background(), *cart, "staff-1")
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != StatusReserved || record.Total != 1600 {
		t.Fatalf("record: %+v", record)
	}
	got, err := service.Activate(context.Background(), record.ID)
	if err != nil || got.Status != StatusActive {
		t.Fatalf("activate: %v %+v", err, got)
	}
	got, err = service.Return(context.Background(), record.ID, "2026-08-12")
	if err != nil || got.Status != StatusReturned {
		t.Fatalf("return: %v %+v", err, got)
	}
	available, err := catalogService.Get(context.Background(), item.ID)
	if err != nil || available.Available != 2 {
		t.Fatalf("availability: %v %+v", err, available)
	}
}

func TestRentalRejectsMaintenance(t *testing.T) {
	repo, err := storage.Open("file:rental-maintenance-regression?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	catalogService := catalog.NewService(repo)
	item, err := catalogService.Create(context.Background(), catalog.ItemInput{ID: "tent-maint", SKU: "T-MAINT", Name: "Repair Tent", Category: catalog.CategoryTent, DailyRate: 1100, Deposit: 4000, Stock: 1, StorageBin: "R-01", MaintenanceStatus: catalog.StatusMaintenance})
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(repo, catalogService)
	cart := NewCart("rent-maint", "Mina", "2026-08-10", "2026-08-11")
	if err := service.AddToCart(context.Background(), cart, item.ID, 1, 1); err == nil {
		t.Fatal("maintenance equipment must be rejected")
	}
}
