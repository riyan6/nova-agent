package collector

import (
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/host"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
	"google.golang.org/protobuf/types/known/timestamppb"
	"nova-agent/config"
	"nova-agent/pb"
)

var (
	lastNetIOTime      time.Time
	lastNetInTransfer  uint64
	lastNetOutTransfer uint64
)

func GetStatus(cfg *config.Config) (*pb.StatusRequest, error) {
	now := time.Now()

	// Host 信息
	hostInfo, _ := host.Info()
	platform := hostInfo.Platform
	arch := hostInfo.KernelArch
	bootTime := int64(hostInfo.BootTime)

	// CPU 信息
	cpuPercent, _ := cpu.Percent(0, false)
	cpuInfo, _ := cpu.Info()
	cpuModels := make([]string, 0, len(cpuInfo))
	for _, ci := range cpuInfo {
		cpuModels = append(cpuModels, ci.ModelName)
	}

	// 内存信息
	memInfo, _ := mem.VirtualMemory()

	// 磁盘信息（根目录）
	diskInfo, _ := disk.Usage("/")

	// 网络信息
	netIO, _ := net.IOCounters(false)
	var netInSpeed, netOutSpeed uint64
	var netInTransfer, netOutTransfer uint64
	if len(netIO) > 0 {
		netInTransfer = netIO[0].BytesRecv
		netOutTransfer = netIO[0].BytesSent
		if !lastNetIOTime.IsZero() {
			duration := now.Sub(lastNetIOTime).Seconds()
			if duration > 0 {
				netInSpeed = uint64(float64(netInTransfer-lastNetInTransfer) / duration)
				netOutSpeed = uint64(float64(netOutTransfer-lastNetOutTransfer) / duration)
			}
		}
		lastNetInTransfer = netInTransfer
		lastNetOutTransfer = netOutTransfer
		lastNetIOTime = now
	}

	// 负载信息
	loadAvg, _ := load.Avg()

	// 网络连接信息
	conns, _ := net.Connections("all")
	tcpConnCount := 0
	udpConnCount := 0
	for _, conn := range conns {
		if conn.Type == 1 { // TCP
			tcpConnCount++
		} else if conn.Type == 2 { // UDP
			udpConnCount++
		}
	}

	// 进程数
	procs, _ := process.Processes()
	processCount := len(procs)

	// 运行时间
	uptime := uint64(hostInfo.Uptime)

	return &pb.StatusRequest{
		Id: cfg.AgentID, // 从配置中获取整数 ID
		Host: &pb.HostInfo{
			Platform:  platform,
			Cpu:       cpuModels,
			MemTotal:  memInfo.Total,
			DiskTotal: diskInfo.Total,
			Arch:      arch,
			BootTime:  bootTime,
		},
		State: &pb.StateInfo{
			Cpu:            cpuPercent[0],
			MemUsed:        memInfo.Used,
			DiskUsed:       diskInfo.Used,
			NetInTransfer:  netInTransfer,
			NetOutTransfer: netOutTransfer,
			NetInSpeed:     netInSpeed,
			NetOutSpeed:    netOutSpeed,
			Uptime:         uptime,
			Load_5:         loadAvg.Load5,
			TcpConnCount:   int32(tcpConnCount),
			UdpConnCount:   int32(udpConnCount),
			ProcessCount:   int32(processCount),
		},
		LastActive: timestamppb.New(now),
	}, nil
}
