package main

import (
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn *net.UDPConn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		// Lis les requêtes entrantes
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error while reading incoming packet ", err)
		}

		fmt.Print("-> ", string(buffer[0:n-1]))

		if strings.TrimSpace(string(buffer[0:n])) == "STOP" {
			fmt.Println("Exiting UDP server!")
			return
		}

		data := []byte("Hello world")
		fmt.Printf("data: %s\n", string(data))
		_, err = conn.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func main() {
	PORT := ":2053"
	server, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		panic(err)
	}
	fmt.Println("Server running on port", server.Port)

	// Ouvre une socket pour recevoir les connexions
	conn, err := net.ListenUDP("udp4", server)
	if err != nil {
		fmt.Println(err)
		return
	}
	handleConnection(conn)
}
