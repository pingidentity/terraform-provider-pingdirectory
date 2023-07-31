package serverinstance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverInstanceResource{}
	_ resource.ResourceWithConfigure   = &serverInstanceResource{}
	_ resource.ResourceWithImportState = &serverInstanceResource{}
)

// Create a Server Instance resource
func NewServerInstanceResource() resource.Resource {
	return &serverInstanceResource{}
}

// serverInstanceResource is the resource implementation.
type serverInstanceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *serverInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_server_instance"
}

// Configure adds the provider configured client to the resource.
func (r *serverInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type serverInstanceResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *serverInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Server Instance.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Server Instance resource. Options are ['proxy', 'metrics-engine', 'authorize', 'directory', 'sync']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"proxy", "metrics-engine", "authorize", "directory", "sync"}...),
				},
			},
			"server_instance_type": schema.StringAttribute{
				Description: "Specifies the type of server installation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"replication_set_name": schema.StringAttribute{
				Description: "The name of the replication set assigned to this Directory Server. Restricted domains are only replicated within instances using the same replication set name.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancing_algorithm_name": schema.SetAttribute{
				Description: "The name of the configuration object for a load-balancing algorithm that should include this server.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server_instance_name": schema.StringAttribute{
				Description: "The name of this Server Instance. The instance name needs to be unique if this server will be part of a topology of servers that are connected to each other. Once set, it may not be changed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cluster_name": schema.StringAttribute{
				Description: "The name of the cluster to which this Server Instance belongs. Server instances within the same cluster will share the same cluster-wide configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_instance_location": schema.StringAttribute{
				Description: "Specifies the location for the Server Instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": schema.StringAttribute{
				Description: "The name of the host where this Server Instance is installed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_root": schema.StringAttribute{
				Description: "The file system path where this Server Instance is installed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_version": schema.StringAttribute{
				Description: "The version of the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inter_server_certificate": schema.StringAttribute{
				Description: "The public component of the certificate used by this instance to protect inter-server communication and to perform server-specific encryption. This will generally be managed by the server and should only be altered by administrators under explicit direction from Ping Identity support personnel.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for LDAP connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ldaps_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for LDAP secure connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"http_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for HTTP connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"https_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for HTTPS connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"replication_port": schema.Int64Attribute{
				Description: "The replication TCP port.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"replication_server_id": schema.Int64Attribute{
				Description: "Specifies a unique identifier for the replication server on this server instance.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"replication_domain_server_id": schema.SetAttribute{
				Description: "Specifies a unique identifier for the Directory Server within the replication domain.",
				Optional:    true,
				Computed:    true,
				ElementType: types.Int64Type,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"jmx_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for JMX connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"jmxs_port": schema.Int64Attribute{
				Description: "The TCP port on which this server is listening for JMX secure connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"preferred_security": schema.StringAttribute{
				Description: "Specifies the preferred mechanism to use for securing connections to the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_tls_enabled": schema.BoolAttribute{
				Description: "Indicates whether StartTLS is enabled on this server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"base_dn": schema.SetAttribute{
				Description: "The set of base DNs under the root DSE.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"member_of_server_group": schema.SetAttribute{
				Description: "The set of groups of which this server is a member.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *serverInstanceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model serverInstanceResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.JmxPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'jmx_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'jmx_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.StartTLSEnabled) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'start_tls_enabled' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'start_tls_enabled', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ReplicationSetName) && model.Type.ValueString() != "directory" {
		resp.Diagnostics.AddError("Attribute 'replication_set_name' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'replication_set_name', the 'type' attribute must be one of ['directory']")
	}
	if internaltypes.IsDefined(model.MemberOfServerGroup) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'member_of_server_group' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'member_of_server_group', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.HttpsPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'https_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'https_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.PreferredSecurity) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'preferred_security' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'preferred_security', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ReplicationPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'replication_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'replication_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.LoadBalancingAlgorithmName) && model.Type.ValueString() != "directory" {
		resp.Diagnostics.AddError("Attribute 'load_balancing_algorithm_name' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'load_balancing_algorithm_name', the 'type' attribute must be one of ['directory']")
	}
	if internaltypes.IsDefined(model.ServerRoot) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'server_root' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_root', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ReplicationDomainServerID) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'replication_domain_server_id' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'replication_domain_server_id', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ServerInstanceType) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'server_instance_type' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_instance_type', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.LdapsPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'ldaps_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ldaps_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.BaseDN) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'base_dn' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'base_dn', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ClusterName) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'cluster_name' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'cluster_name', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.Hostname) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'hostname' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'hostname', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.HttpPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'http_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'http_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ServerInstanceLocation) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'server_instance_location' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_instance_location', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ReplicationServerID) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'replication_server_id' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'replication_server_id', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ServerVersion) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'server_version' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_version', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.LdapPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'ldap_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ldap_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.ServerInstanceName) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'server_instance_name' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_instance_name', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.JmxsPort) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'jmxs_port' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'jmxs_port', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
	if internaltypes.IsDefined(model.InterServerCertificate) && model.Type.ValueString() != "proxy" && model.Type.ValueString() != "authorize" && model.Type.ValueString() != "directory" && model.Type.ValueString() != "sync" {
		resp.Diagnostics.AddError("Attribute 'inter_server_certificate' not supported by pingdirectory_server_instance resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'inter_server_certificate', the 'type' attribute must be one of ['proxy', 'authorize', 'directory', 'sync']")
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateServerInstanceUnknownValues(ctx context.Context, model *serverInstanceResourceModel) {
	if model.ReplicationDomainServerID.ElementType(ctx) == nil {
		model.ReplicationDomainServerID = types.SetNull(types.Int64Type)
	}
	if model.BaseDN.ElementType(ctx) == nil {
		model.BaseDN = types.SetNull(types.StringType)
	}
	if model.LoadBalancingAlgorithmName.ElementType(ctx) == nil {
		model.LoadBalancingAlgorithmName = types.SetNull(types.StringType)
	}
	if model.MemberOfServerGroup.ElementType(ctx) == nil {
		model.MemberOfServerGroup = types.SetNull(types.StringType)
	}
}

// Read a ProxyServerInstanceResponse object into the model struct
func readProxyServerInstanceResponse(ctx context.Context, r *client.ProxyServerInstanceResponse, state *serverInstanceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("proxy")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), true)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, true)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, true)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, true)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, true)
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
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), true)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateServerInstanceUnknownValues(ctx, state)
}

