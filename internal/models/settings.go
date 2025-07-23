package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Settings struct {
	gorm.Model
	UUID             *uuid.UUID   `gorm:"type:uuid;default:uuid_generate_v4()"`
	OrganisationUUID *uuid.UUID   `gorm:"type:uuid;not null;index"`
	Organization     Organization `gorm:"foreignKey:OrganisationUUID;references:UUID"`
	NewClient        bool         `gorm:"default:false" validate:"omitempty"`
	PDN              int64        `gorm:"default:0;check:pdn>=0 AND pdn <= 80" validate:"min=0,max=80"`
	HasDebt          bool         `gorm:"default:false;" validate:"omitempty"`
}
