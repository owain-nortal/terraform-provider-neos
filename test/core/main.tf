terraform {
  required_providers {
    neos = {
      source = "registry.terraform.io/owain-nortal/neos"
    }
  }
}

provider "neos" {
  username  = "neosadmin"
  password  = "**"
  hub_host  = "owain10.neosdata.cloud"
  core_host = "owain10.neosdata.cloud"
  account   = "root"
  partition = "ksa"
}

resource "neos_registry_core" "testcore1" {
  partition = "ksa"
  name      = "testcore-1"
  account   = "owain10"
}

data "neos_registry_core" "cores" {
  account = "owain10"
  name   = "testcore-1"
}

output "foo" {
  value = data.neos_registry_core.cores.registry_cores[0].id
}


# output "access_key_id" {
#   value = neos_registry_core.testcore1.access_key_id
# }

# output "secret_key" {
#   value = neos_registry_core.testcore1.secret_key
# }

# output "urn" {
#   value = neos_registry_core.testcore1.urn
# }


