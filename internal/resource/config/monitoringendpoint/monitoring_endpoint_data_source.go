package monitoringendpoint

import (
	"context"

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
	_ datasource.DataSource              = &monitoringEndpointDataSource{}
	_ datasource.DataSourceWithConfigure = &monitoringEndpointDataSource{}
)

// Create a Monitoring Endpoint data source
func NewMonitoringEndpointDataSource() datasource.DataSource {
	return &monitoringEndpointDataSource{}
}

// monitoringEndpointDataSource is the datasource implementation.
type monitoringEndpointDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *monitoringEndpointDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitoring_endpoint"
}

// Configure adds the provider configured client to the data source.
func (r *monitoringEndpointDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type monitoringEndpointDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Hostname             types.String `tfsdk:"hostname"`
	ServerPort           types.Int64  `tfsdk:"server_port"`
	ConnectionType       types.String `tfsdk:"connection_type"`
	TrustManagerProvider types.String `tfsdk:"trust_manager_provider"`
	AdditionalTags       types.Set    `tfsdk:"additional_tags"`
	Enabled              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *monitoringEndpointDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Monitoring Endpoint.",
		Attributes: map[string]schema.Attribute{
			"hostname": schema.StringAttribute{
				Description: "The name of the host where this StatsD Monitoring Endpoint should send metric data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the endpoint where metric data should be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_type": schema.StringAttribute{
				Description: "Specifies the protocol and security that this StatsD Monitoring Endpoint should use to connect to the configured endpoint.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use if SSL over TCP is to be used for connection-level security.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"additional_tags": schema.SetAttribute{
				Description: "Specifies any optional additional tags to include in StatsD messages. Any additional tags will be appended to the end of each StatsD message, separated by commas. Tags should be written in a [key]:[value] format (\"host:server1\", for example).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Monitoring Endpoint is enabled for use in the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a StatsdMonitoringEndpointResponse object into the model struct
func readStatsdMonitoringEndpointResponseDataSource(ctx context.Context, r *client.StatsdMonitoringEndpointResponse, state *monitoringEndpointDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Hostname = types.StringValue(r.Hostname)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.ConnectionType = types.StringValue(r.ConnectionType.String())
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.AdditionalTags = internaltypes.GetStringSet(r.AdditionalTags)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *monitoringEndpointDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state monitoringEndpointDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitoringEndpointApi.GetMonitoringEndpoint(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitoring Endpoint", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readStatsdMonitoringEndpointResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
