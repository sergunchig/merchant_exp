package offerDto

import (
	"context"

	"github.com/sergunchig/merchant_exp.git/internal/entity"
)

type repoOffers interface {
	GetOffers(ctx context.Context) ([]entity.Offer, error)
}

type OfferService struct {
	repo repoOffers
}

func New(repo repoOffers) OfferService {
	return OfferService{
		repo: repo,
	}
}

func (o *OfferService) Offers(offers []entity.Offer) []OfferDto {
	offersDto := make([]OfferDto, 0, len(offers))
	for _, offer := range offers {
		offersDto = append(offersDto, OfferDto(offer))
	}
	return offersDto
}

func (o *OfferService) OffersAsync(ctx context.Context, in chan entity.Offer, errCh ) ([]OfferDto, error) {
	offersDto := make([]OfferDto, 0, 1000)
	var 
	for  {
		select {
		case offer := <-in:
			offersDto = append(offersDto, MakeOfferDisplay(offer))
		case <-ctx.Done():
			
		}
	}

	return offersDto, nil
}

type OfferDto struct {
	OfferId   int     `json:"offer_id"`
	Name      string  `json:"name"`
	Price     float64 `json:"price"`
	Available bool    `json:"available"`
}

func (os *OfferDto) ToOffer() entity.Offer {
	return entity.Offer{
		OfferId:   os.OfferId,
		Name:      os.Name,
		Price:     os.Price,
		Available: os.Available,
	}
}

func MakeOfferDisplay(offer entity.Offer) OfferDto {
	return OfferDto{
		OfferId:   offer.OfferId,
		Name:      offer.Name,
		Price:     offer.Price,
		Available: offer.Available,
	}
}
