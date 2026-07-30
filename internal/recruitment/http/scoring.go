package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/domain"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateScoringModel(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.ScoringModel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ScoringSvc.CreateModel(c.Request.Context(), companyID, &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListScoringModels(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ScoringSvc.ListModels(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetScoringModel(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ScoringSvc.GetModel(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Scoring model not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateScoringModel(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.ScoringModel
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ScoringSvc.UpdateModel(c.Request.Context(), companyID, c.Param("id"), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteScoringModel(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.ScoringSvc.DeleteModel(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "scoring model deleted"})
}

func (h *Handler) ListScoringCriteria(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ScoringSvc.ListCriteria(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddScoringCriterion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.ScoringCriterion
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ScoringSvc.AddCriterion(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateScoringCriterion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req domain.ScoringCriterion
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.ScoringSvc.UpdateCriterion(c.Request.Context(), companyID, c.Param("criterionId"), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "criterion updated"})
}

func (h *Handler) DeleteScoringCriterion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.ScoringSvc.DeleteCriterion(c.Request.Context(), companyID, c.Query("model_id"), c.Param("criterionId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "criterion deleted"})
}

func (h *Handler) ScoreCandidate(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req struct {
		CandidateID string `json:"candidate_id" binding:"required"`
		PositionID  string `json:"position_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.ScoringSvc.ScoreCandidate(c.Request.Context(), companyID, req.CandidateID, req.PositionID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetMatchingResult(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.ScoringSvc.ScoreCandidate(c.Request.Context(), companyID, c.Param("candidateId"), c.Param("positionId"))
	if err != nil {
		response.NotFound(c, "Matching result not found")
		return
	}
	response.Success(c, data)
}
