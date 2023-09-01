package alerthandler

import (
	"context"

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
	_ resource.Resource                = &alertHandlerResource{}
	_ resource.ResourceWithConfigure   = &alertHandlerResource{}
	_ resource.ResourceWithImportState = &alertHandlerResource{}
	_ resource.Resource                = &defaultAlertHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultAlertHandlerResource{}
	_ resource.ResourceWithImportState = &defaultAlertHandlerResource{}
)

// Create a Alert Handler resource
func NewAlertHandlerResource() resource.Resource {
	return &alertHandlerResource{}
}

func NewDefaultAlertHandlerResource() resource.Resource {
	return &defaultAlertHandlerResource{}
}

// alertHandlerResource is the resource implementation.
type alertHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAlertHandlerResource is the resource implementation.
type defaultAlertHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *alertHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_handler"
}

func (r *defaultAlertHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_alert_handler"
}

// Configure adds the provider configured client to the resource.
func (r *alertHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultAlertHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type alertHandlerResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	Command                           types.String `tfsdk:"command"`
	ScriptClass                       types.String `tfsdk:"script_class"`
	HttpProxyExternalServer           types.String `tfsdk:"http_proxy_external_server"`
	TwilioAccountSID                  types.String `tfsdk:"twilio_account_sid"`
	TwilioAuthToken                   types.String `tfsdk:"twilio_auth_token"`
	TwilioAuthTokenPassphraseProvider types.String `tfsdk:"twilio_auth_token_passphrase_provider"`
	SenderPhoneNumber                 types.Set    `tfsdk:"sender_phone_number"`
	RecipientPhoneNumber              types.Set    `tfsdk:"recipient_phone_number"`
	LongMessageBehavior               types.String `tfsdk:"long_message_behavior"`
	ServerHostName                    types.String `tfsdk:"server_host_name"`
	ServerPort                        types.Int64  `tfsdk:"server_port"`
	CommunityName                     types.String `tfsdk:"community_name"`
	ScriptArgument                    types.Set    `tfsdk:"script_argument"`
	SenderAddress                     types.String `tfsdk:"sender_address"`
	RecipientAddress                  types.Set    `tfsdk:"recipient_address"`
	MessageSubject                    types.String `tfsdk:"message_subject"`
	MessageBody                       types.String `tfsdk:"message_body"`
	IncludeMonitorDataFilter          types.String `tfsdk:"include_monitor_data_filter"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
	Asynchronous                      types.Bool   `tfsdk:"asynchronous"`
	EnabledAlertSeverity              types.Set    `tfsdk:"enabled_alert_severity"`
	EnabledAlertType                  types.Set    `tfsdk:"enabled_alert_type"`
	DisabledAlertType                 types.Set    `tfsdk:"disabled_alert_type"`
}

type defaultAlertHandlerResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	Command                           types.String `tfsdk:"command"`
	ScriptClass                       types.String `tfsdk:"script_class"`
	HttpProxyExternalServer           types.String `tfsdk:"http_proxy_external_server"`
	TwilioAccountSID                  types.String `tfsdk:"twilio_account_sid"`
	TwilioAuthToken                   types.String `tfsdk:"twilio_auth_token"`
	TwilioAuthTokenPassphraseProvider types.String `tfsdk:"twilio_auth_token_passphrase_provider"`
	SenderPhoneNumber                 types.Set    `tfsdk:"sender_phone_number"`
	RecipientPhoneNumber              types.Set    `tfsdk:"recipient_phone_number"`
	LongMessageBehavior               types.String `tfsdk:"long_message_behavior"`
	ServerHostName                    types.String `tfsdk:"server_host_name"`
	ServerPort                        types.Int64  `tfsdk:"server_port"`
	CommunityName                     types.String `tfsdk:"community_name"`
	ScriptArgument                    types.Set    `tfsdk:"script_argument"`
	OutputLocation                    types.String `tfsdk:"output_location"`
	SenderAddress                     types.String `tfsdk:"sender_address"`
	RecipientAddress                  types.Set    `tfsdk:"recipient_address"`
	MessageSubject                    types.String `tfsdk:"message_subject"`
	MessageBody                       types.String `tfsdk:"message_body"`
	IncludeMonitorDataFilter          types.String `tfsdk:"include_monitor_data_filter"`
	OutputFormat                      types.String `tfsdk:"output_format"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
	Asynchronous                      types.Bool   `tfsdk:"asynchronous"`
	EnabledAlertSeverity              types.Set    `tfsdk:"enabled_alert_severity"`
	EnabledAlertType                  types.Set    `tfsdk:"enabled_alert_type"`
	DisabledAlertType                 types.Set    `tfsdk:"disabled_alert_type"`
}

// GetSchema defines the schema for the resource.
func (r *alertHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	alertHandlerSchema(ctx, req, resp, false)
}

func (r *defaultAlertHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	alertHandlerSchema(ctx, req, resp, true)
}

func alertHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Alert Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Alert Handler resource. Options are ['output', 'smtp', 'jmx', 'groovy-scripted', 'custom', 'snmp', 'twilio', 'error-log', 'snmp-sub-agent', 'exec', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"smtp", "jmx", "groovy-scripted", "snmp", "twilio", "error-log", "snmp-sub-agent", "exec", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Alert Handler.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Alert Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"command": schema.StringAttribute{
				Description: "Specifies the path of the command to execute, without any arguments. It must be an absolute path for reasons of security and reliability.",
				Optional:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Alert Handler.",
				Optional:    true,
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
			"sender_phone_number": schema.SetAttribute{
				Description: "The outgoing phone number to use for the messages. Values must be phone numbers you have obtained for use with your Twilio account.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"recipient_phone_number": schema.SetAttribute{
				Description: "The phone number to which alert notifications should be delivered.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"long_message_behavior": schema.StringAttribute{
				Description: "The behavior to use for alert messages that are longer than the 160-character size limit for SMS messages.",
				Optional:    true,
				Computed:    true,
			},
			"server_host_name": schema.StringAttribute{
				Description: "Specifies the address of the SNMP agent to which traps will be sent.",
				Optional:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the SNMP agent to which traps will be sent.",
				Optional:    true,
				Computed:    true,
			},
			"community_name": schema.StringAttribute{
				Description: "Specifies the name of the community to which the traps will be sent.",
				Optional:    true,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Alert Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"sender_address": schema.StringAttribute{
				Description: "Specifies the email address to use as the sender for messages generated by this alert handler.",
				Optional:    true,
			},
			"recipient_address": schema.SetAttribute{
				Description: "Specifies an email address to which the messages should be sent.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"message_subject": schema.StringAttribute{
				Description: "Specifies the subject that should be used for email messages generated by this alert handler.",
				Optional:    true,
				Computed:    true,
			},
			"message_body": schema.StringAttribute{
				Description: "Specifies the body that should be used for email messages generated by this alert handler.",
				Optional:    true,
				Computed:    true,
			},
			"include_monitor_data_filter": schema.StringAttribute{
				Description: "Include monitor entries that match this filter.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Alert Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Alert Handler is enabled.",
				Required:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`output`, `groovy-scripted`, `custom`, `error-log`, `third-party`]: Indicates whether the server should attempt to invoke this Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `smtp`: Indicates whether the server should attempt to invoke this SMTP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `jmx`: Indicates whether the server should attempt to invoke this JMX Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `snmp`: Indicates whether the server should attempt to invoke this SNMP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `twilio`: Indicates whether the server should attempt to invoke this Twilio Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `snmp-sub-agent`: Indicates whether the server should attempt to invoke this SNMP Sub Agent Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `exec`: Indicates whether the server should attempt to invoke this Exec Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`output`, `groovy-scripted`, `custom`, `error-log`, `third-party`]: Indicates whether the server should attempt to invoke this Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `smtp`: Indicates whether the server should attempt to invoke this SMTP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `jmx`: Indicates whether the server should attempt to invoke this JMX Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `snmp`: Indicates whether the server should attempt to invoke this SNMP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `twilio`: Indicates whether the server should attempt to invoke this Twilio Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `snmp-sub-agent`: Indicates whether the server should attempt to invoke this SNMP Sub Agent Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `exec`: Indicates whether the server should attempt to invoke this Exec Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.",
				Optional:            true,
				Computed:            true,
			},
			"enabled_alert_severity": schema.SetAttribute{
				Description: "Specifies the alert severities for which this alert handler should be used. If no values are provided, then this alert handler will be enabled for alerts with any severity.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"enabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that are enabled for this alert handler.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"disabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that are disabled for this alert handler.",
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
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"output", "smtp", "jmx", "groovy-scripted", "custom", "snmp", "twilio", "error-log", "snmp-sub-agent", "exec", "third-party"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["output_location"] = schema.StringAttribute{
			Description: "The location to which alert messages will be written.",
		}
		schemaDef.Attributes["output_format"] = schema.StringAttribute{
			Description: "The format to use when writing the alert messages.",
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *alertHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
	var planModel, configModel alertHandlerResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for smtp type
	if resourceType == "smtp" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MessageSubject) {
			defaultVal := types.StringValue("%%alert-severity%% alert from %%instance-name%%: %%alert-type%%")
			if !planModel.MessageSubject.Equal(defaultVal) {
				planModel.MessageSubject = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MessageBody) {
			defaultVal := types.StringValue("Alert Type: %%alert-type%%\\nAlert Severity: %%alert-severity%%\\nInstance Name: %%instance-name%%\\nAlert ID: %%alert-id%%\\nAlert OID: %%alert-oid%%\\n\\nAlert Message:\\n%%alert-message%%")
			if !planModel.MessageBody.Equal(defaultVal) {
				planModel.MessageBody = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for jmx type
	if resourceType == "jmx" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for groovy-scripted type
	if resourceType == "groovy-scripted" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(false)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for snmp type
	if resourceType == "snmp" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(162)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CommunityName) {
			defaultVal := types.StringValue("public")
			if !planModel.CommunityName.Equal(defaultVal) {
				planModel.CommunityName = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for twilio type
	if resourceType == "twilio" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LongMessageBehavior) {
			defaultVal := types.StringValue("truncate")
			if !planModel.LongMessageBehavior.Equal(defaultVal) {
				planModel.LongMessageBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for error-log type
	if resourceType == "error-log" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(false)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for snmp-sub-agent type
	if resourceType == "snmp-sub-agent" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for exec type
	if resourceType == "exec" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for third-party type
	if resourceType == "third-party" {
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(false)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultAlertHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanAlertHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanAlertHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model defaultAlertHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsAlertHandler() []resource.ConfigValidator {
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
			path.MatchRoot("sender_address"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("recipient_address"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("message_subject"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("message_body"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_monitor_data_filter"),
			path.MatchRoot("type"),
			[]string{"smtp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_host_name"),
			path.MatchRoot("type"),
			[]string{"snmp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_port"),
			path.MatchRoot("type"),
			[]string{"snmp"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("community_name"),
			path.MatchRoot("type"),
			[]string{"snmp"},
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
			path.MatchRoot("sender_phone_number"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("recipient_phone_number"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("long_message_behavior"),
			path.MatchRoot("type"),
			[]string{"twilio"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("command"),
			path.MatchRoot("type"),
			[]string{"exec"},
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
			"smtp",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("sender_address"), path.MatchRoot("recipient_address")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"jmx",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"snmp",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("server_host_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"twilio",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("twilio_account_sid"), path.MatchRoot("sender_phone_number"), path.MatchRoot("recipient_phone_number")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"error-log",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"snmp-sub-agent",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"exec",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("command")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r alertHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsAlertHandler()
}

// Add config validators
func (r defaultAlertHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("output_location"),
			path.MatchRoot("type"),
			[]string{"output"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("output_format"),
			path.MatchRoot("type"),
			[]string{"output"},
		),
	}
	return append(configValidatorsAlertHandler(), validators...)
}

// Add optional fields to create request for smtp alert-handler
func addOptionalSmtpAlertHandlerFields(ctx context.Context, addRequest *client.AddSmtpAlertHandlerRequest, plan alertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageSubject) {
		addRequest.MessageSubject = plan.MessageSubject.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MessageBody) {
		addRequest.MessageBody = plan.MessageBody.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IncludeMonitorDataFilter) {
		addRequest.IncludeMonitorDataFilter = plan.IncludeMonitorDataFilter.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for jmx alert-handler
func addOptionalJmxAlertHandlerFields(ctx context.Context, addRequest *client.AddJmxAlertHandlerRequest, plan alertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for groovy-scripted alert-handler
func addOptionalGroovyScriptedAlertHandlerFields(ctx context.Context, addRequest *client.AddGroovyScriptedAlertHandlerRequest, plan alertHandlerResourceModel) error {
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
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for snmp alert-handler
func addOptionalSnmpAlertHandlerFields(ctx context.Context, addRequest *client.AddSnmpAlertHandlerRequest, plan alertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CommunityName) {
		addRequest.CommunityName = plan.CommunityName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for twilio alert-handler
func addOptionalTwilioAlertHandlerFields(ctx context.Context, addRequest *client.AddTwilioAlertHandlerRequest, plan alertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
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
	if internaltypes.IsNonEmptyString(plan.LongMessageBehavior) {
		longMessageBehavior, err := client.NewEnumalertHandlerLongMessageBehaviorPropFromValue(plan.LongMessageBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LongMessageBehavior = longMessageBehavior
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for error-log alert-handler
func addOptionalErrorLogAlertHandlerFields(ctx context.Context, addRequest *client.AddErrorLogAlertHandlerRequest, plan alertHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for snmp-sub-agent alert-handler
func addOptionalSnmpSubAgentAlertHandlerFields(ctx context.Context, addRequest *client.AddSnmpSubAgentAlertHandlerRequest, plan alertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for exec alert-handler
func addOptionalExecAlertHandlerFields(ctx context.Context, addRequest *client.AddExecAlertHandlerRequest, plan alertHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Add optional fields to create request for third-party alert-handler
func addOptionalThirdPartyAlertHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyAlertHandlerRequest, plan alertHandlerResourceModel) error {
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
	if internaltypes.IsDefined(plan.EnabledAlertSeverity) {
		var slice []string
		plan.EnabledAlertSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.EnabledAlertType) {
		var slice []string
		plan.EnabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerEnabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerEnabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.EnabledAlertType = enumSlice
	}
	if internaltypes.IsDefined(plan.DisabledAlertType) {
		var slice []string
		plan.DisabledAlertType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumalertHandlerDisabledAlertTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumalertHandlerDisabledAlertTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DisabledAlertType = enumSlice
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateAlertHandlerUnknownValues(model *alertHandlerResourceModel) {
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RecipientAddress.IsUnknown() || model.RecipientAddress.IsNull() {
		model.RecipientAddress, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SenderPhoneNumber.IsUnknown() || model.SenderPhoneNumber.IsNull() {
		model.SenderPhoneNumber, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RecipientPhoneNumber.IsUnknown() || model.RecipientPhoneNumber.IsNull() {
		model.RecipientPhoneNumber, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MessageSubject.IsUnknown() || model.MessageSubject.IsNull() {
		model.MessageSubject = types.StringValue("")
	}
	if model.LongMessageBehavior.IsUnknown() || model.LongMessageBehavior.IsNull() {
		model.LongMessageBehavior = types.StringValue("")
	}
	if model.CommunityName.IsUnknown() || model.CommunityName.IsNull() {
		model.CommunityName = types.StringValue("")
	}
	if model.MessageBody.IsUnknown() || model.MessageBody.IsNull() {
		model.MessageBody = types.StringValue("")
	}
	if model.TwilioAuthToken.IsUnknown() {
		model.TwilioAuthToken = types.StringNull()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateAlertHandlerUnknownValuesDefault(model *defaultAlertHandlerResourceModel) {
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RecipientAddress.IsUnknown() || model.RecipientAddress.IsNull() {
		model.RecipientAddress, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SenderPhoneNumber.IsUnknown() || model.SenderPhoneNumber.IsNull() {
		model.SenderPhoneNumber, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RecipientPhoneNumber.IsUnknown() || model.RecipientPhoneNumber.IsNull() {
		model.RecipientPhoneNumber, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.TwilioAccountSID.IsUnknown() || model.TwilioAccountSID.IsNull() {
		model.TwilioAccountSID = types.StringValue("")
	}
	if model.SenderAddress.IsUnknown() || model.SenderAddress.IsNull() {
		model.SenderAddress = types.StringValue("")
	}
	if model.MessageSubject.IsUnknown() || model.MessageSubject.IsNull() {
		model.MessageSubject = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.OutputFormat.IsUnknown() || model.OutputFormat.IsNull() {
		model.OutputFormat = types.StringValue("")
	}
	if model.HttpProxyExternalServer.IsUnknown() || model.HttpProxyExternalServer.IsNull() {
		model.HttpProxyExternalServer = types.StringValue("")
	}
	if model.OutputLocation.IsUnknown() || model.OutputLocation.IsNull() {
		model.OutputLocation = types.StringValue("")
	}
	if model.TwilioAuthTokenPassphraseProvider.IsUnknown() || model.TwilioAuthTokenPassphraseProvider.IsNull() {
		model.TwilioAuthTokenPassphraseProvider = types.StringValue("")
	}
	if model.TwilioAuthToken.IsUnknown() || model.TwilioAuthToken.IsNull() {
		model.TwilioAuthToken = types.StringValue("")
	}
	if model.LongMessageBehavior.IsUnknown() || model.LongMessageBehavior.IsNull() {
		model.LongMessageBehavior = types.StringValue("")
	}
	if model.Command.IsUnknown() || model.Command.IsNull() {
		model.Command = types.StringValue("")
	}
	if model.ServerHostName.IsUnknown() || model.ServerHostName.IsNull() {
		model.ServerHostName = types.StringValue("")
	}
	if model.CommunityName.IsUnknown() || model.CommunityName.IsNull() {
		model.CommunityName = types.StringValue("")
	}
	if model.IncludeMonitorDataFilter.IsUnknown() || model.IncludeMonitorDataFilter.IsNull() {
		model.IncludeMonitorDataFilter = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
	if model.MessageBody.IsUnknown() || model.MessageBody.IsNull() {
		model.MessageBody = types.StringValue("")
	}
}

// Read a OutputAlertHandlerResponse object into the model struct
func readOutputAlertHandlerResponseDefault(ctx context.Context, r *client.OutputAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("output")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumalertHandlerOutputLocationProp(r.OutputLocation), true)
	state.OutputFormat = internaltypes.StringTypeOrNil(
		client.StringPointerEnumalertHandlerOutputFormatProp(r.OutputFormat), true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a SmtpAlertHandlerResponse object into the model struct
func readSmtpAlertHandlerResponse(ctx context.Context, r *client.SmtpAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.RecipientAddress = internaltypes.GetStringSet(r.RecipientAddress)
	state.MessageSubject = types.StringValue(r.MessageSubject)
	state.MessageBody = types.StringValue(r.MessageBody)
	state.IncludeMonitorDataFilter = internaltypes.StringTypeOrNil(r.IncludeMonitorDataFilter, internaltypes.IsEmptyString(expectedValues.IncludeMonitorDataFilter))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a SmtpAlertHandlerResponse object into the model struct
func readSmtpAlertHandlerResponseDefault(ctx context.Context, r *client.SmtpAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.RecipientAddress = internaltypes.GetStringSet(r.RecipientAddress)
	state.MessageSubject = types.StringValue(r.MessageSubject)
	state.MessageBody = types.StringValue(r.MessageBody)
	state.IncludeMonitorDataFilter = internaltypes.StringTypeOrNil(r.IncludeMonitorDataFilter, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a JmxAlertHandlerResponse object into the model struct
func readJmxAlertHandlerResponse(ctx context.Context, r *client.JmxAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jmx")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a JmxAlertHandlerResponse object into the model struct
func readJmxAlertHandlerResponseDefault(ctx context.Context, r *client.JmxAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jmx")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a GroovyScriptedAlertHandlerResponse object into the model struct
func readGroovyScriptedAlertHandlerResponse(ctx context.Context, r *client.GroovyScriptedAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a GroovyScriptedAlertHandlerResponse object into the model struct
func readGroovyScriptedAlertHandlerResponseDefault(ctx context.Context, r *client.GroovyScriptedAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a CustomAlertHandlerResponse object into the model struct
func readCustomAlertHandlerResponseDefault(ctx context.Context, r *client.CustomAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a SnmpAlertHandlerResponse object into the model struct
func readSnmpAlertHandlerResponse(ctx context.Context, r *client.SnmpAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("snmp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.CommunityName = types.StringValue(r.CommunityName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a SnmpAlertHandlerResponse object into the model struct
func readSnmpAlertHandlerResponseDefault(ctx context.Context, r *client.SnmpAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("snmp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.CommunityName = types.StringValue(r.CommunityName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a TwilioAlertHandlerResponse object into the model struct
func readTwilioAlertHandlerResponse(ctx context.Context, r *client.TwilioAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("twilio")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.TwilioAccountSID = types.StringValue(r.TwilioAccountSID)
	state.TwilioAuthTokenPassphraseProvider = internaltypes.StringTypeOrNil(r.TwilioAuthTokenPassphraseProvider, internaltypes.IsEmptyString(expectedValues.TwilioAuthTokenPassphraseProvider))
	state.SenderPhoneNumber = internaltypes.GetStringSet(r.SenderPhoneNumber)
	state.RecipientPhoneNumber = internaltypes.GetStringSet(r.RecipientPhoneNumber)
	state.LongMessageBehavior = types.StringValue(r.LongMessageBehavior.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a TwilioAlertHandlerResponse object into the model struct
func readTwilioAlertHandlerResponseDefault(ctx context.Context, r *client.TwilioAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("twilio")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, true)
	state.TwilioAccountSID = types.StringValue(r.TwilioAccountSID)
	state.TwilioAuthTokenPassphraseProvider = internaltypes.StringTypeOrNil(r.TwilioAuthTokenPassphraseProvider, true)
	state.SenderPhoneNumber = internaltypes.GetStringSet(r.SenderPhoneNumber)
	state.RecipientPhoneNumber = internaltypes.GetStringSet(r.RecipientPhoneNumber)
	state.LongMessageBehavior = types.StringValue(r.LongMessageBehavior.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a ErrorLogAlertHandlerResponse object into the model struct
func readErrorLogAlertHandlerResponse(ctx context.Context, r *client.ErrorLogAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error-log")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a ErrorLogAlertHandlerResponse object into the model struct
func readErrorLogAlertHandlerResponseDefault(ctx context.Context, r *client.ErrorLogAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error-log")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a SnmpSubAgentAlertHandlerResponse object into the model struct
func readSnmpSubAgentAlertHandlerResponse(ctx context.Context, r *client.SnmpSubAgentAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("snmp-sub-agent")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a SnmpSubAgentAlertHandlerResponse object into the model struct
func readSnmpSubAgentAlertHandlerResponseDefault(ctx context.Context, r *client.SnmpSubAgentAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("snmp-sub-agent")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a ExecAlertHandlerResponse object into the model struct
func readExecAlertHandlerResponse(ctx context.Context, r *client.ExecAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exec")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Command = types.StringValue(r.Command)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a ExecAlertHandlerResponse object into the model struct
func readExecAlertHandlerResponseDefault(ctx context.Context, r *client.ExecAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exec")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Command = types.StringValue(r.Command)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Read a ThirdPartyAlertHandlerResponse object into the model struct
func readThirdPartyAlertHandlerResponse(ctx context.Context, r *client.ThirdPartyAlertHandlerResponse, state *alertHandlerResourceModel, expectedValues *alertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValues(state)
}

// Read a ThirdPartyAlertHandlerResponse object into the model struct
func readThirdPartyAlertHandlerResponseDefault(ctx context.Context, r *client.ThirdPartyAlertHandlerResponse, state *defaultAlertHandlerResourceModel, expectedValues *defaultAlertHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAlertHandlerUnknownValuesDefault(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *defaultAlertHandlerResourceModel) setStateValuesNotReturnedByAPI(expectedValues *defaultAlertHandlerResourceModel) {
	if !expectedValues.TwilioAuthToken.IsUnknown() {
		state.TwilioAuthToken = expectedValues.TwilioAuthToken
	}
}

func (state *alertHandlerResourceModel) setStateValuesNotReturnedByAPI(expectedValues *alertHandlerResourceModel) {
	if !expectedValues.TwilioAuthToken.IsUnknown() {
		state.TwilioAuthToken = expectedValues.TwilioAuthToken
	}
}

// Create any update operations necessary to make the state match the plan
func createAlertHandlerOperations(plan alertHandlerResourceModel, state alertHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Command, state.Command, "command")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAccountSID, state.TwilioAccountSID, "twilio-account-sid")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthToken, state.TwilioAuthToken, "twilio-auth-token")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthTokenPassphraseProvider, state.TwilioAuthTokenPassphraseProvider, "twilio-auth-token-passphrase-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SenderPhoneNumber, state.SenderPhoneNumber, "sender-phone-number")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientPhoneNumber, state.RecipientPhoneNumber, "recipient-phone-number")
	operations.AddStringOperationIfNecessary(&ops, plan.LongMessageBehavior, state.LongMessageBehavior, "long-message-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.CommunityName, state.CommunityName, "community-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.SenderAddress, state.SenderAddress, "sender-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientAddress, state.RecipientAddress, "recipient-address")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageSubject, state.MessageSubject, "message-subject")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageBody, state.MessageBody, "message-body")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeMonitorDataFilter, state.IncludeMonitorDataFilter, "include-monitor-data-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertSeverity, state.EnabledAlertSeverity, "enabled-alert-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertType, state.EnabledAlertType, "enabled-alert-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledAlertType, state.DisabledAlertType, "disabled-alert-type")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createAlertHandlerOperationsDefault(plan defaultAlertHandlerResourceModel, state defaultAlertHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Command, state.Command, "command")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAccountSID, state.TwilioAccountSID, "twilio-account-sid")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthToken, state.TwilioAuthToken, "twilio-auth-token")
	operations.AddStringOperationIfNecessary(&ops, plan.TwilioAuthTokenPassphraseProvider, state.TwilioAuthTokenPassphraseProvider, "twilio-auth-token-passphrase-provider")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SenderPhoneNumber, state.SenderPhoneNumber, "sender-phone-number")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientPhoneNumber, state.RecipientPhoneNumber, "recipient-phone-number")
	operations.AddStringOperationIfNecessary(&ops, plan.LongMessageBehavior, state.LongMessageBehavior, "long-message-behavior")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.CommunityName, state.CommunityName, "community-name")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputLocation, state.OutputLocation, "output-location")
	operations.AddStringOperationIfNecessary(&ops, plan.SenderAddress, state.SenderAddress, "sender-address")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RecipientAddress, state.RecipientAddress, "recipient-address")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageSubject, state.MessageSubject, "message-subject")
	operations.AddStringOperationIfNecessary(&ops, plan.MessageBody, state.MessageBody, "message-body")
	operations.AddStringOperationIfNecessary(&ops, plan.IncludeMonitorDataFilter, state.IncludeMonitorDataFilter, "include-monitor-data-filter")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputFormat, state.OutputFormat, "output-format")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertSeverity, state.EnabledAlertSeverity, "enabled-alert-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.EnabledAlertType, state.EnabledAlertType, "enabled-alert-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DisabledAlertType, state.DisabledAlertType, "disabled-alert-type")
	return ops
}

// Create a smtp alert-handler
func (r *alertHandlerResource) CreateSmtpAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	var RecipientAddressSlice []string
	plan.RecipientAddress.ElementsAs(ctx, &RecipientAddressSlice, false)
	addRequest := client.NewAddSmtpAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumsmtpAlertHandlerSchemaUrn{client.ENUMSMTPALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERSMTP},
		plan.SenderAddress.ValueString(),
		RecipientAddressSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSmtpAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddSmtpAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readSmtpAlertHandlerResponse(ctx, addResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a jmx alert-handler
func (r *alertHandlerResource) CreateJmxAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddJmxAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumjmxAlertHandlerSchemaUrn{client.ENUMJMXALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERJMX},
		plan.Enabled.ValueBool())
	err := addOptionalJmxAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddJmxAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readJmxAlertHandlerResponse(ctx, addResponse.JmxAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted alert-handler
func (r *alertHandlerResource) CreateGroovyScriptedAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedAlertHandlerSchemaUrn{client.ENUMGROOVYSCRIPTEDALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddGroovyScriptedAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readGroovyScriptedAlertHandlerResponse(ctx, addResponse.GroovyScriptedAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a snmp alert-handler
func (r *alertHandlerResource) CreateSnmpAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddSnmpAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumsnmpAlertHandlerSchemaUrn{client.ENUMSNMPALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERSNMP},
		plan.ServerHostName.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalSnmpAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddSnmpAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readSnmpAlertHandlerResponse(ctx, addResponse.SnmpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a twilio alert-handler
func (r *alertHandlerResource) CreateTwilioAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	var SenderPhoneNumberSlice []string
	plan.SenderPhoneNumber.ElementsAs(ctx, &SenderPhoneNumberSlice, false)
	var RecipientPhoneNumberSlice []string
	plan.RecipientPhoneNumber.ElementsAs(ctx, &RecipientPhoneNumberSlice, false)
	addRequest := client.NewAddTwilioAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumtwilioAlertHandlerSchemaUrn{client.ENUMTWILIOALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERTWILIO},
		plan.TwilioAccountSID.ValueString(),
		SenderPhoneNumberSlice,
		RecipientPhoneNumberSlice,
		plan.Enabled.ValueBool())
	err := addOptionalTwilioAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddTwilioAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readTwilioAlertHandlerResponse(ctx, addResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a error-log alert-handler
func (r *alertHandlerResource) CreateErrorLogAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddErrorLogAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumerrorLogAlertHandlerSchemaUrn{client.ENUMERRORLOGALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERERROR_LOG},
		plan.Enabled.ValueBool())
	err := addOptionalErrorLogAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddErrorLogAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readErrorLogAlertHandlerResponse(ctx, addResponse.ErrorLogAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a snmp-sub-agent alert-handler
func (r *alertHandlerResource) CreateSnmpSubAgentAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddSnmpSubAgentAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumsnmpSubAgentAlertHandlerSchemaUrn{client.ENUMSNMPSUBAGENTALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERSNMP_SUB_AGENT},
		plan.Enabled.ValueBool())
	err := addOptionalSnmpSubAgentAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddSnmpSubAgentAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readSnmpSubAgentAlertHandlerResponse(ctx, addResponse.SnmpSubAgentAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a exec alert-handler
func (r *alertHandlerResource) CreateExecAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddExecAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumexecAlertHandlerSchemaUrn{client.ENUMEXECALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLEREXEC},
		plan.Command.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalExecAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddExecAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readExecAlertHandlerResponse(ctx, addResponse.ExecAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party alert-handler
func (r *alertHandlerResource) CreateThirdPartyAlertHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan alertHandlerResourceModel) (*alertHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyAlertHandlerRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyAlertHandlerSchemaUrn{client.ENUMTHIRDPARTYALERTHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ALERT_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyAlertHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Alert Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AlertHandlerApi.AddAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAlertHandlerRequest(
		client.AddThirdPartyAlertHandlerRequestAsAddAlertHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AlertHandlerApi.AddAlertHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Alert Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state alertHandlerResourceModel
	readThirdPartyAlertHandlerResponse(ctx, addResponse.ThirdPartyAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *alertHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan alertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *alertHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "smtp" {
		state, err = r.CreateSmtpAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "jmx" {
		state, err = r.CreateJmxAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "snmp" {
		state, err = r.CreateSnmpAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "twilio" {
		state, err = r.CreateTwilioAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "error-log" {
		state, err = r.CreateErrorLogAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "snmp-sub-agent" {
		state, err = r.CreateSnmpSubAgentAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "exec" {
		state, err = r.CreateExecAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyAlertHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
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
func (r *defaultAlertHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultAlertHandlerResourceModel
	if readResponse.OutputAlertHandlerResponse != nil {
		readOutputAlertHandlerResponseDefault(ctx, readResponse.OutputAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SmtpAlertHandlerResponse != nil {
		readSmtpAlertHandlerResponseDefault(ctx, readResponse.SmtpAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JmxAlertHandlerResponse != nil {
		readJmxAlertHandlerResponseDefault(ctx, readResponse.JmxAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAlertHandlerResponse != nil {
		readGroovyScriptedAlertHandlerResponseDefault(ctx, readResponse.GroovyScriptedAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomAlertHandlerResponse != nil {
		readCustomAlertHandlerResponseDefault(ctx, readResponse.CustomAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpAlertHandlerResponse != nil {
		readSnmpAlertHandlerResponseDefault(ctx, readResponse.SnmpAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TwilioAlertHandlerResponse != nil {
		readTwilioAlertHandlerResponseDefault(ctx, readResponse.TwilioAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogAlertHandlerResponse != nil {
		readErrorLogAlertHandlerResponseDefault(ctx, readResponse.ErrorLogAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpSubAgentAlertHandlerResponse != nil {
		readSnmpSubAgentAlertHandlerResponseDefault(ctx, readResponse.SnmpSubAgentAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExecAlertHandlerResponse != nil {
		readExecAlertHandlerResponseDefault(ctx, readResponse.ExecAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAlertHandlerResponse != nil {
		readThirdPartyAlertHandlerResponseDefault(ctx, readResponse.ThirdPartyAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AlertHandlerApi.UpdateAlertHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createAlertHandlerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.OutputAlertHandlerResponse != nil {
			readOutputAlertHandlerResponseDefault(ctx, updateResponse.OutputAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SmtpAlertHandlerResponse != nil {
			readSmtpAlertHandlerResponseDefault(ctx, updateResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JmxAlertHandlerResponse != nil {
			readJmxAlertHandlerResponseDefault(ctx, updateResponse.JmxAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedAlertHandlerResponse != nil {
			readGroovyScriptedAlertHandlerResponseDefault(ctx, updateResponse.GroovyScriptedAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomAlertHandlerResponse != nil {
			readCustomAlertHandlerResponseDefault(ctx, updateResponse.CustomAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SnmpAlertHandlerResponse != nil {
			readSnmpAlertHandlerResponseDefault(ctx, updateResponse.SnmpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TwilioAlertHandlerResponse != nil {
			readTwilioAlertHandlerResponseDefault(ctx, updateResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ErrorLogAlertHandlerResponse != nil {
			readErrorLogAlertHandlerResponseDefault(ctx, updateResponse.ErrorLogAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SnmpSubAgentAlertHandlerResponse != nil {
			readSnmpSubAgentAlertHandlerResponseDefault(ctx, updateResponse.SnmpSubAgentAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExecAlertHandlerResponse != nil {
			readExecAlertHandlerResponseDefault(ctx, updateResponse.ExecAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAlertHandlerResponse != nil {
			readThirdPartyAlertHandlerResponseDefault(ctx, updateResponse.ThirdPartyAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *alertHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state alertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Alert Handler", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Alert Handler", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SmtpAlertHandlerResponse != nil {
		readSmtpAlertHandlerResponse(ctx, readResponse.SmtpAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JmxAlertHandlerResponse != nil {
		readJmxAlertHandlerResponse(ctx, readResponse.JmxAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAlertHandlerResponse != nil {
		readGroovyScriptedAlertHandlerResponse(ctx, readResponse.GroovyScriptedAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpAlertHandlerResponse != nil {
		readSnmpAlertHandlerResponse(ctx, readResponse.SnmpAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TwilioAlertHandlerResponse != nil {
		readTwilioAlertHandlerResponse(ctx, readResponse.TwilioAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogAlertHandlerResponse != nil {
		readErrorLogAlertHandlerResponse(ctx, readResponse.ErrorLogAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpSubAgentAlertHandlerResponse != nil {
		readSnmpSubAgentAlertHandlerResponse(ctx, readResponse.SnmpSubAgentAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ExecAlertHandlerResponse != nil {
		readExecAlertHandlerResponse(ctx, readResponse.ExecAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAlertHandlerResponse != nil {
		readThirdPartyAlertHandlerResponse(ctx, readResponse.ThirdPartyAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultAlertHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultAlertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerApi.GetAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Alert Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.OutputAlertHandlerResponse != nil {
		readOutputAlertHandlerResponseDefault(ctx, readResponse.OutputAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomAlertHandlerResponse != nil {
		readCustomAlertHandlerResponseDefault(ctx, readResponse.CustomAlertHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *alertHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan alertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state alertHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AlertHandlerApi.UpdateAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createAlertHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SmtpAlertHandlerResponse != nil {
			readSmtpAlertHandlerResponse(ctx, updateResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JmxAlertHandlerResponse != nil {
			readJmxAlertHandlerResponse(ctx, updateResponse.JmxAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedAlertHandlerResponse != nil {
			readGroovyScriptedAlertHandlerResponse(ctx, updateResponse.GroovyScriptedAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SnmpAlertHandlerResponse != nil {
			readSnmpAlertHandlerResponse(ctx, updateResponse.SnmpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TwilioAlertHandlerResponse != nil {
			readTwilioAlertHandlerResponse(ctx, updateResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ErrorLogAlertHandlerResponse != nil {
			readErrorLogAlertHandlerResponse(ctx, updateResponse.ErrorLogAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SnmpSubAgentAlertHandlerResponse != nil {
			readSnmpSubAgentAlertHandlerResponse(ctx, updateResponse.SnmpSubAgentAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExecAlertHandlerResponse != nil {
			readExecAlertHandlerResponse(ctx, updateResponse.ExecAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAlertHandlerResponse != nil {
			readThirdPartyAlertHandlerResponse(ctx, updateResponse.ThirdPartyAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
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

func (r *defaultAlertHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultAlertHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultAlertHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AlertHandlerApi.UpdateAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createAlertHandlerOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AlertHandlerApi.UpdateAlertHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Alert Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.OutputAlertHandlerResponse != nil {
			readOutputAlertHandlerResponseDefault(ctx, updateResponse.OutputAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SmtpAlertHandlerResponse != nil {
			readSmtpAlertHandlerResponseDefault(ctx, updateResponse.SmtpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JmxAlertHandlerResponse != nil {
			readJmxAlertHandlerResponseDefault(ctx, updateResponse.JmxAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedAlertHandlerResponse != nil {
			readGroovyScriptedAlertHandlerResponseDefault(ctx, updateResponse.GroovyScriptedAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomAlertHandlerResponse != nil {
			readCustomAlertHandlerResponseDefault(ctx, updateResponse.CustomAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SnmpAlertHandlerResponse != nil {
			readSnmpAlertHandlerResponseDefault(ctx, updateResponse.SnmpAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TwilioAlertHandlerResponse != nil {
			readTwilioAlertHandlerResponseDefault(ctx, updateResponse.TwilioAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ErrorLogAlertHandlerResponse != nil {
			readErrorLogAlertHandlerResponseDefault(ctx, updateResponse.ErrorLogAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SnmpSubAgentAlertHandlerResponse != nil {
			readSnmpSubAgentAlertHandlerResponseDefault(ctx, updateResponse.SnmpSubAgentAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ExecAlertHandlerResponse != nil {
			readExecAlertHandlerResponseDefault(ctx, updateResponse.ExecAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAlertHandlerResponse != nil {
			readThirdPartyAlertHandlerResponseDefault(ctx, updateResponse.ThirdPartyAlertHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultAlertHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *alertHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state alertHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AlertHandlerApi.DeleteAlertHandlerExecute(r.apiClient.AlertHandlerApi.DeleteAlertHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Alert Handler", err, httpResp)
		return
	}
}

func (r *alertHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAlertHandler(ctx, req, resp)
}

func (r *defaultAlertHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAlertHandler(ctx, req, resp)
}

func importAlertHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
