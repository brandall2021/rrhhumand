package performancehttp

import (
	"time"

	"github.com/rrhhumand/api/internal/performance/domain"
)

type CreateCycleRequest struct {
	Name                  string     `json:"name" binding:"required"`
	Description           *string    `json:"description,omitempty"`
	CycleType             string     `json:"cycle_type,omitempty"`
	StartDate             *time.Time `json:"start_date,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty"`
	EvaluationStartDate   *time.Time `json:"evaluation_start_date,omitempty"`
	EvaluationEndDate     *time.Time `json:"evaluation_end_date,omitempty"`
	ReviewStartDate       *time.Time `json:"review_start_date,omitempty"`
	ReviewEndDate         *time.Time `json:"review_end_date,omitempty"`
	CalibrationStartDate  *time.Time `json:"calibration_start_date,omitempty"`
	CalibrationEndDate    *time.Time `json:"calibration_end_date,omitempty"`
	TemplateID            *string    `json:"template_id,omitempty"`
	ObjectiveWeight       *float64   `json:"objective_weight,omitempty"`
	CompetencyWeight      *float64   `json:"competency_weight,omitempty"`
	MinAnonymousResponses *int       `json:"min_anonymous_responses,omitempty"`
}

type UpdateCycleRequest struct {
	Name                  *string    `json:"name,omitempty"`
	Description           *string    `json:"description,omitempty"`
	CycleType             *string    `json:"cycle_type,omitempty"`
	StartDate             *time.Time `json:"start_date,omitempty"`
	EndDate               *time.Time `json:"end_date,omitempty"`
	EvaluationStartDate   *time.Time `json:"evaluation_start_date,omitempty"`
	EvaluationEndDate     *time.Time `json:"evaluation_end_date,omitempty"`
	ReviewStartDate       *time.Time `json:"review_start_date,omitempty"`
	ReviewEndDate         *time.Time `json:"review_end_date,omitempty"`
	CalibrationStartDate  *time.Time `json:"calibration_start_date,omitempty"`
	CalibrationEndDate    *time.Time `json:"calibration_end_date,omitempty"`
	TemplateID            *string    `json:"template_id,omitempty"`
	ObjectiveWeight       *float64   `json:"objective_weight,omitempty"`
	CompetencyWeight      *float64   `json:"competency_weight,omitempty"`
}

type CycleStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateObjectiveRequest struct {
	CycleID          string     `json:"cycle_id" binding:"required"`
	EmployeeID       string     `json:"employee_id" binding:"required"`
	ParentObjectiveID *string   `json:"parent_objective_id,omitempty"`
	Title            string     `json:"title" binding:"required"`
	Description      *string    `json:"description,omitempty"`
	ObjectiveType    *string    `json:"objective_type,omitempty"`
	Weight           *float64   `json:"weight,omitempty"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	TargetValue      *float64   `json:"target_value,omitempty"`
	Unit             *string    `json:"unit,omitempty"`
	ProgressType     *string    `json:"progress_type,omitempty"`
	KeyResults       []CreateKeyResultRequest `json:"key_results,omitempty"`
}

type CreateKeyResultRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description *string  `json:"description,omitempty"`
	Weight      *float64 `json:"weight,omitempty"`
	TargetValue *float64 `json:"target_value,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	SortOrder   *int     `json:"sort_order,omitempty"`
}

type UpdateObjectiveRequest struct {
	Title       *string   `json:"title,omitempty"`
	Description *string   `json:"description,omitempty"`
	Weight      *float64  `json:"weight,omitempty"`
	Status      *string   `json:"status,omitempty"`
	TargetValue *float64  `json:"target_value,omitempty"`
	CurrentValue *float64 `json:"current_value,omitempty"`
	Unit        *string   `json:"unit,omitempty"`
	Notes       *string   `json:"notes,omitempty"`
	RiskNotes   *string   `json:"risk_notes,omitempty"`
}

type UpdateProgressRequest struct {
	CurrentValue float64 `json:"current_value" binding:"required"`
}

type AssignEvaluatorsRequest struct {
	CycleID    string                   `json:"cycle_id" binding:"required"`
	Assignments []EvaluatorAssignment   `json:"assignments" binding:"required"`
}

