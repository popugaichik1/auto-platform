package ws

import (
	"errors"
	"testing"

	"github.com/google/uuid"
)

// fakeConn — fake-реализация Conn, без реального сокета.
type fakeConn struct {
	written []ServerFrame
	failNext bool
	closed  bool
}

func (c *fakeConn) WriteJSON(v any) error {
	if c.failNext {
		return errors.New("write failed")
	}
	c.written = append(c.written, v.(ServerFrame))
	return nil
}

func (c *fakeConn) Close() error {
	c.closed = true
	return nil
}

func TestHub_SendTo_NotConnected(t *testing.T) {
	hub := NewHub()

	if hub.SendTo(uuid.New(), ServerFrame{Type: "message"}) {
		t.Fatalf("SendTo() = true, want false for a user with no connection")
	}
}

func TestHub_RegisterThenSendTo(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	conn := &fakeConn{}

	hub.Register(userID, conn)

	frame := ServerFrame{Type: "message", Payload: "hello"}
	if !hub.SendTo(userID, frame) {
		t.Fatalf("SendTo() = false, want true for a registered connection")
	}
	if len(conn.written) != 1 || conn.written[0] != frame {
		t.Fatalf("conn.written = %+v, want [%+v]", conn.written, frame)
	}
}

func TestHub_SendTo_WriteFailure(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	conn := &fakeConn{failNext: true}

	hub.Register(userID, conn)

	if hub.SendTo(userID, ServerFrame{Type: "message"}) {
		t.Fatalf("SendTo() = true, want false when the underlying write fails")
	}
}

func TestHub_Unregister_RemovesConnection(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	conn := &fakeConn{}

	hub.Register(userID, conn)
	hub.Unregister(userID, conn)

	if hub.SendTo(userID, ServerFrame{Type: "message"}) {
		t.Fatalf("SendTo() = true after Unregister(), want false")
	}
}

func TestHub_Unregister_StaleConnectionDoesNotEvictNewOne(t *testing.T) {
	hub := NewHub()
	userID := uuid.New()
	oldConn := &fakeConn{}
	newConn := &fakeConn{}

	hub.Register(userID, oldConn)
	hub.Register(userID, newConn) // пользователь переподключился новым соединением

	// "протухший" Unregister от старого соединения не должен затереть новое
	hub.Unregister(userID, oldConn)

	frame := ServerFrame{Type: "message", Payload: "hi"}
	if !hub.SendTo(userID, frame) {
		t.Fatalf("SendTo() = false, want true — new connection should still be registered")
	}
	if len(newConn.written) != 1 {
		t.Fatalf("expected message delivered to the new connection, got %d writes", len(newConn.written))
	}
	if len(oldConn.written) != 0 {
		t.Fatalf("expected no messages delivered to the stale old connection, got %d writes", len(oldConn.written))
	}
}
