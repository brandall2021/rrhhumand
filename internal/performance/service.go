package performance

import (
	"context"
	"fmt"
	"time"
)

type Service struct {
	repo   *Repository
	engine *ScoreEngine
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo:   repo,
		engine: NewScoreEngine(repo),
	}
}

// Cycles
func (s *Service) CreateCycle(ctx context.Context, companyID, createdBy string, req *CreateCycleRequest) (*PerformanceCycle, error) {
	return s.repo.CreateCycle(ctx, companyID, createdBy, req)
}

func (s *Service) GetCycle(ctx context.Context, companyID, id string) (*PerformanceCycle, error) {
	return s.repo.GetCycle(ctx, companyID, id)
}

func (s *Service) ListCycles(ctx context.Context, companyID string) ([]PerformanceCycle, error) {
	return s.repo.ListCycles(ctx, companyID)
}

func (s *Service) UpdateCycle(ctx context.Context, companyID, id string, req *UpdateCycleRequest) (*PerformanceCycle, error) {
	cycle, err := s.repo.GetCycle(ctx, companyID, id)
	if err != nil { return nil, err }
	if cycle.Status != "DRAFT" {
		return nil, fmt.Errorf("can only update DRAFT cycles")
	}
	return s.repo.UpdateCycle(ctx, companyID, id, req)
}

func (s *Service) OpenCycle(ctx context.Context, companyID, id string) error {
	cycle, err := s.repo.GetCycle(ctx, companyID, id)
	if err != nil { return err }
	if cycle.Status != "DRAFT" {
		return fmt.Errorf("can only open DRAFT cycles")
	}
	return s.repo.UpdateCycleStatus(ctx, companyID, id, "OPEN")
}

func (s *Service) CloseCycle(ctx context.Context, companyID, id string) error {
	cycle, err := s.repo.GetCycle(ctx, companyID, id)
	if err != nil { return err }
	if cycle.Status == "CLOSED" || cycle.Status == "CANCELLED" {
		return fmt.Errorf("cycle already closed or cancelled")
	}
	return s.repo.UpdateCycleStatus(ctx, companyID, id, "CLOSED")
}

// Templates
func (s *Service) CreateTemplate(ctx context.Context, companyID string, req *CreateTemplateRequest) (*EvaluationTemplate, error) {
	t, err := s.repo.CreateTemplate(ctx, companyID, req)
	if err != nil { return nil, err }

	for _, sec := range req.Sections {
		s2, err := s.repo.CreateTemplateSection(ctx, t.ID, &sec)
		if err != nil { continue }
		for _, item := range sec.Items {
			s.repo.CreateTemplateSectionItem(ctx, s2.ID, &item)
		}
	}
	return t, nil
}

func (s *Service) ListTemplates(ctx context.Context, companyID string) ([]EvaluationTemplate, error) {
	return s.repo.ListTemplates(ctx, companyID)
}

// Scales
func (s *Service) CreateScale(ctx context.Context, companyID string, req *CreateScaleRequest) (*RatingScale, error) {
	scale, err := s.repo.CreateScale(ctx, companyID, req)
	if err != nil { return nil, err }
	for _, lvl := range req.Levels {
		s.repo.CreateScaleLevel(ctx, scale.ID, &lvl)
	}
	return scale, nil
}

func (s *Service) ListScales(ctx context.Context, companyID string) ([]RatingScale, error) {
	return s.repo.ListScales(ctx, companyID)
}

// Competencies
func (s *Service) CreateCompetency(ctx context.Context, companyID string, req *CreateCompetencyRequest) (*Competency, error) {
	return s.repo.CreateCompetency(ctx, companyID, req)
}

func (s *Service) ListCompetencies(ctx context.Context, companyID string) ([]Competency, error) {
	return s.repo.ListCompetencies(ctx, companyID)
}

func (s *Service) UpdateCompetency(ctx context.Context, companyID, id string, req *UpdateCompetencyRequest) (*Competency, error) {
	return s.repo.UpdateCompetency(ctx, companyID, id, req)
}

// Objectives
func (s *Service) CreateObjective(ctx context.Context, companyID, createdBy string, req *CreateObjectiveRequest) (*PerformanceObjective, error) {
	return s.repo.CreateObjective(ctx, companyID, createdBy, req)
}

func (s *Service) ListObjectives(ctx context.Context, companyID string, filters PerformanceFilters) ([]PerformanceObjective, error) {
	return s.repo.ListObjectives(ctx, companyID, filters)
}

func (s *Service) GetObjective(ctx context.Context, companyID, id string) (*PerformanceObjective, error) {
	return s.repo.GetObjective(ctx, companyID, id)
}

func (s *Service) UpdateObjective(ctx context.Context, companyID, id string, req *UpdateObjectiveRequest) (*PerformanceObjective, error) {
	return s.repo.UpdateObjective(ctx, companyID, id, req)
}

