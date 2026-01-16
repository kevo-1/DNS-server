package server

import (
	"DNS-server/internal/protocol"
	"DNS-server/models"
	"DNS-server/pkg/resolver"
	"fmt"
	"log"
)

type Handler struct {
	resolver *resolver.Resolver
	config   *Config
}

func NewHandler(config *Config) *Handler {
	return &Handler{
		resolver: resolver.GetInstance(),
		config:   config,
	}
}

func (h *Handler) HandleRequest(data []byte) ([]byte, error) {
	request, err := protocol.ParseMessage(data)
	if err != nil {
		log.Printf("Failed to parse DNS request: %v", err)
		return nil, fmt.Errorf("parse request: %w", err)
	}

	if len(request.Questions) > 0 {
		q := request.Questions[0]
		log.Printf("DNS Query: %s (Type: %s, Class: %s)",
			q.Name,
			protocol.TypeToString(q.Type),
			protocol.ClassToString(q.Class))
	}

	var response *protocol.Message
	if h.config.EnableRecursion {
		response = h.handleRecursiveRequest(request)
	} else {
		response = h.handleIterativeRequest(request)
	}

	responseData, err := protocol.BuildMessage(response)
	if err != nil {
		log.Printf("Failed to build DNS response: %v", err)
		return nil, fmt.Errorf("build response: %w", err)
	}

	log.Printf("DNS Response: %d answers, RCODE: %s",
		len(response.Answers),
		protocol.RCodeToString(response.Header.Flags&0x0F))

	return responseData, nil
}

func (h *Handler) handleRecursiveRequest(request *protocol.Message) *protocol.Message {
	if len(request.Questions) == 0 {
		return protocol.CreateErrorResponse(request, protocol.RCodeFormErr)
	}

	question := request.Questions[0]

	if question.Type != protocol.TypeA {
		return protocol.CreateErrorResponse(request, protocol.RCodeNotImpl)
	}

	ip, err := h.resolver.ResolveA(question.Name)
	if err != nil {
		log.Printf("Resolution failed for %s: %v", question.Name, err)
		return protocol.CreateErrorResponse(request, protocol.RCodeServFail)
	}

	aRecord, err := protocol.CreateARecord(question.Name, ip, 300)
	if err != nil {
		log.Printf("Failed to create A record: %v", err)
		return protocol.CreateErrorResponse(request, protocol.RCodeServFail)
	}

	response := protocol.CreateResponse(request, []protocol.ResourceRecord{aRecord})
	response.Header.Flags |= protocol.FlagRA

	return response
}

func (h *Handler) handleIterativeRequest(request *protocol.Message) *protocol.Message {
	return protocol.CreateErrorResponse(request, protocol.RCodeNotImpl)
}

func (h *Handler) HandleError(request *protocol.Message, rcode uint16) ([]byte, error) {
	response := protocol.CreateErrorResponse(request, rcode)
	return protocol.BuildMessage(response)
}

func (h *Handler) GetStats() models.CacheStatistics {
	return h.resolver.GetStats()
}
