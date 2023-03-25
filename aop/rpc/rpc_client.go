package rpc

import (
	"net/rpc"
)

type RpcClient struct {
	pool *Pool
	Addr string
}

func NewRpcClient(addr string) *RpcClient {
	rpcClient := &RpcClient{
		pool: &Pool{
			MaxIdle:         1,
			IdleTimeout:     0,
			MaxConnLifetime: 0,
			Dial:            func() (*rpc.Client, error) { return rpc.Dial("tcp", addr) },
		},
		Addr: addr,
	}
	return rpcClient
}

// Call 远程调用
func (c *RpcClient) Call(method string, args interface{}, reply interface{}) error {
	rpcClient, err := c.pool.Get()
	if err != nil {
		return err
	}
	err = rpcClient.Call(method, args, reply)
	if err == rpc.ErrShutdown {
		return err
	}
	rpcClient.Close()
	return err
}
