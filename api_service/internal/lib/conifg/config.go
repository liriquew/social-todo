package conifg

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port string `yaml:"port" env-default:"4040"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" {
		path = "./config/config.yaml"
	}

	return MustLoadByPath(path)
}

func fetchConfigPath() string {
	return os.Getenv("CONF_PATH")
}

func MustLoadByPath(path string) *Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config: file not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("error while reading config" + err.Error())
	}

	return &cfg
}
