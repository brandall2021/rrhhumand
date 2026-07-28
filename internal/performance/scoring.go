package performance

import (
	"context"
	"fmt"
	"math"
)

type ScoreEngine struct {
	repo *Repository
}

func NewScoreEngine(repo *Repository) *ScoreEngine {
	return &ScoreEngine{repo: repo}
}

func (e *ScoreEngine) Calculate(ctx context.Context, companyID, cycleID, employeeID string) (*PerformanceScore, error) {
	rules, err := e.repo.GetScoringRules(ctx, companyID)
	if err != nil { return nil, err }

	result := &PerformanceScore{}

	objScore, _ := e.calculateObjectiveScore(ctx, cycleID, employeeID)
	compScore, _ := e.calculateCompetencyScore(ctx, companyID, cycleID, employeeID)
	kpiScore, _ := e.calculateKPIScore(ctx, cycleID, employeeID)
	selfScore, _ := e.calculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "SELF")
	managerScore, _ := e.calculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "MANAGER")
	peerScore, _ := e.calculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "PEER")
	hrScore, _ := e.calculateEvaluatorScore(ctx, companyID, cycleID, employeeID, "HR")

	result.ObjectiveScore = objScore
	result.CompetencyScore = compScore
	result.KPIScore = kpiScore
	result.SelfScore = selfScore
	result.ManagerScore = managerScore
	result.PeerScore = peerScore
	result.HRScore = hrScore

	objW := rules.ObjectiveWeight / 100.0
	compW := rules.CompetencyWeight / 100.0
	kpiW := rules.KPIWeight / 100.0

	competencyFinal := (selfScore*rules.SelfEvalWeight/100.0 +
		managerScore*rules.ManagerWeight/100.0 +
		peerScore*rules.PeerWeight/100.0 +
		hrScore*rules.HRWeight/100.0)

	result.FinalScore = math.Round((objScore*objW+compScore*compW+kpiScore*kpiW)*100) / 100
	_ = competencyFinal

	result.Rating, result.RatingLabel = e.determineRating(result.FinalScore)

	return result, nil
}

func (e *ScoreEngine) calculateObjectiveScore(ctx context.Context, cycleID, employeeID string) (float64, error) {
	objectives, err := e.repo.ListObjectives(ctx, "", PerformanceFilters{CycleID: cycleID, EmployeeID: employeeID})
	if err != nil || len(objectives) == 0 { return 0, nil }

	totalWeight := 0.0
	weightedScore := 0.0
	for _, obj := range objectives {
		if obj.TargetValue != nil && obj.CurrentValue != nil && *obj.TargetValue > 0 {
			pct := (*obj.CurrentValue / *obj.TargetValue) * 100
			if pct > 100 { pct = 100 }
			weightedScore += pct * (obj.Weight / 100.0)
			totalWeight += obj.Weight
		}
	}

	if totalWeight == 0 { return 0, nil }
	return math.Round((weightedScore/totalWeight)*100*100) / 100, nil
}

func (e *ScoreEngine) calculateCompetencyScore(ctx context.Context, companyID, cycleID, employeeID string) (float64, error) {
	evaluations, err := e.repo.GetEvaluationsByType(ctx, companyID, cycleID, employeeID, "MANAGER")
	if err != nil || len(evaluations) == 0 { return 0, nil }

	total := 0.0
	count := 0
	for _, eval := range evaluations {
		if eval.OverallScore != nil {
			total += *eval.OverallScore
			count++
		}
	}

	if count == 0 { return 0, nil }
	return math.Round((total/float64(count))*100) / 100, nil
}

func (e *ScoreEngine) calculateKPIScore(ctx context.Context, cycleID, employeeID string) (float64, error) {
	kpis, err := e.repo.ListKPIs(ctx, "", PerformanceFilters{CycleID: cycleID, EmployeeID: employeeID})
	if err != nil || len(kpis) == 0 { return 0, nil }

	totalWeight := 0.0
	weightedScore := 0.0
	for _, kpi := range kpis {
		if kpi.TargetValue != nil && kpi.CurrentValue != nil && *kpi.TargetValue > 0 {
			pct := (*kpi.CurrentValue / *kpi.TargetValue) * 100
			if pct > 100 { pct = 100 }
			weightedScore += pct * (kpi.Weight / 100.0)
			totalWeight += kpi.Weight
		}
	}

	if totalWeight == 0 { return 0, nil }
	return math.Round((weightedScore/totalWeight)*100*100) / 100, nil
}

func (e *ScoreEngine) calculateEvaluatorScore(ctx context.Context, companyID, cycleID, employeeID, evaluatorType string) (float64, error) {
	evaluations, err := e.repo.GetEvaluationsByType(ctx, companyID, cycleID, employeeID, evaluatorType)
	if err != nil || len(evaluations) == 0 { return 0, nil }

	total := 0.0
	count := 0
	for _, eval := range evaluations {
		if eval.OverallScore != nil {
			total += *eval.OverallScore
			count++
		}
	}

	if count == 0 { return 0, nil }
	return math.Round((total/float64(count))*100) / 100, nil
}

func (e *ScoreEngine) determineRating(score float64) (string, string) {
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

func (e *ScoreEngine) GetEmployeePerformanceSummary(ctx context.Context, companyID, cycleID, employeeID string) (map[string]interface{}, error) {
	result, err := e.repo.GetResult(ctx, companyID, cycleID, employeeID)
	if err != nil { return nil, fmt.Errorf("no result found") }

	objectives, _ := e.repo.ListObjectives(ctx, companyID, PerformanceFilters{CycleID: cycleID, EmployeeID: employeeID})
	kpis, _ := e.repo.ListKPIs(ctx, companyID, PerformanceFilters{CycleID: cycleID, EmployeeID: employeeID})
	feedbacks, _ := e.repo.ListFeedback(ctx, companyID, employeeID)
	impPlans, _ := e.repo.ListImprovementPlans(ctx, companyID, PerformanceFilters{EmployeeID: employeeID})

	return map[string]interface{}{
		"result":            result,
		"objectives":        objectives,
		"kpis":              kpis,
		"feedback_count":    len(feedbacks),
		"improvement_plans": len(impPlans),
	}, nil
}
