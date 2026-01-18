package handlers

import (
	"context"

	filepb "github.com/Ernestgio/Hangout-Planner/pkg/shared/proto/gen/go/file"
	"github.com/Ernestgio/Hangout-Planner/services/file/internal/services"
)

type FileHandler struct {
	filepb.UnimplementedFileServiceServer
	fileService services.FileService
}

func NewFileHandler(fileService services.FileService) *FileHandler {
	return &FileHandler{
		fileService: fileService,
	}
}

func (h *FileHandler) GenerateUploadURLs(ctx context.Context, req *filepb.GenerateUploadURLsRequest) (*filepb.GenerateUploadURLsResponse, error) {
	resp, err := h.fileService.GenerateUploadURLs(ctx, req)
	if err != nil {
		return nil, mapErrorToGRPCStatus(err)
	}
	return resp, nil
}

func (h *FileHandler) ConfirmUpload(ctx context.Context, req *filepb.ConfirmUploadRequest) (*filepb.ConfirmUploadResponse, error) {
	resp, err := h.fileService.ConfirmUpload(ctx, req)
	if err != nil {
		return nil, mapErrorToGRPCStatus(err)
	}
	return resp, nil
}

func (h *FileHandler) GetFileByMemoryID(ctx context.Context, req *filepb.GetFileByMemoryIDRequest) (*filepb.GetFileByMemoryIDResponse, error) {
	resp, err := h.fileService.GetFileByMemoryID(ctx, req)
	if err != nil {
		return nil, mapErrorToGRPCStatus(err)
	}
	return resp, nil
}

func (h *FileHandler) GetFilesByMemoryIDs(ctx context.Context, req *filepb.GetFilesByMemoryIDsRequest) (*filepb.GetFilesByMemoryIDsResponse, error) {
	resp, err := h.fileService.GetFilesByMemoryIDs(ctx, req)
	if err != nil {
		return nil, mapErrorToGRPCStatus(err)
	}
	return resp, nil
}

func (h *FileHandler) DeleteFile(ctx context.Context, req *filepb.DeleteFileRequest) (*filepb.DeleteFileResponse, error) {
	resp, err := h.fileService.DeleteFile(ctx, req)
	if err != nil {
		return nil, mapErrorToGRPCStatus(err)
	}
	return resp, nil
}
