package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/domain"
)

const (
	TravelStatusDraft     = "DRAFT"
	TravelStatusRequested = "REQUESTED"
	TravelStatusApproved  = "APPROVED"
	TravelStatusRejected  = "REJECTED"
	TravelStatusCompleted = "COMPLETED"
	TravelStatusCancelled = "CANCELLED"
)

type TravelRepository interface {
	Create(ctx context.Context, travel *domain.Travel) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Travel, error)
	List(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.Travel, error)
	Update(ctx context.Context, travel *domain.Travel) error

	AddParticipant(ctx context.Context, participant *domain.TravelParticipant) error
	RemoveParticipant(ctx context.Context, travelID, employeeID uuid.UUID) error
	ListParticipants(ctx context.Context, travelID uuid.UUID) ([]domain.TravelParticipant, error)
}

type TravelService struct {
	travelRepo TravelRepository
	auditRepo  AuditRepository
}

func NewTravelService(travelRepo TravelRepository, auditRepo AuditRepository) *TravelService {
	return &TravelService{travelRepo: travelRepo, auditRepo: auditRepo}
}

func (s *TravelService) CreateTravel(ctx context.Context, companyID, employeeID, userID uuid.UUID, t *domain.Travel) (*domain.Travel, error) {
	const op = "CreateTravel"
	now := time.Now()
	t.ID = uuid.New()
	t.CompanyID = companyID
	t.EmployeeID = employeeID
	t.CreatedBy = userID
	t.CreatedAt = now
	t.UpdatedAt = now
	if t.Status == "" {
		t.Status = TravelStatusDraft
	}
	if err := s.travelRepo.Create(ctx, t); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: companyID, UserID: userID,
		Action: "travel.created", EntityType: "travel", EntityID: t.ID, CreatedAt: now,
	})
	return t, nil
}

func (s *TravelService) GetTravel(ctx context.Context, id uuid.UUID) (*domain.Travel, error) {
	const op = "GetTravel"
	travel, err := s.travelRepo.GetByID(ctx, id)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return travel, nil
}

func (s *TravelService) ListTravels(ctx context.Context, companyID uuid.UUID, employeeID *uuid.UUID, status *string, limit, offset int) ([]domain.Travel, error) {
	const op = "ListTravels"
	travels, err := s.travelRepo.List(ctx, companyID, employeeID, status, limit, offset)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return travels, nil
}

func (s *TravelService) UpdateTravel(ctx context.Context, userID uuid.UUID, t *domain.Travel) (*domain.Travel, error) {
	const op = "UpdateTravel"
	existing, err := s.travelRepo.GetByID(ctx, t.ID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	t.CreatedAt = existing.CreatedAt
	t.CreatedBy = existing.CreatedBy
	t.UpdatedAt = time.Now()
	if err := s.travelRepo.Update(ctx, t); err != nil {
		return nil, svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: t.CompanyID, UserID: userID,
		Action: "travel.updated", EntityType: "travel", EntityID: t.ID, CreatedAt: time.Now(),
	})
	return t, nil
}

func (s *TravelService) RequestTravel(ctx context.Context, id, userID uuid.UUID) error {
	const op = "RequestTravel"
	travel, err := s.travelRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if travel.Status != TravelStatusDraft {
		return svcErr(op, domain.ErrInvalidInput)
	}
	travel.Status = TravelStatusRequested
	travel.UpdatedAt = time.Now()
	if err := s.travelRepo.Update(ctx, travel); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: travel.CompanyID, UserID: userID,
		Action: "travel.requested", EntityType: "travel", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *TravelService) ApproveTravel(ctx context.Context, id, approverID uuid.UUID, comment string) error {
	const op = "ApproveTravel"
	travel, err := s.travelRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if travel.Status != TravelStatusRequested {
		return svcErr(op, domain.ErrInvalidInput)
	}
	travel.Status = TravelStatusApproved
	travel.UpdatedAt = time.Now()
	if err := s.travelRepo.Update(ctx, travel); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: travel.CompanyID, UserID: approverID,
		Action: "travel.approved", EntityType: "travel", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *TravelService) RejectTravel(ctx context.Context, id, approverID uuid.UUID, reason string) error {
	const op = "RejectTravel"
	travel, err := s.travelRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if travel.Status != TravelStatusRequested {
		return svcErr(op, domain.ErrInvalidInput)
	}
	travel.Status = TravelStatusRejected
	travel.RejectionReason = &reason
	travel.UpdatedAt = time.Now()
	if err := s.travelRepo.Update(ctx, travel); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: travel.CompanyID, UserID: approverID,
		Action: "travel.rejected", EntityType: "travel", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *TravelService) CompleteTravel(ctx context.Context, id uuid.UUID) error {
	const op = "CompleteTravel"
	travel, err := s.travelRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if travel.Status != TravelStatusApproved {
		return svcErr(op, domain.ErrInvalidInput)
	}
	travel.Status = TravelStatusCompleted
	travel.UpdatedAt = time.Now()
	if err := s.travelRepo.Update(ctx, travel); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: travel.CompanyID,
		Action: "travel.completed", EntityType: "travel", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *TravelService) CancelTravel(ctx context.Context, id uuid.UUID) error {
	const op = "CancelTravel"
	travel, err := s.travelRepo.GetByID(ctx, id)
	if err != nil {
		return svcErr(op, err)
	}
	if travel.Status == TravelStatusCompleted || travel.Status == TravelStatusCancelled {
		return svcErr(op, domain.ErrInvalidInput)
	}
	travel.Status = TravelStatusCancelled
	travel.UpdatedAt = time.Now()
	if err := s.travelRepo.Update(ctx, travel); err != nil {
		return svcErr(op, err)
	}
	s.auditRepo.Log(ctx, &domain.ExpenseAuditLog{
		ID: uuid.New(), CompanyID: travel.CompanyID,
		Action: "travel.cancelled", EntityType: "travel", EntityID: id, CreatedAt: time.Now(),
	})
	return nil
}

func (s *TravelService) AddParticipant(ctx context.Context, travelID, employeeID uuid.UUID, role string) (*domain.TravelParticipant, error) {
	const op = "AddParticipant"
	p := &domain.TravelParticipant{
		ID:         uuid.New(),
		TravelID:   travelID,
		EmployeeID: employeeID,
		Role:       role,
		CreatedAt:  time.Now(),
	}
	if err := s.travelRepo.AddParticipant(ctx, p); err != nil {
		return nil, svcErr(op, err)
	}
	return p, nil
}

func (s *TravelService) RemoveParticipant(ctx context.Context, travelID, employeeID uuid.UUID) error {
	const op = "RemoveParticipant"
	if err := s.travelRepo.RemoveParticipant(ctx, travelID, employeeID); err != nil {
		return svcErr(op, err)
	}
	return nil
}

func (s *TravelService) ListParticipants(ctx context.Context, travelID uuid.UUID) ([]domain.TravelParticipant, error) {
	const op = "ListParticipants"
	participants, err := s.travelRepo.ListParticipants(ctx, travelID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return participants, nil
}
