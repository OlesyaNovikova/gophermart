package main

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	h "github.com/OlesyaNovikova/gophermart/internal/handlers"
	ac "github.com/OlesyaNovikova/gophermart/internal/integrations/accruals"
	m "github.com/OlesyaNovikova/gophermart/internal/middlewares"
	a "github.com/OlesyaNovikova/gophermart/internal/models/auth"
	s "github.com/OlesyaNovikova/gophermart/internal/store"
)

func main() {
	parseFlags()

	if envKey := os.Getenv("KEY"); envKey != "" {
		a.InitAuth(envKey)
	} else {
		a.InitAuth("default")
	}
	//временное хранилище
	store, err := s.NewStore(dbAddr)
	if err != nil {
		panic(err)
	}
	h.InitStore(&store)
	ac.InitAccruals(accrualAddr, &store)

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()
	sugar := *logger.Sugar()

	router := mux.NewRouter()
	router.HandleFunc("/api/user/register", m.WithLog(sugar, m.WithGzip(h.Register()))).Methods("POST")
	router.HandleFunc("/api/user/login", m.WithLog(sugar, m.WithGzip(h.Login()))).Methods("POST")
	router.HandleFunc("/api/user/orders", m.WithLog(sugar, m.WithGzip(h.OrdersPost()))).Methods("POST")
	router.HandleFunc("/api/user/orders", m.WithLog(sugar, m.WithGzip(h.OrdersGet()))).Methods("GET")
	//router.HandleFunc("/api/user/balance", m.WithLog(sugar, m.WithGzip(h.Balance()))).Methods("GET")
	//router.HandleFunc("/api/user/balance/withdraw", m.WithLog(sugar, m.WithGzip(h.Withdraw()))).Methods("POST")
	//router.HandleFunc("/api/user/withdraws", m.WithLog(sugar, m.WithGzip(h.Withdraws()))).Methods("Get")

	sugar.Infow("Starting server", "addr", runAddr)

	err = http.ListenAndServe(runAddr, router)
	if err != nil {
		sugar.Fatalf("start server error: %v", err)
		panic(err)
	}
}
