package performancehttp

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/performance/domain"
	"github.com/rrhhumand/api/internal/performance/repository"
	"github.com/rrhhumand/api/internal/performance/application/service"
	"github.com/rrhhumand/api/internal/performance/application/scoring"
	"github.com/rrhhumand/api/internal/performance/application/workflow"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type Handler struct {
	cycles       *service.CycleService
	objectives   *service.ObjectiveService
	evaluations  *service.EvaluationService
	feedback     *service.FeedbackService
	calibration  *service.CalibrationService
	plans        *service.PlanService

	cycleRepo    repository.CycleRepository
	templateRepo repository.TemplateRepository
	scaleRepo    repository.ScaleRepository
	compRepo     repository.CompetencyRepository
	objectiveRepo repository.ObjectiveRepository
	participantRepo repository.ParticipantRepository
	evalRepo     repository.EvaluationRepository
	reviewRepo   repository.ReviewRepository
	feedbackRepo repository.FeedbackRepository
	checkInRepo  repository.CheckInRepository
	calibRepo    repository.CalibrationRepository
	impRepo      repository.ImprovementPlanRepository
	devRepo      repository.DevelopmentPlanRepository
	evidenceRepo repository.EvidenceRepository
	resultRepo   repository.ResultRepository
	dashRepo     repository.DashboardRepository

	scorer *scoring.Scorer
}

