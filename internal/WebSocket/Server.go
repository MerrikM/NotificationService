package WebSocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnections обрабатывает подключение клиента через WebSocket
func HandleConnections(writer http.ResponseWriter, request *http.Request) error {

	// Аутентификация (пример через query-параметр) для теста, в продакшн - header
	userIDStr := request.URL.Query().Get("user-ID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID <= 0 {
		return fmt.Errorf("неверный user_id")
	}

	// Обновление до WebSocket-соединения
	connection, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return fmt.Errorf("ошибка соединения: %v", err)
	}
	defer connection.Close()

	client := &Client{
		UserID:     userID,
		Connection: connection,
	}
	clients[userID] = client
	defer func() {
		log.Printf("Клиент с userID=%d отключен", userID)
		RemoveClient(userID)
	}()

	log.Printf("Клиент подключен: user-ID=%d", userID)

	for {
		_, _, err := connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Ошибка соединения: user_id=%d, error: %v", userID, err)
			}
			break
		}
	}
	return nil
}

func StartWebSocketServer(port string) error {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		if err := HandleConnections(w, r); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
		}
	})

	// Запуск канала для обработки уведомлений
	go HandleBroadcasts()

	log.Printf("WebSocket сервер запущен на порту %s", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("ошибка запуска сервера: %v", err)
	}

	return nil
}
