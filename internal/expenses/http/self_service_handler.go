package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) MyExpenses(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	expenses, err := h.ExpenseSvc.ListExpenses(c.Request.Context(), companyID(c), &empID, nil, nil, nil, nil, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, expenses)
}

func (h *Handler) MyExpenseCreate(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		CategoryID       uuid.UUID  `json:"category_id" binding:"required"`
		ExpenseDate      string     `json:"expense_date" binding:"required"`
		Description      string     `json:"description" binding:"required"`
		OriginalAmount   float64    `json:"original_amount" binding:"required"`
		OriginalCurrency string     `json:"original_currency" binding:"required"`
		MerchantName     *string    `json:"merchant_name"`
		MerchantTaxID    *string    `json:"merchant_tax_id"`
		ReceiptNumber    *string    `json:"receipt_number"`
		PaymentMethodID  *uuid.UUID `json:"payment_method_id"`
		CostCenterID     *uuid.UUID `json:"cost_center_id"`
		ProjectID        *uuid.UUID `json:"project_id"`
	}
	if !bindJSON(c, &req) {
		return
	}
	expense, err := h.ExpenseSvc.UploadReceipt(c.Request.Context(), companyID(c), empID, userID(c), req.Description, req.OriginalCurrency, []byte{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_ = expense
	c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
}

func (h *Handler) MyExpenseSubmit(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ExpenseSvc.SubmitExpense(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	expense, _ := h.ExpenseSvc.GetExpense(c.Request.Context(), companyID(c), id)
	success(c, expense)
}

func (h *Handler) MyTravels(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	travels, err := h.TravelSvc.ListTravels(c.Request.Context(), companyID(c), &empID, nil, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, travels)
}

func (h *Handler) MyTravelCreate(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Purpose     *string `json:"purpose"`
		Origin      string  `json:"origin" binding:"required"`
		Destination string  `json:"destination" binding:"required"`
		DepartureDate string `json:"departure_date" binding:"required"`
		ReturnDate  string  `json:"return_date" binding:"required"`
		Currency    string  `json:"currency" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
}

func (h *Handler) MyTravelRequest(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
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

func (h *Handler) MyReports(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	reports, err := h.ReportSvc.ListReports(c.Request.Context(), companyID(c), &empID, nil, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, reports)
}

func (h *Handler) MyReportCreate(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
}

func (h *Handler) MyReportSubmit(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.ReportSvc.SubmitReport(c.Request.Context(), id, userID(c)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	report, _ := h.ReportSvc.GetReport(c.Request.Context(), id)
	success(c, report)
}

func (h *Handler) MyAdvances(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	advances, err := h.AdvanceSvc.ListAdvances(c.Request.Context(), companyID(c), &empID, nil, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, advances)
}

func (h *Handler) MyAdvanceCreate(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "not implemented"})
}

func (h *Handler) MyReimbursements(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	reimbursements, err := h.ReimbursementSvc.ListReimbursements(c.Request.Context(), companyID(c), &empID, nil, 0, 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, reimbursements)
}
