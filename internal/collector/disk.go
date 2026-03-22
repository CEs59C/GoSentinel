package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

type DiskInfo struct {
	Total       uint64
	Used        uint64
	Free        uint64
	UsedPercent float64 // процент использованного
	FreePercent float64 // процент свободного
	InodesUsed  float64 // процент использованных inodes
}

func GetDiskInfo() (DiskInfo, error) {
	usage, err := disk.Usage("/")
	if err != nil {
		return DiskInfo{}, fmt.Errorf("failed to get disck info: %w", err)
	}
	const bytesToGB = 1024 * 1024 * 1024

	d := DiskInfo{
		Total:       usage.Total / bytesToGB,
		Used:        usage.Used / bytesToGB,
		Free:        usage.Free / bytesToGB,
		UsedPercent: usage.UsedPercent,
		FreePercent: 100.0 - usage.UsedPercent,
		InodesUsed:  usage.InodesUsedPercent,
	}
	fmt.Println(d.String())
	return d, err
}

func (d DiskInfo) String() string {
	return fmt.Sprintf(
		"Disk: Total=%dGB, Used=%dGB, (%.1f%%) Free=%dGB, (%.1f%%), Inodes=%.1f%%.",
		d.Total, d.Used, d.UsedPercent, d.Free, d.FreePercent, d.InodesUsed,
	)
}
