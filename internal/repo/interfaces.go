package repo

import (
	"context"
	"merchant_exp/internal/entity"
)

type OfferReader interface {
	GetOffers(ctx context.Context) ([]entity.Offer, error)
	GetOffer(offer_id int) (entity.Offer, error)
}

type OfferCreator interface {
	Create(ctx context.Context, Offer entity.Offer) error
	CreateOffers(ctx context.Context, offers []entity.Offer) error
}
type OfferCreatorAsync interface {
	CreateOffersPipe(ctx context.Context, in <-chan entity.Offer) <-chan entity.Offer
}

type ExcelReader interface {
	Read(ctx context.Context) ([]entity.Offer, error)
	ReadAsync(file string) (<-chan entity.Offer, error)
}
