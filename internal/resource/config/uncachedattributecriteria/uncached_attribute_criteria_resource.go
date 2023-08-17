package uncachedattributecriteria

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &uncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &uncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &uncachedAttributeCriteriaResource{}
	_ resource.Resource                = &defaultUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultUncachedAttributeCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultUncachedAttributeCriteriaResource{}
)

// Create a Uncached Attribute Criteria resource
func NewUncachedAttributeCriteriaResource() resource.Resource {
	return &uncachedAttributeCriteriaResource{}
}

func NewDefaultUncachedAttributeCriteriaResource() resource.Resource {
	return &defaultUncachedAttributeCriteriaResource{}
}

// uncachedAttributeCriteriaResource is the resource implementation.
type uncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultUncachedAttributeCriteriaResource is the resource implementation.
type defaultUncachedAttributeCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *uncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uncached_attribute_criteria"
}

func (r *defaultUncachedAttributeCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_uncached_attribute_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *uncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultUncachedAttributeCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type uncachedAttributeCriteriaResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	Type              types.String `tfsdk:"type"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	AttributeType     types.Set    `tfsdk:"attribute_type"`
	MinValueCount     types.Int64  `tfsdk:"min_value_count"`
	MinTotalValueSize types.String `tfsdk:"min_total_value_size"`
	ScriptClass       types.String `tfsdk:"script_class"`
	ScriptArgument    types.Set    `tfsdk:"script_argument"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *uncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	uncachedAttributeCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultUncachedAttributeCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	uncachedAttributeCriteriaSchema(ctx, req, resp, true)
}

func uncachedAttributeCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Uncached Attribute Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Uncached Attribute Criteria resource. Options are ['default', 'groovy-scripted', 'simple', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"default", "groovy-scripted", "simple", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Uncached Attribute Criteria.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"attribute_type": schema.SetAttribute{
				Description: "Specifies the attribute types for attributes that may be written to the uncached-id2entry database.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"min_value_count": schema.Int64Attribute{
				Description: "Specifies the minimum number of values that an attribute must have before it will be written into the uncached-id2entry database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_total_value_size": schema.StringAttribute{
				Description: "Specifies the minimum total value size (i.e., the sum of the sizes of all values) that an attribute must have before it will be written into the uncached-id2entry database.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Uncached Attribute Criteria.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Uncached Attribute Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Uncached Attribute Criteria",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Uncached Attribute Criteria is enabled for use in the server.",
				Required:    true,
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

// Add config validators that apply to both default_ and non-default_
func configValidatorsUncachedAttributeCriteria() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("attribute_type"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_value_count"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_total_value_size"),
			path.MatchRoot("type"),
			[]string{"simple"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
	}
}

// Add config validators
func (r uncachedAttributeCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsUncachedAttributeCriteria()
}

// Add config validators
func (r defaultUncachedAttributeCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsUncachedAttributeCriteria()
}

// Add optional fields to create request for default uncached-attribute-criteria
func addOptionalDefaultUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddDefaultUncachedAttributeCriteriaRequest, plan uncachedAttributeCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for groovy-scripted uncached-attribute-criteria
func addOptionalGroovyScriptedUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddGroovyScriptedUncachedAttributeCriteriaRequest, plan uncachedAttributeCriteriaResourceModel) {
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

// Add optional fields to create request for simple uncached-attribute-criteria
func addOptionalSimpleUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddSimpleUncachedAttributeCriteriaRequest, plan uncachedAttributeCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.MinValueCount) {
		addRequest.MinValueCount = plan.MinValueCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinTotalValueSize) {
		addRequest.MinTotalValueSize = plan.MinTotalValueSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party uncached-attribute-criteria
func addOptionalThirdPartyUncachedAttributeCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartyUncachedAttributeCriteriaRequest, plan uncachedAttributeCriteriaResourceModel) {
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
func populateUncachedAttributeCriteriaUnknownValues(model *uncachedAttributeCriteriaResourceModel) {
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AttributeType.IsUnknown() || model.AttributeType.IsNull() {
		model.AttributeType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MinTotalValueSize.IsUnknown() || model.MinTotalValueSize.IsNull() {
		model.MinTotalValueSize = types.StringValue("")
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *uncachedAttributeCriteriaResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
}

// Read a DefaultUncachedAttributeCriteriaResponse object into the model struct
func readDefaultUncachedAttributeCriteriaResponse(ctx context.Context, r *client.DefaultUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaResourceModel, expectedValues *uncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("default")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedAttributeCriteriaUnknownValues(state)
}

// Read a GroovyScriptedUncachedAttributeCriteriaResponse object into the model struct
func readGroovyScriptedUncachedAttributeCriteriaResponse(ctx context.Context, r *client.GroovyScriptedUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaResourceModel, expectedValues *uncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedAttributeCriteriaUnknownValues(state)
}

// Read a SimpleUncachedAttributeCriteriaResponse object into the model struct
func readSimpleUncachedAttributeCriteriaResponse(ctx context.Context, r *client.SimpleUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaResourceModel, expectedValues *uncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("simple")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.MinValueCount = internaltypes.Int64TypeOrNil(r.MinValueCount)
	state.MinTotalValueSize = internaltypes.StringTypeOrNil(r.MinTotalValueSize, true)
	config.CheckMismatchedPDFormattedAttributes("min_total_value_size",
		expectedValues.MinTotalValueSize, state.MinTotalValueSize, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedAttributeCriteriaUnknownValues(state)
}

// Read a ThirdPartyUncachedAttributeCriteriaResponse object into the model struct
func readThirdPartyUncachedAttributeCriteriaResponse(ctx context.Context, r *client.ThirdPartyUncachedAttributeCriteriaResponse, state *uncachedAttributeCriteriaResourceModel, expectedValues *uncachedAttributeCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedAttributeCriteriaUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createUncachedAttributeCriteriaOperations(plan uncachedAttributeCriteriaResourceModel, state uncachedAttributeCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeType, state.AttributeType, "attribute-type")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinValueCount, state.MinValueCount, "min-value-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MinTotalValueSize, state.MinTotalValueSize, "min-total-value-size")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a default uncached-attribute-criteria
func (r *uncachedAttributeCriteriaResource) CreateDefaultUncachedAttributeCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedAttributeCriteriaResourceModel) (*uncachedAttributeCriteriaResourceModel, error) {
	addRequest := client.NewAddDefaultUncachedAttributeCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumdefaultUncachedAttributeCriteriaSchemaUrn{client.ENUMDEFAULTUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIADEFAULT},
		plan.Enabled.ValueBool())
	addOptionalDefaultUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddDefaultUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Attribute Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedAttributeCriteriaResourceModel
	readDefaultUncachedAttributeCriteriaResponse(ctx, addResponse.DefaultUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted uncached-attribute-criteria
func (r *uncachedAttributeCriteriaResource) CreateGroovyScriptedUncachedAttributeCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedAttributeCriteriaResourceModel) (*uncachedAttributeCriteriaResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedUncachedAttributeCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedUncachedAttributeCriteriaSchemaUrn{client.ENUMGROOVYSCRIPTEDUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIAGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddGroovyScriptedUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Attribute Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedAttributeCriteriaResourceModel
	readGroovyScriptedUncachedAttributeCriteriaResponse(ctx, addResponse.GroovyScriptedUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a simple uncached-attribute-criteria
func (r *uncachedAttributeCriteriaResource) CreateSimpleUncachedAttributeCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedAttributeCriteriaResourceModel) (*uncachedAttributeCriteriaResourceModel, error) {
	var AttributeTypeSlice []string
	plan.AttributeType.ElementsAs(ctx, &AttributeTypeSlice, false)
	addRequest := client.NewAddSimpleUncachedAttributeCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumsimpleUncachedAttributeCriteriaSchemaUrn{client.ENUMSIMPLEUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIASIMPLE},
		AttributeTypeSlice,
		plan.Enabled.ValueBool())
	addOptionalSimpleUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddSimpleUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Attribute Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedAttributeCriteriaResourceModel
	readSimpleUncachedAttributeCriteriaResponse(ctx, addResponse.SimpleUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party uncached-attribute-criteria
func (r *uncachedAttributeCriteriaResource) CreateThirdPartyUncachedAttributeCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedAttributeCriteriaResourceModel) (*uncachedAttributeCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartyUncachedAttributeCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyUncachedAttributeCriteriaSchemaUrn{client.ENUMTHIRDPARTYUNCACHEDATTRIBUTECRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ATTRIBUTE_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyUncachedAttributeCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedAttributeCriteriaRequest(
		client.AddThirdPartyUncachedAttributeCriteriaRequestAsAddUncachedAttributeCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.AddUncachedAttributeCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Attribute Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedAttributeCriteriaResourceModel
	readThirdPartyUncachedAttributeCriteriaResponse(ctx, addResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *uncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan uncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *uncachedAttributeCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "default" {
		state, err = r.CreateDefaultUncachedAttributeCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedUncachedAttributeCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "simple" {
		state, err = r.CreateSimpleUncachedAttributeCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyUncachedAttributeCriteria(ctx, req, resp, plan)
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
func (r *defaultUncachedAttributeCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan uncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Attribute Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state uncachedAttributeCriteriaResourceModel
	if readResponse.DefaultUncachedAttributeCriteriaResponse != nil {
		readDefaultUncachedAttributeCriteriaResponse(ctx, readResponse.DefaultUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedUncachedAttributeCriteriaResponse != nil {
		readGroovyScriptedUncachedAttributeCriteriaResponse(ctx, readResponse.GroovyScriptedUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleUncachedAttributeCriteriaResponse != nil {
		readSimpleUncachedAttributeCriteriaResponse(ctx, readResponse.SimpleUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyUncachedAttributeCriteriaResponse != nil {
		readThirdPartyUncachedAttributeCriteriaResponse(ctx, readResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DefaultUncachedAttributeCriteriaResponse != nil {
			readDefaultUncachedAttributeCriteriaResponse(ctx, updateResponse.DefaultUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedUncachedAttributeCriteriaResponse != nil {
			readGroovyScriptedUncachedAttributeCriteriaResponse(ctx, updateResponse.GroovyScriptedUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimpleUncachedAttributeCriteriaResponse != nil {
			readSimpleUncachedAttributeCriteriaResponse(ctx, updateResponse.SimpleUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyUncachedAttributeCriteriaResponse != nil {
			readThirdPartyUncachedAttributeCriteriaResponse(ctx, updateResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *uncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultUncachedAttributeCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readUncachedAttributeCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state uncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.GetUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Attribute Criteria", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Attribute Criteria", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DefaultUncachedAttributeCriteriaResponse != nil {
		readDefaultUncachedAttributeCriteriaResponse(ctx, readResponse.DefaultUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedUncachedAttributeCriteriaResponse != nil {
		readGroovyScriptedUncachedAttributeCriteriaResponse(ctx, readResponse.GroovyScriptedUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SimpleUncachedAttributeCriteriaResponse != nil {
		readSimpleUncachedAttributeCriteriaResponse(ctx, readResponse.SimpleUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyUncachedAttributeCriteriaResponse != nil {
		readThirdPartyUncachedAttributeCriteriaResponse(ctx, readResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *uncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUncachedAttributeCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUncachedAttributeCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateUncachedAttributeCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan uncachedAttributeCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state uncachedAttributeCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createUncachedAttributeCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedAttributeCriteriaApi.UpdateUncachedAttributeCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Uncached Attribute Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DefaultUncachedAttributeCriteriaResponse != nil {
			readDefaultUncachedAttributeCriteriaResponse(ctx, updateResponse.DefaultUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedUncachedAttributeCriteriaResponse != nil {
			readGroovyScriptedUncachedAttributeCriteriaResponse(ctx, updateResponse.GroovyScriptedUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SimpleUncachedAttributeCriteriaResponse != nil {
			readSimpleUncachedAttributeCriteriaResponse(ctx, updateResponse.SimpleUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyUncachedAttributeCriteriaResponse != nil {
			readThirdPartyUncachedAttributeCriteriaResponse(ctx, updateResponse.ThirdPartyUncachedAttributeCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultUncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *uncachedAttributeCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state uncachedAttributeCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteriaExecute(r.apiClient.UncachedAttributeCriteriaApi.DeleteUncachedAttributeCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Uncached Attribute Criteria", err, httpResp)
		return
	}
}

func (r *uncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUncachedAttributeCriteria(ctx, req, resp)
}

func (r *defaultUncachedAttributeCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUncachedAttributeCriteria(ctx, req, resp)
}

func importUncachedAttributeCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
