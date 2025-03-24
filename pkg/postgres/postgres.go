package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	_defaultMaxPoolSize  = 1
	_defaultConnAttempts = 10
	_defaultConnTimeOut  = time.Second
)

type Postgress struct {
	maxPoolSize  int
	connAttempts int
	connTimeOut  time.Duration

	Pool *pgxpool.Pool
}

func New(url string) (*Postgress, error) {

	pool, err := pgxpool.New(context.Background(), url)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	pg := &Postgress{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeOut:  _defaultConnTimeOut,
	}
	pg.Pool = pool
	return pg, nil
}

func (p *Postgress) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
