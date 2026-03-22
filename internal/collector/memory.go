package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"
)

func memoryInfo() {
	fmt.Println("=======memoryInfo=======")
	v, _ := mem.VirtualMemory()

	fmt.Printf("Всего: %d MB\n", v.Total/1024/1024)
	fmt.Printf("Доступно: %d MB\n", v.Available/1024/1024)
	fmt.Printf("Использовано: %d MB (%.2f%%)\n", v.Used/1024/1024, v.UsedPercent)
	fmt.Printf("Свободно: %d MB\n", v.Free/1024/1024)

	swapInfo()
}

func swapInfo() {
	fmt.Println("=======swapInfo=======")
	swap, _ := mem.SwapMemory()

	fmt.Printf("Swap Total: %d MB\n", swap.Total/1024/1024)
	fmt.Printf("Swap Used: %d MB (%.2f%%)\n", swap.Used/1024/1024, swap.UsedPercent)
	fmt.Printf("Swap Free: %d MB\n", swap.Free/1024/1024)
}
