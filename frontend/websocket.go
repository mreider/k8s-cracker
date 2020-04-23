package main

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type WebsocketConnections struct {
	connections map[*websocket.Conn]struct{}
	mu          sync.Mutex
}

func NewWebsocketConnections() *WebsocketConnections {
	return &WebsocketConnections{
		connections: make(map[*websocket.Conn]struct{}),
	}
}

func (s *WebsocketConnections) addConnection(conn *websocket.Conn) {
	s.mu.Lock()
	s.connections[conn] = struct{}{}
	s.mu.Unlock()
}

func (s *WebsocketConnections) closeConnection(conn *websocket.Conn) {
	s.mu.Lock()
	delete(s.connections, conn)
	s.mu.Unlock()

	_ = conn.Close(websocket.StatusNormalClosure, "")
}

func (s *WebsocketConnections) Broadcast(ctx context.Context, message []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for conn := range s.connections {
		conn := conn
		go func() {
			err := conn.Write(ctx, websocket.MessageText, message)
			if err != nil {
				log.Error().Err(err).Send()
			}
		}()
	}
}
