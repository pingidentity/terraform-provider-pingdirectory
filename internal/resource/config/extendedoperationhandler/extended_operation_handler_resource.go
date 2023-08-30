package extendedoperationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &extendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &extendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &extendedOperationHandlerResource{}
	_ resource.Resource                = &defaultExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultExtendedOperationHandlerResource{}
)

// Create a Extended Operation Handler resource
func NewExtendedOperationHandlerResource() resource.Resource {
	return &extendedOperationHandlerResource{}
}

func NewDefaultExtendedOperationHandlerResource() resource.Resource {
	return &defaultExtendedOperationHandlerResource{}
}

// extendedOperationHandlerResource is the resource implementation.
type extendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultExtendedOperationHandlerResource is the resource implementation.
type defaultExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *extendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_extended_operation_handler"
}

func (r *defaultExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_extended_operation_handler"
}

// Configure adds the provider configured client to the resource.
func (r *extendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type extendedOperationHandlerResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Name                                  types.String `tfsdk:"name"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	Type                                  types.String `tfsdk:"type"`
	ExtensionClass                        types.String `tfsdk:"extension_class"`
	ExtensionArgument                     types.Set    `tfsdk:"extension_argument"`
	DefaultTokenDeliveryMechanism         types.Set    `tfsdk:"default_token_delivery_mechanism"`
	PasswordResetTokenValidityDuration    types.String `tfsdk:"password_reset_token_validity_duration"`
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

type defaultExtendedOperationHandlerResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	Name                                  types.String `tfsdk:"name"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
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

// GetSchema defines the schema for the resource.
func (r *extendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	extendedOperationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	extendedOperationHandlerSchema(ctx, req, resp, true)
}

func extendedOperationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Extended Operation Handler resource. Options are ['cancel', 'validate-totp-password', 'replace-certificate', 'get-connection-id', 'multi-update', 'notification-subscription', 'password-modify', 'custom', 'collect-support-data', 'export-reversible-passwords', 'batched-transactions', 'get-changelog-batch', 'get-supported-otp-delivery-mechanisms', 'single-use-tokens', 'generate-password', 'who-am-i', 'start-tls', 'deliver-password-reset-token', 'password-policy-state', 'get-password-quality-requirements', 'deliver-otp', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"validate-totp-password", "replace-certificate", "collect-support-data", "export-reversible-passwords", "single-use-tokens", "deliver-password-reset-token", "deliver-otp", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Extended Operation Handler.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Extended Operation Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"default_token_delivery_mechanism": schema.SetAttribute{
				Description: "The set of delivery mechanisms that may be used to deliver password reset tokens to users for requests that do not specify one or more preferred delivery mechanisms.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"password_reset_token_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a password reset token should be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_generator": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `single-use-tokens`: The password generator that will be used to create the single-use token values to be delivered to the end user. When the `type` attribute is set to `deliver-password-reset-token`: The password generator that will be used to create the password reset token values to be delivered to the end user. When the `type` attribute is set to `deliver-otp`: The password generator that will be used to create the one-time password values to be delivered to the end user.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `single-use-tokens`: The password generator that will be used to create the single-use token values to be delivered to the end user.\n  - `deliver-password-reset-token`: The password generator that will be used to create the password reset token values to be delivered to the end user.\n  - `deliver-otp`: The password generator that will be used to create the one-time password values to be delivered to the end user.",
				Optional:            true,
			},
			"default_otp_delivery_mechanism": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `single-use-tokens`: The set of delivery mechanisms that may be used to deliver single-use tokens to users in requests that do not specify one or more preferred delivery mechanisms. When the `type` attribute is set to `deliver-otp`: The set of delivery mechanisms that may be used to deliver one-time passwords to users in requests that do not specify one or more preferred delivery mechanisms.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `single-use-tokens`: The set of delivery mechanisms that may be used to deliver single-use tokens to users in requests that do not specify one or more preferred delivery mechanisms.\n  - `deliver-otp`: The set of delivery mechanisms that may be used to deliver one-time passwords to users in requests that do not specify one or more preferred delivery mechanisms.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"default_single_use_token_validity_duration": schema.StringAttribute{
				Description: "The default length of time that a single-use token will be considered valid by the server if the client doesn't specify a duration in the deliver single-use token request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `password-modify`: Specifies the name of the identity mapper that should be used in conjunction with the password modify extended operation. When the `type` attribute is set to `deliver-otp`: The identity mapper that should be used to identify the user(s) targeted by the authentication identity contained in the extended request. This will only be used for \"u:\"-style authentication identities.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `password-modify`: Specifies the name of the identity mapper that should be used in conjunction with the password modify extended operation.\n  - `deliver-otp`: The identity mapper that should be used to identify the user(s) targeted by the authentication identity contained in the extended request. This will only be used for \"u:\"-style authentication identities.",
				Optional:            true,
			},
			"allow_remotely_provided_certificates": schema.BoolAttribute{
				Description: "Indicates whether clients should be allowed to directly provide a new listener or inter-server certificate chain in the extended request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_operation": schema.SetAttribute{
				Description: "The types of replace certificate operations that clients will be allowed to request.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_criteria": schema.StringAttribute{
				Description: "A set of criteria that client connections must satisfy before they will be allowed to request the associated extended operations.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "A set of criteria that the extended requests must satisfy before they will be processed by the server.",
				Optional:    true,
			},
			"shared_secret_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute that will be used to hold the shared secret key used during TOTP processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"time_interval_duration": schema.StringAttribute{
				Description: "The duration of the time interval used for TOTP processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"adjacent_intervals_to_check": schema.Int64Attribute{
				Description: "The number of adjacent time intervals (both before and after the current time) that should be checked when performing authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"prevent_totp_reuse": schema.BoolAttribute{
				Description: "Indicates whether to prevent clients from re-using TOTP passwords.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Extended Operation Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Extended Operation Handler is enabled (that is, whether the types of extended operations are allowed in the server).",
				Required:    true,
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
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"cancel", "validate-totp-password", "replace-certificate", "get-connection-id", "multi-update", "notification-subscription", "password-modify", "custom", "collect-support-data", "export-reversible-passwords", "batched-transactions", "get-changelog-batch", "get-supported-otp-delivery-mechanisms", "single-use-tokens", "generate-password", "who-am-i", "start-tls", "deliver-password-reset-token", "password-policy-state", "get-password-quality-requirements", "deliver-otp", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["default_password_policy"] = schema.StringAttribute{
			Description: "The default password policy that should be used when generating and validating passwords if the request does not specify an alternate policy. If this is not provided, then this Generate Password Extended Operation Handler will use the default password policy defined in the global configuration.",
		}
		schemaDef.Attributes["default_password_generator"] = schema.StringAttribute{
			Description: "The default password generator that will be used if the selected password policy is not configured with a password generator.",
		}
		schemaDef.Attributes["maximum_passwords_per_request"] = schema.Int64Attribute{
			Description: "The maximum number of passwords that may be generated and returned to the client for a single request.",
		}
		schemaDef.Attributes["maximum_validation_attempts_per_password"] = schema.Int64Attribute{
			Description: "The maximum number of attempts that the server may use to generate a password that passes validation.",
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *extendedOperationHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var model extendedOperationHandlerResourceModel
	req.Plan.Get(ctx, &model)
	resourceType := model.Type.ValueString()
	// Set defaults for validate-totp-password type
	if resourceType == "validate-totp-password" {
		if !internaltypes.IsDefined(model.SharedSecretAttributeType) {
			model.SharedSecretAttributeType = types.StringValue("ds-auth-totp-shared-secret")
		}
		if !internaltypes.IsDefined(model.AdjacentIntervalsToCheck) {
			model.AdjacentIntervalsToCheck = types.Int64Value(2)
		}
		if !internaltypes.IsDefined(model.PreventTOTPReuse) {
			model.PreventTOTPReuse = types.BoolValue(false)
		}
	}
	// Set defaults for replace-certificate type
	if resourceType == "replace-certificate" {
		if !internaltypes.IsDefined(model.AllowRemotelyProvidedCertificates) {
			model.AllowRemotelyProvidedCertificates = types.BoolValue(false)
		}
	}
	resp.Plan.Set(ctx, &model)
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsExtendedOperationHandler() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("shared_secret_attribute_type"),
			path.MatchRoot("type"),
			[]string{"validate-totp-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("time_interval_duration"),
			path.MatchRoot("type"),
			[]string{"validate-totp-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("adjacent_intervals_to_check"),
			path.MatchRoot("type"),
			[]string{"validate-totp-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("prevent_totp_reuse"),
			path.MatchRoot("type"),
			[]string{"validate-totp-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_remotely_provided_certificates"),
			path.MatchRoot("type"),
			[]string{"replace-certificate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_operation"),
			path.MatchRoot("type"),
			[]string{"replace-certificate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("connection_criteria"),
			path.MatchRoot("type"),
			[]string{"replace-certificate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_criteria"),
			path.MatchRoot("type"),
			[]string{"replace-certificate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("identity_mapper"),
			path.MatchRoot("type"),
			[]string{"password-modify", "deliver-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password_generator"),
			path.MatchRoot("type"),
			[]string{"single-use-tokens", "deliver-password-reset-token", "deliver-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_otp_delivery_mechanism"),
			path.MatchRoot("type"),
			[]string{"single-use-tokens", "deliver-otp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_single_use_token_validity_duration"),
			path.MatchRoot("type"),
			[]string{"single-use-tokens"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_token_delivery_mechanism"),
			path.MatchRoot("type"),
			[]string{"deliver-password-reset-token"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password_reset_token_validity_duration"),
			path.MatchRoot("type"),
			[]string{"deliver-password-reset-token"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"validate-totp-password",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"replace-certificate",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"collect-support-data",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"export-reversible-passwords",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"single-use-tokens",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("password_generator"), path.MatchRoot("default_otp_delivery_mechanism")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"deliver-password-reset-token",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("password_generator"), path.MatchRoot("default_token_delivery_mechanism")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"deliver-otp",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("identity_mapper"), path.MatchRoot("password_generator"), path.MatchRoot("default_otp_delivery_mechanism")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r extendedOperationHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsExtendedOperationHandler()
}

// Add config validators
func (r defaultExtendedOperationHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_password_policy"),
			path.MatchRoot("type"),
			[]string{"generate-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_password_generator"),
			path.MatchRoot("type"),
			[]string{"generate-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_passwords_per_request"),
			path.MatchRoot("type"),
			[]string{"generate-password"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("maximum_validation_attempts_per_password"),
			path.MatchRoot("type"),
			[]string{"generate-password"},
		),
	}
	return append(configValidatorsExtendedOperationHandler(), validators...)
}

// Add optional fields to create request for validate-totp-password extended-operation-handler
func addOptionalValidateTotpPasswordExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddValidateTotpPasswordExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SharedSecretAttributeType) {
		addRequest.SharedSecretAttributeType = plan.SharedSecretAttributeType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeIntervalDuration) {
		addRequest.TimeIntervalDuration = plan.TimeIntervalDuration.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AdjacentIntervalsToCheck) {
		addRequest.AdjacentIntervalsToCheck = plan.AdjacentIntervalsToCheck.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.PreventTOTPReuse) {
		addRequest.PreventTOTPReuse = plan.PreventTOTPReuse.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for replace-certificate extended-operation-handler
func addOptionalReplaceCertificateExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddReplaceCertificateExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.AllowRemotelyProvidedCertificates) {
		addRequest.AllowRemotelyProvidedCertificates = plan.AllowRemotelyProvidedCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowedOperation) {
		var slice []string
		plan.AllowedOperation.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumextendedOperationHandlerAllowedOperationProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumextendedOperationHandlerAllowedOperationPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedOperation = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		addRequest.ConnectionCriteria = plan.ConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		addRequest.RequestCriteria = plan.RequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for collect-support-data extended-operation-handler
func addOptionalCollectSupportDataExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddCollectSupportDataExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for export-reversible-passwords extended-operation-handler
func addOptionalExportReversiblePasswordsExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddExportReversiblePasswordsExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for single-use-tokens extended-operation-handler
func addOptionalSingleUseTokensExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddSingleUseTokensExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultSingleUseTokenValidityDuration) {
		addRequest.DefaultSingleUseTokenValidityDuration = plan.DefaultSingleUseTokenValidityDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for deliver-password-reset-token extended-operation-handler
func addOptionalDeliverPasswordResetTokenExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddDeliverPasswordResetTokenExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordResetTokenValidityDuration) {
		addRequest.PasswordResetTokenValidityDuration = plan.PasswordResetTokenValidityDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for deliver-otp extended-operation-handler
func addOptionalDeliverOtpExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddDeliverOtpExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party extended-operation-handler
func addOptionalThirdPartyExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyExtendedOperationHandlerRequest, plan extendedOperationHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateExtendedOperationHandlerUnknownValues(model *extendedOperationHandlerResourceModel) {
	if model.DefaultTokenDeliveryMechanism.IsUnknown() || model.DefaultTokenDeliveryMechanism.IsNull() {
		model.DefaultTokenDeliveryMechanism, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedOperation.IsUnknown() || model.AllowedOperation.IsNull() {
		model.AllowedOperation, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DefaultOTPDeliveryMechanism.IsUnknown() || model.DefaultOTPDeliveryMechanism.IsNull() {
		model.DefaultOTPDeliveryMechanism, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TimeIntervalDuration.IsUnknown() || model.TimeIntervalDuration.IsNull() {
		model.TimeIntervalDuration = types.StringValue("")
	}
	if model.PasswordResetTokenValidityDuration.IsUnknown() || model.PasswordResetTokenValidityDuration.IsNull() {
		model.PasswordResetTokenValidityDuration = types.StringValue("")
	}
	if model.SharedSecretAttributeType.IsUnknown() || model.SharedSecretAttributeType.IsNull() {
		model.SharedSecretAttributeType = types.StringValue("")
	}
	if model.DefaultSingleUseTokenValidityDuration.IsUnknown() || model.DefaultSingleUseTokenValidityDuration.IsNull() {
		model.DefaultSingleUseTokenValidityDuration = types.StringValue("")
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateExtendedOperationHandlerUnknownValuesDefault(model *defaultExtendedOperationHandlerResourceModel) {
	if model.DefaultTokenDeliveryMechanism.IsUnknown() || model.DefaultTokenDeliveryMechanism.IsNull() {
		model.DefaultTokenDeliveryMechanism, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedOperation.IsUnknown() || model.AllowedOperation.IsNull() {
		model.AllowedOperation, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DefaultOTPDeliveryMechanism.IsUnknown() || model.DefaultOTPDeliveryMechanism.IsNull() {
		model.DefaultOTPDeliveryMechanism, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TimeIntervalDuration.IsUnknown() || model.TimeIntervalDuration.IsNull() {
		model.TimeIntervalDuration = types.StringValue("")
	}
	if model.PasswordResetTokenValidityDuration.IsUnknown() || model.PasswordResetTokenValidityDuration.IsNull() {
		model.PasswordResetTokenValidityDuration = types.StringValue("")
	}
	if model.RequestCriteria.IsUnknown() || model.RequestCriteria.IsNull() {
		model.RequestCriteria = types.StringValue("")
	}
	if model.SharedSecretAttributeType.IsUnknown() || model.SharedSecretAttributeType.IsNull() {
		model.SharedSecretAttributeType = types.StringValue("")
	}
	if model.IdentityMapper.IsUnknown() || model.IdentityMapper.IsNull() {
		model.IdentityMapper = types.StringValue("")
	}
	if model.DefaultSingleUseTokenValidityDuration.IsUnknown() || model.DefaultSingleUseTokenValidityDuration.IsNull() {
		model.DefaultSingleUseTokenValidityDuration = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.PasswordGenerator.IsUnknown() || model.PasswordGenerator.IsNull() {
		model.PasswordGenerator = types.StringValue("")
	}
	if model.ConnectionCriteria.IsUnknown() || model.ConnectionCriteria.IsNull() {
		model.ConnectionCriteria = types.StringValue("")
	}
	if model.DefaultPasswordPolicy.IsUnknown() || model.DefaultPasswordPolicy.IsNull() {
		model.DefaultPasswordPolicy = types.StringValue("")
	}
	if model.DefaultPasswordGenerator.IsUnknown() || model.DefaultPasswordGenerator.IsNull() {
		model.DefaultPasswordGenerator = types.StringValue("")
	}
}

// Read a CancelExtendedOperationHandlerResponse object into the model struct
func readCancelExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.CancelExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("cancel")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a ValidateTotpPasswordExtendedOperationHandlerResponse object into the model struct
func readValidateTotpPasswordExtendedOperationHandlerResponse(ctx context.Context, r *client.ValidateTotpPasswordExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("validate-totp-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SharedSecretAttributeType = internaltypes.StringTypeOrNil(r.SharedSecretAttributeType, true)
	state.TimeIntervalDuration = internaltypes.StringTypeOrNil(r.TimeIntervalDuration, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval_duration",
		expectedValues.TimeIntervalDuration, state.TimeIntervalDuration, diagnostics)
	state.AdjacentIntervalsToCheck = internaltypes.Int64TypeOrNil(r.AdjacentIntervalsToCheck)
	state.PreventTOTPReuse = internaltypes.BoolTypeOrNil(r.PreventTOTPReuse)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a ValidateTotpPasswordExtendedOperationHandlerResponse object into the model struct
func readValidateTotpPasswordExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.ValidateTotpPasswordExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("validate-totp-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SharedSecretAttributeType = internaltypes.StringTypeOrNil(r.SharedSecretAttributeType, true)
	state.TimeIntervalDuration = internaltypes.StringTypeOrNil(r.TimeIntervalDuration, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval_duration",
		expectedValues.TimeIntervalDuration, state.TimeIntervalDuration, diagnostics)
	state.AdjacentIntervalsToCheck = internaltypes.Int64TypeOrNil(r.AdjacentIntervalsToCheck)
	state.PreventTOTPReuse = internaltypes.BoolTypeOrNil(r.PreventTOTPReuse)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a ReplaceCertificateExtendedOperationHandlerResponse object into the model struct
func readReplaceCertificateExtendedOperationHandlerResponse(ctx context.Context, r *client.ReplaceCertificateExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replace-certificate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowRemotelyProvidedCertificates = internaltypes.BoolTypeOrNil(r.AllowRemotelyProvidedCertificates)
	state.AllowedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumextendedOperationHandlerAllowedOperationProp(r.AllowedOperation))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a ReplaceCertificateExtendedOperationHandlerResponse object into the model struct
func readReplaceCertificateExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.ReplaceCertificateExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("replace-certificate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowRemotelyProvidedCertificates = internaltypes.BoolTypeOrNil(r.AllowRemotelyProvidedCertificates)
	state.AllowedOperation = internaltypes.GetStringSet(
		client.StringSliceEnumextendedOperationHandlerAllowedOperationProp(r.AllowedOperation))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, true)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a GetConnectionIdExtendedOperationHandlerResponse object into the model struct
func readGetConnectionIdExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.GetConnectionIdExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-connection-id")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a MultiUpdateExtendedOperationHandlerResponse object into the model struct
func readMultiUpdateExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.MultiUpdateExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("multi-update")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a NotificationSubscriptionExtendedOperationHandlerResponse object into the model struct
func readNotificationSubscriptionExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.NotificationSubscriptionExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("notification-subscription")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a PasswordModifyExtendedOperationHandlerResponse object into the model struct
func readPasswordModifyExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.PasswordModifyExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-modify")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a CustomExtendedOperationHandlerResponse object into the model struct
func readCustomExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.CustomExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a CollectSupportDataExtendedOperationHandlerResponse object into the model struct
func readCollectSupportDataExtendedOperationHandlerResponse(ctx context.Context, r *client.CollectSupportDataExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("collect-support-data")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a CollectSupportDataExtendedOperationHandlerResponse object into the model struct
func readCollectSupportDataExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.CollectSupportDataExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("collect-support-data")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a ExportReversiblePasswordsExtendedOperationHandlerResponse object into the model struct
func readExportReversiblePasswordsExtendedOperationHandlerResponse(ctx context.Context, r *client.ExportReversiblePasswordsExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("export-reversible-passwords")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a ExportReversiblePasswordsExtendedOperationHandlerResponse object into the model struct
func readExportReversiblePasswordsExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.ExportReversiblePasswordsExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("export-reversible-passwords")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a BatchedTransactionsExtendedOperationHandlerResponse object into the model struct
func readBatchedTransactionsExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.BatchedTransactionsExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("batched-transactions")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a GetChangelogBatchExtendedOperationHandlerResponse object into the model struct
func readGetChangelogBatchExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.GetChangelogBatchExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-changelog-batch")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse object into the model struct
func readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-supported-otp-delivery-mechanisms")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a SingleUseTokensExtendedOperationHandlerResponse object into the model struct
func readSingleUseTokensExtendedOperationHandlerResponse(ctx context.Context, r *client.SingleUseTokensExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("single-use-tokens")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.DefaultSingleUseTokenValidityDuration = internaltypes.StringTypeOrNil(r.DefaultSingleUseTokenValidityDuration, true)
	config.CheckMismatchedPDFormattedAttributes("default_single_use_token_validity_duration",
		expectedValues.DefaultSingleUseTokenValidityDuration, state.DefaultSingleUseTokenValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a SingleUseTokensExtendedOperationHandlerResponse object into the model struct
func readSingleUseTokensExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.SingleUseTokensExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("single-use-tokens")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.DefaultSingleUseTokenValidityDuration = internaltypes.StringTypeOrNil(r.DefaultSingleUseTokenValidityDuration, true)
	config.CheckMismatchedPDFormattedAttributes("default_single_use_token_validity_duration",
		expectedValues.DefaultSingleUseTokenValidityDuration, state.DefaultSingleUseTokenValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a GeneratePasswordExtendedOperationHandlerResponse object into the model struct
func readGeneratePasswordExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.GeneratePasswordExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generate-password")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DefaultPasswordPolicy = internaltypes.StringTypeOrNil(r.DefaultPasswordPolicy, true)
	state.DefaultPasswordGenerator = types.StringValue(r.DefaultPasswordGenerator)
	state.MaximumPasswordsPerRequest = internaltypes.Int64TypeOrNil(r.MaximumPasswordsPerRequest)
	state.MaximumValidationAttemptsPerPassword = internaltypes.Int64TypeOrNil(r.MaximumValidationAttemptsPerPassword)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a WhoAmIExtendedOperationHandlerResponse object into the model struct
func readWhoAmIExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.WhoAmIExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("who-am-i")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a StartTlsExtendedOperationHandlerResponse object into the model struct
func readStartTlsExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.StartTlsExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("start-tls")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a DeliverPasswordResetTokenExtendedOperationHandlerResponse object into the model struct
func readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx context.Context, r *client.DeliverPasswordResetTokenExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deliver-password-reset-token")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultTokenDeliveryMechanism = internaltypes.GetStringSet(r.DefaultTokenDeliveryMechanism)
	state.PasswordResetTokenValidityDuration = types.StringValue(r.PasswordResetTokenValidityDuration)
	config.CheckMismatchedPDFormattedAttributes("password_reset_token_validity_duration",
		expectedValues.PasswordResetTokenValidityDuration, state.PasswordResetTokenValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a DeliverPasswordResetTokenExtendedOperationHandlerResponse object into the model struct
func readDeliverPasswordResetTokenExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.DeliverPasswordResetTokenExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deliver-password-reset-token")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultTokenDeliveryMechanism = internaltypes.GetStringSet(r.DefaultTokenDeliveryMechanism)
	state.PasswordResetTokenValidityDuration = types.StringValue(r.PasswordResetTokenValidityDuration)
	config.CheckMismatchedPDFormattedAttributes("password_reset_token_validity_duration",
		expectedValues.PasswordResetTokenValidityDuration, state.PasswordResetTokenValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a PasswordPolicyStateExtendedOperationHandlerResponse object into the model struct
func readPasswordPolicyStateExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.PasswordPolicyStateExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("password-policy-state")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a GetPasswordQualityRequirementsExtendedOperationHandlerResponse object into the model struct
func readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("get-password-quality-requirements")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a DeliverOtpExtendedOperationHandlerResponse object into the model struct
func readDeliverOtpExtendedOperationHandlerResponse(ctx context.Context, r *client.DeliverOtpExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deliver-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a DeliverOtpExtendedOperationHandlerResponse object into the model struct
func readDeliverOtpExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.DeliverOtpExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("deliver-otp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Read a ThirdPartyExtendedOperationHandlerResponse object into the model struct
func readThirdPartyExtendedOperationHandlerResponse(ctx context.Context, r *client.ThirdPartyExtendedOperationHandlerResponse, state *extendedOperationHandlerResourceModel, expectedValues *extendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValues(state)
}

// Read a ThirdPartyExtendedOperationHandlerResponse object into the model struct
func readThirdPartyExtendedOperationHandlerResponseDefault(ctx context.Context, r *client.ThirdPartyExtendedOperationHandlerResponse, state *defaultExtendedOperationHandlerResourceModel, expectedValues *defaultExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateExtendedOperationHandlerUnknownValuesDefault(state)
}

// Create any update operations necessary to make the state match the plan
func createExtendedOperationHandlerOperations(plan extendedOperationHandlerResourceModel, state extendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultTokenDeliveryMechanism, state.DefaultTokenDeliveryMechanism, "default-token-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordResetTokenValidityDuration, state.PasswordResetTokenValidityDuration, "password-reset-token-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordGenerator, state.PasswordGenerator, "password-generator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultOTPDeliveryMechanism, state.DefaultOTPDeliveryMechanism, "default-otp-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultSingleUseTokenValidityDuration, state.DefaultSingleUseTokenValidityDuration, "default-single-use-token-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowRemotelyProvidedCertificates, state.AllowRemotelyProvidedCertificates, "allow-remotely-provided-certificates")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedOperation, state.AllowedOperation, "allowed-operation")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SharedSecretAttributeType, state.SharedSecretAttributeType, "shared-secret-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeIntervalDuration, state.TimeIntervalDuration, "time-interval-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.AdjacentIntervalsToCheck, state.AdjacentIntervalsToCheck, "adjacent-intervals-to-check")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventTOTPReuse, state.PreventTOTPReuse, "prevent-totp-reuse")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createExtendedOperationHandlerOperationsDefault(plan defaultExtendedOperationHandlerResourceModel, state defaultExtendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultPasswordPolicy, state.DefaultPasswordPolicy, "default-password-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultTokenDeliveryMechanism, state.DefaultTokenDeliveryMechanism, "default-token-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordResetTokenValidityDuration, state.PasswordResetTokenValidityDuration, "password-reset-token-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultPasswordGenerator, state.DefaultPasswordGenerator, "default-password-generator")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumPasswordsPerRequest, state.MaximumPasswordsPerRequest, "maximum-passwords-per-request")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumValidationAttemptsPerPassword, state.MaximumValidationAttemptsPerPassword, "maximum-validation-attempts-per-password")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordGenerator, state.PasswordGenerator, "password-generator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultOTPDeliveryMechanism, state.DefaultOTPDeliveryMechanism, "default-otp-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultSingleUseTokenValidityDuration, state.DefaultSingleUseTokenValidityDuration, "default-single-use-token-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowRemotelyProvidedCertificates, state.AllowRemotelyProvidedCertificates, "allow-remotely-provided-certificates")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedOperation, state.AllowedOperation, "allowed-operation")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SharedSecretAttributeType, state.SharedSecretAttributeType, "shared-secret-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeIntervalDuration, state.TimeIntervalDuration, "time-interval-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.AdjacentIntervalsToCheck, state.AdjacentIntervalsToCheck, "adjacent-intervals-to-check")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventTOTPReuse, state.PreventTOTPReuse, "prevent-totp-reuse")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a validate-totp-password extended-operation-handler
func (r *extendedOperationHandlerResource) CreateValidateTotpPasswordExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	addRequest := client.NewAddValidateTotpPasswordExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumvalidateTotpPasswordExtendedOperationHandlerSchemaUrn{client.ENUMVALIDATETOTPPASSWORDEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERVALIDATE_TOTP_PASSWORD},
		plan.Enabled.ValueBool())
	err := addOptionalValidateTotpPasswordExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddValidateTotpPasswordExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, addResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a single-use-tokens extended-operation-handler
func (r *extendedOperationHandlerResource) CreateSingleUseTokensExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	var DefaultOTPDeliveryMechanismSlice []string
	plan.DefaultOTPDeliveryMechanism.ElementsAs(ctx, &DefaultOTPDeliveryMechanismSlice, false)
	addRequest := client.NewAddSingleUseTokensExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumsingleUseTokensExtendedOperationHandlerSchemaUrn{client.ENUMSINGLEUSETOKENSEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERSINGLE_USE_TOKENS},
		plan.PasswordGenerator.ValueString(),
		DefaultOTPDeliveryMechanismSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSingleUseTokensExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddSingleUseTokensExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readSingleUseTokensExtendedOperationHandlerResponse(ctx, addResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a deliver-password-reset-token extended-operation-handler
func (r *extendedOperationHandlerResource) CreateDeliverPasswordResetTokenExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	var DefaultTokenDeliveryMechanismSlice []string
	plan.DefaultTokenDeliveryMechanism.ElementsAs(ctx, &DefaultTokenDeliveryMechanismSlice, false)
	addRequest := client.NewAddDeliverPasswordResetTokenExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumdeliverPasswordResetTokenExtendedOperationHandlerSchemaUrn{client.ENUMDELIVERPASSWORDRESETTOKENEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERDELIVER_PASSWORD_RESET_TOKEN},
		plan.PasswordGenerator.ValueString(),
		DefaultTokenDeliveryMechanismSlice,
		plan.Enabled.ValueBool())
	err := addOptionalDeliverPasswordResetTokenExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddDeliverPasswordResetTokenExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, addResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a replace-certificate extended-operation-handler
func (r *extendedOperationHandlerResource) CreateReplaceCertificateExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	addRequest := client.NewAddReplaceCertificateExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumreplaceCertificateExtendedOperationHandlerSchemaUrn{client.ENUMREPLACECERTIFICATEEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERREPLACE_CERTIFICATE},
		plan.Enabled.ValueBool())
	err := addOptionalReplaceCertificateExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddReplaceCertificateExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readReplaceCertificateExtendedOperationHandlerResponse(ctx, addResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a collect-support-data extended-operation-handler
func (r *extendedOperationHandlerResource) CreateCollectSupportDataExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	addRequest := client.NewAddCollectSupportDataExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumcollectSupportDataExtendedOperationHandlerSchemaUrn{client.ENUMCOLLECTSUPPORTDATAEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERCOLLECT_SUPPORT_DATA},
		plan.Enabled.ValueBool())
	err := addOptionalCollectSupportDataExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddCollectSupportDataExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readCollectSupportDataExtendedOperationHandlerResponse(ctx, addResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a export-reversible-passwords extended-operation-handler
func (r *extendedOperationHandlerResource) CreateExportReversiblePasswordsExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	addRequest := client.NewAddExportReversiblePasswordsExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumexportReversiblePasswordsExtendedOperationHandlerSchemaUrn{client.ENUMEXPORTREVERSIBLEPASSWORDSEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLEREXPORT_REVERSIBLE_PASSWORDS},
		plan.Enabled.ValueBool())
	err := addOptionalExportReversiblePasswordsExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddExportReversiblePasswordsExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readExportReversiblePasswordsExtendedOperationHandlerResponse(ctx, addResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a deliver-otp extended-operation-handler
func (r *extendedOperationHandlerResource) CreateDeliverOtpExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	var DefaultOTPDeliveryMechanismSlice []string
	plan.DefaultOTPDeliveryMechanism.ElementsAs(ctx, &DefaultOTPDeliveryMechanismSlice, false)
	addRequest := client.NewAddDeliverOtpExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumdeliverOtpExtendedOperationHandlerSchemaUrn{client.ENUMDELIVEROTPEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERDELIVER_OTP},
		plan.IdentityMapper.ValueString(),
		plan.PasswordGenerator.ValueString(),
		DefaultOTPDeliveryMechanismSlice,
		plan.Enabled.ValueBool())
	err := addOptionalDeliverOtpExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddDeliverOtpExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readDeliverOtpExtendedOperationHandlerResponse(ctx, addResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party extended-operation-handler
func (r *extendedOperationHandlerResource) CreateThirdPartyExtendedOperationHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan extendedOperationHandlerResourceModel) (*extendedOperationHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyExtendedOperationHandlerRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyExtendedOperationHandlerSchemaUrn{client.ENUMTHIRDPARTYEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyExtendedOperationHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Extended Operation Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddThirdPartyExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Extended Operation Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state extendedOperationHandlerResourceModel
	readThirdPartyExtendedOperationHandlerResponse(ctx, addResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *extendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan extendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *extendedOperationHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "validate-totp-password" {
		state, err = r.CreateValidateTotpPasswordExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "single-use-tokens" {
		state, err = r.CreateSingleUseTokensExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "deliver-password-reset-token" {
		state, err = r.CreateDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "replace-certificate" {
		state, err = r.CreateReplaceCertificateExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "collect-support-data" {
		state, err = r.CreateCollectSupportDataExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "export-reversible-passwords" {
		state, err = r.CreateExportReversiblePasswordsExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "deliver-otp" {
		state, err = r.CreateDeliverOtpExtendedOperationHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyExtendedOperationHandler(ctx, req, resp, plan)
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
func (r *defaultExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultExtendedOperationHandlerResourceModel
	if readResponse.CancelExtendedOperationHandlerResponse != nil {
		readCancelExtendedOperationHandlerResponseDefault(ctx, readResponse.CancelExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
		readValidateTotpPasswordExtendedOperationHandlerResponseDefault(ctx, readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReplaceCertificateExtendedOperationHandlerResponse != nil {
		readReplaceCertificateExtendedOperationHandlerResponseDefault(ctx, readResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetConnectionIdExtendedOperationHandlerResponse != nil {
		readGetConnectionIdExtendedOperationHandlerResponseDefault(ctx, readResponse.GetConnectionIdExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MultiUpdateExtendedOperationHandlerResponse != nil {
		readMultiUpdateExtendedOperationHandlerResponseDefault(ctx, readResponse.MultiUpdateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NotificationSubscriptionExtendedOperationHandlerResponse != nil {
		readNotificationSubscriptionExtendedOperationHandlerResponseDefault(ctx, readResponse.NotificationSubscriptionExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordModifyExtendedOperationHandlerResponse != nil {
		readPasswordModifyExtendedOperationHandlerResponseDefault(ctx, readResponse.PasswordModifyExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomExtendedOperationHandlerResponse != nil {
		readCustomExtendedOperationHandlerResponseDefault(ctx, readResponse.CustomExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CollectSupportDataExtendedOperationHandlerResponse != nil {
		readCollectSupportDataExtendedOperationHandlerResponseDefault(ctx, readResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
		readExportReversiblePasswordsExtendedOperationHandlerResponseDefault(ctx, readResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BatchedTransactionsExtendedOperationHandlerResponse != nil {
		readBatchedTransactionsExtendedOperationHandlerResponseDefault(ctx, readResponse.BatchedTransactionsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetChangelogBatchExtendedOperationHandlerResponse != nil {
		readGetChangelogBatchExtendedOperationHandlerResponseDefault(ctx, readResponse.GetChangelogBatchExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse != nil {
		readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDefault(ctx, readResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SingleUseTokensExtendedOperationHandlerResponse != nil {
		readSingleUseTokensExtendedOperationHandlerResponseDefault(ctx, readResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GeneratePasswordExtendedOperationHandlerResponse != nil {
		readGeneratePasswordExtendedOperationHandlerResponseDefault(ctx, readResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.WhoAmIExtendedOperationHandlerResponse != nil {
		readWhoAmIExtendedOperationHandlerResponseDefault(ctx, readResponse.WhoAmIExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StartTlsExtendedOperationHandlerResponse != nil {
		readStartTlsExtendedOperationHandlerResponseDefault(ctx, readResponse.StartTlsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
		readDeliverPasswordResetTokenExtendedOperationHandlerResponseDefault(ctx, readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyStateExtendedOperationHandlerResponse != nil {
		readPasswordPolicyStateExtendedOperationHandlerResponseDefault(ctx, readResponse.PasswordPolicyStateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse != nil {
		readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDefault(ctx, readResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DeliverOtpExtendedOperationHandlerResponse != nil {
		readDeliverOtpExtendedOperationHandlerResponseDefault(ctx, readResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyExtendedOperationHandlerResponse != nil {
		readThirdPartyExtendedOperationHandlerResponseDefault(ctx, readResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createExtendedOperationHandlerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.CancelExtendedOperationHandlerResponse != nil {
			readCancelExtendedOperationHandlerResponseDefault(ctx, updateResponse.CancelExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
			readValidateTotpPasswordExtendedOperationHandlerResponseDefault(ctx, updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplaceCertificateExtendedOperationHandlerResponse != nil {
			readReplaceCertificateExtendedOperationHandlerResponseDefault(ctx, updateResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetConnectionIdExtendedOperationHandlerResponse != nil {
			readGetConnectionIdExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetConnectionIdExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MultiUpdateExtendedOperationHandlerResponse != nil {
			readMultiUpdateExtendedOperationHandlerResponseDefault(ctx, updateResponse.MultiUpdateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NotificationSubscriptionExtendedOperationHandlerResponse != nil {
			readNotificationSubscriptionExtendedOperationHandlerResponseDefault(ctx, updateResponse.NotificationSubscriptionExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordModifyExtendedOperationHandlerResponse != nil {
			readPasswordModifyExtendedOperationHandlerResponseDefault(ctx, updateResponse.PasswordModifyExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomExtendedOperationHandlerResponse != nil {
			readCustomExtendedOperationHandlerResponseDefault(ctx, updateResponse.CustomExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CollectSupportDataExtendedOperationHandlerResponse != nil {
			readCollectSupportDataExtendedOperationHandlerResponseDefault(ctx, updateResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
			readExportReversiblePasswordsExtendedOperationHandlerResponseDefault(ctx, updateResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BatchedTransactionsExtendedOperationHandlerResponse != nil {
			readBatchedTransactionsExtendedOperationHandlerResponseDefault(ctx, updateResponse.BatchedTransactionsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetChangelogBatchExtendedOperationHandlerResponse != nil {
			readGetChangelogBatchExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetChangelogBatchExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse != nil {
			readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SingleUseTokensExtendedOperationHandlerResponse != nil {
			readSingleUseTokensExtendedOperationHandlerResponseDefault(ctx, updateResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GeneratePasswordExtendedOperationHandlerResponse != nil {
			readGeneratePasswordExtendedOperationHandlerResponseDefault(ctx, updateResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.WhoAmIExtendedOperationHandlerResponse != nil {
			readWhoAmIExtendedOperationHandlerResponseDefault(ctx, updateResponse.WhoAmIExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StartTlsExtendedOperationHandlerResponse != nil {
			readStartTlsExtendedOperationHandlerResponseDefault(ctx, updateResponse.StartTlsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
			readDeliverPasswordResetTokenExtendedOperationHandlerResponseDefault(ctx, updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordPolicyStateExtendedOperationHandlerResponse != nil {
			readPasswordPolicyStateExtendedOperationHandlerResponseDefault(ctx, updateResponse.PasswordPolicyStateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse != nil {
			readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeliverOtpExtendedOperationHandlerResponse != nil {
			readDeliverOtpExtendedOperationHandlerResponseDefault(ctx, updateResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyExtendedOperationHandlerResponse != nil {
			readThirdPartyExtendedOperationHandlerResponseDefault(ctx, updateResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *extendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state extendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Extended Operation Handler", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Extended Operation Handler", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
		readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ReplaceCertificateExtendedOperationHandlerResponse != nil {
		readReplaceCertificateExtendedOperationHandlerResponse(ctx, readResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CollectSupportDataExtendedOperationHandlerResponse != nil {
		readCollectSupportDataExtendedOperationHandlerResponse(ctx, readResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
		readExportReversiblePasswordsExtendedOperationHandlerResponse(ctx, readResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SingleUseTokensExtendedOperationHandlerResponse != nil {
		readSingleUseTokensExtendedOperationHandlerResponse(ctx, readResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
		readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DeliverOtpExtendedOperationHandlerResponse != nil {
		readDeliverOtpExtendedOperationHandlerResponse(ctx, readResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyExtendedOperationHandlerResponse != nil {
		readThirdPartyExtendedOperationHandlerResponse(ctx, readResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
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
		readCancelExtendedOperationHandlerResponseDefault(ctx, readResponse.CancelExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetConnectionIdExtendedOperationHandlerResponse != nil {
		readGetConnectionIdExtendedOperationHandlerResponseDefault(ctx, readResponse.GetConnectionIdExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MultiUpdateExtendedOperationHandlerResponse != nil {
		readMultiUpdateExtendedOperationHandlerResponseDefault(ctx, readResponse.MultiUpdateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.NotificationSubscriptionExtendedOperationHandlerResponse != nil {
		readNotificationSubscriptionExtendedOperationHandlerResponseDefault(ctx, readResponse.NotificationSubscriptionExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordModifyExtendedOperationHandlerResponse != nil {
		readPasswordModifyExtendedOperationHandlerResponseDefault(ctx, readResponse.PasswordModifyExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomExtendedOperationHandlerResponse != nil {
		readCustomExtendedOperationHandlerResponseDefault(ctx, readResponse.CustomExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BatchedTransactionsExtendedOperationHandlerResponse != nil {
		readBatchedTransactionsExtendedOperationHandlerResponseDefault(ctx, readResponse.BatchedTransactionsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetChangelogBatchExtendedOperationHandlerResponse != nil {
		readGetChangelogBatchExtendedOperationHandlerResponseDefault(ctx, readResponse.GetChangelogBatchExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse != nil {
		readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDefault(ctx, readResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GeneratePasswordExtendedOperationHandlerResponse != nil {
		readGeneratePasswordExtendedOperationHandlerResponseDefault(ctx, readResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.WhoAmIExtendedOperationHandlerResponse != nil {
		readWhoAmIExtendedOperationHandlerResponseDefault(ctx, readResponse.WhoAmIExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StartTlsExtendedOperationHandlerResponse != nil {
		readStartTlsExtendedOperationHandlerResponseDefault(ctx, readResponse.StartTlsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.PasswordPolicyStateExtendedOperationHandlerResponse != nil {
		readPasswordPolicyStateExtendedOperationHandlerResponseDefault(ctx, readResponse.PasswordPolicyStateExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse != nil {
		readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDefault(ctx, readResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *extendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan extendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state extendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
			readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplaceCertificateExtendedOperationHandlerResponse != nil {
			readReplaceCertificateExtendedOperationHandlerResponse(ctx, updateResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CollectSupportDataExtendedOperationHandlerResponse != nil {
			readCollectSupportDataExtendedOperationHandlerResponse(ctx, updateResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
			readExportReversiblePasswordsExtendedOperationHandlerResponse(ctx, updateResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SingleUseTokensExtendedOperationHandlerResponse != nil {
			readSingleUseTokensExtendedOperationHandlerResponse(ctx, updateResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
			readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeliverOtpExtendedOperationHandlerResponse != nil {
			readDeliverOtpExtendedOperationHandlerResponse(ctx, updateResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyExtendedOperationHandlerResponse != nil {
			readThirdPartyExtendedOperationHandlerResponse(ctx, updateResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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

func (r *defaultExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultExtendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createExtendedOperationHandlerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.CancelExtendedOperationHandlerResponse != nil {
			readCancelExtendedOperationHandlerResponseDefault(ctx, updateResponse.CancelExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse != nil {
			readValidateTotpPasswordExtendedOperationHandlerResponseDefault(ctx, updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ReplaceCertificateExtendedOperationHandlerResponse != nil {
			readReplaceCertificateExtendedOperationHandlerResponseDefault(ctx, updateResponse.ReplaceCertificateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetConnectionIdExtendedOperationHandlerResponse != nil {
			readGetConnectionIdExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetConnectionIdExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MultiUpdateExtendedOperationHandlerResponse != nil {
			readMultiUpdateExtendedOperationHandlerResponseDefault(ctx, updateResponse.MultiUpdateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.NotificationSubscriptionExtendedOperationHandlerResponse != nil {
			readNotificationSubscriptionExtendedOperationHandlerResponseDefault(ctx, updateResponse.NotificationSubscriptionExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordModifyExtendedOperationHandlerResponse != nil {
			readPasswordModifyExtendedOperationHandlerResponseDefault(ctx, updateResponse.PasswordModifyExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomExtendedOperationHandlerResponse != nil {
			readCustomExtendedOperationHandlerResponseDefault(ctx, updateResponse.CustomExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CollectSupportDataExtendedOperationHandlerResponse != nil {
			readCollectSupportDataExtendedOperationHandlerResponseDefault(ctx, updateResponse.CollectSupportDataExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse != nil {
			readExportReversiblePasswordsExtendedOperationHandlerResponseDefault(ctx, updateResponse.ExportReversiblePasswordsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BatchedTransactionsExtendedOperationHandlerResponse != nil {
			readBatchedTransactionsExtendedOperationHandlerResponseDefault(ctx, updateResponse.BatchedTransactionsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetChangelogBatchExtendedOperationHandlerResponse != nil {
			readGetChangelogBatchExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetChangelogBatchExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse != nil {
			readGetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetSupportedOtpDeliveryMechanismsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SingleUseTokensExtendedOperationHandlerResponse != nil {
			readSingleUseTokensExtendedOperationHandlerResponseDefault(ctx, updateResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GeneratePasswordExtendedOperationHandlerResponse != nil {
			readGeneratePasswordExtendedOperationHandlerResponseDefault(ctx, updateResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.WhoAmIExtendedOperationHandlerResponse != nil {
			readWhoAmIExtendedOperationHandlerResponseDefault(ctx, updateResponse.WhoAmIExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StartTlsExtendedOperationHandlerResponse != nil {
			readStartTlsExtendedOperationHandlerResponseDefault(ctx, updateResponse.StartTlsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse != nil {
			readDeliverPasswordResetTokenExtendedOperationHandlerResponseDefault(ctx, updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.PasswordPolicyStateExtendedOperationHandlerResponse != nil {
			readPasswordPolicyStateExtendedOperationHandlerResponseDefault(ctx, updateResponse.PasswordPolicyStateExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse != nil {
			readGetPasswordQualityRequirementsExtendedOperationHandlerResponseDefault(ctx, updateResponse.GetPasswordQualityRequirementsExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DeliverOtpExtendedOperationHandlerResponse != nil {
			readDeliverOtpExtendedOperationHandlerResponseDefault(ctx, updateResponse.DeliverOtpExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyExtendedOperationHandlerResponse != nil {
			readThirdPartyExtendedOperationHandlerResponseDefault(ctx, updateResponse.ThirdPartyExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *extendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state extendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandlerExecute(r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Extended Operation Handler", err, httpResp)
		return
	}
}

func (r *extendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importExtendedOperationHandler(ctx, req, resp)
}

func (r *defaultExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importExtendedOperationHandler(ctx, req, resp)
}

func importExtendedOperationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
