package drive

import "context"

type Interface interface {
	Push(ctx context.Context, message string) error
	GetMessage(ctx context.Context) (string, error)
	CommitMessage(ctx context.Context) error
}
