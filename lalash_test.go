package lalash

import (
	"context"
	"testing"
)

func TestEval(t *testing.T) {
	tests := []struct {
		name string
		expr string
		err  error
	}{
		{
			name: "echo",
			expr: "echo abc",
			err:  nil,
		},
	}
	for _, tt := range tests {
		cmd := eval.Command(cmdNew())

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		t.Run(tt.name, func(t *testing.T) {
			if err := Eval(ctx, cmd, tt.expr); err != tt.err {
				t.Errorf(err.Error())
			}
		})
	}
}
