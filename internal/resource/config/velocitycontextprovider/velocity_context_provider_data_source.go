package velocitycontextprovider

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
	_ datasource.DataSource              = &velocityContextProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &velocityContextProviderDataSource{}
)

// Create a Velocity Context Provider data source
func NewVelocityContextProviderDataSource() datasource.DataSource {
	return &velocityContextProviderDataSource{}
}

// velocityContextProviderDataSource is the datasource implementation.
type velocityContextProviderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *velocityContextProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_velocity_context_provider"
}

// Configure adds the provider configured client to the data source.
func (r *velocityContextProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type velocityContextProviderDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	Type                     types.String `tfsdk:"type"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
	ExtensionClass           types.String `tfsdk:"extension_class"`
	ExtensionArgument        types.Set    `tfsdk:"extension_argument"`
	RequestTool              types.Set    `tfsdk:"request_tool"`
	SessionTool              types.Set    `tfsdk:"session_tool"`
	ApplicationTool          types.Set    `tfsdk:"application_tool"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	HttpMethod               types.Set    `tfsdk:"http_method"`
	ObjectScope              types.String `tfsdk:"object_scope"`
	IncludedView             types.Set    `tfsdk:"included_view"`
	ExcludedView             types.Set    `tfsdk:"excluded_view"`
	ResponseHeader           types.Set    `tfsdk:"response_header"`
}

// GetSchema defines the schema for the datasource.
func (r *velocityContextProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Velocity Context Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Velocity Context Provider resource. Options are ['velocity-tools', 'custom', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Velocity Context Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Velocity Context Provider. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"request_tool": schema.SetAttribute{
				Description: "The fully-qualified name of a Velocity Tool class that will be initialized for each request. May optionally include a path to a properties file used to configure this tool separated from the class name by a semi-colon (;). The path may absolute or relative to the server root.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"session_tool": schema.SetAttribute{
				Description: "The fully-qualified name of a Velocity Tool class that will be initialized for each session. May optionally include a path to a properties file used to configure this tool separated from the class name by a semi-colon (;). The path may absolute or relative to the server root.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"application_tool": schema.SetAttribute{
				Description: "The fully-qualified name of a Velocity Tool class that will be initialized once for the life of the server. May optionally include a path to a properties file used to configure this tool separated from the class name by a semi-colon (;). The path may absolute or relative to the server root.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Velocity Context Provider is enabled. If set to 'false' this Velocity Context Provider will not contribute context content for any requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_method": schema.SetAttribute{
				Description: "Specifies the set of HTTP methods handled by this Velocity Context Provider, which will perform actions necessary to fulfill the request before updating the context for the response. The values of this property are not case-sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"object_scope": schema.StringAttribute{
				Description: "Scope for context objects contributed by this Velocity Context Provider. Must be either 'request' or 'session' or 'application'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_view": schema.SetAttribute{
				Description: "The name of a view for which this Velocity Context Provider will contribute content.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_view": schema.SetAttribute{
				Description: "The name of a view for which this Velocity Context Provider will not contribute content.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for template page requests to which this Velocity Context Provider contributes content.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a VelocityToolsVelocityContextProviderResponse object into the model struct
func readVelocityToolsVelocityContextProviderResponseDataSource(ctx context.Context, r *client.VelocityToolsVelocityContextProviderResponse, state *velocityContextProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("velocity-tools")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequestTool = internaltypes.GetStringSet(r.RequestTool)
	state.SessionTool = internaltypes.GetStringSet(r.SessionTool)
	state.ApplicationTool = internaltypes.GetStringSet(r.ApplicationTool)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.ObjectScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvelocityContextProviderObjectScopeProp(r.ObjectScope), false)
	state.IncludedView = internaltypes.GetStringSet(r.IncludedView)
	state.ExcludedView = internaltypes.GetStringSet(r.ExcludedView)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
}

// Read a CustomVelocityContextProviderResponse object into the model struct
func readCustomVelocityContextProviderResponseDataSource(ctx context.Context, r *client.CustomVelocityContextProviderResponse, state *velocityContextProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.ObjectScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvelocityContextProviderObjectScopeProp(r.ObjectScope), false)
	state.IncludedView = internaltypes.GetStringSet(r.IncludedView)
	state.ExcludedView = internaltypes.GetStringSet(r.ExcludedView)
	state.HttpMethod = internaltypes.GetStringSet(r.HttpMethod)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
}

// Read a ThirdPartyVelocityContextProviderResponse object into the model struct
func readThirdPartyVelocityContextProviderResponseDataSource(ctx context.Context, r *client.ThirdPartyVelocityContextProviderResponse, state *velocityContextProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.ObjectScope = internaltypes.StringTypeOrNil(
		client.StringPointerEnumvelocityContextProviderObjectScopeProp(r.ObjectScope), false)
	state.IncludedView = internaltypes.GetStringSet(r.IncludedView)
	state.ExcludedView = internaltypes.GetStringSet(r.ExcludedView)
	state.HttpMethod = internaltypes.GetStringSet(r.HttpMethod)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
}

// Read resource information
func (r *velocityContextProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state velocityContextProviderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.VelocityContextProviderApi.GetVelocityContextProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.HttpServletExtensionName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Velocity Context Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.VelocityToolsVelocityContextProviderResponse != nil {
		readVelocityToolsVelocityContextProviderResponseDataSource(ctx, readResponse.VelocityToolsVelocityContextProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomVelocityContextProviderResponse != nil {
		readCustomVelocityContextProviderResponseDataSource(ctx, readResponse.CustomVelocityContextProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyVelocityContextProviderResponse != nil {
		readThirdPartyVelocityContextProviderResponseDataSource(ctx, readResponse.ThirdPartyVelocityContextProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
