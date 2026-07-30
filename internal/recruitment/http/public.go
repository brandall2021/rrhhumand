package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/recruitment/application"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) ListPublicPostings(c *gin.Context) {
	data, err := h.PostingSvc.ListPublic(c.Request.Context())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetPublicPosting(c *gin.Context) {
	data, err := h.PostingSvc.GetPublicByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.NotFound(c, "Job posting not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) PublicApply(c *gin.Context) {
	var req application.PublicApplyReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	posting, err := h.PostingSvc.GetPublicByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		response.NotFound(c, "Job posting not found")
		return
	}
	candidate, err := h.CandidateSvc.Create(c.Request.Context(), posting.CompanyID, &application.CreateCandidateReq{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	data, err := h.ApplicationSvc.Create(c.Request.Context(), posting.CompanyID, candidate.ID, posting.ID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}
