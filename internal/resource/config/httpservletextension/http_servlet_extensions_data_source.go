package httpservletextension

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
	_ datasource.DataSource              = &httpServletExtensionsDataSource{}
	_ datasource.DataSourceWithConfigure = &httpServletExtensionsDataSource{}
)

// Create a Http Servlet Extensions data source
func NewHttpServletExtensionsDataSource() datasource.DataSource {
	return &httpServletExtensionsDataSource{}
}

// httpServletExtensionsDataSource is the datasource implementation.
type httpServletExtensionsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *httpServletExtensionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_servlet_extensions"
}

// Configure adds the provider configured client to the data source.
func (r *httpServletExtensionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type httpServletExtensionsDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *httpServletExtensionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Http Servlet Extension objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder name of this object required by Terraform.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "Http Servlet Extension objects found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: internaltypes.ObjectsObjectType(),
			},
		},
	}
}

// Read resource information
func (r *httpServletExtensionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state httpServletExtensionsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.HttpServletExtensionApi.ListHttpServletExtensions(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.ListHttpServletExtensionsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Http Servlet Extension objects", err, httpResp)
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
		if response.StandardHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.StandardHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("standard")
		}
		if response.DelegatedAdminHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.DelegatedAdminHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("delegated-admin")
		}
		if response.QuickstartHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.QuickstartHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("quickstart")
		}
		if response.AvailabilityStateHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.AvailabilityStateHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("availability-state")
		}
		if response.PrometheusMonitoringHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.PrometheusMonitoringHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("prometheus-monitoring")
		}
		if response.VelocityHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.VelocityHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("velocity")
		}
		if response.ConsentHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.ConsentHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("consent")
		}
		if response.LdapMappedScimHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.LdapMappedScimHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("ldap-mapped-scim")
		}
		if response.GroovyScriptedHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.GroovyScriptedHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("groovy-scripted")
		}
		if response.OpenBankingHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.OpenBankingHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("open-banking")
		}
		if response.PdpEndpointHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.PdpEndpointHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("pdp-endpoint")
		}
		if response.FileServerHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.FileServerHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("file-server")
		}
		if response.JsonPdpApiHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.JsonPdpApiHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("json-pdp-api")
		}
		if response.MetricsHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.MetricsHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("metrics")
		}
		if response.ConfigHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.ConfigHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("config")
		}
		if response.Scim2HttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.Scim2HttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("scim2")
		}
		if response.DirectoryRestApiHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.DirectoryRestApiHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("directory-rest-api")
		}
		if response.GatewayHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.GatewayHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("gateway")
		}
		if response.ThirdPartyHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.ThirdPartyHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("third-party")
		}
		if response.SidebandApiHttpServletExtensionResponse != nil {
			attributes["id"] = types.StringValue(response.SidebandApiHttpServletExtensionResponse.Id)
			attributes["type"] = types.StringValue("sideband-api")
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
