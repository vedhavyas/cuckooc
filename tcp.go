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
var tcpLog = log.New(os.Stderr, "TCP", log.LstdFlags)

// readData reads the from the connection and returns it over dataCh
// if error out, sends it over errCh and returns the chan
func readData(conn net.Conn, dataCh chan<- []byte, errCh chan<- error) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			errCh <- fmt.Errorf("failed to read data from connection: %v", err)
			return
		}

		dataCh <- buf[:n]
	}
}

// handleConnection handles a new connection till the connection gets shutdown
// Additionally, we close an idle connection after idleClose time(0)
// TODO(ved): use idle close time to close the idle connection and need a way to disable this as well
// or else remove this block
func handleConnection(conn net.Conn, idleClose string, reqCh chan<- Executor) {
	dataCh, errCh, respCh := make(chan []byte), make(chan error), make(chan string)

	go readData(conn, dataCh, errCh)
	for {
		select {
		case d := <-dataCh:
			scmds := readCommands(d)
			var results []string
			for _, scmd := range scmds {
				exe, err := parseCommand(scmd, respCh)
				if err != nil {
					results = append(results, fmt.Sprintf("%s(%v)", false, err))
					continue
				}

				reqCh <- exe
				results = append(results, <-respCh)
			}

			conn.Write([]byte(strings.Join(results, "\n")))

		case err := <-errCh:
			log.Printf("read error from %s: %v\n", conn.RemoteAddr().String(), err)
			return
		}
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
func StartTCPServer(ctx context.Context, config Config, wg *sync.WaitGroup, reqCh chan<- Executor) {
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
			tcpLog.Printf("handling a new connection from %s\n", conn.RemoteAddr().String())
			go handleConnection(conn, config.TCP.IdleClose, reqCh)
		}
	}

}
