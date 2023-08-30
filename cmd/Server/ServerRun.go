package main

import (
	"ComputerShopServer/internal"
	"ComputerShopServer/internal/DataBaseImplement/Config"
	"github.com/caarlos0/env/v8"
	"log"
)

func main() {
	cfg := Config.Config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("failed to retrieve env variables, %v", err)
	}

	if err := internal.Run(cfg); err != nil {
		log.Fatal("error running grpc server ", err)
	}

}
