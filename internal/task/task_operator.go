package task

import "github.com/jeremybastin1207/mindia-core/internal/scheduler"

type TaskOperator struct {
	storer scheduler.Storer
}

func NewTaskOperator(taskStorage scheduler.Storer) TaskOperator {
	return TaskOperator{
		storer: taskStorage,
	}
}
