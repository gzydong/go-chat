package utils

import (
	"testing"
)

func TestMtRand(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case1",
			args: args{
				min: 0,
				max: 10,
			},
		},
		{
			name: "case2",
			args: args{
				min: 10,
				max: 30,
			},
		},
		{
			name: "case3",
			args: args{
				min: 100,
				max: 105,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MtRand(tt.args.min, tt.args.max); got < tt.args.min || got > tt.args.max {
				t.Errorf("MtRand() = %v, want %v", got, "[min, max]")
			}
		})
	}
}