func NewHandler(
	cycleRepo repository.CycleRepository,
	templateRepo repository.TemplateRepository,
	scaleRepo repository.ScaleRepository,
	compRepo repository.CompetencyRepository,
	objectiveRepo repository.ObjectiveRepository,
	participantRepo repository.ParticipantRepository,
	evalRepo repository.EvaluationRepository,
	reviewRepo repository.ReviewRepository,
	feedbackRepo repository.FeedbackRepository,
	checkInRepo repository.CheckInRepository,
	calibRepo repository.CalibrationRepository,
	impRepo repository.ImprovementPlanRepository,
	devRepo repository.DevelopmentPlanRepository,
	evidenceRepo repository.EvidenceRepository,
	resultRepo repository.ResultRepository,
	dashRepo repository.DashboardRepository,
	scorer *scoring.Scorer,
) *Handler {
	return &Handler{
		cycles:       service.NewCycleService(cycleRepo, templateRepo, participantRepo),
		objectives:   service.NewObjectiveService(objectiveRepo, cycleRepo),
		evaluations:  service.NewEvaluationService(evalRepo, participantRepo, cycleRepo, reviewRepo),
		feedback:     service.NewFeedbackService(feedbackRepo, checkInRepo),
		calibration:  service.NewCalibrationService(calibRepo, cycleRepo),
		plans:        service.NewPlanService(impRepo, devRepo),
		cycleRepo:    cycleRepo,
		templateRepo: templateRepo,
		scaleRepo:    scaleRepo,
		compRepo:     compRepo,
		objectiveRepo: objectiveRepo,
		participantRepo: participantRepo,
		evalRepo:     evalRepo,
		reviewRepo:   reviewRepo,
		feedbackRepo: feedbackRepo,
		checkInRepo:  checkInRepo,
		calibRepo:    calibRepo,
		impRepo:      impRepo,
		devRepo:      devRepo,
		evidenceRepo: evidenceRepo,
		resultRepo:   resultRepo,
		dashRepo:     dashRepo,
		scorer:       scorer,
	}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	perf := rg.Group("/performance")
	{
		// Cycles
		perf.POST("/cycles", h.CreateCycle)
		perf.GET("/cycles", h.ListCycles)
		perf.GET("/cycles/:id", h.GetCycle)
		perf.PUT("/cycles/:id", h.UpdateCycle)
		perf.POST("/cycles/:id/status", h.UpdateCycleStatus)

		// Templates
		perf.POST("/templates", h.CreateTemplate)
		perf.GET("/templates", h.ListTemplates)
		perf.GET("/templates/:id", h.GetTemplate)
		perf.PUT("/templates/:id", h.UpdateTemplate)
		perf.DELETE("/templates/:id", h.DeleteTemplate)

		// Rating Scales
		perf.POST("/scales", h.CreateScale)
		perf.GET("/scales", h.ListScales)
		perf.GET("/scales/:id", h.GetScale)
		perf.PUT("/scales/:id", h.UpdateScale)
		perf.DELETE("/scales/:id", h.DeleteScale)

		// Competencies
		perf.POST("/competencies", h.CreateCompetency)
		perf.GET("/competencies", h.ListCompetencies)
		perf.GET("/competencies/:id", h.GetCompetency)
		perf.PUT("/competencies/:id", h.UpdateCompetency)
		perf.DELETE("/competencies/:id", h.DeleteCompetency)
		perf.POST("/competencies/:id/levels", h.CreateCompetencyLevel)
		perf.POST("/position-competencies", h.UpsertPositionCompetency)
		perf.POST("/cycle-competencies", h.UpsertCycleCompetency)

		// Objectives
		perf.POST("/objectives", h.CreateObjective)
		perf.GET("/objectives", h.ListObjectives)
		perf.GET("/objectives/:id", h.GetObjective)
		perf.PUT("/objectives/:id", h.UpdateObjective)
		perf.DELETE("/objectives/:id", h.DeleteObjective)
		perf.POST("/objectives/:id/progress", h.UpdateObjectiveProgress)
		perf.POST("/objectives/:id/key-results", h.CreateKeyResult)

		// Participants
		perf.POST("/participants", h.AssignParticipants)
		perf.GET("/participants", h.ListParticipants)
		perf.DELETE("/participants/:id", h.RemoveParticipant)

		// Evaluations
		perf.POST("/evaluations", h.CreateEvaluation)
		perf.GET("/evaluations", h.ListEvaluations)
		perf.GET("/evaluations/:id", h.GetEvaluation)
		perf.PUT("/evaluations/:id", h.UpdateEvaluation)
		perf.POST("/evaluations/:id/submit", h.SubmitEvaluation)
		perf.POST("/evaluations/:id/approve", h.ApproveEvaluation)
		perf.POST("/evaluations/:id/reopen", h.ReopenEvaluation)
		perf.POST("/evaluations/:id/lock", h.LockEvaluation)

		// Evaluation answers
		perf.POST("/evaluations/:id/answers", h.CreateAnswers)
		perf.GET("/evaluations/:id/answers", h.ListAnswers)

		// Objective & Competency evaluations
		perf.POST("/evaluations/:id/objective-evaluations", h.CreateObjectiveEvaluation)
		perf.POST("/evaluations/:id/competency-evaluations", h.CreateCompetencyEvaluation)

		// Reviews
		perf.POST("/reviews", h.CreateReview)
		perf.GET("/reviews", h.ListReviews)
		perf.GET("/reviews/:id", h.GetReview)
		perf.PUT("/reviews/:id", h.UpdateReview)
		perf.POST("/reviews/:id/status", h.UpdateReviewStatus)

		// Feedback
		perf.POST("/feedback", h.CreateFeedback)
		perf.GET("/feedback", h.ListFeedback)
		perf.GET("/feedback/:id", h.GetFeedback)
		perf.DELETE("/feedback/:id", h.DeleteFeedback)

		// Recognitions
		perf.POST("/recognitions", h.CreateRecognition)
		perf.GET("/recognitions", h.ListRecognitions)

		// Check-ins
		perf.POST("/checkins", h.CreateCheckIn)
		perf.GET("/checkins", h.ListCheckIns)
		perf.GET("/checkins/:id", h.GetCheckIn)
		perf.PUT("/checkins/:id", h.UpdateCheckIn)
		perf.POST("/checkins/:id/complete", h.CompleteCheckIn)

		// Calibration
		perf.POST("/calibrations", h.CreateCalibrationSession)
		perf.GET("/calibrations", h.ListCalibrationSessions)
		perf.GET("/calibrations/:id", h.GetCalibrationSession)
		perf.PUT("/calibrations/:id", h.UpdateCalibrationSession)
		perf.POST("/calibrations/:id/status", h.UpdateCalibrationStatus)
		perf.POST("/calibrations/:id/items", h.AddCalibrationItems)
		perf.GET("/calibrations/:id/items", h.ListCalibrationItems)
		perf.PUT("/calibrations/items/:itemId", h.UpdateCalibrationItem)
		perf.POST("/calibrations/items/:itemId/approve", h.ApproveCalibrationItem)

		// Improvement Plans
		perf.POST("/improvement-plans", h.CreateImprovementPlan)
		perf.GET("/improvement-plans", h.ListImprovementPlans)
		perf.GET("/improvement-plans/:id", h.GetImprovementPlan)
		perf.PUT("/improvement-plans/:id", h.UpdateImprovementPlan)
		perf.POST("/improvement-plans/:id/status", h.UpdateImprovementPlanStatus)

		// Development Plans
		perf.POST("/development-plans", h.CreateDevelopmentPlan)
		perf.GET("/development-plans", h.ListDevelopmentPlans)
		perf.GET("/development-plans/:id", h.GetDevelopmentPlan)
		perf.PUT("/development-plans/:id", h.UpdateDevelopmentPlan)
		perf.POST("/development-plans/:id/status", h.UpdateDevelopmentPlanStatus)

		// Evidence
		perf.POST("/evidence", h.CreateEvidence)
		perf.GET("/evidence/:id", h.GetEvidence)
		perf.DELETE("/evidence/:id", h.DeleteEvidence)

		// Results
		perf.POST("/results/calculate", h.CalculateResult)
		perf.GET("/results", h.ListResults)
		perf.GET("/results/:cycleId/:employeeId", h.GetResult)

		// Dashboard
		perf.GET("/dashboard", h.GetDashboard)
	}
}

