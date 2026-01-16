package transport

import (
	"context"
	"net"
)

type UDPServer struct {
	addr    string
	handler func([]byte) ([]byte, error)
}

func NewUDPServer(addr string, handler func([]byte) ([]byte, error)) *UDPServer {
	return &UDPServer{
		addr:    addr,
		handler: handler,
	}
}

func (s *UDPServer) ListenAndServe(ctx context.Context) error {
	conn, err := net.ListenPacket("udp", s.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	buffer := make([]byte, 512)
	
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			n, addr, err := conn.ReadFrom(buffer)
			if err != nil {
				continue
			}
			
			go func(data []byte, clientAddr net.Addr) {
				response, err := s.handler(data)
				if err != nil {
					return
				}
				conn.WriteTo(response, clientAddr)
			}(append([]byte(nil), buffer[:n]...), addr)
		}
	}
}