package velocitycontextprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &velocityContextProvidersDataSource{}
	_ datasource.DataSourceWithConfigure = &velocityContextProvidersDataSource{}
)

// Create a Velocity Context Providers data source
func NewVelocityContextProvidersDataSource() datasource.DataSource {
	return &velocityContextProvidersDataSource{}
}

// velocityContextProvidersDataSource is the datasource implementation.
type velocityContextProvidersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *velocityContextProvidersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_velocity_context_providers"
}

// Configure adds the provider configured client to the data source.
func (r *velocityContextProvidersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type velocityContextProvidersDataSourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Filter                   types.String `tfsdk:"filter"`
	Objects                  types.Set    `tfsdk:"objects"`
	HttpServletExtensionName types.String `tfsdk:"http_servlet_extension_name"`
}

// GetSchema defines the schema for the datasource.
func (r *velocityContextProvidersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists Velocity Context Provider objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"http_servlet_extension_name": schema.StringAttribute{
				Description: "Name of the parent HTTP Servlet Extension",
				Required:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Velocity Context Provider objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read resource information
func (r *velocityContextProvidersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state velocityContextProvidersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.VelocityContextProviderAPI.ListVelocityContextProviders(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.HttpServletExtensionName.ValueString())
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.VelocityContextProviderAPI.ListVelocityContextProvidersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Velocity Context Provider objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	objects := []attr.Value{}
	for _, response := range readResponse.Resources {
		attributes := map[string]attr.Value{}
		if response.VelocityToolsVelocityContextProviderResponse != nil {
			attributes["id"] = types.StringValue(response.VelocityToolsVelocityContextProviderResponse.Id)
			attributes["type"] = types.StringValue("velocity-tools")
		}
		if response.CustomVelocityContextProviderResponse != nil {
			attributes["id"] = types.StringValue(response.CustomVelocityContextProviderResponse.Id)
			attributes["type"] = types.StringValue("custom")
		}
		if response.ThirdPartyVelocityContextProviderResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyVelocityContextProviderResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		obj, diags := types.ObjectValue(internaltypes.ObjectsAttrTypes(), attributes)
		resp.Diagnostics.Append(diags...)
		objects = append(objects, obj)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	state.Objects, diags = types.SetValue(internaltypes.ObjectsObjectType(), objects)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
