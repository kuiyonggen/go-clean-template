// Package consul implements consul connection.
package consul

import (
    "fmt"
    "bytes"
    "errors"
    "strconv"
    "strings"
    
    api "github.com/hashicorp/consul/api"
    "github.com/spf13/viper"
    "gopkg.in/yaml.v2"
)

// Consul _.
type Consul struct {
    // Client is the client to use for sending commands to consul.
    Client *api.Client

    // Address is the address of consul
    Address string

    // Folder is the folder name of kv config in consul.
    Folder string

    // Service is the name to be registered in consul.
    Service string

}

// New _.
func New(address, folder, service string) (*Consul, error) {
    consul := &Consul{
        Client: nil,
        Address: address,
        Folder: folder,
        Service: service,
    }
    config := api.DefaultConfig()
    config.Address = address
    var err error
    consul.Client, err = api.NewClient(config)

    return consul, err
}

// Kv _.
func (c *Consul) Kv(cfg interface{}) error {
    key := strings.Join([]string{
        c.Folder,
        c.Service,
    }, "/")

    kv, _, err := c.Client.KV().Get(key, nil)
    if err != nil {
        return err
    }

    if kv == nil {
        return errors.New(fmt.Sprintf("KV not found for %s.", key))
    } else {
        // only support yaml kv
        kvIO := strings.NewReader(string(kv.Value))
        err = yaml.NewDecoder(kvIO).Decode(cfg)
        if err != nil {
            return err
        }
        viper.SetConfigType("yaml")
        err = viper.ReadConfig(bytes.NewBuffer(kv.Value))
        
        return err
    }
}

// Register -.
func (c *Consul) Register(listenAddress, listenPort, checkApi, interval, timeout string, 
        tags []string) (string, error) {
    check := api.AgentServiceCheck{
        HTTP: fmt.Sprintf("http://%s:%s%s", listenAddress, listenPort, checkApi),
        Interval: interval,
        Timeout: timeout,
        Notes: "Consul check service health status.",
    }

    intPort, _ := strconv.Atoi(listenPort)
    reg := &api.AgentServiceRegistration{
        ID: fmt.Sprintf("%s_%s:%s", c.Service, 
                listenAddress, listenPort),
        Name: c.Service,
        Tags: tags,
        Meta: map[string]string{
            "swagger": fmt.Sprintf("http://%s:%s/swagger/index.html", 
                listenAddress, listenPort)},
        Address: listenAddress,
        Port: intPort,
        Check: &check,
    }

    err := c.Client.Agent().ServiceRegister(reg)

    return reg.ID, err
}

// Deregister ._
func (c *Consul) Deregister(serviceID string) error {
    return c.Client.Agent().ServiceDeregister(serviceID)
}

