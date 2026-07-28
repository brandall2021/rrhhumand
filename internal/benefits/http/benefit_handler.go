package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
)

func (h *Handler) CreateBenefit(c *gin.Context) {
	var req domain.Benefit
	if !bindJSON(c, &req) {
		return
	}
	b, err := h.BenefitSvc.CreateBenefit(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, b)
}

func (h *Handler) ListBenefits(c *gin.Context) {
	var typeID *uuid.UUID
	if tid := c.Query("type_id"); tid != "" {
		id, err := uuid.Parse(tid)
		if err == nil {
			typeID = &id
		}
	}
	status := qs(c, "status")
	visibility := qs(c, "visibility")
	limit := qi(c, "limit", 0)
	offset := qi(c, "offset", 0)
	benefits, err := h.BenefitSvc.ListBenefits(c.Request.Context(), companyID(c), status, typeID, visibility, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, benefits)
}

func (h *Handler) GetBenefit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	b, err := h.BenefitSvc.GetBenefit(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, b)
}

func (h *Handler) UpdateBenefit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.Benefit
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	b, err := h.BenefitSvc.UpdateBenefit(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, b)
}

func (h *Handler) DeleteBenefit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.BenefitSvc.DeleteBenefit(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) SearchBenefits(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query required"})
		return
	}
	benefits, err := h.BenefitSvc.SearchBenefits(c.Request.Context(), companyID(c), query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, benefits)
}

func (h *Handler) CreatePlan(c *gin.Context) {
	benefitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid benefit id"})
		return
	}
	var req domain.BenefitPlan
	if !bindJSON(c, &req) {
		return
	}
	req.BenefitID = benefitID
	p, err := h.BenefitSvc.CreatePlan(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, p)
}

func (h *Handler) ListPlans(c *gin.Context) {
	benefitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid benefit id"})
		return
	}
	plans, err := h.BenefitSvc.ListPlans(c.Request.Context(), companyID(c), benefitID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, plans)
}

func (h *Handler) GetPlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	p, err := h.BenefitSvc.GetPlan(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, p)
}

func (h *Handler) UpdatePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	var req domain.BenefitPlan
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	p, err := h.BenefitSvc.UpdatePlan(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, p)
}

func (h *Handler) DeletePlan(c *gin.Context) {
	id, err := uuid.Parse(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	if err := h.BenefitSvc.DeletePlan(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
