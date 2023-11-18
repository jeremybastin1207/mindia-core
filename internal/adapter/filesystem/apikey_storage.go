package filesystem

import (
	"errors"
	"os"

	"github.com/jeremybastin1207/mindia-core/internal/apikey"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"gopkg.in/yaml.v2"
)

type ApiKeyStorage struct {
	filename string
	apikeys  apikey.ApiKeyMap
}

func NewApiKeyStorage() *ApiKeyStorage {
	return &ApiKeyStorage{
		filename: "apikeys.yml",
		apikeys:  apikey.ApiKeyMap{},
	}
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
	s.load()
	if val, ok := s.apikeys[key]; ok {
		return &val, nil
	}
	return nil, &mindiaerr.Error{ErrCode: mindiaerr.ErrCodeApiKeyNotFound}
}

func (s *ApiKeyStorage) GetAll() (apikey.ApiKeyMap, error) {
	s.load()
	return s.apikeys, nil
}

func (s *ApiKeyStorage) Save(apikey apikey.ApiKey) error {
	s.load()
	s.apikeys[apikey.Name] = apikey
	return s.save()
}

func (s *ApiKeyStorage) Delete(key string) error {
	s.load()
	delete(s.apikeys, key)
	return s.save()
}

func (s *ApiKeyStorage) load() {
	if s.apikeys != nil {
		return
	}
	_, err := os.Stat(s.filename)
	if errors.Is(err, os.ErrNotExist) {
		err = s.save()
		if err != nil {
			panic(err)
		}
	}

	body, err := os.ReadFile(s.filename)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(body, &s.apikeys)
	if err != nil {
		panic(err)
	}
}

func (s *ApiKeyStorage) save() error {
	yamlData, err := yaml.Marshal(s.apikeys)
	if err != nil {
		return err
	}
	return os.WriteFile(s.filename, yamlData, 0644)
}
