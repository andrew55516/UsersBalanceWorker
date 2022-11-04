package entities

import (
	"UsersBalanceWorker/internal/db"
	"fmt"
	"strconv"
)

type User struct {
	ID       int     `json:"user_id"`
	Username string  `json:"username"`
	Balance  float64 `json:"balance"`
}

//type Service struct {
//	ID          int     `json:"service_id"`
//	ServiceName string  `json:"service_name"`
//	Cost        float64 `json:"cost"`
//}

//type Record struct {
//	ID         int     `json:"record_id"`
//	UserFromID int     `json:"user_from_id"`
//	UserToID   int     `json:"user_to_id"`
//	ServiceID  int     `json:"service_id"`
//	Value      float64 `json:"value"`
//	Status     string  `json:"status"`
//}

type Db struct {
	Users    db.UsersInstance
	Services db.ServicesInstance
	Record   db.RecordInstance
}

type Credit struct {
	UserID   int     `json:"user_id"`
	Username string  `json:"username"`
	Value    float64 `json:"value"`
}

type Balance struct {
	UserID int `json:"user_id"`
}

type Service struct {
	UserID    int     `json:"user_id"`
	ServiceID int     `json:"service_id"`
	OrderID   int     `json:"order_id"`
	Cost      float64 `json:"cost"`
}

type Transfer struct {
	UserFromID int     `json:"user_from_id"`
	UserToID   int     `json:"user_to_id"`
	Value      float64 `json:"value"`
}

type OrderStatus struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

type Record struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Account struct {
	ServiceID    int
	ServiceName  string
	TotalRevenue float64
}

func (a *Account) ToSlice() []string {
	id := strconv.Itoa(a.ServiceID)
	revenue := fmt.Sprintf("%f", a.TotalRevenue)
	return []string{id, a.ServiceName, revenue}
}
