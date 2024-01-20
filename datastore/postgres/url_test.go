//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/adelowo/sdump"
	"github.com/stretchr/testify/require"
)

func TestURLRepositoryTable_Create(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
	defer teardownFunc()

	urlStore := NewURLRepositoryTable(client)

	require.NoError(t, urlStore.Create(context.Background(), sdump.NewURLEndpoint()))
}
