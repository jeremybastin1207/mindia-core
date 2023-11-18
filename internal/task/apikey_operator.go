package task

import (
	"time"

	"github.com/jeremybastin1207/mindia-core/internal/apikey"
)

type ApiKeyOperator struct {
	apikeyStorage apikey.Storer
}

func NewApiKeyOperator(apikeyStorage apikey.Storer) ApiKeyOperator {
	return ApiKeyOperator{
		apikeyStorage,
	}
}

func (t *ApiKeyOperator) GetAll() ([]apikey.ApiKey, error) {
	apikeys, err := t.apikeyStorage.GetAll()
	if err != nil {
		return []apikey.ApiKey{}, err
	}
	arr := []apikey.ApiKey{}
	for _, apikey := range apikeys {
		arr = append(arr, apikey)
	}
	return arr, nil
}

func (t *ApiKeyOperator) Create(name string) (*apikey.ApiKey, error) {
	apikey := apikey.ApiKey{
		Name:      name,
		Key:       apikey.GenerateApikey(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := t.apikeyStorage.Save(apikey)
	return &apikey, err
}

func (t *ApiKeyOperator) Delete(apiKey string) error {
	return t.apikeyStorage.Delete(apiKey)
}
