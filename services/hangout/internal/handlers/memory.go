package handlers

import (
	"fmt"
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MemoryHandler interface {
	CreateMemories(c echo.Context) error
	GetMemory(c echo.Context) error
	ListMemories(c echo.Context) error
	DeleteMemory(c echo.Context) error
}

type memoryHandler struct {
	memoryService   services.MemoryService
	responseBuilder *response.Builder
}

func NewMemoryHandler(memoryService services.MemoryService, responseBuilder *response.Builder) MemoryHandler {
	return &memoryHandler{
		memoryService:   memoryService,
		responseBuilder: responseBuilder,
	}
}

// @Summary      Upload Memories
// @Description  Uploads multiple photos/memories for a hangout with concurrent processing.
// @Description
// @Description  **Constraints:**
// @Description  - Maximum 10 files per request
// @Description  - Maximum 10MB per file
// @Description  - Allowed formats: .jpg, .jpeg, .png, .gif, .webp
// @Description  - MIME types: image/jpeg, image/png, image/gif, image/webp
// @Description
// @Description  **Upload Behavior:**
// @Description  - Files are processed concurrently for faster upload
// @Description  - Returns partial success if some files succeed and others fail
// @Description  - Each file gets its own transaction (atomic per file)
// @Description  - If all files fail, returns error
// @Description
// @Description  **Form Field:**
// @Description  - Use field name "files" in multipart/form-data
// @Description  - Can attach multiple files to the same field
// @Tags         Memories
// @Accept       multipart/form-data
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Param        files formData file true "Files to upload (use same field name for multiple files)"
// @Success      201 {object} response.StandardResponse{data=[]dto.MemoryResponse} "Memories uploaded successfully (partial success possible)"
// @Failure      400 {object} response.StandardResponse "Too many files, file too large, invalid format, or no files provided"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "Hangout not found"
// @Failure      500 {object} response.StandardResponse "Internal server error or all files failed"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id}/memories [post]
func (h *memoryHandler) CreateMemories(c echo.Context) error {
	hangoutIDParam := c.Param("hangout_id")
	hangoutID, err := uuid.Parse(hangoutIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidHangoutID))
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	memories, err := h.memoryService.CreateMemories(ctx, userID, hangoutID, files)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case apperrors.ErrTooManyFiles, apperrors.ErrFileTooLarge, apperrors.ErrInvalidFileType:
			statusCode = http.StatusBadRequest
		case apperrors.ErrInvalidHangoutID:
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusCreated, h.responseBuilder.Success(constants.MemoriesUploadedSuccessfully, memories))
}

// @Summary      Get Memory by ID
// @Description  Retrieves a memory by its ID with presigned file URL
// @Tags         Memories
// @Produce      json
// @Param        memory_id path string true "Memory ID"
// @Success      200 {object} response.StandardResponse{data=dto.MemoryResponse} "Memory fetched successfully"
// @Failure      400 {object} response.StandardResponse "Invalid memory ID"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "Memory not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /memories/{memory_id} [get]
func (h *memoryHandler) GetMemory(c echo.Context) error {
	memoryIDParam := c.Param("memory_id")
	memoryID, err := uuid.Parse(memoryIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidMemoryID))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	memory, err := h.memoryService.GetMemory(ctx, userID, memoryID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == apperrors.ErrMemoryNotFound {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.MemoryFetchedSuccessfully, memory))
}

// @Summary      List Memories
// @Description  Lists all memories for a hangout with cursor pagination
// @Tags         Memories
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Param        limit query int false "Number of items to return (default 10, max 100)"
// @Param        after_id query string false "Cursor for pagination (Memory ID)"
// @Param        sort_dir query string false "Sort direction: asc or desc (default desc)"
// @Success      200 {object} response.StandardResponse{data=dto.PaginatedMemories} "Memories listed successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "Hangout not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id}/memories [get]
func (h *memoryHandler) ListMemories(c echo.Context) error {
	hangoutIDParam := c.Param("hangout_id")
	hangoutID, err := uuid.Parse(hangoutIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidHangoutID))
	}

	pagination := &dto.CursorPagination{}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		var limit int
		if _, err := fmt.Sscanf(limitStr, "%d", &limit); err == nil {
			pagination.Limit = limit
		}
	}

	if afterIDStr := c.QueryParam("after_id"); afterIDStr != "" {
		if afterID, err := uuid.Parse(afterIDStr); err == nil {
			pagination.AfterID = &afterID
		}
	}

	if sortDir := c.QueryParam("sort_dir"); sortDir != "" {
		pagination.SortDir = sortDir
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	memories, err := h.memoryService.ListMemories(ctx, userID, hangoutID, pagination)
	if err != nil {
		statusCode := http.StatusInternalServerError
		switch err {
		case apperrors.ErrInvalidHangoutID:
			statusCode = http.StatusNotFound
		case apperrors.ErrInvalidCursorPagination:
			statusCode = http.StatusBadRequest
		}
		return c.JSON(statusCode, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.MemoriesListedSuccessfully, memories))
}

// @Summary      Delete Memory
// @Description  Deletes a memory and its associated file
// @Tags         Memories
// @Produce      json
// @Param        memory_id path string true "Memory ID"
// @Success      200 {object} response.StandardResponse "Memory deleted successfully"
// @Failure      400 {object} response.StandardResponse "Invalid memory ID"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "Memory not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /memories/{memory_id} [delete]
func (h *memoryHandler) DeleteMemory(c echo.Context) error {
	memoryIDParam := c.Param("memory_id")
	memoryID, err := uuid.Parse(memoryIDParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidMemoryID))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	err = h.memoryService.DeleteMemory(ctx, userID, memoryID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err == apperrors.ErrMemoryNotFound {
			statusCode = http.StatusNotFound
		}
		return c.JSON(statusCode, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.MemoryDeletedSuccessfully, nil))
}
