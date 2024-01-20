package httpd

import (
	"net/http"

	"github.com/go-chi/render"
)

type APIStatus struct {
	statusCode int
	Message    string `json:"message"`
}

func (a APIStatus) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, a.statusCode)
	return nil
}

type APIError struct {
	APIStatus
}

func newAPIStatus(code int, s string) APIStatus {
	return APIStatus{
		statusCode: code,
		Message:    s,
	}
}

func newAPIError(code int, s string) APIError {
	return APIError{
		APIStatus: APIStatus{
			statusCode: code,
			Message:    s,
		},
	}
}

type createdURLEndpointResponse struct {
	URL struct {
		FQDN                  string `json:"fqdn,omitempty"`
		Identifier            string `json:"identifier,omitempty"`
		HumanReadableEndpoint string `json:"human_readable_endpoint,omitempty"`
	} `json:"url,omitempty"`
	APIStatus
}
