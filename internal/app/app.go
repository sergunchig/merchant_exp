package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/sergunchig/merchant_exp.git/config"

	"github.com/sergunchig/merchant_exp.git/internal/handlers/offers/home"
	"github.com/sergunchig/merchant_exp.git/internal/handlers/offers/importer"
	"github.com/sergunchig/merchant_exp.git/internal/handlers/offers/reader"
	"github.com/sergunchig/merchant_exp.git/internal/repo/offer"

	readservice "github.com/sergunchig/merchant_exp.git/internal/services/readService"
	"github.com/sergunchig/merchant_exp.git/internal/services/writeservice"

	"github.com/sergunchig/merchant_exp.git/internal/storage"
	"github.com/sergunchig/merchant_exp.git/internal/storage/excelreader"
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

	db := postgres.MustInitPg(cfg.DB.DBCONNECTION)
	defer db.Close()

	offerRepo := offer.New(db, log)
	excelReader := excelreader.New(log)
	readService := readservice.New(offerRepo)
	writeService := writeservice.New(excelReader, offerRepo) // todo тут тоже лога не увидел. лог убрал
	storageService := storage.New(log)

	homeHandler := home.New(log)
	readHandler := reader.New(readService, log) // todo тут надо название поправить и сделать в camelCase, для однообразия лучше сделать alias offerHandler
	importHandler := importer.New(writeService, storageService, log)

	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler.Home)
	mux.HandleFunc("/uploadandimport", importHandler.UploadAndImport)
	mux.HandleFunc("/getoffers/", readHandler.GetOffers)
	mux.HandleFunc("/getoffer", readHandler.GetOffer)

	srv := &http.Server{
		Addr:    cfg.HTTP.HOST,
		Handler: mux,
	}
	go func() {
		err = srv.ListenAndServe()
		// todo ты увереш что это сообщение выведется? кажется что на этой функции горутина подвисает
		//fmt.Println("server was started...") //убираю, она не выводится
		if err != nil {
			panic(fmt.Errorf("server can't start: %w", err))
		}
	}()

	<-ctx.Done()
	// todo минута это долго, давай поставим секунд 10
	shutdownCtx, closeFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer closeFunc()
	//nolint:errcheck
	err = srv.Shutdown(shutdownCtx)
	if err != nil {
		log.Error(fmt.Errorf("error shutdown %w", err).Error())
	}

	fmt.Println("app was stoped")
}
