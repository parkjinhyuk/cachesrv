package handler

import (
	"cachesrv/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CacheHandler struct {
	cacheService *service.CacheService
}

func NewCacheHandler(cacheService *service.CacheService) *CacheHandler {
	return &CacheHandler{
		cacheService: cacheService,
	}
}

func (h *CacheHandler) GetCache(c *gin.Context) {
	cacheKey := c.Query("cache_key")
	apiURL := c.Query("api_url")
	ttl, _ := time.ParseDuration(c.Query("ttl"))

	if apiURL == "" || cacheKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	data, err := h.cacheService.GetOrFetch(c.Request.Context(), cacheKey, apiURL, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Data(http.StatusOK, "application/json", []byte(data))
}
