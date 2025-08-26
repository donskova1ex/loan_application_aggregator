package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"app_aggregator/internal"
	"app_aggregator/internal/domain"

	"github.com/google/uuid"
)

type CreateOrganizationRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateOrganizationRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateLoanApplicationRequest struct {
	IncomingOrganizationName string `json:"incoming_organization_name" validate:"required"`
	IssueOrganizationName    string `json:"issue_organization_name" validate:"required"`
	Value                    int64  `json:"value" validate:"required"`
	Phone                    string `json:"phone" validate:"required"`
	Comment                  string `json:"comment"`
}

type UpdateLoanApplicationRequest struct {
	IncomingOrganizationName string `json:"incoming_organization_name"`
	IssueOrganizationName    string `json:"issue_organization_name"`
	Value                    int64  `json:"value"`
	Phone                    string `json:"phone"`
	Comment                  string `json:"comment"`
}

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

func (h *HTTPOrganizationHandler) handleError(w http.ResponseWriter, err error) {
	switch {
	case err == internal.ErrRecordNoFound:
		h.writeError(w, http.StatusNotFound, "Organization not found")
	default:
		h.writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *HTTPOrganizationHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *HTTPOrganizationHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		h.logger.Error("failed to encode error response", slog.String("error", err.Error()))
	}
}

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

	app := &domain.LoanApplication{
		IncomingOrganizationName: req.IncomingOrganizationName,
		IssueOrganizationName:    req.IssueOrganizationName,
		Value:                    req.Value,
		Phone:                    req.Phone,
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

	app := &domain.LoanApplication{
		IncomingOrganizationName: req.IncomingOrganizationName,
		IssueOrganizationName:    req.IssueOrganizationName,
		Value:                    req.Value,
		Phone:                    req.Phone,
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
	switch {
	case err == internal.ErrRecordNoFound:
		h.writeError(w, http.StatusNotFound, "Loan application not found")
	case err == internal.ErrPhoneNumberExistToday:
		h.writeError(w, http.StatusConflict, "Phone number already exists today")
	default:
		h.writeError(w, http.StatusInternalServerError, "Internal server error")
	}
}

func (h *HTTPLoanApplicationHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("failed to encode JSON response", slog.String("error", err.Error()))
	}
}

func (h *HTTPLoanApplicationHandler) writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := map[string]string{"error": message}
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		h.logger.Error("failed to encode error response", slog.String("error", err.Error()))
	}
}
