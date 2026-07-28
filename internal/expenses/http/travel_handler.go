package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateTravel(c *gin.Context) {
	var req domain.Travel
	if !bindJSON(c, &req) {
		return
	}
	travel, err := h.TravelSvc.CreateTravel(c.Request.Context(), companyID(c), employeeID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, travel)
}

func (h *Handler) ListTravels(c *gin.Context) {
	empID := qs(c, "employee_id")
	var eid *uuid.UUID
	if empID != nil {
		id, err := uuid.Parse(*empID)
		if err == nil {
			eid = &id
		}
	}
	status := qs(c, "status")
	travels, err := h.TravelSvc.ListTravels(c.Request.Context(), companyID(c), eid, status, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, travels)
}

func (h *Handler) GetTravel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	travel, err := h.TravelSvc.GetTravel(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, travel)
}

func (h *Handler) UpdateTravel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.Travel
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	travel, err := h.TravelSvc.UpdateTravel(c.Request.Context(), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, travel)
}

func (h *Handler) RequestTravel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.TravelSvc.RequestTravel(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	travel, _ := h.TravelSvc.GetTravel(c.Request.Context(), id)
	success(c, travel)
}

func (h *Handler) ApproveTravel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.TravelSvc.ApproveTravel(c.Request.Context(), id, userID(c), ""); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	travel, _ := h.TravelSvc.GetTravel(c.Request.Context(), id)
	success(c, travel)
}

func (h *Handler) RejectTravel(c *gin.Context) {
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
	if err := h.TravelSvc.RejectTravel(c.Request.Context(), id, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	travel, _ := h.TravelSvc.GetTravel(c.Request.Context(), id)
	success(c, travel)
}

func (h *Handler) CompleteTravel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.TravelSvc.CompleteTravel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	travel, _ := h.TravelSvc.GetTravel(c.Request.Context(), id)
	success(c, travel)
}

func (h *Handler) CancelTravel(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.TravelSvc.CancelTravel(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	travel, _ := h.TravelSvc.GetTravel(c.Request.Context(), id)
	success(c, travel)
}

func (h *Handler) AddParticipant(c *gin.Context) {
	travelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid travel id"})
		return
	}
	var req struct {
		EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
		Role       string    `json:"role"`
	}
	if !bindJSON(c, &req) {
		return
	}
	p, err := h.TravelSvc.AddParticipant(c.Request.Context(), travelID, req.EmployeeID, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, p)
}

func (h *Handler) RemoveParticipant(c *gin.Context) {
	travelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid travel id"})
		return
	}
	participantID, err := uuid.Parse(c.Param("participant_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid participant id"})
		return
	}
	if err := h.TravelSvc.RemoveParticipant(c.Request.Context(), travelID, participantID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ListParticipants(c *gin.Context) {
	travelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid travel id"})
		return
	}
	participants, err := h.TravelSvc.ListParticipants(c.Request.Context(), travelID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, participants)
}
