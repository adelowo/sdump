package postgres

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	"github.com/uptrace/bun/extra/bunotel"
)

func New(dsn string, logQueries bool) (*bun.DB, error) {
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(pgdb, pgdialect.New())

	db.AddQueryHook(bunotel.NewQueryHook(bunotel.WithDBName("getclaimclaim")))

	if logQueries {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, db.Ping()
}
