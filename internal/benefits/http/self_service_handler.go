package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) MyBenefits(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	cid := companyID(c)
	benefits, err := h.AssignmentSvc.ListEmployeeBenefits(c.Request.Context(), &cid, &empID, nil, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, benefits)
}

func (h *Handler) MyWallets(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	wallets, err := h.WalletSvc.ListEmployeeWallets(c.Request.Context(), companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, wallets)
}

func (h *Handler) MyReimbursements(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	reimbursements, err := h.ReimbursementSvc.ListReimbursements(c.Request.Context(), companyID(c), &empID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, reimbursements)
}

func (h *Handler) MyRequests(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requests, err := h.AssignmentSvc.ListRequests(c.Request.Context(), companyID(c), &empID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, requests)
}

func (h *Handler) MyBonuses(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	bonuses, err := h.BonusSvc.ListBonuses(c.Request.Context(), companyID(c), &empID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, bonuses)
}

func (h *Handler) MyIncentives(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	incentives, err := h.BonusSvc.ListIncentives(c.Request.Context(), companyID(c), &empID, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	success(c, incentives)
}

func (h *Handler) MyRewards(c *gin.Context) {
	empID := employeeID(c)
	if empID == uuid.Nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	snapshot, err := h.RewardsSvc.GetLatestSnapshot(c.Request.Context(), companyID(c), empID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	success(c, snapshot)
}
