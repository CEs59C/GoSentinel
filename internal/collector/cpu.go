package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"
)

func cpuInfo() {
	fmt.Println("=======cpuInfo=======")
	infos, _ := cpu.Info()
	for _, info := range infos {
		fmt.Printf("CPU %d:\n", info.CPU)
		fmt.Printf("  Vendor: %s\n", info.VendorID)
		fmt.Printf("  Model: %s\n", info.ModelName)
		fmt.Printf("  Частота: %.0f MHz\n", info.Mhz)
		fmt.Printf("  Кэш: %d\n", info.CacheSize)
		fmt.Printf("  Ядер: %d\n", info.Cores)
	}
}
