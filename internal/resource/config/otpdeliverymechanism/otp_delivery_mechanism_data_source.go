package otpdeliverymechanism

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &otpDeliveryMechanismDataSource{}
	_ datasource.DataSourceWithConfigure = &otpDeliveryMechanismDataSource{}
)

// Create a Otp Delivery Mechanism data source
func NewOtpDeliveryMechanismDataSource() datasource.DataSource {
	return &otpDeliveryMechanismDataSource{}
}

// otpDeliveryMechanismDataSource is the datasource implementation.
type otpDeliveryMechanismDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *otpDeliveryMechanismDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_otp_delivery_mechanism"
}

// Configure adds the provider configured client to the data source.
func (r *otpDeliveryMechanismDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type otpDeliveryMechanismDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *otpDeliveryMechanismDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Otp Delivery Mechanism.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of OTP Delivery Mechanism resource. Options are ['twilio', 'email', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party OTP Delivery Mechanism.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party OTP Delivery Mechanism. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"email_address_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute that holds the email address to which the message should be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
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
			"sender_address": schema.StringAttribute{
				Description: "The e-mail address to use as the sender for the one-time password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"message_subject": schema.StringAttribute{
				Description: "The subject to use for the e-mail message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Twilio service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"twilio_account_sid": schema.StringAttribute{
				Description: "The unique identifier assigned to the Twilio account that will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"twilio_auth_token": schema.StringAttribute{
				Description: "The auth token for the Twilio account that will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"twilio_auth_token_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider that may be used to obtain the auth token for the Twilio account that will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"phone_number_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute in the user's entry that holds the phone number to which the message should be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"phone_number_json_field": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the phone number to which the message should be sent. The phone number must be contained in a top-level field whose value is a single string.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"phone_number_json_object_filter": schema.StringAttribute{
				Description: "A JSON object filter that may be used to identify which phone number value to use when sending the message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sender_phone_number": schema.SetAttribute{
				Description: "The outgoing phone number to use for the messages. Values must be phone numbers you have obtained for use with your Twilio account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"message_text_before_otp": schema.StringAttribute{
				Description: "Any text that should appear in the message before the one-time password value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"message_text_after_otp": schema.StringAttribute{
				Description: "Any text that should appear in the message after the one-time password value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this OTP Delivery Mechanism",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this OTP Delivery Mechanism is enabled for use in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a TwilioOtpDeliveryMechanismResponse object into the model struct
func readTwilioOtpDeliveryMechanismResponseDataSource(ctx context.Context, r *client.TwilioOtpDeliveryMechanismResponse, state *otpDeliveryMechanismDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("twilio")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.TwilioAccountSID = types.StringValue(r.TwilioAccountSID)
	state.TwilioAuthTokenPassphraseProvider = internaltypes.StringTypeOrNil(r.TwilioAuthTokenPassphraseProvider, false)
	state.PhoneNumberAttributeType = types.StringValue(r.PhoneNumberAttributeType)
	state.PhoneNumberJSONField = internaltypes.StringTypeOrNil(r.PhoneNumberJSONField, false)
	state.PhoneNumberJSONObjectFilter = internaltypes.StringTypeOrNil(r.PhoneNumberJSONObjectFilter, false)
	state.SenderPhoneNumber = internaltypes.GetStringSet(r.SenderPhoneNumber)
	state.MessageTextBeforeOTP = internaltypes.StringTypeOrNil(r.MessageTextBeforeOTP, false)
	state.MessageTextAfterOTP = internaltypes.StringTypeOrNil(r.MessageTextAfterOTP, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a EmailOtpDeliveryMechanismResponse object into the model struct
func readEmailOtpDeliveryMechanismResponseDataSource(ctx context.Context, r *client.EmailOtpDeliveryMechanismResponse, state *otpDeliveryMechanismDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("email")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EmailAddressAttributeType = types.StringValue(r.EmailAddressAttributeType)
	state.EmailAddressJSONField = internaltypes.StringTypeOrNil(r.EmailAddressJSONField, false)
	state.EmailAddressJSONObjectFilter = internaltypes.StringTypeOrNil(r.EmailAddressJSONObjectFilter, false)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.MessageSubject = types.StringValue(r.MessageSubject)
	state.MessageTextBeforeOTP = internaltypes.StringTypeOrNil(r.MessageTextBeforeOTP, false)
	state.MessageTextAfterOTP = internaltypes.StringTypeOrNil(r.MessageTextAfterOTP, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyOtpDeliveryMechanismResponse object into the model struct
func readThirdPartyOtpDeliveryMechanismResponseDataSource(ctx context.Context, r *client.ThirdPartyOtpDeliveryMechanismResponse, state *otpDeliveryMechanismDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *otpDeliveryMechanismDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state otpDeliveryMechanismDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.OtpDeliveryMechanismAPI.GetOtpDeliveryMechanism(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Otp Delivery Mechanism", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TwilioOtpDeliveryMechanismResponse != nil {
		readTwilioOtpDeliveryMechanismResponseDataSource(ctx, readResponse.TwilioOtpDeliveryMechanismResponse, &state, &resp.Diagnostics)
	}
	if readResponse.EmailOtpDeliveryMechanismResponse != nil {
		readEmailOtpDeliveryMechanismResponseDataSource(ctx, readResponse.EmailOtpDeliveryMechanismResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyOtpDeliveryMechanismResponse != nil {
		readThirdPartyOtpDeliveryMechanismResponseDataSource(ctx, readResponse.ThirdPartyOtpDeliveryMechanismResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
