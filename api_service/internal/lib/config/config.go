package config

import (
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port          string        `yaml:"port" env-default:"4040"`
	AuthConfig    ServiceConfig `yaml:"auth_client" env-required:"true"`
	NoteConfig    ServiceConfig `yaml:"note_client" env-required:"true"`
	FriendsConfig ServiceConfig `yaml:"friends_client" env-required:"true"`
}

type ServiceConfig struct {
	Port    string        `yaml:"port" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-defauilt:"1s"`
	Retries int           `yaml:"retries" env-default:"10"`
}

func MustLoad() *Config {
	path := fetchConfigPath()

	if path == "" { // lol idk
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
