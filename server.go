package main

import (
	"bytes"
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

	buffer := make([]byte, 512)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error while reading incoming packet ", err)
		}
		if n > 512 {
			fmt.Println("Exceed DNS Payload size")
		}

		// Lire l'identifiant
		headerSection := [12]byte(buffer[:12])
		requestHeader := dns.NewHeader(headerSection)
		requestHeader.Read()

		questionData := buffer[12:33]
		question, questionError := dns.NewQuestion(questionData)
		if questionError != nil {
			fmt.Println("Couldn't parse the question section of the DNS request:", err)
		}
		question.Read()

		// Answer
		var response bytes.Buffer
		responseHeader := dns.Header{
			ID:      requestHeader.ID,
			QR:      0x01,
			OPCODE:  requestHeader.OPCODE,
			AA:      0x00,
			TC:      0x00,
			RD:      requestHeader.RD,
			RA:      0x00,
			Z:       0x00,
			RCODE:   0x00,
			QDCOUNT: 1,
			ANCOUNT: 0,
			NSCOUNT: 0,
			ARCOUNT: 0,
		}
		respH := responseHeader.Write()
		_, err = response.Write(respH[:])
		if err != nil {
			fmt.Println("Error while writing Header into the response buffer")
		}

		responseQuestion := dns.Question{
			Name:  "\x0ccodecrafters\x02io\x00",
			Type:  1,
			Class: 1,
		}
		respQ := responseQuestion.Write()
		_, err = response.Write(respQ)
		if err != nil {
			fmt.Println("Error while writing Question into the response buffer")
		}

		responseAnswer := dns.Answer{
			Name:     "codecrafters.io",
			Type:     1,
			Class:    1,
			TTL:      60,
			RDLENGTH: 4,
			RDATA:    []byte("8.8.8.8"),
		}
		responseAnswer.Read()
		respA := responseAnswer.Write()
		_, err = response.Write(respA[:])
		if err != nil {
			fmt.Println("Error while writing Answer into the response buffer")
		}

		_, err = conn.WriteToUDP(response.Bytes(), addr)
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
