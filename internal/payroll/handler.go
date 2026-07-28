package payroll

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

// Periods
func (h *Handler) CreatePeriod(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	period, err := h.service.CreatePeriod(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, period)
}

func (h *Handler) GetPeriod(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	period, err := h.service.GetPeriod(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Period not found")
		return
	}
	response.Success(c, period)
}

func (h *Handler) ListPeriods(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	periods, err := h.service.ListPeriods(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, periods)
}

func (h *Handler) UpdatePeriod(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req UpdatePeriodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	period, err := h.service.UpdatePeriod(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, period)
}

func (h *Handler) CalculatePeriod(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	review, err := h.service.CalculatePeriod(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, review)
}

func (h *Handler) GetReview(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	review, err := h.service.GetReview(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Period not found")
		return
	}
	response.Success(c, review)
}

func (h *Handler) ApprovePeriod(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.ApprovePeriod(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "period approved"})
}

func (h *Handler) ClosePeriod(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.ClosePeriod(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "period closed"})
}

// Concepts
func (h *Handler) CreateConcept(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateConceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	concept, err := h.service.CreateConcept(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, concept)
}

func (h *Handler) ListConcepts(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	concepts, err := h.service.ListConcepts(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, concepts)
}

// Compensation
func (h *Handler) SetCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req SetCompensationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	comp, err := h.service.SetCompensation(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, comp)
}

func (h *Handler) GetCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	comp, err := h.service.GetCompensation(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "No active compensation")
		return
	}
	response.Success(c, comp)
}

func (h *Handler) GetCompensationHistory(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	history, err := h.service.GetCompensationHistory(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, history)
}

// Benefits
func (h *Handler) CreateBenefit(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateBenefitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	benefit, err := h.service.CreateBenefit(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, benefit)
}

func (h *Handler) ListBenefits(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	benefits, err := h.service.ListBenefits(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, benefits)
}

func (h *Handler) AssignBenefit(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req AssignBenefitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	eb, err := h.service.AssignBenefit(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, eb)
}

func (h *Handler) GetEmployeeBenefits(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	benefits, err := h.service.GetEmployeeBenefits(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, benefits)
}

// Bonuses
func (h *Handler) CreateBonus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateBonusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	bonus, err := h.service.CreateBonus(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, bonus)
}

func (h *Handler) ListBonuses(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := PayrollFilters{
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
	}
	bonuses, err := h.service.ListBonuses(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, bonuses)
}

func (h *Handler) ApproveBonus(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.ApproveBonus(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "bonus approved"})
}

// Advances
func (h *Handler) CreateAdvance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateAdvanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	advance, err := h.service.CreateAdvance(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, advance)
}

func (h *Handler) ListAdvances(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	filters := PayrollFilters{
		EmployeeID: c.Query("employee_id"),
		Status:     c.Query("status"),
	}
	advances, err := h.service.ListAdvances(c.Request.Context(), companyID, filters)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, advances)
}

func (h *Handler) ApproveAdvance(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	if err := h.service.ApproveAdvance(c.Request.Context(), companyID, c.Param("id"), userID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "advance approved"})
}

// Deductions
func (h *Handler) CreateDeduction(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req CreateDeductionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	deduction, err := h.service.CreateDeduction(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, deduction)
}

// Adjustments
func (h *Handler) CreateAdjustment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	periodID := c.Param("id")
	var req CreateAdjustmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	adj, err := h.service.CreateAdjustment(c.Request.Context(), companyID, periodID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, adj)
}

// Employee payroll
func (h *Handler) GetEmployeePayroll(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")
	periodID := c.Query("period_id")
	if periodID == "" {
		response.BadRequest(c, "period_id is required")
		return
	}
	summary, err := h.service.GetEmployeePayroll(c.Request.Context(), companyID, employeeID, periodID)
	if err != nil {
		response.NotFound(c, "No payroll data")
		return
	}
	response.Success(c, summary)
}

func (h *Handler) GetEmployeeCompensation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	comp, err := h.service.GetCompensation(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "No active compensation")
		return
	}
	response.Success(c, comp)
}

func (h *Handler) GetEmployeeBenefitsView(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	benefits, err := h.service.GetEmployeeBenefits(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, benefits)
}

// Snapshot
func (h *Handler) GetSnapshot(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	snapshot, err := h.service.GetSnapshot(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "No snapshot found")
		return
	}
	response.Success(c, snapshot)
}

// Dashboard
func (h *Handler) GetDashboard(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	periodID := c.Param("id")
	dash, err := h.service.GetDashboard(c.Request.Context(), companyID, periodID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    dash,
	})
}
