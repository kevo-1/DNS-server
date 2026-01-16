package protocol

import (
	"encoding/binary"
	"fmt"
)

type Parser struct {
	data   []byte
	offset int
}

func NewParser(data []byte) *Parser {
	return &Parser{
		data:   data,
		offset: 0,
	}
}

func (p *Parser) ParseMessage() (*Message, error) {
	msg := &Message{}

	if err := p.parseHeader(&msg.Header); err != nil {
		return nil, fmt.Errorf("parse header: %w", err)
	}

	for i := 0; i < int(msg.Header.QuestionCount); i++ {
		q, err := p.parseQuestion()
		if err != nil {
			return nil, fmt.Errorf("parse question %d: %w", i, err)
		}
		msg.Questions = append(msg.Questions, q)
	}
	
	for i := 0; i < int(msg.Header.AnswerCount); i++ {
		rr, err := p.parseResourceRecord()
		if err != nil {
			return nil, fmt.Errorf("parse answer %d: %w", i, err)
		}
		msg.Answers = append(msg.Answers, rr)
	}
	
	for i := 0; i < int(msg.Header.AuthorityCount); i++ {
		rr, err := p.parseResourceRecord()
		if err != nil {
			return nil, fmt.Errorf("parse authority %d: %w", i, err)
		}
		msg.Authorities = append(msg.Authorities, rr)
	}
	
	for i := 0; i < int(msg.Header.AdditionalCount); i++ {
		rr, err := p.parseResourceRecord()
		if err != nil {
			return nil, fmt.Errorf("parse additional %d: %w", i, err)
		}
		msg.Additional = append(msg.Additional, rr)
	}
	
	return msg, nil
}

func (p *Parser) parseHeader(h *Header) error {
	if len(p.data) < 12 {
		return fmt.Errorf("message too short: %d bytes", len(p.data))
	}
	
	h.ID = binary.BigEndian.Uint16(p.data[0:2])
	h.Flags = binary.BigEndian.Uint16(p.data[2:4])
	h.QuestionCount = binary.BigEndian.Uint16(p.data[4:6])
	h.AnswerCount = binary.BigEndian.Uint16(p.data[6:8])
	h.AuthorityCount = binary.BigEndian.Uint16(p.data[8:10])
	h.AdditionalCount = binary.BigEndian.Uint16(p.data[10:12])
	
	p.offset = 12
	return nil
}

func (p *Parser) parseName() (string, error) {
	var name string
	jumped := false
	jumpOffset := 0
	
	for {
		if p.offset >= len(p.data) {
			return "", fmt.Errorf("unexpected end of data")
		}
		
		length := p.data[p.offset]
		
		if length&0xC0 == 0xC0 {
			if p.offset+1 >= len(p.data) {
				return "", fmt.Errorf("invalid compression pointer")
			}
			
			pointer := binary.BigEndian.Uint16(p.data[p.offset:p.offset+2]) & 0x3FFF
			
			if !jumped {
				jumpOffset = p.offset + 2
			}
			
			p.offset = int(pointer)
			jumped = true
			continue
		}
		
		if length == 0 {
			p.offset++
			break
		}
		
		p.offset++
		if p.offset+int(length) > len(p.data) {
			return "", fmt.Errorf("label length exceeds data")
		}
		
		if name != "" {
			name += "."
		}
		name += string(p.data[p.offset : p.offset+int(length)])
		p.offset += int(length)
	}
	
	if jumped {
		p.offset = jumpOffset
	}
	
	return name, nil
}

func (p *Parser) parseQuestion() (Question, error) {
	var q Question
	
	name, err := p.parseName()
	if err != nil {
		return q, fmt.Errorf("parse name: %w", err)
	}
	q.Name = name
	
	if p.offset+4 > len(p.data) {
		return q, fmt.Errorf("question too short")
	}
	
	q.Type = binary.BigEndian.Uint16(p.data[p.offset : p.offset+2])
	q.Class = binary.BigEndian.Uint16(p.data[p.offset+2 : p.offset+4])
	p.offset += 4
	
	return q, nil
}

func (p *Parser) parseResourceRecord() (ResourceRecord, error) {
	var rr ResourceRecord
	
	name, err := p.parseName()
	if err != nil {
		return rr, fmt.Errorf("parse name: %w", err)
	}
	rr.Name = name
	
	if p.offset+10 > len(p.data) {
		return rr, fmt.Errorf("resource record too short")
	}
	
	rr.Type = binary.BigEndian.Uint16(p.data[p.offset : p.offset+2])
	rr.Class = binary.BigEndian.Uint16(p.data[p.offset+2 : p.offset+4])
	rr.TTL = binary.BigEndian.Uint32(p.data[p.offset+4 : p.offset+8])
	rr.RDLength = binary.BigEndian.Uint16(p.data[p.offset+8 : p.offset+10])
	p.offset += 10
	
	if p.offset+int(rr.RDLength) > len(p.data) {
		return rr, fmt.Errorf("rdata length exceeds data")
	}
	
	rr.RData = make([]byte, rr.RDLength)
	copy(rr.RData, p.data[p.offset:p.offset+int(rr.RDLength)])
	p.offset += int(rr.RDLength)
	
	return rr, nil
}