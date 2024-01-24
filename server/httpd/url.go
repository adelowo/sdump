package httpd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/r3labs/sse/v2"
	"github.com/sirupsen/logrus"
)

type urlHandler struct {
	logger     *logrus.Entry
	urlRepo    sdump.URLRepository
	ingestRepo sdump.IngestRepository
	cfg        config.Config
	sseServer  *sse.Server
}

func (u *urlHandler) create(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "urlHandler.create").
		WithField("request_id", retrieveRequestID(r))

	logger.Debug("Creating new url endpoint")

	ctx := r.Context()

	endpoint := sdump.NewURLEndpoint()

	if err := u.urlRepo.Create(ctx, endpoint); err != nil {
		logger.WithError(err).Error("could not create url endpoint")

		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError,
			"an error occurred while generating endpoint"))
		return
	}

	go func() {
		_ = u.sseServer.CreateStream(endpoint.PubChannel())
	}()

	_ = render.Render(w, r, &createdURLEndpointResponse{
		APIStatus: newAPIStatus(http.StatusOK, "created url endpoint"),
		SSE: struct {
			Channel string "json:\"channel,omitempty\""
		}{
			Channel: endpoint.PubChannel(),
		},
		URL: struct {
			FQDN                  string "json:\"fqdn,omitempty\""
			Identifier            string "json:\"identifier,omitempty\""
			HumanReadableEndpoint string "json:\"human_readable_endpoint,omitempty\""
		}{
			FQDN:       u.cfg.HTTP.Domain,
			Identifier: endpoint.Reference,
			HumanReadableEndpoint: fmt.Sprintf("%s/%s",
				u.cfg.HTTP.Domain, endpoint.Reference),
		},
	})
}

func (u *urlHandler) ingest(w http.ResponseWriter, r *http.Request) {
	reference := chi.URLParam(r, "reference")

	logger := u.logger.WithField("request_id", retrieveRequestID(r)).
		WithField("method", "urlHandler.ingest").
		WithField("reference", reference)

	logger.Debug("Ingesting http request")

	r.Body = http.MaxBytesReader(w, r.Body, u.cfg.HTTP.MaxRequestBodySize)

	ctx := r.Context()

	endpoint, err := u.urlRepo.Get(ctx, &sdump.FindURLOptions{
		Reference: reference,
	})
	if errors.Is(err, sdump.ErrURLEndpointNotFound) {
		_ = render.Render(w, r, newAPIError(http.StatusNotFound,
			"Dump url does not exist"))
		return
	}

	if err != nil {
		logger.WithError(err).Error("could not find dump url by reference")
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError,
			"an error occurred while ingesting HTTP request"))
		return
	}

	s := &strings.Builder{}

	size, err := io.Copy(s, r.Body)
	if err != nil {
		msg := "could not copy request body"
		if maxErr, ok := err.(*http.MaxBytesError); ok {
			msg = maxErr.Error()
		}

		logger.WithError(err).Error("could not copy request body")
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError,
			msg))
		return
	}

	ingestedRequest := &sdump.IngestHTTPRequest{
		UrlID: endpoint.ID,
		Request: sdump.RequestDefinition{
			Body:      s.String(),
			Query:     r.URL.Query().Encode(),
			Headers:   r.Header,
			IPAddress: util.GetIP(r),
			Size:      size,
		},
	}

	if err := u.ingestRepo.Create(ctx, ingestedRequest); err != nil {
		logger.WithError(err).Error("could not ingest request")
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError,
			"an error occurred while ingesting request"))
		return
	}

	go func() {
		if !u.sseServer.StreamExists(endpoint.PubChannel()) {
			_ = u.sseServer.CreateStream(endpoint.PubChannel())
		}

		b := new(bytes.Buffer)

		var sseEvent struct {
			Request   sdump.RequestDefinition `json:"request"`
			ID        string                  `json:"id"`
			CreatedAt time.Time               `json:"created_at,omitempty"`
		}

		sseEvent.Request = ingestedRequest.Request
		sseEvent.ID = ingestedRequest.ID.String()
		sseEvent.CreatedAt = ingestedRequest.CreatedAt

		if err := json.NewEncoder(b).Encode(&sseEvent); err != nil {
			logger.WithError(err).Error("could not format SSE event")
			return
		}

		u.sseServer.Publish(endpoint.PubChannel(), &sse.Event{
			ID:   []byte(ingestedRequest.ID.String()),
			Data: b.Bytes(),
		})
	}()

	_ = render.Render(w, r, newAPIStatus(http.StatusAccepted,
		"Request ingested"))
}
