package tui

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/adelowo/sdump"
)

type DumpURLMsg struct {
	URL        string `json:"url,omitempty"`
	SSEChannel string `json:"sse_channel,omitempty"`
}

type ItemMsg struct {
	item item
}

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
	ID      string
	Request sdump.RequestDefinition
}

func (i item) Title() string       { return fmt.Sprintf("%s    %s", i.ID, i.Request.IPAddress) }
func (i item) Description() string { return "Here is my long ass description" }
func (i item) FilterValue() string { return i.ID }
