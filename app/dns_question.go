package main

type DNSQuestion struct {
	Name  []byte
	Type  uint16
	Class uint16
}

// Dump DNSQuestion to Bytes
func (q *DNSQuestion) Dump() []byte {
	buf := make([]byte, 0)
	buf = append(buf, []byte(q.Name)...)
	buf = append(buf, Uint16ToBytesBE(q.Type)...)
	buf = append(buf, Uint16ToBytesBE(q.Class)...)
	return buf
}
func NewDNSQuestion(name []byte, qtype, class uint16) *DNSQuestion {
	return &DNSQuestion{
		Name:  name,
		Type:  qtype,
		Class: class,
	}

}
