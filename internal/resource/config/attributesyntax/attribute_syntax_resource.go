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
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
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
	r.apiClient = providerCfg.ApiClientV9300
}

type attributeSyntaxResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
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
				Optional:    false,
				Required:    false,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"attribute-type-description", "directory-string", "telephone-number", "distinguished-name", "generalized-time", "integer", "uuid", "generic", "json-object", "user-password", "boolean", "hex-string", "bit-string", "ldap-url", "name-and-optional-uid"}...),
				},
			},
			"enable_compaction": schema.BoolAttribute{
				Description: "Indicates whether values of attributes with this syntax should be compacted when stored in a local DB database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_attribute_in_compaction": schema.SetAttribute{
				Description: "Specifies the specific attributes (which should be associated with this syntax) whose values should be compacted. If one or more include attributes are specified, then only those attributes will have their values compacted. If not set then all attributes will have their values compacted. The exclude-attribute-from-compaction property takes precedence over this property.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_attribute_from_compaction": schema.SetAttribute{
				Description: "Specifies the specific attributes (which should be associated with this syntax) whose values should not be compacted. If one or more exclude attributes are specified, then values of those attributes will not have their values compacted. This property takes precedence over the include-attribute-in-compaction property.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"strict_format": schema.BoolAttribute{
				Description: " When the `type` value is one of [`telephone-number`]: Indicates whether to require telephone number values to strictly comply with the standard definition for this syntax. When the `type` value is one of [`ldap-url`]: Indicates whether values for attributes with this syntax will be required to be in the valid LDAP URL format. If this is set to false, then arbitrary strings will be allowed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_zero_length_values": schema.BoolAttribute{
				Description: "Indicates whether zero-length (that is, an empty string) values are allowed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"strip_syntax_min_upper_bound": schema.BoolAttribute{
				Description: "Indicates whether the suggested minimum upper bound appended to an attribute's syntax OID in its schema definition Attribute Type Description should be stripped.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Attribute Syntax is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"require_binary_transfer": schema.BoolAttribute{
				Description: "Indicates whether values of this attribute are required to have a \"binary\" transfer option as described in RFC 4522. Attributes with this syntax will generally be referenced with names including \";binary\" (e.g., \"userCertificate;binary\").",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators
func (r attributeSyntaxResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("strict_format"),
			path.MatchRoot("type"),
			[]string{"telephone-number", "ldap-url"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_attribute_from_compaction"),
			path.MatchRoot("type"),
			[]string{"distinguished-name", "generalized-time", "integer", "uuid", "user-password", "boolean", "bit-string", "name-and-optional-uid"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("enable_compaction"),
			path.MatchRoot("type"),
			[]string{"distinguished-name", "generalized-time", "integer", "uuid", "user-password", "boolean", "bit-string", "name-and-optional-uid"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_attribute_in_compaction"),
			path.MatchRoot("type"),
			[]string{"distinguished-name", "generalized-time", "integer", "uuid", "user-password", "boolean", "bit-string", "name-and-optional-uid"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("strip_syntax_min_upper_bound"),
			path.MatchRoot("type"),
			[]string{"attribute-type-description"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_zero_length_values"),
			path.MatchRoot("type"),
			[]string{"directory-string"},
		),
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateAttributeSyntaxUnknownValues(ctx context.Context, model *attributeSyntaxResourceModel) {
	if model.ExcludeAttributeFromCompaction.ElementType(ctx) == nil {
		model.ExcludeAttributeFromCompaction = types.SetNull(types.StringType)
	}
	if model.IncludeAttributeInCompaction.ElementType(ctx) == nil {
		model.IncludeAttributeInCompaction = types.SetNull(types.StringType)
	}
}

// Read a AttributeTypeDescriptionAttributeSyntaxResponse object into the model struct
func readAttributeTypeDescriptionAttributeSyntaxResponse(ctx context.Context, r *client.AttributeTypeDescriptionAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-type-description")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.StripSyntaxMinUpperBound = internaltypes.BoolTypeOrNil(r.StripSyntaxMinUpperBound)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a DirectoryStringAttributeSyntaxResponse object into the model struct
func readDirectoryStringAttributeSyntaxResponse(ctx context.Context, r *client.DirectoryStringAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory-string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowZeroLengthValues = internaltypes.BoolTypeOrNil(r.AllowZeroLengthValues)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a TelephoneNumberAttributeSyntaxResponse object into the model struct
func readTelephoneNumberAttributeSyntaxResponse(ctx context.Context, r *client.TelephoneNumberAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("telephone-number")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a DistinguishedNameAttributeSyntaxResponse object into the model struct
func readDistinguishedNameAttributeSyntaxResponse(ctx context.Context, r *client.DistinguishedNameAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("distinguished-name")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a GeneralizedTimeAttributeSyntaxResponse object into the model struct
func readGeneralizedTimeAttributeSyntaxResponse(ctx context.Context, r *client.GeneralizedTimeAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generalized-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a IntegerAttributeSyntaxResponse object into the model struct
func readIntegerAttributeSyntaxResponse(ctx context.Context, r *client.IntegerAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("integer")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a UuidAttributeSyntaxResponse object into the model struct
func readUuidAttributeSyntaxResponse(ctx context.Context, r *client.UuidAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("uuid")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a GenericAttributeSyntaxResponse object into the model struct
func readGenericAttributeSyntaxResponse(ctx context.Context, r *client.GenericAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a JsonObjectAttributeSyntaxResponse object into the model struct
func readJsonObjectAttributeSyntaxResponse(ctx context.Context, r *client.JsonObjectAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-object")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a UserPasswordAttributeSyntaxResponse object into the model struct
func readUserPasswordAttributeSyntaxResponse(ctx context.Context, r *client.UserPasswordAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a BooleanAttributeSyntaxResponse object into the model struct
func readBooleanAttributeSyntaxResponse(ctx context.Context, r *client.BooleanAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("boolean")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a HexStringAttributeSyntaxResponse object into the model struct
func readHexStringAttributeSyntaxResponse(ctx context.Context, r *client.HexStringAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("hex-string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a BitStringAttributeSyntaxResponse object into the model struct
func readBitStringAttributeSyntaxResponse(ctx context.Context, r *client.BitStringAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("bit-string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a LdapUrlAttributeSyntaxResponse object into the model struct
func readLdapUrlAttributeSyntaxResponse(ctx context.Context, r *client.LdapUrlAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-url")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.StrictFormat = internaltypes.BoolTypeOrNil(r.StrictFormat)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Read a NameAndOptionalUidAttributeSyntaxResponse object into the model struct
func readNameAndOptionalUidAttributeSyntaxResponse(ctx context.Context, r *client.NameAndOptionalUidAttributeSyntaxResponse, state *attributeSyntaxResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("name-and-optional-uid")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnableCompaction = internaltypes.BoolTypeOrNil(r.EnableCompaction)
	state.IncludeAttributeInCompaction = internaltypes.GetStringSet(r.IncludeAttributeInCompaction)
	state.ExcludeAttributeFromCompaction = internaltypes.GetStringSet(r.ExcludeAttributeFromCompaction)
	state.Enabled = types.BoolValue(r.Enabled)
	state.RequireBinaryTransfer = internaltypes.BoolTypeOrNil(r.RequireBinaryTransfer)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAttributeSyntaxUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createAttributeSyntaxOperations(plan attributeSyntaxResourceModel, state attributeSyntaxResourceModel) []client.Operation {
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
	var plan attributeSyntaxResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AttributeSyntaxApi.GetAttributeSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
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
	var state attributeSyntaxResourceModel
	if readResponse.AttributeTypeDescriptionAttributeSyntaxResponse != nil {
		readAttributeTypeDescriptionAttributeSyntaxResponse(ctx, readResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DirectoryStringAttributeSyntaxResponse != nil {
		readDirectoryStringAttributeSyntaxResponse(ctx, readResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TelephoneNumberAttributeSyntaxResponse != nil {
		readTelephoneNumberAttributeSyntaxResponse(ctx, readResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DistinguishedNameAttributeSyntaxResponse != nil {
		readDistinguishedNameAttributeSyntaxResponse(ctx, readResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GeneralizedTimeAttributeSyntaxResponse != nil {
		readGeneralizedTimeAttributeSyntaxResponse(ctx, readResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.IntegerAttributeSyntaxResponse != nil {
		readIntegerAttributeSyntaxResponse(ctx, readResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UuidAttributeSyntaxResponse != nil {
		readUuidAttributeSyntaxResponse(ctx, readResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericAttributeSyntaxResponse != nil {
		readGenericAttributeSyntaxResponse(ctx, readResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JsonObjectAttributeSyntaxResponse != nil {
		readJsonObjectAttributeSyntaxResponse(ctx, readResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UserPasswordAttributeSyntaxResponse != nil {
		readUserPasswordAttributeSyntaxResponse(ctx, readResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BooleanAttributeSyntaxResponse != nil {
		readBooleanAttributeSyntaxResponse(ctx, readResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HexStringAttributeSyntaxResponse != nil {
		readHexStringAttributeSyntaxResponse(ctx, readResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BitStringAttributeSyntaxResponse != nil {
		readBitStringAttributeSyntaxResponse(ctx, readResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapUrlAttributeSyntaxResponse != nil {
		readLdapUrlAttributeSyntaxResponse(ctx, readResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NameAndOptionalUidAttributeSyntaxResponse != nil {
		readNameAndOptionalUidAttributeSyntaxResponse(ctx, readResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AttributeSyntaxApi.UpdateAttributeSyntax(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
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
		if updateResponse.AttributeTypeDescriptionAttributeSyntaxResponse != nil {
			readAttributeTypeDescriptionAttributeSyntaxResponse(ctx, updateResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.DirectoryStringAttributeSyntaxResponse != nil {
			readDirectoryStringAttributeSyntaxResponse(ctx, updateResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.TelephoneNumberAttributeSyntaxResponse != nil {
			readTelephoneNumberAttributeSyntaxResponse(ctx, updateResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.DistinguishedNameAttributeSyntaxResponse != nil {
			readDistinguishedNameAttributeSyntaxResponse(ctx, updateResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GeneralizedTimeAttributeSyntaxResponse != nil {
			readGeneralizedTimeAttributeSyntaxResponse(ctx, updateResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.IntegerAttributeSyntaxResponse != nil {
			readIntegerAttributeSyntaxResponse(ctx, updateResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.UuidAttributeSyntaxResponse != nil {
			readUuidAttributeSyntaxResponse(ctx, updateResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GenericAttributeSyntaxResponse != nil {
			readGenericAttributeSyntaxResponse(ctx, updateResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.JsonObjectAttributeSyntaxResponse != nil {
			readJsonObjectAttributeSyntaxResponse(ctx, updateResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.UserPasswordAttributeSyntaxResponse != nil {
			readUserPasswordAttributeSyntaxResponse(ctx, updateResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.BooleanAttributeSyntaxResponse != nil {
			readBooleanAttributeSyntaxResponse(ctx, updateResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.HexStringAttributeSyntaxResponse != nil {
			readHexStringAttributeSyntaxResponse(ctx, updateResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.BitStringAttributeSyntaxResponse != nil {
			readBitStringAttributeSyntaxResponse(ctx, updateResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.LdapUrlAttributeSyntaxResponse != nil {
			readLdapUrlAttributeSyntaxResponse(ctx, updateResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.NameAndOptionalUidAttributeSyntaxResponse != nil {
			readNameAndOptionalUidAttributeSyntaxResponse(ctx, updateResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
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
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
	if readResponse.AttributeTypeDescriptionAttributeSyntaxResponse != nil {
		readAttributeTypeDescriptionAttributeSyntaxResponse(ctx, readResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DirectoryStringAttributeSyntaxResponse != nil {
		readDirectoryStringAttributeSyntaxResponse(ctx, readResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TelephoneNumberAttributeSyntaxResponse != nil {
		readTelephoneNumberAttributeSyntaxResponse(ctx, readResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DistinguishedNameAttributeSyntaxResponse != nil {
		readDistinguishedNameAttributeSyntaxResponse(ctx, readResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GeneralizedTimeAttributeSyntaxResponse != nil {
		readGeneralizedTimeAttributeSyntaxResponse(ctx, readResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.IntegerAttributeSyntaxResponse != nil {
		readIntegerAttributeSyntaxResponse(ctx, readResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UuidAttributeSyntaxResponse != nil {
		readUuidAttributeSyntaxResponse(ctx, readResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericAttributeSyntaxResponse != nil {
		readGenericAttributeSyntaxResponse(ctx, readResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JsonObjectAttributeSyntaxResponse != nil {
		readJsonObjectAttributeSyntaxResponse(ctx, readResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.UserPasswordAttributeSyntaxResponse != nil {
		readUserPasswordAttributeSyntaxResponse(ctx, readResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BooleanAttributeSyntaxResponse != nil {
		readBooleanAttributeSyntaxResponse(ctx, readResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HexStringAttributeSyntaxResponse != nil {
		readHexStringAttributeSyntaxResponse(ctx, readResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BitStringAttributeSyntaxResponse != nil {
		readBitStringAttributeSyntaxResponse(ctx, readResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapUrlAttributeSyntaxResponse != nil {
		readLdapUrlAttributeSyntaxResponse(ctx, readResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NameAndOptionalUidAttributeSyntaxResponse != nil {
		readNameAndOptionalUidAttributeSyntaxResponse(ctx, readResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
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
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

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
		if updateResponse.AttributeTypeDescriptionAttributeSyntaxResponse != nil {
			readAttributeTypeDescriptionAttributeSyntaxResponse(ctx, updateResponse.AttributeTypeDescriptionAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.DirectoryStringAttributeSyntaxResponse != nil {
			readDirectoryStringAttributeSyntaxResponse(ctx, updateResponse.DirectoryStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.TelephoneNumberAttributeSyntaxResponse != nil {
			readTelephoneNumberAttributeSyntaxResponse(ctx, updateResponse.TelephoneNumberAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.DistinguishedNameAttributeSyntaxResponse != nil {
			readDistinguishedNameAttributeSyntaxResponse(ctx, updateResponse.DistinguishedNameAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GeneralizedTimeAttributeSyntaxResponse != nil {
			readGeneralizedTimeAttributeSyntaxResponse(ctx, updateResponse.GeneralizedTimeAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.IntegerAttributeSyntaxResponse != nil {
			readIntegerAttributeSyntaxResponse(ctx, updateResponse.IntegerAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.UuidAttributeSyntaxResponse != nil {
			readUuidAttributeSyntaxResponse(ctx, updateResponse.UuidAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.GenericAttributeSyntaxResponse != nil {
			readGenericAttributeSyntaxResponse(ctx, updateResponse.GenericAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.JsonObjectAttributeSyntaxResponse != nil {
			readJsonObjectAttributeSyntaxResponse(ctx, updateResponse.JsonObjectAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.UserPasswordAttributeSyntaxResponse != nil {
			readUserPasswordAttributeSyntaxResponse(ctx, updateResponse.UserPasswordAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.BooleanAttributeSyntaxResponse != nil {
			readBooleanAttributeSyntaxResponse(ctx, updateResponse.BooleanAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.HexStringAttributeSyntaxResponse != nil {
			readHexStringAttributeSyntaxResponse(ctx, updateResponse.HexStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.BitStringAttributeSyntaxResponse != nil {
			readBitStringAttributeSyntaxResponse(ctx, updateResponse.BitStringAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.LdapUrlAttributeSyntaxResponse != nil {
			readLdapUrlAttributeSyntaxResponse(ctx, updateResponse.LdapUrlAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
		if updateResponse.NameAndOptionalUidAttributeSyntaxResponse != nil {
			readNameAndOptionalUidAttributeSyntaxResponse(ctx, updateResponse.NameAndOptionalUidAttributeSyntaxResponse, &state, &resp.Diagnostics)
		}
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
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
