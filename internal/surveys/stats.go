package surveys

import (
	"context"
	"fmt"

	"github.com/rrhhumand/api/internal/models"
)

type SurveyStatsService struct{}

func NewSurveyStatsService() *SurveyStatsService {
	return &SurveyStatsService{}
}

func (ss *SurveyStatsService) CalculateStats(ctx context.Context, repo *SurveyRepository, survey *models.Survey) (*SurveyStats, error) {
	totalResponded, err := repo.GetResponseCount(ctx, survey.ID)
	if err != nil {
		return nil, err
	}

	totalTargeted, err := repo.GetTargetedEmployeeCount(ctx, survey.ID)
	if err != nil {
		return nil, err
	}

	var participationRate float64
	if totalTargeted > 0 {
		participationRate = float64(totalResponded) / float64(totalTargeted) * 100
	}

	stats := &SurveyStats{
		TotalTargeted:    totalTargeted,
		TotalResponded:   totalResponded,
		ParticipationRate: participationRate,
	}

	questions, err := repo.ListQuestionsBySurveyID(ctx, survey.ID)
	if err != nil {
		return nil, err
	}

	for _, q := range questions {
		qStats, err := ss.calculateQuestionStats(ctx, repo, &q)
		if err != nil {
			continue
		}
		stats.Questions = append(stats.Questions, *qStats)
	}

	return stats, nil
}

func (ss *SurveyStatsService) calculateQuestionStats(ctx context.Context, repo *SurveyRepository, q *models.SurveyQuestion) (*QuestionStats, error) {
	qStats := &QuestionStats{
		QuestionID: q.ID,
		Question:   q.Question,
		Type:       q.Type,
	}

	totalAnswers, err := repo.GetTotalAnswerCount(ctx, q.ID)
	if err != nil {
		return nil, err
	}
	qStats.TotalAnswers = totalAnswers

	switch q.Type {
	case "RATING":
		avg, min, max, count, err := repo.GetRatingStats(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		qStats.Average = &avg
		qStats.Min = &min
		qStats.Max = &max
		_ = count

		dist, err := repo.GetOptionDistribution(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		for i := range dist {
			if totalAnswers > 0 {
				dist[i].Percentage = float64(dist[i].Count) / float64(totalAnswers) * 100
			}
		}
		qStats.Distribution = dist

	case "SINGLE_CHOICE":
		dist, err := repo.GetOptionDistribution(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		for i := range dist {
			if totalAnswers > 0 {
				dist[i].Percentage = float64(dist[i].Count) / float64(totalAnswers) * 100
			}
		}
		qStats.Distribution = dist

	case "MULTIPLE_CHOICE":
		dist, err := repo.GetMultipleChoiceDistribution(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		for i := range dist {
			if totalAnswers > 0 {
				dist[i].Percentage = float64(dist[i].Count) / float64(totalAnswers) * 100
			}
		}
		qStats.Distribution = dist

	case "YES_NO":
		yesCount, noCount, err := repo.GetYesNoStats(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		qStats.YesCount = &yesCount
		qStats.NoCount = &noCount
		total := yesCount + noCount
		if total > 0 {
			pct := float64(yesCount) / float64(total) * 100
			qStats.YesPercentage = &pct
		}

	case "NUMBER":
		avg, min, max, count, err := repo.GetNumberStats(ctx, q.ID)
		if err != nil {
			return nil, err
		}
		qStats.Average = &avg
		qStats.Min = &min
		qStats.Max = &max
		_ = count

	case "TEXT":
		samples, err := repo.GetTextSamples(ctx, q.ID, 5)
		if err != nil {
			return nil, err
		}
		qStats.SampleTexts = samples
	}

	return qStats, nil
}

func FormatQuestionStatSummary(qs *QuestionStats) string {
	switch qs.Type {
	case "RATING":
		if qs.Average != nil {
			return fmt.Sprintf("Average: %.1f (Min: %.0f, Max: %.0f)", *qs.Average, *qs.Min, *qs.Max)
		}
		return "No responses"
	case "SINGLE_CHOICE", "MULTIPLE_CHOICE":
		var result string
		for _, d := range qs.Distribution {
			result += fmt.Sprintf("%s: %d (%.1f%%) ", d.OptionText, d.Count, d.Percentage)
		}
		return result
	case "YES_NO":
		if qs.YesPercentage != nil {
			return fmt.Sprintf("Yes: %d (%.1f%%) | No: %d (%.1f%%)",
				*qs.YesCount, *qs.YesPercentage,
				*qs.NoCount, 100-*qs.YesPercentage)
		}
		return "No responses"
	case "NUMBER":
		if qs.Average != nil {
			return fmt.Sprintf("Average: %.1f (Min: %.0f, Max: %.0f)", *qs.Average, *qs.Min, *qs.Max)
		}
		return "No responses"
	case "TEXT":
		return fmt.Sprintf("%d text responses", len(qs.SampleTexts))
	default:
		return ""
	}
}
