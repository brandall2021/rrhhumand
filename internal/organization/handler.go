package organization

import (
	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type OrgHandler struct {
	repo *OrgRepository
}

func NewOrgHandler(repo *OrgRepository) *OrgHandler {
	return &OrgHandler{repo: repo}
}

func (h *OrgHandler) GetTree(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	tree, err := h.repo.GetTree(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, "Failed to build organization tree")
		return
	}

	if tree == nil {
		tree = []*OrgNode{}
	}

	response.Success(c, tree)
}
