package xconsul

import (
	"fmt"
	"net"
	"strconv"

	consul "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

type Client struct {
	*consul.Client

	currentServiceID string
}

func NewClient(consulAddr string) (*Client, error) {
	config := consul.DefaultConfig()
	config.Address = consulAddr
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to connect to consul")
	}
	return &Client{Client: client}, nil
}

func (c *Client) RegisterCurrentService(addr, name string, healthChecks consul.AgentServiceChecks) error {
	host, sport, err := net.SplitHostPort(addr)
	if err != nil {
		return errors.WithMessage(err, "failed to split host and port")
	}

	port, err := strconv.Atoi(sport)
	if err != nil {
		return errors.WithMessage(err, "failed to parse the service port")
	}

	c.currentServiceID = fmt.Sprintf("%s__%s:%d", name, host, port)
	service := consul.AgentServiceRegistration{
		ID:      c.currentServiceID,
		Name:    name,
		Address: host,
		Port:    port,
		Checks:  healthChecks,
	}
	if err = c.Agent().ServiceRegister(&service); err != nil {
		return errors.WithMessage(err, "failed to register service "+name)
	}
	return nil
}

func (c *Client) DeregisterCurrentService() error {
	return c.Agent().ServiceDeregister(c.currentServiceID)
}

func (c *Client) Service(service string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
	addrs, meta, err := c.Health().Service(service, "", true, nil)
	if err != nil {
		return nil, nil, errors.WithMessage(err, "failed to get services")
	}
	return addrs, meta, nil
}
