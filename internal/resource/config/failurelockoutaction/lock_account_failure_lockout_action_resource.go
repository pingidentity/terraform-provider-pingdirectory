package failurelockoutaction

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &lockAccountFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &lockAccountFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &lockAccountFailureLockoutActionResource{}
	_ resource.Resource                = &defaultLockAccountFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &defaultLockAccountFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &defaultLockAccountFailureLockoutActionResource{}
)

// Create a Lock Account Failure Lockout Action resource
func NewLockAccountFailureLockoutActionResource() resource.Resource {
	return &lockAccountFailureLockoutActionResource{}
}

func NewDefaultLockAccountFailureLockoutActionResource() resource.Resource {
	return &defaultLockAccountFailureLockoutActionResource{}
}

// lockAccountFailureLockoutActionResource is the resource implementation.
type lockAccountFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLockAccountFailureLockoutActionResource is the resource implementation.
type defaultLockAccountFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *lockAccountFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_lock_account_failure_lockout_action"
}

func (r *defaultLockAccountFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_lock_account_failure_lockout_action"
}

// Configure adds the provider configured client to the resource.
func (r *lockAccountFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLockAccountFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type lockAccountFailureLockoutActionResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *lockAccountFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	lockAccountFailureLockoutActionSchema(ctx, req, resp, false)
}

func (r *defaultLockAccountFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	lockAccountFailureLockoutActionSchema(ctx, req, resp, true)
}

func lockAccountFailureLockoutActionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Lock Account Failure Lockout Action.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this Failure Lockout Action",
				Optional:    true,
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
func addOptionalLockAccountFailureLockoutActionFields(ctx context.Context, addRequest *client.AddLockAccountFailureLockoutActionRequest, plan lockAccountFailureLockoutActionResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a LockAccountFailureLockoutActionResponse object into the model struct
func readLockAccountFailureLockoutActionResponse(ctx context.Context, r *client.LockAccountFailureLockoutActionResponse, state *lockAccountFailureLockoutActionResourceModel, expectedValues *lockAccountFailureLockoutActionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createLockAccountFailureLockoutActionOperations(plan lockAccountFailureLockoutActionResourceModel, state lockAccountFailureLockoutActionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *lockAccountFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan lockAccountFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLockAccountFailureLockoutActionRequest(plan.Id.ValueString(),
		[]client.EnumlockAccountFailureLockoutActionSchemaUrn{client.ENUMLOCKACCOUNTFAILURELOCKOUTACTIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0FAILURE_LOCKOUT_ACTIONLOCK_ACCOUNT})
	addOptionalLockAccountFailureLockoutActionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.FailureLockoutActionApi.AddFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFailureLockoutActionRequest(
		client.AddLockAccountFailureLockoutActionRequestAsAddFailureLockoutActionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.AddFailureLockoutActionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Lock Account Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state lockAccountFailureLockoutActionResourceModel
	readLockAccountFailureLockoutActionResponse(ctx, addResponse.LockAccountFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLockAccountFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan lockAccountFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Lock Account Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state lockAccountFailureLockoutActionResourceModel
	readLockAccountFailureLockoutActionResponse(ctx, readResponse.LockAccountFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLockAccountFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Lock Account Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLockAccountFailureLockoutActionResponse(ctx, updateResponse.LockAccountFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *lockAccountFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLockAccountFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLockAccountFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLockAccountFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLockAccountFailureLockoutAction(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state lockAccountFailureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Lock Account Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLockAccountFailureLockoutActionResponse(ctx, readResponse.LockAccountFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *lockAccountFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLockAccountFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLockAccountFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLockAccountFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLockAccountFailureLockoutAction(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan lockAccountFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state lockAccountFailureLockoutActionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLockAccountFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Lock Account Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLockAccountFailureLockoutActionResponse(ctx, updateResponse.LockAccountFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLockAccountFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *lockAccountFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state lockAccountFailureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutActionExecute(r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Lock Account Failure Lockout Action", err, httpResp)
		return
	}
}

func (r *lockAccountFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLockAccountFailureLockoutAction(ctx, req, resp)
}

func (r *defaultLockAccountFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLockAccountFailureLockoutAction(ctx, req, resp)
}

func importLockAccountFailureLockoutAction(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
