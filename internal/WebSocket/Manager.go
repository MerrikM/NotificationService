package WebSocket

import "log"

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
