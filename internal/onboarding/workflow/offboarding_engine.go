package workflow

import (
	"context"
	"fmt"
	"math"

	"github.com/rrhhumand/api/internal/onboarding/domain"
	"github.com/rrhhumand/api/internal/onboarding/integration"
)

type OffboardingEngine struct {
	empSvc          integration.EmployeeService
	payrollSvc      integration.PayrollService
	notifSvc        integration.NotificationService
	accessSvc       integration.AccessProvisioningService
	docSvc          integration.DocumentService
	offboardingRepo offboardingRepositoryFull
	taskRepo        offboardingTaskRepository
	sharedRepo      sharedRepository
}

type offboardingRepositoryFull interface {
	GetProcessByID(ctx context.Context, companyID, id string) (*domain.OffboardingProcess, error)
	UpdateStatus(ctx context.Context, companyID, id string, status domain.OffboardingStatus) error
	UpdateProgress(ctx context.Context, companyID, id string, progress float64) error
	Complete(ctx context.Context, companyID, id string) error
	ListTasks(ctx context.Context, offboardingID string) ([]domain.OffboardingTask, error)
	GetTaskCounts(ctx context.Context, offboardingID string) (total, completed int, err error)
}

type offboardingTaskRepository interface {
	GetTask(ctx context.Context, companyID, id string) (*domain.OffboardingTask, error)
	CompleteTask(ctx context.Context, companyID, id, completedBy string) error
}

func NewOffboardingEngine(
	empSvc integration.EmployeeService,
	payrollSvc integration.PayrollService,
	notifSvc integration.NotificationService,
	accessSvc integration.AccessProvisioningService,
	docSvc integration.DocumentService,
	offboardingRepo offboardingRepositoryFull,
	taskRepo offboardingTaskRepository,
	sharedRepo sharedRepository,
) *OffboardingEngine {
	return &OffboardingEngine{
		empSvc:          empSvc,
		payrollSvc:      payrollSvc,
		notifSvc:        notifSvc,
		accessSvc:       accessSvc,
		docSvc:          docSvc,
		offboardingRepo: offboardingRepo,
		taskRepo:        taskRepo,
		sharedRepo:      sharedRepo,
	}
}

func (e *OffboardingEngine) Start(ctx context.Context, companyID, processID string) error {
	p, err := e.offboardingRepo.GetProcessByID(ctx, companyID, processID)
	if err != nil {
		return fmt.Errorf("offboarding not found: %w", err)
	}
	if p.Status != domain.OffboardingApproved {
		return fmt.Errorf("cannot start offboarding in status %s", p.Status)
	}

	if err := e.offboardingRepo.UpdateStatus(ctx, companyID, processID, domain.OffboardingInProgress); err != nil {
		return err
	}

	e.notify(ctx, companyID, p.EmployeeID, "Proceso de salida iniciado",
		"Se ha iniciado tu proceso de desvinculación.",
		"OFFBOARDING_STARTED", "offboarding", processID)
	e.emitEvent(ctx, companyID, processID, "offboarding", string(domain.EventOffboardingStarted), nil)

	return nil
}

func (e *OffboardingEngine) ExecuteTask(ctx context.Context, companyID, taskID, completedBy string) error {
	task, err := e.taskRepo.GetTask(ctx, companyID, taskID)
	if err != nil {
		return fmt.Errorf("offboarding task not found: %w", err)
	}

	if err := e.taskRepo.CompleteTask(ctx, companyID, taskID, completedBy); err != nil {
		return err
	}

	e.updateProgress(ctx, companyID, task.OffboardingID)
	e.emitEvent(ctx, companyID, taskID, "offboarding_task", string(domain.EventOffboardingTaskCompleted), nil)

	return nil
}