type EvaluatorAssignment struct {
	EmployeeID     string `json:"employee_id" binding:"required"`
	EvaluatorID    string `json:"evaluator_id" binding:"required"`
	EvaluationType string `json:"evaluation_type" binding:"required"`
}

type CreateEvaluationRequest struct {
	CycleID        string  `json:"cycle_id" binding:"required"`
	EmployeeID     string  `json:"employee_id" binding:"required"`
	EvaluatorID    string  `json:"evaluator_id" binding:"required"`
	EvaluationType string  `json:"evaluation_type" binding:"required"`
	TemplateID     *string `json:"template_id,omitempty"`
	Answers       []CreateAnswerRequest `json:"answers,omitempty"`
}

type CreateAnswerRequest struct {
	QuestionID    *string  `json:"question_id,omitempty"`
	NumericValue  *float64 `json:"numeric_value,omitempty"`
	TextValue     *string  `json:"text_value,omitempty"`
	SelectedValue *string  `json:"selected_value,omitempty"`
	BooleanValue  *bool    `json:"boolean_value,omitempty"`
}

type CreateFeedbackRequest struct {
	EmployeeID   string  `json:"employee_id" binding:"required"`
	CycleID      *string `json:"cycle_id,omitempty"`
	FeedbackType string  `json:"feedback_type" binding:"required"`
	Visibility   string  `json:"visibility,omitempty"`
	Content      string  `json:"content" binding:"required"`
	IsAnonymous  *bool   `json:"is_anonymous,omitempty"`
}

type CreateCheckInRequest struct {
	EmployeeID string    `json:"employee_id" binding:"required"`
	ManagerID  string    `json:"manager_id" binding:"required"`
	CycleID    *string   `json:"cycle_id,omitempty"`
	ScheduledAt time.Time `json:"scheduled_at" binding:"required"`
}

type CompleteCheckInRequest struct {
	EmployeeNotes *string `json:"employee_notes,omitempty"`
	ManagerNotes  *string `json:"manager_notes,omitempty"`
	Achievements  *string `json:"achievements,omitempty"`
	Blockers      *string `json:"blockers,omitempty"`
	NextSteps     *string `json:"next_steps,omitempty"`
}

type CreateCalibrationRequest struct {
	CycleID     string  `json:"cycle_id" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description,omitempty"`
}

type CreateImprovementPlanRequest struct {
	EmployeeID      string    `json:"employee_id" binding:"required"`
	CycleID         *string   `json:"cycle_id,omitempty"`
	Reason          string    `json:"reason" binding:"required"`
	StartDate       time.Time `json:"start_date" binding:"required"`
	EndDate         time.Time `json:"end_date" binding:"required"`
	SuccessCriteria *string   `json:"success_criteria,omitempty"`
	Actions         []CreatePlanActionRequest `json:"actions,omitempty"`
}

type CreatePlanActionRequest struct {
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description,omitempty"`
	DueDate      *time.Time `json:"due_date,omitempty"`
	ResponsibleID *string `json:"responsible_id,omitempty"`
}

type CreateDevelopmentPlanRequest struct {
	EmployeeID  string    `json:"employee_id" binding:"required"`
	CycleID     *string   `json:"cycle_id,omitempty"`
	Title       string    `json:"title" binding:"required"`
	Description *string   `json:"description,omitempty"`
	CareerGoal  *string   `json:"career_goal,omitempty"`
	CurrentLevel *int     `json:"current_level,omitempty"`
	TargetLevel  *int     `json:"target_level,omitempty"`
	CompetencyID *string  `json:"competency_id,omitempty"`
	Actions     []CreateDevActionRequest `json:"actions,omitempty"`
}

type CreateDevActionRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description,omitempty"`
	ActionType  string  `json:"action_type,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type CreateEvidenceRequest struct {
	EvaluationID *string `json:"evaluation_id,omitempty"`
	ObjectiveID  *string `json:"objective_id,omitempty"`
	FeedbackID   *string `json:"feedback_id,omitempty"`
	Title        string  `json:"title" binding:"required"`
	Description  *string `json:"description,omitempty"`
	EvidenceType string  `json:"evidence_type,omitempty"`
	URL          *string `json:"url,omitempty"`
}

type CreateRecognitionRequest struct {
	EmployeeID      string `json:"employee_id" binding:"required"`
	RecognitionType string `json:"recognition_type" binding:"required"`
	Message         string `json:"message" binding:"required"`
}

type CalculateResultRequest struct {
	CycleID    string `json:"cycle_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
}

type CreateReviewRequest struct {
	CycleID    string `json:"cycle_id" binding:"required"`
	EmployeeID string `json:"employee_id" binding:"required"`
	ManagerID  string `json:"manager_id" binding:"required"`
}

type UpdateReviewRequest struct {
	Summary           *string `json:"summary,omitempty"`
	Strengths         *string `json:"strengths,omitempty"`
	ImprovementAreas  *string `json:"improvement_areas,omitempty"`
	FinalScore        *float64 `json:"final_score,omitempty"`
	FinalRating       *string `json:"final_rating,omitempty"`
	EmployeeComments  *string `json:"employee_comments,omitempty"`
	ManagerComments   *string `json:"manager_comments,omitempty"`
	EmployeeAgreement *string `json:"employee_agreement,omitempty"`
}

func ToDomainCycle(req *CreateCycleRequest, companyID, userID string) *domain.PerformanceCycle {
	c := &domain.PerformanceCycle{
		CompanyID:             companyID,
		Name:                  req.Name,
		Description:           req.Description,
		CycleType:             domain.CycleType(req.CycleType),
		StartDate:             req.StartDate,
		EndDate:               req.EndDate,
		EvaluationStartDate:   req.EvaluationStartDate,
		EvaluationEndDate:     req.EvaluationEndDate,
		ReviewStartDate:       req.ReviewStartDate,
		ReviewEndDate:         req.ReviewEndDate,
		CalibrationStartDate:  req.CalibrationStartDate,
		CalibrationEndDate:    req.CalibrationEndDate,
		TemplateID:            req.TemplateID,
		CreatedBy:             userID,
	}
	if req.ObjectiveWeight != nil {
		c.ObjectiveWeight = *req.ObjectiveWeight
	}
	if req.CompetencyWeight != nil {
		c.CompetencyWeight = *req.CompetencyWeight
	}
	if req.MinAnonymousResponses != nil {
		c.MinAnonymousResponses = *req.MinAnonymousResponses
	}
	return c
}

func ToDomainObjective(req *CreateObjectiveRequest, companyID, userID string) *domain.PerformanceObjective {
	o := &domain.PerformanceObjective{
		CompanyID:         companyID,
		CycleID:           req.CycleID,
		EmployeeID:        req.EmployeeID,
		ParentObjectiveID: req.ParentObjectiveID,
		Title:             req.Title,
		Description:       req.Description,
		StartDate:         req.StartDate,
		DueDate:           req.DueDate,
		TargetValue:       req.TargetValue,
		Unit:              req.Unit,
		CreatedBy:         userID,
	}
	if req.ObjectiveType != nil {
		o.ObjectiveType = domain.ObjectiveType(*req.ObjectiveType)
	}
	if req.Weight != nil {
		o.Weight = *req.Weight
	}
	if req.ProgressType != nil {
		o.ProgressType = domain.ProgressType(*req.ProgressType)
	}
	return o
}

func ToDomainParticipants(req *AssignEvaluatorsRequest, companyID string) []domain.PerformanceParticipant {
	participants := make([]domain.PerformanceParticipant, len(req.Assignments))
	for i, a := range req.Assignments {
		participants[i] = domain.PerformanceParticipant{
			CompanyID:      companyID,
			CycleID:        req.CycleID,
			EmployeeID:     a.EmployeeID,
			EvaluatorID:    a.EvaluatorID,
			EvaluationType: domain.EvaluationType(a.EvaluationType),
		}
	}
	return participants
}

func ToDomainEvaluation(req *CreateEvaluationRequest, companyID string) *domain.PerformanceEvaluation {
	return &domain.PerformanceEvaluation{
		CompanyID:      companyID,
		CycleID:        req.CycleID,
		EmployeeID:     req.EmployeeID,
		EvaluatorID:    req.EvaluatorID,
		EvaluationType: domain.EvaluationType(req.EvaluationType),
		TemplateID:     req.TemplateID,
	}
}

