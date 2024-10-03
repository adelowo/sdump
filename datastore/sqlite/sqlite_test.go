//go:build integration
// +build integration

package sqlite

import (
	"database/sql"
	"fmt"
	"testing"

	testfixtures "github.com/go-testfixtures/testfixtures/v3"
	"github.com/stretchr/testify/require"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/driver/sqliteshim"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
)

func prepareTestDatabase(t *testing.T, dsn string) {
	t.Helper()

	var err error

	db, err := sql.Open(sqliteshim.ShimName, dsn)
	require.NoError(t, err)

	require.NoError(t, db.Ping())

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	require.NoError(t, err)

	migrator, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", "migrations"), "sqllite3", driver)
	require.NoError(t, err)

	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("sqlite3"),
		testfixtures.Directory("testdata/fixtures"),
	)
	require.NoError(t, err)

	require.NoError(t, fixtures.Load())
}

func setupDatabase(t *testing.T) (*bun.DB, func()) {
	t.Helper()

	dsn := "file::memory:?cache=shared"

	prepareTestDatabase(t, dsn)

	client, err := New(dsn, false)
	require.NoError(t, err)

	return client, func() {
	}
}
