package model

import (
	"time"

	"github.com/google/uuid"
)

// Event - сущность "Событие".
type Event struct {
	ID           EventUUID     // Уникальный идентификатор события.
	Header       string        // Заголовок.
	Date         time.Time     // Дата и время события.
	Duration     time.Duration // Длительность события.
	Description  string        // Описание события.
	UserID       UserUUID      // ID пользователя, владельца события.
	NotifyBefore time.Duration // За сколько времени высылать уведомление.
}

// Notification - временная сущность, не хранится в БД, складывается в очередь для рассыльщика.
type Notification struct {
	ID     uuid.UUID // Уникальный идентификатор события.
	Header string    // Заголовок.
	Date   time.Time // Дата и время события.
	UserID string    // ID пользователя, которому отправлять.
}

type (
	EventUUID uuid.UUID
	UserUUID  uuid.UUID
)
