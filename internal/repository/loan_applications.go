package repository

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"errors"
	"gorm.io/gorm"
)

type LoanApplicationsRepository struct {
	Repository *Repository
}

func NewLoanApplicationsRepository(repository *Repository) *LoanApplicationsRepository {
	return &LoanApplicationsRepository{
		Repository: repository,
	}
}
func (r *LoanApplicationsRepository) FindAll() ([]*domain.LoanApplication, error) {
	var loanApplications []*domain.LoanApplication
	//result := r.Repository.db.Table("loan_applications").Find(&loanApplications)
	result := r.Repository.db.
		Table("loan_applications la").
		Select(`
            la.uuid,
            io.name as incoming_organization_name,
            oo.name as issue_organization_name,
            la.value,
            la.phone,
            la.created_at,
            la.updated_at,
            la.deleted_at
        `).
		Joins("LEFT JOIN organizations io ON io.uuid = la.incoming_organization_uuid").
		Joins("LEFT JOIN organizations oo ON oo.uuid = la.issue_organization_uuid").
		Scan(&loanApplications)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}

	return loanApplications, nil
}
func (r *LoanApplicationsRepository) Create(loanApplication *models.LoanApplication) (*models.LoanApplication, error) {
	var count int64
	err := r.Repository.db.Table("loan_applications").Where("phone = ? AND created_at::date = CURRENT_DATE", loanApplication.Phone).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, internal.ErrPhoneNumberExistToday
	}

	result := r.Repository.db.Table("loan_applications").Create(loanApplication)
	if result.Error != nil {
		return nil, result.Error
	}
	return loanApplication, nil
}

func (r *LoanApplicationsRepository) FindOrganizationByName(name string) (*domain.Organization, error) {
	organization := &models.Organization{}
	result := r.Repository.db.Table("organizations").Where("name = ?", name).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	domainOrganization := &domain.Organization{
		CreatedAt: organization.CreatedAt,
		UpdatedAt: organization.UpdatedAt,
		DeletedAt: organization.DeletedAt,
		UUID:      organization.UUID,
		Name:      organization.Name,
	}
	return domainOrganization, nil
}
