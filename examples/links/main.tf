terraform {
  required_providers {
    neos = {
      source = "registry.terraform.io/owain-nortal/neos"
    }
  }
}

variable "password" {
  default = ""
  type    = string
}


provider "neos" {
  username      = "owain.perry"
  password      = var.password
  iam_host      = "sandbox.city3os.com"
  core_host     = "op-02.neosdata.net"
  registry_host = "sandbox.city3os.com"
}


resource "neos_data_system" "tfl-system" {
  name        = "TransportForLondon"
  description = "the xy data system"
  owner       = "pete jones"
  label       = "TFL"
  links       = ["link1", "link2"]
  contact_ids = ["contacts1", "contacts2"]
}

resource "neos_data_source" "tfl-api-source" {
  name        = "TFL-API"
  description = "The transport for london API"
  owner       = "john londoner "
  label       = "TFA"
  links       = ["link1", "link2"]
  contact_ids = ["contacts1", "contacts2"]
}

resource "neos_link_data_system_data_source" "link_ds_ds" {
  parent_identifier = neos_data_system.tfl-system.id
  child_identifier  = neos_data_source.tfl-api-source.id
}

resource "neos_data_unit" "tfl-bikes-unit" {
  name        = "TFL-BIKES"
  description = "The bikes data unit"
  owner       = "dave peleton"
  label       = "TFB"
  links       = ["some line"]
  contact_ids = ["mark davies", "john smith"]
}

resource "neos_link_data_source_data_unit" "tfl_ds_du" {
  parent_identifier = neos_data_source.tfl-api-source.id
  child_identifier  = neos_data_unit.tfl-bikes-unit.id
}


resource "neos_data_product" "tfl-bikes-product" {
  name        = "TFL-bike-report"
  description = "The KPI on bike usage"
  owner       = "Richard C level"
  label       = "TFR"
  links       = [""]
  contact_ids = ["jane C level", "martin C level"]
  schema = {
    fields = [
      {
        "name"        = "string"
        "description" = "string"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = [
            { "foo" : "base" }
          ],
          "column_type" : "STRING"
        },
        "tags" = ["string"]
      }
    ]
  }
}

resource "neos_link_data_unit_data_product" "tfl_du_dp" {
  parent_identifier = neos_data_unit.tfl-bikes-unit.id
  child_identifier  = neos_data_product.tfl-bikes-product.id
}


resource "neos_output" "tfl-bikes-dashboard" {
  name        = "TFL-bike-dashboard"
  description = "The bike usage dashboard"
  owner       = "Mark Board"
  label       = "TFD"
  links       = ["abc"]
  contact_ids = ["eric dashy"]
  output_type = "dashboard"
}

resource "neos_link_data_product_output" "tfl_dp_output" {
  parent_identifier = neos_data_product.tfl-bikes-product.id
  child_identifier  = neos_output.tfl-bikes-dashboard.id
}


resource "neos_output" "tfl-bikes-application" {
  name        = "TFL-bike-application"
  description = "The bike usage application"
  owner       = "Dave App"
  label       = "TFA"
  links       = ["xyz"]
  contact_ids = ["harry appy"]
  output_type = "application"
}

resource "neos_link_data_product_output" "tfl_out_app" {
  parent_identifier = neos_data_product.tfl-bikes-product.id
  child_identifier  = neos_output.tfl-bikes-application.id
}


resource "neos_data_product" "uk-bike-capacity-product" {
  name        = "uk-bike-capacity"
  description = "The uk bike capacity product"
  owner       = "Albert capacity "
  label       = "UKC"
  links       = [""]
  contact_ids = ["bungle", "george"]
  schema = {
    fields = [
      {
        "name"        = "string"
        "description" = "string"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = [
            { "foo" : "base" }
          ],
          "column_type" : "STRING"
        },
        "tags" = ["string"]
      }
    ]
  }

}

resource "neos_link_data_product_data_product" "uk_capacity" {
  parent_identifier = neos_data_product.tfl-bikes-product.id
  child_identifier  = neos_data_product.uk-bike-capacity-product.id
}


# variable "links" {
#   type    = list(any)
#   default = 
# }

# variable "contact_ids" {
#   type    = list(any)
#   default = 
# }




# data "neos_registry_core" "cores" {
# }

# output "neos_cores" {
#   value = data.neos_registry_core.cores
# }



data "neos_links" "links" {}

# output "neos_links" {
#   value = data.neos_links.links
# }
