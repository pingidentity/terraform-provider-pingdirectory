package logfilerotationlistener

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
	_ datasource.DataSource              = &logFileRotationListenerDataSource{}
	_ datasource.DataSourceWithConfigure = &logFileRotationListenerDataSource{}
)

// Create a Log File Rotation Listener data source
func NewLogFileRotationListenerDataSource() datasource.DataSource {
	return &logFileRotationListenerDataSource{}
}

// logFileRotationListenerDataSource is the datasource implementation.
type logFileRotationListenerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logFileRotationListenerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_file_rotation_listener"
}

// Configure adds the provider configured client to the data source.
func (r *logFileRotationListenerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logFileRotationListenerDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	Type              types.String `tfsdk:"type"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	CopyToDirectory   types.String `tfsdk:"copy_to_directory"`
	CompressOnCopy    types.Bool   `tfsdk:"compress_on_copy"`
	OutputDirectory   types.String `tfsdk:"output_directory"`
	Description       types.String `tfsdk:"description"`
	Enabled           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *logFileRotationListenerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log File Rotation Listener.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log File Rotation Listener resource. Options are ['summarize', 'copy', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Log File Rotation Listener.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Log File Rotation Listener. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"copy_to_directory": schema.StringAttribute{
				Description: "The path to the directory to which log files should be copied. It must be different from the directory to which the log file is originally written, and administrators should ensure that the filesystem has sufficient space to hold files as they are copied.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"compress_on_copy": schema.BoolAttribute{
				Description: "Indicates whether the file should be gzip-compressed as it is copied into the destination directory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"output_directory": schema.StringAttribute{
				Description: "The path to the directory in which the summarize-access-log output should be written. If no value is provided, the output file will be written into the same directory as the rotated log file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log File Rotation Listener",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Log File Rotation Listener is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SummarizeLogFileRotationListenerResponse object into the model struct
func readSummarizeLogFileRotationListenerResponseDataSource(ctx context.Context, r *client.SummarizeLogFileRotationListenerResponse, state *logFileRotationListenerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("summarize")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OutputDirectory = internaltypes.StringTypeOrNil(r.OutputDirectory, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CopyLogFileRotationListenerResponse object into the model struct
func readCopyLogFileRotationListenerResponseDataSource(ctx context.Context, r *client.CopyLogFileRotationListenerResponse, state *logFileRotationListenerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("copy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CopyToDirectory = types.StringValue(r.CopyToDirectory)
	state.CompressOnCopy = internaltypes.BoolTypeOrNil(r.CompressOnCopy)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyLogFileRotationListenerResponse object into the model struct
func readThirdPartyLogFileRotationListenerResponseDataSource(ctx context.Context, r *client.ThirdPartyLogFileRotationListenerResponse, state *logFileRotationListenerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *logFileRotationListenerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logFileRotationListenerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SummarizeLogFileRotationListenerResponse != nil {
		readSummarizeLogFileRotationListenerResponseDataSource(ctx, readResponse.SummarizeLogFileRotationListenerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CopyLogFileRotationListenerResponse != nil {
		readCopyLogFileRotationListenerResponseDataSource(ctx, readResponse.CopyLogFileRotationListenerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyLogFileRotationListenerResponse != nil {
		readThirdPartyLogFileRotationListenerResponseDataSource(ctx, readResponse.ThirdPartyLogFileRotationListenerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
