package mqtt

import (
	"net"
	"sync"

	"github.com/jyk1987/mqtt/system"
)

// EstablishFunc is a callback function for establishing new clients.
type EstablishFunc func(id string, c net.Conn, auth Auth) error

// CloseFunc is a callback function for closing all listener clients.
type CloseFunc func(id string)

// Listener is an interface for network listeners. A network listener listens
// for incoming client connections and adds them to the server.
type Listener interface {
	Listen(s *system.Info) error // open the network address.
	Serve(EstablishFunc) error   // starting actively listening for new connections.
	ID() string                  // return the id of the listener.
	Auth() Auth
	Close(CloseFunc) // stop and close the listener.
}

// Listeners contains the network listeners for the broker.
type Listeners struct {
	sync.RWMutex
	wg       sync.WaitGroup      // a waitgroup that waits for all listeners to finish.
	internal map[string]Listener // a map of active listeners.
	system   *system.Info        // pointers to system info.
}

// New returns a new instance of Listeners.
func NewListeners(s *system.Info) *Listeners {
	return &Listeners{
		internal: map[string]Listener{},
		system:   s,
	}
}

// Add adds a new listener to the listeners map, keyed on id.
func (l *Listeners) Add(val Listener) {
	l.Lock()
	l.internal[val.ID()] = val
	l.Unlock()
}

// Get returns the value of a listener if it exists.
func (l *Listeners) Get(id string) (Listener, bool) {
	l.RLock()
	val, ok := l.internal[id]
	l.RUnlock()
	return val, ok
}

// Len returns the length of the listeners map.
func (l *Listeners) Len() int {
	l.RLock()
	val := len(l.internal)
	l.RUnlock()
	return val
}

// Delete removes a listener from the internal map.
func (l *Listeners) Delete(id string) {
	l.Lock()
	delete(l.internal, id)
	l.Unlock()
}

// Serve starts a listener serving from the internal map.
func (l *Listeners) Serve(id string, establisher EstablishFunc) {
	l.RLock()
	listener := l.internal[id]
	l.RUnlock()

	go func(e EstablishFunc) {
		defer l.wg.Done()
		l.wg.Add(1)
		err := listener.Serve(e)
		if err != nil {

		}
	}(establisher)
}

// ServeAll starts all listeners serving from the internal map.
func (l *Listeners) ServeAll(establisher EstablishFunc) {
	l.RLock()
	i := 0
	ids := make([]string, len(l.internal))
	for id := range l.internal {
		ids[i] = id
		i++
	}
	l.RUnlock()

	for _, id := range ids {
		l.Serve(id, establisher)
	}
}

// Close stops a listener from the internal map.
func (l *Listeners) Close(id string, closer CloseFunc) {
	l.RLock()
	listener := l.internal[id]
	l.RUnlock()
	listener.Close(closer)
}

// CloseAll iterates and closes all registered listeners.
func (l *Listeners) CloseAll(closer CloseFunc) {
	l.RLock()
	i := 0
	ids := make([]string, len(l.internal))
	for id := range l.internal {
		ids[i] = id
		i++
	}
	l.RUnlock()

	for _, id := range ids {
		l.Close(id, closer)
	}
	l.wg.Wait()
}
