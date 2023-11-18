package redis

import (
	"encoding/json"
	"errors"
	"fmt"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/jeremybastin1207/mindia-core/internal/apikey"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/nitishm/go-rejson"
)

const api_keys_key = "internal:configuration:api_keys"

type ApiKeyStorage struct {
	rejsonHandler *rejson.Handler
}

func NewApiKeyStorage(redisPool *redigo.Pool) *ApiKeyStorage {
	rejsonHandler := rejson.NewReJSONHandler()
	rejsonHandler.SetRedigoClient(redisPool.Get())

	s := ApiKeyStorage{
		rejsonHandler: rejsonHandler,
	}
	s.init()
	return &s
}

func (s *ApiKeyStorage) init() error {
	res, err := s.rejsonHandler.JSONGet(api_keys_key, ".")
	if err != nil {
		return err
	}
	if res == nil {
		_, err := s.rejsonHandler.JSONSet(api_keys_key, ".", apikey.ApiKeyMap{})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *ApiKeyStorage) GetByName(name string) (*apikey.ApiKey, error) {
	apikeys, err := s.GetAll()
	if err != nil {
		return nil, err
	}
	if val, ok := apikeys[name]; ok {
		return &val, nil
	}
	return nil, errors.New("apikey not found")
}

func (s *ApiKeyStorage) GetByKey(key string) (*apikey.ApiKey, error) {
	res, err := s.rejsonHandler.JSONGet(api_keys_key, fmt.Sprintf(".:%s", key))
	if err != nil {
		return nil, &mindiaerr.Error{
			ErrCode: mindiaerr.ErrCodeApiKeyNotFound,
			Msg:     fmt.Errorf("key: %v", key),
		}
	}
	var apikey apikey.ApiKey
	err = json.Unmarshal(res.([]byte), &apikey)
	if err != nil {
		return nil, err
	}
	return &apikey, nil
}

func (s *ApiKeyStorage) GetAll() (apikey.ApiKeyMap, error) {
	res, err := s.rejsonHandler.JSONGet(api_keys_key, ".")
	if err != nil {
		return nil, err
	}
	if res == nil {
		return apikey.ApiKeyMap{}, nil
	}
	var apikeys apikey.ApiKeyMap
	err = json.Unmarshal(res.([]byte), &apikeys)
	if err != nil {
		return nil, err
	}
	return apikeys, nil
}

func (s *ApiKeyStorage) Save(apikey apikey.ApiKey) error {
	s.init()
	_, err := s.rejsonHandler.JSONSet(api_keys_key, fmt.Sprintf(".%s", apikey.Key), apikey)
	return err
}

func (s *ApiKeyStorage) Delete(key string) error {
	_, err := s.rejsonHandler.JSONDel(api_keys_key, fmt.Sprintf(".%s", key))
	return err
}
