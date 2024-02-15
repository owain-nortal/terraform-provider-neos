terraform {
  required_providers {
    neos = {
      source = "registry.terraform.io/owain-nortal/neos"
    }
  }
}

provider "neos" {
  username  = "neosadmin"
  password  = "**Gwen11"
  hub_host  = "owain10.neosdata.cloud"
  core_host = "owain10.neosdata.cloud"
  account   = "root"
  partition = "ksa"
}

resource "neos_registry_core" "testcore1" {
  partition = "ksa"
  name      = "testcore-6a"
}

output "access_key_id" {
  value = neos_registry_core.testcore1.access_key_id
}

output "secret_key" {
  value = neos_registry_core.testcore1.secret_key
}

output "urn" {
  value = neos_registry_core.testcore1.urn
}


