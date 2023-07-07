package serverinstancelistener

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &serverInstanceListenerResource{}
	_ resource.ResourceWithConfigure   = &serverInstanceListenerResource{}
	_ resource.ResourceWithImportState = &serverInstanceListenerResource{}
)

// Create a Server Instance Listener resource
func NewServerInstanceListenerResource() resource.Resource {
	return &serverInstanceListenerResource{}
}

// serverInstanceListenerResource is the resource implementation.
type serverInstanceListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *serverInstanceListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_server_instance_listener"
}

// Configure adds the provider configured client to the resource.
func (r *serverInstanceListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type serverInstanceListenerResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	Notifications       types.Set    `tfsdk:"notifications"`
	RequiredActions     types.Set    `tfsdk:"required_actions"`
	Type                types.String `tfsdk:"type"`
	ServerInstanceName  types.String `tfsdk:"server_instance_name"`
	ListenAddress       types.String `tfsdk:"listen_address"`
	ServerHTTPPort      types.Int64  `tfsdk:"server_http_port"`
	ServerLDAPPort      types.Int64  `tfsdk:"server_ldap_port"`
	ConnectionSecurity  types.String `tfsdk:"connection_security"`
	ListenerCertificate types.String `tfsdk:"listener_certificate"`
	Purpose             types.Set    `tfsdk:"purpose"`
}

// GetSchema defines the schema for the resource.
func (r *serverInstanceListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Server Instance Listener.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Server Instance Listener resource. Options are ['ldap', 'http']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ldap", "http"}...),
				},
			},
			"server_instance_name": schema.StringAttribute{
				Description: "Name of the parent Server Instance",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"listen_address": schema.StringAttribute{
				Description: "If the server is listening on a particular address different from the hostname, then this property may be used to specify the address on which to listen for connections from HTTP clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_http_port": schema.Int64Attribute{
				Description: "The TCP port number on which the HTTP server is listening.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_ldap_port": schema.Int64Attribute{
				Description: "The TCP port number on which the LDAP server is listening.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"connection_security": schema.StringAttribute{
				Description: "Specifies the mechanism to use for securing connections to the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listener_certificate": schema.StringAttribute{
				Description: "The public component of the certificate that the listener is expected to present to clients. When establishing a connection to this server, only the certificate(s) listed here will be trusted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"purpose": schema.SetAttribute{
				Description: "Identifies the purpose of this Server Instance Listener.",
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
func (r *serverInstanceListenerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model serverInstanceListenerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.ServerHTTPPort) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'server_http_port' not supported by pingdirectory_server_instance_listener resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_http_port', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.ListenAddress) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'listen_address' not supported by pingdirectory_server_instance_listener resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'listen_address', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.ServerLDAPPort) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'server_ldap_port' not supported by pingdirectory_server_instance_listener resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_ldap_port', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.ListenerCertificate) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'listener_certificate' not supported by pingdirectory_server_instance_listener resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'listener_certificate', the 'type' attribute must be one of ['ldap']")
	}
}

// Read a LdapServerInstanceListenerResponse object into the model struct
func readLdapServerInstanceListenerResponse(ctx context.Context, r *client.LdapServerInstanceListenerResponse, state *serverInstanceListenerResourceModel, expectedValues *serverInstanceListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceName = expectedValues.ServerInstanceName
	state.ServerLDAPPort = internaltypes.Int64TypeOrNil(r.ServerLDAPPort)
	state.ConnectionSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceListenerLdapConnectionSecurityProp(r.ConnectionSecurity), true)
	state.ListenerCertificate = internaltypes.StringTypeOrNil(r.ListenerCertificate, true)
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnumserverInstanceListenerPurposeProp(r.Purpose))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a HttpServerInstanceListenerResponse object into the model struct
func readHttpServerInstanceListenerResponse(ctx context.Context, r *client.HttpServerInstanceListenerResponse, state *serverInstanceListenerResourceModel, expectedValues *serverInstanceListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http")
	state.Id = types.StringValue(r.Id)
	state.ServerInstanceName = expectedValues.ServerInstanceName
	state.ListenAddress = internaltypes.StringTypeOrNil(r.ListenAddress, true)
	state.ServerHTTPPort = internaltypes.Int64TypeOrNil(r.ServerHTTPPort)
	state.ConnectionSecurity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumserverInstanceListenerHttpConnectionSecurityProp(r.ConnectionSecurity), true)
	state.Purpose = internaltypes.GetStringSet(
		client.StringSliceEnumserverInstanceListenerPurposeProp(r.Purpose))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createServerInstanceListenerOperations(plan serverInstanceListenerResourceModel, state serverInstanceListenerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ListenAddress, state.ListenAddress, "listen-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerHTTPPort, state.ServerHTTPPort, "server-http-port")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerLDAPPort, state.ServerLDAPPort, "server-ldap-port")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionSecurity, state.ConnectionSecurity, "connection-security")
	operations.AddStringOperationIfNecessary(&ops, plan.ListenerCertificate, state.ListenerCertificate, "listener-certificate")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.Purpose, state.Purpose, "purpose")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *serverInstanceListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan serverInstanceListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceListenerApi.GetServerInstanceListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Server Instance Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state serverInstanceListenerResourceModel
	if plan.Type.ValueString() == "ldap" {
		readLdapServerInstanceListenerResponse(ctx, readResponse.LdapServerInstanceListenerResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "http" {
		readHttpServerInstanceListenerResponse(ctx, readResponse.HttpServerInstanceListenerResponse, &state, &plan, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ServerInstanceListenerApi.UpdateServerInstanceListener(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceName.ValueString())
	ops := createServerInstanceListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ServerInstanceListenerApi.UpdateServerInstanceListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "ldap" {
			readLdapServerInstanceListenerResponse(ctx, updateResponse.LdapServerInstanceListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "http" {
			readHttpServerInstanceListenerResponse(ctx, updateResponse.HttpServerInstanceListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *serverInstanceListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state serverInstanceListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ServerInstanceListenerApi.GetServerInstanceListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.ServerInstanceName.ValueString()).Execute()
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
		readLdapServerInstanceListenerResponse(ctx, readResponse.LdapServerInstanceListenerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HttpServerInstanceListenerResponse != nil {
		readHttpServerInstanceListenerResponse(ctx, readResponse.HttpServerInstanceListenerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *serverInstanceListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan serverInstanceListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state serverInstanceListenerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ServerInstanceListenerApi.UpdateServerInstanceListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.ServerInstanceName.ValueString())

	// Determine what update operations are necessary
	ops := createServerInstanceListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ServerInstanceListenerApi.UpdateServerInstanceListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Server Instance Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "ldap" {
			readLdapServerInstanceListenerResponse(ctx, updateResponse.LdapServerInstanceListenerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "http" {
			readHttpServerInstanceListenerResponse(ctx, updateResponse.HttpServerInstanceListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *serverInstanceListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *serverInstanceListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [server-instance-name]/[server-instance-listener-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("server_instance_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}
