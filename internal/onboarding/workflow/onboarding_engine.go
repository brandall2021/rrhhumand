package workflow

import (
	"context"
	"fmt"
	"math"

	"github.com/rrhhumand/api/internal/onboarding/domain"
	"github.com/rrhhumand/api/internal/onboarding/integration"
)

type OnboardingConfig struct {
	NotifyAboutOverdueTasks bool
	DefaultProbationDays    int
}

type OnboardingEngine struct {
	config             OnboardingConfig
	empSvc             integration.EmployeeService
	docSvc             integration.DocumentService
	assetSvc           integration.AssetService
	trainingSvc        integration.TrainingService
	notifSvc           integration.NotificationService
	atsSvc             integration.ATSIntegration
	accessSvc          integration.AccessProvisioningService
	signSvc            integration.SignatureService
	calendarSvc        integration.CalendarService
	taskEngine         *TaskEngine
	onboardingRepo     onboardingRepository
	taskRepo           taskRepository
	docRepo            documentRepository
	offboardingRepo    offboardingRepository
	sharedRepo         sharedRepository
}

type onboardingRepository interface {
	GetByID(ctx context.Context, companyID, id string) (*domain.OnboardingProcess, error)
	UpdateStatus(ctx context.Context, companyID, id string, status domain.OnboardingStatus) error
	UpdateProgress(ctx context.Context, companyID, id string, progress float64) error
	Complete(ctx context.Context, companyID, id string) error
	Cancel(ctx context.Context, companyID, id, reason string) error
	HasActiveProcess(ctx context.Context, companyID, employeeID string) (bool, error)
}

type taskRepository interface {
	GetAssignment(ctx context.Context, companyID, id string) (*domain.OnboardingTaskAssignment, error)
	ListAssignments(ctx context.Context, onboardingID string) ([]domain.OnboardingTaskAssignment, error)
	CompleteAssignment(ctx context.Context, id, completedBy string) error
	BlockAssignment(ctx context.Context, id, comments string) error
	UpdateAssignmentStatus(ctx context.Context, id string, status domain.TaskStatus) error
	GetCounts(ctx context.Context, onboardingID string) (total, completed int, err error)
	AreDependenciesMet(ctx context.Context, taskID string) (bool, error)
}

type documentRepository interface {
	GetByID(ctx context.Context, companyID, id string) (*domain.OnboardingDocument, error)
	ListByOnboarding(ctx context.Context, onboardingID string) ([]domain.OnboardingDocument, error)
	UpdateStatus(ctx context.Context, companyID, id string, status domain.DocStatus, reviewedBy *string) error
}

type offboardingRepository interface {
	GetProcessByID(ctx context.Context, companyID, id string) (*domain.OffboardingProcess, error)
}

type sharedRepository interface {
	CreateAuditLog(ctx context.Context, companyID, userID, action, entityType, entityID, ipAddress, userAgent string, oldVal, newVal interface{}) error
	CreateNotification(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) error
	CreateOutboxEvent(ctx context.Context, e *domain.OutboxEvent) error
}

func NewOnboardingEngine(
	config OnboardingConfig,
	empSvc integration.EmployeeService,
	docSvc integration.DocumentService,
	assetSvc integration.AssetService,
	trainingSvc integration.TrainingService,
	notifSvc integration.NotificationService,
	atsSvc integration.ATSIntegration,
	accessSvc integration.AccessProvisioningService,
	signSvc integration.SignatureService,
	calendarSvc integration.CalendarService,
	taskEngine *TaskEngine,
	onboardingRepo onboardingRepository,
	taskRepo taskRepository,
	docRepo documentRepository,
	offboardingRepo offboardingRepository,
	sharedRepo sharedRepository,
) *OnboardingEngine {
	return &OnboardingEngine{
		config:          config,
		empSvc:          empSvc,
		docSvc:          docSvc,
		assetSvc:        assetSvc,
		trainingSvc:     trainingSvc,
		notifSvc:        notifSvc,
		atsSvc:          atsSvc,
		accessSvc:       accessSvc,
		signSvc:         signSvc,
		calendarSvc:     calendarSvc,
		taskEngine:      taskEngine,
		onboardingRepo:  onboardingRepo,
		taskRepo:        taskRepo,
		docRepo:         docRepo,
		offboardingRepo: offboardingRepo,
		sharedRepo:      sharedRepo,
	}
}

