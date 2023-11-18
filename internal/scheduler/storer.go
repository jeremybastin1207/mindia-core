package scheduler

import "github.com/google/uuid"

type Storer interface {
	GetAll() ([]Task, error)
	Get(id uuid.UUID) (*Task, error)
	GetEnqueued() (*Task, error)
	GetEnqueuedByName(name string) (*Task, error)
	Save(task Task) error
	Delete(id uuid.UUID) error
}
