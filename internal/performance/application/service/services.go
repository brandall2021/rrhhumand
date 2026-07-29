package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rrhhumand/api/internal/performance/domain"
	"github.com/rrhhumand/api/internal/performance/repository"
)

type CycleService struct {
	cycles       repository.CycleRepository
	templates    repository.TemplateRepository
	participants repository.ParticipantRepository
}

func NewCycleService(
	cycles repository.CycleRepository,
	templates repository.TemplateRepository,
	participants repository.ParticipantRepository,
) *CycleService {
	return &CycleService{
		cycles:       cycles,
		templates:    templates,
		participants: participants,
	}
}

func (s *CycleService) Create(ctx context.Context, c *domain.PerformanceCycle) error {
	if c.Status == "" {
		c.Status = domain.CycleStatusDraft
	}
	if c.CycleType == "" {
		c.CycleType = domain.CycleTypeAnnual
	}
	if c.MinAnonymousResponses == 0 {
		c.MinAnonymousResponses = 3
	}
	if c.ObjectiveWeight == 0 {
		c.ObjectiveWeight = 60
	}
	if c.CompetencyWeight == 0 {
		c.CompetencyWeight = 40
	}
	return s.cycles.Create(ctx, c)
}

func (s *CycleService) GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceCycle, error) {
	return s.cycles.GetByID(ctx, companyID, id)
}

func (s *CycleService) List(ctx context.Context, filter domain.PerformanceCycleFilter) ([]domain.PerformanceCycle, error) {
	return s.cycles.List(ctx, filter)
}

func (s *CycleService) Update(ctx context.Context, c *domain.PerformanceCycle) error {
	existing, err := s.cycles.GetByID(ctx, c.CompanyID, c.ID)
	if err != nil {
		return err
	}
	if existing.Status != domain.CycleStatusDraft {
		return fmt.Errorf("solo se pueden actualizar ciclos en estado DRAFT")
	}
	return s.cycles.Update(ctx, c)
}

func (s *CycleService) UpdateStatus(ctx context.Context, companyID, id string, status domain.CycleStatus) error {
	existing, err := s.cycles.GetByID(ctx, companyID, id)
	if err != nil {
		return err
	}
	if !isValidTransition(existing.Status, status) {
		return fmt.Errorf("transición inválida de %s a %s", existing.Status, status)
	}
	return s.cycles.UpdateStatus(ctx, companyID, id, status)
}

type ObjectiveService struct {
	objectives repository.ObjectiveRepository
	cycles     repository.CycleRepository
}

func NewObjectiveService(
	objectives repository.ObjectiveRepository,
	cycles repository.CycleRepository,
) *ObjectiveService {
	return &ObjectiveService{objectives: objectives, cycles: cycles}
}

func (s *ObjectiveService) Create(ctx context.Context, o *domain.PerformanceObjective) error {
	if o.Status == "" {
		o.Status = domain.ObjectiveStatusNotStarted
	}
	cycle, err := s.cycles.GetByID(ctx, o.CompanyID, o.CycleID)
	if err != nil {
		return fmt.Errorf("ciclo no encontrado")
	}
	if cycle.Status == domain.CycleStatusClosed || cycle.Status == domain.CycleStatusCancelled {
		return fmt.Errorf("no se pueden crear objetivos en un ciclo cerrado o cancelado")
	}
	return s.objectives.Create(ctx, o)
}

func (s *ObjectiveService) UpdateProgress(ctx context.Context, companyID, id string, currentValue float64) (*domain.PerformanceObjective, error) {
	o, err := s.objectives.GetByID(ctx, companyID, id)
	if err != nil {
		return nil, err
	}
	o.CurrentValue = &currentValue
	if o.TargetValue != nil && *o.TargetValue > 0 {
		progress := (currentValue / *o.TargetValue) * 100
		if progress > 100 {
			progress = 100
		}
		o.Progress = progress
		if progress >= 100 {
			o.Status = domain.ObjectiveStatusCompleted
		} else {
			o.Status = domain.ObjectiveStatusInProgress
		}
	}
	return o, s.objectives.Update(ctx, o)
}

type EvaluationService struct {
	evaluations  repository.EvaluationRepository
	participants repository.ParticipantRepository
	cycles       repository.CycleRepository
	reviews      repository.ReviewRepository
}

func NewEvaluationService(
	evaluations repository.EvaluationRepository,
	participants repository.ParticipantRepository,
	cycles repository.CycleRepository,
	reviews repository.ReviewRepository,
) *EvaluationService {
	return &EvaluationService{
		evaluations:  evaluations,
		participants: participants,
		cycles:       cycles,
		reviews:      reviews,
	}
}

