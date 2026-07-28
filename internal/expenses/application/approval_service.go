package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

type ApprovalRepository interface {
	Create(ctx context.Context, approval *domain.ExpenseApproval) error
	GetPendingByApprover(ctx context.Context, approverID uuid.UUID) ([]domain.ExpenseApproval, error)
	Update(ctx context.Context, approval *domain.ExpenseApproval) error
}

type WorkflowRepository interface {
	GetByEntityType(ctx context.Context, companyID uuid.UUID, entityType string) (*domain.ExpenseWorkflow, error)
	GetSteps(ctx context.Context, workflowID uuid.UUID) ([]domain.ExpenseWorkflowStep, error)
}

type ApprovalService struct {
	approvalRepo  ApprovalRepository
	workflowRepo  WorkflowRepository
	expenseRepo   ExpenseRepository
	travelRepo    TravelRepository
	reportRepo    ReportRepository
	advanceRepo   AdvanceRepository
	auditRepo     AuditRepository
}

func NewApprovalService(
	approvalRepo ApprovalRepository,
	workflowRepo WorkflowRepository,
	expenseRepo ExpenseRepository,
	travelRepo TravelRepository,
	reportRepo ReportRepository,
	advanceRepo AdvanceRepository,
	auditRepo AuditRepository,
) *ApprovalService {
	return &ApprovalService{
		approvalRepo: approvalRepo,
		workflowRepo: workflowRepo,
		expenseRepo:  expenseRepo,
		travelRepo:   travelRepo,
		reportRepo:   reportRepo,
		advanceRepo:  advanceRepo,
		auditRepo:    auditRepo,
	}
}

func (s *ApprovalService) GetPendingApprovals(ctx context.Context, approverID uuid.UUID) ([]domain.ExpenseApproval, error) {
	const op = "GetPendingApprovals"
	approvals, err := s.approvalRepo.GetPendingByApprover(ctx, approverID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return approvals, nil
}

func (s *ApprovalService) ApproveEntity(ctx context.Context, entityType string, entityID, approverID uuid.UUID, comment string) error {
	const op = "ApproveEntity"
	now := time.Now()

	switch entityType {
	case "expense":
		expense, err := s.expenseRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if expense.Status != ExpenseStatusSubmitted {
			return svcErr(op, domain.ErrInvalidInput)
		}
		expense.Status = ExpenseStatusApproved
		expense.UpdatedAt = now
		if err := s.expenseRepo.Update(ctx, expense); err != nil {
			return svcErr(op, err)
		}

	case "travel":
		travel, err := s.travelRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if travel.Status != TravelStatusRequested {
			return svcErr(op, domain.ErrInvalidInput)
		}
		travel.Status = TravelStatusApproved
		travel.UpdatedAt = now
		if err := s.travelRepo.Update(ctx, travel); err != nil {
			return svcErr(op, err)
		}

	case "report":
		report, err := s.reportRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if report.Status != ReportStatusSubmitted {
			return svcErr(op, domain.ErrInvalidInput)
		}
		report.Status = ReportStatusApproved
		report.ApprovedAt = &now
		report.UpdatedAt = now
		if err := s.reportRepo.Update(ctx, report); err != nil {
			return svcErr(op, err)
		}

	case "advance":
		advance, err := s.advanceRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if advance.Status != AdvanceStatusRequested {
			return svcErr(op, domain.ErrInvalidInput)
		}
		advance.Status = AdvanceStatusApproved
		advance.ApprovedDate = &now
		advance.UpdatedAt = now
		if err := s.advanceRepo.Update(ctx, advance); err != nil {
			return svcErr(op, err)
		}

	default:
		return svcErr(op, domain.ErrInvalidInput)
	}

	approval := &domain.ExpenseApproval{
		ID:         uuid.New(),
		EntityType: entityType,
		EntityID:   entityID,
		ApproverID: approverID,
		Status:     "approved",
		Comment:    &comment,
		ApprovedAt: &now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.approvalRepo.Create(ctx, approval); err != nil {
		return svcErr(op, err)
	}

	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), UserID: approverID,
		Action: entityType + ".approved", EntityType: entityType, EntityID: entityID, CreatedAt: now,
	})
	return nil
}

