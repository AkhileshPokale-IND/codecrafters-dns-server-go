package main

import "encoding/binary"

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func ParseDNSHeader(bytes []byte) *DNSHeader {
	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(bytes[0:2]),
		Flags:   binary.BigEndian.Uint16(bytes[2:4]),
		QDCOUNT: binary.BigEndian.Uint16(bytes[4:6]),
		ANCOUNT: binary.BigEndian.Uint16(bytes[6:8]),
		NSCOUNT: binary.BigEndian.Uint16(bytes[8:10]),
		ARCOUNT: binary.BigEndian.Uint16(bytes[10:12]),
	}
}
func (hdr *DNSHeader) Dump() []byte {
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], hdr.ID)
	binary.BigEndian.PutUint16(buf[2:4], hdr.Flags)
	binary.BigEndian.PutUint16(buf[4:6], hdr.QDCOUNT)
	binary.BigEndian.PutUint16(buf[6:8], hdr.ANCOUNT)
	binary.BigEndian.PutUint16(buf[8:10], hdr.NSCOUNT)
	binary.BigEndian.PutUint16(buf[10:12], hdr.ARCOUNT)
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
func NewDNSHeader(id, flags, qdcount, ancount, nscount, arcount uint16) *DNSHeader {
	return &DNSHeader{
		ID:      id,
		Flags:   flags,
		QDCOUNT: qdcount,
		ANCOUNT: ancount,
		NSCOUNT: nscount,
		ARCOUNT: arcount,
	}
}
