//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package offer

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	offerServices "github.com/sergunchig/merchant_exp.git/internal/Services"
	"github.com/sergunchig/merchant_exp.git/internal/entity"
)

type readOfferService interface {
	GetOffersAsync(ctx context.Context) ([]offerServices.OfferDto, error)
}

type excelOffersReader interface {
	Read(file string) ([]entity.Offer, error)
}
type offerLogger interface {
	Error(msg string)
}

type Handler struct {
	service readOfferService
	reader  excelOffersReader
	log     offerLogger
}

func New(service readOfferService, reader excelOffersReader, log offerLogger) *Handler {
	return &Handler{
		//offers: repo,
		service: service,
		reader:  reader,
		log:     log,
	}
}

func (s *Handler) HomeHandler(rw http.ResponseWriter, r *http.Request) {
	html := `<html>
	<body>
	<div>
		<form action="/UploadAndImport" method="post" enctype="multipart/form-data">
			Excel file: <input type="file" name="my_file">
			<input type="submit" value="Import">
		</form>
	</div>
	</body>
</html>`

	tmpl := template.Must(template.New("loadpage").Parse(html))

	err := tmpl.Execute(rw, nil)

	if err != nil {
		s.log.Error(fmt.Errorf("cant execute template 'home' %w", err).Error())
	}
}

func (h *Handler) GetOffersAsync(rw http.ResponseWriter, r *http.Request) {
	offers, err := h.service.GetOffersAsync(r.Context())

	if err != nil {
		h.log.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	json, err := json.Marshal(offers)
	rw.Write(json)
}

// func (h *Handler) UploadAndImportHandler(rw http.ResponseWriter, r *http.Request) {
// 	uploadData, _, err := r.FormFile("my_file")
// 	if err != nil {
// 		h.log.Error(fmt.Errorf("cant parse file %w", err).Error())
// 		http.Error(rw, "request error", http.StatusInternalServerError)
// 		return
// 	}
// 	defer uploadData.Close()

// 	file := "./storage/excelfile.xlsx"
// 	err = storage.SaveFile(uploadData, file) //mock
// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.Write([]byte(err.Error()))
// 		return
// 	}
// 	offers, err := h.reader.Read(file)

// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.Write([]byte(err.Error()))
// 		return
// 	}

// 	err = h.offers.CreateOffers(r.Context(), offers)
// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.Write([]byte(err.Error()))
// 		return
// 	}
// 	rw.Write([]byte("Offers import is successfully"))
// }

// todo viewmodel
// func (h *Handler) GetOffers(rw http.ResponseWriter, r *http.Request) {

// 	offers, err := h.offers.GetOffers(r.Context())
// 	if err != nil {
// 		h.log.Error(err.Error())
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		rw.Write([]byte("status 500"))
// 		return
// 	}
// 	data, err := json.Marshal(offers)
// 	if err != nil {
// 		h.log.Error(fmt.Errorf("json marshal error: %w", err).Error())
// 		rw.WriteHeader(http.StatusInternalServerError)
// 		rw.Write([]byte("status 500"))
// 		return
// 	}
// 	rw.Header().Set("Content-Type", "application/json")
// 	rw.WriteHeader(http.StatusOK)
// 	rw.Write(data)
// }
