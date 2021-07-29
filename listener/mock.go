package listener

import (
	"fmt"

	"net"
	"sync"

	"github.com/snple/mqtt"
	"github.com/snple/mqtt/system"
)

// MockCloser is a function signature which can be used in testing.
func MockCloser(id string) {}

// MockEstablisher is a function signature which can be used in testing.
func MockEstablisher(id string, c net.Conn, auth mqtt.Auth) error {
	return nil
}

// MockListener is a mock listener for establishing client connections.
type MockListener struct {
	sync.RWMutex
	id        string    // the id of the listener.
	address   string    // the network address the listener binds to.
	Listening bool      // indiciate the listener is listening.
	Serving   bool      // indicate the listener is serving.
	done      chan bool // indicate the listener is done.
	ErrListen bool      // throw an error on listen.
}

// NewMockListener returns a new instance of MockListener
func NewMockListener(id, address string) *MockListener {
	return &MockListener{
		id:      id,
		address: address,
		done:    make(chan bool),
	}
}

// Serve serves the mock listener.
func (l *MockListener) Serve(establisher mqtt.EstablishFunc) error {
	l.Lock()
	l.Serving = true
	l.Unlock()
	for {
		select {
		case <-l.done:
			return nil
		}
	}
}

// SetConfig sets the configuration values of the mock listener.
func (l *MockListener) Listen(s *system.Info) error {
	if l.ErrListen {
		return fmt.Errorf("listen failure")
	}

	l.Lock()
	l.Listening = true
	l.Unlock()
	return nil
}

// ID returns the id of the mock listener.
func (l *MockListener) ID() string {
	l.RLock()
	id := l.id
	l.RUnlock()
	return id
}

// Close closes the mock listener.
func (l *MockListener) Close(closer mqtt.CloseFunc) {
	l.Lock()
	defer l.Unlock()
	l.Serving = false
	closer(l.id)
	close(l.done)
}

// IsServing indicates whether the mock listener is serving.
func (l *MockListener) IsServing() bool {
	l.Lock()
	defer l.Unlock()
	return l.Serving
}

// IsServing indicates whether the mock listener is listening.
func (l *MockListener) IsListening() bool {
	l.Lock()
	defer l.Unlock()
	return l.Listening
}
