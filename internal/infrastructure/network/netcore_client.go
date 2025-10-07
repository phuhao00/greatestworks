package network

import (
	"fmt"
	"net"
	"time"

	"greatestworks/internal/infrastructure/logging"
)

// NetCoreClient 网络核心客户端
type NetCoreClient struct {
	conn   net.Conn
	logger logging.Logger
}

// NewNetCoreClient 创建网络核心客户端
func NewNetCoreClient(conn net.Conn, logger logging.Logger) *NetCoreClient {
	return &NetCoreClient{
		conn:   conn,
		logger: logger,
	}
}

// Connect 连接到服务器
func (c *NetCoreClient) Connect(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("连接失败: %w", err)
	}

	c.conn = conn
	c.logger.Info("已连接到服务器", map[string]interface{}{
		"address": address,
	})

	return nil
}

// Send 发送数据
func (c *NetCoreClient) Send(data []byte) error {
	if c.conn == nil {
		return fmt.Errorf("连接未建立")
	}

	_, err := c.conn.Write(data)
	if err != nil {
		return fmt.Errorf("发送数据失败: %w", err)
	}

	c.logger.Info("数据已发送", map[string]interface{}{
		"data_length": len(data),
	})

	return nil
}

// Receive 接收数据
func (c *NetCoreClient) Receive() ([]byte, error) {
	if c.conn == nil {
		return nil, fmt.Errorf("连接未建立")
	}

	buffer := make([]byte, 4096)
	n, err := c.conn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("接收数据失败: %w", err)
	}

	c.logger.Info("数据已接收", map[string]interface{}{
		"data_length": n,
	})

	return buffer[:n], nil
}

// Close 关闭连接
func (c *NetCoreClient) Close() error {
	if c.conn == nil {
		return nil
	}

	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("关闭连接失败: %w", err)
	}

	c.logger.Info("连接已关闭")
	return nil
}

// SetReadTimeout 设置读取超时
func (c *NetCoreClient) SetReadTimeout(timeout time.Duration) error {
	if c.conn == nil {
		return fmt.Errorf("连接未建立")
	}

	return c.conn.SetReadDeadline(time.Now().Add(timeout))
}

// SetWriteTimeout 设置写入超时
func (c *NetCoreClient) SetWriteTimeout(timeout time.Duration) error {
	if c.conn == nil {
		return fmt.Errorf("连接未建立")
	}

	return c.conn.SetWriteDeadline(time.Now().Add(timeout))
}

// IsConnected 检查是否已连接
func (c *NetCoreClient) IsConnected() bool {
	return c.conn != nil
}

// GetRemoteAddr 获取远程地址
func (c *NetCoreClient) GetRemoteAddr() string {
	if c.conn == nil {
		return ""
	}
	return c.conn.RemoteAddr().String()
}

// GetLocalAddr 获取本地地址
func (c *NetCoreClient) GetLocalAddr() string {
	if c.conn == nil {
		return ""
	}
	return c.conn.LocalAddr().String()
}
