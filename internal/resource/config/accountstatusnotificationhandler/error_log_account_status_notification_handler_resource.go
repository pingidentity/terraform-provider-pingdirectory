package accountstatusnotificationhandler

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
	_ resource.Resource                = &errorLogAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &errorLogAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &errorLogAccountStatusNotificationHandlerResource{}
	_ resource.Resource                = &defaultErrorLogAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultErrorLogAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultErrorLogAccountStatusNotificationHandlerResource{}
)

// Create a Error Log Account Status Notification Handler resource
func NewErrorLogAccountStatusNotificationHandlerResource() resource.Resource {
	return &errorLogAccountStatusNotificationHandlerResource{}
}

func NewDefaultErrorLogAccountStatusNotificationHandlerResource() resource.Resource {
	return &defaultErrorLogAccountStatusNotificationHandlerResource{}
}

// errorLogAccountStatusNotificationHandlerResource is the resource implementation.
type errorLogAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultErrorLogAccountStatusNotificationHandlerResource is the resource implementation.
type defaultErrorLogAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *errorLogAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_error_log_account_status_notification_handler"
}

func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_error_log_account_status_notification_handler"
}

// Configure adds the provider configured client to the resource.
func (r *errorLogAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type errorLogAccountStatusNotificationHandlerResourceModel struct {
	Id                                         types.String `tfsdk:"id"`
	LastUpdated                                types.String `tfsdk:"last_updated"`
	Notifications                              types.Set    `tfsdk:"notifications"`
	RequiredActions                            types.Set    `tfsdk:"required_actions"`
	AccountStatusNotificationType              types.Set    `tfsdk:"account_status_notification_type"`
	Description                                types.String `tfsdk:"description"`
	Enabled                                    types.Bool   `tfsdk:"enabled"`
	Asynchronous                               types.Bool   `tfsdk:"asynchronous"`
	AccountCreationNotificationRequestCriteria types.String `tfsdk:"account_creation_notification_request_criteria"`
	AccountUpdateNotificationRequestCriteria   types.String `tfsdk:"account_update_notification_request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *errorLogAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	errorLogAccountStatusNotificationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	errorLogAccountStatusNotificationHandlerSchema(ctx, req, resp, true)
}

func errorLogAccountStatusNotificationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Error Log Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
			"account_status_notification_type": schema.SetAttribute{
				Description: "Indicates which types of event can trigger an account status notification.",
				Required:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Account Status Notification Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Account Status Notification Handler is enabled. Only enabled handlers are invoked whenever a related event occurs in the server.",
				Required:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to invoke this Account Status Notification Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver a message) will not delay processing for the operation that triggered the notification.",
				Optional:    true,
				Computed:    true,
			},
			"account_creation_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which add requests should result in account creation notifications for this handler.",
				Optional:    true,
				Computed:    true,
			},
			"account_update_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which modify and modify DN requests should result in account update notifications for this handler.",
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
func addOptionalErrorLogAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddErrorLogAccountStatusNotificationHandlerRequest, plan errorLogAccountStatusNotificationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		boolVal := plan.Asynchronous.ValueBool()
		addRequest.Asynchronous = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		stringVal := plan.AccountCreationNotificationRequestCriteria.ValueString()
		addRequest.AccountCreationNotificationRequestCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		stringVal := plan.AccountUpdateNotificationRequestCriteria.ValueString()
		addRequest.AccountUpdateNotificationRequestCriteria = &stringVal
	}
}

// Read a ErrorLogAccountStatusNotificationHandlerResponse object into the model struct
func readErrorLogAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.ErrorLogAccountStatusNotificationHandlerResponse, state *errorLogAccountStatusNotificationHandlerResourceModel, expectedValues *errorLogAccountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AccountStatusNotificationType = internaltypes.GetStringSet(
		client.StringSliceEnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp(r.AccountStatusNotificationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createErrorLogAccountStatusNotificationHandlerOperations(plan errorLogAccountStatusNotificationHandlerResourceModel, state errorLogAccountStatusNotificationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccountStatusNotificationType, state.AccountStatusNotificationType, "account-status-notification-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreationNotificationRequestCriteria, state.AccountCreationNotificationRequestCriteria, "account-creation-notification-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdateNotificationRequestCriteria, state.AccountUpdateNotificationRequestCriteria, "account-update-notification-request-criteria")
	return ops
}

// Create a new resource
func (r *errorLogAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan errorLogAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var AccountStatusNotificationTypeSlice []client.EnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp
	plan.AccountStatusNotificationType.ElementsAs(ctx, &AccountStatusNotificationTypeSlice, false)
	addRequest := client.NewAddErrorLogAccountStatusNotificationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumerrorLogAccountStatusNotificationHandlerSchemaUrn{client.ENUMERRORLOGACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERERROR_LOG},
		AccountStatusNotificationTypeSlice,
		plan.Enabled.ValueBool())
	addOptionalErrorLogAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddErrorLogAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Error Log Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state errorLogAccountStatusNotificationHandlerResourceModel
	readErrorLogAccountStatusNotificationHandlerResponse(ctx, addResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan errorLogAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Error Log Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state errorLogAccountStatusNotificationHandlerResourceModel
	readErrorLogAccountStatusNotificationHandlerResponse(ctx, readResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createErrorLogAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Error Log Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readErrorLogAccountStatusNotificationHandlerResponse(ctx, updateResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *errorLogAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readErrorLogAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readErrorLogAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readErrorLogAccountStatusNotificationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state errorLogAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Error Log Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readErrorLogAccountStatusNotificationHandlerResponse(ctx, readResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *errorLogAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateErrorLogAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateErrorLogAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateErrorLogAccountStatusNotificationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan errorLogAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state errorLogAccountStatusNotificationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createErrorLogAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Error Log Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readErrorLogAccountStatusNotificationHandlerResponse(ctx, updateResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultErrorLogAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *errorLogAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state errorLogAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandlerExecute(r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Error Log Account Status Notification Handler", err, httpResp)
		return
	}
}

func (r *errorLogAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importErrorLogAccountStatusNotificationHandler(ctx, req, resp)
}

func (r *defaultErrorLogAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importErrorLogAccountStatusNotificationHandler(ctx, req, resp)
}

func importErrorLogAccountStatusNotificationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
