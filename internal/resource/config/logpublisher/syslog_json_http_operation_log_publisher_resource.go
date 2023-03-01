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
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &syslogJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &syslogJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithImportState = &syslogJsonHttpOperationLogPublisherResource{}
	_ resource.Resource                = &defaultSyslogJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultSyslogJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultSyslogJsonHttpOperationLogPublisherResource{}
)

// Create a Syslog Json Http Operation Log Publisher resource
func NewSyslogJsonHttpOperationLogPublisherResource() resource.Resource {
	return &syslogJsonHttpOperationLogPublisherResource{}
}

func NewDefaultSyslogJsonHttpOperationLogPublisherResource() resource.Resource {
	return &defaultSyslogJsonHttpOperationLogPublisherResource{}
}

// syslogJsonHttpOperationLogPublisherResource is the resource implementation.
type syslogJsonHttpOperationLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSyslogJsonHttpOperationLogPublisherResource is the resource implementation.
type defaultSyslogJsonHttpOperationLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *syslogJsonHttpOperationLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_syslog_json_http_operation_log_publisher"
}

func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_syslog_json_http_operation_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *syslogJsonHttpOperationLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type syslogJsonHttpOperationLogPublisherResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	SyslogExternalServer                  types.Set    `tfsdk:"syslog_external_server"`
	SyslogFacility                        types.String `tfsdk:"syslog_facility"`
	SyslogSeverity                        types.String `tfsdk:"syslog_severity"`
	SyslogMessageHostName                 types.String `tfsdk:"syslog_message_host_name"`
	SyslogMessageApplicationName          types.String `tfsdk:"syslog_message_application_name"`
	QueueSize                             types.Int64  `tfsdk:"queue_size"`
	LogRequests                           types.Bool   `tfsdk:"log_requests"`
	LogResults                            types.Bool   `tfsdk:"log_results"`
	IncludeProductName                    types.Bool   `tfsdk:"include_product_name"`
	IncludeInstanceName                   types.Bool   `tfsdk:"include_instance_name"`
	IncludeStartupID                      types.Bool   `tfsdk:"include_startup_id"`
	IncludeThreadID                       types.Bool   `tfsdk:"include_thread_id"`
	IncludeRequestDetailsInResultMessages types.Bool   `tfsdk:"include_request_details_in_result_messages"`
	LogRequestHeaders                     types.String `tfsdk:"log_request_headers"`
	SuppressedRequestHeaderName           types.Set    `tfsdk:"suppressed_request_header_name"`
	LogResponseHeaders                    types.String `tfsdk:"log_response_headers"`
	SuppressedResponseHeaderName          types.Set    `tfsdk:"suppressed_response_header_name"`
	LogRequestAuthorizationType           types.Bool   `tfsdk:"log_request_authorization_type"`
	LogRequestCookieNames                 types.Bool   `tfsdk:"log_request_cookie_names"`
	LogResponseCookieNames                types.Bool   `tfsdk:"log_response_cookie_names"`
	LogRequestParameters                  types.String `tfsdk:"log_request_parameters"`
	SuppressedRequestParameterName        types.Set    `tfsdk:"suppressed_request_parameter_name"`
	LogRequestProtocol                    types.Bool   `tfsdk:"log_request_protocol"`
	LogRedirectURI                        types.Bool   `tfsdk:"log_redirect_uri"`
	WriteMultiLineMessages                types.Bool   `tfsdk:"write_multi_line_messages"`
	Description                           types.String `tfsdk:"description"`
	Enabled                               types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                  types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *syslogJsonHttpOperationLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogJsonHttpOperationLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	syslogJsonHttpOperationLogPublisherSchema(ctx, req, resp, true)
}

func syslogJsonHttpOperationLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Syslog Json Http Operation Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"syslog_external_server": schema.SetAttribute{
				Description: "The syslog server to which messages should be sent.",
				Required:    true,
				ElementType: types.StringType,
			},
			"syslog_facility": schema.StringAttribute{
				Description: "The syslog facility to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_severity": schema.StringAttribute{
				Description: "The syslog severity to use for the messages that are logged by this Syslog JSON HTTP Operation Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_message_host_name": schema.StringAttribute{
				Description: "The local host name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"syslog_message_application_name": schema.StringAttribute{
				Description: "The application name that will be included in syslog messages that are logged by this Syslog JSON HTTP Operation Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
			},
			"log_requests": schema.BoolAttribute{
				Description: "Indicates whether to record a log message with information about requests received from the client.",
				Optional:    true,
				Computed:    true,
			},
			"log_results": schema.BoolAttribute{
				Description: "Indicates whether to record a log message with information about the result of processing a requested HTTP operation.",
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
			"include_request_details_in_result_messages": schema.BoolAttribute{
				Description: "Indicates whether result log messages should include all of the elements of request log messages. This may be used to record a single message per operation with details about both the request and response.",
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
			"suppressed_request_parameter_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request parameters that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-parameters property has a value of parameter-names or parameter-names-and-values.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_request_protocol": schema.BoolAttribute{
				Description: "Indicates whether request log messages should include information about the HTTP version specified in the request.",
				Optional:    true,
				Computed:    true,
			},
			"log_redirect_uri": schema.BoolAttribute{
				Description: "Indicates whether the redirect URI (i.e., the value of the \"Location\" header from responses) should be included in response log messages.",
				Optional:    true,
				Computed:    true,
			},
			"write_multi_line_messages": schema.BoolAttribute{
				Description: "Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools.",
				Optional:    true,
				Computed:    true,
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
	if setOptionalToComputed {
		config.SetOptionalAttributesToComputed(&schema)
	}
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalSyslogJsonHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddSyslogJsonHttpOperationLogPublisherRequest, plan syslogJsonHttpOperationLogPublisherResourceModel) error {
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
		stringVal := plan.SyslogMessageHostName.ValueString()
		addRequest.SyslogMessageHostName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SyslogMessageApplicationName) {
		stringVal := plan.SyslogMessageApplicationName.ValueString()
		addRequest.SyslogMessageApplicationName = &stringVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	if internaltypes.IsDefined(plan.LogRequests) {
		boolVal := plan.LogRequests.ValueBool()
		addRequest.LogRequests = &boolVal
	}
	if internaltypes.IsDefined(plan.LogResults) {
		boolVal := plan.LogResults.ValueBool()
		addRequest.LogResults = &boolVal
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
	if internaltypes.IsDefined(plan.IncludeRequestDetailsInResultMessages) {
		boolVal := plan.IncludeRequestDetailsInResultMessages.ValueBool()
		addRequest.IncludeRequestDetailsInResultMessages = &boolVal
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
		boolVal := plan.LogRequestAuthorizationType.ValueBool()
		addRequest.LogRequestAuthorizationType = &boolVal
	}
	if internaltypes.IsDefined(plan.LogRequestCookieNames) {
		boolVal := plan.LogRequestCookieNames.ValueBool()
		addRequest.LogRequestCookieNames = &boolVal
	}
	if internaltypes.IsDefined(plan.LogResponseCookieNames) {
		boolVal := plan.LogResponseCookieNames.ValueBool()
		addRequest.LogResponseCookieNames = &boolVal
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
		boolVal := plan.LogRequestProtocol.ValueBool()
		addRequest.LogRequestProtocol = &boolVal
	}
	if internaltypes.IsDefined(plan.LogRedirectURI) {
		boolVal := plan.LogRedirectURI.ValueBool()
		addRequest.LogRedirectURI = &boolVal
	}
	if internaltypes.IsDefined(plan.WriteMultiLineMessages) {
		boolVal := plan.WriteMultiLineMessages.ValueBool()
		addRequest.WriteMultiLineMessages = &boolVal
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

// Read a SyslogJsonHttpOperationLogPublisherResponse object into the model struct
func readSyslogJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.SyslogJsonHttpOperationLogPublisherResponse, state *syslogJsonHttpOperationLogPublisherResourceModel, expectedValues *syslogJsonHttpOperationLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createSyslogJsonHttpOperationLogPublisherOperations(plan syslogJsonHttpOperationLogPublisherResourceModel, state syslogJsonHttpOperationLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SyslogExternalServer, state.SyslogExternalServer, "syslog-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogFacility, state.SyslogFacility, "syslog-facility")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogSeverity, state.SyslogSeverity, "syslog-severity")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageHostName, state.SyslogMessageHostName, "syslog-message-host-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SyslogMessageApplicationName, state.SyslogMessageApplicationName, "syslog-message-application-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequests, state.LogRequests, "log-requests")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResults, state.LogResults, "log-results")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeProductName, state.IncludeProductName, "include-product-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeInstanceName, state.IncludeInstanceName, "include-instance-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStartupID, state.IncludeStartupID, "include-startup-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeThreadID, state.IncludeThreadID, "include-thread-id")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeRequestDetailsInResultMessages, state.IncludeRequestDetailsInResultMessages, "include-request-details-in-result-messages")
	operations.AddStringOperationIfNecessary(&ops, plan.LogRequestHeaders, state.LogRequestHeaders, "log-request-headers")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedRequestHeaderName, state.SuppressedRequestHeaderName, "suppressed-request-header-name")
	operations.AddStringOperationIfNecessary(&ops, plan.LogResponseHeaders, state.LogResponseHeaders, "log-response-headers")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedResponseHeaderName, state.SuppressedResponseHeaderName, "suppressed-response-header-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestAuthorizationType, state.LogRequestAuthorizationType, "log-request-authorization-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestCookieNames, state.LogRequestCookieNames, "log-request-cookie-names")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogResponseCookieNames, state.LogResponseCookieNames, "log-response-cookie-names")
	operations.AddStringOperationIfNecessary(&ops, plan.LogRequestParameters, state.LogRequestParameters, "log-request-parameters")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedRequestParameterName, state.SuppressedRequestParameterName, "suppressed-request-parameter-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestProtocol, state.LogRequestProtocol, "log-request-protocol")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRedirectURI, state.LogRedirectURI, "log-redirect-uri")
	operations.AddBoolOperationIfNecessary(&ops, plan.WriteMultiLineMessages, state.WriteMultiLineMessages, "write-multi-line-messages")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *syslogJsonHttpOperationLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogJsonHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var SyslogExternalServerSlice []string
	plan.SyslogExternalServer.ElementsAs(ctx, &SyslogExternalServerSlice, false)
	addRequest := client.NewAddSyslogJsonHttpOperationLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumsyslogJsonHttpOperationLogPublisherSchemaUrn{client.ENUMSYSLOGJSONHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERSYSLOG_JSON_HTTP_OPERATION},
		SyslogExternalServerSlice,
		plan.Enabled.ValueBool())
	err := addOptionalSyslogJsonHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Syslog Json Http Operation Log Publisher", err.Error())
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
		client.AddSyslogJsonHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Syslog Json Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state syslogJsonHttpOperationLogPublisherResourceModel
	readSyslogJsonHttpOperationLogPublisherResponse(ctx, addResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan syslogJsonHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Json Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state syslogJsonHttpOperationLogPublisherResourceModel
	readSyslogJsonHttpOperationLogPublisherResponse(ctx, readResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSyslogJsonHttpOperationLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Json Http Operation Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogJsonHttpOperationLogPublisherResponse(ctx, updateResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *syslogJsonHttpOperationLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSyslogJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSyslogJsonHttpOperationLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state syslogJsonHttpOperationLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Syslog Json Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSyslogJsonHttpOperationLogPublisherResponse(ctx, readResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *syslogJsonHttpOperationLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSyslogJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSyslogJsonHttpOperationLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan syslogJsonHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state syslogJsonHttpOperationLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSyslogJsonHttpOperationLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Syslog Json Http Operation Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSyslogJsonHttpOperationLogPublisherResponse(ctx, updateResponse.SyslogJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSyslogJsonHttpOperationLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *syslogJsonHttpOperationLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state syslogJsonHttpOperationLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Syslog Json Http Operation Log Publisher", err, httpResp)
		return
	}
}

func (r *syslogJsonHttpOperationLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogJsonHttpOperationLogPublisher(ctx, req, resp)
}

func (r *defaultSyslogJsonHttpOperationLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSyslogJsonHttpOperationLogPublisher(ctx, req, resp)
}

func importSyslogJsonHttpOperationLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
