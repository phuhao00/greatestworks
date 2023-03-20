package consul

import (
	"github.com/hashicorp/consul/api"
	"path"
	"strings"
)

type Client struct {
	client *api.KV
}

func New(config *Config) (*Client, error) {
	conf := api.DefaultConfig()

	conf.Scheme = config.Scheme

	if len(config.Nodes) > 0 {
		conf.Address = config.Nodes[0]
	}

	if config.BasicAuth {
		conf.HttpAuth = &api.HttpBasicAuth{
			Username: config.Username,
			Password: config.Password,
		}
	}

	if config.ClientCert != "" && config.ClientKey != "" {
		conf.TLSConfig.CertFile = config.ClientCert
		conf.TLSConfig.KeyFile = config.ClientKey
	}
	if config.ClientCaKeys != "" {
		conf.TLSConfig.CAFile = config.ClientCaKeys
	}

	client, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}
	return &Client{client.KV()}, nil
}

// GetValues queries Consul for keys
func (c *Client) GetValues(keys []string) (map[string]string, error) {
	vars := make(map[string]string)
	for _, key := range keys {
		key := strings.TrimPrefix(key, "/")
		pairs, _, err := c.client.List(key, nil)
		if err != nil {
			return vars, err
		}
		for _, p := range pairs {
			vars[path.Join("/", p.Key)] = string(p.Value)
		}
	}
	return vars, nil
}

type watchResponse struct {
	waitIndex uint64
	err       error
}

func (c *Client) WatchPrefix(prefix string, keys []string, waitIndex uint64, stopChan chan bool) (uint64, error) {
	respChan := make(chan watchResponse)
	go func() {
		opts := api.QueryOptions{
			WaitIndex: waitIndex,
		}
		_, meta, err := c.client.List(prefix, &opts)
		if err != nil {
			respChan <- watchResponse{waitIndex, err}
			return
		}
		respChan <- watchResponse{meta.LastIndex, err}
	}()

	select {
	case <-stopChan:
		return waitIndex, nil
	case r := <-respChan:
		return r.waitIndex, r.err
	}
}
