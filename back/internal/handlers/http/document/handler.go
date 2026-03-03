package document

import (
	"errors"
	"keden-service/back/internal/middleware"
	"keden-service/back/internal/services/document"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DocumentHandler struct {
	docService *document.DocumentService
}

func NewDocumentHandler(ds *document.DocumentService) *DocumentHandler {
	return &DocumentHandler{docService: ds}
}

func (h *DocumentHandler) Upload(c *gin.Context) {
	userID := middleware.GetUserID(c)

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	defer file.Close()

	doc, err := h.docService.UploadDocument(c.Request.Context(), userID, file, header)
	if err != nil {
		if errors.Is(err, document.ErrNotPDF) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "only PDF files are allowed"})
			return
		}
		if errors.Is(err, document.ErrFileTooLarge) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file size exceeds 50MB limit"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload document"})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

func (h *DocumentHandler) GetDocuments(c *gin.Context) {
	userID := middleware.GetUserID(c)

	docs, err := h.docService.GetDocuments(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
		return
	}

	c.JSON(http.StatusOK, docs)
}

func (h *DocumentHandler) GetDocumentByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	doc, err := h.docService.GetDocumentByID(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, document.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get document"})
		return
	}

	c.JSON(http.StatusOK, doc)
}

func (h *DocumentHandler) DownloadExcel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document id"})
		return
	}

	data, filename, err := h.docService.DownloadExcel(c.Request.Context(), uint(id), userID)
	if err != nil {
		if errors.Is(err, document.ErrDocumentNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
			return
		}
		if errors.Is(err, document.ErrExcelNotReady) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "excel file is not ready yet"})
			return
		}
		if errors.Is(err, document.ErrAccessDenied) {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to download excel"})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *DocumentHandler) GetAllDocuments(c *gin.Context) {
	docs, err := h.docService.GetAllDocuments(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get documents"})
		return
	}

	c.JSON(http.StatusOK, docs)
}
