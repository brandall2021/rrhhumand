package onboarding

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Templates

func (s *Service) CreateTemplate(ctx context.Context, companyID, userID string, req *CreateTemplateRequest) (*OnboardingTemplate, error) {
	t := &OnboardingTemplate{
		CompanyID: companyID,
		Name:     req.Name,
		Description: req.Description,
		CreatedBy:  userID,
	}
	if req.DefaultDurationDays != nil {
		t.DefaultDurationDays = *req.DefaultDurationDays
	}
	if err := s.repo.CreateTemplate(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTemplate(ctx context.Context, companyID, id string) (*OnboardingTemplate, error) {
	return s.repo.GetTemplate(ctx, companyID, id)
}

func (s *Service) ListTemplates(ctx context.Context, companyID string) ([]OnboardingTemplate, error) {
	return s.repo.ListTemplates(ctx, companyID)
}

func (s *Service) UpdateTemplate(ctx context.Context, companyID, id string, req *UpdateTemplateRequest) (*OnboardingTemplate, error) {
	return s.repo.UpdateTemplate(ctx, companyID, id, req)
}

func (s *Service) DeleteTemplate(ctx context.Context, companyID, id string) error {
	return s.repo.DeleteTemplate(ctx, companyID, id)
}

// Template Tasks

func (s *Service) CreateTemplateTask(ctx context.Context, templateID string, req *CreateTemplateTaskRequest) (*OnboardingTemplateTask, error) {
	required := true
	if req.Required != nil {
		required = *req.Required
	}
	sortOrder := req.DaysOffset
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}
	t := &OnboardingTemplateTask{
		TemplateID:        templateID,
		Title:             req.Title,
		Description:       req.Description,
		Category:          req.Category,
		ResponsibleType:   req.ResponsibleType,
		ResponsibleUserID: req.ResponsibleUserID,
		Required:          required,
		DaysOffset:        req.DaysOffset,
		SortOrder:         sortOrder,
		EstimatedMinutes:  req.EstimatedMinutes,
	}
	if err := s.repo.CreateTemplateTask(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) ListTemplateTasks(ctx context.Context, templateID string) ([]OnboardingTemplateTask, error) {
	return s.repo.ListTemplateTasks(ctx, templateID)
}

func (s *Service) UpdateTemplateTask(ctx context.Context, id string, req *UpdateTemplateTaskRequest) (*OnboardingTemplateTask, error) {
	return s.repo.UpdateTemplateTask(ctx, id, req)
}

func (s *Service) DeleteTemplateTask(ctx context.Context, id string) error {
	return s.repo.DeleteTemplateTask(ctx, id)
}

// Processes

func (s *Service) CreateOnboarding(ctx context.Context, companyID, userID string, req *CreateOnboardingRequest) (*OnboardingProcess, error) {
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date: %w", err)
	}

	durationDays := 90
	if req.TemplateID != nil {
		t, err := s.repo.GetTemplate(ctx, companyID, *req.TemplateID)
		if err == nil {
			durationDays = t.DefaultDurationDays
		}
	}

	policy := "STRICT"
	if req.CompletionPolicy != nil {
		policy = *req.CompletionPolicy
	}

	p := &OnboardingProcess{
		CompanyID:            companyID,
		EmployeeID:           req.EmployeeID,
		TemplateID:           req.TemplateID,
		StartDate:            startDate,
		TargetCompletionDate: startDate.AddDate(0, 0, durationDays),
		CompletionPolicy:     policy,
		CreatedBy:            userID,
	}

	if err := s.repo.CreateProcess(ctx, p); err != nil {
		return nil, err
	}

	if req.TemplateID != nil {
		if err := s.generateTasksFromTemplate(ctx, companyID, p.ID, req.EmployeeID, p.StartDate, *req.TemplateID); err != nil {
			return nil, fmt.Errorf("failed to generate tasks from template: %w", err)
		}
	}

	s.emitEvent(ctx, companyID, p.ID, "onboarding", "onboarding.created", nil)
	s.createNotification(ctx, companyID, userID, "Onboarding creado",
		fmt.Sprintf("Se ha creado el proceso de onboarding para el empleado %s", req.EmployeeID),
		"ONBOARDING_CREATED", "onboarding", p.ID)

	return p, nil
}

func (s *Service) generateTasksFromTemplate(ctx context.Context, companyID, onboardingID, employeeID string, startDate time.Time, templateID string) error {
	tasks, err := s.repo.ListTemplateTasks(ctx, templateID)
	if err != nil {
		return err
	}

	for _, tt := range tasks {
		dueDate := startDate.AddDate(0, 0, tt.DaysOffset)
		task := &OnboardingTask{
			OnboardingID:    onboardingID,
			CompanyID:       companyID,
			EmployeeID:      employeeID,
			Title:           tt.Title,
			Description:     tt.Description,
			Category:        tt.Category,
			ResponsibleType: tt.ResponsibleType,
			ResponsibleID:   tt.ResponsibleUserID,
			DueDate:         dueDate,
			Required:        tt.Required,
			SortOrder:       tt.SortOrder,
			EstimatedMinutes: tt.EstimatedMinutes,
		}
		if err := s.repo.CreateTask(ctx, task); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetProcess(ctx context.Context, companyID, id string) (*OnboardingProcess, error) {
	return s.repo.GetProcess(ctx, companyID, id)
}

func (s *Service) GetProcessByEmployee(ctx context.Context, companyID, employeeID string) (*OnboardingProcess, error) {
	return s.repo.GetProcessByEmployee(ctx, companyID, employeeID)
}

func (s *Service) ListProcesses(ctx context.Context, companyID string, filters OnboardingFilters) ([]OnboardingProcess, error) {
	return s.repo.ListProcesses(ctx, companyID, filters)
}

func (s *Service) UpdateProcess(ctx context.Context, companyID, id string, req *UpdateOnboardingRequest) (*OnboardingProcess, error) {
	if err := s.repo.UpdateProcess(ctx, companyID, id, req); err != nil {
		return nil, err
	}
	return s.repo.GetProcess(ctx, companyID, id)
}

func (s *Service) StartOnboarding(ctx context.Context, companyID, id string) error {
	p, err := s.repo.GetProcess(ctx, companyID, id)
	if err != nil {
		return err
	}
	if p.Status != "NOT_STARTED" {
		return fmt.Errorf("only NOT_STARTED onboarding can be started, current status: %s", p.Status)
	}
	if err := s.repo.UpdateProcessStatus(ctx, companyID, id, "IN_PROGRESS"); err != nil {
		return err
	}
	s.emitEvent(ctx, companyID, id, "onboarding", "onboarding.started", nil)
	return nil
}

func (s *Service) HoldOnboarding(ctx context.Context, companyID, id string) error {
	p, err := s.repo.GetProcess(ctx, companyID, id)
	if err != nil {
		return err
	}
	if p.Status != "IN_PROGRESS" {
		return fmt.Errorf("only IN_PROGRESS onboarding can be put on hold, current status: %s", p.Status)
	}
	return s.repo.UpdateProcessStatus(ctx, companyID, id, "ON_HOLD")
}

func (s *Service) ResumeOnboarding(ctx context.Context, companyID, id string) error {
	p, err := s.repo.GetProcess(ctx, companyID, id)
	if err != nil {
		return err
	}
	if p.Status != "ON_HOLD" {
		return fmt.Errorf("only ON_HOLD onboarding can be resumed, current status: %s", p.Status)
	}
	return s.repo.UpdateProcessStatus(ctx, companyID, id, "IN_PROGRESS")
}

func (s *Service) calculateProgress(ctx context.Context, onboardingID string) (int, error) {
	total, completed, err := s.repo.GetRequiredTaskCount(ctx, onboardingID)
	if err != nil {
		return 0, err
	}
	if total == 0 {
		return 0, nil
	}
	return completed * 100 / total, nil
}

func (s *Service) CompleteOnboarding(ctx context.Context, companyID, id string) error {
	p, err := s.repo.GetProcess(ctx, companyID, id)
	if err != nil {
		return err
	}
	if p.Status == "COMPLETED" {
		return fmt.Errorf("onboarding is already completed")
	}
	if p.Status == "CANCELLED" {
		return fmt.Errorf("onboarding is cancelled")
	}

	if p.CompletionPolicy == "STRICT" {
		tasks, err := s.repo.ListTasks(ctx, id)
		if err != nil {
			return err
		}
		for _, t := range tasks {
			if t.Required && t.Status != "COMPLETED" {
				return fmt.Errorf("required task not completed: %s", t.Title)
			}
		}

		docs, err := s.repo.ListDocuments(ctx, id)
		if err != nil {
			return err
		}
		for _, d := range docs {
			if d.Required && d.Status != "APPROVED" {
				return fmt.Errorf("required document not approved: %s (%s)", d.DocumentType, d.Status)
			}
		}
	} else if p.CompletionPolicy == "FLEXIBLE" {
		// Allow exceptions for non-completed items
	}

	if err := s.repo.CompleteProcess(ctx, companyID, id); err != nil {
		return err
	}

	s.emitEvent(ctx, companyID, id, "onboarding", "onboarding.completed", nil)
	s.createNotification(ctx, companyID, p.EmployeeID, "Onboarding completado",
		"¡Felicitaciones! Has completado tu proceso de onboarding.",
		"ONBOARDING_COMPLETED", "onboarding", id)

	return nil
}

func (s *Service) CancelOnboarding(ctx context.Context, companyID, id, reason string) error {
	p, err := s.repo.GetProcess(ctx, companyID, id)
	if err != nil {
		return err
	}
	if p.Status == "COMPLETED" {
		return fmt.Errorf("cannot cancel a completed onboarding")
	}
	return s.repo.CancelProcess(ctx, companyID, id, reason)
}

// Tasks

func (s *Service) CreateTask(ctx context.Context, companyID string, onboardingID string, req *CreateTaskRequest) (*OnboardingTask, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	dueDate, err := parseDate(req.DueDate)
	if err != nil {
		return nil, fmt.Errorf("invalid due_date: %w", err)
	}

	required := true
	if req.Required != nil {
		required = *req.Required
	}

	t := &OnboardingTask{
		OnboardingID:    onboardingID,
		CompanyID:       companyID,
		EmployeeID:      p.EmployeeID,
		Title:           req.Title,
		Description:     req.Description,
		Category:        req.Category,
		ResponsibleType: req.ResponsibleType,
		ResponsibleID:   req.ResponsibleID,
		DueDate:         dueDate,
		Required:        required,
		SortOrder:       req.SortOrder,
		EstimatedMinutes: req.EstimatedMinutes,
	}
	if err := s.repo.CreateTask(ctx, t); err != nil {
		return nil, err
	}

	s.updateProgress(ctx, companyID, onboardingID)
	return t, nil
}

func (s *Service) GetTask(ctx context.Context, companyID, taskID string) (*OnboardingTask, error) {
	return s.repo.GetTask(ctx, companyID, taskID)
}

func (s *Service) ListTasks(ctx context.Context, onboardingID string) ([]OnboardingTask, error) {
	return s.repo.ListTasks(ctx, onboardingID)
}

func (s *Service) UpdateTask(ctx context.Context, companyID, id string, req *UpdateTaskRequest) (*OnboardingTask, error) {
	if err := s.repo.UpdateTask(ctx, companyID, id, req); err != nil {
		return nil, err
	}
	return s.repo.GetTask(ctx, companyID, id)
}

func (s *Service) CompleteTask(ctx context.Context, companyID, id string) error {
	if err := s.repo.CompleteTask(ctx, companyID, id); err != nil {
		return err
	}
	t, err := s.repo.GetTask(ctx, companyID, id)
	if err != nil {
		return err
	}
	s.updateProgress(ctx, companyID, t.OnboardingID)
	s.emitEvent(ctx, companyID, id, "onboarding_task", "task.completed", nil)
	return nil
}

func (s *Service) StartTask(ctx context.Context, companyID, id string) error {
	return s.repo.StartTask(ctx, companyID, id)
}

func (s *Service) BlockTask(ctx context.Context, companyID, id, reason string) error {
	return s.repo.BlockTask(ctx, companyID, id, reason)
}

func (s *Service) updateProgress(ctx context.Context, companyID, onboardingID string) {
	progress, err := s.calculateProgress(ctx, onboardingID)
	if err != nil {
		return
	}
	s.repo.UpdateProcessProgress(ctx, companyID, onboardingID, progress)
}

// Documents

func (s *Service) CreateDocumentRequirement(ctx context.Context, companyID, onboardingID string, req *UploadDocumentRequest) (*OnboardingDocument, error) {
	required := true
	if req.Required != nil {
		required = *req.Required
	}

	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	d := &OnboardingDocument{
		OnboardingID: onboardingID,
		CompanyID:    companyID,
		EmployeeID:   p.EmployeeID,
		DocumentType: req.DocumentType,
		Required:     required,
	}
	if err := s.repo.CreateDocument(ctx, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) GetDocument(ctx context.Context, companyID, id string) (*OnboardingDocument, error) {
	return s.repo.GetDocument(ctx, companyID, id)
}

func (s *Service) ListDocuments(ctx context.Context, onboardingID string) ([]OnboardingDocument, error) {
	return s.repo.ListDocuments(ctx, onboardingID)
}

func (s *Service) ApproveDocument(ctx context.Context, companyID, id, reviewedBy string) error {
	if err := s.repo.ApproveDocument(ctx, companyID, id, reviewedBy); err != nil {
		return err
	}
	d, err := s.repo.GetDocument(ctx, companyID, id)
	if err != nil {
		return err
	}
	s.emitEvent(ctx, companyID, id, "onboarding_document", "document.approved", nil)
	s.updateProgress(ctx, companyID, d.OnboardingID)
	return nil
}

func (s *Service) RejectDocument(ctx context.Context, companyID, id, reviewedBy, reason string) error {
	return s.repo.RejectDocument(ctx, companyID, id, reviewedBy, reason)
}

// Assets

func (s *Service) CreateAsset(ctx context.Context, companyID, onboardingID string, req *CreateAssetRequest) (*OnboardingAsset, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	a := &OnboardingAsset{
		OnboardingID: onboardingID,
		CompanyID:    companyID,
		EmployeeID:   p.EmployeeID,
		AssetType:    req.AssetType,
		Description:  req.Description,
		SerialNumber: req.SerialNumber,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateAsset(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetAsset(ctx context.Context, companyID, id string) (*OnboardingAsset, error) {
	return s.repo.GetAsset(ctx, companyID, id)
}

func (s *Service) ListAssets(ctx context.Context, onboardingID string) ([]OnboardingAsset, error) {
	return s.repo.ListAssets(ctx, onboardingID)
}

func (s *Service) AssignAsset(ctx context.Context, companyID, id, assignedBy string) error {
	return s.repo.AssignAsset(ctx, companyID, id, assignedBy)
}

func (s *Service) DeliverAsset(ctx context.Context, companyID, id string) error {
	return s.repo.DeliverAsset(ctx, companyID, id)
}

func (s *Service) ReturnAsset(ctx context.Context, companyID, id string) error {
	return s.repo.ReturnAsset(ctx, companyID, id)
}

// Access Requests

func (s *Service) CreateAccessRequest(ctx context.Context, companyID, onboardingID string, req *CreateAccessRequest) (*AccessRequest, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	ar := &AccessRequest{
		OnboardingID: onboardingID,
		CompanyID:    companyID,
		EmployeeID:   p.EmployeeID,
		SystemName:   req.SystemName,
		AccessType:   req.AccessType,
		Notes:        req.Notes,
	}
	if err := s.repo.CreateAccessRequest(ctx, ar); err != nil {
		return nil, err
	}
	return ar, nil
}

func (s *Service) GetAccessRequest(ctx context.Context, companyID, id string) (*AccessRequest, error) {
	return s.repo.GetAccessRequest(ctx, companyID, id)
}

func (s *Service) ListAccessRequests(ctx context.Context, onboardingID string) ([]AccessRequest, error) {
	return s.repo.ListAccessRequests(ctx, onboardingID)
}

func (s *Service) ApproveAccess(ctx context.Context, companyID, id, approvedBy string) error {
	if err := s.repo.ApproveAccess(ctx, companyID, id, approvedBy); err != nil {
		return err
	}
	s.emitEvent(ctx, companyID, id, "access_request", "access.approved", nil)
	return nil
}

func (s *Service) RejectAccess(ctx context.Context, companyID, id string) error {
	return s.repo.RejectAccess(ctx, companyID, id)
}

func (s *Service) ActivateAccess(ctx context.Context, companyID, id string) error {
	return s.repo.ActivateAccess(ctx, companyID, id)
}

func (s *Service) RevokeAccess(ctx context.Context, companyID, id string) error {
	return s.repo.RevokeAccess(ctx, companyID, id)
}

// Milestones

func (s *Service) CreateMilestone(ctx context.Context, companyID, onboardingID string, req *CreateMilestoneRequest) (*OnboardingMilestone, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	daysOffset := req.DaysOffset
	if daysOffset == 0 {
		daysOffset = 30
	}
	responsibleType := req.ResponsibleType
	if responsibleType == "" {
		responsibleType = "MANAGER"
	}

	m := &OnboardingMilestone{
		OnboardingID:    onboardingID,
		CompanyID:       companyID,
		EmployeeID:      p.EmployeeID,
		MilestoneType:   req.MilestoneType,
		Title:           req.Title,
		Description:     req.Description,
		DaysOffset:      daysOffset,
		DueDate:         p.StartDate.AddDate(0, 0, daysOffset),
		ResponsibleType: responsibleType,
		ResponsibleID:   req.ResponsibleID,
	}
	if err := s.repo.CreateMilestone(ctx, m); err != nil {
		return nil, err
	}
	return m, nil
}

func (s *Service) ListMilestones(ctx context.Context, onboardingID string) ([]OnboardingMilestone, error) {
	return s.repo.ListMilestones(ctx, onboardingID)
}

func (s *Service) UpdateMilestone(ctx context.Context, companyID, id string, req *UpdateMilestoneRequest) (*OnboardingMilestone, error) {
	return s.repo.UpdateMilestone(ctx, companyID, id, req)
}

func (s *Service) CompleteMilestone(ctx context.Context, companyID, id string) error {
	if err := s.repo.CompleteMilestone(ctx, companyID, id); err != nil {
		return err
	}
	s.emitEvent(ctx, companyID, id, "onboarding_milestone", "milestone.completed", nil)
	return nil
}

// Feedback

func (s *Service) CreateFeedback(ctx context.Context, companyID, onboardingID, userID string, req *CreateFeedbackRequest) (*OnboardingFeedback, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	f := &OnboardingFeedback{
		OnboardingID:      onboardingID,
		CompanyID:         companyID,
		EmployeeID:        p.EmployeeID,
		FeedbackType:      req.FeedbackType,
		SubmittedBy:       userID,
		AdaptationScore:   req.AdaptationScore,
		TeamScore:         req.TeamScore,
		KnowledgeScore:    req.KnowledgeScore,
		CommunicationScore: req.CommunicationScore,
		OverallScore:      req.OverallScore,
		Comments:          req.Comments,
	}
	if err := s.repo.CreateFeedback(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

func (s *Service) ListFeedback(ctx context.Context, onboardingID string) ([]OnboardingFeedback, error) {
	return s.repo.ListFeedback(ctx, onboardingID)
}

// Buddies

func (s *Service) AssignBuddy(ctx context.Context, companyID, onboardingID, userID string, req *AssignBuddyRequest) (*OnboardingBuddy, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	b := &OnboardingBuddy{
		OnboardingID:    onboardingID,
		CompanyID:       companyID,
		EmployeeID:      p.EmployeeID,
		BuddyEmployeeID: req.BuddyEmployeeID,
		StartDate:       p.StartDate,
		Notes:           req.Notes,
	}
	if err := s.repo.AssignBuddy(ctx, b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetBuddy(ctx context.Context, onboardingID string) (*OnboardingBuddy, error) {
	return s.repo.GetBuddy(ctx, onboardingID)
}

// Exceptions

func (s *Service) CreateException(ctx context.Context, companyID, onboardingID, userID string, req *CreateExceptionRequest) (*OnboardingException, error) {
	e := &OnboardingException{
		OnboardingID: onboardingID,
		CompanyID:    companyID,
		EntityType:   req.EntityType,
		EntityID:     req.EntityID,
		Reason:       req.Reason,
		CreatedBy:    userID,
	}
	if err := s.repo.CreateException(ctx, e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) ListExceptions(ctx context.Context, onboardingID string) ([]OnboardingException, error) {
	return s.repo.ListExceptions(ctx, onboardingID)
}

// Training

func (s *Service) CreateTraining(ctx context.Context, companyID, onboardingID string, req *CreateTrainingAssignmentRequest) (*TrainingAssignment, error) {
	p, err := s.repo.GetProcess(ctx, companyID, onboardingID)
	if err != nil {
		return nil, err
	}

	trainingType := req.TrainingType
	if trainingType == "" {
		trainingType = "MANDATORY"
	}

	var dueDate *time.Time
	if req.DueDate != nil {
		d, err := parseDate(*req.DueDate)
		if err == nil {
			dueDate = &d
		}
	}

	t := &TrainingAssignment{
		OnboardingID:     onboardingID,
		CompanyID:        companyID,
		EmployeeID:       p.EmployeeID,
		CourseName:       req.CourseName,
		Description:      req.Description,
		TrainingType:     trainingType,
		DueDate:          dueDate,
		ExternalProvider: req.ExternalProvider,
		ExternalCourseID: req.ExternalCourseID,
	}
	if err := s.repo.CreateTraining(ctx, t); err != nil {
		return nil, err
	}
	s.emitEvent(ctx, companyID, t.ID, "training_assignment", "training.assigned", nil)
	return t, nil
}

func (s *Service) ListTraining(ctx context.Context, onboardingID string) ([]TrainingAssignment, error) {
	return s.repo.ListTraining(ctx, onboardingID)
}

// Dashboard

func (s *Service) GetDashboard(ctx context.Context, companyID string) (*OnboardingDashboard, error) {
	return s.repo.GetDashboard(ctx, companyID)
}

func (s *Service) GetEmployeeDashboard(ctx context.Context, companyID, employeeID string) (*EmployeeDashboard, error) {
	return s.repo.GetEmployeeDashboard(ctx, companyID, employeeID)
}

// Process with details

func (s *Service) GetProcessWithDetails(ctx context.Context, companyID, id string) (*ProcessWithDetails, error) {
	p, err := s.repo.GetProcess(ctx, companyID, id)
	if err != nil {
		return nil, err
	}

	pd := &ProcessWithDetails{Process: *p}

	pd.Tasks, _ = s.repo.ListTasks(ctx, id)
	pd.Documents, _ = s.repo.ListDocuments(ctx, id)
	pd.Assets, _ = s.repo.ListAssets(ctx, id)
	pd.Access, _ = s.repo.ListAccessRequests(ctx, id)
	pd.Milestones, _ = s.repo.ListMilestones(ctx, id)
	pd.Feedback, _ = s.repo.ListFeedback(ctx, id)
	pd.Buddies, _ = s.repo.ListBuddies(ctx, id)
	pd.Training, _ = s.repo.ListTraining(ctx, id)
	pd.Progress, _ = s.calculateProgress(ctx, id)

	return pd, nil
}

// Template with tasks

func (s *Service) GetTemplateWithTasks(ctx context.Context, companyID, id string) (*TemplateWithTasks, error) {
	t, err := s.repo.GetTemplate(ctx, companyID, id)
	if err != nil {
		return nil, err
	}

	tt := &TemplateWithTasks{Template: *t}
	tt.Tasks, _ = s.repo.ListTemplateTasks(ctx, id)

	return tt, nil
}

// Events

func (s *Service) emitEvent(ctx context.Context, companyID, aggregateID, aggregateType, eventType string, payload interface{}) {
	var rawPayload json.RawMessage
	if payload != nil {
		if b, err := json.Marshal(payload); err == nil {
			rawPayload = b
		}
	}
	e := &DomainEvent{
		EventType:     eventType,
		CompanyID:     companyID,
		AggregateID:   aggregateID,
		AggregateType: aggregateType,
		Payload:       rawPayload,
	}
	s.repo.CreateEvent(ctx, e)
}

// Notifications

func (s *Service) createNotification(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) {
	n := &Notification{
		CompanyID:        companyID,
		UserID:           userID,
		Title:            title,
		Body:             &body,
		NotificationType: notifType,
		Channel:          "IN_APP",
		ReferenceType:    &refType,
		ReferenceID:      &refID,
	}
	s.repo.CreateNotification(ctx, n)
}

// Handle candidate.hired event from FASE 15
func (s *Service) HandleCandidateHired(ctx context.Context, companyID string, event *CandidateHiredEvent) (*OnboardingProcess, error) {
	startDate, err := parseDate(event.StartDate)
	if err != nil {
		startDate = time.Now()
	}

	employeeID := event.EmployeeID

	policy := "STRICT"
	req := &CreateOnboardingRequest{
		EmployeeID:       employeeID,
		StartDate:        startDate.Format("2006-01-02"),
		CompletionPolicy: &policy,
	}

	return s.CreateOnboarding(ctx, companyID, event.EmployeeID, req)
}

// IA Assistant - Generate template proposal based on position/department
func (s *Service) GenerateTemplateProposal(ctx context.Context, req *IAOnboardingRequest) (*IATemplateProposal, error) {
	proposal := &IATemplateProposal{}

	switch req.Position {
	case "Desarrollador Backend", "Desarrollador Backend Go", "Backend Developer":
		proposal.Tasks = []IATaskProposal{
			{Title: "Revisar documentación técnica", Description: "Revisar documentación del proyecto y guías de desarrollo", Category: "IT", ResponsibleType: "EMPLOYEE", DaysOffset: -7},
			{Title: "Preparar notebook", Description: "Preparar equipo con herramientas de desarrollo", Category: "IT", ResponsibleType: "IT", DaysOffset: -7},
			{Title: "Crear usuario corporativo", Description: "Crear usuario en sistemas corporativos", Category: "IT", ResponsibleType: "IT", DaysOffset: -7},
			{Title: "Crear email corporativo", Description: "Configurar email corporativo", Category: "IT", ResponsibleType: "IT", DaysOffset: -2},
			{Title: "Configurar VPN", Description: "Configurar acceso VPN corporativo", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: -2},
			{Title: "Crear acceso a Git", Description: "Crear usuario en GitHub/GitLab corporativo", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: -2},
			{Title: "Firma de documentación", Description: "Firmar contrato y documentación legal", Category: "LEGAL", ResponsibleType: "HR", DaysOffset: 0},
			{Title: "Entrega de notebook", Description: "Recibir equipo asignado", Category: "EQUIPMENT", ResponsibleType: "IT", DaysOffset: 0},
			{Title: "Presentación al equipo", Description: "Presentación con el equipo de trabajo", Category: "TEAM", ResponsibleType: "MANAGER", DaysOffset: 1},
			{Title: "Acceso a sistemas", Description: "Configurar acceso a todos los sistemas necesarios", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: 1},
			{Title: "Capacitación de seguridad", Description: "Completar capacitación de seguridad informática", Category: "TRAINING", ResponsibleType: "EMPLOYEE", DaysOffset: 7},
			{Title: "Feedback 30 días", Description: "Reunión de feedback con manager", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 30},
			{Title: "Seguimiento 60 días", Description: "Reunión de seguimiento con manager", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 60},
			{Title: "Evaluación 90 días", Description: "Evaluación de desempeño inicial", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 90},
		}
	case "Desarrollador Frontend", "Frontend Developer":
		proposal.Tasks = []IATaskProposal{
			{Title: "Revisar documentación técnica", Description: "Revisar guías de frontend y diseño", Category: "IT", ResponsibleType: "EMPLOYEE", DaysOffset: -7},
			{Title: "Preparar notebook", Description: "Preparar equipo con herramientas de desarrollo", Category: "IT", ResponsibleType: "IT", DaysOffset: -7},
			{Title: "Crear email corporativo", Description: "Configurar email corporativo", Category: "IT", ResponsibleType: "IT", DaysOffset: -2},
			{Title: "Configurar VPN", Description: "Configurar acceso VPN", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: -2},
			{Title: "Crear acceso a Git", Description: "Crear usuario en repositorio corporativo", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: -2},
			{Title: "Firma de documentación", Description: "Firmar contrato y documentación legal", Category: "LEGAL", ResponsibleType: "HR", DaysOffset: 0},
			{Title: "Entrega de notebook", Description: "Recibir equipo asignado", Category: "EQUIPMENT", ResponsibleType: "IT", DaysOffset: 0},
			{Title: "Presentación al equipo", Description: "Presentación con el equipo de trabajo", Category: "TEAM", ResponsibleType: "MANAGER", DaysOffset: 1},
			{Title: "Acceso a sistemas de diseño", Description: "Configurar acceso a Figma/Adobe", Category: "ACCESS", ResponsibleType: "IT", DaysOffset: 1},
			{Title: "Capacitación de seguridad", Description: "Completar capacitación de seguridad", Category: "TRAINING", ResponsibleType: "EMPLOYEE", DaysOffset: 7},
			{Title: "Feedback 30 días", Description: "Reunión de feedback con manager", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 30},
			{Title: "Evaluación 90 días", Description: "Evaluación de desempeño inicial", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 90},
		}
	default:
		proposal.Tasks = []IATaskProposal{
			{Title: "Completar datos personales", Description: "Completar formulario de datos personales", Category: "PERSONAL", ResponsibleType: "EMPLOYEE", DaysOffset: 0},
			{Title: "Firma de documentación", Description: "Firmar contrato y documentación legal", Category: "LEGAL", ResponsibleType: "HR", DaysOffset: 0},
			{Title: "Presentación al equipo", Description: "Presentación con el equipo de trabajo", Category: "TEAM", ResponsibleType: "MANAGER", DaysOffset: 1},
			{Title: "Revisar políticas internas", Description: "Revisar políticas y procedimientos de la empresa", Category: "HR", ResponsibleType: "EMPLOYEE", DaysOffset: 1},
			{Title: "Capacitación inicial", Description: "Completar capacitación de inducción", Category: "TRAINING", ResponsibleType: "EMPLOYEE", DaysOffset: 7},
			{Title: "Feedback 30 días", Description: "Reunión de feedback con manager", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 30},
			{Title: "Seguimiento 60 días", Description: "Reunión de seguimiento con manager", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 60},
			{Title: "Evaluación 90 días", Description: "Evaluación de desempeño inicial", Category: "MANAGER", ResponsibleType: "MANAGER", DaysOffset: 90},
		}
	}

	return proposal, nil
}
