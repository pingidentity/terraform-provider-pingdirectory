package logpublisher

import (
	"context"
	"time"

	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &adminAlertAccessLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &adminAlertAccessLogPublisherResource{}
	_ resource.ResourceWithImportState = &adminAlertAccessLogPublisherResource{}
)

// Create a Admin Alert Access Log Publisher resource
func NewAdminAlertAccessLogPublisherResource() resource.Resource {
	return &adminAlertAccessLogPublisherResource{}
}

// adminAlertAccessLogPublisherResource is the resource implementation.
type adminAlertAccessLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *adminAlertAccessLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_admin_alert_access_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *adminAlertAccessLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type adminAlertAccessLogPublisherResourceModel struct {
	Id                                                  types.String `tfsdk:"id"`
	LastUpdated                                         types.String `tfsdk:"last_updated"`
	Notifications                                       types.Set    `tfsdk:"notifications"`
	RequiredActions                                     types.Set    `tfsdk:"required_actions"`
	LogConnects                                         types.Bool   `tfsdk:"log_connects"`
	LogDisconnects                                      types.Bool   `tfsdk:"log_disconnects"`
	LogClientCertificates                               types.Bool   `tfsdk:"log_client_certificates"`
	LogRequests                                         types.Bool   `tfsdk:"log_requests"`
	LogResults                                          types.Bool   `tfsdk:"log_results"`
	LogSearchEntries                                    types.Bool   `tfsdk:"log_search_entries"`
	LogSearchReferences                                 types.Bool   `tfsdk:"log_search_references"`
	SearchEntryCriteria                                 types.String `tfsdk:"search_entry_criteria"`
	SearchReferenceCriteria                             types.String `tfsdk:"search_reference_criteria"`
	CorrelateRequestsAndResults                         types.Bool   `tfsdk:"correlate_requests_and_results"`
	AutoFlush                                           types.Bool   `tfsdk:"auto_flush"`
	Asynchronous                                        types.Bool   `tfsdk:"asynchronous"`
	QueueSize                                           types.Int64  `tfsdk:"queue_size"`
	IncludeRequestDetailsInResultMessages               types.Bool   `tfsdk:"include_request_details_in_result_messages"`
	LogAssuranceCompleted                               types.Bool   `tfsdk:"log_assurance_completed"`
	IncludeProductName                                  types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                                 types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                                    types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                                     types.Bool   `tfsdk:"include_thread_id"`
	IncludeRequesterDN                                  types.Bool   `tfsdk:"include_requester_dn"`
	IncludeRequesterIPAddress                           types.Bool   `tfsdk:"include_requester_ip_address"`
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
	LogIntermediateResponses                            types.Bool   `tfsdk:"log_intermediate_responses"`
	SuppressInternalOperations                          types.Bool   `tfsdk:"suppress_internal_operations"`
	SuppressReplicationOperations                       types.Bool   `tfsdk:"suppress_replication_operations"`
	ConnectionCriteria                                  types.String `tfsdk:"connection_criteria"`
	RequestCriteria                                     types.String `tfsdk:"request_criteria"`
	ResultCriteria                                      types.String `tfsdk:"result_criteria"`
	Description                                         types.String `tfsdk:"description"`
	Enabled                                             types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                                types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *adminAlertAccessLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Admin Alert Access Log Publisher.",
		Attributes: map[string]schema.Attribute{
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
			"log_client_certificates": schema.BoolAttribute{
				Description: "Indicates whether to log information about any client certificates presented to the server.",
				Optional:    true,
				Computed:    true,
			},
			"log_requests": schema.BoolAttribute{
				Description: "Indicates whether to log information about requests received from clients.",
				Optional:    true,
				Computed:    true,
			},
			"log_results": schema.BoolAttribute{
				Description: "Indicates whether to log information about the results of client requests.",
				Optional:    true,
				Computed:    true,
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
			"search_entry_criteria": schema.StringAttribute{
				Description: "Specifies a set of search entry criteria that must match the associated search result entry in order for that it to be logged by this Admin Alert Access Log Publisher.",
				Optional:    true,
			},
			"search_reference_criteria": schema.StringAttribute{
				Description: "Specifies a set of search reference criteria that must match the associated search result reference in order for that it to be logged by this Admin Alert Access Log Publisher.",
				Optional:    true,
			},
			"correlate_requests_and_results": schema.BoolAttribute{
				Description: "Indicates whether to automatically log result messages for any operation in which the corresponding request was logged. In such cases, the result, entry, and reference criteria will be ignored, although the log-responses, log-search-entries, and log-search-references properties will be honored.",
				Optional:    true,
				Computed:    true,
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the Admin Alert Access Log Publisher will publish records asynchronously.",
				Required:    true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
			},
			"include_request_details_in_result_messages": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation results should include information about both the request and the result.",
				Optional:    true,
				Computed:    true,
			},
			"log_assurance_completed": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of replication assurance processing.",
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
				Description: "Indicates whether log messages for operation requests should include the DN of the authenticated user for the client connection on which the operation was requested.",
				Optional:    true,
				Computed:    true,
			},
			"include_requester_ip_address": schema.BoolAttribute{
				Description: "Indicates whether log messages for operation requests should include the IP address of the client that requested the operation.",
				Optional:    true,
				Computed:    true,
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
			"generify_message_strings_when_possible": schema.BoolAttribute{
				Description: "Indicates whether to use generified version of certain message strings, including diagnostic messages, additional information messages, authentication failure reasons, and disconnect messages. Generified versions of those strings may use placeholders (like %s for a string or %d for an integer) rather than the version of the string with those placeholders replaced with specific values.",
				Optional:    true,
				Computed:    true,
			},
			"max_string_length": schema.Int64Attribute{
				Description: "Specifies the maximum number of characters that may be included in any string in a log message before that string is truncated and replaced with a placeholder indicating the number of characters that were omitted. This can help prevent extremely long log messages from being written.",
				Optional:    true,
				Computed:    true,
			},
			"log_field_behavior": schema.StringAttribute{
				Description: "The behavior to use for determining which fields to log and whether to transform the values of those fields in any way.",
				Optional:    true,
			},
			"log_security_negotiation": schema.BoolAttribute{
				Description: "Indicates whether to log information about the result of any security negotiation (e.g., SSL handshake) processing that has been performed.",
				Optional:    true,
				Computed:    true,
			},
			"log_intermediate_responses": schema.BoolAttribute{
				Description: "Indicates whether to log information about intermediate responses sent to the client.",
				Optional:    true,
				Computed:    true,
			},
			"suppress_internal_operations": schema.BoolAttribute{
				Description: "Indicates whether internal operations (for example, operations that are initiated by plugins) should be logged along with the operations that are requested by users.",
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
				Description: "Indicates whether the Log Publisher is enabled for use.",
				Required:    true,
			},
			"logging_error_behavior": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit if an error occurs during logging processing.",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalAdminAlertAccessLogPublisherFields(ctx context.Context, addRequest *client.AddAdminAlertAccessLogPublisherRequest, plan adminAlertAccessLogPublisherResourceModel) error {
	if internaltypes.IsDefined(plan.LogConnects) {
		boolVal := plan.LogConnects.ValueBool()
		addRequest.LogConnects = &boolVal
	}
	if internaltypes.IsDefined(plan.LogDisconnects) {
		boolVal := plan.LogDisconnects.ValueBool()
		addRequest.LogDisconnects = &boolVal
	}
	if internaltypes.IsDefined(plan.LogClientCertificates) {
		boolVal := plan.LogClientCertificates.ValueBool()
		addRequest.LogClientCertificates = &boolVal
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		boolVal := plan.LogRequests.ValueBool()
		addRequest.LogRequests = &boolVal
	}
	if internaltypes.IsDefined(plan.LogResults) {
		boolVal := plan.LogResults.ValueBool()
		addRequest.LogResults = &boolVal
	}
	if internaltypes.IsDefined(plan.LogSearchEntries) {
		boolVal := plan.LogSearchEntries.ValueBool()
		addRequest.LogSearchEntries = &boolVal
	}
	if internaltypes.IsDefined(plan.LogSearchReferences) {
		boolVal := plan.LogSearchReferences.ValueBool()
		addRequest.LogSearchReferences = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchEntryCriteria) {
		stringVal := plan.SearchEntryCriteria.ValueString()
		addRequest.SearchEntryCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SearchReferenceCriteria) {
		stringVal := plan.SearchReferenceCriteria.ValueString()
		addRequest.SearchReferenceCriteria = &stringVal
	}
	if internaltypes.IsDefined(plan.CorrelateRequestsAndResults) {
		boolVal := plan.CorrelateRequestsAndResults.ValueBool()
		addRequest.CorrelateRequestsAndResults = &boolVal
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		boolVal := plan.AutoFlush.ValueBool()
		addRequest.AutoFlush = &boolVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		boolVal := plan.IncludeRequestDetailsInResultMessages.ValueBool()
		addRequest.IncludeRequestDetailsInResultMessages = &boolVal
	}
	if internaltypes.IsDefined(plan.LogAssuranceCompleted) {
		boolVal := plan.LogAssuranceCompleted.ValueBool()
		addRequest.LogAssuranceCompleted = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeProductName) {
		boolVal := plan.IncludeProductName.ValueBool()
		addRequest.IncludeProductName = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeInstanceName) {
		boolVal := plan.IncludeInstanceName.ValueBool()
		addRequest.IncludeInstanceName = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeStartupID) {
		boolVal := plan.IncludeStartupID.ValueBool()
		addRequest.IncludeStartupID = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeThreadID) {
		boolVal := plan.IncludeThreadID.ValueBool()
		addRequest.IncludeThreadID = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequesterDN) {
		boolVal := plan.IncludeRequesterDN.ValueBool()
		addRequest.IncludeRequesterDN = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequesterIPAddress) {
		boolVal := plan.IncludeRequesterIPAddress.ValueBool()
		addRequest.IncludeRequesterIPAddress = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchEntryMessages) {
		boolVal := plan.IncludeRequestDetailsInSearchEntryMessages.ValueBool()
		addRequest.IncludeRequestDetailsInSearchEntryMessages = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInSearchReferenceMessages) {
		boolVal := plan.IncludeRequestDetailsInSearchReferenceMessages.ValueBool()
		addRequest.IncludeRequestDetailsInSearchReferenceMessages = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInIntermediateResponseMessages) {
		boolVal := plan.IncludeRequestDetailsInIntermediateResponseMessages.ValueBool()
		addRequest.IncludeRequestDetailsInIntermediateResponseMessages = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeResultCodeNames) {
		boolVal := plan.IncludeResultCodeNames.ValueBool()
		addRequest.IncludeResultCodeNames = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeExtendedSearchRequestDetails) {
		boolVal := plan.IncludeExtendedSearchRequestDetails.ValueBool()
		addRequest.IncludeExtendedSearchRequestDetails = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeAddAttributeNames) {
		boolVal := plan.IncludeAddAttributeNames.ValueBool()
		addRequest.IncludeAddAttributeNames = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeModifyAttributeNames) {
		boolVal := plan.IncludeModifyAttributeNames.ValueBool()
		addRequest.IncludeModifyAttributeNames = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeSearchEntryAttributeNames) {
		boolVal := plan.IncludeSearchEntryAttributeNames.ValueBool()
		addRequest.IncludeSearchEntryAttributeNames = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeRequestControls) {
		boolVal := plan.IncludeRequestControls.ValueBool()
		addRequest.IncludeRequestControls = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeResponseControls) {
		boolVal := plan.IncludeResponseControls.ValueBool()
		addRequest.IncludeResponseControls = &boolVal
	}
	if internaltypes.IsDefined(plan.IncludeReplicationChangeID) {
		boolVal := plan.IncludeReplicationChangeID.ValueBool()
		addRequest.IncludeReplicationChangeID = &boolVal
	}
	if internaltypes.IsDefined(plan.GenerifyMessageStringsWhenPossible) {
		boolVal := plan.GenerifyMessageStringsWhenPossible.ValueBool()
		addRequest.GenerifyMessageStringsWhenPossible = &boolVal
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		intVal := int32(plan.MaxStringLength.ValueInt64())
		addRequest.MaxStringLength = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFieldBehavior) {
		stringVal := plan.LogFieldBehavior.ValueString()
		addRequest.LogFieldBehavior = &stringVal
	}
	if internaltypes.IsDefined(plan.LogSecurityNegotiation) {
		boolVal := plan.LogSecurityNegotiation.ValueBool()
		addRequest.LogSecurityNegotiation = &boolVal
	}
	if internaltypes.IsDefined(plan.LogIntermediateResponses) {
		boolVal := plan.LogIntermediateResponses.ValueBool()
		addRequest.LogIntermediateResponses = &boolVal
	}
	if internaltypes.IsDefined(plan.SuppressInternalOperations) {
		boolVal := plan.SuppressInternalOperations.ValueBool()
		addRequest.SuppressInternalOperations = &boolVal
	}
	if internaltypes.IsDefined(plan.SuppressReplicationOperations) {
		boolVal := plan.SuppressReplicationOperations.ValueBool()
		addRequest.SuppressReplicationOperations = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionCriteria) {
		stringVal := plan.ConnectionCriteria.ValueString()
		addRequest.ConnectionCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.RequestCriteria) {
		stringVal := plan.RequestCriteria.ValueString()
		addRequest.RequestCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResultCriteria) {
		stringVal := plan.ResultCriteria.ValueString()
		addRequest.ResultCriteria = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
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

// Read a AdminAlertAccessLogPublisherResponse object into the model struct
func readAdminAlertAccessLogPublisherResponse(ctx context.Context, r *client.AdminAlertAccessLogPublisherResponse, state *adminAlertAccessLogPublisherResourceModel, expectedValues *adminAlertAccessLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createAdminAlertAccessLogPublisherOperations(plan adminAlertAccessLogPublisherResourceModel, state adminAlertAccessLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.LogConnects, state.LogConnects, "log-connects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogDisconnects, state.LogDisconnects, "log-disconnects")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogClientCertificates, state.LogClientCertificates, "log-client-certificates")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequests, state.LogRequests, "log-requests")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResults, state.LogResults, "log-results")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchEntries, state.LogSearchEntries, "log-search-entries")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogSearchReferences, state.LogSearchReferences, "log-search-references")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchEntryCriteria, state.SearchEntryCriteria, "search-entry-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SearchReferenceCriteria, state.SearchReferenceCriteria, "search-reference-criteria")
	operations.AddBoolOperationIfNecessary(&ops, plan.CorrelateRequestsAndResults, state.CorrelateRequestsAndResults, "correlate-requests-and-results")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInResultMessages, state.IncludeRequestDetailsInResultMessages, "include-request-details-in-result-messages")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogAssuranceCompleted, state.LogAssuranceCompleted, "log-assurance-completed")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterDN, state.IncludeRequesterDN, "include-requester-dn")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequesterIPAddress, state.IncludeRequesterIPAddress, "include-requester-ip-address")
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
	operations.AddBoolOperationIfNecessary(&ops, plan.LogIntermediateResponses, state.LogIntermediateResponses, "log-intermediate-responses")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressInternalOperations, state.SuppressInternalOperations, "suppress-internal-operations")
	operations.AddBoolOperationIfNecessary(&ops, plan.SuppressReplicationOperations, state.SuppressReplicationOperations, "suppress-replication-operations")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionCriteria, state.ConnectionCriteria, "connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.RequestCriteria, state.RequestCriteria, "request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.ResultCriteria, state.ResultCriteria, "result-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *adminAlertAccessLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan adminAlertAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAdminAlertAccessLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumadminAlertAccessLogPublisherSchemaUrn{client.ENUMADMINALERTACCESSLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERADMIN_ALERT_ACCESS},
		plan.Asynchronous.ValueBool(),
		plan.Enabled.ValueBool())
	err := addOptionalAdminAlertAccessLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Admin Alert Access Log Publisher", err.Error())
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
		client.AddAdminAlertAccessLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Admin Alert Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state adminAlertAccessLogPublisherResourceModel
	readAdminAlertAccessLogPublisherResponse(ctx, addResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *adminAlertAccessLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state adminAlertAccessLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Admin Alert Access Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAdminAlertAccessLogPublisherResponse(ctx, readResponse.AdminAlertAccessLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *adminAlertAccessLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan adminAlertAccessLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state adminAlertAccessLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAdminAlertAccessLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Admin Alert Access Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAdminAlertAccessLogPublisherResponse(ctx, updateResponse.AdminAlertAccessLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *adminAlertAccessLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state adminAlertAccessLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Admin Alert Access Log Publisher", err, httpResp)
		return
	}
}

func (r *adminAlertAccessLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
