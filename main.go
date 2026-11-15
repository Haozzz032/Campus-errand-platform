package main

import (
	"Campus-errand-platform/config"
	"Campus-errand-platform/middleware"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/")
	cfg := config.Loadconfig()
	user := server.Group("/api/user", middleware.Auth(&cfg))
	{

	}

	OptionalAuth := server.Group("/", middleware.OptionalAuth(&cfg))
	{

	}

	if err := server.Run(cfg.Addr()); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
