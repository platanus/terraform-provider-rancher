package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	rancher "github.com/rancher/go-rancher/client"
)

func TestAccRancherRegistry(t *testing.T) {
	var registry rancher.Registry

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherRegistryDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherRegistryConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherRegistryExists("rancher_registry.foo", &registry),
					testAccCheckRancherRegistryAttributes(&registry, "foo", "registry test", "http://foo.com:8080"),
				),
			},
			// resource.TestStep{
			// 	Config: testAccRancherRegistryUpdateConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckRancherRegistryExists("rancher_registry.foo", &registry),
			// 		testAccCheckRancherRegistryAttributes(&registry, "foo2", "Terraform acc test group - updated", "swarm"),
			// 	),
			// },
			// resource.TestStep{
			// 	Config: testAccRancherRegistryRecreateConfig,
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckRancherRegistryExists("rancher_registry.foo", &registry),
			// 		testAccCheckRancherRegistryAttributes(&registry, "foo2", "Terraform acc test group - updated", "swarm"),
			// 	),
			// },
		},
	})
}

func testAccCheckRancherRegistryExists(n string, reg *rancher.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Name is set")
		}

		rancherServer := testAccProvider.Meta().(*RancherServer)
		client := rancherServer.GetEnvironmentClient("1a26")

		foundReg, err := client.Registry.ById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundReg.Resource.Id != rs.Primary.ID {
			return fmt.Errorf("Environment not found")
		}

		*reg = *foundReg

		return nil
	}
}

func testAccCheckRancherRegistryAttributes(reg *rancher.Registry, regName string, regDesc string, regServerAddress string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if reg.Name != regName {
			return fmt.Errorf("Bad name: %s shoud be: %s", reg.Name, regName)
		}

		if reg.Description != regDesc {
			return fmt.Errorf("Bad description: %s shoud be: %s", reg.Description, regDesc)
		}

		if reg.ServerAddress != regServerAddress {
			return fmt.Errorf("Bad server_address: %s shoud be: %s", reg.ServerAddress, regServerAddress)
		}

		return nil
	}
}

func testAccCheckRancherRegistryDestroy(s *terraform.State) error {
	rancherServer := testAccProvider.Meta().(*RancherServer)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_registry" {
			continue
		}
		client := rancherServer.GetEnvironmentClient("1a26")
		reg, err := client.Registry.ById(rs.Primary.ID)

		if err == nil {
			if reg != nil &&
				reg.Resource.Id == rs.Primary.ID &&
				reg.State != "removed" {
				return fmt.Errorf("Registry still exists")
			}
		}

		return nil
	}
	return nil
}

const testAccRancherRegistryConfig = //`
// resource "rancher_environment" "foo_registry" {
//   name = "registry test"
//   description = "environment to test registries"
// }
//
// resource "rancher_environment" "foo_registry2" {
//   name = "alternative registry test"
//   description = "other environment to test registries"
// }
`
resource "rancher_registry" "foo" {
  name = "foo"
  description = "registry test"
  server_address = "http://foo.com:8080"
  environment_id = "1a26"
}
`

//
// const testAccRancherRegistryUpdateConfig = `
// resource "rancher_environment" "foo_registry" {
//   name = "registry test"
//   description = "environment to test registries"
// }
//
// resource "rancher_environment" "foo_registry2" {
//   name = "alternative registry test"
//   description = "other environment to test registries"
// }
//
// resource "rancher_registry" "foo" {
//   name = "foo2"
//   description = "registry test - updated"
//   server_address = "http://foo.updated.com:8080"
//   environment_id = "${rancher_environment.foo_registry.id}"
// }
// `
//
// const testAccRancherRegistryRecreateConfig = `
// resource "rancher_environment" "foo_registry" {
//   name = "registry test"
//   description = "environment to test registries"
// }
//
// resource "rancher_environment" "foo_registry2" {
//   name = "alternative registry test"
//   description = "other environment to test registries"
// }
//
// resource "rancher_registry" "foo" {
//   name = "foo"
//   description = "registry test"
//   server_address = "http://foo.com:8080"
//   environment_id = "${rancher_environment.foo_registry.id}"
// }
// `
