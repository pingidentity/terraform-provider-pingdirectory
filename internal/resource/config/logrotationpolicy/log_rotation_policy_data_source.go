package logrotationpolicy

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
	_ datasource.DataSource              = &logRotationPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &logRotationPolicyDataSource{}
)

// Create a Log Rotation Policy data source
func NewLogRotationPolicyDataSource() datasource.DataSource {
	return &logRotationPolicyDataSource{}
}

// logRotationPolicyDataSource is the datasource implementation.
type logRotationPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logRotationPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_rotation_policy"
}

// Configure adds the provider configured client to the data source.
func (r *logRotationPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logRotationPolicyDataSourceModel struct {
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Type             types.String `tfsdk:"type"`
	FileSizeLimit    types.String `tfsdk:"file_size_limit"`
	TimeOfDay        types.Set    `tfsdk:"time_of_day"`
	RotationInterval types.String `tfsdk:"rotation_interval"`
	Description      types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the datasource.
func (r *logRotationPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log Rotation Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Rotation Policy resource. Options are ['time-limit', 'fixed-time', 'never-rotate', 'size-limit']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"file_size_limit": schema.StringAttribute{
				Description: "Specifies the maximum size that a log file can reach before it is rotated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"time_of_day": schema.SetAttribute{
				Description: "Specifies the time of day at which log rotation should occur.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"rotation_interval": schema.StringAttribute{
				Description: "Specifies the time interval between rotations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Rotation Policy",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a TimeLimitLogRotationPolicyResponse object into the model struct
func readTimeLimitLogRotationPolicyResponseDataSource(ctx context.Context, r *client.TimeLimitLogRotationPolicyResponse, state *logRotationPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("time-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RotationInterval = types.StringValue(r.RotationInterval)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a FixedTimeLogRotationPolicyResponse object into the model struct
func readFixedTimeLogRotationPolicyResponseDataSource(ctx context.Context, r *client.FixedTimeLogRotationPolicyResponse, state *logRotationPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("fixed-time")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.TimeOfDay = internaltypes.GetStringSet(r.TimeOfDay)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a NeverRotateLogRotationPolicyResponse object into the model struct
func readNeverRotateLogRotationPolicyResponseDataSource(ctx context.Context, r *client.NeverRotateLogRotationPolicyResponse, state *logRotationPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("never-rotate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SizeLimitLogRotationPolicyResponse object into the model struct
func readSizeLimitLogRotationPolicyResponseDataSource(ctx context.Context, r *client.SizeLimitLogRotationPolicyResponse, state *logRotationPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("size-limit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.FileSizeLimit = types.StringValue(r.FileSizeLimit)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *logRotationPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logRotationPolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRotationPolicyApi.GetLogRotationPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Rotation Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TimeLimitLogRotationPolicyResponse != nil {
		readTimeLimitLogRotationPolicyResponseDataSource(ctx, readResponse.TimeLimitLogRotationPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FixedTimeLogRotationPolicyResponse != nil {
		readFixedTimeLogRotationPolicyResponseDataSource(ctx, readResponse.FixedTimeLogRotationPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NeverRotateLogRotationPolicyResponse != nil {
		readNeverRotateLogRotationPolicyResponseDataSource(ctx, readResponse.NeverRotateLogRotationPolicyResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SizeLimitLogRotationPolicyResponse != nil {
		readSizeLimitLogRotationPolicyResponseDataSource(ctx, readResponse.SizeLimitLogRotationPolicyResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
