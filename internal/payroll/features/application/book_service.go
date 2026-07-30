package application

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/payroll/features/domain"
	"github.com/rrhhumand/api/internal/payroll/features/repository"
)

type BookService struct {
	bookRepo *repository.BookRepo
}

func NewBookService(bookRepo *repository.BookRepo) *BookService {
	return &BookService{bookRepo: bookRepo}
}

func bookSvcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("book_svc.%s: %w", op, err)
}

func (s *BookService) GenerateBookEntries(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, userID uuid.UUID) ([]domain.BookEntry, error) {
	e := &domain.BookEntry{
		ID:                uuid.New(),
		CompanyID:         companyID,
		RunID:             runID,
		EntryType:         "REGULAR",
		GrossRemunerative: decimal.Zero,
		GrossNonRemunerative: decimal.Zero,
		DeductionsTotal:   decimal.Zero,
		ContributionsTotal: decimal.Zero,
		NetAmount:         decimal.Zero,
		EmployerCost:      decimal.Zero,
		DaysWorked:        0,
		HoursWorked:       decimal.Zero,
		Absences:          0,
		Status:            "GENERATED",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := s.bookRepo.CreateEntry(ctx, e); err != nil {
		return nil, bookSvcErr("GenerateBookEntries", err)
	}
	return []domain.BookEntry{*e}, nil
}

func (s *BookService) GetBookEntries(ctx context.Context, companyID uuid.UUID, runID uuid.UUID) ([]domain.BookEntry, error) {
	return s.bookRepo.ListByRun(ctx, runID)
}

func (s *BookService) GetEmployeeBookEntry(ctx context.Context, companyID uuid.UUID, runID uuid.UUID, employeeID uuid.UUID) (*domain.BookEntry, error) {
	entries, err := s.bookRepo.ListByEmployee(ctx, companyID, employeeID, 1, 0)
	if err != nil {
		return nil, bookSvcErr("GetEmployeeBookEntry", err)
	}
	if len(entries) == 0 {
		return nil, fmt.Errorf("book_svc.GetEmployeeBookEntry: entry not found")
	}
	return &entries[0], nil
}

func (s *BookService) ExportBook(ctx context.Context, companyID uuid.UUID, periodID *uuid.UUID, year, month int, format string, userID uuid.UUID) (*domain.BookExport, error) {
	e := &domain.BookExport{
		ID:            uuid.New(),
		CompanyID:     companyID,
		PeriodID:      periodID,
		Year:          year,
		Month:         month,
		ExportType:    format,
		FileName:      fmt.Sprintf("LIBRO_SUELDOS_%d_%02d.%s", year, month, format),
		Status:        "GENERATED",
		EmployeeCount: 0,
		TotalGross:    decimal.Zero,
		TotalDeductions: decimal.Zero,
		TotalNet:      decimal.Zero,
		GeneratedBy:   userID,
		GeneratedAt:   time.Now(),
		CreatedAt:     time.Now(),
	}
	if err := s.bookRepo.CreateExport(ctx, e); err != nil {
		return nil, bookSvcErr("ExportBook", err)
	}
	return e, nil
}

func (s *BookService) GetBookExports(ctx context.Context, companyID uuid.UUID) ([]domain.BookExport, error) {
	return s.bookRepo.ListExports(ctx, companyID, 0, 0)
}

func (s *BookService) GetEntry(ctx context.Context, companyID, id uuid.UUID) (*domain.BookEntry, error) {
	return s.bookRepo.GetEntry(ctx, companyID, id)
}

func (s *BookService) GetExport(ctx context.Context, companyID, id uuid.UUID) (*domain.BookExport, error) {
	return s.bookRepo.GetExport(ctx, companyID, id)
}

func (s *BookService) GetBookSummary(ctx context.Context, runID uuid.UUID) (*repository.BookSummary, error) {
	return s.bookRepo.GetBookSummary(ctx, runID)
}
