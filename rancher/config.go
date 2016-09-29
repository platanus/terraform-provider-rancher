package rancher

import (
	"log"

	"github.com/rancher/go-rancher/catalog"
	rancher "github.com/rancher/go-rancher/client"
)

type Config struct {
	*rancher.RancherClient
	APIURL    string
	AccessKey string
	SecretKey string
}

// Create creates a generic Rancher client
func (c *Config) CreateClient() error {
	if c.APIURL == "" || c.AccessKey == "" || c.SecretKey == "" {
		return nil
	}

	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       c.APIURL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return err
	}

	log.Printf("[INFO] Rancher Client configured for url: %s", c.APIURL)

	c.RancherClient = client

	return nil
}

func (c *Config) EnvironmentClient(env string) (*rancher.RancherClient, error) {
	if c.APIURL == "" || c.AccessKey == "" || c.SecretKey == "" {
		return nil, nil
	}

	url := c.APIURL + "/projects/" + env + "/schemas"
	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       url,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Client configured for url: %s", url)

	return client, nil
}

func (c *Config) RegistryClient(id string) (*rancher.RancherClient, error) {
	reg, err := c.Registry.ById(id)
	if err != nil {
		return nil, err
	}

	return c.EnvironmentClient(reg.AccountId)
}

func (c *Config) CatalogClient() (*catalog.RancherClient, error) {
	if c.APIURL == "" || c.AccessKey == "" || c.SecretKey == "" {
		return nil, nil
	}

	url := c.APIURL + "-catalog/schemas"
	client, err := catalog.NewRancherClient(&catalog.ClientOpts{
		Url:       url,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Catalog Client configured for url: %s", url)

	return client, nil
}
