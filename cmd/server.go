package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v11"
	"github.com/estebangarcia/cm3070-final-project/pkg/api"
	"github.com/estebangarcia/cm3070-final-project/pkg/config"
	"github.com/estebangarcia/cm3070-final-project/pkg/helpers"
)

func NewSigKillContext() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}

func main() {
	var cfg config.AppConfig

	if err := env.Parse(&cfg); err != nil {
		log.Fatal(err)
	}

	ctx := NewSigKillContext()

	dbClient, err := helpers.GetDBClient(ctx, &cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	r, err := api.NewRouter(ctx, cfg, dbClient)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Run(ctx, fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.Fatal(err)
	}
}
