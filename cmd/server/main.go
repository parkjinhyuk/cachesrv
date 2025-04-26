package main

import (
	"cachesrv/internal/handler"
	"cachesrv/internal/repository"
	"cachesrv/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	redisAddr := "localhost:32768"
	redisPassword := "redispw"
	redisDB := 0

	cacheRepo := repository.NewRedisRepository(redisAddr, redisPassword, redisDB)
	cacheService := service.NewCacheService(cacheRepo)
	cacheHandler := handler.NewCacheHandler(cacheService)

	r := gin.Default()
	r.GET("/cache", cacheHandler.GetCache)

	port := 8080
	addr := fmt.Sprintf(":%d", port)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}
