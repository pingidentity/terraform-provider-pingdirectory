package logpublisher

import (
	"context"
	"strconv"

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
	_ datasource.DataSource              = &logPublisherDataSource{}
	_ datasource.DataSourceWithConfigure = &logPublisherDataSource{}
)

// Create a Log Publisher data source
func NewLogPublisherDataSource() datasource.DataSource {
	return &logPublisherDataSource{}
}

// logPublisherDataSource is the datasource implementation.
type logPublisherDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logPublisherDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_publisher"
}

// Configure adds the provider configured client to the data source.
func (r *logPublisherDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logPublisherDataSourceModel struct {
	Id                                                  types.String `tfsdk:"id"`
	Name                                                types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *logPublisherDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Publisher resource. Options are ['syslog-json-audit', 'syslog-based-error', 'third-party-file-based-access', 'operation-timing-access', 'third-party-http-operation', 'admin-alert-access', 'file-based-trace', 'jdbc-based-error', 'jdbc-based-access', 'common-log-file-http-operation', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'debug-access', 'syslog-json-http-operation', 'third-party-access', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'groovy-scripted-file-based-error', 'syslog-json-access', 'groovy-scripted-access', 'third-party-file-based-error', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'groovy-scripted-error', 'file-based-json-http-operation', 'syslog-json-error', 'groovy-scripted-http-operation']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `groovy-scripted-file-based-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Access Log Publisher. When the `type` attribute is set to `groovy-scripted-file-based-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Error Log Publisher. When the `type` attribute is set to `groovy-scripted-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Access Log Publisher. When the `type` attribute is set to `groovy-scripted-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Error Log Publisher. When the `type` attribute is set to `groovy-scripted-http-operation`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Operation Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `groovy-scripted-file-based-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Access Log Publisher.\n  - `groovy-scripted-file-based-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Error Log Publisher.\n  - `groovy-scripted-access`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Access Log Publisher.\n  - `groovy-scripted-error`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Error Log Publisher.\n  - `groovy-scripted-http-operation`: The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted HTTP Operation Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"server": schema.StringAttribute{
				Description: "The JDBC-based Database Server to use for a connection.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_field_mapping": schema.StringAttribute{
				Description: "The log field mapping associates loggable fields to database column names. The table name is not part of this mapping.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_table_name": schema.StringAttribute{
				Description: "The table name to log entries to the database server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"output_location": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `console-json-error`: Specifies the output stream to which JSON-formatted error log messages should be written. When the `type` attribute is set to `console-json-audit`: Specifies the output stream to which JSON-formatted audit log messages should be written. When the `type` attribute is set to `console-json-http-operation`: Specifies the output stream to which JSON-formatted HTTP operation log messages should be written. When the `type` attribute is set to `console-json-access`: Specifies the output stream to which JSON-formatted access log messages should be written.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `console-json-error`: Specifies the output stream to which JSON-formatted error log messages should be written.\n  - `console-json-audit`: Specifies the output stream to which JSON-formatted audit log messages should be written.\n  - `console-json-http-operation`: Specifies the output stream to which JSON-formatted HTTP operation log messages should be written.\n  - `console-json-access`: Specifies the output stream to which JSON-formatted access log messages should be written.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_file": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The file name to use for the log files generated by the Third Party File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `operation-timing-access`: The file name to use for the log files generated by the Operation Timing Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-trace`: The file name to use for the log files generated by the File Based Trace Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `common-log-file-http-operation`: The file name to use for the log files generated by the Common Log File HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-json-audit`: The file name to use for the log files generated by the File Based JSON Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-debug`: The file name to use for the log files generated by the File Based Debug Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-error`: The file name to use for the log files generated by the File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `detailed-http-operation`: The file name to use for the log files generated by the Detailed HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `json-access`: The file name to use for the log files generated by the JSON Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `debug-access`: The file name to use for the log files generated by the Debug Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-audit`: The file name to use for the log files generated by the File Based Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `json-error`: The file name to use for the log files generated by the JSON Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `groovy-scripted-file-based-access`: The file name to use for the log files generated by the Scripted File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `groovy-scripted-file-based-error`: The file name to use for the log files generated by the Scripted File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `third-party-file-based-error`: The file name to use for the log files generated by the Third Party File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-access`: The file name to use for the log files generated by the File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path. When the `type` attribute is set to `file-based-json-http-operation`: The file name to use for the log files generated by the File Based JSON HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The file name to use for the log files generated by the Third Party File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `operation-timing-access`: The file name to use for the log files generated by the Operation Timing Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-trace`: The file name to use for the log files generated by the File Based Trace Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `common-log-file-http-operation`: The file name to use for the log files generated by the Common Log File HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-json-audit`: The file name to use for the log files generated by the File Based JSON Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-debug`: The file name to use for the log files generated by the File Based Debug Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-error`: The file name to use for the log files generated by the File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `detailed-http-operation`: The file name to use for the log files generated by the Detailed HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `json-access`: The file name to use for the log files generated by the JSON Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `debug-access`: The file name to use for the log files generated by the Debug Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-audit`: The file name to use for the log files generated by the File Based Audit Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `json-error`: The file name to use for the log files generated by the JSON Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `groovy-scripted-file-based-access`: The file name to use for the log files generated by the Scripted File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `groovy-scripted-file-based-error`: The file name to use for the log files generated by the Scripted File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `third-party-file-based-error`: The file name to use for the log files generated by the Third Party File Based Error Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-access`: The file name to use for the log files generated by the File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.\n  - `file-based-json-http-operation`: The file name to use for the log files generated by the File Based JSON HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The UNIX permissions of the log files created by this Third Party File Based Access Log Publisher. When the `type` attribute is set to `operation-timing-access`: The UNIX permissions of the log files created by this Operation Timing Access Log Publisher. When the `type` attribute is set to `file-based-trace`: The UNIX permissions of the log files created by this File Based Trace Log Publisher. When the `type` attribute is set to `common-log-file-http-operation`: The UNIX permissions of the log files created by this Common Log File HTTP Operation Log Publisher. When the `type` attribute is set to `file-based-json-audit`: The UNIX permissions of the log files created by this File Based JSON Audit Log Publisher. When the `type` attribute is set to `file-based-debug`: The UNIX permissions of the log files created by this File Based Debug Log Publisher. When the `type` attribute is set to `file-based-error`: The UNIX permissions of the log files created by this File Based Error Log Publisher. When the `type` attribute is set to `detailed-http-operation`: The UNIX permissions of the log files created by this Detailed HTTP Operation Log Publisher. When the `type` attribute is set to `json-access`: The UNIX permissions of the log files created by this JSON Access Log Publisher. When the `type` attribute is set to `debug-access`: The UNIX permissions of the log files created by this Debug Access Log Publisher. When the `type` attribute is set to `file-based-audit`: The UNIX permissions of the log files created by this File Based Audit Log Publisher. When the `type` attribute is set to `json-error`: The UNIX permissions of the log files created by this JSON Error Log Publisher. When the `type` attribute is set to `groovy-scripted-file-based-access`: The UNIX permissions of the log files created by this Scripted File Based Access Log Publisher. When the `type` attribute is set to `groovy-scripted-file-based-error`: The UNIX permissions of the log files created by this Scripted File Based Error Log Publisher. When the `type` attribute is set to `third-party-file-based-error`: The UNIX permissions of the log files created by this Third Party File Based Error Log Publisher. When the `type` attribute is set to `file-based-access`: The UNIX permissions of the log files created by this File Based Access Log Publisher. When the `type` attribute is set to `file-based-json-http-operation`: The UNIX permissions of the log files created by this File Based JSON HTTP Operation Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The UNIX permissions of the log files created by this Third Party File Based Access Log Publisher.\n  - `operation-timing-access`: The UNIX permissions of the log files created by this Operation Timing Access Log Publisher.\n  - `file-based-trace`: The UNIX permissions of the log files created by this File Based Trace Log Publisher.\n  - `common-log-file-http-operation`: The UNIX permissions of the log files created by this Common Log File HTTP Operation Log Publisher.\n  - `file-based-json-audit`: The UNIX permissions of the log files created by this File Based JSON Audit Log Publisher.\n  - `file-based-debug`: The UNIX permissions of the log files created by this File Based Debug Log Publisher.\n  - `file-based-error`: The UNIX permissions of the log files created by this File Based Error Log Publisher.\n  - `detailed-http-operation`: The UNIX permissions of the log files created by this Detailed HTTP Operation Log Publisher.\n  - `json-access`: The UNIX permissions of the log files created by this JSON Access Log Publisher.\n  - `debug-access`: The UNIX permissions of the log files created by this Debug Access Log Publisher.\n  - `file-based-audit`: The UNIX permissions of the log files created by this File Based Audit Log Publisher.\n  - `json-error`: The UNIX permissions of the log files created by this JSON Error Log Publisher.\n  - `groovy-scripted-file-based-access`: The UNIX permissions of the log files created by this Scripted File Based Access Log Publisher.\n  - `groovy-scripted-file-based-error`: The UNIX permissions of the log files created by this Scripted File Based Error Log Publisher.\n  - `third-party-file-based-error`: The UNIX permissions of the log files created by this Third Party File Based Error Log Publisher.\n  - `file-based-access`: The UNIX permissions of the log files created by this File Based Access Log Publisher.\n  - `file-based-json-http-operation`: The UNIX permissions of the log files created by this File Based JSON HTTP Operation Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"rotation_policy": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The rotation policy to use for the Third Party File Based Access Log Publisher . When the `type` attribute is set to `operation-timing-access`: The rotation policy to use for the Operation Timing Access Log Publisher . When the `type` attribute is set to `file-based-trace`: The rotation policy to use for the File Based Trace Log Publisher . When the `type` attribute is set to `common-log-file-http-operation`: The rotation policy to use for the Common Log File HTTP Operation Log Publisher . When the `type` attribute is set to `file-based-json-audit`: The rotation policy to use for the File Based JSON Audit Log Publisher . When the `type` attribute is set to `file-based-debug`: The rotation policy to use for the File Based Debug Log Publisher . When the `type` attribute is set to `file-based-error`: The rotation policy to use for the File Based Error Log Publisher . When the `type` attribute is set to `detailed-http-operation`: The rotation policy to use for the Detailed HTTP Operation Log Publisher . When the `type` attribute is set to `json-access`: The rotation policy to use for the JSON Access Log Publisher . When the `type` attribute is set to `debug-access`: The rotation policy to use for the Debug Access Log Publisher . When the `type` attribute is set to `file-based-audit`: The rotation policy to use for the File Based Audit Log Publisher . When the `type` attribute is set to `json-error`: The rotation policy to use for the JSON Error Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-access`: The rotation policy to use for the Scripted File Based Access Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-error`: The rotation policy to use for the Scripted File Based Error Log Publisher . When the `type` attribute is set to `third-party-file-based-error`: The rotation policy to use for the Third Party File Based Error Log Publisher . When the `type` attribute is set to `file-based-access`: The rotation policy to use for the File Based Access Log Publisher . When the `type` attribute is set to `file-based-json-http-operation`: The rotation policy to use for the File Based JSON HTTP Operation Log Publisher .",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The rotation policy to use for the Third Party File Based Access Log Publisher .\n  - `operation-timing-access`: The rotation policy to use for the Operation Timing Access Log Publisher .\n  - `file-based-trace`: The rotation policy to use for the File Based Trace Log Publisher .\n  - `common-log-file-http-operation`: The rotation policy to use for the Common Log File HTTP Operation Log Publisher .\n  - `file-based-json-audit`: The rotation policy to use for the File Based JSON Audit Log Publisher .\n  - `file-based-debug`: The rotation policy to use for the File Based Debug Log Publisher .\n  - `file-based-error`: The rotation policy to use for the File Based Error Log Publisher .\n  - `detailed-http-operation`: The rotation policy to use for the Detailed HTTP Operation Log Publisher .\n  - `json-access`: The rotation policy to use for the JSON Access Log Publisher .\n  - `debug-access`: The rotation policy to use for the Debug Access Log Publisher .\n  - `file-based-audit`: The rotation policy to use for the File Based Audit Log Publisher .\n  - `json-error`: The rotation policy to use for the JSON Error Log Publisher .\n  - `groovy-scripted-file-based-access`: The rotation policy to use for the Scripted File Based Access Log Publisher .\n  - `groovy-scripted-file-based-error`: The rotation policy to use for the Scripted File Based Error Log Publisher .\n  - `third-party-file-based-error`: The rotation policy to use for the Third Party File Based Error Log Publisher .\n  - `file-based-access`: The rotation policy to use for the File Based Access Log Publisher .\n  - `file-based-json-http-operation`: The rotation policy to use for the File Based JSON HTTP Operation Log Publisher .",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The retention policy to use for the Third Party File Based Access Log Publisher . When the `type` attribute is set to `operation-timing-access`: The retention policy to use for the Operation Timing Access Log Publisher . When the `type` attribute is set to `file-based-trace`: The retention policy to use for the File Based Trace Log Publisher . When the `type` attribute is set to `common-log-file-http-operation`: The retention policy to use for the Common Log File HTTP Operation Log Publisher . When the `type` attribute is set to `file-based-json-audit`: The retention policy to use for the File Based JSON Audit Log Publisher . When the `type` attribute is set to `file-based-debug`: The retention policy to use for the File Based Debug Log Publisher . When the `type` attribute is set to `file-based-error`: The retention policy to use for the File Based Error Log Publisher . When the `type` attribute is set to `detailed-http-operation`: The retention policy to use for the Detailed HTTP Operation Log Publisher . When the `type` attribute is set to `json-access`: The retention policy to use for the JSON Access Log Publisher . When the `type` attribute is set to `debug-access`: The retention policy to use for the Debug Access Log Publisher . When the `type` attribute is set to `file-based-audit`: The retention policy to use for the File Based Audit Log Publisher . When the `type` attribute is set to `json-error`: The retention policy to use for the JSON Error Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-access`: The retention policy to use for the Scripted File Based Access Log Publisher . When the `type` attribute is set to `groovy-scripted-file-based-error`: The retention policy to use for the Scripted File Based Error Log Publisher . When the `type` attribute is set to `third-party-file-based-error`: The retention policy to use for the Third Party File Based Error Log Publisher . When the `type` attribute is set to `file-based-access`: The retention policy to use for the File Based Access Log Publisher . When the `type` attribute is set to `file-based-json-http-operation`: The retention policy to use for the File Based JSON HTTP Operation Log Publisher .",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The retention policy to use for the Third Party File Based Access Log Publisher .\n  - `operation-timing-access`: The retention policy to use for the Operation Timing Access Log Publisher .\n  - `file-based-trace`: The retention policy to use for the File Based Trace Log Publisher .\n  - `common-log-file-http-operation`: The retention policy to use for the Common Log File HTTP Operation Log Publisher .\n  - `file-based-json-audit`: The retention policy to use for the File Based JSON Audit Log Publisher .\n  - `file-based-debug`: The retention policy to use for the File Based Debug Log Publisher .\n  - `file-based-error`: The retention policy to use for the File Based Error Log Publisher .\n  - `detailed-http-operation`: The retention policy to use for the Detailed HTTP Operation Log Publisher .\n  - `json-access`: The retention policy to use for the JSON Access Log Publisher .\n  - `debug-access`: The retention policy to use for the Debug Access Log Publisher .\n  - `file-based-audit`: The retention policy to use for the File Based Audit Log Publisher .\n  - `json-error`: The retention policy to use for the JSON Error Log Publisher .\n  - `groovy-scripted-file-based-access`: The retention policy to use for the Scripted File Based Access Log Publisher .\n  - `groovy-scripted-file-based-error`: The retention policy to use for the Scripted File Based Error Log Publisher .\n  - `third-party-file-based-error`: The retention policy to use for the Third Party File Based Error Log Publisher .\n  - `file-based-access`: The retention policy to use for the File Based Access Log Publisher .\n  - `file-based-json-http-operation`: The retention policy to use for the File Based JSON HTTP Operation Log Publisher .",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `groovy-scripted-file-based-access`: The set of arguments used to customize the behavior for the Scripted File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-file-based-error`: The set of arguments used to customize the behavior for the Scripted File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-access`: The set of arguments used to customize the behavior for the Scripted Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-error`: The set of arguments used to customize the behavior for the Scripted Error Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `groovy-scripted-http-operation`: The set of arguments used to customize the behavior for the Scripted HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `groovy-scripted-file-based-access`: The set of arguments used to customize the behavior for the Scripted File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-file-based-error`: The set of arguments used to customize the behavior for the Scripted File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-access`: The set of arguments used to customize the behavior for the Scripted Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-error`: The set of arguments used to customize the behavior for the Scripted Error Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `groovy-scripted-http-operation`: The set of arguments used to customize the behavior for the Scripted HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"sign_log": schema.BoolAttribute{
				Description: "Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"timestamp_precision": schema.StringAttribute{
				Description: "Specifies the smallest time unit to be included in timestamps.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encrypt_log": schema.BoolAttribute{
				Description: "Indicates whether log files should be encrypted so that their content is not available to unauthorized users.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"obscure_sensitive_content": schema.BoolAttribute{
				Description: "Indicates whether the resulting log file should attempt to obscure content that may be considered sensitive. This primarily includes the credentials for bind requests, the values of password modify extended requests and responses, and the values of any attributes specified in the obscure-attribute property. Note that the use of this option does not guarantee no sensitive information will be exposed, so the log output should still be carefully guarded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Access Log Publisher. When the `type` attribute is set to `third-party-http-operation`: The fully-qualified name of the Java class providing the logic for the Third Party HTTP Operation Log Publisher. When the `type` attribute is set to `third-party-error`: The fully-qualified name of the Java class providing the logic for the Third Party Error Log Publisher. When the `type` attribute is set to `third-party-access`: The fully-qualified name of the Java class providing the logic for the Third Party Access Log Publisher. When the `type` attribute is set to `third-party-file-based-error`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Access Log Publisher.\n  - `third-party-http-operation`: The fully-qualified name of the Java class providing the logic for the Third Party HTTP Operation Log Publisher.\n  - `third-party-error`: The fully-qualified name of the Java class providing the logic for the Third Party Error Log Publisher.\n  - `third-party-access`: The fully-qualified name of the Java class providing the logic for the Third Party Access Log Publisher.\n  - `third-party-file-based-error`: The fully-qualified name of the Java class providing the logic for the Third Party File Based Error Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"debug_aci_enabled": schema.BoolAttribute{
				Description: "Indicates whether to include debugging information about ACIs being used by the operations being logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_debug_level": schema.StringAttribute{
				Description: "The lowest severity level of debug messages to log when none of the defined targets match the message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_request_headers": schema.StringAttribute{
				Description: "Indicates whether request log messages should include information about HTTP headers included in the request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"suppressed_request_header_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request headers that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-headers property has a value of true.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_response_headers": schema.StringAttribute{
				Description: "Indicates whether response log messages should include information about HTTP headers included in the response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"suppressed_response_header_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of response headers that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-response-headers property has a value of true.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_request_authorization_type": schema.BoolAttribute{
				Description: "Indicates whether to log the type of credentials given if an \"Authorization\" header was included in the request. Logging the authorization type may be useful, and is much more secure than logging the entire value of the \"Authorization\" header.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_request_cookie_names": schema.BoolAttribute{
				Description: "Indicates whether to log the names of any cookies included in an HTTP request. Logging cookie names may be useful and is much more secure than logging the entire content of the cookies (which may include sensitive information).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_response_cookie_names": schema.BoolAttribute{
				Description: "Indicates whether to log the names of any cookies set in an HTTP response. Logging cookie names may be useful and is much more secure than logging the entire content of the cookies (which may include sensitive information).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_request_parameters": schema.StringAttribute{
				Description: "Indicates what (if any) information about request parameters should be included in request log messages. Note that this will only be used for requests with a method other than GET, since GET request parameters will be included in the request URL.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_request_protocol": schema.BoolAttribute{
				Description: "Indicates whether request log messages should include information about the HTTP version specified in the request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"suppressed_request_parameter_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request parameters that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-parameters property has a value of parameter-names or parameter-names-and-values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_redirect_uri": schema.BoolAttribute{
				Description: "Indicates whether the redirect URI (i.e., the value of the \"Location\" header from responses) should be included in response log messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_debug_category": schema.SetAttribute{
				Description: "The debug message categories to be logged when none of the defined targets match the message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"default_omit_method_entry_arguments": schema.BoolAttribute{
				Description: "Indicates whether to include method arguments in debug messages logged by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_omit_method_return_value": schema.BoolAttribute{
				Description: "Indicates whether to include the return value in debug messages logged by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_include_throwable_cause": schema.BoolAttribute{
				Description: "Indicates whether to include the cause of exceptions in exception thrown and caught messages logged by default.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_throwable_stack_frames": schema.Int64Attribute{
				Description: "Indicates the number of stack frames to include in the stack trace for method entry and exception thrown messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party-file-based-access`: The set of arguments used to customize the behavior for the Third Party File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-http-operation`: The set of arguments used to customize the behavior for the Third Party HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-error`: The set of arguments used to customize the behavior for the Third Party Error Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-access`: The set of arguments used to customize the behavior for the Third Party Access Log Publisher. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-file-based-error`: The set of arguments used to customize the behavior for the Third Party File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party-file-based-access`: The set of arguments used to customize the behavior for the Third Party File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-http-operation`: The set of arguments used to customize the behavior for the Third Party HTTP Operation Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-error`: The set of arguments used to customize the behavior for the Third Party Error Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-access`: The set of arguments used to customize the behavior for the Third Party Access Log Publisher. Each configuration property should be given in the form 'name=value'.\n  - `third-party-file-based-error`: The set of arguments used to customize the behavior for the Third Party File Based Error Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"syslog_external_server": schema.SetAttribute{
				Description: "The syslog server to which messages should be sent.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_request_details_in_result_messages": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation results should include information about both the request and the result. When the `type` attribute is set to  one of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether result log messages should include all of the elements of request log messages. This may be used to record a single message per operation with details about both the request and response.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation results should include information about both the request and the result.\n  - One of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether result log messages should include all of the elements of request log messages. This may be used to record a single message per operation with details about both the request and response.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_assurance_completed": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of replication assurance processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"debug_message_type": schema.SetAttribute{
				Description: "Specifies the debug message types which can be logged. Note that enabling these may result in sensitive information being logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"http_message_type": schema.SetAttribute{
				Description: "Specifies the HTTP message types which can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"access_token_validator_message_type": schema.SetAttribute{
				Description: "Specifies the access token validator message types that can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"id_token_validator_message_type": schema.SetAttribute{
				Description: "Specifies the ID token validator message types that can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"scim_message_type": schema.SetAttribute{
				Description: "Specifies the SCIM message types which can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"consent_message_type": schema.SetAttribute{
				Description: "Specifies the consent message types that can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"directory_rest_api_message_type": schema.SetAttribute{
				Description: "Specifies the Directory REST API message types which can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"extension_message_type": schema.SetAttribute{
				Description: "Specifies the Server SDK extension message types that can be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_path_pattern": schema.SetAttribute{
				Description: "Specifies a set of HTTP request URL paths to determine whether log messages are included for a HTTP request. Log messages are included for a HTTP request if the request path does not match any exclude-path-pattern, and the request path does match an include-path-pattern (or no include-path-pattern is specified).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"exclude_path_pattern": schema.SetAttribute{
				Description: "Specifies a set of HTTP request URL paths to determine whether log messages are excluded for a HTTP request. Log messages are included for a HTTP request if the request path does not match any exclude-path-pattern, and the request path does match an include-path-pattern (or no include-path-pattern is specified).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"server_host_name": schema.StringAttribute{
				Description: "Specifies the hostname or IP address of the syslogd host to log to. It is highly recommend to use localhost.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the syslogd host to log to.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_included_operation_processing_time": schema.StringAttribute{
				Description: "The minimum processing time (i.e., \"etime\") for operations that should be logged by this Operation Timing Access Log Publisher",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"min_included_phase_time_nanos": schema.Int64Attribute{
				Description: "The minimum length of time in nanoseconds that an operation phase should take before it is included in a log message.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"time_interval": schema.StringAttribute{
				Description: "Specifies the interval at which to check whether the log files need to be rotated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_request_details_in_search_entry_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result entries should include information about the associated search request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_request_details_in_search_reference_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result references should include information about the associated search request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_request_details_in_intermediate_response_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for intermediate responses should include information about the associated operation request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_result_code_names": schema.BoolAttribute{
				Description: "Indicates whether result log messages should include human-readable names for result codes in addition to their numeric values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_extended_search_request_details": schema.BoolAttribute{
				Description: "Indicates whether log messages for search requests should include extended information from the request, including the requested size limit, time limit, alias dereferencing behavior, and types only behavior.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_add_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for add requests should include a list of the names of the attributes included in the entry to add.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_modify_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for modify requests should include a list of the names of the attributes to be modified.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_search_entry_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result entries should include a list of the names of the attributes included in the entry that was returned.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_connects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections established to the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_disconnects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections that have been closed by the client or terminated by the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_string_length": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`operation-timing-access`, `admin-alert-access`, `file-based-trace`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Specifies the maximum number of characters that may be included in any string in a log message before that string is truncated and replaced with a placeholder indicating the number of characters that were omitted. This can help prevent extremely long log messages from being written. When the `type` attribute is set to `detailed-http-operation`: Specifies the maximum length of any individual string that should be logged. If a log message includes a string longer than this number of characters, it will be truncated. A value of zero indicates that no truncation will be used.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`operation-timing-access`, `admin-alert-access`, `file-based-trace`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Specifies the maximum number of characters that may be included in any string in a log message before that string is truncated and replaced with a placeholder indicating the number of characters that were omitted. This can help prevent extremely long log messages from being written.\n  - `detailed-http-operation`: Specifies the maximum length of any individual string that should be logged. If a log message includes a string longer than this number of characters, it will be truncated. A value of zero indicates that no truncation will be used.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"generify_message_strings_when_possible": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values. When the `type` attribute is set to  one of [`console-json-error`, `syslog-text-error`, `file-based-error`, `json-error`, `syslog-json-error`]: Indicates whether to use the generified version of the log message string (which may use placeholders like %s for a string or %d for an integer), rather than the version of the message with those placeholders replaced with specific values that would normally be written to the log.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`admin-alert-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `syslog-json-access`, `console-json-access`, `file-based-access`]: Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values.\n  - One of [`console-json-error`, `syslog-text-error`, `file-based-error`, `json-error`, `syslog-json-error`]: Indicates whether to use the generified version of the log message string (which may use placeholders like %s for a string or %d for an integer), rather than the version of the message with those placeholders replaced with specific values that would normally be written to the log.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"syslog_facility": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The syslog facility to use for the messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-based-error`: Specifies the syslog facility to use for this Syslog Based Error Log Publisher When the `type` attribute is set to `syslog-text-error`: The syslog facility to use for the messages that are logged by this Syslog Text Error Log Publisher. When the `type` attribute is set to `syslog-based-access`: Specifies the syslog facility to use for this Syslog Based Access Log Publisher When the `type` attribute is set to `syslog-text-access`: The syslog facility to use for the messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The syslog facility to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The syslog facility to use for the messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The syslog facility to use for the messages that are logged by this Syslog JSON Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The syslog facility to use for the messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-based-error`: Specifies the syslog facility to use for this Syslog Based Error Log Publisher\n  - `syslog-text-error`: The syslog facility to use for the messages that are logged by this Syslog Text Error Log Publisher.\n  - `syslog-based-access`: Specifies the syslog facility to use for this Syslog Based Access Log Publisher\n  - `syslog-text-access`: The syslog facility to use for the messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The syslog facility to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The syslog facility to use for the messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The syslog facility to use for the messages that are logged by this Syslog JSON Error Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_field_behavior": schema.StringAttribute{
				Description: "The behavior to use for determining which fields to log and whether to transform the values of those fields in any way.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_client_certificates": schema.BoolAttribute{
				Description: "Indicates whether to log information about any client certificates presented to the server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_requests": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about requests received from clients. When the `type` attribute is set to  one of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about requests received from the client.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about requests received from clients.\n  - One of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about requests received from the client.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_results": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about the results of client requests. When the `type` attribute is set to  one of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about the result of processing a requested HTTP operation.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `admin-alert-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Indicates whether to log information about the results of client requests.\n  - One of [`detailed-http-operation`, `syslog-json-http-operation`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether to record a log message with information about the result of processing a requested HTTP operation.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"log_search_entries": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result entries sent to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_search_references": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result references sent to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_intermediate_responses": schema.BoolAttribute{
				Description: "Indicates whether to log information about intermediate responses sent to the client.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-based-access`, `syslog-text-access`, `file-based-access`]: Indicates whether the Writer Based Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `third-party-file-based-access`: Indicates whether the Third Party File Based Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `operation-timing-access`: Indicates whether the Operation Timing Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `admin-alert-access`: Indicates whether the Admin Alert Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-trace`: Indicates whether the Writer Based Trace Log Publisher will publish records asynchronously. When the `type` attribute is set to `common-log-file-http-operation`: Indicates whether the Common Log File HTTP Operation Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-json-audit`: Indicates whether the File Based JSON Audit Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-debug`: Indicates whether the File Based Debug Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-error`: Indicates whether the File Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `detailed-http-operation`: Indicates whether the Detailed HTTP Operation Log Publisher will publish records asynchronously. When the `type` attribute is set to `json-access`: Indicates whether the JSON Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `debug-access`: Indicates whether the Debug Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-audit`: Indicates whether the File Based Audit Log Publisher will publish records asynchronously. When the `type` attribute is set to `json-error`: Indicates whether the JSON Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `groovy-scripted-file-based-access`: Indicates whether the Scripted File Based Access Log Publisher will publish records asynchronously. When the `type` attribute is set to `groovy-scripted-file-based-error`: Indicates whether the Scripted File Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `third-party-file-based-error`: Indicates whether the Third Party File Based Error Log Publisher will publish records asynchronously. When the `type` attribute is set to `file-based-json-http-operation`: Indicates whether the File Based JSON HTTP Operation Log Publisher will publish records asynchronously.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-based-access`, `syslog-text-access`, `file-based-access`]: Indicates whether the Writer Based Access Log Publisher will publish records asynchronously.\n  - `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher will publish records asynchronously.\n  - `third-party-file-based-access`: Indicates whether the Third Party File Based Access Log Publisher will publish records asynchronously.\n  - `operation-timing-access`: Indicates whether the Operation Timing Access Log Publisher will publish records asynchronously.\n  - `admin-alert-access`: Indicates whether the Admin Alert Access Log Publisher will publish records asynchronously.\n  - `file-based-trace`: Indicates whether the Writer Based Trace Log Publisher will publish records asynchronously.\n  - `common-log-file-http-operation`: Indicates whether the Common Log File HTTP Operation Log Publisher will publish records asynchronously.\n  - `file-based-json-audit`: Indicates whether the File Based JSON Audit Log Publisher will publish records asynchronously.\n  - `file-based-debug`: Indicates whether the File Based Debug Log Publisher will publish records asynchronously.\n  - `file-based-error`: Indicates whether the File Based Error Log Publisher will publish records asynchronously.\n  - `detailed-http-operation`: Indicates whether the Detailed HTTP Operation Log Publisher will publish records asynchronously.\n  - `json-access`: Indicates whether the JSON Access Log Publisher will publish records asynchronously.\n  - `debug-access`: Indicates whether the Debug Access Log Publisher will publish records asynchronously.\n  - `file-based-audit`: Indicates whether the File Based Audit Log Publisher will publish records asynchronously.\n  - `json-error`: Indicates whether the JSON Error Log Publisher will publish records asynchronously.\n  - `groovy-scripted-file-based-access`: Indicates whether the Scripted File Based Access Log Publisher will publish records asynchronously.\n  - `groovy-scripted-file-based-error`: Indicates whether the Scripted File Based Error Log Publisher will publish records asynchronously.\n  - `third-party-file-based-error`: Indicates whether the Third Party File Based Error Log Publisher will publish records asynchronously.\n  - `file-based-json-http-operation`: Indicates whether the File Based JSON HTTP Operation Log Publisher will publish records asynchronously.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"correlate_requests_and_results": schema.BoolAttribute{
				Description: "Indicates whether to automatically log result messages for any operation in which the corresponding request was logged. In such cases, the result, entry, and reference criteria will be ignored, although the log-responses, log-search-entries, and log-search-references properties will be honored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"syslog_severity": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The syslog severity to use for the messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-text-error`: The syslog severity to use for the messages that are logged by this Syslog Text Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message. When the `type` attribute is set to `syslog-text-access`: The syslog severity to use for the messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The syslog severity to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The syslog severity to use for the messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The syslog severity to use for the messages that are logged by this Syslog JSON Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The syslog severity to use for the messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-text-error`: The syslog severity to use for the messages that are logged by this Syslog Text Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.\n  - `syslog-text-access`: The syslog severity to use for the messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The syslog severity to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The syslog severity to use for the messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The syslog severity to use for the messages that are logged by this Syslog JSON Error Log Publisher. If this is not specified, then the severity for each syslog message will be automatically based on the severity for the associated log message.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"search_entry_criteria": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher. When the `type` attribute is set to `admin-alert-access`: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Admin Alert Access Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher.\n  - `admin-alert-access`: Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Admin Alert Access Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"search_reference_criteria": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher. When the `type` attribute is set to `admin-alert-access`: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Admin Alert Access Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`third-party-file-based-access`, `jdbc-based-access`, `syslog-based-access`, `syslog-text-access`, `json-access`, `debug-access`, `third-party-access`, `groovy-scripted-file-based-access`, `syslog-json-access`, `groovy-scripted-access`, `console-json-access`, `file-based-access`]: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher.\n  - `admin-alert-access`: Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Admin Alert Access Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"syslog_message_host_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-text-error`: The local host name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher. When the `type` attribute is set to `syslog-text-access`: The local host name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The local host name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-text-error`: The local host name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher.\n  - `syslog-text-access`: The local host name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The local host name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The local host name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"syslog_message_application_name": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `syslog-json-audit`: The application name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher. When the `type` attribute is set to `syslog-text-error`: The application name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher. When the `type` attribute is set to `syslog-text-access`: The application name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher. When the `type` attribute is set to `syslog-json-http-operation`: The application name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher. When the `type` attribute is set to `syslog-json-access`: The application name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher. When the `type` attribute is set to `syslog-json-error`: The application name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `syslog-json-audit`: The application name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.\n  - `syslog-text-error`: The application name that will be included in syslog messages that are logged by this Syslog Text Error Log Publisher.\n  - `syslog-text-access`: The application name that will be included in syslog messages that are logged by this Syslog Text Access Log Publisher.\n  - `syslog-json-http-operation`: The application name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher.\n  - `syslog-json-access`: The application name that will be included in syslog messages that are logged by this Syslog JSON Access Log Publisher.\n  - `syslog-json-error`: The application name that will be included in syslog messages that are logged by this Syslog JSON Error Log Publisher.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"write_multi_line_messages": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `file-based-json-audit`, `syslog-json-http-operation`, `console-json-audit`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools. When the `type` attribute is set to  one of [`console-json-error`, `json-access`, `json-error`, `console-json-access`]: Indicates whether the JSON objects should be formatted to span multiple lines with a single element on each line. The multi-line format is potentially more user friendly (if administrators may need to look at the log files), but each message will be larger because of the additional spaces and end-of-line markers.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `file-based-json-audit`, `syslog-json-http-operation`, `console-json-audit`, `console-json-http-operation`, `file-based-json-http-operation`]: Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools.\n  - One of [`console-json-error`, `json-access`, `json-error`, `console-json-access`]: Indicates whether the JSON objects should be formatted to span multiple lines with a single element on each line. The multi-line format is potentially more user friendly (if administrators may need to look at the log files), but each message will be larger because of the additional spaces and end-of-line markers.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"use_reversible_form": schema.BoolAttribute{
				Description: "Indicates whether the audit log should be written in reversible form so that it is possible to revert the changes if desired.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"soft_delete_entry_audit_behavior": schema.StringAttribute{
				Description: "Specifies the audit behavior for delete and modify operations on soft-deleted entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_operation_purpose_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any operation purpose request control that may have been included in the request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_intermediate_client_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any intermediate client request control that may have been included in the request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"obscure_attribute": schema.SetAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `file-based-json-audit`, `file-based-audit`, `console-json-audit`]: Specifies the names of any attribute types that should have their values obscured in the audit log because they may be considered sensitive. When the `type` attribute is set to `debug-access`: Specifies the names of any attribute types that should have their values obscured if the obscure-sensitive-content property has a value of true.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `file-based-json-audit`, `file-based-audit`, `console-json-audit`]: Specifies the names of any attribute types that should have their values obscured in the audit log because they may be considered sensitive.\n  - `debug-access`: Specifies the names of any attribute types that should have their values obscured if the obscure-sensitive-content property has a value of true.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"exclude_attribute": schema.SetAttribute{
				Description: "Specifies the names of any attribute types that should be excluded from the audit log.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_product_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the product name for the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_instance_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the instance name for the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_startup_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the startup ID for the Directory Server, which is a value assigned to the server instance at startup and may be used to identify when the server has been restarted.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_thread_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the thread ID for the Directory Server in each log message. This ID can be used to correlate log messages from the same thread within a single log as well as generated by the same thread across different types of log files. More information about the thread with a specific ID can be obtained using the cn=JVM Stack Trace,cn=monitor entry.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_requester_dn": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `admin-alert-access`, `syslog-based-access`, `file-based-json-audit`, `syslog-text-access`, `json-access`, `file-based-audit`, `syslog-json-access`, `console-json-audit`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested. When the `type` attribute is set to `operation-timing-access`: Indicates whether log messages should include the DN of the authenticated user for the client connection on which the operation was requested.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `admin-alert-access`, `syslog-based-access`, `file-based-json-audit`, `syslog-text-access`, `json-access`, `file-based-audit`, `syslog-json-access`, `console-json-audit`, `console-json-access`, `file-based-access`]: Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested.\n  - `operation-timing-access`: Indicates whether log messages should include the DN of the authenticated user for the client connection on which the operation was requested.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"include_requester_ip_address": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the IP address of the client that requested the operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_request_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include a list of the OIDs of any controls included in the request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_response_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation results should include a list of the OIDs of any controls included in the result.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_replication_change_id": schema.BoolAttribute{
				Description: "Indicates whether to log information about the replication change ID.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"suppress_replication_operations": schema.BoolAttribute{
				Description: "Indicates whether access messages that are generated by replication operations should be suppressed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"connection_criteria": schema.StringAttribute{
				Description: "Specifies a set of connection criteria that must match the associated client connection in order for a connect, disconnect, request, or result message to be logged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"request_criteria": schema.StringAttribute{
				Description: "Specifies a set of request criteria that must match the associated operation request in order for a request or result to be logged by this Access Log Publisher.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"result_criteria": schema.StringAttribute{
				Description: "Specifies a set of result criteria that must match the associated operation result in order for that result to be logged by this Access Log Publisher.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`syslog-json-audit`, `third-party-file-based-access`, `operation-timing-access`, `third-party-http-operation`, `admin-alert-access`, `file-based-trace`, `jdbc-based-error`, `jdbc-based-access`, `common-log-file-http-operation`, `syslog-text-error`, `file-based-json-audit`, `file-based-debug`, `file-based-error`, `third-party-error`, `syslog-text-access`, `detailed-http-operation`, `json-access`, `debug-access`, `syslog-json-http-operation`, `third-party-access`, `file-based-audit`, `json-error`, `groovy-scripted-file-based-access`, `groovy-scripted-file-based-error`, `syslog-json-access`, `groovy-scripted-access`, `third-party-file-based-error`, `file-based-access`, `groovy-scripted-error`, `file-based-json-http-operation`, `syslog-json-error`, `groovy-scripted-http-operation`]: Indicates whether the Log Publisher is enabled for use. When the `type` attribute is set to `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher is enabled for use. When the `type` attribute is set to `console-json-error`: Indicates whether the Console JSON Error Log Publisher is enabled for use. When the `type` attribute is set to `syslog-based-access`: Indicates whether the Syslog Based Access Log Publisher is enabled for use. When the `type` attribute is set to `console-json-audit`: Indicates whether the Console JSON Audit Log Publisher is enabled for use. When the `type` attribute is set to `console-json-http-operation`: Indicates whether the Console JSON HTTP Operation Log Publisher is enabled for use. When the `type` attribute is set to `console-json-access`: Indicates whether the Console JSON Access Log Publisher is enabled for use.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`syslog-json-audit`, `third-party-file-based-access`, `operation-timing-access`, `third-party-http-operation`, `admin-alert-access`, `file-based-trace`, `jdbc-based-error`, `jdbc-based-access`, `common-log-file-http-operation`, `syslog-text-error`, `file-based-json-audit`, `file-based-debug`, `file-based-error`, `third-party-error`, `syslog-text-access`, `detailed-http-operation`, `json-access`, `debug-access`, `syslog-json-http-operation`, `third-party-access`, `file-based-audit`, `json-error`, `groovy-scripted-file-based-access`, `groovy-scripted-file-based-error`, `syslog-json-access`, `groovy-scripted-access`, `third-party-file-based-error`, `file-based-access`, `groovy-scripted-error`, `file-based-json-http-operation`, `syslog-json-error`, `groovy-scripted-http-operation`]: Indicates whether the Log Publisher is enabled for use.\n  - `syslog-based-error`: Indicates whether the Syslog Based Error Log Publisher is enabled for use.\n  - `console-json-error`: Indicates whether the Console JSON Error Log Publisher is enabled for use.\n  - `syslog-based-access`: Indicates whether the Syslog Based Access Log Publisher is enabled for use.\n  - `console-json-audit`: Indicates whether the Console JSON Audit Log Publisher is enabled for use.\n  - `console-json-http-operation`: Indicates whether the Console JSON HTTP Operation Log Publisher is enabled for use.\n  - `console-json-access`: Indicates whether the Console JSON Access Log Publisher is enabled for use.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SyslogJsonAuditLogPublisherResponse object into the model struct
func readSyslogJsonAuditLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogJsonAuditLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, false)
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), false)
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogBasedErrorLogPublisherResponse object into the model struct
func readSyslogBasedErrorLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogBasedErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ThirdPartyFileBasedAccessLogPublisherResponse object into the model struct
func readThirdPartyFileBasedAccessLogPublisherResponseDataSource(ctx context.Context, r *client.ThirdPartyFileBasedAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-file-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a OperationTimingAccessLogPublisherResponse object into the model struct
func readOperationTimingAccessLogPublisherResponseDataSource(ctx context.Context, r *client.OperationTimingAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("operation-timing-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequesterIPAddress = internaltypes.BoolTypeOrNil(r.IncludeRequesterIPAddress)
	state.IncludeRequesterDN = internaltypes.BoolTypeOrNil(r.IncludeRequesterDN)
	state.MinIncludedOperationProcessingTime = internaltypes.StringTypeOrNil(r.MinIncludedOperationProcessingTime, false)
	state.MinIncludedPhaseTimeNanos = internaltypes.Int64TypeOrNil(r.MinIncludedPhaseTimeNanos)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ThirdPartyHttpOperationLogPublisherResponse object into the model struct
func readThirdPartyHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.ThirdPartyHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a AdminAlertAccessLogPublisherResponse object into the model struct
func readAdminAlertAccessLogPublisherResponseDataSource(ctx context.Context, r *client.AdminAlertAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
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
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedTraceLogPublisherResponse object into the model struct
func readFileBasedTraceLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedTraceLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-trace")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a JdbcBasedErrorLogPublisherResponse object into the model struct
func readJdbcBasedErrorLogPublisherResponseDataSource(ctx context.Context, r *client.JdbcBasedErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a JdbcBasedAccessLogPublisherResponse object into the model struct
func readJdbcBasedAccessLogPublisherResponseDataSource(ctx context.Context, r *client.JdbcBasedAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a CommonLogFileHttpOperationLogPublisherResponse object into the model struct
func readCommonLogFileHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.CommonLogFileHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("common-log-file-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ConsoleJsonErrorLogPublisherResponse object into the model struct
func readConsoleJsonErrorLogPublisherResponseDataSource(ctx context.Context, r *client.ConsoleJsonErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), false)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogTextErrorLogPublisherResponse object into the model struct
func readSyslogTextErrorLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogTextErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-text-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), false)
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, false)
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, false)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogBasedAccessLogPublisherResponse object into the model struct
func readSyslogBasedAccessLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogBasedAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedJsonAuditLogPublisherResponse object into the model struct
func readFileBasedJsonAuditLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedJsonAuditLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherFileBasedJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), false)
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedDebugLogPublisherResponse object into the model struct
func readFileBasedDebugLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedDebugLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-debug")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), false)
	state.DefaultDebugLevel = types.StringValue(r.DefaultDebugLevel.String())
	state.DefaultDebugCategory = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultDebugCategoryProp(r.DefaultDebugCategory))
	state.DefaultOmitMethodEntryArguments = internaltypes.BoolTypeOrNil(r.DefaultOmitMethodEntryArguments)
	state.DefaultOmitMethodReturnValue = internaltypes.BoolTypeOrNil(r.DefaultOmitMethodReturnValue)
	state.DefaultIncludeThrowableCause = internaltypes.BoolTypeOrNil(r.DefaultIncludeThrowableCause)
	state.DefaultThrowableStackFrames = internaltypes.Int64TypeOrNil(r.DefaultThrowableStackFrames)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedErrorLogPublisherResponse object into the model struct
func readFileBasedErrorLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), false)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ThirdPartyErrorLogPublisherResponse object into the model struct
func readThirdPartyErrorLogPublisherResponseDataSource(ctx context.Context, r *client.ThirdPartyErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogTextAccessLogPublisherResponse object into the model struct
func readSyslogTextAccessLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogTextAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-text-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, false)
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, false)
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
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), false)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a DetailedHttpOperationLogPublisherResponse object into the model struct
func readDetailedHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.DetailedHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("detailed-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), false)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), false)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), false)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a JsonAccessLogPublisherResponse object into the model struct
func readJsonAccessLogPublisherResponseDataSource(ctx context.Context, r *client.JsonAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogAssuranceCompleted = internaltypes.BoolTypeOrNil(r.LogAssuranceCompleted)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
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
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a DebugAccessLogPublisherResponse object into the model struct
func readDebugAccessLogPublisherResponseDataSource(ctx context.Context, r *client.DebugAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ObscureSensitiveContent = internaltypes.BoolTypeOrNil(r.ObscureSensitiveContent)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.DebugACIEnabled = internaltypes.BoolTypeOrNil(r.DebugACIEnabled)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.LogConnects = internaltypes.BoolTypeOrNil(r.LogConnects)
	state.LogDisconnects = internaltypes.BoolTypeOrNil(r.LogDisconnects)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogJsonHttpOperationLogPublisherResponse object into the model struct
func readSyslogJsonHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogJsonHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, false)
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), false)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), false)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), false)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ThirdPartyAccessLogPublisherResponse object into the model struct
func readThirdPartyAccessLogPublisherResponseDataSource(ctx context.Context, r *client.ThirdPartyAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedAuditLogPublisherResponse object into the model struct
func readFileBasedAuditLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedAuditLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
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
		client.StringPointerEnumlogPublisherFileBasedAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), false)
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), false)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a JsonErrorLogPublisherResponse object into the model struct
func readJsonErrorLogPublisherResponseDataSource(ctx context.Context, r *client.JsonErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a GroovyScriptedFileBasedAccessLogPublisherResponse object into the model struct
func readGroovyScriptedFileBasedAccessLogPublisherResponseDataSource(ctx context.Context, r *client.GroovyScriptedFileBasedAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a GroovyScriptedFileBasedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedFileBasedErrorLogPublisherResponseDataSource(ctx context.Context, r *client.GroovyScriptedFileBasedErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogJsonAccessLogPublisherResponse object into the model struct
func readSyslogJsonAccessLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogJsonAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, false)
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, false)
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
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a GroovyScriptedAccessLogPublisherResponse object into the model struct
func readGroovyScriptedAccessLogPublisherResponseDataSource(ctx context.Context, r *client.GroovyScriptedAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ThirdPartyFileBasedErrorLogPublisherResponse object into the model struct
func readThirdPartyFileBasedErrorLogPublisherResponseDataSource(ctx context.Context, r *client.ThirdPartyFileBasedErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-file-based-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ConsoleJsonAuditLogPublisherResponse object into the model struct
func readConsoleJsonAuditLogPublisherResponseDataSource(ctx context.Context, r *client.ConsoleJsonAuditLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), false)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherConsoleJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), false)
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ConsoleJsonHttpOperationLogPublisherResponse object into the model struct
func readConsoleJsonHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.ConsoleJsonHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), false)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), false)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), false)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), false)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a ConsoleJsonAccessLogPublisherResponse object into the model struct
func readConsoleJsonAccessLogPublisherResponseDataSource(ctx context.Context, r *client.ConsoleJsonAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), false)
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
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
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
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedAccessLogPublisherResponse object into the model struct
func readFileBasedAccessLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedAccessLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), false)
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
	state.LogFieldBehavior = internaltypes.StringTypeOrNil(r.LogFieldBehavior, false)
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogClientCertificates = internaltypes.BoolTypeOrNil(r.LogClientCertificates)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.LogSearchEntries = internaltypes.BoolTypeOrNil(r.LogSearchEntries)
	state.LogSearchReferences = internaltypes.BoolTypeOrNil(r.LogSearchReferences)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.CorrelateRequestsAndResults = internaltypes.BoolTypeOrNil(r.CorrelateRequestsAndResults)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, false)
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, false)
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, false)
	state.SearchEntryCriteria = internaltypes.StringTypeOrNil(r.SearchEntryCriteria, false)
	state.SearchReferenceCriteria = internaltypes.StringTypeOrNil(r.SearchReferenceCriteria, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a GroovyScriptedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedErrorLogPublisherResponseDataSource(ctx context.Context, r *client.GroovyScriptedErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a FileBasedJsonHttpOperationLogPublisherResponse object into the model struct
func readFileBasedJsonHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.FileBasedJsonHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), false)
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, false)
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), false)
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), false)
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), false)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a SyslogJsonErrorLogPublisherResponse object into the model struct
func readSyslogJsonErrorLogPublisherResponseDataSource(ctx context.Context, r *client.SyslogJsonErrorLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-error")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), false)
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, false)
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, false)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read a GroovyScriptedHttpOperationLogPublisherResponse object into the model struct
func readGroovyScriptedHttpOperationLogPublisherResponseDataSource(ctx context.Context, r *client.GroovyScriptedHttpOperationLogPublisherResponse, state *logPublisherDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), false)
}

