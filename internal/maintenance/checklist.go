package maintenance

import (
	"fmt"
	"strings"
)

type Inspection struct {
	StructuralOK bool
	Clean        bool
	Accessories  int
	Notes        string
}

type InspectionDecision struct {
	Ready   bool
	Reasons []string
}

func EvaluateInspection(inspection Inspection) InspectionDecision {
	decision := InspectionDecision{Ready: true, Reasons: make([]string, 0)}
	if !inspection.StructuralOK {
		decision.Ready = false
		decision.Reasons = append(decision.Reasons, "structural inspection failed")
	}
	if !inspection.Clean {
		decision.Ready = false
		decision.Reasons = append(decision.Reasons, "cleaning required")
	}
	if inspection.Accessories < 1 {
		decision.Ready = false
		decision.Reasons = append(decision.Reasons, "accessories incomplete")
	}
	return decision
}

func InspectionNote(inspection Inspection) (string, error) {
	if strings.TrimSpace(inspection.Notes) == "" {
		return "", fmt.Errorf("inspection notes are required")
	}
	decision := EvaluateInspection(inspection)
	if decision.Ready {
		return "ready: " + strings.TrimSpace(inspection.Notes), nil
	}
	return "hold: " + strings.Join(decision.Reasons, "; "), nil
}
