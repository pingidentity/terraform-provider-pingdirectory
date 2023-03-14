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
	_ resource.Resource                = &multiPartEmailAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &multiPartEmailAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &multiPartEmailAccountStatusNotificationHandlerResource{}
	_ resource.Resource                = &defaultMultiPartEmailAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultMultiPartEmailAccountStatusNotificationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultMultiPartEmailAccountStatusNotificationHandlerResource{}
)

// Create a Multi Part Email Account Status Notification Handler resource
func NewMultiPartEmailAccountStatusNotificationHandlerResource() resource.Resource {
	return &multiPartEmailAccountStatusNotificationHandlerResource{}
}

func NewDefaultMultiPartEmailAccountStatusNotificationHandlerResource() resource.Resource {
	return &defaultMultiPartEmailAccountStatusNotificationHandlerResource{}
}

// multiPartEmailAccountStatusNotificationHandlerResource is the resource implementation.
type multiPartEmailAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultMultiPartEmailAccountStatusNotificationHandlerResource is the resource implementation.
type defaultMultiPartEmailAccountStatusNotificationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *multiPartEmailAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_multi_part_email_account_status_notification_handler"
}

func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_multi_part_email_account_status_notification_handler"
}

// Configure adds the provider configured client to the resource.
func (r *multiPartEmailAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type multiPartEmailAccountStatusNotificationHandlerResourceModel struct {
	Id                                             types.String `tfsdk:"id"`
	LastUpdated                                    types.String `tfsdk:"last_updated"`
	Notifications                                  types.Set    `tfsdk:"notifications"`
	RequiredActions                                types.Set    `tfsdk:"required_actions"`
	AccountTemporarilyFailureLockedMessageTemplate types.String `tfsdk:"account_temporarily_failure_locked_message_template"`
	AccountPermanentlyFailureLockedMessageTemplate types.String `tfsdk:"account_permanently_failure_locked_message_template"`
	AccountIdleLockedMessageTemplate               types.String `tfsdk:"account_idle_locked_message_template"`
	AccountResetLockedMessageTemplate              types.String `tfsdk:"account_reset_locked_message_template"`
	AccountUnlockedMessageTemplate                 types.String `tfsdk:"account_unlocked_message_template"`
	AccountDisabledMessageTemplate                 types.String `tfsdk:"account_disabled_message_template"`
	AccountEnabledMessageTemplate                  types.String `tfsdk:"account_enabled_message_template"`
	AccountNotYetActiveMessageTemplate             types.String `tfsdk:"account_not_yet_active_message_template"`
	AccountExpiredMessageTemplate                  types.String `tfsdk:"account_expired_message_template"`
	PasswordExpiredMessageTemplate                 types.String `tfsdk:"password_expired_message_template"`
	PasswordExpiringMessageTemplate                types.String `tfsdk:"password_expiring_message_template"`
	PasswordResetMessageTemplate                   types.String `tfsdk:"password_reset_message_template"`
	PasswordChangedMessageTemplate                 types.String `tfsdk:"password_changed_message_template"`
	AccountCreatedMessageTemplate                  types.String `tfsdk:"account_created_message_template"`
	AccountUpdatedMessageTemplate                  types.String `tfsdk:"account_updated_message_template"`
	BindPasswordFailedValidationMessageTemplate    types.String `tfsdk:"bind_password_failed_validation_message_template"`
	MustChangePasswordMessageTemplate              types.String `tfsdk:"must_change_password_message_template"`
	Description                                    types.String `tfsdk:"description"`
	Enabled                                        types.Bool   `tfsdk:"enabled"`
	Asynchronous                                   types.Bool   `tfsdk:"asynchronous"`
	AccountCreationNotificationRequestCriteria     types.String `tfsdk:"account_creation_notification_request_criteria"`
	AccountUpdateNotificationRequestCriteria       types.String `tfsdk:"account_update_notification_request_criteria"`
}

// GetSchema defines the schema for the resource.
func (r *multiPartEmailAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	multiPartEmailAccountStatusNotificationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	multiPartEmailAccountStatusNotificationHandlerSchema(ctx, req, resp, true)
}

func multiPartEmailAccountStatusNotificationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Multi Part Email Account Status Notification Handler.",
		Attributes: map[string]schema.Attribute{
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
			"account_created_message_template": schema.StringAttribute{
				Description: "The path to a file containing the template to use to generate the email message to send in the event that a new account is created in an add request that matches the criteria provided in the account-creation-notification-request-criteria property.",
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
func addOptionalMultiPartEmailAccountStatusNotificationHandlerFields(ctx context.Context, addRequest *client.AddMultiPartEmailAccountStatusNotificationHandlerRequest, plan multiPartEmailAccountStatusNotificationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountTemporarilyFailureLockedMessageTemplate) {
		stringVal := plan.AccountTemporarilyFailureLockedMessageTemplate.ValueString()
		addRequest.AccountTemporarilyFailureLockedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountPermanentlyFailureLockedMessageTemplate) {
		stringVal := plan.AccountPermanentlyFailureLockedMessageTemplate.ValueString()
		addRequest.AccountPermanentlyFailureLockedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountIdleLockedMessageTemplate) {
		stringVal := plan.AccountIdleLockedMessageTemplate.ValueString()
		addRequest.AccountIdleLockedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountResetLockedMessageTemplate) {
		stringVal := plan.AccountResetLockedMessageTemplate.ValueString()
		addRequest.AccountResetLockedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUnlockedMessageTemplate) {
		stringVal := plan.AccountUnlockedMessageTemplate.ValueString()
		addRequest.AccountUnlockedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountDisabledMessageTemplate) {
		stringVal := plan.AccountDisabledMessageTemplate.ValueString()
		addRequest.AccountDisabledMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountEnabledMessageTemplate) {
		stringVal := plan.AccountEnabledMessageTemplate.ValueString()
		addRequest.AccountEnabledMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountNotYetActiveMessageTemplate) {
		stringVal := plan.AccountNotYetActiveMessageTemplate.ValueString()
		addRequest.AccountNotYetActiveMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountExpiredMessageTemplate) {
		stringVal := plan.AccountExpiredMessageTemplate.ValueString()
		addRequest.AccountExpiredMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordExpiredMessageTemplate) {
		stringVal := plan.PasswordExpiredMessageTemplate.ValueString()
		addRequest.PasswordExpiredMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordExpiringMessageTemplate) {
		stringVal := plan.PasswordExpiringMessageTemplate.ValueString()
		addRequest.PasswordExpiringMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordResetMessageTemplate) {
		stringVal := plan.PasswordResetMessageTemplate.ValueString()
		addRequest.PasswordResetMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordChangedMessageTemplate) {
		stringVal := plan.PasswordChangedMessageTemplate.ValueString()
		addRequest.PasswordChangedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountCreatedMessageTemplate) {
		stringVal := plan.AccountCreatedMessageTemplate.ValueString()
		addRequest.AccountCreatedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccountUpdatedMessageTemplate) {
		stringVal := plan.AccountUpdatedMessageTemplate.ValueString()
		addRequest.AccountUpdatedMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindPasswordFailedValidationMessageTemplate) {
		stringVal := plan.BindPasswordFailedValidationMessageTemplate.ValueString()
		addRequest.BindPasswordFailedValidationMessageTemplate = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MustChangePasswordMessageTemplate) {
		stringVal := plan.MustChangePasswordMessageTemplate.ValueString()
		addRequest.MustChangePasswordMessageTemplate = &stringVal
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

// Read a MultiPartEmailAccountStatusNotificationHandlerResponse object into the model struct
func readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx context.Context, r *client.MultiPartEmailAccountStatusNotificationHandlerResponse, state *multiPartEmailAccountStatusNotificationHandlerResourceModel, expectedValues *multiPartEmailAccountStatusNotificationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
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
	state.AccountCreatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountCreatedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountCreatedMessageTemplate))
	state.AccountUpdatedMessageTemplate = internaltypes.StringTypeOrNil(r.AccountUpdatedMessageTemplate, internaltypes.IsEmptyString(expectedValues.AccountUpdatedMessageTemplate))
	state.BindPasswordFailedValidationMessageTemplate = internaltypes.StringTypeOrNil(r.BindPasswordFailedValidationMessageTemplate, internaltypes.IsEmptyString(expectedValues.BindPasswordFailedValidationMessageTemplate))
	state.MustChangePasswordMessageTemplate = internaltypes.StringTypeOrNil(r.MustChangePasswordMessageTemplate, internaltypes.IsEmptyString(expectedValues.MustChangePasswordMessageTemplate))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.AccountCreationNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountCreationNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountCreationNotificationRequestCriteria))
	state.AccountUpdateNotificationRequestCriteria = internaltypes.StringTypeOrNil(r.AccountUpdateNotificationRequestCriteria, internaltypes.IsEmptyString(expectedValues.AccountUpdateNotificationRequestCriteria))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createMultiPartEmailAccountStatusNotificationHandlerOperations(plan multiPartEmailAccountStatusNotificationHandlerResourceModel, state multiPartEmailAccountStatusNotificationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
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
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreatedMessageTemplate, state.AccountCreatedMessageTemplate, "account-created-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdatedMessageTemplate, state.AccountUpdatedMessageTemplate, "account-updated-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.BindPasswordFailedValidationMessageTemplate, state.BindPasswordFailedValidationMessageTemplate, "bind-password-failed-validation-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.MustChangePasswordMessageTemplate, state.MustChangePasswordMessageTemplate, "must-change-password-message-template")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountCreationNotificationRequestCriteria, state.AccountCreationNotificationRequestCriteria, "account-creation-notification-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AccountUpdateNotificationRequestCriteria, state.AccountUpdateNotificationRequestCriteria, "account-update-notification-request-criteria")
	return ops
}