// Read resource information
func (r *logPublisherDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logPublisherDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherAPI.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SyslogJsonAuditLogPublisherResponse != nil {
		readSyslogJsonAuditLogPublisherResponseDataSource(ctx, readResponse.SyslogJsonAuditLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogBasedErrorLogPublisherResponse != nil {
		readSyslogBasedErrorLogPublisherResponseDataSource(ctx, readResponse.SyslogBasedErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyFileBasedAccessLogPublisherResponse != nil {
		readThirdPartyFileBasedAccessLogPublisherResponseDataSource(ctx, readResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.OperationTimingAccessLogPublisherResponse != nil {
		readOperationTimingAccessLogPublisherResponseDataSource(ctx, readResponse.OperationTimingAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyHttpOperationLogPublisherResponse != nil {
		readThirdPartyHttpOperationLogPublisherResponseDataSource(ctx, readResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AdminAlertAccessLogPublisherResponse != nil {
		readAdminAlertAccessLogPublisherResponseDataSource(ctx, readResponse.AdminAlertAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedTraceLogPublisherResponse != nil {
		readFileBasedTraceLogPublisherResponseDataSource(ctx, readResponse.FileBasedTraceLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcBasedErrorLogPublisherResponse != nil {
		readJdbcBasedErrorLogPublisherResponseDataSource(ctx, readResponse.JdbcBasedErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JdbcBasedAccessLogPublisherResponse != nil {
		readJdbcBasedAccessLogPublisherResponseDataSource(ctx, readResponse.JdbcBasedAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CommonLogFileHttpOperationLogPublisherResponse != nil {
		readCommonLogFileHttpOperationLogPublisherResponseDataSource(ctx, readResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonErrorLogPublisherResponse != nil {
		readConsoleJsonErrorLogPublisherResponseDataSource(ctx, readResponse.ConsoleJsonErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogTextErrorLogPublisherResponse != nil {
		readSyslogTextErrorLogPublisherResponseDataSource(ctx, readResponse.SyslogTextErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogBasedAccessLogPublisherResponse != nil {
		readSyslogBasedAccessLogPublisherResponseDataSource(ctx, readResponse.SyslogBasedAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedJsonAuditLogPublisherResponse != nil {
		readFileBasedJsonAuditLogPublisherResponseDataSource(ctx, readResponse.FileBasedJsonAuditLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedDebugLogPublisherResponse != nil {
		readFileBasedDebugLogPublisherResponseDataSource(ctx, readResponse.FileBasedDebugLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedErrorLogPublisherResponse != nil {
		readFileBasedErrorLogPublisherResponseDataSource(ctx, readResponse.FileBasedErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyErrorLogPublisherResponse != nil {
		readThirdPartyErrorLogPublisherResponseDataSource(ctx, readResponse.ThirdPartyErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogTextAccessLogPublisherResponse != nil {
		readSyslogTextAccessLogPublisherResponseDataSource(ctx, readResponse.SyslogTextAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DetailedHttpOperationLogPublisherResponse != nil {
		readDetailedHttpOperationLogPublisherResponseDataSource(ctx, readResponse.DetailedHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JsonAccessLogPublisherResponse != nil {
		readJsonAccessLogPublisherResponseDataSource(ctx, readResponse.JsonAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.DebugAccessLogPublisherResponse != nil {
		readDebugAccessLogPublisherResponseDataSource(ctx, readResponse.DebugAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonHttpOperationLogPublisherResponse != nil {
		readSyslogJsonHttpOperationLogPublisherResponseDataSource(ctx, readResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccessLogPublisherResponse != nil {
		readThirdPartyAccessLogPublisherResponseDataSource(ctx, readResponse.ThirdPartyAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedAuditLogPublisherResponse != nil {
		readFileBasedAuditLogPublisherResponseDataSource(ctx, readResponse.FileBasedAuditLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JsonErrorLogPublisherResponse != nil {
		readJsonErrorLogPublisherResponseDataSource(ctx, readResponse.JsonErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse != nil {
		readGroovyScriptedFileBasedAccessLogPublisherResponseDataSource(ctx, readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse != nil {
		readGroovyScriptedFileBasedErrorLogPublisherResponseDataSource(ctx, readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonAccessLogPublisherResponse != nil {
		readSyslogJsonAccessLogPublisherResponseDataSource(ctx, readResponse.SyslogJsonAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedAccessLogPublisherResponse != nil {
		readGroovyScriptedAccessLogPublisherResponseDataSource(ctx, readResponse.GroovyScriptedAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyFileBasedErrorLogPublisherResponse != nil {
		readThirdPartyFileBasedErrorLogPublisherResponseDataSource(ctx, readResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonAuditLogPublisherResponse != nil {
		readConsoleJsonAuditLogPublisherResponseDataSource(ctx, readResponse.ConsoleJsonAuditLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonHttpOperationLogPublisherResponse != nil {
		readConsoleJsonHttpOperationLogPublisherResponseDataSource(ctx, readResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConsoleJsonAccessLogPublisherResponse != nil {
		readConsoleJsonAccessLogPublisherResponseDataSource(ctx, readResponse.ConsoleJsonAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedAccessLogPublisherResponse != nil {
		readFileBasedAccessLogPublisherResponseDataSource(ctx, readResponse.FileBasedAccessLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedErrorLogPublisherResponse != nil {
		readGroovyScriptedErrorLogPublisherResponseDataSource(ctx, readResponse.GroovyScriptedErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedJsonHttpOperationLogPublisherResponse != nil {
		readFileBasedJsonHttpOperationLogPublisherResponseDataSource(ctx, readResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SyslogJsonErrorLogPublisherResponse != nil {
		readSyslogJsonErrorLogPublisherResponseDataSource(ctx, readResponse.SyslogJsonErrorLogPublisherResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedHttpOperationLogPublisherResponse != nil {
		readGroovyScriptedHttpOperationLogPublisherResponseDataSource(ctx, readResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
