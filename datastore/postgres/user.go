package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adelowo/sdump"
	"github.com/uptrace/bun"
)

type userRepositoryTable struct {
	inner *bun.DB
}

func NewUserRepositoryTable(db *bun.DB) sdump.UserRepository {
	return &userRepositoryTable{
		inner: db,
	}
}

func (u *userRepositoryTable) Create(ctx context.Context,
	model *sdump.User,
) error {
	_, err := bun.NewInsertQuery(u.inner).Model(model).
		Exec(ctx)
	return err
}

func (u *userRepositoryTable) Find(ctx context.Context,
	opts *sdump.FindUserOptions,
) (*sdump.User, error) {
	res := new(sdump.User)

	query := bun.NewSelectQuery(u.inner).Model(res)

	err := query.Scan(ctx, res)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, sdump.ErrUserNotFound
	}

	return res, err
}
