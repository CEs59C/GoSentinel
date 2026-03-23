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
	v, err := mem.VirtualMemory()
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("failed to get virtual memory: %w", err)
	}

	s, err := mem.SwapMemory()
	if err != nil {
		return MemoryInfo{}, fmt.Errorf("failed to get swap memory: %w", err)
	}
	m := MemoryInfo{
		Total:           v.Total,
		Available:       v.Available,
		Used:            v.Used,
		Free:            v.Free,
		UsedPercent:     v.UsedPercent,
		SwapTotal:       s.Total,
		SwapUsed:        s.Used,
		SwapFree:        s.Free,
		SwapUsedPercent: s.UsedPercent,
	}
	fmt.Println(m)
	return m, nil
}

func (m MemoryInfo) String() string {
	const mb = 1024 * 1024

	return fmt.Sprintf("Memory:\t\tTotal=%dMB, Available=%dMB, Used=%dMB (%.2f%%), Free=%dMB\n"+
		"Swap:\t\tTotal=%dMB, Used=%dMB (%.2f%%), Free=%dMB",
		m.Total/mb, m.Available/mb, m.Used/mb, m.UsedPercent, m.Free/mb,
		m.SwapTotal/mb, m.SwapUsed/mb, m.SwapUsedPercent, m.SwapFree/mb,
	)
}
