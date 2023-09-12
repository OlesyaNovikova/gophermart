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
			`CREATE TABLE IF NOT EXISTS users("user_name" varchar(50) UNIQUE,"pass" bytea)`)
		return err
	})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы users: %v \n", err)
		return nil, err
	}

	err = retry(ctx, func(ctx context.Context) error {
		_, err = db.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS orders
			("number" varchar(256) UNIQUE,
			"user_name" varchar(50),
			"status" varchar(12),
			"accrual" double precision,
			"date_str" varchar(30),
			"date" date)`)
		return err
	})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы orders: %v \n", err)
		return nil, err
	}

	err = retry(ctx, func(ctx context.Context) error {
		_, err = db.ExecContext(ctx,
			`CREATE TABLE IF NOT EXISTS withdraws 
			("number" varchar(256) UNIQUE,
			"user_name" varchar(50),
			"summa" double precision,
			"date_str" varchar(30),
			"date" date)`)
		return err
	})
	if err != nil {
		fmt.Printf("Ошибка создания таблицы withdraws: %v \n", err)
		return nil, err
	}

	pdb := PostgresDB{
		db: db,
	}
	return &pdb, nil
}

func (p *PostgresDB) AddUser(ctx context.Context, userName string, pass []byte) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO users (user_name, pass)
		VALUES($1,$2)`, userName, pass)
	return err
}

func (p *PostgresDB) GetPass(ctx context.Context, userName string) ([]byte, error) {
	row := p.db.QueryRowContext(ctx,
		"SELECT pass FROM users WHERE user_name = $1", userName)
	var pass []byte
	err := row.Scan(&pass)
	if err != nil {
		return nil, err
	}
	return pass, nil
}

func (p *PostgresDB) AddOrder(ctx context.Context, order j.Orders) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO orders (number,user_name,status,accrual,date_str,date)
		VALUES($1,$2,$3,$4,$5,$6)`,
		order.Number, order.UserName, order.Status, order.Accrual, order.DateStr, order.Date)
	return err
}

func (p *PostgresDB) GetOrder(ctx context.Context, number string) (j.Orders, error) {
	row := p.db.QueryRowContext(ctx,
		"SELECT * FROM orders WHERE number = $1", number)
	var or j.Orders
	err := row.Scan(&or.Number, &or.UserName, &or.Status, &or.Accrual, &or.DateStr, &or.Date)
	if err != nil {
		fmt.Println(err)
		return j.Orders{}, err
	}
	return or, nil
}

func (p *PostgresDB) GetOrders(ctx context.Context, userName string) ([]j.Orders, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT * FROM orders WHERE user_name=$1 ORDER BY date", userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userOrders []j.Orders
	for rows.Next() {
		var or j.Orders
		err = rows.Scan(&or.Number, &or.UserName, &or.Status, &or.Accrual, &or.DateStr, &or.Date)
		if err != nil {
			return nil, err
		}
		userOrders = append(userOrders, or)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return userOrders, nil
}

func (p *PostgresDB) UpdateOrder(ctx context.Context, order j.Orders) error {
	_, err := p.db.ExecContext(ctx,
		`UPDATE orders SET status=$1,accrual=$2,date_str=$3,date=$4 WHERE number =$5`,
		order.Status, order.Accrual, order.DateStr, order.Date, order.Number)
	return err
}

func (p *PostgresDB) GetOrdersForUpd(ctx context.Context) ([]j.Orders, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT * FROM orders WHERE status NOT IN ($1,$2)", j.StatInvalid, j.StatProcessed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updOrders []j.Orders
	for rows.Next() {
		var or j.Orders
		err = rows.Scan(&or.Number, &or.UserName, &or.Status, &or.Accrual, &or.DateStr, &or.Date)
		if err != nil {
			return nil, err
		}
		updOrders = append(updOrders, or)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return updOrders, nil
}

func (p *PostgresDB) GetBalance(ctx context.Context, userName string) (j.Balance, error) {
	var balance j.Balance
	rows, err := p.db.QueryContext(ctx, "SELECT accrual FROM orders WHERE user_name=$1", userName)
	if err != nil {
		return balance, err
	}
	defer rows.Close()
	var accrual float64
	for rows.Next() {
		var acc float64
		err = rows.Scan(&acc)
		if err != nil {
			return balance, err
		}
		accrual += acc
	}
	err = rows.Err()
	if err != nil {
		return balance, err
	}
	balance.Current = accrual

	rows, err = p.db.QueryContext(ctx, "SELECT summa FROM withdraws WHERE user_name=$1", userName)
	if err != nil {
		return balance, err
	}
	defer rows.Close()
	var withdraw float64
	for rows.Next() {
		var wd float64
		err = rows.Scan(&wd)
		if err != nil {
			return balance, err
		}
		withdraw += wd
	}
	err = rows.Err()
	if err != nil {
		return balance, err
	}

	balance.Current -= withdraw
	balance.Withdrawn = withdraw
	return balance, nil
}

func (p *PostgresDB) AddWithdraw(ctx context.Context, withdraw j.Withdraws) error {
	_, err := p.db.ExecContext(ctx,
		`INSERT INTO withdraws(number,user_name,summa,date_str,date)
		VALUES($1,$2,$3,$4,$5)`,
		withdraw.Order, withdraw.UserName, withdraw.Sum, withdraw.DateStr, withdraw.Date)
	return err
}

func (p *PostgresDB) GetWithdraws(ctx context.Context, userName string) ([]j.Withdraws, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT * FROM withdraws WHERE user_name=$1 ORDER BY date", userName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var userWithdraws []j.Withdraws
	for rows.Next() {
		var wd j.Withdraws
		err = rows.Scan(&wd.Order, &wd.UserName, &wd.Sum, &wd.DateStr, &wd.Date)
		if err != nil {
			return nil, err
		}
		userWithdraws = append(userWithdraws, wd)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return userWithdraws, nil
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
