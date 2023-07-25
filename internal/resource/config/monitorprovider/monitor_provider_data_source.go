package monitorprovider

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
	_ datasource.DataSource              = &monitorProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &monitorProviderDataSource{}
)

// Create a Monitor Provider data source
func NewMonitorProviderDataSource() datasource.DataSource {
	return &monitorProviderDataSource{}
}

// monitorProviderDataSource is the datasource implementation.
type monitorProviderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *monitorProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_monitor_provider"
}

// Configure adds the provider configured client to the data source.
func (r *monitorProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type monitorProviderDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	Type                    types.String `tfsdk:"type"`
	ExtensionClass          types.String `tfsdk:"extension_class"`
	ExtensionArgument       types.Set    `tfsdk:"extension_argument"`
	CheckFrequency          types.String `tfsdk:"check_frequency"`
	ProlongedOutageDuration types.String `tfsdk:"prolonged_outage_duration"`
	ProlongedOutageBehavior types.String `tfsdk:"prolonged_outage_behavior"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *monitorProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Monitor Provider.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Monitor Provider resource. Options are ['memory-usage', 'stack-trace', 'encryption-settings-database-accessibility', 'custom', 'active-operations', 'ssl-context', 'version', 'host-system', 'general', 'disk-space-usage', 'system-info', 'client-connection', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Monitor Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Monitor Provider. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"check_frequency": schema.StringAttribute{
				Description: "The frequency with which this monitor provider should confirm the ability to access the server's encryption settings database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prolonged_outage_duration": schema.StringAttribute{
				Description: "The minimum length of time that an outage should persist before it is considered a prolonged outage. If an outage lasts at least as long as this duration, then the server will take the action indicated by the prolonged-outage-behavior property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prolonged_outage_behavior": schema.StringAttribute{
				Description: "The behavior that the server should exhibit after a prolonged period of time when the encryption settings database remains unreadable.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Monitor Provider",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Monitor Provider is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse object into the model struct
func readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDataSource(ctx context.Context, r *client.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings-database-accessibility")
	state.Id = types.StringValue(r.Id)
	state.CheckFrequency = types.StringValue(r.CheckFrequency)
	state.ProlongedOutageDuration = internaltypes.StringTypeOrNil(r.ProlongedOutageDuration, false)
	state.ProlongedOutageBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnummonitorProviderProlongedOutageBehaviorProp(r.ProlongedOutageBehavior), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyMonitorProviderResponse object into the model struct
func readThirdPartyMonitorProviderResponseDataSource(ctx context.Context, r *client.ThirdPartyMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *monitorProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state monitorProviderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.MonitorProviderApi.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Monitor Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
		readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDataSource(ctx, readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyMonitorProviderResponse != nil {
		readThirdPartyMonitorProviderResponseDataSource(ctx, readResponse.ThirdPartyMonitorProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
