package protocol

import (
	"encoding/binary"
	"strings"
)

type Builder struct {
	data []byte
}

func NewBuilder() *Builder {
	return &Builder{
		data: make([]byte, 0, 512),
	}
}

func (b *Builder) BuildMessage(msg *Message) ([]byte, error) {
	b.data = make([]byte, 0, 512)
	
	b.buildHeader(msg.Header)
	
	for _, q := range msg.Questions {
		b.buildQuestion(q)
	}
	
	for _, rr := range msg.Answers {
		b.buildResourceRecord(rr)
	}
	
	for _, rr := range msg.Authorities {
		b.buildResourceRecord(rr)
	}
	
	for _, rr := range msg.Additional {
		b.buildResourceRecord(rr)
	}
	
	return b.data, nil
}

func (b *Builder) buildHeader(h Header) {
	header := make([]byte, 12)
	
	binary.BigEndian.PutUint16(header[0:2], h.ID)
	binary.BigEndian.PutUint16(header[2:4], h.Flags)
	binary.BigEndian.PutUint16(header[4:6], h.QuestionCount)
	binary.BigEndian.PutUint16(header[6:8], h.AnswerCount)
	binary.BigEndian.PutUint16(header[8:10], h.AuthorityCount)
	binary.BigEndian.PutUint16(header[10:12], h.AdditionalCount)
	
	b.data = append(b.data, header...)
}

func (b *Builder) buildName(name string) {
	name = strings.TrimSuffix(name, ".")
	labels := strings.Split(name, ".")
	
	for _, label := range labels {
		b.data = append(b.data, byte(len(label)))
		b.data = append(b.data, []byte(label)...)
	}
	
	b.data = append(b.data, 0)
}

func (b *Builder) buildQuestion(q Question) {
	b.buildName(q.Name)
	
	qtype := make([]byte, 2)
	qclass := make([]byte, 2)
	
	binary.BigEndian.PutUint16(qtype, q.Type)
	binary.BigEndian.PutUint16(qclass, q.Class)
	
	b.data = append(b.data, qtype...)
	b.data = append(b.data, qclass...)
}

func (b *Builder) buildResourceRecord(rr ResourceRecord) {
	b.buildName(rr.Name)
	
	rrType := make([]byte, 2)
	rrClass := make([]byte, 2)
	rrTTL := make([]byte, 4)
	rrLength := make([]byte, 2)
	
	binary.BigEndian.PutUint16(rrType, rr.Type)
	binary.BigEndian.PutUint16(rrClass, rr.Class)
	binary.BigEndian.PutUint32(rrTTL, rr.TTL)
	binary.BigEndian.PutUint16(rrLength, rr.RDLength)
	
	b.data = append(b.data, rrType...)
	b.data = append(b.data, rrClass...)
	b.data = append(b.data, rrTTL...)
	b.data = append(b.data, rrLength...)
	
	b.data = append(b.data, rr.RData...)
}