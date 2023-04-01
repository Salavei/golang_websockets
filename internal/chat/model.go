package chat

import "time"

type Message struct {
	ID     string        `json:"id" bson:"_id,omitempty"`
	Text   string        `json:"message" bson:"message"`
	UserID string        `json:"user_id" bson:"user_id"`
	Date   time.Duration `json:"date" bson:"date"`
}
