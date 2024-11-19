package app

import (
	"clearWayTest/config"
	"clearWayTest/pkg/handler"
	"clearWayTest/pkg/httpserver"
	"clearWayTest/pkg/repository"
	"clearWayTest/pkg/service"
	"context"
	"fmt"
)

func Run(cfg *config.Config) {
	connDb, err := repository.NewPostgresDB(cfg)
	if err != nil {
		fmt.Printf("\n Can't connect to db: %s", err)
	}

	defer connDb.Close(context.Background())

	fmt.Printf("Start server %s \n", cfg.HTTP.Port)

	r := repository.NewRepository(connDb)
	s := service.NewService(r)
	h := handler.NewHandler(s)

	httpserver.Server(cfg, h)
}
