package WebSocket

import "github.com/gorilla/websocket"

type Notification struct {
	UserID  int64
	Message string
}

type Client struct {
	UserID     int64
	Connection *websocket.Conn
}

var (
	clients   = make(map[int64]*Client) // Мапа подключенных клиентов
	Broadcast = make(chan Notification) // Канал для рассылки уведомлений
)
