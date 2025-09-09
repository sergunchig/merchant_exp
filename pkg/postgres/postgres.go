package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	defaultMaxPoolSize  = 1
	defaultConnAttempts = 10
	defaultConnTimeOut  = time.Second
)

type Postgress struct {
	maxPoolSize  int
	connAttempts int
	connTimeOut  time.Duration

	Pool *pgxpool.Pool
}

func MustInitPg(url string) *Postgress {
	db, err := New(url)
	if err != nil {
		panic(err)
	}
	return db
}

func New(url string) (*Postgress, error) {

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("postgress pool create error %w", err)
	}
	pg := &Postgress{
		maxPoolSize:  defaultMaxPoolSize,
		connAttempts: defaultConnAttempts,
		connTimeOut:  defaultConnTimeOut,
	}
	pg.Pool = pool
	return pg, nil
}

func (p *Postgress) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
