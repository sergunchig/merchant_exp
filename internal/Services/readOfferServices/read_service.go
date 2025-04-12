//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package readOfferServices

import (
	"context"

	"github.com/sergunchig/merchant_exp.git/dto"
	"github.com/sergunchig/merchant_exp.git/internal/entity"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
)

type repoOffers interface {
	ReadAsync(ctx context.Context) (<-chan entity.Offer, <-chan error)
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

func (o *OfferService) GetOffersAsync(ctx context.Context) ([]dto.OfferDto, error) {
	offerCh, errCh := o.repo.ReadAsync(ctx)
	offers, err := o.offersAsync(ctx, offerCh, errCh)
	return offers, err
}

func (o *OfferService) offersAsync(ctx context.Context, out <-chan entity.Offer, errCh <-chan error) ([]dto.OfferDto, error) {
	offersDto := make([]dto.OfferDto, 0, 1000)

	var inprocess bool = false
	var err error = nil

	for inprocess {
		select {
		case offer := <-out:
			offersDto = append(offersDto, dto.MakeOfferDisplay(offer))
		case err = <-errCh:
			inprocess = false
		case <-ctx.Done():
			inprocess = false
		}
	}

	return offersDto, err
}
