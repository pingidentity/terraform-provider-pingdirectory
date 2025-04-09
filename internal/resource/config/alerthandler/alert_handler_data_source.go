// Copyright © 2025 Ping Identity Corporation

package alerthandler

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
	_ datasource.DataSource              = &alertHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &alertHandlerDataSource{}
)

// Create a Alert Handler data source
func NewAlertHandlerDataSource() datasource.DataSource {
	return &alertHandlerDataSource{}
}

// alertHandlerDataSource is the datasource implementation.
type alertHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *alertHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_alert_handler"
}

// Configure adds the provider configured client to the data source.
func (r *alertHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type alertHandlerDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	Command                           types.String `tfsdk:"command"`
	CommandTimeout                    types.String `tfsdk:"command_timeout"`
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

// GetSchema defines the schema for the datasource.
func (r *alertHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Alert Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Alert Handler resource. Options are ['output', 'smtp', 'jmx', 'groovy-scripted', 'custom', 'snmp', 'twilio', 'error-log', 'snmp-sub-agent', 'exec', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Alert Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Alert Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"command": schema.StringAttribute{
				Description: "Specifies the path of the command to execute, without any arguments. It must be an absolute path for reasons of security and reliability.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"command_timeout": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 10.1.0.0+. The maximum length of time this server will wait for the executed command to finish executing before forcibly terminating it.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Alert Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the Twilio service.",
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
			"sender_phone_number": schema.SetAttribute{
				Description: "The outgoing phone number to use for the messages. Values must be phone numbers you have obtained for use with your Twilio account.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"recipient_phone_number": schema.SetAttribute{
				Description: "The phone number to which alert notifications should be delivered.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"long_message_behavior": schema.StringAttribute{
				Description: "The behavior to use for alert messages that are longer than the 160-character size limit for SMS messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_host_name": schema.StringAttribute{
				Description: "Specifies the address of the SNMP agent to which traps will be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the SNMP agent to which traps will be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"community_name": schema.StringAttribute{
				Description: "Specifies the name of the community to which the traps will be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Alert Handler. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"output_location": schema.StringAttribute{
				Description: "The location to which alert messages will be written.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"sender_address": schema.StringAttribute{
				Description: "Specifies the email address to use as the sender for messages generated by this alert handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"recipient_address": schema.SetAttribute{
				Description: "Specifies an email address to which the messages should be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"message_subject": schema.StringAttribute{
				Description: "Specifies the subject that should be used for email messages generated by this alert handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"message_body": schema.StringAttribute{
				Description: "Specifies the body that should be used for email messages generated by this alert handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_monitor_data_filter": schema.StringAttribute{
				Description: "Include monitor entries that match this filter.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"output_format": schema.StringAttribute{
				Description: "The format to use when writing the alert messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Alert Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Alert Handler is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`output`, `groovy-scripted`, `custom`, `error-log`, `third-party`]: Indicates whether the server should attempt to invoke this Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `smtp`: Indicates whether the server should attempt to invoke this SMTP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `jmx`: Indicates whether the server should attempt to invoke this JMX Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `snmp`: Indicates whether the server should attempt to invoke this SNMP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `twilio`: Indicates whether the server should attempt to invoke this Twilio Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `snmp-sub-agent`: Indicates whether the server should attempt to invoke this SNMP Sub Agent Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated. When the `type` attribute is set to `exec`: Indicates whether the server should attempt to invoke this Exec Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`output`, `groovy-scripted`, `custom`, `error-log`, `third-party`]: Indicates whether the server should attempt to invoke this Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `smtp`: Indicates whether the server should attempt to invoke this SMTP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `jmx`: Indicates whether the server should attempt to invoke this JMX Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `snmp`: Indicates whether the server should attempt to invoke this SNMP Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `twilio`: Indicates whether the server should attempt to invoke this Twilio Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `snmp-sub-agent`: Indicates whether the server should attempt to invoke this SNMP Sub Agent Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.\n  - `exec`: Indicates whether the server should attempt to invoke this Exec Alert Handler in a background thread so that any potentially-expensive processing (e.g., performing network communication to deliver the alert notification) will not delay whatever processing the server was performing when the alert was generated.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"enabled_alert_severity": schema.SetAttribute{
				Description: "Specifies the alert severities for which this alert handler should be used. If no values are provided, then this alert handler will be enabled for alerts with any severity.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"enabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that are enabled for this alert handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"disabled_alert_type": schema.SetAttribute{
				Description: "Specifies the names of the alert types that are disabled for this alert handler.",
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

// Read a OutputAlertHandlerResponse object into the model struct
func readOutputAlertHandlerResponseDataSource(ctx context.Context, r *client.OutputAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("output")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumalertHandlerOutputLocationProp(r.OutputLocation), false)
	state.OutputFormat = internaltypes.StringTypeOrNil(
		client.StringPointerEnumalertHandlerOutputFormatProp(r.OutputFormat), false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a SmtpAlertHandlerResponse object into the model struct
func readSmtpAlertHandlerResponseDataSource(ctx context.Context, r *client.SmtpAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("smtp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.SenderAddress = types.StringValue(r.SenderAddress)
	state.RecipientAddress = internaltypes.GetStringSet(r.RecipientAddress)
	state.MessageSubject = types.StringValue(r.MessageSubject)
	state.MessageBody = types.StringValue(r.MessageBody)
	state.IncludeMonitorDataFilter = internaltypes.StringTypeOrNil(r.IncludeMonitorDataFilter, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a JmxAlertHandlerResponse object into the model struct
func readJmxAlertHandlerResponseDataSource(ctx context.Context, r *client.JmxAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jmx")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a GroovyScriptedAlertHandlerResponse object into the model struct
func readGroovyScriptedAlertHandlerResponseDataSource(ctx context.Context, r *client.GroovyScriptedAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a CustomAlertHandlerResponse object into the model struct
func readCustomAlertHandlerResponseDataSource(ctx context.Context, r *client.CustomAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a SnmpAlertHandlerResponse object into the model struct
func readSnmpAlertHandlerResponseDataSource(ctx context.Context, r *client.SnmpAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("snmp")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.CommunityName = types.StringValue(r.CommunityName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a TwilioAlertHandlerResponse object into the model struct
func readTwilioAlertHandlerResponseDataSource(ctx context.Context, r *client.TwilioAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("twilio")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.TwilioAccountSID = types.StringValue(r.TwilioAccountSID)
	state.TwilioAuthTokenPassphraseProvider = internaltypes.StringTypeOrNil(r.TwilioAuthTokenPassphraseProvider, false)
	state.SenderPhoneNumber = internaltypes.GetStringSet(r.SenderPhoneNumber)
	state.RecipientPhoneNumber = internaltypes.GetStringSet(r.RecipientPhoneNumber)
	state.LongMessageBehavior = types.StringValue(r.LongMessageBehavior.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a ErrorLogAlertHandlerResponse object into the model struct
func readErrorLogAlertHandlerResponseDataSource(ctx context.Context, r *client.ErrorLogAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("error-log")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a SnmpSubAgentAlertHandlerResponse object into the model struct
func readSnmpSubAgentAlertHandlerResponseDataSource(ctx context.Context, r *client.SnmpSubAgentAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("snmp-sub-agent")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a ExecAlertHandlerResponse object into the model struct
func readExecAlertHandlerResponseDataSource(ctx context.Context, r *client.ExecAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exec")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Command = types.StringValue(r.Command)
	state.CommandTimeout = internaltypes.StringTypeOrNil(r.CommandTimeout, false)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read a ThirdPartyAlertHandlerResponse object into the model struct
func readThirdPartyAlertHandlerResponseDataSource(ctx context.Context, r *client.ThirdPartyAlertHandlerResponse, state *alertHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Asynchronous = internaltypes.BoolTypeOrNil(r.Asynchronous)
	state.EnabledAlertSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertSeverityProp(r.EnabledAlertSeverity))
	state.EnabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerEnabledAlertTypeProp(r.EnabledAlertType))
	state.DisabledAlertType = internaltypes.GetStringSet(
		client.StringSliceEnumalertHandlerDisabledAlertTypeProp(r.DisabledAlertType))
}

// Read resource information
func (r *alertHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state alertHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AlertHandlerAPI.GetAlertHandler(
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
		readOutputAlertHandlerResponseDataSource(ctx, readResponse.OutputAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SmtpAlertHandlerResponse != nil {
		readSmtpAlertHandlerResponseDataSource(ctx, readResponse.SmtpAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JmxAlertHandlerResponse != nil {
		readJmxAlertHandlerResponseDataSource(ctx, readResponse.JmxAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAlertHandlerResponse != nil {
		readGroovyScriptedAlertHandlerResponseDataSource(ctx, readResponse.GroovyScriptedAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CustomAlertHandlerResponse != nil {
		readCustomAlertHandlerResponseDataSource(ctx, readResponse.CustomAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpAlertHandlerResponse != nil {
		readSnmpAlertHandlerResponseDataSource(ctx, readResponse.SnmpAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TwilioAlertHandlerResponse != nil {
		readTwilioAlertHandlerResponseDataSource(ctx, readResponse.TwilioAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ErrorLogAlertHandlerResponse != nil {
		readErrorLogAlertHandlerResponseDataSource(ctx, readResponse.ErrorLogAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SnmpSubAgentAlertHandlerResponse != nil {
		readSnmpSubAgentAlertHandlerResponseDataSource(ctx, readResponse.SnmpSubAgentAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ExecAlertHandlerResponse != nil {
		readExecAlertHandlerResponseDataSource(ctx, readResponse.ExecAlertHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAlertHandlerResponse != nil {
		readThirdPartyAlertHandlerResponseDataSource(ctx, readResponse.ThirdPartyAlertHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
