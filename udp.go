package cuckooc

import (
	"context"
	"log"
	"net"
	"os"
	"sync"
)

// udpLog with specific prefix set
var udpLog = log.New(os.Stderr, "UDP: ", log.LstdFlags)

// udpMessage holds the commands and address the commands received from
type udpMessage struct {
	addr net.Addr
	cmds []string
	done bool
	resp string
}

func processUDPMessage(msg udpMessage, cmdCh chan<- Executor, result chan<- udpMessage) {
	respCh := make(chan string)
	defer close(respCh)
	msg.resp = executeMessages(msg.cmds, respCh, cmdCh, udpLog)
	msg.done = true
	result <- msg
}

func handleUDPConn(conn net.PacketConn, msgCh chan<- udpMessage) {
	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			udpLog.Printf("failed to read packet: %v\n", err)
			continue
		}

		scmds := readCommands(buf[:n])
		msgCh <- udpMessage{addr: addr, cmds: scmds}
	}
}

// StartUDPServer starts a UDP server
func StartUDPServer(ctx context.Context, config Config, wg *sync.WaitGroup, cmdCh chan<- Executor) {
	defer wg.Done()

	if config.UDP == "" {
		udpLog.Println("UDP transport disabled...")
		return
	}

	conn, err := net.ListenPacket("udp", config.UDP)
	if err != nil {
		udpLog.Fatalf("failed to start UDP server: %v\n", err)
	}

	udpLog.Printf("starting UDP server on %s\n", config.UDP)
	msgCh := make(chan udpMessage)
	go handleUDPConn(conn, msgCh)

	for {
		select {
		case <-ctx.Done():
			udpLog.Println("shutting down UDP server...")
			return
		case msg := <-msgCh:
			if msg.done {
				n, err := conn.WriteTo([]byte(msg.resp), msg.addr)
				if err != nil {
					udpLog.Printf("failed to write UDP packet: %v\n", err)
					continue
				}

				udpLog.Printf("%d bytes written to socket %s\n", n, msg.addr.String())
				continue
			}

			go processUDPMessage(msg, cmdCh, msgCh)
		}
	}

}
