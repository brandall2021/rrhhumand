package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rrhhumand/api/internal/benefits/application"
)

type Handler struct {
	CatalogSvc       *application.CatalogService
	BenefitSvc       *application.BenefitService
	EligibilitySvc   *application.EligibilityService
	WorkflowSvc      *application.WorkflowService
	AssignmentSvc    *application.AssignmentService
	WalletSvc        *application.WalletService
	ReimbursementSvc *application.ReimbursementService
	BonusSvc         *application.BonusService
	RewardsSvc       *application.RewardsService
	CostSvc          *application.CostService
}

func NewHandler(
	cs *application.CatalogService,
	bs *application.BenefitService,
	es *application.EligibilityService,
	ws *application.WorkflowService,
	as *application.AssignmentService,
	wls *application.WalletService,
	rs *application.ReimbursementService,
	bns *application.BonusService,
	rws *application.RewardsService,
	cst *application.CostService,
) *Handler {
	return &Handler{
		CatalogSvc:       cs,
		BenefitSvc:       bs,
		EligibilitySvc:   es,
		WorkflowSvc:      ws,
		AssignmentSvc:    as,
		WalletSvc:        wls,
		ReimbursementSvc: rs,
		BonusSvc:         bns,
		RewardsSvc:       rws,
		CostSvc:          cst,
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
