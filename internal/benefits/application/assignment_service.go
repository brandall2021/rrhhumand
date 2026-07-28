package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/rrhhumand/api/internal/benefits/repository"
)

type AssignmentService struct {
	assignmentRepo *repository.AssignmentRepo
	benefitRepo    *repository.BenefitRepo
	walletRepo     *repository.WalletRepo
}

func NewAssignmentService(assignmentRepo *repository.AssignmentRepo, benefitRepo *repository.BenefitRepo, walletRepo *repository.WalletRepo) *AssignmentService {
	return &AssignmentService{
		assignmentRepo: assignmentRepo,
		benefitRepo:    benefitRepo,
		walletRepo:     walletRepo,
	}
}

func (s *AssignmentService) EnrollEmployee(ctx context.Context, companyID, employeeID, benefitID, planID *uuid.UUID, userID uuid.UUID, source string) (*domain.EmployeeBenefit, error) {
	now := time.Now()
	eb := &domain.EmployeeBenefit{
		ID:         uuid.New(),
		CompanyID:  *companyID,
		EmployeeID: *employeeID,
		BenefitID:  *benefitID,
		PlanID:     planID,
		Status:     "ACTIVE",
		EnrollmentDate: now,
		EnrolledAt: now,
		EnrolledBy: userID,
		Source:     source,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.assignmentRepo.Create(ctx, eb); err != nil {
		return nil, svcErr("EnrollEmployee", err)
	}
	return eb, nil
}

func (s *AssignmentService) GetEmployeeBenefit(ctx context.Context, companyID, id uuid.UUID) (*domain.EmployeeBenefit, error) {
	return s.assignmentRepo.Get(ctx, companyID, id)
}

func (s *AssignmentService) ListEmployeeBenefits(ctx context.Context, companyID, employeeID, benefitID *uuid.UUID, status *string) ([]domain.EmployeeBenefit, error) {
	return s.assignmentRepo.List(ctx, companyID, employeeID, benefitID, status)
}

func (s *AssignmentService) UpdateEmployeeBenefit(ctx context.Context, companyID uuid.UUID, eb *domain.EmployeeBenefit) (*domain.EmployeeBenefit, error) {
	eb.CompanyID = companyID
	eb.UpdatedAt = time.Now()
	if err := s.assignmentRepo.Update(ctx, eb); err != nil {
		return nil, svcErr("UpdateEmployeeBenefit", err)
	}
	return eb, nil
}

func (s *AssignmentService) CancelEmployeeBenefit(ctx context.Context, companyID, id uuid.UUID, reason string, cancelledBy uuid.UUID) error {
	eb, err := s.assignmentRepo.Get(ctx, companyID, id)
	if err != nil {
		return svcErr("CancelEmployeeBenefit", err)
	}
	now := time.Now()
	eb.Status = "CANCELLED"
	eb.CancellationDate = &now
	eb.CancellationReason = &reason
	eb.UpdatedAt = now
	if err := s.assignmentRepo.Update(ctx, eb); err != nil {
		return svcErr("CancelEmployeeBenefit", err)
	}
	return nil
}

func (s *AssignmentService) GetHistory(ctx context.Context, employeeBenefitID uuid.UUID) ([]domain.EmployeeBenefitHistory, error) {
	return s.assignmentRepo.ListHistory(ctx, &employeeBenefitID, nil)
}

func (s *AssignmentService) ListHistoryByEmployee(ctx context.Context, employeeID uuid.UUID) ([]domain.EmployeeBenefitHistory, error) {
	return s.assignmentRepo.ListHistory(ctx, nil, &employeeID)
}

func (s *AssignmentService) CreateRequest(ctx context.Context, companyID, employeeID, benefitID uuid.UUID, reqType string, data map[string]any) (*domain.BenefitRequest, error) {
	req := &domain.BenefitRequest{
		ID:          uuid.New(),
		CompanyID:   companyID,
		EmployeeID:  employeeID,
		BenefitID:   benefitID,
		RequestType: reqType,
		Status:      "DRAFT",
		RequestData: data,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if err := s.assignmentRepo.CreateRequest(ctx, req); err != nil {
		return nil, svcErr("CreateRequest", err)
	}
	return req, nil
}

func (s *AssignmentService) GetRequest(ctx context.Context, companyID, id uuid.UUID) (*domain.BenefitRequest, error) {
	return s.assignmentRepo.GetRequest(ctx, companyID, id)
}

func (s *AssignmentService) ListRequests(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string) ([]domain.BenefitRequest, error) {
	return s.assignmentRepo.ListRequests(ctx, &companyID, employeeID, nil, status)
}

func (s *AssignmentService) SubmitRequest(ctx context.Context, id, submittedBy uuid.UUID) error {
	req, err := s.assignmentRepo.GetRequest(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("SubmitRequest", err)
	}
	now := time.Now()
	req.Status = "SUBMITTED"
	req.SubmittedBy = &submittedBy
	req.SubmittedAt = &now
	req.UpdatedAt = now
	if err := s.assignmentRepo.UpdateRequest(ctx, req); err != nil {
		return svcErr("SubmitRequest", err)
	}
	return nil
}

func (s *AssignmentService) ApproveRequest(ctx context.Context, id, reviewerID uuid.UUID, comment string) error {
	req, err := s.assignmentRepo.GetRequest(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("ApproveRequest", err)
	}
	now := time.Now()
	req.Status = "APPROVED"
	req.ResolvedBy = &reviewerID
	req.ResolvedAt = &now
	req.ResolutionNotes = &comment
	req.UpdatedAt = now
	if err := s.assignmentRepo.UpdateRequest(ctx, req); err != nil {
		return svcErr("ApproveRequest", err)
	}
	rev := &domain.BenefitRequestReview{
		ID:         uuid.New(),
		RequestID:  id,
		ReviewerID: reviewerID,
		ReviewType: "APPROVAL",
		Comment:    &comment,
		ReviewedAt: now,
	}
	if err := s.assignmentRepo.CreateReview(ctx, rev); err != nil {
		return svcErr("ApproveRequest", err)
	}
	return nil
}

func (s *AssignmentService) RejectRequest(ctx context.Context, id, reviewerID uuid.UUID, reason string) error {
	req, err := s.assignmentRepo.GetRequest(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("RejectRequest", err)
	}
	now := time.Now()
	req.Status = "REJECTED"
	req.ResolvedBy = &reviewerID
	req.ResolvedAt = &now
	req.ResolutionNotes = &reason
	req.UpdatedAt = now
	if err := s.assignmentRepo.UpdateRequest(ctx, req); err != nil {
		return svcErr("RejectRequest", err)
	}
	rev := &domain.BenefitRequestReview{
		ID:         uuid.New(),
		RequestID:  id,
		ReviewerID: reviewerID,
		ReviewType: "REJECTION",
		Comment:    &reason,
		ReviewedAt: now,
	}
	if err := s.assignmentRepo.CreateReview(ctx, rev); err != nil {
		return svcErr("RejectRequest", err)
	}
	return nil
}

func (s *AssignmentService) CancelRequest(ctx context.Context, id uuid.UUID) error {
	req, err := s.assignmentRepo.GetRequest(ctx, uuid.UUID{}, id)
	if err != nil {
		return svcErr("CancelRequest", err)
	}
	req.Status = "CANCELLED"
	req.UpdatedAt = time.Now()
	if err := s.assignmentRepo.UpdateRequest(ctx, req); err != nil {
		return svcErr("CancelRequest", err)
	}
	return nil
}

func (s *AssignmentService) ListReviews(ctx context.Context, requestID uuid.UUID) ([]domain.BenefitRequestReview, error) {
	return s.assignmentRepo.ListReviews(ctx, requestID)
}
