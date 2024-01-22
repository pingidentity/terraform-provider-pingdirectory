package velocitytemplateloader

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
	_ datasource.DataSource              = &velocityTemplateLoaderDataSource{}
	_ datasource.DataSourceWithConfigure = &velocityTemplateLoaderDataSource{}
)

// Create a Velocity Template Loader data source
func NewVelocityTemplateLoaderDataSource() datasource.DataSource {
	return &velocityTemplateLoaderDataSource{}
}

// velocityTemplateLoaderDataSource is the datasource implementation.
type velocityTemplateLoaderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *velocityTemplateLoaderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_velocity_template_loader"
}

// Configure adds the provider configured client to the data source.
func (r *velocityTemplateLoaderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type velocityTemplateLoaderDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Type                     types.String `tfsdk:"type"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	EvaluationOrderIndex     types.Int64  `tfsdk:"evaluation_order_index"`
	MimeTypeMatcher          types.String `tfsdk:"mime_type_matcher"`
	MimeType                 types.String `tfsdk:"mime_type"`
	TemplateSuffix           types.String `tfsdk:"template_suffix"`
	TemplateDirectory        types.String `tfsdk:"template_directory"`
}

// GetSchema defines the schema for the datasource.
func (r *velocityTemplateLoaderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Velocity Template Loader.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Velocity Template Loader resource. Options are ['velocity-template-loader']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Velocity Template Loader is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "This property determines the evaluation order for determining the correct Velocity Template Loader to load a template for generating content for a particular request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mime_type_matcher": schema.StringAttribute{
				Description: "Specifies a media type for matching Accept request-header values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mime_type": schema.StringAttribute{
				Description: "Specifies a the value that will be used in the response's Content-Type header that indicates the type of content to return.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"template_suffix": schema.StringAttribute{
				Description: "Specifies the suffix to append to the requested resource name when searching for the template file with which to form a response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"template_directory": schema.StringAttribute{
				Description: "Specifies the directory in which to search for the template files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a VelocityTemplateLoaderResponse object into the model struct
func readVelocityTemplateLoaderResponseDataSource(ctx context.Context, r *client.VelocityTemplateLoaderResponse, state *velocityTemplateLoaderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("velocity-template-loader")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.MimeTypeMatcher = types.StringValue(r.MimeTypeMatcher)
	state.MimeType = internaltypes.StringTypeOrNil(r.MimeType, false)
	state.TemplateSuffix = internaltypes.StringTypeOrNil(r.TemplateSuffix, false)
	state.TemplateDirectory = internaltypes.StringTypeOrNil(r.TemplateDirectory, false)
}

// Read resource information
func (r *velocityTemplateLoaderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state velocityTemplateLoaderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VelocityTemplateLoaderAPI.GetVelocityTemplateLoader(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Template Loader", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readVelocityTemplateLoaderResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
