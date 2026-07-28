package feed

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type FeedHandler struct {
	service *FeedService
}

func NewFeedHandler(service *FeedService) *FeedHandler {
	return &FeedHandler{service: service}
}

func (h *FeedHandler) ListPosts(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)
	search := c.Query("search")

	posts, total, err := h.service.ListPosts(c.Request.Context(), companyID, params, search)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	meta := params.ToMeta(total)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    posts,
		"meta":    meta,
	})
}

func (h *FeedHandler) CreatePost(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	authorID := c.GetString("employee_id")

	var req CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	post, err := h.service.CreatePost(c.Request.Context(), companyID, authorID, &req)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, post)
}

func (h *FeedHandler) GetPost(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	post, err := h.service.GetPostByID(c.Request.Context(), id, companyID, "")
	if err != nil {
		if err.Error() == "post not found" {
			response.NotFound(c, "Post not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, post)
}

func (h *FeedHandler) UpdatePost(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	authorID := c.GetString("employee_id")
	id := c.Param("id")

	var req UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	post, err := h.service.UpdatePost(c.Request.Context(), id, companyID, authorID, &req)
	if err != nil {
		if err.Error() == "post not found" {
			response.NotFound(c, "Post not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, post)
}

func (h *FeedHandler) DeletePost(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	authorID := c.GetString("employee_id")
	id := c.Param("id")

	if err := h.service.DeletePost(c.Request.Context(), id, companyID, authorID); err != nil {
		if err.Error() == "post not found" {
			response.NotFound(c, "Post not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

func (h *FeedHandler) AddComment(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	authorID := c.GetString("employee_id")
	postID := c.Param("id")

	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	comment, err := h.service.AddComment(c.Request.Context(), postID, companyID, authorID, &req)
	if err != nil {
		if err.Error() == "post not found" {
			response.NotFound(c, "Post not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, comment)
}

func (h *FeedHandler) AddReaction(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.GetString("employee_id")
	postID := c.Param("id")

	var req ReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.AddReaction(c.Request.Context(), postID, companyID, employeeID, &req); err != nil {
		if err.Error() == "post not found" {
			response.NotFound(c, "Post not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "reaction added"})
}

func (h *FeedHandler) RemoveReaction(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.GetString("employee_id")
	postID := c.Param("id")
	reactionType := c.Param("type")

	if err := h.service.RemoveReaction(c.Request.Context(), postID, companyID, employeeID, reactionType); err != nil {
		if err.Error() == "post not found" {
			response.NotFound(c, "Post not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}