func (e *OffboardingEngine) Complete(ctx context.Context, companyID, processID string) error {
	p, err := e.offboardingRepo.GetProcessByID(ctx, companyID, processID)
	if err != nil {
		return err
	}
	if p.Status == domain.OffboardingCompleted {
		return domain.ErrOffboardingCompleted
	}

	tasks, err := e.offboardingRepo.ListTasks(ctx, processID)
	if err != nil {
		return err
	}
	for _, t := range tasks {
		if t.Required && t.Status != domain.OffTaskCompleted && t.Status != domain.OffTaskCancelled {
			return fmt.Errorf("required task not completed: %s", t.Title)
		}
	}

	if err := e.offboardingRepo.Complete(ctx, companyID, processID); err != nil {
		return err
	}

	e.emitEvent(ctx, companyID, processID, "offboarding", string(domain.EventOffboardingCompleted),
		map[string]string{
			"offboarding_id":   processID,
			"employee_id":      p.EmployeeID,
			"company_id":       companyID,
			"last_working_date": p.LastWorkingDate,
		})

	e.notify(ctx, companyID, p.RequestedBy, "Offboarding completado",
		"El proceso de desvinculación ha sido completado.",
		"OFFBOARDING_COMPLETED", "offboarding", processID)

	return nil
}

func (e *OffboardingEngine) Approve(ctx context.Context, companyID, processID string) error {
	p, err := e.offboardingRepo.GetProcessByID(ctx, companyID, processID)
	if err != nil {
		return err
	}
	if p.Status != domain.OffboardingPendingApproval {
		return fmt.Errorf("cannot approve offboarding in status %s", p.Status)
	}

	if err := e.offboardingRepo.UpdateStatus(ctx, companyID, processID, domain.OffboardingApproved); err != nil {
		return err
	}

	e.emitEvent(ctx, companyID, processID, "offboarding", string(domain.EventOffboardingApproved), nil)
	return nil
}

func (e *OffboardingEngine) Cancel(ctx context.Context, companyID, processID string) error {
	if err := e.offboardingRepo.UpdateStatus(ctx, companyID, processID, domain.OffboardingCancelled); err != nil {
		return err
	}
	return nil
}

func (e *OffboardingEngine) RequestFinalSettlement(ctx context.Context, companyID, processID string) error {
	p, err := e.offboardingRepo.GetProcessByID(ctx, companyID, processID)
	if err != nil {
		return err
	}

	if err := e.payrollSvc.StartFinalSettlement(ctx, companyID, p.EmployeeID,
		string(p.TerminationType), p.LastWorkingDate); err != nil {
		return fmt.Errorf("failed to request final settlement: %w", err)
	}

	e.emitEvent(ctx, companyID, processID, "offboarding", string(domain.EventOffboardingFinalSettlementRequested), nil)
	return nil
}

func (e *OffboardingEngine) ChangeEmployeeStatus(ctx context.Context, companyID, processID string) error {
	p, err := e.offboardingRepo.GetProcessByID(ctx, companyID, processID)
	if err != nil {
		return err
	}

	if _, err := e.empSvc.Update(ctx, p.EmployeeID, companyID, map[string]string{
		"status":            p.EmployeeStatusAfter,
		"termination_date": p.LastWorkingDate,
	}); err != nil {
		return fmt.Errorf("failed to update employee status: %w", err)
	}

	e.emitEvent(ctx, companyID, processID, "employee", string(domain.EventEmployeeTerminated),
		map[string]string{
			"employee_id": p.EmployeeID,
			"company_id":  companyID,
			"status":      p.EmployeeStatusAfter,
		})

	return nil
}

func (e *OffboardingEngine) RevokeAllAccess(ctx context.Context, companyID, processID string, systems []string) error {
	for _, system := range systems {
		if err := e.accessSvc.RevokeAccess(ctx, processID, system); err != nil {
			e.sharedRepo.CreateAuditLog(ctx, companyID, "", "ACCESS_REVOCATION_FAILED",
				"offboarding", processID, "", "", nil,
				map[string]string{"system": system, "error": err.Error()})
		}
	}
	return nil
}

func (e *OffboardingEngine) updateProgress(ctx context.Context, companyID, offboardingID string) {
	total, completed, err := e.offboardingRepo.GetTaskCounts(ctx, offboardingID)
	if err != nil || total == 0 {
		return
	}
	progress := math.Round(float64(completed) / float64(total) * 100)
	e.offboardingRepo.UpdateProgress(ctx, companyID, offboardingID, progress)
}

func (e *OffboardingEngine) notify(ctx context.Context, companyID, userID, title, body, notifType, refType, refID string) {
	e.notifSvc.Send(ctx, companyID, userID, title, body, notifType, refType, refID)
	e.sharedRepo.CreateNotification(ctx, companyID, userID, title, body, notifType, refType, refID)
}

func (e *OffboardingEngine) emitEvent(ctx context.Context, companyID, aggregateID, aggregateType, eventType string, payload interface{}) {
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
