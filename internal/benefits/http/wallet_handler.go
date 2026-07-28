package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/domain"
	"github.com/shopspring/decimal"
)

func (h *Handler) CreateWallet(c *gin.Context) {
	var req struct {
		EmployeeID uuid.UUID       `json:"employee_id" binding:"required"`
		WalletType string          `json:"wallet_type" binding:"required"`
		Balance    decimal.Decimal `json:"balance"`
		Currency   string          `json:"currency"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Currency == "" {
		req.Currency = "ARS"
	}
	w, err := h.WalletSvc.CreateWallet(c.Request.Context(), companyID(c), req.EmployeeID, req.WalletType, req.Balance, req.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, w)
}

func (h *Handler) GetWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	w, err := h.WalletSvc.GetWallet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, w)
}

func (h *Handler) ListEmployeeWallets(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	wallets, err := h.WalletSvc.ListEmployeeWallets(c.Request.Context(), companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, wallets)
}

func (h *Handler) CreditWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Amount      decimal.Decimal `json:"amount" binding:"required"`
		TxType      string          `json:"transaction_type" binding:"required"`
		Description string          `json:"description"`
	}
	if !bindJSON(c, &req) {
		return
	}
	uid := userID(c)
	if err := h.WalletSvc.CreditWallet(c.Request.Context(), id, req.Amount, req.TxType, req.Description, &uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "wallet credited"})
}

func (h *Handler) DebitWallet(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Amount      decimal.Decimal `json:"amount" binding:"required"`
		TxType      string          `json:"transaction_type" binding:"required"`
		Description string          `json:"description"`
	}
	if !bindJSON(c, &req) {
		return
	}
	uid := userID(c)
	if err := h.WalletSvc.DebitWallet(c.Request.Context(), id, req.Amount, req.TxType, req.Description, &uid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, gin.H{"message": "wallet debited"})
}

func (h *Handler) ListTransactions(c *gin.Context) {
	walletID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet id"})
		return
	}
	limit := qi(c, "limit", 50)
	offset := qi(c, "offset", 0)
	transactions, err := h.WalletSvc.ListTransactions(c.Request.Context(), walletID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, transactions)
}

func (h *Handler) CreateFlexiblePlan(c *gin.Context) {
	var req struct {
		Name                string          `json:"name" binding:"required"`
		Description         *string         `json:"description"`
		PlanType            string          `json:"plan_type" binding:"required"`
		AnnualAmount        decimal.Decimal `json:"annual_amount"`
		MonthlyAmount       decimal.Decimal `json:"monthly_amount"`
		Currency            string          `json:"currency"`
		EmployerContribution decimal.Decimal `json:"employer_contribution"`
		EmployeeContribution decimal.Decimal `json:"employee_contribution"`
		MaxRolloverAmount   decimal.Decimal `json:"max_rollover_amount"`
		AllowReimbursement  bool            `json:"allow_reimbursement"`
		TaxExempt           bool            `json:"tax_exempt"`
		EligibleCategories  []string        `json:"eligible_categories"`
	}
	if !bindJSON(c, &req) {
		return
	}
	if req.Currency == "" {
		req.Currency = "ARS"
	}
	p := &domain.BenefitFlexiblePlan{
		Name:                req.Name,
		Description:         req.Description,
		PlanType:            req.PlanType,
		AnnualAmount:        req.AnnualAmount,
		MonthlyAmount:       req.MonthlyAmount,
		Currency:            req.Currency,
		EmployerContribution: req.EmployerContribution,
		EmployeeContribution: req.EmployeeContribution,
		MaxRolloverAmount:   req.MaxRolloverAmount,
		AllowReimbursement:  req.AllowReimbursement,
		TaxExempt:           req.TaxExempt,
		EligibleCategories:  req.EligibleCategories,
	}
	plan, err := h.WalletSvc.CreateFlexiblePlan(c.Request.Context(), companyID(c), userID(c), p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, plan)
}

func (h *Handler) ListFlexiblePlans(c *gin.Context) {
	plans, err := h.WalletSvc.ListFlexiblePlans(c.Request.Context(), companyID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, plans)
}

func (h *Handler) CreateBudget(c *gin.Context) {
	var req struct {
		EmployeeID  uuid.UUID       `json:"employee_id" binding:"required"`
		PlanID      uuid.UUID       `json:"plan_id" binding:"required"`
		FiscalYear  int             `json:"fiscal_year" binding:"required"`
		TotalAmount decimal.Decimal `json:"total_amount" binding:"required"`
	}
	if !bindJSON(c, &req) {
		return
	}
	b, err := h.WalletSvc.CreateBudget(c.Request.Context(), companyID(c), req.EmployeeID, req.PlanID, req.FiscalYear, req.TotalAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	created(c, b)
}

func (h *Handler) GetBudget(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	planID, err := uuid.Parse(c.Param("plan_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan id"})
		return
	}
	fiscalYear := qi(c, "fiscal_year", 0)
	if fiscalYear == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fiscal_year required"})
		return
	}
	b, err := h.WalletSvc.GetBudget(c.Request.Context(), companyID(c), empID, planID, fiscalYear)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, b)
}

func (h *Handler) ListEmployeeBudgets(c *gin.Context) {
	empID, err := uuid.Parse(c.Param("employee_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee id"})
		return
	}
	budgets, err := h.WalletSvc.ListEmployeeBudgets(c.Request.Context(), companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, budgets)
}
