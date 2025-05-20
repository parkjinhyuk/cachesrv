package main

import (
	"cachesrv/config"
	"cachesrv/internal/handler"
	"cachesrv/internal/repository"
	"cachesrv/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	cfg := config.LoadConfig()

	cacheRepo := repository.NewRedisRepository(cfg.Redis)
	cacheService := service.NewCacheService(cacheRepo)
	cacheHandler := handler.NewCacheHandler(cacheService)

	r := gin.Default()
	r.GET("/cache", cacheHandler.GetCache)

	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
