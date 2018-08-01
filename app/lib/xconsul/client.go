package xconsul

// import (
// 	"fmt"
// 	"strconv"
// 	"strings"

// 	consul "github.com/hashicorp/consul/api"
// 	"github.com/pkg/errors"
// )

// type Client struct {
// 	*consul.Client
// }

// func NewClient(consulAddr string) (*Client, error) {
// 	config := consul.DefaultConfig()
// 	config.Address = consulAddr
// 	client, err := consul.NewClient(config)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &Client{Client: client}, nil
// }

// func (c *Client) RegisterService(addr, serviceName string) (string, error) {
// 	addrs := strings.Split(addr, ":")
// 	if len(addrs) != 2 {
// 		return "", errors.Errorf("incorrect address format: %s", addr)
// 	}

// 	host := addrs[0]
// 	port, err := strconv.Atoi(addrs[1])
// 	if err != nil {
// 		return "", err
// 	}

// 	serviceID := fmt.Sprintf("%s__%s:%d", serviceName, host, port)
// 	service := consul.AgentServiceRegistration{
// 		ID:      serviceID,
// 		Name:    serviceName,
// 		Address: host,
// 		Port:    port,
// 	}
// 	return serviceID, c.Agent().ServiceRegister(&service)
// }

// func (c *Client) DeregisterService(serviceID string) error {
// 	return c.Agent().ServiceDeregister(serviceID)
// }

// func (c *Client) Service(service string) ([]*consul.ServiceEntry, *consul.QueryMeta, error) {
// 	addrs, meta, err := c.Health().Service(service, "", true, nil)
// 	if err != nil {
// 		return nil, nil, errors.Wrap(err, "failed to get services")
// 	}
// 	return addrs, meta, nil
// }
