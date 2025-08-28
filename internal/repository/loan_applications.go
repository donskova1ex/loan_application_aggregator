package repository

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"app_aggregator/internal/models"
	"context"
	"errors"
	"fmt"

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
	result := r.Repository.db.WithContext(ctx).
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
	result := r.Repository.db.WithContext(ctx).
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
	incomingOrg, err := r.FindOrganizationByName(ctx, loanApplication.IncomingOrganizationName)
	if err != nil {
		return nil, err
	}
	doverixClientId, err := findClient(ctx, loanApplication.Phone, r.Repository.doverixDb)
	if err != nil {
		return nil, err
	}
	kassaClientId, err := findClient(ctx, loanApplication.Phone, r.Repository.kassaDb)
	if err != nil {
		return nil, err
	}
	denedClientId, err := findClient(ctx, loanApplication.Phone, r.Repository.deDb)
	if err != nil {
		return nil, err
	}

	var issueOrg *domain.Organization

	if doverixClientId == "" && kassaClientId == "" && denedClientId == "" {
		issueOrg, err = r.FindOrganizationByName(ctx, "dened.ru")
		if err != nil {
			return nil, err
		}
	}

	fmt.Println("doverixClientId: ", doverixClientId)
	fmt.Println("kassaClientId: ", kassaClientId)
	fmt.Println("denedClientId: ", denedClientId)

	model := loanApplication.ToModel()
	model.IncomingOrganizationUuid = incomingOrg.UUID
	model.IssueOrganizationUuid = issueOrg.UUID

	var count int64
	err = r.Repository.db.WithContext(ctx).Table("loan_applications").Where("phone = ? AND created_at::date = CURRENT_DATE", model.Phone).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, internal.ErrPhoneNumberExistToday
	}

	result := r.Repository.db.WithContext(ctx).Table("loan_applications").Create(model)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.Repository.db.WithContext(ctx).
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
	result := r.Repository.db.WithContext(ctx).
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
		incomingOrg, err := r.FindOrganizationByName(ctx, loanApplication.IncomingOrganizationName)
		if err != nil {
			return nil, err
		}
		existingApplication.IncomingOrganizationUuid = incomingOrg.UUID
	}
	if loanApplication.IssueOrganizationName != "" {
		issueOrg, err := r.FindOrganizationByName(ctx, loanApplication.IssueOrganizationName)
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

	result = r.Repository.db.WithContext(ctx).Table("loan_applications").Save(existingApplication)
	if result.Error != nil {
		return nil, result.Error
	}

	result = r.Repository.db.WithContext(ctx).
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

	result := r.Repository.db.WithContext(ctx).Table("loan_applications").Where("uuid = ?", id).First(loan_application)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return internal.ErrRecordNoFound
		}
		return result.Error
	}

	result = r.Repository.db.WithContext(ctx).Table("loan_applications").Delete(loan_application)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *LoanApplicationsRepository) FindOrganizationByName(ctx context.Context, name string) (*domain.Organization, error) {
	organization := &models.Organization{}
	result := r.Repository.db.WithContext(ctx).Table("organizations").Where("name = ?", name).First(organization)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, internal.ErrRecordNoFound
		}
		return nil, result.Error
	}
	return domain.FromModel(organization), nil
}

func findClient(ctx context.Context, phoneNumber string, db *gorm.DB) (string, error) {
	type clientResult struct {
		ClientID    string `gorm:"column:ClientId"`
		ClearNumber string `gorm:"column:ClearNumber"`
	}

	var result clientResult
	query := `
		SELECT pn.ClientId, pn.ClearNumber 
		FROM PhoneNumbers pn 
		INNER JOIN Clients c ON pn.Client_Id = c.id 
		  AND c.PrimaryPhoneNumberId = pn.Id
		WHERE pn.ClearNumber = ?
	`

	err := db.WithContext(ctx).Raw(query, phoneNumber).First(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil
		}
		return "", err
	}

	return result.ClientID, nil
}

func clientActiveLoanCheck(ctx context.Context, db *gorm.DB, id string) (bool, error) {

	var activeLoansCount int64
	result := db.WithContext(ctx).Table("Loans").Where("ClientId = ? AND IsActive = ?", id, 1).Count(&activeLoansCount)
	if result.Error != nil {
		return false, result.Error
	}
	return activeLoansCount > 0, nil
}

func clientHasLoans(ctx context.Context, db *gorm.DB, id string) (bool, error) {
	var activeLoansCount int64
	result := db.WithContext(ctx).Table("Loans").Where("ClientId = ?", id).Count(&activeLoansCount)
	if result.Error != nil {
		return false, result.Error
	}
	return activeLoansCount > 0, nil
}

func clientActiveLoanNumber(ctx context.Context, db *gorm.DB, id string) (string, error) {
	var activeLoanNumber string
	result := db.WithContext(ctx).Table("Loans").Select("Number").Where("ClientId = ? AND IsActive = ?", id, 1).First(&activeLoanNumber)
	if result.Error != nil {
		return "", result.Error
	}
	return activeLoanNumber, nil
}

func clientLastPdn(ctx context.Context, db *gorm.DB, id string) (string, error) {
	var lastPdn string
	result := db.WithContext(ctx).Table("Loans").Select("Pdn").Where("ClientId = ?", id).Order("created_at DESC").First(&lastPdn)
	if result.Error != nil {
		return "", result.Error
	}
	return lastPdn, nil
}
