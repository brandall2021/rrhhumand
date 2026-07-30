package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

const (
	ExpenseStatusDraft     = "DRAFT"
	ExpenseStatusSubmitted = "SUBMITTED"
	ExpenseStatusApproved  = "APPROVED"
	ExpenseStatusRejected  = "REJECTED"
	ExpenseStatusObserved  = "OBSERVED"
	ExpenseStatusCancelled = "CANCELLED"
)

type ExpenseRepository interface {
	Create(ctx context.Context, expense *domain.Expense) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Expense, error)
	List(ctx context.Context, companyID uuid.UUID, employeeID, travelID *uuid.UUID, status *string, dateFrom, dateTo *time.Time, limit, offset int) ([]domain.Expense, error)
	Update(ctx context.Context, expense *domain.Expense) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ReceiptRepository interface {
	Create(ctx context.Context, receipt *domain.ExpenseReceipt) error
	Get(ctx context.Context, id uuid.UUID) (*domain.ExpenseReceipt, error)
	ListByExpense(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseReceipt, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DuplicateRepository interface {
	Create(ctx context.Context, check *domain.ExpenseDuplicateCheck) error
	FindByExpenseID(ctx context.Context, expenseID uuid.UUID) ([]domain.ExpenseDuplicateCheck, error)
}

type PolicyEvaluator interface {
	Evaluate(ctx context.Context, expenseCtx domain.ExpenseContext) (domain.PolicyResult, error)
}

type ExpenseService struct {
	expenseRepo     ExpenseRepository
	auditRepo       AuditRepository
	receiptRepo     ReceiptRepository
	duplicateRepo   DuplicateRepository
	policyEvaluator PolicyEvaluator
}

func NewExpenseService(
	expenseRepo ExpenseRepository,
	auditRepo AuditRepository,
	receiptRepo ReceiptRepository,
	duplicateRepo DuplicateRepository,
	policyEvaluator PolicyEvaluator,
) *ExpenseService {
	return &ExpenseService{
		expenseRepo:     expenseRepo,
		auditRepo:       auditRepo,
		receiptRepo:     receiptRepo,
		duplicateRepo:   duplicateRepo,
		policyEvaluator: policyEvaluator,
	}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, companyID, employeeID, userID uuid.UUID, e *domain.Expense) (*domain.Expense, error) {
	const op = "CreateExpense"
	now := time.Now()
	e.ID = uuid.New()
	e.CompanyID = companyID
	e.EmployeeID = employeeID
	e.CreatedBy = userID
	e.CreatedAt = now
	e.UpdatedAt = now

	if e.Status == "" {
		e.Status = ExpenseStatusDraft
	}

	if e.Status == ExpenseStatusSubmitted {
		result, err := s.policyEvaluator.Evaluate(ctx, domain.ExpenseContext{
			EmployeeID:    employeeID.String(),
			CompanyID:     companyID.String(),
			Category:      e.CategoryID.String(),
			Amount:        e.BaseAmount,
			Currency:      e.BaseCurrency,
			ExpenseDate:   e.ExpenseDate.Format(time.RFC3339),
			HasReceipt:    false,
		})
		if err != nil {
			return nil, svcErr(op, err)
		}
		e.IsPolicyCompliant = result.Compliant
		if result.Compliant {
			e.PolicyStatus = "COMPLIANT"
		} else {
			e.PolicyStatus = "VIOLATED"
		}
	} else if e.Status != ExpenseStatusDraft {
		return nil, svcErr(op, domain.ErrInvalidInput)
	}

	if err := s.expenseRepo.Create(ctx, e); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID, EmployeeID: &employeeID,
		Action: "expense.created", EntityType: "expense", EntityID: e.ID, CreatedAt: now,
	})
	return e, nil
}

func (s *ExpenseService) GetExpense(ctx context.Context, companyID, id uuid.UUID) (*domain.Expense, error) {
	const op = "GetExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if expense.CompanyID != companyID {
		return nil, svcErr(op, domain.ErrNotFound)
	}
	return expense, nil
}

func (s *ExpenseService) ListExpenses(ctx context.Context, companyID uuid.UUID, employeeID, travelID *uuid.UUID, status *string, dateFrom, dateTo *time.Time, limit, offset int) ([]domain.Expense, error) {
	const op = "ListExpenses"
	expenses, err := s.expenseRepo.List(ctx, companyID, employeeID, travelID, status, dateFrom, dateTo, limit, offset)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return expenses, nil
}

func (s *ExpenseService) UpdateExpense(ctx context.Context, companyID, userID uuid.UUID, e *domain.Expense) (*domain.Expense, error) {
	const op = "UpdateExpense"
	existing, err := s.expenseRepo.GetByID(ctx, e.ID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if existing.CompanyID != companyID {
		return nil, svcErr(op, domain.ErrNotFound)
	}
	if existing.Status != ExpenseStatusDraft {
		return nil, svcErr(op, domain.ErrExpenseNotEditable)
	}
	e.CompanyID = companyID
	e.CreatedAt = existing.CreatedAt
	e.CreatedBy = existing.CreatedBy
	e.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, e); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "expense.updated", EntityType: "expense", EntityID: e.ID, CreatedAt: time.Now(),
	})
	return e, nil
}

