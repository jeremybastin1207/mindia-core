package config

type S3AdapterConfig struct {
	AccessKeyId     string `yaml:"access_key_id,omitempty" validate:"required"`
	SecretAccessKey string `yaml:"secret_access_key,omitempty" validate:"required"`
	Endpoint        string `yaml:"endpoint,omitempty" validate:"required"`
	Region          string `yaml:"region,omitempty" validate:"required"`
}

type RedisAdapterConfig struct {
	Host     string `yaml:"host,omitempty" validate:"required"`
	Port     int    `yaml:"port,omitempty" validate:"required"`
	Password string `yaml:"password,omitempty" validate:""`
}

type AdapatersConfig struct {
	S3    *S3AdapterConfig    `yaml:"s3,omitempty"`
	Redis *RedisAdapterConfig `yaml:"redis,omitempty"`
}
