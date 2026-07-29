package http

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/onboarding/domain"
	"github.com/rrhhumand/api/internal/onboarding/integration"
	"github.com/rrhhumand/api/internal/onboarding/workflow"
)

type Handler struct {
	onboardingEngine  *workflow.OnboardingEngine
	offboardingEngine *workflow.OffboardingEngine
	taskEngine        *workflow.TaskEngine
	onbRepo           onboardingRepo
	offbRepo          offboardingRepo
	taskRepo          taskRepo
	docRepo           docRepo
	sharedRepo        sharedRepo
	empSvc            integration.EmployeeService
	atsSvc            integration.ATSIntegration
}

type onboardingRepo interface {
	GetByID(ctx context.Context, companyID, id string) (*domain.OnboardingProcess, error)
	List(ctx context.Context, companyID string, status, employeeID, search string) ([]domain.OnboardingProcess, error)
	HasActiveProcess(ctx context.Context, companyID, employeeID string) (bool, error)
	UpdateProbation(ctx context.Context, companyID, id string, status domain.ProbationStatus) error
	GetDashboardStats(ctx context.Context, companyID string) (active, pending, completed, overdue int, avgProgress float64, err error)
}

type offboardingRepo interface {
	GetProcessByID(ctx context.Context, companyID, id string) (*domain.OffboardingProcess, error)
	ListProcesses(ctx context.Context, companyID string, status, employeeID string) ([]domain.OffboardingProcess, error)
	HasActiveProcess(ctx context.Context, companyID, employeeID string) (bool, error)
	ListExitReasons(ctx context.Context, companyID string) ([]domain.TerminationReason, error)
	GetExitInterview(ctx context.Context, companyID, offboardingID string) (*domain.ExitInterview, error)
	CreateExitInterview(ctx context.Context, e *domain.ExitInterview) error
	CompleteExitInterview(ctx context.Context, companyID, offboardingID string, reason, feedback string, recommendation *string, rating *float64) error
	ListAssets(ctx context.Context, offboardingID string) ([]domain.OffboardingAsset, error)
	CreateAsset(ctx context.Context, a *domain.OffboardingAsset) error
	UpdateAssetStatus(ctx context.Context, companyID, id string, status domain.OffboardingAssetStatus, conditionOnReturn *string) error
	ListAccessRevocations(ctx context.Context, offboardingID string) ([]domain.AccessRevocation, error)
	CreateAccessRevocation(ctx context.Context, a *domain.AccessRevocation) error
	UpdateAccessRevocation(ctx context.Context, companyID, id string, status string, performedBy *string, errMsg *string) error
	CreateHandover(ctx context.Context, h *domain.EmployeeHandover) error
	GetHandover(ctx context.Context, companyID, offboardingID string) (*domain.EmployeeHandover, error)
	CompleteHandover(ctx context.Context, companyID, id string) error
	GetTaskCounts(ctx context.Context, offboardingID string) (total, completed int, err error)
	ListTasks(ctx context.Context, offboardingID string) ([]domain.OffboardingTask, error)
	CreateTask(ctx context.Context, t *domain.OffboardingTask) error
	GetTask(ctx context.Context, companyID, id string) (*domain.OffboardingTask, error)
	GetDashboardStats(ctx context.Context, companyID string) (active, pending, completed, overdue int, err error)
}

type taskRepo interface {
	ListAssignments(ctx context.Context, onboardingID string) ([]domain.OnboardingTaskAssignment, error)
	GetCounts(ctx context.Context, onboardingID string) (total, completed int, err error)
}

type docRepo interface {
	ListByOnboarding(ctx context.Context, onboardingID string) ([]domain.OnboardingDocument, error)
	CountPendingReview(ctx context.Context, companyID string) (int, error)
}

type sharedRepo interface {
	CreateWorkflowRule(ctx context.Context, rule *domain.WorkflowRule) error
	ListWorkflowRules(ctx context.Context, companyID string, workflowType domain.WorkflowType) ([]domain.WorkflowRule, error)
}

