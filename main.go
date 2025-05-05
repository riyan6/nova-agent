package main

import (
	"log"

	"nova-agent/client"
	"nova-agent/config"
)

func main() {
	log.Println("Agent 程序启动...")

	// 加载配置
	cfg := config.NewConfig()

	// 初始化 gRPC 客户端
	agentClient, err := client.NewClient(cfg)
	if err != nil {
		log.Fatalf("初始化客户端失败: %v", err)
	}
	defer agentClient.Close()

	log.Println("Agent 程序启动成功")

	// 启动状态上报
	if err := agentClient.Run(); err != nil {
		log.Fatalf("运行客户端失败: %v", err)
	}
}
