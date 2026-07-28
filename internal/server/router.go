package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rrhhumand/api/internal/attendance"
	"github.com/rrhhumand/api/internal/auth"
	"github.com/rrhhumand/api/internal/branches"
	"github.com/rrhhumand/api/internal/companies"
	"github.com/rrhhumand/api/internal/departments"
	"github.com/rrhhumand/api/internal/document_categories"
	"github.com/rrhhumand/api/internal/documents"
	"github.com/rrhhumand/api/internal/employees"
	"github.com/rrhhumand/api/internal/feed"
	"github.com/rrhhumand/api/internal/handlers"
	"github.com/rrhhumand/api/internal/leave"
	"github.com/rrhhumand/api/internal/middleware"
	"github.com/rrhhumand/api/internal/organization"
	"github.com/rrhhumand/api/internal/overtime"
	"github.com/rrhhumand/api/internal/payroll"
	"github.com/rrhhumand/api/internal/performance"
	"github.com/rrhhumand/api/internal/positions"
	"github.com/rrhhumand/api/internal/profile"
	"github.com/rrhhumand/api/internal/recruitment"
	"github.com/rrhhumand/api/internal/roles"
	"github.com/rrhhumand/api/internal/scheduling"
	"github.com/rrhhumand/api/internal/surveys"
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
	perfHandler *performance.Handler,
	recHandler *recruitment.Handler,
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

		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
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

			// Payroll
			protected.POST("/payroll/periods", payHandler.CreatePeriod)
			protected.GET("/payroll/periods", payHandler.ListPeriods)
			protected.GET("/payroll/periods/:id", payHandler.GetPeriod)
			protected.PUT("/payroll/periods/:id", payHandler.UpdatePeriod)
			protected.POST("/payroll/periods/:id/calculate", payHandler.CalculatePeriod)
			protected.GET("/payroll/periods/:id/review", payHandler.GetReview)
			protected.POST("/payroll/periods/:id/approve", payHandler.ApprovePeriod)
			protected.POST("/payroll/periods/:id/close", payHandler.ClosePeriod)
			protected.GET("/payroll/periods/:id/dashboard", payHandler.GetDashboard)
			protected.POST("/payroll/periods/:id/adjustments", payHandler.CreateAdjustment)
			protected.GET("/payroll/periods/:id/snapshot", payHandler.GetSnapshot)
			protected.GET("/payroll/periods/:id/ledger", payHandler.GetDashboard)

			protected.POST("/payroll/concepts", payHandler.CreateConcept)
			protected.GET("/payroll/concepts", payHandler.ListConcepts)

			protected.POST("/payroll/compensation", payHandler.SetCompensation)
			protected.GET("/employees/:id/compensation", payHandler.GetCompensation)
			protected.GET("/employees/:id/compensation/history", payHandler.GetCompensationHistory)

			protected.POST("/payroll/benefits", payHandler.CreateBenefit)
			protected.GET("/payroll/benefits", payHandler.ListBenefits)
			protected.POST("/payroll/benefits/assign", payHandler.AssignBenefit)
			protected.GET("/employees/:id/benefits", payHandler.GetEmployeeBenefitsView)

			protected.POST("/payroll/bonuses", payHandler.CreateBonus)
			protected.GET("/payroll/bonuses", payHandler.ListBonuses)
			protected.POST("/payroll/bonuses/:id/approve", payHandler.ApproveBonus)

			protected.POST("/payroll/advances", payHandler.CreateAdvance)
			protected.GET("/payroll/advances", payHandler.ListAdvances)
			protected.POST("/payroll/advances/:id/approve", payHandler.ApproveAdvance)

			protected.POST("/payroll/deductions", payHandler.CreateDeduction)

			protected.GET("/employees/:id/payroll", payHandler.GetEmployeePayroll)

			// Performance
			protected.POST("/performance/cycles", perfHandler.CreateCycle)
			protected.GET("/performance/cycles", perfHandler.ListCycles)
			protected.GET("/performance/cycles/:id", perfHandler.GetCycle)
			protected.PUT("/performance/cycles/:id", perfHandler.UpdateCycle)
			protected.POST("/performance/cycles/:id/open", perfHandler.OpenCycle)
			protected.POST("/performance/cycles/:id/close", perfHandler.CloseCycle)

			protected.POST("/performance/templates", perfHandler.CreateTemplate)
			protected.GET("/performance/templates", perfHandler.ListTemplates)

			protected.POST("/performance/scales", perfHandler.CreateScale)
			protected.GET("/performance/scales", perfHandler.ListScales)

			protected.POST("/performance/competencies", perfHandler.CreateCompetency)
			protected.GET("/performance/competencies", perfHandler.ListCompetencies)
			protected.PUT("/performance/competencies/:id", perfHandler.UpdateCompetency)

			protected.POST("/performance/objectives", perfHandler.CreateObjective)
			protected.GET("/performance/objectives", perfHandler.ListObjectives)
			protected.GET("/performance/objectives/:id", perfHandler.GetObjective)
			protected.PUT("/performance/objectives/:id", perfHandler.UpdateObjective)
			protected.DELETE("/performance/objectives/:id", perfHandler.DeleteObjective)
			protected.POST("/performance/objectives/:id/progress", perfHandler.UpdateObjectiveProgress)

			protected.POST("/performance/kpis", perfHandler.CreateKPI)
			protected.GET("/performance/kpis", perfHandler.ListKPIs)
			protected.PUT("/performance/kpis/:id", perfHandler.UpdateKPI)

			protected.POST("/performance/evaluators", perfHandler.AssignEvaluators)
			protected.GET("/performance/evaluators", perfHandler.ListEvaluators)

			protected.POST("/performance/evaluations", perfHandler.CreateEvaluation)
			protected.GET("/performance/evaluations", perfHandler.ListEvaluations)
			protected.GET("/performance/evaluations/:id", perfHandler.GetEvaluation)
			protected.PUT("/performance/evaluations/:id", perfHandler.UpdateEvaluation)
			protected.POST("/performance/evaluations/:id/submit", perfHandler.SubmitEvaluation)
			protected.POST("/performance/evaluations/:id/reopen", perfHandler.ReopenEvaluation)
			protected.POST("/performance/evaluations/:id/approve", perfHandler.ApproveEvaluation)

			protected.POST("/performance/evaluations/:id/answers", perfHandler.CreateAnswer)
			protected.GET("/performance/evaluations/:id/answers", perfHandler.ListAnswers)

			protected.POST("/performance/evaluations/:id/evidence", perfHandler.CreateEvidence)
			protected.GET("/performance/evaluations/:id/evidence", perfHandler.ListEvidence)

			protected.POST("/performance/feedback", perfHandler.CreateFeedback)
			protected.GET("/performance/feedback", perfHandler.ListFeedback)
			protected.GET("/performance/feedback/:id", perfHandler.GetFeedback)

			protected.POST("/performance/results/calculate", perfHandler.CalculateResult)
			protected.GET("/performance/results", perfHandler.ListResults)
			protected.GET("/performance/results/:id", perfHandler.GetResult)

			protected.GET("/performance/scoring-rules", perfHandler.GetScoringRules)
			protected.PUT("/performance/scoring-rules", perfHandler.UpdateScoringRules)

			protected.POST("/performance/improvement-plans", perfHandler.CreateImprovementPlan)
			protected.GET("/performance/improvement-plans", perfHandler.ListImprovementPlans)
			protected.GET("/performance/improvement-plans/:id", perfHandler.GetImprovementPlan)
			protected.PUT("/performance/improvement-plans/:id", perfHandler.UpdateImprovementPlan)
			protected.POST("/performance/improvement-plans/:id/complete", perfHandler.CompleteImprovementPlan)

			protected.POST("/performance/development-plans", perfHandler.CreateDevelopmentPlan)
			protected.GET("/performance/development-plans", perfHandler.ListDevelopmentPlans)
			protected.GET("/performance/development-plans/:id", perfHandler.GetDevelopmentPlan)
			protected.PUT("/performance/development-plans/:id", perfHandler.UpdateDevelopmentPlan)

			protected.GET("/performance/dashboard", perfHandler.GetDashboard)

			// Recruitment (ATS)
			protected.POST("/recruitment/requisitions", recHandler.CreateRequisition)
			protected.GET("/recruitment/requisitions", recHandler.ListRequisitions)
			protected.GET("/recruitment/requisitions/:id", recHandler.GetRequisition)
			protected.PUT("/recruitment/requisitions/:id", recHandler.UpdateRequisition)
			protected.POST("/recruitment/requisitions/:id/submit", recHandler.SubmitRequisition)
			protected.POST("/recruitment/requisitions/:id/approve", recHandler.ApproveRequisition)
			protected.POST("/recruitment/requisitions/:id/open", recHandler.OpenRequisition)
			protected.POST("/recruitment/requisitions/:id/close", recHandler.CloseRequisition)

			protected.POST("/recruitment/jobs", recHandler.CreatePosting)
			protected.GET("/recruitment/jobs", recHandler.ListPostings)
			protected.GET("/recruitment/jobs/:id", recHandler.GetPosting)
			protected.POST("/recruitment/jobs/:id/publish", recHandler.PublishPosting)
			protected.POST("/recruitment/jobs/:id/close", recHandler.ClosePosting)

			protected.POST("/recruitment/candidates", recHandler.CreateCandidate)
			protected.GET("/recruitment/candidates", recHandler.ListCandidates)
			protected.GET("/recruitment/candidates/:id", recHandler.GetCandidate)
			protected.PUT("/recruitment/candidates/:id", recHandler.UpdateCandidate)

			protected.POST("/recruitment/applications", recHandler.CreateApplication)
			protected.GET("/recruitment/applications/:id", recHandler.GetApplication)
			protected.GET("/recruitment/applications", recHandler.ListApplications)
			protected.POST("/recruitment/applications/:id/stage", recHandler.MoveStage)
			protected.POST("/recruitment/applications/:id/reject", recHandler.RejectApplication)
			protected.POST("/recruitment/applications/:id/withdraw", recHandler.WithdrawApplication)
			protected.GET("/recruitment/applications/:id/history", recHandler.GetStageHistory)

			protected.POST("/recruitment/interviews", recHandler.CreateInterview)
			protected.GET("/recruitment/interviews", recHandler.ListInterviews)
			protected.GET("/recruitment/interviews/:id", recHandler.GetInterview)
			protected.PUT("/recruitment/interviews/:id", recHandler.UpdateInterview)
			protected.POST("/recruitment/interviews/:id/feedback", recHandler.CreateInterviewFeedback)
			protected.GET("/recruitment/interviews/:id/feedback", recHandler.ListInterviewFeedback)

			protected.POST("/recruitment/assessments", recHandler.CreateAssessment)
			protected.GET("/recruitment/assessments/:id", recHandler.ListAssessments)

			protected.POST("/recruitment/screening", recHandler.CreateScreeningQuestion)
			protected.GET("/recruitment/screening/:id", recHandler.ListScreeningQuestions)

			protected.POST("/recruitment/offers", recHandler.CreateOffer)
			protected.GET("/recruitment/offers/:id", recHandler.GetOffer)
			protected.POST("/recruitment/offers/:id/send", recHandler.SendOffer)
			protected.POST("/recruitment/offers/:id/accept", recHandler.AcceptOffer)
			protected.POST("/recruitment/offers/:id/reject", recHandler.RejectOffer)

			protected.POST("/recruitment/referrals", recHandler.CreateReferral)
			protected.GET("/recruitment/referrals", recHandler.ListReferrals)

			protected.POST("/recruitment/applications/:id/hire", recHandler.HireCandidate)

			protected.GET("/recruitment/dashboard", recHandler.GetDashboard)

			protected.GET("/ping", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"message": "pong"}})
			})
		}
	}

	return router
}
