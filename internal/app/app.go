package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sergunchig/merchant_exp.git/config"
	offer "github.com/sergunchig/merchant_exp.git/internal/handlers/offers"
	"github.com/sergunchig/merchant_exp.git/internal/repo"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
	"github.com/sergunchig/merchant_exp.git/pkg/postgres"
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
	offerRepo := repo.New(db, log)
	offerhandler := offer.New(offerRepo, log)

	mux := http.NewServeMux()

	mux.HandleFunc("/", offerhandler.HomeHandler)
	mux.HandleFunc("/upload", offerhandler.UploadHandler)
	fmt.Println(cfg.HTTP.HOST)
	srv := &http.Server{
		Addr:    cfg.HTTP.HOST,
		Handler: mux,
	}
	err = srv.ListenAndServe()
	fmt.Println("server was started...")
	if err != nil {
		panic(fmt.Errorf("server can't start %w", err))
	}

	<-ctx.Done()

	srv.Shutdown(ctx)

	fmt.Println("app was stoped")
}
