package svc

import (
	"github.com/pafthang/pocketagent/internal/memo/auth"
	"github.com/pafthang/pocketagent/internal/memo/store"
	"github.com/pafthang/pocketagent/pkgs/service"
)

// RegisterRoutes wires the internal memo HTTP API.
func RegisterRoutes(s *service.Server, mgr *store.Manager, serviceToken string) {
	api := s.Echo.Group("", auth.RequireServiceToken(serviceToken))
	api.GET("/documents", listDocuments(mgr))
	api.GET("/documents/:id", getDocument(mgr))
	api.DELETE("/documents/:id", deleteDocument(mgr))
	api.POST("/documents", addDocument(mgr))
	api.POST("/search", searchDocuments(mgr))
	api.GET("/stats", collectionStats(mgr))
}