package repository

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	Repository *Repository
}

func NewOrganizationRepository(r *Repository) *Organization {
	return &Organization{
		Repository: r,
	}
}

func (o *Organization) GetAll() ([]*domain.Organization, error) {
	var organizations []*domain.Organization
	result := o.Repository.db.Table("organizations").Find(&organizations)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	return organizations, nil
}
func (o *Organization) Create(organization *models.Organization) (*domain.Organization, error) {

	result := o.Repository.db.Table("organizations").Create(organization)
	if result.Error != nil {
		return nil, result.Error
	}
	newOrganization := &domain.Organization{
		CreatedAt: organization.CreatedAt,
		UpdatedAt: organization.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
		UUID:      organization.UUID,
		Name:      organization.Name,
	}
	return newOrganization, nil
}

func (o *Organization) Update(organization *models.Organization) (*domain.Organization, error) {

	existingOrganization := &models.Organization{}
	result := o.Repository.db.Table("organizations").Where("uuid = ?", organization.UUID).First(existingOrganization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}

	existingOrganization.Name = organization.Name

	result = o.Repository.db.Table("organizations").Save(existingOrganization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	newOrganization := &domain.Organization{
		UpdatedAt: organization.UpdatedAt,
		DeletedAt: gorm.DeletedAt{},
		UUID:      organization.UUID,
		Name:      organization.Name,
	}
	return newOrganization, nil
}
func (o *Organization) Delete(uuid *uuid.UUID) error {
	organization := &models.Organization{}
	result := o.Repository.db.Table("organizations").Where("uuid = ?", uuid).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return internal.ErrRecordNoFound
		}
		return result.Error
	}

	result = o.Repository.db.Table("organizations").Delete(organization)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (o *Organization) FindByName(name string) (*domain.Organization, error) {
	organization := &models.Organization{}
	result := o.Repository.db.Table("organizations").Where("name = ?", name).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	domainOrganization := &domain.Organization{
		UUID: organization.UUID,
		Name: organization.Name,
	}
	return domainOrganization, nil
}
