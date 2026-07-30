package application

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type HiringService struct {
	hiringProcessRepo *repository.HiringProcessRepo
	offerRepo         *repository.OfferRepo
	applicationRepo   *repository.ApplicationRepo
	candidateRepo     *repository.CandidateRepo
}

func NewHiringService(
	hiringProcessRepo *repository.HiringProcessRepo,
	offerRepo *repository.OfferRepo,
	applicationRepo *repository.ApplicationRepo,
	candidateRepo *repository.CandidateRepo,
) *HiringService {
	return &HiringService{
		hiringProcessRepo: hiringProcessRepo,
		offerRepo:         offerRepo,
		applicationRepo:   applicationRepo,
		candidateRepo:     candidateRepo,
	}
}

func (s *HiringService) Create(ctx context.Context, companyID, offerID, createdBy string) (*domain.HiringProcess, error) {
	const op = "CreateHiringProcess"
	offer, err := s.offerRepo.GetByID(ctx, companyID, offerID)
	if err != nil {
		return nil, svcErr(op, err)
	}
	if offer.Status != domain.OfferStatusAccepted {
		return nil, svcErr(op, domain.ErrInvalidStatus)
	}
	now := time.Now()
	candidateID := ""
	app, err := s.applicationRepo.GetByID(ctx, companyID, offer.ApplicationID)
	if err == nil {
		candidateID = app.CandidateID
	}
	hp := &domain.HiringProcess{
		ID:                      uuid.New().String(),
		CompanyID:               companyID,
		OfferID:                 &offerID,
		ApplicationID:           offer.ApplicationID,
		CandidateID:             candidateID,
		Status:                  domain.HireStatusPending,
		BackgroundCheckStatus:   "PENDING",
		MedicalCheckStatus:      "PENDING",
		DocVerificationStatus:   "PENDING",
		OnboardingStatus:        "NOT_STARTED",
		CreatedBy:               &createdBy,
		CreatedAt:               now,
		UpdatedAt:               now,
	}
	result, err := s.hiringProcessRepo.Create(ctx, companyID, hp)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *HiringService) GetByID(ctx context.Context, companyID, id string) (*domain.HiringProcess, error) {
	const op = "GetHiringProcess"
	return s.hiringProcessRepo.GetByID(ctx, companyID, id)
}

func (s *HiringService) ListByCompany(ctx context.Context, companyID, status string) ([]domain.HiringProcess, error) {
	const op = "ListHiringProcesses"
	return s.hiringProcessRepo.ListByCompany(ctx, companyID, status)
}

func (s *HiringService) UpdateBackgroundCheck(ctx context.Context, companyID, id, status, result string) error {
	const op = "UpdateBackgroundCheck"
	var r *string
	if result != "" {
		r = &result
	}
	return s.hiringProcessRepo.UpdateBackgroundCheck(ctx, companyID, id, status, r)
}

func (s *HiringService) UpdateMedicalCheck(ctx context.Context, companyID, id, status, result string) error {
	const op = "UpdateMedicalCheck"
	var r *string
	if result != "" {
		r = &result
	}
	return s.hiringProcessRepo.UpdateMedicalCheck(ctx, companyID, id, status, r)
}

func (s *HiringService) UpdateDocVerification(ctx context.Context, companyID, id, status string) error {
	const op = "UpdateDocVerification"
	return s.hiringProcessRepo.UpdateDocVerification(ctx, companyID, id, status)
}

func (s *HiringService) AddTask(ctx context.Context, companyID, processID string, task domain.HiringProcessTask) (*domain.HiringProcessTask, error) {
	const op = "AddHiringTask"
	task.ID = uuid.New().String()
	task.ProcessID = processID
	task.Status = "PENDING"
	task.CreatedAt = time.Now()
	result, err := s.hiringProcessRepo.AddTask(ctx, &task)
	if err != nil {
		return nil, svcErr(op, err)
	}
	return result, nil
}

func (s *HiringService) CompleteTask(ctx context.Context, companyID, processID, taskID string) error {
	const op = "CompleteHiringTask"
	task := &domain.HiringProcessTask{Status: "COMPLETED"}
	_, err := s.hiringProcessRepo.UpdateTask(ctx, taskID, task)
	return err
}

func (s *HiringService) ListTasks(ctx context.Context, companyID, processID string) ([]domain.HiringProcessTask, error) {
	const op = "ListHiringTasks"
	return s.hiringProcessRepo.ListTasks(ctx, processID)
}

func (s *HiringService) LinkOnboarding(ctx context.Context, companyID, processID, onboardingID string) error {
	const op = "LinkOnboarding"
	return s.hiringProcessRepo.UpdateOnboarding(ctx, companyID, processID, "IN_PROGRESS", &onboardingID)
}

func (s *HiringService) Complete(ctx context.Context, companyID, id string) error {
	const op = "CompleteHiringProcess"
	return s.hiringProcessRepo.UpdateStatus(ctx, companyID, id, string(domain.HireStatusCompleted))
}

func (s *HiringService) Cancel(ctx context.Context, companyID, id string) error {
	const op = "CancelHiringProcess"
	return s.hiringProcessRepo.UpdateStatus(ctx, companyID, id, string(domain.HireStatusCancelled))
}
