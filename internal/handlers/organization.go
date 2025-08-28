package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"app_aggregator/internal/domain"

	"github.com/google/uuid"
)

type HTTPOrganizationHandler struct {
	service domain.OrganizationService
	logger  *slog.Logger
}

func NewHTTPOrganizationHandler(service domain.OrganizationService, logger *slog.Logger) *HTTPOrganizationHandler {
	return &HTTPOrganizationHandler{
		service: service,
		logger:  logger,
	}
}

func (h *HTTPOrganizationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	organizations, err := h.service.GetAll(ctx)
	if err != nil {
		h.logger.Error("failed to get organizations", slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, organizations)
}

// GetByID возвращает организацию по ID
func (h *HTTPOrganizationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	organization, err := h.service.GetByID(ctx, id)
	if err != nil {
		h.logger.Error("failed to get organization", slog.String("uuid", id.String()), slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, organization)
}

// Create создает новую организацию
func (h *HTTPOrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req CreateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	org := &domain.Organization{
		Name: req.Name,
	}

	createdOrg, err := h.service.Create(ctx, org)
	if err != nil {
		h.logger.Error("failed to create organization", slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusCreated, createdOrg)
}

// Update обновляет существующую организацию
func (h *HTTPOrganizationHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req UpdateOrganizationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request body", slog.String("error", err.Error()))
		h.writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	org := &domain.Organization{
		Name: req.Name,
	}

	updatedOrg, err := h.service.Update(ctx, id, org)
	if err != nil {
		h.logger.Error("failed to update organization", slog.String("uuid", id.String()), slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	h.writeJSON(w, http.StatusOK, updatedOrg)
}

// Delete удаляет организацию
func (h *HTTPOrganizationHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
		h.logger.Error("failed to delete organization", slog.String("uuid", id.String()), slog.String("error", err.Error()))
		h.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handleError обрабатывает ошибки для организаций
func (h *HTTPOrganizationHandler) handleError(w http.ResponseWriter, err error) {
	baseHandler := NewBaseHandler(h.logger)
	baseHandler.handleOrganizationError(w, err)
}

// writeJSON отправляет JSON ответ
func (h *HTTPOrganizationHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	baseHandler := NewBaseHandler(h.logger)
	baseHandler.writeJSON(w, status, data)
}

// writeError отправляет JSON ошибку
func (h *HTTPOrganizationHandler) writeError(w http.ResponseWriter, status int, message string) {
	baseHandler := NewBaseHandler(h.logger)
	baseHandler.writeError(w, status, message)
}
