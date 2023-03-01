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
	_ resource.Resource                = &adminAlertAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &adminAlertAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &adminAlertAccountStatusNotificationHandlerResource{}
	_ resource.Resource                = &defaultAdminAlertAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultAdminAlertAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultAdminAlertAccountStatusNotificationHandlerResource{}
)

// Create a Admin Alert Account Status Notification Handler resource
func NewAdminAlertAccountStatusNotificationHandlerResource() resource.Resource {
	return &adminAlertAccountStatusNotificationHandlerResource{}
}

func NewDefaultAdminAlertAccountStatusNotificationHandlerResource() resource.Resource {
	return &defaultAdminAlertAccountStatusNotificationHandlerResource{}
}

// adminAlertAccountStatusNotificationHandlerResource is the resource implementation.
type adminAlertAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAdminAlertAccountStatusNotificationHandlerResource is the resource implementation.
type defaultAdminAlertAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *adminAlertAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admin_alert_account_status_notification_handler"
}

func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_admin_alert_account_status_notification_handler"
}

// Configure adds the provider configured client to the resource.
func (r *adminAlertAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type adminAlertAccountStatusNotificationHandlerResourceModel struct {
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
func (r *adminAlertAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	adminAlertAccountStatusNotificationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	adminAlertAccountStatusNotificationHandlerSchema(ctx, req, resp, true)
}

func adminAlertAccountStatusNotificationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Admin Alert Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
			"account_status_notification_type": schema.SetAttribute{
				Description: "The types of account status notifications that should result in administrative alerts.",
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
	config.AddCommonSchema(&schema, true)
	if setOptionalToComputed {
		config.SetOptionalAttributesToComputed(&schema)
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalAdminAlertAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddAdminAlertAccountStatusNotificationHandlerRequest, plan adminAlertAccountStatusNotificationHandlerResourceModel) {
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

// Read a AdminAlertAccountStatusNotificationHandlerResponse object into the model struct
func readAdminAlertAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.AdminAlertAccountStatusNotificationHandlerResponse, state *adminAlertAccountStatusNotificationHandlerResourceModel, expectedValues *adminAlertAccountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
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
func createAdminAlertAccountStatusNotificationHandlerOperations(plan adminAlertAccountStatusNotificationHandlerResourceModel, state adminAlertAccountStatusNotificationHandlerResourceModel) []client.Operation {
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
func (r *adminAlertAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan adminAlertAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var AccountStatusNotificationTypeSlice []client.EnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp
	plan.AccountStatusNotificationType.ElementsAs(ctx, &AccountStatusNotificationTypeSlice, false)
	addRequest := client.NewAddAdminAlertAccountStatusNotificationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumadminAlertAccountStatusNotificationHandlerSchemaUrn{client.ENUMADMINALERTACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERADMIN_ALERT},
		AccountStatusNotificationTypeSlice,
		plan.Enabled.ValueBool())
	addOptionalAdminAlertAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddAdminAlertAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Admin Alert Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state adminAlertAccountStatusNotificationHandlerResourceModel
	readAdminAlertAccountStatusNotificationHandlerResponse(ctx, addResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan adminAlertAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Admin Alert Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state adminAlertAccountStatusNotificationHandlerResourceModel
	readAdminAlertAccountStatusNotificationHandlerResponse(ctx, readResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAdminAlertAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Admin Alert Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAdminAlertAccountStatusNotificationHandlerResponse(ctx, updateResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *adminAlertAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAdminAlertAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAdminAlertAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAdminAlertAccountStatusNotificationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state adminAlertAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Admin Alert Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAdminAlertAccountStatusNotificationHandlerResponse(ctx, readResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *adminAlertAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAdminAlertAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAdminAlertAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAdminAlertAccountStatusNotificationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan adminAlertAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state adminAlertAccountStatusNotificationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAdminAlertAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Admin Alert Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAdminAlertAccountStatusNotificationHandlerResponse(ctx, updateResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *adminAlertAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state adminAlertAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandlerExecute(r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Admin Alert Account Status Notification Handler", err, httpResp)
		return
	}
}

func (r *adminAlertAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAdminAlertAccountStatusNotificationHandler(ctx, req, resp)
}

func (r *defaultAdminAlertAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAdminAlertAccountStatusNotificationHandler(ctx, req, resp)
}

func importAdminAlertAccountStatusNotificationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
