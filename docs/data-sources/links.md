---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "neos_links Data Source - terraform-provider-neos"
subcategory: ""
description: |-
  
---

# neos_links (Data Source)





<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `links` (Attributes List) (see [below for nested schema](#nestedatt--links))

<a id="nestedatt--links"></a>
### Nested Schema for `links`

Read-Only:

- `child` (Attributes) (see [below for nested schema](#nestedatt--links--child))
- `parent` (Attributes) (see [below for nested schema](#nestedatt--links--parent))
- `tmp` (String)

<a id="nestedatt--links--child"></a>
### Nested Schema for `links.child`

Read-Only:

- `created_at` (String)
- `description` (String)
- `entity_type` (String)
- `identifier` (String)
- `is_system` (Boolean)
- `label` (String)
- `name` (String)
- `output_type` (String)
- `owner` (String)
- `urn` (String)


<a id="nestedatt--links--parent"></a>
### Nested Schema for `links.parent`

Read-Only:

- `created_at` (String)
- `description` (String)
- `entity_type` (String)
- `identifier` (String)
- `is_system` (Boolean)
- `label` (String)
- `name` (String)
- `output_type` (String)
- `owner` (String)
- `urn` (String)
