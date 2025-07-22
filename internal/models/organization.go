package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Organization struct {
	gorm.Model
	UUID *uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();not null;uniqueIndex" validate:"omitempty,uuid4"`
	Name string     `gorm:"type:varchar(255);not null;uniqueIndex" validate:"required,min=2,max=150,name"`
}
