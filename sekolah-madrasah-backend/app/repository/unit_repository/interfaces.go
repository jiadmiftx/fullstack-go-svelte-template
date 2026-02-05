package unit_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"
)

type UnitRepository interface {
	GetUnit(ctx context.Context, filter UnitFilter) (Unit, int, error)
	GetUnits(ctx context.Context, filter UnitFilter, paginate *paginate_utils.PaginateData) ([]Unit, int, error)
	CreateUnit(ctx context.Context, perumahan Unit) (Unit, int, error)
	UpdateUnit(ctx context.Context, filter UnitFilter, perumahan Unit) (int, error)
	DeleteUnit(ctx context.Context, filter UnitFilter) (int, error)
}
