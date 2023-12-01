package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("DNS Server Started")
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}
	// Start server
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
		// Decode received DNS query into its constituent parts
		receivedDNS := DNSFromBytes(buf[:size])
		flags := receivedDNS.Header.GetFlags()
		fmt.Printf("Received %d bytes from %s\n", size, source)
		fmt.Printf("Header: %+v\n", receivedDNS.Header)
		fmt.Printf("Flags: %+v\n", flags)

		fmt.Printf("Questions: %+v\n", receivedDNS.Questions)
		// Create blank response
		var response DNS
		response.Header = *NewDNSHeader(receivedDNS.Header.ID, uint16(len(receivedDNS.Questions)), 0, 0, 0)
		rd := uint16(0)
		if flags.RD {
			rd = 1
		}
		rcode := uint16(0)
		if flags.OpCode != 0 {
			rcode = 4
		}

		response.Header.SetFlags(1, uint16(flags.OpCode), 0, 0, rd, 0, 0, rcode)
		// Copy questions to response
		response.Questions = receivedDNS.Questions
		for _, question := range receivedDNS.Questions {
			if question.QType == A {
				var answer DNSAnswer
				answer.Name = question.QName
				answer.Type = A
				answer.Class = IN
				answer.TTL = 60
				answer.Length = 4
				if strings.Join(question.QName, ".") == "localhost" {
					answer.Data = []byte{127, 0, 0, 1}
				} else {
					answer.Data = []byte{8, 8, 8, 8}
				}
				response.Answers = append(response.Answers, answer)
				response.Header.ANCount++
			}
		}
		fmt.Printf("Response: %+v\n", response)
		resData := response.ToBytes()
		_, err = udpConn.WriteToUDP(resData, source)
		if err != nil {
			fmt.Println("Failed to send response: ", err)
		}
	}
}
