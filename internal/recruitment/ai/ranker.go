package ai

import (
	"context"
	"sort"

	"github.com/rrhhumand/api/internal/recruitment/engine"
)

type RankedResult struct {
	Result     *engine.ScoringResult `json:"result"`
	FinalScore float64              `json:"final_score"`
	Rank       int                  `json:"rank"`
}

type Ranker struct{}

func NewRanker() *Ranker {
	return &Ranker{}
}

func (r *Ranker) Rank(ctx context.Context, results []*engine.ScoringResult) ([]*RankedResult, error) {
	var ranked []*RankedResult
	for _, res := range results {
		ranked = append(ranked, &RankedResult{
			Result:     res,
			FinalScore: res.OverallScore,
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].FinalScore > ranked[j].FinalScore
	})

	for i := range ranked {
		ranked[i].Rank = i + 1
	}

	return ranked, nil
}

func (r *Ranker) ApplyWeights(results []*engine.ScoringResult, weights map[string]float64) ([]*RankedResult, error) {
	var ranked []*RankedResult

	for _, res := range results {
		finalScore := 0.0
		totalWeight := 0.0

		if w, ok := weights["skill"]; ok {
			finalScore += res.SkillScore * w
			totalWeight += w
		}
		if w, ok := weights["experience"]; ok {
			finalScore += res.ExperienceScore * w
			totalWeight += w
		}
		if w, ok := weights["education"]; ok {
			finalScore += res.EducationScore * w
			totalWeight += w
		}
		if w, ok := weights["culture"]; ok {
			finalScore += res.CultureScore * w
			totalWeight += w
		}

		if totalWeight > 0 {
			finalScore = finalScore / totalWeight
		}

		ranked = append(ranked, &RankedResult{
			Result:     res,
			FinalScore: finalScore,
		})
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].FinalScore > ranked[j].FinalScore
	})

	for i := range ranked {
		ranked[i].Rank = i + 1
	}

	return ranked, nil
}
