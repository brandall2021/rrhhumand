package scoring

import (
	"context"
	"math"

	"github.com/rrhhumand/api/internal/performance/domain"
	"github.com/rrhhumand/api/internal/performance/repository"
)

type objectiveScorer struct {
	repo repository.ObjectiveRepository
}

func NewObjectiveScorer(repo repository.ObjectiveRepository) ObjectiveScorer {
	return &objectiveScorer{repo: repo}
}

func (s *objectiveScorer) CalculateObjectiveScore(ctx context.Context, companyID, cycleID, employeeID string) (float64, error) {
	objectives, err := s.repo.List(ctx, domain.ObjectiveFilter{
		CompanyID:  companyID,
		CycleID:    cycleID,
		EmployeeID: employeeID,
	})
	if err != nil {
		return 0, err
	}
	if len(objectives) == 0 {
		return 0, nil
	}
	totalWeight := 0.0
	weightedScore := 0.0
	for _, o := range objectives {
		w := 1.0
		if o.Weight != nil && *o.Weight > 0 {
			w = *o.Weight
		}
		score := 0.0
		if o.CurrentValue != nil && o.TargetValue != nil && *o.TargetValue > 0 {
			score = math.Min((*o.CurrentValue / *o.TargetValue)*100, 100)
		} else if o.Progress > 0 {
			score = o.Progress
		}
		weightedScore += score * w
		totalWeight += w
	}
	if totalWeight == 0 {
		return 0, nil
	}
	return math.Round((weightedScore/totalWeight)*100) / 100, nil
}

type competencyScorer struct {
	repo repository.EvaluationRepository
}

func NewCompetencyScorer(repo repository.EvaluationRepository) CompetencyScorer {
	return &competencyScorer{repo: repo}
}

func (s *competencyScorer) CalculateCompetencyScore(ctx context.Context, companyID, cycleID, employeeID string) (float64, error) {
	evaluations, err := s.repo.List(ctx, domain.EvaluationFilter{
		CompanyID:  companyID,
		CycleID:    cycleID,
		EmployeeID: employeeID,
	})
	if err != nil {
		return 0, err
	}
	if len(evaluations) == 0 {
		return 0, nil
	}
	total := 0.0
	count := 0
	for _, e := range evaluations {
		ceList, err := s.repo.ListCompetencyEvaluations(ctx, e.ID)
		if err != nil {
			continue
		}
		for _, ce := range ceList {
			if ce.Score != nil {
				total += *ce.Score
				count++
			}
		}
	}
	if count == 0 {
		return 0, nil
	}
	return math.Round((total/float64(count))*100) / 100, nil
}

type evaluationScorer struct {
	repo repository.EvaluationRepository
}

func NewEvaluationScorer(repo repository.EvaluationRepository) EvaluationScorer {
	return &evaluationScorer{repo: repo}
}

func (s *evaluationScorer) CalculateEvaluatorScore(ctx context.Context, companyID, cycleID, employeeID, evalType string) (float64, error) {
	evaluations, err := s.repo.List(ctx, domain.EvaluationFilter{
		CompanyID:    companyID,
		CycleID:      cycleID,
		EmployeeID:   employeeID,
		EvaluationType: domain.EvaluationType(evalType),
	})
	if err != nil {
		return 0, err
	}
	if len(evaluations) == 0 {
		return 0, nil
	}
	total := 0.0
	count := 0
	for _, e := range evaluations {
		if e.FinalScore != nil {
			total += *e.FinalScore
			count++
		}
	}
	if count == 0 {
		return 0, nil
	}
	return math.Round((total/float64(count))*100) / 100, nil
}
