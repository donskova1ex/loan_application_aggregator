package handlers

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"app_aggregator/internal/repository"
	"app_aggregator/pkg/validators"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
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

	normalizedPhone, err := validators.PhoneNormalization(domainLoanApplication.Phone)
	if err != nil {
		switch {
		case errors.Is(err, internal.ErrEmptyPhoneNumber):
			apiErr := HandleError(http.StatusBadRequest, "phone in loan application is empty", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		case errors.Is(err, internal.ErrInvalidPhoneNumber):
			apiErr := HandleError(http.StatusBadRequest, "phone in loan application invalid", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
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

	testComment := "Redirecting"

	modelLoanApplication := &models.LoanApplication{
		IncomingOrganizationUuid: incomingOrganization.UUID,
		IssueOrganizationUuid:    issueOrganization.UUID,
		Value:                    domainLoanApplication.Value,
		Phone:                    normalizedPhone,
		Comment:                  testComment,
	}

	//testResult, err := h.repo.FindClientHistory(modelLoanApplication)
	//if err != nil {
	//	switch errors.Is(err, internal.ErrRecordNoFound) {
	//	case true:
	//		apiErr := HandleError(http.StatusNotFound, "loan application clients history not found", err)
	//		c.AbortWithStatusJSON(apiErr.Status, apiErr)
	//		return
	//	}
	//}
	//fmt.Println(testResult)
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

	newLoanApplication := &domain.LoanApplication{
		UUID:                     createdLoanApplication.UUID,
		IncomingOrganizationName: incomingOrganization.Name,
		IssueOrganizationName:    issueOrganization.Name,
		Value:                    createdLoanApplication.Value,
		Phone:                    createdLoanApplication.Phone,
		CreatedAt:                createdLoanApplication.CreatedAt,
		UpdatedAt:                createdLoanApplication.UpdatedAt,
		DeletedAt:                createdLoanApplication.DeletedAt,
		Comment:                  createdLoanApplication.Comment,
	}

	c.JSON(http.StatusCreated, newLoanApplication)
}
