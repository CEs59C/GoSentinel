package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

type NetInfo struct {
	Name string
	Port uint32
	Pid  int32
}

func GetNetInfo() ([]NetInfo, error) {
	connections, err := net.Connections("tcp")
	if err != nil {
		return nil, fmt.Errorf("failed to get TCP connections: %w", err)
	}

	var listening []NetInfo

	for _, conn := range connections {
		if conn.Status == "LISTEN" {
			info := NetInfo{
				Port: conn.Laddr.Port,
				Pid:  conn.Pid,
				Name: "Unknown",
			}

			if conn.Pid > 0 {
				p, err := process.NewProcess(conn.Pid)
				if err == nil {
					if name, err := p.Name(); err == nil {
						info.Name = name
					}
				}
			}

			listening = append(listening, info)
		}
	}

	return listening, nil
}

func (n NetInfo) String() string {
	return fmt.Sprintf("Process: %-20s Port: %-5d PID: %d", n.Name, n.Port, n.Pid)
}
