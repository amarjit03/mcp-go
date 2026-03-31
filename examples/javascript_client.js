#!/usr/bin/env node
/**
 * JavaScript client for MCP Dev Assistant Server
 * Demonstrates how to communicate with the server using the MCP protocol
 */

const net = require('net');
const readline = require('readline');

class MCPClient {
  constructor(host = '127.0.0.1', port = 9090) {
    this.host = host;
    this.port = port;
    this.requestId = 0;
  }

  getNextId() {
    return ++this.requestId;
  }

  sendRequest(method, params = {}) {
    return new Promise((resolve, reject) => {
      try {
        const socket = net.createConnection({ 
          port: this.port, 
          host: this.host 
        });

        const request = {
          jsonrpc: '2.0',
          method: method,
          params: params,
          id: this.getNextId()
        };

        socket.write(JSON.stringify(request));

        let data = '';
        socket.on('data', (chunk) => {
          data += chunk.toString();
        });

        socket.on('end', () => {
          try {
            const response = JSON.parse(data);
            socket.destroy();
            resolve(response);
          } catch (e) {
            reject(new Error(`Failed to parse response: ${e.message}`));
          }
        });

        socket.on('error', (err) => {
          socket.destroy();
          reject(err);
        });

        socket.setTimeout(10000, () => {
          socket.destroy();
          reject(new Error('Request timeout'));
        });

      } catch (err) {
        reject(err);
      }
    });
  }

  initialize(clientName = 'js-client') {
    return this.sendRequest('initialize', {
      protocolVersion: '2024-11-05',
      clientInfo: {
        name: clientName,
        version: '1.0'
      }
    });
  }

  listTools() {
    return this.sendRequest('tools/list');
  }

  callTool(toolName, arguments) {
    return this.sendRequest('tools/call', {
      name: toolName,
      arguments: arguments
    });
  }

  // Convenience methods
  readFile(path) {
    return this.callTool('read_file', { path });
  }

  writeFile(path, content, append = false) {
    return this.callTool('write_file', { path, content, append });
  }

  listDirectory(path) {
    return this.callTool('list_directory', { path });
  }

  executeCommand(command) {
    return this.callTool('execute_command', { command });
  }

  getCPUUsage() {
    return this.callTool('get_cpu_usage', {});
  }

  getMemoryUsage() {
    return this.callTool('get_memory_usage', {});
  }

  checkPort(port, protocol = 'tcp') {
    return this.callTool('check_port', { port, protocol });
  }

  getProcessInfo(pidOrName) {
    return this.callTool('get_process_info', { pid_or_name: pidOrName });
  }

  healthCheck(ports = []) {
    return this.callTool('health_check', { ports });
  }

  readLogs(path, lines = 50) {
    return this.callTool('read_logs', { path, lines });
  }
}

function prettyPrint(obj, indent = 0) {
  const prefix = '  '.repeat(indent);
  
  if (Array.isArray(obj)) {
    obj.forEach((item, i) => {
      console.log(`${prefix}[${i}]:`);
      prettyPrint(item, indent + 1);
    });
  } else if (typeof obj === 'object' && obj !== null) {
    Object.entries(obj).forEach(([key, value]) => {
      if (typeof value === 'object' && value !== null) {
        console.log(`${prefix}${key}:`);
        prettyPrint(value, indent + 1);
      } else {
        console.log(`${prefix}${key}: ${value}`);
      }
    });
  } else {
    console.log(`${prefix}${obj}`);
  }
}

async function main() {
  console.log('='.repeat(60));
  console.log('MCP Dev Assistant - JavaScript Client Demo');
  console.log('='.repeat(60));
  console.log();

  const client = new MCPClient();

  try {
    // Initialize
    console.log('[1] Initializing connection...');
    let response = await client.initialize();
    if (response.error) {
      console.error(`Error: ${response.error.message}`);
      process.exit(1);
    }

    console.log(`✓ Connected to ${response.result.serverInfo.name}`);
    console.log(`  Version: ${response.result.serverInfo.version}`);
    console.log(`  Protocol: ${response.result.protocolVersion}`);
    console.log();

    // List tools
    console.log('[2] Listing available tools...');
    response = await client.listTools();
    if (!response.error) {
      const tools = response.result.tools;
      console.log(`✓ Found ${tools.length} tools:`);
      tools.forEach(tool => {
        console.log(`   - ${tool.name}: ${tool.description}`);
      });
    }
    console.log();

    // System health check
    console.log('[3] Performing system health check...');
    response = await client.healthCheck([22, 80, 443]);
    if (!response.error) {
      const result = response.result.content[0].data;
      const cpu = result.cpu;
      const mem = result.memory;
      const status = result.status;

      console.log(`✓ System Status: ${status.toUpperCase()}`);
      console.log(`  CPU: ${cpu.percent.toFixed(1)}% (${cpu.count} cores)`);
      console.log(`  Memory: ${mem.usedPercent.toFixed(1)}% (${(mem.used / (1024**3)).toFixed(1)}GB / ${(mem.total / (1024**3)).toFixed(1)}GB)`);

      if (result.ports && result.ports.length > 0) {
        console.log('  Ports:');
        result.ports.forEach(port => {
          console.log(`    ${port.port}/${port.protocol}: ${port.state}`);
        });
      }
    }
    console.log();

    // Execute command
    console.log('[4] Executing command: "uname -a"');
    response = await client.executeCommand('uname -a');
    if (!response.error) {
      const result = response.result.content[0].data;
      if (result.success) {
        console.log(`✓ Output: ${result.output.trim()}`);
      } else {
        console.log(`✗ Error: ${result.error}`);
      }
    }
    console.log();

    // Try to read a file
    console.log('[5] Reading file: "/etc/hostname"');
    response = await client.readFile('/etc/hostname');
    if (!response.error) {
      const result = response.result.content[0].data;
      console.log(`✓ Content: ${result.content.trim()}`);
    } else {
      console.log(`✗ Error: ${response.error.message}`);
    }
    console.log();

    console.log('='.repeat(60));
    console.log('Demo completed!');
    console.log('='.repeat(60));

  } catch (err) {
    console.error('Error:', err.message);
    process.exit(1);
  }
}

main();
