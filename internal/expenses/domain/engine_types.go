package domain

import "github.com/shopspring/decimal"

type PolicyResult struct {
	Compliant bool
	Reasons   []string
	Violations []PolicyViolation
}

type PolicyViolation struct {
	RuleID      string
	Category    string
	Limit       decimal.Decimal
	Actual      decimal.Decimal
	Excess      decimal.Decimal
	Message     string
}

type SettlementResult struct {
	TotalExpenses  decimal.Decimal
	AdvanceAmount  decimal.Decimal
	CompanyOwes    decimal.Decimal
	EmployeeOwes   decimal.Decimal
	Currency       string
}

type ExpenseContext struct {
	EmployeeID     string
	CompanyID      string
	Category       string
	Amount         decimal.Decimal
	Currency       string
	ExpenseDate    string
	PaymentMethod  string
	EmployeeCategory string
	HasReceipt     bool
}
