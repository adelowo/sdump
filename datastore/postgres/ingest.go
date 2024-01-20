package postgres

import (
	"context"

	"github.com/adelowo/sdump"
	"github.com/uptrace/bun"
)

type ingestRepository struct {
	inner *bun.DB
}

func NewIngestRepository(db *bun.DB) sdump.IngestRepository {
	return &ingestRepository{
		inner: db,
	}
}

func (u *ingestRepository) Create(ctx context.Context,
	model *sdump.IngestHTTPRequest,
) error {
	_, err := bun.NewInsertQuery(u.inner).Model(model).
		Exec(ctx)
	return err
}
