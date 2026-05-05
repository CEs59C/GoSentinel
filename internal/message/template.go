package message

import (
	"bytes"
	"goSentinel/internal/collector"
	"goSentinel/internal/report"
	"html/template"
	"time"
)

type DataViewTemplate struct {
	Time   string
	CPU    collector.CPUInfo
	Disk   collector.DiskInfo
	Host   collector.HostInfo
	Memory collector.MemoryInfo
	Users  []collector.UserInfo
	Uptime string
	Net    []collector.NetInfo
}

func RenderHTMLReport(r report.SystemReport) (string, error) {
	tmpl, err := template.ParseFiles("templates/report.html")
	if err != nil {
		return "", err
	}

	data := DataViewTemplate{
		Time:   time.Now().Format("2006-01-02 15:04:05"),
		CPU:    r.CPU,
		Disk:   r.Disk,
		Host:   r.Host,
		Memory: r.Memory,
		Uptime: (time.Duration(r.Host.Uptime) * time.Second).String(),
		Net:    r.Net,
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
