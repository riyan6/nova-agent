package main

import (
	"context"
	"log"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"

	"google.golang.org/grpc"
	"nova-agent/pb"
)

const (
	serverAddr = "localhost:50051" // 修改为你的主控端地址
	agentID    = "agent-123"
)

var (
	lastNetIOTime   time.Time
	lastNetIO       uint64
	lastNetDownload uint64
)

func getStatus() (*pb.StatusRequest, error) {

	// CPU 使用率
	cpuPercent, _ := cpu.Percent(0, false)

	// CPU 型号
	cpuInfo, _ := cpu.Info()
	cpuModel := ""
	if len(cpuInfo) > 0 {
		cpuModel = cpuInfo[0].ModelName
	}

	// 内存使用率
	memInfo, _ := mem.VirtualMemory()

	// 磁盘使用率（根目录）
	diskInfo, _ := disk.Usage("/")

	// 网络速率（实时）
	netIO, _ := net.IOCounters(false)
	now := time.Now()
	var uploadKbps, downloadKbps float64

	if len(netIO) > 0 {
		if !lastNetIOTime.IsZero() {
			duration := now.Sub(lastNetIOTime).Seconds()
			uploadKbps = float64(netIO[0].BytesSent-lastNetIO) / 1024.0 / duration
			downloadKbps = float64(netIO[0].BytesRecv-lastNetDownload) / 1024.0 / duration
		}
		lastNetIO = netIO[0].BytesSent
		lastNetDownload = netIO[0].BytesRecv
		lastNetIOTime = now
	}

	return &pb.StatusRequest{
		AgentId:           agentID,
		CpuPercent:        cpuPercent[0],
		MemoryPercent:     memInfo.UsedPercent,
		UploadSpeedKbps:   uploadKbps,
		DownloadSpeedKbps: downloadKbps,
		Timestamp:         now.Unix(),
		CpuModel:          cpuModel,
		DiskPercent:       float32(diskInfo.UsedPercent),
	}, nil
}

func main() {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("连接失败: %v", err)
	}
	defer conn.Close()

	client := pb.NewVpsClient(conn)
	stream, err := client.ReportStatus(context.Background())
	if err != nil {
		log.Fatalf("无法开始状态上报: %v", err)
	}

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		status, _ := getStatus()
		if err := stream.Send(status); err != nil {
			log.Printf("发送失败: %v", err)
		}
	}
}
