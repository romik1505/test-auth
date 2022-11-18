package main

import (
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
	app.Run()
}
