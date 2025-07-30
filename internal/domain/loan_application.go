package domain

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type LoanApplication struct {
	UUID                     *uuid.UUID     `json:"uuid"`
	IncomingOrganizationName string         `json:"incoming_organization_name" validate:"required"`
	IssueOrganizationName    string         `json:"issue_organization_name" validate:"required"`
	Value                    int64          `json:"value" db:"value" validate:"required"`
	Phone                    string         `json:"phone" db:"phone" validate:"required"`
	CreatedAt                time.Time      `json:"created_at"`
	UpdatedAt                time.Time      `json:"updated_at"`
	DeletedAt                gorm.DeletedAt `json:"deleted_at"`
	Comment                  string         `json:"comment"`
}
