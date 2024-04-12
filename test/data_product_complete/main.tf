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

resource "neos_data_unit" "foo_unit" {
  name        = "foobar"
  description = "description"
  owner       = "foo team"
  label       = "FOO"
  links       = []
  contact_ids = []
  config_json = <<EOH
    {
      "configuration": 
        {
          "data_unit_type": "csv",
          "delimiter": ";",
          "escape_char": null,
          "has_header": true,
          "path": "data/foo.csv",
          "quote_char": null
        }
    }
  EOH

}

resource "neos_data_source" "foo_source" {
  name            = "foosrc"
  description     = "foo source"
  owner           = "fooybar"
  label           = "BAR"
  links           = []
  contact_ids     = []
  connection_json = jsonencode({ "connection" : { "access_key" : { "env_key" : "MINIO_S3_ACCESS" }, "access_secret" : { "env_key" : "MINIO_S3_SECRET" }, "connection_type" : "s3", "url" : "http://minio.neos-core-minio.svc.cluster.local" } })
  secret_values = {
    MINIO_S3_ACCESS = "abc123"
    MINIO_S3_SECRET = "xyz456"
  }
}


resource "neos_link_data_source_data_unit" "link_foo_source_foo_unit" {
  parent_identifier = neos_data_source.foo_source.id
  child_identifier  = neos_data_unit.foo_unit.id
}

resource "neos_data_product" "test-dp" {
  name        = "test-full-dp"
  description = "full dp test"
  owner       = "Richard C level"
  label       = "FDP"
  links       = [""]
  contact_ids = ["jane C level"]
  schema = {
    product_type = "stored"
    fields = [
      {
        "name"        = "application_status"
        "description" = ""
        "primary"     = false
        "optional"    = true
        "data_type" = {
          meta = {}
          "column_type" : "VARCHAR"
        }
      },
    ]
  }
}


resource "neos_data_product_builder" "test-dp" {
  id                         = neos_data_product.test-dp.id
  dataunit_datasource_linkids = [
    neos_link_data_source_data_unit.link_foo_source_foo_unit.id
  ]
  builder_json = jsonencode(
    {
      "config" : {
        "docker_tag" : "v0.3.64",
        "executor_cores" : 2,
        "executor_instances" : 2,
        "min_executor_instances" : 1,
        "max_executor_instances" : null,
        "executor_memory" : "1024m",
        "driver_cores" : 2,
        "driver_core_limit" : "2400m",
        "driver_memory" : "1024m"
      },
      "inputs" : {
        format("input_%s", replace(neos_data_unit.foo_unit.id, "-", "_")) : {
          "input_type" : "data_unit",
          "identifier" : neos_data_unit.foo_unit.id,
          "preview_limit" : 10
        }
      },
      "transformations" : [
        {
          "transform" : "select_columns",
          "input" : format("input_%s", replace(neos_data_unit.foo_unit.id, "-", "_")),
          "output" : "after_selectrec",
          "columns" : [
            "application_status",
            "on_boarding_candidates",
            "most_advanced_candidate_flag"
          ]
        }
      ],
      "finalisers" : {
        "input" : "after_selectrec",
        "enable_quality" : true,
        "write_config" : {
          "mode" : "overwrite"
        },
        "enable_profiling" : true,
        "enable_classification" : false
      },

    }
  )
}
