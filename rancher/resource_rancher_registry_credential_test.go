package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	rancher "github.com/rancher/go-rancher/client"
)

func TestAccRancherRegistryCredential(t *testing.T) {
	var registry rancher.RegistryCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherRegistryCredentialDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherRegistryCredentialConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherRegistryCredentialExists("rancher_registry_credential.foo", &registry),
					testAccCheckRancherRegistryCredentialAttributes(&registry, "foo", "registry credential test", "user"),
				),
			},
			resource.TestStep{
				Config: testAccRancherRegistryCredentialUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherRegistryCredentialExists("rancher_registry_credential.foo", &registry),
					testAccCheckRancherRegistryCredentialAttributes(&registry, "foo2", "registry credential test - updated", "user2"),
				),
			},
		},
	})
}

func testAccCheckRancherRegistryCredentialExists(n string, reg *rancher.RegistryCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Name is set")
		}

		client := testAccProvider.Meta().(*Config)

		foundReg, err := client.RegistryCredential.ById(rs.Primary.ID)
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

func testAccCheckRancherRegistryCredentialAttributes(reg *rancher.RegistryCredential, regName string, regDesc string, regPublic string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if reg.Name != regName {
			return fmt.Errorf("Bad name: %s shoud be: %s", reg.Name, regName)
		}

		if reg.Description != regDesc {
			return fmt.Errorf("Bad description: %s shoud be: %s", reg.Description, regDesc)
		}

		if reg.PublicValue != regPublic {
			return fmt.Errorf("Bad public_value: %s shoud be: %s", reg.PublicValue, regPublic)
		}

		return nil
	}
}

func testAccCheckRancherRegistryCredentialDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_registry_credential" {
			continue
		}
		reg, err := client.RegistryCredential.ById(rs.Primary.ID)

		if err == nil {
			if reg != nil &&
				reg.Resource.Id == rs.Primary.ID &&
				reg.State != "removed" {
				return fmt.Errorf("RegistryCredential still exists")
			}
		}

		return nil
	}
	return nil
}

const testAccRancherRegistryCredentialConfig = `
resource "rancher_registry" "foo" {
  name = "foo"
  description = "registry test"
  server_address = "http://bar.com:8080"
  environment_id = "1a5"
}

resource "rancher_registry_credential" "foo" {
	name = "foo"
	description = "registry credential test"
	registry_id = "${rancher_registry.foo.id}"
	email = "registry@credential.com"
	public_value = "user"
	secret_value = "pass"
}
`

const testAccRancherRegistryCredentialUpdateConfig = `
resource "rancher_registry" "foo" {
  name = "foo"
  description = "registry test"
  server_address = "http://bar.com:8080"
  environment_id = "1a5"
}

resource "rancher_registry_credential" "foo" {
	name = "foo2"
	description = "registry credential test - updated"
	registry_id = "${rancher_registry.foo.id}"
	email = "registry@credential.com"
	public_value = "user2"
	secret_value = "pass"
}
 `
