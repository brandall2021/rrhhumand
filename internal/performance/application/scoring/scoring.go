package scoring

import (
	"context"
	"math"

	"github.com/rrhhumand/api/internal/performance/domain"
)

type Weights struct {
	Objective  float64
	Competency float64
	Self       float64
	Manager    float64
	Peer       float64
	HR         float64
	PeerCount  int
}

type EvaluationStats struct {
	ObjectiveScore  float64
	CompetencyScore float64
	SelfScore       float64
	ManagerScore    float64
	PeerScore       float64
	HRScore         float64
}

type ObjectiveScorer interface {
	CalculateObjectiveScore(ctx context.Context, companyID, cycleID, employeeID string) (float64, error)
}

type CompetencyScorer interface {
	CalculateCompetencyScore(ctx context.Context, companyID, cycleID, employeeID string) (float64, error)
}

type EvaluationScorer interface {
	CalculateEvaluatorScore(ctx context.Context, companyID, cycleID, employeeID, evalType string) (float64, error)
}

type RatingScaleEngine interface {
	DetermineRating(score float64) (string, string)
}

type Scorer struct {
	objectiveScorer  ObjectiveScorer
	competencyScorer CompetencyScorer
	evaluationScorer EvaluationScorer
	ratingEngine     RatingScaleEngine
	weights          Weights
}

func NewScorer(
	objectiveScorer ObjectiveScorer,
	competencyScorer CompetencyScorer,
	evaluationScorer EvaluationScorer,
	ratingEngine RatingScaleEngine,
	weights Weights,
) *Scorer {
	return &Scorer{
		objectiveScorer:  objectiveScorer,
		competencyScorer: competencyScorer,
		evaluationScorer: evaluationScorer,
		ratingEngine:     ratingEngine,
		weights:          weights,
	}
}

func (s *Scorer) Calculate(ctx context.Context, companyID, cycleID, employeeID string) (*domain.PerformanceResult, error) {
	objScore, _ := s.objectiveScorer.CalculateObjectiveScore(ctx, companyID, cycleID, employeeID)
	compScore, _ := s.competencyScorer.CalculateCompetencyScore(ctx, companyID, cycleID, employeeID)
	selfScore, _ := s.evaluationScorer.CalculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "SELF")
	managerScore, _ := s.evaluationScorer.CalculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "MANAGER")
	peerScore, _ := s.evaluationScorer.CalculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "PEER")
	hrScore, _ := s.evaluationScorer.CalculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "HR")

	objW := s.weights.Objective / 100.0
	compW := s.weights.Competency / 100.0
	compEvalW := (managerScore * s.weights.Manager / 100.0) +
		(peerScore * s.weights.Peer / 100.0) +
		(selfScore * s.weights.Self / 100.0) +
		(hrScore * s.weights.HR / 100.0)

	finalScore := math.Round((objScore*objW+compEvalW*compW)*100) / 100
	rating, ratingLabel := s.ratingEngine.DetermineRating(finalScore)

	return &domain.PerformanceResult{
		ObjectiveScore:   &objScore,
		CompetencyScore:  &compScore,
		SelfScore:        &selfScore,
		ManagerScore:     &managerScore,
		PeerScore:        &peerScore,
		HRScore:          &hrScore,
		FinalScore:       &finalScore,
		FinalRating:      &rating,
		FinalRatingLabel: &ratingLabel,
	}, nil
}

type DefaultRatingScale struct{}

func (DefaultRatingScale) DetermineRating(score float64) (string, string) {
	switch {
	case score >= 90:
		return "EXCEEDS_EXPECTATIONS", "Supera las expectativas"
	case score >= 75:
		return "MEETS_EXPECTATIONS", "Cumple las expectativas"
	case score >= 60:
		return "NEEDS_IMPROVEMENT", "Necesita mejora"
	case score >= 40:
		return "BELOW_EXPECTATIONS", "Por debajo de las expectativas"
	default:
		return "UNSATISFACTORY", "Insatisfactorio"
	}
}
