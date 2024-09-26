package sqlite

import (
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func New(dsn string, logQueries bool) (*bun.DB, error) {
	sqlite, err := sql.Open(sqliteshim.ShimName, dsn)
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqlite, sqlitedialect.New())

	if logQueries {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}

	return db, db.Ping()
}
