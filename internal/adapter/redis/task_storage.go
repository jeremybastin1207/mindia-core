package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jeremybastin1207/mindia-core/internal/scheduler"
)

const tasksQueueKey = "internal:taskQueue"

type TaskStorage struct {
	redisPool *redis.Pool
}

func NewTaskStorage(redisPool *redigo.Pool) *TaskStorage {
	s := &TaskStorage{
		redisPool,
	}
	return s
}

func (s *TaskStorage) EnqueueTask(task *scheduler.Task) error {
	conn := s.redisPool.Get()
	defer conn.Close()

	taskJSON, err := json.Marshal(task)
	if err != nil {
		return err
	}

	_, err = conn.Do("RPUSH", tasksQueueKey, taskJSON)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskStorage) DequeueTask() (*scheduler.Task, error) {
	conn := s.redisPool.Get()
	defer conn.Close()

	taskJSON, err := conn.Do("LPOP", tasksQueueKey)
	if err != nil {
		return nil, err
	}
	if taskJSON == nil {
		return nil, nil
	}

	var task scheduler.Task
	err = json.Unmarshal(taskJSON.([]byte), &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}
