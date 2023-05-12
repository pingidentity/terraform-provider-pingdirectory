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
	_ resource.Resource                = &delayBindResponseFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &delayBindResponseFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &delayBindResponseFailureLockoutActionResource{}
	_ resource.Resource                = &defaultDelayBindResponseFailureLockoutActionResource{}
	_ resource.ResourceWithConfigure   = &defaultDelayBindResponseFailureLockoutActionResource{}
	_ resource.ResourceWithImportState = &defaultDelayBindResponseFailureLockoutActionResource{}
)

// Create a Delay Bind Response Failure Lockout Action resource
func NewDelayBindResponseFailureLockoutActionResource() resource.Resource {
	return &delayBindResponseFailureLockoutActionResource{}
}

func NewDefaultDelayBindResponseFailureLockoutActionResource() resource.Resource {
	return &defaultDelayBindResponseFailureLockoutActionResource{}
}

// delayBindResponseFailureLockoutActionResource is the resource implementation.
type delayBindResponseFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDelayBindResponseFailureLockoutActionResource is the resource implementation.
type defaultDelayBindResponseFailureLockoutActionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delayBindResponseFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delay_bind_response_failure_lockout_action"
}

func (r *defaultDelayBindResponseFailureLockoutActionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_delay_bind_response_failure_lockout_action"
}

// Configure adds the provider configured client to the resource.
func (r *delayBindResponseFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultDelayBindResponseFailureLockoutActionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type delayBindResponseFailureLockoutActionResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Delay                             types.String `tfsdk:"delay"`
	AllowBlockingDelay                types.Bool   `tfsdk:"allow_blocking_delay"`
	GenerateAccountStatusNotification types.Bool   `tfsdk:"generate_account_status_notification"`
	Description                       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *delayBindResponseFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delayBindResponseFailureLockoutActionSchema(ctx, req, resp, false)
}

func (r *defaultDelayBindResponseFailureLockoutActionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	delayBindResponseFailureLockoutActionSchema(ctx, req, resp, true)
}

func delayBindResponseFailureLockoutActionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Delay Bind Response Failure Lockout Action.",
		Attributes: map[string]schema.Attribute{
			"delay": schema.StringAttribute{
				Description: "The length of time to delay the bind response for accounts with too many failed authentication attempts.",
				Required:    true,
			},
			"allow_blocking_delay": schema.BoolAttribute{
				Description: "Indicates whether to delay the response for authentication attempts even if that delay may block the thread being used to process the attempt.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"generate_account_status_notification": schema.BoolAttribute{
				Description: "Indicates whether to generate an account status notification for cases in which a bind response is delayed because of failure lockout.",
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
func addOptionalDelayBindResponseFailureLockoutActionFields(ctx context.Context, addRequest *client.AddDelayBindResponseFailureLockoutActionRequest, plan delayBindResponseFailureLockoutActionResourceModel) {
	if internaltypes.IsDefined(plan.AllowBlockingDelay) {
		addRequest.AllowBlockingDelay = plan.AllowBlockingDelay.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerateAccountStatusNotification) {
		addRequest.GenerateAccountStatusNotification = plan.GenerateAccountStatusNotification.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a DelayBindResponseFailureLockoutActionResponse object into the model struct
func readDelayBindResponseFailureLockoutActionResponse(ctx context.Context, r *client.DelayBindResponseFailureLockoutActionResponse, state *delayBindResponseFailureLockoutActionResourceModel, expectedValues *delayBindResponseFailureLockoutActionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Delay = types.StringValue(r.Delay)
	config.CheckMismatchedPDFormattedAttributes("delay",
		expectedValues.Delay, state.Delay, diagnostics)
	state.AllowBlockingDelay = internaltypes.BoolTypeOrNil(r.AllowBlockingDelay)
	state.GenerateAccountStatusNotification = internaltypes.BoolTypeOrNil(r.GenerateAccountStatusNotification)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDelayBindResponseFailureLockoutActionOperations(plan delayBindResponseFailureLockoutActionResourceModel, state delayBindResponseFailureLockoutActionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Delay, state.Delay, "delay")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowBlockingDelay, state.AllowBlockingDelay, "allow-blocking-delay")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerateAccountStatusNotification, state.GenerateAccountStatusNotification, "generate-account-status-notification")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *delayBindResponseFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delayBindResponseFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddDelayBindResponseFailureLockoutActionRequest(plan.Id.ValueString(),
		[]client.EnumdelayBindResponseFailureLockoutActionSchemaUrn{client.ENUMDELAYBINDRESPONSEFAILURELOCKOUTACTIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0FAILURE_LOCKOUT_ACTIONDELAY_BIND_RESPONSE},
		plan.Delay.ValueString())
	addOptionalDelayBindResponseFailureLockoutActionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.FailureLockoutActionApi.AddFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddFailureLockoutActionRequest(
		client.AddDelayBindResponseFailureLockoutActionRequestAsAddFailureLockoutActionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.AddFailureLockoutActionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Delay Bind Response Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state delayBindResponseFailureLockoutActionResourceModel
	readDelayBindResponseFailureLockoutActionResponse(ctx, addResponse.DelayBindResponseFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultDelayBindResponseFailureLockoutActionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delayBindResponseFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delay Bind Response Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delayBindResponseFailureLockoutActionResourceModel
	readDelayBindResponseFailureLockoutActionResponse(ctx, readResponse.DelayBindResponseFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDelayBindResponseFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delay Bind Response Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelayBindResponseFailureLockoutActionResponse(ctx, updateResponse.DelayBindResponseFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *delayBindResponseFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelayBindResponseFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelayBindResponseFailureLockoutActionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDelayBindResponseFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDelayBindResponseFailureLockoutAction(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state delayBindResponseFailureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.FailureLockoutActionApi.GetFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delay Bind Response Failure Lockout Action", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelayBindResponseFailureLockoutActionResponse(ctx, readResponse.DelayBindResponseFailureLockoutActionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *delayBindResponseFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelayBindResponseFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDelayBindResponseFailureLockoutActionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDelayBindResponseFailureLockoutAction(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDelayBindResponseFailureLockoutAction(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan delayBindResponseFailureLockoutActionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delayBindResponseFailureLockoutActionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.FailureLockoutActionApi.UpdateFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDelayBindResponseFailureLockoutActionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.FailureLockoutActionApi.UpdateFailureLockoutActionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delay Bind Response Failure Lockout Action", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelayBindResponseFailureLockoutActionResponse(ctx, updateResponse.DelayBindResponseFailureLockoutActionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDelayBindResponseFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delayBindResponseFailureLockoutActionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state delayBindResponseFailureLockoutActionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutActionExecute(r.apiClient.FailureLockoutActionApi.DeleteFailureLockoutAction(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Delay Bind Response Failure Lockout Action", err, httpResp)
		return
	}
}

func (r *delayBindResponseFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelayBindResponseFailureLockoutAction(ctx, req, resp)
}

func (r *defaultDelayBindResponseFailureLockoutActionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDelayBindResponseFailureLockoutAction(ctx, req, resp)
}

func importDelayBindResponseFailureLockoutAction(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
