package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"app_aggregator/internal/domain"
	"app_aggregator/pkg/validators"

	"github.com/google/uuid"
)

type HTTPLoanApplicationHandler struct {
	service domain.LoanApplicationService
	logger  *slog.Logger
}

func NewHTTPLoanApplicationHandler(service domain.LoanApplicationService, logger *slog.Logger) *HTTPLoanApplicationHandler {
	return &HTTPLoanApplicationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *HTTPLoanApplicationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	applications, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("failed to get loan applications", slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, applications)
}

func (h *HTTPLoanApplicationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uuidStr := r.PathValue("uuid")
	if uuidStr == "" {
		h.logger.Error("missing UUID parameter")
		h.writeError(w, http.StatusBadRequest, "Missing UUID parameter")
		return
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		h.logger.Error("invalid UUID format", slog.String("uuid", uuidStr), slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	application, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get loan application", slog.String("uuid", id.String()), slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, application)
}

func (h *HTTPLoanApplicationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateLoanApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if !validators.ValidPhone(req.Phone) {
		h.logger.Error("invalid phone number", slog.String("phone", req.Phone))
		h.writeError(w, http.StatusBadRequest, "Invalid phone number format")
		return
	}

	normalizedPhone, err := validators.PhoneNormalization(req.Phone)
	if err != nil {
		h.logger.Error("failed to normalize phone number", slog.String("phone", req.Phone), slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid phone number")
		return
	}

	app := &domain.LoanApplication{
		IncomingOrganizationName: req.IncomingOrganizationName,
		IssueOrganizationName:    req.IssueOrganizationName,
		Value:                    req.Value,
		Phone:                    normalizedPhone,
		Comment:                  req.Comment,
	}

	createdApp, err := h.service.Create(ctx, app)
	if err != nil {
		h.logger.Error("failed to create loan application", slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, createdApp)
}

func (h *HTTPLoanApplicationHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uuidStr := r.PathValue("uuid")
	if uuidStr == "" {
		h.logger.Error("missing UUID parameter")
		h.writeError(w, http.StatusBadRequest, "Missing UUID parameter")
		return
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		h.logger.Error("invalid UUID format", slog.String("uuid", uuidStr), slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	var req UpdateLoanApplicationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Валидация и нормализация телефона (если передан)
	var normalizedPhone string
	if req.Phone != "" {
		if !validators.ValidPhone(req.Phone) {
			h.logger.Error("invalid phone number", slog.String("phone", req.Phone))
			h.writeError(w, http.StatusBadRequest, "Invalid phone number format")
			return
		}

		var err error
		normalizedPhone, err = validators.PhoneNormalization(req.Phone)
		if err != nil {
			h.logger.Error("failed to normalize phone number", slog.String("phone", req.Phone), slog.String("error", err.Error()))
			h.writeError(w, http.StatusBadRequest, "Invalid phone number")
			return
		}
	}

	app := &domain.LoanApplication{
		IncomingOrganizationName: req.IncomingOrganizationName,
		IssueOrganizationName:    req.IssueOrganizationName,
		Value:                    req.Value,
		Phone:                    normalizedPhone,
		Comment:                  req.Comment,
	}

	updatedApp, err := h.service.Update(ctx, id, app)
	if err != nil {
		h.logger.Error("failed to update loan application", slog.String("uuid", id.String()), slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, updatedApp)
}

func (h *HTTPLoanApplicationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	uuidStr := r.PathValue("uuid")
	if uuidStr == "" {
		h.logger.Error("missing UUID parameter")
		h.writeError(w, http.StatusBadRequest, "Missing UUID parameter")
		return
	}

	id, err := uuid.Parse(uuidStr)
	if err != nil {
		h.logger.Error("invalid UUID format", slog.String("uuid", uuidStr), slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid UUID format")
		return
	}

	if err := h.service.Delete(ctx, id); err != nil {
		h.logger.Error("failed to delete loan application", slog.String("uuid", id.String()), slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *HTTPLoanApplicationHandler) handleError(w http.ResponseWriter, err error) {
	baseHandler := NewBaseHandler(h.logger)
	baseHandler.handleLoanApplicationError(w, err)
}

func (h *HTTPLoanApplicationHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	baseHandler := NewBaseHandler(h.logger)
	baseHandler.writeJSON(w, status, data)
}

func (h *HTTPLoanApplicationHandler) writeError(w http.ResponseWriter, status int, message string) {
	baseHandler := NewBaseHandler(h.logger)
	baseHandler.writeError(w, status, message)
}
