package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) EnrollEmployee(c *gin.Context) {
	var req struct {
		EmployeeID uuid.UUID `json:"employee_id" binding:"required"`
		BenefitID  uuid.UUID `json:"benefit_id" binding:"required"`
		PlanID     *uuid.UUID `json:"plan_id"`
		Source     string    `json:"source"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Source == "" {
		req.Source = "MANUAL"
	}
	cid := companyID(c)
	eb, err := h.AssignmentSvc.EnrollEmployee(c.Request.Context(), &cid, &req.EmployeeID, &req.BenefitID, req.PlanID, userID(c), req.Source)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, eb)
}

func (h *Handler) ListEmployeeBenefits(c *gin.Context) {
	var employeeID, benefitID *uuid.UUID
	if eid := c.Query("employee_id"); eid != "" {
		id, err := uuid.Parse(eid)
		if err == nil {
			employeeID = &id
		}
	}
	if bid := c.Query("benefit_id"); bid != "" {
		id, err := uuid.Parse(bid)
		if err == nil {
			benefitID = &id
		}
	}
	status := qs(c, "status")
	cid := companyID(c)
	benefits, err := h.AssignmentSvc.ListEmployeeBenefits(c.Request.Context(), &cid, employeeID, benefitID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, benefits)
}

func (h *Handler) GetEmployeeBenefit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	eb, err := h.AssignmentSvc.GetEmployeeBenefit(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, eb)
}

func (h *Handler) UpdateEmployeeBenefit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		CoverageLevel *string `json:"coverage_level"`
		Notes         *string `json:"notes"`
	}
	if !bindJSON(c, &req) {
		return
	}
	eb, err := h.AssignmentSvc.GetEmployeeBenefit(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if req.CoverageLevel != nil {
		eb.CoverageLevel = req.CoverageLevel
	}
	if req.Notes != nil {
		eb.Notes = req.Notes
	}
	eb, err = h.AssignmentSvc.UpdateEmployeeBenefit(c.Request.Context(), companyID(c), eb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, eb)
}

func (h *Handler) CancelEmployeeBenefit(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.AssignmentSvc.CancelEmployeeBenefit(c.Request.Context(), companyID(c), id, req.Reason, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "benefit cancelled"})
}

func (h *Handler) GetHistory(c *gin.Context) {
	ebID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	history, err := h.AssignmentSvc.GetHistory(c.Request.Context(), ebID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, history)
}

func (h *Handler) ListHistoryByEmployee(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	history, err := h.AssignmentSvc.ListHistoryByEmployee(c.Request.Context(), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, history)
}

func (h *Handler) CreateRequest(c *gin.Context) {
	var req struct {
		EmployeeID  uuid.UUID              `json:"employee_id" binding:"required"`
		BenefitID   uuid.UUID              `json:"benefit_id" binding:"required"`
		RequestType string                 `json:"request_type" binding:"required"`
		Data        map[string]any         `json:"data"`
	}
	if !bindJSON(c, &req) {
		return
	}
	r, err := h.AssignmentSvc.CreateRequest(c.Request.Context(), companyID(c), req.EmployeeID, req.BenefitID, req.RequestType, req.Data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, r)
}

func (h *Handler) ListRequests(c *gin.Context) {
	var employeeID *uuid.UUID
	if eid := c.Query("employee_id"); eid != "" {
		id, err := uuid.Parse(eid)
		if err == nil {
			employeeID = &id
		}
	}
	status := qs(c, "status")
	requests, err := h.AssignmentSvc.ListRequests(c.Request.Context(), companyID(c), employeeID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, requests)
}

func (h *Handler) GetRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	r, err := h.AssignmentSvc.GetRequest(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, r)
}

func (h *Handler) SubmitRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AssignmentSvc.SubmitRequest(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "request submitted"})
}

func (h *Handler) ApproveRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Comment string `json:"comment"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if err := h.AssignmentSvc.ApproveRequest(c.Request.Context(), id, userID(c), req.Comment); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "request approved"})
}

func (h *Handler) RejectRequest(c *gin.Context) {
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
	if err := h.AssignmentSvc.RejectRequest(c.Request.Context(), id, userID(c), req.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "request rejected"})
}

func (h *Handler) CancelRequest(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.AssignmentSvc.CancelRequest(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "request cancelled"})
}

func (h *Handler) ListReviews(c *gin.Context) {
	requestID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request id"})
		return
	}
	reviews, err := h.AssignmentSvc.ListReviews(c.Request.Context(), requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, reviews)
}
