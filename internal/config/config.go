package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

const (
	dev  = "dev"
	test = "test"
	prod = "prod"
)

type Config struct {
	BotKey           string `envconfig:"BOT_KEY" yaml:"BotKey,omitempty"`
	ConnectionString string `yaml:"ConnectionString"`
	Env              string `envconfig:"ENV" yaml:"Env,omitempty"`
}

func MustGetConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("error while read config: ", err)
	}

	configPath := getConfigPath(os.Getenv("ENV"))

	var config Config
	readEnv(&config)
	readFile(&config, configPath)

	return &config
}

func readFile(cfg *Config, filePath string) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err.Error())
	}

	defer file.Close()

	if err := yaml.NewDecoder(file).Decode(cfg); err != nil {
		panic(err.Error())
	}
}

func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		panic(err.Error())
	}
}

func getConfigPath(env string) string {
	switch env {
	case dev:
		return "dev.yaml"
	case test:
		return "test.yaml"
	case prod:
		return "prod.yaml"
	default:
		panic("no config files found")
	}
}
