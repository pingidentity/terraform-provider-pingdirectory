package httpservletcrossoriginpolicy

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
	_ datasource.DataSource              = &httpServletCrossOriginPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &httpServletCrossOriginPolicyDataSource{}
)

// Create a Http Servlet Cross Origin Policy data source
func NewHttpServletCrossOriginPolicyDataSource() datasource.DataSource {
	return &httpServletCrossOriginPolicyDataSource{}
}

// httpServletCrossOriginPolicyDataSource is the datasource implementation.
type httpServletCrossOriginPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *httpServletCrossOriginPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_servlet_cross_origin_policy"
}

// Configure adds the provider configured client to the data source.
func (r *httpServletCrossOriginPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type httpServletCrossOriginPolicyDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Type                 types.String `tfsdk:"type"`
	Description          types.String `tfsdk:"description"`
	CorsAllowedMethods   types.Set    `tfsdk:"cors_allowed_methods"`
	CorsAllowedOrigins   types.Set    `tfsdk:"cors_allowed_origins"`
	CorsExposedHeaders   types.Set    `tfsdk:"cors_exposed_headers"`
	CorsAllowedHeaders   types.Set    `tfsdk:"cors_allowed_headers"`
	CorsPreflightMaxAge  types.String `tfsdk:"cors_preflight_max_age"`
	CorsAllowCredentials types.Bool   `tfsdk:"cors_allow_credentials"`
}

// GetSchema defines the schema for the datasource.
func (r *httpServletCrossOriginPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Http Servlet Cross Origin Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of HTTP Servlet Cross Origin Policy resource. Options are ['http-servlet-cross-origin-policy']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Cross Origin Policy",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cors_allowed_methods": schema.SetAttribute{
				Description: "A list of HTTP methods allowed for cross-origin access to resources. i.e. one or more of GET, POST, PUT, DELETE, etc.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"cors_allowed_origins": schema.SetAttribute{
				Description: "A list of origins that are allowed to execute cross-origin requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"cors_exposed_headers": schema.SetAttribute{
				Description: "A list of HTTP headers other than the simple response headers that browsers are allowed to access.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"cors_allowed_headers": schema.SetAttribute{
				Description: "A list of HTTP headers that are supported by the resource and can be specified in a cross-origin request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"cors_preflight_max_age": schema.StringAttribute{
				Description: "The maximum amount of time that a preflight request can be cached by a client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cors_allow_credentials": schema.BoolAttribute{
				Description: "Indicates whether the servlet extension allows CORS requests with username/password credentials.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a HttpServletCrossOriginPolicyResponse object into the model struct
func readHttpServletCrossOriginPolicyResponseDataSource(ctx context.Context, r *client.HttpServletCrossOriginPolicyResponse, state *httpServletCrossOriginPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http-servlet-cross-origin-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.CorsAllowedMethods = internaltypes.GetStringSet(r.CorsAllowedMethods)
	state.CorsAllowedOrigins = internaltypes.GetStringSet(r.CorsAllowedOrigins)
	state.CorsExposedHeaders = internaltypes.GetStringSet(r.CorsExposedHeaders)
	state.CorsAllowedHeaders = internaltypes.GetStringSet(r.CorsAllowedHeaders)
	state.CorsPreflightMaxAge = internaltypes.StringTypeOrNil(r.CorsPreflightMaxAge, false)
	state.CorsAllowCredentials = internaltypes.BoolTypeOrNil(r.CorsAllowCredentials)
}

// Read resource information
func (r *httpServletCrossOriginPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state httpServletCrossOriginPolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletCrossOriginPolicyApi.GetHttpServletCrossOriginPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Cross Origin Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpServletCrossOriginPolicyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
