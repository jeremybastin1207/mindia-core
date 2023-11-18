package config

type Config struct {
	MasterKey string          `yaml:"-"`
	Server    ServerConfig    `yaml:"api,omitempty" validate:"required"`
	Storage   StorageConfig   `yaml:"storage,omitempty" validate:"required"`
	Adapters  AdapatersConfig `yaml:"adapter,omitempty" validate:"required"`
}

func NewConfig() Config {
	return Config{
		Server: ServerConfig{
			HttpApiConfig: HttpApiConfig{},
		},
		Storage: StorageConfig{
			MediaStorage: MediaStorageConfig{
				FileStorage:     FileStorageConfig{},
				CacheStorage:    FileStorageConfig{},
				MetadataStorage: MetadataStorageConfig{},
			},
			NamedTransforationStorage: NamedTransforationStorageConfig{},
			ApiKeyStorage:             ApiKeyStorageConfig{},
			TaskStorage:               TaskStorageConfig{},
		},
		Adapters: AdapatersConfig{},
	}
}
