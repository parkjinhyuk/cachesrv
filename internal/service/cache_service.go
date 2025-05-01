package service

import (
	"cachesrv/internal/model"
	"cachesrv/internal/repository"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type CacheService struct {
	repo repository.CacheRepository
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
	if cached != nil {
		return cached, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	cache := &model.Cache{
		Body:         string(body),
		ETag:         res.Header.Get("ETag"),
		LastModified: res.Header.Get("Last-Modified"),
		CachedAt:     time.Now(),
	}

	err = s.repo.Set(ctx, cacheKey, cache, ttl)
	if err != nil {
		fmt.Print(err)
	}

	return cache, nil
}
