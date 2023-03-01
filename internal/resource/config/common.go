package config

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
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

func SetOptionalAttributesToComputed(s *schema.Schema) {
	for key, attribute := range s.Attributes {
		// If more attribute types are used by this provider, this method will need to be updated
		if attribute.IsOptional() {
			stringAttr, ok := attribute.(schema.StringAttribute)
			if ok {
				stringAttr.Computed = true
				continue
			}
			setAttr, ok := attribute.(schema.SetAttribute)
			if ok {
				setAttr.Computed = true
				continue
			}
			boolAttr, ok := attribute.(schema.BoolAttribute)
			if ok {
				boolAttr.Computed = true
				continue
			}
			intAttr, ok := attribute.(schema.Int64Attribute)
			if ok {
				intAttr.Computed = true
				continue
			}
			floatAttr, ok := attribute.(schema.Float64Attribute)
			if ok {
				floatAttr.Computed = true
				continue
			}
			panic("No valid schema attribute type found when setting attributes to computed: " + key)
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
