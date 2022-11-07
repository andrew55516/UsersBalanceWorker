package bind

import (
	"UsersBalanceWorker/entities"
	"testing"
)

func TestRightBindedBalance(t *testing.T) {
	type args struct {
		binded entities.Balance
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "null",
			args:    args{},
			wantErr: DataError,
		},
		{
			name:    "wrong ID",
			args:    args{entities.Balance{UserID: 1}},
			wantErr: UserIDError,
		},
		{
			name:    "right binded",
			args:    args{entities.Balance{UserID: 123}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedBalance(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedBalance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRightBindedCredit(t *testing.T) {
	type args struct {
		binded entities.Credit
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "null ID",
			args: args{entities.Credit{
				Username: "user",
				Value:    123,
			}},
			wantErr: DataError,
		},
		{
			name: "null Value",
			args: args{entities.Credit{
				UserID:   123,
				Username: "user",
			}},
			wantErr: DataError,
		},
		{
			name: "null Username",
			args: args{entities.Credit{
				UserID: 123,
				Value:  123,
			}},
			wantErr: nil,
		},
		{
			name: "wrong UserID",
			args: args{entities.Credit{
				UserID:   1,
				Username: "",
				Value:    123,
			}},
			wantErr: UserIDError,
		},
		{
			name: "wrong Value",
			args: args{entities.Credit{
				UserID:   123,
				Username: "",
				Value:    -123,
			}},
			wantErr: ValueError,
		},
		{
			name: "right binded",
			args: args{entities.Credit{
				UserID:   123,
				Username: "",
				Value:    123,
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedCredit(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedCredit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRightBindedHistory(t *testing.T) {
	type args struct {
		binded entities.History
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "null",
			args:    args{},
			wantErr: DataError,
		},
		{
			name:    "wrong ID",
			args:    args{entities.History{UserID: 1}},
			wantErr: UserIDError,
		},
		{
			name: "right binded",
			args: args{entities.History{
				UserID:  123,
				SortBy:  "",
				Reverse: false,
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedHistory(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRightBindedRecord(t *testing.T) {
	type args struct {
		binded entities.Record
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "null from",
			args:    args{entities.Record{To: "2022-01-01"}},
			wantErr: DataError,
		},
		{
			name:    "null to",
			args:    args{entities.Record{From: "2022-01-01"}},
			wantErr: DataError,
		},
		{
			name: "wrong from",
			args: args{entities.Record{
				From: "123131313",
				To:   "2022-01-01",
			}},
			wantErr: DateError,
		},
		{
			name: "wrong to",
			args: args{entities.Record{
				From: "2022-01-01",
				To:   "123132132",
			}},
			wantErr: DateError,
		},
		{
			name: "right binded",
			args: args{entities.Record{
				From: "2022-01-01",
				To:   "2022-01-01",
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedRecord(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedRecord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRightBindedService(t *testing.T) {
	type args struct {
		binded entities.Service
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "null UserID",
			args: args{entities.Service{
				ServiceID: 1,
				OrderID:   1,
				Cost:      1,
			}},
			wantErr: DataError,
		},
		{
			name: "null ServiceID",
			args: args{entities.Service{
				UserID:  123,
				OrderID: 1,
				Cost:    1,
			}},
			wantErr: DataError,
		},
		{
			name: "null OrderID",
			args: args{entities.Service{
				UserID:    123,
				ServiceID: 1,
				Cost:      1,
			}},
			wantErr: DataError,
		},
		{
			name: "null Cost",
			args: args{entities.Service{
				UserID:    123,
				ServiceID: 1,
				OrderID:   1,
			}},
			wantErr: DataError,
		},
		{
			name: "wrong UserID",
			args: args{entities.Service{
				UserID:    1,
				ServiceID: 1,
				OrderID:   1,
				Cost:      1,
			}},
			wantErr: UserIDError,
		},
		{
			name: "wrong ServiceID",
			args: args{entities.Service{
				UserID:    123,
				ServiceID: -1,
				OrderID:   1,
				Cost:      1,
			}},
			wantErr: ServiceIDError,
		},
		{
			name: "wrong OrderID",
			args: args{entities.Service{
				UserID:    123,
				ServiceID: 1,
				OrderID:   -1,
				Cost:      1,
			}},
			wantErr: OrderIDError,
		},
		{
			name: "wrong Cost",
			args: args{entities.Service{
				UserID:    123,
				ServiceID: 1,
				OrderID:   1,
				Cost:      -1,
			}},
			wantErr: CostError,
		},
		{
			name: "right binded",
			args: args{entities.Service{
				UserID:    123,
				ServiceID: 1,
				OrderID:   1,
				Cost:      1,
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedService(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRightBindedStatus(t *testing.T) {
	type args struct {
		binded entities.OrderStatus
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "null OrderID",
			args: args{entities.OrderStatus{
				Status: "ok",
			}},
			wantErr: DataError,
		},
		{
			name: "null status",
			args: args{entities.OrderStatus{
				OrderID: 1,
			}},
			wantErr: DataError,
		},
		{
			name: "wrong OrderID",
			args: args{entities.OrderStatus{
				OrderID: -1,
				Status:  "ok",
			}},
			wantErr: OrderIDError,
		},
		{
			name: "wrong Status",
			args: args{entities.OrderStatus{
				OrderID: 1,
				Status:  "wrong status",
			}},
			wantErr: StatusError,
		},
		{
			name: "right binded",
			args: args{entities.OrderStatus{
				OrderID: 1,
				Status:  "ok",
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedStatus(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRightBindedTransfer(t *testing.T) {
	type args struct {
		binded entities.Transfer
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "null UserFromID",
			args: args{entities.Transfer{
				UserToID: 123,
				Value:    1,
			}},
			wantErr: DataError,
		},
		{
			name: "null UserToID",
			args: args{entities.Transfer{
				UserFromID: 123,
				Value:      1,
			}},
			wantErr: DataError,
		},
		{
			name: "null Value",
			args: args{entities.Transfer{
				UserFromID: 123,
				UserToID:   122,
			}},
			wantErr: DataError,
		},
		{
			name: "wrong Value",
			args: args{entities.Transfer{
				UserFromID: 123,
				UserToID:   122,
				Value:      -1,
			}},
			wantErr: ValueError,
		},
		{
			name: "wrong UserFromID",
			args: args{entities.Transfer{
				UserFromID: 1,
				UserToID:   122,
				Value:      1,
			}},
			wantErr: UserIDError,
		},
		{
			name: "wrong UserToID",
			args: args{entities.Transfer{
				UserFromID: 123,
				UserToID:   1,
				Value:      1,
			}},
			wantErr: UserIDError,
		},
		{
			name: "same users",
			args: args{entities.Transfer{
				UserFromID: 123,
				UserToID:   123,
				Value:      1,
			}},
			wantErr: SameIDError,
		},
		{
			name: "right binded",
			args: args{entities.Transfer{
				UserFromID: 123,
				UserToID:   122,
				Value:      1,
			}},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RightBindedTransfer(tt.args.binded); err != tt.wantErr {
				t.Errorf("RightBindedTransfer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
