package repo

import (
	"context"
	"fmt"
	"strings"

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
func (r *OfferRepo) CreateOffersBach(ctx context.Context, offers []entity.Offer) error {

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

func (r *OfferRepo) CreateOffers(ctx context.Context, offers []entity.Offer) error {
	baseQuery := "insert into offers (offer_id , \"name\" , price , available ) values \n"
	builder := strings.Builder{}
	builder.WriteString(baseQuery)
	for _, offer := range offers {
		builder.WriteString(fmt.Sprintf("(%d, '%s', %f, %t),", offer.OfferId, offer.Name, offer.Price, offer.Available))
	}
	qs := builder.String()
	r.log.Info(qs)
	runeArr := []rune(builder.String())
	runeArr[len(runeArr)] = ';'

	_, err := r.client.Pool.Exec(ctx, string(runeArr), nil)
	if err != nil {
		err = fmt.Errorf("error write to db %w", err)
		r.log.Error(err.Error())
		return err
	}
	return nil
}

func (r *OfferRepo) CreateOffersAsync(ctx context.Context, in <-chan entity.Offer) error {
	baseQuery := "insert into offers (offer_id , \"name\" , price , available ) values \n"
	builder := strings.Builder{}
	builder.WriteString(baseQuery)
	for offer := range in {
		builder.WriteString(fmt.Sprintf("(%d, '%s', %f, %t),", offer.OfferId, offer.Name, offer.Price, offer.Available))
	}
	qs := builder.String()
	r.log.Info(qs)
	runeArr := []rune(builder.String())
	runeArr[len(runeArr)] = ';'

	_, err := r.client.Pool.Exec(ctx, string(runeArr), nil)
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

func (r *OfferRepo) ReadAsync(ctx context.Context) (<-chan entity.Offer, <-chan error) {
	in := make(chan entity.Offer)
	doneCh := make(chan error)

	query := "select offer_id, \"name\", price, available  from offers"

	go func() {
		defer close(in)
		defer close(doneCh)

		rows, err := r.client.Pool.Query(ctx, query)
		if err != nil {
			doneCh <- fmt.Errorf("error query %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var offer entity.Offer
			if err := rows.Scan(&offer.OfferId, &offer.Name, &offer.Price, &offer.Available); err != nil {
				doneCh <- fmt.Errorf("error row scan %w", err)
			}
			in <- offer
		}
		doneCh <- nil
	}()

	return in, doneCh
}

func (r *OfferRepo) GetOffer(ctx context.Context, offer_id int) (*entity.Offer, error) {
	query := fmt.Sprintf("select o.offer_id, o.\"name\", o.price, o.available  from offers o where o.offer_id = %d", offer_id)
	o := &entity.Offer{}
	err := r.client.Pool.QueryRow(ctx, query).Scan(o.OfferId, o.Name, o.Price, o.Available)
	if err != nil {
		return nil, fmt.Errorf("error select offer_id = %d %w", offer_id, err)
	}
	return o, nil
}
