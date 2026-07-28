package engine

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/shopspring/decimal"
)

type ReceiptGenerator struct {
	templates map[string]string
}

func NewReceiptGenerator() *ReceiptGenerator {
	return &ReceiptGenerator{
		templates: make(map[string]string),
	}
}

type ReceiptData struct {
	CompanyName           string
	CompanyCUIT           string
	CompanyAddress        string
	EmployeeName          string
	EmployeeCUIL          string
	EmployeeCategory      string
	PeriodName            string
	PeriodStart           string
	PeriodEnd             string
	PaymentDate           string
	ReceiptNumber         string
	CUIT                  string
	Items                 []ReceiptItemData
	GrossRemunerative     string
	GrossNonRemunerative  string
	DeductionsTotal       string
	ContributionsTotal    string
	NetAmount             string
	EmployerCost          string
	AmountInWords         string
	Currency              string
}

type ReceiptItemData struct {
	Code           string
	Name           string
	Quantity       string
	UnitValue      string
	Amount         string
	IsDeduction    bool
	IsContribution bool
}

func (g *ReceiptGenerator) GenerateHTML(data *ReceiptData, templateHTML string) (string, error) {
	tmpl, err := template.New("receipt").Funcs(template.FuncMap{
		"add": func(a, b int) int { return a + b },
	}).Parse(templateHTML)
	if err != nil {
		return "", fmt.Errorf("receipt_generator: parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("receipt_generator: execute template: %w", err)
	}
	return buf.String(), nil
}

func (g *ReceiptGenerator) GenerateReceiptNumber(companyCode string, year, month, seq int) string {
	return fmt.Sprintf("%s-%04d-%02d-%04d", companyCode, year, month, seq)
}

func FormatDecimal(d decimal.Decimal) string {
	return d.StringFixed(2)
}
