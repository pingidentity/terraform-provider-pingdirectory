// Copyright © 2025 Ping Identity Corporation

package passwordpolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passwordPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordPolicyDataSource{}
)

// Create a Password Policy data source
func NewPasswordPolicyDataSource() datasource.DataSource {
	return &passwordPolicyDataSource{}
}

// passwordPolicyDataSource is the datasource implementation.
type passwordPolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passwordPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_policy"
}

// Configure adds the provider configured client to the data source.
func (r *passwordPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordPolicyDataSourceModel struct {
	Id                                                        types.String `tfsdk:"id"`
	Name                                                      types.String `tfsdk:"name"`
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
	ReEncodePasswordsOnSchemeConfigChange                     types.Bool   `tfsdk:"re_encode_passwords_on_scheme_config_change"`
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
	SuppressRecentLoginHistoryUpdatesForUnusableAccounts      types.Bool   `tfsdk:"suppress_recent_login_history_updates_for_unusable_accounts"`
	RecentLoginHistorySimilarAttemptBehavior                  types.String `tfsdk:"recent_login_history_similar_attempt_behavior"`
	LastLoginIPAddressAttribute                               types.String `tfsdk:"last_login_ip_address_attribute"`
	LastLoginTimeAttribute                                    types.String `tfsdk:"last_login_time_attribute"`
	LastLoginTimeFormat                                       types.String `tfsdk:"last_login_time_format"`
	PreviousLastLoginTimeFormat                               types.Set    `tfsdk:"previous_last_login_time_format"`
}

// GetSchema defines the schema for the datasource.
func (r *passwordPolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Password Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Policy resource. Options are ['password-policy']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Policy",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_secure_authentication": schema.BoolAttribute{
				Description: "Indicates whether users with the associated password policy are required to authenticate in a secure manner.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_secure_password_changes": schema.BoolAttribute{
				Description: "Indicates whether users with the associated password policy are required to change their password in a secure manner that does not expose the credentials.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"account_status_notification_handler": schema.SetAttribute{
				Description: "Specifies the names of the account status notification handlers that are used with the associated password storage scheme.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"state_update_failure_policy": schema.StringAttribute{
				Description: "Specifies how the server deals with the inability to update password policy state information during an authentication attempt.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_debug": schema.BoolAttribute{
				Description: "Indicates whether to enable debugging for the password policy state.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_attribute": schema.StringAttribute{
				Description: "Specifies the attribute type used to hold user passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_password_storage_scheme": schema.SetAttribute{
				Description: "Specifies the names of the password storage schemes that are used to encode clear-text passwords for this password policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"deprecated_password_storage_scheme": schema.SetAttribute{
				Description: "Specifies the names of the password storage schemes that are considered deprecated for this password policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"re_encode_passwords_on_scheme_config_change": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Indicates whether to re-encode passwords on authentication if the configuration for the underlying password storage scheme has changed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_multiple_password_values": schema.BoolAttribute{
				Description: "Indicates whether user entries can have multiple distinct values for the password attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_pre_encoded_passwords": schema.StringAttribute{
				Description: "Indicates whether users can change their passwords by providing a pre-encoded value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_validator": schema.SetAttribute{
				Description: "Specifies the names of the password validators that are used with the associated password storage scheme.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"bind_password_validator": schema.SetAttribute{
				Description: "Specifies the names of the password validators that should be invoked for bind operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"minimum_bind_password_validation_frequency": schema.StringAttribute{
				Description: "Indicates how frequently password validation should be performed during bind operations for each user to whom this password policy is assigned.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_password_validation_failure_action": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if a bind password fails validation by one or more of the configured bind password validators.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_generator": schema.StringAttribute{
				Description: "Specifies the name of the password generator that is used with the associated password policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_history_count": schema.Int64Attribute{
				Description: "Specifies the maximum number of former passwords to maintain in the password history.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_history_duration": schema.StringAttribute{
				Description: "Specifies the maximum length of time that passwords remain in the password history.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_password_age": schema.StringAttribute{
				Description: "Specifies the minimum length of time after a password change before the user is allowed to change the password again.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_password_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a user can continue using the same password before it must be changed (that is, the password expiration interval).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_expiration_warning_interval": schema.StringAttribute{
				Description: "Specifies the maximum length of time before a user's password actually expires that the server begins to include warning notifications in bind responses for that user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"expire_passwords_without_warning": schema.BoolAttribute{
				Description: "Indicates whether the Directory Server allows a user's password to expire even if that user has never seen an expiration warning notification.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"return_password_expiration_controls": schema.StringAttribute{
				Description: "Indicates whether the server should return the password expiring and password expired response controls (as described in draft-vchu-ldap-pwd-policy).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_expired_password_changes": schema.BoolAttribute{
				Description: "Indicates whether a user whose password is expired is still allowed to change that password using the password modify extended operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"grace_login_count": schema.Int64Attribute{
				Description: "Specifies the number of grace logins that a user is allowed after the account has expired to allow that user to choose a new password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"require_change_by_time": schema.StringAttribute{
				Description: "Specifies the time by which all users with the associated password policy must change their passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"lockout_failure_count": schema.Int64Attribute{
				Description: "Specifies the maximum number of authentication failures that a user is allowed before the account is locked out.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"lockout_duration": schema.StringAttribute{
				Description: "Specifies the length of time that an account is locked after too many authentication failures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"lockout_failure_expiration_interval": schema.StringAttribute{
				Description: "Specifies the length of time before an authentication failure is no longer counted against a user for the purposes of account lockout.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ignore_duplicate_password_failures": schema.BoolAttribute{
				Description: "Indicates whether to ignore subsequent authentication failures using the same password as an earlier failed authentication attempt (within the time frame defined by the lockout failure expiration interval). If this option is \"true\", then multiple failed attempts using the same password will be considered only a single failure. If this option is \"false\", then any failure will be tracked regardless of whether it used the same password as an earlier attempt.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"failure_lockout_action": schema.StringAttribute{
				Description: "The action that the server should take for authentication attempts that target a user with more than the configured number of outstanding authentication failures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"idle_lockout_interval": schema.StringAttribute{
				Description: "Specifies the maximum length of time that an account may remain idle (that is, the associated user does not authenticate to the server) before that user is locked out.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_user_password_changes": schema.BoolAttribute{
				Description: "Indicates whether users can change their own passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_change_requires_current_password": schema.BoolAttribute{
				Description: "Indicates whether user password changes must use the password modify extended operation and must include the user's current password before the change is allowed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_retirement_behavior": schema.SetAttribute{
				Description: "Specifies the conditions under which the server may retire a user's current password in the course of setting a new password for that user (whether via a modify operation or a password modify extended operation).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"max_retired_password_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that a retired password should be considered valid and may be used to authenticate to the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_password_reset_token_use_condition": schema.SetAttribute{
				Description: "The set of conditions under which a user governed by this Password Policy will be permitted to generate a password reset token via the deliver password reset token extended operation, and to use that token in lieu of the current password via the password modify extended operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"force_change_on_add": schema.BoolAttribute{
				Description: "Indicates whether users are forced to change their passwords upon first authenticating to the Directory Server after their account has been created.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"force_change_on_reset": schema.BoolAttribute{
				Description: "Indicates whether users are forced to change their passwords if they are reset by an administrator. If a user's password is changed by any other user, that is considered an administrative password reset.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_password_reset_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that users have to change passwords after they have been reset by an administrator before they become locked.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"skip_validation_for_administrators": schema.BoolAttribute{
				Description: "Indicates whether passwords set by administrators are allowed to bypass the password validation process that is required for user password changes.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_recent_login_history_successful_authentication_count": schema.Int64Attribute{
				Description: "The maximum number of successful authentication attempts to include in the recent login history for each account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_recent_login_history_successful_authentication_duration": schema.StringAttribute{
				Description: "The maximum age of successful authentication attempts to include in the recent login history for each account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_recent_login_history_failed_authentication_count": schema.Int64Attribute{
				Description: "The maximum number of failed authentication attempts to include in the recent login history for each account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_recent_login_history_failed_authentication_duration": schema.StringAttribute{
				Description: "The maximum age of failed authentication attempts to include in the recent login history for each account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"suppress_recent_login_history_updates_for_unusable_accounts": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.3.0.0+. Indicates whether the server should suppress updates to a user's recent login history as a result of authentication attempts that fail because the account is in an unusable state (e.g., if the account is administratively disabled, if the account is locked, or if the password is expired).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"recent_login_history_similar_attempt_behavior": schema.StringAttribute{
				Description: "The behavior that the server will exhibit when multiple similar authentication attempts (with the same values for the successful, authentication-method, client-ip-address, and failure-reason fields) are processed for an account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"last_login_ip_address_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute type that is used to hold the IP address of the client from which the user last authenticated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"last_login_time_attribute": schema.StringAttribute{
				Description: "Specifies the name or OID of the attribute type that is used to hold the last login time for users with the associated password policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"last_login_time_format": schema.StringAttribute{
				Description: "Specifies the format string that is used to generate the last login time value for users with the associated password policy. Last login time values will be written using the UTC (also known as GMT, or Greenwich Mean Time) time zone.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"previous_last_login_time_format": schema.SetAttribute{
				Description: "Specifies the format string(s) that might have been used with the last login time at any point in the past for users associated with the password policy.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a PasswordPolicyResponse object into the model struct
func readPasswordPolicyResponseDataSource(ctx context.Context, r *client.PasswordPolicyResponse, state *passwordPolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.RequireSecureAuthentication = internaltypes.BoolTypeOrNil(r.RequireSecureAuthentication)
	state.RequireSecurePasswordChanges = internaltypes.BoolTypeOrNil(r.RequireSecurePasswordChanges)
	state.AccountStatusNotificationHandler = internaltypes.GetStringSet(r.AccountStatusNotificationHandler)
	state.StateUpdateFailurePolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyStateUpdateFailurePolicyProp(r.StateUpdateFailurePolicy), false)
	state.EnableDebug = internaltypes.BoolTypeOrNil(r.EnableDebug)
	state.PasswordAttribute = types.StringValue(r.PasswordAttribute)
	state.DefaultPasswordStorageScheme = internaltypes.GetStringSet(r.DefaultPasswordStorageScheme)
	state.DeprecatedPasswordStorageScheme = internaltypes.GetStringSet(r.DeprecatedPasswordStorageScheme)
	state.ReEncodePasswordsOnSchemeConfigChange = internaltypes.BoolTypeOrNil(r.ReEncodePasswordsOnSchemeConfigChange)
	state.AllowMultiplePasswordValues = internaltypes.BoolTypeOrNil(r.AllowMultiplePasswordValues)
	state.AllowPreEncodedPasswords = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyAllowPreEncodedPasswordsProp(r.AllowPreEncodedPasswords), false)
	state.PasswordValidator = internaltypes.GetStringSet(r.PasswordValidator)
	state.BindPasswordValidator = internaltypes.GetStringSet(r.BindPasswordValidator)
	state.MinimumBindPasswordValidationFrequency = internaltypes.StringTypeOrNil(r.MinimumBindPasswordValidationFrequency, false)
	state.BindPasswordValidationFailureAction = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyBindPasswordValidationFailureActionProp(r.BindPasswordValidationFailureAction), false)
	state.PasswordGenerator = internaltypes.StringTypeOrNil(r.PasswordGenerator, false)
	state.PasswordHistoryCount = internaltypes.Int64TypeOrNil(r.PasswordHistoryCount)
	state.PasswordHistoryDuration = internaltypes.StringTypeOrNil(r.PasswordHistoryDuration, false)
	state.MinPasswordAge = internaltypes.StringTypeOrNil(r.MinPasswordAge, false)
	state.MaxPasswordAge = internaltypes.StringTypeOrNil(r.MaxPasswordAge, false)
	state.PasswordExpirationWarningInterval = internaltypes.StringTypeOrNil(r.PasswordExpirationWarningInterval, false)
	state.ExpirePasswordsWithoutWarning = internaltypes.BoolTypeOrNil(r.ExpirePasswordsWithoutWarning)
	state.ReturnPasswordExpirationControls = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyReturnPasswordExpirationControlsProp(r.ReturnPasswordExpirationControls), false)
	state.AllowExpiredPasswordChanges = internaltypes.BoolTypeOrNil(r.AllowExpiredPasswordChanges)
	state.GraceLoginCount = internaltypes.Int64TypeOrNil(r.GraceLoginCount)
	state.RequireChangeByTime = internaltypes.StringTypeOrNil(r.RequireChangeByTime, false)
	state.LockoutFailureCount = internaltypes.Int64TypeOrNil(r.LockoutFailureCount)
	state.LockoutDuration = internaltypes.StringTypeOrNil(r.LockoutDuration, false)
	state.LockoutFailureExpirationInterval = internaltypes.StringTypeOrNil(r.LockoutFailureExpirationInterval, false)
	state.IgnoreDuplicatePasswordFailures = internaltypes.BoolTypeOrNil(r.IgnoreDuplicatePasswordFailures)
	state.FailureLockoutAction = internaltypes.StringTypeOrNil(r.FailureLockoutAction, false)
	state.IdleLockoutInterval = internaltypes.StringTypeOrNil(r.IdleLockoutInterval, false)
	state.AllowUserPasswordChanges = internaltypes.BoolTypeOrNil(r.AllowUserPasswordChanges)
	state.PasswordChangeRequiresCurrentPassword = internaltypes.BoolTypeOrNil(r.PasswordChangeRequiresCurrentPassword)
	state.PasswordRetirementBehavior = internaltypes.GetStringSet(
		client.StringSliceEnumpasswordPolicyPasswordRetirementBehaviorProp(r.PasswordRetirementBehavior))
	state.MaxRetiredPasswordAge = internaltypes.StringTypeOrNil(r.MaxRetiredPasswordAge, false)
	state.AllowedPasswordResetTokenUseCondition = internaltypes.GetStringSet(
		client.StringSliceEnumpasswordPolicyAllowedPasswordResetTokenUseConditionProp(r.AllowedPasswordResetTokenUseCondition))
	state.ForceChangeOnAdd = internaltypes.BoolTypeOrNil(r.ForceChangeOnAdd)
	state.ForceChangeOnReset = internaltypes.BoolTypeOrNil(r.ForceChangeOnReset)
	state.MaxPasswordResetAge = internaltypes.StringTypeOrNil(r.MaxPasswordResetAge, false)
	state.SkipValidationForAdministrators = internaltypes.BoolTypeOrNil(r.SkipValidationForAdministrators)
	state.MaximumRecentLoginHistorySuccessfulAuthenticationCount = internaltypes.Int64TypeOrNil(r.MaximumRecentLoginHistorySuccessfulAuthenticationCount)
	state.MaximumRecentLoginHistorySuccessfulAuthenticationDuration = internaltypes.StringTypeOrNil(r.MaximumRecentLoginHistorySuccessfulAuthenticationDuration, false)
	state.MaximumRecentLoginHistoryFailedAuthenticationCount = internaltypes.Int64TypeOrNil(r.MaximumRecentLoginHistoryFailedAuthenticationCount)
	state.MaximumRecentLoginHistoryFailedAuthenticationDuration = internaltypes.StringTypeOrNil(r.MaximumRecentLoginHistoryFailedAuthenticationDuration, false)
	state.SuppressRecentLoginHistoryUpdatesForUnusableAccounts = internaltypes.BoolTypeOrNil(r.SuppressRecentLoginHistoryUpdatesForUnusableAccounts)
	state.RecentLoginHistorySimilarAttemptBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordPolicyRecentLoginHistorySimilarAttemptBehaviorProp(r.RecentLoginHistorySimilarAttemptBehavior), false)
	state.LastLoginIPAddressAttribute = internaltypes.StringTypeOrNil(r.LastLoginIPAddressAttribute, false)
	state.LastLoginTimeAttribute = internaltypes.StringTypeOrNil(r.LastLoginTimeAttribute, false)
	state.LastLoginTimeFormat = internaltypes.StringTypeOrNil(r.LastLoginTimeFormat, false)
	state.PreviousLastLoginTimeFormat = internaltypes.GetStringSet(r.PreviousLastLoginTimeFormat)
}

// Read resource information
func (r *passwordPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passwordPolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordPolicyAPI.GetPasswordPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPasswordPolicyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
