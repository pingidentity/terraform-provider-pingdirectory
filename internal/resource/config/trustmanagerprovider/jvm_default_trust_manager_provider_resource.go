package trustmanagerprovider

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
	_ resource.Resource                = &jvmDefaultTrustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &jvmDefaultTrustManagerProviderResource{}
	_ resource.ResourceWithImportState = &jvmDefaultTrustManagerProviderResource{}
)

// Create a JVM Default Trust Manager Provider resource
func NewJVMDefaultTrustManagerProviderResource() resource.Resource {
	return &jvmDefaultTrustManagerProviderResource{}
}

// jvmDefaultTrustManagerProviderResource is the resource implementation.
type jvmDefaultTrustManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// jvmDefaultTrustManagerProviderResourceModel maps the resource schema data.
type jvmDefaultTrustManagerProviderResourceModel struct {
	// Id field required for acceptance testing framework
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	Enabled         types.Bool   `tfsdk:"enabled"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
}

// Metadata returns the resource type name.
func (r *jvmDefaultTrustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jvm_default_trust_manager_provider"
}

// GetSchema defines the schema for the resource.
func (r *jvmDefaultTrustManagerProviderResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	schema := tfsdk.Schema{
		Description: "Manages a JVM Default Trust Manager Provider.",
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "Name of the Trust Manager Provider.",
				Type:        types.StringType,
				Required:    true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"enabled": {
				Description: "Indicate whether the Trust Manager Provider is enabled for use.",
				Type:        types.BoolType,
				Required:    true,
			},
		},
	}
	config.AddCommonSchema(&schema)
	return schema, nil
}

// Configure adds the provider configured client to the resource.
func (r *jvmDefaultTrustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

// Create a new resource
func (r *jvmDefaultTrustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jvmDefaultTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddJvmDefaultTrustManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumjvmDefaultTrustManagerProviderSchemaUrn{client.ENUMJVMDEFAULTTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERJVM_DEFAULT},
		plan.Enabled.ValueBool())
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddJvmDefaultTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Trust Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := trustManagerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	readJvmDefaultTrustManagerProviderResponse(ctx, trustManagerResponse.JvmDefaultTrustManagerProviderResponse, &plan)

	// Populate Computed attribute values
	plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read a JvmDefaultTrustManagerProviderResponse object into the model struct
func readJvmDefaultTrustManagerProviderResponse(ctx context.Context, r *client.JvmDefaultTrustManagerProviderResponse, state *jvmDefaultTrustManagerProviderResourceModel) {
	// Placeholder Id value for acceptance test framework
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20)
}

// Read resource information
func (r *jvmDefaultTrustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state jvmDefaultTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Trust Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := trustManagerResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJvmDefaultTrustManagerProviderResponse(ctx, trustManagerResponse.JvmDefaultTrustManagerProviderResponse, &state)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create any update operations necessary to make the state match the plan
func createJvmDefaultTrustManagerProviderOperations(plan jvmDefaultTrustManagerProviderResourceModel, state jvmDefaultTrustManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation

	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Update a resource
func (r *jvmDefaultTrustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan jvmDefaultTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jvmDefaultTrustManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createJvmDefaultTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		trustManagerResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Trust Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := trustManagerResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJvmDefaultTrustManagerProviderResponse(ctx, trustManagerResponse.JvmDefaultTrustManagerProviderResponse, &plan)
		// Update computed values
		plan.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *jvmDefaultTrustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jvmDefaultTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProviderExecute(
		r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Trust Manager Provider", err, httpResp)
		return
	}
}

func (r *jvmDefaultTrustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to Name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
