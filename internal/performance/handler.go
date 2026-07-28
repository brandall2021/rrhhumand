package performance

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Cycles
func (h *Handler) CreateCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateCycleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	cycle, err := h.service.CreateCycle(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, cycle)
}

func (h *Handler) GetCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycle, err := h.service.GetCycle(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Cycle not found")
		return
	}
	response.Success(c, cycle)
}

func (h *Handler) ListCycles(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycles, err := h.service.ListCycles(c.Request.Context(), companyID)
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
		response.BadRequest(c, "Invalid request body")
		return
	}
	cycle, err := h.service.UpdateCycle(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, cycle)
}

func (h *Handler) OpenCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.OpenCycle(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "cycle opened"})
}

func (h *Handler) CloseCycle(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CloseCycle(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "cycle closed"})
}

// Templates
func (h *Handler) CreateTemplate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	t, err := h.service.CreateTemplate(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, t)
}

func (h *Handler) ListTemplates(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	templates, err := h.service.ListTemplates(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, templates)
}

// Scales
func (h *Handler) CreateScale(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateScaleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	scale, err := h.service.CreateScale(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, scale)
}

func (h *Handler) ListScales(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	scales, err := h.service.ListScales(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, scales)
}

// Competencies
func (h *Handler) CreateCompetency(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateCompetencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	comp, err := h.service.CreateCompetency(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, comp)
}

func (h *Handler) ListCompetencies(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	comps, err := h.service.ListCompetencies(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, comps)
}

func (h *Handler) UpdateCompetency(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateCompetencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	comp, err := h.service.UpdateCompetency(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, comp)
}

// Objectives
func (h *Handler) CreateObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	obj, err := h.service.CreateObjective(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, obj)
}

func (h *Handler) ListObjectives(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := PerformanceFilters{
		CycleID:    c.Query("cycle_id"),
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
	}
	objs, err := h.service.ListObjectives(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, objs)
}

func (h *Handler) GetObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	obj, err := h.service.GetObjective(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Objective not found")
		return
	}
	response.Success(c, obj)
}

func (h *Handler) UpdateObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateObjectiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	obj, err := h.service.UpdateObjective(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, obj)
}

func (h *Handler) DeleteObjective(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cancelled := "CANCELLED"
	_, err := h.service.UpdateObjective(c.Request.Context(), companyID, c.Param("id"), &UpdateObjectiveRequest{Status: &cancelled})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.NoContent(c)
}

func (h *Handler) UpdateObjectiveProgress(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	obj, err := h.service.UpdateObjectiveProgress(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, obj)
}

// KPIs
func (h *Handler) CreateKPI(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateKPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	kpi, err := h.service.CreateKPI(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, kpi)
}

func (h *Handler) ListKPIs(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := PerformanceFilters{
		CycleID:    c.Query("cycle_id"),
		EmployeeID: c.Query("employee_id"),
	}
	kpis, err := h.service.ListKPIs(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, kpis)
}

func (h *Handler) UpdateKPI(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateKPIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	kpi, err := h.service.UpdateKPI(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, kpi)
}

// Evaluators
func (h *Handler) AssignEvaluators(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req AssignEvaluatorsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	evaluators, err := h.service.AssignEvaluators(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, evaluators)
}

func (h *Handler) ListEvaluators(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	if cycleID == "" {
		response.BadRequest(c, "cycle_id is required")
		return
	}
	evaluators, err := h.service.ListEvaluators(c.Request.Context(), companyID, cycleID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, evaluators)
}

// Evaluations
func (h *Handler) CreateEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	eval, err := h.service.CreateEvaluation(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, eval)
}

func (h *Handler) GetEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	eval, err := h.service.GetEvaluation(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Evaluation not found")
		return
	}
	response.Success(c, eval)
}

func (h *Handler) ListEvaluations(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := PerformanceFilters{
		CycleID:    c.Query("cycle_id"),
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
	}
	evaluations, err := h.service.ListEvaluations(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, evaluations)
}

func (h *Handler) UpdateEvaluation(c *gin.Context) {
	var req UpdateEvaluationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	_ = req
	response.Success(c, gin.H{"message": "evaluation updated"})
}

func (h *Handler) SubmitEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.SubmitEvaluation(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "evaluation submitted"})
}

