package rancher

import (
	"time"

	"github.com/rancher/go-rancher/client"
)

// GetActiveOrchestration get the name of the active orchestration for a environment
func GetActiveOrchestration(project *client.Project) string {
	orch := "cattle"

	switch {
	case project.Swarm:
		orch = "swarm"
	case project.Mesos:
		orch = "mesos"
	case project.Kubernetes:
		orch = "kubernetes"
	}

	return orch
}

// WaitFor waits for a resource to reach a certain state.
func WaitFor(c *client.RancherClient, resource *client.Resource, output interface{}, transitioning func() string) error {
	for {
		transitioning := transitioning()
		if transitioning != "yes" && transitioning == "no" {
			return nil
		}

		time.Sleep(150 * time.Millisecond)

		err := c.Reload(resource, output)
		if err != nil {
			return err
		}
	}
}