func NewHandler(
	onboardingEngine *workflow.OnboardingEngine,
	offboardingEngine *workflow.OffboardingEngine,
	taskEngine *workflow.TaskEngine,
	onbRepo onboardingRepo,
	offbRepo offboardingRepo,
	taskRepo taskRepo,
	docRepo docRepo,
	sharedRepo sharedRepo,
	empSvc integration.EmployeeService,
	atsSvc integration.ATSIntegration,
) *Handler {
	return &Handler{
		onboardingEngine:  onboardingEngine,
		offboardingEngine: offboardingEngine,
		taskEngine:        taskEngine,
		onbRepo:           onbRepo,
		offbRepo:          offbRepo,
		taskRepo:          taskRepo,
		docRepo:           docRepo,
		sharedRepo:        sharedRepo,
		empSvc:            empSvc,
		atsSvc:            atsSvc,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	onb := rg.Group("/onboarding")
	{
		onb.GET("", h.ListOnboardings)
		onb.POST("", h.CreateOnboarding)
		onb.GET("/:id", h.GetOnboarding)
		onb.PUT("/:id", h.UpdateOnboarding)

		onb.POST("/:id/start", h.StartOnboarding)
		onb.POST("/:id/complete", h.CompleteOnboarding)
		onb.POST("/:id/cancel", h.CancelOnboarding)
		onb.POST("/:id/block", h.BlockOnboarding)

		onb.GET("/:id/tasks", h.ListOnboardingTasks)
		onb.POST("/:id/tasks", h.CreateOnboardingTask)

		onb.POST("/tasks/:taskId/start", h.StartTask)
		onb.POST("/tasks/:taskId/complete", h.CompleteTask)
		onb.POST("/tasks/:taskId/block", h.BlockTask)

		onb.GET("/:id/documents", h.ListDocuments)
		onb.POST("/:id/documents", h.CreateDocument)
		onb.POST("/documents/:docId/approve", h.ApproveDocument)
		onb.POST("/documents/:docId/reject", h.RejectDocument)

		onb.GET("/:id/assets", h.ListOnboardingAssets)
		onb.POST("/:id/assets", h.CreateOnboardingAsset)
		onb.POST("/assets/:assetId/assign", h.AssignOnboardingAsset)
		onb.POST("/assets/:assetId/deliver", h.DeliverOnboardingAsset)

		onb.GET("/:id/access", h.ListAccessRequests)
		onb.POST("/:id/access", h.CreateAccessRequest)

		onb.GET("/:id/checklist", h.GetChecklist)
		onb.POST("/checklist/:itemId/complete", h.CompleteChecklistItem)

		onb.GET("/:id/notes", h.ListNotes)
		onb.POST("/:id/notes", h.CreateNote)

		onb.POST("/:id/probation", h.UpdateProbation)

		onb.POST("/:id/buddy", h.AssignBuddy)
		onb.GET("/:id/buddy", h.GetBuddy)

		onb.POST("/:id/training", h.AssignTraining)

		onb.GET("/dashboard", h.GetOnboardingDashboard)
		onb.GET("/employee-dashboard", h.GetEmployeeOnboardingDashboard)

		onb.GET("/candidates-ready", h.GetCandidatesReadyForOnboarding)
	}

	offb := rg.Group("/offboarding")
	{
		offb.GET("", h.ListOffboardings)
		offb.POST("", h.CreateOffboarding)
		offb.GET("/:id", h.GetOffboarding)
		offb.PUT("/:id", h.UpdateOffboarding)

		offb.POST("/:id/approve", h.ApproveOffboarding)
		offb.POST("/:id/start", h.StartOffboarding)
		offb.POST("/:id/complete", h.CompleteOffboarding)
		offb.POST("/:id/cancel", h.CancelOffboarding)

		offb.GET("/:id/tasks", h.ListOffboardingTasks)
		offb.POST("/:id/tasks", h.CreateOffboardingTask)
		offb.POST("/tasks/:taskId/complete", h.CompleteOffboardingTask)

		offb.GET("/:id/assets", h.ListOffboardingAssets)
		offb.POST("/:id/assets", h.CreateOffboardingAsset)
		offb.POST("/assets/:assetId/return", h.ReturnOffboardingAsset)
		offb.POST("/assets/:assetId/report-damaged", h.ReportAssetDamaged)
		offb.POST("/assets/:assetId/report-lost", h.ReportAssetLost)

		offb.GET("/:id/access", h.ListAccessRevocations)
		offb.POST("/:id/access", h.CreateAccessRevocation)
		offb.POST("/access/:accessId/revoke", h.RevokeAccess)
		offb.POST("/access/:accessId/retry", h.RetryAccessRevocation)

		offb.GET("/:id/exit-interview", h.GetExitInterview)
		offb.POST("/:id/exit-interview", h.CreateExitInterview)
		offb.PUT("/:id/exit-interview", h.UpdateExitInterview)
		offb.POST("/:id/exit-interview/complete", h.CompleteExitInterview)

		offb.GET("/:id/handover", h.GetHandover)
		offb.POST("/:id/handover", h.CreateHandover)
		offb.POST("/handover/:handoverId/complete", h.CompleteHandover)

		offb.POST("/:id/final-settlement", h.RequestFinalSettlement)

		offb.GET("/exit-reasons", h.ListExitReasons)

		offb.GET("/dashboard", h.GetOffboardingDashboard)
	}

	rg.POST("/workflow-rules", h.CreateWorkflowRule)
	rg.GET("/workflow-rules", h.ListWorkflowRules)
	rg.POST("/workflow-rules/evaluate", h.EvaluateWorkflowRules)
}
