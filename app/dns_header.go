package main

import "encoding/binary"

// The DNS Header is always 12 bytes long
// ID		Packet Identifier	16 bits
// = A random identifier is assigned to query packets. Response packets must reply with the same id. This is needed to differentiate responses due to the stateless nature of UDP.
// FLAGS	Flags				16 bits
// QDCOUNT	Question Count		16 bits
// = The number of entries in the Question Section
// ANCOUNT	Answer Count		16 bits
// = The number of entries in the Answer Section
// NSCOUNT	Authority Count		16 bits
// = The number of entries in the Authority Section
// ARCOUNT	Additional Count	16 bits
// = The number of entries in the Additional Section
type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCount uint16
	ANCount uint16
	NSCount uint16
	ARCount uint16
}

// construct a DNSHeader struct
func ParseDNSHeader(bytes []byte) *DNSHeader {
	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(bytes[0:2]),
		Flags:   binary.BigEndian.Uint16(bytes[2:4]),
		QDCount: binary.BigEndian.Uint16(bytes[4:6]),
		ANCount: binary.BigEndian.Uint16(bytes[6:8]),
		NSCount: binary.BigEndian.Uint16(bytes[8:10]),
		ARCount: binary.BigEndian.Uint16(bytes[10:12]),
	}
}

// dump DNSHeader to Bytes
func (hdr *DNSHeader) Dump() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], hdr.ID)
	binary.BigEndian.PutUint16(buf[2:4], hdr.Flags)
	binary.BigEndian.PutUint16(buf[4:6], hdr.QDCount)
	binary.BigEndian.PutUint16(buf[6:8], hdr.ANCount)
	binary.BigEndian.PutUint16(buf[8:10], hdr.NSCount)
	binary.BigEndian.PutUint16(buf[10:12], hdr.ARCount)
	return buf
}

// Query/Response Indicator (QR) : 1
// Operation Code (OPCODE) : 4
// Authoritative Answer (AA) : 1
// Truncation (TC) : 1
// Recursion Desired (RD) : 1
// Recursion Available (RA) : 1
// Reserved (Z)	: 3
// Response Code (RCODE) : 4
func (hdr *DNSHeader) SetFlags(qr, opcode, aa, tc, rd, ra, z, rcode uint16) {
	hdr.Flags = qr<<15 | opcode<<11 | aa<<10 | tc<<9 | rd<<8 | ra<<7 | z<<4 | rcode
}

// DNSHeader constructor
func NewDNSHeader(id, flags, qdcount, ancount, nscount, arcount uint16) *DNSHeader {
	return &DNSHeader{
		ID:      id,
		Flags:   flags,
		QDCount: qdcount,
		ANCount: ancount,
		NSCount: nscount,
		ARCount: arcount,
	}

}
