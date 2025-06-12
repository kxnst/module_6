package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	MongoUri      string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	MongoUser     string `env:"MONGO_USER"`
	MongoPass     string `env:"MONGO_PASS"`
	MongoDatabase string `env:"MONGO_DATABASE"`
}

func NewConfig() *Config {
	cfg := &Config{}

	envName := os.Getenv("GO_ENV")
	if envName == "" {
		envName = "dev"
	}

	rootDir, err := filepath.Abs(".")
	if err != nil {
		log.Fatal("Config load failed: ", err)
	}

	for i := 0; i < 3; i++ {
		if _, err := os.Stat(filepath.Join(rootDir, "go.mod")); err == nil {
			break
		}
		rootDir = filepath.Dir(rootDir)
	}

	envFile := filepath.Join(rootDir, ".env")
	if envName == "test" {
		envFile = filepath.Join(rootDir, ".env.test")
	}

	if err := godotenv.Load(envFile); err != nil {
		log.Printf("can't load env file from %s: %v", envFile, err)
	}

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Can't parse env: %v", err)
	}

	return cfg
}
