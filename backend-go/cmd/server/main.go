package main

import (
	"context"
	"log"

	"github.com/williamchand/fullstack-fastapi/backend-go/config"
	"github.com/williamchand/fullstack-fastapi/backend-go/internal/app"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	application, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := application.Run(context.Background()); err != nil {
		log.Fatal(err)
	}
}
