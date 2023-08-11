curl -X 'PUT' \
  'https://op-02.neosdata.net/api/gateway/v2/data_product/95bfe9fb-2aa7-46f5-8ff5-5994dde9e25f/spark/builder' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJROERBdmJQczBWa2hwb0hqaUhnUHJRSlowdmhaN1NQV204R1B1bFhDWV9FIn0.eyJleHAiOjE2OTE0NjkxODUsImlhdCI6MTY5MTQ0MDM4NSwianRpIjoiODY5N2E3YjgtZDUwMi00NDUwLWI2NGYtZGYzMDUyMjk3YTQ3IiwiaXNzIjoiaHR0cHM6Ly9hdXRoLnNhbmRib3guY2l0eTNvcy5jb20vcmVhbG1zL25lb3MiLCJhdWQiOiJhY2NvdW50Iiwic3ViIjoiMjVhMmViODAtMDdiOS00ZGVkLTg3ZDAtNzZlOGY3NzNlYTliIiwidHlwIjoiQmVhcmVyIiwiYXpwIjoibmVvcy1pYW0iLCJzZXNzaW9uX3N0YXRlIjoiNGY3YWFjZmQtMzg0ZC00YWQ1LWEwZTUtNTA1NTQ1MjQ2ZmEwIiwiYWNyIjoiMSIsInJlYWxtX2FjY2VzcyI6eyJyb2xlcyI6WyJvZmZsaW5lX2FjY2VzcyIsInVtYV9hdXRob3JpemF0aW9uIiwiZGVmYXVsdC1yb2xlcy1uZW9zIl19LCJyZXNvdXJjZV9hY2Nlc3MiOnsiYWNjb3VudCI6eyJyb2xlcyI6WyJtYW5hZ2UtYWNjb3VudCIsIm1hbmFnZS1hY2NvdW50LWxpbmtzIiwidmlldy1wcm9maWxlIl19fSwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCIsInNpZCI6IjRmN2FhY2ZkLTM4NGQtNGFkNS1hMGU1LTUwNTU0NTI0NmZhMCIsImVtYWlsX3ZlcmlmaWVkIjpmYWxzZSwibmFtZSI6Ik93YWluIFBlcnJ5IiwicHJlZmVycmVkX3VzZXJuYW1lIjoib3dhaW4ucGVycnkiLCJnaXZlbl9uYW1lIjoiT3dhaW4iLCJmYW1pbHlfbmFtZSI6IlBlcnJ5IiwiZW1haWwiOiJvd2Fpbi5wZXJyeUBuZW9tLmNvbSJ9.cEaHPa5CladtjEs4tw4agUrcQpZ9HymLgEiAOHxs_Goe-iSPFHmAn7bfjKQp-yGBpsqTZSWU2SWDkWOSctO9QILNknh_sYdCZ0zT3xZIoX1TjSrYhkjtt7ga1QlW7MYmu8TtQYmYMzNZl7QJdvycM2lePLo2SaoiK9r7bEJucGiFYJNmFL9i5hGedpbW5CHe4s_E1cHcZVIZR1w-Rn_u6J7zZkPec7eOXcTfZnQkKxW-KhS1tpaF0ZhBLkueW3ht4Sh2mmNXuik7_LGhkqXP4diC8FBEMFlqaODVMKmfOXQjXdpimjicbqCnsUtbjM2P50WyycGqCF3yMvD8j_LlgQ' \
  -H 'Content-Type: application/json' \
  -d '          {
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
                }
            ],
            "finalisers": [
                {
                "finaliser": "save_dataframe",
                "input": "after_cast",
                "write_mode": "overwrite"
                }
            ],
            "preview": false
            }'