package serverinstancelistener

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
	_ datasource.DataSource              = &serverInstanceListenerDataSource{}
	_ datasource.DataSourceWithConfigure = &serverInstanceListenerDataSource{}
)

// Create a Server Instance Listener data source
func NewServerInstanceListenerDataSource() datasource.DataSource {
	return &serverInstanceListenerDataSource{}
}

// serverInstanceListenerDataSource is the datasource implementation.
type serverInstanceListenerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *serverInstanceListenerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_instance_listener"
}

// Configure adds the provider configured client to the data source.
func (r *serverInstanceListenerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type serverInstanceListenerDataSourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	Type                types.String `tfsdk:"type"`
	ServerInstanceName  types.String `tfsdk:"server_instance_name"`
	ListenAddress       types.String `tfsdk:"listen_address"`
	ServerHTTPPort      types.Int64  `tfsdk:"server_http_port"`
	ServerLDAPPort      types.Int64  `tfsdk:"server_ldap_port"`
	ConnectionSecurity  types.String `tfsdk:"connection_security"`
	ListenerCertificate types.String `tfsdk:"listener_certificate"`
	Purpose             types.Set    `tfsdk:"purpose"`
}

// GetSchema defines the schema for the datasource.
func (r *serverInstanceListenerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Server Instance Listener.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Server Instance Listener resource. Options are ['ldap', 'http']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
			},
			"listen_address": schema.StringAttribute{
				Description: "If the server is listening on a particular address different from the hostname, then this property may be used to specify the address on which to listen for connections from HTTP clients.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_http_port": schema.Int64Attribute{
				Description: "The TCP port number on which the HTTP server is listening.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_ldap_port": schema.Int64Attribute{
				Description: "The TCP port number on which the LDAP server is listening.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_security": schema.StringAttribute{
				Description: "Specifies the mechanism to use for securing connections to the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"listener_certificate": schema.StringAttribute{
				Description: "The public component of the certificate that the listener is expected to present to clients. When establishing a connection to this server, only the certificate(s) listed here will be trusted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"purpose": schema.SetAttribute{
				Description: "Identifies the purpose of this Server Instance Listener.",
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

// Read a LdapServerInstanceListenerResponse object into the model struct
func readLdapServerInstanceListenerResponseDataSource(ctx context.Context, r *client.LdapServerInstanceListenerResponse, state *serverInstanceListenerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerLDAPPort = internaltypes.Int64TypeOrNil(r.ServerLDAPPort)
	state.ConnectionSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceListenerLdapConnectionSecurityProp(r.ConnectionSecurity), false)
	state.ListenerCertificate = internaltypes.StringTypeOrNil(r.ListenerCertificate, false)
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnumserverInstanceListenerPurposeProp(r.Purpose))
}

// Read a HttpServerInstanceListenerResponse object into the model struct
func readHttpServerInstanceListenerResponseDataSource(ctx context.Context, r *client.HttpServerInstanceListenerResponse, state *serverInstanceListenerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ListenAddress = internaltypes.StringTypeOrNil(r.ListenAddress, false)
	state.ServerHTTPPort = internaltypes.Int64TypeOrNil(r.ServerHTTPPort)
	state.ConnectionSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceListenerHttpConnectionSecurityProp(r.ConnectionSecurity), false)
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnumserverInstanceListenerPurposeProp(r.Purpose))
}

// Read resource information
func (r *serverInstanceListenerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state serverInstanceListenerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceListenerAPI.GetServerInstanceListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LdapServerInstanceListenerResponse != nil {
		readLdapServerInstanceListenerResponseDataSource(ctx, readResponse.LdapServerInstanceListenerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HttpServerInstanceListenerResponse != nil {
		readHttpServerInstanceListenerResponseDataSource(ctx, readResponse.HttpServerInstanceListenerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
