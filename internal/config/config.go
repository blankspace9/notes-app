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
		Env          string     `yaml:"env" env-default:"local"`
		HTTPServer   HTTPServer `yaml:"http"`
		JWT          JWT        `yaml:"tokens"`
		Storage      Postgres
		SpellChecker SpellChecker
	}

	HTTPServer struct {
		Port    int           `yaml:"port"`
		Timeout time.Duration `yaml:"timeout"`
	}

	JWT struct {
		Secret          []byte        `env:"JWT_SECRET"`
		AccessTokenTTL  time.Duration `yaml:"access_token_ttl" env-required:"true"`
		RefreshTokenTTL time.Duration `yaml:"refresh_token_ttl" env-required:"true"`
	}

	Postgres struct {
		Host     string `env:"POSTGRES_HOST" env-default:"localhost"`
		Port     string `env:"PORT" env-default:"5432"`
		Username string `env:"POSTGRES_USERNAME" env-default:"postgres"`
		DBName   string `env:"POSTGRES_DB"`
		SSLMode  string `env:"POSTGRES_SSLMODE"`
		Password string `env:"POSTGRES_PASSWORD"`
	}

	SpellChecker struct {
		URL string `env:"SPELL_CHECKER_URL"`
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

	if err := cleanenv.ReadEnv(&cfg.SpellChecker); err != nil {
		panic("failed to read config: " + err.Error())
	}

	if err := cleanenv.ReadEnv(&cfg.JWT); err != nil {
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
