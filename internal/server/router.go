package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/attendance"
	"github.com/rrhhumand/api/internal/auth"
	"github.com/rrhhumand/api/internal/branches"
	"github.com/rrhhumand/api/internal/compensation"
	"github.com/rrhhumand/api/internal/companies"
	"github.com/rrhhumand/api/internal/departments"
	"github.com/rrhhumand/api/internal/document_categories"
	"github.com/rrhhumand/api/internal/documents"
	"github.com/rrhhumand/api/internal/employees"
	"github.com/rrhhumand/api/internal/feed"
	"github.com/rrhhumand/api/internal/handlers"
	"github.com/rrhhumand/api/internal/leave"
	"github.com/rrhhumand/api/internal/middleware"
	onbhttp "github.com/rrhhumand/api/internal/onboarding/http"
	"github.com/rrhhumand/api/internal/organization"
	"github.com/rrhhumand/api/internal/overtime"
	"github.com/rrhhumand/api/internal/payroll"
	featureshttp "github.com/rrhhumand/api/internal/payroll/features/http"
	benefitshttp "github.com/rrhhumand/api/internal/benefits/http"
	expenseshttp "github.com/rrhhumand/api/internal/expenses/http"
	perfhttp "github.com/rrhhumand/api/internal/performance/http"
	"github.com/rrhhumand/api/internal/positions"
	"github.com/rrhhumand/api/internal/profile"
	recrhttp "github.com/rrhhumand/api/internal/recruitment/http"
	"github.com/rrhhumand/api/internal/roles"
	"github.com/rrhhumand/api/internal/scheduling"
	"github.com/rrhhumand/api/internal/surveys"
	"github.com/rrhhumand/api/internal/training"
	"github.com/rrhhumand/api/pkg/response"
)