// ---- Cycles ----

func (h *Handler) CreateCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateCycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	cycle := ToDomainCycle(&req, companyID, userID)
	if err := h.cycles.Create(c.Request.Context(), cycle); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, cycle)
}

func (h *Handler) GetCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycle, err := h.cycles.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Ciclo no encontrado")
		return
	}
	response.Success(c, cycle)
}

func (h *Handler) ListCycles(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.PerformanceCycleFilter{
		CompanyID: companyID,
		Status:    domain.CycleStatus(c.Query("status")),
		Type:      domain.CycleType(c.Query("cycle_type")),
	}
	cycles, err := h.cycles.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, cycles)
}

func (h *Handler) UpdateCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateCycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	existing, err := h.cycles.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Ciclo no encontrado")
		return
	}
	if req.Name != nil { existing.Name = *req.Name }
	if req.Description != nil { existing.Description = req.Description }
	if req.CycleType != nil { existing.CycleType = domain.CycleType(*req.CycleType) }
	if req.StartDate != nil { existing.StartDate = req.StartDate }
	if req.EndDate != nil { existing.EndDate = req.EndDate }
	if req.EvaluationStartDate != nil { existing.EvaluationStartDate = req.EvaluationStartDate }
	if req.EvaluationEndDate != nil { existing.EvaluationEndDate = req.EvaluationEndDate }
	if req.ObjectiveWeight != nil { existing.ObjectiveWeight = *req.ObjectiveWeight }
	if req.CompetencyWeight != nil { existing.CompetencyWeight = *req.CompetencyWeight }
	if err := h.cycles.Update(c.Request.Context(), existing); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, existing)
}

func (h *Handler) UpdateCycleStatus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CycleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	status := domain.CycleStatus(req.Status)
	if err := h.cycles.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), status); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Estado actualizado"})
}

// ---- Templates ----

func (h *Handler) CreateTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var t domain.PerformanceTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	t.CompanyID = companyID
	if err := h.templateRepo.Create(c.Request.Context(), &t); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, t)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	templates, err := h.templateRepo.List(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, templates)
}

func (h *Handler) GetTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	t, err := h.templateRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Plantilla no encontrada")
		return
	}
	sections, _ := h.templateRepo.ListSectionsByTemplate(c.Request.Context(), t.ID)
	questions, _ := h.templateRepo.ListQuestionsByTemplate(c.Request.Context(), t.ID)
	t.Sections = sections
	t.Questions = questions
	response.Success(c, t)
}

func (h *Handler) UpdateTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var t domain.PerformanceTemplate
	if err := c.ShouldBindJSON(&t); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	t.CompanyID = companyID
	t.ID = c.Param("id")
	if err := h.templateRepo.Update(c.Request.Context(), &t); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, t)
}

func (h *Handler) DeleteTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.templateRepo.Delete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// ---- Scales ----

func (h *Handler) CreateScale(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var s domain.RatingScale
	if err := c.ShouldBindJSON(&s); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	s.CompanyID = companyID
	if err := h.scaleRepo.Create(c.Request.Context(), &s); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, s)
}

func (h *Handler) ListScales(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	scales, err := h.scaleRepo.List(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, scales)
}

