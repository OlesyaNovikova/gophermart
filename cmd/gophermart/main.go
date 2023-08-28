package main

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	h "github.com/OlesyaNovikova/gophermart/internal/handlers"
	ac "github.com/OlesyaNovikova/gophermart/internal/integrations/accruals"
	m "github.com/OlesyaNovikova/gophermart/internal/middlewares"
	p "github.com/OlesyaNovikova/gophermart/internal/store/pg"
	a "github.com/OlesyaNovikova/gophermart/internal/utils/auth"
)

func main() {
	parseFlags()
	a.InitAuth(key)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	base, err := sql.Open("pgx", dbAddr)
	if err != nil {
		panic(err)
	}
	defer base.Close()
	db, err := p.NewPostgresDB(ctx, base)
	if err != nil {
		panic(err)
	}
	//h.NewMemRepo(db)
	//временное хранилище
	//store, err := s.NewStore(dbAddr)
	//if err != nil {
	//	panic(err)
	//}
	h.InitStore(db)
	ch := ac.InitAccruals(ctx, accrualAddr, db)
	defer close(ch)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/register", m.WithLog(sugar, m.WithGzip(h.Register()))).Methods("POST")
	router.HandleFunc("/api/user/login", m.WithLog(sugar, m.WithGzip(h.Login()))).Methods("POST")
	router.HandleFunc("/api/user/orders", m.WithLog(sugar, m.WithGzip(h.OrdersPost(ch)))).Methods("POST")
	router.HandleFunc("/api/user/orders", m.WithLog(sugar, m.WithGzip(h.OrdersGet()))).Methods("GET")
	router.HandleFunc("/api/user/balance", m.WithLog(sugar, m.WithGzip(h.Balance()))).Methods("GET")
	router.HandleFunc("/api/user/balance/withdraw", m.WithLog(sugar, m.WithGzip(h.Withdraw()))).Methods("POST")
	router.HandleFunc("/api/user/withdrawals", m.WithLog(sugar, m.WithGzip(h.Withdrawals()))).Methods("Get")

	sugar.Infow("Starting server", "addr", runAddr)

	err = http.ListenAndServe(runAddr, router)
	if err != nil {
		sugar.Fatalf("start server error: %v", err)
		panic(err)
	}
}
