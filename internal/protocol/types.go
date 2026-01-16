package protocol

const (
	// Query Types
	TypeA     = 1   // IPv4 address
	TypeNS    = 2   // Name server
	TypeCNAME = 5   // Canonical name
	TypeSOA   = 6   // Start of authority
	TypePTR   = 12  // Domain name pointer
	TypeMX    = 15  // Mail exchange
	TypeTXT   = 16  // Text strings
	TypeAAAA  = 28  // IPv6 address
	TypeSRV   = 33  // Service locator
	
	// Classes
	ClassIN = 1  // Internet
	ClassCS = 2  // CSNET (obsolete)
	ClassCH = 3  // CHAOS
	ClassHS = 4  // Hesiod
	
	// Response Codes (RCODE)
	RCodeNoError  = 0 // No error
	RCodeFormErr  = 1 // Format error
	RCodeServFail = 2 // Server failure
	RCodeNXDomain = 3 // Non-existent domain
	RCodeNotImpl  = 4 // Not implemented
	RCodeRefused  = 5 // Query refused
	
	// Header Flags
	FlagQR = 1 << 15 // Query (0) / Response (1)
	FlagAA = 1 << 10 // Authoritative Answer
	FlagTC = 1 << 9  // Truncated
	FlagRD = 1 << 8  // Recursion Desired
	FlagRA = 1 << 7  // Recursion Available
	FlagZ  = 1 << 6  // Reserved (must be zero)
	FlagAD = 1 << 5  // Authenticated Data
	FlagCD = 1 << 4  // Checking Disabled

    // OpCode values (bits 11-14 of flags)
	OpCodeQuery  = 0 // Standard query
	OpCodeIQuery = 1 // Inverse query (obsolete)
	OpCodeStatus = 2 // Server status request
)

// Type -> string
func TypeToString(t uint16) string {
	switch t {
	case TypeA:
		return "A"
	case TypeNS:
		return "NS"
	case TypeCNAME:
		return "CNAME"
	case TypeSOA:
		return "SOA"
	case TypePTR:
		return "PTR"
	case TypeMX:
		return "MX"
	case TypeTXT:
		return "TXT"
	case TypeAAAA:
		return "AAAA"
	case TypeSRV:
		return "SRV"
	default:
		return "UNKNOWN"
	}
}

// Class -> string
func ClassToString(c uint16) string {
	switch c {
	case ClassIN:
		return "IN"
	case ClassCS:
		return "CS"
	case ClassCH:
		return "CH"
	case ClassHS:
		return "HS"
	default:
		return "UNKNOWN"
	}
}

// Response code -> string
func RCodeToString(rcode uint16) string {
	switch rcode {
	case RCodeNoError:
		return "NOERROR"
	case RCodeFormErr:
		return "FORMERR"
	case RCodeServFail:
		return "SERVFAIL"
	case RCodeNXDomain:
		return "NXDOMAIN"
	case RCodeNotImpl:
		return "NOTIMPL"
	case RCodeRefused:
		return "REFUSED"
	default:
		return "UNKNOWN"
	}
}