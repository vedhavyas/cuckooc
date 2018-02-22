package cuckooc

import (
	"context"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// tcpLog with specific prefix set
var tcpLog = log.New(os.Stderr, "TCP", log.LstdFlags)

// handleConnection handles a new connection till the connection gets shutdown
// Additionally, we close an idle connection after idleClose time(0)
func handleConnection(conn net.Conn, idleClose string) {

}

// listen listens for the tcp connections and sends it across the channel
func listen(l net.Listener, connCh chan<- net.Conn) {
	for {
		conn, err := l.Accept()
		if err != nil {
			tcpLog.Printf("error accepting a connection: %v\n", err)
			continue
		}

		connCh <- conn
	}
}

// StartTCPServer starts a TCP server on the address provided in the configuration. If none is provided, this is a no-op
// blocking call. Should be run on a different go routine
func StartTCPServer(ctx context.Context, config Config, wg *sync.WaitGroup) {
	defer wg.Done()
	addr := strings.TrimSpace(config.TCP.Address)
	if addr == "" {
		tcpLog.Printf("no tcp address given")
		return
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		tcpLog.Fatalf("failed to start TCP server: %v", err)
		return
	}

	tcpLog.Printf("starting TCP server on %s\n", config.TCP)
	connCh := make(chan net.Conn)
	go listen(l, connCh)

	for {
		select {
		case <-ctx.Done():
			tcpLog.Println("shutting down tcp server...")
			return
		case conn := <-connCh:
			tcpLog.Println("handling a new connection...")
			go handleConnection(conn, config.TCP.IdleClose)
		}
	}

}
