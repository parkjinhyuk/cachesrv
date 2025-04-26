package service

import (
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

func (s *CacheService) GetOrFetch(ctx context.Context, cacheKey string, apiUrl string, ttl time.Duration) (string, error) {
	cachedData, err := s.repo.Get(ctx, cacheKey)
	if err == nil {
		return cachedData, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
	if err != nil {
		return "", err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	err = s.repo.Set(ctx, cacheKey, string(body), ttl)
	if err != nil {
		fmt.Print(err)
	}

	return string(body), nil
}
