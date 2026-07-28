package documents

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rrhhumand/api/internal/models"
	"github.com/rrhhumand/api/internal/tenant"
	"github.com/rrhhumand/api/pkg/response"
)

type DocumentHandler struct {
	service *DocumentService
}

func NewDocumentHandler(service *DocumentService) *DocumentHandler {
	return &DocumentHandler{service: service}
}

func (h *DocumentHandler) UploadDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required")
		return
	}
	defer file.Close()

	title := c.PostForm("title")
	if title == "" {
		title = header.Filename
	}

	description := c.PostForm("description")
	categoryID := c.PostForm("category_id")
	employeeID := c.PostForm("employee_id")
	departmentID := c.PostForm("department_id")
	isPublic := c.PostForm("is_public") == "true"

	req := &CreateDocumentRequest{
		Title:        title,
		Description:  &description,
		IsPublic:     &isPublic,
	}

	if categoryID != "" {
		req.CategoryID = &categoryID
	}
	if employeeID != "" {
		req.EmployeeID = &employeeID
	}
	if departmentID != "" {
		req.DepartmentID = &departmentID
	}

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	doc, err := h.service.UploadDocument(c.Request.Context(), companyID, userID, req, file, header.Filename, header.Size, contentType)
	if err != nil {
		switch err.Error() {
		case "file type not allowed":
			response.BadRequest(c, "File type not allowed")
		case "file too large":
			response.BadRequest(c, "File exceeds maximum size")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Created(c, doc)
}

func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	params := models.NewPaginationParams(c)

	filters := DocumentFilters{
		Status:       c.Query("status"),
		CategoryID:   c.Query("category_id"),
		EmployeeID:   c.Query("employee_id"),
		DepartmentID: c.Query("department_id"),
		MimeType:     c.Query("mime_type"),
		Search:       c.Query("q"),
		CreatedFrom:  c.Query("created_from"),
		CreatedTo:    c.Query("created_to"),
		Tag:          c.Query("tag"),
	}

	docs, total, err := h.service.ListDocuments(c.Request.Context(), companyID, filters, params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    docs,
		"meta":    params.ToMeta(total),
	})
}

func (h *DocumentHandler) GetDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	doc, err := h.service.GetDocumentByID(c.Request.Context(), id, companyID)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, doc)
}

func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	doc, err := h.service.UpdateDocument(c.Request.Context(), id, companyID, &req)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, doc)
}

func (h *DocumentHandler) DeleteDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.DeleteDocument(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "document not found":
			response.NotFound(c, "Document not found")
		case "document already deleted":
			response.BadRequest(c, "Document already in trash")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.NoContent(c)
}

func (h *DocumentHandler) RestoreDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.RestoreDocument(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "document not found":
			response.NotFound(c, "Document not found")
		case "document not in trash":
			response.BadRequest(c, "Document is not in trash")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "document restored"})
}

func (h *DocumentHandler) PermanentDelete(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.PermanentDelete(c.Request.Context(), id, companyID); err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.NoContent(c)
}

func (h *DocumentHandler) ArchiveDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	id := c.Param("id")

	if err := h.service.ArchiveDocument(c.Request.Context(), id, companyID); err != nil {
		switch err.Error() {
		case "document not found":
			response.NotFound(c, "Document not found")
		case "document not active":
			response.BadRequest(c, "Document is not active")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "document archived"})
}

