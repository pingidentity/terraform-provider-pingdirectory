package config

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	datasourceschema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
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
func AddCommonResourceSchema(s *schema.Schema, addNameAttribute bool) {
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
	s.Attributes["id"] = schema.StringAttribute{
		Description: "The ID of this resource.",
		Computed:    true,
		Required:    false,
		Optional:    false,
		PlanModifiers: []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		},
	}
	// If name is required (for instantiable config objects) then set it as Required and
	// require replace when changing.
	if addNameAttribute {
		s.Attributes["name"] = schema.StringAttribute{
			Description: "Name of this config object.",
			Required:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.RequiresReplace(),
			},
		}
	}
}

// Get schema elements common to all resources
func AddCommonDataSourceSchema(s *datasourceschema.Schema, addNameAttribute bool) {
	s.Attributes["id"] = datasourceschema.StringAttribute{
		Description: "The ID of this resource.",
		Computed:    true,
		Required:    false,
		Optional:    false,
	}
	// If name is required (for instantiable config objects) then set it as Required and
	// require replace when changing.
	if addNameAttribute {
		s.Attributes["name"] = datasourceschema.StringAttribute{
			Description: "Name of this config object.",
			Required:    true,
		}
	}
}

func SetAttributesToOptionalAndComputedAndRemoveDefaults(s *schema.Schema, exemptAttributes []string) {
	for key, attribute := range s.Attributes {
		// If more attribute types are used by this provider, this method will need to be updated
		if !internaltypes.StringSliceContains(exemptAttributes, key) {
			stringAttr, ok := attribute.(schema.StringAttribute)
			anyOk := ok
			if ok {
				stringAttr.Required = false
				stringAttr.Optional = true
				stringAttr.Computed = true
				stringAttr.Default = nil
				stringAttr.PlanModifiers = append(stringAttr.PlanModifiers, stringplanmodifier.UseStateForUnknown())
				s.Attributes[key] = stringAttr
				continue
			}
			setAttr, ok := attribute.(schema.SetAttribute)
			anyOk = ok || anyOk
			if ok {
				setAttr.Required = false
				setAttr.Optional = true
				setAttr.Computed = true
				setAttr.Default = nil
				setAttr.PlanModifiers = append(setAttr.PlanModifiers, setplanmodifier.UseStateForUnknown())
				s.Attributes[key] = setAttr
				continue
			}
			boolAttr, ok := attribute.(schema.BoolAttribute)
			anyOk = ok || anyOk
			if ok {
				boolAttr.Required = false
				boolAttr.Optional = true
				boolAttr.Computed = true
				boolAttr.Default = nil
				boolAttr.PlanModifiers = append(boolAttr.PlanModifiers, boolplanmodifier.UseStateForUnknown())
				s.Attributes[key] = boolAttr
				continue
			}
			intAttr, ok := attribute.(schema.Int64Attribute)
			anyOk = ok || anyOk
			if ok {
				intAttr.Required = false
				intAttr.Optional = true
				intAttr.Computed = true
				intAttr.Default = nil
				intAttr.PlanModifiers = append(intAttr.PlanModifiers, int64planmodifier.UseStateForUnknown())
				s.Attributes[key] = intAttr
				continue
			}
			floatAttr, ok := attribute.(schema.Float64Attribute)
			anyOk = ok || anyOk
			if ok {
				floatAttr.Required = false
				floatAttr.Optional = true
				floatAttr.Computed = true
				floatAttr.Default = nil
				floatAttr.PlanModifiers = append(floatAttr.PlanModifiers, float64planmodifier.UseStateForUnknown())
				s.Attributes[key] = floatAttr
				continue
			}
			if !anyOk {
				//lintignore:R009
				panic("No valid schema attribute type found when setting attributes to computed: " + key)
			}
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
