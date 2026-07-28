package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

const (
	ReportStatusDraft     = "DRAFT"
	ReportStatusSubmitted = "SUBMITTED"
	ReportStatusApproved  = "APPROVED"
	ReportStatusRejected  = "REJECTED"
	ReportStatusObserved  = "OBSERVED"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.ExpenseReport) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.ExpenseReport, error)
	List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseReport, error)
	Update(ctx context.Context, report *domain.ExpenseReport) error
}

type SettlementCalculator interface {
	Calculate(advance decimal.Decimal, expenses []domain.Expense) domain.SettlementResult
}

type ReportService struct {
	reportRepo           ReportRepository
	expenseRepo          ExpenseRepository
	advanceRepo          AdvanceRepository
	auditRepo            AuditRepository
	settlementCalculator SettlementCalculator
}

func NewReportService(
	reportRepo ReportRepository,
	expenseRepo ExpenseRepository,
	advanceRepo AdvanceRepository,
	auditRepo AuditRepository,
	settlementCalculator SettlementCalculator,
) *ReportService {
	return &ReportService{
		reportRepo:           reportRepo,
		expenseRepo:          expenseRepo,
		advanceRepo:          advanceRepo,
		auditRepo:            auditRepo,
		settlementCalculator: settlementCalculator,
	}
}

func (s *ReportService) CreateReport(ctx context.Context, companyID, employeeID, userID uuid.UUID, r *domain.ExpenseReport) (*domain.ExpenseReport, error) {
	const op = "CreateReport"
	now := time.Now()
	r.ID = uuid.New()
	r.CompanyID = companyID
	r.EmployeeID = employeeID
	r.CreatedBy = userID
	r.CreatedAt = now
	r.UpdatedAt = now
	r.Status = ReportStatusDraft
	if err := s.reportRepo.Create(ctx, r); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "report.created", EntityType: "expense_report", EntityID: r.ID, CreatedAt: now,
	})
	return r, nil
}

func (s *ReportService) GetReport(ctx context.Context, id uuid.UUID) (*domain.ExpenseReport, error) {
	const op = "GetReport"
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return report, nil
}

func (s *ReportService) ListReports(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.ExpenseReport, error) {
	const op = "ListReports"
	reports, err := s.reportRepo.List(ctx, companyID, employeeID, status, limit, offset)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return reports, nil
}

func (s *ReportService) SubmitReport(ctx context.Context, id, userID uuid.UUID) error {
	const op = "SubmitReport"
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if report.Status != ReportStatusDraft {
		return svcErr(op, domain.ErrInvalidInput)
	}

	var advanceAmount decimal.Decimal
	if report.AdvanceID != nil {
		advance, err := s.advanceRepo.GetByID(ctx, *report.AdvanceID)
		if err != nil {
			return svcErr(op, err)
		}
		if advance.ApprovedAmount != nil {
			advanceAmount = *advance.ApprovedAmount
		} else {
			advanceAmount = advance.RequestedAmount
		}
	}

	expenses, err := s.expenseRepo.List(ctx, report.CompanyID, &report.EmployeeID, nil, nil, nil, nil, 0, 0)
	if err != nil {
		return svcErr(op, err)
	}

	var reportExpenses []domain.Expense
	for _, e := range expenses {
		if e.ExpenseReportID != nil && *e.ExpenseReportID == id {
			reportExpenses = append(reportExpenses, e)
		}
	}

	result := s.settlementCalculator.Calculate(advanceAmount, reportExpenses)
	now := time.Now()
	report.TotalAmount = result.TotalExpenses
	report.AdvanceAmount = result.AdvanceAmount
	report.EmployeeRefundAmount = result.EmployeeOwes
	report.CompanyOwesAmount = result.CompanyOwes
	report.ReimbursableAmount = result.CompanyOwes
	report.Currency = result.Currency
	report.Status = ReportStatusSubmitted
	report.SubmittedAt = &now
	report.UpdatedAt = now

	if err := s.reportRepo.Update(ctx, report); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: report.CompanyID, UserID: userID,
		Action: "report.submitted", EntityType: "expense_report", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *ReportService) ApproveReport(ctx context.Context, id, approverID uuid.UUID, comment string) error {
	const op = "ApproveReport"
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if report.Status != ReportStatusSubmitted {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	report.Status = ReportStatusApproved
	report.ApprovedAt = &now
	report.UpdatedAt = now
	if err := s.reportRepo.Update(ctx, report); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: report.CompanyID, UserID: approverID,
		Action: "report.approved", EntityType: "expense_report", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *ReportService) RejectReport(ctx context.Context, id, approverID uuid.UUID, reason string) error {
	const op = "RejectReport"
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if report.Status != ReportStatusSubmitted && report.Status != ReportStatusObserved {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	report.Status = ReportStatusRejected
	report.RejectionReason = &reason
	report.UpdatedAt = now
	if err := s.reportRepo.Update(ctx, report); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: report.CompanyID, UserID: approverID,
		Action: "report.rejected", EntityType: "expense_report", EntityID: id, CreatedAt: now,
	})
	return nil
}

func (s *ReportService) ObserveReport(ctx context.Context, id, observerID uuid.UUID, observation string) error {
	const op = "ObserveReport"
	report, err := s.reportRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if report.Status != ReportStatusSubmitted {
		return svcErr(op, domain.ErrInvalidInput)
	}
	now := time.Now()
	report.Status = ReportStatusObserved
	report.Observation = &observation
	report.UpdatedAt = now
	if err := s.reportRepo.Update(ctx, report); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: report.CompanyID, UserID: observerID,
		Action: "report.observed", EntityType: "expense_report", EntityID: id, CreatedAt: now,
	})
	return nil
}
