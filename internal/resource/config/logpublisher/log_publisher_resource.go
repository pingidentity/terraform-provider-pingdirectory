package logpublisher

import (
	"context"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
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
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
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
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type logPublisherResourceModel struct {
	Id                                                  types.String `tfsdk:"id"`
	LastUpdated                                         types.String `tfsdk:"last_updated"`
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
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted File Based Access Log Publisher.",
				Optional:    true,
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"output_location": schema.StringAttribute{
				Description: "Specifies the output stream to which JSON-formatted error log messages should be written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the Third Party File Based Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Optional:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this Third Party File Based Access Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the Third Party File Based Access Log Publisher .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Third Party File Based Access Log Publisher .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"sign_log": schema.BoolAttribute{
				Description: "Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"timestamp_precision": schema.StringAttribute{
				Description: "Specifies the smallest time unit to be included in timestamps.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encrypt_log": schema.BoolAttribute{
				Description: "Indicates whether log files should be encrypted so that their content is not available to unauthorized users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "Specifies the ID of the encryption settings definition that should be used to encrypt the data. If this is not provided, the server's preferred encryption settings definition will be used. The \"encryption-settings list\" command can be used to obtain a list of the encryption settings definitions available in the server.",
				Optional:    true,
			},
			"append": schema.BoolAttribute{
				Description: "Specifies whether to append to existing log files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"obscure_sensitive_content": schema.BoolAttribute{
				Description: "Indicates whether the resulting log file should attempt to obscure content that may be considered sensitive. This primarily includes the credentials for bind requests, the values of password modify extended requests and responses, and the values of any attributes specified in the obscure-attribute property. Note that the use of this option does not guarantee no sensitive information will be exposed, so the log output should still be carefully guarded.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party File Based Access Log Publisher.",
				Optional:    true,
			},
			"debug_aci_enabled": schema.BoolAttribute{
				Description: "Indicates whether to include debugging information about ACIs being used by the operations being logged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_debug_level": schema.StringAttribute{
				Description: "The lowest severity level of debug messages to log when none of the defined targets match the message.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_request_headers": schema.StringAttribute{
				Description: "Indicates whether request log messages should include information about HTTP headers included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"suppressed_request_header_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request headers that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-headers property has a value of true.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"log_response_headers": schema.StringAttribute{
				Description: "Indicates whether response log messages should include information about HTTP headers included in the response.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"suppressed_response_header_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of response headers that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-response-headers property has a value of true.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"log_request_authorization_type": schema.BoolAttribute{
				Description: "Indicates whether to log the type of credentials given if an \"Authorization\" header was included in the request. Logging the authorization type may be useful, and is much more secure than logging the entire value of the \"Authorization\" header.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_request_cookie_names": schema.BoolAttribute{
				Description: "Indicates whether to log the names of any cookies included in an HTTP request. Logging cookie names may be useful and is much more secure than logging the entire content of the cookies (which may include sensitive information).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_response_cookie_names": schema.BoolAttribute{
				Description: "Indicates whether to log the names of any cookies set in an HTTP response. Logging cookie names may be useful and is much more secure than logging the entire content of the cookies (which may include sensitive information).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_request_parameters": schema.StringAttribute{
				Description: "Indicates what (if any) information about request parameters should be included in request log messages. Note that this will only be used for requests with a method other than GET, since GET request parameters will be included in the request URL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_request_protocol": schema.BoolAttribute{
				Description: "Indicates whether request log messages should include information about the HTTP version specified in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"suppressed_request_parameter_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request parameters that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-parameters property has a value of parameter-names or parameter-names-and-values.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"log_redirect_uri": schema.BoolAttribute{
				Description: "Indicates whether the redirect URI (i.e., the value of the \"Location\" header from responses) should be included in response log messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_debug_category": schema.SetAttribute{
				Description: "The debug message categories to be logged when none of the defined targets match the message.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"default_omit_method_entry_arguments": schema.BoolAttribute{
				Description: "Indicates whether to include method arguments in debug messages logged by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_omit_method_return_value": schema.BoolAttribute{
				Description: "Indicates whether to include the return value in debug messages logged by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_include_throwable_cause": schema.BoolAttribute{
				Description: "Indicates whether to include the cause of exceptions in exception thrown and caught messages logged by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"default_throwable_stack_frames": schema.Int64Attribute{
				Description: "Indicates the number of stack frames to include in the stack trace for method entry and exception thrown messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party File Based Access Log Publisher. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Indicates whether log messages for operation results should include information about both the request and the result.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_assurance_completed": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of replication assurance processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"debug_message_type": schema.SetAttribute{
				Description: "Specifies the debug message types which can be logged. Note that enabling these may result in sensitive information being logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"http_message_type": schema.SetAttribute{
				Description: "Specifies the HTTP message types which can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_validator_message_type": schema.SetAttribute{
				Description: "Specifies the access token validator message types that can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"id_token_validator_message_type": schema.SetAttribute{
				Description: "Specifies the ID token validator message types that can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"scim_message_type": schema.SetAttribute{
				Description: "Specifies the SCIM message types which can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"consent_message_type": schema.SetAttribute{
				Description: "Specifies the consent message types that can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"directory_rest_api_message_type": schema.SetAttribute{
				Description: "Specifies the Directory REST API message types which can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"extension_message_type": schema.SetAttribute{
				Description: "Specifies the Server SDK extension message types that can be logged.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"include_path_pattern": schema.SetAttribute{
				Description: "Specifies a set of HTTP request URL paths to determine whether log messages are included for a HTTP request. Log messages are included for a HTTP request if the request path does not match any exclude-path-pattern, and the request path does match an include-path-pattern (or no include-path-pattern is specified).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_path_pattern": schema.SetAttribute{
				Description: "Specifies a set of HTTP request URL paths to determine whether log messages are excluded for a HTTP request. Log messages are included for a HTTP request if the request path does not match any exclude-path-pattern, and the request path does match an include-path-pattern (or no include-path-pattern is specified).",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"server_host_name": schema.StringAttribute{
				Description: "Specifies the hostname or IP address of the syslogd host to log to. It is highly recommend to use localhost.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_port": schema.Int64Attribute{
				Description: "Specifies the port number of the syslogd host to log to.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_included_operation_processing_time": schema.StringAttribute{
				Description: "The minimum processing time (i.e., \"etime\") for operations that should be logged by this Operation Timing Access Log Publisher",
				Optional:    true,
			},
			"min_included_phase_time_nanos": schema.Int64Attribute{
				Description: "The minimum length of time in nanoseconds that an operation phase should take before it is included in a log message.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_details_in_search_reference_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result references should include information about the associated search request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_details_in_intermediate_response_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for intermediate responses should include information about the associated operation request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_result_code_names": schema.BoolAttribute{
				Description: "Indicates whether result log messages should include human-readable names for result codes in addition to their numeric values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_extended_search_request_details": schema.BoolAttribute{
				Description: "Indicates whether log messages for search requests should include extended information from the request, including the requested size limit, time limit, alias dereferencing behavior, and types only behavior.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_add_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for add requests should include a list of the names of the attributes included in the entry to add.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_modify_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for modify requests should include a list of the names of the attributes to be modified.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_search_entry_attribute_names": schema.BoolAttribute{
				Description: "Indicates whether log messages for search result entries should include a list of the names of the attributes included in the entry that was returned.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_connects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections established to the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_disconnects": schema.BoolAttribute{
				Description: "Indicates whether to log information about connections that have been closed by the client or terminated by the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"max_string_length": schema.Int64Attribute{
				Description: "Specifies the maximum number of characters that may be included in any string in a log message before that string is truncated and replaced with a placeholder indicating the number of characters that were omitted. This can help prevent extremely long log messages from being written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"generify_message_strings_when_possible": schema.BoolAttribute{
				Description: "Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"syslog_facility": schema.StringAttribute{
				Description: "The syslog facility to use for the messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_field_behavior": schema.StringAttribute{
				Description: "The behavior to use for determining which fields to log and whether to transform the values of those fields in any way.",
				Optional:    true,
			},
			"log_client_certificates": schema.BoolAttribute{
				Description: "Indicates whether to log information about any client certificates presented to the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_requests": schema.BoolAttribute{
				Description: "Indicates whether to log information about requests received from clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_results": schema.BoolAttribute{
				Description: "Indicates whether to log information about the results of client requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_search_entries": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result entries sent to the client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_search_references": schema.BoolAttribute{
				Description: "Indicates whether to log information about search result references sent to the client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_intermediate_responses": schema.BoolAttribute{
				Description: "Indicates whether to log information about intermediate responses sent to the client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the Syslog Based Error Log Publisher will publish records asynchronously.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"correlate_requests_and_results": schema.BoolAttribute{
				Description: "Indicates whether to automatically log result messages for any operation in which the corresponding request was logged. In such cases, the result, entry, and reference criteria will be ignored, although the log-responses, log-search-entries, and log-search-references properties will be honored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"syslog_severity": schema.StringAttribute{
				Description: "The syslog severity to use for the messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_severity": schema.SetAttribute{
				Description: "Specifies the default severity levels for the logger.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"override_severity": schema.SetAttribute{
				Description: "Specifies the override severity levels for the logger based on the category of the messages.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"search_entry_criteria": schema.StringAttribute{
				Description: "Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"search_reference_criteria": schema.StringAttribute{
				Description: "Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"syslog_message_host_name": schema.StringAttribute{
				Description: "The local host name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"syslog_message_application_name": schema.StringAttribute{
				Description: "The application name that will be included in syslog messages that are logged by this Syslog JSON Audit Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"write_multi_line_messages": schema.BoolAttribute{
				Description: "Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_reversible_form": schema.BoolAttribute{
				Description: "Indicates whether the audit log should be written in reversible form so that it is possible to revert the changes if desired.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_delete_entry_audit_behavior": schema.StringAttribute{
				Description: "Specifies the audit behavior for delete and modify operations on soft-deleted entries.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"include_operation_purpose_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any operation purpose request control that may have been included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_intermediate_client_request_control": schema.BoolAttribute{
				Description: "Indicates whether to include information about any intermediate client request control that may have been included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"obscure_attribute": schema.SetAttribute{
				Description: "Specifies the names of any attribute types that should have their values obscured in the audit log because they may be considered sensitive.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"exclude_attribute": schema.SetAttribute{
				Description: "Specifies the names of any attribute types that should be excluded from the audit log.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_product_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the product name for the Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_instance_name": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the instance name for the Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_startup_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the startup ID for the Directory Server, which is a value assigned to the server instance at startup and may be used to identify when the server has been restarted.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_thread_id": schema.BoolAttribute{
				Description: "Indicates whether log messages should include the thread ID for the Directory Server in each log message. This ID can be used to correlate log messages from the same thread within a single log as well as generated by the same thread across different types of log files. More information about the thread with a specific ID can be obtained using the cn=JVM Stack Trace,cn=monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_requester_dn": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_requester_ip_address": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the IP address of the client that requested the operation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_request_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include a list of the OIDs of any controls included in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_response_controls": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation results should include a list of the OIDs of any controls included in the result.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"include_replication_change_id": schema.BoolAttribute{
				Description: "Indicates whether to log information about the replication change ID.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"suppress_replication_operations": schema.BoolAttribute{
				Description: "Indicates whether access messages that are generated by replication operations should be suppressed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
				Description: "Indicates whether the Log Publisher is enabled for use.",
				Required:    true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"syslog-json-audit", "syslog-based-error", "third-party-file-based-access", "operation-timing-access", "third-party-http-operation", "admin-alert-access", "file-based-trace", "jdbc-based-error", "jdbc-based-access", "common-log-file-http-operation", "console-json-error", "syslog-text-error", "syslog-based-access", "file-based-json-audit", "file-based-debug", "file-based-error", "third-party-error", "syslog-text-access", "detailed-http-operation", "json-access", "debug-access", "syslog-json-http-operation", "third-party-access", "file-based-audit", "json-error", "groovy-scripted-file-based-access", "groovy-scripted-file-based-error", "syslog-json-access", "groovy-scripted-access", "third-party-file-based-error", "console-json-audit", "console-json-http-operation", "console-json-access", "file-based-access", "groovy-scripted-error", "file-based-json-http-operation", "syslog-json-error", "groovy-scripted-http-operation"}...),
		}
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *logPublisherResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogPublisherResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanLogPublisher(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model logPublisherResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.DirectoryRESTAPIMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'directory_rest_api_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'directory_rest_api_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.IncludeRequestDetailsInResultMessages) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_request_details_in_result_messages' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_request_details_in_result_messages', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'detailed-http-operation', 'syslog-json-http-operation', 'json-access', 'console-json-http-operation', 'admin-alert-access', 'console-json-access', 'file-based-access', 'file-based-json-http-operation', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeRequestDetailsInSearchEntryMessages) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_request_details_in_search_entry_messages' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_request_details_in_search_entry_messages', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.LogRequests) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_requests' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_requests', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'syslog-json-http-operation', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.IncludeRequesterDN) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_requester_dn' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_requester_dn', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-access', 'syslog-json-audit', 'operation-timing-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'admin-alert-access', 'file-based-audit', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.ConnectionCriteria) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'connection_criteria' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'connection_criteria', the 'type' attribute must be one of ['syslog-json-audit', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.RequestCriteria) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'request_criteria' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'request_criteria', the 'type' attribute must be one of ['syslog-json-audit', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ServerHostName) && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'server_host_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_host_name', the 'type' attribute must be one of ['syslog-based-error', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.AccessTokenValidatorMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'access_token_validator_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'access_token_validator_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.LogResponseHeaders) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_response_headers' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_response_headers', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.SuppressedRequestHeaderName) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'suppressed_request_header_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'suppressed_request_header_name', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.SearchReferenceCriteria) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'search_reference_criteria' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'search_reference_criteria', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.LogResponseCookieNames) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_response_cookie_names' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_response_cookie_names', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.SuppressInternalOperations) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'suppress_internal_operations' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'suppress_internal_operations', the 'type' attribute must be one of ['syslog-json-audit', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeResponseControls) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_response_controls' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_response_controls', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'json-access', 'console-json-audit', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.LogRequestHeaders) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_request_headers' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_request_headers', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.MaxStringLength) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'max_string_length' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_string_length', the 'type' attribute must be one of ['syslog-json-access', 'operation-timing-access', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'admin-alert-access', 'file-based-trace', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.EncryptLog) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'encrypt_log' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'encrypt_log', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.WriteMultiLineMessages) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'write_multi_line_messages' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'write_multi_line_messages', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'syslog-json-http-operation', 'json-access', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'json-error', 'console-json-error', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.LogIntermediateResponses) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_intermediate_responses' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_intermediate_responses', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ExtensionClass) && model.Type.ValueString() != "third-party-error" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "third-party-http-operation" && model.Type.ValueString() != "third-party-access" {
		resp.Diagnostics.AddError("Attribute 'extension_class' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_class', the 'type' attribute must be one of ['third-party-error', 'third-party-file-based-access', 'third-party-file-based-error', 'third-party-http-operation', 'third-party-access']")
	}
	if internaltypes.IsDefined(model.GenerifyMessageStringsWhenPossible) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'generify_message_strings_when_possible' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'generify_message_strings_when_possible', the 'type' attribute must be one of ['syslog-json-access', 'file-based-error', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'json-error', 'file-based-access', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.ScimMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'scim_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'scim_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.IncludeResultCodeNames) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_result_code_names' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_result_code_names', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.OverrideSeverity) && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "jdbc-based-error" && model.Type.ValueString() != "groovy-scripted-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'override_severity' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'override_severity', the 'type' attribute must be one of ['syslog-based-error', 'file-based-error', 'third-party-error', 'third-party-file-based-error', 'json-error', 'jdbc-based-error', 'groovy-scripted-error', 'console-json-error', 'syslog-text-error', 'groovy-scripted-file-based-error', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.Server) && model.Type.ValueString() != "jdbc-based-error" && model.Type.ValueString() != "jdbc-based-access" {
		resp.Diagnostics.AddError("Attribute 'server' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server', the 'type' attribute must be one of ['jdbc-based-error', 'jdbc-based-access']")
	}
	if internaltypes.IsDefined(model.MinIncludedPhaseTimeNanos) && model.Type.ValueString() != "operation-timing-access" {
		resp.Diagnostics.AddError("Attribute 'min_included_phase_time_nanos' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'min_included_phase_time_nanos', the 'type' attribute must be one of ['operation-timing-access']")
	}
	if internaltypes.IsDefined(model.DefaultOmitMethodReturnValue) && model.Type.ValueString() != "file-based-debug" {
		resp.Diagnostics.AddError("Attribute 'default_omit_method_return_value' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_omit_method_return_value', the 'type' attribute must be one of ['file-based-debug']")
	}
	if internaltypes.IsDefined(model.CorrelateRequestsAndResults) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'correlate_requests_and_results' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'correlate_requests_and_results', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.SyslogMessageApplicationName) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'syslog_message_application_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'syslog_message_application_name', the 'type' attribute must be one of ['syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'syslog-json-http-operation', 'syslog-text-error', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.ObscureAttribute) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "file-based-audit" {
		resp.Diagnostics.AddError("Attribute 'obscure_attribute' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'obscure_attribute', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'debug-access', 'console-json-audit', 'file-based-audit']")
	}
	if internaltypes.IsDefined(model.SignLog) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'sign_log' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'sign_log', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.BufferSize) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'buffer_size' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'buffer_size', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.HttpMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'http_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'http_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.LogRequestCookieNames) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_request_cookie_names' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_request_cookie_names', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.LogRequestProtocol) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_request_protocol' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_request_protocol', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.SuppressedRequestParameterName) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'suppressed_request_parameter_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'suppressed_request_parameter_name', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.DefaultDebugCategory) && model.Type.ValueString() != "file-based-debug" {
		resp.Diagnostics.AddError("Attribute 'default_debug_category' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_debug_category', the 'type' attribute must be one of ['file-based-debug']")
	}
	if internaltypes.IsDefined(model.ExtensionArgument) && model.Type.ValueString() != "third-party-error" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "third-party-http-operation" && model.Type.ValueString() != "third-party-access" {
		resp.Diagnostics.AddError("Attribute 'extension_argument' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_argument', the 'type' attribute must be one of ['third-party-error', 'third-party-file-based-access', 'third-party-file-based-error', 'third-party-http-operation', 'third-party-access']")
	}
	if internaltypes.IsDefined(model.LogSecurityNegotiation) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_security_negotiation' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_security_negotiation', the 'type' attribute must be one of ['syslog-json-audit', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ScriptArgument) && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "groovy-scripted-error" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "groovy-scripted-http-operation" {
		resp.Diagnostics.AddError("Attribute 'script_argument' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_argument', the 'type' attribute must be one of ['groovy-scripted-access', 'groovy-scripted-file-based-access', 'groovy-scripted-error', 'groovy-scripted-file-based-error', 'groovy-scripted-http-operation']")
	}
	if internaltypes.IsDefined(model.TimeInterval) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'time_interval' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'time_interval', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.LogFile) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_file' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_file', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.IncludeAddAttributeNames) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_add_attribute_names' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_add_attribute_names', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeExtendedSearchRequestDetails) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_extended_search_request_details' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_extended_search_request_details', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.ExcludeAttribute) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "file-based-audit" {
		resp.Diagnostics.AddError("Attribute 'exclude_attribute' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'exclude_attribute', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'console-json-audit', 'file-based-audit']")
	}
	if internaltypes.IsDefined(model.LogRedirectURI) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_redirect_uri' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_redirect_uri', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.UseReversibleForm) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "file-based-audit" {
		resp.Diagnostics.AddError("Attribute 'use_reversible_form' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_reversible_form', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'console-json-audit', 'file-based-audit']")
	}
	if internaltypes.IsDefined(model.ResultCriteria) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'result_criteria' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'result_criteria', the 'type' attribute must be one of ['syslog-json-audit', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ExcludePathPattern) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'exclude_path_pattern' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'exclude_path_pattern', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.EncryptionSettingsDefinitionID) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'encryption_settings_definition_id' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'encryption_settings_definition_id', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.MinIncludedOperationProcessingTime) && model.Type.ValueString() != "operation-timing-access" {
		resp.Diagnostics.AddError("Attribute 'min_included_operation_processing_time' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'min_included_operation_processing_time', the 'type' attribute must be one of ['operation-timing-access']")
	}
	if internaltypes.IsDefined(model.ConsentMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'consent_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'consent_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.IdTokenValidatorMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'id_token_validator_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'id_token_validator_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.DebugACIEnabled) && model.Type.ValueString() != "debug-access" {
		resp.Diagnostics.AddError("Attribute 'debug_aci_enabled' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'debug_aci_enabled', the 'type' attribute must be one of ['debug-access']")
	}
	if internaltypes.IsDefined(model.IncludeRequestControls) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_request_controls' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_request_controls', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'json-access', 'console-json-audit', 'admin-alert-access', 'file-based-audit', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.LogRequestAuthorizationType) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_request_authorization_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_request_authorization_type', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.DebugMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'debug_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'debug_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.LogFieldMapping) && model.Type.ValueString() != "jdbc-based-error" && model.Type.ValueString() != "jdbc-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_field_mapping' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_field_mapping', the 'type' attribute must be one of ['jdbc-based-error', 'jdbc-based-access']")
	}
	if internaltypes.IsDefined(model.SuppressedResponseHeaderName) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'suppressed_response_header_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'suppressed_response_header_name', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.DefaultDebugLevel) && model.Type.ValueString() != "file-based-debug" {
		resp.Diagnostics.AddError("Attribute 'default_debug_level' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_debug_level', the 'type' attribute must be one of ['file-based-debug']")
	}
	if internaltypes.IsDefined(model.LogFieldBehavior) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_field_behavior' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_field_behavior', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.SuppressReplicationOperations) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'suppress_replication_operations' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'suppress_replication_operations', the 'type' attribute must be one of ['syslog-json-audit', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'third-party-access', 'admin-alert-access', 'file-based-audit', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.LogClientCertificates) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_client_certificates' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_client_certificates', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeIntermediateClientRequestControl) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "file-based-audit" {
		resp.Diagnostics.AddError("Attribute 'include_intermediate_client_request_control' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_intermediate_client_request_control', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'console-json-audit', 'file-based-audit']")
	}
	if internaltypes.IsDefined(model.SyslogFacility) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'syslog_facility' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'syslog_facility', the 'type' attribute must be one of ['syslog-json-access', 'syslog-json-audit', 'syslog-based-error', 'syslog-text-access', 'syslog-json-http-operation', 'syslog-text-error', 'syslog-based-access', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.IncludePathPattern) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'include_path_pattern' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_path_pattern', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.LogAssuranceCompleted) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_assurance_completed' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_assurance_completed', the 'type' attribute must be one of ['syslog-json-access', 'debug-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeRequestDetailsInSearchReferenceMessages) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_request_details_in_search_reference_messages' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_request_details_in_search_reference_messages', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeOperationPurposeRequestControl) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "file-based-audit" {
		resp.Diagnostics.AddError("Attribute 'include_operation_purpose_request_control' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_operation_purpose_request_control', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'console-json-audit', 'file-based-audit']")
	}
	if internaltypes.IsDefined(model.IncludeStartupID) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'include_startup_id' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_startup_id', the 'type' attribute must be one of ['syslog-json-audit', 'operation-timing-access', 'syslog-json-http-operation', 'admin-alert-access', 'file-based-audit', 'json-error', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'file-based-json-http-operation', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.DefaultOmitMethodEntryArguments) && model.Type.ValueString() != "file-based-debug" {
		resp.Diagnostics.AddError("Attribute 'default_omit_method_entry_arguments' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_omit_method_entry_arguments', the 'type' attribute must be one of ['file-based-debug']")
	}
	if internaltypes.IsDefined(model.IncludeRequesterIPAddress) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_requester_ip_address' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_requester_ip_address', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-access', 'syslog-json-audit', 'operation-timing-access', 'syslog-text-access', 'json-access', 'console-json-audit', 'admin-alert-access', 'file-based-audit', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeSearchEntryAttributeNames) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_search_entry_attribute_names' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_search_entry_attribute_names', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.DefaultSeverity) && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "jdbc-based-error" && model.Type.ValueString() != "groovy-scripted-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'default_severity' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_severity', the 'type' attribute must be one of ['syslog-based-error', 'file-based-error', 'third-party-error', 'third-party-file-based-error', 'json-error', 'jdbc-based-error', 'groovy-scripted-error', 'console-json-error', 'syslog-text-error', 'groovy-scripted-file-based-error', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.LogResults) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_results' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_results', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'syslog-json-http-operation', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.RotationListener) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'rotation_listener' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'rotation_listener', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.DefaultThrowableStackFrames) && model.Type.ValueString() != "file-based-debug" {
		resp.Diagnostics.AddError("Attribute 'default_throwable_stack_frames' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_throwable_stack_frames', the 'type' attribute must be one of ['file-based-debug']")
	}
	if internaltypes.IsDefined(model.TimestampPrecision) && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-text-error" {
		resp.Diagnostics.AddError("Attribute 'timestamp_precision' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'timestamp_precision', the 'type' attribute must be one of ['file-based-debug', 'file-based-error', 'syslog-text-access', 'file-based-audit', 'file-based-access', 'syslog-text-error']")
	}
	if internaltypes.IsDefined(model.SyslogMessageHostName) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'syslog_message_host_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'syslog_message_host_name', the 'type' attribute must be one of ['syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'syslog-json-http-operation', 'syslog-text-error', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.LogSearchReferences) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_search_references' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_search_references', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.OutputLocation) && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "console-json-error" {
		resp.Diagnostics.AddError("Attribute 'output_location' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'output_location', the 'type' attribute must be one of ['console-json-audit', 'console-json-http-operation', 'console-json-access', 'console-json-error']")
	}
	if internaltypes.IsDefined(model.LogConnects) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_connects' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_connects', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ScriptClass) && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "groovy-scripted-error" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "groovy-scripted-http-operation" {
		resp.Diagnostics.AddError("Attribute 'script_class' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'script_class', the 'type' attribute must be one of ['groovy-scripted-access', 'groovy-scripted-file-based-access', 'groovy-scripted-error', 'groovy-scripted-file-based-error', 'groovy-scripted-http-operation']")
	}
	if internaltypes.IsDefined(model.SyslogSeverity) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'syslog_severity' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'syslog_severity', the 'type' attribute must be one of ['syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'syslog-json-http-operation', 'syslog-text-error', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.CompressionMechanism) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'compression_mechanism' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'compression_mechanism', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.RotationPolicy) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'rotation_policy' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'rotation_policy', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.LogRequestParameters) && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_request_parameters' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_request_parameters', the 'type' attribute must be one of ['detailed-http-operation', 'syslog-json-http-operation', 'console-json-http-operation', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.SearchEntryCriteria) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'search_entry_criteria' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'search_entry_criteria', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ServerPort) && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'server_port' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'server_port', the 'type' attribute must be one of ['syslog-based-error', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.LogTableName) && model.Type.ValueString() != "jdbc-based-error" && model.Type.ValueString() != "jdbc-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_table_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_table_name', the 'type' attribute must be one of ['jdbc-based-error', 'jdbc-based-access']")
	}
	if internaltypes.IsDefined(model.LogSearchEntries) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_search_entries' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_search_entries', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.ObscureSensitiveContent) && model.Type.ValueString() != "debug-access" {
		resp.Diagnostics.AddError("Attribute 'obscure_sensitive_content' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'obscure_sensitive_content', the 'type' attribute must be one of ['debug-access']")
	}
	if internaltypes.IsDefined(model.SoftDeleteEntryAuditBehavior) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "file-based-audit" {
		resp.Diagnostics.AddError("Attribute 'soft_delete_entry_audit_behavior' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'soft_delete_entry_audit_behavior', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-audit', 'console-json-audit', 'file-based-audit']")
	}
	if internaltypes.IsDefined(model.IncludeThreadID) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'include_thread_id' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_thread_id', the 'type' attribute must be one of ['syslog-json-audit', 'operation-timing-access', 'syslog-json-http-operation', 'admin-alert-access', 'file-based-audit', 'json-error', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'file-based-json-http-operation', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.Asynchronous) && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'asynchronous' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'asynchronous', the 'type' attribute must be one of ['syslog-based-error', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'admin-alert-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'syslog-based-access', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.IncludeModifyAttributeNames) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_modify_attribute_names' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_modify_attribute_names', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.RetentionPolicy) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'retention_policy' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'retention_policy', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.DefaultIncludeThrowableCause) && model.Type.ValueString() != "file-based-debug" {
		resp.Diagnostics.AddError("Attribute 'default_include_throwable_cause' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'default_include_throwable_cause', the 'type' attribute must be one of ['file-based-debug']")
	}
	if internaltypes.IsDefined(model.IncludeReplicationChangeID) && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_replication_change_id' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_replication_change_id', the 'type' attribute must be one of ['file-based-json-audit', 'syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'json-access', 'console-json-audit', 'admin-alert-access', 'file-based-audit', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.IncludeProductName) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'include_product_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_product_name', the 'type' attribute must be one of ['syslog-json-audit', 'operation-timing-access', 'syslog-json-http-operation', 'admin-alert-access', 'file-based-audit', 'json-error', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'file-based-json-http-operation', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.IncludeInstanceName) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "console-json-error" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-audit" && model.Type.ValueString() != "console-json-http-operation" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'include_instance_name' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_instance_name', the 'type' attribute must be one of ['syslog-json-audit', 'operation-timing-access', 'syslog-json-http-operation', 'admin-alert-access', 'file-based-audit', 'json-error', 'console-json-error', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'syslog-json-access', 'file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'console-json-audit', 'console-json-http-operation', 'console-json-access', 'file-based-access', 'file-based-json-http-operation', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.ExtensionMessageType) && model.Type.ValueString() != "file-based-trace" {
		resp.Diagnostics.AddError("Attribute 'extension_message_type' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'extension_message_type', the 'type' attribute must be one of ['file-based-trace']")
	}
	if internaltypes.IsDefined(model.SyslogExternalServer) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'syslog_external_server' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'syslog_external_server', the 'type' attribute must be one of ['syslog-json-access', 'syslog-json-audit', 'syslog-text-access', 'syslog-json-http-operation', 'syslog-text-error', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.LogDisconnects) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "third-party-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "groovy-scripted-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" {
		resp.Diagnostics.AddError("Attribute 'log_disconnects' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_disconnects', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'third-party-access', 'admin-alert-access', 'jdbc-based-access', 'groovy-scripted-file-based-access', 'syslog-based-access', 'syslog-json-access', 'groovy-scripted-access', 'syslog-text-access', 'json-access', 'console-json-access', 'file-based-access']")
	}
	if internaltypes.IsDefined(model.AutoFlush) && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'auto_flush' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'auto_flush', the 'type' attribute must be one of ['syslog-based-error', 'third-party-file-based-access', 'debug-access', 'operation-timing-access', 'admin-alert-access', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'syslog-based-access', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.QueueSize) && model.Type.ValueString() != "syslog-json-audit" && model.Type.ValueString() != "syslog-based-error" && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "jdbc-based-error" && model.Type.ValueString() != "jdbc-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "syslog-text-error" && model.Type.ValueString() != "syslog-based-access" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "syslog-json-http-operation" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" && model.Type.ValueString() != "syslog-json-error" {
		resp.Diagnostics.AddError("Attribute 'queue_size' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'queue_size', the 'type' attribute must be one of ['syslog-json-audit', 'syslog-based-error', 'third-party-file-based-access', 'operation-timing-access', 'admin-alert-access', 'file-based-trace', 'jdbc-based-error', 'jdbc-based-access', 'common-log-file-http-operation', 'syslog-text-error', 'syslog-based-access', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'syslog-text-access', 'detailed-http-operation', 'json-access', 'debug-access', 'syslog-json-http-operation', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'groovy-scripted-file-based-error', 'syslog-json-access', 'third-party-file-based-error', 'file-based-access', 'file-based-json-http-operation', 'syslog-json-error']")
	}
	if internaltypes.IsDefined(model.IncludeRequestDetailsInIntermediateResponseMessages) && model.Type.ValueString() != "syslog-json-access" && model.Type.ValueString() != "syslog-text-access" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "admin-alert-access" && model.Type.ValueString() != "console-json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "syslog-based-access" {
		resp.Diagnostics.AddError("Attribute 'include_request_details_in_intermediate_response_messages' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'include_request_details_in_intermediate_response_messages', the 'type' attribute must be one of ['syslog-json-access', 'syslog-text-access', 'json-access', 'admin-alert-access', 'console-json-access', 'file-based-access', 'syslog-based-access']")
	}
	if internaltypes.IsDefined(model.LogFilePermissions) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'log_file_permissions' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'log_file_permissions', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
	if internaltypes.IsDefined(model.Append) && model.Type.ValueString() != "third-party-file-based-access" && model.Type.ValueString() != "debug-access" && model.Type.ValueString() != "operation-timing-access" && model.Type.ValueString() != "file-based-trace" && model.Type.ValueString() != "file-based-audit" && model.Type.ValueString() != "json-error" && model.Type.ValueString() != "groovy-scripted-file-based-access" && model.Type.ValueString() != "common-log-file-http-operation" && model.Type.ValueString() != "groovy-scripted-file-based-error" && model.Type.ValueString() != "file-based-json-audit" && model.Type.ValueString() != "file-based-debug" && model.Type.ValueString() != "file-based-error" && model.Type.ValueString() != "third-party-file-based-error" && model.Type.ValueString() != "detailed-http-operation" && model.Type.ValueString() != "json-access" && model.Type.ValueString() != "file-based-access" && model.Type.ValueString() != "file-based-json-http-operation" {
		resp.Diagnostics.AddError("Attribute 'append' not supported by pingdirectory_log_publisher resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'append', the 'type' attribute must be one of ['third-party-file-based-access', 'debug-access', 'operation-timing-access', 'file-based-trace', 'file-based-audit', 'json-error', 'groovy-scripted-file-based-access', 'common-log-file-http-operation', 'groovy-scripted-file-based-error', 'file-based-json-audit', 'file-based-debug', 'file-based-error', 'third-party-file-based-error', 'detailed-http-operation', 'json-access', 'file-based-access', 'file-based-json-http-operation']")
	}
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
	if internaltypes.IsDefined(plan.SyslogFacility) {
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
	if internaltypes.IsDefined(plan.SyslogFacility) {
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

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateLogPublisherNilSets(ctx context.Context, model *logPublisherResourceModel) {
	if model.IdTokenValidatorMessageType.ElementType(ctx) == nil {
		model.IdTokenValidatorMessageType = types.SetNull(types.StringType)
	}
	if model.SuppressedResponseHeaderName.ElementType(ctx) == nil {
		model.SuppressedResponseHeaderName = types.SetNull(types.StringType)
	}
	if model.DefaultDebugCategory.ElementType(ctx) == nil {
		model.DefaultDebugCategory = types.SetNull(types.StringType)
	}
	if model.ScimMessageType.ElementType(ctx) == nil {
		model.ScimMessageType = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.RotationPolicy.ElementType(ctx) == nil {
		model.RotationPolicy = types.SetNull(types.StringType)
	}
	if model.RetentionPolicy.ElementType(ctx) == nil {
		model.RetentionPolicy = types.SetNull(types.StringType)
	}
	if model.SuppressedRequestHeaderName.ElementType(ctx) == nil {
		model.SuppressedRequestHeaderName = types.SetNull(types.StringType)
	}
	if model.IncludePathPattern.ElementType(ctx) == nil {
		model.IncludePathPattern = types.SetNull(types.StringType)
	}
	if model.ExcludeAttribute.ElementType(ctx) == nil {
		model.ExcludeAttribute = types.SetNull(types.StringType)
	}
	if model.HttpMessageType.ElementType(ctx) == nil {
		model.HttpMessageType = types.SetNull(types.StringType)
	}
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.DefaultSeverity.ElementType(ctx) == nil {
		model.DefaultSeverity = types.SetNull(types.StringType)
	}
	if model.RotationListener.ElementType(ctx) == nil {
		model.RotationListener = types.SetNull(types.StringType)
	}
	if model.ExcludePathPattern.ElementType(ctx) == nil {
		model.ExcludePathPattern = types.SetNull(types.StringType)
	}
	if model.AccessTokenValidatorMessageType.ElementType(ctx) == nil {
		model.AccessTokenValidatorMessageType = types.SetNull(types.StringType)
	}
	if model.ExtensionMessageType.ElementType(ctx) == nil {
		model.ExtensionMessageType = types.SetNull(types.StringType)
	}
	if model.SuppressedRequestParameterName.ElementType(ctx) == nil {
		model.SuppressedRequestParameterName = types.SetNull(types.StringType)
	}
	if model.ConsentMessageType.ElementType(ctx) == nil {
		model.ConsentMessageType = types.SetNull(types.StringType)
	}
	if model.OverrideSeverity.ElementType(ctx) == nil {
		model.OverrideSeverity = types.SetNull(types.StringType)
	}
	if model.DirectoryRESTAPIMessageType.ElementType(ctx) == nil {
		model.DirectoryRESTAPIMessageType = types.SetNull(types.StringType)
	}
	if model.ObscureAttribute.ElementType(ctx) == nil {
		model.ObscureAttribute = types.SetNull(types.StringType)
	}
	if model.SyslogExternalServer.ElementType(ctx) == nil {
		model.SyslogExternalServer = types.SetNull(types.StringType)
	}
	if model.DebugMessageType.ElementType(ctx) == nil {
		model.DebugMessageType = types.SetNull(types.StringType)
	}
}

// Read a SyslogJsonAuditLogPublisherResponse object into the model struct
func readSyslogJsonAuditLogPublisherResponse(ctx context.Context, r *client.SyslogJsonAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-audit")
	state.Id = types.StringValue(r.Id)
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = types.StringValue(r.SyslogSeverity.String())
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), internaltypes.IsEmptyString(expectedValues.SoftDeleteEntryAuditBehavior))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogBasedErrorLogPublisherResponse object into the model struct
func readSyslogBasedErrorLogPublisherResponse(ctx context.Context, r *client.SyslogBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-based-error")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ThirdPartyFileBasedAccessLogPublisherResponse object into the model struct
func readThirdPartyFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.ThirdPartyFileBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-file-based-access")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a OperationTimingAccessLogPublisherResponse object into the model struct
func readOperationTimingAccessLogPublisherResponse(ctx context.Context, r *client.OperationTimingAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("operation-timing-access")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
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
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ThirdPartyHttpOperationLogPublisherResponse object into the model struct
func readThirdPartyHttpOperationLogPublisherResponse(ctx context.Context, r *client.ThirdPartyHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-http-operation")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a AdminAlertAccessLogPublisherResponse object into the model struct
func readAdminAlertAccessLogPublisherResponse(ctx context.Context, r *client.AdminAlertAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("admin-alert-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedTraceLogPublisherResponse object into the model struct
func readFileBasedTraceLogPublisherResponse(ctx context.Context, r *client.FileBasedTraceLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-trace")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a JdbcBasedErrorLogPublisherResponse object into the model struct
func readJdbcBasedErrorLogPublisherResponse(ctx context.Context, r *client.JdbcBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jdbc-based-error")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a JdbcBasedAccessLogPublisherResponse object into the model struct
func readJdbcBasedAccessLogPublisherResponse(ctx context.Context, r *client.JdbcBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jdbc-based-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a CommonLogFileHttpOperationLogPublisherResponse object into the model struct
func readCommonLogFileHttpOperationLogPublisherResponse(ctx context.Context, r *client.CommonLogFileHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("common-log-file-http-operation")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ConsoleJsonErrorLogPublisherResponse object into the model struct
func readConsoleJsonErrorLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-error")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), internaltypes.IsEmptyString(expectedValues.OutputLocation))
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogTextErrorLogPublisherResponse object into the model struct
func readSyslogTextErrorLogPublisherResponse(ctx context.Context, r *client.SyslogTextErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-text-error")
	state.Id = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), internaltypes.IsEmptyString(expectedValues.SyslogSeverity))
	state.SyslogMessageHostName = internaltypes.StringTypeOrNil(r.SyslogMessageHostName, internaltypes.IsEmptyString(expectedValues.SyslogMessageHostName))
	state.SyslogMessageApplicationName = internaltypes.StringTypeOrNil(r.SyslogMessageApplicationName, internaltypes.IsEmptyString(expectedValues.SyslogMessageApplicationName))
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.GenerifyMessageStringsWhenPossible = internaltypes.BoolTypeOrNil(r.GenerifyMessageStringsWhenPossible)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogBasedAccessLogPublisherResponse object into the model struct
func readSyslogBasedAccessLogPublisherResponse(ctx context.Context, r *client.SyslogBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-based-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedJsonAuditLogPublisherResponse object into the model struct
func readFileBasedJsonAuditLogPublisherResponse(ctx context.Context, r *client.FileBasedJsonAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-json-audit")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherFileBasedJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), internaltypes.IsEmptyString(expectedValues.SoftDeleteEntryAuditBehavior))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedDebugLogPublisherResponse object into the model struct
func readFileBasedDebugLogPublisherResponse(ctx context.Context, r *client.FileBasedDebugLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-debug")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedErrorLogPublisherResponse object into the model struct
func readFileBasedErrorLogPublisherResponse(ctx context.Context, r *client.FileBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-error")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
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
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ThirdPartyErrorLogPublisherResponse object into the model struct
func readThirdPartyErrorLogPublisherResponse(ctx context.Context, r *client.ThirdPartyErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-error")
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogTextAccessLogPublisherResponse object into the model struct
func readSyslogTextAccessLogPublisherResponse(ctx context.Context, r *client.SyslogTextAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-text-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a DetailedHttpOperationLogPublisherResponse object into the model struct
func readDetailedHttpOperationLogPublisherResponse(ctx context.Context, r *client.DetailedHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("detailed-http-operation")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
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
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), internaltypes.IsEmptyString(expectedValues.LogRequestHeaders))
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), internaltypes.IsEmptyString(expectedValues.LogResponseHeaders))
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), internaltypes.IsEmptyString(expectedValues.LogRequestParameters))
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.MaxStringLength = internaltypes.Int64TypeOrNil(r.MaxStringLength)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a JsonAccessLogPublisherResponse object into the model struct
func readJsonAccessLogPublisherResponse(ctx context.Context, r *client.JsonAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-access")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
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
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a DebugAccessLogPublisherResponse object into the model struct
func readDebugAccessLogPublisherResponse(ctx context.Context, r *client.DebugAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("debug-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ObscureSensitiveContent = internaltypes.BoolTypeOrNil(r.ObscureSensitiveContent)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.DebugACIEnabled = internaltypes.BoolTypeOrNil(r.DebugACIEnabled)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogJsonHttpOperationLogPublisherResponse object into the model struct
func readSyslogJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.SyslogJsonHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-http-operation")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), internaltypes.IsEmptyString(expectedValues.LogRequestHeaders))
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), internaltypes.IsEmptyString(expectedValues.LogResponseHeaders))
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), internaltypes.IsEmptyString(expectedValues.LogRequestParameters))
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ThirdPartyAccessLogPublisherResponse object into the model struct
func readThirdPartyAccessLogPublisherResponse(ctx context.Context, r *client.ThirdPartyAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedAuditLogPublisherResponse object into the model struct
func readFileBasedAuditLogPublisherResponse(ctx context.Context, r *client.FileBasedAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-audit")
	state.Id = types.StringValue(r.Id)
	state.SuppressInternalOperations = internaltypes.BoolTypeOrNil(r.SuppressInternalOperations)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
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
		client.StringPointerEnumlogPublisherFileBasedAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), internaltypes.IsEmptyString(expectedValues.SoftDeleteEntryAuditBehavior))
	state.IncludeRequestControls = internaltypes.BoolTypeOrNil(r.IncludeRequestControls)
	state.IncludeOperationPurposeRequestControl = internaltypes.BoolTypeOrNil(r.IncludeOperationPurposeRequestControl)
	state.IncludeIntermediateClientRequestControl = internaltypes.BoolTypeOrNil(r.IncludeIntermediateClientRequestControl)
	state.ObscureAttribute = internaltypes.GetStringSet(r.ObscureAttribute)
	state.ExcludeAttribute = internaltypes.GetStringSet(r.ExcludeAttribute)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
	state.LogSecurityNegotiation = internaltypes.BoolTypeOrNil(r.LogSecurityNegotiation)
	state.LogIntermediateResponses = internaltypes.BoolTypeOrNil(r.LogIntermediateResponses)
	state.SuppressReplicationOperations = internaltypes.BoolTypeOrNil(r.SuppressReplicationOperations)
	state.ConnectionCriteria = internaltypes.StringTypeOrNil(r.ConnectionCriteria, internaltypes.IsEmptyString(expectedValues.ConnectionCriteria))
	state.RequestCriteria = internaltypes.StringTypeOrNil(r.RequestCriteria, internaltypes.IsEmptyString(expectedValues.RequestCriteria))
	state.ResultCriteria = internaltypes.StringTypeOrNil(r.ResultCriteria, internaltypes.IsEmptyString(expectedValues.ResultCriteria))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a JsonErrorLogPublisherResponse object into the model struct
func readJsonErrorLogPublisherResponse(ctx context.Context, r *client.JsonErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-error")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a GroovyScriptedFileBasedAccessLogPublisherResponse object into the model struct
func readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedFileBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-file-based-access")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a GroovyScriptedFileBasedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedFileBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-file-based-error")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogJsonAccessLogPublisherResponse object into the model struct
func readSyslogJsonAccessLogPublisherResponse(ctx context.Context, r *client.SyslogJsonAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a GroovyScriptedAccessLogPublisherResponse object into the model struct
func readGroovyScriptedAccessLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-access")
	state.Id = types.StringValue(r.Id)
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ThirdPartyFileBasedErrorLogPublisherResponse object into the model struct
func readThirdPartyFileBasedErrorLogPublisherResponse(ctx context.Context, r *client.ThirdPartyFileBasedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-file-based-error")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ConsoleJsonAuditLogPublisherResponse object into the model struct
func readConsoleJsonAuditLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonAuditLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-audit")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), internaltypes.IsEmptyString(expectedValues.OutputLocation))
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.UseReversibleForm = internaltypes.BoolTypeOrNil(r.UseReversibleForm)
	state.SoftDeleteEntryAuditBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherConsoleJsonAuditSoftDeleteEntryAuditBehaviorProp(r.SoftDeleteEntryAuditBehavior), internaltypes.IsEmptyString(expectedValues.SoftDeleteEntryAuditBehavior))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ConsoleJsonHttpOperationLogPublisherResponse object into the model struct
func readConsoleJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), internaltypes.IsEmptyString(expectedValues.OutputLocation))
	state.LogRequests = internaltypes.BoolTypeOrNil(r.LogRequests)
	state.LogResults = internaltypes.BoolTypeOrNil(r.LogResults)
	state.IncludeProductName = internaltypes.BoolTypeOrNil(r.IncludeProductName)
	state.IncludeInstanceName = internaltypes.BoolTypeOrNil(r.IncludeInstanceName)
	state.IncludeStartupID = internaltypes.BoolTypeOrNil(r.IncludeStartupID)
	state.IncludeThreadID = internaltypes.BoolTypeOrNil(r.IncludeThreadID)
	state.IncludeRequestDetailsInResultMessages = internaltypes.BoolTypeOrNil(r.IncludeRequestDetailsInResultMessages)
	state.LogRequestHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), internaltypes.IsEmptyString(expectedValues.LogRequestHeaders))
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), internaltypes.IsEmptyString(expectedValues.LogResponseHeaders))
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), internaltypes.IsEmptyString(expectedValues.LogRequestParameters))
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a ConsoleJsonAccessLogPublisherResponse object into the model struct
func readConsoleJsonAccessLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console-json-access")
	state.Id = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.OutputLocation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherOutputLocationProp(r.OutputLocation), internaltypes.IsEmptyString(expectedValues.OutputLocation))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedAccessLogPublisherResponse object into the model struct
func readFileBasedAccessLogPublisherResponse(ctx context.Context, r *client.FileBasedAccessLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-access")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
	config.CheckMismatchedPDFormattedAttributes("time_interval",
		expectedValues.TimeInterval, state.TimeInterval, diagnostics)
	state.TimestampPrecision = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherTimestampPrecisionProp(r.TimestampPrecision), internaltypes.IsEmptyString(expectedValues.TimestampPrecision))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a GroovyScriptedErrorLogPublisherResponse object into the model struct
func readGroovyScriptedErrorLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-error")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.OverrideSeverity = internaltypes.GetStringSet(r.OverrideSeverity)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a FileBasedJsonHttpOperationLogPublisherResponse object into the model struct
func readFileBasedJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.FileBasedJsonHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based-json-http-operation")
	state.Id = types.StringValue(r.Id)
	state.LogFile = types.StringValue(r.LogFile)
	state.LogFilePermissions = types.StringValue(r.LogFilePermissions)
	state.RotationPolicy = internaltypes.GetStringSet(r.RotationPolicy)
	state.RotationListener = internaltypes.GetStringSet(r.RotationListener)
	state.RetentionPolicy = internaltypes.GetStringSet(r.RetentionPolicy)
	state.CompressionMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherCompressionMechanismProp(r.CompressionMechanism), internaltypes.IsEmptyString(expectedValues.CompressionMechanism))
	state.SignLog = internaltypes.BoolTypeOrNil(r.SignLog)
	state.EncryptLog = internaltypes.BoolTypeOrNil(r.EncryptLog)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Append = internaltypes.BoolTypeOrNil(r.Append)
	state.Asynchronous = types.BoolValue(r.Asynchronous)
	state.AutoFlush = internaltypes.BoolTypeOrNil(r.AutoFlush)
	state.BufferSize = internaltypes.StringTypeOrNil(r.BufferSize, internaltypes.IsEmptyString(expectedValues.BufferSize))
	config.CheckMismatchedPDFormattedAttributes("buffer_size",
		expectedValues.BufferSize, state.BufferSize, diagnostics)
	state.QueueSize = internaltypes.Int64TypeOrNil(r.QueueSize)
	state.TimeInterval = internaltypes.StringTypeOrNil(r.TimeInterval, internaltypes.IsEmptyString(expectedValues.TimeInterval))
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
		client.StringPointerEnumlogPublisherLogRequestHeadersProp(r.LogRequestHeaders), internaltypes.IsEmptyString(expectedValues.LogRequestHeaders))
	state.SuppressedRequestHeaderName = internaltypes.GetStringSet(r.SuppressedRequestHeaderName)
	state.LogResponseHeaders = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogResponseHeadersProp(r.LogResponseHeaders), internaltypes.IsEmptyString(expectedValues.LogResponseHeaders))
	state.SuppressedResponseHeaderName = internaltypes.GetStringSet(r.SuppressedResponseHeaderName)
	state.LogRequestAuthorizationType = internaltypes.BoolTypeOrNil(r.LogRequestAuthorizationType)
	state.LogRequestCookieNames = internaltypes.BoolTypeOrNil(r.LogRequestCookieNames)
	state.LogResponseCookieNames = internaltypes.BoolTypeOrNil(r.LogResponseCookieNames)
	state.LogRequestParameters = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLogRequestParametersProp(r.LogRequestParameters), internaltypes.IsEmptyString(expectedValues.LogRequestParameters))
	state.SuppressedRequestParameterName = internaltypes.GetStringSet(r.SuppressedRequestParameterName)
	state.LogRequestProtocol = internaltypes.BoolTypeOrNil(r.LogRequestProtocol)
	state.LogRedirectURI = internaltypes.BoolTypeOrNil(r.LogRedirectURI)
	state.WriteMultiLineMessages = internaltypes.BoolTypeOrNil(r.WriteMultiLineMessages)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a SyslogJsonErrorLogPublisherResponse object into the model struct
