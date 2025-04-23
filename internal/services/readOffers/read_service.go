// todo название пакета поправить, тест написать
//
//go:generate mockgen -source ${GOFILE} -destination mocks_test.go -package ${GOPACKAGE}_test
package readOffers

import (
	"context"

	"github.com/sergunchig/merchant_exp.git/dto"
	"github.com/sergunchig/merchant_exp.git/internal/entity"
)

type repoOffers interface {
	Read(ctx context.Context) ([]entity.Offer, error)
	GetOffer(ctx context.Context, offer_id int) (entity.Offer, error) // todo offer_id не goway
}

type OfferService struct {
	repo repoOffers
}

func New(repo repoOffers) *OfferService {
	return &OfferService{
		repo: repo,
	}
}

func (o *OfferService) GetOffers(ctx context.Context) ([]dto.OfferDto, error) {
	offers, err := o.repo.Read(ctx)
	if err != nil {
		return nil, err
	}
	offersDto, err := o.offersDto(offers)
	return offersDto, err
} // todo нужно расстояние между функциями
func (o *OfferService) GetOffer(ctx context.Context, offer_id int) (dto.OfferDto, error) { // todo offer_id не goway
	offer, err := o.repo.GetOffer(ctx, offer_id)
	if err != nil {
		return dto.OfferDto{}, err
	}
	// todo не очень понятно обязанность сервиса? сходить в репозиторий и преобразовать в объект транспорта может и сам хендлер?
	return dto.MakeOfferDisplay(offer), nil
}

// todo это же тоже самое что функция dto.MakeOfferDisplay, только если бы она принимала слайс на вход у тебя две одинаковые только в двух разных местах
func (o *OfferService) offersDto(offers []entity.Offer) ([]dto.OfferDto, error) {
	offersDto := make([]dto.OfferDto, 0, 1000)

	for _, offer := range offers {
		offersDto = append(offersDto, dto.MakeOfferDisplay(offer))
	}

	return offersDto, nil
}
