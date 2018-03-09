package cuckooc

import (
	"context"
	"net"
	"runtime"
	"sync"
	"testing"
)

func initSocket(config Config) (context.Context, *sync.WaitGroup, chan Executor, context.CancelFunc) {
	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())
	reqCh := make(chan Executor)
	gk := NewGatekeeper(reqCh)
	wg.Add(1)
	go gk.Start(ctx, config, wg)
	runtime.Gosched()
	return ctx, wg, reqCh, cancel
}

func testSocket(t *testing.T, c net.Conn) {
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

		{
			cmd:    "test set 1 2 3 4\ntest check 1 2 3 5",
			result: "true true true true\ntrue true true false",
		},

		{
			cmd:    "test backup ./testdata/backups-3",
			result: "true",
		},

		{
			cmd:    "test stop",
			result: "true",
		},
	}

	b := make([]byte, 1024)
	for _, s := range tests {
		_, err := c.Write([]byte(s.cmd))
		if err != nil {
			t.Fatalf("failed to write data on socket: %v", err)
		}

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

func TestUDP_integration(t *testing.T) {
	config := Config{UDP: ":5000"}
	ctx, wg, reqCh, cancel := initSocket(config)
	wg.Add(1)
	defer wg.Wait()
	defer cancel()
	go StartUDPServer(ctx, config, wg, reqCh)
	runtime.Gosched()

	c, err := net.Dial("udp", config.UDP)
	if err != nil {
		t.Fatalf("failed to initiate udp connection: %v", err)
	}

	testSocket(t, c)
}
