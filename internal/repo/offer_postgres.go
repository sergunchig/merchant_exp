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

	// baseQuery := fmt.Sprintf("insert into offers (offer_id , \"name\" , price , available ) values (%d, %s, %f, %t)", offer.OfferId, offer.Name, offer.Price, offer.Available)
	baseQuery := "insert into offers (offer_id , \"name\" , price , available ) values ($1, $2, $3, $4)"
	_, err := r.client.Pool.Exec(ctx, baseQuery, offer.OfferId, offer.Name, offer.Price, offer.Available)
	if err != nil {
		return fmt.Errorf("error insert offer %d in db, %w", offer.OfferId, err)
	}
	return nil
}

func (r *OfferRepo) CreateOffers(ctx context.Context, offers []entity.Offer) error {
	_, err := r.client.Pool.CopyFrom(ctx, pgx.Identifier{"offers"},
		[]string{"offer_id", "name", "price", "available"},
		pgx.CopyFromSlice(len(offers), func(i int) ([]any, error) {
			o := offers[i]
			return []interface{}{o.OfferId, o.Name, o.Price, o.Available}, nil
		}))
	if err != nil {
		err = fmt.Errorf("error write to db %w", err)
		r.log.Error(err.Error())
		return err
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

func (r *OfferRepo) GetOffer(ctx context.Context, offer_id int) (*entity.Offer, error) {
	//query := fmt.Sprintf("select o.offer_id, o.\"name\", o.price, o.available  from offers o where o.offer_id = %d", offer_id)
	//query := "select o.offer_id, o.\"name\", o.price, o.available  from offers o where o.offer_id = $1"
	o := &entity.Offer{}
	row := r.client.Pool.QueryRow(ctx, "select o.offer_id, o.\"name\", o.price, o.available  from offers o where o.offer_id = $1", offer_id)
	err := row.Scan(&o.OfferId, &o.Name, &o.Price, &o.Available)
	if err != nil {
		return nil, fmt.Errorf("error select offer_id = %d %w", offer_id, err)
	}
	return o, nil
}