func (s *Service) UpdateObjectiveProgress(ctx context.Context, companyID, id string, req *UpdateProgressRequest) (*PerformanceObjective, error) {
	return s.repo.UpdateObjective(ctx, companyID, id, &UpdateObjectiveRequest{CurrentValue: &req.CurrentValue})
}

// KPIs
func (s *Service) CreateKPI(ctx context.Context, companyID, createdBy string, req *CreateKPIRequest) (*PerformanceKPI, error) {
	return s.repo.CreateKPI(ctx, companyID, createdBy, req)
}

func (s *Service) ListKPIs(ctx context.Context, companyID string, filters PerformanceFilters) ([]PerformanceKPI, error) {
	return s.repo.ListKPIs(ctx, companyID, filters)
}

func (s *Service) UpdateKPI(ctx context.Context, companyID, id string, req *UpdateKPIRequest) (*PerformanceKPI, error) {
	return s.repo.UpdateKPI(ctx, companyID, id, req)
}

// Evaluators
func (s *Service) AssignEvaluators(ctx context.Context, companyID string, req *AssignEvaluatorsRequest) ([]PerformanceEvaluator, error) {
	cycle, err := s.repo.GetCycle(ctx, companyID, req.CycleID)
	if err != nil { return nil, err }
	if cycle.Status == "CLOSED" || cycle.Status == "CANCELLED" {
		return nil, fmt.Errorf("cannot assign evaluators to closed/cancelled cycle")
	}
	return s.repo.AssignEvaluators(ctx, companyID, req)
}

func (s *Service) ListEvaluators(ctx context.Context, companyID, cycleID string) ([]PerformanceEvaluator, error) {
	return s.repo.ListEvaluators(ctx, companyID, cycleID)
}

// Evaluations
func (s *Service) CreateEvaluation(ctx context.Context, companyID string, req *CreateEvaluationRequest) (*PerformanceEvaluation, error) {
	eval, err := s.repo.CreateEvaluation(ctx, companyID, req)
	if err != nil { return nil, err }
	for _, ans := range req.Answers {
		s.repo.CreateAnswer(ctx, eval.ID, &ans)
	}
	return eval, nil
}

func (s *Service) GetEvaluation(ctx context.Context, companyID, id string) (*PerformanceEvaluation, error) {
	return s.repo.GetEvaluation(ctx, companyID, id)
}

func (s *Service) ListEvaluations(ctx context.Context, companyID string, filters PerformanceFilters) ([]PerformanceEvaluation, error) {
	return s.repo.ListEvaluations(ctx, companyID, filters)
}

func (s *Service) UpdateEvaluation(ctx context.Context, companyID, id string, req *UpdateEvaluationRequest) error {
	return s.repo.UpdateEvaluationStatus(ctx, companyID, id, "DRAFT")
}

func (s *Service) SubmitEvaluation(ctx context.Context, companyID, id string) error {
	eval, err := s.repo.GetEvaluation(ctx, companyID, id)
	if err != nil { return err }
	if eval.Status != "DRAFT" {
		return fmt.Errorf("can only submit DRAFT evaluations")
	}

	answers, _ := s.repo.ListAnswersByEvaluation(ctx, id)
	total := 0.0
	count := 0
	for _, a := range answers {
		if a.Score != nil {
			total += *a.Score
			count++
		}
	}
	if count > 0 {
		score := total / float64(count)
		s.repo.SetEvaluationScore(ctx, companyID, id, score)
	}

	return s.repo.SubmitEvaluation(ctx, companyID, id)
}

func (s *Service) ReopenEvaluation(ctx context.Context, companyID, id string) error {
	eval, err := s.repo.GetEvaluation(ctx, companyID, id)
	if err != nil { return err }
	if eval.Status != "SUBMITTED" {
		return fmt.Errorf("can only reopen SUBMITTED evaluations")
	}
	return s.repo.UpdateEvaluationStatus(ctx, companyID, id, "DRAFT")
}

func (s *Service) ApproveEvaluation(ctx context.Context, companyID, id string) error {
	eval, err := s.repo.GetEvaluation(ctx, companyID, id)
	if err != nil { return err }
	if eval.Status != "SUBMITTED" {
		return fmt.Errorf("can only approve SUBMITTED evaluations")
	}
	return s.repo.UpdateEvaluationStatus(ctx, companyID, id, "APPROVED")
}

// Answers
func (s *Service) CreateAnswer(ctx context.Context, evaluationID string, req *CreateAnswerRequest) (*EvaluationAnswer, error) {
	return s.repo.CreateAnswer(ctx, evaluationID, req)
}

func (s *Service) ListAnswers(ctx context.Context, evaluationID string) ([]EvaluationAnswer, error) {
	return s.repo.ListAnswersByEvaluation(ctx, evaluationID)
}

// Feedback
func (s *Service) CreateFeedback(ctx context.Context, companyID, fromUserID string, req *CreateFeedbackRequest) (*PerformanceFeedback, error) {
	return s.repo.CreateFeedback(ctx, companyID, fromUserID, req)
}

