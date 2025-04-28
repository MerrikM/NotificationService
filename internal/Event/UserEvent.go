package Event

import "notificationservice/internal/DTO"

type UserEvent struct {
	EventType string
	User      *DTO.User
}
