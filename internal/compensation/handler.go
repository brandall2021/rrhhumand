package compensation

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

// --- Structures ---
func (h *Handler) CreateStructure(c *gin.Context) {
	var req CreateStructureRequest
	if !h.bindJSON(c, &req) {
		return
	}
	st, err := h.svc.CreateStructure(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, st)
}

func (h *Handler) UpdateStructure(c *gin.Context) {
	var req UpdateStructureRequest
	if !h.bindJSON(c, &req) {
		return
	}
	st, err := h.svc.UpdateStructure(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *Handler) GetStructure(c *gin.Context) {
	st, err := h.svc.GetStructure(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "structure not found"})
		return
	}
	c.JSON(http.StatusOK, st)
}

func (h *Handler) ListStructures(c *gin.Context) {
	list, err := h.svc.ListStructures(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Grades ---
func (h *Handler) CreateGrade(c *gin.Context) {
	var req CreateGradeRequest
	if !h.bindJSON(c, &req) {
		return
	}
	g, err := h.svc.CreateGrade(c.Request.Context(), h.companyID(c), c.Param("structure_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, g)
}

func (h *Handler) UpdateGrade(c *gin.Context) {
	var req UpdateGradeRequest
	if !h.bindJSON(c, &req) {
		return
	}
	g, err := h.svc.UpdateGrade(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, g)
}

func (h *Handler) ListGrades(c *gin.Context) {
	list, err := h.svc.ListGrades(c.Request.Context(), h.companyID(c), c.Param("structure_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Bands ---
func (h *Handler) CreateBand(c *gin.Context) {
	var req CreateBandRequest
	if !h.bindJSON(c, &req) {
		return
	}
	b, err := h.svc.CreateBand(c.Request.Context(), h.companyID(c), c.Param("structure_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *Handler) UpdateBand(c *gin.Context) {
	var req UpdateBandRequest
	if !h.bindJSON(c, &req) {
		return
	}
	b, err := h.svc.UpdateBand(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) GetBand(c *gin.Context) {
	b, err := h.svc.GetBand(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "band not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) ListBands(c *gin.Context) {
	list, err := h.svc.ListBands(c.Request.Context(), h.companyID(c), c.Param("structure_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Position-Band ---
func (h *Handler) AssignPositionBand(c *gin.Context) {
	var req AssignPositionBandRequest
	if !h.bindJSON(c, &req) {
		return
	}
	pb, err := h.svc.AssignPositionBand(c.Request.Context(), c.Param("position_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pb)
}

func (h *Handler) GetPositionBand(c *gin.Context) {
	pb, err := h.svc.GetPositionBand(c.Request.Context(), c.Param("position_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "position band not found"})
		return
	}
	c.JSON(http.StatusOK, pb)
}

// --- Employee Compensation ---
func (h *Handler) SetEmployeeCompensation(c *gin.Context) {
	var req SetEmployeeCompensationRequest
	if !h.bindJSON(c, &req) {
		return
	}
	ec, err := h.svc.SetEmployeeCompensation(c.Request.Context(), h.companyID(c), c.Param("employee_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ec)
}

func (h *Handler) GetEmployeeCompensation(c *gin.Context) {
	ec, err := h.svc.GetEmployeeCompensation(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "compensation not found"})
		return
	}
	c.JSON(http.StatusOK, ec)
}

func (h *Handler) ListEmployeeCompensations(c *gin.Context) {
	list, err := h.svc.ListEmployeeCompensations(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) GetHistory(c *gin.Context) {
	list, err := h.svc.GetHistory(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Components ---
func (h *Handler) CreateComponent(c *gin.Context) {
	var req CreateComponentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	comp, err := h.svc.CreateComponent(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, comp)
}

func (h *Handler) ListComponents(c *gin.Context) {
	list, err := h.svc.ListComponents(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) AssignComponent(c *gin.Context) {
	var req AssignComponentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	ecc, err := h.svc.AssignComponent(c.Request.Context(), h.companyID(c), c.Param("employee_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ecc)
}

func (h *Handler) ListEmployeeComponents(c *gin.Context) {
	list, err := h.svc.ListEmployeeComponents(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Adjustments ---
func (h *Handler) CreateAdjustment(c *gin.Context) {
	var req CreateAdjustmentRequest
	if !h.bindJSON(c, &req) {
		return
	}
	a, err := h.svc.CreateAdjustment(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, a)
}

func (h *Handler) GetAdjustment(c *gin.Context) {
	a, err := h.svc.GetAdjustment(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "adjustment not found"})
		return
	}
	c.JSON(http.StatusOK, a)
}

func (h *Handler) ListAdjustments(c *gin.Context) {
	filter := AdjustmentFilter{
		EmployeeID: qs(c, "employee_id"),
		Status:     qs(c, "status"),
		Limit:      qi(c, "limit", 20),
		Offset:     qi(c, "offset", 0),
	}
	list, err := h.svc.ListAdjustments(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ApproveAdjustment(c *gin.Context) {
	if err := h.svc.ApproveAdjustment(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "adjustment approved"})
}

func (h *Handler) RejectAdjustment(c *gin.Context) {
	if err := h.svc.RejectAdjustment(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "adjustment rejected"})
}

func (h *Handler) ApplyAdjustment(c *gin.Context) {
	if err := h.svc.ApplyAdjustment(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "adjustment applied"})
}

// --- Proposals ---
func (h *Handler) CreateProposal(c *gin.Context) {
	var req CreateProposalRequest
	if !h.bindJSON(c, &req) {
		return
	}
	p, err := h.svc.CreateProposal(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *Handler) SubmitProposal(c *gin.Context) {
	if err := h.svc.SubmitProposal(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "proposal submitted"})
}

func (h *Handler) ApproveProposal(c *gin.Context) {
	if err := h.svc.ApproveProposal(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "proposal approved"})
}

func (h *Handler) RejectProposal(c *gin.Context) {
	var req struct {
		Reason string `json:"reason"`
	}
	h.bindJSON(c, &req)
	if err := h.svc.RejectProposal(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "proposal rejected"})
}

func (h *Handler) ListProposals(c *gin.Context) {
	filter := ProposalFilter{
		ReviewID:   qs(c, "review_id"),
		EmployeeID: qs(c, "employee_id"),
		Status:     qs(c, "status"),
		Limit:      qi(c, "limit", 20),
		Offset:     qi(c, "offset", 0),
	}
	list, err := h.svc.ListProposals(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Bonuses ---
func (h *Handler) CreateBonus(c *gin.Context) {
	var req CreateBonusRequest
	if !h.bindJSON(c, &req) {
		return
	}
	b, err := h.svc.CreateBonus(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *Handler) GetBonus(c *gin.Context) {
	b, err := h.svc.GetBonus(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "bonus not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) ListBonuses(c *gin.Context) {
	filter := BonusFilter{
		EmployeeID: qs(c, "employee_id"),
		Status:     qs(c, "status"),
		Limit:      qi(c, "limit", 20),
		Offset:     qi(c, "offset", 0),
	}
	list, err := h.svc.ListBonuses(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) ApproveBonus(c *gin.Context) {
	if err := h.svc.ApproveBonus(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bonus approved"})
}

func (h *Handler) RejectBonus(c *gin.Context) {
	if err := h.svc.RejectBonus(c.Request.Context(), h.companyID(c), c.Param("id"), h.userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "bonus rejected"})
}

// --- Bonus Plans ---
func (h *Handler) CreateBonusPlan(c *gin.Context) {
	var req CreateBonusPlanRequest
	if !h.bindJSON(c, &req) {
		return
	}
	bp, err := h.svc.CreateBonusPlan(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, bp)
}

func (h *Handler) ListBonusPlans(c *gin.Context) {
	list, err := h.svc.ListBonusPlans(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Benefits ---
func (h *Handler) CreateBenefit(c *gin.Context) {
	var req CreateBenefitRequest
	if !h.bindJSON(c, &req) {
		return
	}
	b, err := h.svc.CreateBenefit(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *Handler) UpdateBenefit(c *gin.Context) {
	var req UpdateBenefitRequest
	if !h.bindJSON(c, &req) {
		return
	}
	b, err := h.svc.UpdateBenefit(c.Request.Context(), h.companyID(c), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) GetBenefit(c *gin.Context) {
	b, err := h.svc.GetBenefit(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "benefit not found"})
		return
	}
	c.JSON(http.StatusOK, b)
}

func (h *Handler) ListBenefits(c *gin.Context) {
	filter := BenefitFilter{
		Active:      nil,
		BenefitType: qs(c, "benefit_type"),
	}
	if v := qs(c, "active"); v != nil {
		b := *v == "true" || *v == "1"
		filter.Active = &b
	}
	list, err := h.svc.ListBenefits(c.Request.Context(), h.companyID(c), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Employee Benefits ---
func (h *Handler) AssignBenefit(c *gin.Context) {
	var req AssignBenefitRequest
	if !h.bindJSON(c, &req) {
		return
	}
	eb, err := h.svc.AssignBenefit(c.Request.Context(), h.companyID(c), c.Param("employee_id"), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, eb)
}

func (h *Handler) ListEmployeeBenefits(c *gin.Context) {
	list, err := h.svc.ListEmployeeBenefits(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) RemoveEmployeeBenefit(c *gin.Context) {
	if err := h.svc.RemoveEmployeeBenefit(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "benefit removed"})
}

// --- Reviews ---
func (h *Handler) CreateReview(c *gin.Context) {
	var req CreateReviewRequest
	if !h.bindJSON(c, &req) {
		return
	}
	rv, err := h.svc.CreateReview(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, rv)
}

func (h *Handler) GetReview(c *gin.Context) {
	rv, err := h.svc.GetReview(c.Request.Context(), h.companyID(c), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "review not found"})
		return
	}
	c.JSON(http.StatusOK, rv)
}

func (h *Handler) ListReviews(c *gin.Context) {
	list, err := h.svc.ListReviews(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) OpenReview(c *gin.Context) {
	if err := h.svc.OpenReview(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "review opened"})
}

func (h *Handler) CloseReview(c *gin.Context) {
	if err := h.svc.CloseReview(c.Request.Context(), h.companyID(c), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "review closed"})
}

// --- Budgets ---
func (h *Handler) CreateBudget(c *gin.Context) {
	var req CreateBudgetRequest
	if !h.bindJSON(c, &req) {
		return
	}
	b, err := h.svc.CreateBudget(c.Request.Context(), h.companyID(c), h.userID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, b)
}

func (h *Handler) ListBudgets(c *gin.Context) {
	list, err := h.svc.ListBudgets(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

// --- Dashboard ---
func (h *Handler) DashboardStats(c *gin.Context) {
	stats, err := h.svc.GetDashboardStats(c.Request.Context(), h.companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// --- Reports ---
func (h *Handler) BandAnalysisReport(c *gin.Context) {
	analysis, err := h.svc.GetBandAnalysis(c.Request.Context(), h.companyID(c), c.Param("band_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, analysis)
}

func (h *Handler) CompaRatioReport(c *gin.Context) {
	ec, err := h.svc.GetEmployeeCompensation(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "compensation not found"})
		return
	}
	if ec.SalaryBandID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee has no salary band"})
		return
	}
	band, err := h.svc.GetBand(c.Request.Context(), h.companyID(c), *ec.SalaryBandID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "band not found"})
		return
	}
	result, err := h.svc.CalculateCompaRatio(ec.BaseAmount, band.MidpointAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) RangePenetrationReport(c *gin.Context) {
	ec, err := h.svc.GetEmployeeCompensation(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "compensation not found"})
		return
	}
	if ec.SalaryBandID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "employee has no salary band"})
		return
	}
	band, err := h.svc.GetBand(c.Request.Context(), h.companyID(c), *ec.SalaryBandID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "band not found"})
		return
	}
	result, err := h.svc.CalculateRangePenetration(ec.BaseAmount, band.MinimumAmount, band.MaximumAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *Handler) TotalCompensationReport(c *gin.Context) {
	tc, err := h.svc.CalculateTotalCompensation(c.Request.Context(), h.companyID(c), c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tc)
}

func (h *Handler) EquityReport(c *gin.Context) {
	var req EquityAnalysisRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.AnalyzeEquity(c.Request.Context(), h.companyID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// --- AI ---
func (h *Handler) GenerateAIRecommendation(c *gin.Context) {
	var req AIRecommendationRequest
	if !h.bindJSON(c, &req) {
		return
	}
	rec, err := h.svc.GenerateAIRecommendation(c.Request.Context(), h.companyID(c), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, rec)
}

// --- Employee Self-Service ---
func (h *Handler) MyCompensation(c *gin.Context) {
	employeeID := c.GetString("employee_id")
	if employeeID == "" {
		employeeID = c.GetString("user_id")
	}
	ec, err := h.svc.GetEmployeeCompensation(c.Request.Context(), h.companyID(c), employeeID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "compensation not found"})
		return
	}
	c.JSON(http.StatusOK, ec)
}

func (h *Handler) MyHistory(c *gin.Context) {
	employeeID := c.GetString("employee_id")
	if employeeID == "" {
		employeeID = c.GetString("user_id")
	}
	list, err := h.svc.GetHistory(c.Request.Context(), h.companyID(c), employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) MyBenefits(c *gin.Context) {
	employeeID := c.GetString("employee_id")
	if employeeID == "" {
		employeeID = c.GetString("user_id")
	}
	list, err := h.svc.ListEmployeeBenefits(c.Request.Context(), h.companyID(c), employeeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}

func (h *Handler) MyBonuses(c *gin.Context) {
	employeeID := c.GetString("employee_id")
	if employeeID == "" {
		employeeID = c.GetString("user_id")
	}
	list, err := h.svc.ListBonuses(c.Request.Context(), h.companyID(c), BonusFilter{EmployeeID: &employeeID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