func (s *Service) ListFeedback(ctx context.Context, companyID, employeeID string) ([]PerformanceFeedback, error) {
	return s.repo.ListFeedback(ctx, companyID, employeeID)
}

// Evidence
func (s *Service) CreateEvidence(ctx context.Context, companyID, evaluationID, createdBy string, req *CreateEvidenceRequest) (*PerformanceEvidence, error) {
	return s.repo.CreateEvidence(ctx, companyID, evaluationID, createdBy, req)
}

func (s *Service) ListEvidence(ctx context.Context, evaluationID string) ([]PerformanceEvidence, error) {
	return s.repo.ListEvidenceByEvaluation(ctx, evaluationID)
}

// Scoring
func (s *Service) CalculateResult(ctx context.Context, companyID, cycleID, employeeID string) (*PerformanceResult, error) {
	score, err := s.engine.Calculate(ctx, companyID, cycleID, employeeID)
	if err != nil { return nil, err }
	return s.repo.UpsertResult(ctx, companyID, cycleID, employeeID, score)
}

func (s *Service) GetResult(ctx context.Context, companyID, cycleID, employeeID string) (*PerformanceResult, error) {
	return s.repo.GetResult(ctx, companyID, cycleID, employeeID)
}

func (s *Service) ListResults(ctx context.Context, companyID, cycleID string) ([]PerformanceResult, error) {
	return s.repo.ListResults(ctx, companyID, cycleID)
}

func (s *Service) GetScoringRules(ctx context.Context, companyID string) (*ScoringRule, error) {
	return s.repo.GetScoringRules(ctx, companyID)
}

func (s *Service) UpdateScoringRules(ctx context.Context, companyID string, req *UpdateScoringRulesRequest) (*ScoringRule, error) {
	return s.repo.UpdateScoringRules(ctx, companyID, req)
}

// Improvement Plans
func (s *Service) CreateImprovementPlan(ctx context.Context, companyID, createdBy string, req *CreateImprovementPlanRequest) (*ImprovementPlan, error) {
	plan, err := s.repo.CreateImprovementPlan(ctx, companyID, createdBy, req)
	if err != nil { return nil, err }
	for _, act := range req.Actions {
		s.repo.CreatePlanAction(ctx, plan.ID, &act)
	}
	return plan, nil
}

func (s *Service) ListImprovementPlans(ctx context.Context, companyID string, filters PerformanceFilters) ([]ImprovementPlan, error) {
	return s.repo.ListImprovementPlans(ctx, companyID, filters)
}

func (s *Service) GetImprovementPlan(ctx context.Context, companyID, id string) (*ImprovementPlan, error) {
	return s.repo.GetImprovementPlan(ctx, companyID, id)
}

func (s *Service) UpdateImprovementPlan(ctx context.Context, companyID, id string, req *UpdateImprovementPlanRequest) (*ImprovementPlan, error) {
	return s.repo.UpdateImprovementPlan(ctx, companyID, id, req)
}

func (s *Service) CompleteImprovementPlan(ctx context.Context, companyID, id string) error {
	_, err := s.repo.UpdateImprovementPlan(ctx, companyID, id, &UpdateImprovementPlanRequest{
		Status: strPtr("SUCCESS"),
		Outcome: strPtr(fmt.Sprintf("Completed on %s", time.Now().Format("2006-01-02"))),
	})
	return err
}

// Development Plans
func (s *Service) CreateDevelopmentPlan(ctx context.Context, companyID, createdBy string, req *CreateDevelopmentPlanRequest) (*DevelopmentPlan, error) {
	plan, err := s.repo.CreateDevelopmentPlan(ctx, companyID, createdBy, req)
	if err != nil { return nil, err }
	for _, act := range req.Actions {
		s.repo.CreateDevAction(ctx, plan.ID, &act)
	}
	return plan, nil
}

func (s *Service) ListDevelopmentPlans(ctx context.Context, companyID, employeeID string) ([]DevelopmentPlan, error) {
	return s.repo.ListDevelopmentPlans(ctx, companyID, employeeID)
}

func (s *Service) GetDevelopmentPlan(ctx context.Context, companyID, id string) (*DevelopmentPlan, error) {
	return s.repo.GetDevelopmentPlan(ctx, companyID, id)
}

func (s *Service) UpdateDevelopmentPlan(ctx context.Context, companyID, id string, req *UpdateDevelopmentPlanRequest) (*DevelopmentPlan, error) {
	return s.repo.UpdateDevelopmentPlan(ctx, companyID, id, req)
}

// Dashboard
func (s *Service) GetDashboard(ctx context.Context, companyID string) (*PerformanceDashboard, error) {
	return s.repo.GetDashboard(ctx, companyID)
}

// Audit
func (s *Service) CreateAuditLog(ctx context.Context, companyID, userID, entityType, entityID, action string, oldVal, newVal []byte, ip string) error {
	return s.repo.CreateAuditLog(ctx, companyID, userID, entityType, entityID, action, oldVal, newVal, ip)
}

func strPtr(s string) *string { return &s }
