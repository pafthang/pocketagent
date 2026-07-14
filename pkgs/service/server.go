package service

import "github.com/pafthang/pocketagent/pkgs/service/httpsrv"

// Server is a base HTTP microservice.
type Server = httpsrv.Server

// NewServer creates an HTTP server with common middleware.
func NewServer(name, listenAddr, logLevel string) *Server {
	return httpsrv.New(name, listenAddr, logLevel)
}