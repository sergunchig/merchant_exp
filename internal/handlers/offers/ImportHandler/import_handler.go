//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package importHandler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sergunchig/merchant_exp.git/internal/storage"
)

type offerLogger interface {
	Error(msg string)
}
type importServices interface {
	ImportOffers(ctx context.Context, file string) error
}

type WriteHandler struct {
	service importServices
	log     offerLogger
}

func New(service importServices, log offerLogger) *WriteHandler {
	return &WriteHandler{
		service: service,
		log:     log,
	}
}

func (h *WriteHandler) UploadAndImportHandler(rw http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("my_file")
	if err != nil {
		h.log.Error(fmt.Errorf("cant parse file %w", err).Error())
		http.Error(rw, "request error", http.StatusInternalServerError)
		return
	}
	defer uploadData.Close()

	file := "./storage/excelfile.xlsx"
	err = storage.SaveFile(uploadData, file) //mock
	if err != nil {
		h.log.Error(err.Error())
		rw.Write([]byte(err.Error()))
		return
	}
	err = h.service.ImportOffers(r.Context(), file)

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		h.log.Error(err.Error())
		return
	}
}
