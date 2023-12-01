package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port string
}

func (s *Server) init() {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP adddress", err)
		return
	}
	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address", err)
		return
	}
	defer udpConn.Close()
	buf := make([]byte, 512)
	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data", err)
			break
		}
		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
		message := newMessage(buf[:size])
		fmt.Println("Query name:", message.Question.Name)
		message.Header.QueryResponse = true
		message.Header.QDCount = 1
		message.Responses = append(message.Responses, &Response{
			Name:     message.Question.Name,
			Type:     1,
			Class:    1,
			Ttl:      60,
			RDLength: 4,
			RData:    []byte{0x08, 0x08, 0x08, 0x08},
		})
		message.Header.ANCount = 1
		response := message.Marshal()
		var buffer bytes.Buffer
		err = binary.Write(&buffer, binary.BigEndian, response)
		if err != nil {
			fmt.Println("Failed to serialise response", err)
		}
		_, err = udpConn.WriteToUDP(buffer.Bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response", err)
			break
		}
	}
}
