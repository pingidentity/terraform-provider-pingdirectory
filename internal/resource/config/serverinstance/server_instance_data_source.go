package serverinstance

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
	_ datasource.DataSource              = &serverInstanceDataSource{}
	_ datasource.DataSourceWithConfigure = &serverInstanceDataSource{}
)

// Create a Server Instance data source
func NewServerInstanceDataSource() datasource.DataSource {
	return &serverInstanceDataSource{}
}

// serverInstanceDataSource is the datasource implementation.
type serverInstanceDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *serverInstanceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_instance"
}

// Configure adds the provider configured client to the data source.
func (r *serverInstanceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type serverInstanceDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Type                       types.String `tfsdk:"type"`
	ServerInstanceType         types.String `tfsdk:"server_instance_type"`
	ReplicationSetName         types.String `tfsdk:"replication_set_name"`
	LoadBalancingAlgorithmName types.Set    `tfsdk:"load_balancing_algorithm_name"`
	ServerInstanceName         types.String `tfsdk:"server_instance_name"`
	ClusterName                types.String `tfsdk:"cluster_name"`
	ServerInstanceLocation     types.String `tfsdk:"server_instance_location"`
	Hostname                   types.String `tfsdk:"hostname"`
	ServerRoot                 types.String `tfsdk:"server_root"`
	ServerVersion              types.String `tfsdk:"server_version"`
	InterServerCertificate     types.String `tfsdk:"inter_server_certificate"`
	LdapPort                   types.Int64  `tfsdk:"ldap_port"`
	LdapsPort                  types.Int64  `tfsdk:"ldaps_port"`
	HttpPort                   types.Int64  `tfsdk:"http_port"`
	HttpsPort                  types.Int64  `tfsdk:"https_port"`
	ReplicationPort            types.Int64  `tfsdk:"replication_port"`
	ReplicationServerID        types.Int64  `tfsdk:"replication_server_id"`
	ReplicationDomainServerID  types.Set    `tfsdk:"replication_domain_server_id"`
	JmxPort                    types.Int64  `tfsdk:"jmx_port"`
	JmxsPort                   types.Int64  `tfsdk:"jmxs_port"`
	PreferredSecurity          types.String `tfsdk:"preferred_security"`
	StartTLSEnabled            types.Bool   `tfsdk:"start_tls_enabled"`
	BaseDN                     types.Set    `tfsdk:"base_dn"`
	MemberOfServerGroup        types.Set    `tfsdk:"member_of_server_group"`
}

// GetSchema defines the schema for the datasource.
func (r *serverInstanceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Server Instance.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Server Instance resource. Options are ['proxy', 'metrics-engine', 'authorize', 'directory', 'sync']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_instance_type": schema.StringAttribute{
				Description: "Specifies the type of server installation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_set_name": schema.StringAttribute{
				Description: "The name of the replication set assigned to this Directory Server. Restricted domains are only replicated within instances using the same replication set name.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"load_balancing_algorithm_name": schema.SetAttribute{
				Description: "The name of the configuration object for a load-balancing algorithm that should include this server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"server_instance_name": schema.StringAttribute{
				Description: "The name of this Server Instance. The instance name needs to be unique if this server will be part of a topology of servers that are connected to each other. Once set, it may not be changed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cluster_name": schema.StringAttribute{
				Description: "The name of the cluster to which this Server Instance belongs. Server instances within the same cluster will share the same cluster-wide configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_instance_location": schema.StringAttribute{
				Description: "Specifies the location for the Server Instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "The name of the host where this Server Instance is installed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_root": schema.StringAttribute{
				Description: "The file system path where this Server Instance is installed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_version": schema.StringAttribute{
				Description: "The version of the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"inter_server_certificate": schema.StringAttribute{
				Description: "The public component of the certificate used by this instance to protect inter-server communication and to perform server-specific encryption. This will generally be managed by the server and should only be altered by administrators under explicit direction from Ping Identity support personnel.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for LDAP connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldaps_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for LDAP secure connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for HTTP connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"https_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for HTTPS connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_port": schema.Int64Attribute{
				Description: "The replication TCP port.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the replication server on this server instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replication_domain_server_id": schema.SetAttribute{
				Description: "Specifies a unique identifier for the Directory Server within the replication domain.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.Int64Type,
			},
			"jmx_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for JMX connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jmxs_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for JMX secure connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"preferred_security": schema.StringAttribute{
				Description: "Specifies the preferred mechanism to use for securing connections to the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"start_tls_enabled": schema.BoolAttribute{
				Description: "Indicates whether StartTLS is enabled on this server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "The set of base DNs under the root DSE.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"member_of_server_group": schema.SetAttribute{
				Description: "The set of groups of which this server is a member.",
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

// Read a ProxyServerInstanceResponse object into the model struct
func readProxyServerInstanceResponseDataSource(ctx context.Context, r *client.ProxyServerInstanceResponse, state *serverInstanceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("proxy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), false)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, false)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, false)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, false)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, false)
	state.LdapPort = internaltypes.Int64TypeOrNil(r.LdapPort)
	state.LdapsPort = internaltypes.Int64TypeOrNil(r.LdapsPort)
	state.HttpPort = internaltypes.Int64TypeOrNil(r.HttpPort)
	state.HttpsPort = internaltypes.Int64TypeOrNil(r.HttpsPort)
	state.ReplicationPort = internaltypes.Int64TypeOrNil(r.ReplicationPort)
	state.ReplicationServerID = internaltypes.Int64TypeOrNil(r.ReplicationServerID)
	state.ReplicationDomainServerID = internaltypes.GetInt64Set(r.ReplicationDomainServerID)
	state.JmxPort = internaltypes.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = internaltypes.Int64TypeOrNil(r.JmxsPort)
	state.PreferredSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), false)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
}

// Read a MetricsEngineServerInstanceResponse object into the model struct
func readMetricsEngineServerInstanceResponseDataSource(ctx context.Context, r *client.MetricsEngineServerInstanceResponse, state *serverInstanceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("metrics-engine")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), false)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, false)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, false)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, false)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, false)
	state.LdapPort = internaltypes.Int64TypeOrNil(r.LdapPort)
	state.LdapsPort = internaltypes.Int64TypeOrNil(r.LdapsPort)
	state.HttpPort = internaltypes.Int64TypeOrNil(r.HttpPort)
	state.HttpsPort = internaltypes.Int64TypeOrNil(r.HttpsPort)
	state.ReplicationPort = internaltypes.Int64TypeOrNil(r.ReplicationPort)
	state.ReplicationServerID = internaltypes.Int64TypeOrNil(r.ReplicationServerID)
	state.ReplicationDomainServerID = internaltypes.GetInt64Set(r.ReplicationDomainServerID)
	state.JmxPort = internaltypes.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = internaltypes.Int64TypeOrNil(r.JmxsPort)
	state.PreferredSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), false)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
}

