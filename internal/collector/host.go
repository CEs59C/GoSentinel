package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
)

func hostInfo() {
	fmt.Println("=======hostInfo=======")
	info, _ := host.Info()

	fmt.Printf("Имя хоста: %s\n", info.Hostname)
	fmt.Printf("Время работы: %s\n", time.Duration(info.Uptime)*time.Second)
	fmt.Printf("Количество процессов: %d\n", info.Procs)
	fmt.Printf("ОС: %s\n", info.OS)
	fmt.Printf("Платформа: %s\n", info.Platform)
	fmt.Printf("Семейство: %s\n", info.PlatformFamily)
	fmt.Printf("Версия: %s\n", info.PlatformVersion)
	fmt.Printf("Виртуализация: %s (%s)\n", info.VirtualizationSystem, info.VirtualizationRole)
	foo()
}

func foo() {
	avg, _ := load.Avg()

	fmt.Printf("Load Average:\n")
	fmt.Printf("  1 минута: %.2f\n", avg.Load1)
	fmt.Printf("  5 минут: %.2f\n", avg.Load5)
	fmt.Printf("  15 минут: %.2f\n", avg.Load15)

	// Дополнительные метрики (Linux)
	misc, _ := load.Misc()
	fmt.Printf("Запущенных процессов: %d\n", misc.ProcsRunning)
	fmt.Printf("Заблокированных: %d\n", misc.ProcsBlocked)
	fmt.Printf("Создано процессов: %d\n", misc.ProcsCreated)
}
