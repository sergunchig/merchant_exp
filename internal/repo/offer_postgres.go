package repo

import (
	"context"
	"fmt"
	"merchant_exp/internal/entity"
	"merchant_exp/pkg/logger"
	"merchant_exp/pkg/postgres"
	"strings"

	"github.com/jackc/pgx/v5"
)

type OfferRepo struct {
	client *postgres.Postgress
	log    *logger.AppLogger
}

func New(pg *postgres.Postgress, log *logger.AppLogger) *OfferRepo {
	return &OfferRepo{client: pg, log: log}
}

func (r *OfferRepo) Create(ctx context.Context, offer entity.Offer) error {
	query := "insert into offers (offer_id , \"name\" , price , available ) values (@offer_id, @name, @price, @available)"
	args := pgx.NamedArgs{
		"offer_id":  offer.OfferId,
		"name":      offer.Name,
		"price":     offer.Price,
		"available": offer.Available,
	}
	_, err := r.client.Pool.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error insert offer %d in db, %w", offer.OfferId, err)
	}
	return nil
}
func (r *OfferRepo) CreateOffers(ctx context.Context, offers []entity.Offer) error {
	var sb strings.Builder
	sb.WriteString("insert into offers (offer_id , \"name\" , price , available ) values \n")
	for _, offer := range offers {
		sb.WriteString(fmt.Sprintf("(%d, %s, %f, %t),", offer.OfferId, offer.Name, offer.Price, offer.Available))
	}
	sb.WriteString(";")
	// err :=
	// if err != nil {
	// 	fmt.Println(err)
	// }
	return nil
}
func (r *OfferRepo) CreateOffersPipe(ctx context.Context, in <-chan entity.Offer) error {

	for offer := range in {
		err := r.Create(ctx, offer)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
func (r *OfferRepo) GetOffers(ctx context.Context) ([]entity.Offer, error) {
	return make([]entity.Offer, 0), nil
}
func (r *OfferRepo) GetOffer(offer_id int) (entity.Offer, error) {
	return entity.Offer{}, nil
}
