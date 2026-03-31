#!/usr/bin/env python3
"""
Python client for MCP Dev Assistant Server
Demonstrates how to communicate with the server using the MCP protocol
"""

import socket
import json
import sys
from typing import Dict, Any, Optional

class MCPClient:
    """Client for communicating with MCP server"""
    
    def __init__(self, host: str = "127.0.0.1", port: int = 9090):
        self.host = host
        self.port = port
        self.request_id = 0
    
    def _get_next_id(self) -> int:
        self.request_id += 1
        return self.request_id
    
    def _send_request(self, method: str, params: Optional[Dict] = None) -> Dict[str, Any]:
        """Send a request to the MCP server and get response"""
        try:
            sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            sock.settimeout(10)
            sock.connect((self.host, self.port))
            
            request = {
                "jsonrpc": "2.0",
                "method": method,
                "params": params or {},
                "id": self._get_next_id()
            }
            
            # Send request
            request_json = json.dumps(request)
            sock.sendall(request_json.encode())
            
            # Receive response
            response_data = b''
            sock.settimeout(5)
            try:
                while True:
                    chunk = sock.recv(4096)
                    if not chunk:
                        break
                    response_data += chunk
            except socket.timeout:
                pass
            
            sock.close()
            
            return json.loads(response_data.decode())
        
        except Exception as e:
            return {
                "jsonrpc": "2.0",
                "error": {
                    "code": -1,
                    "message": str(e)
                },
                "id": self.request_id
            }
    
    def initialize(self, client_name: str = "python-client") -> Dict[str, Any]:
        """Initialize connection with server"""
        return self._send_request("initialize", {
            "protocolVersion": "2024-11-05",
            "clientInfo": {
                "name": client_name,
                "version": "1.0"
            }
        })
    
    def list_tools(self) -> Dict[str, Any]:
        """List available tools"""
        return self._send_request("tools/list")
    
    def call_tool(self, tool_name: str, arguments: Dict) -> Dict[str, Any]:
        """Call a specific tool"""
        return self._send_request("tools/call", {
            "name": tool_name,
            "arguments": arguments
        })
    
    # Convenience methods for common tools
    
    def read_file(self, path: str) -> Dict[str, Any]:
        """Read a file"""
        return self.call_tool("read_file", {"path": path})
    
    def write_file(self, path: str, content: str, append: bool = False) -> Dict[str, Any]:
        """Write to a file"""
        return self.call_tool("write_file", {
            "path": path,
            "content": content,
            "append": append
        })
    
    def list_directory(self, path: str) -> Dict[str, Any]:
        """List directory contents"""
        return self.call_tool("list_directory", {"path": path})
    
    def execute_command(self, command: str) -> Dict[str, Any]:
        """Execute a shell command"""
        return self.call_tool("execute_command", {"command": command})
    
    def get_cpu_usage(self) -> Dict[str, Any]:
        """Get CPU usage"""
        return self.call_tool("get_cpu_usage", {})
    
    def get_memory_usage(self) -> Dict[str, Any]:
        """Get memory usage"""
        return self.call_tool("get_memory_usage", {})
    
    def check_port(self, port: int, protocol: str = "tcp") -> Dict[str, Any]:
        """Check if a port is open"""
        return self.call_tool("check_port", {
            "port": port,
            "protocol": protocol
        })
    
    def get_process_info(self, pid_or_name: str) -> Dict[str, Any]:
        """Get process information"""
        return self.call_tool("get_process_info", {"pid_or_name": pid_or_name})
    
    def health_check(self, ports: Optional[list] = None) -> Dict[str, Any]:
        """Perform system health check"""
        return self.call_tool("health_check", {
            "ports": ports or []
        })
    
    def read_logs(self, path: str, lines: int = 50) -> Dict[str, Any]:
        """Read log file"""
        return self.call_tool("read_logs", {
            "path": path,
            "lines": lines
        })

def pretty_print(obj: Any, indent: int = 0) -> None:
    """Pretty print JSON objects"""
    prefix = "  " * indent
    if isinstance(obj, dict):
        for key, value in obj.items():
            if isinstance(value, (dict, list)):
                print(f"{prefix}{key}:")
                pretty_print(value, indent + 1)
            else:
                print(f"{prefix}{key}: {value}")
    elif isinstance(obj, list):
        for i, item in enumerate(obj):
            print(f"{prefix}[{i}]:")
            pretty_print(item, indent + 1)
    else:
        print(f"{prefix}{obj}")

def main():
    """Example usage of MCP client"""
    
    print("=" * 60)
    print("MCP Dev Assistant - Python Client Demo")
    print("=" * 60)
    print()
    
    client = MCPClient()
    
    # Initialize
    print("[1] Initializing connection...")
    response = client.initialize()
    if "error" in response:
        print(f"Error: {response['error']['message']}")
        sys.exit(1)
    
    print(f"✓ Connected to {response['result']['serverInfo']['name']}")
    print(f"  Version: {response['result']['serverInfo']['version']}")
    print(f"  Protocol: {response['result']['protocolVersion']}")
    print()
    
    # List tools
    print("[2] Listing available tools...")
    response = client.list_tools()
    if "error" not in response:
        tools = response['result']['tools']
        print(f"✓ Found {len(tools)} tools:")
        for tool in tools:
            print(f"   - {tool['name']}: {tool['description']}")
    print()
    
    # System health check
    print("[3] Performing system health check...")
    response = client.health_check(ports=[22, 80, 443])
    if "error" not in response:
        result = response['result']['content'][0]['data']
        cpu = result['cpu']
        mem = result['memory']
        status = result['status']
        
        print(f"✓ System Status: {status.upper()}")
        print(f"  CPU: {cpu['percent']:.1f}% ({cpu['count']} cores)")
        print(f"  Memory: {mem['usedPercent']:.1f}% ({mem['used'] / (1024**3):.1f}GB / {mem['total'] / (1024**3):.1f}GB)")
        
        if result.get('ports'):
            print("  Ports:")
            for port in result['ports']:
                print(f"    {port['port']}/{port['protocol']}: {port['state']}")
    print()
    
    # Execute command
    print("[4] Executing command: 'uname -a'")
    response = client.execute_command("uname -a")
    if "error" not in response:
        result = response['result']['content'][0]['data']
        if result['success']:
            print(f"✓ Output: {result['output'].strip()}")
        else:
            print(f"✗ Error: {result['error']}")
    print()
    
    # Try to read a file
    print("[5] Reading file: '/etc/hostname'")
    response = client.read_file("/etc/hostname")
    if "error" not in response:
        result = response['result']['content'][0]['data']
        print(f"✓ Content: {result['content'].strip()}")
    else:
        print(f"✗ Error: {response['error']['message']}")
    print()
    
    print("=" * 60)
    print("Demo completed!")
    print("=" * 60)

if __name__ == "__main__":
    main()
