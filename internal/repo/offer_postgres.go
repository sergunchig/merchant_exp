package repo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/sergunchig/merchant_exp.git/internal/entity"
	"github.com/sergunchig/merchant_exp.git/pkg/logger"
	"github.com/sergunchig/merchant_exp.git/pkg/postgres"
)

type OfferRepo struct {
	client *postgres.Postgress
	log    *logger.AppLogger
}

func New(pg *postgres.Postgress, log *logger.AppLogger) *OfferRepo {
	return &OfferRepo{client: pg, log: log}
}

func (r *OfferRepo) Create(ctx context.Context, offer entity.Offer) error {
	// query := "insert into offers (offer_id , \"name\" , price , available ) values (@offer_id, @name, @price, @available)"
	// args := pgx.NamedArgs{
	// 	"offer_id":  offer.OfferId,
	// 	"name":      offer.Name,
	// 	"price":     offer.Price,
	// 	"available": offer.Available,
	// }

	query := fmt.Sprintf("insert into offers (offer_id , \"name\" , price , available ) values (%d, %s, %f, %t)", offer.OfferId, offer.Name, offer.Price, offer.Available)

	_, err := r.client.Pool.Exec(ctx, query, nil)
	if err != nil {
		return fmt.Errorf("error insert offer %d in db, %w", offer.OfferId, err)
	}
	return nil
}
func (r *OfferRepo) CreateOffers(ctx context.Context, offers []entity.Offer) error {

	bach := pgx.Batch{}

	for _, offer := range offers {
		query := fmt.Sprintf("insert into offers (offer_id , \"name\" , price , available ) values (%d, '%s', %f, %t)", offer.OfferId, offer.Name, offer.Price, offer.Available)
		r.log.Info(fmt.Sprintf("add to bach %s", query))
		bach.Queue(query)
	}

	tx, err := r.client.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("Create offers error %w", err)
	}

	defer func(err error) {
		if err != nil {
			tx.Rollback(ctx)
		}
	}(err)

	result := tx.SendBatch(ctx, &bach)
	defer result.Close()

	for i := 0; i < len(offers); i++ {
		_, err := result.Exec()
		if err != nil {
			r.log.Error(fmt.Errorf("error executing %w", err).Error())
		}
	}

	if err := tx.Commit(ctx); err != nil {
		r.log.Error(fmt.Errorf("error committing %w", err).Error())
	}

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
