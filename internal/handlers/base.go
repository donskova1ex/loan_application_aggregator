package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"app_aggregator/internal"
)

type BaseHandler struct {
	logger *slog.Logger
}

func NewBaseHandler(logger *slog.Logger) *BaseHandler {
	return &BaseHandler{
		logger: logger,
	}
}

func (h *BaseHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *BaseHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		h.logger.Error("failed to encode error response", slog.String("error", err.Error()))
	}
}

func (h *BaseHandler) handleOrganizationError(w http.ResponseWriter, err error) {
	switch {
	case err == internal.ErrRecordNoFound:
		h.writeError(w, http.StatusNotFound, "Organization not found")
	default:
		h.writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *BaseHandler) handleLoanApplicationError(w http.ResponseWriter, err error) {
	switch {
	case err == internal.ErrRecordNoFound:
		h.writeError(w, http.StatusNotFound, "Loan application not found")
	case err == internal.ErrPhoneNumberExistToday:
		h.writeError(w, http.StatusConflict, "Phone number already exists today")
	default:
		h.writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}
