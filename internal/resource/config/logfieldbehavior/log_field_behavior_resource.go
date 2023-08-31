package logfieldbehavior

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &logFieldBehaviorResource{}
	_ resource.ResourceWithConfigure   = &logFieldBehaviorResource{}
	_ resource.ResourceWithImportState = &logFieldBehaviorResource{}
	_ resource.Resource                = &defaultLogFieldBehaviorResource{}
	_ resource.ResourceWithConfigure   = &defaultLogFieldBehaviorResource{}
	_ resource.ResourceWithImportState = &defaultLogFieldBehaviorResource{}
)

// Create a Log Field Behavior resource
func NewLogFieldBehaviorResource() resource.Resource {
	return &logFieldBehaviorResource{}
}

func NewDefaultLogFieldBehaviorResource() resource.Resource {
	return &defaultLogFieldBehaviorResource{}
}

// logFieldBehaviorResource is the resource implementation.
type logFieldBehaviorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLogFieldBehaviorResource is the resource implementation.
type defaultLogFieldBehaviorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logFieldBehaviorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_field_behavior"
}

func (r *defaultLogFieldBehaviorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_field_behavior"
}

// Configure adds the provider configured client to the resource.
func (r *logFieldBehaviorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLogFieldBehaviorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logFieldBehaviorResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Notifications                    types.Set    `tfsdk:"notifications"`
	RequiredActions                  types.Set    `tfsdk:"required_actions"`
	Type                             types.String `tfsdk:"type"`
	PreserveField                    types.Set    `tfsdk:"preserve_field"`
	PreserveFieldName                types.Set    `tfsdk:"preserve_field_name"`
	OmitField                        types.Set    `tfsdk:"omit_field"`
	OmitFieldName                    types.Set    `tfsdk:"omit_field_name"`
	RedactEntireValueField           types.Set    `tfsdk:"redact_entire_value_field"`
	RedactEntireValueFieldName       types.Set    `tfsdk:"redact_entire_value_field_name"`
	RedactValueComponentsField       types.Set    `tfsdk:"redact_value_components_field"`
	RedactValueComponentsFieldName   types.Set    `tfsdk:"redact_value_components_field_name"`
	TokenizeEntireValueField         types.Set    `tfsdk:"tokenize_entire_value_field"`
	TokenizeEntireValueFieldName     types.Set    `tfsdk:"tokenize_entire_value_field_name"`
	TokenizeValueComponentsField     types.Set    `tfsdk:"tokenize_value_components_field"`
	TokenizeValueComponentsFieldName types.Set    `tfsdk:"tokenize_value_components_field_name"`
	Description                      types.String `tfsdk:"description"`
	DefaultBehavior                  types.String `tfsdk:"default_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *logFieldBehaviorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logFieldBehaviorSchema(ctx, req, resp, false)
}

func (r *defaultLogFieldBehaviorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logFieldBehaviorSchema(ctx, req, resp, true)
}

func logFieldBehaviorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Log Field Behavior.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Field Behavior resource. Options are ['text-access', 'json-formatted-access']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"text-access", "json-formatted-access"}...),
				},
			},
			"preserve_field": schema.SetAttribute{
				Description: "The log fields whose values should be logged with the intended value. The values for these fields will be preserved, although they may be sanitized for parsability or safety purposes (for example, to escape special characters in the value), and values that are too long may be truncated.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"preserve_field_name": schema.SetAttribute{
				Description: "The names of any custom fields whose values should be preserved. This should generally only be used for fields that are not available through the preserve-field property (for example, custom log fields defined in Server SDK extensions).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"omit_field": schema.SetAttribute{
				Description: "The log fields that should be omitted entirely from log messages. Neither the field name nor value will be included.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"omit_field_name": schema.SetAttribute{
				Description: "The names of any custom fields that should be omitted from log messages. This should generally only be used for fields that are not available through the omit-field property (for example, custom log fields defined in Server SDK extensions).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"redact_entire_value_field": schema.SetAttribute{
				Description: "The log fields whose values should be completely redacted in log messages. The field name will be included, but with a fixed value that does not reflect the actual value for the field.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"redact_entire_value_field_name": schema.SetAttribute{
				Description: "The names of any custom fields whose values should be completely redacted. This should generally only be used for fields that are not available through the redact-entire-value-field property (for example, custom log fields defined in Server SDK extensions).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"redact_value_components_field": schema.SetAttribute{
				Description: "The log fields whose values will include redacted components.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"redact_value_components_field_name": schema.SetAttribute{
				Description: "The names of any custom fields for which to redact components within the value. This should generally only be used for fields that are not available through the redact-value-components-field property (for example, custom log fields defined in Server SDK extensions).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"tokenize_entire_value_field": schema.SetAttribute{
				Description: "The log fields whose values should be completely tokenized in log messages. The field name will be included, but the value will be replaced with a token that does not reveal the actual value, but that is generated from the value.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"tokenize_entire_value_field_name": schema.SetAttribute{
				Description: "The names of any custom fields whose values should be completely tokenized. This should generally only be used for fields that are not available through the tokenize-entire-value-field property (for example, custom log fields defined in Server SDK extensions).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"tokenize_value_components_field": schema.SetAttribute{
				Description: "The log fields whose values will include tokenized components.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"tokenize_value_components_field_name": schema.SetAttribute{
				Description: "The names of any custom fields for which to tokenize components within the value. This should generally only be used for fields that are not available through the tokenize-value-components-field property (for example, custom log fields defined in Server SDK extensions).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Behavior",
				Optional:    true,
			},
			"default_behavior": schema.StringAttribute{
				Description: "The default behavior that the server should exhibit for fields for which no explicit behavior is defined. If no default behavior is defined, the server will fall back to using the default behavior configured for the syntax used for each log field.",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for text-access log-field-behavior
func addOptionalTextAccessLogFieldBehaviorFields(ctx context.Context, addRequest *client.AddTextAccessLogFieldBehaviorRequest, plan logFieldBehaviorResourceModel) error {
	if internaltypes.IsDefined(plan.PreserveField) {
		var slice []string
		plan.PreserveField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorTextAccessPreserveFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorTextAccessPreserveFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PreserveField = enumSlice
	}
	if internaltypes.IsDefined(plan.PreserveFieldName) {
		var slice []string
		plan.PreserveFieldName.ElementsAs(ctx, &slice, false)
		addRequest.PreserveFieldName = slice
	}
	if internaltypes.IsDefined(plan.OmitField) {
		var slice []string
		plan.OmitField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorTextAccessOmitFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorTextAccessOmitFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OmitField = enumSlice
	}
	if internaltypes.IsDefined(plan.OmitFieldName) {
		var slice []string
		plan.OmitFieldName.ElementsAs(ctx, &slice, false)
		addRequest.OmitFieldName = slice
	}
	if internaltypes.IsDefined(plan.RedactEntireValueField) {
		var slice []string
		plan.RedactEntireValueField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorTextAccessRedactEntireValueFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorTextAccessRedactEntireValueFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.RedactEntireValueField = enumSlice
	}
	if internaltypes.IsDefined(plan.RedactEntireValueFieldName) {
		var slice []string
		plan.RedactEntireValueFieldName.ElementsAs(ctx, &slice, false)
		addRequest.RedactEntireValueFieldName = slice
	}
	if internaltypes.IsDefined(plan.RedactValueComponentsField) {
		var slice []string
		plan.RedactValueComponentsField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorTextAccessRedactValueComponentsFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorTextAccessRedactValueComponentsFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.RedactValueComponentsField = enumSlice
	}
	if internaltypes.IsDefined(plan.RedactValueComponentsFieldName) {
		var slice []string
		plan.RedactValueComponentsFieldName.ElementsAs(ctx, &slice, false)
		addRequest.RedactValueComponentsFieldName = slice
	}
	if internaltypes.IsDefined(plan.TokenizeEntireValueField) {
		var slice []string
		plan.TokenizeEntireValueField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorTextAccessTokenizeEntireValueFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorTextAccessTokenizeEntireValueFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.TokenizeEntireValueField = enumSlice
	}
	if internaltypes.IsDefined(plan.TokenizeEntireValueFieldName) {
		var slice []string
		plan.TokenizeEntireValueFieldName.ElementsAs(ctx, &slice, false)
		addRequest.TokenizeEntireValueFieldName = slice
	}
	if internaltypes.IsDefined(plan.TokenizeValueComponentsField) {
		var slice []string
		plan.TokenizeValueComponentsField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorTextAccessTokenizeValueComponentsFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorTextAccessTokenizeValueComponentsFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.TokenizeValueComponentsField = enumSlice
	}
	if internaltypes.IsDefined(plan.TokenizeValueComponentsFieldName) {
		var slice []string
		plan.TokenizeValueComponentsFieldName.ElementsAs(ctx, &slice, false)
		addRequest.TokenizeValueComponentsFieldName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultBehavior) {
		defaultBehavior, err := client.NewEnumlogFieldBehaviorDefaultBehaviorPropFromValue(plan.DefaultBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.DefaultBehavior = defaultBehavior
	}
	return nil
}

// Add optional fields to create request for json-formatted-access log-field-behavior
func addOptionalJsonFormattedAccessLogFieldBehaviorFields(ctx context.Context, addRequest *client.AddJsonFormattedAccessLogFieldBehaviorRequest, plan logFieldBehaviorResourceModel) error {
	if internaltypes.IsDefined(plan.PreserveField) {
		var slice []string
		plan.PreserveField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorJsonFormattedAccessPreserveFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorJsonFormattedAccessPreserveFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PreserveField = enumSlice
	}
	if internaltypes.IsDefined(plan.PreserveFieldName) {
		var slice []string
		plan.PreserveFieldName.ElementsAs(ctx, &slice, false)
		addRequest.PreserveFieldName = slice
	}
	if internaltypes.IsDefined(plan.OmitField) {
		var slice []string
		plan.OmitField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorJsonFormattedAccessOmitFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorJsonFormattedAccessOmitFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.OmitField = enumSlice
	}
	if internaltypes.IsDefined(plan.OmitFieldName) {
		var slice []string
		plan.OmitFieldName.ElementsAs(ctx, &slice, false)
		addRequest.OmitFieldName = slice
	}
	if internaltypes.IsDefined(plan.RedactEntireValueField) {
		var slice []string
		plan.RedactEntireValueField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorJsonFormattedAccessRedactEntireValueFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorJsonFormattedAccessRedactEntireValueFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.RedactEntireValueField = enumSlice
	}
	if internaltypes.IsDefined(plan.RedactEntireValueFieldName) {
		var slice []string
		plan.RedactEntireValueFieldName.ElementsAs(ctx, &slice, false)
		addRequest.RedactEntireValueFieldName = slice
	}
	if internaltypes.IsDefined(plan.RedactValueComponentsField) {
		var slice []string
		plan.RedactValueComponentsField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorJsonFormattedAccessRedactValueComponentsFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorJsonFormattedAccessRedactValueComponentsFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.RedactValueComponentsField = enumSlice
	}
	if internaltypes.IsDefined(plan.RedactValueComponentsFieldName) {
		var slice []string
		plan.RedactValueComponentsFieldName.ElementsAs(ctx, &slice, false)
		addRequest.RedactValueComponentsFieldName = slice
	}
	if internaltypes.IsDefined(plan.TokenizeEntireValueField) {
		var slice []string
		plan.TokenizeEntireValueField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorJsonFormattedAccessTokenizeEntireValueFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorJsonFormattedAccessTokenizeEntireValueFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.TokenizeEntireValueField = enumSlice
	}
	if internaltypes.IsDefined(plan.TokenizeEntireValueFieldName) {
		var slice []string
		plan.TokenizeEntireValueFieldName.ElementsAs(ctx, &slice, false)
		addRequest.TokenizeEntireValueFieldName = slice
	}
	if internaltypes.IsDefined(plan.TokenizeValueComponentsField) {
		var slice []string
		plan.TokenizeValueComponentsField.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogFieldBehaviorJsonFormattedAccessTokenizeValueComponentsFieldProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogFieldBehaviorJsonFormattedAccessTokenizeValueComponentsFieldPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.TokenizeValueComponentsField = enumSlice
	}
	if internaltypes.IsDefined(plan.TokenizeValueComponentsFieldName) {
		var slice []string
		plan.TokenizeValueComponentsFieldName.ElementsAs(ctx, &slice, false)
		addRequest.TokenizeValueComponentsFieldName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultBehavior) {
		defaultBehavior, err := client.NewEnumlogFieldBehaviorDefaultBehaviorPropFromValue(plan.DefaultBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.DefaultBehavior = defaultBehavior
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *logFieldBehaviorResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.DefaultBehavior.IsUnknown() || model.DefaultBehavior.IsNull() {
		model.DefaultBehavior = types.StringValue("")
	}
}

// Read a TextAccessLogFieldBehaviorResponse object into the model struct
func readTextAccessLogFieldBehaviorResponse(ctx context.Context, r *client.TextAccessLogFieldBehaviorResponse, state *logFieldBehaviorResourceModel, expectedValues *logFieldBehaviorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("text-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PreserveField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessPreserveFieldProp(r.PreserveField))
	state.PreserveFieldName = internaltypes.GetStringSet(r.PreserveFieldName)
	state.OmitField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessOmitFieldProp(r.OmitField))
	state.OmitFieldName = internaltypes.GetStringSet(r.OmitFieldName)
	state.RedactEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessRedactEntireValueFieldProp(r.RedactEntireValueField))
	state.RedactEntireValueFieldName = internaltypes.GetStringSet(r.RedactEntireValueFieldName)
	state.RedactValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessRedactValueComponentsFieldProp(r.RedactValueComponentsField))
	state.RedactValueComponentsFieldName = internaltypes.GetStringSet(r.RedactValueComponentsFieldName)
	state.TokenizeEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessTokenizeEntireValueFieldProp(r.TokenizeEntireValueField))
	state.TokenizeEntireValueFieldName = internaltypes.GetStringSet(r.TokenizeEntireValueFieldName)
	state.TokenizeValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessTokenizeValueComponentsFieldProp(r.TokenizeValueComponentsField))
	state.TokenizeValueComponentsFieldName = internaltypes.GetStringSet(r.TokenizeValueComponentsFieldName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldBehaviorDefaultBehaviorProp(r.DefaultBehavior), internaltypes.IsEmptyString(expectedValues.DefaultBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a JsonFormattedAccessLogFieldBehaviorResponse object into the model struct
func readJsonFormattedAccessLogFieldBehaviorResponse(ctx context.Context, r *client.JsonFormattedAccessLogFieldBehaviorResponse, state *logFieldBehaviorResourceModel, expectedValues *logFieldBehaviorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-formatted-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PreserveField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessPreserveFieldProp(r.PreserveField))
	state.PreserveFieldName = internaltypes.GetStringSet(r.PreserveFieldName)
	state.OmitField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessOmitFieldProp(r.OmitField))
	state.OmitFieldName = internaltypes.GetStringSet(r.OmitFieldName)
	state.RedactEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessRedactEntireValueFieldProp(r.RedactEntireValueField))
	state.RedactEntireValueFieldName = internaltypes.GetStringSet(r.RedactEntireValueFieldName)
	state.RedactValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessRedactValueComponentsFieldProp(r.RedactValueComponentsField))
	state.RedactValueComponentsFieldName = internaltypes.GetStringSet(r.RedactValueComponentsFieldName)
	state.TokenizeEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessTokenizeEntireValueFieldProp(r.TokenizeEntireValueField))
	state.TokenizeEntireValueFieldName = internaltypes.GetStringSet(r.TokenizeEntireValueFieldName)
	state.TokenizeValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessTokenizeValueComponentsFieldProp(r.TokenizeValueComponentsField))
	state.TokenizeValueComponentsFieldName = internaltypes.GetStringSet(r.TokenizeValueComponentsFieldName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldBehaviorDefaultBehaviorProp(r.DefaultBehavior), internaltypes.IsEmptyString(expectedValues.DefaultBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLogFieldBehaviorOperations(plan logFieldBehaviorResourceModel, state logFieldBehaviorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PreserveField, state.PreserveField, "preserve-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PreserveFieldName, state.PreserveFieldName, "preserve-field-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OmitField, state.OmitField, "omit-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OmitFieldName, state.OmitFieldName, "omit-field-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RedactEntireValueField, state.RedactEntireValueField, "redact-entire-value-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RedactEntireValueFieldName, state.RedactEntireValueFieldName, "redact-entire-value-field-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RedactValueComponentsField, state.RedactValueComponentsField, "redact-value-components-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RedactValueComponentsFieldName, state.RedactValueComponentsFieldName, "redact-value-components-field-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TokenizeEntireValueField, state.TokenizeEntireValueField, "tokenize-entire-value-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TokenizeEntireValueFieldName, state.TokenizeEntireValueFieldName, "tokenize-entire-value-field-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TokenizeValueComponentsField, state.TokenizeValueComponentsField, "tokenize-value-components-field")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.TokenizeValueComponentsFieldName, state.TokenizeValueComponentsFieldName, "tokenize-value-components-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultBehavior, state.DefaultBehavior, "default-behavior")
	return ops
}

// Create a text-access log-field-behavior
func (r *logFieldBehaviorResource) CreateTextAccessLogFieldBehavior(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFieldBehaviorResourceModel) (*logFieldBehaviorResourceModel, error) {
	addRequest := client.NewAddTextAccessLogFieldBehaviorRequest(plan.Name.ValueString(),
		[]client.EnumtextAccessLogFieldBehaviorSchemaUrn{client.ENUMTEXTACCESSLOGFIELDBEHAVIORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FIELD_BEHAVIORTEXT_ACCESS})
	err := addOptionalTextAccessLogFieldBehaviorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Field Behavior", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFieldBehaviorApi.AddLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFieldBehaviorRequest(
		client.AddTextAccessLogFieldBehaviorRequestAsAddLogFieldBehaviorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFieldBehaviorApi.AddLogFieldBehaviorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Field Behavior", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFieldBehaviorResourceModel
	readTextAccessLogFieldBehaviorResponse(ctx, addResponse.TextAccessLogFieldBehaviorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a json-formatted-access log-field-behavior
func (r *logFieldBehaviorResource) CreateJsonFormattedAccessLogFieldBehavior(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logFieldBehaviorResourceModel) (*logFieldBehaviorResourceModel, error) {
	addRequest := client.NewAddJsonFormattedAccessLogFieldBehaviorRequest(plan.Name.ValueString(),
		[]client.EnumjsonFormattedAccessLogFieldBehaviorSchemaUrn{client.ENUMJSONFORMATTEDACCESSLOGFIELDBEHAVIORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FIELD_BEHAVIORJSON_FORMATTED_ACCESS})
	err := addOptionalJsonFormattedAccessLogFieldBehaviorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Field Behavior", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFieldBehaviorApi.AddLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFieldBehaviorRequest(
		client.AddJsonFormattedAccessLogFieldBehaviorRequestAsAddLogFieldBehaviorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFieldBehaviorApi.AddLogFieldBehaviorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Field Behavior", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logFieldBehaviorResourceModel
	readJsonFormattedAccessLogFieldBehaviorResponse(ctx, addResponse.JsonFormattedAccessLogFieldBehaviorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *logFieldBehaviorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFieldBehaviorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *logFieldBehaviorResourceModel
	var err error
	if plan.Type.ValueString() == "text-access" {
		state, err = r.CreateTextAccessLogFieldBehavior(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "json-formatted-access" {
		state, err = r.CreateJsonFormattedAccessLogFieldBehavior(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultLogFieldBehaviorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logFieldBehaviorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldBehaviorApi.GetLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Behavior", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logFieldBehaviorResourceModel
	if readResponse.TextAccessLogFieldBehaviorResponse != nil {
		readTextAccessLogFieldBehaviorResponse(ctx, readResponse.TextAccessLogFieldBehaviorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JsonFormattedAccessLogFieldBehaviorResponse != nil {
		readJsonFormattedAccessLogFieldBehaviorResponse(ctx, readResponse.JsonFormattedAccessLogFieldBehaviorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFieldBehaviorApi.UpdateLogFieldBehavior(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogFieldBehaviorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFieldBehaviorApi.UpdateLogFieldBehaviorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Field Behavior", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.TextAccessLogFieldBehaviorResponse != nil {
			readTextAccessLogFieldBehaviorResponse(ctx, updateResponse.TextAccessLogFieldBehaviorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JsonFormattedAccessLogFieldBehaviorResponse != nil {
			readJsonFormattedAccessLogFieldBehaviorResponse(ctx, updateResponse.JsonFormattedAccessLogFieldBehaviorResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *logFieldBehaviorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogFieldBehavior(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLogFieldBehaviorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogFieldBehavior(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLogFieldBehavior(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state logFieldBehaviorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogFieldBehaviorApi.GetLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Behavior", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Behavior", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TextAccessLogFieldBehaviorResponse != nil {
		readTextAccessLogFieldBehaviorResponse(ctx, readResponse.TextAccessLogFieldBehaviorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JsonFormattedAccessLogFieldBehaviorResponse != nil {
		readJsonFormattedAccessLogFieldBehaviorResponse(ctx, readResponse.JsonFormattedAccessLogFieldBehaviorResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *logFieldBehaviorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogFieldBehavior(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogFieldBehaviorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogFieldBehavior(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLogFieldBehavior(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan logFieldBehaviorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logFieldBehaviorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogFieldBehaviorApi.UpdateLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogFieldBehaviorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogFieldBehaviorApi.UpdateLogFieldBehaviorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Field Behavior", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.TextAccessLogFieldBehaviorResponse != nil {
			readTextAccessLogFieldBehaviorResponse(ctx, updateResponse.TextAccessLogFieldBehaviorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JsonFormattedAccessLogFieldBehaviorResponse != nil {
			readJsonFormattedAccessLogFieldBehaviorResponse(ctx, updateResponse.JsonFormattedAccessLogFieldBehaviorResponse, &state, &plan, &resp.Diagnostics)
		}
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultLogFieldBehaviorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logFieldBehaviorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state logFieldBehaviorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogFieldBehaviorApi.DeleteLogFieldBehaviorExecute(r.apiClient.LogFieldBehaviorApi.DeleteLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Log Field Behavior", err, httpResp)
		return
	}
}

func (r *logFieldBehaviorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogFieldBehavior(ctx, req, resp)
}

func (r *defaultLogFieldBehaviorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogFieldBehavior(ctx, req, resp)
}

func importLogFieldBehavior(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
