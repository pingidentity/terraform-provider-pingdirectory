package accountstatusnotificationhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &accountStatusNotificationHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &accountStatusNotificationHandlerDataSource{}
)

// Create a Account Status Notification Handler data source
func NewAccountStatusNotificationHandlerDataSource() datasource.DataSource {
	return &accountStatusNotificationHandlerDataSource{}
}

// accountStatusNotificationHandlerDataSource is the datasource implementation.
type accountStatusNotificationHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *accountStatusNotificationHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_account_status_notification_handler"
}

// Configure adds the provider configured client to the data source.
func (r *accountStatusNotificationHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type accountStatusNotificationHandlerDataSourceModel struct {
	Id                                              types.String `tfsdk:"id"`
	Name                                            types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *accountStatusNotificationHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Account Status Notification Handler resource. Options are ['smtp', 'groovy-scripted', 'admin-alert', 'error-log', 'multi-part-email', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Account Status Notification Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Account Status Notification Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"account_temporarily_failure_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an account becomes temporarily locked as a result of too many authentication failures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_permanently_failure_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an account becomes permanently locked as a result of too many authentication failures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_idle_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that authentication attempt fails because it has been too long since the user last successfully authenticated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_reset_locked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that authentication attempt fails because the user failed to choose a new password in a timely manner after an administrative reset.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_unlocked_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's account has been unlocked (e.g., by an administrative password reset).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_disabled_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's account is disabled by an administrator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_enabled_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's account is enabled by an administrator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_not_yet_active_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt fails because the account has an activation time that is in the future.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_expired_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt fails because the account has an expiration time that is in the past.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_expired_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt fails because the account has an expired password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_expiring_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an authentication attempt succeeds, but the user's password is about to expire. This notification will only be generated the first time the user authenticates within the window of time that the server should warn about an upcoming password expiration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_reset_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user's password has been reset by an administrator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_changed_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user changes their own password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_authenticated_message_template": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. The path to a file containing the template to use to generate the email message to send in the event that an account has successfully authenticated in a bind operation that matches the criteria provided in the account-authentication-notification-request-criteria property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_created_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a new account is created in an add request that matches the criteria provided in the account-creation-notification-request-criteria property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_deleted_message_template": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. The path to a file containing the template to use to generate the email message to send in the event that an existing accout has been removed in a delete request that matches the criteria provided in the account-deletion-notification-request-criteria property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_updated_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that an existing account is updated with a modify or modify DN operation that matches the criteria provided in the account-update-notification-request-criteria property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_password_failed_validation_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user authenticated with a password that failed to satisfy the criteria for one or more of the configured password validators.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"must_change_password_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a user successfully authenticates to the server but will be required to choose a new password before they will be allowed to perform any other operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_status_notification_type": schema.SetAttribute{
				Description: " When the `type` value is one of [`admin-alert`]: The types of account status notifications that should result in administrative alerts. When the `type` value is one of [`error-log`]: Indicates which types of event can trigger an account status notification.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Account Status Notification Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Account Status Notification Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_address_attribute_type": schema.SetAttribute{
				Description: "Specifies which attribute in the user's entries may be used to obtain the email address when notifying the end user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_address_json_field": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the email address to which the message should be sent. The email address must be contained in a top-level field whose value is a single string.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"email_address_json_object_filter": schema.StringAttribute{
				Description: "A JSON object filter that may be used to identify which email address value to use when sending the message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"recipient_address": schema.SetAttribute{
				Description: "Specifies an email address to which notification messages are sent, either instead of or in addition to the end user for whom the notification has been generated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"send_message_without_end_user_address": schema.BoolAttribute{
				Description: "Indicates whether an email notification message should be generated and sent to the set of notification recipients even if the user entry does not contain any values for any of the email address attributes (that is, in cases when it is not possible to notify the end user).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sender_address": schema.StringAttribute{
				Description: "Specifies the email address from which the message is sent. Note that this does not necessarily have to be a legitimate email address.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"message_subject": schema.SetAttribute{
				Description: "Specifies the subject that should be used for email messages generated by this account status notification handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"message_template_file": schema.SetAttribute{
				Description: "Specifies the path to the file containing the message template to generate the email notification messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Account Status Notification Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Account Status Notification Handler is enabled. Only enabled handlers are invoked whenever a related event occurs in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to invoke this Account Status Notification Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver a message) will not delay processing for the operation that triggered the notification.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_authentication_notification_result_criteria": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. A result criteria object that identifies which successful bind operations should result in account authentication notifications for this handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_creation_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which add requests should result in account creation notifications for this handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_deletion_notification_request_criteria": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. A request criteria object that identifies which delete requests should result in account deletion notifications for this handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_update_notification_request_criteria": schema.StringAttribute{
				Description: "A request criteria object that identifies which modify and modify DN requests should result in account update notifications for this handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SmtpAccountStatusNotificationHandlerResponse object into the model struct
func readSmtpAccountStatusNotificationHandlerResponseDataSource(ctx context.Context, r *client.SmtpAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EmailAddressAttributeType = internaltypes.GetStringSet(r.EmailAddressAttributeType)
	state.EmailAddressJSONField = internaltypes.StringTypeOrNil(r.EmailAddressJSONField, false)
	state.EmailAddressJSONObjectFilter = internaltypes.StringTypeOrNil(r.EmailAddressJSONObjectFilter, false)
	state.RecipientAddress = internaltypes.GetStringSet(r.RecipientAddress)
	state.SendMessageWithoutEndUserAddress = types.BoolValue(r.SendMessageWithoutEndUserAddress)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.MessageSubject = internaltypes.GetStringSet(r.MessageSubject)
	state.MessageTemplateFile = internaltypes.GetStringSet(r.MessageTemplateFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, false)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, false)
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, false)
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, false)
}

// Read a GroovyScriptedAccountStatusNotificationHandlerResponse object into the model struct
func readGroovyScriptedAccountStatusNotificationHandlerResponseDataSource(ctx context.Context, r *client.GroovyScriptedAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, false)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, false)
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, false)
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, false)
}

