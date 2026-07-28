package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreateRequisitionReq struct {
	PositionID      *string                    `json:"position_id"`
	DepartmentID    *string                    `json:"department_id"`
	HiringManagerID *string                    `json:"hiring_manager_id"`
	Title           string                     `json:"title"`
	Description     *string                    `json:"description"`
	Justification   *string                    `json:"justification"`
	Vacancies       int                        `json:"vacancies"`
	EmploymentType  *string                    `json:"employment_type"`
	WorkMode        *string                    `json:"work_mode"`
	Location        *string                    `json:"location"`
	SalaryMin       *float64                   `json:"salary_min"`
	SalaryMax       *float64                   `json:"salary_max"`
	Currency        *string                    `json:"currency"`
	Urgency         domain.RequisitionUrgency  `json:"urgency"`
	Reason          *string                    `json:"reason"`
	Skills          []CreateRequisitionSkillReq `json:"skills,omitempty"`
}

type CreateRequisitionSkillReq struct {
	Skill    string  `json:"skill"`
	Category *string `json:"category"`
	Required bool    `json:"required"`
	MinYears *int    `json:"min_years"`
}

type RequisitionService struct {
	requisitionRepo *repository.RequisitionRepo
	positionRepo    *repository.PositionRepo
}

func NewRequisitionService(requisitionRepo *repository.RequisitionRepo, positionRepo *repository.PositionRepo) *RequisitionService {
	return &RequisitionService{
		requisitionRepo: requisitionRepo,
		positionRepo:    positionRepo,
	}
}

func (s *RequisitionService) Create(ctx context.Context, companyID, userID string, req *CreateRequisitionReq) (*domain.Requisition, error) {
	const op = "CreateRequisition"
	now := time.Now()
	r := &domain.Requisition{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		PositionID:      req.PositionID,
		DepartmentID:    req.DepartmentID,
		RequestedBy:     userID,
		HiringManagerID: req.HiringManagerID,
		Title:           req.Title,
		Description:     req.Description,
		Justification:   req.Justification,
		Vacancies:       req.Vacancies,
		EmploymentType:  req.EmploymentType,
		WorkMode:        req.WorkMode,
		Location:        req.Location,
		SalaryMin:       req.SalaryMin,
		SalaryMax:       req.SalaryMax,
		Currency:        req.Currency,
		Urgency:         req.Urgency,
		Reason:          req.Reason,
		Status:          domain.ReqStatusDraft,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if r.Urgency == "" {
		r.Urgency = domain.UrgencyNormal
	}
	for _, s := range req.Skills {
		r.Skills = append(r.Skills, domain.RequisitionSkill{
			ID:       uuid.New().String(),
			Skill:    s.Skill,
			Category: s.Category,
			Required: s.Required,
			MinYears: s.MinYears,
		})
	}
	result, err := s.requisitionRepo.Create(ctx, companyID, userID, r)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *RequisitionService) GetByID(ctx context.Context, companyID, id string) (*domain.Requisition, error) {
	const op = "GetRequisition"
	r, err := s.requisitionRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return r, nil
}

func (s *RequisitionService) List(ctx context.Context, companyID, status string) ([]domain.Requisition, error) {
	const op = "ListRequisitions"
	return s.requisitionRepo.List(ctx, companyID, status)
}

func (s *RequisitionService) Update(ctx context.Context, companyID, id string, req *domain.Requisition) (*domain.Requisition, error) {
	const op = "UpdateRequisition"
	req.ID = id
	req.CompanyID = companyID
	req.UpdatedAt = time.Now()
	result, err := s.requisitionRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *RequisitionService) Submit(ctx context.Context, companyID, id string) error {
	const op = "SubmitRequisition"
	r, err := s.requisitionRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != domain.ReqStatusDraft {
		return svcErr(op, domain.ErrInvalidStatus)
	}
	return s.requisitionRepo.UpdateStatus(ctx, companyID, id, domain.ReqStatusPendingApproval)
}

func (s *RequisitionService) Approve(ctx context.Context, companyID, id string) error {
	const op = "ApproveRequisition"
	r, err := s.requisitionRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != domain.ReqStatusPendingApproval {
		return svcErr(op, domain.ErrInvalidStatus)
	}
	now := time.Now()
	r.Status = domain.ReqStatusApproved
	r.ApprovedAt = &now
	r.UpdatedAt = now
	_, err = s.requisitionRepo.Update(ctx, companyID, id, r)
	return err
}

func (s *RequisitionService) Open(ctx context.Context, companyID, id string) error {
	const op = "OpenRequisition"
	r, err := s.requisitionRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != domain.ReqStatusApproved {
		return svcErr(op, domain.ErrInvalidStatus)
	}
	now := time.Now()
	r.Status = domain.ReqStatusOpen
	r.OpenedAt = &now
	r.UpdatedAt = now
	_, err = s.requisitionRepo.Update(ctx, companyID, id, r)
	return err
}

func (s *RequisitionService) Close(ctx context.Context, companyID, id string, reason string) error {
	const op = "CloseRequisition"
	r, err := s.requisitionRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status != domain.ReqStatusOpen {
		return svcErr(op, domain.ErrInvalidStatus)
	}
	now := time.Now()
	r.Status = domain.ReqStatusClosed
	r.ClosedReason = &reason
	r.ClosedAt = &now
	r.UpdatedAt = now
	_, err = s.requisitionRepo.Update(ctx, companyID, id, r)
	return err
}

func (s *RequisitionService) Cancel(ctx context.Context, companyID, id string) error {
	const op = "CancelRequisition"
	r, err := s.requisitionRepo.GetByID(ctx, companyID, id)
	if err != nil {
		return svcErr(op, err)
	}
	if r.Status == domain.ReqStatusApproved || r.Status == domain.ReqStatusCancelled {
		return svcErr(op, domain.ErrInvalidStatus)
	}
	return s.requisitionRepo.UpdateStatus(ctx, companyID, id, domain.ReqStatusCancelled)
}
