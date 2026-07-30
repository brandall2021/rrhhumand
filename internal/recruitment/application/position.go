package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type CreatePositionReq struct {
	RequisitionID    *string  `json:"requisition_id"`
	Title            string   `json:"title"`
	DepartmentID     *string  `json:"department_id"`
	LocationID       *string  `json:"location_id"`
	EmploymentType   *string  `json:"employment_type"`
	WorkMode         *string  `json:"work_mode"`
	Description      *string  `json:"description"`
	Requirements     *string  `json:"requirements"`
	Responsibilities *string  `json:"responsibilities"`
	Benefits         *string  `json:"benefits"`
	SalaryMin        *float64 `json:"salary_min"`
	SalaryMax        *float64 `json:"salary_max"`
	Currency         *string  `json:"currency"`
	Vacancies        int      `json:"vacancies"`
}

type PositionService struct {
	positionRepo    *repository.PositionRepo
	requisitionRepo *repository.RequisitionRepo
}

func NewPositionService(positionRepo *repository.PositionRepo, requisitionRepo *repository.RequisitionRepo) *PositionService {
	return &PositionService{
		positionRepo:    positionRepo,
		requisitionRepo: requisitionRepo,
	}
}

func (s *PositionService) Create(ctx context.Context, companyID string, req *CreatePositionReq) (*domain.Position, error) {
	const op = "CreatePosition"
	now := time.Now()
	p := &domain.Position{
		ID:               uuid.New().String(),
		CompanyID:        companyID,
		RequisitionID:    req.RequisitionID,
		Title:            req.Title,
		DepartmentID:     req.DepartmentID,
		LocationID:       req.LocationID,
		EmploymentType:   req.EmploymentType,
		WorkMode:         req.WorkMode,
		Description:      req.Description,
		Requirements:     req.Requirements,
		Responsibilities: req.Responsibilities,
		Benefits:         req.Benefits,
		SalaryMin:        req.SalaryMin,
		SalaryMax:        req.SalaryMax,
		Currency:         req.Currency,
		Vacancies:        req.Vacancies,
		Status:           domain.PosStatusDraft,
		CreatedAt:        now,
		UpdatedAt:        now,
	}
	result, err := s.positionRepo.Create(ctx, companyID, p)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *PositionService) GetByID(ctx context.Context, companyID, id string) (*domain.Position, error) {
	const op = "GetPosition"
	return s.positionRepo.GetByID(ctx, companyID, id)
}

func (s *PositionService) List(ctx context.Context, companyID, status string) ([]domain.Position, error) {
	const op = "ListPositions"
	return s.positionRepo.List(ctx, companyID, status)
}

func (s *PositionService) Update(ctx context.Context, companyID, id string, req *domain.Position) (*domain.Position, error) {
	const op = "UpdatePosition"
	req.UpdatedAt = time.Now()
	result, err := s.positionRepo.Update(ctx, companyID, id, req)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *PositionService) Close(ctx context.Context, companyID, id string) error {
	const op = "ClosePosition"
	return s.positionRepo.UpdateStatus(ctx, companyID, id, string(domain.PosStatusFilled))
}

func (s *PositionService) AddSkill(ctx context.Context, companyID, positionID string, skill domain.PositionSkill) (*domain.PositionSkill, error) {
	const op = "AddPositionSkill"
	skill.ID = uuid.New().String()
	skill.PositionID = positionID
	result, err := s.positionRepo.AddSkill(ctx, &skill)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *PositionService) RemoveSkill(ctx context.Context, companyID, positionID, skillID string) error {
	const op = "RemovePositionSkill"
	return s.positionRepo.RemoveSkill(ctx, skillID)
}

func (s *PositionService) GetSkills(ctx context.Context, companyID, positionID string) ([]domain.PositionSkill, error) {
	const op = "GetPositionSkills"
	return s.positionRepo.ListSkills(ctx, positionID)
}
