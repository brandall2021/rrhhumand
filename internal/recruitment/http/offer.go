package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

func (h *Handler) CreateOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.OfferSvc.Create(c.Request.Context(), companyID, userID, req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) ListOffers(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.OfferSvc.List(c.Request.Context(), companyID, c.Request.URL.Query())
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) GetOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.OfferSvc.GetByID(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.NotFound(c, "Offer not found")
		return
	}
	response.Success(c, data)
}

func (h *Handler) UpdateOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.OfferSvc.Update(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) SubmitOfferForApproval(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.OfferSvc.SubmitForApproval(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer submitted for approval"})
}

func (h *Handler) ApproveOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.OfferSvc.Approve(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer approved"})
}

func (h *Handler) SendOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.OfferSvc.Send(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer sent"})
}

func (h *Handler) AcceptOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.OfferSvc.Accept(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer accepted"})
}

func (h *Handler) RejectOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	if err := h.OfferSvc.Reject(c.Request.Context(), companyID, c.Param("id"), req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer rejected"})
}

func (h *Handler) WithdrawOffer(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	if err := h.OfferSvc.Withdraw(c.Request.Context(), companyID, c.Param("id")); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{"message": "offer withdrawn"})
}

func (h *Handler) ListOfferNegotiations(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.OfferSvc.ListNegotiations(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddOfferNegotiation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.OfferSvc.AddNegotiation(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}

func (h *Handler) UpdateOfferNegotiation(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.OfferSvc.UpdateNegotiation(c.Request.Context(), companyID, c.Param("negId"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) ListOfferDocuments(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	data, err := h.OfferSvc.ListDocuments(c.Request.Context(), companyID, c.Param("id"))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, data)
}

func (h *Handler) AddOfferDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	var req map[string]any
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}
	data, err := h.OfferSvc.AddDocument(c.Request.Context(), companyID, c.Param("id"), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, data)
}
