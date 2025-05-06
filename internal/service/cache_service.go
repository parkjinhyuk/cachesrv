package service

import (
	"cachesrv/internal/model"
	"cachesrv/internal/repository"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type CacheService struct {
	repo repository.CacheRepository
}

type CacheControl struct {
	NoCache bool
	NoStore bool
}

func NewCacheService(repo repository.CacheRepository) *CacheService {
	return &CacheService{
		repo: repo,
	}
}

func (s *CacheService) GetOrFetch(ctx context.Context, cacheKey string, apiUrl string, ttl time.Duration) (*model.Cache, error) {
	cached, err := s.repo.Get(ctx, cacheKey)
	if err != nil {
		return nil, err
	}

	if isValidCache(cached) {
		return cached, nil
	}

	newCache, err := s.fetchFromOrigin(ctx, apiUrl, cached)
	if err != nil {
		return nil, err
	}

	if newCache.NoStore {
		if err = s.repo.Set(ctx, cacheKey, newCache, ttl); err != nil {
			fmt.Print(err)
		}
	}

	return newCache, nil
}

func (s *CacheService) fetchFromOrigin(ctx context.Context, apiUrl string, cached *model.Cache) (*model.Cache, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	if cached != nil {
		if cached.ETag != "" {
			req.Header.Set("If-None-Match", cached.ETag)
		}
		if cached.LastModified != "" {
			req.Header.Set("If-Modified-Since", cached.LastModified)
		}
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotModified {
		return cached, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	cc := parseCacheControl(res.Header.Get("Cache-Control"))

	return &model.Cache{
		Body:         string(body),
		ETag:         res.Header.Get("ETag"),
		LastModified: res.Header.Get("Last-Modified"),
		CachedAt:     time.Now(),
		NoCache:      cc.NoCache,
		NoStore:      cc.NoStore,
	}, nil
}

func isValidCache(c *model.Cache) bool {
	return c != nil && !c.NoCache
}

func parseCacheControl(header string) CacheControl {
	cc := CacheControl{}
	directives := strings.Split(header, ",")
	for _, d := range directives {
		switch strings.ToLower(strings.TrimSpace(d)) {
		case "no-cache":
			cc.NoCache = true
		case "no-store":
			cc.NoStore = true
		}
	}

	return cc
}
