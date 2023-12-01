package main

import (
	"encoding/binary"
)

type DNSHeader struct {
	ID      uint16 // Packet identifier. Random ID. Response must reply with the same ID.
	Flags   uint16 // See SetFlags
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// DNSHeaderFromBytes returns a DNSHeader with appropriate fields set by their byte numbers
// Used to turn an incoming request into a usable type
func DNSHeaderFromBytes(bytes []byte) *DNSHeader {
	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(bytes[0:2]),
		Flags:   binary.BigEndian.Uint16(bytes[2:4]),
		QDCount: binary.BigEndian.Uint16(bytes[4:6]),
		ANCount: binary.BigEndian.Uint16(bytes[6:8]),
		NSCount: binary.BigEndian.Uint16(bytes[8:10]),
		ARCount: binary.BigEndian.Uint16(bytes[10:12]),
	}
}

// ToBytes returns a byte slice formed of provided DNSHeader fields
// Used to turn an outgoing response into a sendable type
func (h *DNSHeader) ToBytes() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], h.ID)
	binary.BigEndian.PutUint16(buf[2:4], h.Flags)
	binary.BigEndian.PutUint16(buf[4:6], h.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], h.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], h.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], h.ARCount)
	return buf
}
func (h *DNSHeader) SetFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint16) {
	// Made of Query Response, Operation Code, Authoritative Answer, Truncation, Recursion Desired, Recursion Available, Reserved, Response Code
	h.Flags = qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode
}

// NewDNSHeader creates a new DNSHeader with fields set as provided, except Flags which are set to 0
func NewDNSHeader(id, qdcount, ancount, nscount, arcount uint16) *DNSHeader {
	return &DNSHeader{
		ID:      id,
		Flags:   0,
		QDCount: qdcount,
		ANCount: ancount,
		NSCount: nscount,
		ARCount: arcount,
	}
}
