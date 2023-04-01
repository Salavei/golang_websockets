package chat

import (
	"context"
)

type Storage interface {
	SendMessage(ctx context.Context, user string, msg Message) (string, error)
	ShowMessage(ctx context.Context) (msgs []Message, err error)
}
