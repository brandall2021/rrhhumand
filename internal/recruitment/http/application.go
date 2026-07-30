package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		CandidateID string `json:"candidate_id"`
		PostingID   string `json:"posting_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ApplicationSvc.Create(c.Request.Context(), companyID, req.CandidateID, req.PostingID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListApplications(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	q := c.Request.URL.Query()
	data, err := h.ApplicationSvc.List(c.Request.Context(), companyID, q.Get("candidate_id"), q.Get("posting_id"), q.Get("status"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Application not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) MoveStage(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req struct {
		ToStageID string `json:"to_stage_id"`
		Reason    string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.ApplicationSvc.MoveStage(c.Request.Context(), companyID, c.Param("id"), req.ToStageID, userID, req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "stage moved"})
}

func (h *Handler) RejectApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		ReasonID   string `json:"reason_id"`
		ReasonText string `json:"reason_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.ApplicationSvc.Reject(c.Request.Context(), companyID, c.Param("id"), req.ReasonID, req.ReasonText); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "application rejected"})
}

func (h *Handler) WithdrawApplication(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.ApplicationSvc.Withdraw(c.Request.Context(), companyID, c.Param("id"), req.Reason); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "application withdrawn"})
}

func (h *Handler) GetStageHistory(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.GetStageHistory(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ListApplicationNotes(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.ListNotes(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddApplicationNote(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req domain.ApplicationNote
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	req.AuthorID = userID
	data, err := h.ApplicationSvc.AddNote(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateApplicationNote(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.ApplicationNote
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	req.ID = c.Param("noteId")
	if err := h.ApplicationSvc.UpdateNote(c.Request.Context(), companyID, c.Param("id"), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "note updated"})
}

func (h *Handler) ListApplicationRatings(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ApplicationSvc.ListRatings(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddApplicationRating(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req domain.ApplicationRating
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	req.RatedBy = userID
	data, err := h.ApplicationSvc.AddRating(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}
