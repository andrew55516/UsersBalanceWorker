package entities

import "UsersBalanceWorker/internal/db"

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

type Record struct {
	ID         int     `json:"record_id"`
	UserFromID int     `json:"user_from_id"`
	UserToID   int     `json:"user_to_id"`
	ServiceID  int     `json:"service_id"`
	Value      float64 `json:"value"`
	Status     string  `json:"status"`
}
