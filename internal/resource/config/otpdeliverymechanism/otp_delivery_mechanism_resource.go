package otpdeliverymechanism

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &otpDeliveryMechanismResource{}
	_ resource.ResourceWithConfigure   = &otpDeliveryMechanismResource{}
	_ resource.ResourceWithImportState = &otpDeliveryMechanismResource{}
	_ resource.Resource                = &defaultOtpDeliveryMechanismResource{}
	_ resource.ResourceWithConfigure   = &defaultOtpDeliveryMechanismResource{}
	_ resource.ResourceWithImportState = &defaultOtpDeliveryMechanismResource{}
)

// Create a Otp Delivery Mechanism resource
func NewOtpDeliveryMechanismResource() resource.Resource {
	return &otpDeliveryMechanismResource{}
}

func NewDefaultOtpDeliveryMechanismResource() resource.Resource {
	return &defaultOtpDeliveryMechanismResource{}
}

// otpDeliveryMechanismResource is the resource implementation.
type otpDeliveryMechanismResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultOtpDeliveryMechanismResource is the resource implementation.
type defaultOtpDeliveryMechanismResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *otpDeliveryMechanismResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_otp_delivery_mechanism"
}

func (r *defaultOtpDeliveryMechanismResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_otp_delivery_mechanism"
}

