package main

import (
	"encoding/binary"
	"fmt"
	"strings"
)

func concatenate(bufs ...[]byte) []byte {
	var buf []byte
	for _, b := range bufs {
		buf = append(buf, b...)
	}
	return buf
}

type Message struct {
	*Header
	*Question
	Responses []*Response
}

func newMessage(buf []byte) *Message {
	return &Message{
		Header:   newHeader(buf[:12]),
		Question: newQuestion(buf[12:]),
	}
}
func (msg *Message) Marshal() []byte {
	buf := concatenate(msg.Header.Marshal(), msg.Question.Marshal())
	for _, response := range msg.Responses {
		buf = concatenate(buf, response.Marshal())
	}
	return buf
}

type Header struct {
	ID            uint16
	QueryResponse bool
	Opcode        uint16

	AuthoritativeAnswer bool
	Truncation          bool
	RecursionDesired    bool
	RecursionAvailable  bool
	Z                   uint16
	RCode               uint16
	QDCount             uint16
	ANCount             uint16
	NSCount             uint16
	ARCount             uint16
}

func boolToUint16(val bool) uint16 {
	if val {
		return 1
	}

	return 0
}

func (header *Header) GetFlags() uint16 {
	var flags uint16
	flags |= boolToUint16(header.QueryResponse) << 15
	flags |= header.Opcode << 11

	flags |= boolToUint16(header.AuthoritativeAnswer) << 10
	flags |= boolToUint16(header.Truncation) << 9
	flags |= boolToUint16(header.RecursionDesired) << 8
	flags |= boolToUint16(header.RecursionAvailable) << 7
	flags |= header.Z << 4
	flags |= header.RCode << 0

	return flags
}
func newHeader(buf []byte) *Header {
	header := &Header{}
	fmt.Printf("0: %08b\n", buf[0])
	fmt.Printf("1: %08b\n", buf[1])
	fmt.Printf("2: %08b\n", buf[2])
	fmt.Printf("3: %08b\n", buf[3])
	header.ID = binary.BigEndian.Uint16(buf[0:2])
	header.QueryResponse = (buf[2] & (1 << 7)) != 0
	header.Opcode = uint16((buf[2] >> 3) & 15)
	header.AuthoritativeAnswer = (buf[2] & (1 << 2)) != 0
	header.Truncation = (buf[2] & (1 << 1)) != 0
	header.RecursionDesired = buf[2]&1 != 0
	header.RecursionAvailable = (buf[3] & (1 << 7)) != 0
	header.Z = uint16((buf[4] >> 3) & 7)
	header.RCode = uint16(buf[4] & 15)
	header.QDCount = binary.BigEndian.Uint16(buf[4:6])
	header.ANCount = binary.BigEndian.Uint16(buf[6:8])
	header.NSCount = binary.BigEndian.Uint16(buf[8:10])

	header.ARCount = binary.BigEndian.Uint16(buf[10:12])
	return header
}
func (header *Header) Marshal() []byte {
	result := make([]byte, 12)
	binary.BigEndian.PutUint16(result[0:2], header.ID)
	binary.BigEndian.PutUint16(result[2:4], header.GetFlags())
	binary.BigEndian.PutUint16(result[4:6], header.QDCount)
	binary.BigEndian.PutUint16(result[6:8], header.ANCount)
	binary.BigEndian.PutUint16(result[8:10], header.NSCount)
	binary.BigEndian.PutUint16(result[10:12], header.ANCount)
	return result
}

type Question struct {
	Name   string
	QType  uint16
	QClass uint16
}

func newQuestion(buf []byte) *Question {
	var name strings.Builder
	i := 0
	for {
		size := int(buf[i])
		label := buf[i+1 : i+1+size]
		name.Write(label)
		name.WriteRune('.')
		i = i + 1 + size
		if buf[i] == 0 {
			break
		}
	}
	return &Question{
		Name:   name.String(),
		QType:  binary.BigEndian.Uint16(buf[i+1 : i+3]),
		QClass: binary.BigEndian.Uint16(buf[i+3 : i+5]),
	}
}
func (question *Question) Marshal() []byte {
	var name []byte
	fields := strings.FieldsFunc(question.Name, func(r rune) bool {
		return r == '.'
	})
	for _, field := range fields {
		name = concatenate(
			name,
			[]byte{byte(len(field))},
			[]byte(field),
		)
	}
	name = concatenate(
		name,
		[]byte{0},
	)
	return concatenate(
		name,
		uint16ToBytes(question.QType),
		uint16ToBytes(question.QClass),
	)
}
func uint16ToBytes(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}
func uint32ToBytes(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf
}

type Response struct {
	Name     string
	Type     uint16
	Class    uint16
	Ttl      uint32
	RDLength uint16
	RData    []byte
}

func (response *Response) Marshal() []byte {
	var name []byte
	fields := strings.FieldsFunc(response.Name, func(r rune) bool {
		return r == '.'
	})
	for _, field := range fields {
		name = concatenate(
			name,
			[]byte{byte(len(field))},
			[]byte(field),
		)
	}
	name = concatenate(name, []byte{0x00})
	var buf []byte
	buf = concatenate(
		buf,
		name,
		uint16ToBytes(response.Type),
		uint16ToBytes(response.Class),
		uint32ToBytes(response.Ttl),
		uint16ToBytes(response.RDLength),
		response.RData,
	)
	return buf
}
