package surveys

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rrhhumand/api/internal/models"
)

type SurveyService struct {
	repo *SurveyRepository
}

func NewSurveyService(repo *SurveyRepository) *SurveyService {
	return &SurveyService{repo: repo}
}

func (s *SurveyService) CreateSurvey(ctx context.Context, companyID, userID string, req *CreateSurveyRequest) (*models.Survey, error) {
	anonymous := false
	if req.Anonymous != nil {
		anonymous = *req.Anonymous
	}
	multipleResponses := false
	if req.MultipleResponses != nil {
		multipleResponses = *req.MultipleResponses
	}

	survey := &models.Survey{
		ID:                uuid.New().String(),
		CompanyID:         companyID,
		Title:             req.Title,
		Description:       req.Description,
		Type:              req.Type,
		Status:            "DRAFT",
		Anonymous:         anonymous,
		MultipleResponses: multipleResponses,
		StartsAt:          req.StartsAt,
		EndsAt:            req.EndsAt,
		CreatedBy:         userID,
	}

	if err := s.repo.CreateSurvey(ctx, survey); err != nil {
		return nil, err
	}
	return survey, nil
}

func (s *SurveyService) GetSurveyByID(ctx context.Context, id, companyID string) (*models.Survey, error) {
	survey, err := s.repo.GetSurveyByID(ctx, id, companyID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("survey not found")
		}
		return nil, err
	}

	questions, _ := s.repo.ListQuestionsBySurveyID(ctx, survey.ID)
	for i := range questions {
		options, _ := s.repo.ListOptionsByQuestionID(ctx, questions[i].ID)
		questions[i].Options = options
	}
	survey.Questions = questions

	targets, _ := s.repo.ListTargetsBySurveyID(ctx, survey.ID)
	survey.Targets = targets

	responseCount, _ := s.repo.GetResponseCount(ctx, survey.ID)
	survey.ResponseCount = responseCount

	if len(targets) > 0 {
		targetCount, _ := s.repo.GetTargetedEmployeeCount(ctx, survey.ID)
		survey.TargetCount = targetCount
		if targetCount > 0 {
			survey.ParticipationRate = float64(responseCount) / float64(targetCount) * 100
		}
	}

	return survey, nil
}

func (s *SurveyService) ListSurveys(ctx context.Context, companyID string, filters SurveyFilters, params *models.PaginationParams) ([]models.Survey, int64, error) {
	return s.repo.ListSurveys(ctx, companyID, filters, params.Offset, params.PerPage)
}

func (s *SurveyService) UpdateSurvey(ctx context.Context, id, companyID string, req *UpdateSurveyRequest) (*models.Survey, error) {
	survey, err := s.repo.GetSurveyByID(ctx, id, companyID)
	if err != nil {
		return nil, errors.New("survey not found")
	}

	if survey.Status != "DRAFT" {
		return nil, errors.New("survey not editable")
	}

	if req.Title != nil {
		survey.Title = *req.Title
	}
	if req.Description != nil {
		survey.Description = req.Description
	}
	if req.Type != nil {
		survey.Type = *req.Type
	}
	if req.Anonymous != nil {
		survey.Anonymous = *req.Anonymous
	}
	if req.MultipleResponses != nil {
		survey.MultipleResponses = *req.MultipleResponses
	}
	if req.StartsAt != nil {
		survey.StartsAt = req.StartsAt
	}
	if req.EndsAt != nil {
		survey.EndsAt = req.EndsAt
	}

	if err := s.repo.UpdateSurvey(ctx, survey); err != nil {
		return nil, err
	}
	return s.GetSurveyByID(ctx, id, companyID)
}

