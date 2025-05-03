package config

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

func RemoveClient(userID int64) {
	if client, exists := clients[userID]; exists {
		log.Printf("Закрытие соединения для userID=%d", userID)
		client.Connection.Close()
		delete(clients, userID)
	}
}

// HandleBroadcasts рассылает уведомления клиентам
func HandleBroadcasts() {
	for notification := range Broadcast {
		log.Printf("Получение уведомления для userID=%d: %+v", notification.UserID, notification)

		// Проверяем, существует ли клиент в мапе
		if client, exists := clients[notification.UserID]; exists {
			log.Printf("Отправка уведомления для userID=%d через WebSocket", notification.UserID)

			// Отправка уведомления через WebSocket
			err := client.Connection.WriteJSON(notification)
			if err != nil {
				log.Printf("Ошибка отправки уведомления: user_id=%d, error: %v", notification.UserID, err)

				// Закрытие соединения и удаление клиента из мапы
				RemoveClient(notification.UserID)
			} else {
				log.Printf("Уведомление отправлено пользователю с userID=%d", notification.UserID)
			}
		} else {
			log.Printf("Клиент с userID=%d не найден в мапе", notification.UserID)
		}
	}
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
