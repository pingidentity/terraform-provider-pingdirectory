package attributesyntax

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &attributeSyntaxResource{}
	_ resource.ResourceWithConfigure   = &attributeSyntaxResource{}
	_ resource.ResourceWithImportState = &attributeSyntaxResource{}
)

// Create a Attribute Syntax resource
func NewAttributeSyntaxResource() resource.Resource {
	return &attributeSyntaxResource{}
}

// attributeSyntaxResource is the resource implementation.
type attributeSyntaxResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *attributeSyntaxResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_attribute_syntax"
}

// Configure adds the provider configured client to the resource.
func (r *attributeSyntaxResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type attributeSyntaxResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
}

type defaultAttributeSyntaxResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	Type                           types.String `tfsdk:"type"`
	EnableCompaction               types.Bool   `tfsdk:"enable_compaction"`
	IncludeAttributeInCompaction   types.Set    `tfsdk:"include_attribute_in_compaction"`
	ExcludeAttributeFromCompaction types.Set    `tfsdk:"exclude_attribute_from_compaction"`
	StrictFormat                   types.Bool   `tfsdk:"strict_format"`
	AllowZeroLengthValues          types.Bool   `tfsdk:"allow_zero_length_values"`
	StripSyntaxMinUpperBound       types.Bool   `tfsdk:"strip_syntax_min_upper_bound"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	RequireBinaryTransfer          types.Bool   `tfsdk:"require_binary_transfer"`
}

// GetSchema defines the schema for the resource.
func (r *attributeSyntaxResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Attribute Syntax.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Attribute Syntax resource. Options are ['attribute-type-description', 'directory-string', 'telephone-number', 'distinguished-name', 'generalized-time', 'integer', 'uuid', 'generic', 'json-object', 'user-password', 'boolean', 'hex-string', 'bit-string', 'ldap-url', 'name-and-optional-uid']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{}...),
				},
			},
		},
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *attributeSyntaxResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model defaultAttributeSyntaxResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.StrictFormat) && model.Type.ValueString() != "telephone-number" && model.Type.ValueString() != "ldap-url" {
		resp.Diagnostics.AddError("Attribute 'strict_format' not supported by pingdirectory_attribute_syntax resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'strict_format', the 'type' attribute must be one of ['telephone-number', 'ldap-url']")
	}
	if internaltypes.IsDefined(model.ExcludeAttributeFromCompaction) && model.Type.ValueString() != "user-password" && model.Type.ValueString() != "boolean" && model.Type.ValueString() != "bit-string" && model.Type.ValueString() != "distinguished-name" && model.Type.ValueString() != "generalized-time" && model.Type.ValueString() != "integer" && model.Type.ValueString() != "uuid" && model.Type.ValueString() != "name-and-optional-uid" {
		resp.Diagnostics.AddError("Attribute 'exclude_attribute_from_compaction' not supported by pingdirectory_attribute_syntax resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'exclude_attribute_from_compaction', the 'type' attribute must be one of ['user-password', 'boolean', 'bit-string', 'distinguished-name', 'generalized-time', 'integer', 'uuid', 'name-and-optional-uid']")
	}
	if internaltypes.IsDefined(model.EnableCompaction) && model.Type.ValueString() != "user-password" && model.Type.ValueString() != "boolean" && model.Type.ValueString() != "bit-string" && model.Type.ValueString() != "distinguished-name" && model.Type.ValueString() != "generalized-time" && model.Type.ValueString() != "integer" && model.Type.ValueString() != "uuid" && model.Type.ValueString() != "name-and-optional-uid" {
		resp.Diagnostics.AddError("Attribute 'enable_compaction' not supported by pingdirectory_attribute_syntax resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'enable_compaction', the 'type' attribute must be one of ['user-password', 'boolean', 'bit-string', 'distinguished-name', 'generalized-time', 'integer', 'uuid', 'name-and-optional-uid']")
	}
	if internaltypes.IsDefined(model.IncludeAttributeInCompaction) && model.Type.ValueString() != "user-password" && model.Type.ValueString() != "boolean" && model.Type.ValueString() != "bit-string" && model.Type.ValueString() != "distinguished-name" && model.Type.ValueString() != "generalized-time" && model.Type.ValueString() != "integer" && model.Type.ValueString() != "uuid" && model.Type.ValueString() != "name-and-optional-uid" {
		resp.Diagnostics.AddError("Attribute 'include_attribute_in_compaction' not supported by pingdirectory_attribute_syntax resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_attribute_in_compaction', the 'type' attribute must be one of ['user-password', 'boolean', 'bit-string', 'distinguished-name', 'generalized-time', 'integer', 'uuid', 'name-and-optional-uid']")
	}
	if internaltypes.IsDefined(model.StripSyntaxMinUpperBound) && model.Type.ValueString() != "attribute-type-description" {
		resp.Diagnostics.AddError("Attribute 'strip_syntax_min_upper_bound' not supported by pingdirectory_attribute_syntax resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'strip_syntax_min_upper_bound', the 'type' attribute must be one of ['attribute-type-description']")
	}
	if internaltypes.IsDefined(model.AllowZeroLengthValues) && model.Type.ValueString() != "directory-string" {
		resp.Diagnostics.AddError("Attribute 'allow_zero_length_values' not supported by pingdirectory_attribute_syntax resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allow_zero_length_values', the 'type' attribute must be one of ['directory-string']")
	}
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateAttributeSyntaxNilSetsDefault(ctx context.Context, model *defaultAttributeSyntaxResourceModel) {
	if model.ExcludeAttributeFromCompaction.ElementType(ctx) == nil {
		model.ExcludeAttributeFromCompaction = types.SetNull(types.StringType)
	}
	if model.IncludeAttributeInCompaction.ElementType(ctx) == nil {
		model.IncludeAttributeInCompaction = types.SetNull(types.StringType)
	}
}

// Read a AttributeTypeDescriptionAttributeSyntaxResponse object into the model struct
func readAttributeTypeDescriptionAttributeSyntaxResponseDefault(ctx context.Context, r *client.AttributeTypeDescriptionAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-type-description")
	state.Id = types.StringValue(r.Id)
	state.StripSyntaxMinUpperBound = internaltypes.BoolTypeOrNil(r.StripSyntaxMinUpperBound)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a DirectoryStringAttributeSyntaxResponse object into the model struct
func readDirectoryStringAttributeSyntaxResponseDefault(ctx context.Context, r *client.DirectoryStringAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory-string")
	state.Id = types.StringValue(r.Id)
	state.AllowZeroLengthValues = internaltypes.BoolTypeOrNil(r.AllowZeroLengthValues)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a TelephoneNumberAttributeSyntaxResponse object into the model struct
func readTelephoneNumberAttributeSyntaxResponseDefault(ctx context.Context, r *client.TelephoneNumberAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("telephone-number")
	state.Id = types.StringValue(r.Id)
	state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a DistinguishedNameAttributeSyntaxResponse object into the model struct
func readDistinguishedNameAttributeSyntaxResponseDefault(ctx context.Context, r *client.DistinguishedNameAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("distinguished-name")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a GeneralizedTimeAttributeSyntaxResponse object into the model struct
func readGeneralizedTimeAttributeSyntaxResponseDefault(ctx context.Context, r *client.GeneralizedTimeAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generalized-time")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a IntegerAttributeSyntaxResponse object into the model struct
func readIntegerAttributeSyntaxResponseDefault(ctx context.Context, r *client.IntegerAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("integer")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a UuidAttributeSyntaxResponse object into the model struct
func readUuidAttributeSyntaxResponseDefault(ctx context.Context, r *client.UuidAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("uuid")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a GenericAttributeSyntaxResponse object into the model struct
func readGenericAttributeSyntaxResponseDefault(ctx context.Context, r *client.GenericAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a JsonObjectAttributeSyntaxResponse object into the model struct
func readJsonObjectAttributeSyntaxResponseDefault(ctx context.Context, r *client.JsonObjectAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-object")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a UserPasswordAttributeSyntaxResponse object into the model struct
func readUserPasswordAttributeSyntaxResponseDefault(ctx context.Context, r *client.UserPasswordAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-password")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a BooleanAttributeSyntaxResponse object into the model struct
func readBooleanAttributeSyntaxResponseDefault(ctx context.Context, r *client.BooleanAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("boolean")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a HexStringAttributeSyntaxResponse object into the model struct
func readHexStringAttributeSyntaxResponseDefault(ctx context.Context, r *client.HexStringAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("hex-string")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a BitStringAttributeSyntaxResponse object into the model struct
func readBitStringAttributeSyntaxResponseDefault(ctx context.Context, r *client.BitStringAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("bit-string")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a LdapUrlAttributeSyntaxResponse object into the model struct
func readLdapUrlAttributeSyntaxResponseDefault(ctx context.Context, r *client.LdapUrlAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-url")
	state.Id = types.StringValue(r.Id)
	state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Read a NameAndOptionalUidAttributeSyntaxResponse object into the model struct
func readNameAndOptionalUidAttributeSyntaxResponseDefault(ctx context.Context, r *client.NameAndOptionalUidAttributeSyntaxResponse, state *defaultAttributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("name-and-optional-uid")
	state.Id = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxNilSetsDefault(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createAttributeSyntaxOperations(plan attributeSyntaxResourceModel, state attributeSyntaxResourceModel) []client.Operation {
	var ops []client.Operation
	return ops
}

// Create any update operations necessary to make the state match the plan
func createAttributeSyntaxOperationsDefault(plan defaultAttributeSyntaxResourceModel, state defaultAttributeSyntaxResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableCompaction, state.EnableCompaction, "enable-compaction")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeAttributeInCompaction, state.IncludeAttributeInCompaction, "include-attribute-in-compaction")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeAttributeFromCompaction, state.ExcludeAttributeFromCompaction, "exclude-attribute-from-compaction")
	operations.AddBoolOperationIfNecessary(&ops, plan.StrictFormat, state.StrictFormat, "strict-format")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowZeroLengthValues, state.AllowZeroLengthValues, "allow-zero-length-values")
	operations.AddBoolOperationIfNecessary(&ops, plan.StripSyntaxMinUpperBound, state.StripSyntaxMinUpperBound, "strip-syntax-min-upper-bound")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireBinaryTransfer, state.RequireBinaryTransfer, "require-binary-transfer")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *attributeSyntaxResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultAttributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultAttributeSyntaxResourceModel
	if plan.Type.ValueString() == "attribute-type-description" {
		readAttributeTypeDescriptionAttributeSyntaxResponseDefault(ctx, readResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "directory-string" {
		readDirectoryStringAttributeSyntaxResponseDefault(ctx, readResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "telephone-number" {
		readTelephoneNumberAttributeSyntaxResponseDefault(ctx, readResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "distinguished-name" {
		readDistinguishedNameAttributeSyntaxResponseDefault(ctx, readResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "generalized-time" {
		readGeneralizedTimeAttributeSyntaxResponseDefault(ctx, readResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "integer" {
		readIntegerAttributeSyntaxResponseDefault(ctx, readResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "uuid" {
		readUuidAttributeSyntaxResponseDefault(ctx, readResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "generic" {
		readGenericAttributeSyntaxResponseDefault(ctx, readResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "json-object" {
		readJsonObjectAttributeSyntaxResponseDefault(ctx, readResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "user-password" {
		readUserPasswordAttributeSyntaxResponseDefault(ctx, readResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "boolean" {
		readBooleanAttributeSyntaxResponseDefault(ctx, readResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "hex-string" {
		readHexStringAttributeSyntaxResponseDefault(ctx, readResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "bit-string" {
		readBitStringAttributeSyntaxResponseDefault(ctx, readResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "ldap-url" {
		readLdapUrlAttributeSyntaxResponseDefault(ctx, readResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "name-and-optional-uid" {
		readNameAndOptionalUidAttributeSyntaxResponseDefault(ctx, readResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAttributeSyntaxOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Attribute Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "attribute-type-description" {
			readAttributeTypeDescriptionAttributeSyntaxResponseDefault(ctx, updateResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "directory-string" {
			readDirectoryStringAttributeSyntaxResponseDefault(ctx, updateResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "telephone-number" {
			readTelephoneNumberAttributeSyntaxResponseDefault(ctx, updateResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "distinguished-name" {
			readDistinguishedNameAttributeSyntaxResponseDefault(ctx, updateResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "generalized-time" {
			readGeneralizedTimeAttributeSyntaxResponseDefault(ctx, updateResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "integer" {
			readIntegerAttributeSyntaxResponseDefault(ctx, updateResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "uuid" {
			readUuidAttributeSyntaxResponseDefault(ctx, updateResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "generic" {
			readGenericAttributeSyntaxResponseDefault(ctx, updateResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "json-object" {
			readJsonObjectAttributeSyntaxResponseDefault(ctx, updateResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "user-password" {
			readUserPasswordAttributeSyntaxResponseDefault(ctx, updateResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "boolean" {
			readBooleanAttributeSyntaxResponseDefault(ctx, updateResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "hex-string" {
			readHexStringAttributeSyntaxResponseDefault(ctx, updateResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "bit-string" {
			readBitStringAttributeSyntaxResponseDefault(ctx, updateResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldap-url" {
			readLdapUrlAttributeSyntaxResponseDefault(ctx, updateResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "name-and-optional-uid" {
			readNameAndOptionalUidAttributeSyntaxResponseDefault(ctx, updateResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *attributeSyntaxResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state attributeSyntaxResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Attribute Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *attributeSyntaxResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan attributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state attributeSyntaxResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAttributeSyntaxOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntaxExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Attribute Syntax", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *attributeSyntaxResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *attributeSyntaxResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
