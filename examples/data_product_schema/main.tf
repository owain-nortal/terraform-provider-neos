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

variable "links" {
  type    = list(any)
  default = ["link1", "link2"]
}

variable "contact_ids" {
  type    = list(any)
  default = ["contacts1", "contacts2"]
}

variable "driver_memory" {
  default = "512m"
  type    = string
}

variable "executor_memory" {
  default = "512m"
  type    = string
}

variable "driver_core_limit" {
  default = "1200m"
  type    = string
}

provider "neos" {
  username      = "owain.perry"
  password      = var.password
  iam_host      = "sandbox.city3os.com"
  core_host     = "op-02.neosdata.net"
  registry_host = "sandbox.city3os.com"
}

resource "neos_data_system" "json-system" {
  name        = "TransportForLondon"
  description = "the xy data system"
  owner       = "pete jones"
  label       = "DJS"
  links       = ["link1", "link2"]
  contact_ids = ["contacts1", "contacts2"]
}

resource "neos_data_source" "ds-json" {
  name        = "datasource-json"
  description = "Test datasource"
  owner       = "john londoner "
  label       = "TFA"
  links       = ["link1", "link2"]
  contact_ids = ["contacts1", "contacts2"]
}

resource "neos_link_data_system_data_source" "link_ds_ds" {
  parent_identifier = neos_data_system.json-system.id
  child_identifier  = neos_data_source.ds-json.id
}


resource "neos_link_data_source_data_unit" "link_ds_du" {
  parent_identifier = neos_data_source.ds-json.id
  child_identifier  = neos_data_unit.du-json.id
}


resource "neos_data_unit" "du-json" {
  name        = "DUJSON"
  description = "desc test data unit 1"
  owner       = "data unit owner"
  label       = "DUJ"
  links       = var.links
  contact_ids = var.contact_ids
  config_json = jsonencode({
    "configuration" : {
      "data_unit_type" : "query",
      "query" : "select * from ilr.worker"
    }
  })
}

resource "neos_link_data_unit_data_product" "du_dp" {
  parent_identifier = neos_data_unit.du-json.id
  child_identifier  = neos_data_product.dp-json.id
}

resource "neos_data_product" "dp-json" {
  name        = "dp-bike-report-d"
  description = "The KPI on bike usage"
  owner       = "Richard C level"
  label       = "DPB"
  links       = [""]
  contact_ids = ["jane C level", "martin C level"]
  # builder_json = jsonencode(

  #   {
  #     "config" : {
  #       "executor_cores" : 1,
  #       "executor_instances" : 1,
  #       "min_executor_instances" : 1,
  #       "max_executor_instances" : 1,
  #       "executor_memory" : "512m",
  #       "driver_cores" : 1,
  #       "driver_core_limit" : "1200m",
  #       "driver_memory" : "512m",
  #       "docker_tag" : "v0.3.23"
  #     },
  #     "inputs" : {
  #       "tdp" : {
  #         "input_type" : "data_unit",
  #         "identifier" : neos_data_unit.du-json.id,
  #         "preview_limit" : 10
  #       }
  #     },
  #     "transformations" : [
  #       {
  #         "transform" : "select_columns",
  #         "input" : "tdp",
  #         "output" : "after_select",
  #         "columns" : [
  #           "worker_first_name"
  #         ]
  #       }
  #     ],
  #     "finalisers" : [
  #       {
  #         "finaliser" : "run_data_quality",
  #         "input" : "after_select"
  #       },
  #       {
  #         "finaliser" : "run_data_profiling",
  #         "input" : "after_select"
  #       },
  #       {
  #         "finaliser" : "save_dataframe",
  #         "input" : "after_select",
  #         "write_mode" : "overwrite"
  #       }
  #     ],
  #     "preview" : false
  #   }
  # )
  schema = {
    product_type = "stored"
    fields = [
      {
        "name"        = "booking_id"
        "description" = "the booking id of the thing"
        "primary"     = true
        "optional"    = false
        "data_type" = {
          "meta" = {
            "length" : "100"
          }

          "column_type" : "VARCHAR"
        }
      },
      {
        "name"        = "firstname"
        "description" = "the firstname of the person"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = {
            "length" : "20"
          }

          "column_type" : "VARCHAR"
        }
      },
      {
        "name"        = "lastname"
        "description" = "the lastname of the person"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = {
            "length" : "20"
          }

          "column_type" : "VARCHAR"
        }
      },
      {
        "name"        = "age"
        "description" = "the age of the person"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = {
            "length" : "20"
          }

          "column_type" : "VARCHAR"
        }
      },
      {
        "name"        = "passport"
        "description" = "the passport number"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = {
            "length" : "20"
          }

          "column_type" : "VARCHAR"
        }
      },
      {
        "name"        = "address1"
        "description" = "the first line of the address"
        "primary"     = false
        "optional"    = false
        "data_type" = {
          "meta" = {
            "length" : "40"
          }

          "column_type" : "VARCHAR"
        }
      }

    ]
  }
}
