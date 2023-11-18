package main

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	redigo "github.com/gomodule/redigo/redis"
	"github.com/jeremybastin1207/mindia-core/internal/adapter/filesystem"
	"github.com/jeremybastin1207/mindia-core/internal/adapter/prometheus"
	"github.com/jeremybastin1207/mindia-core/internal/adapter/redis"
	"github.com/jeremybastin1207/mindia-core/internal/adapter/s3"
	"github.com/jeremybastin1207/mindia-core/internal/api"
	"github.com/jeremybastin1207/mindia-core/internal/apikey"
	"github.com/jeremybastin1207/mindia-core/internal/config"
	mindiaerr "github.com/jeremybastin1207/mindia-core/internal/error"
	"github.com/jeremybastin1207/mindia-core/internal/logging"
	"github.com/jeremybastin1207/mindia-core/internal/media"
	"github.com/jeremybastin1207/mindia-core/internal/plugin"
	"github.com/jeremybastin1207/mindia-core/internal/scheduler"
	"github.com/jeremybastin1207/mindia-core/internal/task"
	"github.com/jeremybastin1207/mindia-core/internal/transform"
	"github.com/joho/godotenv"
)

func main() {
	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load(".env")
		if err != nil {
			mindiaerr.ExitErrorf("unable to load .env, %v", err)
		}
	}

	var (
		logger                     = logging.New()
		s3Client                   *s3.S3
		redisPool                  *redigo.Pool
		fileStorage                media.FileStorer
		cacheStorage               media.FileStorer
		mediaStorage               media.Storer
		namedTransformationStorage transform.Storer
		apikeyStorage              apikey.Storer
		taskStorage                scheduler.Storer
		analyticsRecorder          = prometheus.NewPrometheusRecorder()
	)

	configStorage := config.NewFilesystemStorage()
	c, err := configStorage.LoadConfig()
	if err != nil {
		mindiaerr.ExitErrorf(err.Error())
	}

	validate := validator.New()
	err = validate.Struct(c)
	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		mindiaerr.ExitErrorf(validationErrors.Error())
	}

	if c.Adapters.S3 != nil {
		s3Client = s3.NewS3(s3.S3Config{
			AccessKeyId:     c.Adapters.S3.AccessKeyId,
			SecretAccessKey: c.Adapters.S3.SecretAccessKey,
			Endpoint:        c.Adapters.S3.Endpoint,
			Region:          c.Adapters.S3.Region,
		})
	}
	if c.Adapters.Redis != nil {
		redisAddr := fmt.Sprintf("%s:%v", c.Adapters.Redis.Host, c.Adapters.Redis.Port)
		redisPool = redis.NewPool(redisAddr)
	}

	if c.Storage.MediaStorage.FileStorage.FilesystemStorageConfig != nil {
		fileStorage = filesystem.NewFileStorage(filesystem.FileStorageConfig{
			MountDir: c.Storage.MediaStorage.FileStorage.FilesystemStorageConfig.MountDir,
		})
	} else if c.Storage.MediaStorage.FileStorage.S3StorageConfig != nil {
		fileStorage = s3.NewFileStorage(s3.FileStorageConfig{
			S3:     s3Client,
			Bucket: c.Storage.MediaStorage.FileStorage.S3StorageConfig.Bucket,
		})
	} else {
		mindiaerr.ExitErrorf("file storage config must be provided")
	}

	if c.Storage.MediaStorage.CacheStorage.FilesystemStorageConfig != nil {
		cacheStorage = filesystem.NewFileStorage(filesystem.FileStorageConfig{
			MountDir: c.Storage.MediaStorage.CacheStorage.FilesystemStorageConfig.MountDir,
		})
	} else if c.Storage.MediaStorage.CacheStorage.S3StorageConfig != nil {
		cacheStorage = s3.NewFileStorage(s3.FileStorageConfig{
			S3:     s3Client,
			Bucket: c.Storage.MediaStorage.CacheStorage.S3StorageConfig.Bucket,
		})
	} else {
		mindiaerr.ExitErrorf("cache storage config must be provided")
	}

	if c.Storage.MediaStorage.MetadataStorage.Redis != nil {
		mediaStorage = redis.NewMediaStorage(redisPool)
	} else {
		mindiaerr.ExitErrorf("media storage config must be provided")
	}

	if c.Storage.NamedTransforationStorage.Filesystem != nil {
		namedTransformationStorage = filesystem.NewNamedTransformationStorage()
	} else if c.Storage.NamedTransforationStorage.Redis != nil {
		namedTransformationStorage = redis.NewNamedTransformationStorage(redisPool)
	} else {
		mindiaerr.ExitErrorf("named transformation storage config must be provided")
	}

	if c.Storage.ApiKeyStorage.Filesystem != nil {
		apikeyStorage = filesystem.NewApiKeyStorage()
	} else if c.Storage.ApiKeyStorage.Redis != nil {
		apikeyStorage = redis.NewApiKeyStorage(redisPool)
	} else {
		mindiaerr.ExitErrorf("apikey storage config must be provided")
	}

	if c.Storage.TaskStorage.Filesystem != nil {
		// TODO
	} else if c.Storage.TaskStorage.Redis != nil {
		taskStorage = redis.NewTaskStorage(redisPool)
	}

	taskScheduler := scheduler.NewTaskScheduler(taskStorage, logger)
	go taskScheduler.ProcessTasks()

	pluginManager := plugin.NewPluginManager(fileStorage, cacheStorage, mediaStorage, taskStorage)

	colorizePlugin := plugin.NewColorizePlugin(&pluginManager)
	pluginManager.RegisterPlugin(&colorizePlugin)
	taskScheduler.RegisterListener(colorizePlugin.Name(), &colorizePlugin)

	tasks := api.Tasks{
		ClearCache:                  task.NewClearCacheTask(cacheStorage, analyticsRecorder),
		NamedTransformationOperator: task.NewNamedTransformationOperator(namedTransformationStorage),
		ApiKeyOperator:              task.NewApiKeyOperator(apikeyStorage),
		AnalyticsOperator:           task.NewAnalyticsOperator(fileStorage, cacheStorage),
		TaskOperator:                task.NewTaskOperator(taskStorage),
		GetMedia:                    task.NewGetMediaTask(mediaStorage, analyticsRecorder),
		DownloadMedia:               task.NewDownloadMediaTask(fileStorage, cacheStorage, mediaStorage, namedTransformationStorage, analyticsRecorder),
		UploadMedia:                 task.NewUploadMediaTask(fileStorage, cacheStorage, mediaStorage, namedTransformationStorage),
		DeleteMedia:                 task.NewDeleteMediaTask(fileStorage, cacheStorage, mediaStorage, analyticsRecorder),
		MoveMedia:                   task.NewMoveMediaTask(fileStorage, cacheStorage, mediaStorage),
		CopyMedia:                   task.NewCopyMediaTask(fileStorage, cacheStorage, mediaStorage),
		TagMedia:                    task.NewTagMediaTask(fileStorage, mediaStorage),
		ColorizeMedia:               task.NewColorizeMediaTask(&pluginManager),
	}

	server := api.NewApiServer(c.MasterKey, c.Server.HttpApiConfig.Host, c.Server.HttpApiConfig.Port, apikeyStorage, logger, tasks)
	server.Serve()
}
