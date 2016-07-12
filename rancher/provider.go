package rancher

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	// "github.com/rancher/go-rancher/client"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CATTLE_URL", nil),
				Description: descriptions["api_url"],
			},
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CATTLE_ACCESS_KEY", nil),
				Description: descriptions["access_key"],
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("CATTLE_SECRET_KEY", nil),
				Description: descriptions["secret_key"],
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"rancher_environment": resourceRancherEnvironment(),
		},

		ConfigureFunc: providerConfigure,
	}
}

var descriptions map[string]string

func init() {
	descriptions = map[string]string{
		"access_key": "API Key used to authenticate with the rancher server",

		"secret_key": "API secret used to authenticate with the rancher server",

		"api_url": "The URL to the rancher API",
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		APIURL:    d.Get("api_url").(string),
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
	}

	return config.Client()
}
