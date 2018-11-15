// +build integration

package redis

import (
	"context"
	"testing"
	"time"

	"github.com/xdefrag/hash-ketchum/pkg/types"
)

var (
	cfg = Config{
		Host: "0.0.0.0",
		Port: 6379,
	}
	hash = types.Hash{
		Login:     "brock",
		Hash:      "0054478C35F2708C5D0BF28696B44F1BCF79832BF716A2BFBA665212BA9B4F09",
		Timestamp: time.Now().Unix(),
	}
)

func TestRedis_Store(t *testing.T) {
	r := New(cfg, nil)

	type args struct {
		ctx  context.Context
		hash types.Hash
	}
	tests := []struct {
		name          string
		args          args
		isDataInRedis func() (bool, error)
		wantErr       bool
	}{
		{
			name: "Should store hash to redis",
			args: args{
				ctx:  nil,
				hash: hash,
			},
			isDataInRedis: func() (bool, error) {
				result, err := r.pool.Get().Do("HGETALL", hash.Login)
				if err != nil {
					return false, err
				}

				return len(result.([]interface{})) == 2, nil
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.Store(tt.args.ctx, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Redis.Store() error = %v, wantErr %v", err, tt.wantErr)
			}

			isDataInRedis, err := tt.isDataInRedis()
			if err != nil {
				t.Fatalf("Redis error: %s", err)
			}

			if !isDataInRedis {
				t.Errorf("Redis.Store() data is not saved to redis")
			}
		})
	}
}
