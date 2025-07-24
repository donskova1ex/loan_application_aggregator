package handlers

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"app_aggregator/internal/repository"
	"app_aggregator/pkg/validators"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoanApplicationsHandler struct {
	repo *repository.LoanApplicationsRepository
}

func NewLoanApplicationsHandler(repo *repository.LoanApplicationsRepository) *LoanApplicationsHandler {
	return &LoanApplicationsHandler{
		repo: repo,
	}
}

func (h *LoanApplicationsHandler) GetAll(c *gin.Context) {
	loanApplications, err := h.repo.FindAll()
	if err != nil {
		switch errors.Is(err, internal.ErrRecordNoFound) {
		case true:
			apiErr := HandleError(http.StatusNotFound, "loan applications not found", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "loan applications search failed", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
	}
	c.JSON(http.StatusOK, loanApplications)
}

func (h *LoanApplicationsHandler) Create(c *gin.Context) {
	domainLoanApplication := &domain.LoanApplication{}
	if err := c.BindJSON(domainLoanApplication); err != nil {
		apiErr := HandleError(http.StatusBadRequest, "loan application invalid", err)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
	}

	if !validators.ValidPhone(domainLoanApplication.Phone) {
		apiErr := HandleError(http.StatusBadRequest, "phone in loan application invalid", internal.ErrPhoneFormat)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}

	incomingOrganization, err := h.repo.FindOrganizationByName(domainLoanApplication.IncomingOrganizationName)
	if err != nil {
		switch errors.Is(err, internal.ErrRecordNoFound) {
		case true:
			apiErr := HandleError(http.StatusNotFound, "loan application incoming organization not found", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "loan application incoming organization search failed", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
	}

	issueOrganization, err := h.repo.FindOrganizationByName(domainLoanApplication.IssueOrganizationName)
	if err != nil {
		switch errors.Is(err, internal.ErrRecordNoFound) {
		case true:
			apiErr := HandleError(http.StatusNotFound, "loan application issue organization not found", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "loan application issue organization search failed", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
	}

	modelLoanApplication := &models.LoanApplication{
		IncomingOrganizationUuid: incomingOrganization.UUID,
		IssueOrganizationUuid:    issueOrganization.UUID,
		Value:                    domainLoanApplication.Value,
		Phone:                    domainLoanApplication.Phone,
	}

	createdLoanApplication, err := h.repo.Create(modelLoanApplication)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrPhoneNumberExistToday):
			apiErr := HandleError(http.StatusConflict, "loan application creation failed", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "loan application creation failed", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}

	}
	domainLoanApplication.CreatedAt = createdLoanApplication.CreatedAt
	domainLoanApplication.UpdatedAt = createdLoanApplication.UpdatedAt
	domainLoanApplication.DeletedAt = createdLoanApplication.DeletedAt
	domainLoanApplication.UUID = createdLoanApplication.UUID

	c.JSON(http.StatusCreated, domainLoanApplication)
}