func (h *Handler) GetScale(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	s, err := h.scaleRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Escala no encontrada")
		return
	}
	levels, _ := h.scaleRepo.ListLevelsByScale(c.Request.Context(), s.ID)
	s.Levels = levels
	response.Success(c, s)
}

func (h *Handler) UpdateScale(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var s domain.RatingScale
	if err := c.ShouldBindJSON(&s); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	s.CompanyID = companyID
	s.ID = c.Param("id")
	if err := h.scaleRepo.Update(c.Request.Context(), &s); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, s)
}

func (h *Handler) DeleteScale(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.scaleRepo.Delete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// ---- Competencies ----

func (h *Handler) CreateCompetency(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var comp domain.Competency
	if err := c.ShouldBindJSON(&comp); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	comp.CompanyID = companyID
	if err := h.compRepo.Create(c.Request.Context(), &comp); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, comp)
}

func (h *Handler) ListCompetencies(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.CompetencyFilter{
		CompanyID: companyID,
		Category:  c.Query("category"),
		Type:      domain.CompetencyType(c.Query("competency_type")),
	}
	comps, err := h.compRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, comps)
}

func (h *Handler) GetCompetency(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	comp, err := h.compRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Competencia no encontrada")
		return
	}
	levels, _ := h.compRepo.ListLevelsByCompetency(c.Request.Context(), comp.ID)
	comp.Levels = levels
	response.Success(c, comp)
}

func (h *Handler) UpdateCompetency(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var comp domain.Competency
	if err := c.ShouldBindJSON(&comp); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	comp.CompanyID = companyID
	comp.ID = c.Param("id")
	if err := h.compRepo.Update(c.Request.Context(), &comp); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, comp)
}

func (h *Handler) DeleteCompetency(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.compRepo.Delete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

func (h *Handler) CreateCompetencyLevel(c *gin.Context) {
	competencyID := c.Param("id")
	var l domain.CompetencyLevel
	if err := c.ShouldBindJSON(&l); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	l.CompetencyID = competencyID
	if err := h.compRepo.CreateLevel(c.Request.Context(), &l); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, l)
}

func (h *Handler) UpsertPositionCompetency(c *gin.Context) {
	var pc domain.PositionCompetency
	if err := c.ShouldBindJSON(&pc); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	pc.CompanyID = tenant.GetCompanyID(c)
	if err := h.compRepo.UpsertPositionCompetency(c.Request.Context(), &pc); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, pc)
}

func (h *Handler) UpsertCycleCompetency(c *gin.Context) {
	var cc domain.CycleCompetency
	if err := c.ShouldBindJSON(&cc); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	if err := h.compRepo.UpsertCycleCompetency(c.Request.Context(), &cc); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cc)
}

// ---- Objectives ----

func (h *Handler) CreateObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	o := ToDomainObjective(&req, companyID, userID)
	if err := h.objectives.Create(c.Request.Context(), o); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	for _, kr := range req.KeyResults {
		dkr := &domain.ObjectiveKeyResult{
			ObjectiveID: o.ID,
			Title:       kr.Title,
			Description: kr.Description,
			TargetValue: kr.TargetValue,
			Unit:        kr.Unit,
		}
		if kr.Weight != nil { dkr.Weight = *kr.Weight }
		if kr.SortOrder != nil { dkr.SortOrder = *kr.SortOrder }
		h.objectiveRepo.CreateKeyResult(c.Request.Context(), dkr)
	}
	response.Created(c, o)
}

func (h *Handler) ListObjectives(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.ObjectiveFilter{
		CompanyID:  companyID,
		CycleID:    c.Query("cycle_id"),
		EmployeeID: c.Query("employee_id"),
		Status:     domain.ObjectiveStatus(c.Query("status")),
	}
	objectives, err := 	h.objectiveRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, objectives)
}

func (h *Handler) GetObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	o, err := 	h.objectiveRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Objetivo no encontrado")
		return
	}
	krs, _ := h.objectiveRepo.ListKeyResultsByObjective(c.Request.Context(), o.ID)
	o.KeyResults = krs
	response.Success(c, o)
}

