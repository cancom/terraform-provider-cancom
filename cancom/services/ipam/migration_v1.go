package ipam

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var resourceSchemaInstanceV0 = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"name_tag": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
		},
		"managed_by": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: false,
			Optional: true,
		},
		"release_wait_time": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

func instanceUpgradeV0() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    resourceSchemaInstanceV0.CoreConfigSchema().ImpliedType(),
		Version: 0,
		Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
			if val, ok := rawState["updated_at"].(string); ok {
				i, err := time.Parse("2006-01-02T15:04:05", val)
				if err != nil {
					return rawState, fmt.Errorf("error converting updated_at to int: %w", err)
				}
				rawState["updated_at"] = int(i.Unix())
			}

			if val, ok := rawState["created_at"].(string); ok {
				i, err := time.Parse("2006-01-02T15:04:05", val)
				if err != nil {
					return rawState, fmt.Errorf("error converting created_at to int: %w", err)
				}
				rawState["created_at"] = int(i.Unix())
			}

			if val, ok := rawState["release_wait_time"].(string); ok {
				i, err := strconv.Atoi(val)
				if err != nil {
					return rawState, fmt.Errorf("failed to convert release_wait_time from string to int: %w", err)
				}
				rawState["release_wait_time"] = i
			}

			return rawState, nil
		},
	}
}

var resourceSchemaNetworkV0 = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"supernet_id": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
			ForceNew: true,
		},
		"name_tag": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: false,
			Optional: true,
		},
		"request": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
			ForceNew: true,
		},
		"host_assign": {
			Type:     schema.TypeBool,
			Computed: false,
			Optional: true,
			Default:  true,
		},
		"prefix_str": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
		},
	},
}

func networkUpgradeV0() schema.StateUpgrader {
	return schema.StateUpgrader{
		Type:    resourceSchemaNetworkV0.CoreConfigSchema().ImpliedType(),
		Version: 0,
		Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
			if val, ok := rawState["updated_at"].(string); ok {
				if val != "" {
					i, err := time.Parse("2006-01-02T15:04:05", val)
					if err != nil {
						return rawState, fmt.Errorf("error converting updated_at to int: %w", err)
					}
					rawState["updated_at"] = int(i.Unix())
				} else {
					rawState["updated_at"] = nil
				}
			}

			if val, ok := rawState["created_at"].(string); ok {
				if val != "" {
					i, err := time.Parse("2006-01-02T15:04:05", val)
					if err != nil {
						return rawState, fmt.Errorf("error converting created_at to int: %w", err)
					}
					rawState["created_at"] = int(i.Unix())
				} else {
					rawState["created_at"] = nil
				}
			}

			return rawState, nil
		},
	}
}

var resourceSchemaSupernetV0 = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"id": {
			Type:     schema.TypeString,
			Computed: true,
		},
		"instance_id": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
		},
		"name_tag": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
		},
		"description": {
			Type:     schema.TypeString,
			Computed: false,
			Optional: true,
		},
		"supernet_cidr": {
			Type:     schema.TypeString,
			Computed: false,
			Required: true,
			ForceNew: true,
		},
		"created_at": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
		"updated_at": {
			Type:     schema.TypeString,
			Computed: true,
			Optional: true,
		},
	},
}

func supernetUpgradeV0() schema.StateUpgrader {
	return schema.StateUpgrader{

		Type:    resourceSchemaSupernetV0.CoreConfigSchema().ImpliedType(),
		Version: 0,
		Upgrade: func(ctx context.Context, rawState map[string]interface{}, meta interface{}) (map[string]interface{}, error) {
			if val, ok := rawState["updated_at"].(string); ok {
				if val != "" {
					i, err := time.Parse("2006-01-02T15:04:05", val)
					if err != nil {
						return rawState, fmt.Errorf("error converting updated_at to int: %w", err)
					}
					rawState["updated_at"] = int(i.Unix())
				} else {

					rawState["updated_at"] = 0
				}

			}

			if val, ok := rawState["created_at"].(string); ok {
				if val != "" {
					i, err := time.Parse("2006-01-02T15:04:05", val)
					if err != nil {
						return rawState, fmt.Errorf("error converting created_at to int: %w", err)
					}
					rawState["created_at"] = int(i.Unix())
				} else {
					rawState["created_at"] = nil
				}
			}

			return rawState, nil
		},
	}
}
