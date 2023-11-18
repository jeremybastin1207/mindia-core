package config

import (
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

const fileName = "config.yml"

type FilesystemStorage struct {
}

func NewFilesystemStorage() FilesystemStorage {
	return FilesystemStorage{}
}

func (c *FilesystemStorage) LoadConfig() (*Config, error) {
	config := NewConfig()

	/* 	body, err := os.ReadFile(fileName)
	   	if err != nil {
	   		return nil, err
	   	}

	   	err = yaml.Unmarshal(body, &config)
	   	if err != nil {
	   		return nil, err
	   	}
	*/
	apiHost, isEnv := os.LookupEnv("API_HOST")
	if isEnv {
		config.Server.HttpApiConfig.Host = apiHost
	}
	apiPort, isEnv := os.LookupEnv("API_PORT")
	if isEnv {
		port, _ := strconv.Atoi(apiPort)
		config.Server.HttpApiConfig.Port = port
	}

	s3AccessKeyId, isEnv := os.LookupEnv("S3_ACCESS_KEY_ID")
	if isEnv {
		config.Adapters.S3 = &S3AdapterConfig{
			AccessKeyId: s3AccessKeyId,
		}
		s3SecretAccessKey, isEnv := os.LookupEnv("S3_SECRET_ACCESS_KEY")
		if isEnv {
			config.Adapters.S3.SecretAccessKey = s3SecretAccessKey
		}
		s3Endpoint, isEnv := os.LookupEnv("S3_ENDPOINT")
		if isEnv {
			config.Adapters.S3.Endpoint = s3Endpoint
		}
		s3Region, isEnv := os.LookupEnv("S3_REGION")
		if isEnv {
			config.Adapters.S3.Region = s3Region
		}
	}

	redisHost, isEnv := os.LookupEnv("REDIS_HOST")
	if isEnv {
		config.Adapters.Redis = &RedisAdapterConfig{
			Host: redisHost,
		}
		redisPort, isEnv := os.LookupEnv("REDIS_PORT")
		if isEnv {
			port, _ := strconv.Atoi(redisPort)
			config.Adapters.Redis.Port = port
		}
		redisPassword, isEnv := os.LookupEnv("REDIS_PASSWORD")
		if isEnv {
			config.Adapters.Redis.Password = redisPassword
		}

		redis := ""
		config.Storage.MediaStorage.MetadataStorage.Redis = &redis
		config.Storage.TaskStorage.Redis = &redis
		config.Storage.NamedTransforationStorage.Redis = &redis
		config.Storage.ApiKeyStorage.Redis = &redis
	}

	fileBucketName, isEnv := os.LookupEnv("FILE_BUCKET_NAME")
	if isEnv {
		config.Storage.MediaStorage.FileStorage.S3StorageConfig = &S3StorageConfig{}
		config.Storage.MediaStorage.FileStorage.S3StorageConfig.Bucket = fileBucketName
	}
	cacheBucketName, isEnv := os.LookupEnv("CACHE_BUCKET_NAME")
	if isEnv {
		config.Storage.MediaStorage.CacheStorage.S3StorageConfig = &S3StorageConfig{}
		config.Storage.MediaStorage.CacheStorage.S3StorageConfig.Bucket = cacheBucketName
	}

	config.MasterKey = os.Getenv("MASTER_KEY")

	return &config, nil
}

func (c *FilesystemStorage) SaveConfig(config *Config) error {
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, yamlData, 0644)
}
