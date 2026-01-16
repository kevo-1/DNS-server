package transport

import (
	"context"
	"io"
	"log"
	"net"
)

type TCPServer struct {
	addr    string
	handler func([]byte) ([]byte, error)
}

func NewTCPServer(addr string, handler func([]byte) ([]byte, error)) *TCPServer {
	return &TCPServer{
		addr:    addr,
		handler: handler,
	}
}

func (s *TCPServer) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("TCP DNS server listening on %s", s.addr)

	done := make(chan struct{})

	go func() {
		<-ctx.Done()
		listener.Close()
		close(done)
	}()

	for {
		select {
		case <-done:
			return nil
		default:
			conn, err := listener.Accept()
			if err != nil {
				select {
				case <-done:
					return nil
				default:
					log.Printf("TCP accept error: %v", err)
					continue
				}
			}

			go s.handleConnection(conn)
		}
	}
}

func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	lengthBuf := make([]byte, 2)

	for {
		if _, err := io.ReadFull(conn, lengthBuf); err != nil {
			if err != io.EOF {
				log.Printf("TCP read length error: %v", err)
			}
			return
		}

		msgLen := int(lengthBuf[0])<<8 | int(lengthBuf[1])

		if msgLen == 0 || msgLen > 65535 {
			log.Printf("Invalid message length: %d", msgLen)
			return
		}

		msgBuf := make([]byte, msgLen)
		if _, err := io.ReadFull(conn, msgBuf); err != nil {
			log.Printf("TCP read message error: %v", err)
			return
		}

		response, err := s.handler(msgBuf)
		if err != nil {
			log.Printf("TCP handler error: %v", err)
			return
		}

		if err := s.writeMessage(conn, response); err != nil {
			log.Printf("TCP write error: %v", err)
			return
		}
	}
}

func (s *TCPServer) writeMessage(conn net.Conn, data []byte) error {
	length := len(data)
	lengthBuf := []byte{byte(length >> 8), byte(length & 0xFF)}

	if _, err := conn.Write(lengthBuf); err != nil {
		return err
	}

	if _, err := conn.Write(data); err != nil {
		return err
	}

	return nil
}