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
	// parse dto
	cacheKey := c.Query("cache_key")
	apiURL := c.Query("api_url")
	ttl, err := time.ParseDuration(c.DefaultQuery("ttl", "60") + "s")
	if apiURL == "" || cacheKey == "" || err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
		return
	}

	// get or fetch
	cache, err := h.cacheService.GetOrFetch(c.Request.Context(), cacheKey, apiURL, ttl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// if-None-Match
	ifNoneMatch := c.GetHeader("If-None-Match")
	if ifNoneMatch != "" && cache.ETag != "" && ifNoneMatch == cache.ETag {
		c.Status(http.StatusNotModified)
		return
	}

	// if-Modified-Since
	ifModifiedSince := c.GetHeader("If-Modified-Since")
	if ifModifiedSince != "" && cache.LastModified != "" {
		clientTime, err := http.ParseTime(ifModifiedSince)
		if err == nil {
			serverTime, err := http.ParseTime(cache.LastModified)
			if err == nil && !serverTime.After(clientTime) {
				c.Status(http.StatusNotModified)
				return
			}
		}
	}

	// response
	if cache.ETag != "" {
		c.Header("ETag", cache.ETag)
	}
	if cache.LastModified != "" {
		c.Header("Last-Modified", cache.LastModified)
	}
	c.Data(http.StatusOK, "application/json", []byte(cache.Body))
}