func (s *ExpenseService) SubmitExpense(ctx context.Context, id, userID uuid.UUID) error {
	const op = "SubmitExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if expense.Status != ExpenseStatusDraft {
		return svcErr(op, domain.ErrInvalidInput)
	}

	result, err := s.policyEvaluator.Evaluate(ctx, domain.ExpenseContext{
		EmployeeID:    expense.EmployeeID.String(),
		CompanyID:     expense.CompanyID.String(),
		Category:      expense.CategoryID.String(),
		Amount:        expense.BaseAmount,
		Currency:      expense.BaseCurrency,
		ExpenseDate:   expense.ExpenseDate.Format(time.RFC3339),
		HasReceipt:    false,
	})
	if err != nil {
		return svcErr(op, err)
	}
	expense.IsPolicyCompliant = result.Compliant
	if result.Compliant {
		expense.PolicyStatus = "COMPLIANT"
	} else {
		expense.PolicyStatus = "VIOLATED"
	}

	expense.Status = ExpenseStatusSubmitted
	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: expense.CompanyID, UserID: userID, EmployeeID: &expense.EmployeeID,
		Action: "expense.submitted", EntityType: "expense", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ExpenseService) ApproveExpense(ctx context.Context, id, approverID uuid.UUID, comment string) error {
	const op = "ApproveExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if expense.Status != ExpenseStatusSubmitted {
		return svcErr(op, domain.ErrInvalidInput)
	}
	expense.Status = ExpenseStatusApproved
	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: expense.CompanyID, UserID: approverID,
		Action: "expense.approved", EntityType: "expense", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ExpenseService) RejectExpense(ctx context.Context, id, approverID uuid.UUID, reason string) error {
	const op = "RejectExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if expense.Status != ExpenseStatusSubmitted && expense.Status != ExpenseStatusObserved {
		return svcErr(op, domain.ErrInvalidInput)
	}
	expense.Status = ExpenseStatusRejected
	expense.RejectionReason = &reason
	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: expense.CompanyID, UserID: approverID,
		Action: "expense.rejected", EntityType: "expense", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ExpenseService) ObserveExpense(ctx context.Context, id, observerID uuid.UUID, observation string) error {
	const op = "ObserveExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if expense.Status != ExpenseStatusSubmitted {
		return svcErr(op, domain.ErrInvalidInput)
	}
	expense.Status = ExpenseStatusObserved
	expense.Observation = &observation
	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: expense.CompanyID, UserID: observerID,
		Action: "expense.observed", EntityType: "expense", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ExpenseService) CancelExpense(ctx context.Context, id, userID uuid.UUID) error {
	const op = "CancelExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if expense.Status == ExpenseStatusApproved || expense.Status == ExpenseStatusCancelled {
		return svcErr(op, domain.ErrInvalidInput)
	}
	expense.Status = ExpenseStatusCancelled
	expense.UpdatedAt = time.Now()
	if err := s.expenseRepo.Update(ctx, expense); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: expense.CompanyID, UserID: userID,
		Action: "expense.cancelled", EntityType: "expense", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, companyID, id uuid.UUID) error {
	const op = "DeleteExpense"
	expense, err := s.expenseRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if expense.CompanyID != companyID {
		return svcErr(op, domain.ErrNotFound)
	}
	if expense.Status != ExpenseStatusDraft {
		return svcErr(op, domain.ErrInvalidInput)
	}
	if err := s.expenseRepo.Delete(ctx, id); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID,
		Action: "expense.deleted", EntityType: "expense", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *ExpenseService) UploadReceipt(ctx context.Context, companyID, expenseID, userID uuid.UUID, filename, contentType string, content []byte) (*domain.ExpenseReceipt, error) {
	const op = "UploadReceipt"
	now := time.Now()
	receipt := &domain.ExpenseReceipt{
		ID:         uuid.New(),
		CompanyID:  companyID,
		ExpenseID:  expenseID,
		Filename:   filename,
		MimeType:   contentType,
		Size:       int64(len(content)),
		UploadedBy: userID,
		UploadedAt: now,
	}
	if err := s.receiptRepo.Create(ctx, receipt); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "receipt.uploaded", EntityType: "expense_receipt", EntityID: receipt.ID, CreatedAt: now,
	})
	return receipt, nil
}

func (s *ExpenseService) ListReceipts(ctx context.Context, companyID, expenseID uuid.UUID) ([]domain.ExpenseReceipt, error) {
	const op = "ListReceipts"
	receipts, err := s.receiptRepo.ListByExpense(ctx, expenseID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return receipts, nil
}

func (s *ExpenseService) DeleteReceipt(ctx context.Context, companyID, id uuid.UUID) error {
	const op = "DeleteReceipt"
	if err := s.receiptRepo.Delete(ctx, id); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID,
		Action: "receipt.deleted", EntityType: "expense_receipt", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}
