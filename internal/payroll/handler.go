package payroll

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) companyID(c *gin.Context) string { return tenant.GetCompanyID(c) }
func (h *Handler) userID(c *gin.Context) string    { return c.GetString("user_id") }
func (h *Handler) employeeID(c *gin.Context) string { return c.GetString("employee_id") }

func (h *Handler) bindJSON(c *gin.Context, obj any) bool {
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

// ========================================================================
// PERIODS
// ========================================================================

func (h *Handler) CreatePeriod(c *gin.Context) {
	var req CreatePeriodRequest
	if !h.bindJSON(c, &req) {
		return
	}
	p, err := h.svc.CreatePeriod(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) UpdatePeriod(c *gin.Context) {
	var req UpdatePeriodRequest
	if !h.bindJSON(c, &req) {
		return
	}
	p, err := h.svc.UpdatePeriod(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) GetPeriod(c *gin.Context) {
	p, err := h.svc.GetPeriod(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "period not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *Handler) ListPeriods(c *gin.Context) {
	limit := qi(c, "limit", 20)
	offset := qi(c, "offset", 0)
	list, err := h.svc.ListPeriods(c.Request.Context(), h.companyID(c), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ClosePeriod(c *gin.Context) {
	if err := h.svc.ClosePeriod(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "period closed"})
}

// ========================================================================
// RUNS
// ========================================================================

func (h *Handler) CreateRun(c *gin.Context) {
	var req CreateRunRequest
	if !h.bindJSON(c, &req) {
		return
	}
	run, err := h.svc.CreateRun(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, run)
}

func (h *Handler) GetRun(c *gin.Context) {
	run, err := h.svc.GetRun(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "run not found"})
		return
	}
	c.JSON(http.StatusOK, run)
}

func (h *Handler) ListRuns(c *gin.Context) {
	filter := RunFilter{
		PeriodID: qs(c, "period_id"),
		RunType:  qs(c, "run_type"),
		Status:   qs(c, "status"),
		Limit:    qi(c, "limit", 20),
		Offset:   qi(c, "offset", 0),
	}
	list, err := h.svc.ListRuns(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) CalculateRun(c *gin.Context) {
	if err := h.svc.CalculateRun(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "calculation started"})
}

func (h *Handler) ValidateRun(c *gin.Context) {
	if err := h.svc.ValidateRun(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "validation passed"})
}

func (h *Handler) ApproveRun(c *gin.Context) {
	if err := h.svc.ApproveRun(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "run approved"})
}

func (h *Handler) CloseRun(c *gin.Context) {
	if err := h.svc.CloseRun(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "run closed"})
}

