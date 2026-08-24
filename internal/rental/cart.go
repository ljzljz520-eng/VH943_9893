package rental

import (
	"fmt"

	"campgear/internal/domain"
)

type Cart struct {
	ID        string
	Customer  string
	StartDate string
	EndDate   string
	Lines     []domain.RentalLine
}

func NewCart(id, customer, startDate, endDate string) *Cart {
	return &Cart{ID: id, Customer: customer, StartDate: startDate, EndDate: endDate, Lines: make([]domain.RentalLine, 0)}
}

func (c *Cart) AddLine(item domain.InventoryItem, quantity, days int) error {
	if err := rentable(item, quantity); err != nil {
		return err
	}
	if days <= 0 {
		return fmt.Errorf("rental days must be positive")
	}
	line := domain.RentalLine{ItemID: item.ID, Quantity: quantity, Days: days, Rate: item.DailyRate, Deposit: item.Deposit, Subtotal: int64(quantity) * int64(days) * item.DailyRate}
	for index := range c.Lines {
		if c.Lines[index].ItemID == item.ID {
			c.Lines[index].Quantity += quantity
			c.Lines[index].Subtotal += line.Subtotal
			return nil
		}
	}
	c.Lines = append(c.Lines, line)
	return nil
}

func (c *Cart) RemoveLine(itemID string) bool {
	for index, line := range c.Lines {
		if line.ItemID == itemID {
			c.Lines = append(c.Lines[:index], c.Lines[index+1:]...)
			return true
		}
	}
	return false
}

func (c Cart) Empty() bool { return len(c.Lines) == 0 }

func (c Cart) Total() int64 {
	var total int64
	for _, line := range c.Lines {
		total += line.Subtotal
	}
	return total
}

func (c Cart) DepositHeld() int64 {
	var total int64
	for _, line := range c.Lines {
		total += int64(line.Quantity) * line.Deposit
	}
	return total
}

func (c Cart) ToRecord(status string) domain.RentalRecord {
	return domain.RentalRecord{ID: c.ID, Customer: c.Customer, StartDate: c.StartDate, EndDate: c.EndDate, Status: status, Lines: append([]domain.RentalLine(nil), c.Lines...), Total: c.Total(), DepositHeld: c.DepositHeld()}
}

// rentable validates that an item may enter a rental. It mirrors
// InventoryItem.CanRent but produces a specific reason so callers can report
// why the rental was refused. Equipment under maintenance is rejected outright
// even when units are technically available, so a repair holds while the item
// stays off the rental desk.
func rentable(item domain.InventoryItem, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}
	if item.MaintenanceStatus == domain.StatusMaintenance {
		return fmt.Errorf("item %s is in maintenance", item.ID)
	}
	if item.MaintenanceStatus == domain.StatusRetired {
		return fmt.Errorf("item %s is retired", item.ID)
	}
	if !item.Listed {
		return fmt.Errorf("item %s is not listed for rental", item.ID)
	}
	if item.Available < quantity {
		return fmt.Errorf("insufficient availability for %s", item.ID)
	}
	return nil
}
