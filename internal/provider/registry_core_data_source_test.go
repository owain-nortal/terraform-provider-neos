package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRegistryCoreDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Read testing
			{
				Config: providerConfig + `data "neos_registry_core" "test" {}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify number of coffees returned
					resource.TestCheckResourceAttr("data.neos_registry_core.test", "registry_cores.#", "9"),
					// Verify the first coffee to ensure all attributes are set
					//resource.TestCheckResourceAttr("data.neos_registry_core.test", "registry_cores.0.description", ""),
					//resource.TestCheckResourceAttr("data.neos_registry_core.test", "registry_cores.0.id", "1"),
					// resource.TestCheckResourceAttr("data.neos_registry_core.test", "registry_cores.0.image", "/hashicorp.png"),
					// resource.TestCheckResourceAttr("data.neos_registry_core.test", "coffees.0.ingredients.#", "1"),
					// resource.TestCheckResourceAttr("data.neos_registry_core.test", "coffees.0.ingredients.0.id", "6"),
					// resource.TestCheckResourceAttr("data.neos_registry_core.test", "coffees.0.name", "HCP Aeropress"),
					// resource.TestCheckResourceAttr("data.neos_registry_core.test", "coffees.0.price", "200"),
					// resource.TestCheckResourceAttr("data.neos_registry_core.test", "coffees.0.teaser", "Automation in a cup"),
					// Verify placeholder id attribute
					resource.TestCheckResourceAttr("data.neos_registry_core.test", "registry_cores.0.urn", "00000000-0000-0000-0000-000000000001"),
				),
			},
		},
	})
}
