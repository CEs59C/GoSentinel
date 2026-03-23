package collector

import (
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
)

type HostInfo struct {
	Hostname             string
	Uptime               uint64
	Procs                uint64
	OS                   string
	Platform             string
	PlatformFamily       string
	PlatformVersion      string
	VirtualizationSystem string
	VirtualizationRole   string
	ProcsRunning         int
	ProcsBlocked         int
	ProcsCreated         int
}

func GetHostInfo() (HostInfo, error) {
	info, err := host.Info()
	if err != nil {
		return HostInfo{}, fmt.Errorf("failed to get host.Info: %w", err)
	}

	misc, err := load.Misc()
	if err != nil {
		return HostInfo{}, fmt.Errorf("failed to get load.Misc: %w", err)
	}

	g := HostInfo{
		Hostname:             info.Hostname,
		Uptime:               info.Uptime, // Просто берем секунды. Переделать наверное сразу на prettyUptime
		Procs:                info.Procs,
		OS:                   info.OS,
		Platform:             info.Platform,
		PlatformFamily:       info.PlatformFamily,
		PlatformVersion:      info.PlatformVersion,
		VirtualizationSystem: info.VirtualizationSystem,
		VirtualizationRole:   info.VirtualizationRole,
		ProcsRunning:         misc.ProcsRunning,
		ProcsBlocked:         misc.ProcsBlocked,
		ProcsCreated:         misc.ProcsCreated,
	}
	fmt.Println(g)
	return g, nil
}

func (h HostInfo) String() string {
	prettyUptime := (time.Duration(h.Uptime) * time.Second).String()

	return fmt.Sprintf(
		"Host:\t\t%s [%s %s], Uptime=%s, Processes=%d, Running=%d, Blocked=%d, Created=%d, VM=%s (%s)",
		h.Hostname,
		h.Platform,
		h.PlatformVersion,
		prettyUptime,
		h.Procs,
		h.ProcsRunning,
		h.ProcsBlocked,
		h.ProcsCreated,
		h.VirtualizationSystem,
		h.VirtualizationRole,
	)
}
