package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreatePolicy(c *gin.Context) {
	var req domain.ExpensePolicy
	if !bindJSON(c, &req) {
		return
	}
	policy, err := h.PolicySvc.CreatePolicy(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, policy)
}

func (h *Handler) ListPolicies(c *gin.Context) {
	policies, err := h.PolicySvc.ListPolicies(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, policies)
}

func (h *Handler) GetPolicy(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	policy, err := h.PolicySvc.GetPolicy(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, policy)
}

func (h *Handler) UpdatePolicy(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ExpensePolicy
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	policy, err := h.PolicySvc.UpdatePolicy(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, policy)
}

func (h *Handler) CreateRule(c *gin.Context) {
	var req domain.ExpensePolicyRule
	if !bindJSON(c, &req) {
		return
	}
	rule, err := h.PolicySvc.CreateRule(c.Request.Context(), req.PolicyID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, rule)
}

func (h *Handler) ListRules(c *gin.Context) {
	policyID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}
	rules, err := h.PolicySvc.ListRules(c.Request.Context(), policyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rules)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}
	var req domain.ExpensePolicyRule
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	rule, err := h.PolicySvc.UpdateRule(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rule)
}

func (h *Handler) DeleteRule(c *gin.Context) {
	policyID, err := uuid.Parse(c.Param("policyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid policy id"})
		return
	}
	ruleID, err := uuid.Parse(c.Param("ruleId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rule id"})
		return
	}
	if err := h.PolicySvc.DeleteRule(c.Request.Context(), policyID, ruleID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) EvaluateExpense(c *gin.Context) {
	var req struct {
		Expense          domain.Expense `json:"expense"`
		EmployeeCategory string         `json:"employee_category"`
	}
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.PolicySvc.EvaluateExpense(c.Request.Context(), companyID(c), &req.Expense, req.EmployeeCategory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, result)
}

func (h *Handler) OverridePolicy(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.PolicySvc.OverridePolicy(c.Request.Context(), id, req.Reason, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "policy override applied"})
}
