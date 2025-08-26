package services

import (
	"app_aggregator/internal"
	"app_aggregator/internal/domain"
	"context"

	"github.com/google/uuid"
)

type OrganizationService struct {
	repo domain.OrganizationRepository
}

func NewOrganizationService(repo domain.OrganizationRepository) *OrganizationService {
	return &OrganizationService{
		repo: repo,
	}
}

func (s *OrganizationService) GetAll(ctx context.Context) ([]*domain.Organization, error) {
	return s.repo.GetAll(ctx)
}

func (s *OrganizationService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Organization, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *OrganizationService) Create(ctx context.Context, org *domain.Organization) (*domain.Organization, error) {
	if org.Name == "" {
		return nil, internal.ErrInvalidOrganizationName
	}

	return s.repo.Create(ctx, org)
}

func (s *OrganizationService) Update(ctx context.Context, id uuid.UUID, org *domain.Organization) (*domain.Organization, error) {
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	existing.Name = org.Name

	return s.repo.Update(ctx, existing)
}

func (s *OrganizationService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, id)
}
