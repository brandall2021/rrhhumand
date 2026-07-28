package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/payroll/domain"
)

type Handler struct {
	svc *domain.PayrollService
}

func NewHandler(svc *domain.PayrollService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) companyID(c *gin.Context) string { return c.GetString("company_id") }
func (h *Handler) userID(c *gin.Context) string    { return c.GetString("user_id") }
func (h *Handler) employeeID(c *gin.Context) string {
	if eid := c.GetString("employee_id"); eid != "" {
		return eid
	}
	return c.GetString("user_id")
}

func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func qs(c *gin.Context, key string) *string {
	v := c.Query(key)
	if v == "" {
		return nil
	}
	return &v
}

func qi(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}

// ========================================================================
// PERIODS
// ========================================================================

func (h *Handler) CreatePeriod(c *gin.Context) {
	var req CreatePeriodReq
	if !bindJSON(c, &req) {
		return
	}
	p, err := h.svc.CreatePeriod(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, p)
}

func (h *Handler) UpdatePeriod(c *gin.Context) {
	var req UpdatePeriodReq
	if !bindJSON(c, &req) {
		return
	}
	p, err := h.svc.UpdatePeriod(c.Request.Context(), h.companyID(c), c.Param("id"), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, p)
}

func (h *Handler) GetPeriod(c *gin.Context) {
	p, err := h.svc.GetPeriod(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "period not found"})
		return
	}
	success(c, p)
}

func (h *Handler) ListPeriods(c *gin.Context) {
	list, err := h.svc.ListPeriods(c.Request.Context(), h.companyID(c), qi(c, "limit", 20), qi(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) ClosePeriod(c *gin.Context) {
	if err := h.svc.ClosePeriod(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "period closed"})
}

// ========================================================================
// RUNS
// ========================================================================

func (h *Handler) CreateRun(c *gin.Context) {
	var req CreateRunReq
	if !bindJSON(c, &req) {
		return
	}
	run, err := h.svc.CreateRun(c.Request.Context(), h.companyID(c), c.Param("period_id"), h.userID(c), req.RunType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, run)
}

func (h *Handler) GetRun(c *gin.Context) {
	run, err := h.svc.GetRun(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "run not found"})
		return
	}
	success(c, run)
}

func (h *Handler) ListRuns(c *gin.Context) {
	list, err := h.svc.ListRuns(c.Request.Context(), h.companyID(c), qs(c, "period_id"), qs(c, "run_type"), qs(c, "status"), qi(c, "limit", 20), qi(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) CalculateRun(c *gin.Context) {
	if err := h.svc.CalculateRun(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "calculation started"})
}

func (h *Handler) ValidateRun(c *gin.Context) {
	if err := h.svc.ValidateRun(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "validation passed"})
}

func (h *Handler) ApproveRun(c *gin.Context) {
	if err := h.svc.ApproveRun(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "run approved"})
}

func (h *Handler) CloseRun(c *gin.Context) {
	if err := h.svc.CloseRun(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "run closed"})
}

// ========================================================================
// RUN EMPLOYEES
// ========================================================================

