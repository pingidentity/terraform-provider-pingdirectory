package gauge

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &gaugeDataSource{}
	_ datasource.DataSourceWithConfigure = &gaugeDataSource{}
)

// Create a Gauge data source
func NewGaugeDataSource() datasource.DataSource {
	return &gaugeDataSource{}
}

// gaugeDataSource is the datasource implementation.
type gaugeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *gaugeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gauge"
}

// Configure adds the provider configured client to the data source.
func (r *gaugeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type gaugeDataSourceModel struct {
	Id                             types.String  `tfsdk:"id"`
	Name                           types.String  `tfsdk:"name"`
	Type                           types.String  `tfsdk:"type"`
	GaugeDataSource                types.String  `tfsdk:"gauge_data_source"`
	CriticalValue                  types.String  `tfsdk:"critical_value"`
	CriticalExitValue              types.Float64 `tfsdk:"critical_exit_value"`
	MajorValue                     types.String  `tfsdk:"major_value"`
	MajorExitValue                 types.Float64 `tfsdk:"major_exit_value"`
	MinorValue                     types.String  `tfsdk:"minor_value"`
	MinorExitValue                 types.Float64 `tfsdk:"minor_exit_value"`
	WarningValue                   types.String  `tfsdk:"warning_value"`
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

// GetSchema defines the schema for the datasource.
func (r *gaugeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Gauge.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Gauge resource. Options are ['indicator', 'numeric']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"gauge_data_source": schema.StringAttribute{
				Description: " When the `type` value is one of [`indicator`]: Specifies the source of data to use in determining this Indicator Gauge's severity and status. When the `type` value is one of [`numeric`]: Specifies the source of data to use in determining this gauge's current severity.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"critical_value": schema.StringAttribute{
				Description: " When the `type` value is one of [`indicator`]: A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity should be critical. When the `type` value is one of [`numeric`]: A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'critical'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"critical_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'critical'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"major_value": schema.StringAttribute{
				Description: " When the `type` value is one of [`indicator`]: A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'major'. When the `type` value is one of [`numeric`]: A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'major'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"major_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'major'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minor_value": schema.StringAttribute{
				Description: " When the `type` value is one of [`indicator`]: A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'minor'. When the `type` value is one of [`numeric`]: A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'minor'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minor_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'minor'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"warning_value": schema.StringAttribute{
				Description: " When the `type` value is one of [`indicator`]: A regular expression pattern that is used to determine whether the current monitored value indicates this gauge's severity will be 'warning'. When the `type` value is one of [`numeric`]: A value that is used to determine whether the current monitored value indicates this gauge's severity should be 'warning'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"warning_exit_value": schema.Float64Attribute{
				Description: "A value that is used to determine whether the current monitored value indicates this gauge's severity should no longer be 'warning'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Gauge",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Gauge is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"override_severity": schema.StringAttribute{
				Description: "When defined, causes this Gauge to assume the specified severity, overriding its computed severity. This is useful for testing alarms generated by Gauges as well as suppressing alarms for known conditions.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alert_level": schema.StringAttribute{
				Description: "Specifies the level at which alerts are sent for alarms raised by this Gauge.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"update_interval": schema.StringAttribute{
				Description: "The frequency with which this Gauge is updated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"samples_per_update_interval": schema.Int64Attribute{
				Description: "Indicates the number of times the monitor data source value will be collected during the update interval.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_resource": schema.SetAttribute{
				Description: "Specifies set of resources to be monitored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_resource": schema.SetAttribute{
				Description: "Specifies resources to exclude from being monitored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"server_unavailable_severity_level": schema.StringAttribute{
				Description: "Specifies the alarm severity level at or above which the server is considered unavailable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_degraded_severity_level": schema.StringAttribute{
				Description: "Specifies the alarm severity level at or above which the server is considered degraded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a IndicatorGaugeResponse object into the model struct
func readIndicatorGaugeResponseDataSource(ctx context.Context, r *client.IndicatorGaugeResponse, state *gaugeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("indicator")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.GaugeDataSource = types.StringValue(r.GaugeDataSource)
	state.CriticalValue = internaltypes.StringTypeOrNil(r.CriticalValue, false)
	state.MajorValue = internaltypes.StringTypeOrNil(r.MajorValue, false)
	state.MinorValue = internaltypes.StringTypeOrNil(r.MinorValue, false)
	state.WarningValue = internaltypes.StringTypeOrNil(r.WarningValue, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OverrideSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeOverrideSeverityProp(r.OverrideSeverity), false)
	state.AlertLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeAlertLevelProp(r.AlertLevel), false)
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, false)
	state.SamplesPerUpdateInterval = internaltypes.Int64TypeOrNil(r.SamplesPerUpdateInterval)
	state.IncludeResource = internaltypes.GetStringSet(r.IncludeResource)
	state.ExcludeResource = internaltypes.GetStringSet(r.ExcludeResource)
	state.ServerUnavailableSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerUnavailableSeverityLevelProp(r.ServerUnavailableSeverityLevel), false)
	state.ServerDegradedSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerDegradedSeverityLevelProp(r.ServerDegradedSeverityLevel), false)
}

// Read a NumericGaugeResponse object into the model struct
func readNumericGaugeResponseDataSource(ctx context.Context, r *client.NumericGaugeResponse, state *gaugeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("numeric")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.GaugeDataSource = types.StringValue(r.GaugeDataSource)
	if r.CriticalValue == nil {
		state.CriticalValue = types.StringNull()
	} else {
		state.CriticalValue = types.StringValue(strconv.FormatFloat(*r.CriticalValue, 'f', -1, 64))
	}
	state.CriticalExitValue = internaltypes.Float64TypeOrNil(r.CriticalExitValue)
	if r.MajorValue == nil {
		state.MajorValue = types.StringNull()
	} else {
		state.MajorValue = types.StringValue(strconv.FormatFloat(*r.MajorValue, 'f', -1, 64))
	}
	state.MajorExitValue = internaltypes.Float64TypeOrNil(r.MajorExitValue)
	if r.MinorValue == nil {
		state.MinorValue = types.StringNull()
	} else {
		state.MinorValue = types.StringValue(strconv.FormatFloat(*r.MinorValue, 'f', -1, 64))
	}
	state.MinorExitValue = internaltypes.Float64TypeOrNil(r.MinorExitValue)
	if r.WarningValue == nil {
		state.WarningValue = types.StringNull()
	} else {
		state.WarningValue = types.StringValue(strconv.FormatFloat(*r.WarningValue, 'f', -1, 64))
	}
	state.WarningExitValue = internaltypes.Float64TypeOrNil(r.WarningExitValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OverrideSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeOverrideSeverityProp(r.OverrideSeverity), false)
	state.AlertLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeAlertLevelProp(r.AlertLevel), false)
	state.UpdateInterval = internaltypes.StringTypeOrNil(r.UpdateInterval, false)
	state.SamplesPerUpdateInterval = internaltypes.Int64TypeOrNil(r.SamplesPerUpdateInterval)
	state.IncludeResource = internaltypes.GetStringSet(r.IncludeResource)
	state.ExcludeResource = internaltypes.GetStringSet(r.ExcludeResource)
	state.ServerUnavailableSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerUnavailableSeverityLevelProp(r.ServerUnavailableSeverityLevel), false)
	state.ServerDegradedSeverityLevel = internaltypes.StringTypeOrNil(
		client.StringPointerEnumgaugeServerDegradedSeverityLevelProp(r.ServerDegradedSeverityLevel), false)
}

// Read resource information
func (r *gaugeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state gaugeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.GaugeApi.GetGauge(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Gauge", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.IndicatorGaugeResponse != nil {
		readIndicatorGaugeResponseDataSource(ctx, readResponse.IndicatorGaugeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NumericGaugeResponse != nil {
		readNumericGaugeResponseDataSource(ctx, readResponse.NumericGaugeResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
