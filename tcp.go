package cuckooc

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

// tcpLog with specific prefix set
var tcpLog = log.New(os.Stderr, "TCP: ", log.LstdFlags)

// handleConnection handles a new connection
func handleConnection(conn net.Conn, reqCh chan<- Executor) {
	respCh := make(chan string)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("read error from %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}

		scmds := readCommands(buf[:n])
		var results []string
		for _, scmd := range scmds {
			exe, err := parseCommand(scmd, respCh)
			if err != nil {
				tcpLog.Printf("failed to parse command: %v", err)
				results = append(results, fmt.Sprintf("%s(%v)", notOk, err))
				continue
			}

			tcpLog.Printf("sending request to gatekeper...")
			reqCh <- exe
			res := <-respCh
			results = append(results, res)
			tcpLog.Printf("response received: %s\n", res)
		}

		n, err = conn.Write([]byte(strings.Join(results, "\n")))
		if err != nil {
			tcpLog.Printf("failed to write response to the socket: %v\n", err)
			return
		}

		tcpLog.Printf("%d bytes written to the socket\n", n)
	}
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
func StartTCPServer(ctx context.Context, config Config, wg *sync.WaitGroup, cmdCh chan<- Executor) {
	defer wg.Done()
	addr := strings.TrimSpace(config.TCP)
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
			tcpLog.Printf("handling a new connection from %s\n", conn.RemoteAddr().String())
			go handleConnection(conn, cmdCh)
		}
	}

}
