package tui

import (
	"net"
	"net/http"
	"time"
)

type IncomingHTTPRequest struct {
	Method    string    `json:"method,omitempty"`
	ID        string    `json:"id,omitempty"`
	IPAddress net.IP    `json:"ip_address,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`

	Headers http.Header `json:"headers,omitempty"`
}

type IncomingHTTPRequests []IncomingHTTPRequest

type incomingHTTPRequestMsg struct {
	data IncomingHTTPRequests
}

type item struct {
	title, desc, ip string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }
