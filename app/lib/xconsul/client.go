package xconsul

import (
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Client struct {
	*consul.Client

	serviceID string
}

func NewClient(consulAddr string) (*Client, error) {
	config := consul.DefaultConfig()
	config.Address = consulAddr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &Client{Client: client}, nil
}

func (c *Client) RegisterService(addr, serviceName string) error {
	addrs := strings.Split(addr, ":")
	if len(addrs) != 2 {
		return errors.Errorf("incorrect address format: %s", addr)
	}

	host := addrs[0]
	port, err := strconv.Atoi(addrs[1])
	if err != nil {
		return err
	}

	c.serviceID = fmt.Sprintf("%s__%s:%d", serviceName, host, port)
	service := consul.AgentServiceRegistration{
		ID:      c.serviceID,
		Name:    serviceName,
		Address: host,
		Port:    port,
	}
	return c.Agent().ServiceRegister(&service)
}

func (c *Client) DeregisterService() error {
	return c.Agent().ServiceDeregister(c.serviceID)
}

func (c *Client) GetServices(service string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	addrs, meta, err := c.Health().Service(service, "", true, nil)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to get a service")
	}
	return addrs, meta, nil
}
