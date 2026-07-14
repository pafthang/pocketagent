package rbac

import (
	"sync"
	"time"

	pbclient "github.com/pafthang/pocketagent/internal/pocket/client"
	"github.com/pafthang/pocketagent/pkgs/models"
)

type authorizeCacheKey struct {
	userID  string
	spaceID string
	action  string
}

type authorizeCacheEntry struct {
	resp      models.AuthorizeResponse
	expiresAt time.Time
}

type authorizeCache struct {
	ttl     time.Duration
	mu      sync.RWMutex
	entries map[authorizeCacheKey]authorizeCacheEntry
}

func newAuthorizeCache(ttl time.Duration) *authorizeCache {
	if ttl <= 0 {
		return nil
	}
	return &authorizeCache{
		ttl:     ttl,
		entries: make(map[authorizeCacheKey]authorizeCacheEntry),
	}
}

func (c *authorizeCache) get(userID, spaceID, action string) (models.AuthorizeResponse, bool) {
	if c == nil {
		return models.AuthorizeResponse{}, false
	}
	key := authorizeCacheKey{userID: userID, spaceID: spaceID, action: action}
	now := time.Now()

	c.mu.RLock()
	entry, ok := c.entries[key]
	c.mu.RUnlock()
	if !ok || now.After(entry.expiresAt) {
		if ok {
			c.mu.Lock()
			delete(c.entries, key)
			c.mu.Unlock()
		}
		return models.AuthorizeResponse{}, false
	}
	return entry.resp, true
}

func (c *authorizeCache) set(userID, spaceID, action string, resp models.AuthorizeResponse) {
	if c == nil {
		return
	}
	key := authorizeCacheKey{userID: userID, spaceID: spaceID, action: action}
	c.mu.Lock()
	c.entries[key] = authorizeCacheEntry{
		resp:      resp,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

// CachedAuthorizer evaluates RBAC in-process with optional TTL caching.
type CachedAuthorizer struct {
	inner *Authorizer
	cache *authorizeCache
}

// NewCachedAuthorizer wraps the PocketBase-backed authorizer.
func NewCachedAuthorizer(pb *pbclient.Client, cacheTTL time.Duration) *CachedAuthorizer {
	return &CachedAuthorizer{
		inner: NewAuthorizer(pb),
		cache: newAuthorizeCache(cacheTTL),
	}
}

// Authorize checks permission, using the TTL cache when configured.
func (c *CachedAuthorizer) Authorize(userID, spaceID, action string) (models.AuthorizeResponse, error) {
	if userID != "" {
		if resp, ok := c.cache.get(userID, spaceID, action); ok {
			return resp, nil
		}
	}

	resp, err := c.inner.Authorize(userID, spaceID, action)
	if err != nil {
		return resp, err
	}
	if userID != "" {
		c.cache.set(userID, spaceID, action, resp)
	}
	return resp, nil
}
