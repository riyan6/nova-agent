package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	ServerAddr            string `yaml:"server_addr"`
	AgentID               int32  `yaml:"agent_id"`
	ReportIntervalSeconds int    `yaml:"report_interval_seconds"`
}

func NewConfig() *Config {
	cfg := &Config{
		ServerAddr:            "localhost:50051",
		AgentID:               0,
		ReportIntervalSeconds: 5,
	}

	// 尝试从 config.yml 加载配置
	data, err := os.ReadFile("config.yml")
	if err != nil {
		log.Printf("无法读取 config.yml，使用默认配置: %v", err)
		return cfg
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		log.Printf("解析 config.yml 失败，使用默认配置: %v", err)
		return cfg
	}

	log.Printf("成功加载 config.yml: ServerAddr=%s, AgentID=%s, AgentIDInt=%d", cfg.ServerAddr, cfg.AgentID, cfg.AgentID)
	return cfg
}
