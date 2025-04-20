package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sergunchig/merchant_exp.git/config"
	handler "github.com/sergunchig/merchant_exp.git/internal/handlers/offers"
	importHandler "github.com/sergunchig/merchant_exp.git/internal/handlers/offers/ImportHandler"
	"github.com/sergunchig/merchant_exp.git/internal/handlers/offers/readHandler"
	"github.com/sergunchig/merchant_exp.git/internal/repo"
	"github.com/sergunchig/merchant_exp.git/internal/services/readOffers"
	"github.com/sergunchig/merchant_exp.git/internal/services/writeOffers"
	"github.com/sergunchig/merchant_exp.git/internal/storage"
	"github.com/sergunchig/merchant_exp.git/internal/storage/excelReader"
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
	excelReader := excelReader.New(log)
	readService := readOffers.New(offerRepo, log)
	writeService := writeOffers.New(excelReader, offerRepo, log)
	storageService := storage.New(log)

	offerhandler := handler.New(log)
	readHandler := readHandler.New(readService, log)
	importHandler := importHandler.New(writeService, storageService, log)

	mux := http.NewServeMux()

	mux.HandleFunc("/", offerhandler.HomeHandler)
	mux.HandleFunc("/uploadandimport", importHandler.UploadAndImportHandler)
	mux.HandleFunc("/getoffers/", readHandler.GetOffers)

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
	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error(fmt.Errorf("error shutdown %w", err).Error())
	}

	fmt.Println("app was stoped")
}
