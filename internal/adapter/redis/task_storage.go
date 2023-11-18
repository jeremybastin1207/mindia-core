package redis

import (
	"encoding/json"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/scheduler"
	"github.com/nitishm/go-rejson"
	"github.com/rs/zerolog/log"
)

const tasks_key_prefix = "internal:task"
const task_index_name = "task_index"

type TaskStorage struct {
	rejsonHandler    *rejson.Handler
	redisearchClient *redisearch.Client
}

func NewTaskStorage(redisPool *redigo.Pool) *TaskStorage {
	s := &TaskStorage{
		rejsonHandler:    rejson.NewReJSONHandler(),
		redisearchClient: redisearch.NewClientFromPool(redisPool, task_index_name),
	}
	s.rejsonHandler.SetRedigoClient(redisPool.Get())
	s.init()
	return s
}

func (s *TaskStorage) init() {
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("$.id")).
		AddField(redisearch.NewTextField("$.name")).
		AddField(redisearch.NewTextField("$.status"))

	def := redisearch.NewIndexDefinition()
	def.IndexOn = "JSON"
	def.Prefix = []string{fmt.Sprintf("%s:", tasks_key_prefix)}

	if err := s.redisearchClient.DropIndex(false); err != nil {
		log.Warn().Err(err)
	}
	if err := s.redisearchClient.CreateIndexWithIndexDefinition(schema, def); err != nil {
		mindiaerr.ExitErrorf(err.Error())
	}
}

func parseJSON(jsonContent interface{}) (*scheduler.Task, error) {
	var task scheduler.Task
	err := json.Unmarshal(jsonContent.([]byte), &task)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (s *TaskStorage) GetAll() ([]scheduler.Task, error) {
	tasks := []scheduler.Task{}
	results, _, err := s.redisearchClient.Search(redisearch.NewQuery("*"))
	if err != nil {
		return nil, err
	}
	for _, r := range results {
		m, err := s.get(r.Id)
		if err != nil {
			continue
		}
		tasks = append(tasks, *m)
	}
	return tasks, nil
}

func (s *TaskStorage) get(key string) (*scheduler.Task, error) {
	res, err := s.rejsonHandler.JSONGet(key, ".")
	if err != nil {
		return nil, err
	}
	return parseJSON(res)
}

func (s *TaskStorage) Get(id uuid.UUID) (*scheduler.Task, error) {
	res, err := s.rejsonHandler.JSONGet(fmt.Sprintf("%s:%s", tasks_key_prefix, id.String()), ".")
	if err != nil {
		return nil, err
	}
	return parseJSON(res)
}

func (s *TaskStorage) GetEnqueued() (*scheduler.Task, error) {
	res, err := s.rejsonHandler.JSONGet(tasks_key_prefix, "status:enqueud")
	if err != nil {
		return nil, err
	}
	var task scheduler.Task
	err = json.Unmarshal(res.([]byte), &task)
	if err != nil {
		return nil, err
	}
	return parseJSON(res)
}

func (s *TaskStorage) GetEnqueuedByName(name string) (*scheduler.Task, error) {
	res, err := s.rejsonHandler.JSONGet(tasks_key_prefix, "status:enqueud")
	if err != nil {
		return nil, err
	}
	var task scheduler.Task
	err = json.Unmarshal(res.([]byte), &task)
	if err != nil {
		return nil, err
	}
	return parseJSON(res)
}

func (s *TaskStorage) Save(task scheduler.Task) error {
	_, err := s.rejsonHandler.JSONSet(fmt.Sprintf("%s:%s", tasks_key_prefix, task.Id.String()), ".", task)
	return err
}

func (s *TaskStorage) Delete(id uuid.UUID) error {
	_, err := s.rejsonHandler.JSONDel(fmt.Sprintf("%s:%s", tasks_key_prefix, id.String()), ".")
	return err
}
