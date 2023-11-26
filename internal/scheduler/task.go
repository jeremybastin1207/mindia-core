package scheduler

import (
	"time"

	"github.com/google/uuid"
)

type TaskStatus string

const (
	Enqueued   TaskStatus = "enqueued"
	Processing TaskStatus = "processing"
	Finished   TaskStatus = "finished"
	Canceled   TaskStatus = "canceled"
)

type Task struct {
	Id         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	Status     TaskStatus  `json:"status"`
	Details    interface{} `json:"details,omitempty"`
	Error      string      `json:"error,omitempty"`
	EnqueuedAt time.Time   `json:"enqueued_at,omitempty"`
	StartedAt  time.Time   `json:"started_at,omitempty"`
	FinishedAt time.Time   `json:"finished_at,omitempty"`
}

func NewTask(name string, details interface{}) Task {
	return Task{
		Id:      uuid.New(),
		Name:    name,
		Details: details,
	}
}

type TaskFunc interface {
	Execute(task *Task) (*Task, error)
}
