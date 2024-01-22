package changesubscriptionhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &changeSubscriptionHandlerResource{}
	_ resource.ResourceWithConfigure   = &changeSubscriptionHandlerResource{}
	_ resource.ResourceWithImportState = &changeSubscriptionHandlerResource{}
	_ resource.Resource                = &defaultChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultChangeSubscriptionHandlerResource{}
)

// Create a Change Subscription Handler resource
func NewChangeSubscriptionHandlerResource() resource.Resource {
	return &changeSubscriptionHandlerResource{}
}

func NewDefaultChangeSubscriptionHandlerResource() resource.Resource {
	return &defaultChangeSubscriptionHandlerResource{}
}

// changeSubscriptionHandlerResource is the resource implementation.
type changeSubscriptionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultChangeSubscriptionHandlerResource is the resource implementation.
type defaultChangeSubscriptionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *changeSubscriptionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_change_subscription_handler"
}

func (r *defaultChangeSubscriptionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_change_subscription_handler"
}

// Configure adds the provider configured client to the resource.
func (r *changeSubscriptionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultChangeSubscriptionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type changeSubscriptionHandlerResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	Type               types.String `tfsdk:"type"`
	ExtensionClass     types.String `tfsdk:"extension_class"`
	ExtensionArgument  types.Set    `tfsdk:"extension_argument"`
	LogFile            types.String `tfsdk:"log_file"`
	ScriptClass        types.String `tfsdk:"script_class"`
	ScriptArgument     types.Set    `tfsdk:"script_argument"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ChangeSubscription types.Set    `tfsdk:"change_subscription"`
}

// GetSchema defines the schema for the resource.
func (r *changeSubscriptionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	changeSubscriptionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultChangeSubscriptionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	changeSubscriptionHandlerSchema(ctx, req, resp, true)
}

func changeSubscriptionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Change Subscription Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Change Subscription Handler resource. Options are ['groovy-scripted', 'logging', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"groovy-scripted", "logging", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Change Subscription Handler.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Change Subscription Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"log_file": schema.StringAttribute{
				Description: "Specifies the log file in which the change notification messages will be written.",
				Optional:    true,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Change Subscription Handler.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Change Subscription Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Change Subscription Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this change subscription handler is enabled within the server.",
				Required:    true,
			},
			"change_subscription": schema.SetAttribute{
				Description: "The set of change subscriptions for which this change subscription handler should be notified. If no values are provided then it will be notified for all change subscriptions defined in the server.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *changeSubscriptionHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel changeSubscriptionHandlerResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for logging type
	if resourceType == "logging" {
		if !internaltypes.IsDefined(configModel.LogFile) {
			defaultVal := types.StringValue("logs/change-notifications.log")
			if !planModel.LogFile.Equal(defaultVal) {
				planModel.LogFile = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (model *changeSubscriptionHandlerResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "groovy-scripted" {
		model.LogFile = types.StringNull()
	}
	if resourceType == "third-party" {
		model.LogFile = types.StringNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsChangeSubscriptionHandler() []resource.ConfigValidator {
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
			path.MatchRoot("log_file"),
			path.MatchRoot("type"),
			[]string{"logging"},
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
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted",
			[]path.Expression{path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r changeSubscriptionHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsChangeSubscriptionHandler()
}

// Add config validators
func (r defaultChangeSubscriptionHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsChangeSubscriptionHandler()
}

// Add optional fields to create request for groovy-scripted change-subscription-handler
func addOptionalGroovyScriptedChangeSubscriptionHandlerFields(ctx context.Context, addRequest *client.AddGroovyScriptedChangeSubscriptionHandlerRequest, plan changeSubscriptionHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ChangeSubscription) {
		var slice []string
		plan.ChangeSubscription.ElementsAs(ctx, &slice, false)
		addRequest.ChangeSubscription = slice
	}
}

// Add optional fields to create request for logging change-subscription-handler
func addOptionalLoggingChangeSubscriptionHandlerFields(ctx context.Context, addRequest *client.AddLoggingChangeSubscriptionHandlerRequest, plan changeSubscriptionHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFile) {
		addRequest.LogFile = plan.LogFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ChangeSubscription) {
		var slice []string
		plan.ChangeSubscription.ElementsAs(ctx, &slice, false)
		addRequest.ChangeSubscription = slice
	}
}

// Add optional fields to create request for third-party change-subscription-handler
func addOptionalThirdPartyChangeSubscriptionHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyChangeSubscriptionHandlerRequest, plan changeSubscriptionHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ChangeSubscription) {
		var slice []string
		plan.ChangeSubscription.ElementsAs(ctx, &slice, false)
		addRequest.ChangeSubscription = slice
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateChangeSubscriptionHandlerUnknownValues(model *changeSubscriptionHandlerResourceModel) {
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *changeSubscriptionHandlerResourceModel) populateAllComputedStringAttributes() {
	if model.LogFile.IsUnknown() || model.LogFile.IsNull() {
		model.LogFile = types.StringValue("")
	}
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

// Read a GroovyScriptedChangeSubscriptionHandlerResponse object into the model struct
func readGroovyScriptedChangeSubscriptionHandlerResponse(ctx context.Context, r *client.GroovyScriptedChangeSubscriptionHandlerResponse, state *changeSubscriptionHandlerResourceModel, expectedValues *changeSubscriptionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateChangeSubscriptionHandlerUnknownValues(state)
}

// Read a LoggingChangeSubscriptionHandlerResponse object into the model struct
func readLoggingChangeSubscriptionHandlerResponse(ctx context.Context, r *client.LoggingChangeSubscriptionHandlerResponse, state *changeSubscriptionHandlerResourceModel, expectedValues *changeSubscriptionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("logging")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateChangeSubscriptionHandlerUnknownValues(state)
}

// Read a ThirdPartyChangeSubscriptionHandlerResponse object into the model struct
func readThirdPartyChangeSubscriptionHandlerResponse(ctx context.Context, r *client.ThirdPartyChangeSubscriptionHandlerResponse, state *changeSubscriptionHandlerResourceModel, expectedValues *changeSubscriptionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateChangeSubscriptionHandlerUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createChangeSubscriptionHandlerOperations(plan changeSubscriptionHandlerResourceModel, state changeSubscriptionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangeSubscription, state.ChangeSubscription, "change-subscription")
	return ops
}

// Create a groovy-scripted change-subscription-handler
func (r *changeSubscriptionHandlerResource) CreateGroovyScriptedChangeSubscriptionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan changeSubscriptionHandlerResourceModel) (*changeSubscriptionHandlerResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedChangeSubscriptionHandlerRequest([]client.EnumgroovyScriptedChangeSubscriptionHandlerSchemaUrn{client.ENUMGROOVYSCRIPTEDCHANGESUBSCRIPTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CHANGE_SUBSCRIPTION_HANDLERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalGroovyScriptedChangeSubscriptionHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ChangeSubscriptionHandlerAPI.AddChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddChangeSubscriptionHandlerRequest(
		client.AddGroovyScriptedChangeSubscriptionHandlerRequestAsAddChangeSubscriptionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.AddChangeSubscriptionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Change Subscription Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state changeSubscriptionHandlerResourceModel
	readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, addResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a logging change-subscription-handler
func (r *changeSubscriptionHandlerResource) CreateLoggingChangeSubscriptionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan changeSubscriptionHandlerResourceModel) (*changeSubscriptionHandlerResourceModel, error) {
	addRequest := client.NewAddLoggingChangeSubscriptionHandlerRequest([]client.EnumloggingChangeSubscriptionHandlerSchemaUrn{client.ENUMLOGGINGCHANGESUBSCRIPTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CHANGE_SUBSCRIPTION_HANDLERLOGGING},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalLoggingChangeSubscriptionHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ChangeSubscriptionHandlerAPI.AddChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddChangeSubscriptionHandlerRequest(
		client.AddLoggingChangeSubscriptionHandlerRequestAsAddChangeSubscriptionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.AddChangeSubscriptionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Change Subscription Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state changeSubscriptionHandlerResourceModel
	readLoggingChangeSubscriptionHandlerResponse(ctx, addResponse.LoggingChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party change-subscription-handler
func (r *changeSubscriptionHandlerResource) CreateThirdPartyChangeSubscriptionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan changeSubscriptionHandlerResourceModel) (*changeSubscriptionHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyChangeSubscriptionHandlerRequest([]client.EnumthirdPartyChangeSubscriptionHandlerSchemaUrn{client.ENUMTHIRDPARTYCHANGESUBSCRIPTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CHANGE_SUBSCRIPTION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	addOptionalThirdPartyChangeSubscriptionHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ChangeSubscriptionHandlerAPI.AddChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddChangeSubscriptionHandlerRequest(
		client.AddThirdPartyChangeSubscriptionHandlerRequestAsAddChangeSubscriptionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.AddChangeSubscriptionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Change Subscription Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state changeSubscriptionHandlerResourceModel
	readThirdPartyChangeSubscriptionHandlerResponse(ctx, addResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *changeSubscriptionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan changeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *changeSubscriptionHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedChangeSubscriptionHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "logging" {
		state, err = r.CreateLoggingChangeSubscriptionHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyChangeSubscriptionHandler(ctx, req, resp, plan)
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
func (r *defaultChangeSubscriptionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan changeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state changeSubscriptionHandlerResourceModel
	if readResponse.GroovyScriptedChangeSubscriptionHandlerResponse != nil {
		readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, readResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LoggingChangeSubscriptionHandlerResponse != nil {
		readLoggingChangeSubscriptionHandlerResponse(ctx, readResponse.LoggingChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyChangeSubscriptionHandlerResponse != nil {
		readThirdPartyChangeSubscriptionHandlerResponse(ctx, readResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ChangeSubscriptionHandlerAPI.UpdateChangeSubscriptionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createChangeSubscriptionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.UpdateChangeSubscriptionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Change Subscription Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GroovyScriptedChangeSubscriptionHandlerResponse != nil {
			readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, updateResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LoggingChangeSubscriptionHandlerResponse != nil {
			readLoggingChangeSubscriptionHandlerResponse(ctx, updateResponse.LoggingChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyChangeSubscriptionHandlerResponse != nil {
			readThirdPartyChangeSubscriptionHandlerResponse(ctx, updateResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *changeSubscriptionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultChangeSubscriptionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readChangeSubscriptionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state changeSubscriptionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ChangeSubscriptionHandlerAPI.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Change Subscription Handler", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Change Subscription Handler", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GroovyScriptedChangeSubscriptionHandlerResponse != nil {
		readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, readResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LoggingChangeSubscriptionHandlerResponse != nil {
		readLoggingChangeSubscriptionHandlerResponse(ctx, readResponse.LoggingChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyChangeSubscriptionHandlerResponse != nil {
		readThirdPartyChangeSubscriptionHandlerResponse(ctx, readResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *changeSubscriptionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultChangeSubscriptionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateChangeSubscriptionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan changeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state changeSubscriptionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ChangeSubscriptionHandlerAPI.UpdateChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createChangeSubscriptionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ChangeSubscriptionHandlerAPI.UpdateChangeSubscriptionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Change Subscription Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GroovyScriptedChangeSubscriptionHandlerResponse != nil {
			readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, updateResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LoggingChangeSubscriptionHandlerResponse != nil {
			readLoggingChangeSubscriptionHandlerResponse(ctx, updateResponse.LoggingChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyChangeSubscriptionHandlerResponse != nil {
			readThirdPartyChangeSubscriptionHandlerResponse(ctx, updateResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultChangeSubscriptionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *changeSubscriptionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state changeSubscriptionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ChangeSubscriptionHandlerAPI.DeleteChangeSubscriptionHandlerExecute(r.apiClient.ChangeSubscriptionHandlerAPI.DeleteChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Change Subscription Handler", err, httpResp)
		return
	}
}

func (r *changeSubscriptionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importChangeSubscriptionHandler(ctx, req, resp)
}

func (r *defaultChangeSubscriptionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importChangeSubscriptionHandler(ctx, req, resp)
}

func importChangeSubscriptionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