func (h *Handler) UpdateObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	o, err := h.objectiveRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Objetivo no encontrado")
		return
	}
	if req.Title != nil { o.Title = *req.Title }
	if req.Description != nil { o.Description = req.Description }
	if req.Weight != nil { o.Weight = *req.Weight }
	if req.Status != nil { o.Status = domain.ObjectiveStatus(*req.Status) }
	if req.TargetValue != nil { o.TargetValue = req.TargetValue }
	if req.CurrentValue != nil { o.CurrentValue = req.CurrentValue }
	if req.Unit != nil { o.Unit = req.Unit }
	if req.Notes != nil { o.Notes = req.Notes }
	if req.RiskNotes != nil { o.RiskNotes = req.RiskNotes }
	if err := h.objectiveRepo.Update(c.Request.Context(), o); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, o)
}

func (h *Handler) DeleteObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	o, err := h.objectiveRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Objetivo no encontrado")
		return
	}
	o.Status = domain.ObjectiveStatusCancelled
	if err := h.objectiveRepo.Update(c.Request.Context(), o); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

func (h *Handler) UpdateObjectiveProgress(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	o, err := h.objectives.UpdateProgress(c.Request.Context(), companyID, c.Param("id"), req.CurrentValue)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, o)
}

func (h *Handler) CreateKeyResult(c *gin.Context) {
	var kr domain.ObjectiveKeyResult
	if err := c.ShouldBindJSON(&kr); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	kr.ObjectiveID = c.Param("id")
	if err := h.objectiveRepo.CreateKeyResult(c.Request.Context(), &kr); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, kr)
}

// ---- Participants ----

func (h *Handler) AssignParticipants(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req AssignEvaluatorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	participants := ToDomainParticipants(&req, companyID)
	if err := h.evaluations.AssignParticipants(c.Request.Context(), participants); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, participants)
}

func (h *Handler) ListParticipants(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	employeeID := c.Query("employee_id")
	if employeeID != "" {
		participants, err := h.participantRepo.ListByEmployee(c.Request.Context(), companyID, cycleID, employeeID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		response.Success(c, participants)
		return
	}
	participants, err := h.participantRepo.ListByCycle(c.Request.Context(), companyID, cycleID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, participants)
}

func (h *Handler) RemoveParticipant(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.participantRepo.Delete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// ---- Evaluations ----

func (h *Handler) CreateEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	eval := ToDomainEvaluation(&req, companyID)
	if err := h.evalRepo.Create(c.Request.Context(), eval); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if len(req.Answers) > 0 {
		answers := ToDomainAnswers(eval.ID, req.Answers)
		h.evalRepo.BulkCreateAnswers(c.Request.Context(), answers)
	}
	response.Created(c, eval)
}

func (h *Handler) GetEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	eval, err := h.evalRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Evaluación no encontrada")
		return
	}
	answers, _ := h.evalRepo.ListAnswersByEvaluation(c.Request.Context(), eval.ID)
	eval.Answers = answers
	response.Success(c, eval)
}

func (h *Handler) ListEvaluations(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.EvaluationFilter{
		CompanyID:    companyID,
		CycleID:      c.Query("cycle_id"),
		EmployeeID:   c.Query("employee_id"),
		EvaluatorID:  c.Query("evaluator_id"),
		Status:       domain.EvaluationStatus(c.Query("status")),
		EvaluationType: domain.EvaluationType(c.Query("evaluation_type")),
	}
	evaluations, err := h.evalRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, evaluations)
}

func (h *Handler) UpdateEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var eval domain.PerformanceEvaluation
	if err := c.ShouldBindJSON(&eval); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	eval.CompanyID = companyID
	eval.ID = c.Param("id")
	if err := h.evalRepo.Update(c.Request.Context(), &eval); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, eval)
}

func (h *Handler) SubmitEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.evaluations.SubmitEvaluation(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Evaluación enviada"})
}

func (h *Handler) ApproveEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := workflow.ValidateEvaluationTransition(domain.EvaluationStatusSubmitted, domain.EvaluationStatusApproved); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.evalRepo.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), domain.EvaluationStatusApproved); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Evaluación aprobada"})
}

func (h *Handler) ReopenEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	eval, err := h.evalRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Evaluación no encontrada")
		return
	}
	if err := workflow.ValidateEvaluationTransition(eval.Status, domain.EvaluationStatusReopened); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.evalRepo.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), domain.EvaluationStatusReopened); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Evaluación reabierta"})
}

func (h *Handler) LockEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	eval, err := h.evalRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Evaluación no encontrada")
		return
	}
	if err := workflow.ValidateEvaluationTransition(eval.Status, domain.EvaluationStatusLocked); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.evalRepo.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), domain.EvaluationStatusLocked); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Evaluación bloqueada"})
}

