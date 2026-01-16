package protocol

import (
	"fmt"
	"net"
	"strings"
	"time"
)

func ParseMessage(data []byte) (*Message, error) {
	parser := NewParser(data)
	return parser.ParseMessage()
}

func BuildMessage(msg *Message) ([]byte, error) {
	msg.Header.QuestionCount = uint16(len(msg.Questions))
	msg.Header.AnswerCount = uint16(len(msg.Answers))
	msg.Header.AuthorityCount = uint16(len(msg.Authorities))
	msg.Header.AdditionalCount = uint16(len(msg.Additional))

	builder := NewBuilder()
	return builder.BuildMessage(msg)
}

func ParseTTL(ttl int) time.Duration {
	return time.Duration(ttl) * time.Second
}

func (rr *ResourceRecord) GetStringData() (string, error) {
	switch rr.Type {
	case TypeA:
		if len(rr.RData) != 4 {
			return "", fmt.Errorf("invalid A record data length: %d", len(rr.RData))
		}
		return net.IP(rr.RData).String(), nil

	case TypeAAAA:
		if len(rr.RData) != 16 {
			return "", fmt.Errorf("invalid AAAA record data length: %d", len(rr.RData))
		}
		return net.IP(rr.RData).String(), nil

	case TypeNS, TypeCNAME, TypePTR:
		parser := NewParser(rr.RData)
		name, err := parser.parseName()
		
		if err != nil {
			return "", fmt.Errorf("failed to parse domain name: %w", err)
		}
		return name, nil

	case TypeMX:
		if len(rr.RData) < 3 {
			return "", fmt.Errorf("invalid MX record data length: %d", len(rr.RData))
		}
		parser := NewParser(rr.RData[2:])
		name, err := parser.parseName()
		
		if err != nil {
			return "", fmt.Errorf("failed to parse MX domain: %w", err)
		}
		return name, nil

	case TypeTXT:
		if len(rr.RData) == 0 {
			return "", nil
		}
		length := int(rr.RData[0])
		
		if len(rr.RData) < length+1 {
			return "", fmt.Errorf("invalid TXT record data")
		}
		return string(rr.RData[1 : length+1]), nil

	default:
		return "", fmt.Errorf("unsupported record type: %d", rr.Type)
	}
}

func CreateARecord(name string, ip string, ttl uint32) (ResourceRecord, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return ResourceRecord{}, fmt.Errorf("invalid IP address: %s", ip)
	}

	ipv4 := parsedIP.To4()
	if ipv4 == nil {
		return ResourceRecord{}, fmt.Errorf("not an IPv4 address: %s", ip)
	}

	return ResourceRecord{
		Name:     name,
		Type:     TypeA,
		Class:    ClassIN,
		TTL:      ttl,
		RDLength: 4,
		RData:    []byte(ipv4),
	}, nil
}

func CreateResponse(query *Message, answers []ResourceRecord) *Message {
	response := &Message{
		Header: Header{
			ID:    query.Header.ID,
			Flags: FlagQR | FlagRD | FlagRA,
		},
		Questions: query.Questions,
		Answers:   answers,
	}

	if len(answers) > 0 {
		response.Header.Flags |= RCodeNoError
	} else {
		response.Header.Flags |= RCodeNXDomain
	}

	return response
}


func CreateErrorResponse(query *Message, rcode uint16) *Message {
	return &Message{
		Header: Header{
			ID:    query.Header.ID,
			Flags: FlagQR | (rcode & 0x0F),
		},
		Questions: query.Questions,
	}
}

func DomainToLabels(domain string) []string {
	domain = strings.TrimSuffix(domain, ".")
	return strings.Split(domain, ".")
}

func LabelsToDomain(labels []string) string {
	return strings.Join(labels, ".")
}

func EncodeIPv4(ip string) ([]byte, error) {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return nil, fmt.Errorf("invalid IP: %s", ip)
	}

	ipv4 := parsed.To4()
	if ipv4 == nil {
		return nil, fmt.Errorf("not IPv4: %s", ip)
	}

	return []byte(ipv4), nil
}

func DecodeIPv4(data []byte) (string, error) {
	if len(data) != 4 {
		return "", fmt.Errorf("invalid IPv4 data length: %d", len(data))
	}
	return net.IP(data).String(), nil
}

func EncodeDomainName(domain string) []byte {
	builder := NewBuilder()
	builder.buildName(domain)
	return builder.data
}
