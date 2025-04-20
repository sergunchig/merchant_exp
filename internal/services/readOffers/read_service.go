//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package readOffers

import (
	"context"

	"github.com/sergunchig/merchant_exp.git/dto"
	"github.com/sergunchig/merchant_exp.git/internal/entity"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
)

type repoOffers interface {
	Read(ctx context.Context) ([]entity.Offer, error)
}

type OfferService struct {
	repo repoOffers
	log  *logger.AppLogger
}

func New(repo repoOffers, log *logger.AppLogger) *OfferService {
	return &OfferService{
		repo: repo,
		log:  log,
	}
}

func (o *OfferService) GetOffers(ctx context.Context) ([]dto.OfferDto, error) {
	offers, err := o.repo.Read(ctx)
	if err != nil {
		return nil, err
	}
	offersDto, err := o.offersDto(offers)
	return offersDto, err
}

func (o *OfferService) offersDto(offers []entity.Offer) ([]dto.OfferDto, error) {
	offersDto := make([]dto.OfferDto, 0, 1000)

	for _, offer := range offers {
		offersDto = append(offersDto, dto.MakeOfferDisplay(offer))
	}

	return offersDto, nil
}
