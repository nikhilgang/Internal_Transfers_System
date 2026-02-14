package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/InternalTransfer/internal/apperror"
	"github.com/InternalTransfer/internal/dto"
	"github.com/InternalTransfer/internal/service"
)

type TransactionHandler struct {
	transferSvc *service.TransferService
	logger      *slog.Logger
}

func NewTransactionHandler(transferSvc *service.TransferService, logger *slog.Logger) *TransactionHandler {
	return &TransactionHandler{
		transferSvc: transferSvc,
		logger:      logger,
	}
}

func (h *TransactionHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Code: apperror.CodeValidation, Message: "Invalid request format. Please check your input and try again"})
		return
	}

	if err := h.transferSvc.Transfer(r.Context(), req.SourceAccountID, req.DestinationAccountID, req.Amount); err != nil {
		mapErrorToResponse(w, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
