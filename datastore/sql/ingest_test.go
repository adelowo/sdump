//go:build integration
// +build integration

package sql

import (
	"context"
	"testing"

	"github.com/adelowo/sdump"
	"github.com/stretchr/testify/require"
)

func TestIngestRepository_Create(t *testing.T) {
	client, teardownFunc := setupPostgresDatabase(t)
	defer teardownFunc()

	ingestStore := NewIngestRepository(client)

	urlStore := NewURLRepositoryTable(client)

	endpoint, err := urlStore.Get(context.Background(), &sdump.FindURLOptions{
		Reference: "cmltfm6g330l5l1vq110", // see fixtures/urls.yml
	})
	require.NoError(t, err)

	require.NoError(t, ingestStore.Create(context.Background(), &sdump.IngestHTTPRequest{
		UrlID: endpoint.ID,
		Request: sdump.RequestDefinition{
			Body: "{}",
		},
	}))
}

func TestIngestRepository_Sqllite_Create(t *testing.T) {
	client, teardownFunc := setupSqliteDatabase(t)
	defer teardownFunc()

	ingestStore := NewIngestRepository(client)

	urlStore := NewURLRepositoryTable(client)

	// the models use uuid, tricky bit now is setting up sqlite3 uuid support
	// This currently fails as the migrations fail

	endpoint, err := urlStore.Get(context.Background(), &sdump.FindURLOptions{
		Reference: "cmltfm6g330l5l1vq110", // see fixtures/urls.yml
	})
	require.NoError(t, err)

	require.NoError(t, ingestStore.Create(context.Background(), &sdump.IngestHTTPRequest{
		UrlID: endpoint.ID,
		Request: sdump.RequestDefinition{
			Body: "{}",
		},
	}))
}
