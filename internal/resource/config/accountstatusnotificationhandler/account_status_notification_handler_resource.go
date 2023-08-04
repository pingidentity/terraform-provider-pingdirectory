package accountstatusnotificationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &accountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &accountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &accountStatusNotificationHandlerResource{}
	_ resource.Resource                = &defaultAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultAccountStatusNotificationHandlerResource{}
)

// Create a Account Status Notification Handler resource
func NewAccountStatusNotificationHandlerResource() resource.Resource {
	return &accountStatusNotificationHandlerResource{}
}

func NewDefaultAccountStatusNotificationHandlerResource() resource.Resource {
	return &defaultAccountStatusNotificationHandlerResource{}
}

// accountStatusNotificationHandlerResource is the resource implementation.
type accountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAccountStatusNotificationHandlerResource is the resource implementation.
type defaultAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *accountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_status_notification_handler"
}

func (r *defaultAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_account_status_notification_handler"
}

// Configure adds the provider configured client to the resource.
func (r *accountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type accountStatusNotificationHandlerResourceModel struct {
	Id                                              types.String `tfsdk:"id"`
	Name                                            types.String `tfsdk:"name"`
	LastUpdated                                     types.String `tfsdk:"last_updated"`
	Notifications                                   types.Set    `tfsdk:"notifications"`
	RequiredActions                                 types.Set    `tfsdk:"required_actions"`
	Type                                            types.String `tfsdk:"type"`
	ExtensionClass                                  types.String `tfsdk:"extension_class"`
	ExtensionArgument                               types.Set    `tfsdk:"extension_argument"`
	AccountTemporarilyFailureLockedMessageTemplate  types.String `tfsdk:"account_temporarily_failure_locked_message_template"`
	AccountPermanentlyFailureLockedMessageTemplate  types.String `tfsdk:"account_permanently_failure_locked_message_template"`
	AccountIdleLockedMessageTemplate                types.String `tfsdk:"account_idle_locked_message_template"`
	AccountResetLockedMessageTemplate               types.String `tfsdk:"account_reset_locked_message_template"`
	AccountUnlockedMessageTemplate                  types.String `tfsdk:"account_unlocked_message_template"`
	AccountDisabledMessageTemplate                  types.String `tfsdk:"account_disabled_message_template"`
	AccountEnabledMessageTemplate                   types.String `tfsdk:"account_enabled_message_template"`
	AccountNotYetActiveMessageTemplate              types.String `tfsdk:"account_not_yet_active_message_template"`
	AccountExpiredMessageTemplate                   types.String `tfsdk:"account_expired_message_template"`
	PasswordExpiredMessageTemplate                  types.String `tfsdk:"password_expired_message_template"`
	PasswordExpiringMessageTemplate                 types.String `tfsdk:"password_expiring_message_template"`
	PasswordResetMessageTemplate                    types.String `tfsdk:"password_reset_message_template"`
	PasswordChangedMessageTemplate                  types.String `tfsdk:"password_changed_message_template"`
	AccountAuthenticatedMessageTemplate             types.String `tfsdk:"account_authenticated_message_template"`
	AccountCreatedMessageTemplate                   types.String `tfsdk:"account_created_message_template"`
	AccountDeletedMessageTemplate                   types.String `tfsdk:"account_deleted_message_template"`
	AccountUpdatedMessageTemplate                   types.String `tfsdk:"account_updated_message_template"`
	BindPasswordFailedValidationMessageTemplate     types.String `tfsdk:"bind_password_failed_validation_message_template"`
	MustChangePasswordMessageTemplate               types.String `tfsdk:"must_change_password_message_template"`
	AccountStatusNotificationType                   types.Set    `tfsdk:"account_status_notification_type"`
	ScriptClass                                     types.String `tfsdk:"script_class"`
	ScriptArgument                                  types.Set    `tfsdk:"script_argument"`
	EmailAddressAttributeType                       types.Set    `tfsdk:"email_address_attribute_type"`
	EmailAddressJSONField                           types.String `tfsdk:"email_address_json_field"`
	EmailAddressJSONObjectFilter                    types.String `tfsdk:"email_address_json_object_filter"`
	RecipientAddress                                types.Set    `tfsdk:"recipient_address"`
	SendMessageWithoutEndUserAddress                types.Bool   `tfsdk:"send_message_without_end_user_address"`
	SenderAddress                                   types.String `tfsdk:"sender_address"`
	MessageSubject                                  types.Set    `tfsdk:"message_subject"`
	MessageTemplateFile                             types.Set    `tfsdk:"message_template_file"`
	Description                                     types.String `tfsdk:"description"`
	Enabled                                         types.Bool   `tfsdk:"enabled"`
	Asynchronous                                    types.Bool   `tfsdk:"asynchronous"`
	AccountAuthenticationNotificationResultCriteria types.String `tfsdk:"account_authentication_notification_result_criteria"`
	AccountCreationNotificationRequestCriteria      types.String `tfsdk:"account_creation_notification_request_criteria"`
	AccountDeletionNotificationRequestCriteria      types.String `tfsdk:"account_deletion_notification_request_criteria"`
	AccountUpdateNotificationRequestCriteria        types.String `tfsdk:"account_update_notification_request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *accountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	accountStatusNotificationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	accountStatusNotificationHandlerSchema(ctx, req, resp, true)
}

func accountStatusNotificationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Account Status Notification Handler resource. Options are ['smtp', 'groovy-scripted', 'admin-alert', 'error-log', 'multi-part-email', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"smtp", "groovy-scripted", "admin-alert", "error-log", "multi-part-email", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Account Status Notification Handler.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Account Status Notification Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"account_temporarily_failure_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an account becomes temporarily locked as a result of too many authentication failures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_permanently_failure_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an account becomes permanently locked as a result of too many authentication failures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_idle_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that authentication attempt fails because it has been too long since the user last successfully authenticated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_reset_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that authentication attempt fails because the user failed to choose a new password in a timely manner after an administrative reset.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_unlocked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's account has been unlocked (e.g., by an administrative password reset).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_disabled_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's account is disabled by an administrator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_enabled_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's account is enabled by an administrator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_not_yet_active_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt fails because the account has an activation time that is in the future.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_expired_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt fails because the account has an expiration time that is in the past.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_expired_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt fails because the account has an expired password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_expiring_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt succeeds, but the user's password is about to expire. This notification will only be generated the first time the user authenticates within the window of time that the server should warn about an upcoming password expiration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_reset_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's password has been reset by an administrator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_changed_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user changes their own password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_authenticated_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an account has successfully authenticated in a bind operation that matches the criteria provided in the account-authentication-notification-request-criteria property. Supported in PingDirectory product version 9.3.0.0+.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_created_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a new account is created in an add request that matches the criteria provided in the account-creation-notification-request-criteria property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_deleted_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an existing accout has been removed in a delete request that matches the criteria provided in the account-deletion-notification-request-criteria property. Supported in PingDirectory product version 9.3.0.0+.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_updated_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an existing account is updated with a modify or modify DN operation that matches the criteria provided in the account-update-notification-request-criteria property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_password_failed_validation_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user authenticated with a password that failed to satisfy the criteria for one or more of the configured password validators.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"must_change_password_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user successfully authenticates to the server but will be required to choose a new password before they will be allowed to perform any other operations.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"account_status_notification_type": schema.SetAttribute{
				Description: "The types of account status notifications that should result in administrative alerts.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Account Status Notification Handler.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Account Status Notification Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"email_address_attribute_type": schema.SetAttribute{
				Description: "Specifies which attribute in the user's entries may be used to obtain the email address when notifying the end user.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				Optional:    true,
			},
			"message_subject": schema.SetAttribute{
				Description: "Specifies the subject that should be used for email messages generated by this account status notification handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"message_template_file": schema.SetAttribute{
				Description: "Specifies the path to the file containing the message template to generate the email notification messages.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
			"account_authentication_notification_result_criteria": schema.StringAttribute{
				Description: "A result criteria object that identifies which successful bind operations should result in account authentication notifications for this handler. Supported in PingDirectory product version 9.3.0.0+.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"account_deletion_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which delete requests should result in account deletion notifications for this handler. Supported in PingDirectory product version 9.3.0.0+.",
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
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputed(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *accountStatusNotificationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAccountStatusNotificationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanAccountStatusNotificationHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model accountStatusNotificationHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.SenderAddress) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'sender_address' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'sender_address', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.AccountUnlockedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_unlocked_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_unlocked_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.MustChangePasswordMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'must_change_password_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'must_change_password_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.BindPasswordFailedValidationMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'bind_password_failed_validation_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'bind_password_failed_validation_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.EmailAddressJSONField) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'email_address_json_field' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'email_address_json_field', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.EmailAddressAttributeType) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'email_address_attribute_type' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'email_address_attribute_type', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.AccountStatusNotificationType) && model.Type.ValueString() != "admin-alert" && model.Type.ValueString() != "error-log" {
		resp.Diagnostics.AddError("Attribute 'account_status_notification_type' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_status_notification_type', the 'type' attribute must be one of ['admin-alert', 'error-log']")
	}
	if internaltypes.IsDefined(model.AccountDeletedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_deleted_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_deleted_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.MessageTemplateFile) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'message_template_file' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'message_template_file', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.AccountExpiredMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_expired_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_expired_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.PasswordExpiringMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'password_expiring_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_expiring_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountIdleLockedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_idle_locked_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_idle_locked_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountNotYetActiveMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_not_yet_active_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_not_yet_active_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptArgument) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_argument' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_argument', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.AccountTemporarilyFailureLockedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_temporarily_failure_locked_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_temporarily_failure_locked_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountPermanentlyFailureLockedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_permanently_failure_locked_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_permanently_failure_locked_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.PasswordExpiredMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'password_expired_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_expired_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.PasswordResetMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'password_reset_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_reset_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountResetLockedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_reset_locked_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_reset_locked_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountEnabledMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_enabled_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_enabled_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.RecipientAddress) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'recipient_address' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'recipient_address', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.EmailAddressJSONObjectFilter) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'email_address_json_object_filter' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'email_address_json_object_filter', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.MessageSubject) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'message_subject' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'message_subject', the 'type' attribute must be one of ['smtp']")
	}
	if internaltypes.IsDefined(model.PasswordChangedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'password_changed_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'password_changed_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountAuthenticatedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_authenticated_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_authenticated_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountCreatedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_created_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_created_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountDisabledMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_disabled_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_disabled_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.AccountUpdatedMessageTemplate) && model.Type.ValueString() != "multi-part-email" {
		resp.Diagnostics.AddError("Attribute 'account_updated_message_template' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'account_updated_message_template', the 'type' attribute must be one of ['multi-part-email']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party']")
	}
	if internaltypes.IsDefined(model.ScriptClass) && model.Type.ValueString() != "groovy-scripted" {
		resp.Diagnostics.AddError("Attribute 'script_class' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_class', the 'type' attribute must be one of ['groovy-scripted']")
	}
	if internaltypes.IsDefined(model.SendMessageWithoutEndUserAddress) && model.Type.ValueString() != "smtp" {
		resp.Diagnostics.AddError("Attribute 'send_message_without_end_user_address' not supported by pingdirectory_account_status_notification_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'send_message_without_end_user_address', the 'type' attribute must be one of ['smtp']")
	}
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9300)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	if internaltypes.IsNonEmptyString(model.AccountAuthenticationNotificationResultCriteria) {
		resp.Diagnostics.AddError("Attribute 'account_authentication_notification_result_criteria' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.AccountDeletionNotificationRequestCriteria) {
		resp.Diagnostics.AddError("Attribute 'account_deletion_notification_request_criteria' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.AccountAuthenticatedMessageTemplate) {
		resp.Diagnostics.AddError("Attribute 'account_authenticated_message_template' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsNonEmptyString(model.AccountDeletedMessageTemplate) {
		resp.Diagnostics.AddError("Attribute 'account_deleted_message_template' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add optional fields to create request for smtp account-status-notification-handler
func addOptionalSmtpAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddSmtpAccountStatusNotificationHandlerRequest, plan accountStatusNotificationHandlerResourceModel) {
	if internaltypes.IsDefined(plan.EmailAddressAttributeType) {
		var slice []string
		plan.EmailAddressAttributeType.ElementsAs(ctx, &slice, false)
		addRequest.EmailAddressAttributeType = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressJSONField) {
		addRequest.EmailAddressJSONField = plan.EmailAddressJSONField.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressJSONObjectFilter) {
		addRequest.EmailAddressJSONObjectFilter = plan.EmailAddressJSONObjectFilter.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RecipientAddress) {
		var slice []string
		plan.RecipientAddress.ElementsAs(ctx, &slice, false)
		addRequest.RecipientAddress = slice
	}
	if internaltypes.IsDefined(plan.SendMessageWithoutEndUserAddress) {
		addRequest.SendMessageWithoutEndUserAddress = plan.SendMessageWithoutEndUserAddress.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticationNotificationResultCriteria) {
		addRequest.AccountAuthenticationNotificationResultCriteria = plan.AccountAuthenticationNotificationResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletionNotificationRequestCriteria) {
		addRequest.AccountDeletionNotificationRequestCriteria = plan.AccountDeletionNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Add optional fields to create request for groovy-scripted account-status-notification-handler
func addOptionalGroovyScriptedAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddGroovyScriptedAccountStatusNotificationHandlerRequest, plan accountStatusNotificationHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticationNotificationResultCriteria) {
		addRequest.AccountAuthenticationNotificationResultCriteria = plan.AccountAuthenticationNotificationResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletionNotificationRequestCriteria) {
		addRequest.AccountDeletionNotificationRequestCriteria = plan.AccountDeletionNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Add optional fields to create request for admin-alert account-status-notification-handler
func addOptionalAdminAlertAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddAdminAlertAccountStatusNotificationHandlerRequest, plan accountStatusNotificationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticationNotificationResultCriteria) {
		addRequest.AccountAuthenticationNotificationResultCriteria = plan.AccountAuthenticationNotificationResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletionNotificationRequestCriteria) {
		addRequest.AccountDeletionNotificationRequestCriteria = plan.AccountDeletionNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Add optional fields to create request for error-log account-status-notification-handler
func addOptionalErrorLogAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddErrorLogAccountStatusNotificationHandlerRequest, plan accountStatusNotificationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticationNotificationResultCriteria) {
		addRequest.AccountAuthenticationNotificationResultCriteria = plan.AccountAuthenticationNotificationResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletionNotificationRequestCriteria) {
		addRequest.AccountDeletionNotificationRequestCriteria = plan.AccountDeletionNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Add optional fields to create request for multi-part-email account-status-notification-handler
func addOptionalMultiPartEmailAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddMultiPartEmailAccountStatusNotificationHandlerRequest, plan accountStatusNotificationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountTemporarilyFailureLockedMessageTemplate) {
		addRequest.AccountTemporarilyFailureLockedMessageTemplate = plan.AccountTemporarilyFailureLockedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountPermanentlyFailureLockedMessageTemplate) {
		addRequest.AccountPermanentlyFailureLockedMessageTemplate = plan.AccountPermanentlyFailureLockedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountIdleLockedMessageTemplate) {
		addRequest.AccountIdleLockedMessageTemplate = plan.AccountIdleLockedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountResetLockedMessageTemplate) {
		addRequest.AccountResetLockedMessageTemplate = plan.AccountResetLockedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUnlockedMessageTemplate) {
		addRequest.AccountUnlockedMessageTemplate = plan.AccountUnlockedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDisabledMessageTemplate) {
		addRequest.AccountDisabledMessageTemplate = plan.AccountDisabledMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountEnabledMessageTemplate) {
		addRequest.AccountEnabledMessageTemplate = plan.AccountEnabledMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountNotYetActiveMessageTemplate) {
		addRequest.AccountNotYetActiveMessageTemplate = plan.AccountNotYetActiveMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountExpiredMessageTemplate) {
		addRequest.AccountExpiredMessageTemplate = plan.AccountExpiredMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordExpiredMessageTemplate) {
		addRequest.PasswordExpiredMessageTemplate = plan.PasswordExpiredMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordExpiringMessageTemplate) {
		addRequest.PasswordExpiringMessageTemplate = plan.PasswordExpiringMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordResetMessageTemplate) {
		addRequest.PasswordResetMessageTemplate = plan.PasswordResetMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordChangedMessageTemplate) {
		addRequest.PasswordChangedMessageTemplate = plan.PasswordChangedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticatedMessageTemplate) {
		addRequest.AccountAuthenticatedMessageTemplate = plan.AccountAuthenticatedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreatedMessageTemplate) {
		addRequest.AccountCreatedMessageTemplate = plan.AccountCreatedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletedMessageTemplate) {
		addRequest.AccountDeletedMessageTemplate = plan.AccountDeletedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdatedMessageTemplate) {
		addRequest.AccountUpdatedMessageTemplate = plan.AccountUpdatedMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindPasswordFailedValidationMessageTemplate) {
		addRequest.BindPasswordFailedValidationMessageTemplate = plan.BindPasswordFailedValidationMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MustChangePasswordMessageTemplate) {
		addRequest.MustChangePasswordMessageTemplate = plan.MustChangePasswordMessageTemplate.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticationNotificationResultCriteria) {
		addRequest.AccountAuthenticationNotificationResultCriteria = plan.AccountAuthenticationNotificationResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletionNotificationRequestCriteria) {
		addRequest.AccountDeletionNotificationRequestCriteria = plan.AccountDeletionNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party account-status-notification-handler
func addOptionalThirdPartyAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyAccountStatusNotificationHandlerRequest, plan accountStatusNotificationHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountAuthenticationNotificationResultCriteria) {
		addRequest.AccountAuthenticationNotificationResultCriteria = plan.AccountAuthenticationNotificationResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreationNotificationRequestCriteria) {
		addRequest.AccountCreationNotificationRequestCriteria = plan.AccountCreationNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDeletionNotificationRequestCriteria) {
		addRequest.AccountDeletionNotificationRequestCriteria = plan.AccountDeletionNotificationRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdateNotificationRequestCriteria) {
		addRequest.AccountUpdateNotificationRequestCriteria = plan.AccountUpdateNotificationRequestCriteria.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateAccountStatusNotificationHandlerUnknownValues(ctx context.Context, model *accountStatusNotificationHandlerResourceModel) {
	if model.MessageTemplateFile.ElementType(ctx) == nil {
		model.MessageTemplateFile = types.SetNull(types.StringType)
	}
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.MessageSubject.ElementType(ctx) == nil {
		model.MessageSubject = types.SetNull(types.StringType)
	}
	if model.RecipientAddress.ElementType(ctx) == nil {
		model.RecipientAddress = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.EmailAddressAttributeType.ElementType(ctx) == nil {
		model.EmailAddressAttributeType = types.SetNull(types.StringType)
	}
	if model.AccountStatusNotificationType.ElementType(ctx) == nil {
		model.AccountStatusNotificationType = types.SetNull(types.StringType)
	}
}

// Read a SmtpAccountStatusNotificationHandlerResponse object into the model struct
func readSmtpAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.SmtpAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerResourceModel, expectedValues *accountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, internaltypes.IsEmptyString(expectedValues.AccountAuthenticationNotificationResultCriteria))
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountDeletionNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccountStatusNotificationHandlerUnknownValues(ctx, state)
}

// Read a GroovyScriptedAccountStatusNotificationHandlerResponse object into the model struct
func readGroovyScriptedAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.GroovyScriptedAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerResourceModel, expectedValues *accountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, internaltypes.IsEmptyString(expectedValues.AccountAuthenticationNotificationResultCriteria))
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountDeletionNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccountStatusNotificationHandlerUnknownValues(ctx, state)
}

// Read a AdminAlertAccountStatusNotificationHandlerResponse object into the model struct
func readAdminAlertAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.AdminAlertAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerResourceModel, expectedValues *accountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("admin-alert")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccountStatusNotificationType = internaltypes.GetStringSet(
		client.StringSliceEnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp(r.AccountStatusNotificationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, internaltypes.IsEmptyString(expectedValues.AccountAuthenticationNotificationResultCriteria))
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountDeletionNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccountStatusNotificationHandlerUnknownValues(ctx, state)
}

// Read a ErrorLogAccountStatusNotificationHandlerResponse object into the model struct
func readErrorLogAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.ErrorLogAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerResourceModel, expectedValues *accountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error-log")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccountStatusNotificationType = internaltypes.GetStringSet(
		client.StringSliceEnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp(r.AccountStatusNotificationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, internaltypes.IsEmptyString(expectedValues.AccountAuthenticationNotificationResultCriteria))
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountDeletionNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccountStatusNotificationHandlerUnknownValues(ctx, state)
}

// Read a MultiPartEmailAccountStatusNotificationHandlerResponse object into the model struct
func readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.MultiPartEmailAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerResourceModel, expectedValues *accountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("multi-part-email")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccountTemporarilyFailureLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountTemporarilyFailureLockedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountTemporarilyFailureLockedMessageTemplate))
	state.AccountPermanentlyFailureLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountPermanentlyFailureLockedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountPermanentlyFailureLockedMessageTemplate))
	state.AccountIdleLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountIdleLockedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountIdleLockedMessageTemplate))
	state.AccountResetLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountResetLockedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountResetLockedMessageTemplate))
	state.AccountUnlockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountUnlockedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountUnlockedMessageTemplate))
	state.AccountDisabledMessageTemplate = internaltypes.StringTypeOrNil(r.AccountDisabledMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountDisabledMessageTemplate))
	state.AccountEnabledMessageTemplate = internaltypes.StringTypeOrNil(r.AccountEnabledMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountEnabledMessageTemplate))
	state.AccountNotYetActiveMessageTemplate = internaltypes.StringTypeOrNil(r.AccountNotYetActiveMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountNotYetActiveMessageTemplate))
	state.AccountExpiredMessageTemplate = internaltypes.StringTypeOrNil(r.AccountExpiredMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountExpiredMessageTemplate))
	state.PasswordExpiredMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordExpiredMessageTemplate, internaltypes.IsEmptyString(expectedValues.PasswordExpiredMessageTemplate))
	state.PasswordExpiringMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordExpiringMessageTemplate, internaltypes.IsEmptyString(expectedValues.PasswordExpiringMessageTemplate))
	state.PasswordResetMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordResetMessageTemplate, internaltypes.IsEmptyString(expectedValues.PasswordResetMessageTemplate))
	state.PasswordChangedMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordChangedMessageTemplate, internaltypes.IsEmptyString(expectedValues.PasswordChangedMessageTemplate))
	state.AccountAuthenticatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountAuthenticatedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountAuthenticatedMessageTemplate))
	state.AccountCreatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountCreatedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountCreatedMessageTemplate))
	state.AccountDeletedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountDeletedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountDeletedMessageTemplate))
	state.AccountUpdatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountUpdatedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountUpdatedMessageTemplate))
	state.BindPasswordFailedValidationMessageTemplate = internaltypes.StringTypeOrNil(r.BindPasswordFailedValidationMessageTemplate, internaltypes.IsEmptyString(expectedValues.BindPasswordFailedValidationMessageTemplate))
	state.MustChangePasswordMessageTemplate = internaltypes.StringTypeOrNil(r.MustChangePasswordMessageTemplate, internaltypes.IsEmptyString(expectedValues.MustChangePasswordMessageTemplate))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, internaltypes.IsEmptyString(expectedValues.AccountAuthenticationNotificationResultCriteria))
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountDeletionNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccountStatusNotificationHandlerUnknownValues(ctx, state)
}

// Read a ThirdPartyAccountStatusNotificationHandlerResponse object into the model struct
func readThirdPartyAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.ThirdPartyAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerResourceModel, expectedValues *accountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, internaltypes.IsEmptyString(expectedValues.AccountAuthenticationNotificationResultCriteria))
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountDeletionNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccountStatusNotificationHandlerUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createAccountStatusNotificationHandlerOperations(plan accountStatusNotificationHandlerResourceModel, state accountStatusNotificationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountTemporarilyFailureLockedMessageTemplate, state.AccountTemporarilyFailureLockedMessageTemplate, "account-temporarily-failure-locked-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountPermanentlyFailureLockedMessageTemplate, state.AccountPermanentlyFailureLockedMessageTemplate, "account-permanently-failure-locked-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountIdleLockedMessageTemplate, state.AccountIdleLockedMessageTemplate, "account-idle-locked-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountResetLockedMessageTemplate, state.AccountResetLockedMessageTemplate, "account-reset-locked-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUnlockedMessageTemplate, state.AccountUnlockedMessageTemplate, "account-unlocked-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountDisabledMessageTemplate, state.AccountDisabledMessageTemplate, "account-disabled-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountEnabledMessageTemplate, state.AccountEnabledMessageTemplate, "account-enabled-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountNotYetActiveMessageTemplate, state.AccountNotYetActiveMessageTemplate, "account-not-yet-active-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountExpiredMessageTemplate, state.AccountExpiredMessageTemplate, "account-expired-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordExpiredMessageTemplate, state.PasswordExpiredMessageTemplate, "password-expired-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordExpiringMessageTemplate, state.PasswordExpiringMessageTemplate, "password-expiring-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordResetMessageTemplate, state.PasswordResetMessageTemplate, "password-reset-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordChangedMessageTemplate, state.PasswordChangedMessageTemplate, "password-changed-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountAuthenticatedMessageTemplate, state.AccountAuthenticatedMessageTemplate, "account-authenticated-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreatedMessageTemplate, state.AccountCreatedMessageTemplate, "account-created-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountDeletedMessageTemplate, state.AccountDeletedMessageTemplate, "account-deleted-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdatedMessageTemplate, state.AccountUpdatedMessageTemplate, "account-updated-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.BindPasswordFailedValidationMessageTemplate, state.BindPasswordFailedValidationMessageTemplate, "bind-password-failed-validation-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.MustChangePasswordMessageTemplate, state.MustChangePasswordMessageTemplate, "must-change-password-message-template")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccountStatusNotificationType, state.AccountStatusNotificationType, "account-status-notification-type")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
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
	operations.AddStringOperationIfNecessary(&ops, plan.AccountAuthenticationNotificationResultCriteria, state.AccountAuthenticationNotificationResultCriteria, "account-authentication-notification-result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreationNotificationRequestCriteria, state.AccountCreationNotificationRequestCriteria, "account-creation-notification-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountDeletionNotificationRequestCriteria, state.AccountDeletionNotificationRequestCriteria, "account-deletion-notification-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdateNotificationRequestCriteria, state.AccountUpdateNotificationRequestCriteria, "account-update-notification-request-criteria")
	return ops
}

// Create a smtp account-status-notification-handler
func (r *accountStatusNotificationHandlerResource) CreateSmtpAccountStatusNotificationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accountStatusNotificationHandlerResourceModel) (*accountStatusNotificationHandlerResourceModel, error) {
	var MessageSubjectSlice []string
	plan.MessageSubject.ElementsAs(ctx, &MessageSubjectSlice, false)
	var MessageTemplateFileSlice []string
	plan.MessageTemplateFile.ElementsAs(ctx, &MessageTemplateFileSlice, false)
	addRequest := client.NewAddSmtpAccountStatusNotificationHandlerRequest(plan.Name.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Account Status Notification Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accountStatusNotificationHandlerResourceModel
	readSmtpAccountStatusNotificationHandlerResponse(ctx, addResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted account-status-notification-handler
func (r *accountStatusNotificationHandlerResource) CreateGroovyScriptedAccountStatusNotificationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accountStatusNotificationHandlerResourceModel) (*accountStatusNotificationHandlerResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedAccountStatusNotificationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedAccountStatusNotificationHandlerSchemaUrn{client.ENUMGROOVYSCRIPTEDACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddGroovyScriptedAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Account Status Notification Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accountStatusNotificationHandlerResourceModel
	readGroovyScriptedAccountStatusNotificationHandlerResponse(ctx, addResponse.GroovyScriptedAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a admin-alert account-status-notification-handler
func (r *accountStatusNotificationHandlerResource) CreateAdminAlertAccountStatusNotificationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accountStatusNotificationHandlerResourceModel) (*accountStatusNotificationHandlerResourceModel, error) {
	var AccountStatusNotificationTypeSlice []client.EnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp
	plan.AccountStatusNotificationType.ElementsAs(ctx, &AccountStatusNotificationTypeSlice, false)
	addRequest := client.NewAddAdminAlertAccountStatusNotificationHandlerRequest(plan.Name.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Account Status Notification Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accountStatusNotificationHandlerResourceModel
	readAdminAlertAccountStatusNotificationHandlerResponse(ctx, addResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a error-log account-status-notification-handler
func (r *accountStatusNotificationHandlerResource) CreateErrorLogAccountStatusNotificationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accountStatusNotificationHandlerResourceModel) (*accountStatusNotificationHandlerResourceModel, error) {
	var AccountStatusNotificationTypeSlice []client.EnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp
	plan.AccountStatusNotificationType.ElementsAs(ctx, &AccountStatusNotificationTypeSlice, false)
	addRequest := client.NewAddErrorLogAccountStatusNotificationHandlerRequest(plan.Name.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Account Status Notification Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accountStatusNotificationHandlerResourceModel
	readErrorLogAccountStatusNotificationHandlerResponse(ctx, addResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a multi-part-email account-status-notification-handler
func (r *accountStatusNotificationHandlerResource) CreateMultiPartEmailAccountStatusNotificationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accountStatusNotificationHandlerResourceModel) (*accountStatusNotificationHandlerResourceModel, error) {
	addRequest := client.NewAddMultiPartEmailAccountStatusNotificationHandlerRequest(plan.Name.ValueString(),
		[]client.EnummultiPartEmailAccountStatusNotificationHandlerSchemaUrn{client.ENUMMULTIPARTEMAILACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERMULTI_PART_EMAIL},
		plan.Enabled.ValueBool())
	addOptionalMultiPartEmailAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddMultiPartEmailAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Account Status Notification Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accountStatusNotificationHandlerResourceModel
	readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, addResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party account-status-notification-handler
func (r *accountStatusNotificationHandlerResource) CreateThirdPartyAccountStatusNotificationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accountStatusNotificationHandlerResourceModel) (*accountStatusNotificationHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyAccountStatusNotificationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyAccountStatusNotificationHandlerSchemaUrn{client.ENUMTHIRDPARTYACCOUNTSTATUSNOTIFICATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCOUNT_STATUS_NOTIFICATION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyAccountStatusNotificationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccountStatusNotificationHandlerRequest(
		client.AddThirdPartyAccountStatusNotificationHandlerRequestAsAddAccountStatusNotificationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.AddAccountStatusNotificationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Account Status Notification Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accountStatusNotificationHandlerResourceModel
	readThirdPartyAccountStatusNotificationHandlerResponse(ctx, addResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *accountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *accountStatusNotificationHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "smtp" {
		state, err = r.CreateSmtpAccountStatusNotificationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedAccountStatusNotificationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "admin-alert" {
		state, err = r.CreateAdminAlertAccountStatusNotificationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "error-log" {
		state, err = r.CreateErrorLogAccountStatusNotificationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "multi-part-email" {
		state, err = r.CreateMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyAccountStatusNotificationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state accountStatusNotificationHandlerResourceModel
	if plan.Type.ValueString() == "smtp" {
		readSmtpAccountStatusNotificationHandlerResponse(ctx, readResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		readGroovyScriptedAccountStatusNotificationHandlerResponse(ctx, readResponse.GroovyScriptedAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "admin-alert" {
		readAdminAlertAccountStatusNotificationHandlerResponse(ctx, readResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "error-log" {
		readErrorLogAccountStatusNotificationHandlerResponse(ctx, readResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "multi-part-email" {
		readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party" {
		readThirdPartyAccountStatusNotificationHandlerResponse(ctx, readResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "smtp" {
			readSmtpAccountStatusNotificationHandlerResponse(ctx, updateResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedAccountStatusNotificationHandlerResponse(ctx, updateResponse.GroovyScriptedAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "admin-alert" {
			readAdminAlertAccountStatusNotificationHandlerResponse(ctx, updateResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "error-log" {
			readErrorLogAccountStatusNotificationHandlerResponse(ctx, updateResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "multi-part-email" {
			readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, updateResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyAccountStatusNotificationHandlerResponse(ctx, updateResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *accountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAccountStatusNotificationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state accountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SmtpAccountStatusNotificationHandlerResponse != nil {
		readSmtpAccountStatusNotificationHandlerResponse(ctx, readResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAccountStatusNotificationHandlerResponse != nil {
		readGroovyScriptedAccountStatusNotificationHandlerResponse(ctx, readResponse.GroovyScriptedAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AdminAlertAccountStatusNotificationHandlerResponse != nil {
		readAdminAlertAccountStatusNotificationHandlerResponse(ctx, readResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogAccountStatusNotificationHandlerResponse != nil {
		readErrorLogAccountStatusNotificationHandlerResponse(ctx, readResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse != nil {
		readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccountStatusNotificationHandlerResponse != nil {
		readThirdPartyAccountStatusNotificationHandlerResponse(ctx, readResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *accountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAccountStatusNotificationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan accountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state accountStatusNotificationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "smtp" {
			readSmtpAccountStatusNotificationHandlerResponse(ctx, updateResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted" {
			readGroovyScriptedAccountStatusNotificationHandlerResponse(ctx, updateResponse.GroovyScriptedAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "admin-alert" {
			readAdminAlertAccountStatusNotificationHandlerResponse(ctx, updateResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "error-log" {
			readErrorLogAccountStatusNotificationHandlerResponse(ctx, updateResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "multi-part-email" {
			readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, updateResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party" {
			readThirdPartyAccountStatusNotificationHandlerResponse(ctx, updateResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *accountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state accountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandlerExecute(r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Account Status Notification Handler", err, httpResp)
		return
	}
}

func (r *accountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAccountStatusNotificationHandler(ctx, req, resp)
}

func (r *defaultAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAccountStatusNotificationHandler(ctx, req, resp)
}

func importAccountStatusNotificationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
