package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

func (h *Handler) CreateRule(c *gin.Context) {
	var req domain.BenefitEligibilityRule
	if !bindJSON(c, &req) {
		return
	}
	rule, err := h.EligibilitySvc.CreateRule(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, rule)
}

func (h *Handler) ListRules(c *gin.Context) {
	benefitID, err := uuid.Parse(c.Param("benefitId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid benefit id"})
		return
	}
	rules, err := h.EligibilitySvc.ListRules(c.Request.Context(), companyID(c), benefitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rules)
}

func (h *Handler) GetRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	rule, err := h.EligibilitySvc.GetRule(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, rule)
}

func (h *Handler) UpdateRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.BenefitEligibilityRule
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	rule, err := h.EligibilitySvc.UpdateRule(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, rule)
}

func (h *Handler) DeleteRule(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.EligibilitySvc.DeleteRule(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) EvaluateEmployee(c *gin.Context) {
	var req struct {
		EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
		BenefitID  uuid.UUID `json:"benefit_id" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	eligible, reasons, err := h.EligibilitySvc.EvaluateEmployee(c.Request.Context(), companyID(c), req.EmployeeID, req.BenefitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{
		"eligible": eligible,
		"reasons":  reasons,
	})
}

func (h *Handler) ListEligibleBenefits(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	benefits, err := h.EligibilitySvc.ListEligibleBenefits(c.Request.Context(), companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, benefits)
}
