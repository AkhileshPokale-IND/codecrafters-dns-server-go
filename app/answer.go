package main

import "encoding/binary"

type DNSAnswer struct {
	Name   []string
	Type   QType
	Class  QClass
	TTL    uint32
	Length uint16
	Data   []byte
}

// DNSAnswerFromBytes returns a DNSAnswer from a byte slice and a byte slice of anything left over
// Used to turn an incoming request into a usable type
func DNSAnswerFromBytes(b []byte) (*DNSAnswer, []byte) {
	var answer DNSAnswer
	// Reads the length from the first byte, then iterates until reaching the '0' byte at the end
	// Each iteration adds a new name to the answer
	for length := b[0]; length != 0; length = b[0] {
		answer.Name = append(answer.Name, string(b[1:length+1]))
	}
	b = b[1:]
	answer.Type = QType(binary.BigEndian.Uint16((b[0:2])))
	answer.Class = QClass(binary.BigEndian.Uint16((b[2:4])))
	answer.TTL = binary.BigEndian.Uint32(b[4:8])
	answer.Length = binary.BigEndian.Uint16(b[8:10])
	answer.Data = b[10 : 10+answer.Length]
	if len(b) > 10+int(answer.Length) {
		return &answer, b[10+answer.Length:]
	}

	return &answer, nil
}

// ToBytes returns a byte slice from a DNSAnswer
// Used to turn an outgoing message into a sendable type
func (a *DNSAnswer) ToBytes() []byte {
	var buf []byte
	// For each name, add a length byte, then all subsequent bytes, followed by the '0' byte to indicate the end
	for _, name := range a.Name {
		buf = append(buf, byte(len(name)))
		buf = append(buf, []byte(name)...)
	}
	buf = append(buf, 0)
	buf = binary.BigEndian.AppendUint16(buf, uint16(a.Type))
	buf = binary.BigEndian.AppendUint16(buf, uint16(a.Class))
	buf = binary.BigEndian.AppendUint32(buf, a.TTL)
	buf = binary.BigEndian.AppendUint16(buf, a.Length)
	buf = append(buf, a.Data...)

	return buf
}
