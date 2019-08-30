package main

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/xild/youn26b/cmd/ynab"
)

func TestExecutor_leftAntiJoin(t *testing.T) {
	type args struct {
		left  []ynab.Transaction
		right []ynab.Transaction
	}
	tests := []struct {
		args       args
		wantResult []ynab.Transaction
	}{
		{
			args: args{
				left: []ynab.Transaction{
					{
						Date:   "2019-01-01",
						Amount: -1000,
					},
					{
						Date:   "2019-01-01",
						Amount: -2000,
					},
					{
						Date:   "2019-01-01",
						Amount: -3000,
					},
				},
				right: []ynab.Transaction{
					{
						Date:   "2019-01-01",
						Amount: -2000,
					},
					{
						Date:   "2019-01-01",
						Amount: -3000,
					},
				},
			},
			wantResult: []ynab.Transaction{
				{
					Date:   "2019-01-01",
					Amount: -1000,
				},
			},
		},
	}
	for i, tt := range tests {
		t.Run("test"+strconv.Itoa(i), func(t *testing.T) {
			e := &Executor{}
			if gotResult := e.leftAntiJoin(tt.args.left, tt.args.right); !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("Executor.leftAntiJoin() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
