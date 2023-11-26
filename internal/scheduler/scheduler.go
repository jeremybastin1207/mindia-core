package scheduler

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/logging"
)

const tickDuration = 10 * time.Second

type TaskScheduler struct {
	taskStorage Storer
	logger      logging.Logger
	listeners   map[string]TaskFunc
}

func NewTaskScheduler(taskStorage Storer, logger logging.Logger) TaskScheduler {
	return TaskScheduler{
		taskStorage: taskStorage,
		logger:      logger,
		listeners:   make(map[string]TaskFunc),
	}
}

func (s *TaskScheduler) RegisterListener(taskName string, taskFunc TaskFunc) {
	if taskName == "" || taskFunc == nil {
		s.logger.Error("invalid taskName or taskFunc")
		return
	}
	s.listeners[taskName] = taskFunc
}

func (s *TaskScheduler) ProcessTasks() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	ticker := time.NewTicker(tickDuration)

	for {
		select {
		case <-ticker.C:
			s.processTask()
		case <-quit:
			s.logger.Info("task scheduler stopped")
			ticker.Stop()
			return
		}
	}
}

func (s *TaskScheduler) processTask() {
	t, err := s.taskStorage.DequeueTask()
	if err != nil {
		return
	}
	if t == nil {
		return
	}
	l := s.listeners[t.Name]
	go func(t *Task) {
		t2, err := l.Execute(t)
		if err != nil {
			s.logger.Error(fmt.Sprintf("Failed to execute task: %v", err))
			return
		}
		if t2 != nil {
			err := s.taskStorage.EnqueueTask(t2)
			if err != nil {
				s.logger.Error(fmt.Sprintf("Failed to enqueue task: %v", err))
				return
			}
		}
	}(t)
}
