package redis

import (
	"encoding/json"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
	"github.com/nitishm/go-rejson"
)

const named_transformations_key = "internal:configuration:named_transformations"

type NamedTransformationStorage struct {
	rejsonHandler *rejson.Handler
}

func NewNamedTransformationStorage(redisPool *redigo.Pool) *NamedTransformationStorage {
	rejsonHandler := rejson.NewReJSONHandler()
	rejsonHandler.SetRedigoClient(redisPool.Get())

	s := NamedTransformationStorage{
		rejsonHandler: rejsonHandler,
	}
	s.init()
	return &s
}

func (s *NamedTransformationStorage) init() error {
	res, err := s.rejsonHandler.JSONGet(named_transformations_key, ".")
	if err != nil {
		return err
	}
	if res == nil {
		_, err := s.rejsonHandler.JSONSet(named_transformations_key, ".", transform.NamedTransformationMap{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *NamedTransformationStorage) GetAll() (transform.NamedTransformationMap, error) {
	res, err := s.rejsonHandler.JSONGet(named_transformations_key, ".")
	if err != nil {
		return nil, err
	}
	if res == nil {
		return transform.NamedTransformationMap{}, nil
	}
	var namedTransformations transform.NamedTransformationMap
	err = json.Unmarshal(res.([]byte), &namedTransformations)
	if err != nil {
		return nil, err
	}
	return namedTransformations, nil
}

func (s *NamedTransformationStorage) Get(name string) (*transform.NamedTransformation, error) {
	res, err := s.rejsonHandler.JSONGet(named_transformations_key, fmt.Sprintf(".:%s", name))
	if err != nil {
		return nil, &mindiaerr.Error{
			ErrCode: mindiaerr.ErrCodeNamedTransformationNotFound,
			Msg:     fmt.Errorf("name: %v", name),
		}
	}
	var namedTransformation transform.NamedTransformation
	err = json.Unmarshal(res.([]byte), &namedTransformation)
	if err != nil {
		return nil, err
	}
	return &namedTransformation, nil
}

func (s *NamedTransformationStorage) Save(namedTransformation transform.NamedTransformation) error {
	s.init()
	_, err := s.rejsonHandler.JSONSet(named_transformations_key, fmt.Sprintf(".%s", namedTransformation.Name), namedTransformation)
	return err
}

func (s *NamedTransformationStorage) Delete(name string) error {
	_, err := s.rejsonHandler.JSONDel(named_transformations_key, fmt.Sprintf(".%s", name))
	return err
}

func (s *NamedTransformationStorage) DeleteAll() error {
	_, err := s.rejsonHandler.JSONDel(named_transformations_key, ".")
	return err
}
