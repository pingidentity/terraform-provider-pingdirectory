package logpublisher

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	_ resource.Resource                = &logPublisherResource{}
	_ resource.ResourceWithConfigure   = &logPublisherResource{}
	_ resource.ResourceWithImportState = &logPublisherResource{}
	_ resource.Resource                = &defaultLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultLogPublisherResource{}
)

// Create a Log Publisher resource
func NewLogPublisherResource() resource.Resource {
	return &logPublisherResource{}
}

func NewDefaultLogPublisherResource() resource.Resource {
	return &defaultLogPublisherResource{}
}

// logPublisherResource is the resource implementation.
type logPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLogPublisherResource is the resource implementation.
type defaultLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *logPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_publisher"
}

func (r *defaultLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *logPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logPublisherResourceModel struct {
	Id                                                  types.String `tfsdk:"id"`
	Name                                                types.String `tfsdk:"name"`
	Notifications                                       types.Set    `tfsdk:"notifications"`
	RequiredActions                                     types.Set    `tfsdk:"required_actions"`
	Type                                                types.String `tfsdk:"type"`
	ScriptClass                                         types.String `tfsdk:"script_class"`
	Server                                              types.String `tfsdk:"server"`
	LogFieldMapping                                     types.String `tfsdk:"log_field_mapping"`
	LogTableName                                        types.String `tfsdk:"log_table_name"`
	OutputLocation                                      types.String `tfsdk:"output_location"`
	LogFile                                             types.String `tfsdk:"log_file"`
	LogFilePermissions                                  types.String `tfsdk:"log_file_permissions"`
	RotationPolicy                                      types.Set    `tfsdk:"rotation_policy"`
	RotationListener                                    types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                                     types.Set    `tfsdk:"retention_policy"`
	CompressionMechanism                                types.String `tfsdk:"compression_mechanism"`
	ScriptArgument                                      types.Set    `tfsdk:"script_argument"`
	SignLog                                             types.Bool   `tfsdk:"sign_log"`
	TimestampPrecision                                  types.String `tfsdk:"timestamp_precision"`
	EncryptLog                                          types.Bool   `tfsdk:"encrypt_log"`
	EncryptionSettingsDefinitionID                      types.String `tfsdk:"encryption_settings_definition_id"`
	Append                                              types.Bool   `tfsdk:"append"`
	ObscureSensitiveContent                             types.Bool   `tfsdk:"obscure_sensitive_content"`
	ExtensionClass                                      types.String `tfsdk:"extension_class"`
	DebugACIEnabled                                     types.Bool   `tfsdk:"debug_aci_enabled"`
	DefaultDebugLevel                                   types.String `tfsdk:"default_debug_level"`
	LogRequestHeaders                                   types.String `tfsdk:"log_request_headers"`
	SuppressedRequestHeaderName                         types.Set    `tfsdk:"suppressed_request_header_name"`
	LogResponseHeaders                                  types.String `tfsdk:"log_response_headers"`
	SuppressedResponseHeaderName                        types.Set    `tfsdk:"suppressed_response_header_name"`
	LogRequestAuthorizationType                         types.Bool   `tfsdk:"log_request_authorization_type"`
	LogRequestCookieNames                               types.Bool   `tfsdk:"log_request_cookie_names"`
	LogResponseCookieNames                              types.Bool   `tfsdk:"log_response_cookie_names"`
	LogRequestParameters                                types.String `tfsdk:"log_request_parameters"`
	LogRequestProtocol                                  types.Bool   `tfsdk:"log_request_protocol"`
	SuppressedRequestParameterName                      types.Set    `tfsdk:"suppressed_request_parameter_name"`
	LogRedirectURI                                      types.Bool   `tfsdk:"log_redirect_uri"`
	DefaultDebugCategory                                types.Set    `tfsdk:"default_debug_category"`
	DefaultOmitMethodEntryArguments                     types.Bool   `tfsdk:"default_omit_method_entry_arguments"`
	DefaultOmitMethodReturnValue                        types.Bool   `tfsdk:"default_omit_method_return_value"`
	DefaultIncludeThrowableCause                        types.Bool   `tfsdk:"default_include_throwable_cause"`
	DefaultThrowableStackFrames                         types.Int64  `tfsdk:"default_throwable_stack_frames"`
	ExtensionArgument                                   types.Set    `tfsdk:"extension_argument"`
	SyslogExternalServer                                types.Set    `tfsdk:"syslog_external_server"`
	IncludeRequestDetailsInResultMessages               types.Bool   `tfsdk:"include_request_details_in_result_messages"`
	LogAssuranceCompleted                               types.Bool   `tfsdk:"log_assurance_completed"`
	DebugMessageType                                    types.Set    `tfsdk:"debug_message_type"`
	HttpMessageType                                     types.Set    `tfsdk:"http_message_type"`
	AccessTokenValidatorMessageType                     types.Set    `tfsdk:"access_token_validator_message_type"`
	IdTokenValidatorMessageType                         types.Set    `tfsdk:"id_token_validator_message_type"`
	ScimMessageType                                     types.Set    `tfsdk:"scim_message_type"`
	ConsentMessageType                                  types.Set    `tfsdk:"consent_message_type"`
	DirectoryRESTAPIMessageType                         types.Set    `tfsdk:"directory_rest_api_message_type"`
	ExtensionMessageType                                types.Set    `tfsdk:"extension_message_type"`
	IncludePathPattern                                  types.Set    `tfsdk:"include_path_pattern"`
	ExcludePathPattern                                  types.Set    `tfsdk:"exclude_path_pattern"`
	ServerHostName                                      types.String `tfsdk:"server_host_name"`
	BufferSize                                          types.String `tfsdk:"buffer_size"`
	ServerPort                                          types.Int64  `tfsdk:"server_port"`
	MinIncludedOperationProcessingTime                  types.String `tfsdk:"min_included_operation_processing_time"`
	MinIncludedPhaseTimeNanos                           types.Int64  `tfsdk:"min_included_phase_time_nanos"`
	TimeInterval                                        types.String `tfsdk:"time_interval"`
	IncludeRequestDetailsInSearchEntryMessages          types.Bool   `tfsdk:"include_request_details_in_search_entry_messages"`
	IncludeRequestDetailsInSearchReferenceMessages      types.Bool   `tfsdk:"include_request_details_in_search_reference_messages"`
	IncludeRequestDetailsInIntermediateResponseMessages types.Bool   `tfsdk:"include_request_details_in_intermediate_response_messages"`
	IncludeResultCodeNames                              types.Bool   `tfsdk:"include_result_code_names"`
	IncludeExtendedSearchRequestDetails                 types.Bool   `tfsdk:"include_extended_search_request_details"`
	IncludeAddAttributeNames                            types.Bool   `tfsdk:"include_add_attribute_names"`
	IncludeModifyAttributeNames                         types.Bool   `tfsdk:"include_modify_attribute_names"`
	IncludeSearchEntryAttributeNames                    types.Bool   `tfsdk:"include_search_entry_attribute_names"`
	LogConnects                                         types.Bool   `tfsdk:"log_connects"`
	LogDisconnects                                      types.Bool   `tfsdk:"log_disconnects"`
	MaxStringLength                                     types.Int64  `tfsdk:"max_string_length"`
	GenerifyMessageStringsWhenPossible                  types.Bool   `tfsdk:"generify_message_strings_when_possible"`
	SyslogFacility                                      types.String `tfsdk:"syslog_facility"`
	LogFieldBehavior                                    types.String `tfsdk:"log_field_behavior"`
	LogClientCertificates                               types.Bool   `tfsdk:"log_client_certificates"`
	LogRequests                                         types.Bool   `tfsdk:"log_requests"`
	LogResults                                          types.Bool   `tfsdk:"log_results"`
	LogSearchEntries                                    types.Bool   `tfsdk:"log_search_entries"`
	LogSearchReferences                                 types.Bool   `tfsdk:"log_search_references"`
	LogIntermediateResponses                            types.Bool   `tfsdk:"log_intermediate_responses"`
	AutoFlush                                           types.Bool   `tfsdk:"auto_flush"`
	Asynchronous                                        types.Bool   `tfsdk:"asynchronous"`
	CorrelateRequestsAndResults                         types.Bool   `tfsdk:"correlate_requests_and_results"`
	SyslogSeverity                                      types.String `tfsdk:"syslog_severity"`
	DefaultSeverity                                     types.Set    `tfsdk:"default_severity"`
	OverrideSeverity                                    types.Set    `tfsdk:"override_severity"`
	SearchEntryCriteria                                 types.String `tfsdk:"search_entry_criteria"`
	SearchReferenceCriteria                             types.String `tfsdk:"search_reference_criteria"`
	SyslogMessageHostName                               types.String `tfsdk:"syslog_message_host_name"`
	SyslogMessageApplicationName                        types.String `tfsdk:"syslog_message_application_name"`
	QueueSize                                           types.Int64  `tfsdk:"queue_size"`
	WriteMultiLineMessages                              types.Bool   `tfsdk:"write_multi_line_messages"`
	UseReversibleForm                                   types.Bool   `tfsdk:"use_reversible_form"`
	SoftDeleteEntryAuditBehavior                        types.String `tfsdk:"soft_delete_entry_audit_behavior"`
	IncludeOperationPurposeRequestControl               types.Bool   `tfsdk:"include_operation_purpose_request_control"`
	IncludeIntermediateClientRequestControl             types.Bool   `tfsdk:"include_intermediate_client_request_control"`
	ObscureAttribute                                    types.Set    `tfsdk:"obscure_attribute"`
	ExcludeAttribute                                    types.Set    `tfsdk:"exclude_attribute"`
	SuppressInternalOperations                          types.Bool   `tfsdk:"suppress_internal_operations"`
	IncludeProductName                                  types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                                 types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                                    types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                                     types.Bool   `tfsdk:"include_thread_id"`
	IncludeRequesterDN                                  types.Bool   `tfsdk:"include_requester_dn"`
	IncludeRequesterIPAddress                           types.Bool   `tfsdk:"include_requester_ip_address"`
	IncludeRequestControls                              types.Bool   `tfsdk:"include_request_controls"`
	IncludeResponseControls                             types.Bool   `tfsdk:"include_response_controls"`
	IncludeReplicationChangeID                          types.Bool   `tfsdk:"include_replication_change_id"`
	LogSecurityNegotiation                              types.Bool   `tfsdk:"log_security_negotiation"`
	SuppressReplicationOperations                       types.Bool   `tfsdk:"suppress_replication_operations"`
	ConnectionCriteria                                  types.String `tfsdk:"connection_criteria"`
	RequestCriteria                                     types.String `tfsdk:"request_criteria"`
	ResultCriteria                                      types.String `tfsdk:"result_criteria"`
	Description                                         types.String `tfsdk:"description"`
	Enabled                                             types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                                types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *logPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logPublisherSchema(ctx, req, resp, false)
}

func (r *defaultLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	logPublisherSchema(ctx, req, resp, true)
}

func logPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Publisher resource. Options are ['syslog-json-audit', 'syslog-based-error', 'third-party-file-based-access', 'operation-timing-access', 'third-party-http-operation', 'admin-alert-access', 'file-based-trace', 'jdbc-based-error', 'jdbc-based-access', 'common-log-file-http-operation', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'debug-access', 'syslog-json-http-operation', 'third-party-access', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'groovy-scripted-file-based-error', 'syslog-json-access', 'groovy-scripted-access', 'third-party-file-based-error', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'groovy-scripted-error', 'file-based-json-http-operation', 'syslog-json-error', 'groovy-scripted-http-operation']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"syslog-json-audit", "syslog-based-error", "third-party-file-based-access", "operation-timing-access", "third-party-http-operation", "admin-alert-access", "file-based-trace", "jdbc-based-error", "jdbc-based-access", "common-log-file-http-operation", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-debug", "file-based-error", "third-party-error", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "syslog-json-http-operation", "third-party-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "syslog-json-access", "groovy-scripted-access", "third-party-file-based-error", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-access", "groovy-scripted-error", "file-based-json-http-operation", "syslog-json-error", "groovy-scripted-http-operation"}...),
				},
			},
			"script_class": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `groovy-scripted-file-based-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Access Log Publisher. When the `type` attribute is set to `groovy-scripted-file-based-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Error Log Publisher. When the `type` attribute is set to `groovy-scripted-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Access Log Publisher. When the `type` attribute is set to `groovy-scripted-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Error Log Publisher. When the `type` attribute is set to `groovy-scripted-http-operation`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Operation Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `groovy-scripted-file-based-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Access Log Publisher.\n  - `groovy-scripted-file-based-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Error Log Publisher.\n  - `groovy-scripted-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Access Log Publisher.\n  - `groovy-scripted-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Error Log Publisher.\n  - `groovy-scripted-http-operation`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Operation Log Publisher.",
				Optional:            true,
			},
			"server": schema.StringAttribute{
				Description: "The JDBC-based Database Server to use for a connection.",
				Optional:    true,
			},
			"log_field_mapping": schema.StringAttribute{
				Description: "The log field mapping associates loggable fields to database column names. The table name is not part of this mapping.",
				Optional:    true,
			},
			"log_table_name": schema.StringAttribute{
				Description: "The table name to log entries to the database server.",
				Optional:    true,
				Computed:    true,
			},
			"output_location": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `console-json-error`: Specifies the output stream to which JSON-formatted error log messages should be written. When the `type` attribute is set to `console-json-audit`: Specifies the output stream to which JSON-formatted audit log messages should be written. When the `type` attribute is set to `console-json-http-operation`: Specifies the output stream to which JSON-formatted HTTP operation log messages should be written. When the `type` attribute is set to `console-json-access`: Specifies the output stream to which JSON-formatted access log messages should be written.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `console-json-error`: Specifies the output stream to which JSON-formatted error log messages should be written.\n  - `console-json-audit`: Specifies the output stream to which JSON-formatted audit log messages should be written.\n  - `console-json-http-operation`: Specifies the output stream to which JSON-formatted HTTP operation log messages should be written.\n  - `console-json-access`: Specifies the output stream to which JSON-formatted access log messages should be written.",
				Optional:            true,
				Computed:            true,
			},
			"log_file": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The file name to use for the log files generated by the Third Party File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `operation-timing-access`: The file name to use for the log files generated by the Operation Timing Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-trace`: The file name to use for the log files generated by the File Based Trace Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `common-log-file-http-operation`: The file name to use for the log files generated by the Common Log File HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-json-audit`: The file name to use for the log files generated by the File Based JSON Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-debug`: The file name to use for the log files generated by the File Based Debug Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-error`: The file name to use for the log files generated by the File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `detailed-http-operation`: The file name to use for the log files generated by the Detailed HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `json-access`: The file name to use for the log files generated by the JSON Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `debug-access`: The file name to use for the log files generated by the Debug Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-audit`: The file name to use for the log files generated by the File Based Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `json-error`: The file name to use for the log files generated by the JSON Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `groovy-scripted-file-based-access`: The file name to use for the log files generated by the Scripted File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `groovy-scripted-file-based-error`: The file name to use for the log files generated by the Scripted File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `third-party-file-based-error`: The file name to use for the log files generated by the Third Party File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-access`: The file name to use for the log files generated by the File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-json-http-operation`: The file name to use for the log files generated by the File Based JSON HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The file name to use for the log files generated by the Third Party File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `operation-timing-access`: The file name to use for the log files generated by the Operation Timing Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-trace`: The file name to use for the log files generated by the File Based Trace Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `common-log-file-http-operation`: The file name to use for the log files generated by the Common Log File HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-json-audit`: The file name to use for the log files generated by the File Based JSON Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-debug`: The file name to use for the log files generated by the File Based Debug Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-error`: The file name to use for the log files generated by the File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `detailed-http-operation`: The file name to use for the log files generated by the Detailed HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `json-access`: The file name to use for the log files generated by the JSON Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `debug-access`: The file name to use for the log files generated by the Debug Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-audit`: The file name to use for the log files generated by the File Based Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `json-error`: The file name to use for the log files generated by the JSON Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `groovy-scripted-file-based-access`: The file name to use for the log files generated by the Scripted File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `groovy-scripted-file-based-error`: The file name to use for the log files generated by the Scripted File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `third-party-file-based-error`: The file name to use for the log files generated by the Third Party File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-access`: The file name to use for the log files generated by the File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-json-http-operation`: The file name to use for the log files generated by the File Based JSON HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Optional:            true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The UNIX permissions of the log files created by this Third Party File Based Access Log Publisher. When the `type` attribute is set to `operation-timing-access`: The UNIX permissions of the log files created by this Operation Timing Access Log Publisher. When the `type` attribute is set to `file-based-trace`: The UNIX permissions of the log files created by this File Based Trace Log Publisher. When the `type` attribute is set to `common-log-file-http-operation`: The UNIX permissions of the log files created by this Common Log File HTTP Operation Log Publisher. When the `type` attribute is set to `file-based-json-audit`: The UNIX permissions of the log files created by this File Based JSON Audit Log Publisher. When the `type` attribute is set to `file-based-debug`: The UNIX permissions of the log files created by this File Based Debug Log Publisher. When the `type` attribute is set to `file-based-error`: The UNIX permissions of the log files created by this File Based Error Log Publisher. When the `type` attribute is set to `detailed-http-operation`: The UNIX permissions of the log files created by this Detailed HTTP Operation Log Publisher. When the `type` attribute is set to `json-access`: The UNIX permissions of the log files created by this JSON Access Log Publisher. When the `type` attribute is set to `debug-access`: The UNIX permissions of the log files created by this Debug Access Log Publisher. When the `type` attribute is set to `file-based-audit`: The UNIX permissions of the log files created by this File Based Audit Log Publisher. When the `type` attribute is set to `json-error`: The UNIX permissions of the log files created by this JSON Error Log Publisher. When the `type` attribute is set to `groovy-scripted-file-based-access`: The UNIX permissions of the log files created by this Scripted File Based Access Log Publisher. When the `type` attribute is set to `groovy-scripted-file-based-error`: The UNIX permissions of the log files created by this Scripted File Based Error Log Publisher. When the `type` attribute is set to `third-party-file-based-error`: The UNIX permissions of the log files created by this Third Party File Based Error Log Publisher. When the `type` attribute is set to `file-based-access`: The UNIX permissions of the log files created by this File Based Access Log Publisher. When the `type` attribute is set to `file-based-json-http-operation`: The UNIX permissions of the log files created by this File Based JSON HTTP Operation Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The UNIX permissions of the log files created by this Third Party File Based Access Log Publisher.\n  - `operation-timing-access`: The UNIX permissions of the log files created by this Operation Timing Access Log Publisher.\n  - `file-based-trace`: The UNIX permissions of the log files created by this File Based Trace Log Publisher.\n  - `common-log-file-http-operation`: The UNIX permissions of the log files created by this Common Log File HTTP Operation Log Publisher.\n  - `file-based-json-audit`: The UNIX permissions of the log files created by this File Based JSON Audit Log Publisher.\n  - `file-based-debug`: The UNIX permissions of the log files created by this File Based Debug Log Publisher.\n  - `file-based-error`: The UNIX permissions of the log files created by this File Based Error Log Publisher.\n  - `detailed-http-operation`: The UNIX permissions of the log files created by this Detailed HTTP Operation Log Publisher.\n  - `json-access`: The UNIX permissions of the log files created by this JSON Access Log Publisher.\n  - `debug-access`: The UNIX permissions of the log files created by this Debug Access Log Publisher.\n  - `file-based-audit`: The UNIX permissions of the log files created by this File Based Audit Log Publisher.\n  - `json-error`: The UNIX permissions of the log files created by this JSON Error Log Publisher.\n  - `groovy-scripted-file-based-access`: The UNIX permissions of the log files created by this Scripted File Based Access Log Publisher.\n  - `groovy-scripted-file-based-error`: The UNIX permissions of the log files created by this Scripted File Based Error Log Publisher.\n  - `third-party-file-based-error`: The UNIX permissions of the log files created by this Third Party File Based Error Log Publisher.\n  - `file-based-access`: The UNIX permissions of the log files created by this File Based Access Log Publisher.\n  - `file-based-json-http-operation`: The UNIX permissions of the log files created by this File Based JSON HTTP Operation Log Publisher.",
				Optional:            true,
				Computed:            true,
			},
			"rotation_policy": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The rotation policy to use for the Third Party File Based Access Log Publisher . When the `type` attribute is set to `operation-timing-access`: The rotation policy to use for the Operation Timing Access Log Publisher . When the `type` attribute is set to `file-based-trace`: The rotation policy to use for the File Based Trace Log Publisher . When the `type` attribute is set to `common-log-file-http-operation`: The rotation policy to use for the Common Log File HTTP Operation Log Publisher . When the `type` attribute is set to `file-based-json-audit`: The rotation policy to use for the File Based JSON Audit Log Publisher . When the `type` attribute is set to `file-based-debug`: The rotation policy to use for the File Based Debug Log Publisher . When the `type` attribute is set to `file-based-error`: The rotation policy to use for the File Based Error Log Publisher . When the `type` attribute is set to `detailed-http-operation`: The rotation policy to use for the Detailed HTTP Operation Log Publisher . When the `type` attribute is set to `json-access`: The rotation policy to use for the JSON Access Log Publisher . When the `type` attribute is set to `debug-access`: The rotation policy to use for the Debug Access Log Publisher . When the `type` attribute is set to `file-based-audit`: The rotation policy to use for the File Based Audit Log Publisher . When the `type` attribute is set to `json-error`: The rotation policy to use for the JSON Error Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-access`: The rotation policy to use for the Scripted File Based Access Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-error`: The rotation policy to use for the Scripted File Based Error Log Publisher . When the `type` attribute is set to `third-party-file-based-error`: The rotation policy to use for the Third Party File Based Error Log Publisher . When the `type` attribute is set to `file-based-access`: The rotation policy to use for the File Based Access Log Publisher . When the `type` attribute is set to `file-based-json-http-operation`: The rotation policy to use for the File Based JSON HTTP Operation Log Publisher .",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The rotation policy to use for the Third Party File Based Access Log Publisher .\n  - `operation-timing-access`: The rotation policy to use for the Operation Timing Access Log Publisher .\n  - `file-based-trace`: The rotation policy to use for the File Based Trace Log Publisher .\n  - `common-log-file-http-operation`: The rotation policy to use for the Common Log File HTTP Operation Log Publisher .\n  - `file-based-json-audit`: The rotation policy to use for the File Based JSON Audit Log Publisher .\n  - `file-based-debug`: The rotation policy to use for the File Based Debug Log Publisher .\n  - `file-based-error`: The rotation policy to use for the File Based Error Log Publisher .\n  - `detailed-http-operation`: The rotation policy to use for the Detailed HTTP Operation Log Publisher .\n  - `json-access`: The rotation policy to use for the JSON Access Log Publisher .\n  - `debug-access`: The rotation policy to use for the Debug Access Log Publisher .\n  - `file-based-audit`: The rotation policy to use for the File Based Audit Log Publisher .\n  - `json-error`: The rotation policy to use for the JSON Error Log Publisher .\n  - `groovy-scripted-file-based-access`: The rotation policy to use for the Scripted File Based Access Log Publisher .\n  - `groovy-scripted-file-based-error`: The rotation policy to use for the Scripted File Based Error Log Publisher .\n  - `third-party-file-based-error`: The rotation policy to use for the Third Party File Based Error Log Publisher .\n  - `file-based-access`: The rotation policy to use for the File Based Access Log Publisher .\n  - `file-based-json-http-operation`: The rotation policy to use for the File Based JSON HTTP Operation Log Publisher .",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The retention policy to use for the Third Party File Based Access Log Publisher . When the `type` attribute is set to `operation-timing-access`: The retention policy to use for the Operation Timing Access Log Publisher . When the `type` attribute is set to `file-based-trace`: The retention policy to use for the File Based Trace Log Publisher . When the `type` attribute is set to `common-log-file-http-operation`: The retention policy to use for the Common Log File HTTP Operation Log Publisher . When the `type` attribute is set to `file-based-json-audit`: The retention policy to use for the File Based JSON Audit Log Publisher . When the `type` attribute is set to `file-based-debug`: The retention policy to use for the File Based Debug Log Publisher . When the `type` attribute is set to `file-based-error`: The retention policy to use for the File Based Error Log Publisher . When the `type` attribute is set to `detailed-http-operation`: The retention policy to use for the Detailed HTTP Operation Log Publisher . When the `type` attribute is set to `json-access`: The retention policy to use for the JSON Access Log Publisher . When the `type` attribute is set to `debug-access`: The retention policy to use for the Debug Access Log Publisher . When the `type` attribute is set to `file-based-audit`: The retention policy to use for the File Based Audit Log Publisher . When the `type` attribute is set to `json-error`: The retention policy to use for the JSON Error Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-access`: The retention policy to use for the Scripted File Based Access Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-error`: The retention policy to use for the Scripted File Based Error Log Publisher . When the `type` attribute is set to `third-party-file-based-error`: The retention policy to use for the Third Party File Based Error Log Publisher . When the `type` attribute is set to `file-based-access`: The retention policy to use for the File Based Access Log Publisher . When the `type` attribute is set to `file-based-json-http-operation`: The retention policy to use for the File Based JSON HTTP Operation Log Publisher .",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The retention policy to use for the Third Party File Based Access Log Publisher .\n  - `operation-timing-access`: The retention policy to use for the Operation Timing Access Log Publisher .\n  - `file-based-trace`: The retention policy to use for the File Based Trace Log Publisher .\n  - `common-log-file-http-operation`: The retention policy to use for the Common Log File HTTP Operation Log Publisher .\n  - `file-based-json-audit`: The retention policy to use for the File Based JSON Audit Log Publisher .\n  - `file-based-debug`: The retention policy to use for the File Based Debug Log Publisher .\n  - `file-based-error`: The retention policy to use for the File Based Error Log Publisher .\n  - `detailed-http-operation`: The retention policy to use for the Detailed HTTP Operation Log Publisher .\n  - `json-access`: The retention policy to use for the JSON Access Log Publisher .\n  - `debug-access`: The retention policy to use for the Debug Access Log Publisher .\n  - `file-based-audit`: The retention policy to use for the File Based Audit Log Publisher .\n  - `json-error`: The retention policy to use for the JSON Error Log Publisher .\n  - `groovy-scripted-file-based-access`: The retention policy to use for the Scripted File Based Access Log Publisher .\n  - `groovy-scripted-file-based-error`: The retention policy to use for the Scripted File Based Error Log Publisher .\n  - `third-party-file-based-error`: The retention policy to use for the Third Party File Based Error Log Publisher .\n  - `file-based-access`: The retention policy to use for the File Based Access Log Publisher .\n  - `file-based-json-http-operation`: The retention policy to use for the File Based JSON HTTP Operation Log Publisher .",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Optional:    true,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `groovy-scripted-file-based-access`: The set of arguments used to customize the behavior for the Scripted File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-file-based-error`: The set of arguments used to customize the behavior for the Scripted File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-access`: The set of arguments used to customize the behavior for the Scripted Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-error`: The set of arguments used to customize the behavior for the Scripted Error Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-http-operation`: The set of arguments used to customize the behavior for the Scripted HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `groovy-scripted-file-based-access`: The set of arguments used to customize the behavior for the Scripted File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-file-based-error`: The set of arguments used to customize the behavior for the Scripted File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-access`: The set of arguments used to customize the behavior for the Scripted Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-error`: The set of arguments used to customize the behavior for the Scripted Error Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-http-operation`: The set of arguments used to customize the behavior for the Scripted HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"sign_log": schema.BoolAttribute{
				Description: "Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.",
				Optional:    true,
				Computed:    true,
			},
			"timestamp_precision": schema.StringAttribute{
				Description: "Specifies the smallest time unit to be included in timestamps.",
				Optional:    true,
				Computed:    true,
			},
			"encrypt_log": schema.BoolAttribute{
				Description: "Indicates whether log files should be encrypted so that their content is not available to unauthorized users.",
				Optional:    true,
				Computed:    true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
				Optional:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Optional:    true,
				Computed:    true,
			},
			"obscure_sensitive_content": schema.BoolAttribute{
				Description: "Indicates whether the resulting log file should attempt to obscure content that may be considered sensitive. This primarily includes the credentials for bind requests, the values of password modify extended requests and responses, and the values of any attributes specified in the obscure-attribute property. Note that the use of this option does not guarantee no sensitive information will be exposed, so the log output should still be carefully guarded.",
				Optional:    true,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Access Log Publisher. When the `type` attribute is set to `third-party-http-operation`: The fully-qualified name of the Java class providing the logic for the Third Party HTTP Operation Log Publisher. When the `type` attribute is set to `third-party-error`: The fully-qualified name of the Java class providing the logic for the Third Party Error Log Publisher. When the `type` attribute is set to `third-party-access`: The fully-qualified name of the Java class providing the logic for the Third Party Access Log Publisher. When the `type` attribute is set to `third-party-file-based-error`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Access Log Publisher.\n  - `third-party-http-operation`: The fully-qualified name of the Java class providing the logic for the Third Party HTTP Operation Log Publisher.\n  - `third-party-error`: The fully-qualified name of the Java class providing the logic for the Third Party Error Log Publisher.\n  - `third-party-access`: The fully-qualified name of the Java class providing the logic for the Third Party Access Log Publisher.\n  - `third-party-file-based-error`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Error Log Publisher.",
				Optional:            true,
			},
			"debug_aci_enabled": schema.BoolAttribute{
				Description: "Indicates whether to include debugging information about ACIs being used by the operations being logged.",
				Optional:    true,
				Computed:    true,
			},
			"default_debug_level": schema.StringAttribute{
				Description: "The lowest severity level of debug messages to log when none of the defined targets match the message.",
				Optional:    true,
				Computed:    true,
			},
			"log_request_headers": schema.StringAttribute{
				Description: "Indicates whether request log messages should include information about HTTP headers included in the request.",
				Optional:    true,
				Computed:    true,
			},
			"suppressed_request_header_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request headers that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-headers property has a value of true.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_response_headers": schema.StringAttribute{
				Description: "Indicates whether response log messages should include information about HTTP headers included in the response.",
				Optional:    true,
				Computed:    true,
			},
			"suppressed_response_header_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of response headers that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-response-headers property has a value of true.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_request_authorization_type": schema.BoolAttribute{
				Description: "Indicates whether to log the type of credentials given if an \"Authorization\" header was included in the request. Logging the authorization type may be useful, and is much more secure than logging the entire value of the \"Authorization\" header.",
				Optional:    true,
				Computed:    true,
			},
			"log_request_cookie_names": schema.BoolAttribute{
				Description: "Indicates whether to log the names of any cookies included in an HTTP request. Logging cookie names may be useful and is much more secure than logging the entire content of the cookies (which may include sensitive information).",
				Optional:    true,
				Computed:    true,
			},
			"log_response_cookie_names": schema.BoolAttribute{
				Description: "Indicates whether to log the names of any cookies set in an HTTP response. Logging cookie names may be useful and is much more secure than logging the entire content of the cookies (which may include sensitive information).",
				Optional:    true,
				Computed:    true,
			},
			"log_request_parameters": schema.StringAttribute{
				Description: "Indicates what (if any) information about request parameters should be included in request log messages. Note that this will only be used for requests with a method other than GET, since GET request parameters will be included in the request URL.",
				Optional:    true,
				Computed:    true,
			},
			"log_request_protocol": schema.BoolAttribute{
				Description: "Indicates whether request log messages should include information about the HTTP version specified in the request.",
				Optional:    true,
				Computed:    true,
			},
			"suppressed_request_parameter_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request parameters that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-parameters property has a value of parameter-names or parameter-names-and-values.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"log_redirect_uri": schema.BoolAttribute{
				Description: "Indicates whether the redirect URI (i.e., the value of the \"Location\" header from responses) should be included in response log messages.",
				Optional:    true,
				Computed:    true,
			},
			"default_debug_category": schema.SetAttribute{
				Description: "The debug message categories to be logged when none of the defined targets match the message.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"default_omit_method_entry_arguments": schema.BoolAttribute{
				Description: "Indicates whether to include method arguments in debug messages logged by default.",
				Optional:    true,
				Computed:    true,
			},
			"default_omit_method_return_value": schema.BoolAttribute{
				Description: "Indicates whether to include the return value in debug messages logged by default.",
				Optional:    true,
				Computed:    true,
			},
			"default_include_throwable_cause": schema.BoolAttribute{
				Description: "Indicates whether to include the cause of exceptions in exception thrown and caught messages logged by default.",
				Optional:    true,
				Computed:    true,
			},
			"default_throwable_stack_frames": schema.Int64Attribute{
				Description: "Indicates the number of stack frames to include in the stack trace for method entry and exception thrown messages.",
				Optional:    true,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The set of arguments used to customize the behavior for the Third Party File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-http-operation`: The set of arguments used to customize the behavior for the Third Party HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-error`: The set of arguments used to customize the behavior for the Third Party Error Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-access`: The set of arguments used to customize the behavior for the Third Party Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-file-based-error`: The set of arguments used to customize the behavior for the Third Party File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The set of arguments used to customize the behavior for the Third Party File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-http-operation`: The set of arguments used to customize the behavior for the Third Party HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-error`: The set of arguments used to customize the behavior for the Third Party Error Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-access`: The set of arguments used to customize the behavior for the Third Party Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-file-based-error`: The set of arguments used to customize the behavior for the Third Party File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"syslog_external_server": schema.SetAttribute{
				Description: "The syslog server to which messages should be sent.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_details_in_result_messages": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation results should include information about both the request and the result. When the `type` attribute is set to  one of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether result log messages should include all of the elements of request log messages. This may be used to record a single message per operation with details about both the request and response.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation results should include information about both the request and the result.\n  - One of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether result log messages should include all of the elements of request log messages. This may be used to record a single message per operation with details about both the request and response.",
				Optional:            true,
				Computed:            true,
			},
			"log_assurance_completed": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of replication assurance processing.",
				Optional:    true,
				Computed:    true,
			},
			"debug_message_type": schema.SetAttribute{
				Description: "Specifies the debug message types which can be logged. Note that enabling these may result in sensitive information being logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"http_message_type": schema.SetAttribute{
				Description: "Specifies the HTTP message types which can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"access_token_validator_message_type": schema.SetAttribute{
				Description: "Specifies the access token validator message types that can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"id_token_validator_message_type": schema.SetAttribute{
				Description: "Specifies the ID token validator message types that can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"scim_message_type": schema.SetAttribute{
				Description: "Specifies the SCIM message types which can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"consent_message_type": schema.SetAttribute{
				Description: "Specifies the consent message types that can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"directory_rest_api_message_type": schema.SetAttribute{
				Description: "Specifies the Directory REST API message types which can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"extension_message_type": schema.SetAttribute{
				Description: "Specifies the Server SDK extension message types that can be logged.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"include_path_pattern": schema.SetAttribute{
				Description: "Specifies a set of HTTP request URL paths to determine whether log messages are included for a HTTP request. Log messages are included for a HTTP request if the request path does not match any exclude-path-pattern, and the request path does match an include-path-pattern (or no include-path-pattern is specified).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"exclude_path_pattern": schema.SetAttribute{
				Description: "Specifies a set of HTTP request URL paths to determine whether log messages are excluded for a HTTP request. Log messages are included for a HTTP request if the request path does not match any exclude-path-pattern, and the request path does match an include-path-pattern (or no include-path-pattern is specified).",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"server_host_name": schema.StringAttribute{
				Description: "Specifies the hostname or IP address of the syslogd host to log to. It is highly recommend to use localhost.",
				Optional:    true,
				Computed:    true,
			},
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
				Optional:    true,
				Computed:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the syslogd host to log to.",
				Optional:    true,
				Computed:    true,
			},
			"min_included_operation_processing_time": schema.StringAttribute{
				Description: "The minimum processing time (i.e., \"etime\") for operations that should be logged by this Operation Timing Access Log Publisher",
				Optional:    true,
			},
			"min_included_phase_time_nanos": schema.Int64Attribute{
				Description: "The minimum length of time in nanoseconds that an operation phase should take before it is included in a log message.",
				Optional:    true,
			},
			"time_interval": schema.StringAttribute{
				Description: "Specifies the interval at which to check whether the log files need to be rotated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_details_in_search_entry_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result entries should include information about the associated search request.",
				Optional:    true,
				Computed:    true,
			},
			"include_request_details_in_search_reference_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result references should include information about the associated search request.",
				Optional:    true,
				Computed:    true,
			},
			"include_request_details_in_intermediate_response_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for intermediate responses should include information about the associated operation request.",
				Optional:    true,
				Computed:    true,
			},
			"include_result_code_names": schema.BoolAttribute{
				Description: "Indicates whether result log messages should include human-readable names for result codes in addition to their numeric values.",
				Optional:    true,
				Computed:    true,
			},
			"include_extended_search_request_details": schema.BoolAttribute{
				Description: "Indicates whether log messages for search requests should include extended information from the request, including the requested size limit, time limit, alias dereferencing behavior, and types only behavior.",
				Optional:    true,
				Computed:    true,
			},
			"include_add_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for add requests should include a list of the names of the attributes included in the entry to add.",
				Optional:    true,
				Computed:    true,
			},
			"include_modify_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for modify requests should include a list of the names of the attributes to be modified.",
				Optional:    true,
				Computed:    true,
			},
			"include_search_entry_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result entries should include a list of the names of the attributes included in the entry that was returned.",
				Optional:    true,
				Computed:    true,
			},
			"log_connects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections established to the server.",
				Optional:    true,
				Computed:    true,
			},
			"log_disconnects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections that have been closed by the client or terminated by the server.",
				Optional:    true,
				Computed:    true,
			},
			"max_string_length": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`operation-timing-access`, `admin-alert-access`, `file-based-trace`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Specifies the maximum number of characters that may be included in any string in a log message before that string is truncated and replaced with a placeholder indicating the number of characters that were omitted. This can help prevent extremely long log messages from being written. When the `type` attribute is set to `detailed-http-operation`: Specifies the maximum length of any individual string that should be logged. If a log message includes a string longer than this number of characters, it will be truncated. A value of zero indicates that no truncation will be used.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`operation-timing-access`, `admin-alert-access`, `file-based-trace`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Specifies the maximum number of characters that may be included in any string in a log message before that string is truncated and replaced with a placeholder indicating the number of characters that were omitted. This can help prevent extremely long log messages from being written.\n  - `detailed-http-operation`: Specifies the maximum length of any individual string that should be logged. If a log message includes a string longer than this number of characters, it will be truncated. A value of zero indicates that no truncation will be used.",
				Optional:            true,
				Computed:            true,
			},
			"generify_message_strings_when_possible": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values. When the `type` attribute is set to  one of [`console-json-error`, `syslog-text-error`, `file-based-error`, `json-error`, `syslog-json-error`]: Indicates whether to use the generified version of the log message string (which may use placeholders like %s for a string or %d for an integer), rather than the version of the message with those placeholders replaced with specific values that would normally be written to the log.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values.\n  - One of [`console-json-error`, `syslog-text-error`, `file-based-error`, `json-error`, `syslog-json-error`]: Indicates whether to use the generified version of the log message string (which may use placeholders like %s for a string or %d for an integer), rather than the version of the message with those placeholders replaced with specific values that would normally be written to the log.",
				Optional:            true,
				Computed:            true,
			},
			"syslog_facility": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The syslog facility to use for the messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-based-error`: Specifies the syslog facility to use for this Syslog Based Error Log Publisher When the `type` attribute is set to `syslog-text-error`: The syslog facility to use for the messages that are logged by this Syslog Text Error Log Publisher. When the `type` attribute is set to `syslog-based-access`: Specifies the syslog facility to use for this Syslog Based Access Log Publisher When the `type` attribute is set to `syslog-text-access`: The syslog facility to use for the messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The syslog facility to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The syslog facility to use for the messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The syslog facility to use for the messages that are logged by this Syslog JSON Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The syslog facility to use for the messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-based-error`: Specifies the syslog facility to use for this Syslog Based Error Log Publisher\n  - `syslog-text-error`: The syslog facility to use for the messages that are logged by this Syslog Text Error Log Publisher.\n  - `syslog-based-access`: Specifies the syslog facility to use for this Syslog Based Access Log Publisher\n  - `syslog-text-access`: The syslog facility to use for the messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The syslog facility to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The syslog facility to use for the messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The syslog facility to use for the messages that are logged by this Syslog JSON Error Log Publisher.",
				Optional:            true,
				Computed:            true,
			},
			"log_field_behavior": schema.StringAttribute{
				Description: "The behavior to use for determining which fields to log and whether to transform the values of those fields in any way.",
				Optional:    true,
			},
			"log_client_certificates": schema.BoolAttribute{
				Description: "Indicates whether to log information about any client certificates presented to the server.",
				Optional:    true,
				Computed:    true,
			},
			"log_requests": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about requests received from clients. When the `type` attribute is set to  one of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about requests received from the client.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about requests received from clients.\n  - One of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about requests received from the client.",
				Optional:            true,
				Computed:            true,
			},
			"log_results": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about the results of client requests. When the `type` attribute is set to  one of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about the result of processing a requested HTTP operation.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about the results of client requests.\n  - One of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about the result of processing a requested HTTP operation.",
				Optional:            true,
				Computed:            true,
			},
			"log_search_entries": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result entries sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"log_search_references": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result references sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"log_intermediate_responses": schema.BoolAttribute{
				Description: "Indicates whether to log information about intermediate responses sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-based-access`, `syslog-text-access`, `file-based-access`]: Indicates whether the Writer Based Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `third-party-file-based-access`: Indicates whether the Third Party File Based Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `operation-timing-access`: Indicates whether the Operation Timing Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `admin-alert-access`: Indicates whether the Admin Alert Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-trace`: Indicates whether the Writer Based Trace Log Publisher will publish records asynchronously. When the `type` attribute is set to `common-log-file-http-operation`: Indicates whether the Common Log File HTTP Operation Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-json-audit`: Indicates whether the File Based JSON Audit Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-debug`: Indicates whether the File Based Debug Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-error`: Indicates whether the File Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `detailed-http-operation`: Indicates whether the Detailed HTTP Operation Log Publisher will publish records asynchronously. When the `type` attribute is set to `json-access`: Indicates whether the JSON Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `debug-access`: Indicates whether the Debug Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-audit`: Indicates whether the File Based Audit Log Publisher will publish records asynchronously. When the `type` attribute is set to `json-error`: Indicates whether the JSON Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `groovy-scripted-file-based-access`: Indicates whether the Scripted File Based Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `groovy-scripted-file-based-error`: Indicates whether the Scripted File Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `third-party-file-based-error`: Indicates whether the Third Party File Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-json-http-operation`: Indicates whether the File Based JSON HTTP Operation Log Publisher will publish records asynchronously.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-based-access`, `syslog-text-access`, `file-based-access`]: Indicates whether the Writer Based Access Log Publisher will publish records asynchronously.\n  - `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher will publish records asynchronously.\n  - `third-party-file-based-access`: Indicates whether the Third Party File Based Access Log Publisher will publish records asynchronously.\n  - `operation-timing-access`: Indicates whether the Operation Timing Access Log Publisher will publish records asynchronously.\n  - `admin-alert-access`: Indicates whether the Admin Alert Access Log Publisher will publish records asynchronously.\n  - `file-based-trace`: Indicates whether the Writer Based Trace Log Publisher will publish records asynchronously.\n  - `common-log-file-http-operation`: Indicates whether the Common Log File HTTP Operation Log Publisher will publish records asynchronously.\n  - `file-based-json-audit`: Indicates whether the File Based JSON Audit Log Publisher will publish records asynchronously.\n  - `file-based-debug`: Indicates whether the File Based Debug Log Publisher will publish records asynchronously.\n  - `file-based-error`: Indicates whether the File Based Error Log Publisher will publish records asynchronously.\n  - `detailed-http-operation`: Indicates whether the Detailed HTTP Operation Log Publisher will publish records asynchronously.\n  - `json-access`: Indicates whether the JSON Access Log Publisher will publish records asynchronously.\n  - `debug-access`: Indicates whether the Debug Access Log Publisher will publish records asynchronously.\n  - `file-based-audit`: Indicates whether the File Based Audit Log Publisher will publish records asynchronously.\n  - `json-error`: Indicates whether the JSON Error Log Publisher will publish records asynchronously.\n  - `groovy-scripted-file-based-access`: Indicates whether the Scripted File Based Access Log Publisher will publish records asynchronously.\n  - `groovy-scripted-file-based-error`: Indicates whether the Scripted File Based Error Log Publisher will publish records asynchronously.\n  - `third-party-file-based-error`: Indicates whether the Third Party File Based Error Log Publisher will publish records asynchronously.\n  - `file-based-json-http-operation`: Indicates whether the File Based JSON HTTP Operation Log Publisher will publish records asynchronously.",
				Optional:            true,
				Computed:            true,
			},
			"correlate_requests_and_results": schema.BoolAttribute{
				Description: "Indicates whether to automatically log result messages for any operation in which the corresponding request was logged. In such cases, the result, entry, and reference criteria will be ignored, although the log-responses, log-search-entries, and log-search-references properties will be honored.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_severity": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The syslog severity to use for the messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-text-error`: The syslog severity to use for the messages that are logged by this Syslog Text Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message. When the `type` attribute is set to `syslog-text-access`: The syslog severity to use for the messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The syslog severity to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The syslog severity to use for the messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The syslog severity to use for the messages that are logged by this Syslog JSON Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The syslog severity to use for the messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-text-error`: The syslog severity to use for the messages that are logged by this Syslog Text Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.\n  - `syslog-text-access`: The syslog severity to use for the messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The syslog severity to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The syslog severity to use for the messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The syslog severity to use for the messages that are logged by this Syslog JSON Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.",
				Optional:            true,
				Computed:            true,
			},
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"search_entry_criteria": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher. When the `type` attribute is set to `admin-alert-access`: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Admin Alert Access Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher.\n  - `admin-alert-access`: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Admin Alert Access Log Publisher.",
				Optional:            true,
			},
			"search_reference_criteria": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher. When the `type` attribute is set to `admin-alert-access`: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Admin Alert Access Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher.\n  - `admin-alert-access`: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Admin Alert Access Log Publisher.",
				Optional:            true,
			},
			"syslog_message_host_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-text-error`: The local host name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher. When the `type` attribute is set to `syslog-text-access`: The local host name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The local host name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-text-error`: The local host name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher.\n  - `syslog-text-access`: The local host name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The local host name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				Optional:            true,
			},
			"syslog_message_application_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The application name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-text-error`: The application name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher. When the `type` attribute is set to `syslog-text-access`: The application name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The application name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The application name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The application name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The application name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-text-error`: The application name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher.\n  - `syslog-text-access`: The application name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The application name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The application name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The application name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				Optional:            true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
			},
			"write_multi_line_messages": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `file-based-json-audit`, `syslog-json-http-operation`, `console-json-audit`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools. When the `type` attribute is set to  one of [`console-json-error`, `json-access`, `json-error`, `console-json-access`]: Indicates whether the JSON objects should be formatted to span multiple lines with a single element on each line. The multi-line format is potentially more user friendly (if administrators may need to look at the log files), but each message will be larger because of the additional spaces and end-of-line markers.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `file-based-json-audit`, `syslog-json-http-operation`, `console-json-audit`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools.\n  - One of [`console-json-error`, `json-access`, `json-error`, `console-json-access`]: Indicates whether the JSON objects should be formatted to span multiple lines with a single element on each line. The multi-line format is potentially more user friendly (if administrators may need to look at the log files), but each message will be larger because of the additional spaces and end-of-line markers.",
				Optional:            true,
				Computed:            true,
			},
			"use_reversible_form": schema.BoolAttribute{
				Description: "Indicates whether the audit log should be written in reversible form so that it is possible to revert the changes if desired.",
				Optional:    true,
				Computed:    true,
			},
			"soft_delete_entry_audit_behavior": schema.StringAttribute{
				Description: "Specifies the audit behavior for delete and modify operations on soft-deleted entries.",
				Optional:    true,
				Computed:    true,
			},
			"include_operation_purpose_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any operation purpose request control that may have been included in the request.",
				Optional:    true,
				Computed:    true,
			},
			"include_intermediate_client_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any intermediate client request control that may have been included in the request.",
				Optional:    true,
				Computed:    true,
			},
			"obscure_attribute": schema.SetAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `file-based-json-audit`, `file-based-audit`, `console-json-audit`]: Specifies the names of any attribute types that should have their values obscured in the audit log because they may be considered sensitive. When the `type` attribute is set to `debug-access`: Specifies the names of any attribute types that should have their values obscured if the obscure-sensitive-content property has a value of true.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `file-based-json-audit`, `file-based-audit`, `console-json-audit`]: Specifies the names of any attribute types that should have their values obscured in the audit log because they may be considered sensitive.\n  - `debug-access`: Specifies the names of any attribute types that should have their values obscured if the obscure-sensitive-content property has a value of true.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"exclude_attribute": schema.SetAttribute{
				Description: "Specifies the names of any attribute types that should be excluded from the audit log.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
				Optional:    true,
				Computed:    true,
			},
			"include_product_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the product name for the Directory Server.",
				Optional:    true,
				Computed:    true,
			},
			"include_instance_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the instance name for the Directory Server.",
				Optional:    true,
				Computed:    true,
			},
			"include_startup_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the startup ID for the Directory Server, which is a value assigned to the server instance at startup and may be used to identify when the server has been restarted.",
				Optional:    true,
				Computed:    true,
			},
			"include_thread_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the thread ID for the Directory Server in each log message. This ID can be used to correlate log messages from the same thread within a single log as well as generated by the same thread across different types of log files. More information about the thread with a specific ID can be obtained using the cn=JVM Stack Trace,cn=monitor entry.",
				Optional:    true,
				Computed:    true,
			},
			"include_requester_dn": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `admin-alert-access`, `syslog-based-access`, `file-based-json-audit`, `syslog-text-access`, `json-access`, `file-based-audit`, `syslog-json-access`, `console-json-audit`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested. When the `type` attribute is set to `operation-timing-access`: Indicates whether log messages should include the DN of the authenticated user for the client connection on which the operation was requested.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `admin-alert-access`, `syslog-based-access`, `file-based-json-audit`, `syslog-text-access`, `json-access`, `file-based-audit`, `syslog-json-access`, `console-json-audit`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested.\n  - `operation-timing-access`: Indicates whether log messages should include the DN of the authenticated user for the client connection on which the operation was requested.",
				Optional:            true,
				Computed:            true,
			},
			"include_requester_ip_address": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the IP address of the client that requested the operation.",
				Optional:    true,
				Computed:    true,
			},
			"include_request_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include a list of the OIDs of any controls included in the request.",
				Optional:    true,
				Computed:    true,
			},
			"include_response_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation results should include a list of the OIDs of any controls included in the result.",
				Optional:    true,
				Computed:    true,
			},
			"include_replication_change_id": schema.BoolAttribute{
				Description: "Indicates whether to log information about the replication change ID.",
				Optional:    true,
				Computed:    true,
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Optional:    true,
				Computed:    true,
			},
			"suppress_replication_operations": schema.BoolAttribute{
				Description: "Indicates whether access messages that are generated by replication operations should be suppressed.",
				Optional:    true,
				Computed:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a set of connection criteria that must match the associated client connection in order for a connect, disconnect, request, or result message to be logged.",
				Optional:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a set of request criteria that must match the associated operation request in order for a request or result to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"result_criteria": schema.StringAttribute{
				Description: "Specifies a set of result criteria that must match the associated operation result in order for that result to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `third-party-file-based-access`, `operation-timing-access`, `third-party-http-operation`, `admin-alert-access`, `file-based-trace`, `jdbc-based-error`, `jdbc-based-access`, `common-log-file-http-operation`, `syslog-text-error`, `file-based-json-audit`, `file-based-debug`, `file-based-error`, `third-party-error`, `syslog-text-access`, `detailed-http-operation`, `json-access`, `debug-access`, `syslog-json-http-operation`, `third-party-access`, `file-based-audit`, `json-error`, `groovy-scripted-file-based-access`, `groovy-scripted-file-based-error`, `syslog-json-access`, `groovy-scripted-access`, `third-party-file-based-error`, `file-based-access`, `groovy-scripted-error`, `file-based-json-http-operation`, `syslog-json-error`, `groovy-scripted-http-operation`]: Indicates whether the Log Publisher is enabled for use. When the `type` attribute is set to `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher is enabled for use. When the `type` attribute is set to `console-json-error`: Indicates whether the Console JSON Error Log Publisher is enabled for use. When the `type` attribute is set to `syslog-based-access`: Indicates whether the Syslog Based Access Log Publisher is enabled for use. When the `type` attribute is set to `console-json-audit`: Indicates whether the Console JSON Audit Log Publisher is enabled for use. When the `type` attribute is set to `console-json-http-operation`: Indicates whether the Console JSON HTTP Operation Log Publisher is enabled for use. When the `type` attribute is set to `console-json-access`: Indicates whether the Console JSON Access Log Publisher is enabled for use.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `third-party-file-based-access`, `operation-timing-access`, `third-party-http-operation`, `admin-alert-access`, `file-based-trace`, `jdbc-based-error`, `jdbc-based-access`, `common-log-file-http-operation`, `syslog-text-error`, `file-based-json-audit`, `file-based-debug`, `file-based-error`, `third-party-error`, `syslog-text-access`, `detailed-http-operation`, `json-access`, `debug-access`, `syslog-json-http-operation`, `third-party-access`, `file-based-audit`, `json-error`, `groovy-scripted-file-based-access`, `groovy-scripted-file-based-error`, `syslog-json-access`, `groovy-scripted-access`, `third-party-file-based-error`, `file-based-access`, `groovy-scripted-error`, `file-based-json-http-operation`, `syslog-json-error`, `groovy-scripted-http-operation`]: Indicates whether the Log Publisher is enabled for use.\n  - `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher is enabled for use.\n  - `console-json-error`: Indicates whether the Console JSON Error Log Publisher is enabled for use.\n  - `syslog-based-access`: Indicates whether the Syslog Based Access Log Publisher is enabled for use.\n  - `console-json-audit`: Indicates whether the Console JSON Audit Log Publisher is enabled for use.\n  - `console-json-http-operation`: Indicates whether the Console JSON HTTP Operation Log Publisher is enabled for use.\n  - `console-json-access`: Indicates whether the Console JSON Access Log Publisher is enabled for use.",
				Required:            true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("standard-error"),
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
	} else {
		// Add RequiresReplace modifier for read-only attributes
		compressionMechanismAttr := schemaDef.Attributes["compression_mechanism"].(schema.StringAttribute)
		compressionMechanismAttr.PlanModifiers = append(compressionMechanismAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["compression_mechanism"] = compressionMechanismAttr
		encryptLogAttr := schemaDef.Attributes["encrypt_log"].(schema.BoolAttribute)
		encryptLogAttr.PlanModifiers = append(encryptLogAttr.PlanModifiers, boolplanmodifier.RequiresReplace())
		schemaDef.Attributes["encrypt_log"] = encryptLogAttr
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *logPublisherResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var planModel, configModel logPublisherResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for syslog-json-audit type
	if resourceType == "syslog-json-audit" {
		if !internaltypes.IsDefined(configModel.UseReversibleForm) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseReversibleForm.Equal(defaultVal) {
				planModel.UseReversibleForm = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeOperationPurposeRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeOperationPurposeRequestControl.Equal(defaultVal) {
				planModel.IncludeOperationPurposeRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeIntermediateClientRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeIntermediateClientRequestControl.Equal(defaultVal) {
				planModel.IncludeIntermediateClientRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ExcludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("ds-sync-hist")})
			if !planModel.ExcludeAttribute.Equal(defaultVal) {
				planModel.ExcludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-based-error type
	if resourceType == "syslog-based-error" {
		if !internaltypes.IsDefined(configModel.ServerHostName) {
			defaultVal := types.StringValue("localhost")
			if !planModel.ServerHostName.Equal(defaultVal) {
				planModel.ServerHostName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(514)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("1")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for third-party-file-based-access type
	if resourceType == "third-party-file-based-access" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for operation-timing-access type
	if resourceType == "operation-timing-access" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for admin-alert-access type
	if resourceType == "admin-alert-access" {
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(false)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(1000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchEntryMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchEntryMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchEntryMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchReferenceMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchReferenceMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchReferenceMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInIntermediateResponseMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInIntermediateResponseMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInIntermediateResponseMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResultCodeNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResultCodeNames.Equal(defaultVal) {
				planModel.IncludeResultCodeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtendedSearchRequestDetails) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtendedSearchRequestDetails.Equal(defaultVal) {
				planModel.IncludeExtendedSearchRequestDetails = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAddAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAddAttributeNames.Equal(defaultVal) {
				planModel.IncludeAddAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeModifyAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeModifyAttributeNames.Equal(defaultVal) {
				planModel.IncludeModifyAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeSearchEntryAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeSearchEntryAttributeNames.Equal(defaultVal) {
				planModel.IncludeSearchEntryAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(2000)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-trace type
	if resourceType == "file-based-trace" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(50000)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for jdbc-based-error type
	if resourceType == "jdbc-based-error" {
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for jdbc-based-access type
	if resourceType == "jdbc-based-access" {
		if !internaltypes.IsDefined(configModel.LogTableName) {
			defaultVal := types.StringValue("access_log")
			if !planModel.LogTableName.Equal(defaultVal) {
				planModel.LogTableName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for common-log-file-http-operation type
	if resourceType == "common-log-file-http-operation" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-text-error type
	if resourceType == "syslog-text-error" {
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-error"), types.StringValue("severe-warning"), types.StringValue("notice")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("system-daemons")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TimestampPrecision) {
			defaultVal := types.StringValue("milliseconds")
			if !planModel.TimestampPrecision.Equal(defaultVal) {
				planModel.TimestampPrecision = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(100000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-based-access type
	if resourceType == "syslog-based-access" {
		if !internaltypes.IsDefined(configModel.ServerHostName) {
			defaultVal := types.StringValue("localhost")
			if !planModel.ServerHostName.Equal(defaultVal) {
				planModel.ServerHostName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ServerPort) {
			defaultVal := types.Int64Value(514)
			if !planModel.ServerPort.Equal(defaultVal) {
				planModel.ServerPort = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("1")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(500)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchEntryMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchEntryMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchEntryMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchReferenceMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchReferenceMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchReferenceMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInIntermediateResponseMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInIntermediateResponseMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInIntermediateResponseMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResultCodeNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResultCodeNames.Equal(defaultVal) {
				planModel.IncludeResultCodeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtendedSearchRequestDetails) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtendedSearchRequestDetails.Equal(defaultVal) {
				planModel.IncludeExtendedSearchRequestDetails = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAddAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAddAttributeNames.Equal(defaultVal) {
				planModel.IncludeAddAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeModifyAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeModifyAttributeNames.Equal(defaultVal) {
				planModel.IncludeModifyAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeSearchEntryAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeSearchEntryAttributeNames.Equal(defaultVal) {
				planModel.IncludeSearchEntryAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-json-audit type
	if resourceType == "file-based-json-audit" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseReversibleForm) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseReversibleForm.Equal(defaultVal) {
				planModel.UseReversibleForm = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SoftDeleteEntryAuditBehavior) {
			defaultVal := types.StringValue("included")
			if !planModel.SoftDeleteEntryAuditBehavior.Equal(defaultVal) {
				planModel.SoftDeleteEntryAuditBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeOperationPurposeRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeOperationPurposeRequestControl.Equal(defaultVal) {
				planModel.IncludeOperationPurposeRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeIntermediateClientRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeIntermediateClientRequestControl.Equal(defaultVal) {
				planModel.IncludeIntermediateClientRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ExcludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("ds-sync-hist")})
			if !planModel.ExcludeAttribute.Equal(defaultVal) {
				planModel.ExcludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-debug type
	if resourceType == "file-based-debug" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TimestampPrecision) {
			defaultVal := types.StringValue("milliseconds")
			if !planModel.TimestampPrecision.Equal(defaultVal) {
				planModel.TimestampPrecision = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultDebugLevel) {
			defaultVal := types.StringValue("disabled")
			if !planModel.DefaultDebugLevel.Equal(defaultVal) {
				planModel.DefaultDebugLevel = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultOmitMethodEntryArguments) {
			defaultVal := types.BoolValue(false)
			if !planModel.DefaultOmitMethodEntryArguments.Equal(defaultVal) {
				planModel.DefaultOmitMethodEntryArguments = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultOmitMethodReturnValue) {
			defaultVal := types.BoolValue(false)
			if !planModel.DefaultOmitMethodReturnValue.Equal(defaultVal) {
				planModel.DefaultOmitMethodReturnValue = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultIncludeThrowableCause) {
			defaultVal := types.BoolValue(true)
			if !planModel.DefaultIncludeThrowableCause.Equal(defaultVal) {
				planModel.DefaultIncludeThrowableCause = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultThrowableStackFrames) {
			defaultVal := types.Int64Value(2147483647)
			if !planModel.DefaultThrowableStackFrames.Equal(defaultVal) {
				planModel.DefaultThrowableStackFrames = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-error type
	if resourceType == "file-based-error" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(false)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TimestampPrecision) {
			defaultVal := types.StringValue("milliseconds")
			if !planModel.TimestampPrecision.Equal(defaultVal) {
				planModel.TimestampPrecision = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for third-party-error type
	if resourceType == "third-party-error" {
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-text-access type
	if resourceType == "syslog-text-access" {
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("system-daemons")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogSeverity) {
			defaultVal := types.StringValue("informational")
			if !planModel.SyslogSeverity.Equal(defaultVal) {
				planModel.SyslogSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(100000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchEntryMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchEntryMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchEntryMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchReferenceMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchReferenceMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchReferenceMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInIntermediateResponseMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInIntermediateResponseMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInIntermediateResponseMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResultCodeNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResultCodeNames.Equal(defaultVal) {
				planModel.IncludeResultCodeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtendedSearchRequestDetails) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtendedSearchRequestDetails.Equal(defaultVal) {
				planModel.IncludeExtendedSearchRequestDetails = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAddAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAddAttributeNames.Equal(defaultVal) {
				planModel.IncludeAddAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeModifyAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeModifyAttributeNames.Equal(defaultVal) {
				planModel.IncludeModifyAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeSearchEntryAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeSearchEntryAttributeNames.Equal(defaultVal) {
				planModel.IncludeSearchEntryAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(500)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TimestampPrecision) {
			defaultVal := types.StringValue("milliseconds")
			if !planModel.TimestampPrecision.Equal(defaultVal) {
				planModel.TimestampPrecision = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for detailed-http-operation type
	if resourceType == "detailed-http-operation" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogRequestHeaders.Equal(defaultVal) {
				planModel.LogRequestHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedRequestHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Authorization"), types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Cookie")})
			if !planModel.SuppressedRequestHeaderName.Equal(defaultVal) {
				planModel.SuppressedRequestHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogResponseHeaders.Equal(defaultVal) {
				planModel.LogResponseHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedResponseHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Location"), types.StringValue("Set-Cookie")})
			if !planModel.SuppressedResponseHeaderName.Equal(defaultVal) {
				planModel.SuppressedResponseHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestAuthorizationType) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestAuthorizationType.Equal(defaultVal) {
				planModel.LogRequestAuthorizationType = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestCookieNames.Equal(defaultVal) {
				planModel.LogRequestCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResponseCookieNames.Equal(defaultVal) {
				planModel.LogResponseCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestParameters) {
			defaultVal := types.StringValue("parameter-names")
			if !planModel.LogRequestParameters.Equal(defaultVal) {
				planModel.LogRequestParameters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRedirectURI) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRedirectURI.Equal(defaultVal) {
				planModel.LogRedirectURI = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(2000)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for json-access type
	if resourceType == "json-access" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchEntryMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchEntryMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchEntryMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchReferenceMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchReferenceMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchReferenceMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInIntermediateResponseMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInIntermediateResponseMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInIntermediateResponseMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResultCodeNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResultCodeNames.Equal(defaultVal) {
				planModel.IncludeResultCodeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtendedSearchRequestDetails) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtendedSearchRequestDetails.Equal(defaultVal) {
				planModel.IncludeExtendedSearchRequestDetails = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAddAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAddAttributeNames.Equal(defaultVal) {
				planModel.IncludeAddAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeModifyAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeModifyAttributeNames.Equal(defaultVal) {
				planModel.IncludeModifyAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeSearchEntryAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeSearchEntryAttributeNames.Equal(defaultVal) {
				planModel.IncludeSearchEntryAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(2000)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for debug-access type
	if resourceType == "debug-access" {
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ObscureSensitiveContent) {
			defaultVal := types.BoolValue(true)
			if !planModel.ObscureSensitiveContent.Equal(defaultVal) {
				planModel.ObscureSensitiveContent = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ObscureAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("userPassword"), types.StringValue("authPassword")})
			if !planModel.ObscureAttribute.Equal(defaultVal) {
				planModel.ObscureAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DebugACIEnabled) {
			defaultVal := types.BoolValue(false)
			if !planModel.DebugACIEnabled.Equal(defaultVal) {
				planModel.DebugACIEnabled = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-json-http-operation type
	if resourceType == "syslog-json-http-operation" {
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("system-daemons")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogSeverity) {
			defaultVal := types.StringValue("informational")
			if !planModel.SyslogSeverity.Equal(defaultVal) {
				planModel.SyslogSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(100000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogRequestHeaders.Equal(defaultVal) {
				planModel.LogRequestHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedRequestHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Authorization"), types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Cookie")})
			if !planModel.SuppressedRequestHeaderName.Equal(defaultVal) {
				planModel.SuppressedRequestHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogResponseHeaders.Equal(defaultVal) {
				planModel.LogResponseHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedResponseHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Location"), types.StringValue("Set-Cookie")})
			if !planModel.SuppressedResponseHeaderName.Equal(defaultVal) {
				planModel.SuppressedResponseHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestAuthorizationType) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestAuthorizationType.Equal(defaultVal) {
				planModel.LogRequestAuthorizationType = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestCookieNames.Equal(defaultVal) {
				planModel.LogRequestCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResponseCookieNames.Equal(defaultVal) {
				planModel.LogResponseCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestParameters) {
			defaultVal := types.StringValue("parameter-names")
			if !planModel.LogRequestParameters.Equal(defaultVal) {
				planModel.LogRequestParameters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestProtocol) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestProtocol.Equal(defaultVal) {
				planModel.LogRequestProtocol = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRedirectURI) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRedirectURI.Equal(defaultVal) {
				planModel.LogRedirectURI = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for third-party-access type
	if resourceType == "third-party-access" {
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-audit type
	if resourceType == "file-based-audit" {
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseReversibleForm) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseReversibleForm.Equal(defaultVal) {
				planModel.UseReversibleForm = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SoftDeleteEntryAuditBehavior) {
			defaultVal := types.StringValue("commented")
			if !planModel.SoftDeleteEntryAuditBehavior.Equal(defaultVal) {
				planModel.SoftDeleteEntryAuditBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeOperationPurposeRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeOperationPurposeRequestControl.Equal(defaultVal) {
				planModel.IncludeOperationPurposeRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeIntermediateClientRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeIntermediateClientRequestControl.Equal(defaultVal) {
				planModel.IncludeIntermediateClientRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ExcludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("ds-sync-hist")})
			if !planModel.ExcludeAttribute.Equal(defaultVal) {
				planModel.ExcludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TimestampPrecision) {
			defaultVal := types.StringValue("milliseconds")
			if !planModel.TimestampPrecision.Equal(defaultVal) {
				planModel.TimestampPrecision = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for json-error type
	if resourceType == "json-error" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(false)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for groovy-scripted-file-based-access type
	if resourceType == "groovy-scripted-file-based-access" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for groovy-scripted-file-based-error type
	if resourceType == "groovy-scripted-file-based-error" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-json-access type
	if resourceType == "syslog-json-access" {
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("system-daemons")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogSeverity) {
			defaultVal := types.StringValue("informational")
			if !planModel.SyslogSeverity.Equal(defaultVal) {
				planModel.SyslogSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(100000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchEntryMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchEntryMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchEntryMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchReferenceMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchReferenceMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchReferenceMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInIntermediateResponseMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInIntermediateResponseMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInIntermediateResponseMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResultCodeNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResultCodeNames.Equal(defaultVal) {
				planModel.IncludeResultCodeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtendedSearchRequestDetails) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtendedSearchRequestDetails.Equal(defaultVal) {
				planModel.IncludeExtendedSearchRequestDetails = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAddAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAddAttributeNames.Equal(defaultVal) {
				planModel.IncludeAddAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeModifyAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeModifyAttributeNames.Equal(defaultVal) {
				planModel.IncludeModifyAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeSearchEntryAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeSearchEntryAttributeNames.Equal(defaultVal) {
				planModel.IncludeSearchEntryAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(2000)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for groovy-scripted-access type
	if resourceType == "groovy-scripted-access" {
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for third-party-file-based-error type
	if resourceType == "third-party-file-based-error" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for console-json-audit type
	if resourceType == "console-json-audit" {
		if !internaltypes.IsDefined(configModel.OutputLocation) {
			defaultVal := types.StringValue("standard-output")
			if !planModel.OutputLocation.Equal(defaultVal) {
				planModel.OutputLocation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseReversibleForm) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseReversibleForm.Equal(defaultVal) {
				planModel.UseReversibleForm = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SoftDeleteEntryAuditBehavior) {
			defaultVal := types.StringValue("included")
			if !planModel.SoftDeleteEntryAuditBehavior.Equal(defaultVal) {
				planModel.SoftDeleteEntryAuditBehavior = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeOperationPurposeRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeOperationPurposeRequestControl.Equal(defaultVal) {
				planModel.IncludeOperationPurposeRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeIntermediateClientRequestControl) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeIntermediateClientRequestControl.Equal(defaultVal) {
				planModel.IncludeIntermediateClientRequestControl = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ExcludeAttribute) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("ds-sync-hist")})
			if !planModel.ExcludeAttribute.Equal(defaultVal) {
				planModel.ExcludeAttribute = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for console-json-http-operation type
	if resourceType == "console-json-http-operation" {
		if !internaltypes.IsDefined(configModel.OutputLocation) {
			defaultVal := types.StringValue("standard-output")
			if !planModel.OutputLocation.Equal(defaultVal) {
				planModel.OutputLocation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogRequestHeaders.Equal(defaultVal) {
				planModel.LogRequestHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedRequestHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Authorization"), types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Cookie")})
			if !planModel.SuppressedRequestHeaderName.Equal(defaultVal) {
				planModel.SuppressedRequestHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogResponseHeaders.Equal(defaultVal) {
				planModel.LogResponseHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedResponseHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Location"), types.StringValue("Set-Cookie")})
			if !planModel.SuppressedResponseHeaderName.Equal(defaultVal) {
				planModel.SuppressedResponseHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestAuthorizationType) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestAuthorizationType.Equal(defaultVal) {
				planModel.LogRequestAuthorizationType = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestCookieNames.Equal(defaultVal) {
				planModel.LogRequestCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResponseCookieNames.Equal(defaultVal) {
				planModel.LogResponseCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestParameters) {
			defaultVal := types.StringValue("parameter-names")
			if !planModel.LogRequestParameters.Equal(defaultVal) {
				planModel.LogRequestParameters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestProtocol) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestProtocol.Equal(defaultVal) {
				planModel.LogRequestProtocol = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRedirectURI) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRedirectURI.Equal(defaultVal) {
				planModel.LogRedirectURI = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-access type
	if resourceType == "file-based-access" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.TimestampPrecision) {
			defaultVal := types.StringValue("milliseconds")
			if !planModel.TimestampPrecision.Equal(defaultVal) {
				planModel.TimestampPrecision = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogConnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogConnects.Equal(defaultVal) {
				planModel.LogConnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogDisconnects) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogDisconnects.Equal(defaultVal) {
				planModel.LogDisconnects = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogAssuranceCompleted) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogAssuranceCompleted.Equal(defaultVal) {
				planModel.LogAssuranceCompleted = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterDN) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterDN.Equal(defaultVal) {
				planModel.IncludeRequesterDN = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequesterIPAddress) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequesterIPAddress.Equal(defaultVal) {
				planModel.IncludeRequesterIPAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchEntryMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchEntryMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchEntryMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInSearchReferenceMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInSearchReferenceMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInSearchReferenceMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInIntermediateResponseMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestDetailsInIntermediateResponseMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInIntermediateResponseMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResultCodeNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeResultCodeNames.Equal(defaultVal) {
				planModel.IncludeResultCodeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeExtendedSearchRequestDetails) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeExtendedSearchRequestDetails.Equal(defaultVal) {
				planModel.IncludeExtendedSearchRequestDetails = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeAddAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeAddAttributeNames.Equal(defaultVal) {
				planModel.IncludeAddAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeModifyAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeModifyAttributeNames.Equal(defaultVal) {
				planModel.IncludeModifyAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeSearchEntryAttributeNames) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeSearchEntryAttributeNames.Equal(defaultVal) {
				planModel.IncludeSearchEntryAttributeNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeRequestControls.Equal(defaultVal) {
				planModel.IncludeRequestControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeResponseControls) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeResponseControls.Equal(defaultVal) {
				planModel.IncludeResponseControls = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeReplicationChangeID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeReplicationChangeID.Equal(defaultVal) {
				planModel.IncludeReplicationChangeID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxStringLength) {
			defaultVal := types.Int64Value(2000)
			if !planModel.MaxStringLength.Equal(defaultVal) {
				planModel.MaxStringLength = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSecurityNegotiation) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSecurityNegotiation.Equal(defaultVal) {
				planModel.LogSecurityNegotiation = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogClientCertificates) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogClientCertificates.Equal(defaultVal) {
				planModel.LogClientCertificates = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchEntries) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchEntries.Equal(defaultVal) {
				planModel.LogSearchEntries = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogSearchReferences) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogSearchReferences.Equal(defaultVal) {
				planModel.LogSearchReferences = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogIntermediateResponses) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogIntermediateResponses.Equal(defaultVal) {
				planModel.LogIntermediateResponses = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressInternalOperations) {
			defaultVal := types.BoolValue(true)
			if !planModel.SuppressInternalOperations.Equal(defaultVal) {
				planModel.SuppressInternalOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressReplicationOperations) {
			defaultVal := types.BoolValue(false)
			if !planModel.SuppressReplicationOperations.Equal(defaultVal) {
				planModel.SuppressReplicationOperations = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelateRequestsAndResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.CorrelateRequestsAndResults.Equal(defaultVal) {
				planModel.CorrelateRequestsAndResults = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for groovy-scripted-error type
	if resourceType == "groovy-scripted-error" {
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-warning"), types.StringValue("severe-error")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for file-based-json-http-operation type
	if resourceType == "file-based-json-http-operation" {
		if !internaltypes.IsDefined(configModel.LogFilePermissions) {
			defaultVal := types.StringValue("600")
			if !planModel.LogFilePermissions.Equal(defaultVal) {
				planModel.LogFilePermissions = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CompressionMechanism) {
			defaultVal := types.StringValue("none")
			if !planModel.CompressionMechanism.Equal(defaultVal) {
				planModel.CompressionMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SignLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.SignLog.Equal(defaultVal) {
				planModel.SignLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EncryptLog) {
			defaultVal := types.BoolValue(false)
			if !planModel.EncryptLog.Equal(defaultVal) {
				planModel.EncryptLog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Append) {
			defaultVal := types.BoolValue(true)
			if !planModel.Append.Equal(defaultVal) {
				planModel.Append = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.Asynchronous) {
			defaultVal := types.BoolValue(true)
			if !planModel.Asynchronous.Equal(defaultVal) {
				planModel.Asynchronous = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoFlush) {
			defaultVal := types.BoolValue(true)
			if !planModel.AutoFlush.Equal(defaultVal) {
				planModel.AutoFlush = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.BufferSize) {
			defaultVal := types.StringValue("64 kb")
			if !planModel.BufferSize.Equal(defaultVal) {
				planModel.BufferSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(10000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequests) {
			defaultVal := types.BoolValue(false)
			if !planModel.LogRequests.Equal(defaultVal) {
				planModel.LogRequests = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResults) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResults.Equal(defaultVal) {
				planModel.LogResults = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(false)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeRequestDetailsInResultMessages) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeRequestDetailsInResultMessages.Equal(defaultVal) {
				planModel.IncludeRequestDetailsInResultMessages = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogRequestHeaders.Equal(defaultVal) {
				planModel.LogRequestHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedRequestHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Authorization"), types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Cookie")})
			if !planModel.SuppressedRequestHeaderName.Equal(defaultVal) {
				planModel.SuppressedRequestHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseHeaders) {
			defaultVal := types.StringValue("none")
			if !planModel.LogResponseHeaders.Equal(defaultVal) {
				planModel.LogResponseHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SuppressedResponseHeaderName) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("Content-Length"), types.StringValue("Content-Type"), types.StringValue("Location"), types.StringValue("Set-Cookie")})
			if !planModel.SuppressedResponseHeaderName.Equal(defaultVal) {
				planModel.SuppressedResponseHeaderName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestAuthorizationType) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestAuthorizationType.Equal(defaultVal) {
				planModel.LogRequestAuthorizationType = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestCookieNames.Equal(defaultVal) {
				planModel.LogRequestCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogResponseCookieNames) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogResponseCookieNames.Equal(defaultVal) {
				planModel.LogResponseCookieNames = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestParameters) {
			defaultVal := types.StringValue("parameter-names")
			if !planModel.LogRequestParameters.Equal(defaultVal) {
				planModel.LogRequestParameters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRequestProtocol) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRequestProtocol.Equal(defaultVal) {
				planModel.LogRequestProtocol = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LogRedirectURI) {
			defaultVal := types.BoolValue(true)
			if !planModel.LogRedirectURI.Equal(defaultVal) {
				planModel.LogRedirectURI = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.WriteMultiLineMessages) {
			defaultVal := types.BoolValue(false)
			if !planModel.WriteMultiLineMessages.Equal(defaultVal) {
				planModel.WriteMultiLineMessages = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for syslog-json-error type
	if resourceType == "syslog-json-error" {
		if !internaltypes.IsDefined(configModel.DefaultSeverity) {
			defaultVal, _ := types.SetValue(types.StringType, []attr.Value{types.StringValue("fatal-error"), types.StringValue("severe-error"), types.StringValue("severe-warning"), types.StringValue("notice")})
			if !planModel.DefaultSeverity.Equal(defaultVal) {
				planModel.DefaultSeverity = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SyslogFacility) {
			defaultVal := types.StringValue("system-daemons")
			if !planModel.SyslogFacility.Equal(defaultVal) {
				planModel.SyslogFacility = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.QueueSize) {
			defaultVal := types.Int64Value(100000)
			if !planModel.QueueSize.Equal(defaultVal) {
				planModel.QueueSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeProductName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeProductName.Equal(defaultVal) {
				planModel.IncludeProductName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeInstanceName) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeInstanceName.Equal(defaultVal) {
				planModel.IncludeInstanceName = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeStartupID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeStartupID.Equal(defaultVal) {
				planModel.IncludeStartupID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IncludeThreadID) {
			defaultVal := types.BoolValue(true)
			if !planModel.IncludeThreadID.Equal(defaultVal) {
				planModel.IncludeThreadID = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.GenerifyMessageStringsWhenPossible) {
			defaultVal := types.BoolValue(false)
			if !planModel.GenerifyMessageStringsWhenPossible.Equal(defaultVal) {
				planModel.GenerifyMessageStringsWhenPossible = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (model *logPublisherResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "syslog-json-audit" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.LogRedirectURI = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.LogClientCertificates = types.BoolNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "syslog-based-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "third-party-file-based-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "operation-timing-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "third-party-http-operation" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "admin-alert-access" {
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "file-based-trace" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "jdbc-based-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "jdbc-based-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "common-log-file-http-operation" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "syslog-text-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "syslog-based-access" {
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "file-based-json-audit" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.LogRedirectURI = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.LogClientCertificates = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "file-based-debug" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "file-based-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "third-party-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "syslog-text-access" {
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.OutputLocation = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "detailed-http-operation" {
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "json-access" {
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "debug-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "syslog-json-http-operation" {
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.CompressionMechanism = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "third-party-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "file-based-audit" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.LogRedirectURI = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.LogClientCertificates = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "json-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.LogIntermediateResponses = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "groovy-scripted-file-based-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "groovy-scripted-file-based-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "syslog-json-access" {
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.Asynchronous = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "groovy-scripted-access" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "third-party-file-based-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "console-json-audit" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.LogRedirectURI = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.LogClientCertificates = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "console-json-http-operation" {
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "file-based-access" {
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.OutputLocation = types.StringNull()
		model.SyslogSeverity = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if resourceType == "groovy-scripted-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "file-based-json-http-operation" {
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.MaxStringLength = types.Int64Null()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogSecurityNegotiation = types.BoolNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
	}
	if resourceType == "syslog-json-error" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
	if resourceType == "groovy-scripted-http-operation" {
		model.IncludeRequestDetailsInResultMessages = types.BoolNull()
		model.IncludeRequestDetailsInSearchEntryMessages = types.BoolNull()
		model.LogRequests = types.BoolNull()
		model.IncludeRequesterDN = types.BoolNull()
		model.ServerHostName = types.StringNull()
		model.LogResponseHeaders = types.StringNull()
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResponseCookieNames = types.BoolNull()
		model.SuppressInternalOperations = types.BoolNull()
		model.IncludeResponseControls = types.BoolNull()
		model.LogRequestHeaders = types.StringNull()
		model.MaxStringLength = types.Int64Null()
		model.EncryptLog = types.BoolNull()
		model.WriteMultiLineMessages = types.BoolNull()
		model.LogIntermediateResponses = types.BoolNull()
		model.GenerifyMessageStringsWhenPossible = types.BoolNull()
		model.IncludeResultCodeNames = types.BoolNull()
		model.DefaultOmitMethodReturnValue = types.BoolNull()
		model.CorrelateRequestsAndResults = types.BoolNull()
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.SignLog = types.BoolNull()
		model.BufferSize = types.StringNull()
		model.LogRequestCookieNames = types.BoolNull()
		model.LogRequestProtocol = types.BoolNull()
		model.LogSecurityNegotiation = types.BoolNull()
		model.TimeInterval = types.StringNull()
		model.IncludeAddAttributeNames = types.BoolNull()
		model.IncludeExtendedSearchRequestDetails = types.BoolNull()
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogRedirectURI = types.BoolNull()
		model.UseReversibleForm = types.BoolNull()
		model.DebugACIEnabled = types.BoolNull()
		model.IncludeRequestControls = types.BoolNull()
		model.LogRequestAuthorizationType = types.BoolNull()
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
		model.DefaultDebugLevel = types.StringNull()
		model.SuppressReplicationOperations = types.BoolNull()
		model.LogClientCertificates = types.BoolNull()
		model.IncludeIntermediateClientRequestControl = types.BoolNull()
		model.SyslogFacility = types.StringNull()
		model.LogAssuranceCompleted = types.BoolNull()
		model.IncludeRequestDetailsInSearchReferenceMessages = types.BoolNull()
		model.IncludeOperationPurposeRequestControl = types.BoolNull()
		model.IncludeStartupID = types.BoolNull()
		model.DefaultOmitMethodEntryArguments = types.BoolNull()
		model.IncludeRequesterIPAddress = types.BoolNull()
		model.IncludeSearchEntryAttributeNames = types.BoolNull()
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogResults = types.BoolNull()
		model.DefaultThrowableStackFrames = types.Int64Null()
		model.TimestampPrecision = types.StringNull()
		model.LogSearchReferences = types.BoolNull()
		model.OutputLocation = types.StringNull()
		model.LogConnects = types.BoolNull()
		model.SyslogSeverity = types.StringNull()
		model.CompressionMechanism = types.StringNull()
		model.LogRequestParameters = types.StringNull()
		model.ServerPort = types.Int64Null()
		model.LogTableName = types.StringNull()
		model.LogSearchEntries = types.BoolNull()
		model.ObscureSensitiveContent = types.BoolNull()
		model.SoftDeleteEntryAuditBehavior = types.StringNull()
		model.IncludeThreadID = types.BoolNull()
		model.Asynchronous = types.BoolNull()
		model.IncludeModifyAttributeNames = types.BoolNull()
		model.DefaultIncludeThrowableCause = types.BoolNull()
		model.IncludeReplicationChangeID = types.BoolNull()
		model.IncludeProductName = types.BoolNull()
		model.IncludeInstanceName = types.BoolNull()
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
		model.LogDisconnects = types.BoolNull()
		model.AutoFlush = types.BoolNull()
		model.QueueSize = types.Int64Null()
		model.IncludeRequestDetailsInIntermediateResponseMessages = types.BoolNull()
		model.LogFilePermissions = types.StringNull()
		model.Append = types.BoolNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsLogPublisher() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("syslog_external_server"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "syslog-text-error", "syslog-text-access", "syslog-json-http-operation", "syslog-json-access", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("syslog_facility"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "syslog-based-error", "syslog-text-error", "syslog-based-access", "syslog-text-access", "syslog-json-http-operation", "syslog-json-access", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("syslog_severity"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "syslog-text-error", "syslog-text-access", "syslog-json-http-operation", "syslog-json-access", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("syslog_message_host_name"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "syslog-text-error", "syslog-text-access", "syslog-json-http-operation", "syslog-json-access", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("syslog_message_application_name"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "syslog-text-error", "syslog-text-access", "syslog-json-http-operation", "syslog-json-access", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("queue_size"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "syslog-based-error", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "file-based-trace", "jdbc-based-error", "jdbc-based-access", "common-log-file-http-operation", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-debug", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "syslog-json-http-operation", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "syslog-json-access", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("write_multi_line_messages"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "console-json-error", "file-based-json-audit", "json-access", "syslog-json-http-operation", "json-error", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_reversible_form"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "file-based-json-audit", "file-based-audit", "console-json-audit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("soft_delete_entry_audit_behavior"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "file-based-json-audit", "file-based-audit", "console-json-audit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_operation_purpose_request_control"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "file-based-json-audit", "file-based-audit", "console-json-audit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_intermediate_client_request_control"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "file-based-json-audit", "file-based-audit", "console-json-audit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("obscure_attribute"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "file-based-json-audit", "debug-access", "file-based-audit", "console-json-audit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_attribute"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "file-based-json-audit", "file-based-audit", "console-json-audit"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("suppress_internal_operations"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_product_name"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "operation-timing-access", "admin-alert-access", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "syslog-json-http-operation", "file-based-audit", "json-error", "syslog-json-access", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_instance_name"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "operation-timing-access", "admin-alert-access", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "syslog-json-http-operation", "file-based-audit", "json-error", "syslog-json-access", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_startup_id"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "operation-timing-access", "admin-alert-access", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "syslog-json-http-operation", "file-based-audit", "json-error", "syslog-json-access", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_thread_id"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "operation-timing-access", "admin-alert-access", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "syslog-json-http-operation", "file-based-audit", "json-error", "syslog-json-access", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_requester_dn"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "operation-timing-access", "admin-alert-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "file-based-audit", "syslog-json-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_requester_ip_address"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "operation-timing-access", "admin-alert-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "file-based-audit", "syslog-json-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_request_controls"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "admin-alert-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "file-based-audit", "syslog-json-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_response_controls"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "admin-alert-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "syslog-json-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_replication_change_id"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "admin-alert-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "file-based-audit", "syslog-json-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_security_negotiation"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("suppress_replication_operations"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("connection_criteria"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_criteria"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("result_criteria"),
			path.MatchRoot("type"),
			[]string{"syslog-json-audit", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "file-based-json-audit", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-audit", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_host_name"),
			path.MatchRoot("type"),
			[]string{"syslog-based-error", "syslog-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server_port"),
			path.MatchRoot("type"),
			[]string{"syslog-based-error", "syslog-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("auto_flush"),
			path.MatchRoot("type"),
			[]string{"syslog-based-error", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "common-log-file-http-operation", "syslog-based-access", "file-based-json-audit", "file-based-debug", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("asynchronous"),
			path.MatchRoot("type"),
			[]string{"syslog-based-error", "third-party-file-based-access", "operation-timing-access", "admin-alert-access", "file-based-trace", "common-log-file-http-operation", "syslog-based-access", "file-based-json-audit", "file-based-debug", "file-based-error", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_severity"),
			path.MatchRoot("type"),
			[]string{"syslog-based-error", "jdbc-based-error", "console-json-error", "syslog-text-error", "file-based-error", "third-party-error", "json-error", "groovy-scripted-file-based-error", "third-party-file-based-error", "groovy-scripted-error", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("override_severity"),
			path.MatchRoot("type"),
			[]string{"syslog-based-error", "jdbc-based-error", "console-json-error", "syslog-text-error", "file-based-error", "third-party-error", "json-error", "groovy-scripted-file-based-error", "third-party-file-based-error", "groovy-scripted-error", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_file"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_file_permissions"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("rotation_policy"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("rotation_listener"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("retention_policy"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("compression_mechanism"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sign_log"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encrypt_log"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encryption_settings_definition_id"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("append"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "third-party-http-operation", "third-party-error", "third-party-access", "third-party-file-based-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "third-party-http-operation", "third-party-error", "third-party-access", "third-party-file-based-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("buffer_size"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("time_interval"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "file-based-trace", "common-log-file-http-operation", "file-based-json-audit", "file-based-debug", "file-based-error", "detailed-http-operation", "json-access", "debug-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "third-party-file-based-error", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_connects"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_disconnects"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_client_certificates"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_requests"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "syslog-json-http-operation", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_results"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "syslog-json-http-operation", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_search_entries"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_search_references"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_intermediate_responses"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "operation-timing-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "file-based-audit", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("correlate_requests_and_results"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_entry_criteria"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("search_reference_criteria"),
			path.MatchRoot("type"),
			[]string{"third-party-file-based-access", "admin-alert-access", "jdbc-based-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "third-party-access", "groovy-scripted-file-based-access", "syslog-json-access", "groovy-scripted-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_included_operation_processing_time"),
			path.MatchRoot("type"),
			[]string{"operation-timing-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("min_included_phase_time_nanos"),
			path.MatchRoot("type"),
			[]string{"operation-timing-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_string_length"),
			path.MatchRoot("type"),
			[]string{"operation-timing-access", "admin-alert-access", "file-based-trace", "syslog-based-access", "syslog-text-access", "detailed-http-operation", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_request_details_in_result_messages"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "detailed-http-operation", "json-access", "syslog-json-http-operation", "syslog-json-access", "console-json-http-operation", "console-json-access", "file-based-access", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_assurance_completed"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "debug-access", "syslog-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_request_details_in_search_entry_messages"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_request_details_in_search_reference_messages"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_request_details_in_intermediate_response_messages"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_result_code_names"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_extended_search_request_details"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_add_attribute_names"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_modify_attribute_names"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_search_entry_attribute_names"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("generify_message_strings_when_possible"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-error", "syslog-text-access", "json-access", "json-error", "syslog-json-access", "console-json-access", "file-based-access", "syslog-json-error"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_behavior"),
			path.MatchRoot("type"),
			[]string{"admin-alert-access", "syslog-based-access", "syslog-text-access", "json-access", "syslog-json-access", "console-json-access", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("debug_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("access_token_validator_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("id_token_validator_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("scim_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("consent_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("directory_rest_api_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_message_type"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_path_pattern"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("exclude_path_pattern"),
			path.MatchRoot("type"),
			[]string{"file-based-trace"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("server"),
			path.MatchRoot("type"),
			[]string{"jdbc-based-error", "jdbc-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_field_mapping"),
			path.MatchRoot("type"),
			[]string{"jdbc-based-error", "jdbc-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_table_name"),
			path.MatchRoot("type"),
			[]string{"jdbc-based-error", "jdbc-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("output_location"),
			path.MatchRoot("type"),
			[]string{"console-json-error", "console-json-audit", "console-json-http-operation", "console-json-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("timestamp_precision"),
			path.MatchRoot("type"),
			[]string{"syslog-text-error", "file-based-debug", "file-based-error", "syslog-text-access", "file-based-audit", "file-based-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_debug_level"),
			path.MatchRoot("type"),
			[]string{"file-based-debug"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_debug_category"),
			path.MatchRoot("type"),
			[]string{"file-based-debug"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_omit_method_entry_arguments"),
			path.MatchRoot("type"),
			[]string{"file-based-debug"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_omit_method_return_value"),
			path.MatchRoot("type"),
			[]string{"file-based-debug"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_include_throwable_cause"),
			path.MatchRoot("type"),
			[]string{"file-based-debug"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_throwable_stack_frames"),
			path.MatchRoot("type"),
			[]string{"file-based-debug"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_request_headers"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("suppressed_request_header_name"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_response_headers"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("suppressed_response_header_name"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_request_authorization_type"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_request_cookie_names"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_response_cookie_names"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_request_parameters"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_request_protocol"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("suppressed_request_parameter_name"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_redirect_uri"),
			path.MatchRoot("type"),
			[]string{"detailed-http-operation", "syslog-json-http-operation", "console-json-http-operation", "file-based-json-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("obscure_sensitive_content"),
			path.MatchRoot("type"),
			[]string{"debug-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("debug_aci_enabled"),
			path.MatchRoot("type"),
			[]string{"debug-access"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "groovy-scripted-access", "groovy-scripted-error", "groovy-scripted-http-operation"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "groovy-scripted-access", "groovy-scripted-error", "groovy-scripted-http-operation"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-json-audit",
			[]path.Expression{path.MatchRoot("syslog_external_server"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-text-error",
			[]path.Expression{path.MatchRoot("syslog_external_server"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-text-access",
			[]path.Expression{path.MatchRoot("syslog_external_server"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-json-http-operation",
			[]path.Expression{path.MatchRoot("syslog_external_server"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-json-access",
			[]path.Expression{path.MatchRoot("syslog_external_server"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-json-error",
			[]path.Expression{path.MatchRoot("syslog_external_server"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-based-error",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party-file-based-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"operation-timing-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party-http-operation",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"admin-alert-access",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-trace",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"jdbc-based-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("server"), path.MatchRoot("log_field_mapping")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"jdbc-based-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("server"), path.MatchRoot("log_field_mapping")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"common-log-file-http-operation",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"syslog-based-access",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-json-audit",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-debug",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"detailed-http-operation",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"json-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"debug-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-audit",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"json-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted-file-based-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file"), path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted-file-based-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file"), path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party-file-based-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"console-json-audit",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"console-json-http-operation",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-access",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted-error",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("script_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based-json-http-operation",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("log_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"groovy-scripted-http-operation",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("script_class")},
		),
	}
}

// Add config validators
func (r logPublisherResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogPublisher()
}

// Add config validators
func (r defaultLogPublisherResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsLogPublisher()
}

// Add optional fields to create request for syslog-json-audit log-publisher
func addOptionalSyslogJsonAuditLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonAuditLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		addRequest.SyslogMessageHostName = plan.SyslogMessageHostName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		addRequest.SyslogMessageApplicationName = plan.SyslogMessageApplicationName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseReversibleForm) {
		addRequest.UseReversibleForm = plan.UseReversibleForm.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteEntryAuditBehavior) {
		softDeleteEntryAuditBehavior, err := client.NewEnumlogPublisherSyslogJsonAuditSoftDeleteEntryAuditBehaviorPropFromValue(plan.SoftDeleteEntryAuditBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SoftDeleteEntryAuditBehavior = softDeleteEntryAuditBehavior
	}
	if internaltypes.IsDefined(plan.IncludeOperationPurposeRequestControl) {
		addRequest.IncludeOperationPurposeRequestControl = plan.IncludeOperationPurposeRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeIntermediateClientRequestControl) {
		addRequest.IncludeIntermediateClientRequestControl = plan.IncludeIntermediateClientRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureAttribute) {
		var slice []string
		plan.ObscureAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ObscureAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludeAttribute) {
		var slice []string
		plan.ExcludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-based-error log-publisher
func addOptionalSyslogBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogBasedErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerHostName) {
		addRequest.ServerHostName = plan.ServerHostName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		intVal, err := strconv.ParseInt(plan.SyslogFacility.ValueString(), 0, 64)
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = &intVal
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for third-party-file-based-access log-publisher
func addOptionalThirdPartyFileBasedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddThirdPartyFileBasedAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for operation-timing-access log-publisher
func addOptionalOperationTimingAccessLogPublisherFields(ctx context.Context, addRequest *client.AddOperationTimingAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinIncludedOperationProcessingTime) {
		addRequest.MinIncludedOperationProcessingTime = plan.MinIncludedOperationProcessingTime.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MinIncludedPhaseTimeNanos) {
		addRequest.MinIncludedPhaseTimeNanos = plan.MinIncludedPhaseTimeNanos.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for third-party-http-operation log-publisher
func addOptionalThirdPartyHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddThirdPartyHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for admin-alert-access log-publisher
func addOptionalAdminAlertAccessLogPublisherFields(ctx context.Context, addRequest *client.AddAdminAlertAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		addRequest.IncludeRequestDetailsInSearchEntryMessages = plan.IncludeRequestDetailsInSearchEntryMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		addRequest.IncludeResultCodeNames = plan.IncludeResultCodeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		addRequest.IncludeExtendedSearchRequestDetails = plan.IncludeExtendedSearchRequestDetails.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		addRequest.IncludeAddAttributeNames = plan.IncludeAddAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		addRequest.IncludeModifyAttributeNames = plan.IncludeModifyAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		addRequest.IncludeSearchEntryAttributeNames = plan.IncludeSearchEntryAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		addRequest.LogFieldBehavior = plan.LogFieldBehavior.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-trace log-publisher
func addOptionalFileBasedTraceLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedTraceLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DebugMessageType) {
		var slice []string
		plan.DebugMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDebugMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDebugMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DebugMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.HttpMessageType) {
		var slice []string
		plan.HttpMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherHttpMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherHttpMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.HttpMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.AccessTokenValidatorMessageType) {
		var slice []string
		plan.AccessTokenValidatorMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherAccessTokenValidatorMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherAccessTokenValidatorMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AccessTokenValidatorMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.IdTokenValidatorMessageType) {
		var slice []string
		plan.IdTokenValidatorMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherIdTokenValidatorMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherIdTokenValidatorMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.IdTokenValidatorMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.ScimMessageType) {
		var slice []string
		plan.ScimMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherScimMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherScimMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ScimMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.ConsentMessageType) {
		var slice []string
		plan.ConsentMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherConsentMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherConsentMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ConsentMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.DirectoryRESTAPIMessageType) {
		var slice []string
		plan.DirectoryRESTAPIMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDirectoryRESTAPIMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDirectoryRESTAPIMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DirectoryRESTAPIMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.ExtensionMessageType) {
		var slice []string
		plan.ExtensionMessageType.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherExtensionMessageTypeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherExtensionMessageTypePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.ExtensionMessageType = enumSlice
	}
	if internaltypes.IsDefined(plan.IncludePathPattern) {
		var slice []string
		plan.IncludePathPattern.ElementsAs(ctx, &slice, false)
		addRequest.IncludePathPattern = slice
	}
	if internaltypes.IsDefined(plan.ExcludePathPattern) {
		var slice []string
		plan.ExcludePathPattern.ElementsAs(ctx, &slice, false)
		addRequest.ExcludePathPattern = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for jdbc-based-error log-publisher
func addOptionalJdbcBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddJdbcBasedErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogTableName) {
		addRequest.LogTableName = plan.LogTableName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for jdbc-based-access log-publisher
func addOptionalJdbcBasedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddJdbcBasedAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogTableName) {
		addRequest.LogTableName = plan.LogTableName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for common-log-file-http-operation log-publisher
func addOptionalCommonLogFileHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddCommonLogFileHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-text-error log-publisher
func addOptionalSyslogTextErrorLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogTextErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		addRequest.SyslogMessageHostName = plan.SyslogMessageHostName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		addRequest.SyslogMessageApplicationName = plan.SyslogMessageApplicationName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-based-access log-publisher
func addOptionalSyslogBasedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogBasedAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ServerHostName) {
		addRequest.ServerHostName = plan.ServerHostName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		intVal, err := strconv.ParseInt(plan.SyslogFacility.ValueString(), 0, 64)
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = &intVal
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		addRequest.IncludeRequestDetailsInSearchEntryMessages = plan.IncludeRequestDetailsInSearchEntryMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		addRequest.IncludeResultCodeNames = plan.IncludeResultCodeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		addRequest.IncludeExtendedSearchRequestDetails = plan.IncludeExtendedSearchRequestDetails.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		addRequest.IncludeAddAttributeNames = plan.IncludeAddAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		addRequest.IncludeModifyAttributeNames = plan.IncludeModifyAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		addRequest.IncludeSearchEntryAttributeNames = plan.IncludeSearchEntryAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		addRequest.LogFieldBehavior = plan.LogFieldBehavior.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-json-audit log-publisher
func addOptionalFileBasedJsonAuditLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedJsonAuditLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseReversibleForm) {
		addRequest.UseReversibleForm = plan.UseReversibleForm.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteEntryAuditBehavior) {
		softDeleteEntryAuditBehavior, err := client.NewEnumlogPublisherFileBasedJsonAuditSoftDeleteEntryAuditBehaviorPropFromValue(plan.SoftDeleteEntryAuditBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SoftDeleteEntryAuditBehavior = softDeleteEntryAuditBehavior
	}
	if internaltypes.IsDefined(plan.IncludeOperationPurposeRequestControl) {
		addRequest.IncludeOperationPurposeRequestControl = plan.IncludeOperationPurposeRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeIntermediateClientRequestControl) {
		addRequest.IncludeIntermediateClientRequestControl = plan.IncludeIntermediateClientRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureAttribute) {
		var slice []string
		plan.ObscureAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ObscureAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludeAttribute) {
		var slice []string
		plan.ExcludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-debug log-publisher
func addOptionalFileBasedDebugLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedDebugLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultDebugLevel) {
		defaultDebugLevel, err := client.NewEnumlogPublisherDefaultDebugLevelPropFromValue(plan.DefaultDebugLevel.ValueString())
		if err != nil {
			return err
		}
		addRequest.DefaultDebugLevel = defaultDebugLevel
	}
	if internaltypes.IsDefined(plan.DefaultDebugCategory) {
		var slice []string
		plan.DefaultDebugCategory.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultDebugCategoryProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultDebugCategoryPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultDebugCategory = enumSlice
	}
	if internaltypes.IsDefined(plan.DefaultOmitMethodEntryArguments) {
		addRequest.DefaultOmitMethodEntryArguments = plan.DefaultOmitMethodEntryArguments.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DefaultOmitMethodReturnValue) {
		addRequest.DefaultOmitMethodReturnValue = plan.DefaultOmitMethodReturnValue.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DefaultIncludeThrowableCause) {
		addRequest.DefaultIncludeThrowableCause = plan.DefaultIncludeThrowableCause.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DefaultThrowableStackFrames) {
		addRequest.DefaultThrowableStackFrames = plan.DefaultThrowableStackFrames.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-error log-publisher
func addOptionalFileBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for third-party-error log-publisher
func addOptionalThirdPartyErrorLogPublisherFields(ctx context.Context, addRequest *client.AddThirdPartyErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-text-access log-publisher
func addOptionalSyslogTextAccessLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogTextAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		addRequest.SyslogMessageHostName = plan.SyslogMessageHostName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		addRequest.SyslogMessageApplicationName = plan.SyslogMessageApplicationName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		addRequest.IncludeRequestDetailsInSearchEntryMessages = plan.IncludeRequestDetailsInSearchEntryMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		addRequest.IncludeResultCodeNames = plan.IncludeResultCodeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		addRequest.IncludeExtendedSearchRequestDetails = plan.IncludeExtendedSearchRequestDetails.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		addRequest.IncludeAddAttributeNames = plan.IncludeAddAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		addRequest.IncludeModifyAttributeNames = plan.IncludeModifyAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		addRequest.IncludeSearchEntryAttributeNames = plan.IncludeSearchEntryAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		addRequest.LogFieldBehavior = plan.LogFieldBehavior.ValueStringPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for detailed-http-operation log-publisher
func addOptionalDetailedHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddDetailedHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestHeaders) {
		logRequestHeaders, err := client.NewEnumlogPublisherLogRequestHeadersPropFromValue(plan.LogRequestHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestHeaders = logRequestHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedRequestHeaderName) {
		var slice []string
		plan.SuppressedRequestHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestHeaderName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogResponseHeaders) {
		logResponseHeaders, err := client.NewEnumlogPublisherLogResponseHeadersPropFromValue(plan.LogResponseHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogResponseHeaders = logResponseHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedResponseHeaderName) {
		var slice []string
		plan.SuppressedResponseHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedResponseHeaderName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestAuthorizationType) {
		addRequest.LogRequestAuthorizationType = plan.LogRequestAuthorizationType.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequestCookieNames) {
		addRequest.LogRequestCookieNames = plan.LogRequestCookieNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResponseCookieNames) {
		addRequest.LogResponseCookieNames = plan.LogResponseCookieNames.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestParameters) {
		logRequestParameters, err := client.NewEnumlogPublisherLogRequestParametersPropFromValue(plan.LogRequestParameters.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestParameters = logRequestParameters
	}
	if internaltypes.IsDefined(plan.LogRequestProtocol) {
		addRequest.LogRequestProtocol = plan.LogRequestProtocol.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressedRequestParameterName) {
		var slice []string
		plan.SuppressedRequestParameterName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestParameterName = slice
	}
	if internaltypes.IsDefined(plan.LogRedirectURI) {
		addRequest.LogRedirectURI = plan.LogRedirectURI.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for json-access log-publisher
func addOptionalJsonAccessLogPublisherFields(ctx context.Context, addRequest *client.AddJsonAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		addRequest.IncludeRequestDetailsInSearchEntryMessages = plan.IncludeRequestDetailsInSearchEntryMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		addRequest.IncludeResultCodeNames = plan.IncludeResultCodeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		addRequest.IncludeExtendedSearchRequestDetails = plan.IncludeExtendedSearchRequestDetails.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		addRequest.IncludeAddAttributeNames = plan.IncludeAddAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		addRequest.IncludeModifyAttributeNames = plan.IncludeModifyAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		addRequest.IncludeSearchEntryAttributeNames = plan.IncludeSearchEntryAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		addRequest.LogFieldBehavior = plan.LogFieldBehavior.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for debug-access log-publisher
func addOptionalDebugAccessLogPublisherFields(ctx context.Context, addRequest *client.AddDebugAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureSensitiveContent) {
		addRequest.ObscureSensitiveContent = plan.ObscureSensitiveContent.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureAttribute) {
		var slice []string
		plan.ObscureAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ObscureAttribute = slice
	}
	if internaltypes.IsDefined(plan.DebugACIEnabled) {
		addRequest.DebugACIEnabled = plan.DebugACIEnabled.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-json-http-operation log-publisher
func addOptionalSyslogJsonHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		addRequest.SyslogMessageHostName = plan.SyslogMessageHostName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		addRequest.SyslogMessageApplicationName = plan.SyslogMessageApplicationName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestHeaders) {
		logRequestHeaders, err := client.NewEnumlogPublisherLogRequestHeadersPropFromValue(plan.LogRequestHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestHeaders = logRequestHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedRequestHeaderName) {
		var slice []string
		plan.SuppressedRequestHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestHeaderName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogResponseHeaders) {
		logResponseHeaders, err := client.NewEnumlogPublisherLogResponseHeadersPropFromValue(plan.LogResponseHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogResponseHeaders = logResponseHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedResponseHeaderName) {
		var slice []string
		plan.SuppressedResponseHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedResponseHeaderName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestAuthorizationType) {
		addRequest.LogRequestAuthorizationType = plan.LogRequestAuthorizationType.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequestCookieNames) {
		addRequest.LogRequestCookieNames = plan.LogRequestCookieNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResponseCookieNames) {
		addRequest.LogResponseCookieNames = plan.LogResponseCookieNames.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestParameters) {
		logRequestParameters, err := client.NewEnumlogPublisherLogRequestParametersPropFromValue(plan.LogRequestParameters.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestParameters = logRequestParameters
	}
	if internaltypes.IsDefined(plan.SuppressedRequestParameterName) {
		var slice []string
		plan.SuppressedRequestParameterName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestParameterName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestProtocol) {
		addRequest.LogRequestProtocol = plan.LogRequestProtocol.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRedirectURI) {
		addRequest.LogRedirectURI = plan.LogRedirectURI.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for third-party-access log-publisher
func addOptionalThirdPartyAccessLogPublisherFields(ctx context.Context, addRequest *client.AddThirdPartyAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-audit log-publisher
func addOptionalFileBasedAuditLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedAuditLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseReversibleForm) {
		addRequest.UseReversibleForm = plan.UseReversibleForm.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteEntryAuditBehavior) {
		softDeleteEntryAuditBehavior, err := client.NewEnumlogPublisherFileBasedAuditSoftDeleteEntryAuditBehaviorPropFromValue(plan.SoftDeleteEntryAuditBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SoftDeleteEntryAuditBehavior = softDeleteEntryAuditBehavior
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeOperationPurposeRequestControl) {
		addRequest.IncludeOperationPurposeRequestControl = plan.IncludeOperationPurposeRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeIntermediateClientRequestControl) {
		addRequest.IncludeIntermediateClientRequestControl = plan.IncludeIntermediateClientRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureAttribute) {
		var slice []string
		plan.ObscureAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ObscureAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludeAttribute) {
		var slice []string
		plan.ExcludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for json-error log-publisher
func addOptionalJsonErrorLogPublisherFields(ctx context.Context, addRequest *client.AddJsonErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for groovy-scripted-file-based-access log-publisher
func addOptionalGroovyScriptedFileBasedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddGroovyScriptedFileBasedAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for groovy-scripted-file-based-error log-publisher
func addOptionalGroovyScriptedFileBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddGroovyScriptedFileBasedErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-json-access log-publisher
func addOptionalSyslogJsonAccessLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		addRequest.SyslogMessageHostName = plan.SyslogMessageHostName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		addRequest.SyslogMessageApplicationName = plan.SyslogMessageApplicationName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		addRequest.IncludeRequestDetailsInSearchEntryMessages = plan.IncludeRequestDetailsInSearchEntryMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		addRequest.IncludeResultCodeNames = plan.IncludeResultCodeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		addRequest.IncludeExtendedSearchRequestDetails = plan.IncludeExtendedSearchRequestDetails.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		addRequest.IncludeAddAttributeNames = plan.IncludeAddAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		addRequest.IncludeModifyAttributeNames = plan.IncludeModifyAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		addRequest.IncludeSearchEntryAttributeNames = plan.IncludeSearchEntryAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		addRequest.LogFieldBehavior = plan.LogFieldBehavior.ValueStringPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for groovy-scripted-access log-publisher
func addOptionalGroovyScriptedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddGroovyScriptedAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for third-party-file-based-error log-publisher
func addOptionalThirdPartyFileBasedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddThirdPartyFileBasedErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for console-json-audit log-publisher
func addOptionalConsoleJsonAuditLogPublisherFields(ctx context.Context, addRequest *client.AddConsoleJsonAuditLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OutputLocation) {
		outputLocation, err := client.NewEnumlogPublisherOutputLocationPropFromValue(plan.OutputLocation.ValueString())
		if err != nil {
			return err
		}
		addRequest.OutputLocation = outputLocation
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseReversibleForm) {
		addRequest.UseReversibleForm = plan.UseReversibleForm.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteEntryAuditBehavior) {
		softDeleteEntryAuditBehavior, err := client.NewEnumlogPublisherConsoleJsonAuditSoftDeleteEntryAuditBehaviorPropFromValue(plan.SoftDeleteEntryAuditBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.SoftDeleteEntryAuditBehavior = softDeleteEntryAuditBehavior
	}
	if internaltypes.IsDefined(plan.IncludeOperationPurposeRequestControl) {
		addRequest.IncludeOperationPurposeRequestControl = plan.IncludeOperationPurposeRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeIntermediateClientRequestControl) {
		addRequest.IncludeIntermediateClientRequestControl = plan.IncludeIntermediateClientRequestControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ObscureAttribute) {
		var slice []string
		plan.ObscureAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ObscureAttribute = slice
	}
	if internaltypes.IsDefined(plan.ExcludeAttribute) {
		var slice []string
		plan.ExcludeAttribute.ElementsAs(ctx, &slice, false)
		addRequest.ExcludeAttribute = slice
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for console-json-http-operation log-publisher
func addOptionalConsoleJsonHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddConsoleJsonHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OutputLocation) {
		outputLocation, err := client.NewEnumlogPublisherOutputLocationPropFromValue(plan.OutputLocation.ValueString())
		if err != nil {
			return err
		}
		addRequest.OutputLocation = outputLocation
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestHeaders) {
		logRequestHeaders, err := client.NewEnumlogPublisherLogRequestHeadersPropFromValue(plan.LogRequestHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestHeaders = logRequestHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedRequestHeaderName) {
		var slice []string
		plan.SuppressedRequestHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestHeaderName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogResponseHeaders) {
		logResponseHeaders, err := client.NewEnumlogPublisherLogResponseHeadersPropFromValue(plan.LogResponseHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogResponseHeaders = logResponseHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedResponseHeaderName) {
		var slice []string
		plan.SuppressedResponseHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedResponseHeaderName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestAuthorizationType) {
		addRequest.LogRequestAuthorizationType = plan.LogRequestAuthorizationType.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequestCookieNames) {
		addRequest.LogRequestCookieNames = plan.LogRequestCookieNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResponseCookieNames) {
		addRequest.LogResponseCookieNames = plan.LogResponseCookieNames.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestParameters) {
		logRequestParameters, err := client.NewEnumlogPublisherLogRequestParametersPropFromValue(plan.LogRequestParameters.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestParameters = logRequestParameters
	}
	if internaltypes.IsDefined(plan.SuppressedRequestParameterName) {
		var slice []string
		plan.SuppressedRequestParameterName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestParameterName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestProtocol) {
		addRequest.LogRequestProtocol = plan.LogRequestProtocol.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRedirectURI) {
		addRequest.LogRedirectURI = plan.LogRedirectURI.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-access log-publisher
func addOptionalFileBasedAccessLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedAccessLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimestampPrecision) {
		timestampPrecision, err := client.NewEnumlogPublisherTimestampPrecisionPropFromValue(plan.TimestampPrecision.ValueString())
		if err != nil {
			return err
		}
		addRequest.TimestampPrecision = timestampPrecision
	}
	if internaltypes.IsDefined(plan.LogConnects) {
		addRequest.LogConnects = plan.LogConnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		addRequest.LogDisconnects = plan.LogDisconnects.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		addRequest.LogAssuranceCompleted = plan.LogAssuranceCompleted.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		addRequest.IncludeRequesterDN = plan.IncludeRequesterDN.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		addRequest.IncludeRequesterIPAddress = plan.IncludeRequesterIPAddress.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		addRequest.IncludeRequestDetailsInSearchEntryMessages = plan.IncludeRequestDetailsInSearchEntryMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		addRequest.IncludeResultCodeNames = plan.IncludeResultCodeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		addRequest.IncludeExtendedSearchRequestDetails = plan.IncludeExtendedSearchRequestDetails.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		addRequest.IncludeAddAttributeNames = plan.IncludeAddAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		addRequest.IncludeModifyAttributeNames = plan.IncludeModifyAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		addRequest.IncludeSearchEntryAttributeNames = plan.IncludeSearchEntryAttributeNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		addRequest.IncludeRequestControls = plan.IncludeRequestControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		addRequest.IncludeResponseControls = plan.IncludeResponseControls.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		addRequest.IncludeReplicationChangeID = plan.IncludeReplicationChangeID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		addRequest.MaxStringLength = plan.MaxStringLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		addRequest.LogFieldBehavior = plan.LogFieldBehavior.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		addRequest.LogSecurityNegotiation = plan.LogSecurityNegotiation.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		addRequest.LogClientCertificates = plan.LogClientCertificates.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		addRequest.LogSearchEntries = plan.LogSearchEntries.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		addRequest.LogSearchReferences = plan.LogSearchReferences.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		addRequest.LogIntermediateResponses = plan.LogIntermediateResponses.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		addRequest.SuppressInternalOperations = plan.SuppressInternalOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		addRequest.SuppressReplicationOperations = plan.SuppressReplicationOperations.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		addRequest.CorrelateRequestsAndResults = plan.CorrelateRequestsAndResults.ValueBoolPointer()
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
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		addRequest.ResultCriteria = plan.ResultCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		addRequest.SearchEntryCriteria = plan.SearchEntryCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		addRequest.SearchReferenceCriteria = plan.SearchReferenceCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for groovy-scripted-error log-publisher
func addOptionalGroovyScriptedErrorLogPublisherFields(ctx context.Context, addRequest *client.AddGroovyScriptedErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for file-based-json-http-operation log-publisher
func addOptionalFileBasedJsonHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddFileBasedJsonHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		addRequest.LogFilePermissions = plan.LogFilePermissions.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.RotationPolicy) {
		var slice []string
		plan.RotationPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RotationPolicy = slice
	}
	if internaltypes.IsDefined(plan.RotationListener) {
		var slice []string
		plan.RotationListener.ElementsAs(ctx, &slice, false)
		addRequest.RotationListener = slice
	}
	if internaltypes.IsDefined(plan.RetentionPolicy) {
		var slice []string
		plan.RetentionPolicy.ElementsAs(ctx, &slice, false)
		addRequest.RetentionPolicy = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CompressionMechanism) {
		compressionMechanism, err := client.NewEnumlogPublisherCompressionMechanismPropFromValue(plan.CompressionMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.CompressionMechanism = compressionMechanism
	}
	if internaltypes.IsDefined(plan.SignLog) {
		addRequest.SignLog = plan.SignLog.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		addRequest.EncryptLog = plan.EncryptLog.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.Append) {
		addRequest.Append = plan.Append.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		addRequest.Asynchronous = plan.Asynchronous.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		addRequest.AutoFlush = plan.AutoFlush.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		addRequest.BufferSize = plan.BufferSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		addRequest.TimeInterval = plan.TimeInterval.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		addRequest.LogRequests = plan.LogRequests.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResults) {
		addRequest.LogResults = plan.LogResults.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		addRequest.IncludeRequestDetailsInResultMessages = plan.IncludeRequestDetailsInResultMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestHeaders) {
		logRequestHeaders, err := client.NewEnumlogPublisherLogRequestHeadersPropFromValue(plan.LogRequestHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestHeaders = logRequestHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedRequestHeaderName) {
		var slice []string
		plan.SuppressedRequestHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestHeaderName = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogResponseHeaders) {
		logResponseHeaders, err := client.NewEnumlogPublisherLogResponseHeadersPropFromValue(plan.LogResponseHeaders.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogResponseHeaders = logResponseHeaders
	}
	if internaltypes.IsDefined(plan.SuppressedResponseHeaderName) {
		var slice []string
		plan.SuppressedResponseHeaderName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedResponseHeaderName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestAuthorizationType) {
		addRequest.LogRequestAuthorizationType = plan.LogRequestAuthorizationType.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRequestCookieNames) {
		addRequest.LogRequestCookieNames = plan.LogRequestCookieNames.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogResponseCookieNames) {
		addRequest.LogResponseCookieNames = plan.LogResponseCookieNames.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogRequestParameters) {
		logRequestParameters, err := client.NewEnumlogPublisherLogRequestParametersPropFromValue(plan.LogRequestParameters.ValueString())
		if err != nil {
			return err
		}
		addRequest.LogRequestParameters = logRequestParameters
	}
	if internaltypes.IsDefined(plan.SuppressedRequestParameterName) {
		var slice []string
		plan.SuppressedRequestParameterName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestParameterName = slice
	}
	if internaltypes.IsDefined(plan.LogRequestProtocol) {
		addRequest.LogRequestProtocol = plan.LogRequestProtocol.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.LogRedirectURI) {
		addRequest.LogRedirectURI = plan.LogRedirectURI.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		addRequest.WriteMultiLineMessages = plan.WriteMultiLineMessages.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for syslog-json-error log-publisher
func addOptionalSyslogJsonErrorLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonErrorLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.DefaultSeverity) {
		var slice []string
		plan.DefaultSeverity.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumlogPublisherDefaultSeverityProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumlogPublisherDefaultSeverityPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefaultSeverity = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogFacility) {
		syslogFacility, err := client.NewEnumlogPublisherSyslogFacilityPropFromValue(plan.SyslogFacility.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogFacility = syslogFacility
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogSeverity) {
		syslogSeverity, err := client.NewEnumlogPublisherSyslogSeverityPropFromValue(plan.SyslogSeverity.ValueString())
		if err != nil {
			return err
		}
		addRequest.SyslogSeverity = syslogSeverity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageHostName) {
		addRequest.SyslogMessageHostName = plan.SyslogMessageHostName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		addRequest.SyslogMessageApplicationName = plan.SyslogMessageApplicationName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		addRequest.QueueSize = plan.QueueSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		addRequest.IncludeProductName = plan.IncludeProductName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		addRequest.IncludeInstanceName = plan.IncludeInstanceName.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		addRequest.IncludeStartupID = plan.IncludeStartupID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		addRequest.IncludeThreadID = plan.IncludeThreadID.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		addRequest.GenerifyMessageStringsWhenPossible = plan.GenerifyMessageStringsWhenPossible.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.OverrideSeverity) {
		var slice []string
		plan.OverrideSeverity.ElementsAs(ctx, &slice, false)
		addRequest.OverrideSeverity = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Add optional fields to create request for groovy-scripted-http-operation log-publisher
func addOptionalGroovyScriptedHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddGroovyScriptedHttpOperationLogPublisherRequest, plan logPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LoggingErrorBehavior) {
		loggingErrorBehavior, err := client.NewEnumlogPublisherLoggingErrorBehaviorPropFromValue(plan.LoggingErrorBehavior.ValueString())
		if err != nil {
			return err
		}
		addRequest.LoggingErrorBehavior = loggingErrorBehavior
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateLogPublisherUnknownValues(model *logPublisherResourceModel) {
	if model.IdTokenValidatorMessageType.IsUnknown() || model.IdTokenValidatorMessageType.IsNull() {
		model.IdTokenValidatorMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SuppressedResponseHeaderName.IsUnknown() || model.SuppressedResponseHeaderName.IsNull() {
		model.SuppressedResponseHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DefaultDebugCategory.IsUnknown() || model.DefaultDebugCategory.IsNull() {
		model.DefaultDebugCategory, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ScimMessageType.IsUnknown() || model.ScimMessageType.IsNull() {
		model.ScimMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RotationPolicy.IsUnknown() || model.RotationPolicy.IsNull() {
		model.RotationPolicy, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RetentionPolicy.IsUnknown() || model.RetentionPolicy.IsNull() {
		model.RetentionPolicy, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SuppressedRequestHeaderName.IsUnknown() || model.SuppressedRequestHeaderName.IsNull() {
		model.SuppressedRequestHeaderName, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.IncludePathPattern.IsUnknown() || model.IncludePathPattern.IsNull() {
		model.IncludePathPattern, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludeAttribute.IsUnknown() || model.ExcludeAttribute.IsNull() {
		model.ExcludeAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.HttpMessageType.IsUnknown() || model.HttpMessageType.IsNull() {
		model.HttpMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ScriptArgument.IsUnknown() || model.ScriptArgument.IsNull() {
		model.ScriptArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DefaultSeverity.IsUnknown() || model.DefaultSeverity.IsNull() {
		model.DefaultSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.RotationListener.IsUnknown() || model.RotationListener.IsNull() {
		model.RotationListener, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExcludePathPattern.IsUnknown() || model.ExcludePathPattern.IsNull() {
		model.ExcludePathPattern, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AccessTokenValidatorMessageType.IsUnknown() || model.AccessTokenValidatorMessageType.IsNull() {
		model.AccessTokenValidatorMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ExtensionMessageType.IsUnknown() || model.ExtensionMessageType.IsNull() {
		model.ExtensionMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SuppressedRequestParameterName.IsUnknown() || model.SuppressedRequestParameterName.IsNull() {
		model.SuppressedRequestParameterName, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ConsentMessageType.IsUnknown() || model.ConsentMessageType.IsNull() {
		model.ConsentMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.OverrideSeverity.IsUnknown() || model.OverrideSeverity.IsNull() {
		model.OverrideSeverity, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DirectoryRESTAPIMessageType.IsUnknown() || model.DirectoryRESTAPIMessageType.IsNull() {
		model.DirectoryRESTAPIMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ObscureAttribute.IsUnknown() || model.ObscureAttribute.IsNull() {
		model.ObscureAttribute, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SyslogExternalServer.IsUnknown() || model.SyslogExternalServer.IsNull() {
		model.SyslogExternalServer, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DebugMessageType.IsUnknown() || model.DebugMessageType.IsNull() {
		model.DebugMessageType, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *logPublisherResourceModel) populateAllComputedStringAttributes() {
	if model.LogRequestParameters.IsUnknown() || model.LogRequestParameters.IsNull() {
		model.LogRequestParameters = types.StringValue("")
	}
	if model.SyslogFacility.IsUnknown() || model.SyslogFacility.IsNull() {
		model.SyslogFacility = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.Server.IsUnknown() || model.Server.IsNull() {
		model.Server = types.StringValue("")
	}
	if model.SoftDeleteEntryAuditBehavior.IsUnknown() || model.SoftDeleteEntryAuditBehavior.IsNull() {
		model.SoftDeleteEntryAuditBehavior = types.StringValue("")
	}
	if model.CompressionMechanism.IsUnknown() || model.CompressionMechanism.IsNull() {
		model.CompressionMechanism = types.StringValue("")
	}
	if model.EncryptionSettingsDefinitionID.IsUnknown() || model.EncryptionSettingsDefinitionID.IsNull() {
		model.EncryptionSettingsDefinitionID = types.StringValue("")
	}
	if model.MinIncludedOperationProcessingTime.IsUnknown() || model.MinIncludedOperationProcessingTime.IsNull() {
		model.MinIncludedOperationProcessingTime = types.StringValue("")
	}
	if model.ConnectionCriteria.IsUnknown() || model.ConnectionCriteria.IsNull() {
		model.ConnectionCriteria = types.StringValue("")
	}
	if model.SyslogMessageHostName.IsUnknown() || model.SyslogMessageHostName.IsNull() {
		model.SyslogMessageHostName = types.StringValue("")
	}
	if model.LogFile.IsUnknown() || model.LogFile.IsNull() {
		model.LogFile = types.StringValue("")
	}
	if model.SearchEntryCriteria.IsUnknown() || model.SearchEntryCriteria.IsNull() {
		model.SearchEntryCriteria = types.StringValue("")
	}
	if model.LogFieldBehavior.IsUnknown() || model.LogFieldBehavior.IsNull() {
		model.LogFieldBehavior = types.StringValue("")
	}
	if model.RequestCriteria.IsUnknown() || model.RequestCriteria.IsNull() {
		model.RequestCriteria = types.StringValue("")
	}
	if model.LogFilePermissions.IsUnknown() || model.LogFilePermissions.IsNull() {
		model.LogFilePermissions = types.StringValue("")
	}
	if model.SearchReferenceCriteria.IsUnknown() || model.SearchReferenceCriteria.IsNull() {
		model.SearchReferenceCriteria = types.StringValue("")
	}
	if model.ServerHostName.IsUnknown() || model.ServerHostName.IsNull() {
		model.ServerHostName = types.StringValue("")
	}
	if model.DefaultDebugLevel.IsUnknown() || model.DefaultDebugLevel.IsNull() {
		model.DefaultDebugLevel = types.StringValue("")
	}
	if model.TimestampPrecision.IsUnknown() || model.TimestampPrecision.IsNull() {
		model.TimestampPrecision = types.StringValue("")
	}
	if model.SyslogMessageApplicationName.IsUnknown() || model.SyslogMessageApplicationName.IsNull() {
		model.SyslogMessageApplicationName = types.StringValue("")
	}
	if model.LogFieldMapping.IsUnknown() || model.LogFieldMapping.IsNull() {
		model.LogFieldMapping = types.StringValue("")
	}
	if model.SyslogSeverity.IsUnknown() || model.SyslogSeverity.IsNull() {
		model.SyslogSeverity = types.StringValue("")
	}
	if model.TimeInterval.IsUnknown() || model.TimeInterval.IsNull() {
		model.TimeInterval = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.OutputLocation.IsUnknown() || model.OutputLocation.IsNull() {
		model.OutputLocation = types.StringValue("")
	}
	if model.LogTableName.IsUnknown() || model.LogTableName.IsNull() {
		model.LogTableName = types.StringValue("")
	}
	if model.BufferSize.IsUnknown() || model.BufferSize.IsNull() {
		model.BufferSize = types.StringValue("")
	}
	if model.LogResponseHeaders.IsUnknown() || model.LogResponseHeaders.IsNull() {
		model.LogResponseHeaders = types.StringValue("")
	}
	if model.LoggingErrorBehavior.IsUnknown() || model.LoggingErrorBehavior.IsNull() {
		model.LoggingErrorBehavior = types.StringValue("")
	}
	if model.LogRequestHeaders.IsUnknown() || model.LogRequestHeaders.IsNull() {
		model.LogRequestHeaders = types.StringValue("")
	}
	if model.ResultCriteria.IsUnknown() || model.ResultCriteria.IsNull() {
		model.ResultCriteria = types.StringValue("")
	}
	if model.ScriptClass.IsUnknown() || model.ScriptClass.IsNull() {
		model.ScriptClass = types.StringValue("")
	}
}

// Read a SyslogJsonAuditLogPublisherResponse object into the model struct
func readSyslogJsonAuditLogPublisherResponse(ctx context.Context, r *client.SyslogJsonAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), true)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogBasedErrorLogPublisherResponse object into the model struct
func readSyslogBasedErrorLogPublisherResponse(ctx context.Context, r *client.SyslogBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.SyslogFacility = types.StringValue(strconv.FormatInt(r.SyslogFacility, 10))
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ThirdPartyFileBasedAccessLogPublisherResponse object into the model struct
func readThirdPartyFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.ThirdPartyFileBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-file-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a OperationTimingAccessLogPublisherResponse object into the model struct
func readOperationTimingAccessLogPublisherResponse(ctx context.Context, r *client.OperationTimingAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("operation-timing-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.MinIncludedOperationProcessingTime = internaltypes.StringTypeOrNil(r.MinIncludedOperationProcessingTime, internaltypes.IsEmptyString(expectedValues.MinIncludedOperationProcessingTime))
	config.CheckMismatchedPDFormattedAttributes("min_included_operation_processing_time",
		expectedValues.MinIncludedOperationProcessingTime, state.MinIncludedOperationProcessingTime, diagnostics)
	state.MinIncludedPhaseTimeNanos = internaltypes.Int64TypeOrNil(r.MinIncludedPhaseTimeNanos)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ThirdPartyHttpOperationLogPublisherResponse object into the model struct
func readThirdPartyHttpOperationLogPublisherResponse(ctx context.Context, r *client.ThirdPartyHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a AdminAlertAccessLogPublisherResponse object into the model struct
func readAdminAlertAccessLogPublisherResponse(ctx context.Context, r *client.AdminAlertAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("admin-alert-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedTraceLogPublisherResponse object into the model struct
func readFileBasedTraceLogPublisherResponse(ctx context.Context, r *client.FileBasedTraceLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-trace")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.DebugMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDebugMessageTypeProp(r.DebugMessageType))
	state.HttpMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherHttpMessageTypeProp(r.HttpMessageType))
	state.AccessTokenValidatorMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherAccessTokenValidatorMessageTypeProp(r.AccessTokenValidatorMessageType))
	state.IdTokenValidatorMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherIdTokenValidatorMessageTypeProp(r.IdTokenValidatorMessageType))
	state.ScimMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherScimMessageTypeProp(r.ScimMessageType))
	state.ConsentMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherConsentMessageTypeProp(r.ConsentMessageType))
	state.DirectoryRESTAPIMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDirectoryRESTAPIMessageTypeProp(r.DirectoryRESTAPIMessageType))
	state.ExtensionMessageType = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherExtensionMessageTypeProp(r.ExtensionMessageType))
	state.IncludePathPattern = internaltypes.GetStringSet(r.IncludePathPattern)
	state.ExcludePathPattern = internaltypes.GetStringSet(r.ExcludePathPattern)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a JdbcBasedErrorLogPublisherResponse object into the model struct
func readJdbcBasedErrorLogPublisherResponse(ctx context.Context, r *client.JdbcBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jdbc-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = types.StringValue(r.Server)
	state.LogFieldMapping = types.StringValue(r.LogFieldMapping)
	state.LogTableName = types.StringValue(r.LogTableName)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a JdbcBasedAccessLogPublisherResponse object into the model struct
func readJdbcBasedAccessLogPublisherResponse(ctx context.Context, r *client.JdbcBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jdbc-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Server = types.StringValue(r.Server)
	state.LogFieldMapping = types.StringValue(r.LogFieldMapping)
	state.LogTableName = types.StringValue(r.LogTableName)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a CommonLogFileHttpOperationLogPublisherResponse object into the model struct
func readCommonLogFileHttpOperationLogPublisherResponse(ctx context.Context, r *client.CommonLogFileHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("common-log-file-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ConsoleJsonErrorLogPublisherResponse object into the model struct
func readConsoleJsonErrorLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), true)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogTextErrorLogPublisherResponse object into the model struct
func readSyslogTextErrorLogPublisherResponse(ctx context.Context, r *client.SyslogTextErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-text-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), true)
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), true)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogBasedAccessLogPublisherResponse object into the model struct
func readSyslogBasedAccessLogPublisherResponse(ctx context.Context, r *client.SyslogBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.SyslogFacility = types.StringValue(strconv.FormatInt(r.SyslogFacility, 10))
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedJsonAuditLogPublisherResponse object into the model struct
func readFileBasedJsonAuditLogPublisherResponse(ctx context.Context, r *client.FileBasedJsonAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherFileBasedJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), true)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedDebugLogPublisherResponse object into the model struct
func readFileBasedDebugLogPublisherResponse(ctx context.Context, r *client.FileBasedDebugLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-debug")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), true)
	state.DefaultDebugLevel = types.StringValue(r.DefaultDebugLevel.String())
	state.DefaultDebugCategory = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultDebugCategoryProp(r.DefaultDebugCategory))
	state.DefaultOmitMethodEntryArguments = internaltypes.BoolTypeOrNil(r.DefaultOmitMethodEntryArguments)
	state.DefaultOmitMethodReturnValue = internaltypes.BoolTypeOrNil(r.DefaultOmitMethodReturnValue)
	state.DefaultIncludeThrowableCause = internaltypes.BoolTypeOrNil(r.DefaultIncludeThrowableCause)
	state.DefaultThrowableStackFrames = internaltypes.Int64TypeOrNil(r.DefaultThrowableStackFrames)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedErrorLogPublisherResponse object into the model struct
func readFileBasedErrorLogPublisherResponse(ctx context.Context, r *client.FileBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), true)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ThirdPartyErrorLogPublisherResponse object into the model struct
func readThirdPartyErrorLogPublisherResponse(ctx context.Context, r *client.ThirdPartyErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogTextAccessLogPublisherResponse object into the model struct
func readSyslogTextAccessLogPublisherResponse(ctx context.Context, r *client.SyslogTextAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-text-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), true)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a DetailedHttpOperationLogPublisherResponse object into the model struct
func readDetailedHttpOperationLogPublisherResponse(ctx context.Context, r *client.DetailedHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("detailed-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), true)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), true)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), true)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a JsonAccessLogPublisherResponse object into the model struct
func readJsonAccessLogPublisherResponse(ctx context.Context, r *client.JsonAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a DebugAccessLogPublisherResponse object into the model struct
func readDebugAccessLogPublisherResponse(ctx context.Context, r *client.DebugAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("debug-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ObscureSensitiveContent = internaltypes.BoolTypeOrNil(r.ObscureSensitiveContent)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.DebugACIEnabled = internaltypes.BoolTypeOrNil(r.DebugACIEnabled)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogJsonHttpOperationLogPublisherResponse object into the model struct
func readSyslogJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.SyslogJsonHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), true)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), true)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), true)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ThirdPartyAccessLogPublisherResponse object into the model struct
func readThirdPartyAccessLogPublisherResponse(ctx context.Context, r *client.ThirdPartyAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedAuditLogPublisherResponse object into the model struct
func readFileBasedAuditLogPublisherResponse(ctx context.Context, r *client.FileBasedAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherFileBasedAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), true)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), true)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a JsonErrorLogPublisherResponse object into the model struct
func readJsonErrorLogPublisherResponse(ctx context.Context, r *client.JsonErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a GroovyScriptedFileBasedAccessLogPublisherResponse object into the model struct
func readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedFileBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-file-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a GroovyScriptedFileBasedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedFileBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-file-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogJsonAccessLogPublisherResponse object into the model struct
func readSyslogJsonAccessLogPublisherResponse(ctx context.Context, r *client.SyslogJsonAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a GroovyScriptedAccessLogPublisherResponse object into the model struct
func readGroovyScriptedAccessLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ThirdPartyFileBasedErrorLogPublisherResponse object into the model struct
func readThirdPartyFileBasedErrorLogPublisherResponse(ctx context.Context, r *client.ThirdPartyFileBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-file-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ConsoleJsonAuditLogPublisherResponse object into the model struct
func readConsoleJsonAuditLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), true)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherConsoleJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), true)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ConsoleJsonHttpOperationLogPublisherResponse object into the model struct
func readConsoleJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), true)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), true)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), true)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), true)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a ConsoleJsonAccessLogPublisherResponse object into the model struct
func readConsoleJsonAccessLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), true)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedAccessLogPublisherResponse object into the model struct
func readFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.FileBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), true)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.IncludeRequestDetailsInSearchEntryMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchEntryMessages)
	state.IncludeRequestDetailsInSearchReferenceMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInSearchReferenceMessages)
	state.IncludeRequestDetailsInIntermediateResponseMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInIntermediateResponseMessages)
	state.IncludeResultCodeNames = internaltypes.BoolTypeOrNil(r.IncludeResultCodeNames)
	state.IncludeExtendedSearchRequestDetails = internaltypes.BoolTypeOrNil(r.IncludeExtendedSearchRequestDetails)
	state.IncludeAddAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeAddAttributeNames)
	state.IncludeModifyAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeModifyAttributeNames)
	state.IncludeSearchEntryAttributeNames = internaltypes.BoolTypeOrNil(r.IncludeSearchEntryAttributeNames)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeResponseControls = internaltypes.BoolTypeOrNil(r.IncludeResponseControls)
	state.IncludeReplicationChangeID = internaltypes.BoolTypeOrNil(r.IncludeReplicationChangeID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, internaltypes.IsEmptyString(expectedValues.LogFieldBehavior))
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, internaltypes.IsEmptyString(expectedValues.SearchEntryCriteria))
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, internaltypes.IsEmptyString(expectedValues.SearchReferenceCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a GroovyScriptedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedErrorLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a FileBasedJsonHttpOperationLogPublisherResponse object into the model struct
func readFileBasedJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.FileBasedJsonHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), true)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, true)
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, true)
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), true)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), true)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), true)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a SyslogJsonErrorLogPublisherResponse object into the model struct
func readSyslogJsonErrorLogPublisherResponse(ctx context.Context, r *client.SyslogJsonErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), true)
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Read a GroovyScriptedHttpOperationLogPublisherResponse object into the model struct
func readGroovyScriptedHttpOperationLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createLogPublisherOperations(plan logPublisherResourceModel, state logPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringOperationIfNecessary(&ops, plan.Server, state.Server, "server")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldMapping, state.LogFieldMapping, "log-field-mapping")
	operations.AddStringOperationIfNecessary(&ops, plan.LogTableName, state.LogTableName, "log-table-name")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputLocation, state.OutputLocation, "output-location")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.CompressionMechanism, state.CompressionMechanism, "compression-mechanism")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddBoolOperationIfNecessary(&ops, plan.SignLog, state.SignLog, "sign-log")
	operations.AddStringOperationIfNecessary(&ops, plan.TimestampPrecision, state.TimestampPrecision, "timestamp-precision")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptLog, state.EncryptLog, "encrypt-log")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddBoolOperationIfNecessary(&ops, plan.ObscureSensitiveContent, state.ObscureSensitiveContent, "obscure-sensitive-content")
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddBoolOperationIfNecessary(&ops, plan.DebugACIEnabled, state.DebugACIEnabled, "debug-aci-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultDebugLevel, state.DefaultDebugLevel, "default-debug-level")
	operations.AddStringOperationIfNecessary(&ops, plan.LogRequestHeaders, state.LogRequestHeaders, "log-request-headers")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedRequestHeaderName, state.SuppressedRequestHeaderName, "suppressed-request-header-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogResponseHeaders, state.LogResponseHeaders, "log-response-headers")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedResponseHeaderName, state.SuppressedResponseHeaderName, "suppressed-response-header-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestAuthorizationType, state.LogRequestAuthorizationType, "log-request-authorization-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestCookieNames, state.LogRequestCookieNames, "log-request-cookie-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResponseCookieNames, state.LogResponseCookieNames, "log-response-cookie-names")
	operations.AddStringOperationIfNecessary(&ops, plan.LogRequestParameters, state.LogRequestParameters, "log-request-parameters")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestProtocol, state.LogRequestProtocol, "log-request-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedRequestParameterName, state.SuppressedRequestParameterName, "suppressed-request-parameter-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRedirectURI, state.LogRedirectURI, "log-redirect-uri")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultDebugCategory, state.DefaultDebugCategory, "default-debug-category")
	operations.AddBoolOperationIfNecessary(&ops, plan.DefaultOmitMethodEntryArguments, state.DefaultOmitMethodEntryArguments, "default-omit-method-entry-arguments")
	operations.AddBoolOperationIfNecessary(&ops, plan.DefaultOmitMethodReturnValue, state.DefaultOmitMethodReturnValue, "default-omit-method-return-value")
	operations.AddBoolOperationIfNecessary(&ops, plan.DefaultIncludeThrowableCause, state.DefaultIncludeThrowableCause, "default-include-throwable-cause")
	operations.AddInt64OperationIfNecessary(&ops, plan.DefaultThrowableStackFrames, state.DefaultThrowableStackFrames, "default-throwable-stack-frames")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SyslogExternalServer, state.SyslogExternalServer, "syslog-external-server")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInResultMessages, state.IncludeRequestDetailsInResultMessages, "include-request-details-in-result-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogAssuranceCompleted, state.LogAssuranceCompleted, "log-assurance-completed")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DebugMessageType, state.DebugMessageType, "debug-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HttpMessageType, state.HttpMessageType, "http-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidatorMessageType, state.AccessTokenValidatorMessageType, "access-token-validator-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IdTokenValidatorMessageType, state.IdTokenValidatorMessageType, "id-token-validator-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScimMessageType, state.ScimMessageType, "scim-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ConsentMessageType, state.ConsentMessageType, "consent-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DirectoryRESTAPIMessageType, state.DirectoryRESTAPIMessageType, "directory-rest-api-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionMessageType, state.ExtensionMessageType, "extension-message-type")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.IncludePathPattern, state.IncludePathPattern, "include-path-pattern")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludePathPattern, state.ExcludePathPattern, "exclude-path-pattern")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddStringOperationIfNecessary(&ops, plan.BufferSize, state.BufferSize, "buffer-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.MinIncludedOperationProcessingTime, state.MinIncludedOperationProcessingTime, "min-included-operation-processing-time")
	operations.AddInt64OperationIfNecessary(&ops, plan.MinIncludedPhaseTimeNanos, state.MinIncludedPhaseTimeNanos, "min-included-phase-time-nanos")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeInterval, state.TimeInterval, "time-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInSearchEntryMessages, state.IncludeRequestDetailsInSearchEntryMessages, "include-request-details-in-search-entry-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInSearchReferenceMessages, state.IncludeRequestDetailsInSearchReferenceMessages, "include-request-details-in-search-reference-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInIntermediateResponseMessages, state.IncludeRequestDetailsInIntermediateResponseMessages, "include-request-details-in-intermediate-response-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResultCodeNames, state.IncludeResultCodeNames, "include-result-code-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeExtendedSearchRequestDetails, state.IncludeExtendedSearchRequestDetails, "include-extended-search-request-details")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeAddAttributeNames, state.IncludeAddAttributeNames, "include-add-attribute-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeModifyAttributeNames, state.IncludeModifyAttributeNames, "include-modify-attribute-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeSearchEntryAttributeNames, state.IncludeSearchEntryAttributeNames, "include-search-entry-attribute-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogConnects, state.LogConnects, "log-connects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogDisconnects, state.LogDisconnects, "log-disconnects")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxStringLength, state.MaxStringLength, "max-string-length")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerifyMessageStringsWhenPossible, state.GenerifyMessageStringsWhenPossible, "generify-message-strings-when-possible")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogFacility, state.SyslogFacility, "syslog-facility")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldBehavior, state.LogFieldBehavior, "log-field-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogClientCertificates, state.LogClientCertificates, "log-client-certificates")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequests, state.LogRequests, "log-requests")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResults, state.LogResults, "log-results")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchEntries, state.LogSearchEntries, "log-search-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchReferences, state.LogSearchReferences, "log-search-references")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogIntermediateResponses, state.LogIntermediateResponses, "log-intermediate-responses")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddBoolOperationIfNecessary(&ops, plan.CorrelateRequestsAndResults, state.CorrelateRequestsAndResults, "correlate-requests-and-results")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogSeverity, state.SyslogSeverity, "syslog-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultSeverity, state.DefaultSeverity, "default-severity")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OverrideSeverity, state.OverrideSeverity, "override-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchEntryCriteria, state.SearchEntryCriteria, "search-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchReferenceCriteria, state.SearchReferenceCriteria, "search-reference-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageHostName, state.SyslogMessageHostName, "syslog-message-host-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageApplicationName, state.SyslogMessageApplicationName, "syslog-message-application-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddBoolOperationIfNecessary(&ops, plan.WriteMultiLineMessages, state.WriteMultiLineMessages, "write-multi-line-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseReversibleForm, state.UseReversibleForm, "use-reversible-form")
	operations.AddStringOperationIfNecessary(&ops, plan.SoftDeleteEntryAuditBehavior, state.SoftDeleteEntryAuditBehavior, "soft-delete-entry-audit-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeOperationPurposeRequestControl, state.IncludeOperationPurposeRequestControl, "include-operation-purpose-request-control")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeIntermediateClientRequestControl, state.IncludeIntermediateClientRequestControl, "include-intermediate-client-request-control")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ObscureAttribute, state.ObscureAttribute, "obscure-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExcludeAttribute, state.ExcludeAttribute, "exclude-attribute")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressInternalOperations, state.SuppressInternalOperations, "suppress-internal-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterDN, state.IncludeRequesterDN, "include-requester-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterIPAddress, state.IncludeRequesterIPAddress, "include-requester-ip-address")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestControls, state.IncludeRequestControls, "include-request-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResponseControls, state.IncludeResponseControls, "include-response-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeReplicationChangeID, state.IncludeReplicationChangeID, "include-replication-change-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSecurityNegotiation, state.LogSecurityNegotiation, "log-security-negotiation")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressReplicationOperations, state.SuppressReplicationOperations, "suppress-replication-operations")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCriteria, state.ResultCriteria, "result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a syslog-json-audit log-publisher
func (r *logPublisherResource) CreateSyslogJsonAuditLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonAuditLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogJsonAuditLogPublisherSchemaUrn{client.ENUMSYSLOGJSONAUDITLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_AUDIT},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonAuditLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogJsonAuditLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogJsonAuditLogPublisherResponse(ctx, addResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-based-error log-publisher
func (r *logPublisherResource) CreateSyslogBasedErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddSyslogBasedErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogBasedErrorLogPublisherSchemaUrn{client.ENUMSYSLOGBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_BASED_ERROR},
		plan.Enabled.ValueBool())
	err := addOptionalSyslogBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogBasedErrorLogPublisherResponse(ctx, addResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party-file-based-access log-publisher
func (r *logPublisherResource) CreateThirdPartyFileBasedAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddThirdPartyFileBasedAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyFileBasedAccessLogPublisherSchemaUrn{client.ENUMTHIRDPARTYFILEBASEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERTHIRD_PARTY_FILE_BASED_ACCESS},
		plan.LogFile.ValueString(),
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyFileBasedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddThirdPartyFileBasedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readThirdPartyFileBasedAccessLogPublisherResponse(ctx, addResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a operation-timing-access log-publisher
func (r *logPublisherResource) CreateOperationTimingAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddOperationTimingAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumoperationTimingAccessLogPublisherSchemaUrn{client.ENUMOPERATIONTIMINGACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHEROPERATION_TIMING_ACCESS},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalOperationTimingAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddOperationTimingAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readOperationTimingAccessLogPublisherResponse(ctx, addResponse.OperationTimingAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party-http-operation log-publisher
func (r *logPublisherResource) CreateThirdPartyHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddThirdPartyHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyHttpOperationLogPublisherSchemaUrn{client.ENUMTHIRDPARTYHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERTHIRD_PARTY_HTTP_OPERATION},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddThirdPartyHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readThirdPartyHttpOperationLogPublisherResponse(ctx, addResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a admin-alert-access log-publisher
func (r *logPublisherResource) CreateAdminAlertAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddAdminAlertAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumadminAlertAccessLogPublisherSchemaUrn{client.ENUMADMINALERTACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERADMIN_ALERT_ACCESS},
		plan.Enabled.ValueBool())
	err := addOptionalAdminAlertAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddAdminAlertAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readAdminAlertAccessLogPublisherResponse(ctx, addResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-trace log-publisher
func (r *logPublisherResource) CreateFileBasedTraceLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedTraceLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedTraceLogPublisherSchemaUrn{client.ENUMFILEBASEDTRACELOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_TRACE},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedTraceLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedTraceLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedTraceLogPublisherResponse(ctx, addResponse.FileBasedTraceLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a jdbc-based-error log-publisher
func (r *logPublisherResource) CreateJdbcBasedErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddJdbcBasedErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumjdbcBasedErrorLogPublisherSchemaUrn{client.ENUMJDBCBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERJDBC_BASED_ERROR},
		plan.Server.ValueString(),
		plan.LogFieldMapping.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalJdbcBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddJdbcBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readJdbcBasedErrorLogPublisherResponse(ctx, addResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a jdbc-based-access log-publisher
func (r *logPublisherResource) CreateJdbcBasedAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddJdbcBasedAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumjdbcBasedAccessLogPublisherSchemaUrn{client.ENUMJDBCBASEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERJDBC_BASED_ACCESS},
		plan.Server.ValueString(),
		plan.LogFieldMapping.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalJdbcBasedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddJdbcBasedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readJdbcBasedAccessLogPublisherResponse(ctx, addResponse.JdbcBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a common-log-file-http-operation log-publisher
func (r *logPublisherResource) CreateCommonLogFileHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddCommonLogFileHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumcommonLogFileHttpOperationLogPublisherSchemaUrn{client.ENUMCOMMONLOGFILEHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERCOMMON_LOG_FILE_HTTP_OPERATION},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalCommonLogFileHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddCommonLogFileHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readCommonLogFileHttpOperationLogPublisherResponse(ctx, addResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-text-error log-publisher
func (r *logPublisherResource) CreateSyslogTextErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogTextErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogTextErrorLogPublisherSchemaUrn{client.ENUMSYSLOGTEXTERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_TEXT_ERROR},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogTextErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogTextErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogTextErrorLogPublisherResponse(ctx, addResponse.SyslogTextErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-based-access log-publisher
func (r *logPublisherResource) CreateSyslogBasedAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddSyslogBasedAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogBasedAccessLogPublisherSchemaUrn{client.ENUMSYSLOGBASEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_BASED_ACCESS},
		plan.Enabled.ValueBool())
	err := addOptionalSyslogBasedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogBasedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogBasedAccessLogPublisherResponse(ctx, addResponse.SyslogBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-json-audit log-publisher
func (r *logPublisherResource) CreateFileBasedJsonAuditLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedJsonAuditLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedJsonAuditLogPublisherSchemaUrn{client.ENUMFILEBASEDJSONAUDITLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_JSON_AUDIT},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedJsonAuditLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedJsonAuditLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedJsonAuditLogPublisherResponse(ctx, addResponse.FileBasedJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-debug log-publisher
func (r *logPublisherResource) CreateFileBasedDebugLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedDebugLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedDebugLogPublisherSchemaUrn{client.ENUMFILEBASEDDEBUGLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_DEBUG},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedDebugLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedDebugLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedDebugLogPublisherResponse(ctx, addResponse.FileBasedDebugLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-error log-publisher
func (r *logPublisherResource) CreateFileBasedErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedErrorLogPublisherSchemaUrn{client.ENUMFILEBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_ERROR},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedErrorLogPublisherResponse(ctx, addResponse.FileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party-error log-publisher
func (r *logPublisherResource) CreateThirdPartyErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddThirdPartyErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyErrorLogPublisherSchemaUrn{client.ENUMTHIRDPARTYERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERTHIRD_PARTY_ERROR},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddThirdPartyErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readThirdPartyErrorLogPublisherResponse(ctx, addResponse.ThirdPartyErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-text-access log-publisher
func (r *logPublisherResource) CreateSyslogTextAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogTextAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogTextAccessLogPublisherSchemaUrn{client.ENUMSYSLOGTEXTACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_TEXT_ACCESS},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogTextAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogTextAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogTextAccessLogPublisherResponse(ctx, addResponse.SyslogTextAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a detailed-http-operation log-publisher
func (r *logPublisherResource) CreateDetailedHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddDetailedHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumdetailedHttpOperationLogPublisherSchemaUrn{client.ENUMDETAILEDHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERDETAILED_HTTP_OPERATION},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalDetailedHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddDetailedHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readDetailedHttpOperationLogPublisherResponse(ctx, addResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a json-access log-publisher
func (r *logPublisherResource) CreateJsonAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddJsonAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumjsonAccessLogPublisherSchemaUrn{client.ENUMJSONACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERJSON_ACCESS},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalJsonAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddJsonAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readJsonAccessLogPublisherResponse(ctx, addResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a debug-access log-publisher
func (r *logPublisherResource) CreateDebugAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddDebugAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumdebugAccessLogPublisherSchemaUrn{client.ENUMDEBUGACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERDEBUG_ACCESS},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalDebugAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddDebugAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readDebugAccessLogPublisherResponse(ctx, addResponse.DebugAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-json-http-operation log-publisher
func (r *logPublisherResource) CreateSyslogJsonHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogJsonHttpOperationLogPublisherSchemaUrn{client.ENUMSYSLOGJSONHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_HTTP_OPERATION},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogJsonHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogJsonHttpOperationLogPublisherResponse(ctx, addResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party-access log-publisher
func (r *logPublisherResource) CreateThirdPartyAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddThirdPartyAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyAccessLogPublisherSchemaUrn{client.ENUMTHIRDPARTYACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERTHIRD_PARTY_ACCESS},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddThirdPartyAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readThirdPartyAccessLogPublisherResponse(ctx, addResponse.ThirdPartyAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-audit log-publisher
func (r *logPublisherResource) CreateFileBasedAuditLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedAuditLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedAuditLogPublisherSchemaUrn{client.ENUMFILEBASEDAUDITLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_AUDIT},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedAuditLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedAuditLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedAuditLogPublisherResponse(ctx, addResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a json-error log-publisher
func (r *logPublisherResource) CreateJsonErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddJsonErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumjsonErrorLogPublisherSchemaUrn{client.ENUMJSONERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERJSON_ERROR},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalJsonErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddJsonErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readJsonErrorLogPublisherResponse(ctx, addResponse.JsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted-file-based-access log-publisher
func (r *logPublisherResource) CreateGroovyScriptedFileBasedAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedFileBasedAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedFileBasedAccessLogPublisherSchemaUrn{client.ENUMGROOVYSCRIPTEDFILEBASEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERGROOVY_SCRIPTED_FILE_BASED_ACCESS},
		plan.ScriptClass.ValueString(),
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedFileBasedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddGroovyScriptedFileBasedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, addResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted-file-based-error log-publisher
func (r *logPublisherResource) CreateGroovyScriptedFileBasedErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedFileBasedErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedFileBasedErrorLogPublisherSchemaUrn{client.ENUMGROOVYSCRIPTEDFILEBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERGROOVY_SCRIPTED_FILE_BASED_ERROR},
		plan.ScriptClass.ValueString(),
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedFileBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddGroovyScriptedFileBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, addResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-json-access log-publisher
func (r *logPublisherResource) CreateSyslogJsonAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogJsonAccessLogPublisherSchemaUrn{client.ENUMSYSLOGJSONACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_ACCESS},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogJsonAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogJsonAccessLogPublisherResponse(ctx, addResponse.SyslogJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted-access log-publisher
func (r *logPublisherResource) CreateGroovyScriptedAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedAccessLogPublisherSchemaUrn{client.ENUMGROOVYSCRIPTEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERGROOVY_SCRIPTED_ACCESS},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddGroovyScriptedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readGroovyScriptedAccessLogPublisherResponse(ctx, addResponse.GroovyScriptedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party-file-based-error log-publisher
func (r *logPublisherResource) CreateThirdPartyFileBasedErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddThirdPartyFileBasedErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyFileBasedErrorLogPublisherSchemaUrn{client.ENUMTHIRDPARTYFILEBASEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERTHIRD_PARTY_FILE_BASED_ERROR},
		plan.LogFile.ValueString(),
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalThirdPartyFileBasedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddThirdPartyFileBasedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readThirdPartyFileBasedErrorLogPublisherResponse(ctx, addResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a console-json-audit log-publisher
func (r *logPublisherResource) CreateConsoleJsonAuditLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddConsoleJsonAuditLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumconsoleJsonAuditLogPublisherSchemaUrn{client.ENUMCONSOLEJSONAUDITLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERCONSOLE_JSON_AUDIT},
		plan.Enabled.ValueBool())
	err := addOptionalConsoleJsonAuditLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddConsoleJsonAuditLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readConsoleJsonAuditLogPublisherResponse(ctx, addResponse.ConsoleJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a console-json-http-operation log-publisher
func (r *logPublisherResource) CreateConsoleJsonHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddConsoleJsonHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumconsoleJsonHttpOperationLogPublisherSchemaUrn{client.ENUMCONSOLEJSONHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERCONSOLE_JSON_HTTP_OPERATION},
		plan.Enabled.ValueBool())
	err := addOptionalConsoleJsonHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddConsoleJsonHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readConsoleJsonHttpOperationLogPublisherResponse(ctx, addResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-access log-publisher
func (r *logPublisherResource) CreateFileBasedAccessLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedAccessLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedAccessLogPublisherSchemaUrn{client.ENUMFILEBASEDACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_ACCESS},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedAccessLogPublisherResponse(ctx, addResponse.FileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted-error log-publisher
func (r *logPublisherResource) CreateGroovyScriptedErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedErrorLogPublisherSchemaUrn{client.ENUMGROOVYSCRIPTEDERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERGROOVY_SCRIPTED_ERROR},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddGroovyScriptedErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readGroovyScriptedErrorLogPublisherResponse(ctx, addResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based-json-http-operation log-publisher
func (r *logPublisherResource) CreateFileBasedJsonHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddFileBasedJsonHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedJsonHttpOperationLogPublisherSchemaUrn{client.ENUMFILEBASEDJSONHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERFILE_BASED_JSON_HTTP_OPERATION},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalFileBasedJsonHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddFileBasedJsonHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readFileBasedJsonHttpOperationLogPublisherResponse(ctx, addResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a syslog-json-error log-publisher
func (r *logPublisherResource) CreateSyslogJsonErrorLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonErrorLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumsyslogJsonErrorLogPublisherSchemaUrn{client.ENUMSYSLOGJSONERRORLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_ERROR},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonErrorLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddSyslogJsonErrorLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readSyslogJsonErrorLogPublisherResponse(ctx, addResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a groovy-scripted-http-operation log-publisher
func (r *logPublisherResource) CreateGroovyScriptedHttpOperationLogPublisher(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan logPublisherResourceModel) (*logPublisherResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedHttpOperationLogPublisherRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedHttpOperationLogPublisherSchemaUrn{client.ENUMGROOVYSCRIPTEDHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERGROOVY_SCRIPTED_HTTP_OPERATION},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalGroovyScriptedHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Log Publisher", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogPublisherApi.AddLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogPublisherRequest(
		client.AddGroovyScriptedHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Log Publisher", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state logPublisherResourceModel
	readGroovyScriptedHttpOperationLogPublisherResponse(ctx, addResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *logPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *logPublisherResourceModel
	var err error
	if plan.Type.ValueString() == "syslog-json-audit" {
		state, err = r.CreateSyslogJsonAuditLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-based-error" {
		state, err = r.CreateSyslogBasedErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party-file-based-access" {
		state, err = r.CreateThirdPartyFileBasedAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "operation-timing-access" {
		state, err = r.CreateOperationTimingAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party-http-operation" {
		state, err = r.CreateThirdPartyHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "admin-alert-access" {
		state, err = r.CreateAdminAlertAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-trace" {
		state, err = r.CreateFileBasedTraceLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "jdbc-based-error" {
		state, err = r.CreateJdbcBasedErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "jdbc-based-access" {
		state, err = r.CreateJdbcBasedAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "common-log-file-http-operation" {
		state, err = r.CreateCommonLogFileHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-text-error" {
		state, err = r.CreateSyslogTextErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-based-access" {
		state, err = r.CreateSyslogBasedAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-json-audit" {
		state, err = r.CreateFileBasedJsonAuditLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-debug" {
		state, err = r.CreateFileBasedDebugLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-error" {
		state, err = r.CreateFileBasedErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party-error" {
		state, err = r.CreateThirdPartyErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-text-access" {
		state, err = r.CreateSyslogTextAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "detailed-http-operation" {
		state, err = r.CreateDetailedHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "json-access" {
		state, err = r.CreateJsonAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "debug-access" {
		state, err = r.CreateDebugAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-json-http-operation" {
		state, err = r.CreateSyslogJsonHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party-access" {
		state, err = r.CreateThirdPartyAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-audit" {
		state, err = r.CreateFileBasedAuditLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "json-error" {
		state, err = r.CreateJsonErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted-file-based-access" {
		state, err = r.CreateGroovyScriptedFileBasedAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted-file-based-error" {
		state, err = r.CreateGroovyScriptedFileBasedErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-json-access" {
		state, err = r.CreateSyslogJsonAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted-access" {
		state, err = r.CreateGroovyScriptedAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party-file-based-error" {
		state, err = r.CreateThirdPartyFileBasedErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "console-json-audit" {
		state, err = r.CreateConsoleJsonAuditLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "console-json-http-operation" {
		state, err = r.CreateConsoleJsonHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-access" {
		state, err = r.CreateFileBasedAccessLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted-error" {
		state, err = r.CreateGroovyScriptedErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based-json-http-operation" {
		state, err = r.CreateFileBasedJsonHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "syslog-json-error" {
		state, err = r.CreateSyslogJsonErrorLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "groovy-scripted-http-operation" {
		state, err = r.CreateGroovyScriptedHttpOperationLogPublisher(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state logPublisherResourceModel
	if readResponse.SyslogJsonAuditLogPublisherResponse != nil {
		readSyslogJsonAuditLogPublisherResponse(ctx, readResponse.SyslogJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogBasedErrorLogPublisherResponse != nil {
		readSyslogBasedErrorLogPublisherResponse(ctx, readResponse.SyslogBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyFileBasedAccessLogPublisherResponse != nil {
		readThirdPartyFileBasedAccessLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OperationTimingAccessLogPublisherResponse != nil {
		readOperationTimingAccessLogPublisherResponse(ctx, readResponse.OperationTimingAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyHttpOperationLogPublisherResponse != nil {
		readThirdPartyHttpOperationLogPublisherResponse(ctx, readResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AdminAlertAccessLogPublisherResponse != nil {
		readAdminAlertAccessLogPublisherResponse(ctx, readResponse.AdminAlertAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedTraceLogPublisherResponse != nil {
		readFileBasedTraceLogPublisherResponse(ctx, readResponse.FileBasedTraceLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcBasedErrorLogPublisherResponse != nil {
		readJdbcBasedErrorLogPublisherResponse(ctx, readResponse.JdbcBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcBasedAccessLogPublisherResponse != nil {
		readJdbcBasedAccessLogPublisherResponse(ctx, readResponse.JdbcBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CommonLogFileHttpOperationLogPublisherResponse != nil {
		readCommonLogFileHttpOperationLogPublisherResponse(ctx, readResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonErrorLogPublisherResponse != nil {
		readConsoleJsonErrorLogPublisherResponse(ctx, readResponse.ConsoleJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogTextErrorLogPublisherResponse != nil {
		readSyslogTextErrorLogPublisherResponse(ctx, readResponse.SyslogTextErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogBasedAccessLogPublisherResponse != nil {
		readSyslogBasedAccessLogPublisherResponse(ctx, readResponse.SyslogBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedJsonAuditLogPublisherResponse != nil {
		readFileBasedJsonAuditLogPublisherResponse(ctx, readResponse.FileBasedJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedDebugLogPublisherResponse != nil {
		readFileBasedDebugLogPublisherResponse(ctx, readResponse.FileBasedDebugLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedErrorLogPublisherResponse != nil {
		readFileBasedErrorLogPublisherResponse(ctx, readResponse.FileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyErrorLogPublisherResponse != nil {
		readThirdPartyErrorLogPublisherResponse(ctx, readResponse.ThirdPartyErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogTextAccessLogPublisherResponse != nil {
		readSyslogTextAccessLogPublisherResponse(ctx, readResponse.SyslogTextAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DetailedHttpOperationLogPublisherResponse != nil {
		readDetailedHttpOperationLogPublisherResponse(ctx, readResponse.DetailedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JsonAccessLogPublisherResponse != nil {
		readJsonAccessLogPublisherResponse(ctx, readResponse.JsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DebugAccessLogPublisherResponse != nil {
		readDebugAccessLogPublisherResponse(ctx, readResponse.DebugAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonHttpOperationLogPublisherResponse != nil {
		readSyslogJsonHttpOperationLogPublisherResponse(ctx, readResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccessLogPublisherResponse != nil {
		readThirdPartyAccessLogPublisherResponse(ctx, readResponse.ThirdPartyAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedAuditLogPublisherResponse != nil {
		readFileBasedAuditLogPublisherResponse(ctx, readResponse.FileBasedAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JsonErrorLogPublisherResponse != nil {
		readJsonErrorLogPublisherResponse(ctx, readResponse.JsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse != nil {
		readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse != nil {
		readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonAccessLogPublisherResponse != nil {
		readSyslogJsonAccessLogPublisherResponse(ctx, readResponse.SyslogJsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAccessLogPublisherResponse != nil {
		readGroovyScriptedAccessLogPublisherResponse(ctx, readResponse.GroovyScriptedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyFileBasedErrorLogPublisherResponse != nil {
		readThirdPartyFileBasedErrorLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonAuditLogPublisherResponse != nil {
		readConsoleJsonAuditLogPublisherResponse(ctx, readResponse.ConsoleJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonHttpOperationLogPublisherResponse != nil {
		readConsoleJsonHttpOperationLogPublisherResponse(ctx, readResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonAccessLogPublisherResponse != nil {
		readConsoleJsonAccessLogPublisherResponse(ctx, readResponse.ConsoleJsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedAccessLogPublisherResponse != nil {
		readFileBasedAccessLogPublisherResponse(ctx, readResponse.FileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedErrorLogPublisherResponse != nil {
		readGroovyScriptedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedJsonHttpOperationLogPublisherResponse != nil {
		readFileBasedJsonHttpOperationLogPublisherResponse(ctx, readResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonErrorLogPublisherResponse != nil {
		readSyslogJsonErrorLogPublisherResponse(ctx, readResponse.SyslogJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedHttpOperationLogPublisherResponse != nil {
		readGroovyScriptedHttpOperationLogPublisherResponse(ctx, readResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SyslogJsonAuditLogPublisherResponse != nil {
			readSyslogJsonAuditLogPublisherResponse(ctx, updateResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogBasedErrorLogPublisherResponse != nil {
			readSyslogBasedErrorLogPublisherResponse(ctx, updateResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse != nil {
			readThirdPartyFileBasedAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OperationTimingAccessLogPublisherResponse != nil {
			readOperationTimingAccessLogPublisherResponse(ctx, updateResponse.OperationTimingAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyHttpOperationLogPublisherResponse != nil {
			readThirdPartyHttpOperationLogPublisherResponse(ctx, updateResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AdminAlertAccessLogPublisherResponse != nil {
			readAdminAlertAccessLogPublisherResponse(ctx, updateResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedTraceLogPublisherResponse != nil {
			readFileBasedTraceLogPublisherResponse(ctx, updateResponse.FileBasedTraceLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JdbcBasedErrorLogPublisherResponse != nil {
			readJdbcBasedErrorLogPublisherResponse(ctx, updateResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JdbcBasedAccessLogPublisherResponse != nil {
			readJdbcBasedAccessLogPublisherResponse(ctx, updateResponse.JdbcBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CommonLogFileHttpOperationLogPublisherResponse != nil {
			readCommonLogFileHttpOperationLogPublisherResponse(ctx, updateResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonErrorLogPublisherResponse != nil {
			readConsoleJsonErrorLogPublisherResponse(ctx, updateResponse.ConsoleJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogTextErrorLogPublisherResponse != nil {
			readSyslogTextErrorLogPublisherResponse(ctx, updateResponse.SyslogTextErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogBasedAccessLogPublisherResponse != nil {
			readSyslogBasedAccessLogPublisherResponse(ctx, updateResponse.SyslogBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedJsonAuditLogPublisherResponse != nil {
			readFileBasedJsonAuditLogPublisherResponse(ctx, updateResponse.FileBasedJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedDebugLogPublisherResponse != nil {
			readFileBasedDebugLogPublisherResponse(ctx, updateResponse.FileBasedDebugLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedErrorLogPublisherResponse != nil {
			readFileBasedErrorLogPublisherResponse(ctx, updateResponse.FileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyErrorLogPublisherResponse != nil {
			readThirdPartyErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogTextAccessLogPublisherResponse != nil {
			readSyslogTextAccessLogPublisherResponse(ctx, updateResponse.SyslogTextAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DetailedHttpOperationLogPublisherResponse != nil {
			readDetailedHttpOperationLogPublisherResponse(ctx, updateResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JsonAccessLogPublisherResponse != nil {
			readJsonAccessLogPublisherResponse(ctx, updateResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DebugAccessLogPublisherResponse != nil {
			readDebugAccessLogPublisherResponse(ctx, updateResponse.DebugAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogJsonHttpOperationLogPublisherResponse != nil {
			readSyslogJsonHttpOperationLogPublisherResponse(ctx, updateResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAccessLogPublisherResponse != nil {
			readThirdPartyAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedAuditLogPublisherResponse != nil {
			readFileBasedAuditLogPublisherResponse(ctx, updateResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JsonErrorLogPublisherResponse != nil {
			readJsonErrorLogPublisherResponse(ctx, updateResponse.JsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedFileBasedAccessLogPublisherResponse != nil {
			readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedFileBasedErrorLogPublisherResponse != nil {
			readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogJsonAccessLogPublisherResponse != nil {
			readSyslogJsonAccessLogPublisherResponse(ctx, updateResponse.SyslogJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedAccessLogPublisherResponse != nil {
			readGroovyScriptedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyFileBasedErrorLogPublisherResponse != nil {
			readThirdPartyFileBasedErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonAuditLogPublisherResponse != nil {
			readConsoleJsonAuditLogPublisherResponse(ctx, updateResponse.ConsoleJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonHttpOperationLogPublisherResponse != nil {
			readConsoleJsonHttpOperationLogPublisherResponse(ctx, updateResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonAccessLogPublisherResponse != nil {
			readConsoleJsonAccessLogPublisherResponse(ctx, updateResponse.ConsoleJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedAccessLogPublisherResponse != nil {
			readFileBasedAccessLogPublisherResponse(ctx, updateResponse.FileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedErrorLogPublisherResponse != nil {
			readGroovyScriptedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedJsonHttpOperationLogPublisherResponse != nil {
			readFileBasedJsonHttpOperationLogPublisherResponse(ctx, updateResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogJsonErrorLogPublisherResponse != nil {
			readSyslogJsonErrorLogPublisherResponse(ctx, updateResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedHttpOperationLogPublisherResponse != nil {
			readGroovyScriptedHttpOperationLogPublisherResponse(ctx, updateResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *logPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state logPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Log Publisher", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Publisher", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SyslogJsonAuditLogPublisherResponse != nil {
		readSyslogJsonAuditLogPublisherResponse(ctx, readResponse.SyslogJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogBasedErrorLogPublisherResponse != nil {
		readSyslogBasedErrorLogPublisherResponse(ctx, readResponse.SyslogBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyFileBasedAccessLogPublisherResponse != nil {
		readThirdPartyFileBasedAccessLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.OperationTimingAccessLogPublisherResponse != nil {
		readOperationTimingAccessLogPublisherResponse(ctx, readResponse.OperationTimingAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyHttpOperationLogPublisherResponse != nil {
		readThirdPartyHttpOperationLogPublisherResponse(ctx, readResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AdminAlertAccessLogPublisherResponse != nil {
		readAdminAlertAccessLogPublisherResponse(ctx, readResponse.AdminAlertAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedTraceLogPublisherResponse != nil {
		readFileBasedTraceLogPublisherResponse(ctx, readResponse.FileBasedTraceLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcBasedErrorLogPublisherResponse != nil {
		readJdbcBasedErrorLogPublisherResponse(ctx, readResponse.JdbcBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcBasedAccessLogPublisherResponse != nil {
		readJdbcBasedAccessLogPublisherResponse(ctx, readResponse.JdbcBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CommonLogFileHttpOperationLogPublisherResponse != nil {
		readCommonLogFileHttpOperationLogPublisherResponse(ctx, readResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonErrorLogPublisherResponse != nil {
		readConsoleJsonErrorLogPublisherResponse(ctx, readResponse.ConsoleJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogTextErrorLogPublisherResponse != nil {
		readSyslogTextErrorLogPublisherResponse(ctx, readResponse.SyslogTextErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogBasedAccessLogPublisherResponse != nil {
		readSyslogBasedAccessLogPublisherResponse(ctx, readResponse.SyslogBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedJsonAuditLogPublisherResponse != nil {
		readFileBasedJsonAuditLogPublisherResponse(ctx, readResponse.FileBasedJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedDebugLogPublisherResponse != nil {
		readFileBasedDebugLogPublisherResponse(ctx, readResponse.FileBasedDebugLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedErrorLogPublisherResponse != nil {
		readFileBasedErrorLogPublisherResponse(ctx, readResponse.FileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyErrorLogPublisherResponse != nil {
		readThirdPartyErrorLogPublisherResponse(ctx, readResponse.ThirdPartyErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogTextAccessLogPublisherResponse != nil {
		readSyslogTextAccessLogPublisherResponse(ctx, readResponse.SyslogTextAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DetailedHttpOperationLogPublisherResponse != nil {
		readDetailedHttpOperationLogPublisherResponse(ctx, readResponse.DetailedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JsonAccessLogPublisherResponse != nil {
		readJsonAccessLogPublisherResponse(ctx, readResponse.JsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.DebugAccessLogPublisherResponse != nil {
		readDebugAccessLogPublisherResponse(ctx, readResponse.DebugAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonHttpOperationLogPublisherResponse != nil {
		readSyslogJsonHttpOperationLogPublisherResponse(ctx, readResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccessLogPublisherResponse != nil {
		readThirdPartyAccessLogPublisherResponse(ctx, readResponse.ThirdPartyAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedAuditLogPublisherResponse != nil {
		readFileBasedAuditLogPublisherResponse(ctx, readResponse.FileBasedAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JsonErrorLogPublisherResponse != nil {
		readJsonErrorLogPublisherResponse(ctx, readResponse.JsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse != nil {
		readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse != nil {
		readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonAccessLogPublisherResponse != nil {
		readSyslogJsonAccessLogPublisherResponse(ctx, readResponse.SyslogJsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAccessLogPublisherResponse != nil {
		readGroovyScriptedAccessLogPublisherResponse(ctx, readResponse.GroovyScriptedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyFileBasedErrorLogPublisherResponse != nil {
		readThirdPartyFileBasedErrorLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonAuditLogPublisherResponse != nil {
		readConsoleJsonAuditLogPublisherResponse(ctx, readResponse.ConsoleJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonHttpOperationLogPublisherResponse != nil {
		readConsoleJsonHttpOperationLogPublisherResponse(ctx, readResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonAccessLogPublisherResponse != nil {
		readConsoleJsonAccessLogPublisherResponse(ctx, readResponse.ConsoleJsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedAccessLogPublisherResponse != nil {
		readFileBasedAccessLogPublisherResponse(ctx, readResponse.FileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedErrorLogPublisherResponse != nil {
		readGroovyScriptedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedJsonHttpOperationLogPublisherResponse != nil {
		readFileBasedJsonHttpOperationLogPublisherResponse(ctx, readResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonErrorLogPublisherResponse != nil {
		readSyslogJsonErrorLogPublisherResponse(ctx, readResponse.SyslogJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedHttpOperationLogPublisherResponse != nil {
		readGroovyScriptedHttpOperationLogPublisherResponse(ctx, readResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *logPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan logPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state logPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SyslogJsonAuditLogPublisherResponse != nil {
			readSyslogJsonAuditLogPublisherResponse(ctx, updateResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogBasedErrorLogPublisherResponse != nil {
			readSyslogBasedErrorLogPublisherResponse(ctx, updateResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse != nil {
			readThirdPartyFileBasedAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.OperationTimingAccessLogPublisherResponse != nil {
			readOperationTimingAccessLogPublisherResponse(ctx, updateResponse.OperationTimingAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyHttpOperationLogPublisherResponse != nil {
			readThirdPartyHttpOperationLogPublisherResponse(ctx, updateResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AdminAlertAccessLogPublisherResponse != nil {
			readAdminAlertAccessLogPublisherResponse(ctx, updateResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedTraceLogPublisherResponse != nil {
			readFileBasedTraceLogPublisherResponse(ctx, updateResponse.FileBasedTraceLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JdbcBasedErrorLogPublisherResponse != nil {
			readJdbcBasedErrorLogPublisherResponse(ctx, updateResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JdbcBasedAccessLogPublisherResponse != nil {
			readJdbcBasedAccessLogPublisherResponse(ctx, updateResponse.JdbcBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CommonLogFileHttpOperationLogPublisherResponse != nil {
			readCommonLogFileHttpOperationLogPublisherResponse(ctx, updateResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonErrorLogPublisherResponse != nil {
			readConsoleJsonErrorLogPublisherResponse(ctx, updateResponse.ConsoleJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogTextErrorLogPublisherResponse != nil {
			readSyslogTextErrorLogPublisherResponse(ctx, updateResponse.SyslogTextErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogBasedAccessLogPublisherResponse != nil {
			readSyslogBasedAccessLogPublisherResponse(ctx, updateResponse.SyslogBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedJsonAuditLogPublisherResponse != nil {
			readFileBasedJsonAuditLogPublisherResponse(ctx, updateResponse.FileBasedJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedDebugLogPublisherResponse != nil {
			readFileBasedDebugLogPublisherResponse(ctx, updateResponse.FileBasedDebugLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedErrorLogPublisherResponse != nil {
			readFileBasedErrorLogPublisherResponse(ctx, updateResponse.FileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyErrorLogPublisherResponse != nil {
			readThirdPartyErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogTextAccessLogPublisherResponse != nil {
			readSyslogTextAccessLogPublisherResponse(ctx, updateResponse.SyslogTextAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DetailedHttpOperationLogPublisherResponse != nil {
			readDetailedHttpOperationLogPublisherResponse(ctx, updateResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JsonAccessLogPublisherResponse != nil {
			readJsonAccessLogPublisherResponse(ctx, updateResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.DebugAccessLogPublisherResponse != nil {
			readDebugAccessLogPublisherResponse(ctx, updateResponse.DebugAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogJsonHttpOperationLogPublisherResponse != nil {
			readSyslogJsonHttpOperationLogPublisherResponse(ctx, updateResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAccessLogPublisherResponse != nil {
			readThirdPartyAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedAuditLogPublisherResponse != nil {
			readFileBasedAuditLogPublisherResponse(ctx, updateResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JsonErrorLogPublisherResponse != nil {
			readJsonErrorLogPublisherResponse(ctx, updateResponse.JsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedFileBasedAccessLogPublisherResponse != nil {
			readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedFileBasedErrorLogPublisherResponse != nil {
			readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogJsonAccessLogPublisherResponse != nil {
			readSyslogJsonAccessLogPublisherResponse(ctx, updateResponse.SyslogJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedAccessLogPublisherResponse != nil {
			readGroovyScriptedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyFileBasedErrorLogPublisherResponse != nil {
			readThirdPartyFileBasedErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonAuditLogPublisherResponse != nil {
			readConsoleJsonAuditLogPublisherResponse(ctx, updateResponse.ConsoleJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonHttpOperationLogPublisherResponse != nil {
			readConsoleJsonHttpOperationLogPublisherResponse(ctx, updateResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConsoleJsonAccessLogPublisherResponse != nil {
			readConsoleJsonAccessLogPublisherResponse(ctx, updateResponse.ConsoleJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedAccessLogPublisherResponse != nil {
			readFileBasedAccessLogPublisherResponse(ctx, updateResponse.FileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedErrorLogPublisherResponse != nil {
			readGroovyScriptedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedJsonHttpOperationLogPublisherResponse != nil {
			readFileBasedJsonHttpOperationLogPublisherResponse(ctx, updateResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SyslogJsonErrorLogPublisherResponse != nil {
			readSyslogJsonErrorLogPublisherResponse(ctx, updateResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GroovyScriptedHttpOperationLogPublisherResponse != nil {
			readGroovyScriptedHttpOperationLogPublisherResponse(ctx, updateResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *logPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state logPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Log Publisher", err, httpResp)
		return
	}
}

func (r *logPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogPublisher(ctx, req, resp)
}

func (r *defaultLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLogPublisher(ctx, req, resp)
}

func importLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
