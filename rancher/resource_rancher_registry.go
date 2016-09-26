package rancher

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	rancher "github.com/rancher/go-rancher/client"
)

func resourceRancherRegistry() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherRegistryCreate,
		Read:   resourceRancherRegistryRead,
		Update: resourceRancherRegistryUpdate,
		Delete: resourceRancherRegistryDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceRancherRegistryCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Registry: %s", d.Id())
	rancherServer := meta.(*RancherServer)
	client := rancherServer.GetEnvironmentClient(d.Get("environment_id").(string))

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	serverAddress := d.Get("server_address").(string)

	registry := rancher.Registry{
		Name:          name,
		Description:   description,
		ServerAddress: serverAddress,
	}
	newRegistry, err := client.Registry.Create(&registry)
	if err != nil {
		return err
	}

	d.SetId(newRegistry.Id)
	log.Printf("[INFO] Registry ID: %s", d.Id())

	return resourceRancherRegistryRead(d, meta)
}

func resourceRancherRegistryRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Registry: %s", d.Id())
	rancherServer := meta.(*RancherServer)
	client := rancherServer.GetEnvironmentClient(d.Get("environment_id").(string))

	env, err := client.Registry.ById(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Registry Name: %s", env.Name)

	d.Set("description", env.Description)
	d.Set("name", env.Name)
	d.Set("server_address", env.ServerAddress)

	return nil
}

func resourceRancherRegistryUpdate(d *schema.ResourceData, meta interface{}) error {
	rancherServer := meta.(*RancherServer)
	client := rancherServer.GetEnvironmentClient(d.Get("environment_id").(string))

	registry, err := client.Registry.ById(d.Id())
	if err != nil {
		return err
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)

	registry.Name = name
	registry.Description = description
	client.Registry.Update(registry, &registry)

	return resourceRancherRegistryRead(d, meta)
}

func resourceRancherRegistryDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Registry: %s", d.Id())
	id := d.Id()
	rancherServer := meta.(*RancherServer)
	client := rancherServer.GetEnvironmentClient(d.Get("environment_id").(string))

	reg, err := client.Registry.ById(id)
	if err != nil {
		return err
	}

	client.Registry.ActionDeactivate(reg)

	WaitFor(client, &reg.Resource, reg, func() string {
		return reg.Transitioning
	})

	if err := client.Registry.Delete(reg); err != nil {
		return fmt.Errorf("Error deleting Registry: %s", err)
	}

	log.Printf("[DEBUG] Waiting for registry (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"removed"},
		Refresh:    RegistryStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for registry (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// RegistryStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher Environment.
func RegistryStateRefreshFunc(client *rancher.RancherClient, registryID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		env, err := client.Registry.ById(registryID)

		if err != nil {
			return nil, "", err
		}

		return env, env.State, nil
	}
}
