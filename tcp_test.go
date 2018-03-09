package cuckooc

import (
	"net"
	"runtime"
	"testing"
)

func TestTCP_integration(t *testing.T) {
	config := Config{TCP: ":4000"}
	ctx, wg, reqCh, cancel := initSocket(config)
	defer wg.Wait()
	defer cancel()
	wg.Add(1)
	go StartTCPServer(ctx, config, wg, reqCh)
	runtime.Gosched()

	c, err := net.Dial("tcp", config.TCP)
	if err != nil {
		t.Fatalf("failed to initiate tcp connection: %v", err)
	}

	testSocket(t, c)
}
