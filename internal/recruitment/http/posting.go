package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreatePosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PostingSvc.Create(c.Request.Context(), companyID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListPostings(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.PostingSvc.List(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetPosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.PostingSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Posting not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdatePosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PostingSvc.Update(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) PublishPosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.PostingSvc.Publish(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "posting published"})
}

func (h *Handler) ClosePosting(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.PostingSvc.Close(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "posting closed"})
}

func (h *Handler) ListScreeningQuestions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.PostingSvc.ListScreeningQuestions(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddScreeningQuestion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PostingSvc.AddScreeningQuestion(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateScreeningQuestion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.PostingSvc.UpdateScreeningQuestion(c.Request.Context(), companyID, c.Param("questionId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) DeleteScreeningQuestion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.PostingSvc.DeleteScreeningQuestion(c.Request.Context(), companyID, c.Param("questionId")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "question deleted"})
}
