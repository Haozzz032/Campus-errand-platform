package main

import (
	"Campus-errand-platform/config"
	"Campus-errand-platform/database"
	"Campus-errand-platform/middleware"
	"Campus-errand-platform/model"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()
	server.GET("/")
	cfg := config.Loadconfig()
	db := database.InitDB(cfg)

	//需要登录
	Auth := server.Group("/", middleware.Auth(&cfg))
	{

	}
	//可选登录
	OptionalAuth := server.Group("/", middleware.OptionalAuth(&cfg))
	{

	}
	//需要用户身份登录
	UserAuth := server.Group("/", middleware.RequireRole(model.RoleUser, &cfg, db))
	{

	}
	//需要跑者身份登录
	RunnerAuth := server.Group("/", middleware.RequireRole(model.RoleRunner, &cfg, db))
	{

	}
	//需要管理员身份登录
	AdminAuth := server.Group("/", middleware.RequireRole(model.RoleAdmin, &cfg, db))
	{

	}

	if err := server.Run(cfg.Addr()); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
