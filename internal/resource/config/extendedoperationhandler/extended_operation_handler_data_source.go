package extendedoperationhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &extendedOperationHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &extendedOperationHandlerDataSource{}
)

// Create a Extended Operation Handler data source
func NewExtendedOperationHandlerDataSource() datasource.DataSource {
	return &extendedOperationHandlerDataSource{}
}

// extendedOperationHandlerDataSource is the datasource implementation.
type extendedOperationHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *extendedOperationHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extended_operation_handler"
}

// Configure adds the provider configured client to the data source.
func (r *extendedOperationHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type extendedOperationHandlerDataSourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Name                                  types.String `tfsdk:"name"`
	Type                                  types.String `tfsdk:"type"`
	ExtensionClass                        types.String `tfsdk:"extension_class"`
	ExtensionArgument                     types.Set    `tfsdk:"extension_argument"`
	DefaultPasswordPolicy                 types.String `tfsdk:"default_password_policy"`
	DefaultTokenDeliveryMechanism         types.Set    `tfsdk:"default_token_delivery_mechanism"`
	PasswordResetTokenValidityDuration    types.String `tfsdk:"password_reset_token_validity_duration"`
	DefaultPasswordGenerator              types.String `tfsdk:"default_password_generator"`
	MaximumPasswordsPerRequest            types.Int64  `tfsdk:"maximum_passwords_per_request"`
	MaximumValidationAttemptsPerPassword  types.Int64  `tfsdk:"maximum_validation_attempts_per_password"`
	PasswordGenerator                     types.String `tfsdk:"password_generator"`
	DefaultOTPDeliveryMechanism           types.Set    `tfsdk:"default_otp_delivery_mechanism"`
	DefaultSingleUseTokenValidityDuration types.String `tfsdk:"default_single_use_token_validity_duration"`
	IdentityMapper                        types.String `tfsdk:"identity_mapper"`
	AllowRemotelyProvidedCertificates     types.Bool   `tfsdk:"allow_remotely_provided_certificates"`
	AllowedOperation                      types.Set    `tfsdk:"allowed_operation"`
	ConnectionCriteria                    types.String `tfsdk:"connection_criteria"`
	RequestCriteria                       types.String `tfsdk:"request_criteria"`
	SharedSecretAttributeType             types.String `tfsdk:"shared_secret_attribute_type"`
	TimeIntervalDuration                  types.String `tfsdk:"time_interval_duration"`
	AdjacentIntervalsToCheck              types.Int64  `tfsdk:"adjacent_intervals_to_check"`
	PreventTOTPReuse                      types.Bool   `tfsdk:"prevent_totp_reuse"`
	Description                           types.String `tfsdk:"description"`
	Enabled                               types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *extendedOperationHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Extended Operation Handler resource. Options are ['cancel', 'validate-totp-password', 'replace-certificate', 'get-connection-id', 'multi-update', 'notification-subscription', 'password-modify', 'custom', 'collect-support-data', 'export-reversible-passwords', 'batched-transactions', 'get-changelog-batch', 'get-supported-otp-delivery-mechanisms', 'single-use-tokens', 'generate-password', 'who-am-i', 'start-tls', 'deliver-password-reset-token', 'password-policy-state', 'get-password-quality-requirements', 'deliver-otp', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Extended Operation Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Extended Operation Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_password_policy": schema.StringAttribute{
				Description: "The default password policy that should be used when generating and validating passwords if the request does not specify an alternate policy. If this is not provided, then this Generate Password Extended Operation Handler will use the default password policy defined in the global configuration.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_token_delivery_mechanism": schema.SetAttribute{
				Description: "The set of delivery mechanisms that may be used to deliver password reset tokens to users for requests that do not specify one or more preferred delivery mechanisms.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"password_reset_token_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a password reset token should be considered valid.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_password_generator": schema.StringAttribute{
				Description: "The default password generator that will be used if the selected password policy is not configured with a password generator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_passwords_per_request": schema.Int64Attribute{
				Description: "The maximum number of passwords that may be generated and returned to the client for a single request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_validation_attempts_per_password": schema.Int64Attribute{
				Description: "The maximum number of attempts that the server may use to generate a password that passes validation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_generator": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `single-use-tokens`: The password generator that will be used to create the single-use token values to be delivered to the end user. When the `type` attribute is set to `deliver-password-reset-token`: The password generator that will be used to create the password reset token values to be delivered to the end user. When the `type` attribute is set to `deliver-otp`: The password generator that will be used to create the one-time password values to be delivered to the end user.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `single-use-tokens`: The password generator that will be used to create the single-use token values to be delivered to the end user.\n  - `deliver-password-reset-token`: The password generator that will be used to create the password reset token values to be delivered to the end user.\n  - `deliver-otp`: The password generator that will be used to create the one-time password values to be delivered to the end user.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"default_otp_delivery_mechanism": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `single-use-tokens`: The set of delivery mechanisms that may be used to deliver single-use tokens to users in requests that do not specify one or more preferred delivery mechanisms. When the `type` attribute is set to `deliver-otp`: The set of delivery mechanisms that may be used to deliver one-time passwords to users in requests that do not specify one or more preferred delivery mechanisms.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `single-use-tokens`: The set of delivery mechanisms that may be used to deliver single-use tokens to users in requests that do not specify one or more preferred delivery mechanisms.\n  - `deliver-otp`: The set of delivery mechanisms that may be used to deliver one-time passwords to users in requests that do not specify one or more preferred delivery mechanisms.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"default_single_use_token_validity_duration": schema.StringAttribute{
				Description: "The default length of time that a single-use token will be considered valid by the server if the client doesn't specify a duration in the deliver single-use token request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `password-modify`: Specifies the name of the identity mapper that should be used in conjunction with the password modify extended operation. When the `type` attribute is set to `deliver-otp`: The identity mapper that should be used to identify the user(s) targeted by the authentication identity contained in the extended request. This will only be used for \"u:\"-style authentication identities.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `password-modify`: Specifies the name of the identity mapper that should be used in conjunction with the password modify extended operation.\n  - `deliver-otp`: The identity mapper that should be used to identify the user(s) targeted by the authentication identity contained in the extended request. This will only be used for \"u:\"-style authentication identities.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"allow_remotely_provided_certificates": schema.BoolAttribute{
				Description: "Indicates whether clients should be allowed to directly provide a new listener or inter-server certificate chain in the extended request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_operation": schema.SetAttribute{
				Description: "The types of replace certificate operations that clients will be allowed to request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A set of criteria that client connections must satisfy before they will be allowed to request the associated extended operations.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "A set of criteria that the extended requests must satisfy before they will be processed by the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"shared_secret_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute that will be used to hold the shared secret key used during TOTP processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"time_interval_duration": schema.StringAttribute{
				Description: "The duration of the time interval used for TOTP processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"adjacent_intervals_to_check": schema.Int64Attribute{
				Description: "The number of adjacent time intervals (both before and after the current time) that should be checked when performing authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prevent_totp_reuse": schema.BoolAttribute{
				Description: "Indicates whether to prevent clients from re-using TOTP passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Extended Operation Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Extended Operation Handler is enabled (that is, whether the types of extended operations are allowed in the server).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a CancelExtendedOperationHandlerResponse object into the model struct
func readCancelExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.CancelExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("cancel")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ValidateTotpPasswordExtendedOperationHandlerResponse object into the model struct
func readValidateTotpPasswordExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.ValidateTotpPasswordExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("validate-totp-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SharedSecretAttributeType = internaltypes.StringTypeOrNil(r.SharedSecretAttributeType, false)
	state.TimeIntervalDuration = internaltypes.StringTypeOrNil(r.TimeIntervalDuration, false)
	state.AdjacentIntervalsToCheck = internaltypes.Int64TypeOrNil(r.AdjacentIntervalsToCheck)
	state.PreventTOTPReuse = internaltypes.BoolTypeOrNil(r.PreventTOTPReuse)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ReplaceCertificateExtendedOperationHandlerResponse object into the model struct
func readReplaceCertificateExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.ReplaceCertificateExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replace-certificate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowRemotelyProvidedCertificates = internaltypes.BoolTypeOrNil(r.AllowRemotelyProvidedCertificates)
	state.AllowedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumextendedOperationHandlerAllowedOperationProp(r.AllowedOperation))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GetConnectionIdExtendedOperationHandlerResponse object into the model struct
func readGetConnectionIdExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.GetConnectionIdExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-connection-id")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a MultiUpdateExtendedOperationHandlerResponse object into the model struct
func readMultiUpdateExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.MultiUpdateExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("multi-update")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a NotificationSubscriptionExtendedOperationHandlerResponse object into the model struct
func readNotificationSubscriptionExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.NotificationSubscriptionExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("notification-subscription")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PasswordModifyExtendedOperationHandlerResponse object into the model struct
func readPasswordModifyExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.PasswordModifyExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-modify")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CustomExtendedOperationHandlerResponse object into the model struct
func readCustomExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.CustomExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CollectSupportDataExtendedOperationHandlerResponse object into the model struct
func readCollectSupportDataExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.CollectSupportDataExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("collect-support-data")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ExportReversiblePasswordsExtendedOperationHandlerResponse object into the model struct
func readExportReversiblePasswordsExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.ExportReversiblePasswordsExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("export-reversible-passwords")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a BatchedTransactionsExtendedOperationHandlerResponse object into the model struct
func readBatchedTransactionsExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.BatchedTransactionsExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("batched-transactions")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GetChangelogBatchExtendedOperationHandlerResponse object into the model struct
func readGetChangelogBatchExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.GetChangelogBatchExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-changelog-batch")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse object into the model struct
func readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-supported-otp-delivery-mechanisms")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SingleUseTokensExtendedOperationHandlerResponse object into the model struct
func readSingleUseTokensExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.SingleUseTokensExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("single-use-tokens")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.DefaultSingleUseTokenValidityDuration = internaltypes.StringTypeOrNil(r.DefaultSingleUseTokenValidityDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GeneratePasswordExtendedOperationHandlerResponse object into the model struct
func readGeneratePasswordExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.GeneratePasswordExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generate-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DefaultPasswordPolicy = internaltypes.StringTypeOrNil(r.DefaultPasswordPolicy, false)
	state.DefaultPasswordGenerator = types.StringValue(r.DefaultPasswordGenerator)
	state.MaximumPasswordsPerRequest = internaltypes.Int64TypeOrNil(r.MaximumPasswordsPerRequest)
	state.MaximumValidationAttemptsPerPassword = internaltypes.Int64TypeOrNil(r.MaximumValidationAttemptsPerPassword)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a WhoAmIExtendedOperationHandlerResponse object into the model struct
func readWhoAmIExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.WhoAmIExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("who-am-i")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a StartTlsExtendedOperationHandlerResponse object into the model struct
func readStartTlsExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.StartTlsExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("start-tls")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a DeliverPasswordResetTokenExtendedOperationHandlerResponse object into the model struct
func readDeliverPasswordResetTokenExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.DeliverPasswordResetTokenExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deliver-password-reset-token")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultTokenDeliveryMechanism = internaltypes.GetStringSet(r.DefaultTokenDeliveryMechanism)
	state.PasswordResetTokenValidityDuration = types.StringValue(r.PasswordResetTokenValidityDuration)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PasswordPolicyStateExtendedOperationHandlerResponse object into the model struct
func readPasswordPolicyStateExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.PasswordPolicyStateExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy-state")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GetPasswordQualityRequirementsExtendedOperationHandlerResponse object into the model struct
func readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-password-quality-requirements")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a DeliverOtpExtendedOperationHandlerResponse object into the model struct
func readDeliverOtpExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.DeliverOtpExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deliver-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyExtendedOperationHandlerResponse object into the model struct
func readThirdPartyExtendedOperationHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartyExtendedOperationHandlerResponse, state *extendedOperationHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *extendedOperationHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state extendedOperationHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerAPI.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.CancelExtendedOperationHandlerResponse != nil {
		readCancelExtendedOperationHandlerResponseDataSource(ctx, readResponse.CancelExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
		readValidateTotpPasswordExtendedOperationHandlerResponseDataSource(ctx, readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ReplaceCertificateExtendedOperationHandlerResponse != nil {
		readReplaceCertificateExtendedOperationHandlerResponseDataSource(ctx, readResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GetConnectionIdExtendedOperationHandlerResponse != nil {
		readGetConnectionIdExtendedOperationHandlerResponseDataSource(ctx, readResponse.GetConnectionIdExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.MultiUpdateExtendedOperationHandlerResponse != nil {
		readMultiUpdateExtendedOperationHandlerResponseDataSource(ctx, readResponse.MultiUpdateExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.NotificationSubscriptionExtendedOperationHandlerResponse != nil {
		readNotificationSubscriptionExtendedOperationHandlerResponseDataSource(ctx, readResponse.NotificationSubscriptionExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordModifyExtendedOperationHandlerResponse != nil {
		readPasswordModifyExtendedOperationHandlerResponseDataSource(ctx, readResponse.PasswordModifyExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomExtendedOperationHandlerResponse != nil {
		readCustomExtendedOperationHandlerResponseDataSource(ctx, readResponse.CustomExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CollectSupportDataExtendedOperationHandlerResponse != nil {
		readCollectSupportDataExtendedOperationHandlerResponseDataSource(ctx, readResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
		readExportReversiblePasswordsExtendedOperationHandlerResponseDataSource(ctx, readResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BatchedTransactionsExtendedOperationHandlerResponse != nil {
		readBatchedTransactionsExtendedOperationHandlerResponseDataSource(ctx, readResponse.BatchedTransactionsExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GetChangelogBatchExtendedOperationHandlerResponse != nil {
		readGetChangelogBatchExtendedOperationHandlerResponseDataSource(ctx, readResponse.GetChangelogBatchExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse != nil {
		readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDataSource(ctx, readResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SingleUseTokensExtendedOperationHandlerResponse != nil {
		readSingleUseTokensExtendedOperationHandlerResponseDataSource(ctx, readResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GeneratePasswordExtendedOperationHandlerResponse != nil {
		readGeneratePasswordExtendedOperationHandlerResponseDataSource(ctx, readResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.WhoAmIExtendedOperationHandlerResponse != nil {
		readWhoAmIExtendedOperationHandlerResponseDataSource(ctx, readResponse.WhoAmIExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.StartTlsExtendedOperationHandlerResponse != nil {
		readStartTlsExtendedOperationHandlerResponseDataSource(ctx, readResponse.StartTlsExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
		readDeliverPasswordResetTokenExtendedOperationHandlerResponseDataSource(ctx, readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyStateExtendedOperationHandlerResponse != nil {
		readPasswordPolicyStateExtendedOperationHandlerResponseDataSource(ctx, readResponse.PasswordPolicyStateExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse != nil {
		readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDataSource(ctx, readResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DeliverOtpExtendedOperationHandlerResponse != nil {
		readDeliverOtpExtendedOperationHandlerResponseDataSource(ctx, readResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyExtendedOperationHandlerResponse != nil {
		readThirdPartyExtendedOperationHandlerResponseDataSource(ctx, readResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
