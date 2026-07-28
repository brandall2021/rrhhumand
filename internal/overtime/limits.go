package overtime

type Limits struct {
	rounding *Rounding
}

func NewLimits() *Limits {
	return &Limits{rounding: NewRounding()}
}

type LimitCheck struct {
	DailyAllowed    int  `json:"daily_allowed"`
	DailyExcess     int  `json:"daily_excess"`
	WeeklyAllowed   int  `json:"weekly_allowed"`
	WeeklyExcess    int  `json:"weekly_excess"`
	MonthlyAllowed  int  `json:"monthly_allowed"`
	MonthlyExcess   int  `json:"monthly_excess"`
	ExceedsPolicy   bool `json:"exceeds_policy"`
}

func (l *Limits) CheckDaily(overtimeMinutes int, policy *OvertimePolicy) (int, int) {
	if policy.MaxDailyMinutes <= 0 {
		return overtimeMinutes, 0
	}
	if overtimeMinutes <= policy.MaxDailyMinutes {
		return overtimeMinutes, 0
	}
	return policy.MaxDailyMinutes, overtimeMinutes - policy.MaxDailyMinutes
}

func (l *Limits) CheckWeekly(dailyMinutes int, weeklyTotal int, policy *OvertimePolicy) (int, int) {
	if policy.MaxWeeklyMinutes <= 0 {
		return dailyMinutes, 0
	}
	projected := weeklyTotal + dailyMinutes
	if projected <= policy.MaxWeeklyMinutes {
		return dailyMinutes, 0
	}
	remaining := policy.MaxWeeklyMinutes - weeklyTotal
	if remaining < 0 {
		remaining = 0
	}
	return remaining, dailyMinutes - remaining
}

func (l *Limits) ApplyLimits(overtimeMinutes int, weeklyTotal int, monthlyTotal int, policy *OvertimePolicy) (allowed int, excess int) {
	dailyAllowed, dailyExcess := l.CheckDaily(overtimeMinutes, policy)
	_ = dailyExcess

	weeklyAllowed, weeklyExcess := l.CheckWeekly(dailyAllowed, weeklyTotal, policy)
	_ = weeklyExcess

	allowed = weeklyAllowed
	if policy.MaxMonthlyMinutes > 0 {
		projected := monthlyTotal + allowed
		if projected > policy.MaxMonthlyMinutes {
			monthlyRemaining := policy.MaxMonthlyMinutes - monthlyTotal
			if monthlyRemaining < 0 {
				monthlyRemaining = 0
			}
			if monthlyRemaining < allowed {
				allowed = monthlyRemaining
			}
		}
	}

	excess = overtimeMinutes - allowed
	if excess < 0 {
		excess = 0
	}

	rounded := l.rounding.RoundDown(allowed, policy.RoundingMinutes)
	if rounded > allowed {
		rounded = allowed
	}

	return rounded, excess
}
