package rancher

import rancher "github.com/rancher/go-rancher/client"

type RancherServerConfig struct {
	APIURL    string
	AccessKey string
	SecretKey string
}

//
type RancherServer struct {
	config *RancherServerConfig
}

// // Client returns a new Client for accessing Rancher.
// func (r *RancherServer) Server() (*RancherServer, error) {
// 	if r.config.APIURL == "" || r.config.AccessKey == "" || r.config.SecretKey == "" {
// 		return nil, nil
// 	}
//
// 	log.Printf("[INFO] Rancher Server configured for url: %s", r.config.APIURL)
//
// 	return client, nil
// }

// GetClient
func (r *RancherServer) GetClient() *rancher.RancherClient {
	rancherClient, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       r.config.APIURL,
		AccessKey: r.config.AccessKey,
		SecretKey: r.config.SecretKey,
	})
	if err != nil {
		return nil
	}

	return rancherClient
}

// GetEnvironmentClient
func (r *RancherServer) GetEnvironmentClient(environmentID string) *rancher.RancherClient {
	rancherClient, err := rancher.NewRancherClient(&rancher.ClientOpts{
		Url:       r.config.APIURL + "/projects/" + environmentID + "/schemas",
		AccessKey: r.config.AccessKey,
		SecretKey: r.config.SecretKey,
	})
	if err != nil {
		return nil
	}

	return rancherClient
}
