package main

type DNS struct {
	Header    DNSHeader
	Questions []DNSQuestion
	Answers   []DNSAnswer
}

func DNSFromBytes(data []byte) DNS {
	header := DNSHeaderFromBytes(data[:12])
	data = data[12:]
	dns := DNS{*header, nil, nil}
	var i uint16
	var question *DNSQuestion
	for i = 0; i < header.QDCount; i++ {
		question, data = DNSQuestionFromBytes(data)
		dns.Questions = append(dns.Questions, *question)
	}
	var answer *DNSAnswer
	for i = 0; i < header.ANCount; i++ {
		answer, data = DNSAnswerFromBytes(data)
		dns.Answers = append(dns.Answers, *answer)

	}
	return dns
}

func (dns DNS) ToBytes() []byte {
	var bytes []byte
	bytes = append(bytes, dns.Header.ToBytes()...)
	for _, question := range dns.Questions {
		bytes = append(bytes, question.ToBytes()...)
	}
	for _, answer := range dns.Answers {
		bytes = append(bytes, answer.ToBytes()...)
	}
	return bytes
}
