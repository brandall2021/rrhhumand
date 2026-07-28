package engine

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type SummaryReportData struct {
	CompanyID           uuid.UUID
	CompanyName         string
	PeriodName          string
	PeriodStart         time.Time
	PeriodEnd           time.Time
	TotalEmployees      int
	TotalRemunerative   decimal.Decimal
	TotalNonRemunerative decimal.Decimal
	TotalDeductions     decimal.Decimal
	TotalContributions  decimal.Decimal
	TotalNetAmount      decimal.Decimal
	TotalEmployerCost   decimal.Decimal
	GeneratedAt         time.Time
}

type DetailedReportItem struct {
	EmployeeName        string
	EmployeeCUIL        string
	EmployeeCategory    string
	RemunerativeAmount  decimal.Decimal
	NonRemunerativeAmount decimal.Decimal
	DeductionsTotal     decimal.Decimal
	ContributionsTotal  decimal.Decimal
	NetAmount           decimal.Decimal
	EmployerCost        decimal.Decimal
	Items               []DetailedConceptItem
}

type DetailedConceptItem struct {
	ConceptCode string
	ConceptName string
	Amount      decimal.Decimal
}

type DetailedReportData struct {
	CompanyID   uuid.UUID
	CompanyName string
	PeriodName  string
	PeriodStart time.Time
	PeriodEnd   time.Time
	Employees   []DetailedReportItem
	GeneratedAt time.Time
}

type ComparativePeriodData struct {
	PeriodName        string
	TotalEmployees    int
	TotalRemunerative decimal.Decimal
	TotalDeductions   decimal.Decimal
	TotalContributions decimal.Decimal
	TotalNetAmount    decimal.Decimal
}

type ComparativeReportData struct {
	CompanyID   uuid.UUID
	CompanyName string
	Periods     []ComparativePeriodData
	GeneratedAt time.Time
}

type ReportGenerator struct{}

func NewReportGenerator() *ReportGenerator {
	return &ReportGenerator{}
}

func (g *ReportGenerator) GenerateSummaryReport(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, data SummaryReportData) (string, error) {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Summary Report - %s\r\n", data.PeriodName))
	b.WriteString(fmt.Sprintf("Company: %s\r\n", data.CompanyName))
	b.WriteString(fmt.Sprintf("Period: %s - %s\r\n", data.PeriodStart.Format("02/01/2006"), data.PeriodEnd.Format("02/01/2006")))
	b.WriteString(fmt.Sprintf("Generated: %s\r\n", data.GeneratedAt.Format("02/01/2006 15:04:05")))
	b.WriteString("---\r\n")
	b.WriteString(fmt.Sprintf("Total Employees: %d\r\n", data.TotalEmployees))
	b.WriteString(fmt.Sprintf("Total Remunerative: %s\r\n", data.TotalRemunerative.StringFixed(2)))
	b.WriteString(fmt.Sprintf("Total Non-Remunerative: %s\r\n", data.TotalNonRemunerative.StringFixed(2)))
	b.WriteString(fmt.Sprintf("Total Deductions: %s\r\n", data.TotalDeductions.StringFixed(2)))
	b.WriteString(fmt.Sprintf("Total Contributions: %s\r\n", data.TotalContributions.StringFixed(2)))
	b.WriteString(fmt.Sprintf("Total Net Amount: %s\r\n", data.TotalNetAmount.StringFixed(2)))
	b.WriteString(fmt.Sprintf("Total Employer Cost: %s\r\n", data.TotalEmployerCost.StringFixed(2)))
	return b.String(), nil
}

func (g *ReportGenerator) GenerateDetailedReport(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, data DetailedReportData) (string, error) {
	var b strings.Builder
	b.WriteString("Employee;CUIL;Category;Remunerative;NonRemunerative;Deductions;Contributions;NetAmount;EmployerCost\r\n")
	for _, emp := range data.Employees {
		b.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%s;%s\r\n",
			emp.EmployeeName,
			emp.EmployeeCUIL,
			emp.EmployeeCategory,
			emp.RemunerativeAmount.StringFixed(2),
			emp.NonRemunerativeAmount.StringFixed(2),
			emp.DeductionsTotal.StringFixed(2),
			emp.ContributionsTotal.StringFixed(2),
			emp.NetAmount.StringFixed(2),
			emp.EmployerCost.StringFixed(2),
		))
	}
	return b.String(), nil
}

