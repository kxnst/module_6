package config

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	MongoUri  string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	MongoUser string `env:"MONGO_USER"`
	MongoPass string `env:"MONGO_PASS"`
}

func NewConfig() *Config {
	c := Config{}
	//add argument resolver for math to env file (need for tests and here)
	if err := godotenv.Load(); err != nil {
		log.Println("env file not loaded")
	}

	env.Parse(&c)

	return &c
}
