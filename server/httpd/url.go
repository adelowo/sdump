package httpd

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type urlHandler struct {
	logger *logrus.Entry
}

func (u *urlHandler) create(w http.ResponseWriter, r *http.Request) {
	logger := u.logger.WithField("method", "urlHandler.create").
		WithField("request_id", retrieveRequestID(r))

	logger.Debug("Creating new url endpoint")
}
