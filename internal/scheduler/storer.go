package scheduler

type Storer interface {
	EnqueueTask(task *Task) error
	DequeueTask() (*Task, error)
}
