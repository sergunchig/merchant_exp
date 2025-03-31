package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sergunchig/merchant_exp.git/config"
	offer "github.com/sergunchig/merchant_exp.git/internal/handlers/offers"
	"github.com/sergunchig/merchant_exp.git/internal/repo"
	"github.com/sergunchig/merchant_exp.git/internal/storage/excel_reader"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
	"github.com/sergunchig/merchant_exp.git/pkg/postgres"
)

func Run(cfg *config.Config) {
	fmt.Println("start project...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	log, err := logger.NewLogger(cfg.Log.PATH)
	if err != nil {
		panic(fmt.Errorf("logger error, %w", err))
	}

	db, err := postgres.New(cfg.DB.DBCONNECTION)
	if err != nil {
		panic(fmt.Errorf("potgres error, %w", err))
	}
	defer db.Close()

	offerRepo := repo.New(db, log)
	//пакеты кэмелкейсом
	excelReader := excel_reader.New(log)
	offerhandler := offer.New(offerRepo, excelReader, log)

	mux := http.NewServeMux()

	mux.HandleFunc("/", offerhandler.HomeHandler)
	mux.HandleFunc("/upload_and_import", offerhandler.UploadAndImportHandler)
	mux.HandleFunc("/get_offers", offerhandler.GetOffers)

	srv := &http.Server{
		Addr:    cfg.HTTP.HOST,
		Handler: mux,
	}
	go func() {
		err = srv.ListenAndServe()
		fmt.Println("server was started...")
		if err != nil {
			panic(fmt.Errorf("server can't start %w", err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, closeFunc := context.WithTimeout(context.Background(), time.Minute)
	defer closeFunc()
	//nolint:errcheck
	srv.Shutdown(shutdownCtx)
	//log

	fmt.Println("app was stoped")
}
