package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUInfo struct {
	Model  string
	Vendor string
	Cores  int
	Usage  float64 // % загрузки прямо сейчас
}

func GetCPUInfo() (CPUInfo, error) {
	infos, err := cpu.Info()
	if err != nil {
		return CPUInfo{}, fmt.Errorf("failed to get CPU info: %w", err)
	}

	usage, err := cpu.Percent(time.Second, false)
	if err != nil {
		return CPUInfo{}, fmt.Errorf("failed to get CPU usage: %w", err)
	}

	c := CPUInfo{
		Model:  "Unknown",
		Vendor: "Unknown",
		Cores:  -1,
		Usage:  -1,
	}

	if len(usage) > 0 {
		c.Usage = usage[0]
	}
	if len(infos) > 0 {
		c.Model = infos[0].ModelName
		c.Vendor = infos[0].VendorID
		c.Cores = int(infos[0].Cores)
	}

	return c, err
}
func (c CPUInfo) String() string {
	return fmt.Sprintf(
		"CPU Info:\tModel=%s, Vendor=%s, Cores=%d, Usage=%.2f%%.",
		c.Model, c.Vendor, c.Cores, c.Usage,
	)
}
