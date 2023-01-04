package config

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdata-config-api-go-client"
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
func AddCommonSchema(schema *tfsdk.Schema) {
	schema.Attributes["last_updated"] = tfsdk.Attribute{
		Description: "Timestamp of the last Terraform update of this resource.",
		Type:        types.StringType,
		Computed:    true,
		Required:    false,
		Optional:    false,
	}
	schema.Attributes["notifications"] = tfsdk.Attribute{
		Description: "Notifications returned by the PingDirectory Configuration API.",
		Type: types.SetType{
			ElemType: types.StringType,
		},
		Computed: true,
		Required: false,
		Optional: false,
	}
	schema.Attributes["required_actions"] = tfsdk.Attribute{
		Description: "Required actions returned by the PingDirectory Configuration API.",
		Type: types.SetType{
			ElemType: GetRequiredActionsObjectType(),
		},
		Computed: true,
		Required: false,
		Optional: false,
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
