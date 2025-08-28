package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	GetAll(ctx context.Context) ([]*Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	Create(ctx context.Context, org *Organization) (*Organization, error)
	Update(ctx context.Context, org *Organization) (*Organization, error)
	Delete(ctx context.Context, id uuid.UUID) error
	FindByName(ctx context.Context, name string) (*Organization, error)
}

type LoanApplicationRepository interface {
	GetAll(ctx context.Context) ([]*LoanApplication, error)
	GetByID(ctx context.Context, id uuid.UUID) (*LoanApplication, error)
	Create(ctx context.Context, app *LoanApplication) (*LoanApplication, error)
	Update(ctx context.Context, app *LoanApplication) (*LoanApplication, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type OrganizationService interface {
	GetAll(ctx context.Context) ([]*Organization, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	Create(ctx context.Context, org *Organization) (*Organization, error)
	Update(ctx context.Context, id uuid.UUID, org *Organization) (*Organization, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type LoanApplicationService interface {
	GetAll(ctx context.Context) ([]*LoanApplication, error)
	GetByID(ctx context.Context, id uuid.UUID) (*LoanApplication, error)
	Create(ctx context.Context, app *LoanApplication) (*LoanApplication, error)
	Update(ctx context.Context, id uuid.UUID, app *LoanApplication) (*LoanApplication, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
