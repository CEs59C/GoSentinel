package report

import (
	"fmt"
	"goSentinel/internal/collector"
	"goSentinel/internal/email"
	"log"
	"strings"
	"time"
)

type SystemReport struct {
	CPU    collector.CPUInfo
	Disk   collector.DiskInfo
	Host   collector.HostInfo
	Memory collector.MemoryInfo
	Users  []collector.UserInfo
	Net    []collector.NetInfo
	Errors map[string]error
}

func Report() SystemReport {
	sr := SystemReport{}
	sr.CPU, _ = collector.GetCPUInfo() // wait 1 sec
	sr.Disk, _ = collector.GetDiskInfo()
	sr.Host, _ = collector.GetHostInfo()
	sr.Memory, _ = collector.GetMemoryInfo()
	sr.Users, _ = collector.GetUserInfo() // []
	sr.Net, _ = collector.GetNetInfo()    // []
	err := email.SendYandexEmail(sr.String())
	if err != nil {
		log.Println(err)
	}
	return sr
}

func (r SystemReport) String() string {
	var sb strings.Builder

	//sb.WriteString("=== System Report ===\n")
	sb.WriteString(fmt.Sprintf("Time%s\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("%s\n", r.CPU))
	sb.WriteString(fmt.Sprintf("%s\n", r.Disk))
	sb.WriteString(fmt.Sprintf("%s\n", r.Host))
	sb.WriteString(fmt.Sprintf("%s\n", r.Memory))
	for _, u := range r.Users {
		sb.WriteString(fmt.Sprintf("%s\n", u))
	}
	for _, n := range r.Net {
		sb.WriteString(fmt.Sprintf("%s\n", n))
	}

	if len(r.Errors) > 0 {
		sb.WriteString("\nErrors:\n")
		for key, err := range r.Errors {
			sb.WriteString(fmt.Sprintf("  %s: %v\n", key, err))
		}
	}

	return sb.String()
}
