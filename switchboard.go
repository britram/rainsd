package rains

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	OUTBOX_LEN = 20
	INBOX_LEN  = 20
)

// Switchboard manages connections to peers

type Switchboard struct {
	outboxes  map[string]chan string
	inbox     chan string
	LocalName string
}

func writeOutbox(w io.WriteCloser, peerName string, outbox <-chan string) {
	for msg := range outbox {
		fmt.Fprintln(w, msg)
	}

	w.Close()
}

func readInbox(r io.Reader, peerName string, inbox chan<- string) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		inbox <- scanner.Text()
	}
	// FIXME need to tell EOF from other errors
}

func (sb *Switchboard) setupConnection(peerName string, conn net.Conn) {
	outbox := make(chan string, OUTBOX_LEN)
	sb.outboxes[peerName] = outbox
	go writeOutbox(conn, peerName, outbox)
	go readInbox(conn, peerName, sb.inbox)
}

// Send a message to a peer identified by name.
// Open a connection to that peer if none exists.
func (sb *Switchboard) SendMessage(peerName string, address string, message string) error {
	var outbox chan string
	var ok bool

	outbox, ok = sb.outboxes[peerName]

	if !ok {
		// open a connection and start handlers on it
		conn, err := net.Dial("tcp", address)
		if err != nil {
			return err
		}

		sb.setupConnection(peerName, conn)
		outbox <- "HELLO " + sb.LocalName
		log.Printf("switchboard opened connection to %v at %v\n", peerName, conn.RemoteAddr)

	}

	outbox <- message
	return nil
}

func (sb *Switchboard) NextMessage() string {
	return <-sb.inbox
}

func (sb *Switchboard) Listen(port int) error {

	listener, err := net.Listen("tcp", fmt.Sprintf("*:%d", port))
	if err != nil {
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go func() {
			scanner := bufio.NewScanner(conn)
			if scanner.Scan() {
				hello := scanner.Text()
				peerName := strings.Split(hello, " ")[1]
				sb.setupConnection(peerName, conn)
				log.Printf("switchboard got connection from %v at %v\n", peerName, conn.RemoteAddr)
			}
		}()
	}

}

func NewSwitchboard() *Switchboard {
	var sb Switchboard
	sb.inbox = make(chan string, INBOX_LEN)
	return &sb
}

func (sb *Switchboard) reapConnection(peerName string) {
	outbox, ok := sb.outboxes[peerName]
	if ok {
		delete(sb.outboxes, peerName)
		close(outbox)
	}
}
