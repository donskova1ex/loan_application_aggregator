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

type LoanApplicationsRepository struct {
	Repository *Repository
}

func NewLoanApplicationsRepository(repository *Repository) *LoanApplicationsRepository {
	return &LoanApplicationsRepository{
		Repository: repository,
	}
}

func (r *LoanApplicationsRepository) GetAll(ctx context.Context) ([]*domain.LoanApplication, error) {
	var loanApplications []*models.LoanApplication
	result := r.Repository.db.
		Table("loan_applications").
		Preload("IncomingOrganization").
		Preload("IssueOrganization").
		Find(&loanApplications)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}

	domainApplications := make([]*domain.LoanApplication, len(loanApplications))
	for i, app := range loanApplications {
		domainApplications[i] = domain.LoanApplicationFromModel(app)
	}

	return domainApplications, nil
}

func (r *LoanApplicationsRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.LoanApplication, error) {
	loanApplication := &models.LoanApplication{}
	result := r.Repository.db.
		Table("loan_applications").
		Preload("IncomingOrganization").
		Preload("IssueOrganization").
		Where("uuid = ?", id).
		First(loanApplication)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	return domain.LoanApplicationFromModel(loanApplication), nil
}

func (r *LoanApplicationsRepository) Create(ctx context.Context, loanApplication *domain.LoanApplication) (*domain.LoanApplication, error) {
	incomingOrg, err := r.FindOrganizationByName(loanApplication.IncomingOrganizationName)
	if err != nil {
		return nil, err
	}

	issueOrg, err := r.FindOrganizationByName(loanApplication.IssueOrganizationName)
	if err != nil {
		return nil, err
	}

	model := loanApplication.ToModel()
	model.IncomingOrganizationUuid = incomingOrg.UUID
	model.IssueOrganizationUuid = issueOrg.UUID

	var count int64
	err = r.Repository.db.Table("loan_applications").Where("phone = ? AND created_at::date = CURRENT_DATE", model.Phone).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, internal.ErrPhoneNumberExistToday
	}

	result := r.Repository.db.Table("loan_applications").Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.Repository.db.
		Table("loan_applications").
		Preload("IncomingOrganization").
		Preload("IssueOrganization").
		Where("uuid = ?", model.UUID).
		First(model)
	if result.Error != nil {
		return nil, result.Error
	}

	return domain.LoanApplicationFromModel(model), nil
}

func (r *LoanApplicationsRepository) Update(ctx context.Context, loanApplication *domain.LoanApplication) (*domain.LoanApplication, error) {
	model := loanApplication.ToModel()

	existingApplication := &models.LoanApplication{}
	result := r.Repository.db.
		Table("loan_applications").
		Preload("IncomingOrganization").
		Preload("IssueOrganization").
		Where("uuid = ?", model.UUID).
		First(existingApplication)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}

	if loanApplication.IncomingOrganizationName != "" {
		incomingOrg, err := r.FindOrganizationByName(loanApplication.IncomingOrganizationName)
		if err != nil {
			return nil, err
		}
		existingApplication.IncomingOrganizationUuid = incomingOrg.UUID
	}
	if loanApplication.IssueOrganizationName != "" {
		issueOrg, err := r.FindOrganizationByName(loanApplication.IssueOrganizationName)
		if err != nil {
			return nil, err
		}
		existingApplication.IssueOrganizationUuid = issueOrg.UUID
	}
	if loanApplication.Value != 0 {
		existingApplication.Value = loanApplication.Value
	}
	if loanApplication.Phone != "" {
		existingApplication.Phone = loanApplication.Phone
	}
	if loanApplication.Comment != "" {
		existingApplication.Comment = loanApplication.Comment
	}

	result = r.Repository.db.Table("loan_applications").Save(existingApplication)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.Repository.db.
		Table("loan_applications").
		Preload("IncomingOrganization").
		Preload("IssueOrganization").
		Where("uuid = ?", existingApplication.UUID).
		First(existingApplication)
	if result.Error != nil {
		return nil, result.Error
	}

	return domain.LoanApplicationFromModel(existingApplication), nil
}

func (r *LoanApplicationsRepository) Delete(ctx context.Context, id uuid.UUID) error {
	loan_application := &models.LoanApplication{}

	result := r.Repository.db.Table("loan_applications").Where("uuid = ?", id).First(loan_application)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return internal.ErrRecordNoFound
		}
		return result.Error
	}

	result = r.Repository.db.Table("loan_applications").Delete(loan_application)
	if result.Error != nil {
		return result.Error
	}
	return nil
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
	return domain.FromModel(organization), nil
}

func (r *LoanApplicationsRepository) FindClientHistory(loanApplication *domain.LoanApplication) (*map[string]interface{}, error) {
	return nil, nil
}

func (r *LoanApplicationsRepository) checkKassaHistory(phoneNumber string) *map[string]interface{} {
	mapData := make(map[string]interface{})
	resultKassa := r.Repository.
		kassaDb.
		Table("Clients").
		Select("id, MobileNumber").
		Where("MobileNumber = ?", phoneNumber).
		Scan(&mapData)
	if resultKassa.Error != nil {
		if errors.Is(resultKassa.Error, gorm.ErrRecordNotFound) {
			return nil
		}
	}
	return &mapData
}

func (r *LoanApplicationsRepository) checkDoverixHistory(phoneNumber string) *map[string]interface{} {
	return nil
}
