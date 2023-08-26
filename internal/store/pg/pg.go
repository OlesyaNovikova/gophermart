package pg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"

	j "github.com/OlesyaNovikova/gophermart/internal/models/json"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(ctx context.Context, db *sql.DB) (*PostgresDB, error) {
	err := db.PingContext(ctx)
	if err != nil {
		fmt.Printf("Ошибка соединения с базой: %v \n", err)
		return nil, err
	}

	err = retry(ctx, func(ctx context.Context) error {
		_, err = db.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS users("user_name" varchar(50) UNIQUE,"pass" varchar(100))`)
		return err
	})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы users: %v \n", err)
		return nil, err
	}

	err = retry(ctx, func(ctx context.Context) error {
		_, err = db.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS orders("user_name" varchar(50) UNIQUE,"pass" varchar(100))`)
		return err
	})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы orders: %v \n", err)
		return nil, err
	}

	pdb := PostgresDB{
		db: db,
	}
	return &pdb, nil
}

func AddOrder(ctx context.Context, order j.Orders) error {
	return nil
}

func retry(ctx context.Context, f func(ctx context.Context) error) error {
	var err error
	delay := [3]int{1, 3, 5}
	err = f(ctx)
	if err != nil {
		for _, t := range delay {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgerrcode.IsConnectionException(pgErr.Code) {
				time.Sleep(time.Duration(t) * time.Second)
				err = f(ctx)
				if err == nil {
					return err
				}
			} else {
				return err
			}
		}
	}
	return err
}
