package engine

import (
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/rrhhumand/api/internal/payroll/features/domain"
)

type BankGenerator struct{}

func NewBankGenerator() *BankGenerator {
	return &BankGenerator{}
}

func (g *BankGenerator) GenerateCSV(items []domain.BankBatchItem) (string, error) {
	var b strings.Builder
	writer := csv.NewWriter(&b)
	writer.Write([]string{"CUIL", "Surname", "Name", "CBU", "AccountType", "Amount", "Concept"})
	for _, item := range items {
		cbu := ""
		if item.CBU != nil {
			cbu = *item.CBU
		}
		accType := ""
		if item.AccountType != nil {
			accType = *item.AccountType
		}
		concept := ""
		if item.Concept != nil {
			concept = *item.Concept
		}
		writer.Write([]string{
			item.CUIL,
			item.Surname,
			item.Name,
			cbu,
			accType,
			item.Amount.StringFixed(2),
			concept,
		})
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("bank_generator: csv write: %w", err)
	}
	return b.String(), nil
}

func (g *BankGenerator) GenerateTXT(items []domain.BankBatchItem) (string, error) {
	var b strings.Builder
	for _, item := range items {
		cbu := ""
		if item.CBU != nil {
			cbu = *item.CBU
		}
		accType := ""
		if item.AccountType != nil {
			accType = *item.AccountType
		}
		line := fmt.Sprintf("%s%s%s%s%s%s",
			padRight(item.CUIL, 11),
			padRight(item.Surname, 30),
			padRight(item.Name, 30),
			padRight(cbu, 22),
			padRight(accType, 10),
			padRight(item.Amount.StringFixed(2), 15),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	return b.String(), nil
}

func (g *BankGenerator) GenerateNacion(items []domain.BankBatchItem) (string, error) {
	var b strings.Builder
	for _, item := range items {
		cbu := ""
		if item.CBU != nil {
			cbu = *item.CBU
		}
		line := fmt.Sprintf("01%s%s%015s%s",
			padRight(cbu, 22),
			padRight(item.Surname+", "+item.Name, 40),
			item.Amount.StringFixed(2),
			padRight(item.CUIL, 11),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	return b.String(), nil
}

func (g *BankGenerator) GenerateSantander(items []domain.BankBatchItem) (string, error) {
	var b strings.Builder
	for _, item := range items {
		cbu := ""
		if item.CBU != nil {
			cbu = *item.CBU
		}
		line := fmt.Sprintf("S%s%s%015s%s",
			padRight(cbu, 22),
			padRight(item.CUIL, 11),
			item.Amount.StringFixed(2),
			padRight(item.Surname+" "+item.Name, 40),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	return b.String(), nil
}

func (g *BankGenerator) GenerateGalicia(items []domain.BankBatchItem) (string, error) {
	var b strings.Builder
	for _, item := range items {
		cbu := ""
		if item.CBU != nil {
			cbu = *item.CBU
		}
		line := fmt.Sprintf("G%s%s%015s%s",
			padRight(item.CUIL, 11),
			padRight(cbu, 22),
			item.Amount.StringFixed(2),
			padRight(item.Surname+", "+item.Name, 40),
		)
		b.WriteString(line)
		b.WriteString("\r\n")
	}
	return b.String(), nil
}
