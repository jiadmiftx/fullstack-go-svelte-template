package unit_use_case

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type UnitUseCase interface {
	GetUnit(ctx context.Context, id uuid.UUID) (Unit, int, error)
	GetUnits(ctx context.Context, filter UnitFilter, paginate *paginate_utils.PaginateData) ([]Unit, int, error)
	CreateUnit(ctx context.Context, organizationId uuid.UUID, req CreateUnitRequest) (Unit, int, error)
	UpdateUnit(ctx context.Context, id uuid.UUID, req UpdateUnitRequest) (Unit, int, error)
	DeleteUnit(ctx context.Context, id uuid.UUID) (int, error)
}
