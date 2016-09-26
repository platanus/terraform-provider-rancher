package rancher

import (
	"log"

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

	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       c.APIURL + "/projects/" + env + "/schemas",
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Client configured for url: %s", c.APIURL)

	return client, nil
}
