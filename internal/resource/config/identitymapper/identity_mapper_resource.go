package identitymapper

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	_ resource.Resource                = &identityMapperResource{}
	_ resource.ResourceWithConfigure   = &identityMapperResource{}
	_ resource.ResourceWithImportState = &identityMapperResource{}
	_ resource.Resource                = &defaultIdentityMapperResource{}
	_ resource.ResourceWithConfigure   = &defaultIdentityMapperResource{}
	_ resource.ResourceWithImportState = &defaultIdentityMapperResource{}
)

// Create a Identity Mapper resource
func NewIdentityMapperResource() resource.Resource {
	return &identityMapperResource{}
}

func NewDefaultIdentityMapperResource() resource.Resource {
	return &defaultIdentityMapperResource{}
}

// identityMapperResource is the resource implementation.
type identityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultIdentityMapperResource is the resource implementation.
type defaultIdentityMapperResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *identityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_mapper"
}

func (r *defaultIdentityMapperResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_identity_mapper"
}

// Configure adds the provider configured client to the resource.
func (r *identityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultIdentityMapperResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type identityMapperResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	Type                      types.String `tfsdk:"type"`
	ExtensionClass            types.String `tfsdk:"extension_class"`
	ExtensionArgument         types.Set    `tfsdk:"extension_argument"`
	AllIncludedIdentityMapper types.Set    `tfsdk:"all_included_identity_mapper"`
	AnyIncludedIdentityMapper types.Set    `tfsdk:"any_included_identity_mapper"`
	ScriptClass               types.String `tfsdk:"script_class"`
	ScriptArgument            types.Set    `tfsdk:"script_argument"`
	MatchAttribute            types.Set    `tfsdk:"match_attribute"`
	MatchPattern              types.String `tfsdk:"match_pattern"`
	ReplacePattern            types.String `tfsdk:"replace_pattern"`
	MatchBaseDN               types.Set    `tfsdk:"match_base_dn"`
	MatchFilter               types.String `tfsdk:"match_filter"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *identityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	identityMapperSchema(ctx, req, resp, false)
}

func (r *defaultIdentityMapperResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	identityMapperSchema(ctx, req, resp, true)
}

func identityMapperSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Identity Mapper.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Identity Mapper resource. Options are ['exact-match', 'groovy-scripted', 'regular-expression', 'aggregate', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"exact-match", "groovy-scripted", "regular-expression", "aggregate", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Identity Mapper.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Identity Mapper. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"all_included_identity_mapper": schema.SetAttribute{
				Description: "The set of identity mappers that must all match the target entry. Each identity mapper must uniquely match the same target entry. If any of the identity mappers match multiple entries, if any of them match zero entries, or if any of them match different entries, then the mapping will fail.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"any_included_identity_mapper": schema.SetAttribute{
				Description: "The set of identity mappers that will be used to identify the target entry. At least one identity mapper must uniquely match an entry. If multiple identity mappers match entries, then they must all uniquely match the same entry. If none of the identity mappers match any entries, if any of them match multiple entries, or if any of them match different entries, then the mapping will fail.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Identity Mapper.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Identity Mapper. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"match_attribute": schema.SetAttribute{
				Description: "Specifies the attribute whose value should exactly match the ID string provided to this identity mapper.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"match_pattern": schema.StringAttribute{
				Description: "Specifies the regular expression pattern that is used to identify portions of the ID string that will be replaced.",
				Optional:    true,
			},
			"replace_pattern": schema.StringAttribute{
				Description: "Specifies the replacement pattern that should be used for substrings in the ID string that match the provided regular expression pattern.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"match_base_dn": schema.SetAttribute{
				Description: "Specifies the set of base DNs below which to search for users.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"match_filter": schema.StringAttribute{
				Description: "An optional filter that mapped users must match.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Identity Mapper",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Identity Mapper is enabled for use.",
				Required:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputed(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsIdentityMapper() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"aggregate"},
			resourcevalidator.ExactlyOneOf(
				path.MatchRoot("any_included_identity_mapper"),
				path.MatchRoot("all_included_identity_mapper"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_included_identity_mapper"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("replace_pattern"),
			path.MatchRoot("type"),
			[]string{"regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_filter"),
			path.MatchRoot("type"),
			[]string{"exact-match", "regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_base_dn"),
			path.MatchRoot("type"),
			[]string{"exact-match", "regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_attribute"),
			path.MatchRoot("type"),
			[]string{"exact-match", "regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("match_pattern"),
			path.MatchRoot("type"),
			[]string{"regular-expression"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_included_identity_mapper"),
			path.MatchRoot("type"),
			[]string{"aggregate"},
		),
	}
}

// Add config validators
func (r identityMapperResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsIdentityMapper()
}

// Add config validators
func (r defaultIdentityMapperResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsIdentityMapper()
}

// Add optional fields to create request for exact-match identity-mapper
func addOptionalExactMatchIdentityMapperFields(ctx context.Context, addRequest *client.AddExactMatchIdentityMapperRequest, plan identityMapperResourceModel) {
	if internaltypes.IsDefined(plan.MatchAttribute) {
		var slice []string
		plan.MatchAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MatchAttribute = slice
	}
	if internaltypes.IsDefined(plan.MatchBaseDN) {
		var slice []string
		plan.MatchBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.MatchBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MatchFilter) {
		addRequest.MatchFilter = plan.MatchFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for groovy-scripted identity-mapper
func addOptionalGroovyScriptedIdentityMapperFields(ctx context.Context, addRequest *client.AddGroovyScriptedIdentityMapperRequest, plan identityMapperResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for regular-expression identity-mapper
func addOptionalRegularExpressionIdentityMapperFields(ctx context.Context, addRequest *client.AddRegularExpressionIdentityMapperRequest, plan identityMapperResourceModel) {
	if internaltypes.IsDefined(plan.MatchAttribute) {
		var slice []string
		plan.MatchAttribute.ElementsAs(ctx, &slice, false)
		addRequest.MatchAttribute = slice
	}
	if internaltypes.IsDefined(plan.MatchBaseDN) {
		var slice []string
		plan.MatchBaseDN.ElementsAs(ctx, &slice, false)
		addRequest.MatchBaseDN = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MatchFilter) {
		addRequest.MatchFilter = plan.MatchFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReplacePattern) {
		addRequest.ReplacePattern = plan.ReplacePattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for aggregate identity-mapper
func addOptionalAggregateIdentityMapperFields(ctx context.Context, addRequest *client.AddAggregateIdentityMapperRequest, plan identityMapperResourceModel) {
	if internaltypes.IsDefined(plan.AllIncludedIdentityMapper) {
		var slice []string
		plan.AllIncludedIdentityMapper.ElementsAs(ctx, &slice, false)
		addRequest.AllIncludedIdentityMapper = slice
	}
	if internaltypes.IsDefined(plan.AnyIncludedIdentityMapper) {
		var slice []string
		plan.AnyIncludedIdentityMapper.ElementsAs(ctx, &slice, false)
		addRequest.AnyIncludedIdentityMapper = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party identity-mapper
func addOptionalThirdPartyIdentityMapperFields(ctx context.Context, addRequest *client.AddThirdPartyIdentityMapperRequest, plan identityMapperResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateIdentityMapperUnknownValues(ctx context.Context, model *identityMapperResourceModel) {
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.MatchBaseDN.ElementType(ctx) == nil {
		model.MatchBaseDN = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.AllIncludedIdentityMapper.ElementType(ctx) == nil {
		model.AllIncludedIdentityMapper = types.SetNull(types.StringType)
	}
	if model.MatchAttribute.ElementType(ctx) == nil {
		model.MatchAttribute = types.SetNull(types.StringType)
	}
	if model.AnyIncludedIdentityMapper.ElementType(ctx) == nil {
		model.AnyIncludedIdentityMapper = types.SetNull(types.StringType)
	}
}

// Read a ExactMatchIdentityMapperResponse object into the model struct
func readExactMatchIdentityMapperResponse(ctx context.Context, r *client.ExactMatchIdentityMapperResponse, state *identityMapperResourceModel, expectedValues *identityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exact-match")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.MatchBaseDN = internaltypes.GetStringSet(r.MatchBaseDN)
	state.MatchFilter = internaltypes.StringTypeOrNil(r.MatchFilter, internaltypes.IsEmptyString(expectedValues.MatchFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateIdentityMapperUnknownValues(ctx, state)
}

// Read a GroovyScriptedIdentityMapperResponse object into the model struct
func readGroovyScriptedIdentityMapperResponse(ctx context.Context, r *client.GroovyScriptedIdentityMapperResponse, state *identityMapperResourceModel, expectedValues *identityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateIdentityMapperUnknownValues(ctx, state)
}

// Read a RegularExpressionIdentityMapperResponse object into the model struct
func readRegularExpressionIdentityMapperResponse(ctx context.Context, r *client.RegularExpressionIdentityMapperResponse, state *identityMapperResourceModel, expectedValues *identityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("regular-expression")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.MatchBaseDN = internaltypes.GetStringSet(r.MatchBaseDN)
	state.MatchFilter = internaltypes.StringTypeOrNil(r.MatchFilter, internaltypes.IsEmptyString(expectedValues.MatchFilter))
	state.MatchPattern = types.StringValue(r.MatchPattern)
	state.ReplacePattern = internaltypes.StringTypeOrNil(r.ReplacePattern, internaltypes.IsEmptyString(expectedValues.ReplacePattern))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateIdentityMapperUnknownValues(ctx, state)
}

// Read a AggregateIdentityMapperResponse object into the model struct
func readAggregateIdentityMapperResponse(ctx context.Context, r *client.AggregateIdentityMapperResponse, state *identityMapperResourceModel, expectedValues *identityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedIdentityMapper = internaltypes.GetStringSet(r.AllIncludedIdentityMapper)
	state.AnyIncludedIdentityMapper = internaltypes.GetStringSet(r.AnyIncludedIdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateIdentityMapperUnknownValues(ctx, state)
}

// Read a ThirdPartyIdentityMapperResponse object into the model struct
func readThirdPartyIdentityMapperResponse(ctx context.Context, r *client.ThirdPartyIdentityMapperResponse, state *identityMapperResourceModel, expectedValues *identityMapperResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateIdentityMapperUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createIdentityMapperOperations(plan identityMapperResourceModel, state identityMapperResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllIncludedIdentityMapper, state.AllIncludedIdentityMapper, "all-included-identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyIncludedIdentityMapper, state.AnyIncludedIdentityMapper, "any-included-identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MatchAttribute, state.MatchAttribute, "match-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.MatchPattern, state.MatchPattern, "match-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplacePattern, state.ReplacePattern, "replace-pattern")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MatchBaseDN, state.MatchBaseDN, "match-base-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.MatchFilter, state.MatchFilter, "match-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a exact-match identity-mapper
func (r *identityMapperResource) CreateExactMatchIdentityMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan identityMapperResourceModel) (*identityMapperResourceModel, error) {
	addRequest := client.NewAddExactMatchIdentityMapperRequest(plan.Name.ValueString(),
		[]client.EnumexactMatchIdentityMapperSchemaUrn{client.ENUMEXACTMATCHIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPEREXACT_MATCH},
		plan.Enabled.ValueBool())
	addOptionalExactMatchIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddExactMatchIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Identity Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state identityMapperResourceModel
	readExactMatchIdentityMapperResponse(ctx, addResponse.ExactMatchIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted identity-mapper
func (r *identityMapperResource) CreateGroovyScriptedIdentityMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan identityMapperResourceModel) (*identityMapperResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedIdentityMapperRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedIdentityMapperSchemaUrn{client.ENUMGROOVYSCRIPTEDIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddGroovyScriptedIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Identity Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state identityMapperResourceModel
	readGroovyScriptedIdentityMapperResponse(ctx, addResponse.GroovyScriptedIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a regular-expression identity-mapper
func (r *identityMapperResource) CreateRegularExpressionIdentityMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan identityMapperResourceModel) (*identityMapperResourceModel, error) {
	addRequest := client.NewAddRegularExpressionIdentityMapperRequest(plan.Name.ValueString(),
		[]client.EnumregularExpressionIdentityMapperSchemaUrn{client.ENUMREGULAREXPRESSIONIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERREGULAR_EXPRESSION},
		plan.MatchPattern.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalRegularExpressionIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddRegularExpressionIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Identity Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state identityMapperResourceModel
	readRegularExpressionIdentityMapperResponse(ctx, addResponse.RegularExpressionIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aggregate identity-mapper
func (r *identityMapperResource) CreateAggregateIdentityMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan identityMapperResourceModel) (*identityMapperResourceModel, error) {
	addRequest := client.NewAddAggregateIdentityMapperRequest(plan.Name.ValueString(),
		[]client.EnumaggregateIdentityMapperSchemaUrn{client.ENUMAGGREGATEIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERAGGREGATE},
		plan.Enabled.ValueBool())
	addOptionalAggregateIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddAggregateIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Identity Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state identityMapperResourceModel
	readAggregateIdentityMapperResponse(ctx, addResponse.AggregateIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party identity-mapper
func (r *identityMapperResource) CreateThirdPartyIdentityMapper(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan identityMapperResourceModel) (*identityMapperResourceModel, error) {
	addRequest := client.NewAddThirdPartyIdentityMapperRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyIdentityMapperSchemaUrn{client.ENUMTHIRDPARTYIDENTITYMAPPERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0IDENTITY_MAPPERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyIdentityMapperFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.IdentityMapperApi.AddIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddIdentityMapperRequest(
		client.AddThirdPartyIdentityMapperRequestAsAddIdentityMapperRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.IdentityMapperApi.AddIdentityMapperExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Identity Mapper", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state identityMapperResourceModel
	readThirdPartyIdentityMapperResponse(ctx, addResponse.ThirdPartyIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *identityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan identityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *identityMapperResourceModel
	var err error
	if plan.Type.ValueString() == "exact-match" {
		state, err = r.CreateExactMatchIdentityMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedIdentityMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "regular-expression" {
		state, err = r.CreateRegularExpressionIdentityMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aggregate" {
		state, err = r.CreateAggregateIdentityMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyIdentityMapper(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultIdentityMapperResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan identityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state identityMapperResourceModel
	if readResponse.ExactMatchIdentityMapperResponse != nil {
		readExactMatchIdentityMapperResponse(ctx, readResponse.ExactMatchIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedIdentityMapperResponse != nil {
		readGroovyScriptedIdentityMapperResponse(ctx, readResponse.GroovyScriptedIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.RegularExpressionIdentityMapperResponse != nil {
		readRegularExpressionIdentityMapperResponse(ctx, readResponse.RegularExpressionIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateIdentityMapperResponse != nil {
		readAggregateIdentityMapperResponse(ctx, readResponse.AggregateIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyIdentityMapperResponse != nil {
		readThirdPartyIdentityMapperResponse(ctx, readResponse.ThirdPartyIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.IdentityMapperApi.UpdateIdentityMapper(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ExactMatchIdentityMapperResponse != nil {
			readExactMatchIdentityMapperResponse(ctx, updateResponse.ExactMatchIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedIdentityMapperResponse != nil {
			readGroovyScriptedIdentityMapperResponse(ctx, updateResponse.GroovyScriptedIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.RegularExpressionIdentityMapperResponse != nil {
			readRegularExpressionIdentityMapperResponse(ctx, updateResponse.RegularExpressionIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateIdentityMapperResponse != nil {
			readAggregateIdentityMapperResponse(ctx, updateResponse.AggregateIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyIdentityMapperResponse != nil {
			readThirdPartyIdentityMapperResponse(ctx, updateResponse.ThirdPartyIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *identityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdentityMapperResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIdentityMapper(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state identityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ExactMatchIdentityMapperResponse != nil {
		readExactMatchIdentityMapperResponse(ctx, readResponse.ExactMatchIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedIdentityMapperResponse != nil {
		readGroovyScriptedIdentityMapperResponse(ctx, readResponse.GroovyScriptedIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.RegularExpressionIdentityMapperResponse != nil {
		readRegularExpressionIdentityMapperResponse(ctx, readResponse.RegularExpressionIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateIdentityMapperResponse != nil {
		readAggregateIdentityMapperResponse(ctx, readResponse.AggregateIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyIdentityMapperResponse != nil {
		readThirdPartyIdentityMapperResponse(ctx, readResponse.ThirdPartyIdentityMapperResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *identityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIdentityMapperResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIdentityMapper(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIdentityMapper(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan identityMapperResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state identityMapperResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.IdentityMapperApi.UpdateIdentityMapper(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createIdentityMapperOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.IdentityMapperApi.UpdateIdentityMapperExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Identity Mapper", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ExactMatchIdentityMapperResponse != nil {
			readExactMatchIdentityMapperResponse(ctx, updateResponse.ExactMatchIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedIdentityMapperResponse != nil {
			readGroovyScriptedIdentityMapperResponse(ctx, updateResponse.GroovyScriptedIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.RegularExpressionIdentityMapperResponse != nil {
			readRegularExpressionIdentityMapperResponse(ctx, updateResponse.RegularExpressionIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AggregateIdentityMapperResponse != nil {
			readAggregateIdentityMapperResponse(ctx, updateResponse.AggregateIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyIdentityMapperResponse != nil {
			readThirdPartyIdentityMapperResponse(ctx, updateResponse.ThirdPartyIdentityMapperResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultIdentityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *identityMapperResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state identityMapperResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.IdentityMapperApi.DeleteIdentityMapperExecute(r.apiClient.IdentityMapperApi.DeleteIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Identity Mapper", err, httpResp)
		return
	}
}

func (r *identityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdentityMapper(ctx, req, resp)
}

func (r *defaultIdentityMapperResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIdentityMapper(ctx, req, resp)
}

func importIdentityMapper(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
