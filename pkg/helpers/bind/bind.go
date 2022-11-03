package bind

import (
	"UsersBalanceWorker/entities"
	"errors"
)

var DataError = errors.New("not enough or wrong data")

func IsBindedCredit(binded entities.Credit) error {
	def := entities.Credit{}
	if binded.UserID == def.UserID ||
		binded.Username == def.Username ||
		binded.Value == def.Value {
		return DataError
	}

	return nil
}

func IsBindedBalance(binded entities.Balance) error {
	def := entities.Balance{}
	if binded.UserID == def.UserID {
		return DataError
	}

	return nil
}

func IsBindedService(binded entities.Service) error {
	def := entities.Service{}
	if binded.UserID == def.UserID ||
		binded.ServiceID == def.ServiceID ||
		binded.OrderID == def.OrderID ||
		binded.Cost == def.Cost {
		return DataError
	}

	return nil
}
