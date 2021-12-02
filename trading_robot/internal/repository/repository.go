package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/keruch/tfs-go-hw/trading_robot/internal/domain"
	"github.com/keruch/tfs-go-hw/trading_robot/pkg/log"
)

type PostgreSQLPool struct {
	pool   *pgxpool.Pool
	logger *log.Logger
}

func NewPostgreSQLPool(url string, logger *log.Logger) (*PostgreSQLPool, error) {
	pool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &PostgreSQLPool{
		pool:   pool,
		logger: logger,
	}, nil
}

const insertTemplate = `insert into orders 
(order_id, TS, order_type, symbol, status, side, quantity, price) 
values ('%s', '%s', '%s', '%s', '%s', '%s', %v,  %v);`

func (p *PostgreSQLPool) StoreToDB(ctx context.Context, r domain.CreateOrderResponse) error {
	insertCommand := fmt.Sprintf(insertTemplate,
		r.OrderID, r.ReceivedTime, r.OrderType, r.Symbol, r.Status, r.Side, r.Size, r.LimitPrice)
	if _, err := p.pool.Query(ctx, insertCommand); err != nil {
		return err
	}
	return nil
}