// ---- Answers ----

func (h *Handler) CreateAnswers(c *gin.Context) {
	evaluationID := c.Param("id")
	var reqs []CreateAnswerRequest
	if err := c.ShouldBindJSON(&reqs); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	answers := ToDomainAnswers(evaluationID, reqs)
	if err := h.evalRepo.BulkCreateAnswers(c.Request.Context(), answers); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, answers)
}

func (h *Handler) ListAnswers(c *gin.Context) {
	answers, err := h.evalRepo.ListAnswersByEvaluation(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, answers)
}

// ---- Objective/Competency Evaluations ----

func (h *Handler) CreateObjectiveEvaluation(c *gin.Context) {
	evaluationID := c.Param("id")
	var oe domain.ObjectiveEvaluation
	if err := c.ShouldBindJSON(&oe); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	oe.EvaluationID = evaluationID
	if err := h.evalRepo.CreateObjectiveEvaluation(c.Request.Context(), &oe); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, oe)
}

func (h *Handler) CreateCompetencyEvaluation(c *gin.Context) {
	evaluationID := c.Param("id")
	var ce domain.CompetencyEvaluation
	if err := c.ShouldBindJSON(&ce); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	ce.EvaluationID = evaluationID
	if err := h.evalRepo.CreateCompetencyEvaluation(c.Request.Context(), &ce); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, ce)
}

// ---- Reviews ----

func (h *Handler) CreateReview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	rev := ToDomainReview(&req, companyID)
	if err := h.evaluations.CreateReview(c.Request.Context(), rev); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, rev)
}

func (h *Handler) ListReviews(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	reviews, err := h.reviewRepo.ListByCycle(c.Request.Context(), companyID, cycleID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, reviews)
}

func (h *Handler) GetReview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	rev, err := h.reviewRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Revisión no encontrada")
		return
	}
	response.Success(c, rev)
}

func (h *Handler) UpdateReview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	rev, err := h.reviewRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Revisión no encontrada")
		return
	}
	if req.Summary != nil { rev.Summary = req.Summary }
	if req.Strengths != nil { rev.Strengths = req.Strengths }
	if req.ImprovementAreas != nil { rev.ImprovementAreas = req.ImprovementAreas }
	if req.FinalScore != nil { rev.FinalScore = req.FinalScore }
	if req.FinalRating != nil { rev.FinalRating = req.FinalRating }
	if req.EmployeeComments != nil { rev.EmployeeComments = req.EmployeeComments }
	if req.ManagerComments != nil { rev.ManagerComments = req.ManagerComments }
	if req.EmployeeAgreement != nil { rev.EmployeeAgreement = *req.EmployeeAgreement }
	if err := h.reviewRepo.Update(c.Request.Context(), rev); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rev)
}

func (h *Handler) UpdateReviewStatus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CycleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	if err := h.reviewRepo.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), domain.EvaluationStatus(req.Status)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Estado actualizado"})
}

// ---- Feedback ----

func (h *Handler) CreateFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	fb := ToDomainFeedback(&req, companyID, userID)
	if err := h.feedbackRepo.Create(c.Request.Context(), fb); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, fb)
}

func (h *Handler) ListFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.FeedbackFilter{
		CompanyID:    companyID,
		EmployeeID:   c.Query("employee_id"),
		AuthorID:     c.Query("author_id"),
		FeedbackType: domain.FeedbackType(c.Query("feedback_type")),
		Visibility:   domain.FeedbackVisibility(c.Query("visibility")),
	}
	feedbacks, err := h.feedbackRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, feedbacks)
}

func (h *Handler) GetFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	fb, err := h.feedbackRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Feedback no encontrado")
		return
	}
	response.Success(c, fb)
}

func (h *Handler) DeleteFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.feedbackRepo.Delete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// ---- Recognitions ----

func (h *Handler) CreateRecognition(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateRecognitionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	rec := ToDomainRecognition(&req, companyID, userID)
	if err := h.feedbackRepo.CreateRecognition(c.Request.Context(), rec); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, rec)
}

func (h *Handler) ListRecognitions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Query("employee_id")
	recognitions, err := h.feedbackRepo.ListRecognitionsByEmployee(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, recognitions)
}

// ---- Check-ins ----

func (h *Handler) CreateCheckIn(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	ci := ToDomainCheckIn(&req, companyID)
	if err := h.checkInRepo.Create(c.Request.Context(), ci); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, ci)
}

