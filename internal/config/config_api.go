package config

type HttpApiConfig struct {
	Host string `yaml:"host" validate:"required"`
	Port int    `yaml:"port" validate:"required"`
}

type ServerConfig struct {
	HttpApiConfig HttpApiConfig `yaml:"http,omitempty" validate:"required"`
}
