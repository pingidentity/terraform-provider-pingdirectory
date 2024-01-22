package monitorprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
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
	r.apiClient = providerCfg.ApiClient
}

type monitorProviderDataSourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Type                                 types.String `tfsdk:"type"`
	ExtensionClass                       types.String `tfsdk:"extension_class"`
	ExtensionArgument                    types.Set    `tfsdk:"extension_argument"`
	LowSpaceWarningSizeThreshold         types.String `tfsdk:"low_space_warning_size_threshold"`
	LowSpaceWarningPercentThreshold      types.Int64  `tfsdk:"low_space_warning_percent_threshold"`
	LowSpaceErrorSizeThreshold           types.String `tfsdk:"low_space_error_size_threshold"`
	LowSpaceErrorPercentThreshold        types.Int64  `tfsdk:"low_space_error_percent_threshold"`
	OutOfSpaceErrorSizeThreshold         types.String `tfsdk:"out_of_space_error_size_threshold"`
	OutOfSpaceErrorPercentThreshold      types.Int64  `tfsdk:"out_of_space_error_percent_threshold"`
	AlertFrequency                       types.String `tfsdk:"alert_frequency"`
	CheckFrequency                       types.String `tfsdk:"check_frequency"`
	DiskDevices                          types.Set    `tfsdk:"disk_devices"`
	NetworkDevices                       types.Set    `tfsdk:"network_devices"`
	SystemUtilizationMonitorLogDirectory types.String `tfsdk:"system_utilization_monitor_log_directory"`
	ProlongedOutageDuration              types.String `tfsdk:"prolonged_outage_duration"`
	ProlongedOutageBehavior              types.String `tfsdk:"prolonged_outage_behavior"`
	Description                          types.String `tfsdk:"description"`
	Enabled                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *monitorProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Monitor Provider.",
		Attributes: map[string]schema.Attribute{
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
			"low_space_warning_size_threshold": schema.StringAttribute{
				Description: "Specifies the low space warning threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will begin generating warning alert notifications.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"low_space_warning_percent_threshold": schema.Int64Attribute{
				Description: "Specifies the low space warning threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will begin generating warning alert notifications.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"low_space_error_size_threshold": schema.StringAttribute{
				Description: "Specifies the low space error threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will start rejecting operations requested by non-root users.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"low_space_error_percent_threshold": schema.Int64Attribute{
				Description: "Specifies the low space error threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will start rejecting operations requested by non-root users.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"out_of_space_error_size_threshold": schema.StringAttribute{
				Description: "Specifies the out of space error threshold value as an absolute amount of space. If the amount of usable disk space drops below this amount, then the Directory Server will shut itself down to avoid problems that may occur from complete exhaustion of usable space.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"out_of_space_error_percent_threshold": schema.Int64Attribute{
				Description: "Specifies the out of space error threshold value as a percentage of total space. If the amount of usable disk space drops below this amount, then the Directory Server will shut itself down to avoid problems that may occur from complete exhaustion of usable space.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"alert_frequency": schema.StringAttribute{
				Description: "Specifies the length of time between administrative alerts generated in response to lack of usable disk space. Administrative alerts will be generated whenever the amount of usable space drops below any threshold, and they will also be generated at regular intervals as long as the amount of usable space remains below the threshold value. A value of zero indicates that alerts should only be generated when the amount of usable space drops below a configured threshold.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"check_frequency": schema.StringAttribute{
				Description: "The frequency with which this monitor provider should confirm the ability to access the server's encryption settings database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disk_devices": schema.SetAttribute{
				Description: "Specifies which disk devices to monitor for I/O activity. Should be the device name as displayed by iostat -d.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"network_devices": schema.SetAttribute{
				Description: "Specifies which network interfaces to monitor for I/O activity. Should be the device name as displayed by netstat -i.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"system_utilization_monitor_log_directory": schema.StringAttribute{
				Description: "Specifies a relative or absolute path to the directory on the local filesystem containing the log files used by the system utilization monitor. The path must exist, and it must be a writable directory by the server process.",
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
				Description:         "When the `type` attribute is set to  one of [`memory-usage`, `stack-trace`, `encryption-settings-database-accessibility`, `custom`, `active-operations`, `ssl-context`, `version`, `general`, `disk-space-usage`, `system-info`, `client-connection`, `third-party`]: Indicates whether the Monitor Provider is enabled for use. When the `type` attribute is set to `host-system`: Indicates whether the Host System Monitor Provider is enabled for use.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`memory-usage`, `stack-trace`, `encryption-settings-database-accessibility`, `custom`, `active-operations`, `ssl-context`, `version`, `general`, `disk-space-usage`, `system-info`, `client-connection`, `third-party`]: Indicates whether the Monitor Provider is enabled for use.\n  - `host-system`: Indicates whether the Host System Monitor Provider is enabled for use.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a MemoryUsageMonitorProviderResponse object into the model struct
func readMemoryUsageMonitorProviderResponseDataSource(ctx context.Context, r *client.MemoryUsageMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("memory-usage")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a StackTraceMonitorProviderResponse object into the model struct
func readStackTraceMonitorProviderResponseDataSource(ctx context.Context, r *client.StackTraceMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("stack-trace")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse object into the model struct
func readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDataSource(ctx context.Context, r *client.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("encryption-settings-database-accessibility")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CheckFrequency = types.StringValue(r.CheckFrequency)
	state.ProlongedOutageDuration = internaltypes.StringTypeOrNil(r.ProlongedOutageDuration, false)
	state.ProlongedOutageBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnummonitorProviderProlongedOutageBehaviorProp(r.ProlongedOutageBehavior), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CustomMonitorProviderResponse object into the model struct
func readCustomMonitorProviderResponseDataSource(ctx context.Context, r *client.CustomMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ActiveOperationsMonitorProviderResponse object into the model struct
func readActiveOperationsMonitorProviderResponseDataSource(ctx context.Context, r *client.ActiveOperationsMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("active-operations")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SslContextMonitorProviderResponse object into the model struct
func readSslContextMonitorProviderResponseDataSource(ctx context.Context, r *client.SslContextMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ssl-context")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a VersionMonitorProviderResponse object into the model struct
func readVersionMonitorProviderResponseDataSource(ctx context.Context, r *client.VersionMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("version")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a HostSystemMonitorProviderResponse object into the model struct
func readHostSystemMonitorProviderResponseDataSource(ctx context.Context, r *client.HostSystemMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("host-system")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.DiskDevices = internaltypes.GetStringSet(r.DiskDevices)
	state.NetworkDevices = internaltypes.GetStringSet(r.NetworkDevices)
	state.SystemUtilizationMonitorLogDirectory = types.StringValue(r.SystemUtilizationMonitorLogDirectory)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a GeneralMonitorProviderResponse object into the model struct
func readGeneralMonitorProviderResponseDataSource(ctx context.Context, r *client.GeneralMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("general")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a DiskSpaceUsageMonitorProviderResponse object into the model struct
func readDiskSpaceUsageMonitorProviderResponseDataSource(ctx context.Context, r *client.DiskSpaceUsageMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("disk-space-usage")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LowSpaceWarningSizeThreshold = internaltypes.StringTypeOrNil(r.LowSpaceWarningSizeThreshold, false)
	state.LowSpaceWarningPercentThreshold = internaltypes.Int64TypeOrNil(r.LowSpaceWarningPercentThreshold)
	state.LowSpaceErrorSizeThreshold = internaltypes.StringTypeOrNil(r.LowSpaceErrorSizeThreshold, false)
	state.LowSpaceErrorPercentThreshold = internaltypes.Int64TypeOrNil(r.LowSpaceErrorPercentThreshold)
	state.OutOfSpaceErrorSizeThreshold = internaltypes.StringTypeOrNil(r.OutOfSpaceErrorSizeThreshold, false)
	state.OutOfSpaceErrorPercentThreshold = internaltypes.Int64TypeOrNil(r.OutOfSpaceErrorPercentThreshold)
	state.AlertFrequency = types.StringValue(r.AlertFrequency)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SystemInfoMonitorProviderResponse object into the model struct
func readSystemInfoMonitorProviderResponseDataSource(ctx context.Context, r *client.SystemInfoMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("system-info")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ClientConnectionMonitorProviderResponse object into the model struct
func readClientConnectionMonitorProviderResponseDataSource(ctx context.Context, r *client.ClientConnectionMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("client-connection")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyMonitorProviderResponse object into the model struct
func readThirdPartyMonitorProviderResponseDataSource(ctx context.Context, r *client.ThirdPartyMonitorProviderResponse, state *monitorProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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

	readResponse, httpResp, err := r.apiClient.MonitorProviderAPI.GetMonitorProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
	if readResponse.MemoryUsageMonitorProviderResponse != nil {
		readMemoryUsageMonitorProviderResponseDataSource(ctx, readResponse.MemoryUsageMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.StackTraceMonitorProviderResponse != nil {
		readStackTraceMonitorProviderResponseDataSource(ctx, readResponse.StackTraceMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse != nil {
		readEncryptionSettingsDatabaseAccessibilityMonitorProviderResponseDataSource(ctx, readResponse.EncryptionSettingsDatabaseAccessibilityMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomMonitorProviderResponse != nil {
		readCustomMonitorProviderResponseDataSource(ctx, readResponse.CustomMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ActiveOperationsMonitorProviderResponse != nil {
		readActiveOperationsMonitorProviderResponseDataSource(ctx, readResponse.ActiveOperationsMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SslContextMonitorProviderResponse != nil {
		readSslContextMonitorProviderResponseDataSource(ctx, readResponse.SslContextMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VersionMonitorProviderResponse != nil {
		readVersionMonitorProviderResponseDataSource(ctx, readResponse.VersionMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HostSystemMonitorProviderResponse != nil {
		readHostSystemMonitorProviderResponseDataSource(ctx, readResponse.HostSystemMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GeneralMonitorProviderResponse != nil {
		readGeneralMonitorProviderResponseDataSource(ctx, readResponse.GeneralMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DiskSpaceUsageMonitorProviderResponse != nil {
		readDiskSpaceUsageMonitorProviderResponseDataSource(ctx, readResponse.DiskSpaceUsageMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SystemInfoMonitorProviderResponse != nil {
		readSystemInfoMonitorProviderResponseDataSource(ctx, readResponse.SystemInfoMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ClientConnectionMonitorProviderResponse != nil {
		readClientConnectionMonitorProviderResponseDataSource(ctx, readResponse.ClientConnectionMonitorProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyMonitorProviderResponse != nil {
		readThirdPartyMonitorProviderResponseDataSource(ctx, readResponse.ThirdPartyMonitorProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
