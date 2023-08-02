package customloggedstats

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	_ resource.Resource                = &customLoggedStatsResource{}
	_ resource.ResourceWithConfigure   = &customLoggedStatsResource{}
	_ resource.ResourceWithImportState = &customLoggedStatsResource{}
	_ resource.Resource                = &defaultCustomLoggedStatsResource{}
	_ resource.ResourceWithConfigure   = &defaultCustomLoggedStatsResource{}
	_ resource.ResourceWithImportState = &defaultCustomLoggedStatsResource{}
)

// Create a Custom Logged Stats resource
func NewCustomLoggedStatsResource() resource.Resource {
	return &customLoggedStatsResource{}
}

func NewDefaultCustomLoggedStatsResource() resource.Resource {
	return &defaultCustomLoggedStatsResource{}
}

// customLoggedStatsResource is the resource implementation.
type customLoggedStatsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCustomLoggedStatsResource is the resource implementation.
type defaultCustomLoggedStatsResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *customLoggedStatsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_custom_logged_stats"
}

func (r *defaultCustomLoggedStatsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_custom_logged_stats"
}

// Configure adds the provider configured client to the resource.
func (r *customLoggedStatsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultCustomLoggedStatsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type customLoggedStatsResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	LastUpdated            types.String `tfsdk:"last_updated"`
	Notifications          types.Set    `tfsdk:"notifications"`
	RequiredActions        types.Set    `tfsdk:"required_actions"`
	PluginName             types.String `tfsdk:"plugin_name"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
	MonitorObjectclass     types.String `tfsdk:"monitor_objectclass"`
	IncludeFilter          types.String `tfsdk:"include_filter"`
	AttributeToLog         types.Set    `tfsdk:"attribute_to_log"`
	ColumnName             types.Set    `tfsdk:"column_name"`
	StatisticType          types.Set    `tfsdk:"statistic_type"`
	HeaderPrefix           types.String `tfsdk:"header_prefix"`
	HeaderPrefixAttribute  types.String `tfsdk:"header_prefix_attribute"`
	RegexPattern           types.String `tfsdk:"regex_pattern"`
	RegexReplacement       types.String `tfsdk:"regex_replacement"`
	DivideValueBy          types.String `tfsdk:"divide_value_by"`
	DivideValueByAttribute types.String `tfsdk:"divide_value_by_attribute"`
	DecimalFormat          types.String `tfsdk:"decimal_format"`
	NonZeroImpliesNotIdle  types.Bool   `tfsdk:"non_zero_implies_not_idle"`
}

// GetSchema defines the schema for the resource.
func (r *customLoggedStatsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	customLoggedStatsSchema(ctx, req, resp, false)
}

func (r *defaultCustomLoggedStatsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	customLoggedStatsSchema(ctx, req, resp, true)
}

func customLoggedStatsSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Custom Logged Stats.",
		Attributes: map[string]schema.Attribute{
			"plugin_name": schema.StringAttribute{
				Description: "Name of the parent Plugin",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Custom Logged Stats",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Custom Logged Stats object is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"monitor_objectclass": schema.StringAttribute{
				Description: "The objectclass name of the monitor entries to examine for generating these statistics.",
				Required:    true,
			},
			"include_filter": schema.StringAttribute{
				Description: "An optional LDAP filter that can be used restrict which monitor entries are used to produce the output.",
				Optional:    true,
			},
			"attribute_to_log": schema.SetAttribute{
				Description: "Specifies the attributes on the monitor entries that should be included in the output.",
				Required:    true,
				ElementType: types.StringType,
			},
			"column_name": schema.SetAttribute{
				Description: "Optionally, specifies an explicit name for each column header instead of having these names automatically generated from the monitored attribute name.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"statistic_type": schema.SetAttribute{
				Description: "Specifies the type of statistic to include in the output for each monitored attribute.",
				Required:    true,
				ElementType: types.StringType,
			},
			"header_prefix": schema.StringAttribute{
				Description: "An optional prefix that is included in the header before the column name.",
				Optional:    true,
			},
			"header_prefix_attribute": schema.StringAttribute{
				Description: "An optional attribute from the monitor entry that is included as a prefix before the column name in the column header.",
				Optional:    true,
			},
			"regex_pattern": schema.StringAttribute{
				Description: "An optional regular expression pattern, that when used in conjunction with regex-replacement, can alter the value of the attribute being monitored.",
				Optional:    true,
			},
			"regex_replacement": schema.StringAttribute{
				Description: "An optional regular expression replacement value, that when used in conjunction with regex-pattern, can alter the value of the attribute being monitored.",
				Optional:    true,
			},
			"divide_value_by": schema.StringAttribute{
				Description: "An optional floating point value that can be used to scale the resulting value.",
				Optional:    true,
			},
			"divide_value_by_attribute": schema.StringAttribute{
				Description: "An optional property that can scale the resulting value by another attribute in the monitored entry.",
				Optional:    true,
			},
			"decimal_format": schema.StringAttribute{
				Description: "This provides a way to format the monitored attribute value in the output to control the precision for instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"non_zero_implies_not_idle": schema.BoolAttribute{
				Description: "If this property is set to true, then the value of any of the monitored attributes here can contribute to whether an interval is considered \"idle\" by the Periodic Stats Logger.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputed(&schemaDef, []string{"plugin_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators
func (r customLoggedStatsResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.Implies(
			path.MatchRoot("regex_pattern"),
			path.MatchRoot("regex_replacement"),
		),
	}
}

// Add optional fields to create request for custom-logged-stats custom-logged-stats
func addOptionalCustomLoggedStatsFields(ctx context.Context, addRequest *client.AddCustomLoggedStatsRequest, plan customLoggedStatsResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IncludeFilter) {
		addRequest.IncludeFilter = plan.IncludeFilter.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ColumnName) {
		var slice []string
		plan.ColumnName.ElementsAs(ctx, &slice, false)
		addRequest.ColumnName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HeaderPrefix) {
		addRequest.HeaderPrefix = plan.HeaderPrefix.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HeaderPrefixAttribute) {
		addRequest.HeaderPrefixAttribute = plan.HeaderPrefixAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RegexPattern) {
		addRequest.RegexPattern = plan.RegexPattern.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RegexReplacement) {
		addRequest.RegexReplacement = plan.RegexReplacement.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DivideValueBy) {
		addRequest.DivideValueBy = plan.DivideValueBy.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DivideValueByAttribute) {
		addRequest.DivideValueByAttribute = plan.DivideValueByAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DecimalFormat) {
		addRequest.DecimalFormat = plan.DecimalFormat.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.NonZeroImpliesNotIdle) {
		addRequest.NonZeroImpliesNotIdle = plan.NonZeroImpliesNotIdle.ValueBoolPointer()
	}
}

// Read a CustomLoggedStatsResponse object into the model struct
func readCustomLoggedStatsResponse(ctx context.Context, r *client.CustomLoggedStatsResponse, state *customLoggedStatsResourceModel, expectedValues *customLoggedStatsResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.MonitorObjectclass = types.StringValue(r.MonitorObjectclass)
	state.IncludeFilter = internaltypes.StringTypeOrNil(r.IncludeFilter, internaltypes.IsEmptyString(expectedValues.IncludeFilter))
	state.AttributeToLog = internaltypes.GetStringSet(r.AttributeToLog)
	state.ColumnName = internaltypes.GetStringSet(r.ColumnName)
	state.StatisticType = internaltypes.GetStringSet(
		client.StringSliceEnumcustomLoggedStatsStatisticTypeProp(r.StatisticType))
	state.HeaderPrefix = internaltypes.StringTypeOrNil(r.HeaderPrefix, internaltypes.IsEmptyString(expectedValues.HeaderPrefix))
	state.HeaderPrefixAttribute = internaltypes.StringTypeOrNil(r.HeaderPrefixAttribute, internaltypes.IsEmptyString(expectedValues.HeaderPrefixAttribute))
	state.RegexPattern = internaltypes.StringTypeOrNil(r.RegexPattern, internaltypes.IsEmptyString(expectedValues.RegexPattern))
	state.RegexReplacement = internaltypes.StringTypeOrNil(r.RegexReplacement, internaltypes.IsEmptyString(expectedValues.RegexReplacement))
	state.DivideValueBy = internaltypes.StringTypeOrNil(r.DivideValueBy, internaltypes.IsEmptyString(expectedValues.DivideValueBy))
	state.DivideValueByAttribute = internaltypes.StringTypeOrNil(r.DivideValueByAttribute, internaltypes.IsEmptyString(expectedValues.DivideValueByAttribute))
	state.DecimalFormat = internaltypes.StringTypeOrNil(r.DecimalFormat, internaltypes.IsEmptyString(expectedValues.DecimalFormat))
	state.NonZeroImpliesNotIdle = internaltypes.BoolTypeOrNil(r.NonZeroImpliesNotIdle)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *customLoggedStatsResourceModel) setStateValuesNotReturnedByAPI(expectedValues *customLoggedStatsResourceModel) {
	if !expectedValues.PluginName.IsUnknown() {
		state.PluginName = expectedValues.PluginName
	}
}

// Create any update operations necessary to make the state match the plan
func createCustomLoggedStatsOperations(plan customLoggedStatsResourceModel, state customLoggedStatsResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.MonitorObjectclass, state.MonitorObjectclass, "monitor-objectclass")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeFilter, state.IncludeFilter, "include-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AttributeToLog, state.AttributeToLog, "attribute-to-log")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ColumnName, state.ColumnName, "column-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.StatisticType, state.StatisticType, "statistic-type")
	operations.AddStringOperationIfNecessary(&ops, plan.HeaderPrefix, state.HeaderPrefix, "header-prefix")
	operations.AddStringOperationIfNecessary(&ops, plan.HeaderPrefixAttribute, state.HeaderPrefixAttribute, "header-prefix-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.RegexPattern, state.RegexPattern, "regex-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.RegexReplacement, state.RegexReplacement, "regex-replacement")
	operations.AddStringOperationIfNecessary(&ops, plan.DivideValueBy, state.DivideValueBy, "divide-value-by")
	operations.AddStringOperationIfNecessary(&ops, plan.DivideValueByAttribute, state.DivideValueByAttribute, "divide-value-by-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.DecimalFormat, state.DecimalFormat, "decimal-format")
	operations.AddBoolOperationIfNecessary(&ops, plan.NonZeroImpliesNotIdle, state.NonZeroImpliesNotIdle, "non-zero-implies-not-idle")
	return ops
}

// Create a custom-logged-stats custom-logged-stats
func (r *customLoggedStatsResource) CreateCustomLoggedStats(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan customLoggedStatsResourceModel) (*customLoggedStatsResourceModel, error) {
	var AttributeToLogSlice []string
	plan.AttributeToLog.ElementsAs(ctx, &AttributeToLogSlice, false)
	var StatisticTypeSlice []client.EnumcustomLoggedStatsStatisticTypeProp
	plan.StatisticType.ElementsAs(ctx, &StatisticTypeSlice, false)
	addRequest := client.NewAddCustomLoggedStatsRequest(plan.Name.ValueString(),
		[]client.EnumcustomLoggedStatsSchemaUrn{client.ENUMCUSTOMLOGGEDSTATSSCHEMAURN_STATSCUSTOM},
		plan.MonitorObjectclass.ValueString(),
		AttributeToLogSlice,
		StatisticTypeSlice)
	addOptionalCustomLoggedStatsFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CustomLoggedStatsApi.AddCustomLoggedStats(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.PluginName.ValueString())
	apiAddRequest = apiAddRequest.AddCustomLoggedStatsRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.CustomLoggedStatsApi.AddCustomLoggedStatsExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Custom Logged Stats", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state customLoggedStatsResourceModel
	readCustomLoggedStatsResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *customLoggedStatsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan customLoggedStatsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateCustomLoggedStats(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultCustomLoggedStatsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan customLoggedStatsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CustomLoggedStatsApi.GetCustomLoggedStats(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.PluginName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Custom Logged Stats", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state customLoggedStatsResourceModel
	readCustomLoggedStatsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CustomLoggedStatsApi.UpdateCustomLoggedStats(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.PluginName.ValueString())
	ops := createCustomLoggedStatsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CustomLoggedStatsApi.UpdateCustomLoggedStatsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Custom Logged Stats", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCustomLoggedStatsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *customLoggedStatsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCustomLoggedStats(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCustomLoggedStatsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCustomLoggedStats(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCustomLoggedStats(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state customLoggedStatsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CustomLoggedStatsApi.GetCustomLoggedStats(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.PluginName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Custom Logged Stats", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCustomLoggedStatsResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *customLoggedStatsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCustomLoggedStats(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCustomLoggedStatsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCustomLoggedStats(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCustomLoggedStats(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan customLoggedStatsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state customLoggedStatsResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CustomLoggedStatsApi.UpdateCustomLoggedStats(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.PluginName.ValueString())

	// Determine what update operations are necessary
	ops := createCustomLoggedStatsOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CustomLoggedStatsApi.UpdateCustomLoggedStatsExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Custom Logged Stats", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCustomLoggedStatsResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCustomLoggedStatsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *customLoggedStatsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state customLoggedStatsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CustomLoggedStatsApi.DeleteCustomLoggedStatsExecute(r.apiClient.CustomLoggedStatsApi.DeleteCustomLoggedStats(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.PluginName.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Custom Logged Stats", err, httpResp)
		return
	}
}

func (r *customLoggedStatsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCustomLoggedStats(ctx, req, resp)
}

func (r *defaultCustomLoggedStatsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCustomLoggedStats(ctx, req, resp)
}

func importCustomLoggedStats(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [plugin-name]/[custom-logged-stats-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("plugin_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
