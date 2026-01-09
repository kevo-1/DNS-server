package data

type NameServer struct {
	name string
	IPv4 string
	IPv6 string
}

var RootServers = []NameServer{
	{name: "a.root-servers.net", IPv4: "198.41.0.4", IPv6: "2001:503:ba3e::2:30"},
	{name: "b.root-servers.net", IPv4: "170.247.170.2", IPv6: "2801:1b8:10::b"},
	{name: "c.root-servers.net", IPv4: "192.33.4.12", IPv6: "2001:500:2::c"},
	{name: "d.root-servers.net", IPv4: "199.7.91.13", IPv6: "2001:500:2d::d"},
	{name: "e.root-servers.net", IPv4: "192.203.230.10", IPv6: "2001:500:a8::e"},
	{name: "f.root-servers.net", IPv4: "192.5.5.241", IPv6: "2001:500:2f::f"},
	{name: "g.root-servers.net", IPv4: "192.112.36.4", IPv6: "2001:500:12::d0d"},
	{name: "h.root-servers.net", IPv4: "198.97.190.53", IPv6: "2001:500:1::53"},
	{name: "i.root-servers.net", IPv4: "192.36.148.17", IPv6: "2001:7fe::53"},
	{name: "j.root-servers.net", IPv4: "192.58.128.30", IPv6: "2001:503:c27::2:30"},
	{name: "k.root-servers.net", IPv4: "193.0.14.129", IPv6: "2001:7fd::1"},
	{name: "l.root-servers.net", IPv4: "199.7.83.42", IPv6: "2001:500:9f::42"},
	{name: "m.root-servers.net", IPv4: "202.12.27.33", IPv6: "2001:dc3::35"},
}