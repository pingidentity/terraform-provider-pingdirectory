package serverinstance

import (
	"context"
	"terraform-provider-pingdirectory/internal/operations"
	"terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "terraform-provider-pingdirectory/internal/types"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	_ resource.Resource                = &syncServerInstanceResource{}
	_ resource.ResourceWithConfigure   = &syncServerInstanceResource{}
	_ resource.ResourceWithImportState = &syncServerInstanceResource{}
)

// Create a Sync Server Instance resource
func NewSyncServerInstanceResource() resource.Resource {
	return &syncServerInstanceResource{}
}

// syncServerInstanceResource is the resource implementation.
type syncServerInstanceResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *syncServerInstanceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sync_server_instance"
}

// GetSchema defines the schema for the resource.
func (r *syncServerInstanceResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return GetCommonServerInstanceSchema("Manages a Sync Server Instance.")
}

// Configure adds the provider configured client to the resource.
func (r *syncServerInstanceResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
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
func (r *syncServerInstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan CommonServerInstanceResourceModel
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
	var state CommonServerInstanceResourceModel
	readSyncServerInstanceResponse(ctx, getResp.SyncServerInstanceResponse, &state)

	// Determine what changes need to be made to match the plan
	updateInstanceRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(config.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString())
	ops := CreateCommonServerInstanceOperations(plan, state)

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
		readSyncServerInstanceResponse(ctx, instanceResp.SyncServerInstanceResponse, &plan)
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

// Read a SyncServerInstanceResponse object into the model struct.
// Use empty string for nils since everything is marked as computed.
func readSyncServerInstanceResponse(ctx context.Context, r *client.SyncServerInstanceResponse, state *CommonServerInstanceResourceModel) {
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
	if r.PreferredSecurity != nil {
		state.PreferredSecurity = types.StringValue(string(*r.PreferredSecurity))
	} else {
		state.PreferredSecurity = types.StringValue("")
	}
	state.StartTLSEnabled = internaltypes.BoolTypeOrNil(r.StartTLSEnabled)
	state.BaseDN = internaltypes.GetStringSet(r.BaseDN)
	state.MemberOfServerGroup = internaltypes.GetStringSet(r.MemberOfServerGroup)
	// Report any notifications from the Config API
	if r.Urnpingidentityschemasconfigurationmessages20 != nil {
		state.Notifications = internaltypes.GetStringSet(r.Urnpingidentityschemasconfigurationmessages20.Notifications)
		config.LogNotifications(ctx, r.Urnpingidentityschemasconfigurationmessages20)
	} else {
		state.Notifications, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Read resource information
func (r *syncServerInstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state CommonServerInstanceResourceModel
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
	readSyncServerInstanceResponse(ctx, serverInstanceResponse.SyncServerInstanceResponse, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *syncServerInstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan CommonServerInstanceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state CommonServerInstanceResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ServerInstanceApi.UpdateServerInstance(config.BasicAuthContext(ctx, r.providerConfig), plan.ServerInstanceName.ValueString())

	// Determine what update operations are necessary
	ops := CreateCommonServerInstanceOperations(plan, state)
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
		readSyncServerInstanceResponse(ctx, serverInstanceResponse.SyncServerInstanceResponse, &state)
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
func (r *syncServerInstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *syncServerInstanceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
