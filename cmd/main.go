package main

import (
	"log"

	"github.com/Komilov31/calendar-service/cmd/api"
	"github.com/Komilov31/calendar-service/internal/config"
)

func main() {
	apiServer := api.NewServer(config.Envs.Port)
	if err := apiServer.Run(); err != nil {
		log.Fatal(err)
	}
}
