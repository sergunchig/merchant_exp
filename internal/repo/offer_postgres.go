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

func (r *OfferRepo) Read(ctx context.Context) ([]entity.Offer, error) {
	query := "select offer_id, \"name\", price, available  from offers"

	rows, err := r.client.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error query %w", err)
	}

	defer rows.Close()
	offers := make([]entity.Offer, 0, 10000)

	for rows.Next() {
		var offer entity.Offer
		if err := rows.Scan(&offer.OfferId, &offer.Name, &offer.Price, &offer.Available); err != nil {
			return nil, fmt.Errorf("error row scan %w", err)
		}
		offers = append(offers, offer)
	}

	return offers, nil
}

func (r *OfferRepo) ReadAsync(ctx context.Context) (chan<- entity.Offer, chan<- error) {
	in := make(chan entity.Offer)
	errCh := make(chan error)

	query := "select offer_id, \"name\", price, available  from offers"

	go func() {
		rows, err := r.client.Pool.Query(ctx, query)
		if err != nil {
			errCh <- fmt.Errorf("error query %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var offer entity.Offer
			if err := rows.Scan(&offer.OfferId, &offer.Name, &offer.Price, &offer.Available); err != nil {
				errCh <- fmt.Errorf("error row scan %w", err)
			}
			in <- offer
		}
	}()

	return in, errCh
}

func (r *OfferRepo) GetOffer(offer_id int) (entity.Offer, error) {
	return entity.Offer{}, nil
}
