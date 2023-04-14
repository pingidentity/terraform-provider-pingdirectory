package gauge

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &indicatorGaugeResource{}
	_ resource.ResourceWithConfigure   = &indicatorGaugeResource{}
	_ resource.ResourceWithImportState = &indicatorGaugeResource{}
	_ resource.Resource                = &defaultIndicatorGaugeResource{}
	_ resource.ResourceWithConfigure   = &defaultIndicatorGaugeResource{}
	_ resource.ResourceWithImportState = &defaultIndicatorGaugeResource{}
)

// Create a Indicator Gauge resource
func NewIndicatorGaugeResource() resource.Resource {
	return &indicatorGaugeResource{}
}

func NewDefaultIndicatorGaugeResource() resource.Resource {
	return &defaultIndicatorGaugeResource{}
}

// indicatorGaugeResource is the resource implementation.
type indicatorGaugeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultIndicatorGaugeResource is the resource implementation.
type defaultIndicatorGaugeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *indicatorGaugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_indicator_gauge"
}

func (r *defaultIndicatorGaugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_indicator_gauge"
}

// Configure adds the provider configured client to the resource.
func (r *indicatorGaugeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultIndicatorGaugeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type indicatorGaugeResourceModel struct {
	Id                             types.String `tfsdk:"id"`
	LastUpdated                    types.String `tfsdk:"last_updated"`
	Notifications                  types.Set    `tfsdk:"notifications"`
	RequiredActions                types.Set    `tfsdk:"required_actions"`
	GaugeDataSource                types.String `tfsdk:"gauge_data_source"`
	CriticalValue                  types.String `tfsdk:"critical_value"`
	MajorValue                     types.String `tfsdk:"major_value"`
	MinorValue                     types.String `tfsdk:"minor_value"`
	WarningValue                   types.String `tfsdk:"warning_value"`
	Description                    types.String `tfsdk:"description"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	OverrideSeverity               types.String `tfsdk:"override_severity"`
	AlertLevel                     types.String `tfsdk:"alert_level"`
	UpdateInterval                 types.String `tfsdk:"update_interval"`
	SamplesPerUpdateInterval       types.Int64  `tfsdk:"samples_per_update_interval"`
	IncludeResource                types.Set    `tfsdk:"include_resource"`
	ExcludeResource                types.Set    `tfsdk:"exclude_resource"`
	ServerUnavailableSeverityLevel types.String `tfsdk:"server_unavailable_severity_level"`
	ServerDegradedSeverityLevel    types.String `tfsdk:"server_degraded_severity_level"`
}

// GetSchema defines the schema for the resource.
func (r *indicatorGaugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	indicatorGaugeSchema(ctx, req, resp, false)
}

func (r *defaultIndicatorGaugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	indicatorGaugeSchema(ctx, req, resp, true)
}

func indicatorGaugeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Indicator Gauge.",
		Attributes: map[string]schema.Attribute{
			"gauge_data_source": schema.StringAttribute{
				Description: "Specifies the source of data to use in determining this Indicator Gauge's severity and status.",
				Required:    true,
			},
			"critical_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity should be critical.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"major_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'major'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"minor_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'minor'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"warning_value": schema.StringAttribute{
				Description: "A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'warning'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Gauge",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Gauge is enabled.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"override_severity": schema.StringAttribute{
				Description: "When defined, causes this Gauge to assume the specified severity, overriding its computed severity. This is useful for testing alarms generated by Gauges as well as suppressing alarms for known conditions.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"alert_level": schema.StringAttribute{
				Description: "Specifies the level at which alerts are sent for alarms raised by this Gauge.",
				Optional:    true,
			},
			"update_interval": schema.StringAttribute{
				Description: "The frequency with which this Gauge is updated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"samples_per_update_interval": schema.Int64Attribute{
				Description: "Indicates the number of times the monitor data source value will be collected during the update interval.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"include_resource": schema.SetAttribute{
				Description: "Specifies set of resources to be monitored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"exclude_resource": schema.SetAttribute{
				Description: "Specifies resources to exclude from being monitored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"server_unavailable_severity_level": schema.StringAttribute{
				Description: "Specifies the alarm severity level at or above which the server is considered unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_degraded_severity_level": schema.StringAttribute{
				Description: "Specifies the alarm severity level at or above which the server is considered degraded.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalIndicatorGaugeFields(ctx context.Context, addRequest *client.AddIndicatorGaugeRequest, plan indicatorGaugeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CriticalValue) {
		addRequest.CriticalValue = plan.CriticalValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MajorValue) {
		addRequest.MajorValue = plan.MajorValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinorValue) {
		addRequest.MinorValue = plan.MinorValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WarningValue) {
		addRequest.WarningValue = plan.WarningValue.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Enabled) {
		addRequest.Enabled = plan.Enabled.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OverrideSeverity) {
		overrideSeverity, err := client.NewEnumgaugeOverrideSeverityPropFromValue(plan.OverrideSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.OverrideSeverity = overrideSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AlertLevel) {
		alertLevel, err := client.NewEnumgaugeAlertLevelPropFromValue(plan.AlertLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.AlertLevel = alertLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.UpdateInterval) {
		addRequest.UpdateInterval = plan.UpdateInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SamplesPerUpdateInterval) {
		addRequest.SamplesPerUpdateInterval = plan.SamplesPerUpdateInterval.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeResource) {
		var slice []string
		plan.IncludeResource.ElementsAs(ctx, &slice, false)
		addRequest.IncludeResource = slice
	}
	if internaltypes.IsDefined(plan.ExcludeResource) {
		var slice []string
		plan.ExcludeResource.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeResource = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerUnavailableSeverityLevel) {
		serverUnavailableSeverityLevel, err := client.NewEnumgaugeServerUnavailableSeverityLevelPropFromValue(plan.ServerUnavailableSeverityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerUnavailableSeverityLevel = serverUnavailableSeverityLevel
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerDegradedSeverityLevel) {
		serverDegradedSeverityLevel, err := client.NewEnumgaugeServerDegradedSeverityLevelPropFromValue(plan.ServerDegradedSeverityLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.ServerDegradedSeverityLevel = serverDegradedSeverityLevel
	}
	return nil
}

// Read a IndicatorGaugeResponse object into the model struct
func readIndicatorGaugeResponse(ctx context.Context, r *client.IndicatorGaugeResponse, state *indicatorGaugeResourceModel, expectedValues *indicatorGaugeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.GaugeDataSource = types.StringValue(r.GaugeDataSource)
	state.CriticalValue = internaltypes.StringTypeOrNil(r.CriticalValue, internaltypes.IsEmptyString(expectedValues.CriticalValue))
	state.MajorValue = internaltypes.StringTypeOrNil(r.MajorValue, internaltypes.IsEmptyString(expectedValues.MajorValue))
	state.MinorValue = internaltypes.StringTypeOrNil(r.MinorValue, internaltypes.IsEmptyString(expectedValues.MinorValue))
	state.WarningValue = internaltypes.StringTypeOrNil(r.WarningValue, internaltypes.IsEmptyString(expectedValues.WarningValue))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.OverrideSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeOverrideSeverityProp(r.OverrideSeverity), internaltypes.IsEmptyString(expectedValues.OverrideSeverity))
	state.AlertLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeAlertLevelProp(r.AlertLevel), internaltypes.IsEmptyString(expectedValues.AlertLevel))
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, internaltypes.IsEmptyString(expectedValues.UpdateInterval))
	config.CheckMismatchedPDFormattedAttributes("update_interval",
		expectedValues.UpdateInterval, state.UpdateInterval, diagnostics)
	state.SamplesPerUpdateInterval = internaltypes.Int64TypeOrNil(r.SamplesPerUpdateInterval)
	state.IncludeResource = internaltypes.GetStringSet(r.IncludeResource)
	state.ExcludeResource = internaltypes.GetStringSet(r.ExcludeResource)
	state.ServerUnavailableSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerUnavailableSeverityLevelProp(r.ServerUnavailableSeverityLevel), internaltypes.IsEmptyString(expectedValues.ServerUnavailableSeverityLevel))
	state.ServerDegradedSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerDegradedSeverityLevelProp(r.ServerDegradedSeverityLevel), internaltypes.IsEmptyString(expectedValues.ServerDegradedSeverityLevel))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createIndicatorGaugeOperations(plan indicatorGaugeResourceModel, state indicatorGaugeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.GaugeDataSource, state.GaugeDataSource, "gauge-data-source")
	operations.AddStringOperationIfNecessary(&ops, plan.CriticalValue, state.CriticalValue, "critical-value")
	operations.AddStringOperationIfNecessary(&ops, plan.MajorValue, state.MajorValue, "major-value")
	operations.AddStringOperationIfNecessary(&ops, plan.MinorValue, state.MinorValue, "minor-value")
	operations.AddStringOperationIfNecessary(&ops, plan.WarningValue, state.WarningValue, "warning-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.AlertLevel, state.AlertLevel, "alert-level")
	operations.AddStringOperationIfNecessary(&ops, plan.UpdateInterval, state.UpdateInterval, "update-interval")
	operations.AddInt64OperationIfNecessary(&ops, plan.SamplesPerUpdateInterval, state.SamplesPerUpdateInterval, "samples-per-update-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludeResource, state.IncludeResource, "include-resource")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeResource, state.ExcludeResource, "exclude-resource")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerUnavailableSeverityLevel, state.ServerUnavailableSeverityLevel, "server-unavailable-severity-level")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerDegradedSeverityLevel, state.ServerDegradedSeverityLevel, "server-degraded-severity-level")
	return ops
}

// Create a new resource
func (r *indicatorGaugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan indicatorGaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddIndicatorGaugeRequest(plan.Id.ValueString(),
		[]client.EnumindicatorGaugeSchemaUrn{client.ENUMINDICATORGAUGESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGEINDICATOR},
		plan.GaugeDataSource.ValueString())
	err := addOptionalIndicatorGaugeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Indicator Gauge", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.GaugeApi.AddGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGaugeRequest(
		client.AddIndicatorGaugeRequestAsAddGaugeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.GaugeApi.AddGaugeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Indicator Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state indicatorGaugeResourceModel
	readIndicatorGaugeResponse(ctx, addResponse.IndicatorGaugeResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultIndicatorGaugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan indicatorGaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Indicator Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state indicatorGaugeResourceModel
	readIndicatorGaugeResponse(ctx, readResponse.IndicatorGaugeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GaugeApi.UpdateGauge(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createIndicatorGaugeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GaugeApi.UpdateGaugeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Indicator Gauge", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readIndicatorGaugeResponse(ctx, updateResponse.IndicatorGaugeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *indicatorGaugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIndicatorGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIndicatorGaugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readIndicatorGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readIndicatorGauge(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state indicatorGaugeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Indicator Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readIndicatorGaugeResponse(ctx, readResponse.IndicatorGaugeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *indicatorGaugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIndicatorGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultIndicatorGaugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateIndicatorGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateIndicatorGauge(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan indicatorGaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state indicatorGaugeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.GaugeApi.UpdateGauge(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createIndicatorGaugeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.GaugeApi.UpdateGaugeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Indicator Gauge", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readIndicatorGaugeResponse(ctx, updateResponse.IndicatorGaugeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultIndicatorGaugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *indicatorGaugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state indicatorGaugeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GaugeApi.DeleteGaugeExecute(r.apiClient.GaugeApi.DeleteGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Indicator Gauge", err, httpResp)
		return
	}
}

func (r *indicatorGaugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIndicatorGauge(ctx, req, resp)
}

func (r *defaultIndicatorGaugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importIndicatorGauge(ctx, req, resp)
}

func importIndicatorGauge(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
