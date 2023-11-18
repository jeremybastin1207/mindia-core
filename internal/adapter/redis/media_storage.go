package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/RediSearch/redisearch-go/redisearch"
	redigo "github.com/gomodule/redigo/redis"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	rejson "github.com/nitishm/go-rejson/v4"
	"github.com/rs/zerolog/log"
)

const medias_key = "media"

type mediaWithTimestamp struct {
	media.Media
	Timestamp int64 `json:"timestamp"`
}

type MediaStorage struct {
	rejsonHandler    *rejson.Handler
	redisearchClient *redisearch.Client
}

func NewMediaStorage(redisPool *redigo.Pool) *MediaStorage {
	rejsonHandler := rejson.NewReJSONHandler()
	rejsonHandler.SetRedigoClient(redisPool.Get())
	redisearchClient := redisearch.NewClientFromPool(redisPool, "media_idx")

	s := MediaStorage{
		rejsonHandler:    rejsonHandler,
		redisearchClient: redisearchClient,
	}
	s.init()
	return &s
}

func (s *MediaStorage) init() {
	schema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("$.path")).
		AddField(redisearch.NewTextField("$.content_type")).
		AddField(redisearch.NewNumericField("$.content_length")).
		AddField(redisearch.NewNumericFieldOptions("$.timestamp", redisearch.NumericFieldOptions{Sortable: true}))

	def := redisearch.NewIndexDefinition()
	def.IndexOn = "JSON"
	def.Prefix = []string{fmt.Sprintf("%s:", medias_key)}

	if err := s.redisearchClient.DropIndex(false); err != nil {
		log.Warn().Err(err)
	}
	if err := s.redisearchClient.CreateIndexWithIndexDefinition(schema, def); err != nil {
		mindiaerr.ExitErrorf(err.Error())
	}
}

func (s *MediaStorage) Get(path media.Path) (*media.Media, error) {
	id := fmt.Sprintf("%s:%s", medias_key, path.ToString())
	return s.get(id)
}

func (s *MediaStorage) get(id string) (*media.Media, error) {
	res, err := s.rejsonHandler.JSONGet(id, ".")
	if err != nil {
		return nil, err
	}
	if res == nil {
		return nil, mindiaerr.New(mindiaerr.ErrCodeMediaNotFound)
	}
	var media mediaWithTimestamp
	err = json.Unmarshal(res.([]byte), &media)
	if err != nil {
		return nil, &mindiaerr.Error{ErrCode: mindiaerr.ErrCodeInternal, Msg: err}
	}
	return &media.Media, nil
}

func (s *MediaStorage) GetMultiple(path media.Path, offset int, limit int, sortBy string, asc bool) ([]media.Media, error) {
	var medias = []media.Media{}

	switch sortBy {
	case "created_at":
		sortBy = "$.timestamp"
	case "content_length":
		sortBy = "$.content_length"
	default:
		return medias, errors.New("sortby not supported")
	}

	results, _, err := s.redisearchClient.Search(
		redisearch.
			NewQuery(fmt.Sprintf("^%s*", path.ToString())).
			SetInFields("$.path").
			SetSortBy(sortBy, asc).
			Limit(offset, limit))
	if err != nil {
		return nil, err
	}

	for _, r := range results {
		m, err := s.get(r.Id)
		if err != nil {
			continue
		}
		medias = append(medias, *m)
	}
	return medias, nil
}

func (s *MediaStorage) Save(m *media.Media) error {
	mi := mediaWithTimestamp{
		Media:     *m,
		Timestamp: m.CreatedAt.Unix(),
	}
	_, err := s.rejsonHandler.JSONSet(fmt.Sprintf("%s:%s", medias_key, m.Path.ToString()), ".", mi)
	return err
}

func (s *MediaStorage) Delete(path media.Path) error {
	_, err := s.rejsonHandler.JSONDel(fmt.Sprintf("%s:%s", medias_key, path.ToString()), ".")
	return err
}
