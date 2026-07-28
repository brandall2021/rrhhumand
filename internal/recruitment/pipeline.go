package recruitment

import (
	"fmt"
)

var validStages = []string{
	"NEW",
	"SCREENING",
	"PHONE_INTERVIEW",
	"TECHNICAL_INTERVIEW",
	"HR_INTERVIEW",
	"ASSESSMENT",
	"FINALIST",
	"OFFER",
	"HIRED",
}

var validExits = map[string][]string{
	"NEW":               {"SCREENING", "REJECTED", "WITHDRAWN"},
	"SCREENING":         {"PHONE_INTERVIEW", "REJECTED", "WITHDRAWN", "ON_HOLD"},
	"PHONE_INTERVIEW":   {"TECHNICAL_INTERVIEW", "REJECTED", "WITHDRAWN", "ON_HOLD"},
	"TECHNICAL_INTERVIEW": {"HR_INTERVIEW", "REJECTED", "WITHDRAWN", "ON_HOLD"},
	"HR_INTERVIEW":      {"ASSESSMENT", "FINALIST", "REJECTED", "WITHDRAWN", "ON_HOLD"},
	"ASSESSMENT":        {"FINALIST", "REJECTED", "WITHDRAWN"},
	"FINALIST":          {"OFFER", "REJECTED", "WITHDRAWN"},
	"OFFER":             {"HIRED", "REJECTED", "WITHDRAWN"},
	"HIRED":             {},
	"REJECTED":          {},
	"WITHDRAWN":         {},
	"ON_HOLD":           {"SCREENING", "PHONE_INTERVIEW", "TECHNICAL_INTERVIEW", "HR_INTERVIEW", "ASSESSMENT", "FINALIST", "OFFER", "REJECTED"},
}

func ValidateStageTransition(from, to string) error {
	allowed, ok := validExits[from]
	if !ok {
		return fmt.Errorf("unknown stage: %s", from)
	}
	for _, s := range allowed {
		if s == to {
			return nil
		}
	}
	return fmt.Errorf("invalid transition from %s to %s", from, to)
}

func IsTerminalStage(stage string) bool {
	return stage == "HIRED" || stage == "REJECTED" || stage == "WITHDRAWN"
}

func GetStages() []string {
	out := make([]string, len(validStages))
	copy(out, validStages)
	return out
}
