package prometheusmonitorattributemetric

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &prometheusMonitorAttributeMetricsDataSource{}
	_ datasource.DataSourceWithConfigure = &prometheusMonitorAttributeMetricsDataSource{}
)

// Create a Prometheus Monitor Attribute Metrics data source
func NewPrometheusMonitorAttributeMetricsDataSource() datasource.DataSource {
	return &prometheusMonitorAttributeMetricsDataSource{}
}

// prometheusMonitorAttributeMetricsDataSource is the datasource implementation.
type prometheusMonitorAttributeMetricsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *prometheusMonitorAttributeMetricsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_prometheus_monitor_attribute_metrics"
}

// Configure adds the provider configured client to the data source.
func (r *prometheusMonitorAttributeMetricsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type prometheusMonitorAttributeMetricsDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Filter                   types.String `tfsdk:"filter"`
	Ids                      types.Set    `tfsdk:"ids"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
}

// GetSchema defines the schema for the datasource.
func (r *prometheusMonitorAttributeMetricsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Prometheus Monitor Attribute Metric objects in the server configuration. Supported in PingDirectory product version 9.2.0.0+.",
		Attributes: map[string]schema.Attribute{
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"ids": schema.SetAttribute{
				Description: "Prometheus Monitor Attribute Metric IDs found in the configuration",
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
func (r *prometheusMonitorAttributeMetricsDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if r.providerConfig.ProductVersion != "" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory9200,
			r.providerConfig.ProductVersion, "pingdirectory_prometheus_monitor_attribute_metrics")
	}
}

// Read resource information
func (r *prometheusMonitorAttributeMetricsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state prometheusMonitorAttributeMetricsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.PrometheusMonitorAttributeMetricAPI.ListPrometheusMonitorAttributeMetrics(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.HttpServletExtensionName.ValueString())
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.PrometheusMonitorAttributeMetricAPI.ListPrometheusMonitorAttributeMetricsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Prometheus Monitor Attribute Metric objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	ids := []attr.Value{}
	for _, response := range readResponse.Resources {
		ids = append(ids, types.StringValue(response.Id))
	}

	state.Ids, diags = types.SetValue(types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
