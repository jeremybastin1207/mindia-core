package config

type S3StorageConfig struct {
	Bucket string `yaml:"bucket,omitempty" validate:"required"`
}

type FilesystemStorageConfig struct {
	MountDir string `yaml:"mount_dir,omitempty" validate:"required"`
}

type FileStorageConfig struct {
	FilesystemStorageConfig *FilesystemStorageConfig `yaml:"filesystem,omitempty"`
	S3StorageConfig         *S3StorageConfig         `yaml:"s3,omitempty"`
}

type MetadataStorageConfig struct {
	Redis *string `yaml:"redis"`
}

type MediaStorageConfig struct {
	FileStorage     FileStorageConfig     `yaml:"file" validate:"required"`
	CacheStorage    FileStorageConfig     `yaml:"cache" validate:"required"`
	MetadataStorage MetadataStorageConfig `yaml:"metadata" validate:"required"`
}

type NamedTransforationStorageConfig struct {
	Filesystem *string `yaml:"filesystem"`
	Redis      *string `yaml:"redis"`
}

type ApiKeyStorageConfig struct {
	Filesystem *string `yaml:"filesystem"`
	Redis      *string `yaml:"redis"`
}

type TaskStorageConfig struct {
	Filesystem *string `yaml:"filesystem"`
	Redis      *string `yaml:"redis"`
}

type StorageConfig struct {
	MediaStorage              MediaStorageConfig              `yaml:"media" validate:"required"`
	NamedTransforationStorage NamedTransforationStorageConfig `yaml:"named_transformation" validate:"required"`
	ApiKeyStorage             ApiKeyStorageConfig             `yaml:"apikey" validate:"required"`
	TaskStorage               TaskStorageConfig               `yaml:"task" validate:"required"`
}
