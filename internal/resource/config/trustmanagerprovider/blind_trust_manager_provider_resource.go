package trustmanagerprovider

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &blindTrustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &blindTrustManagerProviderResource{}
	_ resource.ResourceWithImportState = &blindTrustManagerProviderResource{}
	_ resource.Resource                = &defaultBlindTrustManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultBlindTrustManagerProviderResource{}
	_ resource.ResourceWithImportState = &defaultBlindTrustManagerProviderResource{}
)

// Create a Blind Trust Manager Provider resource
func NewBlindTrustManagerProviderResource() resource.Resource {
	return &blindTrustManagerProviderResource{}
}

func NewDefaultBlindTrustManagerProviderResource() resource.Resource {
	return &defaultBlindTrustManagerProviderResource{}
}

// blindTrustManagerProviderResource is the resource implementation.
type blindTrustManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultBlindTrustManagerProviderResource is the resource implementation.
type defaultBlindTrustManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *blindTrustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_blind_trust_manager_provider"
}

func (r *defaultBlindTrustManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_blind_trust_manager_provider"
}

// Configure adds the provider configured client to the resource.
func (r *blindTrustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultBlindTrustManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type blindTrustManagerProviderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
	IncludeJVMDefaultIssuers types.Bool   `tfsdk:"include_jvm_default_issuers"`
}

// GetSchema defines the schema for the resource.
func (r *blindTrustManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	blindTrustManagerProviderSchema(ctx, req, resp, false)
}

func (r *defaultBlindTrustManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	blindTrustManagerProviderSchema(ctx, req, resp, true)
}

func blindTrustManagerProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Blind Trust Manager Provider.",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Indicate whether the Trust Manager Provider is enabled for use.",
				Required:    true,
			},
			"include_jvm_default_issuers": schema.BoolAttribute{
				Description: "Indicates whether certificates issued by an authority included in the JVM's set of default issuers should be automatically trusted, even if they would not otherwise be trusted by this provider.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalBlindTrustManagerProviderFields(ctx context.Context, addRequest *client.AddBlindTrustManagerProviderRequest, plan blindTrustManagerProviderResourceModel) {
	if internaltypes.IsDefined(plan.IncludeJVMDefaultIssuers) {
		boolVal := plan.IncludeJVMDefaultIssuers.ValueBool()
		addRequest.IncludeJVMDefaultIssuers = &boolVal
	}
}

// Read a BlindTrustManagerProviderResponse object into the model struct
func readBlindTrustManagerProviderResponse(ctx context.Context, r *client.BlindTrustManagerProviderResponse, state *blindTrustManagerProviderResourceModel, expectedValues *blindTrustManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IncludeJVMDefaultIssuers = internaltypes.BoolTypeOrNil(r.IncludeJVMDefaultIssuers)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createBlindTrustManagerProviderOperations(plan blindTrustManagerProviderResourceModel, state blindTrustManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeJVMDefaultIssuers, state.IncludeJVMDefaultIssuers, "include-jvm-default-issuers")
	return ops
}

// Create a new resource
func (r *blindTrustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan blindTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddBlindTrustManagerProviderRequest(plan.Id.ValueString(),
		[]client.EnumblindTrustManagerProviderSchemaUrn{client.ENUMBLINDTRUSTMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TRUST_MANAGER_PROVIDERBLIND},
		plan.Enabled.ValueBool())
	addOptionalBlindTrustManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TrustManagerProviderApi.AddTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddTrustManagerProviderRequest(
		client.AddBlindTrustManagerProviderRequestAsAddTrustManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.AddTrustManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Blind Trust Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state blindTrustManagerProviderResourceModel
	readBlindTrustManagerProviderResponse(ctx, addResponse.BlindTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultBlindTrustManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan blindTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Blind Trust Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state blindTrustManagerProviderResourceModel
	readBlindTrustManagerProviderResponse(ctx, readResponse.BlindTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createBlindTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Blind Trust Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readBlindTrustManagerProviderResponse(ctx, updateResponse.BlindTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *blindTrustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readBlindTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultBlindTrustManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readBlindTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readBlindTrustManagerProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state blindTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TrustManagerProviderApi.GetTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Blind Trust Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readBlindTrustManagerProviderResponse(ctx, readResponse.BlindTrustManagerProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *blindTrustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateBlindTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultBlindTrustManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateBlindTrustManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateBlindTrustManagerProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan blindTrustManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state blindTrustManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TrustManagerProviderApi.UpdateTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createBlindTrustManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TrustManagerProviderApi.UpdateTrustManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Blind Trust Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readBlindTrustManagerProviderResponse(ctx, updateResponse.BlindTrustManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultBlindTrustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *blindTrustManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state blindTrustManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProviderExecute(r.apiClient.TrustManagerProviderApi.DeleteTrustManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Blind Trust Manager Provider", err, httpResp)
		return
	}
}

func (r *blindTrustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importBlindTrustManagerProvider(ctx, req, resp)
}

func (r *defaultBlindTrustManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importBlindTrustManagerProvider(ctx, req, resp)
}

func importBlindTrustManagerProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
