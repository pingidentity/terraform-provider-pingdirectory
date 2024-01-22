package externalserver

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
	_ datasource.DataSource              = &externalServersDataSource{}
	_ datasource.DataSourceWithConfigure = &externalServersDataSource{}
)

// Create a External Servers data source
func NewExternalServersDataSource() datasource.DataSource {
	return &externalServersDataSource{}
}

// externalServersDataSource is the datasource implementation.
type externalServersDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *externalServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_external_servers"
}

// Configure adds the provider configured client to the data source.
func (r *externalServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type externalServersDataSourceModel struct {
	Id      types.String `tfsdk:"id"`
	Filter  types.String `tfsdk:"filter"`
	Objects types.Set    `tfsdk:"objects"`
}

// GetSchema defines the schema for the datasource.
func (r *externalServersDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Lists External Server objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"objects": schema.SetAttribute{
				Description: "External Server objects found in the configuration",
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
func (r *externalServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state externalServersDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.ExternalServerAPI.ListExternalServers(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerAPI.ListExternalServersExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the External Server objects", err, httpResp)
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
		if response.ConsentServiceExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.ConsentServiceExternalServerResponse.Id)
			attributes["type"] = types.StringValue("consent-service")
		}
		if response.ScimExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.ScimExternalServerResponse.Id)
			attributes["type"] = types.StringValue("scim")
		}
		if response.NokiaDsExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.NokiaDsExternalServerResponse.Id)
			attributes["type"] = types.StringValue("nokia-ds")
		}
		if response.PingIdentityDsExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.PingIdentityDsExternalServerResponse.Id)
			attributes["type"] = types.StringValue("ping-identity-ds")
		}
		if response.MetricsEngineExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.MetricsEngineExternalServerResponse.Id)
			attributes["type"] = types.StringValue("metrics-engine")
		}
		if response.JdbcExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.JdbcExternalServerResponse.Id)
			attributes["type"] = types.StringValue("jdbc")
		}
		if response.SyslogExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.SyslogExternalServerResponse.Id)
			attributes["type"] = types.StringValue("syslog")
		}
		if response.PingIdentityProxyServerExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.PingIdentityProxyServerExternalServerResponse.Id)
			attributes["type"] = types.StringValue("ping-identity-proxy-server")
		}
		if response.NokiaProxyServerExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.NokiaProxyServerExternalServerResponse.Id)
			attributes["type"] = types.StringValue("nokia-proxy-server")
		}
		if response.SunDsExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.SunDsExternalServerResponse.Id)
			attributes["type"] = types.StringValue("sun-ds")
		}
		if response.OpendjExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.OpendjExternalServerResponse.Id)
			attributes["type"] = types.StringValue("opendj")
		}
		if response.LdapExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.LdapExternalServerResponse.Id)
			attributes["type"] = types.StringValue("ldap")
		}
		if response.PingOneHttpExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.PingOneHttpExternalServerResponse.Id)
			attributes["type"] = types.StringValue("ping-one-http")
		}
		if response.ApiExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.ApiExternalServerResponse.Id)
			attributes["type"] = types.StringValue("api")
		}
		if response.RedHatDsExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.RedHatDsExternalServerResponse.Id)
			attributes["type"] = types.StringValue("red-hat-ds")
		}
		if response.SyncServerExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.SyncServerExternalServerResponse.Id)
			attributes["type"] = types.StringValue("sync-server")
		}
		if response.VaultExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.VaultExternalServerResponse.Id)
			attributes["type"] = types.StringValue("vault")
		}
		if response.PolicyExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.PolicyExternalServerResponse.Id)
			attributes["type"] = types.StringValue("policy")
		}
		if response.SmtpExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.SmtpExternalServerResponse.Id)
			attributes["type"] = types.StringValue("smtp")
		}
		if response.ActiveDirectoryExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.ActiveDirectoryExternalServerResponse.Id)
			attributes["type"] = types.StringValue("active-directory")
		}
		if response.BrokerExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.BrokerExternalServerResponse.Id)
			attributes["type"] = types.StringValue("broker")
		}
		if response.HttpProxyExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.HttpProxyExternalServerResponse.Id)
			attributes["type"] = types.StringValue("http-proxy")
		}
		if response.KafkaClusterExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.KafkaClusterExternalServerResponse.Id)
			attributes["type"] = types.StringValue("kafka-cluster")
		}
		if response.HttpExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.HttpExternalServerResponse.Id)
			attributes["type"] = types.StringValue("http")
		}
		if response.MockExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.MockExternalServerResponse.Id)
			attributes["type"] = types.StringValue("mock")
		}
		if response.OracleUnifiedDirectoryExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.OracleUnifiedDirectoryExternalServerResponse.Id)
			attributes["type"] = types.StringValue("oracle-unified-directory")
		}
		if response.ConjurExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.ConjurExternalServerResponse.Id)
			attributes["type"] = types.StringValue("conjur")
		}
		if response.AmazonAwsExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.AmazonAwsExternalServerResponse.Id)
			attributes["type"] = types.StringValue("amazon-aws")
		}
		if response.Scim2ExternalServerResponse != nil {
			attributes["id"] = types.StringValue(response.Scim2ExternalServerResponse.Id)
			attributes["type"] = types.StringValue("scim2")
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
