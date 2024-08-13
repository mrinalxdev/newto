package main

import (
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	psnet "github.com/shirou/gopsutil/net"
)

type NetworkData struct {
	Bandwidth     uint64    `json:"bandwidth"`
	CPUUsage      float64   `json:"cpuUsage"`
	MemoryUsage   float64   `json:"memoryUsage"`
	DiskUsage     float64   `json:"diskUsage"`
	PingLatency   float64   `json:"pingLatency"`
	Timestamp     time.Time `json:"timestamp"`
	HostInfo      HostInfo  `json:"hostInfo"`
	NetworkInfo   NetInfo   `json:"networkInfo"`
}

type HostInfo struct {
	Hostname string `json:"hostname"`
	Platform string `json:"platform"`
	OS       string `json:"os"`
	Uptime   uint64 `json:"uptime"`
}

type NetInfo struct {
	IPAddress  string `json:"ipAddress"`
	MacAddress string `json:"macAddress"`
}

func getNetworkData() NetworkData {
	return NetworkData{
		Bandwidth:     getBandwidthUsage(),
		CPUUsage:      getCPUUsage(),
		MemoryUsage:   getMemoryUsage(),
		DiskUsage:     getDiskUsage(),
		PingLatency:   getPingLatency("google.com"),
		Timestamp:     time.Now(),
		HostInfo:      getHostInfo(),
		NetworkInfo:   getNetworkInfo(),
	}
}

func getBandwidthUsage() uint64 {
	stats, _ := psnet.IOCounters(false)
	return stats[0].BytesSent + stats[0].BytesRecv
}

func getCPUUsage() float64 {
	percentage, _ := cpu.Percent(time.Second, false)
	if len(percentage) > 0 {
		return percentage[0]
	}
	return 0
}

func getMemoryUsage() float64 {
	vm, _ := mem.VirtualMemory()
	return vm.UsedPercent
}

func getDiskUsage() float64 {
	usage, _ := disk.Usage("/")
	return usage.UsedPercent
}

func getPingLatency(host string) float64 {
	start := time.Now()
	_, err := http.Get("http://" + host)
	if err != nil {
		return -1
	}
	return float64(time.Since(start).Milliseconds())
}

func getHostInfo() HostInfo {
	hostInfo, _ := host.Info()
	return HostInfo{
		Hostname: hostInfo.Hostname,
		Platform: hostInfo.Platform,
		OS:       fmt.Sprintf("%s %s", hostInfo.OS, runtime.GOARCH),
		Uptime:   hostInfo.Uptime,
	}
}

func getNetworkInfo() NetInfo {
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						return NetInfo{
							IPAddress:  ipnet.IP.String(),
							MacAddress: iface.HardwareAddr.String(),
						}
					}
				}
			}
		}
	}
	return NetInfo{}
}