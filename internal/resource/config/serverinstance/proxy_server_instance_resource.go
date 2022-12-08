package serverinstance

import (
	"context"
	"terraform-provider-pingdirectory/internal/operations"
	"terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "terraform-provider-pingdirectory/internal/types"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdata-config-api-go-client"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &proxyServerInstanceResource{}
	_ resource.ResourceWithConfigure   = &proxyServerInstanceResource{}
	_ resource.ResourceWithImportState = &proxyServerInstanceResource{}
)

// Create a Proxy Server Instance resource
func NewProxyServerInstanceResource() resource.Resource {
	return &proxyServerInstanceResource{}
}

// proxyServerInstanceResource is the resource implementation.
type proxyServerInstanceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// proxyServerInstanceResourceModel maps the resource schema data.
type proxyServerInstanceResourceModel struct {
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
}

// Metadata returns the resource type name.
func (r *proxyServerInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_proxy_server_instance"
}

// GetSchema defines the schema for the resource.
func (r *proxyServerInstanceResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a Proxy Server Instance.",
		Attributes: map[string]tfsdk.Attribute{
			// All are considered computed, since we are importing the existing server
			// instance from a server, rather than "creating" a server instance
			// like a typical Terraform resource.
			"server_instance_name": {
				Description: "The name of this Server Instance. The instance name needs to be unique if this server will be part of a topology of servers that are connected to each other. Once set, it may not be changed.",
				Type:        types.StringType,
				Required:    true,
			},
			"cluster_name": {
				Description: "The name of the cluster to which this Server Instance belongs. Server instances within the same cluster will share the same cluster-wide configuration.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"server_instance_location": {
				Description: "Specifies the location for the Server Instance.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"hostname": {
				Description: "The name of the host where this Server Instance is installed.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"server_root": {
				Description: "The file system path where this Server Instance is installed.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"server_version": {
				Description: "The version of the server.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"inter_server_certificate": {
				Description: "The public component of the certificate used by this instance to protect inter-server communication and to perform server-specific encryption. This will generally be managed by the server and should only be altered by administrators under explicit direction from Ping Identity support personnel.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"ldap_port": {
				Description: "The TCP port on which this server is listening for LDAP connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"ldaps_port": {
				Description: "The TCP port on which this server is listening for LDAP secure connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"http_port": {
				Description: "The TCP port on which this server is listening for HTTP connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"https_port": {
				Description: "The TCP port on which this server is listening for HTTPS connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_port": {
				Description: "The replication TCP port.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_server_id": {
				Description: "Specifies a unique identifier for the replication server on this server instance.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"replication_domain_server_id": {
				Description: "Specifies a unique identifier for the Directory Server within the replication domain.",
				Type: types.SetType{
					ElemType: types.Int64Type,
				},
				Optional: true,
				Computed: true,
			},
			"jmx_port": {
				Description: "The TCP port on which this server is listening for JMX connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"jmxs_port": {
				Description: "The TCP port on which this server is listening for JMX secure connections.",
				Type:        types.Int64Type,
				Optional:    true,
				Computed:    true,
			},
			"preferred_security": {
				Description: "Specifies the preferred mechanism to use for securing connections to the server.",
				Type:        types.StringType,
				Optional:    true,
				Computed:    true,
			},
			"start_tls_enabled": {
				Description: "Indicates whether StartTLS is enabled on this server.",
				Type:        types.BoolType,
				Optional:    true,
				Computed:    true,
			},
			"base_dn": {
				Description: "The set of base DNs under the root DSE.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"member_of_server_group": {
				Description: "The set of groups of which this server is a member.",
				Type: types.SetType{
					ElemType: types.StringType,
				},
				Optional: true,
				Computed: true,
			},
			"last_updated": {
				Description: "Timestamp of the last Terraform update of the Server Instance.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

// Configure adds the provider configured client to the resource.
func (r *proxyServerInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Create a new resource
// For server instances, create doesn't actually "create" anything - it "adopts" the existing
// server instance into management by terraform. This method reads the existing server instance
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *proxyServerInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan proxyServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	getResp, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(config.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := getResp.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read existing config
	var state proxyServerInstanceResourceModel
	readProxyServerInstanceResponse(getResp.ProxyServerInstanceResponse, &state)

	// Determine what changes need to be made to match the plan
	updateInstanceRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(config.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString())
	ops := createProxyServerInstanceOperations(plan, state)

	if len(ops) > 0 {
		updateInstanceRequest = updateInstanceRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)
		instanceResp, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateInstanceRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := instanceResp.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readProxyServerInstanceResponse(instanceResp.ProxyServerInstanceResponse, &plan)
		// Populate Computed attribute values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		// Just put the initial read into the plan
		plan = state
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a ProxyServerInstanceResponse object into the model struct.
// Use empty string for nils since everything is marked as computed.
func readProxyServerInstanceResponse(r *client.ProxyServerInstanceResponse, state *proxyServerInstanceResourceModel) {
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
	/*
		if r.ReplicationDomainServerID != nil {
			state.ReplicationDomainServerID = utils.GetInt64Set(*r.ReplicationDomainServerID)
		} else {
			state.ReplicationDomainServerID, _ = types.SetValue(types.Int64Type, []attr.Value{})
		}*/
	state.JmxPort = internaltypes.Int64TypeOrNil(r.JmxPort)
	state.JmxsPort = internaltypes.Int64TypeOrNil(r.JmxsPort)
	if r.PreferredSecurity != nil {
		state.PreferredSecurity = types.StringValue(string(*r.PreferredSecurity))
	} else {
		state.PreferredSecurity = types.StringValue("")
	}
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
}

// Read resource information
func (r *proxyServerInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state proxyServerInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverInstanceResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(config.BasicAuthContext(ctx, r.providerConfig), state.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := serverInstanceResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readProxyServerInstanceResponse(serverInstanceResponse.ProxyServerInstanceResponse, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create any update operations necessary to make the state match the plan
func createProxyServerInstanceOperations(plan proxyServerInstanceResourceModel, state proxyServerInstanceResourceModel) []client.Operation {
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

// Update a resource
func (r *proxyServerInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan proxyServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state proxyServerInstanceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(config.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString())

	// Determine what update operations are necessary
	ops := createProxyServerInstanceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		serverInstanceResponse, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := serverInstanceResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readProxyServerInstanceResponse(serverInstanceResponse.ProxyServerInstanceResponse, &state)
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
// Terraform can't actually delete server instances, so this method does nothing.
// Terraform will just "forget" about the server instance config, and it can be managed elsewhere.
func (r *proxyServerInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *proxyServerInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