func (s *SurveyService) DeleteSurvey(ctx context.Context, id, companyID string) error {
	survey, err := s.repo.GetSurveyByID(ctx, id, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "DRAFT" {
		return errors.New("only draft surveys can be deleted")
	}
	return s.repo.DeleteSurvey(ctx, id, companyID)
}

func (s *SurveyService) PublishSurvey(ctx context.Context, id, companyID string) error {
	survey, err := s.repo.GetSurveyByID(ctx, id, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "DRAFT" {
		return errors.New("survey invalid transition")
	}
	if len(survey.Questions) == 0 {
		return errors.New("survey no questions")
	}
	return s.repo.UpdateSurveyStatus(ctx, id, companyID, "PUBLISHED")
}

func (s *SurveyService) CloseSurvey(ctx context.Context, id, companyID string) error {
	survey, err := s.repo.GetSurveyByID(ctx, id, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "PUBLISHED" {
		return errors.New("survey invalid transition")
	}
	return s.repo.UpdateSurveyStatus(ctx, id, companyID, "CLOSED")
}

func (s *SurveyService) ArchiveSurvey(ctx context.Context, id, companyID string) error {
	survey, err := s.repo.GetSurveyByID(ctx, id, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "CLOSED" {
		return errors.New("survey invalid transition")
	}
	return s.repo.UpdateSurveyStatus(ctx, id, companyID, "ARCHIVED")
}

func (s *SurveyService) AddQuestion(ctx context.Context, surveyID, companyID string, req *CreateQuestionRequest) (*models.SurveyQuestion, error) {
	survey, err := s.repo.GetSurveyByID(ctx, surveyID, companyID)
	if err != nil {
		return nil, errors.New("survey not found")
	}
	if survey.Status != "DRAFT" {
		return nil, errors.New("survey not editable")
	}

	questions, _ := s.repo.ListQuestionsBySurveyID(ctx, surveyID)
	position := len(questions) + 1
	if req.Position != nil {
		position = *req.Position
	}

	required := false
	if req.Required != nil {
		required = *req.Required
	}

	q := &models.SurveyQuestion{
		ID:       uuid.New().String(),
		SurveyID: surveyID,
		Question: req.Question,
		Type:     req.Type,
		Position: position,
		Required: required,
	}

	if err := s.repo.CreateQuestion(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *SurveyService) UpdateQuestion(ctx context.Context, questionID string, req *UpdateQuestionRequest) (*models.SurveyQuestion, error) {
	q, err := s.repo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return nil, errors.New("question not found")
	}

	survey, err := s.repo.GetSurveyByID(ctx, q.SurveyID, "")
	if err == nil && survey.Status != "DRAFT" {
		count, _ := s.repo.GetQuestionResponseCount(ctx, questionID)
		if count > 0 {
			return nil, errors.New("survey question has responses")
		}
	}

	if req.Question != nil {
		q.Question = *req.Question
	}
	if req.Type != nil {
		q.Type = *req.Type
	}
	if req.Position != nil {
		q.Position = *req.Position
	}
	if req.Required != nil {
		q.Required = *req.Required
	}

	if err := s.repo.UpdateQuestion(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *SurveyService) DeleteQuestion(ctx context.Context, questionID, companyID string) error {
	q, err := s.repo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return errors.New("question not found")
	}

	survey, err := s.repo.GetSurveyByID(ctx, q.SurveyID, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "DRAFT" {
		return errors.New("survey not editable")
	}

	return s.repo.DeleteQuestion(ctx, questionID)
}

func (s *SurveyService) AddOption(ctx context.Context, questionID string, req *CreateOptionRequest) (*models.SurveyOption, error) {
	q, err := s.repo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return nil, errors.New("question not found")
	}

	survey, err := s.repo.GetSurveyByID(ctx, q.SurveyID, "")
	if err == nil && survey.Status != "DRAFT" {
		return nil, errors.New("survey not editable")
	}

	existingOptions, _ := s.repo.ListOptionsByQuestionID(ctx, questionID)
	position := len(existingOptions) + 1
	if req.Position != nil {
		position = *req.Position
	}

	o := &models.SurveyOption{
		ID:         uuid.New().String(),
		QuestionID: questionID,
		OptionText: req.OptionText,
		Position:   position,
	}

	if err := s.repo.CreateOption(ctx, o); err != nil {
		return nil, err
	}
	return o, nil
}

func (s *SurveyService) DeleteOption(ctx context.Context, optionID, companyID string) error {
	o, err := s.repo.GetQuestionByID(ctx, optionID)
	if err != nil {
		return errors.New("option not found")
	}

	survey, err := s.repo.GetSurveyByID(ctx, o.SurveyID, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "DRAFT" {
		return errors.New("survey not editable")
	}

	return s.repo.DeleteOption(ctx, optionID)
}

func (s *SurveyService) SetTargets(ctx context.Context, surveyID, companyID string, req *SetTargetsRequest) error {
	survey, err := s.repo.GetSurveyByID(ctx, surveyID, companyID)
	if err != nil {
		return errors.New("survey not found")
	}
	if survey.Status != "DRAFT" {
		return errors.New("survey not editable")
	}

	var targets []models.SurveyTarget
	for _, t := range req.Targets {
		targets = append(targets, models.SurveyTarget{
			SurveyID:   surveyID,
			TargetType: t.TargetType,
			TargetID:   t.TargetID,
		})
	}

	return s.repo.SetTargets(ctx, surveyID, targets)
}

func (s *SurveyService) ListTargets(ctx context.Context, surveyID, companyID string) ([]models.SurveyTarget, error) {
	_, err := s.repo.GetSurveyByID(ctx, surveyID, companyID)
	if err != nil {
		return nil, errors.New("survey not found")
	}
	return s.repo.ListTargetsBySurveyID(ctx, surveyID)
}

func (s *SurveyService) ListAvailableSurveys(ctx context.Context, companyID, employeeID string) ([]models.Survey, error) {
	return s.repo.ListAvailableSurveys(ctx, companyID, employeeID)
}

func (s *SurveyService) RespondSurvey(ctx context.Context, surveyID, companyID, employeeID string, req *RespondSurveyRequest, ipAddress, userAgent string) error {
	survey, err := s.repo.GetSurveyByID(ctx, surveyID, companyID)
	if err != nil {
		return errors.New("survey not found")
	}

	if survey.Status != "PUBLISHED" {
		return errors.New("survey closed")
	}

	now := time.Now()
	if survey.StartsAt != nil && now.Before(*survey.StartsAt) {
		return errors.New("survey not active")
	}
	if survey.EndsAt != nil && now.After(*survey.EndsAt) {
		return errors.New("survey not active")
	}

	isTargeted, err := s.repo.IsEmployeeTargeted(ctx, surveyID, employeeID)
	if err != nil {
		return err
	}
	if !isTargeted {
		return errors.New("employee not targeted")
	}

	if !survey.MultipleResponses {
		alreadyAnswered, err := s.repo.HasEmployeeResponded(ctx, surveyID, employeeID)
		if err != nil {
			return err
		}
		if alreadyAnswered {
			return errors.New("survey already answered")
		}
	}

	response := &models.SurveyResponse{
		ID:         uuid.New().String(),
		SurveyID:   surveyID,
		EmployeeID: employeeID,
	}
	if err := s.repo.CreateResponse(ctx, response); err != nil {
		return err
	}

	for _, ans := range req.Answers {
		answer := &models.SurveyAnswer{
			ID:         uuid.New().String(),
			ResponseID: response.ID,
			QuestionID: ans.QuestionID,
			TextValue:  ans.Text,
			NumberValue: ans.Number,
			OptionID:   ans.OptionID,
		}
		if err := s.repo.CreateAnswer(ctx, answer); err != nil {
			return err
		}

		if len(ans.OptionIDs) > 0 {
			for _, optID := range ans.OptionIDs {
				_ = s.repo.CreateAnswerOption(ctx, &models.SurveyAnswerOption{
					AnswerID: answer.ID,
					OptionID: optID,
				})
			}
		}
	}

	return nil
}

func (s *SurveyService) GetResults(ctx context.Context, surveyID, companyID string, statsService *SurveyStatsService) (*SurveyStats, error) {
	survey, err := s.repo.GetSurveyByID(ctx, surveyID, companyID)
	if err != nil {
		return nil, errors.New("survey not found")
	}

	return statsService.CalculateStats(ctx, s.repo, survey)
}

func (s *SurveyService) ExportCSV(ctx context.Context, surveyID, companyID string) ([][]string, error) {
	survey, err := s.repo.GetSurveyByID(ctx, surveyID, companyID)
	if err != nil {
		return nil, errors.New("survey not found")
	}

	questions, err := s.repo.ListQuestionsBySurveyID(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	responses, err := s.repo.GetResponsesByEmployeeForExport(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	var responseIDs []string
	for _, r := range responses {
		responseIDs = append(responseIDs, r.ID)
	}

	answersMap, err := s.repo.GetAnswersForExport(ctx, responseIDs)
	if err != nil {
		return nil, err
	}

	headers := []string{"response_id"}
	if !survey.Anonymous {
		headers = append(headers, "employee_id", "employee_name")
	}
	for _, q := range questions {
		headers = append(headers, q.Question)
	}

	var rows [][]string
	rows = append(rows, headers)

	for _, resp := range responses {
		row := []string{resp.ID}
		if !survey.Anonymous {
			empName, _ := s.repo.GetEmployeeName(ctx, resp.EmployeeID)
			row = append(row, resp.EmployeeID, empName)
		}

		answers := answersMap[resp.ID]
		answerByQuestion := make(map[string]models.SurveyAnswer)
		for _, a := range answers {
			answerByQuestion[a.QuestionID] = a
		}

		for _, q := range questions {
			a, exists := answerByQuestion[q.ID]
			if !exists {
				row = append(row, "")
				continue
			}

			switch q.Type {
			case "TEXT":
				if a.TextValue != nil {
					row = append(row, *a.TextValue)
				} else {
					row = append(row, "")
				}
			case "NUMBER", "RATING":
				if a.NumberValue != nil {
					row = append(row, fmt.Sprintf("%.0f", *a.NumberValue))
				} else {
					row = append(row, "")
				}
			case "SINGLE_CHOICE", "YES_NO":
				if a.OptionID != nil {
					opts, _ := s.repo.ListOptionsByQuestionID(ctx, q.ID)
					for _, o := range opts {
						if o.ID == *a.OptionID {
							row = append(row, o.OptionText)
							break
						}
					}
				} else {
					row = append(row, "")
				}
			case "MULTIPLE_CHOICE":
				answerOpts, _ := s.repo.GetAnswerOptions(ctx, a.ID)
				var texts []string
				for _, ao := range answerOpts {
					opts, _ := s.repo.ListOptionsByQuestionID(ctx, q.ID)
					for _, o := range opts {
						if o.ID == ao.OptionID {
							texts = append(texts, o.OptionText)
							break
						}
					}
				}
				if len(texts) > 0 {
					row = append(row, fmt.Sprintf("%v", texts))
				} else {
					row = append(row, "")
				}
			default:
				row = append(row, "")
			}
		}

		rows = append(rows, row)
	}

	return rows, nil
}
