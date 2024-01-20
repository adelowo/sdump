package httpd

import (
	"fmt"
	"net/http"

	"github.com/adelowo/sdump"
	"github.com/adelowo/sdump/config"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"
)

type urlHandler struct {
	logger  *logrus.Entry
	urlRepo sdump.URLRepository
	cfg     config.Config
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