func (h *Handler) ListCheckIns(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Query("employee_id")
	managerID := c.Query("manager_id")
	var checkins []domain.PerformanceCheckIn
	var err error
	if employeeID != "" {
		checkins, err = h.checkInRepo.ListByEmployee(c.Request.Context(), companyID, employeeID)
	} else if managerID != "" {
		checkins, err = h.checkInRepo.ListByManager(c.Request.Context(), companyID, managerID)
	} else {
		response.BadRequest(c, "Se requiere employee_id o manager_id")
		return
	}
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, checkins)
}

func (h *Handler) GetCheckIn(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	ci, err := h.checkInRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Check-in no encontrado")
		return
	}
	response.Success(c, ci)
}

func (h *Handler) UpdateCheckIn(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var ci domain.PerformanceCheckIn
	if err := c.ShouldBindJSON(&ci); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	ci.CompanyID = companyID
	ci.ID = c.Param("id")
	if err := h.checkInRepo.Update(c.Request.Context(), &ci); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, ci)
}

func (h *Handler) CompleteCheckIn(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CompleteCheckInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	notes := map[string]*string{
		"employee_notes": req.EmployeeNotes,
		"manager_notes":  req.ManagerNotes,
		"achievements":   req.Achievements,
		"blockers":       req.Blockers,
		"next_steps":     req.NextSteps,
	}
	if err := h.checkInRepo.Complete(c.Request.Context(), companyID, c.Param("id"), notes); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Check-in completado"})
}

// ---- Calibration ----

func (h *Handler) CreateCalibrationSession(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateCalibrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	s := ToDomainCalibration(&req, companyID, userID)
	if err := h.calibRepo.CreateSession(c.Request.Context(), s); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, s)
}

func (h *Handler) ListCalibrationSessions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	sessions, err := h.calibRepo.ListSessionsByCycle(c.Request.Context(), companyID, cycleID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, sessions)
}

func (h *Handler) GetCalibrationSession(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	s, err := h.calibRepo.GetSessionByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Sesión de calibración no encontrada")
		return
	}
	items, _ := h.calibRepo.ListItemsBySession(c.Request.Context(), s.ID)
	s.Items = items
	response.Success(c, s)
}

func (h *Handler) UpdateCalibrationSession(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var s domain.CalibrationSession
	if err := c.ShouldBindJSON(&s); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	s.CompanyID = companyID
	s.ID = c.Param("id")
	if err := h.calibRepo.UpdateSession(c.Request.Context(), &s); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, s)
}

func (h *Handler) UpdateCalibrationStatus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CycleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	if err := h.calibRepo.UpdateSessionStatus(c.Request.Context(), companyID, c.Param("id"), domain.CalibrationStatus(req.Status)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Estado actualizado"})
}

func (h *Handler) AddCalibrationItems(c *gin.Context) {
	sessionID := c.Param("id")
	var items []domain.CalibrationItem
	if err := c.ShouldBindJSON(&items); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	for i := range items {
		items[i].SessionID = sessionID
	}
	if err := h.calibRepo.BulkCreateItems(c.Request.Context(), items); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, items)
}

func (h *Handler) ListCalibrationItems(c *gin.Context) {
	items, err := h.calibRepo.ListItemsBySession(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *Handler) UpdateCalibrationItem(c *gin.Context) {
	var item domain.CalibrationItem
	if err := c.ShouldBindJSON(&item); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	item.ID = c.Param("itemId")
	if err := h.calibRepo.UpdateItem(c.Request.Context(), &item); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, item)
}

func (h *Handler) ApproveCalibrationItem(c *gin.Context) {
	userID := tenant.GetUserID(c)
	if err := h.calibRepo.ApproveItem(c.Request.Context(), c.Param("itemId"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Item aprobado"})
}

// ---- Improvement Plans ----

func (h *Handler) CreateImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateImprovementPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	p := ToDomainImprovementPlan(&req, companyID, userID)
	if err := h.plans.CreateImprovement(c.Request.Context(), p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	for _, a := range req.Actions {
		da := &domain.ImprovementPlanAction{
			PlanID:       p.ID,
			Title:        a.Title,
			Description:  a.Description,
			ResponsibleID: a.ResponsibleID,
			DueDate:      a.DueDate,
		}
		h.impRepo.CreateAction(c.Request.Context(), da)
	}
	response.Created(c, p)
}

func (h *Handler) ListImprovementPlans(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.PlanFilter{
		CompanyID:  companyID,
		EmployeeID: c.Query("employee_id"),
		Status:     domain.PlanStatus(c.Query("status")),
	}
	plans, err := h.impRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, plans)
}

func (h *Handler) GetImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	p, err := h.impRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Plan no encontrado")
		return
	}
	actions, _ := h.impRepo.ListActionsByPlan(c.Request.Context(), p.ID)
	p.Actions = actions
	response.Success(c, p)
}

