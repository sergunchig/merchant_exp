package models

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type OfferDao interface {
	Create(ctx context.Context, Offer *Offer) error
	CreateOffers(ctx context.Context, offers []Offer) error
	CreateOffersPipe(ctx context.Context, in <-chan Offer) <-chan Offer
	GetOffers(ctx context.Context) ([]Offer, error)
	GetOffer(offer_id int) (Offer, error)
}

type Dao struct {
	db *pgx.Conn
}

func NewDao(connStr string) (*Dao, error) {
	db, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, errors.New("error connection database")
	}
	return &Dao{
		db: db,
	}, nil
}

func (d *Dao) Create(ctx context.Context, offer Offer) error {
	query := "insert into offers (offer_id , \"name\" , price , available ) values (@offer_id, @name, @price, @available)"
	args := pgx.NamedArgs{
		"offer_id":  offer.OfferId,
		"name":      offer.Name,
		"price":     offer.Price,
		"available": offer.Available,
	}

	_, err := d.db.Exec(ctx, query, args)
	if err != nil {
		return fmt.Errorf("error insert offer %d in db, %w", offer.OfferId, err)
	}
	return nil
}
func (d *Dao) CreateOffers(ctx context.Context, offers []Offer) error {

	for _, offer := range offers {
		err := d.Create(ctx, offer)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (d *Dao) CreateOffersPipe(ctx context.Context, in <-chan Offer) error {

	for offer := range in {
		err := d.Create(ctx, offer)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func (d *Dao) GetOffers(ctx context.Context) ([]Offer, error) {
	return nil, nil
}
