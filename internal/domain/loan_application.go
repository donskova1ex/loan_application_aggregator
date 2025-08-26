package domain

import (
	"app_aggregator/internal/models"
	"time"

	"github.com/google/uuid"
)

type LoanApplication struct {
	UUID                     uuid.UUID `json:"uuid"`
	IncomingOrganizationName string    `json:"incoming_organization_name" validate:"required"`
	IssueOrganizationName    string    `json:"issue_organization_name" validate:"required"`
	Value                    int64     `json:"value" validate:"required"`
	Phone                    string    `json:"phone" validate:"required"`
	Comment                  string    `json:"comment"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

func NewLoanApplication(incomingOrgName, issueOrgName, phone string, value int64, comment string) *LoanApplication {
	return &LoanApplication{
		UUID:                     uuid.New(),
		IncomingOrganizationName: incomingOrgName,
		IssueOrganizationName:    issueOrgName,
		Value:                    value,
		Phone:                    phone,
		Comment:                  comment,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}
}

func LoanApplicationFromModel(model *models.LoanApplication) *LoanApplication {
	if model == nil {
		return nil
	}

	app := &LoanApplication{
		IncomingOrganizationName: model.IncomingOrganization.Name,
		IssueOrganizationName:    model.IssueOrganization.Name,
		Value:                    model.Value,
		Phone:                    model.Phone,
		Comment:                  model.Comment,
		CreatedAt:                model.CreatedAt,
		UpdatedAt:                model.UpdatedAt,
	}

	if model.UUID != nil {
		app.UUID = *model.UUID
	}

	return app
}

func (la *LoanApplication) ToModel() *models.LoanApplication {
	model := &models.LoanApplication{
		Value:   la.Value,
		Phone:   la.Phone,
		Comment: la.Comment,
	}

	if la.UUID != uuid.Nil {
		model.UUID = &la.UUID
	}

	return model
}

func (la *LoanApplication) UpdateFromModel(model *models.LoanApplication) {
	if model == nil {
		return
	}

	la.IncomingOrganizationName = model.IncomingOrganization.Name
	la.IssueOrganizationName = model.IssueOrganization.Name
	la.Value = model.Value
	la.Phone = model.Phone
	la.Comment = model.Comment
	la.UpdatedAt = model.UpdatedAt
}
