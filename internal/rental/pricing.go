package rental

import (
	"fmt"

	"campgear/internal/domain"
)

type PriceSummary struct {
	Subtotal   int64 `json:"subtotal"`
	Deposit    int64 `json:"deposit"`
	ServiceFee int64 `json:"service_fee"`
	Total      int64 `json:"total"`
}

func Quote(cart Cart) (PriceSummary, error) {
	if cart.Empty() {
		return PriceSummary{}, fmt.Errorf("cannot quote empty cart")
	}
	summary := PriceSummary{Subtotal: cart.Total(), Deposit: cart.DepositHeld()}
	if summary.Subtotal >= 5000 {
		summary.ServiceFee = summary.Subtotal / 20
	} else {
		summary.ServiceFee = 300
	}
	summary.Total = summary.Subtotal + summary.ServiceFee + summary.Deposit
	return summary, nil
}

func LineFor(item domain.InventoryItem, quantity, days int) (domain.RentalLine, error) {
	if !item.CanRent(quantity) {
		return domain.RentalLine{}, fmt.Errorf("item is not rentable")
	}
	if days <= 0 {
		return domain.RentalLine{}, fmt.Errorf("rental days must be positive")
	}
	return domain.RentalLine{ItemID: item.ID, Quantity: quantity, Days: days, Rate: item.DailyRate, Deposit: item.Deposit, Subtotal: int64(quantity) * int64(days) * item.DailyRate}, nil
}
