package logretentionpolicy

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
	_ datasource.DataSource              = &logRetentionPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &logRetentionPolicyDataSource{}
)

// Create a Log Retention Policy data source
func NewLogRetentionPolicyDataSource() datasource.DataSource {
	return &logRetentionPolicyDataSource{}
}

// logRetentionPolicyDataSource is the datasource implementation.
type logRetentionPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logRetentionPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_retention_policy"
}

// Configure adds the provider configured client to the data source.
func (r *logRetentionPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type logRetentionPolicyDataSourceModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Type           types.String `tfsdk:"type"`
	DiskSpaceUsed  types.String `tfsdk:"disk_space_used"`
	FreeDiskSpace  types.String `tfsdk:"free_disk_space"`
	NumberOfFiles  types.Int64  `tfsdk:"number_of_files"`
	RetainDuration types.String `tfsdk:"retain_duration"`
	Description    types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *logRetentionPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log Retention Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Retention Policy resource. Options are ['time-limit', 'never-delete', 'file-count', 'free-disk-space', 'size-limit']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"disk_space_used": schema.StringAttribute{
				Description: "Specifies the maximum total disk space used by the log files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"free_disk_space": schema.StringAttribute{
				Description: "Specifies the minimum amount of free disk space that should be available on the file system on which the archived log files are stored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"number_of_files": schema.Int64Attribute{
				Description: "Specifies the number of archived log files to retain before the oldest ones are cleaned.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"retain_duration": schema.StringAttribute{
				Description: "Specifies the desired minimum length of time that each log file should be retained.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Retention Policy",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a TimeLimitLogRetentionPolicyResponse object into the model struct
func readTimeLimitLogRetentionPolicyResponseDataSource(ctx context.Context, r *client.TimeLimitLogRetentionPolicyResponse, state *logRetentionPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("time-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RetainDuration = types.StringValue(r.RetainDuration)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a NeverDeleteLogRetentionPolicyResponse object into the model struct
func readNeverDeleteLogRetentionPolicyResponseDataSource(ctx context.Context, r *client.NeverDeleteLogRetentionPolicyResponse, state *logRetentionPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("never-delete")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a FileCountLogRetentionPolicyResponse object into the model struct
func readFileCountLogRetentionPolicyResponseDataSource(ctx context.Context, r *client.FileCountLogRetentionPolicyResponse, state *logRetentionPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-count")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.NumberOfFiles = types.Int64Value(r.NumberOfFiles)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a FreeDiskSpaceLogRetentionPolicyResponse object into the model struct
func readFreeDiskSpaceLogRetentionPolicyResponseDataSource(ctx context.Context, r *client.FreeDiskSpaceLogRetentionPolicyResponse, state *logRetentionPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("free-disk-space")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.FreeDiskSpace = types.StringValue(r.FreeDiskSpace)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SizeLimitLogRetentionPolicyResponse object into the model struct
func readSizeLimitLogRetentionPolicyResponseDataSource(ctx context.Context, r *client.SizeLimitLogRetentionPolicyResponse, state *logRetentionPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("size-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DiskSpaceUsed = types.StringValue(r.DiskSpaceUsed)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *logRetentionPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logRetentionPolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TimeLimitLogRetentionPolicyResponse != nil {
		readTimeLimitLogRetentionPolicyResponseDataSource(ctx, readResponse.TimeLimitLogRetentionPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NeverDeleteLogRetentionPolicyResponse != nil {
		readNeverDeleteLogRetentionPolicyResponseDataSource(ctx, readResponse.NeverDeleteLogRetentionPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileCountLogRetentionPolicyResponse != nil {
		readFileCountLogRetentionPolicyResponseDataSource(ctx, readResponse.FileCountLogRetentionPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FreeDiskSpaceLogRetentionPolicyResponse != nil {
		readFreeDiskSpaceLogRetentionPolicyResponseDataSource(ctx, readResponse.FreeDiskSpaceLogRetentionPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SizeLimitLogRetentionPolicyResponse != nil {
		readSizeLimitLogRetentionPolicyResponseDataSource(ctx, readResponse.SizeLimitLogRetentionPolicyResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
