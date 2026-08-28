package storage

import (
	"context"
	"path/filepath"
	"testing"

	"campgear/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "campgear.db")
	repo, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	item := domain.InventoryItem{ID: "light-1", SKU: "L-10", Name: "Beacon Light", Category: domain.CategoryLight, DailyRate: 300, Deposit: 1000, Stock: 4, Available: 4, MaintenanceStatus: domain.StatusReady, StorageBin: "C-02", Listed: true, Version: 1}
	if err := repo.InsertItem(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	if err := repo.Close(); err != nil {
		t.Fatal(err)
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	got, err := reopened.GetItem(context.Background(), item.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.SKU != item.SKU || got.Available != item.Available || !got.Listed {
		t.Fatalf("persistence mismatch: %+v", got)
	}
}
