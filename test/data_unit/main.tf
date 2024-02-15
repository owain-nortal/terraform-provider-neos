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

resource "neos_data_unit" "du-json-a6" {
  name        = "DUJSON-a6"
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


# resource "neos_data_unit" "du-json1" {
#   name        = "DUJSON1"
#   description = "desc test data unit 1"
#   owner       = "data unit owner"
#   label       = "DU1"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }



# resource "neos_data_unit" "du-json2" {
#   name        = "DUJSON2"
#   description = "desc test data unit 2"
#   owner       = "data unit owner"
#   label       = "DU2"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }

# resource "neos_data_unit" "du-json3" {
#   name        = "DUJSON3"
#   description = "desc test data unit 3"
#   owner       = "data unit owner"
#   label       = "DU3"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }

# resource "neos_data_unit" "du-json4" {
#   name        = "DUJSON4"
#   description = "desc test data unit 4"
#   owner       = "data unit owner"
#   label       = "DU4"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }

# resource "neos_data_unit" "du-json6" {
#   name        = "DUJSON6"
#   description = "desc test data unit 6"
#   owner       = "data unit owner"
#   label       = "DU6"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }


# resource "neos_data_unit" "du-json7" {
#   name        = "DUJSON7"
#   description = "desc test data unit 7"
#   owner       = "data unit owner"
#   label       = "DU7"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }


# resource "neos_data_unit" "du-json8" {
#   name        = "DUJSON8"
#   description = "desc test data unit 8"
#   owner       = "data unit owner"
#   label       = "DU8"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }

# resource "neos_data_unit" "du-json9" {
#   name        = "DUJSON9"
#   description = "desc test data unit 9"
#   owner       = "data unit owner"
#   label       = "DU9"
#   links       = var.links
#   contact_ids = var.contact_ids
#   config_json = <<EOH
#     {
#         "configuration": {
#             "data_unit_type": "csv",
#             "path": "/some/path",
#             "has_header": true,
#             "delimiter": ";",
#             "quote_char": "'",
#             "escape_char": "/"
#         }
#     }

#   EOH
# }