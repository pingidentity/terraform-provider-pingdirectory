package serverinstance

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &authorizeServerInstanceResource{}
	_ resource.ResourceWithConfigure   = &authorizeServerInstanceResource{}
	_ resource.ResourceWithImportState = &authorizeServerInstanceResource{}
)

// Create a Authorize Server Instance resource
func NewAuthorizeServerInstanceResource() resource.Resource {
	return &authorizeServerInstanceResource{}
}

// authorizeServerInstanceResource is the resource implementation.
type authorizeServerInstanceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *authorizeServerInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_authorize_server_instance"
}

// Configure adds the provider configured client to the resource.
func (r *authorizeServerInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type authorizeServerInstanceResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	ServerInstanceType        types.String `tfsdk:"server_instance_type"`
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
}

// GetSchema defines the schema for the resource.
func (r *authorizeServerInstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Authorize Server Instance.",
		Attributes: map[string]schema.Attribute{
			"server_instance_type": schema.StringAttribute{
				Description: "Specifies the type of server installation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.Int64Type,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"member_of_server_group": schema.SetAttribute{
				Description: "The set of groups of which this server is a member.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a AuthorizeServerInstanceResponse object into the model struct
func readAuthorizeServerInstanceResponse(ctx context.Context, r *client.AuthorizeServerInstanceResponse, state *authorizeServerInstanceResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createAuthorizeServerInstanceOperations(plan authorizeServerInstanceResourceModel, state authorizeServerInstanceResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ServerInstanceType, state.ServerInstanceType, "server-instance-type")
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
func (r *authorizeServerInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan authorizeServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Authorize Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state authorizeServerInstanceResourceModel
	readAuthorizeServerInstanceResponse(ctx, readResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAuthorizeServerInstanceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Authorize Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAuthorizeServerInstanceResponse(ctx, updateResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
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
func (r *authorizeServerInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state authorizeServerInstanceResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceApi.GetServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Authorize Server Instance", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAuthorizeServerInstanceResponse(ctx, readResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *authorizeServerInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan authorizeServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state authorizeServerInstanceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAuthorizeServerInstanceOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ServerInstanceApi.UpdateServerInstanceExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Authorize Server Instance", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAuthorizeServerInstanceResponse(ctx, updateResponse.AuthorizeServerInstanceResponse, &state, &resp.Diagnostics)
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
func (r *authorizeServerInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *authorizeServerInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