func (s *EvaluationService) AssignParticipants(ctx context.Context, participants []domain.PerformanceParticipant) error {
	if len(participants) == 0 {
		return fmt.Errorf("debe asignar al menos un evaluador")
	}
	cycle, err := s.cycles.GetByID(ctx, participants[0].CompanyID, participants[0].CycleID)
	if err != nil {
		return err
	}
	if cycle.Status == domain.CycleStatusClosed || cycle.Status == domain.CycleStatusCancelled {
		return fmt.Errorf("no se pueden asignar evaluadores a un ciclo cerrado")
	}
	for i := range participants {
		participants[i].Status = domain.EvaluationStatusPending
		participants[i].AssignedAt = time.Now()
	}
	return s.participants.BulkCreate(ctx, participants)
}

func (s *EvaluationService) SubmitEvaluation(ctx context.Context, companyID, id string) error {
	eval, err := s.evaluations.GetByID(ctx, companyID, id)
	if err != nil {
		return err
	}
	if eval.Status != domain.EvaluationStatusDraft {
		return fmt.Errorf("solo se pueden enviar evaluaciones en estado DRAFT")
	}
	answers, err := s.evaluations.ListAnswersByEvaluation(ctx, id)
	if err != nil {
		return err
	}
	total := 0.0
	count := 0
	for _, a := range answers {
		if a.NumericValue != nil {
			total += *a.NumericValue
			count++
		}
	}
	score := 0.0
	if count > 0 {
		score = total / float64(count)
	}
	if err := s.participants.UpdateStatus(ctx, id, domain.EvaluationStatusSubmitted); err != nil {
		s.participants.UpdateStatus(ctx, id, domain.EvaluationStatusSubmitted)
	}
	return s.evaluations.Submit(ctx, companyID, id, score)
}

func (s *EvaluationService) CreateReview(ctx context.Context, rev *domain.PerformanceReview) error {
	if rev.Status == "" {
		rev.Status = domain.EvaluationStatusPending
	}
	return s.reviews.Create(ctx, rev)
}

type FeedbackService struct {
	feedback repository.FeedbackRepository
	checkins repository.CheckInRepository
}

func NewFeedbackService(
	feedback repository.FeedbackRepository,
	checkins repository.CheckInRepository,
) *FeedbackService {
	return &FeedbackService{feedback: feedback, checkins: checkins}
}

type CalibrationService struct {
	calibrations repository.CalibrationRepository
	cycles       repository.CycleRepository
}

func NewCalibrationService(
	calibrations repository.CalibrationRepository,
	cycles repository.CycleRepository,
) *CalibrationService {
	return &CalibrationService{calibrations: calibrations, cycles: cycles}
}

type PlanService struct {
	improvementPlans repository.ImprovementPlanRepository
	developmentPlans repository.DevelopmentPlanRepository
}

func NewPlanService(
	improvementPlans repository.ImprovementPlanRepository,
	developmentPlans repository.DevelopmentPlanRepository,
) *PlanService {
	return &PlanService{
		improvementPlans: improvementPlans,
		developmentPlans: developmentPlans,
	}
}

func (s *PlanService) CreateImprovement(ctx context.Context, p *domain.ImprovementPlan) error {
	if p.Status == "" {
		p.Status = domain.PlanStatusDraft
	}
	return s.improvementPlans.Create(ctx, p)
}

func (s *PlanService) CreateDevelopment(ctx context.Context, p *domain.DevelopmentPlan) error {
	if p.Status == "" {
		p.Status = domain.PlanStatusActive
	}
	return s.developmentPlans.Create(ctx, p)
}

func isValidTransition(from, to domain.CycleStatus) bool {
	valid := map[domain.CycleStatus][]domain.CycleStatus{
		domain.CycleStatusDraft:       {domain.CycleStatusOpen},
		domain.CycleStatusOpen:        {domain.CycleStatusInProgress, domain.CycleStatusCancelled},
		domain.CycleStatusInProgress:  {domain.CycleStatusReview, domain.CycleStatusCancelled},
		domain.CycleStatusReview:      {domain.CycleStatusCalibration, domain.CycleStatusClosed, domain.CycleStatusCancelled},
		domain.CycleStatusCalibration: {domain.CycleStatusClosed, domain.CycleStatusCancelled},
	}
	allowed, ok := valid[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}

func IsValidCycleTransition(from, to domain.CycleStatus) bool {
	return isValidTransition(from, to)
}