// Configure adds the provider configured client to the resource.
func (r *otpDeliveryMechanismResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultOtpDeliveryMechanismResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type otpDeliveryMechanismResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	EmailAddressAttributeType         types.String `tfsdk:"email_address_attribute_type"`
	EmailAddressJSONField             types.String `tfsdk:"email_address_json_field"`
	EmailAddressJSONObjectFilter      types.String `tfsdk:"email_address_json_object_filter"`
	SenderAddress                     types.String `tfsdk:"sender_address"`
	MessageSubject                    types.String `tfsdk:"message_subject"`
	HttpProxyExternalServer           types.String `tfsdk:"http_proxy_external_server"`
	TwilioAccountSID                  types.String `tfsdk:"twilio_account_sid"`
	TwilioAuthToken                   types.String `tfsdk:"twilio_auth_token"`
	TwilioAuthTokenPassphraseProvider types.String `tfsdk:"twilio_auth_token_passphrase_provider"`
	PhoneNumberAttributeType          types.String `tfsdk:"phone_number_attribute_type"`
	PhoneNumberJSONField              types.String `tfsdk:"phone_number_json_field"`
	PhoneNumberJSONObjectFilter       types.String `tfsdk:"phone_number_json_object_filter"`
	SenderPhoneNumber                 types.Set    `tfsdk:"sender_phone_number"`
	MessageTextBeforeOTP              types.String `tfsdk:"message_text_before_otp"`
	MessageTextAfterOTP               types.String `tfsdk:"message_text_after_otp"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *otpDeliveryMechanismResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	otpDeliveryMechanismSchema(ctx, req, resp, false)
}

func (r *defaultOtpDeliveryMechanismResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	otpDeliveryMechanismSchema(ctx, req, resp, true)
}

func otpDeliveryMechanismSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Otp Delivery Mechanism.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of OTP Delivery Mechanism resource. Options are ['twilio', 'email', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"twilio", "email", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party OTP Delivery Mechanism.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party OTP Delivery Mechanism. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"email_address_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute that holds the email address to which the message should be sent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"sender_address": schema.StringAttribute{
				Description: "The e-mail address to use as the sender for the one-time password.",
				Optional:    true,
			},
			"message_subject": schema.StringAttribute{
				Description: "The subject to use for the e-mail message.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Twilio service.",
				Optional:    true,
			},
			"twilio_account_sid": schema.StringAttribute{
				Description: "The unique identifier assigned to the Twilio account that will be used.",
				Optional:    true,
			},
			"twilio_auth_token": schema.StringAttribute{
				Description: "The auth token for the Twilio account that will be used.",
				Optional:    true,
				Sensitive:   true,
			},
			"twilio_auth_token_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the auth token for the Twilio account that will be used.",
				Optional:    true,
			},
			"phone_number_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute in the user's entry that holds the phone number to which the message should be sent.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"phone_number_json_field": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the phone number to which the message should be sent. The phone number must be contained in a top-level field whose value is a single string.",
				Optional:    true,
			},
			"phone_number_json_object_filter": schema.StringAttribute{
				Description: "A JSON object filter that may be used to identify which phone number value to use when sending the message.",
				Optional:    true,
			},
			"sender_phone_number": schema.SetAttribute{
				Description: "The outgoing phone number to use for the messages. Values must be phone numbers you have obtained for use with your Twilio account.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"message_text_before_otp": schema.StringAttribute{
				Description: "Any text that should appear in the message before the one-time password value.",
				Optional:    true,
			},
			"message_text_after_otp": schema.StringAttribute{
				Description: "Any text that should appear in the message after the one-time password value.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this OTP Delivery Mechanism",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this OTP Delivery Mechanism is enabled for use in the server.",
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
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *otpDeliveryMechanismResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanOtpDeliveryMechanism(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOtpDeliveryMechanismResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanOtpDeliveryMechanism(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanOtpDeliveryMechanism(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model otpDeliveryMechanismResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsOtpDeliveryMechanism() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"twilio"},
			resourcevalidator.ExactlyOneOf(
				path.MatchRoot("twilio_auth_token"),
				path.MatchRoot("twilio_auth_token_passphrase_provider"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_proxy_external_server"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("twilio_account_sid"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("twilio_auth_token"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("twilio_auth_token_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("phone_number_attribute_type"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("phone_number_json_field"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("phone_number_json_object_filter"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sender_phone_number"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("message_text_before_otp"),
			path.MatchRoot("type"),
			[]string{"twilio", "email"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("message_text_after_otp"),
			path.MatchRoot("type"),
			[]string{"twilio", "email"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("email_address_attribute_type"),
			path.MatchRoot("type"),
			[]string{"email"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("email_address_json_field"),
			path.MatchRoot("type"),
			[]string{"email"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("email_address_json_object_filter"),
			path.MatchRoot("type"),
			[]string{"email"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sender_address"),
			path.MatchRoot("type"),
			[]string{"email"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("message_subject"),
			path.MatchRoot("type"),
			[]string{"email"},
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
	}
}

// Add config validators
func (r otpDeliveryMechanismResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsOtpDeliveryMechanism()
}

// Add config validators
func (r defaultOtpDeliveryMechanismResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsOtpDeliveryMechanism()
}

// Add optional fields to create request for twilio otp-delivery-mechanism
func addOptionalTwilioOtpDeliveryMechanismFields(ctx context.Context, addRequest *client.AddTwilioOtpDeliveryMechanismRequest, plan otpDeliveryMechanismResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TwilioAuthToken) {
		addRequest.TwilioAuthToken = plan.TwilioAuthToken.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TwilioAuthTokenPassphraseProvider) {
		addRequest.TwilioAuthTokenPassphraseProvider = plan.TwilioAuthTokenPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PhoneNumberAttributeType) {
		addRequest.PhoneNumberAttributeType = plan.PhoneNumberAttributeType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PhoneNumberJSONField) {
		addRequest.PhoneNumberJSONField = plan.PhoneNumberJSONField.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PhoneNumberJSONObjectFilter) {
		addRequest.PhoneNumberJSONObjectFilter = plan.PhoneNumberJSONObjectFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageTextBeforeOTP) {
		addRequest.MessageTextBeforeOTP = plan.MessageTextBeforeOTP.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageTextAfterOTP) {
		addRequest.MessageTextAfterOTP = plan.MessageTextAfterOTP.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for email otp-delivery-mechanism
func addOptionalEmailOtpDeliveryMechanismFields(ctx context.Context, addRequest *client.AddEmailOtpDeliveryMechanismRequest, plan otpDeliveryMechanismResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressAttributeType) {
		addRequest.EmailAddressAttributeType = plan.EmailAddressAttributeType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressJSONField) {
		addRequest.EmailAddressJSONField = plan.EmailAddressJSONField.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EmailAddressJSONObjectFilter) {
		addRequest.EmailAddressJSONObjectFilter = plan.EmailAddressJSONObjectFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageSubject) {
		addRequest.MessageSubject = plan.MessageSubject.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageTextBeforeOTP) {
		addRequest.MessageTextBeforeOTP = plan.MessageTextBeforeOTP.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageTextAfterOTP) {
		addRequest.MessageTextAfterOTP = plan.MessageTextAfterOTP.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party otp-delivery-mechanism
func addOptionalThirdPartyOtpDeliveryMechanismFields(ctx context.Context, addRequest *client.AddThirdPartyOtpDeliveryMechanismRequest, plan otpDeliveryMechanismResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateOtpDeliveryMechanismUnknownValues(model *otpDeliveryMechanismResourceModel) {
	if model.SenderPhoneNumber.IsUnknown() || model.SenderPhoneNumber.IsNull() {
		model.SenderPhoneNumber, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MessageSubject.IsUnknown() || model.MessageSubject.IsNull() {
		model.MessageSubject = types.StringValue("")
	}
	if model.PhoneNumberAttributeType.IsUnknown() || model.PhoneNumberAttributeType.IsNull() {
		model.PhoneNumberAttributeType = types.StringValue("")
	}
	if model.EmailAddressAttributeType.IsUnknown() || model.EmailAddressAttributeType.IsNull() {
		model.EmailAddressAttributeType = types.StringValue("")
	}
	if model.TwilioAuthToken.IsUnknown() {
		model.TwilioAuthToken = types.StringNull()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *otpDeliveryMechanismResourceModel) populateAllComputedStringAttributes() {
	if model.TwilioAccountSID.IsUnknown() || model.TwilioAccountSID.IsNull() {
		model.TwilioAccountSID = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.SenderAddress.IsUnknown() || model.SenderAddress.IsNull() {
		model.SenderAddress = types.StringValue("")
	}
	if model.MessageTextAfterOTP.IsUnknown() || model.MessageTextAfterOTP.IsNull() {
		model.MessageTextAfterOTP = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.HttpProxyExternalServer.IsUnknown() || model.HttpProxyExternalServer.IsNull() {
		model.HttpProxyExternalServer = types.StringValue("")
	}
	if model.MessageTextBeforeOTP.IsUnknown() || model.MessageTextBeforeOTP.IsNull() {
		model.MessageTextBeforeOTP = types.StringValue("")
	}
	if model.PhoneNumberJSONField.IsUnknown() || model.PhoneNumberJSONField.IsNull() {
		model.PhoneNumberJSONField = types.StringValue("")
	}
	if model.PhoneNumberJSONObjectFilter.IsUnknown() || model.PhoneNumberJSONObjectFilter.IsNull() {
		model.PhoneNumberJSONObjectFilter = types.StringValue("")
	}
	if model.TwilioAuthTokenPassphraseProvider.IsUnknown() || model.TwilioAuthTokenPassphraseProvider.IsNull() {
		model.TwilioAuthTokenPassphraseProvider = types.StringValue("")
	}
	if model.EmailAddressJSONField.IsUnknown() || model.EmailAddressJSONField.IsNull() {
		model.EmailAddressJSONField = types.StringValue("")
	}
	if model.EmailAddressJSONObjectFilter.IsUnknown() || model.EmailAddressJSONObjectFilter.IsNull() {
		model.EmailAddressJSONObjectFilter = types.StringValue("")
	}
}

// Read a TwilioOtpDeliveryMechanismResponse object into the model struct
func readTwilioOtpDeliveryMechanismResponse(ctx context.Context, r *client.TwilioOtpDeliveryMechanismResponse, state *otpDeliveryMechanismResourceModel, expectedValues *otpDeliveryMechanismResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("twilio")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.TwilioAccountSID = types.StringValue(r.TwilioAccountSID)
	state.TwilioAuthTokenPassphraseProvider = internaltypes.StringTypeOrNil(r.TwilioAuthTokenPassphraseProvider, internaltypes.IsEmptyString(expectedValues.TwilioAuthTokenPassphraseProvider))
	state.PhoneNumberAttributeType = types.StringValue(r.PhoneNumberAttributeType)
	state.PhoneNumberJSONField = internaltypes.StringTypeOrNil(r.PhoneNumberJSONField, internaltypes.IsEmptyString(expectedValues.PhoneNumberJSONField))
	state.PhoneNumberJSONObjectFilter = internaltypes.StringTypeOrNil(r.PhoneNumberJSONObjectFilter, internaltypes.IsEmptyString(expectedValues.PhoneNumberJSONObjectFilter))
	state.SenderPhoneNumber = internaltypes.GetStringSet(r.SenderPhoneNumber)
	state.MessageTextBeforeOTP = internaltypes.StringTypeOrNil(r.MessageTextBeforeOTP, internaltypes.IsEmptyString(expectedValues.MessageTextBeforeOTP))
	state.MessageTextAfterOTP = internaltypes.StringTypeOrNil(r.MessageTextAfterOTP, internaltypes.IsEmptyString(expectedValues.MessageTextAfterOTP))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateOtpDeliveryMechanismUnknownValues(state)
}

// Read a EmailOtpDeliveryMechanismResponse object into the model struct
func readEmailOtpDeliveryMechanismResponse(ctx context.Context, r *client.EmailOtpDeliveryMechanismResponse, state *otpDeliveryMechanismResourceModel, expectedValues *otpDeliveryMechanismResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("email")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EmailAddressAttributeType = types.StringValue(r.EmailAddressAttributeType)
	state.EmailAddressJSONField = internaltypes.StringTypeOrNil(r.EmailAddressJSONField, internaltypes.IsEmptyString(expectedValues.EmailAddressJSONField))
	state.EmailAddressJSONObjectFilter = internaltypes.StringTypeOrNil(r.EmailAddressJSONObjectFilter, internaltypes.IsEmptyString(expectedValues.EmailAddressJSONObjectFilter))
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.MessageSubject = types.StringValue(r.MessageSubject)
	state.MessageTextBeforeOTP = internaltypes.StringTypeOrNil(r.MessageTextBeforeOTP, internaltypes.IsEmptyString(expectedValues.MessageTextBeforeOTP))
	state.MessageTextAfterOTP = internaltypes.StringTypeOrNil(r.MessageTextAfterOTP, internaltypes.IsEmptyString(expectedValues.MessageTextAfterOTP))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateOtpDeliveryMechanismUnknownValues(state)
}

// Read a ThirdPartyOtpDeliveryMechanismResponse object into the model struct
func readThirdPartyOtpDeliveryMechanismResponse(ctx context.Context, r *client.ThirdPartyOtpDeliveryMechanismResponse, state *otpDeliveryMechanismResourceModel, expectedValues *otpDeliveryMechanismResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateOtpDeliveryMechanismUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *otpDeliveryMechanismResourceModel) setStateValuesNotReturnedByAPI(expectedValues *otpDeliveryMechanismResourceModel) {
	if !expectedValues.TwilioAuthToken.IsUnknown() {
		state.TwilioAuthToken = expectedValues.TwilioAuthToken
	}
}

// Create any update operations necessary to make the state match the plan
func createOtpDeliveryMechanismOperations(plan otpDeliveryMechanismResourceModel, state otpDeliveryMechanismResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.EmailAddressAttributeType, state.EmailAddressAttributeType, "email-address-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.EmailAddressJSONField, state.EmailAddressJSONField, "email-address-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.EmailAddressJSONObjectFilter, state.EmailAddressJSONObjectFilter, "email-address-json-object-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.SenderAddress, state.SenderAddress, "sender-address")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageSubject, state.MessageSubject, "message-subject")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAccountSID, state.TwilioAccountSID, "twilio-account-sid")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthToken, state.TwilioAuthToken, "twilio-auth-token")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthTokenPassphraseProvider, state.TwilioAuthTokenPassphraseProvider, "twilio-auth-token-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.PhoneNumberAttributeType, state.PhoneNumberAttributeType, "phone-number-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.PhoneNumberJSONField, state.PhoneNumberJSONField, "phone-number-json-field")
	operations.AddStringOperationIfNecessary(&ops, plan.PhoneNumberJSONObjectFilter, state.PhoneNumberJSONObjectFilter, "phone-number-json-object-filter")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SenderPhoneNumber, state.SenderPhoneNumber, "sender-phone-number")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageTextBeforeOTP, state.MessageTextBeforeOTP, "message-text-before-otp")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageTextAfterOTP, state.MessageTextAfterOTP, "message-text-after-otp")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a twilio otp-delivery-mechanism
func (r *otpDeliveryMechanismResource) CreateTwilioOtpDeliveryMechanism(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan otpDeliveryMechanismResourceModel) (*otpDeliveryMechanismResourceModel, error) {
	var SenderPhoneNumberSlice []string
	plan.SenderPhoneNumber.ElementsAs(ctx, &SenderPhoneNumberSlice, false)
	addRequest := client.NewAddTwilioOtpDeliveryMechanismRequest(plan.Name.ValueString(),
		[]client.EnumtwilioOtpDeliveryMechanismSchemaUrn{client.ENUMTWILIOOTPDELIVERYMECHANISMSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0OTP_DELIVERY_MECHANISMTWILIO},
		plan.TwilioAccountSID.ValueString(),
		SenderPhoneNumberSlice,
		plan.Enabled.ValueBool())
	addOptionalTwilioOtpDeliveryMechanismFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.OtpDeliveryMechanismApi.AddOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddOtpDeliveryMechanismRequest(
		client.AddTwilioOtpDeliveryMechanismRequestAsAddOtpDeliveryMechanismRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.OtpDeliveryMechanismApi.AddOtpDeliveryMechanismExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Otp Delivery Mechanism", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state otpDeliveryMechanismResourceModel
	readTwilioOtpDeliveryMechanismResponse(ctx, addResponse.TwilioOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a email otp-delivery-mechanism
func (r *otpDeliveryMechanismResource) CreateEmailOtpDeliveryMechanism(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan otpDeliveryMechanismResourceModel) (*otpDeliveryMechanismResourceModel, error) {
	addRequest := client.NewAddEmailOtpDeliveryMechanismRequest(plan.Name.ValueString(),
		[]client.EnumemailOtpDeliveryMechanismSchemaUrn{client.ENUMEMAILOTPDELIVERYMECHANISMSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0OTP_DELIVERY_MECHANISMEMAIL},
		plan.SenderAddress.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalEmailOtpDeliveryMechanismFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.OtpDeliveryMechanismApi.AddOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddOtpDeliveryMechanismRequest(
		client.AddEmailOtpDeliveryMechanismRequestAsAddOtpDeliveryMechanismRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.OtpDeliveryMechanismApi.AddOtpDeliveryMechanismExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Otp Delivery Mechanism", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state otpDeliveryMechanismResourceModel
	readEmailOtpDeliveryMechanismResponse(ctx, addResponse.EmailOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party otp-delivery-mechanism
func (r *otpDeliveryMechanismResource) CreateThirdPartyOtpDeliveryMechanism(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan otpDeliveryMechanismResourceModel) (*otpDeliveryMechanismResourceModel, error) {
	addRequest := client.NewAddThirdPartyOtpDeliveryMechanismRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyOtpDeliveryMechanismSchemaUrn{client.ENUMTHIRDPARTYOTPDELIVERYMECHANISMSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0OTP_DELIVERY_MECHANISMTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyOtpDeliveryMechanismFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.OtpDeliveryMechanismApi.AddOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddOtpDeliveryMechanismRequest(
		client.AddThirdPartyOtpDeliveryMechanismRequestAsAddOtpDeliveryMechanismRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.OtpDeliveryMechanismApi.AddOtpDeliveryMechanismExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Otp Delivery Mechanism", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state otpDeliveryMechanismResourceModel
	readThirdPartyOtpDeliveryMechanismResponse(ctx, addResponse.ThirdPartyOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *otpDeliveryMechanismResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan otpDeliveryMechanismResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *otpDeliveryMechanismResourceModel
	var err error
	if plan.Type.ValueString() == "twilio" {
		state, err = r.CreateTwilioOtpDeliveryMechanism(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "email" {
		state, err = r.CreateEmailOtpDeliveryMechanism(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyOtpDeliveryMechanism(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	state.setStateValuesNotReturnedByAPI(&plan)
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
func (r *defaultOtpDeliveryMechanismResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan otpDeliveryMechanismResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.OtpDeliveryMechanismApi.GetOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Otp Delivery Mechanism", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state otpDeliveryMechanismResourceModel
	if readResponse.TwilioOtpDeliveryMechanismResponse != nil {
		readTwilioOtpDeliveryMechanismResponse(ctx, readResponse.TwilioOtpDeliveryMechanismResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EmailOtpDeliveryMechanismResponse != nil {
		readEmailOtpDeliveryMechanismResponse(ctx, readResponse.EmailOtpDeliveryMechanismResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyOtpDeliveryMechanismResponse != nil {
		readThirdPartyOtpDeliveryMechanismResponse(ctx, readResponse.ThirdPartyOtpDeliveryMechanismResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.OtpDeliveryMechanismApi.UpdateOtpDeliveryMechanism(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createOtpDeliveryMechanismOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.OtpDeliveryMechanismApi.UpdateOtpDeliveryMechanismExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Otp Delivery Mechanism", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.TwilioOtpDeliveryMechanismResponse != nil {
			readTwilioOtpDeliveryMechanismResponse(ctx, updateResponse.TwilioOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EmailOtpDeliveryMechanismResponse != nil {
			readEmailOtpDeliveryMechanismResponse(ctx, updateResponse.EmailOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyOtpDeliveryMechanismResponse != nil {
			readThirdPartyOtpDeliveryMechanismResponse(ctx, updateResponse.ThirdPartyOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *otpDeliveryMechanismResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOtpDeliveryMechanism(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultOtpDeliveryMechanismResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOtpDeliveryMechanism(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readOtpDeliveryMechanism(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state otpDeliveryMechanismResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.OtpDeliveryMechanismApi.GetOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Otp Delivery Mechanism", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Otp Delivery Mechanism", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TwilioOtpDeliveryMechanismResponse != nil {
		readTwilioOtpDeliveryMechanismResponse(ctx, readResponse.TwilioOtpDeliveryMechanismResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.EmailOtpDeliveryMechanismResponse != nil {
		readEmailOtpDeliveryMechanismResponse(ctx, readResponse.EmailOtpDeliveryMechanismResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyOtpDeliveryMechanismResponse != nil {
		readThirdPartyOtpDeliveryMechanismResponse(ctx, readResponse.ThirdPartyOtpDeliveryMechanismResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *otpDeliveryMechanismResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOtpDeliveryMechanism(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOtpDeliveryMechanismResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOtpDeliveryMechanism(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOtpDeliveryMechanism(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan otpDeliveryMechanismResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state otpDeliveryMechanismResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.OtpDeliveryMechanismApi.UpdateOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createOtpDeliveryMechanismOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.OtpDeliveryMechanismApi.UpdateOtpDeliveryMechanismExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Otp Delivery Mechanism", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.TwilioOtpDeliveryMechanismResponse != nil {
			readTwilioOtpDeliveryMechanismResponse(ctx, updateResponse.TwilioOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.EmailOtpDeliveryMechanismResponse != nil {
			readEmailOtpDeliveryMechanismResponse(ctx, updateResponse.EmailOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyOtpDeliveryMechanismResponse != nil {
			readThirdPartyOtpDeliveryMechanismResponse(ctx, updateResponse.ThirdPartyOtpDeliveryMechanismResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultOtpDeliveryMechanismResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *otpDeliveryMechanismResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state otpDeliveryMechanismResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.OtpDeliveryMechanismApi.DeleteOtpDeliveryMechanismExecute(r.apiClient.OtpDeliveryMechanismApi.DeleteOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Otp Delivery Mechanism", err, httpResp)
		return
	}
}

func (r *otpDeliveryMechanismResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOtpDeliveryMechanism(ctx, req, resp)
}

func (r *defaultOtpDeliveryMechanismResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOtpDeliveryMechanism(ctx, req, resp)
}

func importOtpDeliveryMechanism(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
