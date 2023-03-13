package gauge

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
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
	_ resource.Resource                = &numericGaugeResource{}
	_ resource.ResourceWithConfigure   = &numericGaugeResource{}
	_ resource.ResourceWithImportState = &numericGaugeResource{}
	_ resource.Resource                = &defaultNumericGaugeResource{}
	_ resource.ResourceWithConfigure   = &defaultNumericGaugeResource{}
	_ resource.ResourceWithImportState = &defaultNumericGaugeResource{}
)

// Create a Numeric Gauge resource
func NewNumericGaugeResource() resource.Resource {
	return &numericGaugeResource{}
}

func NewDefaultNumericGaugeResource() resource.Resource {
	return &defaultNumericGaugeResource{}
}

// numericGaugeResource is the resource implementation.
type numericGaugeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultNumericGaugeResource is the resource implementation.
type defaultNumericGaugeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *numericGaugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_numeric_gauge"
}

func (r *defaultNumericGaugeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_numeric_gauge"
}

// Configure adds the provider configured client to the resource.
func (r *numericGaugeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultNumericGaugeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type numericGaugeResourceModel struct {
	Id                             types.String  `tfsdk:"id"`
	LastUpdated                    types.String  `tfsdk:"last_updated"`
	Notifications                  types.Set     `tfsdk:"notifications"`
	RequiredActions                types.Set     `tfsdk:"required_actions"`
	GaugeDataSource                types.String  `tfsdk:"gauge_data_source"`
	CriticalValue                  types.Float64 `tfsdk:"critical_value"`
	CriticalExitValue              types.Float64 `tfsdk:"critical_exit_value"`
	MajorValue                     types.Float64 `tfsdk:"major_value"`
	MajorExitValue                 types.Float64 `tfsdk:"major_exit_value"`
	MinorValue                     types.Float64 `tfsdk:"minor_value"`
	MinorExitValue                 types.Float64 `tfsdk:"minor_exit_value"`
	WarningValue                   types.Float64 `tfsdk:"warning_value"`
	WarningExitValue               types.Float64 `tfsdk:"warning_exit_value"`
	Description                    types.String  `tfsdk:"description"`
	Enabled                        types.Bool    `tfsdk:"enabled"`
	OverrideSeverity               types.String  `tfsdk:"override_severity"`
	AlertLevel                     types.String  `tfsdk:"alert_level"`
	UpdateInterval                 types.String  `tfsdk:"update_interval"`
	SamplesPerUpdateInterval       types.Int64   `tfsdk:"samples_per_update_interval"`
	IncludeResource                types.Set     `tfsdk:"include_resource"`
	ExcludeResource                types.Set     `tfsdk:"exclude_resource"`
	ServerUnavailableSeverityLevel types.String  `tfsdk:"server_unavailable_severity_level"`
	ServerDegradedSeverityLevel    types.String  `tfsdk:"server_degraded_severity_level"`
}

// GetSchema defines the schema for the resource.
func (r *numericGaugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	numericGaugeSchema(ctx, req, resp, false)
}

func (r *defaultNumericGaugeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	numericGaugeSchema(ctx, req, resp, true)
}

func numericGaugeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Numeric Gauge.",
		Attributes: map[string]schema.Attribute{
			"gauge_data_source": schema.StringAttribute{
				Description: "Specifies the source of data to use in determining this gauge's current severity.",
				Required:    true,
			},
			"critical_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'critical'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"critical_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'critical'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"major_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'major'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"major_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'major'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"minor_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'minor'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"minor_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'minor'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"warning_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'warning'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"warning_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'warning'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
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
func addOptionalNumericGaugeFields(ctx context.Context, addRequest *client.AddNumericGaugeRequest, plan numericGaugeResourceModel) error {
	if internaltypes.IsDefined(plan.CriticalValue) {
		floatVal := float32(plan.CriticalValue.ValueFloat64())
		addRequest.CriticalValue = &floatVal
	}
	if internaltypes.IsDefined(plan.CriticalExitValue) {
		floatVal := float32(plan.CriticalExitValue.ValueFloat64())
		addRequest.CriticalExitValue = &floatVal
	}
	if internaltypes.IsDefined(plan.MajorValue) {
		floatVal := float32(plan.MajorValue.ValueFloat64())
		addRequest.MajorValue = &floatVal
	}
	if internaltypes.IsDefined(plan.MajorExitValue) {
		floatVal := float32(plan.MajorExitValue.ValueFloat64())
		addRequest.MajorExitValue = &floatVal
	}
	if internaltypes.IsDefined(plan.MinorValue) {
		floatVal := float32(plan.MinorValue.ValueFloat64())
		addRequest.MinorValue = &floatVal
	}
	if internaltypes.IsDefined(plan.MinorExitValue) {
		floatVal := float32(plan.MinorExitValue.ValueFloat64())
		addRequest.MinorExitValue = &floatVal
	}
	if internaltypes.IsDefined(plan.WarningValue) {
		floatVal := float32(plan.WarningValue.ValueFloat64())
		addRequest.WarningValue = &floatVal
	}
	if internaltypes.IsDefined(plan.WarningExitValue) {
		floatVal := float32(plan.WarningExitValue.ValueFloat64())
		addRequest.WarningExitValue = &floatVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.Enabled) {
		boolVal := plan.Enabled.ValueBool()
		addRequest.Enabled = &boolVal
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
		stringVal := plan.UpdateInterval.ValueString()
		addRequest.UpdateInterval = &stringVal
	}
	if internaltypes.IsDefined(plan.SamplesPerUpdateInterval) {
		intVal := int32(plan.SamplesPerUpdateInterval.ValueInt64())
		addRequest.SamplesPerUpdateInterval = &intVal
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

// Read a NumericGaugeResponse object into the model struct
func readNumericGaugeResponse(ctx context.Context, r *client.NumericGaugeResponse, state *numericGaugeResourceModel, expectedValues *numericGaugeResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.GaugeDataSource = types.StringValue(r.GaugeDataSource)
	state.CriticalValue = internaltypes.Float64TypeOrNil(r.CriticalValue)
	state.CriticalExitValue = internaltypes.Float64TypeOrNil(r.CriticalExitValue)
	state.MajorValue = internaltypes.Float64TypeOrNil(r.MajorValue)
	state.MajorExitValue = internaltypes.Float64TypeOrNil(r.MajorExitValue)
	state.MinorValue = internaltypes.Float64TypeOrNil(r.MinorValue)
	state.MinorExitValue = internaltypes.Float64TypeOrNil(r.MinorExitValue)
	state.WarningValue = internaltypes.Float64TypeOrNil(r.WarningValue)
	state.WarningExitValue = internaltypes.Float64TypeOrNil(r.WarningExitValue)
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
func createNumericGaugeOperations(plan numericGaugeResourceModel, state numericGaugeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.GaugeDataSource, state.GaugeDataSource, "gauge-data-source")
	operations.AddFloat64OperationIfNecessary(&ops, plan.CriticalValue, state.CriticalValue, "critical-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.CriticalExitValue, state.CriticalExitValue, "critical-exit-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.MajorValue, state.MajorValue, "major-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.MajorExitValue, state.MajorExitValue, "major-exit-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.MinorValue, state.MinorValue, "minor-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.MinorExitValue, state.MinorExitValue, "minor-exit-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.WarningValue, state.WarningValue, "warning-value")
	operations.AddFloat64OperationIfNecessary(&ops, plan.WarningExitValue, state.WarningExitValue, "warning-exit-value")
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
func (r *numericGaugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan numericGaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddNumericGaugeRequest(plan.Id.ValueString(),
		[]client.EnumnumericGaugeSchemaUrn{client.ENUMNUMERICGAUGESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0GAUGENUMERIC},
		plan.GaugeDataSource.ValueString())
	err := addOptionalNumericGaugeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Numeric Gauge", err.Error())
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
		client.AddNumericGaugeRequestAsAddGaugeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.GaugeApi.AddGaugeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Numeric Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state numericGaugeResourceModel
	readNumericGaugeResponse(ctx, addResponse.NumericGaugeResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultNumericGaugeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan numericGaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Numeric Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state numericGaugeResourceModel
	readNumericGaugeResponse(ctx, readResponse.NumericGaugeResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.GaugeApi.UpdateGauge(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createNumericGaugeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.GaugeApi.UpdateGaugeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Numeric Gauge", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readNumericGaugeResponse(ctx, updateResponse.NumericGaugeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *numericGaugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readNumericGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultNumericGaugeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readNumericGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readNumericGauge(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state numericGaugeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Numeric Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readNumericGaugeResponse(ctx, readResponse.NumericGaugeResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *numericGaugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateNumericGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultNumericGaugeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateNumericGauge(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateNumericGauge(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan numericGaugeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state numericGaugeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.GaugeApi.UpdateGauge(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createNumericGaugeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.GaugeApi.UpdateGaugeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Numeric Gauge", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readNumericGaugeResponse(ctx, updateResponse.NumericGaugeResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultNumericGaugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *numericGaugeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state numericGaugeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.GaugeApi.DeleteGaugeExecute(r.apiClient.GaugeApi.DeleteGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Numeric Gauge", err, httpResp)
		return
	}
}

func (r *numericGaugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importNumericGauge(ctx, req, resp)
}

func (r *defaultNumericGaugeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importNumericGauge(ctx, req, resp)
}

func importNumericGauge(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
