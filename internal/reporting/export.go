package reporting

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"

	"campgear/internal/domain"
	"campgear/internal/storage"
)

func WriteInventoryCSV(ctx context.Context, repo *storage.Repository, writer io.Writer) error {
	items, err := repo.ListItems(ctx, "", false)
	if err != nil {
		return err
	}
	csvWriter := csv.NewWriter(writer)
	if err := csvWriter.Write([]string{"id", "sku", "category", "daily_rate", "stock", "available", "status", "listed"}); err != nil {
		return err
	}
	for _, item := range items {
		row := []string{item.ID, item.SKU, string(item.Category), strconv.FormatInt(item.DailyRate, 10), strconv.Itoa(item.Stock), strconv.Itoa(item.Available), string(item.MaintenanceStatus), strconv.FormatBool(item.Listed)}
		if err := csvWriter.Write(row); err != nil {
			return err
		}
	}
	csvWriter.Flush()
	return csvWriter.Error()
}

func FormatRevenue(report RevenueReport) string {
	return fmt.Sprintf("rentals=%d gross=%d deposits=%d", report.RentalCount, report.Gross, report.Deposits)
}

func CategoryBreakdown(report RevenueReport) []string {
	result := make([]string, 0, len(report.ByCategory))
	for category, amount := range report.ByCategory {
		result = append(result, fmt.Sprintf("%s:%d", category, amount))
	}
	return result
}

func StatusLabel(status domain.MaintenanceStatus) string {
	switch status {
	case domain.StatusReady:
		return "ready for rental"
	case domain.StatusMaintenance:
		return "in maintenance"
	default:
		return "retired"
	}
}
