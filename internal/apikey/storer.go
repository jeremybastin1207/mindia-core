package apikey

type Storer interface {
	GetByName(name string) (*ApiKey, error)
	GetByKey(key string) (*ApiKey, error)
	GetAll() (ApiKeyMap, error)
	Save(apikey ApiKey) error
	Delete(apikey string) error
}