func (e *OnboardingEngine) Start(ctx context.Context, companyID, processID string) error {
	p, err := e.onboardingRepo.GetByID(ctx, companyID, processID)
	if err != nil {
		return fmt.Errorf("onboarding not found: %w", err)
	}
	if p.Status != domain.OnboardingDraft && p.Status != domain.OnboardingPending {
		return fmt.Errorf("cannot start onboarding in status %s", p.Status)
	}

	if err := e.onboardingRepo.UpdateStatus(ctx, companyID, processID, domain.OnboardingInProgress); err != nil {
		return err
	}

	e.notify(ctx, companyID, p.EmployeeID, "Onboarding iniciado",
		"Tu proceso de incorporación ha comenzado. Revisa tus tareas pendientes.",
		"ONBOARDING_STARTED", "onboarding", processID)
	e.emitEvent(ctx, companyID, processID, "onboarding", string(domain.EventOnboardingStarted), nil)

	return nil
}

func (e *OnboardingEngine) ExecuteTask(ctx context.Context, companyID, assignmentID, completedBy string) error {
	assignment, err := e.taskRepo.GetAssignment(ctx, companyID, assignmentID)
	if err != nil {
		return fmt.Errorf("task assignment not found: %w", err)
	}

	if assignment.Status == domain.TaskCompleted {
		return fmt.Errorf("task already completed")
	}

	met, err := e.taskRepo.AreDependenciesMet(ctx, assignment.TaskID)
	if err != nil {
		return err
	}
	if !met {
		return domain.ErrTaskDependencyNotMet
	}

	if err := e.taskRepo.CompleteAssignment(ctx, assignmentID, completedBy); err != nil {
		return err
	}

	e.updateProgress(ctx, companyID, assignment.OnboardingID)
	e.notify(ctx, companyID, completedBy, "Tarea completada",
		"Has completado una tarea del onboarding.",
		"TASK_COMPLETED", "onboarding_task", assignmentID)
	e.emitEvent(ctx, companyID, assignmentID, "onboarding_task", string(domain.EventOnboardingTaskCompleted), nil)

	return nil
}

func (e *OnboardingEngine) BlockTask(ctx context.Context, companyID, assignmentID, reason string) error {
	assignment, err := e.taskRepo.GetAssignment(ctx, companyID, assignmentID)
	if err != nil {
		return err
	}
	if assignment.Status == domain.TaskCompleted {
		return fmt.Errorf("cannot block a completed task")
	}

	if err := e.taskRepo.BlockAssignment(ctx, assignmentID, reason); err != nil {
		return err
	}

	p, err := e.onboardingRepo.GetByID(ctx, companyID, assignment.OnboardingID)
	if err == nil {
		e.notify(ctx, companyID, p.CreatedBy, "Tarea bloqueada",
			fmt.Sprintf("La tarea requiere atención: %s", reason),
			"TASK_BLOCKED", "onboarding_task", assignmentID)
	}

	return nil
}

