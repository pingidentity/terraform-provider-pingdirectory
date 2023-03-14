package accountstatusnotificationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	_ resource.Resource                = &smtpAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &smtpAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &smtpAccountStatusNotificationHandlerResource{}
	_ resource.Resource                = &defaultSmtpAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultSmtpAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultSmtpAccountStatusNotificationHandlerResource{}
)

// Create a Smtp Account Status Notification Handler resource
func NewSmtpAccountStatusNotificationHandlerResource() resource.Resource {
	return &smtpAccountStatusNotificationHandlerResource{}
}

func NewDefaultSmtpAccountStatusNotificationHandlerResource() resource.Resource {
	return &defaultSmtpAccountStatusNotificationHandlerResource{}
}

// smtpAccountStatusNotificationHandlerResource is the resource implementation.
type smtpAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSmtpAccountStatusNotificationHandlerResource is the resource implementation.
type defaultSmtpAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *smtpAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_smtp_account_status_notification_handler"
}

func (r *defaultSmtpAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_smtp_account_status_notification_handler"
}

// Configure adds the provider configured client to the resource.
func (r *smtpAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSmtpAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type smtpAccountStatusNotificationHandlerResourceModel struct {
	Id                                         types.String `tfsdk:"id"`
	LastUpdated                                types.String `tfsdk:"last_updated"`
	Notifications                              types.Set    `tfsdk:"notifications"`
	RequiredActions                            types.Set    `tfsdk:"required_actions"`
	EmailAddressAttributeType                  types.Set    `tfsdk:"email_address_attribute_type"`
	EmailAddressJSONField                      types.String `tfsdk:"email_address_json_field"`
	EmailAddressJSONObjectFilter               types.String `tfsdk:"email_address_json_object_filter"`
	RecipientAddress                           types.Set    `tfsdk:"recipient_address"`
	SendMessageWithoutEndUserAddress           types.Bool   `tfsdk:"send_message_without_end_user_address"`
	SenderAddress                              types.String `tfsdk:"sender_address"`
	MessageSubject                             types.Set    `tfsdk:"message_subject"`
	MessageTemplateFile                        types.Set    `tfsdk:"message_template_file"`
	Description                                types.String `tfsdk:"description"`
	Enabled                                    types.Bool   `tfsdk:"enabled"`
	Asynchronous                               types.Bool   `tfsdk:"asynchronous"`
	AccountCreationNotificationRequestCriteria types.String `tfsdk:"account_creation_notification_request_criteria"`
	AccountUpdateNotificationRequestCriteria   types.String `tfsdk:"account_update_notification_request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *smtpAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	smtpAccountStatusNotificationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultSmtpAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	smtpAccountStatusNotificationHandlerSchema(ctx, req, resp, true)
}

func smtpAccountStatusNotificationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Smtp Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
			"email_address_attribute_type": schema.SetAttribute{
				Description: "Specifies which attribute in the user's entries may be used to obtain the email address when notifying the end user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"email_address_json_field": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the email address to which the message should be sent. The email address must be contained in a top-level field whose value is a single string.",
				Optional:    true,
			},
			"email_address_json_object_filter": schema.StringAttribute{
				Description: "A JSON object filter that may be used to identify which email address value to use when sending the message.",
				Optional:    true,
			},
			"recipient_address": schema.SetAttribute{
				Description: "Specifies an email address to which notification messages are sent, either instead of or in addition to the end user for whom the notification has been generated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"send_message_without_end_user_address": schema.BoolAttribute{
				Description: "Indicates whether an email notification message should be generated and sent to the set of notification recipients even if the user entry does not contain any values for any of the email address attributes (that is, in cases when it is not possible to notify the end user).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"sender_address": schema.StringAttribute{
				Description: "Specifies the email address from which the message is sent. Note that this does not necessarily have to be a legitimate email address.",
				Required:    true,
			},
			"message_subject": schema.SetAttribute{
				Description: "Specifies the subject that should be used for email messages generated by this account status notification handler.",
				Required:    true,
				ElementType: types.StringType,
			},
			"message_template_file": schema.SetAttribute{
				Description: "Specifies the path to the file containing the message template to generate the email notification messages.",
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"account_creation_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which add requests should result in account creation notifications for this handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_update_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which modify and modify DN requests should result in account update notifications for this handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
func addOptionalSmtpAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddSmtpAccountStatusNotificationHandlerRequest, plan smtpAccountStatusNotificationHandlerResourceModel) {
	if internaltypes.IsDefined(plan.EmailAddressAttributeType) {
		var slice []string
		plan.EmailAddressAttributeType.ElementsAs(ctx, &slice, false)
		addRequest.EmailAddressAttributeType = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressJSONField) {
		stringVal := plan.EmailAddressJSONField.ValueString()
		addRequest.EmailAddressJSONField = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressJSONObjectFilter) {
		stringVal := plan.EmailAddressJSONObjectFilter.ValueString()
		addRequest.EmailAddressJSONObjectFilter = &stringVal
	}
	if internaltypes.IsDefined(plan.RecipientAddress) {
		var slice []string
		plan.RecipientAddress.ElementsAs(ctx, &slice, false)
		addRequest.RecipientAddress = slice
	}
	if internaltypes.IsDefined(plan.SendMessageWithoutEndUserAddress) {
		boolVal := plan.SendMessageWithoutEndUserAddress.ValueBool()
		addRequest.SendMessageWithoutEndUserAddress = &boolVal
	}
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

// Read a SmtpAccountStatusNotificationHandlerResponse object into the model struct
func readSmtpAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.SmtpAccountStatusNotificationHandlerResponse, state *smtpAccountStatusNotificationHandlerResourceModel, expectedValues *smtpAccountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.EmailAddressAttributeType = internaltypes.GetStringSet(r.EmailAddressAttributeType)
	state.EmailAddressJSONField = internaltypes.StringTypeOrNil(r.EmailAddressJSONField, internaltypes.IsEmptyString(expectedValues.EmailAddressJSONField))
	state.EmailAddressJSONObjectFilter = internaltypes.StringTypeOrNil(r.EmailAddressJSONObjectFilter, internaltypes.IsEmptyString(expectedValues.EmailAddressJSONObjectFilter))
	state.RecipientAddress = internaltypes.GetStringSet(r.RecipientAddress)
	state.SendMessageWithoutEndUserAddress = types.BoolValue(r.SendMessageWithoutEndUserAddress)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.MessageSubject = internaltypes.GetStringSet(r.MessageSubject)
	state.MessageTemplateFile = internaltypes.GetStringSet(r.MessageTemplateFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSmtpAccountStatusNotificationHandlerOperations(plan smtpAccountStatusNotificationHandlerResourceModel, state smtpAccountStatusNotificationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EmailAddressAttributeType, state.EmailAddressAttributeType, "email-address-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.EmailAddressJSONField, state.EmailAddressJSONField, "email-address-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.EmailAddressJSONObjectFilter, state.EmailAddressJSONObjectFilter, "email-address-json-object-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientAddress, state.RecipientAddress, "recipient-address")
	operations.AddBoolOperationIfNecessary(&ops, plan.SendMessageWithoutEndUserAddress, state.SendMessageWithoutEndUserAddress, "send-message-without-end-user-address")
	operations.AddStringOperationIfNecessary(&ops, plan.SenderAddress, state.SenderAddress, "sender-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MessageSubject, state.MessageSubject, "message-subject")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.MessageTemplateFile, state.MessageTemplateFile, "message-template-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreationNotificationRequestCriteria, state.AccountCreationNotificationRequestCriteria, "account-creation-notification-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdateNotificationRequestCriteria, state.AccountUpdateNotificationRequestCriteria, "account-update-notification-request-criteria")
	return ops
}

// Create a new resource
func (r *smtpAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan smtpAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var MessageSubjectSlice []string
	plan.MessageSubject.ElementsAs(ctx, &MessageSubjectSlice, false)
	var MessageTemplateFileSlice []string
	plan.MessageTemplateFile.ElementsAs(ctx, &MessageTemplateFileSlice, false)
	addRequest := client.NewAddSmtpAccountStatusNotificationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumsmtpAccountStatusNotificationHandlerSchemaUrn{client.ENUMSMTPACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERSMTP},
		plan.SenderAddress.ValueString(),
		MessageSubjectSlice,
		MessageTemplateFileSlice,
		plan.Enabled.ValueBool())
	addOptionalSmtpAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddSmtpAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Smtp Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state smtpAccountStatusNotificationHandlerResourceModel
	readSmtpAccountStatusNotificationHandlerResponse(ctx, addResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSmtpAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan smtpAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Smtp Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state smtpAccountStatusNotificationHandlerResourceModel
	readSmtpAccountStatusNotificationHandlerResponse(ctx, readResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSmtpAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Smtp Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSmtpAccountStatusNotificationHandlerResponse(ctx, updateResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *smtpAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSmtpAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSmtpAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSmtpAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSmtpAccountStatusNotificationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state smtpAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Smtp Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSmtpAccountStatusNotificationHandlerResponse(ctx, readResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *smtpAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSmtpAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSmtpAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSmtpAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSmtpAccountStatusNotificationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan smtpAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state smtpAccountStatusNotificationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSmtpAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Smtp Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSmtpAccountStatusNotificationHandlerResponse(ctx, updateResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSmtpAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *smtpAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state smtpAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandlerExecute(r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Smtp Account Status Notification Handler", err, httpResp)
		return
	}
}

func (r *smtpAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSmtpAccountStatusNotificationHandler(ctx, req, resp)
}

func (r *defaultSmtpAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSmtpAccountStatusNotificationHandler(ctx, req, resp)
}

func importSmtpAccountStatusNotificationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
