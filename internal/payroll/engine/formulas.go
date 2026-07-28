package engine

import (
	"time"

	"github.com/shopspring/decimal"
)

func CalcHourlyValue(baseSalary decimal.Decimal, workDays int) decimal.Decimal {
	if workDays <= 0 {
		workDays = 30
	}
	monthlyHours := decimal.NewFromInt(int64(workDays) * 8)
	if monthlyHours.IsZero() {
		return decimal.Zero
	}
	return baseSalary.Div(monthlyHours).Round(2)
}

func CalcDailyValue(baseSalary decimal.Decimal, workDays int) decimal.Decimal {
	if workDays <= 0 {
		workDays = 30
	}
	return baseSalary.Div(decimal.NewFromInt(int64(workDays))).Round(2)
}

func CalcWorkDaysInPeriod(start, end time.Time) int {
	count := 0
	current := start
	for !current.After(end) {
		if current.Weekday() != time.Sunday {
			count++
		}
		current = current.AddDate(0, 0, 1)
	}
	return count
}

func ProportionalPart(annualAmount decimal.Decimal, monthsWorked int) decimal.Decimal {
	if monthsWorked <= 0 {
		return decimal.Zero
	}
	if monthsWorked > 12 {
		monthsWorked = 12
	}
	return annualAmount.Mul(decimal.NewFromInt(int64(monthsWorked))).Div(decimal.NewFromInt(12)).Round(2)
}

func RoundMoney(amount decimal.Decimal) decimal.Decimal {
	return amount.Round(2)
}
