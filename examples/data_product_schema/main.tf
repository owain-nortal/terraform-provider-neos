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



resource "neos_data_product" "foo-bikes-product-d" {
  name         = "foo-bike-report-d"
  description  = "The KPI on bike usage"
  owner        = "Richard C level"
  label        = "FOO"
  links        = [""]
  contact_ids  = ["jane C level", "martin C level"]
  builder_json = <<EOF
          {
            "config": {
                "executor_cores": 1,
                "executor_instances": 1,
                "min_executor_instances": 1,
                "max_executor_instances": 1,
                "executor_memory": "512m",
                "driver_cores": 1,
                "driver_core_limit": "1200m",
                "driver_memory": "512m",
                "docker_tag": "v0.3.23"
            },
            "inputs": {
                "input": {
                "input_type": "data_unit",
                "identifier": "3d3bc3b6-d1b8-4988-b5d9-933e6a40e67d",
                "preview_limit": 10
                }
            },
            "transformations": [
                {
                "transform": "select_columns",
                "input": "input",
                "output": "after_select",
                "columns": [
                    "foo",
                    "year"
                ]
                },
                {
                "transform": "cast",
                "input": "after_select",
                "output": "after_cast",
                "changes": [
                    {
                    "column": "year",
                    "data_type": "integer"
                    }
                ]
                }
            ],
            "finalisers": [
                {
                "finaliser": "run_data_quality",
                "input": "after_cast"
                },
                {
                "finaliser": "run_data_profiling",
                "input": "after_cast"
                },
                {
                "finaliser": "save_dataframe",
                "input": "after_cast",
                "write_mode": "overwrite"
                }
            ],
            "preview": false
            }
          EOF
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
