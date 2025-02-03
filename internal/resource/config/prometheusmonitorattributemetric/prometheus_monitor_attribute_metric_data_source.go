// Copyright © 2025 Ping Identity Corporation

package prometheusmonitorattributemetric

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &prometheusMonitorAttributeMetricDataSource{}
	_ datasource.DataSourceWithConfigure = &prometheusMonitorAttributeMetricDataSource{}
)

// Create a Prometheus Monitor Attribute Metric data source
func NewPrometheusMonitorAttributeMetricDataSource() datasource.DataSource {
	return &prometheusMonitorAttributeMetricDataSource{}
}

// prometheusMonitorAttributeMetricDataSource is the datasource implementation.
type prometheusMonitorAttributeMetricDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *prometheusMonitorAttributeMetricDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prometheus_monitor_attribute_metric"
}

// Configure adds the provider configured client to the data source.
func (r *prometheusMonitorAttributeMetricDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type prometheusMonitorAttributeMetricDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Type                     types.String `tfsdk:"type"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
	MetricName               types.String `tfsdk:"metric_name"`
	MonitorAttributeName     types.String `tfsdk:"monitor_attribute_name"`
	MonitorObjectClassName   types.String `tfsdk:"monitor_object_class_name"`
	MetricType               types.String `tfsdk:"metric_type"`
	Filter                   types.String `tfsdk:"filter"`
	MetricDescription        types.String `tfsdk:"metric_description"`
	LabelNameValuePair       types.Set    `tfsdk:"label_name_value_pair"`
}

// GetSchema defines the schema for the datasource.
func (r *prometheusMonitorAttributeMetricDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Prometheus Monitor Attribute Metric. Supported in PingDirectory product version 9.2.0.0+.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Prometheus Monitor Attribute Metric resource. Options are ['prometheus-monitor-attribute-metric']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
			},
			"metric_name": schema.StringAttribute{
				Description: "The name that will be used in the metric to be consumed by Prometheus.",
				Required:    true,
			},
			"monitor_attribute_name": schema.StringAttribute{
				Description: "The name of the monitor attribute that contains the numeric value to be published.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"monitor_object_class_name": schema.StringAttribute{
				Description: "The name of the object class for monitor entries that contain the monitor attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"metric_type": schema.StringAttribute{
				Description: "The metric type that should be used for the value of the specified monitor attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "A filter that may be used to restrict the set of monitor entries for which the metric should be generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"metric_description": schema.StringAttribute{
				Description: "A human-readable description that should be published as part of the metric definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"label_name_value_pair": schema.SetAttribute{
				Description: "A set of name-value pairs for labels that should be included in the published metric for the target attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any version restrictions are met
func (r *prometheusMonitorAttributeMetricDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if r.providerConfig.ProductVersion != "" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
			r.providerConfig.ProductVersion, "pingdirectory_prometheus_monitor_attribute_metric")
	}
}

// Read a PrometheusMonitorAttributeMetricResponse object into the model struct
func readPrometheusMonitorAttributeMetricResponseDataSource(ctx context.Context, r *client.PrometheusMonitorAttributeMetricResponse, state *prometheusMonitorAttributeMetricDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("prometheus-monitor-attribute-metric")
	state.Id = types.StringValue(r.Id)
	state.MetricName = types.StringValue(r.MetricName)
	state.MonitorAttributeName = types.StringValue(r.MonitorAttributeName)
	state.MonitorObjectClassName = types.StringValue(r.MonitorObjectClassName)
	state.MetricType = types.StringValue(r.MetricType.String())
	state.Filter = internaltypes.StringTypeOrNil(r.Filter, false)
	state.MetricDescription = internaltypes.StringTypeOrNil(r.MetricDescription, false)
	state.LabelNameValuePair = internaltypes.GetStringSet(r.LabelNameValuePair)
}

// Read resource information
func (r *prometheusMonitorAttributeMetricDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state prometheusMonitorAttributeMetricDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PrometheusMonitorAttributeMetricAPI.GetPrometheusMonitorAttributeMetric(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.MetricName.ValueString(), state.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Prometheus Monitor Attribute Metric", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPrometheusMonitorAttributeMetricResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
