package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

type MemoryInfo struct {
	Total       uint64
	Available   uint64
	Used        uint64
	Free        uint64
	UsedPercent float64

	SwapTotal       uint64
	SwapUsed        uint64
	SwapFree        uint64
	SwapUsedPercent float64
}

func GetMemoryInfo() (MemoryInfo, error) {
	const mb = 1024 * 1024
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	s, err := mem.SwapMemory()
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("failed to get swap memory: %w", err)
	}

	m := MemoryInfo{
		Total:           v.Total / mb,
		Available:       v.Available / mb,
		Used:            v.Used / mb,
		Free:            v.Free / mb,
		UsedPercent:     v.UsedPercent,
		SwapTotal:       s.Total / mb,
		SwapUsed:        s.Used / mb,
		SwapFree:        s.Free / mb,
		SwapUsedPercent: s.UsedPercent,
	}

	return m, nil
}

func (m MemoryInfo) String() string {
	return fmt.Sprintf("Memory:\t\tTotal=%dMB, Available=%dMB, Used=%dMB (%.2f%%), Free=%dMB\n"+
		"Swap:\t\tTotal=%dMB, Used=%dMB (%.2f%%), Free=%dMB",
		m.Total, m.Available, m.Used, m.UsedPercent, m.Free,
		m.SwapTotal, m.SwapUsed, m.SwapUsedPercent, m.SwapFree,
	)
}
