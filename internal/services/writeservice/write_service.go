//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package writeservice

import (
	"context"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
)

type excelReader interface {
	Read(file string) ([]entity.Offer, error)
}
type offerRepo interface {
	CreateOffers(ctx context.Context, offers []entity.Offer) error
}

type WriteService struct {
	excelReader excelReader
	repo        offerRepo
}

func New(reader excelReader, repo offerRepo) *WriteService {
	return &WriteService{
		excelReader: reader,
		repo:        repo,
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
