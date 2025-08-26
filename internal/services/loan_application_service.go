package services

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"context"

	"github.com/google/uuid"
)

type LoanApplicationService struct {
	repo domain.LoanApplicationRepository
}

func NewLoanApplicationService(repo domain.LoanApplicationRepository) *LoanApplicationService {
	return &LoanApplicationService{
		repo: repo,
	}
}

func (s *LoanApplicationService) GetAll(ctx context.Context) ([]*domain.LoanApplication, error) {
	return s.repo.GetAll(ctx)
}

func (s *LoanApplicationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.LoanApplication, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *LoanApplicationService) Create(ctx context.Context, app *domain.LoanApplication) (*domain.LoanApplication, error) {
	if app.IncomingOrganizationName == "" {
		return nil, internal.ErrInvalidLoanApplication
	}

	return s.repo.Create(ctx, app)
}

func (s *LoanApplicationService) Update(ctx context.Context, id uuid.UUID, app *domain.LoanApplication) (*domain.LoanApplication, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if app.IncomingOrganizationName != "" {
		existing.IncomingOrganizationName = app.IncomingOrganizationName
	}
	if app.IssueOrganizationName != "" {
		existing.IssueOrganizationName = app.IssueOrganizationName
	}
	if app.Value != 0 {
		existing.Value = app.Value
	}
	if app.Phone != "" {
		existing.Phone = app.Phone
	}
	if app.Comment != "" {
		existing.Comment = app.Comment
	}

	return s.repo.Update(ctx, existing)
}

func (s *LoanApplicationService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
