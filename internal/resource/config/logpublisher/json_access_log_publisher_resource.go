package logpublisher

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &jsonAccessLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &jsonAccessLogPublisherResource{}
	_ resource.ResourceWithImportState = &jsonAccessLogPublisherResource{}
	_ resource.Resource                = &defaultJsonAccessLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultJsonAccessLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultJsonAccessLogPublisherResource{}
)

// Create a Json Access Log Publisher resource
func NewJsonAccessLogPublisherResource() resource.Resource {
	return &jsonAccessLogPublisherResource{}
}

func NewDefaultJsonAccessLogPublisherResource() resource.Resource {
	return &defaultJsonAccessLogPublisherResource{}
}

// jsonAccessLogPublisherResource is the resource implementation.
type jsonAccessLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultJsonAccessLogPublisherResource is the resource implementation.
type defaultJsonAccessLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *jsonAccessLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_json_access_log_publisher"
}

func (r *defaultJsonAccessLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_json_access_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *jsonAccessLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultJsonAccessLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type jsonAccessLogPublisherResourceModel struct {
	Id                                                  types.String `tfsdk:"id"`
	LastUpdated                                         types.String `tfsdk:"last_updated"`
	Notifications                                       types.Set    `tfsdk:"notifications"`
	RequiredActions                                     types.Set    `tfsdk:"required_actions"`
	LogFile                                             types.String `tfsdk:"log_file"`
	LogFilePermissions                                  types.String `tfsdk:"log_file_permissions"`
	RotationPolicy                                      types.Set    `tfsdk:"rotation_policy"`
	RotationListener                                    types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                                     types.Set    `tfsdk:"retention_policy"`
	CompressionMechanism                                types.String `tfsdk:"compression_mechanism"`
	SignLog                                             types.Bool   `tfsdk:"sign_log"`
	EncryptLog                                          types.Bool   `tfsdk:"encrypt_log"`
	EncryptionSettingsDefinitionID                      types.String `tfsdk:"encryption_settings_definition_id"`
	Append                                              types.Bool   `tfsdk:"append"`
	LogConnects                                         types.Bool   `tfsdk:"log_connects"`
	LogDisconnects                                      types.Bool   `tfsdk:"log_disconnects"`
	LogRequests                                         types.Bool   `tfsdk:"log_requests"`
	LogAssuranceCompleted                               types.Bool   `tfsdk:"log_assurance_completed"`
	Asynchronous                                        types.Bool   `tfsdk:"asynchronous"`
	AutoFlush                                           types.Bool   `tfsdk:"auto_flush"`
	BufferSize                                          types.String `tfsdk:"buffer_size"`
	QueueSize                                           types.Int64  `tfsdk:"queue_size"`
	TimeInterval                                        types.String `tfsdk:"time_interval"`
	WriteMultiLineMessages                              types.Bool   `tfsdk:"write_multi_line_messages"`
	IncludeProductName                                  types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                                 types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                                    types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                                     types.Bool   `tfsdk:"include_thread_id"`
	IncludeRequesterDN                                  types.Bool   `tfsdk:"include_requester_dn"`
	IncludeRequesterIPAddress                           types.Bool   `tfsdk:"include_requester_ip_address"`
	IncludeRequestDetailsInResultMessages               types.Bool   `tfsdk:"include_request_details_in_result_messages"`
	IncludeRequestDetailsInSearchEntryMessages          types.Bool   `tfsdk:"include_request_details_in_search_entry_messages"`
	IncludeRequestDetailsInSearchReferenceMessages      types.Bool   `tfsdk:"include_request_details_in_search_reference_messages"`
	IncludeRequestDetailsInIntermediateResponseMessages types.Bool   `tfsdk:"include_request_details_in_intermediate_response_messages"`
	IncludeResultCodeNames                              types.Bool   `tfsdk:"include_result_code_names"`
	IncludeExtendedSearchRequestDetails                 types.Bool   `tfsdk:"include_extended_search_request_details"`
	IncludeAddAttributeNames                            types.Bool   `tfsdk:"include_add_attribute_names"`
	IncludeModifyAttributeNames                         types.Bool   `tfsdk:"include_modify_attribute_names"`
	IncludeSearchEntryAttributeNames                    types.Bool   `tfsdk:"include_search_entry_attribute_names"`
	IncludeRequestControls                              types.Bool   `tfsdk:"include_request_controls"`
	IncludeResponseControls                             types.Bool   `tfsdk:"include_response_controls"`
	IncludeReplicationChangeID                          types.Bool   `tfsdk:"include_replication_change_id"`
	GenerifyMessageStringsWhenPossible                  types.Bool   `tfsdk:"generify_message_strings_when_possible"`
	MaxStringLength                                     types.Int64  `tfsdk:"max_string_length"`
	LogFieldBehavior                                    types.String `tfsdk:"log_field_behavior"`
	LogSecurityNegotiation                              types.Bool   `tfsdk:"log_security_negotiation"`
	LogClientCertificates                               types.Bool   `tfsdk:"log_client_certificates"`
	LogResults                                          types.Bool   `tfsdk:"log_results"`
	LogSearchEntries                                    types.Bool   `tfsdk:"log_search_entries"`
	LogSearchReferences                                 types.Bool   `tfsdk:"log_search_references"`
	LogIntermediateResponses                            types.Bool   `tfsdk:"log_intermediate_responses"`
	SuppressInternalOperations                          types.Bool   `tfsdk:"suppress_internal_operations"`
	SuppressReplicationOperations                       types.Bool   `tfsdk:"suppress_replication_operations"`
	CorrelateRequestsAndResults                         types.Bool   `tfsdk:"correlate_requests_and_results"`
	ConnectionCriteria                                  types.String `tfsdk:"connection_criteria"`
	RequestCriteria                                     types.String `tfsdk:"request_criteria"`
	ResultCriteria                                      types.String `tfsdk:"result_criteria"`
	SearchEntryCriteria                                 types.String `tfsdk:"search_entry_criteria"`
	SearchReferenceCriteria                             types.String `tfsdk:"search_reference_criteria"`
	Description                                         types.String `tfsdk:"description"`
	Enabled                                             types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                                types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *jsonAccessLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsonAccessLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultJsonAccessLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	jsonAccessLogPublisherSchema(ctx, req, resp, true)
}

func jsonAccessLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Json Access Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the JSON Access Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this JSON Access Log Publisher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the JSON Access Log Publisher .",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the JSON Access Log Publisher .",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"log_requests": schema.BoolAttribute{
				Description: "Indicates whether to log information about requests received from clients.",
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
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the JSON Access Log Publisher will publish records asynchronously.",
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
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
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
			"time_interval": schema.StringAttribute{
				Description: "Specifies the interval at which to check whether the log files need to be rotated.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"write_multi_line_messages": schema.BoolAttribute{
				Description: "Indicates whether the JSON objects should be formatted to span multiple lines with a single element on each line. The multi-line format is potentially more user friendly (if administrators may need to look at the log files), but each message will be larger because of the additional spaces and end-of-line markers.",
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
			"include_request_details_in_result_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation results should include information about both the request and the result.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
			"generify_message_strings_when_possible": schema.BoolAttribute{
				Description: "Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values.",
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
			"log_field_behavior": schema.StringAttribute{
				Description: "The behavior to use for determining which fields to log and whether to transform the values of those fields in any way.",
				Optional:    true,
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_client_certificates": schema.BoolAttribute{
				Description: "Indicates whether to log information about any client certificates presented to the server.",
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
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
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
			"correlate_requests_and_results": schema.BoolAttribute{
				Description: "Indicates whether to automatically log result messages for any operation in which the corresponding request was logged. In such cases, the result, entry, and reference criteria will be ignored, although the log-responses, log-search-entries, and log-search-references properties will be honored.",
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
			"search_entry_criteria": schema.StringAttribute{
				Description: "Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Access Log Publisher.",
				Optional:    true,
			},
			"search_reference_criteria": schema.StringAttribute{
				Description: "Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Access Log Publisher.",
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
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalJsonAccessLogPublisherFields(ctx context.Context, addRequest *client.AddJsonAccessLogPublisherRequest, plan jsonAccessLogPublisherResourceModel) error {
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
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
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
		intVal := int32(plan.MaxStringLength.ValueInt64())
		addRequest.MaxStringLength = &intVal
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

// Read a JsonAccessLogPublisherResponse object into the model struct
func readJsonAccessLogPublisherResponse(ctx context.Context, r *client.JsonAccessLogPublisherResponse, state *jsonAccessLogPublisherResourceModel, expectedValues *jsonAccessLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createJsonAccessLogPublisherOperations(plan jsonAccessLogPublisherResourceModel, state jsonAccessLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFilePermissions, state.LogFilePermissions, "log-file-permissions")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationPolicy, state.RotationPolicy, "rotation-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RotationListener, state.RotationListener, "rotation-listener")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RetentionPolicy, state.RetentionPolicy, "retention-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.CompressionMechanism, state.CompressionMechanism, "compression-mechanism")
	operations.AddBoolOperationIfNecessary(&ops, plan.SignLog, state.SignLog, "sign-log")
	operations.AddBoolOperationIfNecessary(&ops, plan.EncryptLog, state.EncryptLog, "encrypt-log")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.Append, state.Append, "append")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogConnects, state.LogConnects, "log-connects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogDisconnects, state.LogDisconnects, "log-disconnects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequests, state.LogRequests, "log-requests")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogAssuranceCompleted, state.LogAssuranceCompleted, "log-assurance-completed")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddStringOperationIfNecessary(&ops, plan.BufferSize, state.BufferSize, "buffer-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeInterval, state.TimeInterval, "time-interval")
	operations.AddBoolOperationIfNecessary(&ops, plan.WriteMultiLineMessages, state.WriteMultiLineMessages, "write-multi-line-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterDN, state.IncludeRequesterDN, "include-requester-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterIPAddress, state.IncludeRequesterIPAddress, "include-requester-ip-address")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInResultMessages, state.IncludeRequestDetailsInResultMessages, "include-request-details-in-result-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInSearchEntryMessages, state.IncludeRequestDetailsInSearchEntryMessages, "include-request-details-in-search-entry-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInSearchReferenceMessages, state.IncludeRequestDetailsInSearchReferenceMessages, "include-request-details-in-search-reference-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInIntermediateResponseMessages, state.IncludeRequestDetailsInIntermediateResponseMessages, "include-request-details-in-intermediate-response-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResultCodeNames, state.IncludeResultCodeNames, "include-result-code-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeExtendedSearchRequestDetails, state.IncludeExtendedSearchRequestDetails, "include-extended-search-request-details")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeAddAttributeNames, state.IncludeAddAttributeNames, "include-add-attribute-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeModifyAttributeNames, state.IncludeModifyAttributeNames, "include-modify-attribute-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeSearchEntryAttributeNames, state.IncludeSearchEntryAttributeNames, "include-search-entry-attribute-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestControls, state.IncludeRequestControls, "include-request-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResponseControls, state.IncludeResponseControls, "include-response-controls")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeReplicationChangeID, state.IncludeReplicationChangeID, "include-replication-change-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.GenerifyMessageStringsWhenPossible, state.GenerifyMessageStringsWhenPossible, "generify-message-strings-when-possible")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxStringLength, state.MaxStringLength, "max-string-length")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFieldBehavior, state.LogFieldBehavior, "log-field-behavior")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSecurityNegotiation, state.LogSecurityNegotiation, "log-security-negotiation")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogClientCertificates, state.LogClientCertificates, "log-client-certificates")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResults, state.LogResults, "log-results")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchEntries, state.LogSearchEntries, "log-search-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchReferences, state.LogSearchReferences, "log-search-references")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogIntermediateResponses, state.LogIntermediateResponses, "log-intermediate-responses")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressInternalOperations, state.SuppressInternalOperations, "suppress-internal-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressReplicationOperations, state.SuppressReplicationOperations, "suppress-replication-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.CorrelateRequestsAndResults, state.CorrelateRequestsAndResults, "correlate-requests-and-results")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCriteria, state.ResultCriteria, "result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchEntryCriteria, state.SearchEntryCriteria, "search-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchReferenceCriteria, state.SearchReferenceCriteria, "search-reference-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *jsonAccessLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jsonAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddJsonAccessLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumjsonAccessLogPublisherSchemaUrn{client.ENUMJSONACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERJSON_ACCESS},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalJsonAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Json Access Log Publisher", err.Error())
		return
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Json Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state jsonAccessLogPublisherResourceModel
	readJsonAccessLogPublisherResponse(ctx, addResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultJsonAccessLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan jsonAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state jsonAccessLogPublisherResourceModel
	readJsonAccessLogPublisherResponse(ctx, readResponse.JsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createJsonAccessLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Json Access Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJsonAccessLogPublisherResponse(ctx, updateResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *jsonAccessLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJsonAccessLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJsonAccessLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readJsonAccessLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readJsonAccessLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state jsonAccessLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJsonAccessLogPublisherResponse(ctx, readResponse.JsonAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *jsonAccessLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJsonAccessLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultJsonAccessLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateJsonAccessLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateJsonAccessLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan jsonAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state jsonAccessLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createJsonAccessLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Json Access Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readJsonAccessLogPublisherResponse(ctx, updateResponse.JsonAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultJsonAccessLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *jsonAccessLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state jsonAccessLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Json Access Log Publisher", err, httpResp)
		return
	}
}

func (r *jsonAccessLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJsonAccessLogPublisher(ctx, req, resp)
}

func (r *defaultJsonAccessLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importJsonAccessLogPublisher(ctx, req, resp)
}

func importJsonAccessLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
