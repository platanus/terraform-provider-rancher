package rancher

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	rancher "github.com/rancher/go-rancher/client"
)

func resourceRancherStack() *schema.Resource {
	return &schema.Resource{
		Create: resourceRancherStackCreate,
		Read:   resourceRancherStackRead,
		Update: resourceRancherStackUpdate,
		Delete: resourceRancherStackDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"environment_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"docker_compose": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"rancher_compose": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"environment": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"catalog_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Default:  "user",
				Optional: true,
			},
			"start_on_create": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceRancherStackCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Creating Stack: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	data, err := makeStackData(d, meta)
	if err != nil {
		return err
	}

	var newStack rancher.Environment
	if err := client.Create("environment", data, &newStack); err != nil {
		return err
	}

	d.SetId(newStack.Id)
	log.Printf("[INFO] Stack ID: %s", d.Id())

	return resourceRancherStackRead(d, meta)
}

func resourceRancherStackRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Refreshing Stack: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	stack, err := client.Environment.ById(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[INFO] Stack Name: %s", stack.Name)

	d.Set("description", stack.Description)
	d.Set("name", stack.Name)
	d.Set("docker_compose", stack.DockerCompose)
	d.Set("rancher_compose", stack.RancherCompose)
	d.Set("environment", stack.Environment)

	if stack.ExternalId == "" {
		d.Set("scope", "user")
		d.Set("catalog_id", "")
	} else {
		trimmedID := strings.TrimPrefix(stack.ExternalId, "system-")
		if trimmedID == stack.ExternalId {
			d.Set("scope", "user")
		} else {
			d.Set("scope", "system")
		}
		d.Set("catalog_id", strings.TrimPrefix(trimmedID, "catalog://"))
	}

	d.Set("start_on_create", stack.StartOnCreate)

	return nil
}

func resourceRancherStackUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Updating Stack: %s", d.Id())
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	data, err := makeStackData(d, meta)
	if err != nil {
		return err
	}

	// TODO: upgrade stack if docker_compose, rancher_compose or environment have changed

	var newStack rancher.Environment
	stack, err := client.Environment.ById(d.Id())
	if err != nil {
		return err
	}

	if err := client.Update("environment", &stack.Resource, data, &newStack); err != nil {
		return err
	}

	return resourceRancherStackRead(d, meta)
}

func resourceRancherStackDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Deleting Stack: %s", d.Id())
	id := d.Id()
	client, err := meta.(*Config).EnvironmentClient(d.Get("environment_id").(string))
	if err != nil {
		return err
	}

	stack, err := client.Environment.ById(id)
	if err != nil {
		return err
	}

	if err := client.Environment.Delete(stack); err != nil {
		return fmt.Errorf("Error deleting Stack: %s", err)
	}

	log.Printf("[DEBUG] Waiting for stack (%s) to be removed", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"active", "removed", "removing"},
		Target:     []string{"removed"},
		Refresh:    StackStateRefreshFunc(client, id),
		Timeout:    10 * time.Minute,
		Delay:      1 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, waitErr := stateConf.WaitForState()
	if waitErr != nil {
		return fmt.Errorf(
			"Error waiting for stack (%s) to be removed: %s", id, waitErr)
	}

	d.SetId("")
	return nil
}

// StackStateRefreshFunc returns a resource.StateRefreshFunc that is used to watch
// a Rancher Stack.
func StackStateRefreshFunc(client *rancher.RancherClient, stackID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		stack, err := client.Environment.ById(stackID)

		if err != nil {
			return nil, "", err
		}

		return stack, stack.State, nil
	}
}

func environmentFromMap(m map[string]interface{}) map[string]string {
	result := make(map[string]string)
	for k, v := range m {
		result[k] = v.(string)
	}
	return result
}

func makeStackData(d *schema.ResourceData, meta interface{}) (data map[string]interface{}, err error) {
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	var externalID string
	var dockerCompose string
	var rancherCompose string
	var environment map[string]string
	if c, ok := d.GetOk("catalog_id"); ok {
		if scope, ok := d.GetOk("scope"); ok && scope.(string) == "system" {
			externalID = "system-"
		}
		catalogID := c.(string)
		externalID += "catalog://" + catalogID

		catalogClient, err := meta.(*Config).CatalogClient()
		if err != nil {
			return data, err
		}
		template, err := catalogClient.Template.ById(catalogID)
		if err != nil {
			return data, fmt.Errorf("Failed to get catalog template: %s", err)
		}

		dockerCompose = template.Files["docker-compose.yml"].(string)
		// Pretend user provided this
		d.Set("docker_compose", dockerCompose)
		rancherCompose = template.Files["rancher-compose.yml"].(string)
		// Pretend user provided this
		d.Set("rancher_compose", rancherCompose)
	}

	if c, ok := d.GetOk("docker_compose"); ok {
		dockerCompose = c.(string)
	}
	if c, ok := d.GetOk("rancher_compose"); ok {
		rancherCompose = c.(string)
	}
	environment = environmentFromMap(d.Get("environment").(map[string]interface{}))

	startOnCreate := d.Get("start_on_create")

	data = map[string]interface{}{
		"name":           &name,
		"description":    &description,
		"dockerCompose":  &dockerCompose,
		"rancherCompose": &rancherCompose,
		"environment":    &environment,
		"externalId":     &externalID,
		"startOnCreate":  &startOnCreate,
	}

	return data, nil
}
