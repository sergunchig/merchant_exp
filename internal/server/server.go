package server

import (
	"context"
	"fmt"
	"merchant_exp/config"
	"merchant_exp/internal/repo"
	"merchant_exp/pkg/logger"
	"merchant_exp/pkg/postgres"
	"net/http"
)

type Server struct {
	cfg        *config.HTTP
	httpServer *http.Server
	logger     *logger.AppLogger
	db         *repo.OfferRepo
}

func New(cfg *config.Config) (*Server, error) {


	return &srv, nil
}
func (s *Server) createServer() *http.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HomeHandler)
	mux.HandleFunc("/upload", s.UploadHandler)
	fmt.Println(s.cfg.HOST)
	srv := &http.Server{
		Addr:    s.cfg.HOST,
		Handler: mux,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {

		}
	}()

	fmt.Println("server was created")
	return srv
}

func (s *Server) Shutdown() {
	s.httpServer.Shutdown(context.Background())
}
