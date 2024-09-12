package forward

import (
	"sync"

	"github.com/charmbracelet/ssh"
	gossh "golang.org/x/crypto/ssh"
)

type ConnectionInfo struct {
	Port int
	Host string

	// do we even need this really?
	// since we already have the host and port
	// RemoteAddr net.Addr
}

type Forwarder struct {
	// known connections to the ssh server.
	// the key here is the host of the connecting user
	connections map[string]ConnectionInfo
	mutex       sync.RWMutex
}

func New() *Forwarder {
	return &Forwarder{
		mutex:       sync.RWMutex{},
		connections: make(map[string]ConnectionInfo),
	}
}

func (f *Forwarder) AddConnection(key string, data ConnectionInfo) {
	f.mutex.Lock()
	defer f.mutex.Unlock()

	f.connections[key] = data
}

func (f *Forwarder) HandleSSHRequest(ctx ssh.Context,
	srv *ssh.Server,
	req *gossh.Request,
) (bool, []byte) {
	return true, nil
}