func (g *ReportGenerator) GenerateDetailedHTML(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, data DetailedReportData) (string, error) {
	tmpl := template.Must(template.New("detailed").Parse(detailedHTMLTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("report_generator: execute detailed template: %w", err)
	}
	return buf.String(), nil
}

func (g *ReportGenerator) GenerateComparativeReport(ctx context.Context, companyID uuid.UUID, periodIDs []uuid.UUID, data ComparativeReportData) (string, error) {
	var b strings.Builder
	b.WriteString("Period;Employees;Remunerative;Deductions;Contributions;NetAmount\r\n")
	for _, p := range data.Periods {
		b.WriteString(fmt.Sprintf("%s;%d;%s;%s;%s;%s\r\n",
			p.PeriodName,
			p.TotalEmployees,
			p.TotalRemunerative.StringFixed(2),
			p.TotalDeductions.StringFixed(2),
			p.TotalContributions.StringFixed(2),
			p.TotalNetAmount.StringFixed(2),
		))
	}
	return b.String(), nil
}

func (g *ReportGenerator) GenerateSummaryHTML(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, data SummaryReportData) (string, error) {
	tmpl := template.Must(template.New("summary").Parse(summaryHTMLTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("report_generator: execute summary template: %w", err)
	}
	return buf.String(), nil
}

func (g *ReportGenerator) GenerateComparativeHTML(ctx context.Context, companyID uuid.UUID, periodIDs []uuid.UUID, data ComparativeReportData) (string, error) {
	tmpl := template.Must(template.New("comparative").Parse(comparativeHTMLTemplate))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("report_generator: execute comparative template: %w", err)
	}
	return buf.String(), nil
}

var summaryHTMLTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Summary Report</title></head>
<body>
<h1>Summary Report</h1>
<p>Company: {{.CompanyName}}</p>
<p>Period: {{.PeriodName}}</p>
<p>Generated: {{.GeneratedAt.Format "02/01/2006 15:04:05"}}</p>
<hr>
<table border="1">
<tr><th>Metric</th><th>Value</th></tr>
<tr><td>Total Employees</td><td>{{.TotalEmployees}}</td></tr>
<tr><td>Total Remunerative</td><td>{{.TotalRemunerative.StringFixed 2}}</td></tr>
<tr><td>Total Non-Remunerative</td><td>{{.TotalNonRemunerative.StringFixed 2}}</td></tr>
<tr><td>Total Deductions</td><td>{{.TotalDeductions.StringFixed 2}}</td></tr>
<tr><td>Total Contributions</td><td>{{.TotalContributions.StringFixed 2}}</td></tr>
<tr><td>Total Net Amount</td><td>{{.TotalNetAmount.StringFixed 2}}</td></tr>
<tr><td>Total Employer Cost</td><td>{{.TotalEmployerCost.StringFixed 2}}</td></tr>
</table>
</body>
</html>`

var detailedHTMLTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Detailed Report</title></head>
<body>
<h1>Detailed Report</h1>
<p>Company: {{.CompanyName}}</p>
<p>Period: {{.PeriodName}}</p>
<table border="1">
<tr><th>Employee</th><th>CUIL</th><th>Category</th><th>Remunerative</th><th>Net Amount</th></tr>
{{range .Employees}}
<tr>
<td>{{.EmployeeName}}</td>
<td>{{.EmployeeCUIL}}</td>
<td>{{.EmployeeCategory}}</td>
<td>{{.RemunerativeAmount.StringFixed 2}}</td>
<td>{{.NetAmount.StringFixed 2}}</td>
</tr>
{{end}}
</table>
</body>
</html>`

var comparativeHTMLTemplate = `<!DOCTYPE html>
<html>
<head><meta charset="utf-8"><title>Comparative Report</title></head>
<body>
<h1>Comparative Report</h1>
<p>Company: {{.CompanyName}}</p>
<table border="1">
<tr><th>Period</th><th>Employees</th><th>Remunerative</th><th>Deductions</th><th>Contributions</th><th>Net Amount</th></tr>
{{range .Periods}}
<tr>
<td>{{.PeriodName}}</td>
<td>{{.TotalEmployees}}</td>
<td>{{.TotalRemunerative.StringFixed 2}}</td>
<td>{{.TotalDeductions.StringFixed 2}}</td>
<td>{{.TotalContributions.StringFixed 2}}</td>
<td>{{.TotalNetAmount.StringFixed 2}}</td>
</tr>
{{end}}
</table>
</body>
</html>`
