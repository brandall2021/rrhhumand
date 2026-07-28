package engine

import (
	"context"
	"math"

	"github.com/rrhhumand/api/internal/recruitment/domain"
)

type ScoringResult struct {
	SkillScore      float64                `json:"skill_score"`
	ExperienceScore float64                `json:"experience_score"`
	EducationScore  float64                `json:"education_score"`
	CultureScore    float64                `json:"culture_score"`
	OverallScore    float64                `json:"overall_score"`
	CriterionScores []CriterionScoreDetail `json:"criterion_scores,omitempty"`
}

type CriterionScoreDetail struct {
	CriterionName string  `json:"criterion_name"`
	Score         float64 `json:"score"`
	Weight        float64 `json:"weight"`
	WeightedScore float64 `json:"weighted_score"`
	MatchType     string  `json:"match_type"`
}

type ScoringEngine struct{}

func NewScoringEngine() *ScoringEngine {
	return &ScoringEngine{}
}

func (e *ScoringEngine) CalculateScore(ctx context.Context, candidate domain.Candidate, position domain.Position, criteria []domain.ScoringCriterion) (*ScoringResult, error) {
	result := &ScoringResult{}

	for _, c := range criteria {
		var criterionScore float64
		matchType := "unknown"

		switch c.ScoringType {
		case "exact_match":
			criterionScore, matchType = e.evaluateExactMatch(c, candidate, position)
		case "partial_match":
			criterionScore, matchType = e.evaluatePartialMatch(c, candidate, position)
		case "range":
			criterionScore, matchType = e.evaluateRange(c, candidate, position)
		case "boolean":
			criterionScore, matchType = e.evaluateBoolean(c, candidate, position)
		default:
			criterionScore, matchType = e.evaluateCustom(c, candidate, position)
		}

		weighted := criterionScore * c.Weight
		result.CriterionScores = append(result.CriterionScores, CriterionScoreDetail{
			CriterionName: c.Name,
			Score:         criterionScore,
			Weight:        c.Weight,
			WeightedScore: weighted,
			MatchType:     matchType,
		})
	}

	result = e.aggregateScores(result, candidate, position)
	return result, nil
}

func (e *ScoringEngine) evaluateExactMatch(c domain.ScoringCriterion, candidate domain.Candidate, position domain.Position) (float64, string) {
	fieldValue := getFieldValue(c.Field, candidate, position)
	expectedValue := getConfigValue(c.Config, "value")
	if fieldValue == expectedValue {
		return 1.0, "exact"
	}
	return 0.0, "exact"
}

func (e *ScoringEngine) evaluatePartialMatch(c domain.ScoringCriterion, candidate domain.Candidate, position domain.Position) (float64, string) {
	switch c.Field {
	case "skills":
		candidateSkills := make(map[string]bool)
		for _, s := range candidate.Skills {
			candidateSkills[s.Skill] = true
		}
		matched := 0
		for _, s := range position.Skills {
			if candidateSkills[s.Skill] {
				matched++
			}
		}
		if len(position.Skills) == 0 {
			return 0, "partial"
		}
		return float64(matched) / float64(len(position.Skills)), "partial"
	default:
		fieldValue := getFieldValue(c.Field, candidate, position)
		if fieldValue != "" {
			return 0.5, "partial"
		}
		return 0, "partial"
	}
}

func (e *ScoringEngine) evaluateRange(c domain.ScoringCriterion, candidate domain.Candidate, position domain.Position) (float64, string) {
	minVal, maxVal := getRangeConfig(c.Config)
	fieldValue := getNumericFieldValue(c.Field, candidate, position)
	if fieldValue >= minVal && fieldValue <= maxVal {
		if maxVal-minVal == 0 {
			return 1.0, "range"
		}
		position := (fieldValue - minVal) / (maxVal - minVal)
		return math.Min(1.0, math.Max(0.0, position)), "range"
	}
	return 0.0, "range"
}

