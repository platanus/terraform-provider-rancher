package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	rancher "github.com/rancher/go-rancher/client"
)

func TestAccRancherEnvironment(t *testing.T) {
	var environment rancher.Project

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherEnvironmentDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherEnvironmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					testAccCheckRancherEnvironmentAttributes(&environment, "foo", "Terraform acc test group", "cattle"),
				),
			},
			resource.TestStep{
				Config: testAccRancherEnvironmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherEnvironmentExists("rancher_environment.foo", &environment),
					testAccCheckRancherEnvironmentAttributes(&environment, "foo2", "Terraform acc test group - updated", "swarm"),
				),
			},
		},
	})
}

func testAccCheckRancherEnvironmentExists(n string, env *rancher.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Name is set")
		}

		rancherServer := testAccProvider.Meta().(*RancherServer)
		client := rancherServer.GetClient()

		foundEnv, err := client.Project.ById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundEnv.Resource.Id != rs.Primary.ID {
			return fmt.Errorf("Environment not found")
		}

		*env = *foundEnv

		return nil
	}
}

func testAccCheckRancherEnvironmentAttributes(env *rancher.Project, envName string, envDesc string, envOrchestration string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if env.Name != envName {
			return fmt.Errorf("Bad name: %s shoud be: %s", env.Name, envName)
		}

		if env.Description != envDesc {
			return fmt.Errorf("Bad description: %s shoud be: %s", env.Description, envDesc)
		}

		orchestration := GetActiveOrchestration(env)
		if orchestration != envOrchestration {
			return fmt.Errorf("Bad orchestraion: %s shoud be: %s", orchestration, envOrchestration)
		}

		return nil
	}
}

func testAccCheckRancherEnvironmentDestroy(s *terraform.State) error {
	rancherServer := testAccProvider.Meta().(*RancherServer)
	client := rancherServer.GetClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_environment" {
			continue
		}
		env, err := client.Project.ById(rs.Primary.ID)

		if err == nil {
			if env != nil &&
				env.Resource.Id == rs.Primary.ID &&
				env.State != "removed" {
				return fmt.Errorf("Environment still exists")
			}
		}

		return nil
	}
	return nil
}

const testAccRancherEnvironmentConfig = `
resource "rancher_environment" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	orchestration = "cattle"
}
`

const testAccRancherEnvironmentUpdateConfig = `
resource "rancher_environment" "foo" {
	name = "foo2"
	description = "Terraform acc test group - updated"
	orchestration = "swarm"
}
`
