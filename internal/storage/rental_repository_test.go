package storage

import (
	"context"
	"database/sql"
	"testing"

	"campgear/internal/domain"
)

func TestRentalRepository(t *testing.T) {
	repo, err := Open("file:rental-repository?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer repo.Close()
	item := domain.InventoryItem{ID: "bag-1", SKU: "B-1", Name: "Down Bag", Category: domain.CategorySleepingBag, DailyRate: 700, Deposit: 2000, Stock: 2, Available: 2, MaintenanceStatus: domain.StatusReady, StorageBin: "B-01", Listed: true, Version: 1}
	if err := repo.InsertItem(context.Background(), item); err != nil {
		t.Fatal(err)
	}
	record := domain.RentalRecord{ID: "rent-1", Customer: "Ada", StartDate: "2026-08-01", EndDate: "2026-08-03", Status: "reserved", Total: 1400, DepositHeld: 2000}
	err = repo.WithTx(context.Background(), func(tx *sql.Tx) error {
		if err := repo.InsertRental(context.Background(), tx, record); err != nil {
			return err
		}
		return repo.InsertRentalLine(context.Background(), tx, record.ID, domain.RentalLine{ItemID: item.ID, Quantity: 1, Days: 2, Rate: 700, Deposit: 2000, Subtotal: 1400})
	})
	if err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetRental(context.Background(), record.ID)
	if err != nil || len(got.Lines) != 1 {
		t.Fatalf("rental read: %v %#v", err, got)
	}
}
