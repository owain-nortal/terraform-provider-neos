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

# variable "links" {
#   type    = list(any)
#   default = ["link1", "link2"]
# }

# variable "contact_ids" {
#   type    = list(any)
#   default = ["contacts1", "contacts2"]
# }



# resource "neos_data_system" "json-system" {
#   name        = "TransportForLondon"
#   description = "the xy data system"
#   owner       = "pete jones"
#   label       = "DJS"
#   links       = ["link1", "link2"]
#   contact_ids = ["contacts1", "contacts2"]
# }

# resource "neos_data_source" "ds-json" {
#   name        = "datasource-json"
#   description = "Test datasource"
#   owner       = "john londoner "
#   label       = "TFA"
#   links       = ["link1", "link2"]
#   contact_ids = ["contacts1", "contacts2"]
#   connection_json = jsonencode({
#     "connection" : {
#       "connection_type" : "database",
#       "engine" : "postgresql",
#       "host" : "postgres-host",
#       "port" : 5432,
#       "schema" : "public",
#       "database" : "test",
#       "user" : { "env_key" : "DB_USERNAME" },
#       "password" : { "env_key" : "DB_PASSWORD" }
#     }
#   })
#   secret_json = jsonencode({
#     "DB_PASSWORD" : "pass",
#     "DB_USERNAME" : "user"
#   })
# }

# resource "neos_link_data_system_data_source" "link_ds_ds" {
#   parent_identifier = neos_data_system.json-system.id
#   child_identifier  = neos_data_source.ds-json.id
# }


# resource "neos_link_data_source_data_unit" "link_ds_du" {
#   parent_identifier = neos_data_source.ds-json.id
#   child_identifier  = neos_data_unit.du-json11.id
# }


# resource "neos_data_unit" "du-json11" {
#   name        = "DUJSON11"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "DUJ"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = jsonencode({
#     "configuration" : {
#       "data_unit_type" : "query",
#       "query" : "select * from ilr.worker"
#     }
#   })
# }

# resource "neos_link_data_unit_data_product" "du_dp" {
#   parent_identifier = neos_data_unit.du-json11.id
#   child_identifier  = neos_data_product.dp-json.id
# }

resource "neos_data_product" "dp-json" {
  name        = "dp-bike-report-d9"
  description = "The KPI on bike usage"
  owner       = "Richard C level"
  label       = "DPB"
  links       = [""]
  contact_ids = ["jane C level", "martin C level"]
  schema = {}
}