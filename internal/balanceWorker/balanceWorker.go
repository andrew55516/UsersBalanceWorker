package balanceWorker

import "UsersBalanceWorker/entities"

func Credit(cr entities.Credit, db *entities.Db) error {
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
	return db.Users.Balance(b.UserID)
}

func Service(s entities.Service, db *entities.Db) error {
	err := db.Record.ServiceRecord(s.OrderID, s.ServiceID, s.UserID, s.Cost, "in process")
	if err != nil {
		return err
	}

	err = db.Users.UpdateUserBalance(s.UserID, "", -s.Cost)
	if err != nil {
		db.Record.UpdateServiceRecord(s.OrderID, "failed")
		return err
	}
	// TODO: service_record status failed
	err = db.Users.UpdateUserBalance(1, "", s.Cost)
	if err != nil {
		db.Record.UpdateServiceRecord(s.OrderID, "failed")
		return err
	}

	return nil
}
