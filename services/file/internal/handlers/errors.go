package handlers

import (
	"errors"

	"github.com/Ernestgio/Hangout-Planner/services/file/internal/apperrors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func mapErrorToGRPCStatus(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return status.Error(codes.NotFound, apperrors.ErrFileNotFound.Error())
	}

	switch {
	case errors.Is(err, apperrors.ErrInvalidFileSize),
		errors.Is(err, apperrors.ErrFileTooLarge),
		errors.Is(err, apperrors.ErrInvalidFilename),
		errors.Is(err, apperrors.ErrInvalidFileExtension),
		errors.Is(err, apperrors.ErrInvalidMimeType),
		errors.Is(err, apperrors.ErrInvalidMemoryID):
		return status.Error(codes.InvalidArgument, err.Error())
	}

	switch {
	case errors.Is(err, apperrors.ErrFileUploadFailed),
		errors.Is(err, apperrors.ErrFileDeleteFailed),
		errors.Is(err, apperrors.ErrPresignedUploadURLFailed),
		errors.Is(err, apperrors.ErrPresignedDownloadURLFailed):
		return status.Error(codes.Internal, err.Error())
	}

	switch {
	case errors.Is(err, apperrors.ErrFileCreationFailed),
		errors.Is(err, apperrors.ErrFileStatusUpdateFailed):
		return status.Error(codes.Internal, err.Error())
	}

	return status.Error(codes.Internal, apperrors.ErrInternalServer.Error())
}
