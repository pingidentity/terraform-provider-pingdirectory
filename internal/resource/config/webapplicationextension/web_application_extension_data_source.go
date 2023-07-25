package webapplicationextension

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
	_ datasource.DataSource              = &webApplicationExtensionDataSource{}
	_ datasource.DataSourceWithConfigure = &webApplicationExtensionDataSource{}
)

// Create a Web Application Extension data source
func NewWebApplicationExtensionDataSource() datasource.DataSource {
	return &webApplicationExtensionDataSource{}
}

// webApplicationExtensionDataSource is the datasource implementation.
type webApplicationExtensionDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *webApplicationExtensionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_application_extension"
}

// Configure adds the provider configured client to the data source.
func (r *webApplicationExtensionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type webApplicationExtensionDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Type                     types.String `tfsdk:"type"`
	Description              types.String `tfsdk:"description"`
	BaseContextPath          types.String `tfsdk:"base_context_path"`
	WarFile                  types.String `tfsdk:"war_file"`
	DocumentRootDirectory    types.String `tfsdk:"document_root_directory"`
	DeploymentDescriptorFile types.String `tfsdk:"deployment_descriptor_file"`
	TemporaryDirectory       types.String `tfsdk:"temporary_directory"`
	InitParameter            types.Set    `tfsdk:"init_parameter"`
}

// GetSchema defines the schema for the datasource.
func (r *webApplicationExtensionDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Web Application Extension.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of Web Application Extension resource. Options are ['console', 'generic']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Web Application Extension",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and at least one additional character and must represent a valid HTTP context path.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"war_file": schema.StringAttribute{
				Description: "Specifies the path to a standard web application archive (WAR) file.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this Web Application Extension. The path must exist, and it must be a directory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"deployment_descriptor_file": schema.StringAttribute{
				Description: "Specifies the path to the deployment descriptor file when used with document-root-directory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory that may be used to store temporary files such as extracted WAR files and compiled JSP files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"init_parameter": schema.SetAttribute{
				Description: "Specifies an initialization parameter to pass into the web application during startup.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read a GenericWebApplicationExtensionResponse object into the model struct
func readGenericWebApplicationExtensionResponseDataSource(ctx context.Context, r *client.GenericWebApplicationExtensionResponse, state *webApplicationExtensionDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, false)
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, false)
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, false)
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, false)
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
}

// Read resource information
func (r *webApplicationExtensionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state webApplicationExtensionDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GenericWebApplicationExtensionResponse != nil {
		readGenericWebApplicationExtensionResponseDataSource(ctx, readResponse.GenericWebApplicationExtensionResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
