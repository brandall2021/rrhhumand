package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/payroll/features/application"
)

type Handler struct {
	ReceiptSvc    *application.ReceiptService
	ArcaSvc       *application.ArcaService
	BookSvc       *application.BookService
	BankSvc       *application.BankService
	AccountingSvc *application.AccountingService
	ReportSvc     *application.ReportService
}

func NewHandler(rs *application.ReceiptService, as *application.ArcaService, bs *application.BookService, bks *application.BankService, acs *application.AccountingService, rps *application.ReportService) *Handler {
	return &Handler{
		ReceiptSvc:    rs,
		ArcaSvc:       as,
		BookSvc:       bs,
		BankSvc:       bks,
		AccountingSvc: acs,
		ReportSvc:     rps,
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

func uuidPtr(s *string) *uuid.UUID {
	if s == nil {
		return nil
	}
	id, err := uuid.Parse(*s)
	if err != nil {
		return nil
	}
	return &id
}

func success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data})
}

func created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, gin.H{"data": data})
}