func ToDomainAnswers(evaluationID string, reqs []CreateAnswerRequest) []domain.EvaluationAnswer {
	answers := make([]domain.EvaluationAnswer, len(reqs))
	for i, r := range reqs {
		answers[i] = domain.EvaluationAnswer{
			EvaluationID:  evaluationID,
			QuestionID:    r.QuestionID,
			NumericValue:  r.NumericValue,
			TextValue:     r.TextValue,
			SelectedValue: r.SelectedValue,
			BooleanValue:  r.BooleanValue,
		}
	}
	return answers
}

func ToDomainFeedback(req *CreateFeedbackRequest, companyID, authorID string) *domain.PerformanceFeedback {
	return &domain.PerformanceFeedback{
		CompanyID:    companyID,
		CycleID:      req.CycleID,
		EmployeeID:   req.EmployeeID,
		AuthorID:     authorID,
		FeedbackType: domain.FeedbackType(req.FeedbackType),
		Content:      req.Content,
		IsAnonymous:  req.IsAnonymous != nil && *req.IsAnonymous,
	}
}

func ToDomainCheckIn(req *CreateCheckInRequest, companyID string) *domain.PerformanceCheckIn {
	return &domain.PerformanceCheckIn{
		CompanyID:  companyID,
		EmployeeID: req.EmployeeID,
		ManagerID:  req.ManagerID,
		CycleID:    req.CycleID,
		ScheduledAt: req.ScheduledAt,
	}
}

func ToDomainCalibration(req *CreateCalibrationRequest, companyID, userID string) *domain.CalibrationSession {
	return &domain.CalibrationSession{
		CompanyID:   companyID,
		CycleID:     req.CycleID,
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}
}

func ToDomainImprovementPlan(req *CreateImprovementPlanRequest, companyID, userID string) *domain.ImprovementPlan {
	return &domain.ImprovementPlan{
		CompanyID:      companyID,
		EmployeeID:     req.EmployeeID,
		CycleID:        req.CycleID,
		CreatedBy:      userID,
		Reason:         req.Reason,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		SuccessCriteria: req.SuccessCriteria,
	}
}

func ToDomainDevPlan(req *CreateDevelopmentPlanRequest, companyID, userID string) *domain.DevelopmentPlan {
	return &domain.DevelopmentPlan{
		CompanyID:    companyID,
		EmployeeID:   req.EmployeeID,
		CycleID:      req.CycleID,
		CreatedBy:    userID,
		Title:        req.Title,
		Description:  req.Description,
		CareerGoal:   req.CareerGoal,
		CurrentLevel: req.CurrentLevel,
		TargetLevel:  req.TargetLevel,
		CompetencyID: req.CompetencyID,
	}
}

func ToDomainEvidence(req *CreateEvidenceRequest, companyID, userID string) *domain.PerformanceEvidence {
	return &domain.PerformanceEvidence{
		CompanyID:    companyID,
		EvaluationID: req.EvaluationID,
		ObjectiveID:  req.ObjectiveID,
		FeedbackID:   req.FeedbackID,
		Title:        req.Title,
		Description:  req.Description,
		EvidenceType: domain.EvidenceType(req.EvidenceType),
		URL:          req.URL,
		CreatedBy:    userID,
	}
}

func ToDomainRecognition(req *CreateRecognitionRequest, companyID, authorID string) *domain.PerformanceRecognition {
	return &domain.PerformanceRecognition{
		CompanyID:       companyID,
		EmployeeID:      req.EmployeeID,
		AuthorID:        authorID,
		RecognitionType: domain.RecognitionType(req.RecognitionType),
		Message:         req.Message,
	}
}

func ToDomainReview(req *CreateReviewRequest, companyID string) *domain.PerformanceReview {
	return &domain.PerformanceReview{
		CompanyID:  companyID,
		CycleID:    req.CycleID,
		EmployeeID: req.EmployeeID,
		ManagerID:  req.ManagerID,
	}
}
