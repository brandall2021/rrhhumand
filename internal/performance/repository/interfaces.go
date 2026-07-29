package repository

import (
	"context"

	"github.com/rrhhumand/api/internal/performance/domain"
)

type CycleRepository interface {
	Create(ctx context.Context, cycle *domain.PerformanceCycle) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceCycle, error)
	List(ctx context.Context, filter domain.PerformanceCycleFilter) ([]domain.PerformanceCycle, error)
	Update(ctx context.Context, cycle *domain.PerformanceCycle) error
	UpdateStatus(ctx context.Context, companyID, id string, status domain.CycleStatus) error
	Delete(ctx context.Context, companyID, id string) error
}

type TemplateRepository interface {
	Create(ctx context.Context, t *domain.PerformanceTemplate) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceTemplate, error)
	List(ctx context.Context, companyID string) ([]domain.PerformanceTemplate, error)
	Update(ctx context.Context, t *domain.PerformanceTemplate) error
	Delete(ctx context.Context, companyID, id string) error

	CreateSection(ctx context.Context, s *domain.TemplateSection) error
	ListSectionsByTemplate(ctx context.Context, templateID string) ([]domain.TemplateSection, error)
	DeleteSectionsByTemplate(ctx context.Context, templateID string) error

	CreateQuestion(ctx context.Context, q *domain.TemplateQuestion) error
	ListQuestionsByTemplate(ctx context.Context, templateID string) ([]domain.TemplateQuestion, error)
	ListQuestionsBySection(ctx context.Context, sectionID string) ([]domain.TemplateQuestion, error)
	DeleteQuestionsByTemplate(ctx context.Context, templateID string) error
}

type ScaleRepository interface {
	Create(ctx context.Context, s *domain.RatingScale) error
	GetByID(ctx context.Context, companyID, id string) (*domain.RatingScale, error)
	List(ctx context.Context, companyID string) ([]domain.RatingScale, error)
	Update(ctx context.Context, s *domain.RatingScale) error
	Delete(ctx context.Context, companyID, id string) error

	CreateLevel(ctx context.Context, l *domain.RatingScaleLevel) error
	ListLevelsByScale(ctx context.Context, scaleID string) ([]domain.RatingScaleLevel, error)
	DeleteLevelsByScale(ctx context.Context, scaleID string) error
}

type CompetencyRepository interface {
	Create(ctx context.Context, c *domain.Competency) error
	GetByID(ctx context.Context, companyID, id string) (*domain.Competency, error)
	List(ctx context.Context, filter domain.CompetencyFilter) ([]domain.Competency, error)
	Update(ctx context.Context, c *domain.Competency) error
	Delete(ctx context.Context, companyID, id string) error

	CreateLevel(ctx context.Context, l *domain.CompetencyLevel) error
	ListLevelsByCompetency(ctx context.Context, competencyID string) ([]domain.CompetencyLevel, error)
	DeleteLevelsByCompetency(ctx context.Context, competencyID string) error

	UpsertPositionCompetency(ctx context.Context, pc *domain.PositionCompetency) error
	ListByPosition(ctx context.Context, companyID, positionID string) ([]domain.PositionCompetency, error)
	DeletePositionCompetency(ctx context.Context, companyID, positionID, competencyID string) error

	UpsertCycleCompetency(ctx context.Context, cc *domain.CycleCompetency) error
	ListByCycleEmployee(ctx context.Context, cycleID, employeeID string) ([]domain.CycleCompetency, error)
	DeleteCycleCompetency(ctx context.Context, cycleID, employeeID, competencyID string) error
}

type ObjectiveRepository interface {
	Create(ctx context.Context, o *domain.PerformanceObjective) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceObjective, error)
	List(ctx context.Context, filter domain.ObjectiveFilter) ([]domain.PerformanceObjective, error)
	Update(ctx context.Context, o *domain.PerformanceObjective) error
	Delete(ctx context.Context, companyID, id string) error

	CreateKeyResult(ctx context.Context, kr *domain.ObjectiveKeyResult) error
	ListKeyResultsByObjective(ctx context.Context, objectiveID string) ([]domain.ObjectiveKeyResult, error)
	UpdateKeyResult(ctx context.Context, kr *domain.ObjectiveKeyResult) error
	DeleteKeyResult(ctx context.Context, id string) error

	GetWeightTotal(ctx context.Context, cycleID, employeeID string) (float64, error)
}

