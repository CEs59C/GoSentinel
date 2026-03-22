package collector

import (
	"fmt"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

func netInfo() {
	fmt.Println("=======netInfo (Listening Ports)=======")
	connections, _ := net.Connections("tcp")

	for _, conn := range connections {
		// Показываем только порты, которые ждут входящих соединений
		if conn.Status == "LISTEN" {
			p, err := process.NewProcess(conn.Pid)
			name := "Unknown"
			if err == nil {
				name, _ = p.Name()
			}

			fmt.Printf("Процесс: %-12s | Порт: %-5d | PID: %d\n",
				name, conn.Laddr.Port, conn.Pid)
		}
	}
}