func (s *ApprovalService) RejectEntity(ctx context.Context, entityType string, entityID, approverID uuid.UUID, reason string) error {
	const op = "RejectEntity"
	now := time.Now()

	switch entityType {
	case "expense":
		expense, err := s.expenseRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if expense.Status != ExpenseStatusSubmitted && expense.Status != ExpenseStatusObserved {
			return svcErr(op, domain.ErrInvalidInput)
		}
		expense.Status = ExpenseStatusRejected
		expense.RejectionReason = &reason
		expense.UpdatedAt = now
		if err := s.expenseRepo.Update(ctx, expense); err != nil {
			return svcErr(op, err)
		}

	case "travel":
		travel, err := s.travelRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if travel.Status != TravelStatusRequested {
			return svcErr(op, domain.ErrInvalidInput)
		}
		travel.Status = TravelStatusRejected
		travel.RejectionReason = &reason
		travel.UpdatedAt = now
		if err := s.travelRepo.Update(ctx, travel); err != nil {
			return svcErr(op, err)
		}

	case "report":
		report, err := s.reportRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if report.Status != ReportStatusSubmitted && report.Status != ReportStatusObserved {
			return svcErr(op, domain.ErrInvalidInput)
		}
		report.Status = ReportStatusRejected
		report.RejectionReason = &reason
		report.UpdatedAt = now
		if err := s.reportRepo.Update(ctx, report); err != nil {
			return svcErr(op, err)
		}

	case "advance":
		advance, err := s.advanceRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if advance.Status != AdvanceStatusRequested {
			return svcErr(op, domain.ErrInvalidInput)
		}
		advance.Status = AdvanceStatusRejected
		advance.RejectionReason = &reason
		advance.UpdatedAt = now
		if err := s.advanceRepo.Update(ctx, advance); err != nil {
			return svcErr(op, err)
		}

	default:
		return svcErr(op, domain.ErrInvalidInput)
	}

	approval := &domain.ExpenseApproval{
		ID:          uuid.New(),
		EntityType:  entityType,
		EntityID:    entityID,
		ApproverID:  approverID,
		Status:      "rejected",
		RejectedAt:  &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.approvalRepo.Create(ctx, approval); err != nil {
		return svcErr(op, err)
	}

	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), UserID: approverID,
		Action: entityType + ".rejected", EntityType: entityType, EntityID: entityID, CreatedAt: now,
	})
	return nil
}

func (s *ApprovalService) ObserveEntity(ctx context.Context, entityType string, entityID, approverID uuid.UUID, observation string) error {
	const op = "ObserveEntity"
	now := time.Now()

	switch entityType {
	case "expense":
		expense, err := s.expenseRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if expense.Status != ExpenseStatusSubmitted {
			return svcErr(op, domain.ErrInvalidInput)
		}
		expense.Status = ExpenseStatusObserved
		expense.Observation = &observation
		expense.UpdatedAt = now
		if err := s.expenseRepo.Update(ctx, expense); err != nil {
			return svcErr(op, err)
		}

	case "report":
		report, err := s.reportRepo.GetByID(ctx, entityID)
		if err != nil {
			return svcErr(op, err)
		}
		if report.Status != ReportStatusSubmitted {
			return svcErr(op, domain.ErrInvalidInput)
		}
		report.Status = ReportStatusObserved
		report.Observation = &observation
		report.UpdatedAt = now
		if err := s.reportRepo.Update(ctx, report); err != nil {
			return svcErr(op, err)
		}

	default:
		return svcErr(op, domain.ErrInvalidInput)
	}

	approval := &domain.ExpenseApproval{
		ID:          uuid.New(),
		EntityType:  entityType,
		EntityID:    entityID,
		ApproverID:  approverID,
		Status:      "observed",
		ObservedAt:  &now,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := s.approvalRepo.Create(ctx, approval); err != nil {
		return svcErr(op, err)
	}

	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), UserID: approverID,
		Action: entityType + ".observed", EntityType: entityType, EntityID: entityID, CreatedAt: now,
	})
	return nil
}
