package passwordpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordPolicyResource{}
	_ resource.ResourceWithConfigure   = &passwordPolicyResource{}
	_ resource.ResourceWithImportState = &passwordPolicyResource{}
	_ resource.Resource                = &defaultPasswordPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultPasswordPolicyResource{}
	_ resource.ResourceWithImportState = &defaultPasswordPolicyResource{}
)

// Create a Password Policy resource
func NewPasswordPolicyResource() resource.Resource {
	return &passwordPolicyResource{}
}

func NewDefaultPasswordPolicyResource() resource.Resource {
	return &defaultPasswordPolicyResource{}
}

// passwordPolicyResource is the resource implementation.
type passwordPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPasswordPolicyResource is the resource implementation.
type defaultPasswordPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passwordPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_policy"
}

func (r *defaultPasswordPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_password_policy"
}

// Configure adds the provider configured client to the resource.
func (r *passwordPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPasswordPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordPolicyResourceModel struct {
	Id                                                        types.String `tfsdk:"id"`
	Name                                                      types.String `tfsdk:"name"`
	Notifications                                             types.Set    `tfsdk:"notifications"`
	RequiredActions                                           types.Set    `tfsdk:"required_actions"`
	Type                                                      types.String `tfsdk:"type"`
	Description                                               types.String `tfsdk:"description"`
	RequireSecureAuthentication                               types.Bool   `tfsdk:"require_secure_authentication"`
	RequireSecurePasswordChanges                              types.Bool   `tfsdk:"require_secure_password_changes"`
	AccountStatusNotificationHandler                          types.Set    `tfsdk:"account_status_notification_handler"`
	StateUpdateFailurePolicy                                  types.String `tfsdk:"state_update_failure_policy"`
	EnableDebug                                               types.Bool   `tfsdk:"enable_debug"`
	PasswordAttribute                                         types.String `tfsdk:"password_attribute"`
	DefaultPasswordStorageScheme                              types.Set    `tfsdk:"default_password_storage_scheme"`
	DeprecatedPasswordStorageScheme                           types.Set    `tfsdk:"deprecated_password_storage_scheme"`
	AllowMultiplePasswordValues                               types.Bool   `tfsdk:"allow_multiple_password_values"`
	AllowPreEncodedPasswords                                  types.String `tfsdk:"allow_pre_encoded_passwords"`
	PasswordValidator                                         types.Set    `tfsdk:"password_validator"`
	BindPasswordValidator                                     types.Set    `tfsdk:"bind_password_validator"`
	MinimumBindPasswordValidationFrequency                    types.String `tfsdk:"minimum_bind_password_validation_frequency"`
	BindPasswordValidationFailureAction                       types.String `tfsdk:"bind_password_validation_failure_action"`
	PasswordGenerator                                         types.String `tfsdk:"password_generator"`
	PasswordHistoryCount                                      types.Int64  `tfsdk:"password_history_count"`
	PasswordHistoryDuration                                   types.String `tfsdk:"password_history_duration"`
	MinPasswordAge                                            types.String `tfsdk:"min_password_age"`
	MaxPasswordAge                                            types.String `tfsdk:"max_password_age"`
	PasswordExpirationWarningInterval                         types.String `tfsdk:"password_expiration_warning_interval"`
	ExpirePasswordsWithoutWarning                             types.Bool   `tfsdk:"expire_passwords_without_warning"`
	ReturnPasswordExpirationControls                          types.String `tfsdk:"return_password_expiration_controls"`
	AllowExpiredPasswordChanges                               types.Bool   `tfsdk:"allow_expired_password_changes"`
	GraceLoginCount                                           types.Int64  `tfsdk:"grace_login_count"`
	RequireChangeByTime                                       types.String `tfsdk:"require_change_by_time"`
	LockoutFailureCount                                       types.Int64  `tfsdk:"lockout_failure_count"`
	LockoutDuration                                           types.String `tfsdk:"lockout_duration"`
	LockoutFailureExpirationInterval                          types.String `tfsdk:"lockout_failure_expiration_interval"`
	IgnoreDuplicatePasswordFailures                           types.Bool   `tfsdk:"ignore_duplicate_password_failures"`
	FailureLockoutAction                                      types.String `tfsdk:"failure_lockout_action"`
	IdleLockoutInterval                                       types.String `tfsdk:"idle_lockout_interval"`
	AllowUserPasswordChanges                                  types.Bool   `tfsdk:"allow_user_password_changes"`
	PasswordChangeRequiresCurrentPassword                     types.Bool   `tfsdk:"password_change_requires_current_password"`
	PasswordRetirementBehavior                                types.Set    `tfsdk:"password_retirement_behavior"`
	MaxRetiredPasswordAge                                     types.String `tfsdk:"max_retired_password_age"`
	AllowedPasswordResetTokenUseCondition                     types.Set    `tfsdk:"allowed_password_reset_token_use_condition"`
	ForceChangeOnAdd                                          types.Bool   `tfsdk:"force_change_on_add"`
	ForceChangeOnReset                                        types.Bool   `tfsdk:"force_change_on_reset"`
	MaxPasswordResetAge                                       types.String `tfsdk:"max_password_reset_age"`
	SkipValidationForAdministrators                           types.Bool   `tfsdk:"skip_validation_for_administrators"`
	MaximumRecentLoginHistorySuccessfulAuthenticationCount    types.Int64  `tfsdk:"maximum_recent_login_history_successful_authentication_count"`
	MaximumRecentLoginHistorySuccessfulAuthenticationDuration types.String `tfsdk:"maximum_recent_login_history_successful_authentication_duration"`
	MaximumRecentLoginHistoryFailedAuthenticationCount        types.Int64  `tfsdk:"maximum_recent_login_history_failed_authentication_count"`
	MaximumRecentLoginHistoryFailedAuthenticationDuration     types.String `tfsdk:"maximum_recent_login_history_failed_authentication_duration"`
	RecentLoginHistorySimilarAttemptBehavior                  types.String `tfsdk:"recent_login_history_similar_attempt_behavior"`
	LastLoginIPAddressAttribute                               types.String `tfsdk:"last_login_ip_address_attribute"`
	LastLoginTimeAttribute                                    types.String `tfsdk:"last_login_time_attribute"`
	LastLoginTimeFormat                                       types.String `tfsdk:"last_login_time_format"`
	PreviousLastLoginTimeFormat                               types.Set    `tfsdk:"previous_last_login_time_format"`
}

// GetSchema defines the schema for the resource.
func (r *passwordPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordPolicySchema(ctx, req, resp, false)
}

func (r *defaultPasswordPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordPolicySchema(ctx, req, resp, true)
}

func passwordPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Password Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Policy resource. Options are ['password-policy']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("password-policy"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"password-policy"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Policy",
				Optional:    true,
			},
			"require_secure_authentication": schema.BoolAttribute{
				Description: "Indicates whether users with the associated password policy are required to authenticate in a secure manner.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"require_secure_password_changes": schema.BoolAttribute{
				Description: "Indicates whether users with the associated password policy are required to change their password in a secure manner that does not expose the credentials.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"account_status_notification_handler": schema.SetAttribute{
				Description: "Specifies the names of the account status notification handlers that are used with the associated password storage scheme.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"state_update_failure_policy": schema.StringAttribute{
				Description: "Specifies how the server deals with the inability to update password policy state information during an authentication attempt.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("reactive"),
			},
			"enable_debug": schema.BoolAttribute{
				Description: "Indicates whether to enable debugging for the password policy state.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"password_attribute": schema.StringAttribute{
				Description: "Specifies the attribute type used to hold user passwords.",
				Required:    true,
			},
			"default_password_storage_scheme": schema.SetAttribute{
				Description: "Specifies the names of the password storage schemes that are used to encode clear-text passwords for this password policy.",
				Required:    true,
				ElementType: types.StringType,
			},
			"deprecated_password_storage_scheme": schema.SetAttribute{
				Description: "Specifies the names of the password storage schemes that are considered deprecated for this password policy.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"allow_multiple_password_values": schema.BoolAttribute{
				Description: "Indicates whether user entries can have multiple distinct values for the password attribute.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"allow_pre_encoded_passwords": schema.StringAttribute{
				Description: "Indicates whether users can change their passwords by providing a pre-encoded value.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("false"),
			},
			"password_validator": schema.SetAttribute{
				Description: "Specifies the names of the password validators that are used with the associated password storage scheme.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"bind_password_validator": schema.SetAttribute{
				Description: "Specifies the names of the password validators that should be invoked for bind operations.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"minimum_bind_password_validation_frequency": schema.StringAttribute{
				Description: "Indicates how frequently password validation should be performed during bind operations for each user to whom this password policy is assigned.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"bind_password_validation_failure_action": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if a bind password fails validation by one or more of the configured bind password validators.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("force-password-change"),
			},
			"password_generator": schema.StringAttribute{
				Description: "Specifies the name of the password generator that is used with the associated password policy.",
				Optional:    true,
			},
			"password_history_count": schema.Int64Attribute{
				Description: "Specifies the maximum number of former passwords to maintain in the password history.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"password_history_duration": schema.StringAttribute{
				Description: "Specifies the maximum length of time that passwords remain in the password history.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"min_password_age": schema.StringAttribute{
				Description: "Specifies the minimum length of time after a password change before the user is allowed to change the password again.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_password_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a user can continue using the same password before it must be changed (that is, the password expiration interval).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_expiration_warning_interval": schema.StringAttribute{
				Description: "Specifies the maximum length of time before a user's password actually expires that the server begins to include warning notifications in bind responses for that user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"expire_passwords_without_warning": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server allows a user's password to expire even if that user has never seen an expiration warning notification.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"return_password_expiration_controls": schema.StringAttribute{
				Description: "Indicates whether the server should return the password expiring and password expired response controls (as described in draft-vchu-ldap-pwd-policy).",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("unless-password-policy-control-is-used"),
			},
			"allow_expired_password_changes": schema.BoolAttribute{
				Description: "Indicates whether a user whose password is expired is still allowed to change that password using the password modify extended operation.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"grace_login_count": schema.Int64Attribute{
				Description: "Specifies the number of grace logins that a user is allowed after the account has expired to allow that user to choose a new password.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"require_change_by_time": schema.StringAttribute{
				Description: "Specifies the time by which all users with the associated password policy must change their passwords.",
				Optional:    true,
			},
			"lockout_failure_count": schema.Int64Attribute{
				Description: "Specifies the maximum number of authentication failures that a user is allowed before the account is locked out.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"lockout_duration": schema.StringAttribute{
				Description: "Specifies the length of time that an account is locked after too many authentication failures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"lockout_failure_expiration_interval": schema.StringAttribute{
				Description: "Specifies the length of time before an authentication failure is no longer counted against a user for the purposes of account lockout.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ignore_duplicate_password_failures": schema.BoolAttribute{
				Description: "Indicates whether to ignore subsequent authentication failures using the same password as an earlier failed authentication attempt (within the time frame defined by the lockout failure expiration interval). If this option is \"true\", then multiple failed attempts using the same password will be considered only a single failure. If this option is \"false\", then any failure will be tracked regardless of whether it used the same password as an earlier attempt.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"failure_lockout_action": schema.StringAttribute{
				Description: "The action that the server should take for authentication attempts that target a user with more than the configured number of outstanding authentication failures.",
				Optional:    true,
			},
			"idle_lockout_interval": schema.StringAttribute{
				Description: "Specifies the maximum length of time that an account may remain idle (that is, the associated user does not authenticate to the server) before that user is locked out.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_user_password_changes": schema.BoolAttribute{
				Description: "Indicates whether users can change their own passwords.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
			},
			"password_change_requires_current_password": schema.BoolAttribute{
				Description: "Indicates whether user password changes must use the password modify extended operation and must include the user's current password before the change is allowed.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"password_retirement_behavior": schema.SetAttribute{
				Description: "Specifies the conditions under which the server may retire a user's current password in the course of setting a new password for that user (whether via a modify operation or a password modify extended operation).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"max_retired_password_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a retired password should be considered valid and may be used to authenticate to the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_password_reset_token_use_condition": schema.SetAttribute{
				Description: "The set of conditions under which a user governed by this Password Policy will be permitted to generate a password reset token via the deliver password reset token extended operation, and to use that token in lieu of the current password via the password modify extended operation.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"force_change_on_add": schema.BoolAttribute{
				Description: "Indicates whether users are forced to change their passwords upon first authenticating to the Directory Server after their account has been created.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"force_change_on_reset": schema.BoolAttribute{
				Description: "Indicates whether users are forced to change their passwords if they are reset by an administrator. If a user's password is changed by any other user, that is considered an administrative password reset.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"max_password_reset_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that users have to change passwords after they have been reset by an administrator before they become locked.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"skip_validation_for_administrators": schema.BoolAttribute{
				Description: "Indicates whether passwords set by administrators are allowed to bypass the password validation process that is required for user password changes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"maximum_recent_login_history_successful_authentication_count": schema.Int64Attribute{
				Description: "The maximum number of successful authentication attempts to include in the recent login history for each account.",
				Optional:    true,
			},
			"maximum_recent_login_history_successful_authentication_duration": schema.StringAttribute{
				Description: "The maximum age of successful authentication attempts to include in the recent login history for each account.",
				Optional:    true,
			},
			"maximum_recent_login_history_failed_authentication_count": schema.Int64Attribute{
				Description: "The maximum number of failed authentication attempts to include in the recent login history for each account.",
				Optional:    true,
			},
			"maximum_recent_login_history_failed_authentication_duration": schema.StringAttribute{
				Description: "The maximum age of failed authentication attempts to include in the recent login history for each account.",
				Optional:    true,
			},
			"recent_login_history_similar_attempt_behavior": schema.StringAttribute{
				Description: "The behavior that the server will exhibit when multiple similar authentication attempts (with the same values for the successful, authentication-method, client-ip-address, and failure-reason fields) are processed for an account.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("collapse-similar-attempts-on-the-same-date"),
			},
			"last_login_ip_address_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute type that is used to hold the IP address of the client from which the user last authenticated.",
				Optional:    true,
			},
			"last_login_time_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute type that is used to hold the last login time for users with the associated password policy.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("ds-pwp-last-login-time"),
			},
			"last_login_time_format": schema.StringAttribute{
				Description: "Specifies the format string that is used to generate the last login time value for users with the associated password policy. Last login time values will be written using the UTC (also known as GMT, or Greenwich Mean Time) time zone.",
				Optional:    true,
			},
			"previous_last_login_time_format": schema.SetAttribute{
				Description: "Specifies the format string(s) that might have been used with the last login time at any point in the past for users associated with the password policy.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for password-policy password-policy
func addOptionalPasswordPolicyFields(ctx context.Context, addRequest *client.AddPasswordPolicyRequest, plan passwordPolicyResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RequireSecureAuthentication) {
		addRequest.RequireSecureAuthentication = plan.RequireSecureAuthentication.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.RequireSecurePasswordChanges) {
		addRequest.RequireSecurePasswordChanges = plan.RequireSecurePasswordChanges.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AccountStatusNotificationHandler) {
		var slice []string
		plan.AccountStatusNotificationHandler.ElementsAs(ctx, &slice, false)
		addRequest.AccountStatusNotificationHandler = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.StateUpdateFailurePolicy) {
		stateUpdateFailurePolicy, err := client.NewEnumpasswordPolicyStateUpdateFailurePolicyPropFromValue(plan.StateUpdateFailurePolicy.ValueString())
		if err != nil {
			return err
		}
		addRequest.StateUpdateFailurePolicy = stateUpdateFailurePolicy
	}
	if internaltypes.IsDefined(plan.EnableDebug) {
		addRequest.EnableDebug = plan.EnableDebug.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DeprecatedPasswordStorageScheme) {
		var slice []string
		plan.DeprecatedPasswordStorageScheme.ElementsAs(ctx, &slice, false)
		addRequest.DeprecatedPasswordStorageScheme = slice
	}
	if internaltypes.IsDefined(plan.AllowMultiplePasswordValues) {
		addRequest.AllowMultiplePasswordValues = plan.AllowMultiplePasswordValues.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AllowPreEncodedPasswords) {
		allowPreEncodedPasswords, err := client.NewEnumpasswordPolicyAllowPreEncodedPasswordsPropFromValue(plan.AllowPreEncodedPasswords.ValueString())
		if err != nil {
			return err
		}
		addRequest.AllowPreEncodedPasswords = allowPreEncodedPasswords
	}
	if internaltypes.IsDefined(plan.PasswordValidator) {
		var slice []string
		plan.PasswordValidator.ElementsAs(ctx, &slice, false)
		addRequest.PasswordValidator = slice
	}
	if internaltypes.IsDefined(plan.BindPasswordValidator) {
		var slice []string
		plan.BindPasswordValidator.ElementsAs(ctx, &slice, false)
		addRequest.BindPasswordValidator = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinimumBindPasswordValidationFrequency) {
		addRequest.MinimumBindPasswordValidationFrequency = plan.MinimumBindPasswordValidationFrequency.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindPasswordValidationFailureAction) {
		bindPasswordValidationFailureAction, err := client.NewEnumpasswordPolicyBindPasswordValidationFailureActionPropFromValue(plan.BindPasswordValidationFailureAction.ValueString())
		if err != nil {
			return err
		}
		addRequest.BindPasswordValidationFailureAction = bindPasswordValidationFailureAction
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordGenerator) {
		addRequest.PasswordGenerator = plan.PasswordGenerator.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PasswordHistoryCount) {
		addRequest.PasswordHistoryCount = plan.PasswordHistoryCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordHistoryDuration) {
		addRequest.PasswordHistoryDuration = plan.PasswordHistoryDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinPasswordAge) {
		addRequest.MinPasswordAge = plan.MinPasswordAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxPasswordAge) {
		addRequest.MaxPasswordAge = plan.MaxPasswordAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordExpirationWarningInterval) {
		addRequest.PasswordExpirationWarningInterval = plan.PasswordExpirationWarningInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ExpirePasswordsWithoutWarning) {
		addRequest.ExpirePasswordsWithoutWarning = plan.ExpirePasswordsWithoutWarning.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ReturnPasswordExpirationControls) {
		returnPasswordExpirationControls, err := client.NewEnumpasswordPolicyReturnPasswordExpirationControlsPropFromValue(plan.ReturnPasswordExpirationControls.ValueString())
		if err != nil {
			return err
		}
		addRequest.ReturnPasswordExpirationControls = returnPasswordExpirationControls
	}
	if internaltypes.IsDefined(plan.AllowExpiredPasswordChanges) {
		addRequest.AllowExpiredPasswordChanges = plan.AllowExpiredPasswordChanges.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GraceLoginCount) {
		addRequest.GraceLoginCount = plan.GraceLoginCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequireChangeByTime) {
		addRequest.RequireChangeByTime = plan.RequireChangeByTime.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LockoutFailureCount) {
		addRequest.LockoutFailureCount = plan.LockoutFailureCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LockoutDuration) {
		addRequest.LockoutDuration = plan.LockoutDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LockoutFailureExpirationInterval) {
		addRequest.LockoutFailureExpirationInterval = plan.LockoutFailureExpirationInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IgnoreDuplicatePasswordFailures) {
		addRequest.IgnoreDuplicatePasswordFailures = plan.IgnoreDuplicatePasswordFailures.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.FailureLockoutAction) {
		addRequest.FailureLockoutAction = plan.FailureLockoutAction.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdleLockoutInterval) {
		addRequest.IdleLockoutInterval = plan.IdleLockoutInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowUserPasswordChanges) {
		addRequest.AllowUserPasswordChanges = plan.AllowUserPasswordChanges.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.PasswordChangeRequiresCurrentPassword) {
		addRequest.PasswordChangeRequiresCurrentPassword = plan.PasswordChangeRequiresCurrentPassword.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.PasswordRetirementBehavior) {
		var slice []string
		plan.PasswordRetirementBehavior.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpasswordPolicyPasswordRetirementBehaviorProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpasswordPolicyPasswordRetirementBehaviorPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.PasswordRetirementBehavior = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxRetiredPasswordAge) {
		addRequest.MaxRetiredPasswordAge = plan.MaxRetiredPasswordAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowedPasswordResetTokenUseCondition) {
		var slice []string
		plan.AllowedPasswordResetTokenUseCondition.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumpasswordPolicyAllowedPasswordResetTokenUseConditionProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumpasswordPolicyAllowedPasswordResetTokenUseConditionPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedPasswordResetTokenUseCondition = enumSlice
	}
	if internaltypes.IsDefined(plan.ForceChangeOnAdd) {
		addRequest.ForceChangeOnAdd = plan.ForceChangeOnAdd.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ForceChangeOnReset) {
		addRequest.ForceChangeOnReset = plan.ForceChangeOnReset.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxPasswordResetAge) {
		addRequest.MaxPasswordResetAge = plan.MaxPasswordResetAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SkipValidationForAdministrators) {
		addRequest.SkipValidationForAdministrators = plan.SkipValidationForAdministrators.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MaximumRecentLoginHistorySuccessfulAuthenticationCount) {
		addRequest.MaximumRecentLoginHistorySuccessfulAuthenticationCount = plan.MaximumRecentLoginHistorySuccessfulAuthenticationCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumRecentLoginHistorySuccessfulAuthenticationDuration) {
		addRequest.MaximumRecentLoginHistorySuccessfulAuthenticationDuration = plan.MaximumRecentLoginHistorySuccessfulAuthenticationDuration.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaximumRecentLoginHistoryFailedAuthenticationCount) {
		addRequest.MaximumRecentLoginHistoryFailedAuthenticationCount = plan.MaximumRecentLoginHistoryFailedAuthenticationCount.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaximumRecentLoginHistoryFailedAuthenticationDuration) {
		addRequest.MaximumRecentLoginHistoryFailedAuthenticationDuration = plan.MaximumRecentLoginHistoryFailedAuthenticationDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RecentLoginHistorySimilarAttemptBehavior) {
		recentLoginHistorySimilarAttemptBehavior, err := client.NewEnumpasswordPolicyRecentLoginHistorySimilarAttemptBehaviorPropFromValue(plan.RecentLoginHistorySimilarAttemptBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.RecentLoginHistorySimilarAttemptBehavior = recentLoginHistorySimilarAttemptBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LastLoginIPAddressAttribute) {
		addRequest.LastLoginIPAddressAttribute = plan.LastLoginIPAddressAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LastLoginTimeAttribute) {
		addRequest.LastLoginTimeAttribute = plan.LastLoginTimeAttribute.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LastLoginTimeFormat) {
		addRequest.LastLoginTimeFormat = plan.LastLoginTimeFormat.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.PreviousLastLoginTimeFormat) {
		var slice []string
		plan.PreviousLastLoginTimeFormat.ElementsAs(ctx, &slice, false)
		addRequest.PreviousLastLoginTimeFormat = slice
	}
	return nil
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *passwordPolicyResourceModel) populateAllComputedStringAttributes() {
	if model.StateUpdateFailurePolicy.IsUnknown() || model.StateUpdateFailurePolicy.IsNull() {
		model.StateUpdateFailurePolicy = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.LockoutDuration.IsUnknown() || model.LockoutDuration.IsNull() {
		model.LockoutDuration = types.StringValue("")
	}
	if model.RequireChangeByTime.IsUnknown() || model.RequireChangeByTime.IsNull() {
		model.RequireChangeByTime = types.StringValue("")
	}
	if model.MaximumRecentLoginHistoryFailedAuthenticationDuration.IsUnknown() || model.MaximumRecentLoginHistoryFailedAuthenticationDuration.IsNull() {
		model.MaximumRecentLoginHistoryFailedAuthenticationDuration = types.StringValue("")
	}
	if model.PasswordHistoryDuration.IsUnknown() || model.PasswordHistoryDuration.IsNull() {
		model.PasswordHistoryDuration = types.StringValue("")
	}
	if model.BindPasswordValidationFailureAction.IsUnknown() || model.BindPasswordValidationFailureAction.IsNull() {
		model.BindPasswordValidationFailureAction = types.StringValue("")
	}
	if model.IdleLockoutInterval.IsUnknown() || model.IdleLockoutInterval.IsNull() {
		model.IdleLockoutInterval = types.StringValue("")
	}
	if model.LastLoginTimeAttribute.IsUnknown() || model.LastLoginTimeAttribute.IsNull() {
		model.LastLoginTimeAttribute = types.StringValue("")
	}
	if model.MaxPasswordResetAge.IsUnknown() || model.MaxPasswordResetAge.IsNull() {
		model.MaxPasswordResetAge = types.StringValue("")
	}
	if model.MinPasswordAge.IsUnknown() || model.MinPasswordAge.IsNull() {
		model.MinPasswordAge = types.StringValue("")
	}
	if model.PasswordGenerator.IsUnknown() || model.PasswordGenerator.IsNull() {
		model.PasswordGenerator = types.StringValue("")
	}
	if model.LastLoginTimeFormat.IsUnknown() || model.LastLoginTimeFormat.IsNull() {
		model.LastLoginTimeFormat = types.StringValue("")
	}
	if model.MinimumBindPasswordValidationFrequency.IsUnknown() || model.MinimumBindPasswordValidationFrequency.IsNull() {
		model.MinimumBindPasswordValidationFrequency = types.StringValue("")
	}
	if model.PasswordAttribute.IsUnknown() || model.PasswordAttribute.IsNull() {
		model.PasswordAttribute = types.StringValue("")
	}
	if model.MaxRetiredPasswordAge.IsUnknown() || model.MaxRetiredPasswordAge.IsNull() {
		model.MaxRetiredPasswordAge = types.StringValue("")
	}
	if model.PasswordExpirationWarningInterval.IsUnknown() || model.PasswordExpirationWarningInterval.IsNull() {
		model.PasswordExpirationWarningInterval = types.StringValue("")
	}
	if model.MaximumRecentLoginHistorySuccessfulAuthenticationDuration.IsUnknown() || model.MaximumRecentLoginHistorySuccessfulAuthenticationDuration.IsNull() {
		model.MaximumRecentLoginHistorySuccessfulAuthenticationDuration = types.StringValue("")
	}
	if model.RecentLoginHistorySimilarAttemptBehavior.IsUnknown() || model.RecentLoginHistorySimilarAttemptBehavior.IsNull() {
		model.RecentLoginHistorySimilarAttemptBehavior = types.StringValue("")
	}
	if model.FailureLockoutAction.IsUnknown() || model.FailureLockoutAction.IsNull() {
		model.FailureLockoutAction = types.StringValue("")
	}
	if model.MaxPasswordAge.IsUnknown() || model.MaxPasswordAge.IsNull() {
		model.MaxPasswordAge = types.StringValue("")
	}
	if model.LastLoginIPAddressAttribute.IsUnknown() || model.LastLoginIPAddressAttribute.IsNull() {
		model.LastLoginIPAddressAttribute = types.StringValue("")
	}
	if model.AllowPreEncodedPasswords.IsUnknown() || model.AllowPreEncodedPasswords.IsNull() {
		model.AllowPreEncodedPasswords = types.StringValue("")
	}
	if model.ReturnPasswordExpirationControls.IsUnknown() || model.ReturnPasswordExpirationControls.IsNull() {
		model.ReturnPasswordExpirationControls = types.StringValue("")
	}
	if model.LockoutFailureExpirationInterval.IsUnknown() || model.LockoutFailureExpirationInterval.IsNull() {
		model.LockoutFailureExpirationInterval = types.StringValue("")
	}
}

// Read a PasswordPolicyResponse object into the model struct
func readPasswordPolicyResponse(ctx context.Context, r *client.PasswordPolicyResponse, state *passwordPolicyResourceModel, expectedValues *passwordPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.RequireSecureAuthentication = internaltypes.BoolTypeOrNil(r.RequireSecureAuthentication)
	state.RequireSecurePasswordChanges = internaltypes.BoolTypeOrNil(r.RequireSecurePasswordChanges)
	state.AccountStatusNotificationHandler = internaltypes.GetStringSet(r.AccountStatusNotificationHandler)
	state.StateUpdateFailurePolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyStateUpdateFailurePolicyProp(r.StateUpdateFailurePolicy), true)
	state.EnableDebug = internaltypes.BoolTypeOrNil(r.EnableDebug)
	state.PasswordAttribute = types.StringValue(r.PasswordAttribute)
	state.DefaultPasswordStorageScheme = internaltypes.GetStringSet(r.DefaultPasswordStorageScheme)
	state.DeprecatedPasswordStorageScheme = internaltypes.GetStringSet(r.DeprecatedPasswordStorageScheme)
	state.AllowMultiplePasswordValues = internaltypes.BoolTypeOrNil(r.AllowMultiplePasswordValues)
	state.AllowPreEncodedPasswords = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyAllowPreEncodedPasswordsProp(r.AllowPreEncodedPasswords), true)
	state.PasswordValidator = internaltypes.GetStringSet(r.PasswordValidator)
	state.BindPasswordValidator = internaltypes.GetStringSet(r.BindPasswordValidator)
	state.MinimumBindPasswordValidationFrequency = internaltypes.StringTypeOrNil(r.MinimumBindPasswordValidationFrequency, true)
	config.CheckMismatchedPDFormattedAttributes("minimum_bind_password_validation_frequency",
		expectedValues.MinimumBindPasswordValidationFrequency, state.MinimumBindPasswordValidationFrequency, diagnostics)
	state.BindPasswordValidationFailureAction = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyBindPasswordValidationFailureActionProp(r.BindPasswordValidationFailureAction), true)
	state.PasswordGenerator = internaltypes.StringTypeOrNil(r.PasswordGenerator, internaltypes.IsEmptyString(expectedValues.PasswordGenerator))
	state.PasswordHistoryCount = internaltypes.Int64TypeOrNil(r.PasswordHistoryCount)
	state.PasswordHistoryDuration = internaltypes.StringTypeOrNil(r.PasswordHistoryDuration, true)
	config.CheckMismatchedPDFormattedAttributes("password_history_duration",
		expectedValues.PasswordHistoryDuration, state.PasswordHistoryDuration, diagnostics)
	state.MinPasswordAge = internaltypes.StringTypeOrNil(r.MinPasswordAge, true)
	config.CheckMismatchedPDFormattedAttributes("min_password_age",
		expectedValues.MinPasswordAge, state.MinPasswordAge, diagnostics)
	state.MaxPasswordAge = internaltypes.StringTypeOrNil(r.MaxPasswordAge, true)
	config.CheckMismatchedPDFormattedAttributes("max_password_age",
		expectedValues.MaxPasswordAge, state.MaxPasswordAge, diagnostics)
	state.PasswordExpirationWarningInterval = internaltypes.StringTypeOrNil(r.PasswordExpirationWarningInterval, true)
	config.CheckMismatchedPDFormattedAttributes("password_expiration_warning_interval",
		expectedValues.PasswordExpirationWarningInterval, state.PasswordExpirationWarningInterval, diagnostics)
	state.ExpirePasswordsWithoutWarning = internaltypes.BoolTypeOrNil(r.ExpirePasswordsWithoutWarning)
	state.ReturnPasswordExpirationControls = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyReturnPasswordExpirationControlsProp(r.ReturnPasswordExpirationControls), true)
	state.AllowExpiredPasswordChanges = internaltypes.BoolTypeOrNil(r.AllowExpiredPasswordChanges)
	state.GraceLoginCount = internaltypes.Int64TypeOrNil(r.GraceLoginCount)
	state.RequireChangeByTime = internaltypes.StringTypeOrNil(r.RequireChangeByTime, internaltypes.IsEmptyString(expectedValues.RequireChangeByTime))
	state.LockoutFailureCount = internaltypes.Int64TypeOrNil(r.LockoutFailureCount)
	state.LockoutDuration = internaltypes.StringTypeOrNil(r.LockoutDuration, true)
	config.CheckMismatchedPDFormattedAttributes("lockout_duration",
		expectedValues.LockoutDuration, state.LockoutDuration, diagnostics)
	state.LockoutFailureExpirationInterval = internaltypes.StringTypeOrNil(r.LockoutFailureExpirationInterval, true)
	config.CheckMismatchedPDFormattedAttributes("lockout_failure_expiration_interval",
		expectedValues.LockoutFailureExpirationInterval, state.LockoutFailureExpirationInterval, diagnostics)
	state.IgnoreDuplicatePasswordFailures = internaltypes.BoolTypeOrNil(r.IgnoreDuplicatePasswordFailures)
	state.FailureLockoutAction = internaltypes.StringTypeOrNil(r.FailureLockoutAction, internaltypes.IsEmptyString(expectedValues.FailureLockoutAction))
	state.IdleLockoutInterval = internaltypes.StringTypeOrNil(r.IdleLockoutInterval, true)
	config.CheckMismatchedPDFormattedAttributes("idle_lockout_interval",
		expectedValues.IdleLockoutInterval, state.IdleLockoutInterval, diagnostics)
	state.AllowUserPasswordChanges = internaltypes.BoolTypeOrNil(r.AllowUserPasswordChanges)
	state.PasswordChangeRequiresCurrentPassword = internaltypes.BoolTypeOrNil(r.PasswordChangeRequiresCurrentPassword)
	state.PasswordRetirementBehavior = internaltypes.GetStringSet(
		client.StringSliceEnumpasswordPolicyPasswordRetirementBehaviorProp(r.PasswordRetirementBehavior))
	state.MaxRetiredPasswordAge = internaltypes.StringTypeOrNil(r.MaxRetiredPasswordAge, true)
	config.CheckMismatchedPDFormattedAttributes("max_retired_password_age",
		expectedValues.MaxRetiredPasswordAge, state.MaxRetiredPasswordAge, diagnostics)
	state.AllowedPasswordResetTokenUseCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpasswordPolicyAllowedPasswordResetTokenUseConditionProp(r.AllowedPasswordResetTokenUseCondition))
	state.ForceChangeOnAdd = internaltypes.BoolTypeOrNil(r.ForceChangeOnAdd)
	state.ForceChangeOnReset = internaltypes.BoolTypeOrNil(r.ForceChangeOnReset)
	state.MaxPasswordResetAge = internaltypes.StringTypeOrNil(r.MaxPasswordResetAge, true)
	config.CheckMismatchedPDFormattedAttributes("max_password_reset_age",
		expectedValues.MaxPasswordResetAge, state.MaxPasswordResetAge, diagnostics)
	state.SkipValidationForAdministrators = internaltypes.BoolTypeOrNil(r.SkipValidationForAdministrators)
	state.MaximumRecentLoginHistorySuccessfulAuthenticationCount = internaltypes.Int64TypeOrNil(r.MaximumRecentLoginHistorySuccessfulAuthenticationCount)
	state.MaximumRecentLoginHistorySuccessfulAuthenticationDuration = internaltypes.StringTypeOrNil(r.MaximumRecentLoginHistorySuccessfulAuthenticationDuration, internaltypes.IsEmptyString(expectedValues.MaximumRecentLoginHistorySuccessfulAuthenticationDuration))
	config.CheckMismatchedPDFormattedAttributes("maximum_recent_login_history_successful_authentication_duration",
		expectedValues.MaximumRecentLoginHistorySuccessfulAuthenticationDuration, state.MaximumRecentLoginHistorySuccessfulAuthenticationDuration, diagnostics)
	state.MaximumRecentLoginHistoryFailedAuthenticationCount = internaltypes.Int64TypeOrNil(r.MaximumRecentLoginHistoryFailedAuthenticationCount)
	state.MaximumRecentLoginHistoryFailedAuthenticationDuration = internaltypes.StringTypeOrNil(r.MaximumRecentLoginHistoryFailedAuthenticationDuration, internaltypes.IsEmptyString(expectedValues.MaximumRecentLoginHistoryFailedAuthenticationDuration))
	config.CheckMismatchedPDFormattedAttributes("maximum_recent_login_history_failed_authentication_duration",
		expectedValues.MaximumRecentLoginHistoryFailedAuthenticationDuration, state.MaximumRecentLoginHistoryFailedAuthenticationDuration, diagnostics)
	state.RecentLoginHistorySimilarAttemptBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyRecentLoginHistorySimilarAttemptBehaviorProp(r.RecentLoginHistorySimilarAttemptBehavior), true)
	state.LastLoginIPAddressAttribute = internaltypes.StringTypeOrNil(r.LastLoginIPAddressAttribute, internaltypes.IsEmptyString(expectedValues.LastLoginIPAddressAttribute))
	state.LastLoginTimeAttribute = internaltypes.StringTypeOrNil(r.LastLoginTimeAttribute, true)
	state.LastLoginTimeFormat = internaltypes.StringTypeOrNil(r.LastLoginTimeFormat, internaltypes.IsEmptyString(expectedValues.LastLoginTimeFormat))
	state.PreviousLastLoginTimeFormat = internaltypes.GetStringSet(r.PreviousLastLoginTimeFormat)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPasswordPolicyOperations(plan passwordPolicyResourceModel, state passwordPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireSecureAuthentication, state.RequireSecureAuthentication, "require-secure-authentication")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequireSecurePasswordChanges, state.RequireSecurePasswordChanges, "require-secure-password-changes")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccountStatusNotificationHandler, state.AccountStatusNotificationHandler, "account-status-notification-handler")
	operations.AddStringOperationIfNecessary(&ops, plan.StateUpdateFailurePolicy, state.StateUpdateFailurePolicy, "state-update-failure-policy")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableDebug, state.EnableDebug, "enable-debug")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordAttribute, state.PasswordAttribute, "password-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultPasswordStorageScheme, state.DefaultPasswordStorageScheme, "default-password-storage-scheme")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DeprecatedPasswordStorageScheme, state.DeprecatedPasswordStorageScheme, "deprecated-password-storage-scheme")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowMultiplePasswordValues, state.AllowMultiplePasswordValues, "allow-multiple-password-values")
	operations.AddStringOperationIfNecessary(&ops, plan.AllowPreEncodedPasswords, state.AllowPreEncodedPasswords, "allow-pre-encoded-passwords")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PasswordValidator, state.PasswordValidator, "password-validator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.BindPasswordValidator, state.BindPasswordValidator, "bind-password-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.MinimumBindPasswordValidationFrequency, state.MinimumBindPasswordValidationFrequency, "minimum-bind-password-validation-frequency")
	operations.AddStringOperationIfNecessary(&ops, plan.BindPasswordValidationFailureAction, state.BindPasswordValidationFailureAction, "bind-password-validation-failure-action")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordGenerator, state.PasswordGenerator, "password-generator")
	operations.AddInt64OperationIfNecessary(&ops, plan.PasswordHistoryCount, state.PasswordHistoryCount, "password-history-count")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordHistoryDuration, state.PasswordHistoryDuration, "password-history-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.MinPasswordAge, state.MinPasswordAge, "min-password-age")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxPasswordAge, state.MaxPasswordAge, "max-password-age")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordExpirationWarningInterval, state.PasswordExpirationWarningInterval, "password-expiration-warning-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.ExpirePasswordsWithoutWarning, state.ExpirePasswordsWithoutWarning, "expire-passwords-without-warning")
	operations.AddStringOperationIfNecessary(&ops, plan.ReturnPasswordExpirationControls, state.ReturnPasswordExpirationControls, "return-password-expiration-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowExpiredPasswordChanges, state.AllowExpiredPasswordChanges, "allow-expired-password-changes")
	operations.AddInt64OperationIfNecessary(&ops, plan.GraceLoginCount, state.GraceLoginCount, "grace-login-count")
	operations.AddStringOperationIfNecessary(&ops, plan.RequireChangeByTime, state.RequireChangeByTime, "require-change-by-time")
	operations.AddInt64OperationIfNecessary(&ops, plan.LockoutFailureCount, state.LockoutFailureCount, "lockout-failure-count")
	operations.AddStringOperationIfNecessary(&ops, plan.LockoutDuration, state.LockoutDuration, "lockout-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.LockoutFailureExpirationInterval, state.LockoutFailureExpirationInterval, "lockout-failure-expiration-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.IgnoreDuplicatePasswordFailures, state.IgnoreDuplicatePasswordFailures, "ignore-duplicate-password-failures")
	operations.AddStringOperationIfNecessary(&ops, plan.FailureLockoutAction, state.FailureLockoutAction, "failure-lockout-action")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleLockoutInterval, state.IdleLockoutInterval, "idle-lockout-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowUserPasswordChanges, state.AllowUserPasswordChanges, "allow-user-password-changes")
	operations.AddBoolOperationIfNecessary(&ops, plan.PasswordChangeRequiresCurrentPassword, state.PasswordChangeRequiresCurrentPassword, "password-change-requires-current-password")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PasswordRetirementBehavior, state.PasswordRetirementBehavior, "password-retirement-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxRetiredPasswordAge, state.MaxRetiredPasswordAge, "max-retired-password-age")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedPasswordResetTokenUseCondition, state.AllowedPasswordResetTokenUseCondition, "allowed-password-reset-token-use-condition")
	operations.AddBoolOperationIfNecessary(&ops, plan.ForceChangeOnAdd, state.ForceChangeOnAdd, "force-change-on-add")
	operations.AddBoolOperationIfNecessary(&ops, plan.ForceChangeOnReset, state.ForceChangeOnReset, "force-change-on-reset")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxPasswordResetAge, state.MaxPasswordResetAge, "max-password-reset-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.SkipValidationForAdministrators, state.SkipValidationForAdministrators, "skip-validation-for-administrators")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumRecentLoginHistorySuccessfulAuthenticationCount, state.MaximumRecentLoginHistorySuccessfulAuthenticationCount, "maximum-recent-login-history-successful-authentication-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumRecentLoginHistorySuccessfulAuthenticationDuration, state.MaximumRecentLoginHistorySuccessfulAuthenticationDuration, "maximum-recent-login-history-successful-authentication-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumRecentLoginHistoryFailedAuthenticationCount, state.MaximumRecentLoginHistoryFailedAuthenticationCount, "maximum-recent-login-history-failed-authentication-count")
	operations.AddStringOperationIfNecessary(&ops, plan.MaximumRecentLoginHistoryFailedAuthenticationDuration, state.MaximumRecentLoginHistoryFailedAuthenticationDuration, "maximum-recent-login-history-failed-authentication-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.RecentLoginHistorySimilarAttemptBehavior, state.RecentLoginHistorySimilarAttemptBehavior, "recent-login-history-similar-attempt-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.LastLoginIPAddressAttribute, state.LastLoginIPAddressAttribute, "last-login-ip-address-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.LastLoginTimeAttribute, state.LastLoginTimeAttribute, "last-login-time-attribute")
	operations.AddStringOperationIfNecessary(&ops, plan.LastLoginTimeFormat, state.LastLoginTimeFormat, "last-login-time-format")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.PreviousLastLoginTimeFormat, state.PreviousLastLoginTimeFormat, "previous-last-login-time-format")
	return ops
}

