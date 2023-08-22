package json

import "time"

const (
	StatRegistered = "REGISTERED"
	StatNew        = "New"
	StatProcessed  = "PROCESSED"
	StatInvalid    = "INVALID"
	StatProcessing = "PROCESSING"
)

type User struct {
	UserName string `json:"login"`
	Pass     string `json:"password"`
}

type Accrual struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type Orders struct {
	UserName string    `json:"-"`
	Number   string    `json:"number"`
	Status   string    `json:"status"`
	Accrual  float64   `json:"accrual,omitempty"`
	DateStr  string    `json:"uploaded_at"`
	Date     time.Time `json:"-"`
}

type Withdraws struct {
	UserName string    `json:"-"`
	Order    string    `json:"order"`
	Sum      float64   `json:"accrual"`
	DateStr  string    `json:"processed_at,omitempty"`
	Date     time.Time `json:"-"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}
