package mcp

import (
	"encoding/json"
	"fmt"
)

// Protocol handles MCP protocol encoding/decoding
type Protocol struct {
	requestHandlers map[RequestMethod]RequestHandler
}

// RequestHandler is a function that handles a specific request
type RequestHandler func(params json.RawMessage) (interface{}, error)

// NewProtocol creates a new MCP protocol handler
func NewProtocol() *Protocol {
	return &Protocol{
		requestHandlers: make(map[RequestMethod]RequestHandler),
	}
}

// RegisterHandler registers a handler for a request method
func (p *Protocol) RegisterHandler(method RequestMethod, handler RequestHandler) {
	p.requestHandlers[method] = handler
}

// HandleRequest processes an incoming request and returns a response
func (p *Protocol) HandleRequest(data []byte) ([]byte, error) {
	var req Request
	if err := json.Unmarshal(data, &req); err != nil {
		return p.encodeError(nil, ErrorCodeParse, "Parse error", ""), nil
	}

	if req.JSONRPC != "2.0" {
		return p.encodeError(req.ID, ErrorCodeInvalidReq, "Invalid Request", ""), nil
	}

	if req.Method == "" {
		return p.encodeError(req.ID, ErrorCodeInvalidReq, "Invalid Request", ""), nil
	}

	handler, ok := p.requestHandlers[req.Method]
	if !ok {
		return p.encodeError(req.ID, ErrorCodeNotFound, fmt.Sprintf("Method not found: %s", req.Method), ""), nil
	}

	result, err := handler(req.Params)
	if err != nil {
		return p.encodeError(req.ID, ErrorCodeInternal, err.Error(), ""), nil
	}

	return p.encodeResponse(req.ID, result), nil
}

// encodeError creates an error response
func (p *Protocol) encodeError(id interface{}, code int, message string, data string) []byte {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Error: &Error{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	respData, _ := json.Marshal(resp)
	return respData
}

// encodeResponse creates a success response
func (p *Protocol) encodeResponse(id interface{}, result interface{}) []byte {
	resp := Response{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
	respData, _ := json.Marshal(resp)
	return respData
}
