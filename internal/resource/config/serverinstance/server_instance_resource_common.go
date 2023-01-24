package serverinstance

import (
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// commonServerInstanceResource maps the resource schema data common to server instance config objects.
// The structs and functions defined in this file are used by the proxy, authorize, and sync server instance
// resource implementations because they use the exact same model. If they change, they won't be able to share
// these definitions.
type CommonServerInstanceResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	ServerInstanceName        types.String `tfsdk:"server_instance_name"`
	ClusterName               types.String `tfsdk:"cluster_name"`
	ServerInstanceLocation    types.String `tfsdk:"server_instance_location"`
	Hostname                  types.String `tfsdk:"hostname"`
	ServerRoot                types.String `tfsdk:"server_root"`
	ServerVersion             types.String `tfsdk:"server_version"`
	InterServerCertificate    types.String `tfsdk:"inter_server_certificate"`
	LdapPort                  types.Int64  `tfsdk:"ldap_port"`
	LdapsPort                 types.Int64  `tfsdk:"ldaps_port"`
	HttpPort                  types.Int64  `tfsdk:"http_port"`
	HttpsPort                 types.Int64  `tfsdk:"https_port"`
	ReplicationPort           types.Int64  `tfsdk:"replication_port"`
	ReplicationServerID       types.Int64  `tfsdk:"replication_server_id"`
	ReplicationDomainServerID types.Set    `tfsdk:"replication_domain_server_id"`
	JmxPort                   types.Int64  `tfsdk:"jmx_port"`
	JmxsPort                  types.Int64  `tfsdk:"jmxs_port"`
	PreferredSecurity         types.String `tfsdk:"preferred_security"`
	StartTLSEnabled           types.Bool   `tfsdk:"start_tls_enabled"`
	BaseDN                    types.Set    `tfsdk:"base_dn"`
	MemberOfServerGroup       types.Set    `tfsdk:"member_of_server_group"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
}

// GetCommonServerInstanceSchema defines the common schema for server instance resources.
func GetCommonServerInstanceSchema(description string) schema.Schema {
	schema := schema.Schema{
		Description: description,
		Attributes: map[string]schema.Attribute{
			// All are considered computed, since we are importing the existing server
			// instance from a server, rather than "creating" a server instance
			// like a typical Terraform resource.
			"server_instance_name": schema.StringAttribute{
				Description: "The name of this Server Instance. The instance name needs to be unique if this server will be part of a topology of servers that are connected to each other. Once set, it may not be changed.",
				Required:    true,
			},
			"cluster_name": schema.StringAttribute{
				Description: "The name of the cluster to which this Server Instance belongs. Server instances within the same cluster will share the same cluster-wide configuration.",
				Optional:    true,
				Computed:    true,
			},
			"server_instance_location": schema.StringAttribute{
				Description: "Specifies the location for the Server Instance.",
				Optional:    true,
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "The name of the host where this Server Instance is installed.",
				Optional:    true,
				Computed:    true,
			},
			"server_root": schema.StringAttribute{
				Description: "The file system path where this Server Instance is installed.",
				Optional:    true,
				Computed:    true,
			},
			"server_version": schema.StringAttribute{
				Description: "The version of the server.",
				Optional:    true,
				Computed:    true,
			},
			"inter_server_certificate": schema.StringAttribute{
				Description: "The public component of the certificate used by this instance to protect inter-server communication and to perform server-specific encryption. This will generally be managed by the server and should only be altered by administrators under explicit direction from Ping Identity support personnel.",
				Optional:    true,
				Computed:    true,
			},
			"ldap_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for LDAP connections.",
				Optional:    true,
				Computed:    true,
			},
			"ldaps_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for LDAP secure connections.",
				Optional:    true,
				Computed:    true,
			},
			"http_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for HTTP connections.",
				Optional:    true,
				Computed:    true,
			},
			"https_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for HTTPS connections.",
				Optional:    true,
				Computed:    true,
			},
			"replication_port": schema.Int64Attribute{
				Description: "The replication TCP port.",
				Optional:    true,
				Computed:    true,
			},
			"replication_server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the replication server on this server instance.",
				Optional:    true,
				Computed:    true,
			},
			"replication_domain_server_id": schema.SetAttribute{
				Description: "Specifies a unique identifier for the Directory Server within the replication domain.",
				ElementType: types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"jmx_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for JMX connections.",
				Optional:    true,
				Computed:    true,
			},
			"jmxs_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for JMX secure connections.",
				Optional:    true,
				Computed:    true,
			},
			"preferred_security": schema.StringAttribute{
				Description: "Specifies the preferred mechanism to use for securing connections to the server.",
				Optional:    true,
				Computed:    true,
			},
			"start_tls_enabled": schema.BoolAttribute{
				Description: "Indicates whether StartTLS is enabled on this server.",
				Optional:    true,
				Computed:    true,
			},
			"base_dn": schema.SetAttribute{
				Description: "The set of base DNs under the root DSE.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"member_of_server_group": schema.SetAttribute{
				Description: "The set of groups of which this server is a member.",
				ElementType: types.StringType,
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	return schema
}

// Create any update operations necessary to make the state match the plan
func CreateCommonServerInstanceOperations(plan CommonServerInstanceResourceModel, state CommonServerInstanceResourceModel) []client.Operation {
	var ops []client.Operation

	operations.AddStringOperationIfNecessary(&ops, plan.ClusterName, state.ClusterName, "cluster-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerInstanceLocation, state.ServerInstanceLocation, "server-instance-location")
	operations.AddStringOperationIfNecessary(&ops, plan.Hostname, state.Hostname, "hostname")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerRoot, state.ServerRoot, "server-root")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerVersion, state.ServerVersion, "server-version")
	operations.AddStringOperationIfNecessary(&ops, plan.InterServerCertificate, state.InterServerCertificate, "inter-server-certificate")
	operations.AddInt64OperationIfNecessary(&ops, plan.LdapPort, state.LdapPort, "ldap-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.LdapsPort, state.LdapsPort, "ldaps-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.HttpPort, state.HttpPort, "http-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.HttpsPort, state.HttpsPort, "https-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationPort, state.ReplicationPort, "replication-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.ReplicationServerID, state.ReplicationServerID, "replication-server-id")
	operations.AddInt64SetOperationsIfNecessary(&ops, plan.ReplicationDomainServerID, state.ReplicationDomainServerID, "replication-domain-server-id")
	operations.AddInt64OperationIfNecessary(&ops, plan.JmxPort, state.JmxPort, "jmx-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.JmxsPort, state.JmxsPort, "jmxs-port")
	operations.AddStringOperationIfNecessary(&ops, plan.PreferredSecurity, state.PreferredSecurity, "preferred-security")
	operations.AddBoolOperationIfNecessary(&ops, plan.StartTLSEnabled, state.StartTLSEnabled, "start-tls-enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MemberOfServerGroup, state.MemberOfServerGroup, "member-of-server-group")
	operations.AddStringOperationIfNecessary(&ops, plan.LastUpdated, state.LastUpdated, "last-updated")
	return ops
}
