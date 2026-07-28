package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateCategory(c *gin.Context) {
	var req domain.ExpenseCategory
	if !bindJSON(c, &req) {
		return
	}
	cat, err := h.CatalogSvc.CreateCategory(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, cat)
}

func (h *Handler) ListCategories(c *gin.Context) {
	parentID := qs(c, "parent_id")
	var pid *uuid.UUID
	if parentID != nil {
		id, err := uuid.Parse(*parentID)
		if err == nil {
			pid = &id
		}
	}
	categories, err := h.CatalogSvc.ListCategories(c.Request.Context(), companyID(c), pid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, categories)
}

func (h *Handler) GetCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	cat, err := h.CatalogSvc.GetCategory(c.Request.Context(), companyID(c), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, cat)
}

func (h *Handler) UpdateCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req domain.ExpenseCategory
	if !bindJSON(c, &req) {
		return
	}
	req.ID = id
	cat, err := h.CatalogSvc.UpdateCategory(c.Request.Context(), companyID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, cat)
}

func (h *Handler) DeleteCategory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.CatalogSvc.DeleteCategory(c.Request.Context(), companyID(c), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) CreatePaymentMethod(c *gin.Context) {
	var req domain.ExpensePaymentMethod
	if !bindJSON(c, &req) {
		return
	}
	pm, err := h.CatalogSvc.CreatePaymentMethod(c.Request.Context(), companyID(c), userID(c), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, pm)
}

func (h *Handler) ListPaymentMethods(c *gin.Context) {
	methods, err := h.CatalogSvc.ListPaymentMethods(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, methods)
}
