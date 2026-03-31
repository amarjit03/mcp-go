package system

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

// SystemMonitor provides system health information
type SystemMonitor struct{}

// NewSystemMonitor creates a new system monitor
func NewSystemMonitor() *SystemMonitor {
	return &SystemMonitor{}
}

// CPUInfo represents CPU usage information
type CPUInfo struct {
	Percent      float64 `json:"percent"`
	Count        int     `json:"count"`
	CountLogical int     `json:"countLogical"`
	Timestamp    string  `json:"timestamp"`
}

// GetCPUUsage returns current CPU usage
func (sm *SystemMonitor) GetCPUUsage() (*CPUInfo, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU percent: %w", err)
	}

	count, err := cpu.Counts(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get CPU count: %w", err)
	}

	countLogical, err := cpu.Counts(true)
	if err != nil {
		countLogical = count
	}

	return &CPUInfo{
		Percent:      percent[0],
		Count:        count,
		CountLogical: countLogical,
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// MemoryInfo represents memory usage information
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Available   uint64  `json:"available"`
	UsedPercent float64 `json:"usedPercent"`
	Timestamp   string  `json:"timestamp"`
}

// GetMemoryUsage returns current memory usage
func (sm *SystemMonitor) GetMemoryUsage() (*MemoryInfo, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("failed to get memory info: %w", err)
	}

	return &MemoryInfo{
		Total:       v.Total,
		Used:        v.Used,
		Free:        v.Free,
		Available:   v.Available,
		UsedPercent: v.UsedPercent,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// ProcessInfo represents process information
type ProcessInfo struct {
	PID       int32   `json:"pid"`
	Name      string  `json:"name"`
	Status    string  `json:"status"`
	Percent   float64 `json:"percent"`
	MemoryRSS uint64  `json:"memoryRss"`
	Timestamp string  `json:"timestamp"`
}

// GetProcessInfo returns information about a running process
func (sm *SystemMonitor) GetProcessInfo(pidOrName interface{}) (*ProcessInfo, error) {
	var p *process.Process
	var err error

	switch v := pidOrName.(type) {
	case int32:
		p, err = process.NewProcess(v)
	case string:
		// Try to parse as PID first
		if pid, err := strconv.ParseInt(v, 10, 32); err == nil {
			p, err = process.NewProcess(int32(pid))
		} else {
			// Try to find by name
			processes, err := process.Processes()
			if err != nil {
				return nil, fmt.Errorf("failed to get processes: %w", err)
			}

			for _, proc := range processes {
				name, err := proc.Name()
				if err == nil && strings.Contains(name, v) {
					p = proc
					break
				}
			}
			if p == nil {
				return nil, fmt.Errorf("process not found: %s", v)
			}
		}
	default:
		return nil, fmt.Errorf("invalid PID type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get process: %w", err)
	}

	name, _ := p.Name()
	status, _ := p.Status()
	percent, _ := p.CPUPercent()
	memInfo, _ := p.MemoryInfo()

	return &ProcessInfo{
		PID:       p.Pid,
		Name:      name,
		Status:    strings.Join(status, ","),
		Percent:   percent,
		MemoryRSS: memInfo.RSS,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}, nil
}

// PortStatus represents the status of a network port
type PortStatus struct {
	Port      int    `json:"port"`
	Protocol  string `json:"protocol"`
	State     string `json:"state"`
	Timestamp string `json:"timestamp"`
}

// CheckPort checks if a port is open and listening
func (sm *SystemMonitor) CheckPort(port int, protocol string) *PortStatus {
	if protocol == "" {
		protocol = "tcp"
	}

	address := fmt.Sprintf(":%d", port)
	conn, err := net.DialTimeout(protocol, address, 2*time.Second)

	status := &PortStatus{
		Port:      port,
		Protocol:  protocol,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if err == nil {
		conn.Close()
		status.State = "open"
	} else {
		status.State = "closed"
	}

	return status
}

// HealthCheck represents overall system health
type HealthCheck struct {
	CPU       *CPUInfo     `json:"cpu"`
	Memory    *MemoryInfo  `json:"memory"`
	Ports     []PortStatus `json:"ports,omitempty"`
	Status    string       `json:"status"`
	Timestamp string       `json:"timestamp"`
}

// PerformHealthCheck performs a full system health check
func (sm *SystemMonitor) PerformHealthCheck(portsToCheck []int) (*HealthCheck, error) {
	cpuInfo, cpuErr := sm.GetCPUUsage()
	memInfo, memErr := sm.GetMemoryUsage()

	status := "healthy"
	if memInfo != nil && memInfo.UsedPercent > 90 {
		status = "warning"
	}
	if memInfo != nil && memInfo.UsedPercent > 95 {
		status = "critical"
	}

	hc := &HealthCheck{
		CPU:       cpuInfo,
		Memory:    memInfo,
		Status:    status,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	if cpuErr != nil || memErr != nil {
		hc.Status = "error"
	}

	// Check ports
	for _, port := range portsToCheck {
		portStatus := sm.CheckPort(port, "tcp")
		hc.Ports = append(hc.Ports, *portStatus)
	}

	return hc, nil
}
