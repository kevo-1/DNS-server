# DNS Server

![Tech Stack](https://skillicons.dev/icons?i=go)

**Recursive DNS Server with Caching and Iterative Resolution**

A complete DNS server implementation in Go featuring thread-safe LRU caching, full iterative resolution from root servers, and support for both UDP and TCP transports.

> **Note:** This is a fully functional DNS server suitable for learning and local testing.

---

## Features

-   **Iterative Resolution** – Queries root → TLD → authoritative servers
-   **LRU Cache with TTL** – Thread-safe caching with automatic expiration
-   **Dual Transport** – Both UDP (port 53) and TCP support
-   **Graceful Shutdown** – Clean resource cleanup with statistics reporting
---

## Project Evolution

### Initial Structure (Simple)

```
DNS-server/
├── main.go
├── go.mod
│
├── core/
│   ├── resolver
│   └── parser
│
├── data/
│   └── name_server.go
│
├── models/
│   └── url.go
│
└── tests/
    └── parser_test.go
```

### Final Structure (Organized)

```
DNS-server/
├── main.go
├── go.mod
│
├── internal/
│   ├── protocol/
│   │   ├── message.go
│   │   ├── parser.go
│   │   ├── builder.go
│   │   ├── types.go
│   │   └── helpers.go
│   │
│   ├── server/
│   │   ├── server.go
│   │   ├── handler.go
│   │   └── config.go
│   │
│   └── transport/
│       ├── udp.go
│       └── tcp.go
│
├── pkg/
│   ├── resolver/
│   │   ├── resolver.go
│   │   ├── cache.go
│   │   └── iterative.go
│   │
│   └── parser/
│       └── url_parser.go
│
├── data/
│   ├── root_servers.go
│   ├── root_manager.go
│   └── helpers.go
│
├── models/
│   ├── url.go
│   ├── dns.go
│   └── cache.go
│
└── tests/
    ├── parser_test.go
    ├── resolver_test.go
    └── integration_test.go
```

---

## Quick Start

### Prerequisites

-   **Go 1.24+** – [Download Go](https://go.dev/dl/)
-   **Administrator privileges** (for port 53) or use custom port

### Running the Server

```bash
# Clone the repository
git clone <repository-url>
cd DNS-server

# Run directly
go run main.go

# Or build and run
go build -o dns-server.exe .
.\dns-server.exe
```

### Testing DNS Queries

```bash
# Using nslookup (Windows built-in)
nslookup example.com localhost

# Using dig (requires installation)
dig @localhost example.com
```

---

## Configuration

Default settings in `internal/server/config.go`:

| Setting        | Default      | Description              |
| -------------- | ------------ | ------------------------ |
| **UDP Port**   | 53           | DNS UDP listener port    |
| **TCP Port**   | 53           | DNS TCP listener port    |
| **Host**       | 0.0.0.0      | Listen on all interfaces |
| **Cache Size** | 1000 entries | Maximum cached domains   |
| **Cache TTL**  | 5 minutes    | Default time-to-live     |
| **Recursion**  | Enabled      | Perform full resolution  |

### Using a Custom Port

If port 53 requires admin rights, modify `DefaultConfig()`:

```go
func DefaultConfig() *Config {
    return &Config{
        UDPPort: 8053,  // Custom port
        TCPPort: 8053,
        ...
    }
}
```

Then query.

---

## Architecture

### DNS Resolution Flow

```
Client Query → Transport (UDP/TCP)
    ↓
Request Handler
    ↓
Cache Lookup → Hit? → Return cached answer
    ↓ Miss
Iterative Resolver
    ↓
Root Servers → TLD Servers → Authoritative Servers
    ↓
Cache Result
    ↓
Build Response → Transport → Client
```

### Key Components

-   **Protocol Package** – DNS message parsing and building (RFC 1035)
-   **Resolver Package** – Iterative DNS resolution with caching
-   **Transport Package** – UDP and TCP network handlers
-   **Server Package** – Request orchestration and lifecycle management

### Supported Record Types

-   **A** (IPv4 addresses)
-   **CNAME** (Canonical names)
-   **NS** (Nameserver records)

---

## Performance Features

-   **Thread-Safe Cache** – Concurrent read/write with RWMutex (which was heavly inspired by my OS course)
-   **LRU Caching** – Automatic removal of least-used entries
-   **TTL Management** – Background cleanup of expired entries (which I learned about in my Networks course lab)
-   **Connection Pooling** – Efficient upstream queries

---

## Example Output

**Startup:**

```
(Date Time) DNS Server starting...
(Date Time) UDP server listening on 0.0.0.0:53
(Date Time) TCP server listening on 0.0.0.0:53
(Date Time) DNS server started successfully
```

**Query Processing:**

```
(Date Time) DNS Query: example.com (Type: A, Class: IN)
(Date Time) DNS Response: 1 answers, RCODE: NOERROR
```

**Shutdown Statistics:**

```
Server Statistics:
  Cache Hits: 45
  Cache Misses: 12
  Cache Hit Rate: 78.95%
  Cache Evictions: 0
  Total Entries: 12/1000
```

---

## Troubleshooting

### Port 53 Already in Use

```
listen udp :53: bind: Only one usage of each socket address...
```

**Solution:** Stop other DNS services or use a custom port (8053, 5353, etc.)

### Permission Denied

```
listen udp :53: bind: permission denied
```

**Solution:** Run as Administrator or use port > 1024

### Slow Queries

-   Check upstream network connectivity
-   Verify root server accessibility
-   Review cache hit rate in statistics

---

## Dependencies

-   **Go Standard Library** – `net`, `context`, `sync`, `container/list`
-   **No external dependencies** – Pure Go implementation :)

---

## Enhancements & Improvements

Please let me know if you have any suggestions or improvements through GitHub issues or pull requests!

---

**Built with Go**
