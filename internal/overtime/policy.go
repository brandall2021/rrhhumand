package overtime

import (
	"time"
)

type PolicyEngine struct{}

func NewPolicyEngine() *PolicyEngine {
	return &PolicyEngine{}
}

func (pe *PolicyEngine) IsNightShift(startHour, endHour int, nightStart, nightEnd string) bool {
	ns, _ := time.Parse("15:04", nightStart)
	ne, _ := time.Parse("15:04", nightEnd)

	nsHour := ns.Hour()
	neHour := ne.Hour()

	if nsHour > neHour {
		return startHour >= nsHour || endHour <= neHour
	}
	return startHour >= nsHour && endHour <= neHour
}

func (pe *PolicyEngine) CalculateNightMinutes(startMin, endMin, totalMinutes int, nightStartMin, nightEndMin int) int {
	if nightStartMin > nightEndMin {
		night1Start := nightStartMin
		night1End := 24 * 60
		night2Start := 0
		night2End := nightEndMin

		minutes := 0
		if startMin < night1End && endMin > night1Start {
			s := max(startMin, night1Start)
			e := min(endMin, night1End)
			minutes += e - s
		}
		if startMin < night2End && endMin > night2Start {
			s := max(startMin, night2Start)
			e := min(endMin, night2End)
			minutes += e - s
		}
		return minutes
	}

	if startMin < nightEndMin && endMin > nightStartMin {
		s := max(startMin, nightStartMin)
		e := min(endMin, nightEndMin)
		return e - s
	}
	return 0
}

func (pe *PolicyEngine) GetMultiplier(overtimeType string, policy *OvertimePolicy) float64 {
	if policy == nil {
		return 1.0
	}
	switch overtimeType {
	case "WEEKEND":
		return policy.WeekendMultiplier
	case "HOLIDAY":
		return policy.HolidayMultiplier
	case "NIGHT":
		return policy.NightMultiplier
	default:
		return 1.0
	}
}
