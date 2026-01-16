package transport

import (
	"context"
	"net"
)

type UDPTransport struct {
	addr    string
	handler func([]byte) ([]byte, error)
}

func NewUDPTransport(addr string, handler func([]byte) ([]byte, error)) *UDPTransport {
	return &UDPTransport{
		addr:    addr,
		handler: handler,
	}
}

func (s *UDPTransport) Start(ctx context.Context) error {
	conn, err := net.ListenPacket("udp", s.addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	
	go func() {
		<-ctx.Done()
		conn.Close()
	}()
	
	buffer := make([]byte, 512)
	
	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				continue
			}
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
