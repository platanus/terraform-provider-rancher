package rancher

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	rancher "github.com/rancher/go-rancher/client"
)

func TestAccRancherStack(t *testing.T) {
	var stack rancher.Environment

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckRancherStackDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccRancherStackConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherStackExists("rancher_stack.foo", &stack),
					testAccCheckRancherStackAttributes(&stack, "foo", "Terraform acc test group", "", "", "", emptyEnvironment, false),
				),
			},
			resource.TestStep{
				Config: testAccRancherStackUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherStackExists("rancher_stack.foo", &stack),
					testAccCheckRancherStackAttributes(&stack, "foo2", "Terraform acc test group - updated", "", "", "", emptyEnvironment, false),
				),
			},
			resource.TestStep{
				Config: testAccRancherStackComposeConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherStackExists("rancher_stack.compose", &stack),
					testAccCheckRancherStackAttributes(&stack, "compose", "Terraform acc test group - compose", "", "web: { image: nginx }", "web: { scale: 1 }", emptyEnvironment, true),
				),
			},
			resource.TestStep{
				Config: testAccRancherStackSystemCatalogConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRancherStackExists("rancher_stack.catalog", &stack),
					testAccCheckRancherStackAttributes(&stack, "catalog", "Terraform acc test group - catalog", "system-catalog://library:route53:7", route53DockerCompose, route53RancherCompose, route53Environment, false),
				),
			},
		},
	})
}

func testAccCheckRancherStackExists(n string, stack *rancher.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Name is set")
		}

		client := testAccProvider.Meta().(*Config)

		foundStack, err := client.Environment.ById(rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundStack.Resource.Id != rs.Primary.ID {
			return fmt.Errorf("Stack not found")
		}

		*stack = *foundStack

		return nil
	}
}

func testAccCheckRancherStackAttributes(stack *rancher.Environment, stackName string, stackDesc string, externalID string, dockerCompose string, rancherCompose string, environment map[string]string, startOnCreate bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if stack.Name != stackName {
			return fmt.Errorf("Bad name: %s should be: %s", stack.Name, stackName)
		}

		if stack.Description != stackDesc {
			return fmt.Errorf("Bad description: %s should be: %s", stack.Description, stackDesc)
		}

		if stack.ExternalId != externalID {
			return fmt.Errorf("Bad externalID: %s should be: %s", stack.ExternalId, externalID)
		}

		if stack.DockerCompose != dockerCompose {
			return fmt.Errorf("Bad dockerCompose: %s should be: %s", stack.DockerCompose, dockerCompose)
		}

		if stack.RancherCompose != rancherCompose {
			return fmt.Errorf("Bad rancherCompose: %s should be: %s", stack.RancherCompose, rancherCompose)
		}

		if len(stack.Environment) != len(environment) {
			return fmt.Errorf("Bad environment size: %v should be: %v", len(stack.Environment), environment)
		}

		for k, v := range stack.Environment {
			if environment[k] != v {
				return fmt.Errorf("Bad environment value for %s: %s should be: %s", k, environment[k], v)
			}
		}

		if stack.StartOnCreate != startOnCreate {
			return fmt.Errorf("Bad startOnCreate: %s should be: %s", stack.StartOnCreate, startOnCreate)
		}

		return nil
	}
}

func testAccCheckRancherStackDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*Config)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "rancher_stack" {
			continue
		}
		stack, err := client.Environment.ById(rs.Primary.ID)

		if err == nil {
			if stack != nil &&
				stack.Resource.Id == rs.Primary.ID &&
				stack.State != "removed" {
				return fmt.Errorf("Stack still exists")
			}
		}

		return nil
	}
	return nil
}

const testAccRancherStackConfig = `
resource "rancher_stack" "foo" {
	name = "foo"
	description = "Terraform acc test group"
	environment_id = "1a5"
}
`

const testAccRancherStackUpdateConfig = `
resource "rancher_stack" "foo" {
	name = "foo2"
	description = "Terraform acc test group - updated"
	environment_id = "1a5"
}
`

const testAccRancherStackComposeConfig = `
resource "rancher_stack" "compose" {
	name = "compose"
	description = "Terraform acc test group - compose"
	environment_id = "1a5"
	docker_compose = "web: { image: nginx }"
	rancher_compose = "web: { scale: 1 }"
	start_on_create = true
}
`

const testAccRancherStackSystemCatalogConfig = `
resource "rancher_stack" "catalog" {
	name = "catalog"
	description = "Terraform acc test group - catalog"
	environment_id = "1a5"
	catalog_id = "library:route53:7"
	scope = "system"
	environment {
		AWS_ACCESS_KEY = "MYKEY"
		AWS_SECRET_KEY = "MYSECRET"
		AWS_REGION = "eu-central-1"
		ROOT_DOMAIN = "example.com"
	}
}
`

const route53DockerCompose = `route53:
  image: rancher/external-dns:v0.5.0
  expose: 
   - 1000
  environment:
    AWS_ACCESS_KEY: ${AWS_ACCESS_KEY}
    AWS_SECRET_KEY: ${AWS_SECRET_KEY}
    AWS_REGION: ${AWS_REGION}
    ROOT_DOMAIN: ${ROOT_DOMAIN}
    ROUTE53_ZONE_ID: ${ROUTE53_ZONE_ID}
    TTL: ${TTL}
  labels:
    io.rancher.container.create_agent: "true"
    io.rancher.container.agent.role: "external-dns"
`

const route53RancherCompose = `.catalog:
  name: "Route53 DNS"
  version: "v0.5.0-rancher1"
  description: "Rancher External DNS service powered by Amazon Route53. Requires Rancher version 0.44.0"
  minimum_rancher_version: v1.1.0
  questions:
    - variable: "AWS_ACCESS_KEY"
      label: "AWS Access Key ID"
      description: "Access key ID for your AWS account"
      type: "string"
      required: true
    - variable: "AWS_SECRET_KEY"
      label: "AWS Secret Access Key"
      description: "Secret access key for your AWS account"
      type: "string"
      required: true
    - variable: "AWS_REGION"
      label: "AWS Region"
      description: "AWS region name"
      type: "string"
      default: "us-west-2"
      required: true
    - variable: "TTL"
      label: "TTL"
      description: "The resource record cache time to live (TTL), in seconds"
      type: "int"
      default: 60
      required: false
    - variable: "ROOT_DOMAIN"
      label: "Hosted Zone Name"
      description: "Route53 hosted zone name (zone has to be pre-created). DNS entries will be created for <service>.<stack>.<environment>.<hosted zone>"
      type: "string"
      required: true
    - variable: "ROUTE53_ZONE_ID"
      label: "Hosted Zone ID"
      description: "If there are multiple zones with the same name, then you must additionally specify the ID of the hosted zone to use."
      type: "string"
      required: false
    - variable: "HEALTH_CHECK_INTERVAL"
      label: "Health Check Interval"
      description: "The health check interval for this service, in seconds. Raise this value if the total requests from your AWS account exceed the Route53 API rate limits."
      type: "int"
      min: 1
      max: 999
      default: 15
      required: false

route53:
  health_check:
    port: 1000
    interval: ${HEALTH_CHECK_INTERVAL}000
    unhealthy_threshold: 2
    request_line: GET / HTTP/1.0
    healthy_threshold: 2
    response_timeout: 5000
`

var emptyEnvironment = map[string]string{}

var route53Environment = map[string]string{
	"AWS_ACCESS_KEY": "MYKEY",
	"AWS_SECRET_KEY": "MYSECRET",
	"AWS_REGION":     "eu-central-1",
	"ROOT_DOMAIN":    "example.com",
}
