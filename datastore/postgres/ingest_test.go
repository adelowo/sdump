//go:build integration
// +build integration

package postgres

import (
	"context"
	"testing"

	"github.com/adelowo/sdump"
	"github.com/stretchr/testify/require"
)

func TestIngestRepository_Create(t *testing.T) {
	client, teardownFunc := setupDatabase(t)
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