// Read a AuthorizeServerInstanceResponse object into the model struct
func readAuthorizeServerInstanceResponseDataSource(ctx context.Context, r *client.AuthorizeServerInstanceResponse, state *serverInstanceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("authorize")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), false)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, false)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, false)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, false)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, false)
	state.LdapPort = internaltypes.Int64TypeOrNil(r.LdapPort)
	state.LdapsPort = internaltypes.Int64TypeOrNil(r.LdapsPort)
	state.HttpPort = internaltypes.Int64TypeOrNil(r.HttpPort)
	state.HttpsPort = internaltypes.Int64TypeOrNil(r.HttpsPort)
	state.ReplicationPort = internaltypes.Int64TypeOrNil(r.ReplicationPort)
	state.ReplicationServerID = internaltypes.Int64TypeOrNil(r.ReplicationServerID)
	state.ReplicationDomainServerID = internaltypes.GetInt64Set(r.ReplicationDomainServerID)
	state.JmxPort = internaltypes.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = internaltypes.Int64TypeOrNil(r.JmxsPort)
	state.PreferredSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), false)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
}

// Read a DirectoryServerInstanceResponse object into the model struct
func readDirectoryServerInstanceResponseDataSource(ctx context.Context, r *client.DirectoryServerInstanceResponse, state *serverInstanceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), false)
	state.ReplicationSetName = internaltypes.StringTypeOrNil(r.ReplicationSetName, false)
	state.LoadBalancingAlgorithmName = internaltypes.GetStringSet(r.LoadBalancingAlgorithmName)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, false)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, false)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, false)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, false)
	state.LdapPort = internaltypes.Int64TypeOrNil(r.LdapPort)
	state.LdapsPort = internaltypes.Int64TypeOrNil(r.LdapsPort)
	state.HttpPort = internaltypes.Int64TypeOrNil(r.HttpPort)
	state.HttpsPort = internaltypes.Int64TypeOrNil(r.HttpsPort)
	state.ReplicationPort = internaltypes.Int64TypeOrNil(r.ReplicationPort)
	state.ReplicationServerID = internaltypes.Int64TypeOrNil(r.ReplicationServerID)
	state.ReplicationDomainServerID = internaltypes.GetInt64Set(r.ReplicationDomainServerID)
	state.JmxPort = internaltypes.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = internaltypes.Int64TypeOrNil(r.JmxsPort)
	state.PreferredSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), false)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
}

// Read a SyncServerInstanceResponse object into the model struct
func readSyncServerInstanceResponseDataSource(ctx context.Context, r *client.SyncServerInstanceResponse, state *serverInstanceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("sync")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), false)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, false)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, false)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, false)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, false)
	state.LdapPort = internaltypes.Int64TypeOrNil(r.LdapPort)
	state.LdapsPort = internaltypes.Int64TypeOrNil(r.LdapsPort)
	state.HttpPort = internaltypes.Int64TypeOrNil(r.HttpPort)
	state.HttpsPort = internaltypes.Int64TypeOrNil(r.HttpsPort)
	state.ReplicationPort = internaltypes.Int64TypeOrNil(r.ReplicationPort)
	state.ReplicationServerID = internaltypes.Int64TypeOrNil(r.ReplicationServerID)
	state.ReplicationDomainServerID = internaltypes.GetInt64Set(r.ReplicationDomainServerID)
	state.JmxPort = internaltypes.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = internaltypes.Int64TypeOrNil(r.JmxsPort)
	state.PreferredSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), false)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
}

// Read resource information
func (r *serverInstanceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state serverInstanceDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ProxyServerInstanceResponse != nil {
		readProxyServerInstanceResponseDataSource(ctx, readResponse.ProxyServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MetricsEngineServerInstanceResponse != nil {
		readMetricsEngineServerInstanceResponseDataSource(ctx, readResponse.MetricsEngineServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AuthorizeServerInstanceResponse != nil {
		readAuthorizeServerInstanceResponseDataSource(ctx, readResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DirectoryServerInstanceResponse != nil {
		readDirectoryServerInstanceResponseDataSource(ctx, readResponse.DirectoryServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyncServerInstanceResponse != nil {
		readSyncServerInstanceResponseDataSource(ctx, readResponse.SyncServerInstanceResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
