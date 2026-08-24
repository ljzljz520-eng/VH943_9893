package domain

import "testing"

func TestItemInputValidation(t *testing.T) {
	input := ItemInput{ID: "tent-1", SKU: "T-1", Name: "Ridge tent", Category: CategoryTent, DailyRate: 1200, Deposit: 5000, Stock: 2, StorageBin: "A-01"}
	input = input.Normalize()
	if err := input.Validate(); err != nil {
		t.Fatal(err)
	}
	if input.MaintenanceStatus != StatusReady {
		t.Fatalf("expected ready default, got %s", input.MaintenanceStatus)
	}
	if (InventoryItem{MaintenanceStatus: StatusMaintenance, Listed: false}).CanRent(1) {
		t.Fatal("maintenance item must not be rentable")
	}
}
