package config

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Get attrtype map for the requiredActions returned by the config API
func getRequiredActionsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"property": types.StringType,
		"type":     types.StringType,
		"synopsis": types.StringType,
	}
}

// Get the requiredActions ObjectType definition
func GetRequiredActionsObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: getRequiredActionsAttrTypes(),
	}
}

// Get schema elements common to all resources
func AddCommonSchema(s *schema.Schema, idRequired bool) {
	s.Attributes["last_updated"] = schema.StringAttribute{
		Description: "Timestamp of the last Terraform update of this resource.",
		Computed:    true,
		Required:    false,
		Optional:    false,
	}
	s.Attributes["notifications"] = schema.SetAttribute{
		Description: "Notifications returned by the PingDirectory Configuration API.",
		ElementType: types.StringType,
		Computed:    true,
		Required:    false,
		Optional:    false,
	}
	s.Attributes["required_actions"] = schema.SetAttribute{
		Description: "Required actions returned by the PingDirectory Configuration API.",
		ElementType: GetRequiredActionsObjectType(),
		Computed:    true,
		Required:    false,
		Optional:    false,
	}
	// If ID is required (for instantiable config objects) then set it as Required and
	// require replace when changing. Otherwise, mark it as Computed.
	if idRequired {
		s.Attributes["id"] = schema.StringAttribute{
			Description: "Name of this object.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		}
	} else {
		s.Attributes["id"] = schema.StringAttribute{
			Description: "Placeholder name of this object required by Terraform.",
			Computed:    true,
		}
	}
}

// Get the set of required actions from the configuration messages returned by the config API
func GetRequiredActionsSet(messages client.MetaUrnPingidentitySchemasConfigurationMessages20) (types.Set, diag.Diagnostics) {
	setValues := make([]attr.Value, len(messages.RequiredActions))
	for i := 0; i < len(messages.RequiredActions); i++ {
		property := types.StringNull()
		if messages.RequiredActions[i].Property != nil {
			property = types.StringValue(*messages.RequiredActions[i].Property)
		}
		setValues[i], _ = types.ObjectValue(getRequiredActionsAttrTypes(), map[string]attr.Value{
			"property": property,
			"type":     types.StringValue(messages.RequiredActions[i].Type),
			"synopsis": types.StringValue(messages.RequiredActions[i].Synopsis),
		})
	}
	return types.SetValue(GetRequiredActionsObjectType(), setValues)
}
