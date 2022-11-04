package balanceWorker

import (
	"UsersBalanceWorker/entities"
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func Credit(cr entities.Credit, db *entities.Db) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	err := db.Users.UpdateUserBalance(cr.UserID, cr.Username, cr.Value)
	if err != nil {
		return err
	}

	err = db.Record.CreditRecord(cr.UserID, cr.Value, "ok")
	if err != nil {
		return err
	}

	return nil
}

func Balance(b entities.Balance, db *entities.Db) (float64, error) {
	if err := db.Users.Ping(); err != nil {
		return 0, err
	}
	return db.Users.Balance(b.UserID)
}

func Service(s entities.Service, db *entities.Db) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	err := db.Record.ServiceRecord(s.OrderID, s.ServiceID, s.UserID, s.Cost, "in process")
	if err != nil {
		return err
	}

	err = db.Users.UpdateUserBalance(s.UserID, "", -s.Cost)
	if err != nil {
		_, _ = db.Record.UpdateServiceRecord(s.OrderID, "failed")
		return err
	}

	err = db.Users.UpdateUserBalance(1, "", s.Cost)
	if err != nil {
		_, _ = db.Record.UpdateServiceRecord(s.OrderID, "failed")
		return err
	}

	return nil
}

func OrderStatus(ord entities.OrderStatus, db *entities.Db) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	upd, err := db.Record.UpdateServiceRecord(ord.OrderID, ord.Status)
	if err != nil {
		return err
	}
	switch ord.Status {
	case "ok":
		if err := db.Users.UpdateUserBalance(1, "", -upd.Value); err != nil {
			return err
		}
		// transfer money to companies bill
	default:
		if err := db.Users.UpdateUserBalance(upd.UserID, "", upd.Value); err != nil {
			return err
		}

		if err := db.Users.UpdateUserBalance(1, "", -upd.Value); err != nil {
			return err
		}
	}

	return nil
}

func Transfer(t entities.Transfer, db *entities.Db) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	if err := db.Users.UpdateUserBalance(t.UserFromID, "", -t.Value); err != nil {
		return err
	}

	if err := db.Users.UpdateUserBalance(t.UserToID, "", t.Value); err != nil {
		return err
	}

	if err := db.Record.TransferRecord(t.UserFromID, t.UserToID, t.Value, "ok"); err != nil {
		return err
	}

	return nil
}

func Record(r entities.Record, db *entities.Db) (string, error) {
	if err := db.Services.Ping(); err != nil {
		return "", err
	}

	if err := db.Record.Ping(); err != nil {
		return "", err
	}

	records, err := db.Record.Accounting(r.From+"+03", r.To+"+03")
	if err != nil {
		return "", err
	}

	services, err := db.Services.Services()
	if err != nil {
		return "", err
	}

	revenue := make(map[int]float64, 0)

	//for id, s := range serv {
	//	revenue[id] = entities.Account{
	//		ServiceID:    id,
	//		ServiceName:  s,
	//		TotalRevenue: 0,
	//	}
	//}

	for _, rec := range records {
		if _, ok := revenue[rec.ServiceID]; !ok {
			revenue[rec.ServiceID] = rec.Value
		} else {
			revenue[rec.ServiceID] += rec.Value
		}
	}

	account := make([]entities.Account, 0)

	for id, r := range revenue {
		var name string

		if _, ok := services[id]; !ok {
			name = "unknown service"
		} else {
			name = services[id]
		}

		account = append(account, entities.Account{
			ServiceID:    id,
			ServiceName:  name,
			TotalRevenue: r,
		})
	}

	sort.Slice(account, func(i, j int) bool {
		return account[i].ServiceID < account[j].ServiceID
	})

	files, err := os.ReadDir("records")
	log.Println(len(files))
	if err != nil {
		log.Println(err)
		//return "", err
	}

	suf := len(files)
	filePath := filepath.Join("records", "record"+strconv.Itoa(suf)+".csv")
	log.Println(filePath)
	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
		log.Println(err)
		return "", err
	}

	w := csv.NewWriter(f)
	defer w.Flush()

	headers := []string{"SERVICE_ID", "SERVICE_NAME", "TOTAL_REVENUE"}
	if err := w.Write(headers); err != nil {
		return "", err
	}

	for _, a := range account {
		if err := w.Write(a.ToSlice()); err != nil {
			return "", err
		}
	}

	filePath = strings.Replace(filePath, "\\", "/", -1)

	return filePath, nil
}
