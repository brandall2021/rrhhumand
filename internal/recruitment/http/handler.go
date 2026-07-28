package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/application"
)

type Handler struct {
	RequisitionSvc *application.RequisitionService
	PositionSvc    *application.PositionService
	PostingSvc     *application.PostingService
	CandidateSvc   *application.CandidateService
	ApplicationSvc *application.ApplicationService
	InterviewSvc   *application.InterviewService
	AssessmentSvc  *application.AssessmentService
	OfferSvc       *application.OfferService
	HiringSvc      *application.HiringService
	WorkflowSvc    *application.WorkflowService
	ScoringSvc     *application.ScoringService
	EmailSvc       *application.EmailService
	DashboardSvc   *application.DashboardService
	SettingsSvc    *application.SettingsService
}

func NewHandler(
	rs *application.RequisitionService,
	ps *application.PositionService,
	pts *application.PostingService,
	cs *application.CandidateService,
	as *application.ApplicationService,
	is *application.InterviewService,
	ass *application.AssessmentService,
	os *application.OfferService,
	hs *application.HiringService,
	ws *application.WorkflowService,
	scs *application.ScoringService,
	es *application.EmailService,
	ds *application.DashboardService,
	ss *application.SettingsService,
) *Handler {
	return &Handler{
		RequisitionSvc: rs,
		PositionSvc:    ps,
		PostingSvc:     pts,
		CandidateSvc:   cs,
		ApplicationSvc: as,
		InterviewSvc:   is,
		AssessmentSvc:  ass,
		OfferSvc:       os,
		HiringSvc:      hs,
		WorkflowSvc:    ws,
		ScoringSvc:     scs,
		EmailSvc:       es,
		DashboardSvc:   ds,
		SettingsSvc:    ss,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rec := rg.Group("/recruitment")
	{
		rec.POST("/requisitions", h.CreateRequisition)
		rec.GET("/requisitions", h.ListRequisitions)
		rec.GET("/requisitions/:id", h.GetRequisition)
		rec.PUT("/requisitions/:id", h.UpdateRequisition)
		rec.POST("/requisitions/:id/submit", h.SubmitRequisition)
		rec.POST("/requisitions/:id/approve", h.ApproveRequisition)
		rec.POST("/requisitions/:id/open", h.OpenRequisition)
		rec.POST("/requisitions/:id/close", h.CloseRequisition)
		rec.POST("/requisitions/:id/cancel", h.CancelRequisition)
		rec.GET("/requisitions/:id/skills", h.ListRequisitionSkills)
		rec.POST("/requisitions/:id/skills", h.AddRequisitionSkill)
		rec.DELETE("/requisitions/:id/skills/:skillId", h.RemoveRequisitionSkill)

		rec.POST("/positions", h.CreatePosition)
		rec.GET("/positions", h.ListPositions)
		rec.GET("/positions/:id", h.GetPosition)
		rec.PUT("/positions/:id", h.UpdatePosition)
		rec.POST("/positions/:id/close", h.ClosePosition)
		rec.GET("/positions/:id/skills", h.ListPositionSkills)
		rec.POST("/positions/:id/skills", h.AddPositionSkill)
		rec.DELETE("/positions/:id/skills/:skillId", h.RemovePositionSkill)

		rec.POST("/postings", h.CreatePosting)
		rec.GET("/postings", h.ListPostings)
		rec.GET("/postings/:id", h.GetPosting)
		rec.PUT("/postings/:id", h.UpdatePosting)
		rec.POST("/postings/:id/publish", h.PublishPosting)
		rec.POST("/postings/:id/close", h.ClosePosting)
		rec.GET("/postings/:id/questions", h.ListScreeningQuestions)
		rec.POST("/postings/:id/questions", h.AddScreeningQuestion)
		rec.PUT("/postings/questions/:questionId", h.UpdateScreeningQuestion)
		rec.DELETE("/postings/questions/:questionId", h.DeleteScreeningQuestion)

		rec.POST("/candidates", h.CreateCandidate)
		rec.GET("/candidates", h.ListCandidates)
		rec.GET("/candidates/:id", h.GetCandidate)
		rec.PUT("/candidates/:id", h.UpdateCandidate)
		rec.POST("/candidates/:id/blacklist", h.BlacklistCandidate)
		rec.POST("/candidates/:id/unblacklist", h.UnblacklistCandidate)
		rec.GET("/candidates/search", h.SearchCandidates)
		rec.GET("/candidates/:id/education", h.ListCandidateEducation)
		rec.POST("/candidates/:id/education", h.AddCandidateEducation)
		rec.PUT("/candidates/education/:eduId", h.UpdateCandidateEducation)
		rec.DELETE("/candidates/education/:eduId", h.DeleteCandidateEducation)
		rec.GET("/candidates/:id/experience", h.ListCandidateExperience)
		rec.POST("/candidates/:id/experience", h.AddCandidateExperience)
		rec.PUT("/candidates/experience/:expId", h.UpdateCandidateExperience)
		rec.DELETE("/candidates/experience/:expId", h.DeleteCandidateExperience)
		rec.GET("/candidates/:id/skills", h.ListCandidateSkills)
		rec.POST("/candidates/:id/skills", h.AddCandidateSkill)
		rec.PUT("/candidates/skills/:skillId", h.UpdateCandidateSkill)
		rec.DELETE("/candidates/skills/:skillId", h.DeleteCandidateSkill)
		rec.GET("/candidates/:id/certifications", h.ListCandidateCertifications)
		rec.POST("/candidates/:id/certifications", h.AddCandidateCertification)
		rec.DELETE("/candidates/certifications/:certId", h.DeleteCandidateCertification)
		rec.GET("/candidates/:id/languages", h.ListCandidateLanguages)
		rec.POST("/candidates/:id/languages", h.AddCandidateLanguage)
		rec.PUT("/candidates/languages/:langId", h.UpdateCandidateLanguage)
		rec.DELETE("/candidates/languages/:langId", h.DeleteCandidateLanguage)
		rec.GET("/candidates/:id/documents", h.ListCandidateDocuments)
		rec.POST("/candidates/:id/documents", h.AddCandidateDocument)

		rec.POST("/applications", h.CreateApplication)
		rec.GET("/applications", h.ListApplications)
		rec.GET("/applications/:id", h.GetApplication)
		rec.POST("/applications/:id/stage", h.MoveStage)
		rec.POST("/applications/:id/reject", h.RejectApplication)
		rec.POST("/applications/:id/withdraw", h.WithdrawApplication)
		rec.GET("/applications/:id/history", h.GetStageHistory)
		rec.GET("/applications/:id/notes", h.ListApplicationNotes)
		rec.POST("/applications/:id/notes", h.AddApplicationNote)
		rec.PUT("/applications/notes/:noteId", h.UpdateApplicationNote)
		rec.GET("/applications/:id/ratings", h.ListApplicationRatings)
		rec.POST("/applications/:id/ratings", h.AddApplicationRating)

		rec.POST("/interviews", h.CreateInterview)
		rec.GET("/interviews", h.ListInterviews)
		rec.GET("/interviews/:id", h.GetInterview)
		rec.PUT("/interviews/:id", h.UpdateInterview)
		rec.POST("/interviews/:id/cancel", h.CancelInterview)
		rec.POST("/interviews/:id/complete", h.CompleteInterview)
		rec.POST("/interviews/:id/panel", h.AddPanelMember)
		rec.DELETE("/interviews/panel/:panelId", h.RemovePanelMember)
		rec.GET("/interviews/:id/panel", h.ListPanelMembers)
		rec.POST("/interviews/:id/feedback", h.SubmitInterviewFeedback)
		rec.GET("/interviews/:id/feedback", h.ListInterviewFeedback)

		rec.POST("/assessments", h.CreateAssessment)
		rec.GET("/assessments", h.ListAssessments)
		rec.GET("/assessments/:id", h.GetAssessment)
		rec.PUT("/assessments/:id", h.UpdateAssessment)
		rec.POST("/assessments/:id/send", h.SendAssessment)
		rec.POST("/assessments/:id/score", h.ScoreAssessment)
		rec.POST("/assessments/:id/cancel", h.CancelAssessment)
		rec.GET("/assessments/:id/sections", h.ListAssessmentSections)
		rec.POST("/assessments/:id/sections", h.AddAssessmentSection)
		rec.GET("/assessments/:id/results", h.ListAssessmentResults)
		rec.POST("/assessments/:id/results", h.AddAssessmentResult)

		rec.POST("/offers", h.CreateOffer)
		rec.GET("/offers", h.ListOffers)
		rec.GET("/offers/:id", h.GetOffer)
		rec.PUT("/offers/:id", h.UpdateOffer)
		rec.POST("/offers/:id/submit", h.SubmitOfferForApproval)
		rec.POST("/offers/:id/approve", h.ApproveOffer)
		rec.POST("/offers/:id/send", h.SendOffer)
		rec.POST("/offers/:id/accept", h.AcceptOffer)
		rec.POST("/offers/:id/reject", h.RejectOffer)
		rec.POST("/offers/:id/withdraw", h.WithdrawOffer)
		rec.GET("/offers/:id/negotiations", h.ListOfferNegotiations)
		rec.POST("/offers/:id/negotiations", h.AddOfferNegotiation)
		rec.PUT("/offers/negotiations/:negId", h.UpdateOfferNegotiation)
		rec.GET("/offers/:id/documents", h.ListOfferDocuments)
		rec.POST("/offers/:id/documents", h.AddOfferDocument)

		rec.POST("/hiring-processes", h.CreateHiringProcess)
		rec.GET("/hiring-processes", h.ListHiringProcesses)
		rec.GET("/hiring-processes/:id", h.GetHiringProcess)
		rec.POST("/hiring-processes/:id/background-check", h.UpdateBackgroundCheck)
		rec.POST("/hiring-processes/:id/medical-check", h.UpdateMedicalCheck)
		rec.POST("/hiring-processes/:id/document-verification", h.UpdateDocVerification)
		rec.POST("/hiring-processes/:id/complete", h.CompleteHiringProcess)
		rec.POST("/hiring-processes/:id/cancel", h.CancelHiringProcess)
		rec.GET("/hiring-processes/:id/tasks", h.ListHiringTasks)
		rec.POST("/hiring-processes/:id/tasks", h.AddHiringTask)
		rec.POST("/hiring-processes/tasks/:taskId/complete", h.CompleteHiringTask)

		rec.POST("/workflows", h.CreateWorkflow)
		rec.GET("/workflows", h.ListWorkflows)
		rec.GET("/workflows/:id", h.GetWorkflow)
		rec.PUT("/workflows/:id", h.UpdateWorkflow)
		rec.POST("/workflows/:id/activate", h.ActivateWorkflow)
		rec.POST("/workflows/:id/deactivate", h.DeactivateWorkflow)
		rec.GET("/workflows/:id/stages", h.ListWorkflowStages)
		rec.POST("/workflows/:id/stages", h.AddWorkflowStage)
		rec.DELETE("/workflows/stages/:stageId", h.RemoveWorkflowStage)
		rec.POST("/workflows/:id/stages/reorder", h.ReorderWorkflowStages)
		rec.GET("/workflows/:id/rules", h.ListWorkflowRules)
		rec.POST("/workflows/:id/rules", h.AddWorkflowRule)
		rec.PUT("/workflows/rules/:ruleId", h.UpdateWorkflowRule)
		rec.DELETE("/workflows/rules/:ruleId", h.DeleteWorkflowRule)

		rec.POST("/scoring/models", h.CreateScoringModel)
		rec.GET("/scoring/models", h.ListScoringModels)
		rec.GET("/scoring/models/:id", h.GetScoringModel)
		rec.PUT("/scoring/models/:id", h.UpdateScoringModel)
		rec.DELETE("/scoring/models/:id", h.DeleteScoringModel)
		rec.GET("/scoring/models/:id/criteria", h.ListScoringCriteria)
		rec.POST("/scoring/models/:id/criteria", h.AddScoringCriterion)
		rec.PUT("/scoring/models/criteria/:criterionId", h.UpdateScoringCriterion)
		rec.DELETE("/scoring/models/criteria/:criterionId", h.DeleteScoringCriterion)
		rec.POST("/scoring/match", h.ScoreCandidate)
		rec.GET("/scoring/matches/:candidateId/:positionId", h.GetMatchingResult)

		rec.POST("/email/templates", h.CreateEmailTemplate)
		rec.GET("/email/templates", h.ListEmailTemplates)
		rec.GET("/email/templates/:id", h.GetEmailTemplate)
		rec.PUT("/email/templates/:id", h.UpdateEmailTemplate)
		rec.DELETE("/email/templates/:id", h.DeleteEmailTemplate)
		rec.POST("/email/send", h.SendEmail)
		rec.GET("/email/log", h.ListSentEmails)

		rec.GET("/dashboard", h.GetDashboard)
		rec.GET("/dashboard/funnel", h.GetFunnel)
		rec.GET("/dashboard/time-to-hire", h.GetTimeToHire)

		rec.POST("/settings/sources", h.CreateSource)
		rec.GET("/settings/sources", h.ListSources)
		rec.PUT("/settings/sources/:id", h.UpdateSource)
		rec.POST("/settings/stages", h.CreateStage)
		rec.GET("/settings/stages", h.ListStages)
		rec.PUT("/settings/stages/:id", h.UpdateStage)
		rec.POST("/settings/stages/reorder", h.ReorderStages)
		rec.POST("/settings/transitions", h.CreateTransition)
		rec.GET("/settings/transitions", h.ListTransitions)
		rec.DELETE("/settings/transitions/:id", h.DeleteTransition)
		rec.POST("/settings/rejection-reasons", h.CreateRejectionReason)
		rec.GET("/settings/rejection-reasons", h.ListRejectionReasons)
		rec.PUT("/settings/rejection-reasons/:id", h.UpdateRejectionReason)
	}
}

func (h *Handler) RegisterPublicRoutes(rg *gin.RouterGroup) {
	rg.GET("/public/jobs", h.ListPublicPostings)
	rg.GET("/public/jobs/:id", h.GetPublicPosting)
	rg.POST("/public/jobs/:id/apply", h.PublicApply)
}
