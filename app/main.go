package main

import (
	"fmt"
	"net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()
	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %v\n", size, source, receivedData)

		hdr := ParseDNSHeader(buf[:12])
		fmt.Printf("Received Header: %+v\n", hdr)
		// Create an empty response

		response := NewDNSHeader(1234, 0, 1, 0, 0, 0)
		response.SetFlags(1, 0, 0, 0, 0, 0, 0, 0)

		fmt.Printf("Response Header: %+v\n", response)

		question := NewDNSQuestion(LabelSequence("codecrafters.io"), 1, 1)
		r := response.Dump()
		r = append(r, question.Dump()...)
		_, err = udpConn.WriteToUDP(r, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
