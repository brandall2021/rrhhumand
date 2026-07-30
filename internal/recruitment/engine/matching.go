package engine

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/recruitment/repository"
)

type MatchingEngine struct {
	scoringEngine *ScoringEngine
	candidateRepo *repository.CandidateRepo
	positionRepo  *repository.PositionRepo
	scoringRepo   *repository.ScoringRepo
}

func NewMatchingEngine(
	scoringEngine *ScoringEngine,
	candidateRepo *repository.CandidateRepo,
	positionRepo *repository.PositionRepo,
	scoringRepo *repository.ScoringRepo,
) *MatchingEngine {
	return &MatchingEngine{
		scoringEngine: scoringEngine,
		candidateRepo: candidateRepo,
		positionRepo:  positionRepo,
		scoringRepo:   scoringRepo,
	}
}

func (e *MatchingEngine) Match(ctx context.Context, companyID, candidateID, positionID string) (*domain.MatchingResult, error) {
	candidate, err := e.candidateRepo.GetByID(ctx, companyID, candidateID)
	if err != nil {
		return nil, err
	}

	position, err := e.positionRepo.GetByID(ctx, companyID, positionID)
	if err != nil {
		return nil, err
	}

	criteria, err := e.getDefaultCriteria(ctx, position.CompanyID)
	if err != nil {
		criteria = []domain.ScoringCriterion{}
	}

	scoringResult, err := e.scoringEngine.CalculateScore(ctx, *candidate, *position, criteria)
	if err != nil {
		return nil, err
	}

	details := formatDetails(scoringResult)

	result := &domain.MatchingResult{
		ID:              uuid.New().String(),
		CandidateID:     candidateID,
		PositionID:      positionID,
		OverallScore:    &scoringResult.OverallScore,
		SkillScore:      &scoringResult.SkillScore,
		ExperienceScore: &scoringResult.ExperienceScore,
		EducationScore:  &scoringResult.EducationScore,
		CultureScore:    &scoringResult.CultureScore,
		Details:         &details,
		MatchedAt:       time.Now(),
	}

	return result, nil
}

func (e *MatchingEngine) AutoMatch(ctx context.Context, companyID, positionID string) ([]*domain.MatchingResult, error) {
	position, err := e.positionRepo.GetByID(ctx, companyID, positionID)
	if err != nil {
		return nil, err
	}

	criteria, err := e.getDefaultCriteria(ctx, position.CompanyID)
	if err != nil {
		criteria = []domain.ScoringCriterion{}
	}

	candidates, err := e.candidateRepo.List(ctx, position.CompanyID, "", "")
	if err != nil {
		return nil, err
	}

	var results []*domain.MatchingResult
	for i := range candidates {
		if candidates[i].Status == domain.CandStatusBlacklisted {
			continue
		}

		scoringResult, err := e.scoringEngine.CalculateScore(ctx, candidates[i], *position, criteria)
		if err != nil {
			continue
		}

		if scoringResult.OverallScore < 0.3 {
			continue
		}

		details := formatDetails(scoringResult)
		result := &domain.MatchingResult{
			ID:              uuid.New().String(),
			CandidateID:     candidates[i].ID,
			PositionID:      positionID,
			OverallScore:    &scoringResult.OverallScore,
			SkillScore:      &scoringResult.SkillScore,
			ExperienceScore: &scoringResult.ExperienceScore,
			EducationScore:  &scoringResult.EducationScore,
			CultureScore:    &scoringResult.CultureScore,
			Details:         &details,
			MatchedAt:       time.Now(),
		}
		results = append(results, result)
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].OverallScore == nil || results[j].OverallScore == nil {
			return false
		}
		return *results[i].OverallScore > *results[j].OverallScore
	})

	if len(results) > 20 {
		results = results[:20]
	}

	return results, nil
}

func (e *MatchingEngine) getDefaultCriteria(ctx context.Context, companyID string) ([]domain.ScoringCriterion, error) {
	models, err := e.scoringRepo.ListScoringModels(ctx, companyID)
	if err != nil {
		return nil, err
	}
	for _, m := range models {
		if m.IsDefault && m.Active {
			return e.scoringRepo.ListCriteria(ctx, m.ID)
		}
	}
	return []domain.ScoringCriterion{}, nil
}

func formatDetails(result *ScoringResult) string {
	return ""
}
