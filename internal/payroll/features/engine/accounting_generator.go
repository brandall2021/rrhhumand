package engine

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

type AccountingGenerator struct{}

func NewAccountingGenerator() *AccountingGenerator {
	return &AccountingGenerator{}
}

func (g *AccountingGenerator) GenerateCSV(entries []domain.AccountingEntry) (string, error) {
	var b strings.Builder
	writer := csv.NewWriter(&b)
	writer.Write([]string{"AccountCode", "AccountName", "CostCenter", "DebitAmount", "CreditAmount", "EmployeeName", "ConceptName", "Reference"})
	for _, e := range entries {
		acctName := ""
		if e.AccountName != nil {
			acctName = *e.AccountName
		}
		costCenter := ""
		if e.CostCenter != nil {
			costCenter = *e.CostCenter
		}
		empName := ""
		if e.EmployeeName != nil {
			empName = *e.EmployeeName
		}
		conceptName := ""
		if e.ConceptName != nil {
			conceptName = *e.ConceptName
		}
		ref := ""
		if e.Reference != nil {
			ref = *e.Reference
		}
		writer.Write([]string{
			e.AccountCode,
			acctName,
			costCenter,
			e.DebitAmount.StringFixed(2),
			e.CreditAmount.StringFixed(2),
			empName,
			conceptName,
			ref,
		})
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("accounting_generator: csv write: %w", err)
	}
	return b.String(), nil
}

func (g *AccountingGenerator) GenerateSIAPContable(entries []domain.AccountingEntry) (string, error) {
	var b strings.Builder
	for i, e := range entries {
		acctName := ""
		if e.AccountName != nil {
			acctName = *e.AccountName
		}
		conceptCode := ""
		if e.ConceptCode != nil {
			conceptCode = *e.ConceptCode
		}
		line := fmt.Sprintf("AS%d;%s;%s;%s;%s;%s",
			i+1,
			e.AccountCode,
			acctName,
			e.DebitAmount.StringFixed(2),
			e.CreditAmount.StringFixed(2),
			conceptCode,
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	return b.String(), nil
}

func (g *AccountingGenerator) GenerateTXT(entries []domain.AccountingEntry) (string, error) {
	var b strings.Builder
	for _, e := range entries {
		acctName := ""
		if e.AccountName != nil {
			acctName = *e.AccountName
		}
		costCenter := ""
		if e.CostCenter != nil {
			costCenter = *e.CostCenter
		}
		empName := ""
		if e.EmployeeName != nil {
			empName = *e.EmployeeName
		}
		conceptName := ""
		if e.ConceptName != nil {
			conceptName = *e.ConceptName
		}
		line := fmt.Sprintf("%s%s%s%s%s%s%s",
			padRight(e.AccountCode, 20),
			padRight(acctName, 40),
			padRight(costCenter, 10),
			padRight(e.DebitAmount.StringFixed(2), 15),
			padRight(e.CreditAmount.StringFixed(2), 15),
			padRight(empName, 40),
			padRight(conceptName, 40),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	return b.String(), nil
}
