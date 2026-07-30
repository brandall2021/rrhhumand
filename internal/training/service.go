package training

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"github.com/rrhhumand/api/internal/events"
	"github.com/rrhhumand/api/internal/notifications"
)

type Service struct {
	repo   *Repository
	evts   *events.Service
	notif  *notifications.Service
	storage StorageService
	log    *zap.Logger
}

type StorageService interface {
	Upload(ctx context.Context, bucket, key string, data []byte, contentType string) (string, error)
	Delete(ctx context.Context, bucket, key string) error
}

func NewService(repo *Repository, evts *events.Service, notif *notifications.Service, storage StorageService, log *zap.Logger) *Service {
	return &Service{repo: repo, evts: evts, notif: notif, storage: storage, log: log}
}

func svcErr(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("training_svc.%s: %w", op, err)
}

// ---------------------------------------------------------------------------
// Categories
// ---------------------------------------------------------------------------

func (s *Service) CreateCategory(ctx context.Context, companyID, userID string, req CreateCategoryRequest) (*Category, error) {
	c := &Category{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		Active:      true,
		CreatedBy:   userID,
	}
	if err := s.repo.CreateCategory(ctx, c); err != nil {
		return nil, svcErr("CreateCategory", err)
	}
	return c, nil
}

func (s *Service) UpdateCategory(ctx context.Context, companyID, id string, req UpdateCategoryRequest) (*Category, error) {
	c, err := s.repo.GetCategory(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateCategory", err)
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.Active != nil {
		c.Active = *req.Active
	}
	if err := s.repo.UpdateCategory(ctx, c); err != nil {
		return nil, svcErr("UpdateCategory", err)
	}
	return c, nil
}

func (s *Service) GetCategory(ctx context.Context, companyID, id string) (*Category, error) {
	c, err := s.repo.GetCategory(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("GetCategory", err)
	}
	return c, nil
}

func (s *Service) ListCategories(ctx context.Context, companyID string, parentID *string, active *bool) ([]Category, error) {
	return s.repo.ListCategories(ctx, companyID, parentID, active)
}

// ---------------------------------------------------------------------------
// Courses
// ---------------------------------------------------------------------------

func (s *Service) CreateCourse(ctx context.Context, companyID, userID string, req CreateCourseRequest) (*Course, error) {
	c := &Course{
		ID:                     uuid.New().String(),
		CompanyID:              companyID,
		Code:                   req.Code,
		Name:                   req.Name,
		CategoryID:             req.CategoryID,
		ShortDescription:       req.ShortDescription,
		Description:            req.Description,
		Objectives:             req.Objectives,
		Difficulty:             req.Difficulty,
		DurationMinutes:        req.DurationMinutes,
		Modality:               req.Modality,
		Status:                 "draft",
		Mandatory:              false,
		PassingScore:           req.PassingScore,
		CertificateEnabled:     false,
		MinAttendancePercentage: 0,
		CreatedBy:   userID,
	}
	if req.Mandatory != nil {
		c.Mandatory = *req.Mandatory
	}
	if req.CertificateEnabled != nil {
		c.CertificateEnabled = *req.CertificateEnabled
	}
	if req.MinAttendancePercentage != nil {
		c.MinAttendancePercentage = *req.MinAttendancePercentage
	}
	if err := s.repo.CreateCourse(ctx, c); err != nil {
		return nil, svcErr("CreateCourse", err)
	}
	return c, nil
}

func (s *Service) UpdateCourse(ctx context.Context, companyID, id string, req UpdateCourseRequest) (*Course, error) {
	c, err := s.repo.GetCourse(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateCourse", err)
	}
	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.CategoryID != nil {
		c.CategoryID = req.CategoryID
	}
	if req.ShortDescription != nil {
		c.ShortDescription = req.ShortDescription
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.Objectives != nil {
		c.Objectives = req.Objectives
	}
	if req.Difficulty != nil {
		c.Difficulty = *req.Difficulty
	}
	if req.DurationMinutes != nil {
		c.DurationMinutes = *req.DurationMinutes
	}
	if req.Modality != nil {
		c.Modality = *req.Modality
	}
	if req.Status != nil {
		c.Status = *req.Status
	}
	if req.Mandatory != nil {
		c.Mandatory = *req.Mandatory
	}
	if req.PassingScore != nil {
		c.PassingScore = req.PassingScore
	}
	if req.CertificateEnabled != nil {
		c.CertificateEnabled = *req.CertificateEnabled
	}
	if req.MinAttendancePercentage != nil {
		c.MinAttendancePercentage = *req.MinAttendancePercentage
	}
	if err := s.repo.UpdateCourse(ctx, c); err != nil {
		return nil, svcErr("UpdateCourse", err)
	}
	return c, nil
}

func (s *Service) GetCourse(ctx context.Context, companyID, id string) (*Course, error) {
	return s.repo.GetCourse(ctx, companyID, id)
}

func (s *Service) ListCourses(ctx context.Context, companyID string, filter CourseFilter) ([]Course, int, error) {
	return s.repo.ListCourses(ctx, companyID, filter)
}

func (s *Service) PublishCourse(ctx context.Context, companyID, id, userID string) error {
	c, err := s.repo.GetCourse(ctx, companyID, id)
	if err != nil {
		return svcErr("PublishCourse", err)
	}
	if c.Status == "published" {
		return fmt.Errorf("course already published")
	}
	// Verify at least one version exists
	vers, err := s.repo.ListVersions(ctx, id)
	if err != nil {
		return svcErr("PublishCourse", err)
	}
	if len(vers) == 0 {
		return fmt.Errorf("cannot publish course without at least one version")
	}
	if err := s.repo.PublishCourse(ctx, companyID, id, userID); err != nil {
		return svcErr("PublishCourse", err)
	}
	s.createEvent(ctx, companyID, "course_published", "Course Published", fmt.Sprintf("Course %s has been published", c.Name), userID, "", "", "info", nil)
	return nil
}

func (s *Service) GetCourseWithDetails(ctx context.Context, companyID, id string) (*CourseWithDetails, error) {
	return s.repo.GetCourseWithDetails(ctx, companyID, id)
}

// ---------------------------------------------------------------------------
// Versions
// ---------------------------------------------------------------------------

func (s *Service) CreateVersion(ctx context.Context, courseID, userID string, req CreateVersionRequest) (*CourseVersion, error) {
	// Unpublish current active version
	av, err := s.repo.GetActiveVersion(ctx, courseID)
	if err != nil && err.Error() != "training_repo.GetActiveVersion: no rows in result set" {
		return nil, svcErr("CreateVersion", err)
	}
	if av != nil {
		_ = s.repo.PublishVersion(ctx, av.ID) // keep it published, new version will be the active one
	}
	v := &CourseVersion{
		ID:          uuid.New().String(),
		CourseID:    courseID,
		Version:     req.Version,
		Description: req.Description,
		IsPublished: true,
		CreatedBy:   &userID,
	}
	if err := s.repo.CreateVersion(ctx, v); err != nil {
		return nil, svcErr("CreateVersion", err)
	}
	return v, nil
}

func (s *Service) ListVersions(ctx context.Context, courseID string) ([]CourseVersion, error) {
	return s.repo.ListVersions(ctx, courseID)
}

// ---------------------------------------------------------------------------
// Contents
// ---------------------------------------------------------------------------

func (s *Service) CreateContent(ctx context.Context, versionID, userID string, req CreateContentRequest) (*CourseContent, error) {
	c := &CourseContent{
		ID:              uuid.New().String(),
		CourseVersionID: versionID,
		Title:           req.Title,
		Description:     req.Description,
		ContentType:     req.ContentType,
		ExternalURL:     req.ExternalURL,
		DurationSeconds: 0,
		SortOrder:       0,
		Required:        false,
		Published:       false,
		CreatedBy:       &userID,
	}
	if err := s.repo.CreateContent(ctx, c); err != nil {
		return nil, svcErr("CreateContent", err)
	}
	return c, nil
}

func (s *Service) UpdateContent(ctx context.Context, id string, req UpdateContentRequest) (*CourseContent, error) {
	c, err := s.repo.GetContent(ctx, id)
	if err != nil {
		return nil, svcErr("UpdateContent", err)
	}
	if req.Title != nil {
		c.Title = *req.Title
	}
	if req.Description != nil {
		c.Description = req.Description
	}
	if req.ContentType != nil {
		c.ContentType = *req.ContentType
	}
	if req.ExternalURL != nil {
		c.ExternalURL = req.ExternalURL
	}
	if req.DurationSeconds != nil {
		c.DurationSeconds = *req.DurationSeconds
	}
	if req.SortOrder != nil {
		c.SortOrder = *req.SortOrder
	}
	if req.Required != nil {
		c.Required = *req.Required
	}
	if req.Published != nil {
		c.Published = *req.Published
	}
	if err := s.repo.UpdateContent(ctx, c); err != nil {
		return nil, svcErr("UpdateContent", err)
	}
	return c, nil
}

func (s *Service) ListContents(ctx context.Context, versionID string, publishedOnly bool) ([]CourseContent, error) {
	return s.repo.ListContents(ctx, versionID, publishedOnly)
}

func (s *Service) ReorderContents(ctx context.Context, versionID string, contentIDs []string) error {
	return s.repo.ReorderContents(ctx, versionID, contentIDs)
}

// ---------------------------------------------------------------------------
// Offerings
// ---------------------------------------------------------------------------

func (s *Service) CreateOffering(ctx context.Context, companyID, userID string, req CreateOfferingRequest) (*CourseOffering, error) {
	o := &CourseOffering{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		CourseID:        req.CourseID,
		CourseVersionID: req.CourseVersionID,
		Name:            req.Name,
		StartDate:       parseTimePtr(req.StartDate),
		EndDate:         parseTimePtr(req.EndDate),
		Modality:        req.Modality,
		Location:        req.Location,
		MeetingURL:      req.MeetingURL,
		InstructorID:    req.InstructorID,
		ProviderID:      req.ProviderID,
		CostAmount:      req.CostAmount,
		Status:          "draft",
		CreatedBy:       &userID,
	}
	if req.Capacity != nil {
		o.Capacity = *req.Capacity
	}
	if req.CostCurrency != nil {
		o.CostCurrency = *req.CostCurrency
	}
	if err := s.repo.CreateOffering(ctx, o); err != nil {
		return nil, svcErr("CreateOffering", err)
	}
	return o, nil
}

func (s *Service) UpdateOffering(ctx context.Context, companyID, id string, req UpdateOfferingRequest) (*CourseOffering, error) {
	o, err := s.repo.GetOffering(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateOffering", err)
	}
	if req.Name != nil {
		o.Name = *req.Name
	}
	if req.StartDate != nil {
		o.StartDate = parseTimePtr(req.StartDate)
	}
	if req.EndDate != nil {
		o.EndDate = parseTimePtr(req.EndDate)
	}
	if req.Capacity != nil {
		o.Capacity = *req.Capacity
	}
	if req.Modality != nil {
		o.Modality = req.Modality
	}
	if req.Location != nil {
		o.Location = req.Location
	}
	if req.MeetingURL != nil {
		o.MeetingURL = req.MeetingURL
	}
	if req.InstructorID != nil {
		o.InstructorID = req.InstructorID
	}
	if req.ProviderID != nil {
		o.ProviderID = req.ProviderID
	}
	if req.CostAmount != nil {
		o.CostAmount = req.CostAmount
	}
	if req.CostCurrency != nil {
		o.CostCurrency = *req.CostCurrency
	}
	if req.Status != nil {
		o.Status = *req.Status
	}
	if err := s.repo.UpdateOffering(ctx, o); err != nil {
		return nil, svcErr("UpdateOffering", err)
	}
	return o, nil
}

func (s *Service) GetOffering(ctx context.Context, companyID, id string) (*CourseOffering, error) {
	return s.repo.GetOffering(ctx, companyID, id)
}

func (s *Service) ListOfferings(ctx context.Context, companyID string, filter OfferingFilter) ([]CourseOffering, int, error) {
	return s.repo.ListOfferings(ctx, companyID, filter)
}

// ---------------------------------------------------------------------------
// Sessions
// ---------------------------------------------------------------------------

func (s *Service) CreateSession(ctx context.Context, offeringID, userID string, req CreateSessionRequest) (*OfferingSession, error) {
	sessionDate, _ := time.Parse("2006-01-02", req.SessionDate)
	sess := &OfferingSession{
		ID:           uuid.New().String(),
		OfferingID:   offeringID,
		Title:        req.Title,
		SessionDate:  sessionDate,
		StartTime:    req.StartTime,
		EndTime:      req.EndTime,
		Location:     req.Location,
		MeetingURL:   req.MeetingURL,
		InstructorID: req.InstructorID,
		CreatedBy:    userID,
	}
	if err := s.repo.CreateSession(ctx, sess); err != nil {
		return nil, svcErr("CreateSession", err)
	}
	return sess, nil
}

func (s *Service) ListSessions(ctx context.Context, offeringID string) ([]OfferingSession, error) {
	return s.repo.ListSessionsByOffering(ctx, offeringID)
}

// ---------------------------------------------------------------------------
// Enrollments
// ---------------------------------------------------------------------------

func (s *Service) Enroll(ctx context.Context, companyID, offeringID, userID string, req EnrollRequest) (*Enrollment, error) {
	// Check capacity
	o, err := s.repo.GetOffering(ctx, companyID, offeringID)
	if err != nil {
		return nil, svcErr("Enroll", err)
	}
	if o.Capacity > 0 && o.EnrolledCount >= o.Capacity {
		return nil, fmt.Errorf("offering at full capacity")
	}
	e := &Enrollment{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		OfferingID:     offeringID,
		EmployeeID:     req.EmployeeID,
		AssignmentType: req.AssignmentType,
		Status:         "enrolled",
		CreatedBy:      &userID,
	}
	if e.AssignmentType == "" {
		e.AssignmentType = "voluntary"
	}
	if err := s.repo.Enroll(ctx, e); err != nil {
		return nil, svcErr("Enroll", err)
	}
	s.createEvent(ctx, companyID, "enrollment_created", "New Enrollment",
		fmt.Sprintf("Employee enrolled in offering %s", o.Name), userID, e.EmployeeID, "", "info", nil)
	return e, nil
}

func (s *Service) CompleteEnrollment(ctx context.Context, companyID, id string) error {
	e, err := s.repo.GetEnrollment(ctx, companyID, id)
	if err != nil {
		return svcErr("CompleteEnrollment", err)
	}
	e.Status = "completed"
	now := time.Now()
	e.CompletedAt = &now
	if err := s.repo.UpdateEnrollmentStatus(ctx, id, "completed"); err != nil {
		return svcErr("CompleteEnrollment", err)
	}
	// Check if certificate should be issued
	o, err := s.repo.GetOffering(ctx, companyID, e.OfferingID)
	if err != nil {
		return svcErr("CompleteEnrollment", err)
	}
	c, err := s.repo.GetCourse(ctx, companyID, o.CourseID)
	if err != nil {
		return svcErr("CompleteEnrollment", err)
	}
	if c.CertificateEnabled {
		// Generate certificate URL placeholder — actual PDF gen would go here
		certURL := fmt.Sprintf("certificates/%s/%s.pdf", companyID, id)
		if err := s.repo.UpdateCertificate(ctx, id, certURL); err != nil {
			return svcErr("CompleteEnrollment", err)
		}
	}
	createdBy := ""
	if e.CreatedBy != nil {
		createdBy = *e.CreatedBy
	}
	s.createEvent(ctx, companyID, "enrollment_completed", "Course Completed",
		fmt.Sprintf("Employee completed course"), createdBy, e.EmployeeID, "", "success", nil)
	return nil
}

func (s *Service) GetEnrollment(ctx context.Context, companyID, id string) (*Enrollment, error) {
	return s.repo.GetEnrollment(ctx, companyID, id)
}

func (s *Service) ListEnrollments(ctx context.Context, companyID string, filter EnrollmentFilter) ([]Enrollment, int, error) {
	return s.repo.ListEnrollments(ctx, companyID, filter)
}

// ---------------------------------------------------------------------------
// Content Progress
// ---------------------------------------------------------------------------

func (s *Service) UpdateProgress(ctx context.Context, enrollmentID string, req UpdateProgressRequest) error {
	if err := s.repo.UpsertContentProgress(ctx, &ContentProgress{
		ID:                 uuid.New().String(),
		EnrollmentID:       enrollmentID,
		ProgressPercentage: req.ProgressPercentage,
		TimeSpentSeconds:   req.TimeSpentSeconds,
		LastPosition:       req.LastPosition,
	}); err != nil {
		return svcErr("UpdateProgress", err)
	}
	return nil
}

func (s *Service) GetProgress(ctx context.Context, enrollmentID, contentID string) (*ContentProgress, error) {
	return s.repo.GetContentProgress(ctx, enrollmentID, contentID)
}

func (s *Service) ListProgress(ctx context.Context, enrollmentID string) ([]ContentProgress, error) {
	return s.repo.ListContentProgressByEnrollment(ctx, enrollmentID)
}

// ---------------------------------------------------------------------------
// Assignments (massive)
// ---------------------------------------------------------------------------

func (s *Service) CreateAssignment(ctx context.Context, companyID, userID string, req CreateAssignmentRequest) (*Assignment, error) {
	a := &Assignment{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		CourseID:       req.CourseID,
		AssigneeType:   req.AssigneeType,
		AssigneeID:     &req.AssigneeID,
		AssignmentType: req.AssignmentType,
		DueDate:        parseTimePtr(req.DueDate),
		CreatedBy:      userID,
	}
	if err := s.repo.CreateAssignment(ctx, a); err != nil {
		return nil, svcErr("CreateAssignment", err)
	}
	// Auto-enroll
	o, err := s.repo.GetOffering(ctx, companyID, a.CourseID)
	if err != nil && err != pgx.ErrNoRows {
		return nil, svcErr("CreateAssignment", err)
	}
	if err == nil {
		_, err := s.Enroll(ctx, companyID, o.ID, userID, EnrollRequest{
			EmployeeID:     req.AssigneeID,
			AssignmentType: req.AssignmentType,
		})
		if err != nil {
			return nil, svcErr("CreateAssignment", err)
		}
	}
	return a, nil
}

// ---------------------------------------------------------------------------
// Assignment Rules
// ---------------------------------------------------------------------------

func (s *Service) CreateAssignmentRule(ctx context.Context, companyID, userID string, req CreateAssignmentRuleRequest) (*AssignmentRule, error) {
	ar := &AssignmentRule{
		ID:             uuid.New().String(),
		CompanyID:      companyID,
		Name:           req.Name,
		CriteriaField:  req.CriteriaField,
		CriteriaValue:  req.CriteriaValue,
		CourseID:       req.CourseID,
		AssignmentType: req.AssignmentType,
		Active:         true,
		CreatedBy:      userID,
	}
	if err := s.repo.CreateAssignmentRule(ctx, ar); err != nil {
		return nil, svcErr("CreateAssignmentRule", err)
	}
	return ar, nil
}

func (s *Service) ListAssignmentRules(ctx context.Context, companyID string) ([]AssignmentRule, error) {
	return s.repo.ListAssignmentRules(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Assessments
// ---------------------------------------------------------------------------

func (s *Service) CreateAssessment(ctx context.Context, companyID, userID string, req CreateAssessmentRequest) (*Assessment, error) {
	a := &Assessment{
		ID:                uuid.New().String(),
		CompanyID:         companyID,
		Title:             req.Title,
		Description:       req.Description,
		AssessmentType:    req.AssessmentType,
		PassingScore:      req.PassingScore,
		TimeLimitMinutes:  req.TimeLimitMinutes,
		Status:            "draft",
		CreatedBy:         &userID,
	}
	if req.AttemptsAllowed != nil {
		a.AttemptsAllowed = *req.AttemptsAllowed
	}
	if req.RandomizeQuestions != nil {
		a.RandomizeQuestions = *req.RandomizeQuestions
	}
	if req.ShowResults != nil {
		a.ShowResults = *req.ShowResults
	}
	if err := s.repo.CreateAssessment(ctx, a); err != nil {
		return nil, svcErr("CreateAssessment", err)
	}
	return a, nil
}

func (s *Service) UpdateAssessment(ctx context.Context, companyID, id string, req UpdateAssessmentRequest) (*Assessment, error) {
	a, err := s.repo.GetAssessment(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("UpdateAssessment", err)
	}
	if req.Title != nil {
		a.Title = *req.Title
	}
	if req.Description != nil {
		a.Description = req.Description
	}
	if req.AttemptsAllowed != nil {
		a.AttemptsAllowed = *req.AttemptsAllowed
	}
	if req.PassingScore != nil {
		a.PassingScore = req.PassingScore
	}
	if req.TimeLimitMinutes != nil {
		a.TimeLimitMinutes = req.TimeLimitMinutes
	}
	if req.RandomizeQuestions != nil {
		a.RandomizeQuestions = *req.RandomizeQuestions
	}
	if req.ShowResults != nil {
		a.ShowResults = *req.ShowResults
	}
	if req.Status != nil {
		a.Status = *req.Status
	}
	if err := s.repo.UpdateAssessment(ctx, a); err != nil {
		return nil, svcErr("UpdateAssessment", err)
	}
	return a, nil
}

func (s *Service) GetAssessment(ctx context.Context, companyID, id string) (*Assessment, error) {
	return s.repo.GetAssessment(ctx, companyID, id)
}

func (s *Service) ListAssessments(ctx context.Context, companyID, courseID string) ([]Assessment, error) {
	return s.repo.ListAssessments(ctx, companyID, courseID)
}

// ---------------------------------------------------------------------------
// Questions
// ---------------------------------------------------------------------------

func (s *Service) AddQuestion(ctx context.Context, assessmentID string, req CreateQuestionRequest) (*Question, error) {
	// Get count for sort_order
	existing, err := s.repo.ListQuestions(ctx, assessmentID)
	if err != nil {
		return nil, svcErr("AddQuestion", err)
	}
	q := &Question{
		ID:           uuid.New().String(),
		AssessmentID: assessmentID,
		Question:     req.Question,
		QuestionType: req.QuestionType,
		Points:       1.0,
		SortOrder:    len(existing) + 1,
	}
	if req.Points != nil {
		q.Points = *req.Points
	}
	if err := s.repo.CreateQuestion(ctx, q); err != nil {
		return nil, svcErr("AddQuestion", err)
	}
	for i, opt := range req.Options {
		o := &QuestionOption{
			ID:         uuid.New().String(),
			QuestionID: q.ID,
			OptionText: opt.OptionText,
			SortOrder:  i + 1,
		}
		if opt.IsCorrect != nil {
			o.IsCorrect = *opt.IsCorrect
		}
		if err := s.repo.CreateOption(ctx, o); err != nil {
			return nil, svcErr("AddQuestion", err)
		}
	}
	return q, nil
}

func (s *Service) GetQuestions(ctx context.Context, assessmentID string) ([]Question, error) {
	qs, err := s.repo.ListQuestions(ctx, assessmentID)
	if err != nil {
		return nil, svcErr("GetQuestions", err)
	}
	for i := range qs {
		opts, err := s.repo.ListOptionsByQuestion(ctx, qs[i].ID)
		if err != nil {
			return nil, svcErr("GetQuestions", err)
		}
		qs[i].Options = opts
	}
	return qs, nil
}

// ---------------------------------------------------------------------------
// Attempts
// ---------------------------------------------------------------------------

func (s *Service) StartAttempt(ctx context.Context, enrollmentID, assessmentID string) (*Attempt, error) {
	// Check existing attempts
	existing, err := s.repo.ListAttemptsByEnrollment(ctx, enrollmentID)
	if err != nil {
		return nil, svcErr("StartAttempt", err)
	}
	assessment, err := s.repo.GetAssessment(ctx, "", assessmentID)
	if err != nil {
		return nil, svcErr("StartAttempt", err)
	}
	if assessment.AttemptsAllowed > 0 && len(existing) >= assessment.AttemptsAllowed {
		return nil, fmt.Errorf("maximum attempts reached")
	}
	a := &Attempt{
		ID:             uuid.New().String(),
		EnrollmentID:   enrollmentID,
		AssessmentID:   assessmentID,
		AttemptNumber:  len(existing) + 1,
		Status:         "in_progress",
		StartedAt:      strPtr(time.Now().Format(time.RFC3339)),
	}
	if err := s.repo.CreateAttempt(ctx, a); err != nil {
		return nil, svcErr("StartAttempt", err)
	}
	return a, nil
}

func (s *Service) SubmitAttempt(ctx context.Context, attemptID string, req SubmitAttemptRequest) (*AttemptResult, error) {
	attempt, err := s.repo.GetAttempt(ctx, attemptID)
	if err != nil {
		return nil, svcErr("SubmitAttempt", err)
	}
	assessment, err := s.repo.GetAssessment(ctx, "", attempt.AssessmentID)
	if err != nil {
		return nil, svcErr("SubmitAttempt", err)
	}
	questions, err := s.repo.ListQuestions(ctx, attempt.AssessmentID)
	if err != nil {
		return nil, svcErr("SubmitAttempt", err)
	}
	totalPoints := 0.0
	earnedPoints := 0.0
	for _, q := range questions {
		totalPoints += q.Points
	}
	answerMap := make(map[string]AnswerRequest)
	for _, ans := range req.Answers {
		answerMap[ans.QuestionID] = ans
	}
	for _, q := range questions {
		ans, ok := answerMap[q.ID]
		if !ok {
			continue
		}
		isCorrect := false
		pointsEarned := 0.0
		if q.QuestionType == "multiple_choice" || q.QuestionType == "single_choice" {
			if ans.SelectedOptionID != nil {
				opts, err := s.repo.ListOptionsByQuestion(ctx, q.ID)
				if err != nil {
					return nil, svcErr("SubmitAttempt", err)
				}
				for _, opt := range opts {
					if opt.ID == *ans.SelectedOptionID && opt.IsCorrect {
						isCorrect = true
						break
					}
				}
			}
		}
		if isCorrect {
			pointsEarned = q.Points
		}
		earnedPoints += pointsEarned
		answer := &Answer{
			ID:              uuid.New().String(),
			AttemptID:       attemptID,
			QuestionID:      q.ID,
			SelectedOptionID: ans.SelectedOptionID,
			TextAnswer:      ans.TextAnswer,
			NumericAnswer:   ans.NumericAnswer,
			IsCorrect:       &isCorrect,
			PointsEarned:    &pointsEarned,
		}
		if err := s.repo.CreateAnswer(ctx, answer); err != nil {
			return nil, svcErr("SubmitAttempt", err)
		}
	}
	// Calculate score
	score := 0.0
	if totalPoints > 0 {
		score = math.Round(earnedPoints/totalPoints*10000) / 100
	}
	status := "completed"
	passed := false
	if assessment.PassingScore != nil {
		if score >= *assessment.PassingScore {
			passed = true
		}
	}
	if passed {
		status = "passed"
	} else {
		status = "failed"
	}
	if err := s.repo.SubmitAttempt(ctx, attemptID, score, totalPoints, status); err != nil {
		return nil, svcErr("SubmitAttempt", err)
	}
	// If passed, auto-complete enrollment
	if passed {
		enrollmentID := attempt.EnrollmentID
		// Get enrollment company_id
		e, err := s.repo.GetEnrollment(ctx, "", enrollmentID)
		if err == nil {
			_ = s.CompleteEnrollment(context.Background(), e.CompanyID, enrollmentID)
		}
	}
	// Update attempt model for return
	attempt.Score = &score
	attempt.TotalPoints = &totalPoints
	attempt.Status = status
	return &AttemptResult{Attempt: *attempt, Passed: passed}, nil
}

func (s *Service) GetAttempt(ctx context.Context, id string) (*Attempt, error) {
	return s.repo.GetAttempt(ctx, id)
}

func (s *Service) ListAttempts(ctx context.Context, enrollmentID string) ([]Attempt, error) {
	return s.repo.ListAttemptsByEnrollment(ctx, enrollmentID)
}

// ---------------------------------------------------------------------------
// Instructors
// ---------------------------------------------------------------------------

func (s *Service) CreateInstructor(ctx context.Context, companyID, userID string, req CreateInstructorRequest) (*Instructor, error) {
	i := &Instructor{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		EmployeeID:      req.EmployeeID,
		InstructorType:  req.InstructorType,
		Name:            req.Name,
		Email:           req.Email,
		Phone:           req.Phone,
		Specialization:  req.Specialization,
		Bio:             req.Bio,
		Active:          true,
		CreatedBy:       &userID,
	}
	if err := s.repo.CreateInstructor(ctx, i); err != nil {
		return nil, svcErr("CreateInstructor", err)
	}
	return i, nil
}

func (s *Service) GetInstructor(ctx context.Context, companyID, id string) (*Instructor, error) {
	return s.repo.GetInstructor(ctx, companyID, id)
}

func (s *Service) ListInstructors(ctx context.Context, companyID string) ([]Instructor, error) {
	return s.repo.ListInstructors(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Providers
// ---------------------------------------------------------------------------

func (s *Service) CreateProvider(ctx context.Context, companyID, userID string, req CreateProviderRequest) (*TrainingProvider, error) {
	p := &TrainingProvider{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		Name:        req.Name,
		TaxID:       req.TaxID,
		Email:       req.Email,
		Phone:       req.Phone,
		Website:     req.Website,
		ContactName: req.ContactName,
		Notes:       req.Notes,
		Active:      true,
		CreatedBy:   &userID,
	}
	if err := s.repo.CreateProvider(ctx, p); err != nil {
		return nil, svcErr("CreateProvider", err)
	}
	return p, nil
}

func (s *Service) ListProviders(ctx context.Context, companyID string) ([]TrainingProvider, error) {
	return s.repo.ListProviders(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Competencies
// ---------------------------------------------------------------------------

func (s *Service) CreateCompetency(ctx context.Context, companyID, userID string, req CreateCompetencyRequest) (*Competency, error) {
	c := &Competency{
		ID:              uuid.New().String(),
		CompanyID:       companyID,
		Name:            req.Name,
		Description:     req.Description,
		CompetencyType:  req.CompetencyType,
		CreatedBy:       &userID,
	}
	if err := s.repo.CreateCompetency(ctx, c); err != nil {
		return nil, svcErr("CreateCompetency", err)
	}
	for _, l := range req.Levels {
		cl := &CompetencyLevel{
			ID:            uuid.New().String(),
			CompetencyID:  c.ID,
			Level:         l.Level,
			Label:         l.Label,
			Description:   l.Description,
			CreatedBy:     &userID,
		}
		if err := s.repo.CreateCompetencyLevel(ctx, cl); err != nil {
			return nil, svcErr("CreateCompetency", err)
		}
	}
	return c, nil
}

func (s *Service) GetCompetency(ctx context.Context, companyID, id string) (*Competency, error) {
	c, err := s.repo.GetCompetency(ctx, companyID, id)
	if err != nil {
		return nil, svcErr("GetCompetency", err)
	}
	levels, err := s.repo.ListCompetencyLevels(ctx, id)
	if err != nil {
		return nil, svcErr("GetCompetency", err)
	}
	c.Levels = levels
	return c, nil
}

func (s *Service) ListCompetencies(ctx context.Context, companyID string) ([]Competency, error) {
	return s.repo.ListCompetencies(ctx, companyID)
}

func (s *Service) AssignCompetency(ctx context.Context, companyID, employeeID, competencyID string, req AssignCompetencyRequest) error {
	ec := &EmployeeCompetency{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		EmployeeID:    employeeID,
		CompetencyID:  competencyID,
		Level:         req.Level,
		Source:        req.Source,
		VerifiedBy:    req.VerifiedBy,
	}
	if req.Verified != nil {
		ec.Verified = *req.Verified
	}
	if err := s.repo.UpsertEmployeeCompetency(ctx, ec); err != nil {
		return svcErr("AssignCompetency", err)
	}
	return nil
}

func (s *Service) GetEmployeeCompetencies(ctx context.Context, companyID, employeeID string) ([]EmployeeCompetency, error) {
	return s.repo.ListEmployeeCompetencies(ctx, companyID, employeeID)
}

// ---------------------------------------------------------------------------
// Course-Competency mapping
// ---------------------------------------------------------------------------

func (s *Service) AddCourseCompetency(ctx context.Context, courseID, competencyID string, expectedLevel, weight int) error {
	cc := &CourseCompetency{
		ID:             uuid.New().String(),
		CourseID:       courseID,
		CompetencyID:   competencyID,
		ExpectedLevel:  expectedLevel,
		Weight:         float64(weight),
	}
	if err := s.repo.AddCourseCompetency(ctx, cc); err != nil {
		return svcErr("AddCourseCompetency", err)
	}
	return nil
}

func (s *Service) ListCourseCompetencies(ctx context.Context, courseID string) ([]CourseCompetency, error) {
	return s.repo.ListCourseCompetencies(ctx, courseID)
}

// ---------------------------------------------------------------------------
// Training Needs
// ---------------------------------------------------------------------------

func (s *Service) CreateTrainingNeed(ctx context.Context, companyID, userID string, req CreateTrainingNeedRequest) (*TrainingNeed, error) {
	n := &TrainingNeed{
		ID:            uuid.New().String(),
		CompanyID:     companyID,
		EmployeeID:    req.EmployeeID,
		CompetencyID:  req.CompetencyID,
		Title:         req.Title,
		Description:   req.Description,
		Priority:      req.Priority,
		Source:        req.Source,
		SourceID:      req.SourceID,
		Status:        "open",
		CreatedBy:     &userID,
	}
	if err := s.repo.CreateTrainingNeed(ctx, n); err != nil {
		return nil, svcErr("CreateTrainingNeed", err)
	}
	return n, nil
}

func (s *Service) ListTrainingNeeds(ctx context.Context, companyID string) ([]TrainingNeed, error) {
	return s.repo.ListTrainingNeeds(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Training Plans
// ---------------------------------------------------------------------------

func (s *Service) CreatePlan(ctx context.Context, companyID, userID string, req CreateTrainingPlanRequest) (*TrainingPlan, error) {
	p := &TrainingPlan{
		ID:           uuid.New().String(),
		CompanyID:    companyID,
		EmployeeID:   req.EmployeeID,
		Name:         req.Name,
		Description:  req.Description,
		Objectives:   req.Objectives,
		PeriodStart:  parseTimePtr(req.PeriodStart),
		PeriodEnd:    parseTimePtr(req.PeriodEnd),
		BudgetAmount: req.BudgetAmount,
		Status:       "active",
		CreatedBy:    userID,
	}
	if req.BudgetCurrency != nil {
		p.BudgetCurrency = *req.BudgetCurrency
	}
	if err := s.repo.CreatePlan(ctx, p); err != nil {
		return nil, svcErr("CreatePlan", err)
	}
	// Auto-assign courses
	for _, courseID := range req.CourseIDs {
		a := &Assignment{
			ID:             uuid.New().String(),
			CompanyID:      companyID,
			CourseID:       courseID,
			AssigneeType:   "employee",
			AssigneeID:     req.EmployeeID,
			AssignmentType: "planned",
			CreatedBy:      userID,
		}
		if err := s.repo.CreateAssignment(ctx, a); err != nil {
			return nil, svcErr("CreatePlan", err)
		}
	}
	return p, nil
}

func (s *Service) ListPlans(ctx context.Context, companyID string) ([]TrainingPlan, error) {
	return s.repo.ListPlans(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Learning Paths
// ---------------------------------------------------------------------------

func (s *Service) CreateLearningPath(ctx context.Context, companyID, userID string, req CreateLearningPathRequest) (*LearningPath, error) {
	lp := &LearningPath{
		ID:           uuid.New().String(),
		CompanyID:    companyID,
		Name:         req.Name,
		Description:  req.Description,
		Objectives:   req.Objectives,
		DurationDays: req.DurationDays,
		Status:       "active",
		CreatedBy:    userID,
	}
	if err := s.repo.CreateLearningPath(ctx, lp); err != nil {
		return nil, svcErr("CreateLearningPath", err)
	}
	for i, courseID := range req.CourseIDs {
		if err := s.repo.AddPathCourse(ctx, lp.ID, courseID, i+1); err != nil {
			return nil, svcErr("CreateLearningPath", err)
		}
	}
	return lp, nil
}

func (s *Service) ListLearningPaths(ctx context.Context, companyID string) ([]LearningPath, error) {
	return s.repo.ListLearningPaths(ctx, companyID)
}

// ---------------------------------------------------------------------------
// Feedback
// ---------------------------------------------------------------------------

func (s *Service) CreateFeedback(ctx context.Context, enrollmentID, userID string, req CreateFeedbackRequest) (*Feedback, error) {
	f := &Feedback{
		ID:                 uuid.New().String(),
		EnrollmentID:       enrollmentID,
		InstructorRating:   req.InstructorRating,
		ContentRating:      req.ContentRating,
		OrganizationRating: req.OrganizationRating,
		PlatformRating:     req.PlatformRating,
		OverallRating:      req.OverallRating,
		Comments:           req.Comments,
		CreatedBy:          userID,
	}
	if err := s.repo.CreateFeedback(ctx, f); err != nil {
		return nil, svcErr("CreateFeedback", err)
	}
	return f, nil
}

func (s *Service) GetFeedbackByEnrollment(ctx context.Context, enrollmentID string) (*Feedback, error) {
	return s.repo.GetFeedbackByEnrollment(ctx, enrollmentID)
}

// ---------------------------------------------------------------------------
// Attendance
// ---------------------------------------------------------------------------

func (s *Service) CreateAttendance(ctx context.Context, enrollmentID, sessionID, userID string, req CreateAttendanceRequest) (*Attendance, error) {
	a := &Attendance{
		ID:           uuid.New().String(),
		EnrollmentID: enrollmentID,
		SessionID:    sessionID,
		Status:       req.Status,
		CheckIn:      req.CheckIn,
		CheckOut:     req.CheckOut,
		Notes:        req.Notes,
		CreatedBy:    userID,
	}
	if err := s.repo.CreateAttendance(ctx, a); err != nil {
		return nil, svcErr("CreateAttendance", err)
	}
	return a, nil
}

func (s *Service) GetAttendance(ctx context.Context, enrollmentID, sessionID string) (*Attendance, error) {
	return s.repo.GetAttendanceBySession(ctx, enrollmentID, sessionID)
}

func (s *Service) ListAttendance(ctx context.Context, enrollmentID string) ([]Attendance, error) {
	return s.repo.ListAttendanceByEnrollment(ctx, enrollmentID)
}

// ---------------------------------------------------------------------------
// Certificates
// ---------------------------------------------------------------------------

func (s *Service) ListCertificates(ctx context.Context, companyID, employeeID string) ([]Enrollment, error) {
	return s.repo.ListCertificatesByEmployee(ctx, companyID, employeeID)
}

// ---------------------------------------------------------------------------
// Dashboard & Stats
// ---------------------------------------------------------------------------

func (s *Service) GetDashboardStats(ctx context.Context, companyID string) (*DashboardStats, error) {
	return s.repo.GetDashboardStats(ctx, companyID)
}

func (s *Service) GetEmployeeStats(ctx context.Context, companyID, employeeID string) (*EmployeeStats, error) {
	return s.repo.GetEmployeeStats(ctx, companyID, employeeID)
}

// ---------------------------------------------------------------------------
// AI Recommendations
// ---------------------------------------------------------------------------

func (s *Service) GenerateRecommendations(ctx context.Context, companyID string, req AIRecommendationRequest) ([]AIRecommendation, error) {
	// Placeholder: rule-based recommendations
	// In production, this would call an AI service
	needs, err := s.repo.ListTrainingNeeds(ctx, companyID)
	if err != nil {
		return nil, svcErr("GenerateRecommendations", err)
	}
	var recs []AIRecommendation
	for _, need := range needs {
		if need.EmployeeID != nil && *need.EmployeeID == req.EmployeeID {
			// Find courses that map to this competency
			competencies, err := s.repo.ListCourseCompetencies(ctx, "")
			if err != nil {
				continue
			}
			for _, cc := range competencies {
				if need.CompetencyID != nil && cc.CompetencyID == *need.CompetencyID {
					c, err := s.repo.GetCourse(ctx, companyID, cc.CourseID)
					if err != nil {
						continue
					}
					recs = append(recs, AIRecommendation{
						CourseID:   c.ID,
						CourseName: c.Name,
						Reason:     need.Title,
						Priority:   need.Priority,
					})
				}
			}
		}
	}
	// Also recommend based on competency gaps
	empComps, err := s.repo.ListEmployeeCompetencies(ctx, companyID, req.EmployeeID)
	if err == nil {
		empCompMap := make(map[string]int)
		for _, ec := range empComps {
			empCompMap[ec.CompetencyID] = ec.Level
		}
		allComps, _ := s.repo.ListCompetencies(ctx, companyID)
		for _, comp := range allComps {
			currentLevel, exists := empCompMap[comp.ID]
			if !exists || currentLevel < 3 {
				courses, _ := s.repo.ListCourseCompetencies(ctx, "")
				for _, cc := range courses {
					if cc.CompetencyID == comp.ID && cc.ExpectedLevel > currentLevel {
						c, err := s.repo.GetCourse(ctx, companyID, cc.CourseID)
						if err != nil {
							continue
						}
						recs = append(recs, AIRecommendation{
							CourseID:       c.ID,
							CourseName:     c.Name,
							Reason:         fmt.Sprintf("Gap in competency: %s", comp.Name),
							CompetencyID:   comp.ID,
							CompetencyName: comp.Name,
							ExpectedLevel:  cc.ExpectedLevel,
							Priority:       "medium",
						})
					}
				}
			}
		}
	}
	if len(recs) > 0 {
		_ = s.repo.SaveRecommendations(ctx, req.EmployeeID, recs)
	}
	return recs, nil
}

func (s *Service) GetRecommendations(ctx context.Context, employeeID string) ([]AIRecommendation, error) {
	return s.repo.GetRecommendations(ctx, employeeID)
}

// ---------------------------------------------------------------------------
// Events
// ---------------------------------------------------------------------------

func (s *Service) createEvent(ctx context.Context, companyID, eventType, title, description, userID, employeeID, offeringID, severity string, scheduledFor *string) {
	evt := &TrainingEvent{
		ID:          uuid.New().String(),
		CompanyID:   companyID,
		EventType:   eventType,
		Title:       title,
		Description: description,
		EmployeeID:  employeeID,
		OfferingID:  offeringID,
		Severity:    severity,
		ScheduledFor: scheduledFor,
		CreatedBy:   userID,
	}
	if err := s.repo.CreateEvent(ctx, evt); err != nil {
		s.log.Warn("failed to create training event", zap.Error(err))
	}
}

func (s *Service) ProcessPendingEvents(ctx context.Context, companyID string) error {
	events, err := s.repo.ListPendingEvents(ctx, companyID, 50)
	if err != nil {
		return svcErr("ProcessPendingEvents", err)
	}
	for _, e := range events {
		// Process based on type
		switch e.EventType {
		case "enrollment_reminder":
			if e.EmployeeID != "" {
				s.notif.Send(ctx, e.CompanyID, e.EmployeeID, "Training Reminder", e.Title, e.Description)
			}
		case "course_published":
			// Notify relevant employees
			s.log.Info("course published event", zap.String("title", e.Title))
		case "certificate_expiring":
			if e.EmployeeID != "" {
				s.notif.Send(ctx, e.CompanyID, e.EmployeeID, "Certificate Expiring", e.Title, e.Description)
			}
		}
		_ = s.repo.MarkEventProcessed(ctx, e.ID)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Worker helpers
// ---------------------------------------------------------------------------

func (s *Service) CheckOverdueEnrollments(ctx context.Context, companyID string) {
	// Find enrollments in "in_progress" past due date
	offerings, _, err := s.repo.ListOfferings(ctx, companyID, OfferingFilter{Status: strPtr("active")})
	if err != nil {
		s.log.Error("check overdue enrollments", zap.Error(err))
		return
	}
	for _, o := range offerings {
		if o.EndDate != nil {
			if time.Now().After(*o.EndDate) {
				enrollments, _, err := s.repo.ListEnrollments(ctx, companyID, EnrollmentFilter{OfferingID: &o.ID, Status: strPtr("in_progress")})
				if err != nil {
					continue
				}
				for _, e := range enrollments {
					eventUserID := ""
					if e.CreatedBy != nil {
						eventUserID = *e.CreatedBy
					}
					s.createEvent(ctx, companyID, "enrollment_overdue", "Enrollment Overdue",
						fmt.Sprintf("Enrollment %s is past the end date", e.ID),
						eventUserID, e.EmployeeID, o.ID, "warning", nil)
				}
			}
		}
	}
}

func (s *Service) CheckCertificateExpirations(ctx context.Context, companyID string) {
	// Find certificates that are about to expire (no real expiry in schema, placeholder)
	s.log.Info("certificate expiration check", zap.String("company", companyID))
}

func strPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func parseTimePtr(s *string) *time.Time {
	if s == nil {
		return nil
	}
	t, err := time.Parse(time.RFC3339, *s)
	if err != nil {
		return nil
	}
	return &t
}

// ---------------------------------------------------------------------------
// Result types
// ---------------------------------------------------------------------------

type AttemptResult struct {
	Attempt Attempt `json:"attempt"`
	Passed  bool    `json:"passed"`
}

// ---------------------------------------------------------------------------
// Filter types
// ---------------------------------------------------------------------------

type CourseFilter struct {
	CategoryID *string `json:"category_id"`
	Status     *string `json:"status"`
	Modality   *string `json:"modality"`
	Mandatory  *bool   `json:"mandatory"`
	Difficulty *string `json:"difficulty"`
	Search     *string `json:"search"`
	SortBy     string  `json:"sort_by"`
	SortDesc   bool    `json:"sort_desc"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type OfferingFilter struct {
	CourseID     *string `json:"course_id"`
	Status       *string `json:"status"`
	InstructorID *string `json:"instructor_id"`
	FromDate     *string `json:"from_date"`
	ToDate       *string `json:"to_date"`
	Limit        int     `json:"limit"`
	Offset       int     `json:"offset"`
}

type EnrollmentFilter struct {
	OfferingID *string `json:"offering_id"`
	EmployeeID *string `json:"employee_id"`
	Status     *string `json:"status"`
	Limit      int     `json:"limit"`
	Offset     int     `json:"offset"`
}

type DashboardStats struct {
	TotalCourses           int     `json:"total_courses"`
	TotalEnrollments       int     `json:"total_enrollments"`
	CompletedEnrollments   int     `json:"completed_enrollments"`
	InProgressEnrollments  int     `json:"in_progress_enrollments"`
	ActiveOfferings        int     `json:"active_offerings"`
	AverageRating          float64 `json:"average_rating"`
}

type EmployeeStats struct {
	TotalEnrollments  int `json:"total_enrollments"`
	CompletedCourses  int `json:"completed_courses"`
	InProgress        int `json:"in_progress"`
	TotalTrainingHours int `json:"total_training_hours"`
	CertificatesCount  int `json:"certificates_count"`
}

type TrainingEvent struct {
	ID           string  `json:"id"`
	CompanyID    string  `json:"company_id"`
	EventType    string  `json:"event_type"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	EmployeeID   string  `json:"employee_id"`
	EnrollmentID string  `json:"enrollment_id"`
	OfferingID   string  `json:"offering_id"`
	Severity     string  `json:"severity"`
	ScheduledFor *string `json:"scheduled_for"`
	ProcessedAt  *string `json:"processed_at"`
	Metadata     *string `json:"metadata"`
	CreatedBy    string  `json:"created_by"`
	CreatedAt    string  `json:"created_at"`
}
