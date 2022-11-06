package balanceWorker

import (
	"UsersBalanceWorker/entities"
	"UsersBalanceWorker/internal/db"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const timeZone = "+03"

type DB struct {
	Users    db.UsersInstance
	Services db.ServicesInstance
	Record   db.RecordInstance
}

func Credit(cr entities.Credit, db *DB) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	err := db.Users.UpdateUserBalance(cr)
	if err != nil {
		return err
	}

	err = db.Record.CreditRecord(cr, "ok")
	if err != nil {
		return err
	}

	return nil
}

func Balance(b entities.Balance, db *DB) (float64, error) {
	if err := db.Users.Ping(); err != nil {
		return 0, err
	}
	return db.Users.Balance(b)
}

func Service(s entities.Service, db *DB) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	err := db.Record.ServiceRecord(s, "in process")
	if err != nil {
		return err
	}

	err = db.Users.UpdateUserBalance(entities.Credit{
		UserID:   s.UserID,
		Username: "",
		Value:    -s.Cost,
	})
	if err != nil {
		_, _ = db.Record.UpdateServiceRecord(s.OrderID, "failed")
		return err
	}

	err = db.Users.UpdateUserBalance(entities.Credit{
		UserID:   1,
		Username: "",
		Value:    s.Cost,
	})
	if err != nil {
		_ = db.Users.UpdateUserBalance(entities.Credit{
			UserID:   s.UserID,
			Username: "",
			Value:    s.Cost,
		})

		_, _ = db.Record.UpdateServiceRecord(s.OrderID, "failed")
		return err
	}

	comment := fmt.Sprintf("payment for service #%d", s.ServiceID)
	err = db.Record.TransferRecord(entities.Transfer{
		UserFromID: s.UserID,
		UserToID:   1,
		Value:      s.Cost,
	}, comment, "ok")
	if err != nil {
		log.Println(err)
	}

	return nil
}

func OrderStatus(ord entities.OrderStatus, db *DB) error {
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
		if err := db.Users.UpdateUserBalance(entities.Credit{
			UserID:   1,
			Username: "",
			Value:    -upd.Value,
		}); err != nil {
			return err
		}
		// transfer money to companies bill
	default:
		if err := db.Users.UpdateUserBalance(entities.Credit{
			UserID:   upd.UserID,
			Username: "",
			Value:    upd.Value,
		}); err != nil {
			return err
		}

		if err := db.Users.UpdateUserBalance(entities.Credit{
			UserID:   1,
			Username: "",
			Value:    -upd.Value,
		}); err != nil {
			return err
		}

		comment := fmt.Sprintf("refound from service #%d", upd.ServiceID)
		err = db.Record.TransferRecord(entities.Transfer{
			UserFromID: 1,
			UserToID:   upd.UserID,
			Value:      upd.Value,
		}, comment, "ok")
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

func Transfer(t entities.Transfer, db *DB) error {
	if err := db.Users.Ping(); err != nil {
		return err
	}

	if err := db.Record.Ping(); err != nil {
		return err
	}

	if err := db.Users.UpdateUserBalance(entities.Credit{
		UserID:   t.UserFromID,
		Username: "",
		Value:    -t.Value,
	}); err != nil {
		return err
	}

	if err := db.Users.UpdateUserBalance(entities.Credit{
		UserID:   t.UserToID,
		Username: "",
		Value:    t.Value,
	}); err != nil {
		return err
	}

	if err := db.Record.TransferRecord(t, "", "ok"); err != nil {
		return err
	}

	return nil
}

func Record(r entities.Record, db *DB) (string, error) {
	if err := db.Services.Ping(); err != nil {
		log.Println(err)
	}

	if err := db.Record.Ping(); err != nil {
		return "", err
	}

	records, err := db.Record.Accounting(r.From+timeZone, r.To+timeZone)
	if err != nil {
		return "", err
	}

	services, err := db.Services.Services()
	if err != nil {
		log.Println(err)
	}

	revenue := make(map[int]float64, 0)

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
	if err != nil {
		return "", err
	}

	suf := len(files)
	filePath := filepath.Join("records", "record"+strconv.Itoa(suf)+".csv")

	f, err := os.Create(filePath)
	defer f.Close()
	if err != nil {
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

func History(h entities.History, db *DB) ([]entities.Operation, error) {
	if err := db.Users.Ping(); err != nil {
		log.Println(err)
	}

	if err := db.Services.Ping(); err != nil {
		log.Println(err)
	}

	if err := db.Record.Ping(); err != nil {
		return nil, err
	}

	creditHistory, err := db.Record.CreditHistory(h.UserID)
	if err != nil {
		return nil, err
	}

	transferTo, transferFrom, err := db.Record.TransferHistory(h.UserID)
	if err != nil {
		return nil, err
	}

	history := make([]entities.Operation, 0)

	for _, cr := range creditHistory {
		op := entities.Operation{
			Value:   cr.Value,
			Time:    cr.Time,
			Comment: "refilling",
		}
		history = append(history, op)
	}

	for _, t := range transferTo {
		op := entities.Operation{
			Value:   -t.Value,
			Time:    t.Time,
			Comment: t.Comment,
		}
		history = append(history, op)
	}

	for _, t := range transferFrom {
		op := entities.Operation{
			Value:   t.Value,
			Time:    t.Time,
			Comment: t.Comment,
		}
		history = append(history, op)
	}

	switch h.SortBy {
	case "value":
		sort.Slice(history, func(i, j int) bool {
			return math.Abs(history[i].Value) < math.Abs(history[j].Value) == h.Reverse
		})
	default:
		sort.Slice(history, func(i, j int) bool {
			return history[i].Time.Unix() < history[j].Time.Unix() == h.Reverse
		})
	}

	users, err := db.Users.Users()
	if err != nil {
		log.Println(err)
	}

	services, err := db.Services.Services()
	if err != nil {
		log.Println(err)
	}

	for i := range history {
		com := strings.Split(history[i].Comment, "#")

		if len(com) == 2 {
			body := strings.Trim(com[0], " ")

			id, err := strconv.Atoi(com[1])
			if err != nil {
				log.Println(err)
				continue
			}

			switch {

			case strings.Contains(body, "service"):
				if name, ok := services[id]; ok {
					history[i].Comment = fmt.Sprintf("%s: %s", body, name)
				}

			case strings.Contains(body, "user"):
				if name, ok := users[id]; ok && name != "" {
					history[i].Comment = fmt.Sprintf("%s: %s", body, name)
				}
			}
		}
	}

	return history, nil
}