// Read a AdminAlertAccountStatusNotificationHandlerResponse object into the model struct
func readAdminAlertAccountStatusNotificationHandlerResponseDataSource(ctx context.Context, r *client.AdminAlertAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("admin-alert")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccountStatusNotificationType = internaltypes.GetStringSet(
		client.StringSliceEnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp(r.AccountStatusNotificationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, false)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, false)
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, false)
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, false)
}

// Read a ErrorLogAccountStatusNotificationHandlerResponse object into the model struct
func readErrorLogAccountStatusNotificationHandlerResponseDataSource(ctx context.Context, r *client.ErrorLogAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error-log")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccountStatusNotificationType = internaltypes.GetStringSet(
		client.StringSliceEnumaccountStatusNotificationHandlerAccountStatusNotificationTypeProp(r.AccountStatusNotificationType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, false)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, false)
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, false)
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, false)
}

// Read a MultiPartEmailAccountStatusNotificationHandlerResponse object into the model struct
func readMultiPartEmailAccountStatusNotificationHandlerResponseDataSource(ctx context.Context, r *client.MultiPartEmailAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("multi-part-email")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AccountTemporarilyFailureLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountTemporarilyFailureLockedMessageTemplate, false)
	state.AccountPermanentlyFailureLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountPermanentlyFailureLockedMessageTemplate, false)
	state.AccountIdleLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountIdleLockedMessageTemplate, false)
	state.AccountResetLockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountResetLockedMessageTemplate, false)
	state.AccountUnlockedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountUnlockedMessageTemplate, false)
	state.AccountDisabledMessageTemplate = internaltypes.StringTypeOrNil(r.AccountDisabledMessageTemplate, false)
	state.AccountEnabledMessageTemplate = internaltypes.StringTypeOrNil(r.AccountEnabledMessageTemplate, false)
	state.AccountNotYetActiveMessageTemplate = internaltypes.StringTypeOrNil(r.AccountNotYetActiveMessageTemplate, false)
	state.AccountExpiredMessageTemplate = internaltypes.StringTypeOrNil(r.AccountExpiredMessageTemplate, false)
	state.PasswordExpiredMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordExpiredMessageTemplate, false)
	state.PasswordExpiringMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordExpiringMessageTemplate, false)
	state.PasswordResetMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordResetMessageTemplate, false)
	state.PasswordChangedMessageTemplate = internaltypes.StringTypeOrNil(r.PasswordChangedMessageTemplate, false)
	state.AccountAuthenticatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountAuthenticatedMessageTemplate, false)
	state.AccountCreatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountCreatedMessageTemplate, false)
	state.AccountDeletedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountDeletedMessageTemplate, false)
	state.AccountUpdatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountUpdatedMessageTemplate, false)
	state.BindPasswordFailedValidationMessageTemplate = internaltypes.StringTypeOrNil(r.BindPasswordFailedValidationMessageTemplate, false)
	state.MustChangePasswordMessageTemplate = internaltypes.StringTypeOrNil(r.MustChangePasswordMessageTemplate, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, false)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, false)
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, false)
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, false)
}

// Read a ThirdPartyAccountStatusNotificationHandlerResponse object into the model struct
func readThirdPartyAccountStatusNotificationHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartyAccountStatusNotificationHandlerResponse, state *accountStatusNotificationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountAuthenticationNotificationResultCriteria = internaltypes.StringTypeOrNil(r.AccountAuthenticationNotificationResultCriteria, false)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, false)
	state.AccountDeletionNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountDeletionNotificationRequestCriteria, false)
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, false)
}

// Read resource information
func (r *accountStatusNotificationHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state accountStatusNotificationHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
		readSmtpAccountStatusNotificationHandlerResponseDataSource(ctx, readResponse.SmtpAccountStatusNotificationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAccountStatusNotificationHandlerResponse != nil {
		readGroovyScriptedAccountStatusNotificationHandlerResponseDataSource(ctx, readResponse.GroovyScriptedAccountStatusNotificationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AdminAlertAccountStatusNotificationHandlerResponse != nil {
		readAdminAlertAccountStatusNotificationHandlerResponseDataSource(ctx, readResponse.AdminAlertAccountStatusNotificationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogAccountStatusNotificationHandlerResponse != nil {
		readErrorLogAccountStatusNotificationHandlerResponseDataSource(ctx, readResponse.ErrorLogAccountStatusNotificationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse != nil {
		readMultiPartEmailAccountStatusNotificationHandlerResponseDataSource(ctx, readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccountStatusNotificationHandlerResponse != nil {
		readThirdPartyAccountStatusNotificationHandlerResponseDataSource(ctx, readResponse.ThirdPartyAccountStatusNotificationHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
