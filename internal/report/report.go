package report

import (
	"fmt"
	"goSentinel/internal/collector"
	"goSentinel/internal/email"
	"log"
	"strings"
	"sync"
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

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()

		if cpu, err := collector.GetCPUInfo(); err != nil { // wait 1 sec
			sr.Errors["cpu"] = err
		} else {
			sr.CPU = cpu
		}
	}()

	if disk, err := collector.GetDiskInfo(); err != nil {
		sr.Errors["disk"] = err
	} else {
		sr.Disk = disk
	}

	if host, err := collector.GetHostInfo(); err != nil {
		sr.Errors["host"] = err
	} else {
		sr.Host = host
	}

	if memory, err := collector.GetMemoryInfo(); err != nil {
		sr.Errors["memory"] = err
	} else {
		sr.Memory = memory
	}

	if user, err := collector.GetUserInfo(); err != nil {
		sr.Errors["user"] = err
	} else {
		sr.Users = user
	}

	if net, err := collector.GetNetInfo(); err != nil {
		sr.Errors["net"] = err
	} else {
		sr.Net = net
	}

	err := email.SendYandexEmail(sr.String())
	if err != nil {
		log.Println(err)
	}
	return sr
}

func (r SystemReport) String() string {
	var sb strings.Builder

	//sb.WriteString("=== System Report ===\n")
	fmt.Fprintf(&sb, "Time%s\n", time.Now().Format("2006-01-02 15:04:05"))
	fmt.Fprintf(&sb, "%s\n", r.CPU)
	fmt.Fprintf(&sb, "%s\n", r.Disk)
	fmt.Fprintf(&sb, "%s\n", r.Host)
	fmt.Fprintf(&sb, "%s\n", r.Memory)
	for _, u := range r.Users {
		fmt.Fprintf(&sb, "%s\n", u)
	}
	for _, n := range r.Net {
		fmt.Fprintf(&sb, "%s\n", n)
	}

	if len(r.Errors) > 0 {
		sb.WriteString("\nErrors:\n")
		for key, err := range r.Errors {
			fmt.Fprintf(&sb, "  %s: %v\n", key, err)
		}
	}

	return sb.String()
}
