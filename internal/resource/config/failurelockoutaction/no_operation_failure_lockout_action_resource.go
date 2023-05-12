package failurelockoutaction

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &noOperationFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &noOperationFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &noOperationFailureLockoutActionResource{}
	_ resource.Resource                = &defaultNoOperationFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &defaultNoOperationFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &defaultNoOperationFailureLockoutActionResource{}
)

// Create a No Operation Failure Lockout Action resource
func NewNoOperationFailureLockoutActionResource() resource.Resource {
	return &noOperationFailureLockoutActionResource{}
}

func NewDefaultNoOperationFailureLockoutActionResource() resource.Resource {
	return &defaultNoOperationFailureLockoutActionResource{}
}

// noOperationFailureLockoutActionResource is the resource implementation.
type noOperationFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultNoOperationFailureLockoutActionResource is the resource implementation.
type defaultNoOperationFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *noOperationFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_no_operation_failure_lockout_action"
}

func (r *defaultNoOperationFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_no_operation_failure_lockout_action"
}

// Configure adds the provider configured client to the resource.
func (r *noOperationFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultNoOperationFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type noOperationFailureLockoutActionResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	GenerateAccountStatusNotification types.Bool   `tfsdk:"generate_account_status_notification"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *noOperationFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	noOperationFailureLockoutActionSchema(ctx, req, resp, false)
}

func (r *defaultNoOperationFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	noOperationFailureLockoutActionSchema(ctx, req, resp, true)
}

func noOperationFailureLockoutActionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a No Operation Failure Lockout Action.",
		Attributes: map[string]schema.Attribute{
			"generate_account_status_notification": schema.BoolAttribute{
				Description: "Indicates whether to generate an account status notification for cases in which this failure lockout action is invoked for a bind attempt with too many outstanding authentication failures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
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
func addOptionalNoOperationFailureLockoutActionFields(ctx context.Context, addRequest *client.AddNoOperationFailureLockoutActionRequest, plan noOperationFailureLockoutActionResourceModel) {
	if internaltypes.IsDefined(plan.GenerateAccountStatusNotification) {
		addRequest.GenerateAccountStatusNotification = plan.GenerateAccountStatusNotification.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a NoOperationFailureLockoutActionResponse object into the model struct
func readNoOperationFailureLockoutActionResponse(ctx context.Context, r *client.NoOperationFailureLockoutActionResponse, state *noOperationFailureLockoutActionResourceModel, expectedValues *noOperationFailureLockoutActionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.GenerateAccountStatusNotification = internaltypes.BoolTypeOrNil(r.GenerateAccountStatusNotification)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createNoOperationFailureLockoutActionOperations(plan noOperationFailureLockoutActionResourceModel, state noOperationFailureLockoutActionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerateAccountStatusNotification, state.GenerateAccountStatusNotification, "generate-account-status-notification")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *noOperationFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan noOperationFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddNoOperationFailureLockoutActionRequest(plan.Id.ValueString(),
		[]client.EnumnoOperationFailureLockoutActionSchemaUrn{client.ENUMNOOPERATIONFAILURELOCKOUTACTIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0FAILURE_LOCKOUT_ACTIONNO_OPERATION})
	addOptionalNoOperationFailureLockoutActionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.FailureLockoutActionApi.AddFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFailureLockoutActionRequest(
		client.AddNoOperationFailureLockoutActionRequestAsAddFailureLockoutActionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.AddFailureLockoutActionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the No Operation Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state noOperationFailureLockoutActionResourceModel
	readNoOperationFailureLockoutActionResponse(ctx, addResponse.NoOperationFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultNoOperationFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan noOperationFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the No Operation Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state noOperationFailureLockoutActionResourceModel
	readNoOperationFailureLockoutActionResponse(ctx, readResponse.NoOperationFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createNoOperationFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the No Operation Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readNoOperationFailureLockoutActionResponse(ctx, updateResponse.NoOperationFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *noOperationFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readNoOperationFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultNoOperationFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readNoOperationFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readNoOperationFailureLockoutAction(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state noOperationFailureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the No Operation Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readNoOperationFailureLockoutActionResponse(ctx, readResponse.NoOperationFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *noOperationFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateNoOperationFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultNoOperationFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateNoOperationFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateNoOperationFailureLockoutAction(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan noOperationFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state noOperationFailureLockoutActionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createNoOperationFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the No Operation Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readNoOperationFailureLockoutActionResponse(ctx, updateResponse.NoOperationFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultNoOperationFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *noOperationFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state noOperationFailureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutActionExecute(r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the No Operation Failure Lockout Action", err, httpResp)
		return
	}
}

func (r *noOperationFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importNoOperationFailureLockoutAction(ctx, req, resp)
}

func (r *defaultNoOperationFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importNoOperationFailureLockoutAction(ctx, req, resp)
}

func importNoOperationFailureLockoutAction(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
