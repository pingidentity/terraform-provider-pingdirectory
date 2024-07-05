package postldifexporttaskprocessor

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &postLdifExportTaskProcessorDataSource{}
	_ datasource.DataSourceWithConfigure = &postLdifExportTaskProcessorDataSource{}
)

// Create a Post Ldif Export Task Processor data source
func NewPostLdifExportTaskProcessorDataSource() datasource.DataSource {
	return &postLdifExportTaskProcessorDataSource{}
}

// postLdifExportTaskProcessorDataSource is the datasource implementation.
type postLdifExportTaskProcessorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *postLdifExportTaskProcessorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_post_ldif_export_task_processor"
}

// Configure adds the provider configured client to the data source.
func (r *postLdifExportTaskProcessorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type postLdifExportTaskProcessorDataSourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	Name                                 types.String `tfsdk:"name"`
	Type                                 types.String `tfsdk:"type"`
	ExtensionClass                       types.String `tfsdk:"extension_class"`
	ExtensionArgument                    types.Set    `tfsdk:"extension_argument"`
	AwsExternalServer                    types.String `tfsdk:"aws_external_server"`
	S3BucketName                         types.String `tfsdk:"s3_bucket_name"`
	TargetThroughputInMegabitsPerSecond  types.Int64  `tfsdk:"target_throughput_in_megabits_per_second"`
	MaximumConcurrentTransferConnections types.Int64  `tfsdk:"maximum_concurrent_transfer_connections"`
	MaximumFileCountToRetain             types.Int64  `tfsdk:"maximum_file_count_to_retain"`
	MaximumFileAgeToRetain               types.String `tfsdk:"maximum_file_age_to_retain"`
	FileRetentionPattern                 types.String `tfsdk:"file_retention_pattern"`
	Description                          types.String `tfsdk:"description"`
	Enabled                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *postLdifExportTaskProcessorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Post Ldif Export Task Processor. Supported in PingDirectory product version 10.0.0.0+.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Post LDIF Export Task Processor resource. Options are ['upload-to-s3', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Post LDIF Export Task Processor.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Post LDIF Export Task Processor. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS S3 service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"s3_bucket_name": schema.StringAttribute{
				Description: "The name of the S3 bucket into which LDIF files should be copied.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"target_throughput_in_megabits_per_second": schema.Int64Attribute{
				Description: "The target throughput to attempt to achieve for data transfers to or from S3, in megabits per second.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_concurrent_transfer_connections": schema.Int64Attribute{
				Description: "The maximum number of concurrent connections that may be used when transferring data to or from S3.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_file_count_to_retain": schema.Int64Attribute{
				Description: "The maximum number of existing files matching the file retention pattern that should be retained in the S3 bucket after successfully uploading a newly exported file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_file_age_to_retain": schema.StringAttribute{
				Description: "The maximum length of time to retain files matching the file retention pattern that should be retained in the S3 bucket after successfully uploading a newly exported file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"file_retention_pattern": schema.StringAttribute{
				Description: "A regular expression pattern that will be used to identify which files are candidates for automatic removal based on the maximum-file-count-to-retain and maximum-file-age-to-retain properties. By default, all files in the bucket will be eligible for removal by retention processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Post LDIF Export Task Processor",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Post LDIF Export Task Processor is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any version restrictions are met
func (r *postLdifExportTaskProcessorDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	if r.providerConfig.ProductVersion != "" {
		version.CheckResourceSupported(&resp.Diagnostics, version.PingDirectory10000,
			r.providerConfig.ProductVersion, "pingdirectory_post_ldif_export_task_processor")
	}
}

// Read a UploadToS3PostLdifExportTaskProcessorResponse object into the model struct
func readUploadToS3PostLdifExportTaskProcessorResponseDataSource(ctx context.Context, r *client.UploadToS3PostLdifExportTaskProcessorResponse, state *postLdifExportTaskProcessorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("upload-to-s3")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.S3BucketName = types.StringValue(r.S3BucketName)
	state.TargetThroughputInMegabitsPerSecond = internaltypes.Int64TypeOrNil(r.TargetThroughputInMegabitsPerSecond)
	state.MaximumConcurrentTransferConnections = internaltypes.Int64TypeOrNil(r.MaximumConcurrentTransferConnections)
	state.MaximumFileCountToRetain = internaltypes.Int64TypeOrNil(r.MaximumFileCountToRetain)
	state.MaximumFileAgeToRetain = internaltypes.StringTypeOrNil(r.MaximumFileAgeToRetain, false)
	state.FileRetentionPattern = internaltypes.StringTypeOrNil(r.FileRetentionPattern, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyPostLdifExportTaskProcessorResponse object into the model struct
func readThirdPartyPostLdifExportTaskProcessorResponseDataSource(ctx context.Context, r *client.ThirdPartyPostLdifExportTaskProcessorResponse, state *postLdifExportTaskProcessorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *postLdifExportTaskProcessorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state postLdifExportTaskProcessorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PostLdifExportTaskProcessorAPI.GetPostLdifExportTaskProcessor(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Post Ldif Export Task Processor", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.UploadToS3PostLdifExportTaskProcessorResponse != nil {
		readUploadToS3PostLdifExportTaskProcessorResponseDataSource(ctx, readResponse.UploadToS3PostLdifExportTaskProcessorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPostLdifExportTaskProcessorResponse != nil {
		readThirdPartyPostLdifExportTaskProcessorResponseDataSource(ctx, readResponse.ThirdPartyPostLdifExportTaskProcessorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
