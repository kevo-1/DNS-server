package resolver

import (
	"DNS-server/data"
	"DNS-server/internal/protocol"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

const (
	maxIterations = 15
	queryTimeout  = 5 * time.Second
)

var (
	ErrMaxIterationsExceeded = errors.New("maximum iterations exceeded")
	ErrNoAnswer              = errors.New("no answer received")
	ErrInvalidResponse       = errors.New("invalid DNS response")
)

type IterativeResolver struct {
	rootServers []string
	cache       *DNSCache
}

func NewIterativeResolver(cache *DNSCache) *IterativeResolver {
	return &IterativeResolver{
		rootServers: data.GetRootServers(),
		cache:       cache,
	}
}

func (r *IterativeResolver) Resolve(domain string, recordType uint16) (string, error) {
	domain = strings.TrimSuffix(domain, ".")
	domain = strings.ToLower(domain)

	if r.cache != nil {
		if ip, found := r.cache.Get(domain); found {
			return ip, nil
		}
	}

	nameservers := r.rootServers
	iteration := 0

	for iteration < maxIterations {
		iteration++

		response, err := r.queryNameserver(nameservers[0], domain, recordType)
		if err != nil {
			if len(nameservers) > 1 {
				nameservers = nameservers[1:]
				continue
			}
			return "", fmt.Errorf("failed to query nameserver: %w", err)
		}

		if len(response.Answers) > 0 {
			for _, answer := range response.Answers {
				if answer.Type == protocol.TypeA {
					ip, err := answer.GetStringData()
					if err != nil {
						continue
					}
					if r.cache != nil {
						r.cache.Set(domain, ip, time.Duration(answer.TTL)*time.Second)
					}
					return ip, nil
				} else if answer.Type == protocol.TypeCNAME {
					cname, err := answer.GetStringData()
					if err != nil {
						continue
					}
					domain = cname
					nameservers = r.rootServers
					break
				}
			}
		}

		if len(response.Authorities) > 0 {
			newNameservers := make([]string, 0)
			for _, auth := range response.Authorities {
				if auth.Type == protocol.TypeNS {
					nsName, err := auth.GetStringData()
					if err != nil {
						continue
					}
					nsIP, err := r.resolveNameserver(nsName)
					if err == nil {
						newNameservers = append(newNameservers, nsIP)
					}
				}
			}

			if len(newNameservers) > 0 {
				nameservers = newNameservers
				continue
			}
		}

		if len(response.Additional) > 0 {
			newNameservers := make([]string, 0)
			for _, add := range response.Additional {
				if add.Type == protocol.TypeA {
					ip, err := add.GetStringData()
					if err != nil {
						continue
					}
					newNameservers = append(newNameservers, ip)
				}
			}

			if len(newNameservers) > 0 {
				nameservers = newNameservers
				continue
			}
		}

		return "", ErrNoAnswer
	}

	return "", ErrMaxIterationsExceeded
}

func (r *IterativeResolver) queryNameserver(nameserver, domain string, recordType uint16) (*protocol.Message, error) {
	query := &protocol.Message{
		Header: protocol.Header{
			ID:            uint16(time.Now().Unix() & 0xFFFF),
			Flags:         0x0100,
			QuestionCount: 1,
		},
		Questions: []protocol.Question{
			{
				Name:  domain,
				Type:  recordType,
				Class: protocol.ClassIN,
			},
		},
	}

	queryData, err := protocol.BuildMessage(query)
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	conn, err := net.DialTimeout("udp", nameserver+":53", queryTimeout)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nameserver: %w", err)
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(queryTimeout))

	_, err = conn.Write(queryData)
	if err != nil {
		return nil, fmt.Errorf("failed to send query: %w", err)
	}

	buffer := make([]byte, 512)
	n, err := conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	response, err := protocol.ParseMessage(buffer[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response, nil
}

func (r *IterativeResolver) resolveNameserver(nsName string) (string, error) {
	if r.cache != nil {
		if ip, found := r.cache.Get(nsName); found {
			return ip, nil
		}
	}

	ips, err := net.LookupIP(nsName)
	if err != nil {
		return "", err
	}

	if len(ips) == 0 {
		return "", errors.New("no IP addresses found")
	}

	for _, ip := range ips {
		if ip.To4() != nil {
			ipStr := ip.String()
			if r.cache != nil {
				r.cache.Set(nsName, ipStr, 5*time.Minute)
			}
			return ipStr, nil
		}
	}

	return "", errors.New("no IPv4 address found")
}