type ParticipantRepository interface {
	Create(ctx context.Context, p *domain.PerformanceParticipant) error
	BulkCreate(ctx context.Context, participants []domain.PerformanceParticipant) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceParticipant, error)
	ListByCycle(ctx context.Context, companyID, cycleID string) ([]domain.PerformanceParticipant, error)
	ListByEmployee(ctx context.Context, companyID, cycleID, employeeID string) ([]domain.PerformanceParticipant, error)
	ListByEvaluator(ctx context.Context, companyID, cycleID, evaluatorID string) ([]domain.PerformanceParticipant, error)
	UpdateStatus(ctx context.Context, id string, status domain.EvaluationStatus) error
	Delete(ctx context.Context, companyID, id string) error
}

type EvaluationRepository interface {
	Create(ctx context.Context, e *domain.PerformanceEvaluation) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceEvaluation, error)
	List(ctx context.Context, filter domain.EvaluationFilter) ([]domain.PerformanceEvaluation, error)
	Update(ctx context.Context, e *domain.PerformanceEvaluation) error
	UpdateStatus(ctx context.Context, companyID, id string, status domain.EvaluationStatus) error
	Submit(ctx context.Context, companyID, id string, score float64) error
	Delete(ctx context.Context, companyID, id string) error

	CreateAnswer(ctx context.Context, a *domain.EvaluationAnswer) error
	BulkCreateAnswers(ctx context.Context, answers []domain.EvaluationAnswer) error
	ListAnswersByEvaluation(ctx context.Context, evaluationID string) ([]domain.EvaluationAnswer, error)
	UpdateAnswer(ctx context.Context, a *domain.EvaluationAnswer) error
	DeleteAnswersByEvaluation(ctx context.Context, evaluationID string) error

	CreateObjectiveEvaluation(ctx context.Context, oe *domain.ObjectiveEvaluation) error
	ListObjectiveEvaluations(ctx context.Context, evaluationID string) ([]domain.ObjectiveEvaluation, error)
	CreateCompetencyEvaluation(ctx context.Context, ce *domain.CompetencyEvaluation) error
	ListCompetencyEvaluations(ctx context.Context, evaluationID string) ([]domain.CompetencyEvaluation, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, r *domain.PerformanceReview) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceReview, error)
	GetByCycleEmployee(ctx context.Context, companyID, cycleID, employeeID string) (*domain.PerformanceReview, error)
	ListByCycle(ctx context.Context, companyID, cycleID string) ([]domain.PerformanceReview, error)
	Update(ctx context.Context, r *domain.PerformanceReview) error
	UpdateStatus(ctx context.Context, companyID, id string, status domain.EvaluationStatus) error
}

type FeedbackRepository interface {
	Create(ctx context.Context, f *domain.PerformanceFeedback) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceFeedback, error)
	List(ctx context.Context, filter domain.FeedbackFilter) ([]domain.PerformanceFeedback, error)
	Update(ctx context.Context, f *domain.PerformanceFeedback) error
	Delete(ctx context.Context, companyID, id string) error

	CreateRecognition(ctx context.Context, r *domain.PerformanceRecognition) error
	ListRecognitionsByEmployee(ctx context.Context, companyID, employeeID string) ([]domain.PerformanceRecognition, error)
}

type CheckInRepository interface {
	Create(ctx context.Context, ci *domain.PerformanceCheckIn) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceCheckIn, error)
	ListByEmployee(ctx context.Context, companyID, employeeID string) ([]domain.PerformanceCheckIn, error)
	ListByManager(ctx context.Context, companyID, managerID string) ([]domain.PerformanceCheckIn, error)
	Update(ctx context.Context, ci *domain.PerformanceCheckIn) error
	Complete(ctx context.Context, companyID, id string, notes map[string]*string) error
}