func (h *Handler) GetRunSummary(c *gin.Context) {
	s, err := h.svc.GetRunSummary(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

// ========================================================================
// RUN EMPLOYEES
// ========================================================================

func (h *Handler) AddEmployeeToRun(c *gin.Context) {
	var req struct {
		EmployeeID string `json:"employee_id" binding:"required"`
	}
	if !h.bindJSON(c, &req) {
		return
	}
	re, err := h.svc.AddEmployeeToRun(c.Request.Context(), h.companyID(c), c.Param("id"), req.EmployeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, re)
}

func (h *Handler) ListRunEmployees(c *gin.Context) {
	list, err := h.svc.ListRunEmployees(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetEmployeeResult(c *gin.Context) {
	res, err := h.svc.GetEmployeeResult(c.Request.Context(), h.companyID(c), c.Param("id"), c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee result not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetEmployeeItems(c *gin.Context) {
	items, err := h.svc.GetEmployeeItems(c.Request.Context(), h.companyID(c), c.Param("id"), c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "items not found"})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *Handler) GetEmployeeBases(c *gin.Context) {
	re, err := h.svc.GetRunEmployee(c.Request.Context(), h.companyID(c), c.Param("id"), c.Param("eid"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee not found"})
		return
	}
	bases, err := h.svc.GetBases(c.Request.Context(), re.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, bases)
}

// ========================================================================
// CONCEPTS
// ========================================================================

func (h *Handler) CreateConcept(c *gin.Context) {
	var req CreateConceptRequest
	if !h.bindJSON(c, &req) {
		return
	}
	concept, err := h.svc.CreateConcept(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, concept)
}

func (h *Handler) UpdateConcept(c *gin.Context) {
	var req UpdateConceptRequest
	if !h.bindJSON(c, &req) {
		return
	}
	concept, err := h.svc.UpdateConcept(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, concept)
}

func (h *Handler) ListConcepts(c *gin.Context) {
	filter := ConceptFilter{}
	if v := qs(c, "active"); v != nil {
		b := *v == "true"
		filter.Active = &b
	}
	filter.ConceptType = qs(c, "concept_type")
	filter.Taxability = qs(c, "taxability")
	list, err := h.svc.repo.ListConcepts(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// RULES
// ========================================================================

func (h *Handler) CreateRule(c *gin.Context) {
	var req CreateRuleRequest
	if !h.bindJSON(c, &req) {
		return
	}
	rule, err := h.svc.CreateRule(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	var req UpdateRuleRequest
	if !h.bindJSON(c, &req) {
		return
	}
	rule, err := h.svc.UpdateRule(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rule)
}

func (h *Handler) ListRules(c *gin.Context) {
	list, err := h.svc.repo.ListRules(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// NOVELTIES
// ========================================================================

func (h *Handler) CreateNovelty(c *gin.Context) {
	var req CreateNoveltyRequest
	if !h.bindJSON(c, &req) {
		return
	}
	n, err := h.svc.CreateNovelty(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, n)
}

func (h *Handler) UpdateNovelty(c *gin.Context) {
	var req UpdateNoveltyRequest
	if !h.bindJSON(c, &req) {
		return
	}
	n, err := h.svc.UpdateNovelty(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, n)
}

func (h *Handler) ListNovelties(c *gin.Context) {
	filter := NoveltyFilter{
		EmployeeID:  qs(c, "employee_id"),
		PeriodID:    qs(c, "period_id"),
		NoveltyType: qs(c, "novelty_type"),
		Status:      qs(c, "status"),
		Source:      qs(c, "source"),
		Limit:       qi(c, "limit", 20),
		Offset:      qi(c, "offset", 0),
	}
	list, err := h.svc.repo.ListNovelties(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) DeleteNovelty(c *gin.Context) {
	if err := h.svc.repo.DeleteNovelty(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "novelty deleted"})
}

func (h *Handler) ApproveNovelty(c *gin.Context) {
	if err := h.svc.ApproveNovelty(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "novelty approved"})
}

func (h *Handler) ImportNovelties(c *gin.Context) {
	var req ImportNoveltiesRequest
	if !h.bindJSON(c, &req) {
		return
	}
	novelties, err := h.svc.ImportNovelties(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, novelties)
}

// ========================================================================
// ADVANCES
// ========================================================================

func (h *Handler) CreateAdvance(c *gin.Context) {
	var req CreateAdvanceRequest
	if !h.bindJSON(c, &req) {
		return
	}
	a, err := h.svc.CreateAdvance(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) ListAdvances(c *gin.Context) {
	list, err := h.svc.repo.ListAdvances(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// GARNISHMENTS
// ========================================================================

func (h *Handler) CreateGarnishment(c *gin.Context) {
	var req CreateGarnishmentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	g, err := h.svc.CreateGarnishment(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *Handler) ListGarnishments(c *gin.Context) {
	list, err := h.svc.repo.ListGarnishments(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// AGREEMENTS & CATEGORIES
// ========================================================================

func (h *Handler) CreateAgreement(c *gin.Context) {
	var req CreateAgreementRequest
	if !h.bindJSON(c, &req) {
		return
	}
	a, err := h.svc.CreateAgreement(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) ListAgreements(c *gin.Context) {
	list, err := h.svc.ListAgreements(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if !h.bindJSON(c, &req) {
		return
	}
	cat, err := h.svc.CreateCategory(c.Request.Context(), h.companyID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func (h *Handler) ListCategories(c *gin.Context) {
	list, err := h.svc.ListCategories(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// SALARY SCALES
// ========================================================================

func (h *Handler) CreateSalaryScale(c *gin.Context) {
	var req CreateSalaryScaleRequest
	if !h.bindJSON(c, &req) {
		return
	}
	sc, err := h.svc.CreateSalaryScale(c.Request.Context(), h.companyID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sc)
}

func (h *Handler) ListSalaryScales(c *gin.Context) {
	list, err := h.svc.ListSalaryScales(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// DASHBOARD & ERRORS
// ========================================================================

func (h *Handler) GetDashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func (h *Handler) ListErrors(c *gin.Context) {
	list, err := h.svc.GetErrors(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// ========================================================================
// EMPLOYEE SELF-SERVICE
// ========================================================================

func (h *Handler) MyPeriods(c *gin.Context) {
	empID := h.employeeID(c)
	if empID == "" {
		empID = h.userID(c)
	}
	list, err := h.svc.GetMyPeriods(c.Request.Context(), h.companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) MyRunResult(c *gin.Context) {
	empID := h.employeeID(c)
	if empID == "" {
		empID = h.userID(c)
	}
	res, err := h.svc.GetEmployeeResult(c.Request.Context(), h.companyID(c), c.Param("run_id"), empID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "result not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) MyItems(c *gin.Context) {
	empID := h.employeeID(c)
	if empID == "" {
		empID = h.userID(c)
	}
	items, err := h.svc.GetMyItems(c.Request.Context(), h.companyID(c), empID, c.Param("run_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "items not found"})
		return
	}
	c.JSON(http.StatusOK, items)
}
