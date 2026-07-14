package client

import (
	"sync"
	"time"

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
	return &authorizeCache{
		ttl:     ttl,
		entries: make(map[authorizeCacheKey]authorizeCacheEntry),
	}
}

func (c *authorizeCache) get(userID, spaceID, action string) (models.AuthorizeResponse, bool) {
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
	key := authorizeCacheKey{userID: userID, spaceID: spaceID, action: action}
	c.mu.Lock()
	c.entries[key] = authorizeCacheEntry{
		resp:      resp,
		expiresAt: time.Now().Add(c.ttl),
	}
	c.mu.Unlock()
}

// EnableAuthorizeCache caches RBAC authorize responses keyed by user, space, and action.
// TTL should stay short (e.g. 30s) so role changes propagate reasonably quickly.
func (c *Client) EnableAuthorizeCache(ttl time.Duration) {
	if ttl <= 0 {
		c.authorizeCache = nil
		return
	}
	c.authorizeCache = newAuthorizeCache(ttl)
}

// AuthorizeCached checks permission, using a TTL cache when enabled and userID is set.
func (c *Client) AuthorizeCached(userID, token, spaceID, action string) (models.AuthorizeResponse, error) {
	if c.authorizeCache != nil && userID != "" {
		if resp, ok := c.authorizeCache.get(userID, spaceID, action); ok {
			return resp, nil
		}
	}

	resp, err := c.Authorize(token, spaceID, action)
	if err != nil {
		return resp, err
	}
	if c.authorizeCache != nil && userID != "" {
		c.authorizeCache.set(userID, spaceID, action, resp)
	}
	return resp, nil
}
