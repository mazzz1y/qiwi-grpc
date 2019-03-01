package main

import (
	"log"

	"github.com/caarlos0/env"
)

type config struct {
	RPCPort   string `env:"PORT" envDefault:"50051"`
	Debug     bool `env:"DEBUG" envDefault:"false"`
	MongoHost string `env:"MONGO_HOST" envDefault:"localhost"`
	MongoPort string `env:"MONGO_PORT" envDefault:"27017"`
	MongoDB   string `env:"MONGO_DB" envDefault:"qiwi"`
}

var Config = initEnv()

func initEnv() config {
	c := config{}
	err := env.Parse(&c)
	if err != nil {
		log.Fatal(err)
	}
	if c.Debug {
		log.Println("Debug mode enabled")
	}
	return c
}
