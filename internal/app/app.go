package app

import (
	"context"
	"fmt"
	"merchant_exp/config"
	"merchant_exp/internal/repo"
	"merchant_exp/pkg/logger"
	"merchant_exp/pkg/postgres"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

func Run(cfg *config.Config) {
	fmt.Println("start project...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5) //WithCancel(context.Background())
	defer cancel()

	log, err := logger.NewLogger(cfg.Log.FILE)

	if err != nil {
		panic(fmt.Errorf("logger error, %w", err))
	}

	db, err := postgres.New(cfg.DB.DBCONNECTION)
	if err != nil {
		panic(fmt.Errorf("potgres error, %w", err))
	}
	defer db.Close()
	offerRepo := repo.New(db)

	offerhandler := handlers.New(offerRepo)

	mux := http.NewServeMux()

	mux.HandleFunc("/", offerhandler.HomeHandler)
	mux.HandleFunc("/upload", offerhandler.UploadHandler)
	fmt.Println(s.cfg.HOST)
	srv := &http.Server{
		Addr:    cfg.HTTP.HOST,
		Handler: mux,
	}

	fmt.Println("server was started...")
	if err != nil {
		fmt.Println(err)
	}

	<-ctx.Done()

	srv.Shutdown()

	fmt.Println("app was stoped")
}
