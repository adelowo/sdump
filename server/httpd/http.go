package httpd

import (
	"net/http"

	"github.com/adelowo/sdump/config"
)

func New(cfg config.Config) *http.Server {
	return &http.Server{}
}
