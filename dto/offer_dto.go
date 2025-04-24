// todo из названия файла надо удлить _dto
package dto

import "github.com/sergunchig/merchant_exp.git/internal/entity"

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
