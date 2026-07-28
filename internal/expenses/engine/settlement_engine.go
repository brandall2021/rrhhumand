package engine

import (
	"github.com/rrhhumand/api/internal/expenses/domain"
	"github.com/shopspring/decimal"
)

type SettlementEngine struct{}

func NewSettlementEngine() *SettlementEngine {
	return &SettlementEngine{}
}

func (e *SettlementEngine) Calculate(advance decimal.Decimal, expenses []domain.Expense) domain.SettlementResult {
	total := decimal.Zero
	currency := ""
	for _, exp := range expenses {
		total = total.Add(exp.BaseAmount)
		if currency == "" {
			currency = exp.BaseCurrency
		}
	}

	result := domain.SettlementResult{
		TotalExpenses: total,
		AdvanceAmount: advance,
		Currency:      currency,
	}

	if total.GreaterThan(advance) {
		result.CompanyOwes = total.Sub(advance)
		result.EmployeeOwes = decimal.Zero
	} else {
		result.CompanyOwes = decimal.Zero
		result.EmployeeOwes = advance.Sub(total)
	}

	return result
}
