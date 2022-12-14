package bind

import (
	"UsersBalanceWorker/entities"
	"errors"
	"regexp"
)

var status = map[string]struct{}{
	"ok":     struct{}{},
	"failed": struct{}{},
}
var DataError = errors.New("not enough or wrong data")
var UserIDError = errors.New("wrong UserID")
var ServiceIDError = errors.New("wrong ServiceID")
var OrderIDError = errors.New("wrong OrderID")
var ValueError = errors.New("value must be positive")
var CostError = errors.New("cost must be positive")
var DateError = errors.New("wrong date-style")
var StatusError = errors.New("bad order status")
var SameIDError = errors.New("userFrom and userTo can't be the same")

func RightBindedCredit(binded entities.Credit) error {
	def := entities.Credit{}
	if binded.UserID == def.UserID ||
		//binded.Username == def.Username ||
		binded.Value == def.Value {
		return DataError
	}

	if binded.Value <= 0 {
		return ValueError
	}

	if binded.UserID <= 1 {
		return UserIDError
	}

	return nil
}

func RightBindedBalance(binded entities.Balance) error {
	def := entities.Balance{}
	if binded.UserID == def.UserID {
		return DataError
	}

	if binded.UserID <= 1 {
		return UserIDError
	}

	return nil
}

func RightBindedService(binded entities.Service) error {
	def := entities.Service{}
	if binded.UserID == def.UserID ||
		binded.ServiceID == def.ServiceID ||
		binded.OrderID == def.OrderID ||
		binded.Cost == def.Cost {
		return DataError
	}

	if binded.UserID <= 1 {
		return UserIDError
	}

	if binded.ServiceID <= 0 {
		return ServiceIDError
	}

	if binded.OrderID <= 0 {
		return OrderIDError
	}

	if binded.Cost <= 0 {
		return CostError
	}

	return nil
}

func RightBindedStatus(binded entities.OrderStatus) error {
	def := entities.OrderStatus{}
	if binded.OrderID == def.OrderID ||
		binded.Status == def.Status {
		return DataError
	}

	if _, ok := status[binded.Status]; !ok {
		return StatusError
	}

	if binded.OrderID <= 0 {
		return OrderIDError
	}

	return nil
}

func RightBindedTransfer(binded entities.Transfer) error {
	def := entities.Transfer{}
	if binded.UserFromID == def.UserFromID ||
		binded.UserToID == def.UserToID ||
		binded.Value == def.Value {
		return DataError
	}

	if binded.Value <= 0 {
		return ValueError
	}

	if binded.UserFromID <= 1 || binded.UserToID <= 1 {
		return UserIDError
	}

	if binded.UserFromID == binded.UserToID {
		return SameIDError
	}

	return nil
}

func RightBindedRecord(binded entities.Record) error {
	def := entities.Record{}
	if binded.From == def.From || binded.To == def.To {
		return DataError
	}

	re := regexp.MustCompile(`\d{4}-\d{2}-\d{2}`)

	if !(re.MatchString(binded.From) && re.MatchString(binded.To) &&
		len(binded.From) == 10 && len(binded.To) == 10) {
		return DateError
	}

	return nil
}

func RightBindedHistory(binded entities.History) error {
	def := entities.History{}
	if binded.UserID == def.UserID {
		return DataError
	}

	if binded.UserID <= 1 {
		return UserIDError
	}

	return nil
}
