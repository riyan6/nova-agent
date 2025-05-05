package client

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"nova-agent/collector"
	"nova-agent/config"
	"nova-agent/pb"
)

type Client struct {
	cfg       *config.Config
	conn      *grpc.ClientConn
	vpsClient pb.VpsClient
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewClient(cfg *config.Config) (*Client, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client := &Client{
		cfg:    cfg,
		ctx:    ctx,
		cancel: cancel,
	}

	if err := client.connect(); err != nil {
		cancel()
		return nil, err
	}

	return client, nil
}

func (c *Client) connect() error {
	log.Printf("正在连接主控端: %s", c.cfg.ServerAddr)
	conn, err := grpc.Dial(c.cfg.ServerAddr, grpc.WithInsecure())
	if err != nil {
		return err
	}
	c.conn = conn
	c.vpsClient = pb.NewVpsClient(conn)
	log.Println("gRPC 连接主控端成功")
	return nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
		log.Println("gRPC 连接已关闭")
	}
	c.cancel()
}

func (c *Client) Run() error {
	for {
		select {
		case <-c.ctx.Done():
			log.Println("客户端停止运行")
			return nil
		default:
			if err := c.reportStatus(); err != nil {
				log.Printf("状态上报失败: %v", err)
				// 尝试重连
				c.Close()
				if err := c.connect(); err != nil {
					log.Printf("重连失败: %v，将在 5 秒后重试", err)
					time.Sleep(5 * time.Second)
					continue
				}
			}
		}
	}
}

func (c *Client) reportStatus() error {
	stream, err := c.vpsClient.ReportStatus(c.ctx)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(time.Duration(c.cfg.ReportIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			return nil
		case <-ticker.C:
			status, err := collector.GetStatus(c.cfg) // 传入整个 config
			if err != nil {
				log.Printf("获取状态失败: %v", err)
				continue
			}
			if err := stream.Send(status); err != nil {
				log.Printf("发送状态失败: %v", err)
				return err
			}
			log.Printf("成功发送状态 [Agent %d]: CPU:%.2f%%, 内存:%.2f%%", status.Id, status.State.Cpu, float64(status.State.MemUsed)/float64(status.Host.MemTotal)*100)
		}
	}
}
