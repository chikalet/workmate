package tasks

import (
	"context"
	"time"
)

type Status string

const (
	StatusCreated    Status = "создана"
	StatusInProgress Status = "в процессе"
	StatusCompleted  Status = "завершена"
)

type Task struct {
	ID              int       `json:"id"`
	Status          Status    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	DurationSeconds int64     `json:"duration"`

	startedAt time.Time
	cancel    context.CancelFunc
}
