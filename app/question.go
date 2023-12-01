package main

import "encoding/binary"

type DNSQuestion struct {
	QName  []string
	QType  QType
	QClass QClass
}
type QType uint16

const (
	A     QType = iota + 1 // a host address
	NS                     // an authoritative name server
	MD                     // a mail destination (Obsolete - use MX)
	MF                     // a mail forwarder (Obsolete - use MX)
	CNAME                  // the canonical name for an alias
	SOA                    // marks the start of a zone of authority
	MB                     // a mailbox domain name (EXPERIMENTAL)
	MG                     // a mail group member (EXPERIMENTAL)
	MR                     // a mail rename domain name (EXPERIMENTAL)
	NULL                   // a null RR (EXPERIMENTAL)
	WKS                    // a well known service description
	PTR                    // a domain name pointer
	HINFO                  // host information
	MINFO                  // mailbox or mail list information
	MX                     // mail exchange
	TXT                    // text string
	AXFR  QType = 252      // a request for a transfer of an entire zone
	MAILB QType = 253      // a request for mailbox-related records (MB, MG or MR)
	MAILA QType = 254      // a request for mail agent RRs (Obsolete - see MX)
)
const (
	ANY uint16 = 255
)

type QClass uint16

const (
	IN QClass = iota + 1 // the Internet
	CS                   // CSNET class
	CH                   // CHAOS class
	HS                   // Hesiod
)

// ToBytes returns a byte slice from a DNSQuestion
// Used to turn an outgoing response into a sendable type
func (q *DNSQuestion) ToBytes() []byte {
	var buf []byte
	// Appends a byte indicating length, then the name
	for _, name := range q.QName {
		length := len(name)
		buf = append(buf, byte(length))
		buf = append(buf, []byte(name)...)
	}

	buf = append(buf, 0) // appends a 0 byte to indicate the end of the names
	buf = binary.BigEndian.AppendUint16(buf, uint16(q.QType))
	buf = binary.BigEndian.AppendUint16(buf, uint16(q.QClass))

	return buf
}

// DNSQuestionFromBytes returns a DNSQuestion from a byte slice and a byte slice of anything remaining >4 bytes so that the message can continue to be processed
// Used to turn an incoming question into a usable type
func DNSQuestionFromBytes(bytes []byte) (*DNSQuestion, []byte) {
	var qname []string
	for length := bytes[0]; length != 0; length = bytes[0] {
		qname = append(qname, string(bytes[1:length+1]))
		bytes = bytes[length+1:]
	}
	bytes = bytes[1:]

	question := &DNSQuestion{
		QName:  qname,
		QType:  QType(binary.BigEndian.Uint16(bytes[0:2])),
		QClass: QClass(binary.BigEndian.Uint16(bytes[2:4])),
	}
	if len(bytes) > 4 {
		return question, bytes[4:]
	}
	return question, nil
}
