package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/InternalTransfer/internal/apperror"
	"github.com/InternalTransfer/internal/dto"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if v != nil {
		json.NewEncoder(w).Encode(v)
	}
}

func httpStatusForError(code string) int {
	switch code {
	case apperror.CodeValidation:
		return http.StatusBadRequest
	case apperror.CodeNotFound:
		return http.StatusNotFound
	case apperror.CodeConflict:
		return http.StatusConflict
	case apperror.CodeInsufficientBalance:
		return http.StatusUnprocessableEntity
	default:
		return http.StatusInternalServerError
	}
}

func mapErrorToResponse(w http.ResponseWriter, err error, logger *slog.Logger) {
	var appErr apperror.AppError
	if errors.As(err, &appErr) {
		status := httpStatusForError(appErr.Code())
		writeJSON(w, status, dto.ErrorResponse{Code: appErr.Code(), Message: appErr.Error()})
		return
	}

	logger.Error("unhandled error", "error", err)
	writeJSON(w, http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.CodeInternal, Message: "internal server error"})
}
