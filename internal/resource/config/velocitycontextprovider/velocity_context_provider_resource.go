package velocitycontextprovider

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	_ resource.Resource                = &velocityContextProviderResource{}
	_ resource.ResourceWithConfigure   = &velocityContextProviderResource{}
	_ resource.ResourceWithImportState = &velocityContextProviderResource{}
	_ resource.Resource                = &defaultVelocityContextProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultVelocityContextProviderResource{}
	_ resource.ResourceWithImportState = &defaultVelocityContextProviderResource{}
)

// Create a Velocity Context Provider resource
func NewVelocityContextProviderResource() resource.Resource {
	return &velocityContextProviderResource{}
}

func NewDefaultVelocityContextProviderResource() resource.Resource {
	return &defaultVelocityContextProviderResource{}
}

// velocityContextProviderResource is the resource implementation.
type velocityContextProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultVelocityContextProviderResource is the resource implementation.
type defaultVelocityContextProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *velocityContextProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_velocity_context_provider"
}

func (r *defaultVelocityContextProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_velocity_context_provider"
}

// Configure adds the provider configured client to the resource.
func (r *velocityContextProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultVelocityContextProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type velocityContextProviderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Type                     types.String `tfsdk:"type"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
	ExtensionClass           types.String `tfsdk:"extension_class"`
	ExtensionArgument        types.Set    `tfsdk:"extension_argument"`
	RequestTool              types.Set    `tfsdk:"request_tool"`
	SessionTool              types.Set    `tfsdk:"session_tool"`
	ApplicationTool          types.Set    `tfsdk:"application_tool"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	HttpMethod               types.Set    `tfsdk:"http_method"`
	ObjectScope              types.String `tfsdk:"object_scope"`
	IncludedView             types.Set    `tfsdk:"included_view"`
	ExcludedView             types.Set    `tfsdk:"excluded_view"`
	ResponseHeader           types.Set    `tfsdk:"response_header"`
}

// GetSchema defines the schema for the resource.
func (r *velocityContextProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	velocityContextProviderSchema(ctx, req, resp, false)
}

func (r *defaultVelocityContextProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	velocityContextProviderSchema(ctx, req, resp, true)
}

func velocityContextProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Velocity Context Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Velocity Context Provider resource. Options are ['velocity-tools', 'custom', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"velocity-tools", "custom", "third-party"}...),
				},
			},
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Velocity Context Provider.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Velocity Context Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"request_tool": schema.SetAttribute{
				Description: "The fully-qualified name of a Velocity Tool class that will be initialized for each request. May optionally include a path to a properties file used to configure this tool separated from the class name by a semi-colon (;). The path may absolute or relative to the server root.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"session_tool": schema.SetAttribute{
				Description: "The fully-qualified name of a Velocity Tool class that will be initialized for each session. May optionally include a path to a properties file used to configure this tool separated from the class name by a semi-colon (;). The path may absolute or relative to the server root.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"application_tool": schema.SetAttribute{
				Description: "The fully-qualified name of a Velocity Tool class that will be initialized once for the life of the server. May optionally include a path to a properties file used to configure this tool separated from the class name by a semi-colon (;). The path may absolute or relative to the server root.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Velocity Context Provider is enabled. If set to 'false' this Velocity Context Provider will not contribute context content for any requests.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"http_method": schema.SetAttribute{
				Description: "Specifies the set of HTTP methods handled by this Velocity Context Provider, which will perform actions necessary to fulfill the request before updating the context for the response. The values of this property are not case-sensitive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"object_scope": schema.StringAttribute{
				Description: "Scope for context objects contributed by this Velocity Context Provider. Must be either 'request' or 'session' or 'application'.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("application"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"included_view": schema.SetAttribute{
				Description: "The name of a view for which this Velocity Context Provider will contribute content.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"excluded_view": schema.SetAttribute{
				Description: "The name of a view for which this Velocity Context Provider will not contribute content.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for template page requests to which this Velocity Context Provider contributes content.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "http_servlet_extension_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsVelocityContextProvider() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_method"),
			path.MatchRoot("type"),
			[]string{"custom", "third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("session_tool"),
			path.MatchRoot("type"),
			[]string{"velocity-tools"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_tool"),
			path.MatchRoot("type"),
			[]string{"velocity-tools"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("application_tool"),
			path.MatchRoot("type"),
			[]string{"velocity-tools"},
		),
	}
}

// Add config validators
func (r velocityContextProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsVelocityContextProvider()
}

// Add config validators
func (r defaultVelocityContextProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsVelocityContextProvider()
}

// Add optional fields to create request for velocity-tools velocity-context-provider
func addOptionalVelocityToolsVelocityContextProviderFields(ctx context.Context, addRequest *client.AddVelocityToolsVelocityContextProviderRequest, plan velocityContextProviderResourceModel) error {
	if internaltypes.IsDefined(plan.RequestTool) {
		var slice []string
		plan.RequestTool.ElementsAs(ctx, &slice, false)
		addRequest.RequestTool = slice
	}
	if internaltypes.IsDefined(plan.SessionTool) {
		var slice []string
		plan.SessionTool.ElementsAs(ctx, &slice, false)
		addRequest.SessionTool = slice
	}
	if internaltypes.IsDefined(plan.ApplicationTool) {
		var slice []string
		plan.ApplicationTool.ElementsAs(ctx, &slice, false)
		addRequest.ApplicationTool = slice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ObjectScope) {
		objectScope, err := client.NewEnumvelocityContextProviderObjectScopePropFromValue(plan.ObjectScope.ValueString())
		if err != nil {
			return err
		}
		addRequest.ObjectScope = objectScope
	}
	if internaltypes.IsDefined(plan.IncludedView) {
		var slice []string
		plan.IncludedView.ElementsAs(ctx, &slice, false)
		addRequest.IncludedView = slice
	}
	if internaltypes.IsDefined(plan.ExcludedView) {
		var slice []string
		plan.ExcludedView.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedView = slice
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	return nil
}

// Add optional fields to create request for third-party velocity-context-provider
func addOptionalThirdPartyVelocityContextProviderFields(ctx context.Context, addRequest *client.AddThirdPartyVelocityContextProviderRequest, plan velocityContextProviderResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ObjectScope) {
		objectScope, err := client.NewEnumvelocityContextProviderObjectScopePropFromValue(plan.ObjectScope.ValueString())
		if err != nil {
			return err
		}
		addRequest.ObjectScope = objectScope
	}
	if internaltypes.IsDefined(plan.IncludedView) {
		var slice []string
		plan.IncludedView.ElementsAs(ctx, &slice, false)
		addRequest.IncludedView = slice
	}
	if internaltypes.IsDefined(plan.ExcludedView) {
		var slice []string
		plan.ExcludedView.ElementsAs(ctx, &slice, false)
		addRequest.ExcludedView = slice
	}
	if internaltypes.IsDefined(plan.HttpMethod) {
		var slice []string
		plan.HttpMethod.ElementsAs(ctx, &slice, false)
		addRequest.HttpMethod = slice
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateVelocityContextProviderUnknownValues(ctx context.Context, model *velocityContextProviderResourceModel) {
	if model.RequestTool.ElementType(ctx) == nil {
		model.RequestTool = types.SetNull(types.StringType)
	}
	if model.ApplicationTool.ElementType(ctx) == nil {
		model.ApplicationTool = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.SessionTool.ElementType(ctx) == nil {
		model.SessionTool = types.SetNull(types.StringType)
	}
	if model.HttpMethod.ElementType(ctx) == nil {
		model.HttpMethod = types.SetNull(types.StringType)
	}
}

// Read a VelocityToolsVelocityContextProviderResponse object into the model struct
func readVelocityToolsVelocityContextProviderResponse(ctx context.Context, r *client.VelocityToolsVelocityContextProviderResponse, state *velocityContextProviderResourceModel, expectedValues *velocityContextProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("velocity-tools")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestTool = internaltypes.GetStringSet(r.RequestTool)
	state.SessionTool = internaltypes.GetStringSet(r.SessionTool)
	state.ApplicationTool = internaltypes.GetStringSet(r.ApplicationTool)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.ObjectScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvelocityContextProviderObjectScopeProp(r.ObjectScope), true)
	state.IncludedView = internaltypes.GetStringSet(r.IncludedView)
	state.ExcludedView = internaltypes.GetStringSet(r.ExcludedView)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVelocityContextProviderUnknownValues(ctx, state)
}

// Read a CustomVelocityContextProviderResponse object into the model struct
func readCustomVelocityContextProviderResponse(ctx context.Context, r *client.CustomVelocityContextProviderResponse, state *velocityContextProviderResourceModel, expectedValues *velocityContextProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.ObjectScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvelocityContextProviderObjectScopeProp(r.ObjectScope), true)
	state.IncludedView = internaltypes.GetStringSet(r.IncludedView)
	state.ExcludedView = internaltypes.GetStringSet(r.ExcludedView)
	state.HttpMethod = internaltypes.GetStringSet(r.HttpMethod)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVelocityContextProviderUnknownValues(ctx, state)
}

// Read a ThirdPartyVelocityContextProviderResponse object into the model struct
func readThirdPartyVelocityContextProviderResponse(ctx context.Context, r *client.ThirdPartyVelocityContextProviderResponse, state *velocityContextProviderResourceModel, expectedValues *velocityContextProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.ObjectScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvelocityContextProviderObjectScopeProp(r.ObjectScope), true)
	state.IncludedView = internaltypes.GetStringSet(r.IncludedView)
	state.ExcludedView = internaltypes.GetStringSet(r.ExcludedView)
	state.HttpMethod = internaltypes.GetStringSet(r.HttpMethod)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateVelocityContextProviderUnknownValues(ctx, state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *velocityContextProviderResourceModel) setStateValuesNotReturnedByAPI(expectedValues *velocityContextProviderResourceModel) {
	if !expectedValues.HttpServletExtensionName.IsUnknown() {
		state.HttpServletExtensionName = expectedValues.HttpServletExtensionName
	}
}

// Create any update operations necessary to make the state match the plan
func createVelocityContextProviderOperations(plan velocityContextProviderResourceModel, state velocityContextProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RequestTool, state.RequestTool, "request-tool")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SessionTool, state.SessionTool, "session-tool")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ApplicationTool, state.ApplicationTool, "application-tool")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HttpMethod, state.HttpMethod, "http-method")
	operations.AddStringOperationIfNecessary(&ops, plan.ObjectScope, state.ObjectScope, "object-scope")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludedView, state.IncludedView, "included-view")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludedView, state.ExcludedView, "excluded-view")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	return ops
}

