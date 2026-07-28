package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/expenses/application"
)

type Handler struct {
	CatalogSvc       *application.CatalogService
	ExpenseSvc       *application.ExpenseService
	TravelSvc        *application.TravelService
	ReportSvc        *application.ReportService
	AdvanceSvc       *application.AdvanceService
	ReimbursementSvc *application.ReimbursementService
	PolicySvc        *application.PolicyService
	ApprovalSvc      *application.ApprovalService
	BudgetSvc        *application.BudgetService
	ExchangeSvc      *application.ExchangeService
	AllowanceSvc     *application.AllowanceService
}

func NewHandler(
	cs *application.CatalogService,
	es *application.ExpenseService,
	ts *application.TravelService,
	rs *application.ReportService,
	as *application.AdvanceService,
	rms *application.ReimbursementService,
	ps *application.PolicyService,
	aps *application.ApprovalService,
	bs *application.BudgetService,
	exs *application.ExchangeService,
	als *application.AllowanceService,
) *Handler {
	return &Handler{
		CatalogSvc:       cs,
		ExpenseSvc:       es,
		TravelSvc:        ts,
		ReportSvc:        rs,
		AdvanceSvc:       as,
		ReimbursementSvc: rms,
		PolicySvc:        ps,
		ApprovalSvc:      aps,
		BudgetSvc:        bs,
		ExchangeSvc:      exs,
		AllowanceSvc:     als,
	}
}

func companyID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString("company_id"))
	return id
}

func userID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString("user_id"))
	return id
}

func employeeID(c *gin.Context) uuid.UUID {
	id, _ := uuid.Parse(c.GetString("employee_id"))
	if id == uuid.Nil {
		id, _ = uuid.Parse(c.GetString("user_id"))
	}
	return id
}

func bindJSON(c *gin.Context, obj any) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}

func qs(c *gin.Context, key string) *string {
	v := c.Query(key)
	if v == "" {
		return nil
	}
	return &v
}

func qi(c *gin.Context, key string, def int) int {
	v := c.Query(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}
