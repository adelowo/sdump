package httpd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/adelowo/sdump/internal/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type urlHandler struct {
	logger     *logrus.Entry
	urlRepo    sdump.URLRepository
	ingestRepo sdump.IngestRepository
	cfg        config.Config
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

	_ = render.Render(w, r, &createdURLEndpointResponse{
		APIStatus: newAPIStatus(http.StatusOK, "created url endpoint"),
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
		logger.WithError(err).Error("could not copy request body")
		_ = render.Render(w, r, newAPIError(http.StatusInternalServerError,
			"could not copy request body"))
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

	_ = render.Render(w, r, newAPIStatus(http.StatusAccepted,
		"Request ingested"))
}