func (e *OnboardingEngine) Complete(ctx context.Context, companyID, processID string) error {
	p, err := e.onboardingRepo.GetByID(ctx, companyID, processID)
	if err != nil {
		return err
	}
	if p.Status == domain.OnboardingCompleted {
		return domain.ErrOnboardingCompleted
	}
	if p.Status == domain.OnboardingCancelled {
		return domain.ErrOnboardingCancelled
	}

	assignments, err := e.taskRepo.ListAssignments(ctx, processID)
	if err != nil {
		return err
	}
	for _, a := range assignments {
		if a.Status != domain.TaskCompleted && a.Status != domain.TaskCancelled {
			return fmt.Errorf("task not completed: %s", a.ID)
		}
	}

	if err := e.onboardingRepo.Complete(ctx, companyID, processID); err != nil {
		return err
	}

	e.notify(ctx, companyID, p.EmployeeID, "Onboarding completado",
		"¡Felicitaciones! Has completado tu proceso de incorporación.",
		"ONBOARDING_COMPLETED", "onboarding", processID)
	e.emitEvent(ctx, companyID, processID, "onboarding", string(domain.EventOnboardingCompleted),
		map[string]string{
			"onboarding_id": processID,
			"employee_id":   p.EmployeeID,
			"company_id":    companyID,
		})

	return nil
}

func (e *OnboardingEngine) Block(ctx context.Context, companyID, processID string) error {
	if err := e.onboardingRepo.UpdateStatus(ctx, companyID, processID, domain.OnboardingBlocked); err != nil {
		return err
	}
	return nil
}

func (e *OnboardingEngine) Cancel(ctx context.Context, companyID, processID, reason string) error {
	if err := e.onboardingRepo.Cancel(ctx, companyID, processID, reason); err != nil {
		return err
	}
	return nil
}

func (e *OnboardingEngine) ApproveDocument(ctx context.Context, companyID, documentID, reviewedBy string) error {
	doc, err := e.docRepo.GetByID(ctx, companyID, documentID)
	if err != nil {
		return err
	}
	if doc.Status == domain.DocApproved {
		return fmt.Errorf("document already approved")
	}

	if err := e.docRepo.UpdateStatus(ctx, companyID, documentID, domain.DocApproved, &reviewedBy); err != nil {
		return err
	}

	e.updateProgress(ctx, companyID, doc.OnboardingID)
	e.emitEvent(ctx, companyID, documentID, "onboarding_document", string(domain.EventOnboardingDocumentApproved), nil)
	return nil
}

func (e *OnboardingEngine) RejectDocument(ctx context.Context, companyID, documentID, reviewedBy string) error {
	if err := e.docRepo.UpdateStatus(ctx, companyID, documentID, domain.DocRejected, &reviewedBy); err != nil {
		return err
	}
	doc, _ := e.docRepo.GetByID(ctx, companyID, documentID)
	e.emitEvent(ctx, companyID, documentID, "onboarding_document", string(domain.EventOnboardingDocumentRejected), nil)
	e.notify(ctx, companyID, doc.EmployeeID, "Documento rechazado",
		"Tu documento ha sido rechazado. Por favor, súbelo nuevamente.",
		"DOCUMENT_REJECTED", "onboarding_document", documentID)
	return nil
}

func (e *OnboardingEngine) updateProgress(ctx context.Context, companyID, onboardingID string) {
	total, completed, err := e.taskRepo.GetCounts(ctx, onboardingID)
	if err != nil || total == 0 {
		return
	}
	progress := math.Round(float64(completed) / float64(total) * 100)
	e.onboardingRepo.UpdateProgress(ctx, companyID, onboardingID, progress)
}

func (e *OnboardingEngine) notify(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) {
	e.notifSvc.Send(ctx, companyID, userID, title, body, notifType, refType, refID)
	e.sharedRepo.CreateNotification(ctx, companyID, userID, title, body, notifType, refType, refID)
}

func (e *OnboardingEngine) emitEvent(ctx context.Context, companyID, aggregateID, aggregateType, eventType string, payload interface{}) {
	var payloadStr *string
	if payload != nil {
		s := fmt.Sprintf("%v", payload)
		payloadStr = &s
	}
	evt := &domain.OutboxEvent{
		CompanyID:     companyID,
		EventType:     eventType,
		AggregateType: aggregateType,
		AggregateID:   aggregateID,
		Payload:       payloadStr,
		Status:        string(domain.OutboxPending),
	}
	e.sharedRepo.CreateOutboxEvent(ctx, evt)
}
