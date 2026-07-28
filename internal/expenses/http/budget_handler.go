package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateBudget(c *gin.Context) {
	var req domain.ExpenseBudget
	if !bindJSON(c, &req) {
		return
	}
	budget, err := h.BudgetSvc.CreateBudget(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, budget)
}

func (h *Handler) ListBudgets(c *gin.Context) {
	fiscalYear := qi(c, "fiscal_year", 0)
	var fy *int
	if fiscalYear > 0 {
		fy = &fiscalYear
	}
	budgets, err := h.BudgetSvc.ListBudgets(c.Request.Context(), companyID(c), fy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, budgets)
}

func (h *Handler) GetBudget(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	budget, err := h.BudgetSvc.GetBudget(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, budget)
}

func (h *Handler) CheckAvailability(c *gin.Context) {
	var req struct {
		Amount      decimal.Decimal `json:"amount" binding:"required"`
		CostCenterID *uuid.UUID     `json:"cost_center_id"`
		ProjectID   *uuid.UUID     `json:"project_id"`
		CategoryID  *uuid.UUID     `json:"category_id"`
	}
	if !bindJSON(c, &req) {
		return
	}
	ok, remaining, err := h.BudgetSvc.CheckAvailability(c.Request.Context(), companyID(c), req.Amount, req.CostCenterID, req.ProjectID, req.CategoryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"available": ok, "remaining": remaining})
}
