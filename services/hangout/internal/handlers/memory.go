package handlers

import (
	"fmt"
	"net/http"

	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/apperrors"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/constants"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/dto"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/request"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/http/response"
	"github.com/Ernestgio/Hangout-Planner/services/hangout/internal/services"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MemoryHandler interface {
	GenerateUploadURLs(c echo.Context) error
	ConfirmUpload(c echo.Context) error
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

// @Summary      Generate Upload URLs
// @Description  Creates memory records and returns presigned URLs for client-side upload to S3
// @Description  hangout_id is taken from the URL path, not the request body
// @Tags         Memories
// @Accept       json
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Param        request body dto.GenerateUploadURLsRequest true "Files to upload (hangout_id not needed in body)"
// @Success      201 {object} response.StandardResponse{data=dto.MemoryUploadResponse} "Upload URLs generated successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request payload"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      404 {object} response.StandardResponse "Hangout not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id}/memories/upload-urls [post]
func (h *memoryHandler) GenerateUploadURLs(c echo.Context) error {
	hangoutIDStr := c.Param("hangout_id")
	hangoutID, err := uuid.Parse(hangoutIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidHangoutID))
	}

	req, err := request.BindAndValidate[dto.GenerateUploadURLsRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	uploadResponse, err := h.memoryService.GenerateUploadURLs(ctx, userID, hangoutID, req)
	if err != nil {
		if err == apperrors.ErrInvalidHangoutID {
			return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusCreated, h.responseBuilder.Success(constants.UploadURLsGeneratedSuccessfully, uploadResponse))
}

// @Summary      Confirm Upload
// @Description  Confirms that files have been uploaded to S3 and marks them as ready
// @Tags         Memories
// @Accept       json
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Param        request body dto.ConfirmUploadRequest true "Memory IDs to confirm"
// @Success      200 {object} response.StandardResponse "Upload confirmed successfully"
// @Failure      400 {object} response.StandardResponse "Invalid request payload"
// @Failure      401 {object} response.StandardResponse "Unauthorized"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id}/memories/confirm-upload [post]
func (h *memoryHandler) ConfirmUpload(c echo.Context) error {
	req, err := request.BindAndValidate[dto.ConfirmUploadRequest](c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidPayload))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	err = h.memoryService.ConfirmUpload(ctx, userID, req)
	if err != nil {
		if err == apperrors.ErrMemoryNotFound || err == apperrors.ErrInvalidMemoryID {
			return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.UploadConfirmedSuccessfully, nil))
}

// @Summary      Get Memory
// @Description  Retrieves a single memory by ID
// @Tags         Memories
// @Produce      json
// @Param        memory_id path string true "Memory ID"
// @Success      200 {object} response.StandardResponse{data=dto.MemoryResponse} "Memory retrieved successfully"
// @Failure      400 {object} response.StandardResponse "Invalid memory ID"
// @Failure      404 {object} response.StandardResponse "Memory not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /memories/{memory_id} [get]
func (h *memoryHandler) GetMemory(c echo.Context) error {
	memoryIDStr := c.Param("memory_id")
	memoryID, err := uuid.Parse(memoryIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidMemoryID))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	memory, err := h.memoryService.GetMemory(ctx, userID, memoryID)
	if err != nil {
		if err == apperrors.ErrMemoryNotFound {
			return c.JSON(http.StatusNotFound, h.responseBuilder.Error(err))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.MemoryRetrievedSuccessfully, memory))
}

// @Summary      List Memories
// @Description  Lists all memories for a hangout with cursor pagination
// @Tags         Memories
// @Produce      json
// @Param        hangout_id path string true "Hangout ID"
// @Param        after_id query string false "Cursor for pagination (memory ID)"
// @Param        limit query int false "Limit for pagination"
// @Param        sort_dir query string false "Sort direction (asc/desc)"
// @Success      200 {object} response.StandardResponse{data=dto.PaginatedMemories} "Memories retrieved successfully"
// @Failure      400 {object} response.StandardResponse "Invalid hangout ID"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /hangouts/{hangout_id}/memories [get]
func (h *memoryHandler) ListMemories(c echo.Context) error {
	hangoutIDStr := c.Param("hangout_id")
	hangoutID, err := uuid.Parse(hangoutIDStr)
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
		if err == apperrors.ErrInvalidHangoutID {
			return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(err))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.MemoriesRetrievedSuccessfully, memories))
}

// @Summary      Delete Memory
// @Description  Deletes a memory and its associated file
// @Tags         Memories
// @Produce      json
// @Param        memory_id path string true "Memory ID"
// @Success      200 {object} response.StandardResponse "Memory deleted successfully"
// @Failure      400 {object} response.StandardResponse "Invalid memory ID"
// @Failure      404 {object} response.StandardResponse "Memory not found"
// @Failure      500 {object} response.StandardResponse "Internal server error"
// @Security     BearerAuth
// @Router       /memories/{memory_id} [delete]
func (h *memoryHandler) DeleteMemory(c echo.Context) error {
	memoryIDStr := c.Param("memory_id")
	memoryID, err := uuid.Parse(memoryIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, h.responseBuilder.Error(apperrors.ErrInvalidMemoryID))
	}

	userID := c.Get("user_id").(uuid.UUID)
	ctx := c.Request().Context()

	err = h.memoryService.DeleteMemory(ctx, userID, memoryID)
	if err != nil {
		if err == apperrors.ErrMemoryNotFound {
			return c.JSON(http.StatusNotFound, h.responseBuilder.Error(err))
		}
		return c.JSON(http.StatusInternalServerError, h.responseBuilder.Error(err))
	}

	return c.JSON(http.StatusOK, h.responseBuilder.Success(constants.MemoryDeletedSuccessfully, nil))
}