func (h *DocumentHandler) DownloadDocument(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	id := c.Param("id")

	reader, doc, err := h.service.DownloadDocument(c.Request.Context(), id, companyID, userID)
	if err != nil {
		switch err.Error() {
		case "document not found":
			response.NotFound(c, "Document not found")
		case "access denied":
			response.Forbidden(c, "Access denied")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}
	defer reader.Close()

	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", doc.OriginalFilename))
	c.Header("Content-Type", doc.MimeType)
	c.Header("Content-Length", strconv.FormatInt(doc.FileSize, 10))
	c.Status(http.StatusOK)
	io.Copy(c.Writer, reader)
}

func (h *DocumentHandler) CreateVersion(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	documentID := c.Param("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "File is required")
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	version, err := h.service.CreateVersion(c.Request.Context(), documentID, companyID, userID, file, header.Filename, header.Size, contentType)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, version)
}

func (h *DocumentHandler) ListVersions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	documentID := c.Param("id")

	versions, err := h.service.ListVersions(c.Request.Context(), documentID, companyID)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, versions)
}

func (h *DocumentHandler) SetPermissions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	documentID := c.Param("id")

	var req SetDocumentPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.SetPermissions(c.Request.Context(), documentID, companyID, &req); err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "permissions updated"})
}

func (h *DocumentHandler) ListPermissions(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	documentID := c.Param("id")

	perms, err := h.service.ListPermissions(c.Request.Context(), documentID, companyID)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, perms)
}

func (h *DocumentHandler) CreateTag(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	var req CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	tag, err := h.service.CreateTag(c.Request.Context(), companyID, req.Name)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, tag)
}

func (h *DocumentHandler) ListTags(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	tags, err := h.service.ListTags(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, tags)
}

func (h *DocumentHandler) CreateShare(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	documentID := c.Param("id")

	var req ShareDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	share, err := h.service.CreateShare(c.Request.Context(), documentID, companyID, userID, &req)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, share)
}

func (h *DocumentHandler) CreateShareLink(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	userID := tenant.GetUserID(c)
	documentID := c.Param("id")

	var req CreateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	share, err := h.service.CreateShareLink(c.Request.Context(), documentID, companyID, userID, &req)
	if err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Created(c, share)
}

func (h *DocumentHandler) AccessShareLink(c *gin.Context) {
	token := c.Param("token")
	userID := tenant.GetUserID(c)

	doc, err := h.service.AccessShareLink(c.Request.Context(), token, userID)
	if err != nil {
		switch err.Error() {
		case "invalid share link":
			response.NotFound(c, "Invalid share link")
		case "share link expired":
			response.BadRequest(c, "Share link has expired")
		case "share link max uses reached":
			response.BadRequest(c, "Share link has reached maximum uses")
		default:
			response.InternalError(c, err.Error())
		}
		return
	}

	response.Success(c, doc)
}

func (h *DocumentHandler) RevokeShare(c *gin.Context) {
	shareID := c.Param("id")

	if err := h.service.RevokeShare(c.Request.Context(), shareID); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "share revoked"})
}

func (h *DocumentHandler) ListExpiringDocuments(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	days := 30
	if d := c.Query("days"); d != "" {
		if parsed, err := strconv.Atoi(d); err == nil && parsed > 0 {
			days = parsed
		}
	}

	docs, err := h.service.ListExpiringDocuments(c.Request.Context(), companyID, days)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, docs)
}

func (h *DocumentHandler) GetDocumentStats(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)

	stats, err := h.service.GetDocumentStats(c.Request.Context(), companyID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, stats)
}

func (h *DocumentHandler) ListEmployeeDocuments(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	employeeID := c.Param("id")

	filters := DocumentFilters{
		EmployeeID: employeeID,
	}

	params := models.NewPaginationParams(c)
	docs, total, err := h.service.ListDocuments(c.Request.Context(), companyID, filters, params)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    docs,
		"meta":    params.ToMeta(total),
	})
}

func (h *DocumentHandler) SetDocumentTags(c *gin.Context) {
	companyID := tenant.GetCompanyID(c)
	documentID := c.Param("id")

	var req struct {
		Tags []string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Invalid request body")
		return
	}

	if err := h.service.SetDocumentTags(c.Request.Context(), documentID, companyID, req.Tags); err != nil {
		if err.Error() == "document not found" {
			response.NotFound(c, "Document not found")
			return
		}
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{"message": "tags updated"})
}
