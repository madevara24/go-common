package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"
)

type IServer interface {
	StartServer()
}

type server struct {
	httpServer *http.Server
}

type Option func(*server)

func NewServer(handler http.Handler, config Config, opts ...Option) IServer {
	server := server{
		httpServer: &http.Server{
			Handler:      handler,
			Addr:         fmt.Sprintf(":%s", config.Address),
			ReadTimeout:  config.ReadTimeout,
			WriteTimeout: config.WriteTimeout,
		},
	}

	for _, opt := range opts {
		opt(&server)
	}

	return &server
}

func (g *server) StartServer() {
	go func() {
		if err := g.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()

	log.Println("server is running ...")

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := g.httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err.Error())
		os.Exit(1)
	}

	log.Println("Server stopped.")
}