func (h *Handler) UpdateImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var p domain.ImprovementPlan
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	p.CompanyID = companyID
	p.ID = c.Param("id")
	if err := h.impRepo.Update(c.Request.Context(), &p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, p)
}

func (h *Handler) UpdateImprovementPlanStatus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CycleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	if err := h.impRepo.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), domain.PlanStatus(req.Status)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Estado actualizado"})
}

// ---- Development Plans ----

func (h *Handler) CreateDevelopmentPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateDevelopmentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	p := ToDomainDevPlan(&req, companyID, userID)
	if err := h.plans.CreateDevelopment(c.Request.Context(), p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	for _, a := range req.Actions {
		da := &domain.DevelopmentPlanAction{
			PlanID:      p.ID,
			Title:       a.Title,
			Description: a.Description,
			ActionType:  a.ActionType,
			DueDate:     a.DueDate,
		}
		h.devRepo.CreateAction(c.Request.Context(), da)
	}
	response.Created(c, p)
}

func (h *Handler) ListDevelopmentPlans(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filter := domain.PlanFilter{
		CompanyID:  companyID,
		EmployeeID: c.Query("employee_id"),
		Status:     domain.PlanStatus(c.Query("status")),
	}
	plans, err := h.devRepo.List(c.Request.Context(), filter)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, plans)
}

func (h *Handler) GetDevelopmentPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	p, err := h.devRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Plan no encontrado")
		return
	}
	actions, _ := h.devRepo.ListActionsByPlan(c.Request.Context(), p.ID)
	p.Actions = actions
	response.Success(c, p)
}

func (h *Handler) UpdateDevelopmentPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var p domain.DevelopmentPlan
	if err := c.ShouldBindJSON(&p); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	p.CompanyID = companyID
	p.ID = c.Param("id")
	if err := h.devRepo.Update(c.Request.Context(), &p); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, p)
}

func (h *Handler) UpdateDevelopmentPlanStatus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CycleStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	if err := h.devRepo.UpdateStatus(c.Request.Context(), companyID, c.Param("id"), domain.PlanStatus(req.Status)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "Estado actualizado"})
}

// ---- Evidence ----

func (h *Handler) CreateEvidence(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	e := ToDomainEvidence(&req, companyID, userID)
	if err := h.evidenceRepo.Create(c.Request.Context(), e); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, e)
}

func (h *Handler) GetEvidence(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	e, err := h.evidenceRepo.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Evidencia no encontrada")
		return
	}
	response.Success(c, e)
}

func (h *Handler) DeleteEvidence(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.evidenceRepo.Delete(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

// ---- Results ----

func (h *Handler) CalculateResult(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CalculateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Solicitud inválida")
		return
	}
	result, err := h.scorer.Calculate(c.Request.Context(), companyID, req.CycleID, req.EmployeeID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	result.CompanyID = companyID
	result.CycleID = req.CycleID
	result.EmployeeID = req.EmployeeID
	if err := h.resultRepo.Upsert(c.Request.Context(), result); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *Handler) GetResult(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	result, err := h.resultRepo.GetByCycleEmployee(c.Request.Context(), companyID, c.Param("cycleId"), c.Param("employeeId"))
	if err != nil {
		response.NotFound(c, "Resultado no encontrado")
		return
	}
	response.Success(c, result)
}

func (h *Handler) ListResults(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	employeeID := c.Query("employee_id")

	if employeeID != "" {
		results, err := h.resultRepo.ListByEmployee(c.Request.Context(), companyID, employeeID)
		if err != nil {
			response.InternalError(c, err.Error())
			return
		}
		response.Success(c, results)
		return
	}

	results, err := h.resultRepo.ListByCycle(c.Request.Context(), companyID, cycleID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, results)
}

// ---- Dashboard ----

func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	dash, err := h.dashRepo.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dash,
	})
}
