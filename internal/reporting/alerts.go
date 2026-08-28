package reporting

import (
	"context"
	"fmt"
	"sort"

	"campgear/internal/domain"
	"campgear/internal/storage"
)

type AlertLevel string

const (
	AlertInfo     AlertLevel = "info"
	AlertWarning  AlertLevel = "warning"
	AlertCritical AlertLevel = "critical"
)

type OperationalAlert struct {
	Level    AlertLevel `json:"level"`
	Code     string     `json:"code"`
	EntityID string     `json:"entity_id"`
	Message  string     `json:"message"`
}

func BuildAlerts(ctx context.Context, repo *storage.Repository) ([]OperationalAlert, error) {
	items, err := repo.ListItems(ctx, "", false)
	if err != nil {
		return nil, err
	}
	alerts := make([]OperationalAlert, 0)
	for _, item := range items {
		if item.MaintenanceStatus == domain.StatusMaintenance {
			alerts = append(alerts, OperationalAlert{Level: AlertWarning, Code: "MAINTENANCE", EntityID: item.ID, Message: item.Name + " is unavailable for rental"})
		}
		if item.Available == 0 && item.MaintenanceStatus == domain.StatusReady {
			alerts = append(alerts, OperationalAlert{Level: AlertCritical, Code: "OUT_OF_STOCK", EntityID: item.ID, Message: item.Name + " has no available units"})
		}
		if item.Available < item.Stock && item.MaintenanceStatus == domain.StatusReady {
			alerts = append(alerts, OperationalAlert{Level: AlertInfo, Code: "LOW_AVAILABILITY", EntityID: item.ID, Message: fmt.Sprintf("%d of %d units available", item.Available, item.Stock)})
		}
	}
	orders, err := repo.ListMaintenance(ctx, "open")
	if err != nil {
		return nil, err
	}
	for _, order := range orders {
		alerts = append(alerts, OperationalAlert{Level: AlertWarning, Code: "OPEN_WORK_ORDER", EntityID: order.ID, Message: order.Reason})
	}
	sort.SliceStable(alerts, func(a, b int) bool {
		if alerts[a].Level == alerts[b].Level {
			if alerts[a].Code == alerts[b].Code {
				return alerts[a].EntityID < alerts[b].EntityID
			}
			return alerts[a].Code < alerts[b].Code
		}
		return alertRank(alerts[a].Level) > alertRank(alerts[b].Level)
	})
	return alerts, nil
}

func alertRank(level AlertLevel) int {
	switch level {
	case AlertCritical:
		return 3
	case AlertWarning:
		return 2
	default:
		return 1
	}
}

func AlertsByLevel(alerts []OperationalAlert, level AlertLevel) []OperationalAlert {
	result := make([]OperationalAlert, 0)
	for _, alert := range alerts {
		if alert.Level == level {
			result = append(result, alert)
		}
	}
	return result
}
