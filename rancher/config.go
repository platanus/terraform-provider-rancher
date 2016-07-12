package rancher

import (
	"log"

	rancher "github.com/rancher/go-rancher/client"
)

type Config struct {
	APIURL    string
	AccessKey string
	SecretKey string
}

// Client returns a new Client for accessing Rancher.
func (c *Config) Client() (*rancher.RancherClient, error) {
	if c.APIURL == "" || c.AccessKey == "" || c.SecretKey == "" {
		return nil, nil
	}

	client, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       c.APIURL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Client configured for url: %s", c.APIURL)

	return client, nil
}
