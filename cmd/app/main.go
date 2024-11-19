package main

import (
	"clearWayTest/config"
	"clearWayTest/internal/app"
	"fmt"
)

func main() {
	cfg, err := config.NewConfig()

	if err != nil {
		panic(fmt.Sprintf("Can't get config: %s", err))
	}

	app.Run(cfg)
}
