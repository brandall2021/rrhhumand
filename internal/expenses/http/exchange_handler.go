package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"

	"github.com/rrhhumand/api/internal/expenses/domain"
)

func (h *Handler) CreateRate(c *gin.Context) {
	var req struct {
		FromCurrency  string          `json:"from_currency" binding:"required"`
		ToCurrency    string          `json:"to_currency" binding:"required"`
		Rate          decimal.Decimal `json:"rate" binding:"required"`
		EffectiveDate *time.Time      `json:"effective_date"`
		Source        *string         `json:"source"`
	}
	if !bindJSON(c, &req) {
		return
	}
	rate := &domain.ExchangeRate{
		FromCurrency: req.FromCurrency,
		ToCurrency:   req.ToCurrency,
		Rate:         req.Rate,
		Source:       "manual",
	}
	if req.Source != nil {
		rate.Source = *req.Source
	}
	r, err := h.ExchangeSvc.CreateRate(c.Request.Context(), companyID(c), rate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, r)
}

func (h *Handler) GetLatestRate(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	if from == "" || to == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "from and to query params required"})
		return
	}
	rate, err := h.ExchangeSvc.GetLatestRate(c.Request.Context(), from, to)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, rate)
}

func (h *Handler) Convert(c *gin.Context) {
	var req struct {
		FromCurrency string          `json:"from_currency" binding:"required"`
		ToCurrency   string          `json:"to_currency" binding:"required"`
		Amount       decimal.Decimal `json:"amount" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	result, err := h.ExchangeSvc.Convert(c.Request.Context(), req.Amount, req.FromCurrency, req.ToCurrency, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"converted_amount": result})
}
