//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package readHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/sergunchig/merchant_exp.git/dto"
)

type readService interface {
	GetOffers(ctx context.Context) ([]dto.OfferDto, error)
	GetOffer(ctx context.Context, offer_id int) (dto.OfferDto, error)
}
type offerLogger interface {
	Error(msg string)
}

type ReadHandler struct {
	service readService
	log     offerLogger
}

func New(service readService, log offerLogger) *ReadHandler {
	return &ReadHandler{
		service: service,
		log:     log,
	}
}

func (h *ReadHandler) GetOffers(rw http.ResponseWriter, r *http.Request) {
	offers, err := h.service.GetOffers(r.Context())

	if err != nil {
		h.log.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	json, err := json.Marshal(offers)
	if err != nil {
		err = fmt.Errorf("error marshaled offers %w", err)
		h.log.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	rw.Write(json)
}
func (h *ReadHandler) GetOffer(rw http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {
		h.log.Error(err.Error())
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("not id"))
	}

	offer, err := h.service.GetOffer(r.Context(), id)
	if err != nil {
		h.log.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
	}

	json, err := json.Marshal(offer)
	if err != nil {
		err = fmt.Errorf("can't marshal offer id = %d : %w", id, err)
		h.log.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
	}
	rw.WriteHeader(http.StatusOK)
	rw.Write(json)
}