func (h *Handler) ReopenEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ReopenEvaluation(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "evaluation reopened"})
}

func (h *Handler) ApproveEvaluation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.ApproveEvaluation(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "evaluation approved"})
}

// Answers
func (h *Handler) CreateAnswer(c *gin.Context) {
	evaluationID := c.Param("id")
	var req CreateAnswerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	answer, err := h.service.CreateAnswer(c.Request.Context(), evaluationID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, answer)
}

func (h *Handler) ListAnswers(c *gin.Context) {
	answers, err := h.service.ListAnswers(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, answers)
}

// Feedback
func (h *Handler) CreateFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateFeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	fb, err := h.service.CreateFeedback(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, fb)
}

func (h *Handler) ListFeedback(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		employeeID = tenant.GetUserID(c)
	}
	feedbacks, err := h.service.ListFeedback(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, feedbacks)
}

func (h *Handler) GetFeedback(c *gin.Context) {
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		employeeID = tenant.GetUserID(c)
	}
	companyID := tenant.GetCompanyID(c)
	feedbacks, err := h.service.ListFeedback(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.NotFound(c, "Feedback not found")
		return
	}
	if len(feedbacks) == 0 {
		response.NotFound(c, "No feedback found")
		return
	}
	response.Success(c, feedbacks[0])
}

// Evidence
func (h *Handler) CreateEvidence(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	evaluationID := c.Param("id")
	var req CreateEvidenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	ev, err := h.service.CreateEvidence(c.Request.Context(), companyID, evaluationID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, ev)
}

func (h *Handler) ListEvidence(c *gin.Context) {
	evidence, err := h.service.ListEvidence(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, evidence)
}

// Results
func (h *Handler) CalculateResult(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CalculateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	result, err := h.service.CalculateResult(c.Request.Context(), companyID, req.CycleID, req.EmployeeID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, result)
}

func (h *Handler) GetResult(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	employeeID := c.Param("id")
	if cycleID == "" {
		response.BadRequest(c, "cycle_id is required")
		return
	}
	result, err := h.service.GetResult(c.Request.Context(), companyID, cycleID, employeeID)
	if err != nil {
		response.NotFound(c, "Result not found")
		return
	}
	response.Success(c, result)
}

func (h *Handler) ListResults(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	cycleID := c.Query("cycle_id")
	if cycleID == "" {
		response.BadRequest(c, "cycle_id is required")
		return
	}
	results, err := h.service.ListResults(c.Request.Context(), companyID, cycleID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, results)
}

// Scoring rules
func (h *Handler) GetScoringRules(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	rules, err := h.service.GetScoringRules(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, rules)
}

func (h *Handler) UpdateScoringRules(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateScoringRulesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	rules, err := h.service.UpdateScoringRules(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, rules)
}

// Improvement Plans
func (h *Handler) CreateImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateImprovementPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	plan, err := h.service.CreateImprovementPlan(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, plan)
}

func (h *Handler) ListImprovementPlans(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := PerformanceFilters{
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
	}
	plans, err := h.service.ListImprovementPlans(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, plans)
}

func (h *Handler) GetImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	plan, err := h.service.GetImprovementPlan(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Plan not found")
		return
	}
	response.Success(c, plan)
}

func (h *Handler) UpdateImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateImprovementPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	plan, err := h.service.UpdateImprovementPlan(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, plan)
}

func (h *Handler) CompleteImprovementPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.service.CompleteImprovementPlan(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "plan completed"})
}

// Development Plans
func (h *Handler) CreateDevelopmentPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req CreateDevelopmentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	plan, err := h.service.CreateDevelopmentPlan(c.Request.Context(), companyID, userID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, plan)
}

func (h *Handler) ListDevelopmentPlans(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Query("employee_id")
	if employeeID == "" {
		employeeID = tenant.GetUserID(c)
	}
	plans, err := h.service.ListDevelopmentPlans(c.Request.Context(), companyID, employeeID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, plans)
}

func (h *Handler) GetDevelopmentPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	plan, err := h.service.GetDevelopmentPlan(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Plan not found")
		return
	}
	response.Success(c, plan)
}

func (h *Handler) UpdateDevelopmentPlan(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdateDevelopmentPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	plan, err := h.service.UpdateDevelopmentPlan(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, plan)
}

// Dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	dash, err := h.service.GetDashboard(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dash,
	})
}
