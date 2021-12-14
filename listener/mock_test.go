package listener

import (
	"net"
	"testing"
	"time"

	"github.com/jyk1987/mqtt"
	"github.com/stretchr/testify/require"
)

func TestMockEstablisher(t *testing.T) {
	_, w := net.Pipe()
	err := MockEstablisher("t1", w, new(mqtt.AuthAllow))
	require.NoError(t, err)
	w.Close()
}

func TestNewMockListener(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	require.Equal(t, "t1", mocked.id)
	require.Equal(t, ":1882", mocked.address)
}

func TestNewMockListenerListen(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	require.Equal(t, "t1", mocked.id)
	require.Equal(t, ":1882", mocked.address)

	require.Equal(t, false, mocked.IsListening())
	err := mocked.Listen(nil)
	require.NoError(t, err)
	require.Equal(t, true, mocked.IsListening())
}
func TestNewMockListenerListenFailure(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	mocked.ErrListen = true
	err := mocked.Listen(nil)
	require.Error(t, err)
}

func TestMockListenerServe(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	require.Equal(t, false, mocked.IsServing())

	o := make(chan bool)
	go func(o chan bool) {
		mocked.Serve(MockEstablisher)
		o <- true
	}(o)

	time.Sleep(time.Millisecond) // easy non-channel wait for start of serving
	require.Equal(t, true, mocked.IsServing())

	var closed bool
	mocked.Close(func(id string) {
		closed = true
	})
	require.Equal(t, true, closed)
	<-o

	mocked.Listen(nil)
}

func TestMockListenerClose(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	var closed bool
	mocked.Close(func(id string) {
		closed = true
	})
	require.Equal(t, true, closed)
}

func TestNewMockListenerIsListening(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	require.Equal(t, false, mocked.IsListening())
}

func TestNewMockListenerIsServing(t *testing.T) {
	mocked := NewMockListener("t1", ":1882")
	require.Equal(t, false, mocked.IsServing())
}
