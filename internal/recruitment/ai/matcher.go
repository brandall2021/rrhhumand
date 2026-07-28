package ai

import (
	"context"
	"math/rand"
)

type SemanticScore struct {
	Overall    float64            `json:"overall"`
	Skills     float64            `json:"skills"`
	Experience float64            `json:"experience"`
	Education  float64            `json:"education"`
	Culture    float64            `json:"culture"`
	Confidence float64            `json:"confidence"`
	Highlights []string           `json:"highlights,omitempty"`
	Details    map[string]float64 `json:"details,omitempty"`
}

type RankedCandidate struct {
	CandidateID  string        `json:"candidate_id"`
	Score        float64       `json:"score"`
	SemanticScore *SemanticScore `json:"semantic_score,omitempty"`
	Rank         int           `json:"rank"`
}

type AIMatcher struct{}

func NewAIMatcher() *AIMatcher {
	return &AIMatcher{}
}

// TODO: implement AI integration
func (m *AIMatcher) SemanticMatch(ctx context.Context, candidateProfile, jobRequirements string) (*SemanticScore, error) {
	score := &SemanticScore{
		Overall:    rand.Float64() * 100,
		Skills:     rand.Float64() * 100,
		Experience: rand.Float64() * 100,
		Education:  rand.Float64() * 100,
		Culture:    rand.Float64() * 100,
		Confidence: rand.Float64() * 0.5,
		Highlights: []string{"Skill match detected: Go, Python", "Experience overlap: 3 years"},
		Details: map[string]float64{
			"technical":   rand.Float64() * 100,
			"leadership":  rand.Float64() * 100,
			"domain_knowledge": rand.Float64() * 100,
		},
	}
	return score, nil
}

// TODO: implement AI integration
func (m *AIMatcher) RankCandidates(ctx context.Context, candidates []string, positionID string) ([]*RankedCandidate, error) {
	var ranked []*RankedCandidate
	for i, c := range candidates {
		ranked = append(ranked, &RankedCandidate{
			CandidateID: c,
			Score:       rand.Float64() * 100,
			Rank:        i + 1,
		})
	}
	return ranked, nil
}

// TODO: implement AI integration
func (m *AIMatcher) BatchSemanticMatch(ctx context.Context, profiles []string, jobRequirements string) ([]*SemanticScore, error) {
	var scores []*SemanticScore
	for range profiles {
		s, err := m.SemanticMatch(ctx, "", jobRequirements)
		if err != nil {
			return nil, err
		}
		scores = append(scores, s)
	}
	return scores, nil
}

// TODO: implement AI integration
func (m *AIMatcher) ExplainMatch(ctx context.Context, candidateID, positionID string) (string, error) {
	return "Match explanation not available in stub mode.", nil
}
