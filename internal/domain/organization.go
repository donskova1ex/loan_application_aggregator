package domain

import (
	"app_aggregator/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	CreatedAt time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt time.Time      `json:"updated_at,omitempty" db:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" db:"deleted_at"`
	UUID      *uuid.UUID     `json:"uuid" db:"uuid"`
	Name      string         `json:"name" db:"name" validate:"required"`
}

func FromModel(model *models.Organization) *Organization {
	if model == nil {
		return nil
	}

	return &Organization{
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		DeletedAt: model.DeletedAt,
		UUID:      model.UUID,
		Name:      model.Name,
	}
}

func (o *Organization) ToModel() *models.Organization {
	model := &models.Organization{
		Name: o.Name,
	}

	if o.UUID != nil {
		model.UUID = o.UUID
	}

	return model
}
