package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/disk"
)

func discInfo() {
	fmt.Println("=======discInfo=======")
	usage, _ := disk.Usage("/")

	fmt.Printf("Всего: %d GB\n", usage.Total/1024/1024/1024)
	fmt.Printf("Использовано: %d GB\n", usage.Used/1024/1024/1024)
	fmt.Printf("Свободно: %d GB\n", usage.Free/1024/1024/1024)
	fmt.Printf("Использовано: %.2f%%\n", usage.UsedPercent)

	// Дополнительно: inodes (Linux)
	if usage.InodesTotal > 0 {
		fmt.Printf("Inodes used: %.2f%%\n", usage.InodesUsedPercent)
	}
}
