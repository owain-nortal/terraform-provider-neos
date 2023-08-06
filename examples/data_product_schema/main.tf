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



resource "neos_data_product" "foo-bikes-product" {
  name        = "foo-bike-report5"
  description = "The KPI on bike usage"
  owner       = "Richard C level"
  label       = "FOO"
  links       = [""]
  contact_ids = ["jane C level", "martin C level"]
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
