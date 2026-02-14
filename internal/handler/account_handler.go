package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/InternalTransfer/internal/apperror"
	"github.com/InternalTransfer/internal/dto"
	"github.com/InternalTransfer/internal/service"
)

type AccountHandler struct {
	accountSvc *service.AccountService
	logger     *slog.Logger
}

func NewAccountHandler(accountSvc *service.AccountService, logger *slog.Logger) *AccountHandler {
	return &AccountHandler{
		accountSvc: accountSvc,
		logger:     logger,
	}
}

func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Code: apperror.CodeValidation, Message: "Invalid request format. Please check your input and try again"})
		return
	}

	if err := h.accountSvc.Create(r.Context(), req.AccountID, req.InitialBalance); err != nil {
		mapErrorToResponse(w, err, h.logger)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *AccountHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("account_id")
	accountID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, dto.ErrorResponse{Code: apperror.CodeValidation, Message: "Invalid account ID. Please provide a valid account number"})
		return
	}

	account, err := h.accountSvc.GetByID(r.Context(), accountID)
	if err != nil {
		mapErrorToResponse(w, err, h.logger)
		return
	}

	writeJSON(w, http.StatusOK, dto.AccountResponse{
		AccountID: account.AccountID,
		Balance:   account.Balance,
	})
}
