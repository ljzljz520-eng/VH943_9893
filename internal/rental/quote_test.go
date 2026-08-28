package rental

import (
	"testing"

	"campgear/internal/domain"
)

func TestQuote(t *testing.T) {
	cart := NewCart("q-1", "Ada", "2026-08-01", "2026-08-02")
	if err := cart.AddLine(domain.InventoryItem{ID: "l-1", DailyRate: 300, Deposit: 1000, Available: 2, Listed: true, MaintenanceStatus: domain.StatusReady}, 1, 1); err != nil {
		t.Fatal(err)
	}
	quote, err := Quote(*cart)
	if err != nil {
		t.Fatal(err)
	}
	if quote.Subtotal != 300 || quote.ServiceFee != 300 || quote.Total != 1600 {
		t.Fatalf("unexpected quote: %+v", quote)
	}
}
