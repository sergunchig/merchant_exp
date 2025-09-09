// Название ImportHandler надо сделать в go way стиле
//
//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package importer

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

type offerLogger interface {
	Error(msg string)
}
type importServices interface {
	ImportOffers(ctx context.Context, file string) error
}
type storageService interface {
	SaveFile(in io.Reader, fileName string) error
}

type WriteHandler struct {
	service importServices
	log     offerLogger
	storage storageService
}

func New(service importServices, storage storageService, log offerLogger) *WriteHandler {
	return &WriteHandler{
		service: service,
		log:     log,
		storage: storage,
	}
}

func (h *WriteHandler) UploadAndImport(rw http.ResponseWriter, r *http.Request) {
	uploadData, _, err := r.FormFile("file")
	if err != nil {
		h.log.Error(fmt.Errorf("cant parse file %w", err).Error())
		http.Error(rw, "request error", http.StatusInternalServerError)
		return
	}
	defer func() {
		err = uploadData.Close()
		if err != nil {
			h.log.Error(err.Error())
		}
	}()

	file := "./storage/excelfile.xlsx"
	err = h.storage.SaveFile(uploadData, file) //mock
	if err != nil {
		h.log.Error(err.Error())
		rw.Write([]byte(err.Error()))
		return
	}
	err = h.service.ImportOffers(r.Context(), file)
	if err != nil { // todo между err и проверкой не должно быть пустой строки
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		h.log.Error(err.Error())
		return
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("imported"))
}
