package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/internal/util"
	"github.com/uptrace/bun"
)

type planRepositoryTable struct {
	inner *bun.DB
}

func NewPlanRepositoryTable(db *bun.DB) sdump.PlanRepository {
	return &planRepositoryTable{
		inner: db,
	}
}

func (u *planRepositoryTable) Get(ctx context.Context,
	opts *sdump.FindPlanOptions,
) (*sdump.Plan, error) {
	res := new(sdump.Plan)

	query := bun.NewSelectQuery(u.inner).Model(res)

	if !util.IsStringEmpty(opts.HumanReadableName) {
		query = query.Where("human_readable_name = ?", opts.HumanReadableName)
	}

	err := query.Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sdump.ErrURLEndpointNotFound
	}

	return res, err
}