func (e *ScoringEngine) evaluateBoolean(c domain.ScoringCriterion, candidate domain.Candidate, position domain.Position) (float64, string) {
	fieldValue := getBooleanFieldValue(c.Field, candidate, position)
	if fieldValue {
		return 1.0, "boolean"
	}
	return 0.0, "boolean"
}

func (e *ScoringEngine) evaluateCustom(c domain.ScoringCriterion, candidate domain.Candidate, position domain.Position) (float64, string) {
	return 0.5, "custom"
}

func (e *ScoringEngine) aggregateScores(result *ScoringResult, candidate domain.Candidate, position domain.Position) *ScoringResult {
	var totalWeight float64
	var totalWeightedScore float64

	skillTotal := 0.0
	skillWeight := 0.0
	expTotal := 0.0
	expWeight := 0.0
	eduTotal := 0.0
	eduWeight := 0.0
	cultureTotal := 0.0
	cultureWeight := 0.0

	for _, cs := range result.CriterionScores {
		totalWeight += cs.Weight
		totalWeightedScore += cs.WeightedScore

		switch {
		case contains(cs.CriterionName, "skill", "Skill", "habilidad"):
			skillTotal += cs.WeightedScore
			skillWeight += cs.Weight
		case contains(cs.CriterionName, "experience", "Experience", "exp", "años"):
			expTotal += cs.WeightedScore
			expWeight += cs.Weight
		case contains(cs.CriterionName, "education", "Education", "edu", "titulo"):
			eduTotal += cs.WeightedScore
			eduWeight += cs.Weight
		default:
			cultureTotal += cs.WeightedScore
			cultureWeight += cs.Weight
		}
	}

	if skillWeight > 0 {
		result.SkillScore = skillTotal / skillWeight
	}
	if expWeight > 0 {
		result.ExperienceScore = expTotal / expWeight
	}
	if eduWeight > 0 {
		result.EducationScore = eduTotal / eduWeight
	}
	if cultureWeight > 0 {
		result.CultureScore = cultureTotal / cultureWeight
	}
	if totalWeight > 0 {
		result.OverallScore = totalWeightedScore / totalWeight
	}

	return result
}

func getFieldValue(field string, candidate domain.Candidate, position domain.Position) string {
	switch field {
	case "candidate.first_name":
		return candidate.FirstName
	case "candidate.last_name":
		return candidate.LastName
	case "candidate.email":
		return candidate.Email
	case "candidate.location":
		if candidate.Location != nil {
			return *candidate.Location
		}
	case "position.title":
		return position.Title
	case "position.location":
		if position.LocationID != nil {
			return *position.LocationID
		}
	}
	return ""
}

func getNumericFieldValue(field string, candidate domain.Candidate, position domain.Position) float64 {
	switch field {
	case "candidate.years_experience":
		var total float64
		for _, e := range candidate.Experience {
			if e.StartDate != nil && e.EndDate != nil {
				total += e.EndDate.Sub(*e.StartDate).Hours() / (365.25 * 24)
			}
		}
		return total
	case "candidate.salary_expectation_min":
		if candidate.SalaryExpectMin != nil {
			return *candidate.SalaryExpectMin
		}
	case "position.salary_min":
		if position.SalaryMin != nil {
			return *position.SalaryMin
		}
	}
	return 0
}

func getBooleanFieldValue(field string, candidate domain.Candidate, position domain.Position) bool {
	switch field {
	case "candidate.is_referral":
		return candidate.IsReferral
	case "candidate.has_portfolio":
		return candidate.PortfolioURL != nil
	case "position.has_remote":
		return position.WorkMode != nil && *position.WorkMode == "REMOTE"
	}
	return false
}

func getConfigValue(config *string, key string) string {
	if config == nil {
		return ""
	}
	return *config
}

func getRangeConfig(config *string) (float64, float64) {
	return 0, 100
}

func contains(s string, substrs ...string) bool {
	for _, sub := range substrs {
		if s == sub {
			return true
		}
	}
	return false
}
