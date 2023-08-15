package uncachedentrycriteria

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
	_ resource.Resource                = &uncachedEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &uncachedEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &uncachedEntryCriteriaResource{}
	_ resource.Resource                = &defaultUncachedEntryCriteriaResource{}
	_ resource.ResourceWithConfigure   = &defaultUncachedEntryCriteriaResource{}
	_ resource.ResourceWithImportState = &defaultUncachedEntryCriteriaResource{}
)

// Create a Uncached Entry Criteria resource
func NewUncachedEntryCriteriaResource() resource.Resource {
	return &uncachedEntryCriteriaResource{}
}

func NewDefaultUncachedEntryCriteriaResource() resource.Resource {
	return &defaultUncachedEntryCriteriaResource{}
}

// uncachedEntryCriteriaResource is the resource implementation.
type uncachedEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultUncachedEntryCriteriaResource is the resource implementation.
type defaultUncachedEntryCriteriaResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *uncachedEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uncached_entry_criteria"
}

func (r *defaultUncachedEntryCriteriaResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_uncached_entry_criteria"
}

// Configure adds the provider configured client to the resource.
func (r *uncachedEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultUncachedEntryCriteriaResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type uncachedEntryCriteriaResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	ScriptClass                     types.String `tfsdk:"script_class"`
	ScriptArgument                  types.Set    `tfsdk:"script_argument"`
	Filter                          types.String `tfsdk:"filter"`
	FilterIdentifiesUncachedEntries types.Bool   `tfsdk:"filter_identifies_uncached_entries"`
	AccessTimeThreshold             types.String `tfsdk:"access_time_threshold"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *uncachedEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	uncachedEntryCriteriaSchema(ctx, req, resp, false)
}

func (r *defaultUncachedEntryCriteriaResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	uncachedEntryCriteriaSchema(ctx, req, resp, true)
}

func uncachedEntryCriteriaSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Uncached Entry Criteria.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Uncached Entry Criteria resource. Options are ['default', 'last-access-time', 'filter-based', 'groovy-scripted', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"default", "last-access-time", "filter-based", "groovy-scripted", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Uncached Entry Criteria.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Uncached Entry Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Uncached Entry Criteria.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Uncached Entry Criteria. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"filter": schema.StringAttribute{
				Description: "Specifies the search filter that should be used to differentiate entries into cached and uncached sets.",
				Optional:    true,
			},
			"filter_identifies_uncached_entries": schema.BoolAttribute{
				Description: "Indicates whether the associated filter identifies those entries which should be stored in the uncached-id2entry database (if true) or entries which should be stored in the id2entry database (if false).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"access_time_threshold": schema.StringAttribute{
				Description: "Specifies the maximum length of time that has passed since an entry was last accessed that it should still be included in the id2entry database. Entries that have not been accessed in more than this length of time may be written into the uncached-id2entry database.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Uncached Entry Criteria",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Uncached Entry Criteria is enabled for use in the server.",
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
func configValidatorsUncachedEntryCriteria() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("filter"),
			path.MatchRoot("type"),
			[]string{"filter-based"},
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
			path.MatchRoot("access_time_threshold"),
			path.MatchRoot("type"),
			[]string{"last-access-time"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("filter_identifies_uncached_entries"),
			path.MatchRoot("type"),
			[]string{"filter-based"},
		),
	}
}

// Add config validators
func (r uncachedEntryCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsUncachedEntryCriteria()
}

// Add config validators
func (r defaultUncachedEntryCriteriaResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsUncachedEntryCriteria()
}

// Add optional fields to create request for default uncached-entry-criteria
func addOptionalDefaultUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddDefaultUncachedEntryCriteriaRequest, plan uncachedEntryCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for last-access-time uncached-entry-criteria
func addOptionalLastAccessTimeUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddLastAccessTimeUncachedEntryCriteriaRequest, plan uncachedEntryCriteriaResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for filter-based uncached-entry-criteria
func addOptionalFilterBasedUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddFilterBasedUncachedEntryCriteriaRequest, plan uncachedEntryCriteriaResourceModel) {
	if internaltypes.IsDefined(plan.FilterIdentifiesUncachedEntries) {
		addRequest.FilterIdentifiesUncachedEntries = plan.FilterIdentifiesUncachedEntries.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for groovy-scripted uncached-entry-criteria
func addOptionalGroovyScriptedUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddGroovyScriptedUncachedEntryCriteriaRequest, plan uncachedEntryCriteriaResourceModel) {
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

// Add optional fields to create request for third-party uncached-entry-criteria
func addOptionalThirdPartyUncachedEntryCriteriaFields(ctx context.Context, addRequest *client.AddThirdPartyUncachedEntryCriteriaRequest, plan uncachedEntryCriteriaResourceModel) {
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
func populateUncachedEntryCriteriaUnknownValues(ctx context.Context, model *uncachedEntryCriteriaResourceModel) {
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
}

// Read a DefaultUncachedEntryCriteriaResponse object into the model struct
func readDefaultUncachedEntryCriteriaResponse(ctx context.Context, r *client.DefaultUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaResourceModel, expectedValues *uncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("default")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedEntryCriteriaUnknownValues(ctx, state)
}

// Read a LastAccessTimeUncachedEntryCriteriaResponse object into the model struct
func readLastAccessTimeUncachedEntryCriteriaResponse(ctx context.Context, r *client.LastAccessTimeUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaResourceModel, expectedValues *uncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("last-access-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccessTimeThreshold = types.StringValue(r.AccessTimeThreshold)
	config.CheckMismatchedPDFormattedAttributes("access_time_threshold",
		expectedValues.AccessTimeThreshold, state.AccessTimeThreshold, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedEntryCriteriaUnknownValues(ctx, state)
}

// Read a FilterBasedUncachedEntryCriteriaResponse object into the model struct
func readFilterBasedUncachedEntryCriteriaResponse(ctx context.Context, r *client.FilterBasedUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaResourceModel, expectedValues *uncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("filter-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Filter = types.StringValue(r.Filter)
	state.FilterIdentifiesUncachedEntries = types.BoolValue(r.FilterIdentifiesUncachedEntries)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedEntryCriteriaUnknownValues(ctx, state)
}

// Read a GroovyScriptedUncachedEntryCriteriaResponse object into the model struct
func readGroovyScriptedUncachedEntryCriteriaResponse(ctx context.Context, r *client.GroovyScriptedUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaResourceModel, expectedValues *uncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedEntryCriteriaUnknownValues(ctx, state)
}

// Read a ThirdPartyUncachedEntryCriteriaResponse object into the model struct
func readThirdPartyUncachedEntryCriteriaResponse(ctx context.Context, r *client.ThirdPartyUncachedEntryCriteriaResponse, state *uncachedEntryCriteriaResourceModel, expectedValues *uncachedEntryCriteriaResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateUncachedEntryCriteriaUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createUncachedEntryCriteriaOperations(plan uncachedEntryCriteriaResourceModel, state uncachedEntryCriteriaResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Filter, state.Filter, "filter")
	operations.AddBoolOperationIfNecessary(&ops, plan.FilterIdentifiesUncachedEntries, state.FilterIdentifiesUncachedEntries, "filter-identifies-uncached-entries")
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTimeThreshold, state.AccessTimeThreshold, "access-time-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a default uncached-entry-criteria
func (r *uncachedEntryCriteriaResource) CreateDefaultUncachedEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedEntryCriteriaResourceModel) (*uncachedEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddDefaultUncachedEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumdefaultUncachedEntryCriteriaSchemaUrn{client.ENUMDEFAULTUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIADEFAULT},
		plan.Enabled.ValueBool())
	addOptionalDefaultUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddDefaultUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedEntryCriteriaResourceModel
	readDefaultUncachedEntryCriteriaResponse(ctx, addResponse.DefaultUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a last-access-time uncached-entry-criteria
func (r *uncachedEntryCriteriaResource) CreateLastAccessTimeUncachedEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedEntryCriteriaResourceModel) (*uncachedEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddLastAccessTimeUncachedEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumlastAccessTimeUncachedEntryCriteriaSchemaUrn{client.ENUMLASTACCESSTIMEUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIALAST_ACCESS_TIME},
		plan.AccessTimeThreshold.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalLastAccessTimeUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddLastAccessTimeUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedEntryCriteriaResourceModel
	readLastAccessTimeUncachedEntryCriteriaResponse(ctx, addResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a filter-based uncached-entry-criteria
func (r *uncachedEntryCriteriaResource) CreateFilterBasedUncachedEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedEntryCriteriaResourceModel) (*uncachedEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddFilterBasedUncachedEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumfilterBasedUncachedEntryCriteriaSchemaUrn{client.ENUMFILTERBASEDUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIAFILTER_BASED},
		plan.Filter.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFilterBasedUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddFilterBasedUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedEntryCriteriaResourceModel
	readFilterBasedUncachedEntryCriteriaResponse(ctx, addResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted uncached-entry-criteria
func (r *uncachedEntryCriteriaResource) CreateGroovyScriptedUncachedEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedEntryCriteriaResourceModel) (*uncachedEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedUncachedEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedUncachedEntryCriteriaSchemaUrn{client.ENUMGROOVYSCRIPTEDUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIAGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddGroovyScriptedUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedEntryCriteriaResourceModel
	readGroovyScriptedUncachedEntryCriteriaResponse(ctx, addResponse.GroovyScriptedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party uncached-entry-criteria
func (r *uncachedEntryCriteriaResource) CreateThirdPartyUncachedEntryCriteria(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan uncachedEntryCriteriaResourceModel) (*uncachedEntryCriteriaResourceModel, error) {
	addRequest := client.NewAddThirdPartyUncachedEntryCriteriaRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyUncachedEntryCriteriaSchemaUrn{client.ENUMTHIRDPARTYUNCACHEDENTRYCRITERIASCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0UNCACHED_ENTRY_CRITERIATHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyUncachedEntryCriteriaFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddUncachedEntryCriteriaRequest(
		client.AddThirdPartyUncachedEntryCriteriaRequestAsAddUncachedEntryCriteriaRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.AddUncachedEntryCriteriaExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Uncached Entry Criteria", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state uncachedEntryCriteriaResourceModel
	readThirdPartyUncachedEntryCriteriaResponse(ctx, addResponse.ThirdPartyUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *uncachedEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan uncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *uncachedEntryCriteriaResourceModel
	var err error
	if plan.Type.ValueString() == "default" {
		state, err = r.CreateDefaultUncachedEntryCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "last-access-time" {
		state, err = r.CreateLastAccessTimeUncachedEntryCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "filter-based" {
		state, err = r.CreateFilterBasedUncachedEntryCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedUncachedEntryCriteria(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyUncachedEntryCriteria(ctx, req, resp, plan)
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
func (r *defaultUncachedEntryCriteriaResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan uncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Entry Criteria", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state uncachedEntryCriteriaResourceModel
	if readResponse.DefaultUncachedEntryCriteriaResponse != nil {
		readDefaultUncachedEntryCriteriaResponse(ctx, readResponse.DefaultUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LastAccessTimeUncachedEntryCriteriaResponse != nil {
		readLastAccessTimeUncachedEntryCriteriaResponse(ctx, readResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FilterBasedUncachedEntryCriteriaResponse != nil {
		readFilterBasedUncachedEntryCriteriaResponse(ctx, readResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedUncachedEntryCriteriaResponse != nil {
		readGroovyScriptedUncachedEntryCriteriaResponse(ctx, readResponse.GroovyScriptedUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyUncachedEntryCriteriaResponse != nil {
		readThirdPartyUncachedEntryCriteriaResponse(ctx, readResponse.ThirdPartyUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteria(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createUncachedEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Uncached Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DefaultUncachedEntryCriteriaResponse != nil {
			readDefaultUncachedEntryCriteriaResponse(ctx, updateResponse.DefaultUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LastAccessTimeUncachedEntryCriteriaResponse != nil {
			readLastAccessTimeUncachedEntryCriteriaResponse(ctx, updateResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FilterBasedUncachedEntryCriteriaResponse != nil {
			readFilterBasedUncachedEntryCriteriaResponse(ctx, updateResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedUncachedEntryCriteriaResponse != nil {
			readGroovyScriptedUncachedEntryCriteriaResponse(ctx, updateResponse.GroovyScriptedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyUncachedEntryCriteriaResponse != nil {
			readThirdPartyUncachedEntryCriteriaResponse(ctx, updateResponse.ThirdPartyUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *uncachedEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultUncachedEntryCriteriaResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readUncachedEntryCriteria(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state uncachedEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.UncachedEntryCriteriaApi.GetUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Entry Criteria", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Uncached Entry Criteria", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.DefaultUncachedEntryCriteriaResponse != nil {
		readDefaultUncachedEntryCriteriaResponse(ctx, readResponse.DefaultUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LastAccessTimeUncachedEntryCriteriaResponse != nil {
		readLastAccessTimeUncachedEntryCriteriaResponse(ctx, readResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FilterBasedUncachedEntryCriteriaResponse != nil {
		readFilterBasedUncachedEntryCriteriaResponse(ctx, readResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedUncachedEntryCriteriaResponse != nil {
		readGroovyScriptedUncachedEntryCriteriaResponse(ctx, readResponse.GroovyScriptedUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyUncachedEntryCriteriaResponse != nil {
		readThirdPartyUncachedEntryCriteriaResponse(ctx, readResponse.ThirdPartyUncachedEntryCriteriaResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *uncachedEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUncachedEntryCriteriaResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUncachedEntryCriteria(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateUncachedEntryCriteria(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan uncachedEntryCriteriaResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state uncachedEntryCriteriaResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createUncachedEntryCriteriaOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.UncachedEntryCriteriaApi.UpdateUncachedEntryCriteriaExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Uncached Entry Criteria", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.DefaultUncachedEntryCriteriaResponse != nil {
			readDefaultUncachedEntryCriteriaResponse(ctx, updateResponse.DefaultUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LastAccessTimeUncachedEntryCriteriaResponse != nil {
			readLastAccessTimeUncachedEntryCriteriaResponse(ctx, updateResponse.LastAccessTimeUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FilterBasedUncachedEntryCriteriaResponse != nil {
			readFilterBasedUncachedEntryCriteriaResponse(ctx, updateResponse.FilterBasedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedUncachedEntryCriteriaResponse != nil {
			readGroovyScriptedUncachedEntryCriteriaResponse(ctx, updateResponse.GroovyScriptedUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyUncachedEntryCriteriaResponse != nil {
			readThirdPartyUncachedEntryCriteriaResponse(ctx, updateResponse.ThirdPartyUncachedEntryCriteriaResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultUncachedEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *uncachedEntryCriteriaResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state uncachedEntryCriteriaResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.UncachedEntryCriteriaApi.DeleteUncachedEntryCriteriaExecute(r.apiClient.UncachedEntryCriteriaApi.DeleteUncachedEntryCriteria(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Uncached Entry Criteria", err, httpResp)
		return
	}
}

func (r *uncachedEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUncachedEntryCriteria(ctx, req, resp)
}

func (r *defaultUncachedEntryCriteriaResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUncachedEntryCriteria(ctx, req, resp)
}

func importUncachedEntryCriteria(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
