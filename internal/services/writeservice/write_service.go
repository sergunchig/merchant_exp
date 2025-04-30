//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package writeservice

import (
	"context"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
)

type excelReader interface {
	//ReadAsync(ctx context.Context, done func(), file string) (<-chan entity.Offer, error)
	Read(file string) ([]entity.Offer, error)
}
type offerRepo interface {
	//CreateOffersAsync(ctx context.Context, in <-chan entity.Offer) error
	CreateOffers(ctx context.Context, offers []entity.Offer) error
}
type logger interface {
	Error(msg string)
}

type WriteService struct {
	excelReader excelReader
	repo        offerRepo
	log         logger
}

func New(reader excelReader, repo offerRepo, log logger) *WriteService {
	return &WriteService{
		excelReader: reader,
		repo:        repo,
		log:         log,
	}
}

func (ws *WriteService) ImportOffers(ctx context.Context, file string) error {
	offers, err := ws.excelReader.Read(file)
	if err != nil {
		return err
	}
	err = ws.repo.CreateOffers(ctx, offers)
	if err != nil {
		return err
	}
	return nil
}