// Create a password-policy password-policy
func (r *passwordPolicyResource) CreatePasswordPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordPolicyResourceModel) (*passwordPolicyResourceModel, error) {
	var DefaultPasswordStorageSchemeSlice []string
	plan.DefaultPasswordStorageScheme.ElementsAs(ctx, &DefaultPasswordStorageSchemeSlice, false)
	addRequest := client.NewAddPasswordPolicyRequest(plan.Name.ValueString(),
		plan.PasswordAttribute.ValueString(),
		DefaultPasswordStorageSchemeSlice)
	err := addOptionalPasswordPolicyFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Policy", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordPolicyApi.AddPasswordPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordPolicyRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.PasswordPolicyApi.AddPasswordPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordPolicyResourceModel
	readPasswordPolicyResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passwordPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreatePasswordPolicy(ctx, req, resp, plan)
	if err != nil {
		return
	}

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
func (r *defaultPasswordPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordPolicyApi.GetPasswordPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passwordPolicyResourceModel
	readPasswordPolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PasswordPolicyApi.UpdatePasswordPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPasswordPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PasswordPolicyApi.UpdatePasswordPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPasswordPolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *passwordPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordPolicy(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPasswordPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordPolicy(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPasswordPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state passwordPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PasswordPolicyApi.GetPasswordPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Password Policy", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Policy", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPasswordPolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *passwordPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPasswordPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePasswordPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passwordPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passwordPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PasswordPolicyApi.UpdatePasswordPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPasswordPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PasswordPolicyApi.UpdatePasswordPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPasswordPolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPasswordPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passwordPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passwordPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PasswordPolicyApi.DeletePasswordPolicyExecute(r.apiClient.PasswordPolicyApi.DeletePasswordPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Password Policy", err, httpResp)
		return
	}
}

func (r *passwordPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordPolicy(ctx, req, resp)
}

func (r *defaultPasswordPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordPolicy(ctx, req, resp)
}

func importPasswordPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