// Create a velocity-tools velocity-context-provider
func (r *velocityContextProviderResource) CreateVelocityToolsVelocityContextProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan velocityContextProviderResourceModel) (*velocityContextProviderResourceModel, error) {
	addRequest := client.NewAddVelocityToolsVelocityContextProviderRequest(plan.Name.ValueString(),
		[]client.EnumvelocityToolsVelocityContextProviderSchemaUrn{client.ENUMVELOCITYTOOLSVELOCITYCONTEXTPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VELOCITY_CONTEXT_PROVIDERVELOCITY_TOOLS})
	err := addOptionalVelocityToolsVelocityContextProviderFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Velocity Context Provider", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VelocityContextProviderApi.AddVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.HttpServletExtensionName.ValueString())
	apiAddRequest = apiAddRequest.AddVelocityContextProviderRequest(
		client.AddVelocityToolsVelocityContextProviderRequestAsAddVelocityContextProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VelocityContextProviderApi.AddVelocityContextProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Velocity Context Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state velocityContextProviderResourceModel
	readVelocityToolsVelocityContextProviderResponse(ctx, addResponse.VelocityToolsVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party velocity-context-provider
func (r *velocityContextProviderResource) CreateThirdPartyVelocityContextProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan velocityContextProviderResourceModel) (*velocityContextProviderResourceModel, error) {
	addRequest := client.NewAddThirdPartyVelocityContextProviderRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyVelocityContextProviderSchemaUrn{client.ENUMTHIRDPARTYVELOCITYCONTEXTPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0VELOCITY_CONTEXT_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	err := addOptionalThirdPartyVelocityContextProviderFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Velocity Context Provider", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.VelocityContextProviderApi.AddVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.HttpServletExtensionName.ValueString())
	apiAddRequest = apiAddRequest.AddVelocityContextProviderRequest(
		client.AddThirdPartyVelocityContextProviderRequestAsAddVelocityContextProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.VelocityContextProviderApi.AddVelocityContextProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Velocity Context Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state velocityContextProviderResourceModel
	readThirdPartyVelocityContextProviderResponse(ctx, addResponse.ThirdPartyVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *velocityContextProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan velocityContextProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *velocityContextProviderResourceModel
	var err error
	if plan.Type.ValueString() == "velocity-tools" {
		state, err = r.CreateVelocityToolsVelocityContextProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyVelocityContextProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	state.setStateValuesNotReturnedByAPI(&plan)
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
func (r *defaultVelocityContextProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan velocityContextProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VelocityContextProviderApi.GetVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Context Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state velocityContextProviderResourceModel
	if readResponse.VelocityToolsVelocityContextProviderResponse != nil {
		readVelocityToolsVelocityContextProviderResponse(ctx, readResponse.VelocityToolsVelocityContextProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomVelocityContextProviderResponse != nil {
		readCustomVelocityContextProviderResponse(ctx, readResponse.CustomVelocityContextProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyVelocityContextProviderResponse != nil {
		readThirdPartyVelocityContextProviderResponse(ctx, readResponse.ThirdPartyVelocityContextProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.VelocityContextProviderApi.UpdateVelocityContextProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.HttpServletExtensionName.ValueString())
	ops := createVelocityContextProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.VelocityContextProviderApi.UpdateVelocityContextProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Velocity Context Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.VelocityToolsVelocityContextProviderResponse != nil {
			readVelocityToolsVelocityContextProviderResponse(ctx, updateResponse.VelocityToolsVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomVelocityContextProviderResponse != nil {
			readCustomVelocityContextProviderResponse(ctx, updateResponse.CustomVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyVelocityContextProviderResponse != nil {
			readThirdPartyVelocityContextProviderResponse(ctx, updateResponse.ThirdPartyVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *velocityContextProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVelocityContextProvider(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultVelocityContextProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVelocityContextProvider(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readVelocityContextProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state velocityContextProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.VelocityContextProviderApi.GetVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Context Provider", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Context Provider", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.VelocityToolsVelocityContextProviderResponse != nil {
		readVelocityToolsVelocityContextProviderResponse(ctx, readResponse.VelocityToolsVelocityContextProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomVelocityContextProviderResponse != nil {
		readCustomVelocityContextProviderResponse(ctx, readResponse.CustomVelocityContextProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyVelocityContextProviderResponse != nil {
		readThirdPartyVelocityContextProviderResponse(ctx, readResponse.ThirdPartyVelocityContextProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *velocityContextProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVelocityContextProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVelocityContextProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVelocityContextProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVelocityContextProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan velocityContextProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state velocityContextProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.VelocityContextProviderApi.UpdateVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.HttpServletExtensionName.ValueString())

	// Determine what update operations are necessary
	ops := createVelocityContextProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.VelocityContextProviderApi.UpdateVelocityContextProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Velocity Context Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.VelocityToolsVelocityContextProviderResponse != nil {
			readVelocityToolsVelocityContextProviderResponse(ctx, updateResponse.VelocityToolsVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomVelocityContextProviderResponse != nil {
			readCustomVelocityContextProviderResponse(ctx, updateResponse.CustomVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyVelocityContextProviderResponse != nil {
			readThirdPartyVelocityContextProviderResponse(ctx, updateResponse.ThirdPartyVelocityContextProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultVelocityContextProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *velocityContextProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state velocityContextProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.VelocityContextProviderApi.DeleteVelocityContextProviderExecute(r.apiClient.VelocityContextProviderApi.DeleteVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.HttpServletExtensionName.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Velocity Context Provider", err, httpResp)
		return
	}
}

func (r *velocityContextProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVelocityContextProvider(ctx, req, resp)
}

func (r *defaultVelocityContextProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVelocityContextProvider(ctx, req, resp)
}

func importVelocityContextProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [http-servlet-extension-name]/[velocity-context-provider-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("http_servlet_extension_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
