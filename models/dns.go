package models

type NameServer struct {
	Name string
	IPv4 string
	IPv6 string
}

type DNSQuery struct {
	Domain    string
	QueryType uint16
}

type DNSRecord struct {
	Name  string
	Type  uint16
	Value string
	TTL   uint32
}