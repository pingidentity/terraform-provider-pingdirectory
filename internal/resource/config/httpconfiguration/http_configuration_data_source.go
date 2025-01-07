package httpconfiguration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &httpConfigurationDataSource{}
	_ datasource.DataSourceWithConfigure = &httpConfigurationDataSource{}
)

// Create a Http Configuration data source
func NewHttpConfigurationDataSource() datasource.DataSource {
	return &httpConfigurationDataSource{}
}

// httpConfigurationDataSource is the datasource implementation.
type httpConfigurationDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *httpConfigurationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_configuration"
}

// Configure adds the provider configured client to the data source.
func (r *httpConfigurationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type httpConfigurationDataSourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Type                                  types.String `tfsdk:"type"`
	IncludeStackTracesInErrorPages        types.Bool   `tfsdk:"include_stack_traces_in_error_pages"`
	IncludeServletInformationInErrorPages types.Bool   `tfsdk:"include_servlet_information_in_error_pages"`
}

// GetSchema defines the schema for the datasource.
func (r *httpConfigurationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Http Configuration.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of HTTP Configuration resource. Options are ['http-configuration']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_stack_traces_in_error_pages": schema.BoolAttribute{
				Description: "Indicates whether exceptions thrown by servlet or web application extensions will be included in the resulting error page response. Stack traces can be helpful in diagnosing application errors, but in production they may reveal information that might be useful to a malicious attacker.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_servlet_information_in_error_pages": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. Indicates whether to expose servlet information in the error page response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a HttpConfigurationResponse object into the model struct
func readHttpConfigurationResponseDataSource(ctx context.Context, r *client.HttpConfigurationResponse, state *httpConfigurationDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http-configuration")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.IncludeStackTracesInErrorPages = internaltypes.BoolTypeOrNil(r.IncludeStackTracesInErrorPages)
	state.IncludeServletInformationInErrorPages = internaltypes.BoolTypeOrNil(r.IncludeServletInformationInErrorPages)
}

// Read resource information
func (r *httpConfigurationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state httpConfigurationDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpConfigurationAPI.GetHttpConfiguration(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Configuration", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpConfigurationResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
