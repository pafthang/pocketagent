package rbac

import (
	"testing"
	"time"

	"github.com/pafthang/pocketagent/pkgs/models"
)

func TestAuthorizeCacheHitMiss(t *testing.T) {
	cache := newAuthorizeCache(30 * time.Second)
	cache.set("u1", "s1", ActionAgentRead, models.AuthorizeResponse{Allowed: true, Role: models.RoleAdmin})
	if resp, ok := cache.get("u1", "s1", ActionAgentRead); !ok || !resp.Allowed {
		t.Fatal("expected cache hit")
	}
	if _, ok := cache.get("u1", "s1", ActionTaskRead); ok {
		t.Fatal("unexpected cache hit for different action")
	}
}

func TestAuthorizeCacheExpiry(t *testing.T) {
	cache := newAuthorizeCache(20 * time.Millisecond)
	cache.set("u1", "s1", ActionAgentRead, models.AuthorizeResponse{Allowed: true})
	time.Sleep(30 * time.Millisecond)
	if _, ok := cache.get("u1", "s1", ActionAgentRead); ok {
		t.Fatal("expected cache expiry")
	}
}
