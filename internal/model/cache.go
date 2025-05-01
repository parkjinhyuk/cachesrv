package model

import "time"

type Cache struct {
	Body         string    `json:"body"`
	ETag         string    `json:"etag,omitempty"`
	LastModified string    `json:"last_modified,omitempty"`
	CachedAt     time.Time `json:"cached_at"`
}