func readSyslogJsonErrorLogPublisherResponse(ctx context.Context, r *client.SyslogJsonErrorLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("syslog-json-error")
	state.Id = types.StringValue(r.Id)
	state.DefaultSeverity = internaltypes.GetStringSet(
		client.StringSliceEnumlogPublisherDefaultSeverityProp(r.DefaultSeverity))
	state.SyslogExternalServer = internaltypes.GetStringSet(r.SyslogExternalServer)
	state.SyslogFacility = types.StringValue(r.SyslogFacility.String())
	state.SyslogSeverity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherSyslogSeverityProp(r.SyslogSeverity), internaltypes.IsEmptyString(expectedValues.SyslogSeverity))
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
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
}

// Read a GroovyScriptedHttpOperationLogPublisherResponse object into the model struct
func readGroovyScriptedHttpOperationLogPublisherResponse(ctx context.Context, r *client.GroovyScriptedHttpOperationLogPublisherResponse, state *logPublisherResourceModel, expectedValues *logPublisherResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted-http-operation")
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.LoggingErrorBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogPublisherLoggingErrorBehaviorProp(r.LoggingErrorBehavior), internaltypes.IsEmptyString(expectedValues.LoggingErrorBehavior))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateLogPublisherNilSets(ctx, state)
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
	addRequest := client.NewAddSyslogJsonAuditLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogBasedErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddThirdPartyFileBasedAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddOperationTimingAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddThirdPartyHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddAdminAlertAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedTraceLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddJdbcBasedErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddJdbcBasedAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddCommonLogFileHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogTextErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogBasedAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedJsonAuditLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedDebugLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddThirdPartyErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogTextAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddDetailedHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddJsonAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddDebugAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogJsonHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddThirdPartyAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedAuditLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddJsonErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddGroovyScriptedFileBasedAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddGroovyScriptedFileBasedErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogJsonAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddGroovyScriptedAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddThirdPartyFileBasedErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddConsoleJsonAuditLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddConsoleJsonHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedAccessLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddGroovyScriptedErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddFileBasedJsonHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddSyslogJsonErrorLogPublisherRequest(plan.Id.ValueString(),
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
	addRequest := client.NewAddGroovyScriptedHttpOperationLogPublisherRequest(plan.Id.ValueString(),
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
func (r *defaultLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan logPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
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
	if plan.Type.ValueString() == "syslog-json-audit" {
		readSyslogJsonAuditLogPublisherResponse(ctx, readResponse.SyslogJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-based-error" {
		readSyslogBasedErrorLogPublisherResponse(ctx, readResponse.SyslogBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party-file-based-access" {
		readThirdPartyFileBasedAccessLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "operation-timing-access" {
		readOperationTimingAccessLogPublisherResponse(ctx, readResponse.OperationTimingAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party-http-operation" {
		readThirdPartyHttpOperationLogPublisherResponse(ctx, readResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "admin-alert-access" {
		readAdminAlertAccessLogPublisherResponse(ctx, readResponse.AdminAlertAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-trace" {
		readFileBasedTraceLogPublisherResponse(ctx, readResponse.FileBasedTraceLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "jdbc-based-error" {
		readJdbcBasedErrorLogPublisherResponse(ctx, readResponse.JdbcBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "jdbc-based-access" {
		readJdbcBasedAccessLogPublisherResponse(ctx, readResponse.JdbcBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "common-log-file-http-operation" {
		readCommonLogFileHttpOperationLogPublisherResponse(ctx, readResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "console-json-error" {
		readConsoleJsonErrorLogPublisherResponse(ctx, readResponse.ConsoleJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-text-error" {
		readSyslogTextErrorLogPublisherResponse(ctx, readResponse.SyslogTextErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-based-access" {
		readSyslogBasedAccessLogPublisherResponse(ctx, readResponse.SyslogBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-json-audit" {
		readFileBasedJsonAuditLogPublisherResponse(ctx, readResponse.FileBasedJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-debug" {
		readFileBasedDebugLogPublisherResponse(ctx, readResponse.FileBasedDebugLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-error" {
		readFileBasedErrorLogPublisherResponse(ctx, readResponse.FileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party-error" {
		readThirdPartyErrorLogPublisherResponse(ctx, readResponse.ThirdPartyErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-text-access" {
		readSyslogTextAccessLogPublisherResponse(ctx, readResponse.SyslogTextAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "detailed-http-operation" {
		readDetailedHttpOperationLogPublisherResponse(ctx, readResponse.DetailedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "json-access" {
		readJsonAccessLogPublisherResponse(ctx, readResponse.JsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "debug-access" {
		readDebugAccessLogPublisherResponse(ctx, readResponse.DebugAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-json-http-operation" {
		readSyslogJsonHttpOperationLogPublisherResponse(ctx, readResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party-access" {
		readThirdPartyAccessLogPublisherResponse(ctx, readResponse.ThirdPartyAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-audit" {
		readFileBasedAuditLogPublisherResponse(ctx, readResponse.FileBasedAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "json-error" {
		readJsonErrorLogPublisherResponse(ctx, readResponse.JsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted-file-based-access" {
		readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, readResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted-file-based-error" {
		readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-json-access" {
		readSyslogJsonAccessLogPublisherResponse(ctx, readResponse.SyslogJsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted-access" {
		readGroovyScriptedAccessLogPublisherResponse(ctx, readResponse.GroovyScriptedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "third-party-file-based-error" {
		readThirdPartyFileBasedErrorLogPublisherResponse(ctx, readResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "console-json-audit" {
		readConsoleJsonAuditLogPublisherResponse(ctx, readResponse.ConsoleJsonAuditLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "console-json-http-operation" {
		readConsoleJsonHttpOperationLogPublisherResponse(ctx, readResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "console-json-access" {
		readConsoleJsonAccessLogPublisherResponse(ctx, readResponse.ConsoleJsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-access" {
		readFileBasedAccessLogPublisherResponse(ctx, readResponse.FileBasedAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted-error" {
		readGroovyScriptedErrorLogPublisherResponse(ctx, readResponse.GroovyScriptedErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "file-based-json-http-operation" {
		readFileBasedJsonHttpOperationLogPublisherResponse(ctx, readResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "syslog-json-error" {
		readSyslogJsonErrorLogPublisherResponse(ctx, readResponse.SyslogJsonErrorLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "groovy-scripted-http-operation" {
		readGroovyScriptedHttpOperationLogPublisherResponse(ctx, readResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
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
		if plan.Type.ValueString() == "syslog-json-audit" {
			readSyslogJsonAuditLogPublisherResponse(ctx, updateResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-based-error" {
			readSyslogBasedErrorLogPublisherResponse(ctx, updateResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-file-based-access" {
			readThirdPartyFileBasedAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "operation-timing-access" {
			readOperationTimingAccessLogPublisherResponse(ctx, updateResponse.OperationTimingAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-http-operation" {
			readThirdPartyHttpOperationLogPublisherResponse(ctx, updateResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "admin-alert-access" {
			readAdminAlertAccessLogPublisherResponse(ctx, updateResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-trace" {
			readFileBasedTraceLogPublisherResponse(ctx, updateResponse.FileBasedTraceLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "jdbc-based-error" {
			readJdbcBasedErrorLogPublisherResponse(ctx, updateResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "jdbc-based-access" {
			readJdbcBasedAccessLogPublisherResponse(ctx, updateResponse.JdbcBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "common-log-file-http-operation" {
			readCommonLogFileHttpOperationLogPublisherResponse(ctx, updateResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-error" {
			readConsoleJsonErrorLogPublisherResponse(ctx, updateResponse.ConsoleJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-text-error" {
			readSyslogTextErrorLogPublisherResponse(ctx, updateResponse.SyslogTextErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-based-access" {
			readSyslogBasedAccessLogPublisherResponse(ctx, updateResponse.SyslogBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-json-audit" {
			readFileBasedJsonAuditLogPublisherResponse(ctx, updateResponse.FileBasedJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-debug" {
			readFileBasedDebugLogPublisherResponse(ctx, updateResponse.FileBasedDebugLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-error" {
			readFileBasedErrorLogPublisherResponse(ctx, updateResponse.FileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-error" {
			readThirdPartyErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-text-access" {
			readSyslogTextAccessLogPublisherResponse(ctx, updateResponse.SyslogTextAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "detailed-http-operation" {
			readDetailedHttpOperationLogPublisherResponse(ctx, updateResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "json-access" {
			readJsonAccessLogPublisherResponse(ctx, updateResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "debug-access" {
			readDebugAccessLogPublisherResponse(ctx, updateResponse.DebugAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-json-http-operation" {
			readSyslogJsonHttpOperationLogPublisherResponse(ctx, updateResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-access" {
			readThirdPartyAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-audit" {
			readFileBasedAuditLogPublisherResponse(ctx, updateResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "json-error" {
			readJsonErrorLogPublisherResponse(ctx, updateResponse.JsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-file-based-access" {
			readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-file-based-error" {
			readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-json-access" {
			readSyslogJsonAccessLogPublisherResponse(ctx, updateResponse.SyslogJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-access" {
			readGroovyScriptedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-file-based-error" {
			readThirdPartyFileBasedErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-audit" {
			readConsoleJsonAuditLogPublisherResponse(ctx, updateResponse.ConsoleJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-http-operation" {
			readConsoleJsonHttpOperationLogPublisherResponse(ctx, updateResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-access" {
			readConsoleJsonAccessLogPublisherResponse(ctx, updateResponse.ConsoleJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-access" {
			readFileBasedAccessLogPublisherResponse(ctx, updateResponse.FileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-error" {
			readGroovyScriptedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-json-http-operation" {
			readFileBasedJsonHttpOperationLogPublisherResponse(ctx, updateResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-json-error" {
			readSyslogJsonErrorLogPublisherResponse(ctx, updateResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-http-operation" {
			readGroovyScriptedHttpOperationLogPublisherResponse(ctx, updateResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *logPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state logPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
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

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
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
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

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
		if plan.Type.ValueString() == "syslog-json-audit" {
			readSyslogJsonAuditLogPublisherResponse(ctx, updateResponse.SyslogJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-based-error" {
			readSyslogBasedErrorLogPublisherResponse(ctx, updateResponse.SyslogBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-file-based-access" {
			readThirdPartyFileBasedAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "operation-timing-access" {
			readOperationTimingAccessLogPublisherResponse(ctx, updateResponse.OperationTimingAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-http-operation" {
			readThirdPartyHttpOperationLogPublisherResponse(ctx, updateResponse.ThirdPartyHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "admin-alert-access" {
			readAdminAlertAccessLogPublisherResponse(ctx, updateResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-trace" {
			readFileBasedTraceLogPublisherResponse(ctx, updateResponse.FileBasedTraceLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "jdbc-based-error" {
			readJdbcBasedErrorLogPublisherResponse(ctx, updateResponse.JdbcBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "jdbc-based-access" {
			readJdbcBasedAccessLogPublisherResponse(ctx, updateResponse.JdbcBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "common-log-file-http-operation" {
			readCommonLogFileHttpOperationLogPublisherResponse(ctx, updateResponse.CommonLogFileHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-error" {
			readConsoleJsonErrorLogPublisherResponse(ctx, updateResponse.ConsoleJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-text-error" {
			readSyslogTextErrorLogPublisherResponse(ctx, updateResponse.SyslogTextErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-based-access" {
			readSyslogBasedAccessLogPublisherResponse(ctx, updateResponse.SyslogBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-json-audit" {
			readFileBasedJsonAuditLogPublisherResponse(ctx, updateResponse.FileBasedJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-debug" {
			readFileBasedDebugLogPublisherResponse(ctx, updateResponse.FileBasedDebugLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-error" {
			readFileBasedErrorLogPublisherResponse(ctx, updateResponse.FileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-error" {
			readThirdPartyErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-text-access" {
			readSyslogTextAccessLogPublisherResponse(ctx, updateResponse.SyslogTextAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "detailed-http-operation" {
			readDetailedHttpOperationLogPublisherResponse(ctx, updateResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "json-access" {
			readJsonAccessLogPublisherResponse(ctx, updateResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "debug-access" {
			readDebugAccessLogPublisherResponse(ctx, updateResponse.DebugAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-json-http-operation" {
			readSyslogJsonHttpOperationLogPublisherResponse(ctx, updateResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-access" {
			readThirdPartyAccessLogPublisherResponse(ctx, updateResponse.ThirdPartyAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-audit" {
			readFileBasedAuditLogPublisherResponse(ctx, updateResponse.FileBasedAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "json-error" {
			readJsonErrorLogPublisherResponse(ctx, updateResponse.JsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-file-based-access" {
			readGroovyScriptedFileBasedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-file-based-error" {
			readGroovyScriptedFileBasedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-json-access" {
			readSyslogJsonAccessLogPublisherResponse(ctx, updateResponse.SyslogJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-access" {
			readGroovyScriptedAccessLogPublisherResponse(ctx, updateResponse.GroovyScriptedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "third-party-file-based-error" {
			readThirdPartyFileBasedErrorLogPublisherResponse(ctx, updateResponse.ThirdPartyFileBasedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-audit" {
			readConsoleJsonAuditLogPublisherResponse(ctx, updateResponse.ConsoleJsonAuditLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-http-operation" {
			readConsoleJsonHttpOperationLogPublisherResponse(ctx, updateResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "console-json-access" {
			readConsoleJsonAccessLogPublisherResponse(ctx, updateResponse.ConsoleJsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-access" {
			readFileBasedAccessLogPublisherResponse(ctx, updateResponse.FileBasedAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-error" {
			readGroovyScriptedErrorLogPublisherResponse(ctx, updateResponse.GroovyScriptedErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "file-based-json-http-operation" {
			readFileBasedJsonHttpOperationLogPublisherResponse(ctx, updateResponse.FileBasedJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "syslog-json-error" {
			readSyslogJsonErrorLogPublisherResponse(ctx, updateResponse.SyslogJsonErrorLogPublisherResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "groovy-scripted-http-operation" {
			readGroovyScriptedHttpOperationLogPublisherResponse(ctx, updateResponse.GroovyScriptedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
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
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
