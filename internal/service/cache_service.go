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

	if cached != nil {
		if cached.NoCache {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiUrl, nil)
			if err != nil {
				return nil, err
			}
			if cached.ETag != "" {
				req.Header.Set("If-None-Match", cached.ETag)
			}
			if cached.LastModified != "" {
				req.Header.Set("If-Modified-Since", cached.LastModified)
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

			newCache := &model.Cache{
				Body:         string(body),
				ETag:         res.Header.Get("ETag"),
				LastModified: res.Header.Get("Last-Modified"),
				CachedAt:     time.Now(),
				NoCache:      cc.NoCache,
			}

			if cc.NoStore {
				return newCache, nil
			}

			err = s.repo.Set(ctx, cacheKey, newCache, ttl)
			if err != nil {
				fmt.Print(err)
			}

			return newCache, nil
		}

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
	cc := parseCacheControl(res.Header.Get("Cache-Control"))

	cache := &model.Cache{
		Body:         string(body),
		ETag:         res.Header.Get("ETag"),
		LastModified: res.Header.Get("Last-Modified"),
		CachedAt:     time.Now(),
		NoCache:      cc.NoCache,
	}

	if cc.NoStore {
		return cache, nil
	}

	err = s.repo.Set(ctx, cacheKey, cache, ttl)
	if err != nil {
		fmt.Print(err)
	}

	return cache, nil
}

func parseCacheControl(header string) CacheControl {
	cc := CacheControl{}
	directives := strings.Split(header, ",")
	for _, d := range directives {
		if strings.TrimSpace(strings.ToLower(d)) == "no-cache" {
			cc.NoCache = true
		} else if strings.TrimSpace(strings.ToLower(d)) == "no-store" {
			cc.NoStore = true
		}
	}
	return cc
}