// Create a new resource
func (r *multiPartEmailAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan multiPartEmailAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddMultiPartEmailAccountStatusNotificationHandlerRequest(plan.Id.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Multi Part Email Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state multiPartEmailAccountStatusNotificationHandlerResourceModel
	readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, addResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan multiPartEmailAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Multi Part Email Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state multiPartEmailAccountStatusNotificationHandlerResourceModel
	readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createMultiPartEmailAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Multi Part Email Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, updateResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *multiPartEmailAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readMultiPartEmailAccountStatusNotificationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state multiPartEmailAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.GetAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Multi Part Email Account Status Notification Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, readResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *multiPartEmailAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateMultiPartEmailAccountStatusNotificationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan multiPartEmailAccountStatusNotificationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state multiPartEmailAccountStatusNotificationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createMultiPartEmailAccountStatusNotificationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccountStatusNotificationHandlerApi.UpdateAccountStatusNotificationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Multi Part Email Account Status Notification Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readMultiPartEmailAccountStatusNotificationHandlerResponse(ctx, updateResponse.MultiPartEmailAccountStatusNotificationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *multiPartEmailAccountStatusNotificationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state multiPartEmailAccountStatusNotificationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandlerExecute(r.apiClient.AccountStatusNotificationHandlerApi.DeleteAccountStatusNotificationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Multi Part Email Account Status Notification Handler", err, httpResp)
		return
	}
}

func (r *multiPartEmailAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp)
}

func (r *defaultMultiPartEmailAccountStatusNotificationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importMultiPartEmailAccountStatusNotificationHandler(ctx, req, resp)
}

func importMultiPartEmailAccountStatusNotificationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
