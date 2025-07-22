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
	//TODO: Описать поля параметров. Уточнить у Антона еще раз
}
