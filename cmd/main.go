package main

import (
	"log"

	"github.com/romik1505/auth/internal/config"
	"github.com/romik1505/auth/internal/handler"
	"github.com/romik1505/auth/internal/server"
	"github.com/romik1505/auth/internal/services/auth"
)

func main() {
	postgres := config.NewPostgresConnenction()
	auth := auth.NewAuthService(postgres)
	h := handler.NewHandler(auth)
	app := server.NewApp(h.InitRoutes())

	if err := app.Run(); err != nil {
		log.Fatalf("%v", err)
	}
}
