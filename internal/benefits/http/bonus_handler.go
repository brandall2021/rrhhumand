package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

func (h *Handler) CreateBonus(c *gin.Context) {
	var req domain.EmployeeBonus
	if !bindJSON(c, &req) {
		return
	}
	b, err := h.BonusSvc.CreateBonus(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, b)
}

func (h *Handler) ListBonuses(c *gin.Context) {
	var employeeID *uuid.UUID
	if eid := c.Query("employee_id"); eid != "" {
		id, err := uuid.Parse(eid)
		if err == nil {
			employeeID = &id
		}
	}
	status := qs(c, "status")
	bonuses, err := h.BonusSvc.ListBonuses(c.Request.Context(), companyID(c), employeeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, bonuses)
}

func (h *Handler) GetBonus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	b, err := h.BonusSvc.GetBonus(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, b)
}

func (h *Handler) UpdateBonus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.EmployeeBonus
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	b, err := h.BonusSvc.UpdateBonus(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, b)
}

func (h *Handler) ApproveBonus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.BonusSvc.ApproveBonus(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "bonus approved"})
}

func (h *Handler) PayBonus(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		PayrollRunID uuid.UUID `json:"payroll_run_id" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.BonusSvc.PayBonus(c.Request.Context(), id, req.PayrollRunID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "bonus paid"})
}

func (h *Handler) CancelBonus(c *gin.Context) {
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
	if err := h.BonusSvc.CancelBonus(c.Request.Context(), id, req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "bonus cancelled"})
}

func (h *Handler) CreateIncentive(c *gin.Context) {
	var req domain.EmployeeIncentive
	if !bindJSON(c, &req) {
		return
	}
	empID := employeeID(c)
	i, err := h.BonusSvc.CreateIncentive(c.Request.Context(), companyID(c), &empID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, i)
}

func (h *Handler) ListIncentives(c *gin.Context) {
	var employeeID *uuid.UUID
	if eid := c.Query("employee_id"); eid != "" {
		id, err := uuid.Parse(eid)
		if err == nil {
			employeeID = &id
		}
	}
	status := qs(c, "status")
	incentives, err := h.BonusSvc.ListIncentives(c.Request.Context(), companyID(c), employeeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, incentives)
}

func (h *Handler) GetIncentive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	i, err := h.BonusSvc.GetIncentive(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, i)
}

func (h *Handler) RedeemIncentive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.BonusSvc.RedeemIncentive(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "incentive redeemed"})
}

func (h *Handler) CreatePayrollMapping(c *gin.Context) {
	var req domain.BenefitPayrollMapping
	if !bindJSON(c, &req) {
		return
	}
	m, err := h.BonusSvc.CreatePayrollMapping(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, m)
}

func (h *Handler) ListPayrollMappings(c *gin.Context) {
	var benefitID *uuid.UUID
	if bid := c.Query("benefit_id"); bid != "" {
		id, err := uuid.Parse(bid)
		if err == nil {
			benefitID = &id
		}
	}
	mappingType := qs(c, "mapping_type")
	mappings, err := h.BonusSvc.ListPayrollMappings(c.Request.Context(), companyID(c), benefitID, mappingType)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, mappings)
}

func (h *Handler) SyncToPayroll(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.BonusSvc.SyncToPayroll(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "synced to payroll"})
}
