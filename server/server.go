package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
	db         *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server {
	return &Server{
		db: db,
	}
}

func (s *Server) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// HTTP Server
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) {
	s.httpServer.Shutdown(ctx)
	s.db.Close()
}