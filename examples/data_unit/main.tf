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
//data "neos_links" "links" {}

data "neos_data_unit" "data_units" {
}


output "data_unit" {
  value = data.neos_data_unit.data_units
}


# output "neos_links" {
#   value = data.neos_links
# }

# resource "neos_link_data_source_data_unit" {
#   parent_identifier = "2e6841d9-349c-495d-b57b-7ae5f7fc54da"
#   child_identifier  = "0493365a-80f7-4fb9-802e-1d902f8155b2"
# }


# resource "neos_registry_core" "testcore1" {
#   partition = "ksa"
#   name      = "owain-test3"
# }

# output "access_key" {
#   value = neos_registry_core.testcore1.access_key
# }


# output "edu_data_system" {
#   value = data.neos_data_system.edu
# }

variable "links" {
  type    = list(any)
  default = ["link1", "link2"]
}

variable "contact_ids" {
  type    = list(any)
  default = ["contacts1", "contacts2"]
}





# resource "neos_data_system" "ds-test1" {
#   name        = "OneTestDataSystem1"
#   description = "desc test data system 1"
#   owner       = "data system owner"
#   label       = "ODS"
#   links       = var.links
#   contact_ids = var.contact_ids
# }

# resource "neos_data_source" "ds-test1" {
#   name        = "OneDataProduct1"
#   description = "desc test data product 1"
#   owner       = "data source owner"
#   label       = "ODS"
#   links       = var.links
#   contact_ids = var.contact_ids
# }

# resource "neos_data_product" "dp-test1" {
#   name        = "OneDataProduct1"
#   description = "desc test data product 1"
#   owner       = "data product owner"
#   label       = "ODP"
#   links       = var.links
#   contact_ids = var.contact_ids
#   schema = {
#     fields = [
#       {
#         "name"        = "string"
#         "description" = "string"
#         "primary"     = false
#         "optional"    = false
#         "data_type" = {
#           "meta" = [
#             { "foo" : "base" }  
#           ],
#           "column_type" : "STRING"
#         },
#         "tags" = ["string"]
#       }
#     ]
#   }
# }




resource "neos_data_unit" "du-json" {
  name        = "DUJSON"
  description = "desc test data unit 1"
  owner       = "data unit owner"
  label       = "DUJ"
  links       = var.links
  contact_ids = var.contact_ids
  config_json = <<EOH
    {
        "configuration": {
            "data_unit_type": "csv",
            "path": "/some/path",
            "has_header": true,
            "delimiter": ";",
            "quote_char": "'",
            "escape_char": "/"
        }
    }

  EOH
}

# resource "neos_data_unit" "du-table" {
#   name        = "DUTable"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "T01"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config = {
#     data_unit_type = "table"
#     table = {
#       table = "users"
#     }
#   }
# }

# resource "neos_data_unit" "du-csv" {
#   name        = "DUCSV"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "CSV"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config = {
#     data_unit_type = "csv"
#     csv = {
#       path        = "/some/path"
#       delimiter   = ","
#       has_header  = true
#       escape_char = "/"
#       quote_char  = "'"
#     }
#   }
# }

# resource "neos_data_unit" "du-parquet" {
#   name        = "DUParquet"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "PAR"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config = {
#     data_unit_type = "parquet"
#     parquet        = {}
#   }
# }

# resource "neos_data_unit" "du-data_product" {
#   name        = "DU-DP"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "DP1"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config = {
#     data_unit_type = "data_product"
#     data_product = {
#       engine = "postgres"
#       table  = "some-table"
#     }
#   }
# }
# resource "neos_data_unit" "du-test1" {
#   name        = "OneDataUnit1"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "ODU"
#   links       = var.links
#   contact_ids = var.contact_ids
# }
