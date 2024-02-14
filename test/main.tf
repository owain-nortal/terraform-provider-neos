terraform {
  required_providers {
    neos = {
      source = "registry.terraform.io/owain-nortal/neos"
    }
  }
}

provider "neos" {
  username  = "neosadmin"
  password  = ""
  hub_host  = "owain10.neosdata.cloud"
  core_host = "owain10.neosdata.cloud"
  account   = "root"
  partition = "ksa"
}

//data "neos_data_system" "example" {}

# data "neos_registry_core" "cores" {
# }

# output "edu_data_system" {
#   value = data.neos_registry_core.cores
# }

# data "neos_data_system" "edu" {
# }


resource "neos_registry_core" "testcore1" {
  partition = "ksa"
  name      = "testcore-5a"
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



# output "urn" {
#   value = neos_registry_core.testcore1.urn
# }


# output "edu_data_system" {
#   value = data.neos_data_system.edu
# }

# variable "links" {

#   type    = list(any)
#   default = ["link1", "link2"]
# }

# variable "contact_ids" {
#   type    = list(any)
#   default = ["contacts1", "contacts2"]
# }

# resource "neos_data_system" "op-test1" {
#   name        = "APTestDataSystem"
#   description = "desc test data system 2"
#   owner       = "test data system 2 owner"
#   label       = "AP2"
#   links       = var.links
#   contact_ids = var.contact_ids
# }

# resource "neos_data_system" "op-test2" {
#   name        = "APTestDataSystem3"
#   description = "desc test data system 3"
#   owner       = "test data system 3 owner"
#   label       = "AP3"
#   links       = var.links
#   contact_ids = var.contact_ids
# }

# resource "neos_data_system" "dp-test1" {
#   name        = "APTestDataProduct1"
#   description = "desc test data product 1"
#   owner       = "test data product 1"
#   label       = "AP1"
#   links       = var.links
#   contact_ids = var.contact_ids
# }
