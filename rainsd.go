package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
)

const (
	RAINSD_PORT = 1228
	OUTBOX_LEN  = 20
	INBOX_LEN   = 20
)

type Switchboard struct {
	outboxes  map[string]chan string
	inbox     chan string
	localName string
}

func writeOutbox(w io.Writer, peerName string, outbox <-chan string) {
	for msg := range outbox {
		fmt.Fprintln(w, string)
	}

	w.Close()
}

func readInbox(r io.Reader, peerName string, inbox chan<- string) {
	scanner = bufio.NewScanner(r)
	for scanner.Scan() {
		inbox <- scanner.Text()
	}
	// FIXME need to tell EOF from other errors
}

func (sb *Switchboard) setupConnection(peerName string, conn net.Conn) {
	outbox = make(chan []byte, OUTQ_LEN)
	sb.outboxes[peerName] = outbox
	go writeOutbox(conn, peerName, outbox)
	go readInbox(conn, peerName, sb.inbox)
}

func (sb *Switchboard) sendMessage(peerName string, message string) error {
	var outbox chan []byte
	var ok bool

	outbox, ok = sb.outboxes[peerName]

	if !ok {
		// open a connection and start handlers on it
		if conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", address, RAINSD_PORT)); err != nil {
			return err
		}

		sb.setupConnection(peerName, conn)
		outbox <- "HELLO " + localName
		log.Printf("switchboard opened connection to %v at %v\n", peerName, conn.RemoteAddr)

	}

	outbox <- message
	return nil
}

func (sb *Switchboard) listen() error {
	if listener, err := net.Listen("tcp", fmt.Sprintf("*:%d", RAINSD_PORT)); err != nil {
		return error
	}

	for {
		if conn, err := listener.Accept(); err != nil {
			log.Print(err)
			continue
		}

		go func() {
			scanner = bufio.NewScanner(conn)
			if scanner.Scan() {
				hello := scanner.Text()
				peerName := strings.Split(hello, " ")[1]
				sb.setup(peerName, conn)
				log.Printf("switchboard got connection from %v at %v\n", peerName, conn.RemoteAddr)
			}
		}()
	}

}

func (sb *Switchboard) reapConnection(peerName string) {
	outbox, ok := sb.outboxes[peerName]
	if ok {
		delete(sb.outboxes, peerName)
		outbox.Close()
	}
}

func main() {

	sb := new(Switchboard)
	go sb.listen()

}
