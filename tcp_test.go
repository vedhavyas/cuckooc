package cuckooc

import (
	"context"
	"net"
	"runtime"
	"sync"
	"testing"
)

func TestTCP_integration(t *testing.T) {
	wg := new(sync.WaitGroup)
	defer wg.Wait()
	config := Config{TCP: ":4000"}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reqCh := make(chan Executor)
	gk := NewGatekeeper(reqCh)
	wg.Add(2)
	go gk.Start(ctx, config, wg)
	runtime.Gosched()
	go StartTCPServer(ctx, config, wg, reqCh)
	runtime.Gosched()

	tests := []struct {
		cmd    string
		result string
	}{
		{
			cmd:    "test new",
			result: "true",
		},

		{
			cmd:    "test setu a b c d e f g",
			result: "true true true true true true true",
		},

		{
			cmd:    "test check a b 1 2 c d e f g",
			result: "true true false false true true true true true",
		},
	}

	c, err := net.Dial("tcp", config.TCP)
	if err != nil {
		t.Fatalf("failed to initiate tcp connection: %v", err)
	}

	for _, s := range tests {
		_, err := c.Write([]byte(s.cmd))
		if err != nil {
			t.Fatalf("failed to write data on socket: %v", err)
		}

		b := make([]byte, 1024)
		n, err := c.Read(b)
		if err != nil {
			t.Fatalf("failed to read data from socket :%v", err)
		}

		res := string(b[:n])
		if s.result != res {
			t.Fatalf("expected %s but got %s", s.result, res)
		}
	}

	c.Close()
}
