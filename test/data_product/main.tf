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

variable "links" {
  type    = list(any)
  default = ["link1", "link2"]
}

variable "contact_ids" {
  type    = list(any)
  default = ["contacts1", "contacts2"]
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
  connection_json = jsonencode({
    "connection" : {
      "connection_type" : "database",
      "engine" : "postgresql",
      "host" : "postgres-host",
      "port" : 5432,
      "schema" : "public",
      "database" : "test",
      "user" : { "env_key" : "DB_USERNAME" },
      "password" : { "env_key" : "DB_PASSWORD" }
    }
  })
  secret_json = jsonencode({
    "DB_PASSWORD" : "pass",
    "DB_USERNAME" : "user"
  })
}

resource "neos_link_data_system_data_source" "link_ds_ds" {
  parent_identifier = neos_data_system.json-system.id
  child_identifier  = neos_data_source.ds-json.id
}


resource "neos_link_data_source_data_unit" "link_ds_du" {
  parent_identifier = neos_data_source.ds-json.id
  child_identifier  = neos_data_unit.du-json11.id
}


resource "neos_data_unit" "du-json11" {
  name        = "DUJSON11"
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
  parent_identifier = neos_data_unit.du-json11.id
  child_identifier  = neos_data_product.dp-json.id
}

resource "neos_data_product" "dp-json" {
  name        = "dp-bike-report-d9"
  description = "The KPI on bike usage"
  owner       = "Richard C level"
  label       = "DPB"
  links       = [""]
  contact_ids = ["jane C level", "martin C level"]
  builder_json = jsonencode(

    {
      "config" : {
        "docker_tag" : "v0.3.64",
        "executor_cores" : 1,
        "executor_instances" : 2,
        "min_executor_instances" : 1,
        "max_executor_instances" : 3,
        "executor_memory" : "2048m",
        "driver_cores" : 1,
        "driver_core_limit" : "1200m",
        "driver_memory" : "1024m"
      },
      "inputs" : {
        format("input_%s", replace(neos_data_unit.du-json11.id, "-", "_")) : {
          "input_type" : "data_unit",
          "identifier" : neos_data_unit.du-json11.id,
          "preview_limit" : 10
        }
      },
      "transformations" : [
        {
          "transform" : "rename_column",
          "input" : format("input_%s", replace(neos_data_unit.du-json11.id, "-", "_")),
          "output" : "rename_lowercase",
          "changes" : {
            "dept_no" : "dept_no",
            "emp_no" : "emp_no",
            "from_date" : "from_date",
            "to_date" : "to_date"
          }
        },
        {
          "transform" : "select_columns",
          "input" : "rename_lowercase",
          "output" : "select_all",
          "columns" : [
            "emp_no",
            "dept_no",
            "from_date",
            "to_date"
          ]
        }
      ],
      "finalisers" : {
        "input" : "select_all",
        "enable_quality" : true,
        "write_config" : {
          "mode" : "overwrite"
        },
        "enable_profiling" : true,
        "enable_classification" : true
      },
      "preview" : false
    }
  )
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
