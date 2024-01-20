package postgres

import (
	"context"

	"github.com/adelowo/sdump"
	"github.com/uptrace/bun"
)

type urlRepositoryTable struct {
	inner *bun.DB
}

func NewURLRepositoryTable(db *bun.DB) sdump.URLRepository {
	return &urlRepositoryTable{
		inner: db,
	}
}

func (u *urlRepositoryTable) Create(ctx context.Context,
	model *sdump.URLEndpoint,
) error {
	_, err := bun.NewInsertQuery(u.inner).Model(model).
		Exec(ctx)
	return err
}
