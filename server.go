package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ECecillo/dns-server-go/dns"
)

type Server struct {
	wg         sync.WaitGroup
	addr       net.UDPAddr
	shutdown   chan struct{}
	connection chan net.UDPConn
}

func newServer(ip string, port int) (*Server, error) {
	addr := net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}

	return &Server{
		addr:       addr,
		shutdown:   make(chan struct{}),
		connection: make(chan net.UDPConn),
	}, nil
}

func (s *Server) start() {
	s.wg.Add(1)
	go s.listenForConnections()
}

func (s *Server) listenForConnections() {
	defer s.wg.Done()

	for {
		select {
		case <-s.shutdown:
			return
		case conn := <-s.connection:
			// In case we get a new connection run a go routine to handle it
			go s.handleClient(&conn)
		}
	}
}

func (s *Server) handleClient(conn *net.UDPConn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		_, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error while reading incoming packet ", err)
		}

		// Lire l'identifiant
		headerSection := [12]byte(buffer[:12])
		requestHeader := dns.Create(headerSection)
		requestHeader.Read()

		question := buffer[12:]
		fmt.Println("Question:", question)

		// Pour afficher le QNAME correctement
		qname := ""
		for i := 12; buffer[i] != 0; {
			length := int(buffer[i])
			i++
			qname += string(buffer[i:i+length]) + "."
			i += length
		}
		fmt.Println("QNAME:", qname)

		responseHeader := dns.Header{
			ID:      requestHeader.ID,
			QR:      0x01,
			OPCODE:  0x00,
			AA:      0x00,
			TC:      0x00,
			RD:      0x00,
			RA:      0x00,
			Z:       0x00,
			RCODE:   0x00,
			QDCOUNT: 0,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		}
		response := responseHeader.ToByte()
		_, err = conn.WriteToUDP(response[:], addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (s *Server) stop() {
	close(s.shutdown)

	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return
	case <-time.After(time.Second):
		fmt.Println("Timed out waiting for connections to finish.")
		return
	}
}

func main() {
	argIp := flag.String("ip", "127.0.0.1", "IP to listen on")
	argPort := flag.Int("port", 2053, "Port to listen on")

	server, err := newServer(*argIp, *argPort)
	if err != nil {
		panic(err)
	}
	conn, err := net.ListenUDP("udp", &server.addr)
	if err != nil {
		fmt.Println("Error accepting connection:", err)
	}

	server.start()
	fmt.Println("Server started on port", *argPort)
	server.connection <- *conn

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Shutting down server...")
	server.stop()
	fmt.Println("Server stopped.")
}