func (h *Handler) AddEmployee(c *gin.Context) {
	var req struct {
		EmployeeID string `json:"employee_id" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	re, err := h.svc.AddEmployeeToRun(c.Request.Context(), h.companyID(c), c.Param("id"), req.EmployeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, re)
}

func (h *Handler) ListRunEmployees(c *gin.Context) {
	list, err := h.svc.ListRunEmployees(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) GetRunEmployee(c *gin.Context) {
	re, err := h.svc.GetRunEmployee(c.Request.Context(), h.companyID(c), c.Param("run_id"), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}
	success(c, re)
}

func (h *Handler) GetEmployeeItems(c *gin.Context) {
	items, err := h.svc.GetEmployeeItems(c.Request.Context(), h.companyID(c), c.Param("run_id"), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "items not found"})
		return
	}
	success(c, items)
}

// ========================================================================
// CONCEPTS
// ========================================================================

func (h *Handler) CreateConcept(c *gin.Context) {
	var req CreateConceptReq
	if !bindJSON(c, &req) {
		return
	}
	concept, err := h.svc.CreateConcept(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, concept)
}

func (h *Handler) UpdateConcept(c *gin.Context) {
	var req UpdateConceptReq
	if !bindJSON(c, &req) {
		return
	}
	concept, err := h.svc.UpdateConcept(c.Request.Context(), h.companyID(c), c.Param("id"), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, concept)
}

func (h *Handler) ListConcepts(c *gin.Context) {
	active := qs(c, "active")
	var act *bool
	if active != nil {
		b := *active == "true"
		act = &b
	}
	list, err := h.svc.ListConcepts(c.Request.Context(), h.companyID(c), qs(c, "concept_type"), qs(c, "taxability"), act)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// RULES
// ========================================================================

func (h *Handler) CreateRule(c *gin.Context) {
	var req CreateRuleReq
	if !bindJSON(c, &req) {
		return
	}
	rule, err := h.svc.CreateRule(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	var req UpdateRuleReq
	if !bindJSON(c, &req) {
		return
	}
	rule, err := h.svc.UpdateRule(c.Request.Context(), h.companyID(c), c.Param("id"), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rule)
}

func (h *Handler) ListRules(c *gin.Context) {
	list, err := h.svc.ListRules(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// NOVELTIES
// ========================================================================

func (h *Handler) CreateNovelty(c *gin.Context) {
	var req CreateNoveltyReq
	if !bindJSON(c, &req) {
		return
	}
	n, err := h.svc.CreateNovelty(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, n)
}

func (h *Handler) UpdateNovelty(c *gin.Context) {
	var req UpdateNoveltyReq
	if !bindJSON(c, &req) {
		return
	}
	n, err := h.svc.UpdateNovelty(c.Request.Context(), h.companyID(c), c.Param("id"), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, n)
}

func (h *Handler) ListNovelties(c *gin.Context) {
	list, err := h.svc.ListNovelties(c.Request.Context(), h.companyID(c), qs(c, "employee_id"), qs(c, "period_id"), qs(c, "novelty_type"), qs(c, "status"), qs(c, "source"), qi(c, "limit", 20), qi(c, "offset", 0))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) DeleteNovelty(c *gin.Context) {
	if err := h.svc.DeleteNovelty(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "deleted"})
}

func (h *Handler) ApproveNovelty(c *gin.Context) {
	if err := h.svc.ApproveNovelty(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "approved"})
}

func (h *Handler) ImportNovelties(c *gin.Context) {
	var req ImportNoveltiesReq
	if !bindJSON(c, &req) {
		return
	}
	inputs := make([]domain.CreateNoveltyInput, len(req.Novelties))
	for i, nr := range req.Novelties {
		inputs[i] = nr.ToInput()
	}
	novelties, err := h.svc.ImportNovelties(c.Request.Context(), h.companyID(c), h.userID(c), inputs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, novelties)
}

// ========================================================================
// ADVANCES
// ========================================================================

func (h *Handler) CreateAdvance(c *gin.Context) {
	var req CreateAdvanceReq
	if !bindJSON(c, &req) {
		return
	}
	a, err := h.svc.CreateAdvance(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, a)
}

func (h *Handler) ListAdvances(c *gin.Context) {
	list, err := h.svc.ListAdvances(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// GARNISHMENTS
// ========================================================================

func (h *Handler) CreateGarnishment(c *gin.Context) {
	var req CreateGarnishmentReq
	if !bindJSON(c, &req) {
		return
	}
	g, err := h.svc.CreateGarnishment(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, g)
}

func (h *Handler) ListGarnishments(c *gin.Context) {
	list, err := h.svc.ListGarnishments(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// LABOR AGREEMENTS & CATEGORIES
// ========================================================================

func (h *Handler) CreateAgreement(c *gin.Context) {
	var req CreateAgreementReq
	if !bindJSON(c, &req) {
		return
	}
	a, err := h.svc.CreateAgreement(c.Request.Context(), h.companyID(c), h.userID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, a)
}

func (h *Handler) ListAgreements(c *gin.Context) {
	list, err := h.svc.ListAgreements(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryReq
	if !bindJSON(c, &req) {
		return
	}
	cat, err := h.svc.CreateCategory(c.Request.Context(), h.companyID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, cat)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// SALARY SCALES
// ========================================================================

func (h *Handler) CreateSalaryScale(c *gin.Context) {
	var req CreateSalaryScaleReq
	if !bindJSON(c, &req) {
		return
	}
	sc, err := h.svc.CreateSalaryScale(c.Request.Context(), h.companyID(c), req.ToInput())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, sc)
}

func (h *Handler) ListSalaryScales(c *gin.Context) {
	list, err := h.svc.ListSalaryScales(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// DASHBOARD & ERRORS
// ========================================================================

func (h *Handler) GetDashboard(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, stats)
}

func (h *Handler) GetRunSummary(c *gin.Context) {
	s, err := h.svc.GetRunSummary(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, s)
}

func (h *Handler) ListErrors(c *gin.Context) {
	list, err := h.svc.GetErrors(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

// ========================================================================
// EMPLOYEE SELF-SERVICE
// ========================================================================

func (h *Handler) MyPeriods(c *gin.Context) {
	list, err := h.svc.ListPeriods(c.Request.Context(), h.companyID(c), 12, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, list)
}

func (h *Handler) MyRunResult(c *gin.Context) {
	re, err := h.svc.GetEmployeeResult(c.Request.Context(), h.companyID(c), c.Param("run_id"), h.employeeID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	success(c, re)
}

func (h *Handler) MyItems(c *gin.Context) {
	items, err := h.svc.GetEmployeeItems(c.Request.Context(), h.companyID(c), c.Param("run_id"), h.employeeID(c))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	success(c, items)
}
