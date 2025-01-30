package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Event - сущность "Событие".
//
//nolint:tagliatelle
type Event struct {
	ID           EventUUID     `json:"id"`            // Уникальный идентификатор события.
	Header       string        `json:"header"`        // Заголовок.
	Date         time.Time     `json:"date"`          // Дата и время события.
	Duration     time.Duration `json:"duration"`      // Длительность события.
	Description  string        `json:"description"`   // Описание события.
	NotifyBefore time.Duration `json:"notify_before"` // За сколько времени высылать уведомление.
}

// Notification - временная сущность, не хранится в БД, складывается в очередь для рассыльщика.
type Notification struct {
	ID     uuid.UUID // Уникальный идентификатор события.
	Header string    // Заголовок.
	Date   time.Time // Дата и время события.
}

type EventUUID uuid.UUID

func (e *EventUUID) MarshalJSON() ([]byte, error) {
	event := *e
	return json.Marshal(uuid.UUID(event).String())
}

func (e *EventUUID) UnmarshalJSON(data []byte) error {
	var event uuid.UUID
	if err := json.Unmarshal(data, &event); err != nil {
		return fmt.Errorf("failed to json.Unmarshal: %w", err)
	}

	*e = EventUUID(event)

	return nil
}
