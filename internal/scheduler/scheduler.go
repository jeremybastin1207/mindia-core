package scheduler

import (
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/logging"
)

type TaskScheduler struct {
	taskStorage Storer
	logger      *logging.Logger
	listeners   map[string][]TaskFunc
}

func NewTaskScheduler(taskStorage Storer, logger *logging.Logger) TaskScheduler {
	return TaskScheduler{
		taskStorage,
		logger,
		make(map[string][]TaskFunc),
	}
}

func (s *TaskScheduler) RegisterListener(taskName string, taskFunc TaskFunc) {
	s.listeners[taskName] = append(s.listeners[taskName], taskFunc)
}

func (s *TaskScheduler) ProcessTasks() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			s.processTask()
		case <-quit:
			s.logger.Info("task scheduler stopped")
			return
		}
	}
}

func (s *TaskScheduler) processTask() {
	tasks, err := s.taskStorage.GetAll()
	if err != nil {
		return
	}
	for _, t := range tasks {
		if t.Status != Finished {
			listeners := s.listeners[t.Name]
			for _, l := range listeners {
				err := l.Hook(t)
				if err != nil {
					fmt.Println(err)
				}
			}
		} else {
			err := s.taskStorage.Delete(t.Id)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