type CalibrationRepository interface {
	CreateSession(ctx context.Context, s *domain.CalibrationSession) error
	GetSessionByID(ctx context.Context, companyID, id string) (*domain.CalibrationSession, error)
	ListSessionsByCycle(ctx context.Context, companyID, cycleID string) ([]domain.CalibrationSession, error)
	UpdateSession(ctx context.Context, s *domain.CalibrationSession) error
	UpdateSessionStatus(ctx context.Context, companyID, id string, status domain.CalibrationStatus) error

	CreateItem(ctx context.Context, item *domain.CalibrationItem) error
	BulkCreateItems(ctx context.Context, items []domain.CalibrationItem) error
	ListItemsBySession(ctx context.Context, sessionID string) ([]domain.CalibrationItem, error)
	UpdateItem(ctx context.Context, item *domain.CalibrationItem) error
	ApproveItem(ctx context.Context, id, approvedBy string) error
}

type ImprovementPlanRepository interface {
	Create(ctx context.Context, p *domain.ImprovementPlan) error
	GetByID(ctx context.Context, companyID, id string) (*domain.ImprovementPlan, error)
	List(ctx context.Context, filter domain.PlanFilter) ([]domain.ImprovementPlan, error)
	Update(ctx context.Context, p *domain.ImprovementPlan) error
	UpdateStatus(ctx context.Context, companyID, id string, status domain.PlanStatus) error

	CreateAction(ctx context.Context, a *domain.ImprovementPlanAction) error
	ListActionsByPlan(ctx context.Context, planID string) ([]domain.ImprovementPlanAction, error)
	UpdateAction(ctx context.Context, a *domain.ImprovementPlanAction) error
}

type DevelopmentPlanRepository interface {
	Create(ctx context.Context, p *domain.DevelopmentPlan) error
	GetByID(ctx context.Context, companyID, id string) (*domain.DevelopmentPlan, error)
	List(ctx context.Context, filter domain.PlanFilter) ([]domain.DevelopmentPlan, error)
	Update(ctx context.Context, p *domain.DevelopmentPlan) error
	UpdateStatus(ctx context.Context, companyID, id string, status domain.PlanStatus) error

	CreateAction(ctx context.Context, a *domain.DevelopmentPlanAction) error
	ListActionsByPlan(ctx context.Context, planID string) ([]domain.DevelopmentPlanAction, error)
	UpdateAction(ctx context.Context, a *domain.DevelopmentPlanAction) error
}

type EvidenceRepository interface {
	Create(ctx context.Context, e *domain.PerformanceEvidence) error
	GetByID(ctx context.Context, companyID, id string) (*domain.PerformanceEvidence, error)
	ListByEvaluation(ctx context.Context, evaluationID string) ([]domain.PerformanceEvidence, error)
	ListByObjective(ctx context.Context, objectiveID string) ([]domain.PerformanceEvidence, error)
	ListByFeedback(ctx context.Context, feedbackID string) ([]domain.PerformanceEvidence, error)
	Delete(ctx context.Context, companyID, id string) error
}

type ResultRepository interface {
	Upsert(ctx context.Context, r *domain.PerformanceResult) error
	GetByCycleEmployee(ctx context.Context, companyID, cycleID, employeeID string) (*domain.PerformanceResult, error)
	ListByCycle(ctx context.Context, companyID, cycleID string) ([]domain.PerformanceResult, error)
	ListByEmployee(ctx context.Context, companyID, employeeID string) ([]domain.PerformanceResult, error)
	Delete(ctx context.Context, companyID, cycleID, employeeID string) error
}

type DashboardRepository interface {
	GetDashboard(ctx context.Context, companyID string) (*domain.PerformanceDashboard, error)
}

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.PerformanceAuditLog) error
	ListByEntity(ctx context.Context, companyID, entityType, entityID string) ([]domain.PerformanceAuditLog, error)
}

type OutboxRepository interface {
	Create(ctx context.Context, event *domain.OutboxEvent) error
	ListPending(ctx context.Context, limit int) ([]domain.OutboxEvent, error)
	MarkProcessed(ctx context.Context, id string) error
	MarkFailed(ctx context.Context, id string, err error) error
}
