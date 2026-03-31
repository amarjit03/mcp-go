# MCP Protocol Specification

## Overview

The Model Context Protocol (MCP) is a standardized protocol for AI systems to interact with external tools and resources. This document describes the MCP implementation in this Go server.

## Protocol Basics

### Transport Layer
- **Protocol**: TCP/IP
- **Message Format**: JSON-RPC 2.0
- **Encoding**: UTF-8

### Connection Flow

```
1. Client connects to server
2. Client sends "initialize" request
3. Server responds with capabilities and available tools
4. Client can now call tools
5. Server executes tools and returns results
6. Connection can be kept alive for multiple requests
```

## Message Structure

### Request

```json
{
  "jsonrpc": "2.0",
  "method": "<method-name>",
  "params": { ... },
  "id": <unique-integer>
}
```

- **jsonrpc**: Always "2.0"
- **method**: One of: `initialize`, `tools/list`, `tools/call`
- **params**: Method-specific parameters
- **id**: Unique identifier for matching response

### Response

```json
{
  "jsonrpc": "2.0",
  "result": { ... },
  "id": <same-as-request>
}
```

Or on error:

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": <error-code>,
    "message": "<error-message>",
    "data": "<optional-data>"
  },
  "id": <same-as-request>
}
```

## Methods

### 1. initialize

Initialize the connection and retrieve server capabilities.

**Request**:
```json
{
  "jsonrpc": "2.0",
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": { },
    "clientInfo": {
      "name": "string",
      "version": "string"
    }
  },
  "id": 1
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": true
    },
    "serverInfo": {
      "name": "string",
      "version": "string"
    },
    "tools": [ ]
  },
  "id": 1
}
```

### 2. tools/list

List all available tools.

**Request**:
```json
{
  "jsonrpc": "2.0",
  "method": "tools/list",
  "id": 2
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "tools": [
      {
        "name": "string",
        "description": "string",
        "inputSchema": {
          "type": "object",
          "properties": { },
          "required": [ ]
        }
      }
    ]
  },
  "id": 2
}
```

### 3. tools/call

Execute a specific tool.

**Request**:
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "tool-name",
    "arguments": { ... }
  },
  "id": 3
}
```

**Response**:
```json
{
  "jsonrpc": "2.0",
  "result": {
    "content": [
      {
        "type": "text",
        "data": { ... }
      }
    ]
  },
  "id": 3
}
```

## Error Codes

| Code | Message | Meaning |
|------|---------|---------|
| -32700 | Parse error | Invalid JSON received |
| -32600 | Invalid Request | JSON-RPC request is malformed |
| -32601 | Method not found | Method does not exist |
| -32602 | Invalid params | Invalid method parameters |
| -32603 | Internal error | Internal server error |
| -32000 | Server error | Server-specific error |

## Tool Definition Format

Each tool must define:

- **name**: Unique identifier
- **description**: Human-readable description
- **inputSchema**: JSON Schema describing parameters
  - **type**: Always "object"
  - **properties**: Map of parameter names to types
  - **required**: List of required parameters

### Example Tool Definition

```json
{
  "name": "read_file",
  "description": "Read the contents of a file",
  "inputSchema": {
    "type": "object",
    "properties": {
      "path": {
        "type": "string",
        "description": "The file path to read"
      }
    },
    "required": ["path"]
  }
}
```

## Implementation Details

### Error Handling

All errors should be returned as JSON-RPC error responses:

```json
{
  "jsonrpc": "2.0",
  "error": {
    "code": -32603,
    "message": "Tool execution failed",
    "data": "Additional error details..."
  },
  "id": 123
}
```

### Response Content Structure

Tool responses are wrapped in a `content` array:

```json
{
  "content": [
    {
      "type": "text",
      "data": {
        "key": "value"
      }
    }
  ]
}
```

### Timeouts

- Default request timeout: 30 seconds
- Configurable per tool

### Output Size Limiting

- File read limit: 50 MB (configurable)
- Command output limit: 10 MB (configurable)

## Security Considerations

### 1. Input Validation

All inputs are validated:
- Parameter types checked
- File paths normalized and restricted
- Command arguments sanitized

### 2. Path Security

Paths are:
- Converted to absolute paths
- Checked against allowed/denied lists
- Verified to exist and be accessible

### 3. Command Security

Commands are:
- Parsed from shell-safe execution
- Checked against allow/deny lists
- Executed with timeout protection
- Output size-limited

### 4. Rate Limiting

Consider implementing:
- Connection limits
- Request rate limiting
- Resource usage monitoring

## Future Extensions

Possible protocol extensions:

1. **Streaming Responses** - For large file/command output
2. **Resource Definition** - Read external data sources
3. **Prompt Templates** - Provide AI with context
4. **Long-running Tasks** - Async task support
5. **Authentication** - JWT/OAuth support

## References

- [JSON-RPC 2.0 Specification](https://www.jsonrpc.org/specification)
- [Model Context Protocol](https://modelcontextprotocol.io/)
