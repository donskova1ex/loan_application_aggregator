package repository

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"context"
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

func (o *Organization) GetAll(ctx context.Context) ([]*domain.Organization, error) {
	var organizations []*models.Organization
	result := o.Repository.db.WithContext(ctx).Table("organizations").Find(&organizations)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}

	domainOrganizations := make([]*domain.Organization, len(organizations))
	for i, org := range organizations {
		domainOrganizations[i] = domain.FromModel(org)
	}

	return domainOrganizations, nil
}

func (o *Organization) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	organization := &models.Organization{}
	result := o.Repository.db.WithContext(ctx).Table("organizations").Where("uuid = ?", id).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	return domain.FromModel(organization), nil
}

func (o *Organization) Create(ctx context.Context, organization *domain.Organization) (*domain.Organization, error) {
	model := organization.ToModel()

	result := o.Repository.db.WithContext(ctx).Table("organizations").Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return domain.FromModel(model), nil
}

func (o *Organization) Update(ctx context.Context, organization *domain.Organization) (*domain.Organization, error) {
	model := organization.ToModel()

	existingOrganization := &models.Organization{}
	result := o.Repository.db.WithContext(ctx).Table("organizations").Where("uuid = ?", model.UUID).First(existingOrganization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}

	existingOrganization.Name = model.Name

	result = o.Repository.db.WithContext(ctx).Table("organizations").Save(existingOrganization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	return domain.FromModel(existingOrganization), nil
}

func (o *Organization) Delete(ctx context.Context, id uuid.UUID) error {
	organization := &models.Organization{}
	result := o.Repository.db.WithContext(ctx).Table("organizations").Where("uuid = ?", id).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return internal.ErrRecordNoFound
		}
		return result.Error
	}

	result = o.Repository.db.WithContext(ctx).Table("organizations").Delete(organization)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (o *Organization) FindByName(ctx context.Context, name string) (*domain.Organization, error) {
	organization := &models.Organization{}
	result := o.Repository.db.WithContext(ctx).Table("organizations").Where("name = ?", name).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	return domain.FromModel(organization), nil
}
