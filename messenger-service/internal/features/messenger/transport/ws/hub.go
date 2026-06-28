package ws

import (
	"sync"

	"github.com/google/uuid"
)

// Conn — то немногое, что Hub'у нужно от *websocket.Conn. Выделено в
// интерфейс, чтобы hub_test.go мог подставить fake без реального сокета.
type Conn interface {
	WriteJSON(v any) error
	Close() error
}

// Hub — реестр "какой пользователь подключён к этой конкретной реплике
// прямо сейчас". Один пользователь = одно соединение (последнее подключение
// побеждает; переподключение с нового устройства/таба вытесняет старое).
// Используется конкурентно: каждое WS-соединение обслуживается своей
// goroutine, поэтому доступ защищён мьютексом.
type Hub struct {
	mu    sync.RWMutex
	conns map[uuid.UUID]Conn
}

func NewHub() *Hub {
	return &Hub{conns: make(map[uuid.UUID]Conn)}
}

func (h *Hub) Register(userID uuid.UUID, conn Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.conns[userID] = conn
}

// Unregister удаляет соединение пользователя, но только если в реестре
// до сих пор лежит именно то соединение, что передано сюда — если
// пользователь успел переподключиться (новый Register уже произошёл),
// "протухший" Unregister от старого соединения не должен затирать новое.
func (h *Hub) Unregister(userID uuid.UUID, conn Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if existing, ok := h.conns[userID]; ok && existing == conn {
		delete(h.conns, userID)
	}
}

// SendTo пытается доставить сообщение пользователю, если он подключён
// именно к этой реплике. Возвращает false без ошибки, если пользователь
// сейчас не подключён здесь — это штатная ситуация (он на другой реплике
// или офлайн), не повод логировать ошибку выше.
func (h *Hub) SendTo(userID uuid.UUID, frame ServerFrame) bool {
	h.mu.RLock()
	conn, ok := h.conns[userID]
	h.mu.RUnlock()

	if !ok {
		return false
	}

	return conn.WriteJSON(frame) == nil
}
