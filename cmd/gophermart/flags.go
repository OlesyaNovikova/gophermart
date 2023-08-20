package main

import (
	"flag"
	"os"
)

var (
	runAddr     string
	dbAddr      string
	accrualAddr string
)

// parseFlags обрабатывает аргументы командной строки
func parseFlags() {
	flag.StringVar(&runAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&dbAddr, "d", "", "data base DSN")
	flag.StringVar(&accrualAddr, "r", "localhost:8585", "address and port of accrual system")

	flag.Parse()

	if envA := os.Getenv("RUN_ADDRESS"); envA != "" {
		runAddr = envA
	}
	if envD := os.Getenv("DATABASE_URI"); envD != "" {
		dbAddr = envD
	}
	if envR := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envR != "" {
		accrualAddr = envR
	}
}