// Read a MetricsEngineServerInstanceResponse object into the model struct
func readMetricsEngineServerInstanceResponse(ctx context.Context, r *client.MetricsEngineServerInstanceResponse, state *serverInstanceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("metrics-engine")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), true)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, true)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, true)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, true)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, true)
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
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), true)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateServerInstanceUnknownValues(ctx, state)
}

// Read a AuthorizeServerInstanceResponse object into the model struct
func readAuthorizeServerInstanceResponse(ctx context.Context, r *client.AuthorizeServerInstanceResponse, state *serverInstanceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("authorize")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), true)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, true)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, true)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, true)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, true)
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
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), true)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateServerInstanceUnknownValues(ctx, state)
}

// Read a DirectoryServerInstanceResponse object into the model struct
func readDirectoryServerInstanceResponse(ctx context.Context, r *client.DirectoryServerInstanceResponse, state *serverInstanceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("directory")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), true)
	state.ReplicationSetName = internaltypes.StringTypeOrNil(r.ReplicationSetName, true)
	state.LoadBalancingAlgorithmName = internaltypes.GetStringSet(r.LoadBalancingAlgorithmName)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, true)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, true)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, true)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, true)
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
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), true)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateServerInstanceUnknownValues(ctx, state)
}

// Read a SyncServerInstanceResponse object into the model struct
func readSyncServerInstanceResponse(ctx context.Context, r *client.SyncServerInstanceResponse, state *serverInstanceResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("sync")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceType = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceServerInstanceTypeProp(r.ServerInstanceType), true)
	state.ServerInstanceName = types.StringValue(r.ServerInstanceName)
	state.ClusterName = types.StringValue(r.ClusterName)
	state.ServerInstanceLocation = internaltypes.StringTypeOrNil(r.ServerInstanceLocation, true)
	state.Hostname = internaltypes.StringTypeOrNil(r.Hostname, true)
	state.ServerRoot = internaltypes.StringTypeOrNil(r.ServerRoot, true)
	state.ServerVersion = types.StringValue(r.ServerVersion)
	state.InterServerCertificate = internaltypes.StringTypeOrNil(r.InterServerCertificate, true)
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
		client.StringPointerEnumserverInstancePreferredSecurityProp(r.PreferredSecurity), true)
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateServerInstanceUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createServerInstanceOperations(plan serverInstanceResourceModel, state serverInstanceResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ServerInstanceType, state.ServerInstanceType, "server-instance-type")
	operations.AddStringOperationIfNecessary(&ops, plan.ReplicationSetName, state.ReplicationSetName, "replication-set-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.LoadBalancingAlgorithmName, state.LoadBalancingAlgorithmName, "load-balancing-algorithm-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerInstanceName, state.ServerInstanceName, "server-instance-name")
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
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ReplicationDomainServerID, state.ReplicationDomainServerID, "replication-domain-server-id")
	operations.AddInt64OperationIfNecessary(&ops, plan.JmxPort, state.JmxPort, "jmx-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.JmxsPort, state.JmxsPort, "jmxs-port")
	operations.AddStringOperationIfNecessary(&ops, plan.PreferredSecurity, state.PreferredSecurity, "preferred-security")
	operations.AddBoolOperationIfNecessary(&ops, plan.StartTLSEnabled, state.StartTLSEnabled, "start-tls-enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BaseDN, state.BaseDN, "base-dn")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MemberOfServerGroup, state.MemberOfServerGroup, "member-of-server-group")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *serverInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan serverInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state serverInstanceResourceModel
	if plan.Type.ValueString() == "proxy" {
		readProxyServerInstanceResponse(ctx, readResponse.ProxyServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "metrics-engine" {
		readMetricsEngineServerInstanceResponse(ctx, readResponse.MetricsEngineServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "authorize" {
		readAuthorizeServerInstanceResponse(ctx, readResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "directory" {
		readDirectoryServerInstanceResponse(ctx, readResponse.DirectoryServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "sync" {
		readSyncServerInstanceResponse(ctx, readResponse.SyncServerInstanceResponse, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createServerInstanceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "proxy" {
			readProxyServerInstanceResponse(ctx, updateResponse.ProxyServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "metrics-engine" {
			readMetricsEngineServerInstanceResponse(ctx, updateResponse.MetricsEngineServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "authorize" {
			readAuthorizeServerInstanceResponse(ctx, updateResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "directory" {
			readDirectoryServerInstanceResponse(ctx, updateResponse.DirectoryServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "sync" {
			readSyncServerInstanceResponse(ctx, updateResponse.SyncServerInstanceResponse, &state, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *serverInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state serverInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
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
		readProxyServerInstanceResponse(ctx, readResponse.ProxyServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MetricsEngineServerInstanceResponse != nil {
		readMetricsEngineServerInstanceResponse(ctx, readResponse.MetricsEngineServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AuthorizeServerInstanceResponse != nil {
		readAuthorizeServerInstanceResponse(ctx, readResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DirectoryServerInstanceResponse != nil {
		readDirectoryServerInstanceResponse(ctx, readResponse.DirectoryServerInstanceResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyncServerInstanceResponse != nil {
		readSyncServerInstanceResponse(ctx, readResponse.SyncServerInstanceResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *serverInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state serverInstanceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createServerInstanceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "proxy" {
			readProxyServerInstanceResponse(ctx, updateResponse.ProxyServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "metrics-engine" {
			readMetricsEngineServerInstanceResponse(ctx, updateResponse.MetricsEngineServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "authorize" {
			readAuthorizeServerInstanceResponse(ctx, updateResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "directory" {
			readDirectoryServerInstanceResponse(ctx, updateResponse.DirectoryServerInstanceResponse, &state, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "sync" {
			readSyncServerInstanceResponse(ctx, updateResponse.SyncServerInstanceResponse, &state, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *serverInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *serverInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
