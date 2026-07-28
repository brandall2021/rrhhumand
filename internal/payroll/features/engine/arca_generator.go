package engine

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ArcaExportEntry struct {
	CUIL                string
	Surname             string
	Name                string
	ConceptCode         string
	RemunerativeAmount  decimal.Decimal
	DeductionAmount     decimal.Decimal
	ContributionAmount  decimal.Decimal
	DaysWorked          int
	CategoryCode        string
	AgreementCode       string
}

type ArcaGenerator struct{}

func NewArcaGenerator() *ArcaGenerator {
	return &ArcaGenerator{}
}

func (g *ArcaGenerator) GenerateSICOSS(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, entries []ArcaExportEntry) (string, string) {
	var b strings.Builder
	now := time.Now()
	period := now.Format("200601")
	date := now.Format("20060102")

	header := fmt.Sprintf("01%s%s%s", padRight(companyID.String(), 11), period, date)
	b.WriteString(header)
	b.WriteString("\r\n")

	var totalRem, totalDed, totalCont decimal.Decimal
	for _, e := range entries {
		rem := e.RemunerativeAmount.StringFixed(2)
		ded := e.DeductionAmount.StringFixed(2)
		cont := e.ContributionAmount.StringFixed(2)
		line := fmt.Sprintf("02%s%s%s%s%s%s%02d%s%s",
			padRight(e.CUIL, 11),
			padRight(e.Surname, 30),
			padRight(e.Name, 30),
			padRight(rem, 13),
			padRight(ded, 13),
			padRight(cont, 13),
			e.DaysWorked,
			padRight(e.CategoryCode, 6),
			padRight(e.AgreementCode, 6),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
		totalRem = totalRem.Add(e.RemunerativeAmount)
		totalDed = totalDed.Add(e.DeductionAmount)
		totalCont = totalCont.Add(e.ContributionAmount)
	}

	content := b.String()
	checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))

	trailer := fmt.Sprintf("03%06d%s%s%s%s",
		len(entries),
		padRight(totalRem.StringFixed(2), 15),
		padRight(totalDed.StringFixed(2), 15),
		padRight(totalCont.StringFixed(2), 15),
		checksum,
	)
	b.WriteString(trailer)
	b.WriteString("\r\n")

	return b.String(), checksum
}

func (g *ArcaGenerator) GenerateSIAP(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, entries []ArcaExportEntry) (string, string) {
	var b strings.Builder
	b.WriteString("CUIL;Surname;Name;ConceptCode;Remunerative;Deduction;Contribution;Days;Category;Agreement\r\n")
	for _, e := range entries {
		b.WriteString(fmt.Sprintf("%s;%s;%s;%s;%s;%s;%s;%d;%s;%s\r\n",
			e.CUIL, e.Surname, e.Name, e.ConceptCode,
			e.RemunerativeAmount.StringFixed(2),
			e.DeductionAmount.StringFixed(2),
			e.ContributionAmount.StringFixed(2),
			e.DaysWorked, e.CategoryCode, e.AgreementCode,
		))
	}
	content := b.String()
	checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	return content, checksum
}

func (g *ArcaGenerator) GenerateF931(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, entries []ArcaExportEntry) (string, string) {
	var b strings.Builder
	now := time.Now()
	for _, e := range entries {
		line := fmt.Sprintf("%s%s%s%s%s",
			padRight(e.CUIL, 11),
			padRight(e.Surname, 30),
			padRight(e.Name, 30),
			padRight(e.RemunerativeAmount.StringFixed(2), 13),
			now.Format("200601"),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	content := b.String()
	checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	return content, checksum
}

func (g *ArcaGenerator) GenerateF1357(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, entries []ArcaExportEntry) (string, string) {
	var b strings.Builder
	now := time.Now()
	var total decimal.Decimal
	for _, e := range entries {
		amt := e.RemunerativeAmount.Add(e.DeductionAmount)
		line := fmt.Sprintf("%s%s%s%s%s%s",
			padRight(e.CUIL, 11),
			padRight(e.Surname, 30),
			padRight(e.Name, 30),
			padRight(e.ConceptCode, 6),
			padRight(amt.StringFixed(2), 13),
			now.Format("200601"),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
		total = total.Add(amt)
	}
	trailer := fmt.Sprintf("99%06d%s", len(entries), padRight(total.StringFixed(2), 15))
	b.WriteString(trailer)
	b.WriteString("\r\n")
	content := b.String()
	checksum := fmt.Sprintf("%x", sha256.Sum256([]byte(content)))
	return content, checksum
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s[:n]
	}
	return s + strings.Repeat(" ", n-len(s))
}
