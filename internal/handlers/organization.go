package handlers

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"app_aggregator/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type OrganizationHandler struct {
	repo *repository.Organization
}

func NewOrganizationHandler(repo *repository.Organization) *OrganizationHandler {
	return &OrganizationHandler{
		repo: repo,
	}
}

func (o *OrganizationHandler) GetAll(c *gin.Context) {
	organizations, err := o.repo.GetAll()
	if err != nil {
		switch errors.Is(err, internal.ErrRecordNoFound) {
		case true:
			apiErr := HandleError(http.StatusNotFound, "organizations not found", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "error getting organizations", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
	}
	c.JSON(http.StatusOK, organizations)
}

func (o *OrganizationHandler) Create(c *gin.Context) {
	domainOrg := &domain.Organization{}
	if err := c.BindJSON(domainOrg); err != nil {
		apiErr := HandleError(http.StatusBadRequest, "error parsing body", err)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}
	org := &models.Organization{
		Name: domainOrg.Name,
	}
	createdOrg, err := o.repo.Create(org)
	if err != nil {
		apiErr := HandleError(http.StatusInternalServerError, "error creating organization", err)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}
	c.JSON(http.StatusCreated, createdOrg)
}

func (o *OrganizationHandler) Update(c *gin.Context) {
	domainOrg := &domain.Organization{}
	uuidStr := c.Param("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		apiErr := HandleError(http.StatusBadRequest, "error parsing uuid", err)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}
	if err := c.BindJSON(domainOrg); err != nil {
		apiErr := HandleError(http.StatusBadRequest, "error parsing body", err)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}
	org := &models.Organization{
		Name: domainOrg.Name,
		UUID: &uuid,
	}
	updatedOrg, err := o.repo.Update(org)
	if err != nil {
		switch errors.Is(err, internal.ErrRecordNoFound) {
		case true:
			apiErr := HandleError(http.StatusNotFound, "organization not found", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "error updating organization", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
	}
	c.JSON(http.StatusOK, updatedOrg)
}
func (o *OrganizationHandler) Delete(c *gin.Context) {
	uuidStr := c.Param("uuid")
	uuid, err := uuid.Parse(uuidStr)
	if err != nil {
		apiErr := HandleError(http.StatusBadRequest, "error parsing uuid", err)
		c.AbortWithStatusJSON(apiErr.Status, apiErr)
		return
	}
	err = o.repo.Delete(&uuid)
	if err != nil {
		switch errors.Is(err, internal.ErrRecordNoFound) {
		case true:
			apiErr := HandleError(http.StatusNotFound, "organization not found", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		default:
			apiErr := HandleError(http.StatusInternalServerError, "error deleting organization", err)
			c.AbortWithStatusJSON(apiErr.Status, apiErr)
			return
		}
	}
	c.JSON(http.StatusOK, nil)
}
