package cuckooc

import (
	"context"
	"testing"
)

func Test_filterWrapper_integration(t *testing.T) {
	tests := []struct {
		cmd    string
		result string
	}{
		{
			cmd:    "test create 30",
			result: "true",
		},

		{
			cmd:    "test set a b c d e",
			result: "true true true true true",
		},

		{
			cmd:    "test setu a b f g",
			result: "true true true true",
		},

		{
			cmd:    "test check a b 1 2 c d e f g",
			result: "true true false false true true true true true",
		},

		{
			cmd:    "test count",
			result: "7",
		},

		{
			cmd:    "test loadfactor",
			result: "0.2188",
		},

		{
			cmd:    "test delete a c d f z",
			result: "true true true true false",
		},

		{
			cmd:    "test count",
			result: "3",
		},

		{
			cmd:    "test loadfactor",
			result: "0.0938",
		},
	}

	cmdCh := make(chan Executor)
	respCh := make(chan string)
	f := newWrapper("test", cmdCh)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go f.listen(ctx, Config{})
	for _, c := range tests {
		exe, err := parseCommand(c.cmd, respCh)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		cmdCh <- exe
		result := <-respCh

		if result != c.result {
			t.Fatalf("expected result %s but got %s", c.result, result)
		}
	}
}
