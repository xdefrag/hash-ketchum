package config

import (
	"testing"
)

func TestWithDefaultString(t *testing.T) {
	type args struct {
		v string
		d string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Should return value",
			args: args{
				v: "value",
				d: "default",
			},
			want: "value",
		},
		{
			name: "Shoud return default",
			args: args{
				v: "",
				d: "default",
			},
			want: "default",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDefaultString(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("WithDefaultString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithDefaultInt(t *testing.T) {
	type args struct {
		v int
		d int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "Should return value",
			args: args{
				v: 13,
				d: 666,
			},
			want: 13,
		},
		{
			name: "Should return default",
			args: args{
				v: 0,
				d: 666,
			},
			want: 666,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithDefaultInt(tt.args.v, tt.args.d); got != tt.want {
				t.Errorf("WithDefaultInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
