package rental

import (
	"fmt"
	"sort"

	"campgear/internal/domain"
)

type LedgerEntry struct {
	RentalID string `json:"rental_id"`
	ItemID   string `json:"item_id"`
	Date     string `json:"date"`
	Kind     string `json:"kind"`
	Units    int    `json:"units"`
	Amount   int64  `json:"amount"`
}

func BuildLedger(records []domain.RentalRecord) []LedgerEntry {
	entries := make([]LedgerEntry, 0)
	for _, record := range records {
		for _, line := range record.Lines {
			entries = append(entries, LedgerEntry{RentalID: record.ID, ItemID: line.ItemID, Date: record.StartDate, Kind: "reserve", Units: -line.Quantity, Amount: line.Subtotal})
			if record.Status == StatusReturned {
				entries = append(entries, LedgerEntry{RentalID: record.ID, ItemID: line.ItemID, Date: record.EndDate, Kind: "return", Units: line.Quantity, Amount: 0})
			}
		}
	}
	sort.SliceStable(entries, func(a, b int) bool {
		if entries[a].Date == entries[b].Date {
			if entries[a].RentalID == entries[b].RentalID {
				return entries[a].ItemID < entries[b].ItemID
			}
			return entries[a].RentalID < entries[b].RentalID
		}
		return entries[a].Date < entries[b].Date
	})
	return entries
}

func ValidateLedger(entries []LedgerEntry) error {
	balances := make(map[string]int)
	for _, entry := range entries {
		if entry.RentalID == "" || entry.ItemID == "" || entry.Date == "" {
			return fmt.Errorf("ledger entry is incomplete")
		}
		if entry.Kind != "reserve" && entry.Kind != "return" {
			return fmt.Errorf("unsupported ledger kind")
		}
		balances[entry.ItemID] += entry.Units
	}
	for itemID, balance := range balances {
		if balance > 0 {
			return fmt.Errorf("item %s has excess returns", itemID)
		}
	}
	return nil
}

func LedgerAmount(entries []LedgerEntry) int64 {
	var amount int64
	for _, entry := range entries {
		amount += entry.Amount
	}
	return amount
}
