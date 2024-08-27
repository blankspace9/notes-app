package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Env        string     `yaml:"env" env-default:"local"`
		HTTPServer HTTPServer `yaml:"http"`
		Storage    Postgres
	}

	HTTPServer struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	}

	Postgres struct {
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"PORT" env-default:"5432"`
		Username string `env:"POSTGRES_USERNAME" env-default:"postgres"`
		Password string `env:"POSTGRES_PASSWORD"`
		DBName   string `env:"POSTGRES_DB"`
		SSLMode  string `env:"POSTGRES_SSLMODE"`
	}
)

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	return MustLoadByPath(path, ".env")
}

func MustLoadByPath(configPath string, envPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exists: " + configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	if err := godotenv.Load(envPath); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg.Storage); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &cfg
}

func MigrateMustLoad() *Postgres {
	cfg := new(Postgres)

	if err := godotenv.Load(".env"); err != nil {
		panic("failed to load .env file: " + err.Error())
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return cfg
}

// flag > env > default
func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