func NewRouter(
	healthHandler *handlers.HealthHandler,
	authHandler *auth.Handler,
	jwtService *auth.JWTService,
	companyHandler *companies.CompanyHandler,
	branchHandler *branches.BranchHandler,
	departmentHandler *departments.DepartmentHandler,
	positionHandler *positions.PositionHandler,
	employeeHandler *employees.EmployeeHandler,
	orgHandler *organization.OrgHandler,
	profileHandler *profile.ProfileHandler,
	feedHandler *feed.FeedHandler,
	surveyHandler *surveys.SurveyHandler,
	docHandler *documents.DocumentHandler,
	categoryHandler *document_categories.CategoryHandler,
	leaveHandler *leave.Handler,
	attHandler *attendance.Handler,
	schedHandler *scheduling.Handler,
	otHandler *overtime.Handler,
	payHandler *payroll.Handler,
	perfHandler *perfhttp.Handler,
	recHandler *recrhttp.Handler,
	trainingHandler *training.Handler,
	onbFase23Handler *onbhttp.Handler,
	compHandler *compensation.Handler,
	featuresHandler *featureshttp.Handler,
	benefitsHandler *benefitshttp.Handler,
	expensesHandler *expenseshttp.Handler,
	pool *pgxpool.Pool,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)
		v1.GET("/ready", healthHandler.Ready)

		authGroup := v1.Group("/auth")
		{
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.Refresh)
			authGroup.POST("/logout", authHandler.Logout)
		}

		publicGroup := v1.Group("")
		{
			recHandler.RegisterPublicRoutes(publicGroup)
		}

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService, func(ctx context.Context, userID, companyID string) (string, error) {
			var employeeID string
			err := pool.QueryRow(ctx,
				`SELECT e.id FROM user_companies uc
				 JOIN employees e ON e.company_id = uc.company_id AND e.status = 'active'
				 WHERE uc.user_id = $1 AND uc.company_id = $2
				 LIMIT 1`, userID, companyID).Scan(&employeeID)
			if err != nil {
				return "", err
			}
			return employeeID, nil
		}))
		{
			protected.GET("/auth/me", authHandler.Me)

			protected.GET("/roles", func(c *gin.Context) {
				roleRepo := roles.NewRoleRepository(pool)
				roleService := roles.NewRoleService(roleRepo)
				allRoles, err := roleService.GetAll(c.Request.Context())
				if err != nil {
					response.InternalError(c, "Failed to fetch roles")
					return
				}
				response.Success(c, allRoles)
			})

			protected.POST("/companies", companyHandler.Create)
			protected.GET("/companies", companyHandler.List)
			protected.GET("/companies/:id", companyHandler.GetByID)
			protected.PUT("/companies/:id", companyHandler.Update)

			protected.Use(middleware.TenantMiddleware())

			protected.POST("/branches", branchHandler.Create)
			protected.GET("/branches", branchHandler.List)
			protected.GET("/branches/:id", branchHandler.GetByID)
			protected.PUT("/branches/:id", branchHandler.Update)
			protected.DELETE("/branches/:id", branchHandler.Delete)

			protected.POST("/departments", departmentHandler.Create)
			protected.GET("/departments", departmentHandler.List)
			protected.GET("/departments/:id", departmentHandler.GetByID)
			protected.PUT("/departments/:id", departmentHandler.Update)
			protected.DELETE("/departments/:id", departmentHandler.Delete)

			protected.POST("/positions", positionHandler.Create)
			protected.GET("/positions", positionHandler.List)
			protected.GET("/positions/:id", positionHandler.GetByID)
			protected.PUT("/positions/:id", positionHandler.Update)
			protected.DELETE("/positions/:id", positionHandler.Delete)

			protected.POST("/employees", employeeHandler.Create)
			protected.GET("/employees", employeeHandler.List)
			protected.GET("/employees/:id", employeeHandler.GetByID)
			protected.PUT("/employees/:id", employeeHandler.Update)
			protected.DELETE("/employees/:id", employeeHandler.Delete)

			protected.GET("/employees/:id/contacts", employeeHandler.GetContacts)
			protected.PUT("/employees/:id/contacts", employeeHandler.UpsertContacts)
			protected.GET("/employees/:id/addresses", employeeHandler.GetAddresses)
			protected.PUT("/employees/:id/addresses", employeeHandler.UpsertAddresses)
			protected.GET("/employees/:id/emergency-contacts", employeeHandler.GetEmergencyContacts)
			protected.PUT("/employees/:id/emergency-contacts", employeeHandler.UpsertEmergencyContacts)
			protected.GET("/employees/:id/history", employeeHandler.GetHistory)

			protected.GET("/organization/tree", orgHandler.GetTree)

			protected.GET("/me/profile", profileHandler.GetProfile)
			protected.PUT("/me/profile", profileHandler.UpdateProfile)

			protected.GET("/feed", feedHandler.ListPosts)
			protected.POST("/feed", feedHandler.CreatePost)
			protected.GET("/feed/:id", feedHandler.GetPost)
			protected.PUT("/feed/:id", feedHandler.UpdatePost)
			protected.DELETE("/feed/:id", feedHandler.DeletePost)
			protected.POST("/feed/:id/comments", feedHandler.AddComment)
			protected.POST("/feed/:id/reactions", feedHandler.AddReaction)
			protected.DELETE("/feed/:id/reactions/:type", feedHandler.RemoveReaction)

			protected.GET("/surveys", surveyHandler.ListSurveys)
			protected.POST("/surveys", surveyHandler.CreateSurvey)
			protected.GET("/surveys/:id", surveyHandler.GetSurvey)
			protected.PUT("/surveys/:id", surveyHandler.UpdateSurvey)
			protected.DELETE("/surveys/:id", surveyHandler.DeleteSurvey)
			protected.POST("/surveys/:id/publish", surveyHandler.PublishSurvey)
			protected.POST("/surveys/:id/close", surveyHandler.CloseSurvey)
			protected.POST("/surveys/:id/archive", surveyHandler.ArchiveSurvey)
			protected.GET("/surveys/:id/questions", surveyHandler.ListQuestions)
			protected.POST("/surveys/:id/questions", surveyHandler.AddQuestion)
			protected.PUT("/surveys/questions/:id", surveyHandler.UpdateQuestion)
			protected.DELETE("/surveys/questions/:id", surveyHandler.DeleteQuestion)
			protected.POST("/surveys/questions/:id/options", surveyHandler.AddOption)
			protected.DELETE("/surveys/options/:id", surveyHandler.DeleteOption)
			protected.PUT("/surveys/:id/targets", surveyHandler.SetTargets)
			protected.GET("/surveys/:id/targets", surveyHandler.ListTargets)
			protected.GET("/me/surveys", surveyHandler.ListAvailableSurveys)
			protected.POST("/surveys/:id/respond", surveyHandler.RespondSurvey)
			protected.GET("/surveys/:id/results", surveyHandler.GetResults)
			protected.GET("/surveys/:id/export", surveyHandler.ExportSurvey)

			protected.GET("/documents", docHandler.ListDocuments)
			protected.POST("/documents", docHandler.UploadDocument)
			protected.GET("/documents/:id", docHandler.GetDocument)
			protected.PUT("/documents/:id", docHandler.UpdateDocument)
			protected.DELETE("/documents/:id", docHandler.DeleteDocument)
			protected.GET("/documents/:id/download", docHandler.DownloadDocument)
			protected.POST("/documents/:id/versions", docHandler.CreateVersion)
			protected.GET("/documents/:id/versions", docHandler.ListVersions)
			protected.POST("/documents/:id/archive", docHandler.ArchiveDocument)
			protected.POST("/documents/:id/restore", docHandler.RestoreDocument)
			protected.DELETE("/documents/:id/permanent", docHandler.PermanentDelete)
			protected.PUT("/documents/:id/permissions", docHandler.SetPermissions)
			protected.GET("/documents/:id/permissions", docHandler.ListPermissions)
			protected.POST("/documents/:id/share", docHandler.CreateShare)
			protected.POST("/documents/:id/share-link", docHandler.CreateShareLink)
			protected.POST("/documents/:id/tags", docHandler.SetDocumentTags)
			protected.GET("/documents/stats", docHandler.GetDocumentStats)
			protected.GET("/documents/expiring", docHandler.ListExpiringDocuments)
			protected.GET("/employees/:id/documents", docHandler.ListEmployeeDocuments)

			protected.GET("/document-categories", categoryHandler.List)
			protected.POST("/document-categories", categoryHandler.Create)
			protected.GET("/document-categories/:id", categoryHandler.GetByID)
			protected.PUT("/document-categories/:id", categoryHandler.Update)
			protected.DELETE("/document-categories/:id", categoryHandler.Delete)

			protected.GET("/document-tags", docHandler.ListTags)
			protected.POST("/document-tags", docHandler.CreateTag)

			protected.GET("/share/:token", docHandler.AccessShareLink)

			protected.POST("/leave/types", leaveHandler.CreateLeaveType)
			protected.GET("/leave/types", leaveHandler.ListLeaveTypes)
			protected.GET("/leave/types/:id", leaveHandler.GetLeaveType)
			protected.PUT("/leave/types/:id", leaveHandler.UpdateLeaveType)
			protected.DELETE("/leave/types/:id", leaveHandler.DeleteLeaveType)

			protected.POST("/leave/policies", leaveHandler.CreateLeavePolicy)
			protected.GET("/leave/policies", leaveHandler.ListLeavePolicies)

			protected.POST("/leave/holidays", leaveHandler.CreateHoliday)
			protected.GET("/leave/holidays", leaveHandler.ListHolidays)
			protected.DELETE("/leave/holidays/:id", leaveHandler.DeleteHoliday)

			protected.GET("/leave/balance", leaveHandler.GetBalance)
			protected.POST("/leave/balances/adjust", leaveHandler.AdjustBalance)

			protected.POST("/leave/requests", leaveHandler.CreateLeaveRequest)
			protected.GET("/leave/requests", leaveHandler.ListLeaveRequests)
			protected.GET("/leave/requests/:id", leaveHandler.GetLeaveRequest)
			protected.POST("/leave/requests/:id/approve", leaveHandler.ApproveRequest)
			protected.POST("/leave/requests/:id/reject", leaveHandler.RejectRequest)
			protected.POST("/leave/requests/:id/cancel", leaveHandler.CancelRequest)
			protected.GET("/leave/requests/:id/history", leaveHandler.GetRequestHistory)

			protected.GET("/leave/calendar", leaveHandler.GetCalendar)
			protected.GET("/leave/calendar/team", leaveHandler.GetTeamCalendar)

			protected.GET("/leave/reports", leaveHandler.GetReport)

			protected.POST("/attendance/clock-in", attHandler.ClockIn)
			protected.POST("/attendance/clock-out", attHandler.ClockOut)
			protected.POST("/attendance/break/start", attHandler.StartBreak)
			protected.POST("/attendance/break/end", attHandler.EndBreak)
			protected.GET("/attendance/me", attHandler.GetMyAttendance)
			protected.GET("/attendance", attHandler.ListRecords)
			protected.GET("/attendance/:id", attHandler.GetRecord)
			protected.GET("/attendance/dashboard", attHandler.GetDashboard)
			protected.GET("/attendance/dashboard/team", attHandler.GetTeamDashboard)
			protected.GET("/attendance/calendar", attHandler.GetCalendar)
			protected.POST("/attendance/policies", attHandler.CreatePolicy)
			protected.POST("/attendance/locations", attHandler.CreateLocation)
			protected.POST("/attendance/devices", attHandler.CreateDevice)
			protected.POST("/attendance/corrections", attHandler.CreateCorrection)
			protected.GET("/attendance/corrections", attHandler.ListCorrections)
			protected.POST("/attendance/corrections/:id/approve", attHandler.ApproveCorrection)
			protected.POST("/attendance/corrections/:id/reject", attHandler.RejectCorrection)
			protected.GET("/attendance/export", attHandler.ExportCSV)

			// Scheduling
			protected.POST("/schedules", schedHandler.CreateSchedule)
			protected.GET("/schedules", schedHandler.ListSchedules)
			protected.GET("/schedules/:id", schedHandler.GetSchedule)
			protected.PUT("/schedules/:id", schedHandler.UpdateSchedule)
			protected.DELETE("/schedules/:id", schedHandler.DeleteSchedule)

			protected.POST("/shifts", schedHandler.CreateShift)
			protected.GET("/shifts", schedHandler.ListShifts)
			protected.GET("/shifts/:id", schedHandler.GetShift)
			protected.PUT("/shifts/:id", schedHandler.UpdateShift)
			protected.DELETE("/shifts/:id", schedHandler.DeleteShift)

			protected.POST("/employees/:id/schedule", schedHandler.AssignSchedule)
			protected.GET("/employees/:id/schedule", schedHandler.GetEmployeeSchedule)
			protected.POST("/employees/:id/shift", schedHandler.AssignShift)
			protected.GET("/employees/:id/shift", schedHandler.GetEmployeeShift)
			protected.POST("/employees/:id/rotation", schedHandler.AssignRotation)

			protected.GET("/employees/:id/schedule/resolved", schedHandler.GetResolvedSchedule)

			protected.POST("/scheduling/calendar/generate", schedHandler.GenerateCalendar)
			protected.GET("/scheduling/calendar", schedHandler.ListCalendar)

			protected.POST("/scheduling/exceptions", schedHandler.CreateException)
			protected.GET("/scheduling/exceptions", schedHandler.ListExceptions)

			protected.POST("/scheduling/swap", schedHandler.SwapShift)
			protected.POST("/scheduling/swap/:id/approve", schedHandler.ApproveSwap)
			protected.POST("/scheduling/swap/:id/reject", schedHandler.RejectSwap)

			protected.POST("/rotation-templates", schedHandler.CreateRotationTemplate)
			protected.GET("/rotation-templates", schedHandler.ListRotationTemplates)
			protected.GET("/rotation-templates/:id", schedHandler.GetRotationTemplate)

			// Overtime
			protected.GET("/overtime", otHandler.ListRecords)
			protected.GET("/overtime/:id", otHandler.GetRecord)
			protected.POST("/overtime/:id/approve", otHandler.ApproveRecord)
			protected.POST("/overtime/:id/reject", otHandler.RejectRecord)
			protected.POST("/overtime/detect", otHandler.DetectOvertime)
			protected.GET("/overtime/dashboard", otHandler.GetDashboard)

			protected.POST("/overtime/requests", otHandler.CreateRequest)
			protected.GET("/overtime/requests", otHandler.ListRequests)
			protected.GET("/overtime/requests/:id", otHandler.GetRequest)
			protected.POST("/overtime/requests/:id/approve", otHandler.ApproveRequest)
			protected.POST("/overtime/requests/:id/reject", otHandler.RejectRequest)

			protected.POST("/overtime-policies", otHandler.CreatePolicy)
			protected.GET("/overtime-policies", otHandler.ListPolicies)
			protected.GET("/overtime-policies/:id", otHandler.GetPolicy)
			protected.PUT("/overtime-policies/:id", otHandler.UpdatePolicy)
			protected.DELETE("/overtime-policies/:id", otHandler.DeletePolicy)

			protected.POST("/compensations", otHandler.RequestCompensation)
			protected.GET("/compensations", otHandler.ListCompensations)
			protected.POST("/compensations/:id/approve", otHandler.ApproveCompensation)
			protected.POST("/compensations/:id/reject", otHandler.RejectCompensation)
			protected.POST("/compensations/:id/cancel", otHandler.CancelCompensation)

			protected.GET("/employees/:id/overtime", otHandler.GetEmployeeOvertime)
			protected.GET("/employees/:id/time-balance", otHandler.GetBalance)
			protected.POST("/employees/:id/time-balance/adjust", otHandler.AdjustBalance)
			protected.GET("/employees/:id/time-balance/transactions", otHandler.GetBalanceTransactions)

			// Payroll / Liquidación de Haberes (FASE 19A)
			protected.GET("/payroll/dashboard", payHandler.GetDashboardStats)

			protected.POST("/payroll/periods", payHandler.CreatePeriod)
			protected.GET("/payroll/periods", payHandler.ListPeriods)
			protected.GET("/payroll/periods/:id", payHandler.GetPeriod)
			protected.PUT("/payroll/periods/:id", payHandler.UpdatePeriod)
			protected.POST("/payroll/periods/:id/close", payHandler.ClosePeriod)

			protected.GET("/payroll/periods/:id/runs", payHandler.ListRuns)
			protected.POST("/payroll/periods/:id/runs", payHandler.CreateRun)
			protected.GET("/payroll/runs/:id", payHandler.GetRun)
			protected.POST("/payroll/runs/:id/calculate", payHandler.CalculateRun)
			protected.POST("/payroll/runs/:id/validate", payHandler.ValidateRun)
			protected.POST("/payroll/runs/:id/approve", payHandler.ApproveRun)
			protected.POST("/payroll/runs/:id/close", payHandler.CloseRun)
			protected.GET("/payroll/runs/:id/summary", payHandler.GetRunSummary)
			protected.GET("/payroll/runs/:id/errors", payHandler.ListErrors)

			protected.POST("/payroll/runs/:id/employees", payHandler.AddEmployeeToRun)
			protected.GET("/payroll/runs/:id/employees", payHandler.ListRunEmployees)
			protected.GET("/payroll/runs/:id/employees/:eid/result", payHandler.GetEmployeeResult)
			protected.GET("/payroll/runs/:id/employees/:eid/items", payHandler.GetEmployeeItems)
			protected.GET("/payroll/runs/:id/employees/:eid/bases", payHandler.GetEmployeeBases)

			protected.POST("/payroll/concepts", payHandler.CreateConcept)
			protected.GET("/payroll/concepts", payHandler.ListConcepts)
			protected.PUT("/payroll/concepts/:id", payHandler.UpdateConcept)

			protected.POST("/payroll/rules", payHandler.CreateRule)
			protected.GET("/payroll/rules", payHandler.ListRules)
			protected.PUT("/payroll/rules/:id", payHandler.UpdateRule)

			protected.POST("/payroll/novelties", payHandler.CreateNovelty)
			protected.GET("/payroll/novelties", payHandler.ListNovelties)
			protected.PUT("/payroll/novelties/:id", payHandler.UpdateNovelty)
			protected.DELETE("/payroll/novelties/:id", payHandler.DeleteNovelty)
			protected.POST("/payroll/novelties/:id/approve", payHandler.ApproveNovelty)
			protected.POST("/payroll/novelties/import", payHandler.ImportNovelties)

			protected.POST("/payroll/advances", payHandler.CreateAdvance)
			protected.GET("/payroll/employees/:employee_id/advances", payHandler.ListAdvances)

			protected.POST("/payroll/garnishments", payHandler.CreateGarnishment)
			protected.GET("/payroll/employees/:employee_id/garnishments", payHandler.ListGarnishments)

			protected.POST("/payroll/agreements", payHandler.CreateAgreement)
			protected.GET("/payroll/agreements", payHandler.ListAgreements)

			protected.POST("/payroll/categories", payHandler.CreateCategory)
			protected.GET("/payroll/categories", payHandler.ListCategories)

			protected.POST("/payroll/salary-scales", payHandler.CreateSalaryScale)
			protected.GET("/payroll/salary-scales", payHandler.ListSalaryScales)

			// Employee self-service
			protected.GET("/me/payroll/periods", payHandler.MyPeriods)
			protected.GET("/me/payroll/runs/:run_id", payHandler.MyRunResult)
			protected.GET("/me/payroll/runs/:run_id/items", payHandler.MyItems)

			// Payroll Features — FASE 19B (Recibos, ARCA, Libro, Bancos, Contabilidad, Reportes)
			protected.POST("/payroll/receipt-templates", featuresHandler.CreateTemplate)
			protected.GET("/payroll/receipt-templates", featuresHandler.ListTemplates)
			protected.GET("/payroll/receipt-templates/:id", featuresHandler.GetTemplate)
			protected.PUT("/payroll/receipt-templates/:id", featuresHandler.UpdateTemplate)
			protected.DELETE("/payroll/receipt-templates/:id", featuresHandler.DeleteTemplate)
			protected.POST("/payroll/receipts/generate", featuresHandler.GenerateReceipts)
			protected.GET("/payroll/receipts", featuresHandler.ListReceipts)
			protected.GET("/payroll/receipts/:id", featuresHandler.GetReceipt)
			protected.GET("/payroll/receipts/:id/items", featuresHandler.GetReceiptItems)
			protected.POST("/payroll/receipts/:id/acknowledge", featuresHandler.AcknowledgeReceipt)

			protected.POST("/payroll/arca/mappings", featuresHandler.CreateMappingArca)
			protected.GET("/payroll/arca/mappings", featuresHandler.ListMappingsArca)
			protected.PUT("/payroll/arca/mappings/:id", featuresHandler.UpdateMappingArca)
			protected.DELETE("/payroll/arca/mappings/:id", featuresHandler.DeleteMappingArca)
			protected.POST("/payroll/arca/exports/generate", featuresHandler.GenerateExportArca)
			protected.GET("/payroll/arca/exports", featuresHandler.ListExportsArca)
			protected.GET("/payroll/arca/exports/:id", featuresHandler.GetExportArca)
			protected.POST("/payroll/arca/exports/:id/validate", featuresHandler.ValidateExportArca)
			protected.GET("/payroll/arca/exports/:id/download", featuresHandler.DownloadExportArca)

			protected.POST("/payroll/book/entries/generate", featuresHandler.GenerateEntries)
			protected.GET("/payroll/book/entries", featuresHandler.ListEntries)
			protected.GET("/payroll/book/entries/:id", featuresHandler.GetEntry)
			protected.POST("/payroll/book/exports", featuresHandler.ExportBook)
			protected.GET("/payroll/book/exports", featuresHandler.ListExportsBook)
			protected.GET("/payroll/book/exports/:id", featuresHandler.GetExportBook)
			protected.GET("/payroll/book/summary", featuresHandler.GetBookSummary)

			protected.POST("/payroll/bank/batches", featuresHandler.CreateBatch)
			protected.GET("/payroll/bank/batches", featuresHandler.ListBatches)
			protected.GET("/payroll/bank/batches/:id", featuresHandler.GetBatch)
			protected.GET("/payroll/bank/batches/:id/items", featuresHandler.GetBatchItems)
			protected.POST("/payroll/bank/batches/:id/generate-file", featuresHandler.GenerateBankFile)
			protected.POST("/payroll/bank/batches/:id/send", featuresHandler.SendBatch)
			protected.GET("/payroll/bank/batches/:id/summary", featuresHandler.GetBatchSummary)

			protected.POST("/payroll/accounting/mappings", featuresHandler.CreateMappingAccounting)
			protected.GET("/payroll/accounting/mappings", featuresHandler.ListMappingsAccounting)
			protected.PUT("/payroll/accounting/mappings/:id", featuresHandler.UpdateMappingAccounting)
			protected.POST("/payroll/accounting/exports/generate", featuresHandler.GenerateExportAccounting)
			protected.GET("/payroll/accounting/exports", featuresHandler.ListExportsAccounting)
			protected.GET("/payroll/accounting/exports/:id", featuresHandler.GetExportAccounting)
			protected.GET("/payroll/accounting/exports/:id/entries", featuresHandler.GetEntriesAccounting)
			protected.GET("/payroll/accounting/exports/:id/download", featuresHandler.DownloadExportAccounting)

			protected.POST("/payroll/report-templates", featuresHandler.CreateTemplateReport)
			protected.GET("/payroll/report-templates", featuresHandler.ListTemplatesReport)
			protected.GET("/payroll/report-templates/:id", featuresHandler.GetTemplateReport)
			protected.PUT("/payroll/report-templates/:id", featuresHandler.UpdateTemplateReport)
			protected.DELETE("/payroll/report-templates/:id", featuresHandler.DeleteTemplateReport)
			protected.POST("/payroll/reports/generate", featuresHandler.GenerateReport)
			protected.GET("/payroll/reports", featuresHandler.ListExportsReport)
			protected.GET("/payroll/reports/:id", featuresHandler.GetExportReport)

			// Employee self-service (FASE 19B)
			protected.GET("/me/payroll/receipts", featuresHandler.MyReceipts)

			// Benefits & Total Rewards — FASE 20
			// Catalog
			protected.POST("/benefits/categories", benefitsHandler.CreateCategory)
			protected.GET("/benefits/categories", benefitsHandler.ListCategories)
			protected.GET("/benefits/categories/:id", benefitsHandler.GetCategory)
			protected.PUT("/benefits/categories/:id", benefitsHandler.UpdateCategory)
			protected.DELETE("/benefits/categories/:id", benefitsHandler.DeleteCategory)

			protected.POST("/benefits/types", benefitsHandler.CreateType)
			protected.GET("/benefits/types", benefitsHandler.ListTypes)
			protected.GET("/benefits/types/:id", benefitsHandler.GetType)
			protected.PUT("/benefits/types/:id", benefitsHandler.UpdateType)
			protected.DELETE("/benefits/types/:id", benefitsHandler.DeleteType)

			protected.POST("/benefits/providers", benefitsHandler.CreateProvider)
			protected.GET("/benefits/providers", benefitsHandler.ListProviders)
			protected.GET("/benefits/providers/:id", benefitsHandler.GetProvider)
			protected.PUT("/benefits/providers/:id", benefitsHandler.UpdateProvider)
			protected.DELETE("/benefits/providers/:id", benefitsHandler.DeleteProvider)

			// Benefits CRUD
			protected.POST("/benefits", benefitsHandler.CreateBenefit)
			protected.GET("/benefits", benefitsHandler.ListBenefits)
			protected.GET("/benefits/search", benefitsHandler.SearchBenefits)
			protected.GET("/benefits/:id", benefitsHandler.GetBenefit)
			protected.PUT("/benefits/:id", benefitsHandler.UpdateBenefit)
			protected.DELETE("/benefits/:id", benefitsHandler.DeleteBenefit)

			// Plans
			protected.POST("/benefits/:id/plans", benefitsHandler.CreatePlan)
			protected.GET("/benefits/:id/plans", benefitsHandler.ListPlans)
			protected.GET("/benefits/plans/:planId", benefitsHandler.GetPlan)
			protected.PUT("/benefits/plans/:planId", benefitsHandler.UpdatePlan)
			protected.DELETE("/benefits/plans/:planId", benefitsHandler.DeletePlan)

			// Eligibility
			protected.POST("/benefits/eligibility/rules", benefitsHandler.CreateRule)
			protected.GET("/benefits/eligibility/rules", benefitsHandler.ListRules)
			protected.GET("/benefits/eligibility/rules/:id", benefitsHandler.GetRule)
			protected.PUT("/benefits/eligibility/rules/:id", benefitsHandler.UpdateRule)
			protected.DELETE("/benefits/eligibility/rules/:id", benefitsHandler.DeleteRule)
			protected.POST("/benefits/eligibility/evaluate", benefitsHandler.EvaluateEmployee)
			protected.GET("/benefits/eligibility/eligible/:employee_id", benefitsHandler.ListEligibleBenefits)

			// Workflows
			protected.POST("/benefits/workflows", benefitsHandler.CreateWorkflow)
			protected.GET("/benefits/workflows", benefitsHandler.ListWorkflows)
			protected.GET("/benefits/workflows/:id", benefitsHandler.GetWorkflow)
			protected.PUT("/benefits/workflows/:id", benefitsHandler.UpdateWorkflow)
			protected.POST("/benefits/workflows/:id/steps", benefitsHandler.AddStep)
			protected.GET("/benefits/workflows/:id/steps", benefitsHandler.ListSteps)
			protected.PUT("/benefits/workflows/steps/:stepId", benefitsHandler.UpdateStep)
			protected.DELETE("/benefits/workflows/steps/:stepId", benefitsHandler.DeleteStep)

			// Employee assignments
			protected.POST("/benefits/assignments", benefitsHandler.EnrollEmployee)
			protected.GET("/benefits/assignments", benefitsHandler.ListEmployeeBenefits)
			protected.GET("/benefits/assignments/:id", benefitsHandler.GetEmployeeBenefit)
			protected.PUT("/benefits/assignments/:id", benefitsHandler.UpdateEmployeeBenefit)
			protected.POST("/benefits/assignments/:id/cancel", benefitsHandler.CancelEmployeeBenefit)
			protected.GET("/benefits/assignments/:id/history", benefitsHandler.GetHistory)
			protected.GET("/benefits/employees/:employee_id/history", benefitsHandler.ListHistoryByEmployee)

			// Requests
			protected.POST("/benefits/requests", benefitsHandler.CreateRequest)
			protected.GET("/benefits/requests", benefitsHandler.ListRequests)
			protected.GET("/benefits/requests/:id", benefitsHandler.GetRequest)
			protected.POST("/benefits/requests/:id/submit", benefitsHandler.SubmitRequest)
			protected.POST("/benefits/requests/:id/approve", benefitsHandler.ApproveRequest)
			protected.POST("/benefits/requests/:id/reject", benefitsHandler.RejectRequest)
			protected.POST("/benefits/requests/:id/cancel", benefitsHandler.CancelRequest)
			protected.GET("/benefits/requests/:id/reviews", benefitsHandler.ListReviews)

			// Wallets
			protected.POST("/benefits/wallets", benefitsHandler.CreateWallet)
			protected.GET("/benefits/wallets/:id", benefitsHandler.GetWallet)
			protected.GET("/benefits/employees/:employee_id/wallets", benefitsHandler.ListEmployeeWallets)
			protected.POST("/benefits/wallets/:id/credit", benefitsHandler.CreditWallet)
			protected.POST("/benefits/wallets/:id/debit", benefitsHandler.DebitWallet)
			protected.GET("/benefits/wallets/:id/transactions", benefitsHandler.ListTransactions)

			// Flexible plans
			protected.POST("/benefits/flexible-plans", benefitsHandler.CreateFlexiblePlan)
			protected.GET("/benefits/flexible-plans", benefitsHandler.ListFlexiblePlans)
			protected.POST("/benefits/flexible-budgets", benefitsHandler.CreateBudget)
			protected.GET("/benefits/flexible-budgets/:employee_id/:plan_id/:year", benefitsHandler.GetBudget)
			protected.GET("/benefits/employees/:employee_id/flexible-budgets", benefitsHandler.ListEmployeeBudgets)

			// Costs
			protected.POST("/benefits/costs", benefitsHandler.CreateCost)
			protected.GET("/benefits/costs", benefitsHandler.ListCosts)
			protected.PUT("/benefits/costs/:id", benefitsHandler.UpdateCost)
			protected.POST("/benefits/cost-schedules", benefitsHandler.CreateSchedule)
			protected.GET("/benefits/cost-schedules", benefitsHandler.ListSchedules)
			protected.POST("/benefits/cost-schedules/:id/pay", benefitsHandler.MarkSchedulePaid)

			// Reimbursements
			protected.POST("/benefits/reimbursements", benefitsHandler.CreateReimbursement)
			protected.GET("/benefits/reimbursements", benefitsHandler.ListReimbursements)
			protected.GET("/benefits/reimbursements/:id", benefitsHandler.GetReimbursement)
			protected.POST("/benefits/reimbursements/:id/approve", benefitsHandler.ApproveReimbursement)
			protected.POST("/benefits/reimbursements/:id/reject", benefitsHandler.RejectReimbursement)
			protected.POST("/benefits/reimbursements/:id/pay", benefitsHandler.PayReimbursement)
			protected.POST("/benefits/reimbursements/:id/cancel", benefitsHandler.CancelReimbursement)
			protected.POST("/benefits/reimbursements/:id/documents", benefitsHandler.UploadDocument)
			protected.GET("/benefits/reimbursements/:id/documents", benefitsHandler.ListDocuments)

			// Bonuses
			protected.POST("/benefits/bonuses", benefitsHandler.CreateBonus)
			protected.GET("/benefits/bonuses", benefitsHandler.ListBonuses)
			protected.GET("/benefits/bonuses/:id", benefitsHandler.GetBonus)
			protected.PUT("/benefits/bonuses/:id", benefitsHandler.UpdateBonus)
			protected.POST("/benefits/bonuses/:id/approve", benefitsHandler.ApproveBonus)
			protected.POST("/benefits/bonuses/:id/pay", benefitsHandler.PayBonus)
			protected.POST("/benefits/bonuses/:id/cancel", benefitsHandler.CancelBonus)

			// Incentives
			protected.POST("/benefits/incentives", benefitsHandler.CreateIncentive)
			protected.GET("/benefits/incentives", benefitsHandler.ListIncentives)
			protected.GET("/benefits/incentives/:id", benefitsHandler.GetIncentive)
			protected.POST("/benefits/incentives/:id/redeem", benefitsHandler.RedeemIncentive)

			// Payroll integration
			protected.POST("/benefits/payroll-mappings", benefitsHandler.CreatePayrollMapping)
			protected.GET("/benefits/payroll-mappings", benefitsHandler.ListPayrollMappings)
			protected.POST("/benefits/payroll-mappings/:id/sync", benefitsHandler.SyncToPayroll)

			// Total Rewards
			protected.POST("/benefits/rewards-items", benefitsHandler.CreateRewardsItem)
			protected.GET("/benefits/rewards-items", benefitsHandler.ListRewardsItems)
			protected.PUT("/benefits/rewards-items/:id", benefitsHandler.UpdateRewardsItem)
			protected.POST("/benefits/rewards/snapshots/generate", benefitsHandler.GenerateSnapshot)
			protected.GET("/benefits/rewards/snapshots/:id", benefitsHandler.GetLatestSnapshot)
			protected.GET("/benefits/rewards/snapshots", benefitsHandler.ListSnapshots)

			// Reports
			protected.POST("/benefits/report-definitions", benefitsHandler.CreateReportDefinition)
			protected.GET("/benefits/report-definitions", benefitsHandler.ListReportDefinitions)

			// Notifications
			protected.POST("/benefits/notifications", benefitsHandler.LogNotification)
			protected.GET("/benefits/notifications", benefitsHandler.ListNotifications)
			protected.POST("/benefits/notifications/:id/read", benefitsHandler.MarkNotificationRead)

			// Employee self-service
			protected.GET("/me/benefits", benefitsHandler.MyBenefits)
			protected.GET("/me/benefits/wallets", benefitsHandler.MyWallets)
			protected.GET("/me/benefits/reimbursements", benefitsHandler.MyReimbursements)
			protected.GET("/me/benefits/requests", benefitsHandler.MyRequests)
			protected.GET("/me/benefits/bonuses", benefitsHandler.MyBonuses)
			protected.GET("/me/benefits/incentives", benefitsHandler.MyIncentives)
			protected.GET("/me/benefits/rewards-statement", benefitsHandler.MyRewards)

			// Expenses & Travel — FASE 21
			// Categories
			protected.POST("/expenses/categories", expensesHandler.CreateCategory)
			protected.GET("/expenses/categories", expensesHandler.ListCategories)
			protected.GET("/expenses/categories/:id", expensesHandler.GetCategory)
			protected.PUT("/expenses/categories/:id", expensesHandler.UpdateCategory)
			protected.DELETE("/expenses/categories/:id", expensesHandler.DeleteCategory)

			// Payment methods
			protected.POST("/expenses/payment-methods", expensesHandler.CreatePaymentMethod)
			protected.GET("/expenses/payment-methods", expensesHandler.ListPaymentMethods)

			// Expenses
			protected.POST("/expenses", expensesHandler.CreateExpense)
			protected.GET("/expenses", expensesHandler.ListExpenses)
			protected.GET("/expenses/:id", expensesHandler.GetExpense)
			protected.PUT("/expenses/:id", expensesHandler.UpdateExpense)
			protected.POST("/expenses/:id/submit", expensesHandler.SubmitExpense)
			protected.POST("/expenses/:id/approve", expensesHandler.ApproveExpense)
			protected.POST("/expenses/:id/reject", expensesHandler.RejectExpense)
			protected.POST("/expenses/:id/observe", expensesHandler.ObserveExpense)
			protected.POST("/expenses/:id/cancel", expensesHandler.CancelExpense)
			protected.DELETE("/expenses/:id", expensesHandler.DeleteExpense)
			protected.POST("/expenses/:id/override-policy", expensesHandler.OverridePolicy)

			// Receipts
			protected.POST("/expenses/:id/receipts", expensesHandler.UploadReceipt)
			protected.GET("/expenses/:id/receipts", expensesHandler.ListReceipts)
			protected.DELETE("/expenses/:id/receipts/:receiptId", expensesHandler.DeleteReceipt)

			// Travels
			protected.POST("/expenses/travels", expensesHandler.CreateTravel)
			protected.GET("/expenses/travels", expensesHandler.ListTravels)
			protected.GET("/expenses/travels/:id", expensesHandler.GetTravel)
			protected.PUT("/expenses/travels/:id", expensesHandler.UpdateTravel)
			protected.POST("/expenses/travels/:id/request", expensesHandler.RequestTravel)
			protected.POST("/expenses/travels/:id/approve", expensesHandler.ApproveTravel)
			protected.POST("/expenses/travels/:id/reject", expensesHandler.RejectTravel)
			protected.POST("/expenses/travels/:id/complete", expensesHandler.CompleteTravel)
			protected.POST("/expenses/travels/:id/cancel", expensesHandler.CancelTravel)
			protected.POST("/expenses/travels/:id/participants", expensesHandler.AddParticipant)
			protected.DELETE("/expenses/travels/:id/participants/:employeeId", expensesHandler.RemoveParticipant)
			protected.GET("/expenses/travels/:id/participants", expensesHandler.ListParticipants)

			// Expense Reports
			protected.POST("/expenses/reports", expensesHandler.CreateReport)
			protected.GET("/expenses/reports", expensesHandler.ListReports)
			protected.GET("/expenses/reports/:id", expensesHandler.GetReport)
			protected.POST("/expenses/reports/:id/submit", expensesHandler.SubmitReport)
			protected.POST("/expenses/reports/:id/approve", expensesHandler.ApproveReport)
			protected.POST("/expenses/reports/:id/reject", expensesHandler.RejectReport)
			protected.POST("/expenses/reports/:id/observe", expensesHandler.ObserveReport)

			// Advances
			protected.POST("/expenses/advances", expensesHandler.RequestAdvance)
			protected.GET("/expenses/advances", expensesHandler.ListAdvances)
			protected.GET("/expenses/advances/:id", expensesHandler.GetAdvance)
			protected.POST("/expenses/advances/:id/approve", expensesHandler.ApproveAdvance)
			protected.POST("/expenses/advances/:id/pay", expensesHandler.PayAdvance)
			protected.POST("/expenses/advances/:id/reject", expensesHandler.RejectAdvance)
			protected.POST("/expenses/advances/:id/cancel", expensesHandler.CancelAdvance)

			// Reimbursements
			protected.POST("/expenses/reimbursements", expensesHandler.CreateReimbursement)
			protected.GET("/expenses/reimbursements", expensesHandler.ListReimbursements)
			protected.GET("/expenses/reimbursements/:id", expensesHandler.GetReimbursement)
			protected.POST("/expenses/reimbursements/:id/approve", expensesHandler.ApproveReimbursement)
			protected.POST("/expenses/reimbursements/:id/pay", expensesHandler.PayReimbursement)
			protected.POST("/expenses/reimbursements/:id/reject", expensesHandler.RejectReimbursement)
			protected.POST("/expenses/reimbursements/:id/cancel", expensesHandler.CancelReimbursement)

			// Policies
			protected.POST("/expenses/policies", expensesHandler.CreatePolicy)
			protected.GET("/expenses/policies", expensesHandler.ListPolicies)
			protected.GET("/expenses/policies/:id", expensesHandler.GetPolicy)
			protected.PUT("/expenses/policies/:id", expensesHandler.UpdatePolicy)
			protected.POST("/expenses/policies/:id/rules", expensesHandler.CreateRule)
			protected.GET("/expenses/policies/:id/rules", expensesHandler.ListRules)
			protected.PUT("/expenses/policies/:policyId/rules/:ruleId", expensesHandler.UpdateRule)
			protected.DELETE("/expenses/policies/:policyId/rules/:ruleId", expensesHandler.DeleteRule)
			protected.POST("/expenses/policies/evaluate", expensesHandler.EvaluateExpense)

			// Approvals
			protected.GET("/expenses/approvals/pending", expensesHandler.GetPendingApprovals)
			protected.POST("/expenses/approvals/:entityType/:entityId/approve", expensesHandler.ApproveEntity)
			protected.POST("/expenses/approvals/:entityType/:entityId/reject", expensesHandler.RejectEntity)
			protected.POST("/expenses/approvals/:entityType/:entityId/observe", expensesHandler.ObserveEntity)

			// Budgets
			protected.POST("/expenses/budgets", expensesHandler.CreateBudget)
			protected.GET("/expenses/budgets", expensesHandler.ListBudgets)
			protected.GET("/expenses/budgets/:id", expensesHandler.GetBudget)
			protected.POST("/expenses/budgets/check", expensesHandler.CheckAvailability)

			// Exchange rates
			protected.POST("/expenses/exchange-rates", expensesHandler.CreateRate)
			protected.GET("/expenses/exchange-rates/latest", expensesHandler.GetLatestRate)
			protected.POST("/expenses/exchange-rates/convert", expensesHandler.Convert)

			// Daily allowances
			protected.POST("/expenses/allowances", expensesHandler.CreateAllowanceRule)
			protected.GET("/expenses/allowances", expensesHandler.ListAllowanceRules)
			protected.GET("/expenses/allowances/:id", expensesHandler.GetAllowanceRule)
			protected.PUT("/expenses/allowances/:id", expensesHandler.UpdateAllowanceRule)
			protected.POST("/expenses/allowances/calculate", expensesHandler.CalculateAllowance)

			// Self-service
			protected.GET("/me/expenses", expensesHandler.MyExpenses)
			protected.POST("/me/expenses", expensesHandler.MyExpenseCreate)
			protected.POST("/me/expenses/:id/submit", expensesHandler.MyExpenseSubmit)
			protected.GET("/me/travels", expensesHandler.MyTravels)
			protected.POST("/me/travels", expensesHandler.MyTravelCreate)
			protected.POST("/me/travels/:id/request", expensesHandler.MyTravelRequest)
			protected.GET("/me/expense-reports", expensesHandler.MyReports)
			protected.POST("/me/expense-reports", expensesHandler.MyReportCreate)
			protected.POST("/me/expense-reports/:id/submit", expensesHandler.MyReportSubmit)
			protected.GET("/me/advances", expensesHandler.MyAdvances)
			protected.POST("/me/advances", expensesHandler.MyAdvanceCreate)
			protected.GET("/me/reimbursements", expensesHandler.MyReimbursements)

			// Performance — FASE 24 (DDD)
			perfHandler.RegisterRoutes(protected)

			// Recruitment (ATS)
			recHandler.RegisterRoutes(protected)

			protected.GET("/me/career/applications", recHandler.MyApplications)
			protected.GET("/me/career/referrals", recHandler.MyReferrals)
			protected.POST("/me/career/referrals", recHandler.CreateReferral)
			protected.GET("/me/career/interviews", recHandler.MyInterviews)

			// Onboarding & Offboarding (FASE 23)
			onbFase23Handler.RegisterRoutes(protected)

			// Training / LMS (FASE 17)
			protected.GET("/training/dashboard/stats", trainingHandler.DashboardStats)
			protected.GET("/training/employees/:employee_id/stats", trainingHandler.EmployeeStats)

			protected.GET("/training/categories", trainingHandler.ListCategories)
			protected.GET("/training/categories/:id", trainingHandler.GetCategory)
			protected.POST("/training/categories", trainingHandler.CreateCategory)
			protected.PUT("/training/categories/:id", trainingHandler.UpdateCategory)

			protected.GET("/training/courses", trainingHandler.ListCourses)
			protected.GET("/training/courses/:id", trainingHandler.GetCourse)
			protected.GET("/training/courses/:id/details", trainingHandler.GetCourseWithDetails)
			protected.POST("/training/courses", trainingHandler.CreateCourse)
			protected.PUT("/training/courses/:id", trainingHandler.UpdateCourse)
			protected.POST("/training/courses/:id/publish", trainingHandler.PublishCourse)

			protected.GET("/training/courses/:id/versions", trainingHandler.ListVersions)
			protected.POST("/training/courses/:id/versions", trainingHandler.CreateVersion)

			protected.GET("/training/versions/:version_id/contents", trainingHandler.ListContents)
			protected.POST("/training/versions/:version_id/contents", trainingHandler.CreateContent)
			protected.PUT("/training/contents/:id", trainingHandler.UpdateContent)
			protected.PUT("/training/versions/:version_id/contents/reorder", trainingHandler.ReorderContents)

			protected.GET("/training/offerings", trainingHandler.ListOfferings)
			protected.GET("/training/offerings/:id", trainingHandler.GetOffering)
			protected.POST("/training/offerings", trainingHandler.CreateOffering)
			protected.PUT("/training/offerings/:id", trainingHandler.UpdateOffering)

			protected.GET("/training/offerings/:id/sessions", trainingHandler.ListSessions)
			protected.POST("/training/offerings/:id/sessions", trainingHandler.CreateSession)

			protected.GET("/training/enrollments", trainingHandler.ListEnrollments)
			protected.GET("/training/enrollments/:id", trainingHandler.GetEnrollment)
			protected.POST("/training/offerings/:id/enroll", trainingHandler.Enroll)
			protected.POST("/training/enrollments/:id/complete", trainingHandler.CompleteEnrollment)

			protected.GET("/training/enrollments/:id/progress", trainingHandler.ListProgress)
			protected.GET("/training/enrollments/:id/progress/:content_id", trainingHandler.GetProgress)
			protected.PUT("/training/enrollments/:id/progress", trainingHandler.UpdateProgress)

			protected.POST("/training/assignments", trainingHandler.CreateAssignment)

			protected.GET("/training/assignment-rules", trainingHandler.ListAssignmentRules)
			protected.POST("/training/assignment-rules", trainingHandler.CreateAssignmentRule)

			protected.GET("/training/courses/:id/assessments", trainingHandler.ListAssessments)
			protected.GET("/training/assessments/:id", trainingHandler.GetAssessment)
			protected.POST("/training/courses/:id/assessments", trainingHandler.CreateAssessment)
			protected.PUT("/training/assessments/:id", trainingHandler.UpdateAssessment)

			protected.GET("/training/assessments/:id/questions", trainingHandler.GetQuestions)
			protected.POST("/training/assessments/:id/questions", trainingHandler.AddQuestion)

			protected.POST("/training/enrollments/:id/assessments/:assessment_id/start", trainingHandler.StartAttempt)
			protected.POST("/training/attempts/:id/submit", trainingHandler.SubmitAttempt)
			protected.GET("/training/attempts/:id", trainingHandler.GetAttempt)
			protected.GET("/training/enrollments/:id/attempts", trainingHandler.ListAttempts)

			protected.GET("/training/instructors", trainingHandler.ListInstructors)
			protected.GET("/training/instructors/:id", trainingHandler.GetInstructor)
			protected.POST("/training/instructors", trainingHandler.CreateInstructor)

			protected.GET("/training/providers", trainingHandler.ListProviders)
			protected.POST("/training/providers", trainingHandler.CreateProvider)

			protected.GET("/training/competencies", trainingHandler.ListCompetencies)
			protected.GET("/training/competencies/:id", trainingHandler.GetCompetency)
			protected.POST("/training/competencies", trainingHandler.CreateCompetency)
			protected.POST("/training/employees/:employee_id/competencies/:competency_id", trainingHandler.AssignCompetency)
			protected.GET("/training/employees/:employee_id/competencies", trainingHandler.GetEmployeeCompetencies)

			protected.POST("/training/courses/:id/competencies", trainingHandler.AddCourseCompetency)
			protected.GET("/training/courses/:id/competencies", trainingHandler.ListCourseCompetencies)

			protected.GET("/training/training-needs", trainingHandler.ListTrainingNeeds)
			protected.POST("/training/training-needs", trainingHandler.CreateTrainingNeed)

			protected.GET("/training/plans", trainingHandler.ListPlans)
			protected.POST("/training/plans", trainingHandler.CreatePlan)

			protected.GET("/training/learning-paths", trainingHandler.ListLearningPaths)
			protected.POST("/training/learning-paths", trainingHandler.CreateLearningPath)

			protected.GET("/training/enrollments/:id/feedback", trainingHandler.GetFeedbackByEnrollment)
			protected.POST("/training/enrollments/:id/feedback", trainingHandler.CreateFeedback)

			protected.GET("/training/enrollments/:id/attendance", trainingHandler.ListAttendance)
			protected.GET("/training/enrollments/:id/attendance/:session_id", trainingHandler.GetAttendance)
			protected.POST("/training/enrollments/:id/sessions/:session_id/attendance", trainingHandler.CreateAttendance)

			protected.GET("/training/employees/:employee_id/certificates", trainingHandler.ListCertificates)

			protected.POST("/training/ai/recommendations", trainingHandler.GenerateRecommendations)
			protected.GET("/training/employees/:employee_id/recommendations", trainingHandler.GetRecommendations)

			// Compensation (FASE 18)
			protected.GET("/compensation/dashboard", compHandler.DashboardStats)

			protected.GET("/compensation/structures", compHandler.ListStructures)
			protected.POST("/compensation/structures", compHandler.CreateStructure)
			protected.GET("/compensation/structures/:id", compHandler.GetStructure)
			protected.PUT("/compensation/structures/:id", compHandler.UpdateStructure)

			protected.GET("/compensation/structures/:id/grades", compHandler.ListGrades)
			protected.POST("/compensation/structures/:id/grades", compHandler.CreateGrade)
			protected.PUT("/compensation/grades/:id", compHandler.UpdateGrade)

			protected.GET("/compensation/structures/:id/bands", compHandler.ListBands)
			protected.POST("/compensation/structures/:id/bands", compHandler.CreateBand)
			protected.GET("/compensation/bands/:id", compHandler.GetBand)
			protected.PUT("/compensation/bands/:id", compHandler.UpdateBand)

			protected.POST("/compensation/positions/:position_id/band", compHandler.AssignPositionBand)
			protected.GET("/compensation/positions/:position_id/band", compHandler.GetPositionBand)

			protected.POST("/compensation/employees/:employee_id/compensation", compHandler.SetEmployeeCompensation)
			protected.GET("/compensation/employees/:employee_id/compensation", compHandler.GetEmployeeCompensation)
			protected.GET("/compensation/employees", compHandler.ListEmployeeCompensations)
			protected.GET("/compensation/employees/:employee_id/history", compHandler.GetHistory)

			protected.POST("/compensation/components", compHandler.CreateComponent)
			protected.GET("/compensation/components", compHandler.ListComponents)
			protected.POST("/compensation/employees/:employee_id/components", compHandler.AssignComponent)
			protected.GET("/compensation/employees/:employee_id/components", compHandler.ListEmployeeComponents)

			protected.POST("/compensation/adjustments", compHandler.CreateAdjustment)
			protected.GET("/compensation/adjustments", compHandler.ListAdjustments)
			protected.GET("/compensation/adjustments/:id", compHandler.GetAdjustment)
			protected.POST("/compensation/adjustments/:id/approve", compHandler.ApproveAdjustment)
			protected.POST("/compensation/adjustments/:id/reject", compHandler.RejectAdjustment)
			protected.POST("/compensation/adjustments/:id/apply", compHandler.ApplyAdjustment)

			protected.POST("/compensation/proposals", compHandler.CreateProposal)
			protected.POST("/compensation/proposals/:id/submit", compHandler.SubmitProposal)
			protected.POST("/compensation/proposals/:id/approve", compHandler.ApproveProposal)
			protected.POST("/compensation/proposals/:id/reject", compHandler.RejectProposal)
			protected.GET("/compensation/proposals", compHandler.ListProposals)

			protected.POST("/compensation/bonus-plans", compHandler.CreateBonusPlan)
			protected.GET("/compensation/bonus-plans", compHandler.ListBonusPlans)
			protected.POST("/compensation/bonuses", compHandler.CreateBonus)
			protected.GET("/compensation/bonuses", compHandler.ListBonuses)
			protected.GET("/compensation/bonuses/:id", compHandler.GetBonus)
			protected.POST("/compensation/bonuses/:id/approve", compHandler.ApproveBonus)
			protected.POST("/compensation/bonuses/:id/reject", compHandler.RejectBonus)

			protected.POST("/compensation/benefits", compHandler.CreateBenefit)
			protected.GET("/compensation/benefits", compHandler.ListBenefits)
			protected.GET("/compensation/benefits/:id", compHandler.GetBenefit)
			protected.PUT("/compensation/benefits/:id", compHandler.UpdateBenefit)
			protected.POST("/compensation/employees/:employee_id/benefits", compHandler.AssignBenefit)
			protected.GET("/compensation/employees/:employee_id/benefits", compHandler.ListEmployeeBenefits)
			protected.DELETE("/compensation/employee-benefits/:id", compHandler.RemoveEmployeeBenefit)

			protected.POST("/compensation/reviews", compHandler.CreateReview)
			protected.GET("/compensation/reviews", compHandler.ListReviews)
			protected.GET("/compensation/reviews/:id", compHandler.GetReview)
			protected.POST("/compensation/reviews/:id/open", compHandler.OpenReview)
			protected.POST("/compensation/reviews/:id/close", compHandler.CloseReview)

			protected.POST("/compensation/budgets", compHandler.CreateBudget)
			protected.GET("/compensation/budgets", compHandler.ListBudgets)

			protected.GET("/compensation/reports/bands/:band_id/analysis", compHandler.BandAnalysisReport)
			protected.GET("/compensation/reports/employees/:employee_id/compa-ratio", compHandler.CompaRatioReport)
			protected.GET("/compensation/reports/employees/:employee_id/range-penetration", compHandler.RangePenetrationReport)
			protected.GET("/compensation/reports/employees/:employee_id/total-compensation", compHandler.TotalCompensationReport)
			protected.GET("/compensation/reports/equity", compHandler.EquityReport)

			protected.POST("/compensation/ai/recommendations", compHandler.GenerateAIRecommendation)

			protected.GET("/me/compensation", compHandler.MyCompensation)
			protected.GET("/me/compensation/history", compHandler.MyHistory)
			protected.GET("/me/compensation/benefits", compHandler.MyBenefits)
			protected.GET("/me/compensation/bonuses", compHandler.MyBonuses)

			protected.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"message": "pong"}})
			})
		}
	}

	// Serve frontend SPA (build from frontend/dist)
	frontendDir := "./frontend/dist"
	if _, err := http.Dir(frontendDir).Open("."); err == nil {
		router.Static("/assets", frontendDir+"/assets")
		router.NoRoute(func(c *gin.Context) {
			c.File(frontendDir + "/index.html")
		})
	}

	return router
}
