package staff

import (
	"fmt"
	"sort"
)

type Shift struct {
	StaffID string
	Date    string
	Start   string
	End     string
	Area    string
}

func ValidateShift(shift Shift) error {
	if shift.StaffID == "" || shift.Date == "" || shift.Start == "" || shift.End == "" || shift.Area == "" {
		return fmt.Errorf("incomplete shift")
	}
	if shift.Start >= shift.End {
		return fmt.Errorf("shift end must follow start")
	}
	return nil
}

func BuildRoster(shifts []Shift) (map[string][]Shift, error) {
	roster := make(map[string][]Shift)
	for _, shift := range shifts {
		if err := ValidateShift(shift); err != nil {
			return nil, err
		}
		roster[shift.Date] = append(roster[shift.Date], shift)
	}
	for date := range roster {
		sort.SliceStable(roster[date], func(a, b int) bool {
			if roster[date][a].Start == roster[date][b].Start {
				return roster[date][a].StaffID < roster[date][b].StaffID
			}
			return roster[date][a].Start < roster[date][b].Start
		})
	}
	return roster, nil
}

func HasCoverage(roster map[string][]Shift, date, area string) bool {
	for _, shift := range roster[date] {
		if shift.Area == area {
			return true
		}
	}
	return false
}

func AreasForDate(roster map[string][]Shift, date string) []string {
	seen := make(map[string]bool)
	for _, shift := range roster[date] {
		seen[shift.Area] = true
	}
	areas := make([]string, 0, len(seen))
	for area := range seen {
		areas = append(areas, area)
	}
	sort.Strings(areas)
	return areas
}
