package report

import "goSentinel/internal/collector"

func Report() {
	collector.GetCPUInfo()
	collector.GetDiskInfo()
	//collector.MemoryInfo()
	//collector.NetInfo()
	//collector.HostInfo()
	//collector.UserInfo()
}
