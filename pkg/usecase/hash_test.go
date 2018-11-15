package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/xdefrag/hash-ketchum/pkg/adapter/mock"

	"github.com/xdefrag/hash-ketchum/pkg/types"
)

type storage func(*mock.MockHashStorager) *mock.MockHashStorager

var (
	hashGood = types.Hash{
		Hash:      "0054478C35F2708C5D0BF28696B44F1BCF79832BF716A2BFBA665212BA9B4F09",
		Login:     "login",
		Timestamp: time.Now().Unix(),
	}

	hashBad = types.Hash{
		Hash:      "B454478C35F2708C5D0BF28696B44F1BCF79832BF716A2BFBA665212BA9B4F09",
		Login:     "login",
		Timestamp: time.Now().Unix(),
	}
)

func TestHash_Store(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	type args struct {
		ctx  context.Context
		hash types.Hash
	}
	tests := []struct {
		name    string
		storage storage
		args    args
		wantErr bool
	}{
		{
			name: "Should store valid hash",
			storage: func(mock *mock.MockHashStorager) *mock.MockHashStorager {
				mock.EXPECT().Store(context.TODO(), hashGood).Return(nil)
				return mock
			},
			args: args{
				ctx:  context.TODO(),
				hash: hashGood,
			},
			wantErr: false,
		},
		{
			name: "Should not store invalid hash and return error",
			storage: func(mock *mock.MockHashStorager) *mock.MockHashStorager {
				return mock
			},
			args: args{
				ctx:  context.TODO(),
				hash: hashBad,
			},
			wantErr: true,
		},
		{
			name: "Should return error if store not available",
			storage: func(mock *mock.MockHashStorager) *mock.MockHashStorager {
				mock.EXPECT().Store(context.TODO(), hashGood).Return(errors.New("Store not available"))
				return mock
			},
			args: args{
				ctx:  context.TODO(),
				hash: hashGood,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := mock.NewMockHashStorager(mockCtrl)

			h := NewHash(tt.storage(storage))
			if err := h.Store(tt.args.ctx, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Hash.Store() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
