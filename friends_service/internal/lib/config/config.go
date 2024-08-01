package config

import (
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port     int         `yaml:"port" env-default:"4041"`
	Neo4jCfg Neo4jConfig `yaml:"neo4j" env-required:"true"`
}

type Neo4jConfig struct {
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	DBName   string `yaml:"db_name" env-required:"true"`
}

func MustLoad() Config {
	path := fetchConfigPath()

	if path == "" {
		path = "./config/config.yaml"
	}

	return MustLoadPath(path)
}

func fetchConfigPath() string {
	return os.Getenv("CONF_PATH")
}

func MustLoadPath(path string) Config {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config: file not exist")
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("error while reading config" + err.Error())
	}

	return cfg
}
