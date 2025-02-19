---
page_title: "cancom_dynamic_cloud_vpc_project Resource - terraform-provider-cancom"
subcategory: "Dynamic Cloud"
description: |-
  Manage Dynamic Cloud VPC Projects lifecycle
  This creates a Virtual Private Cloud (VPC) Project with the specified name and the optional comment. The parameter users can be used to specify which user should get access to the VPC Project.
  !> Changing the name or comment will force the VPC Project to be recreated, i.e. all resources in the VPC Project will be deleted.
---

# cancom_dynamic_cloud_vpc_project (Resource)

Manage Dynamic Cloud VPC Projects lifecycle

This creates a Virtual Private Cloud (VPC) Project with the specified name and the optional comment. The parameter `users` can be used to specify which user should get access to the VPC Project.

!> Changing the `name` or `comment` will force the VPC Project to be recreated, i.e. all resources in the VPC Project will be deleted.

## Example Usage

```terraform
resource "cancom_dynamic_cloud_vpc_project" "basic_usage_example" {
  name            = "test-cancom-terraform-provider"
  project_comment = "Test VPC Project created with the CANCOM Terraform provider"
  users           = ["crn:cancom::iam:user:example.name@example.domain"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The user defined name used to construct the OpenStack project name with the schema `tenant-name`.  
By changing this value, the old project will be deleted and a new project with the new name will be created.

!> Changing this value will delete all resources in the VPC Project.

### Optional

- `project_comment` (String) A comment to describe what this VPC Project is used for.  
By changing this value, the old project will be deleted and a new project will be created.

!> Changing this value will delete all resources in the VPC Project.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `users` (Set of String) The list of users with access to the VPC Project. The list may only contains CRNs of human iam users.

### Read-Only

- `created_by` (String) The CRN of the user who created the VPC Project.
- `creation_date` (String) The timestamp of the date when the VPC Project was created.
- `id` (String) The ID of this resource.
- `limits` (Map of Number) The resource limits currently configured for this VPC Project.
- `openstack_uuid` (String) The uuid of the OpenStack Project.
- `tenant` (String) The id of the tenant this VPC Project belongs to.

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
