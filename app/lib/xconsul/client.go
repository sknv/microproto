package xconsul

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
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
		return errors.New("incorrect address format")
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
