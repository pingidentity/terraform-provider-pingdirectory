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
	_ resource.Resource                = &detailedHttpOperationLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &detailedHttpOperationLogPublisherResource{}
	_ resource.ResourceWithImportState = &detailedHttpOperationLogPublisherResource{}
)

// Create a Detailed Http Operation Log Publisher resource
func NewDetailedHttpOperationLogPublisherResource() resource.Resource {
	return &detailedHttpOperationLogPublisherResource{}
}

// detailedHttpOperationLogPublisherResource is the resource implementation.
type detailedHttpOperationLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *detailedHttpOperationLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_detailed_http_operation_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *detailedHttpOperationLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type detailedHttpOperationLogPublisherResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	LogFile                               types.String `tfsdk:"log_file"`
	LogFilePermissions                    types.String `tfsdk:"log_file_permissions"`
	RotationPolicy                        types.Set    `tfsdk:"rotation_policy"`
	RotationListener                      types.Set    `tfsdk:"rotation_listener"`
	RetentionPolicy                       types.Set    `tfsdk:"retention_policy"`
	CompressionMechanism                  types.String `tfsdk:"compression_mechanism"`
	SignLog                               types.Bool   `tfsdk:"sign_log"`
	EncryptLog                            types.Bool   `tfsdk:"encrypt_log"`
	EncryptionSettingsDefinitionID        types.String `tfsdk:"encryption_settings_definition_id"`
	Append                                types.Bool   `tfsdk:"append"`
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
	LogRequestProtocol                    types.Bool   `tfsdk:"log_request_protocol"`
	SuppressedRequestParameterName        types.Set    `tfsdk:"suppressed_request_parameter_name"`
	LogRedirectURI                        types.Bool   `tfsdk:"log_redirect_uri"`
	Asynchronous                          types.Bool   `tfsdk:"asynchronous"`
	AutoFlush                             types.Bool   `tfsdk:"auto_flush"`
	BufferSize                            types.String `tfsdk:"buffer_size"`
	MaxStringLength                       types.Int64  `tfsdk:"max_string_length"`
	QueueSize                             types.Int64  `tfsdk:"queue_size"`
	TimeInterval                          types.String `tfsdk:"time_interval"`
	Description                           types.String `tfsdk:"description"`
	Enabled                               types.Bool   `tfsdk:"enabled"`
	LoggingErrorBehavior                  types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *detailedHttpOperationLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Detailed Http Operation Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"log_file": schema.StringAttribute{
				Description: "The file name to use for the log files generated by the Detailed HTTP Operation Log Publisher. The path to the file can be specified either as relative to the server root or as an absolute path.",
				Required:    true,
			},
			"log_file_permissions": schema.StringAttribute{
				Description: "The UNIX permissions of the log files created by this Detailed HTTP Operation Log Publisher.",
				Optional:    true,
				Computed:    true,
			},
			"rotation_policy": schema.SetAttribute{
				Description: "The rotation policy to use for the Detailed HTTP Operation Log Publisher .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"rotation_listener": schema.SetAttribute{
				Description: "A listener that should be notified whenever a log file is rotated out of service.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"retention_policy": schema.SetAttribute{
				Description: "The retention policy to use for the Detailed HTTP Operation Log Publisher .",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"compression_mechanism": schema.StringAttribute{
				Description: "Specifies the type of compression (if any) to use for log files that are written.",
				Optional:    true,
				Computed:    true,
			},
			"sign_log": schema.BoolAttribute{
				Description: "Indicates whether the log should be cryptographically signed so that the log content cannot be altered in an undetectable manner.",
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
			"log_request_protocol": schema.BoolAttribute{
				Description: "Indicates whether request log messages should include information about the HTTP version specified in the request.",
				Optional:    true,
				Computed:    true,
			},
			"suppressed_request_parameter_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request parameters that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-parameters property has a value of parameter-names or parameter-names-and-values.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"log_redirect_uri": schema.BoolAttribute{
				Description: "Indicates whether the redirect URI (i.e., the value of the \"Location\" header from responses) should be included in response log messages.",
				Optional:    true,
				Computed:    true,
			},
			"asynchronous": schema.BoolAttribute{
				Description: "Indicates whether the Detailed HTTP Operation Log Publisher will publish records asynchronously.",
				Optional:    true,
				Computed:    true,
			},
			"auto_flush": schema.BoolAttribute{
				Description: "Specifies whether to flush the writer after every log record.",
				Optional:    true,
				Computed:    true,
			},
			"buffer_size": schema.StringAttribute{
				Description: "Specifies the log file buffer size.",
				Optional:    true,
				Computed:    true,
			},
			"max_string_length": schema.Int64Attribute{
				Description: "Specifies the maximum length of any individual string that should be logged. If a log message includes a string longer than this number of characters, it will be truncated. A value of zero indicates that no truncation will be used.",
				Optional:    true,
				Computed:    true,
			},
			"queue_size": schema.Int64Attribute{
				Description: "The maximum number of log records that can be stored in the asynchronous queue.",
				Optional:    true,
				Computed:    true,
			},
			"time_interval": schema.StringAttribute{
				Description: "Specifies the interval at which to check whether the log files need to be rotated.",
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
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalDetailedHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddDetailedHttpOperationLogPublisherRequest, plan detailedHttpOperationLogPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LogFilePermissions) {
		stringVal := plan.LogFilePermissions.ValueString()
		addRequest.LogFilePermissions = &stringVal
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
		boolVal := plan.SignLog.ValueBool()
		addRequest.SignLog = &boolVal
	}
	if internaltypes.IsDefined(plan.EncryptLog) {
		boolVal := plan.EncryptLog.ValueBool()
		addRequest.EncryptLog = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		stringVal := plan.EncryptionSettingsDefinitionID.ValueString()
		addRequest.EncryptionSettingsDefinitionID = &stringVal
	}
	if internaltypes.IsDefined(plan.Append) {
		boolVal := plan.Append.ValueBool()
		addRequest.Append = &boolVal
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
	if internaltypes.IsDefined(plan.LogRequestProtocol) {
		boolVal := plan.LogRequestProtocol.ValueBool()
		addRequest.LogRequestProtocol = &boolVal
	}
	if internaltypes.IsDefined(plan.SuppressedRequestParameterName) {
		var slice []string
		plan.SuppressedRequestParameterName.ElementsAs(ctx, &slice, false)
		addRequest.SuppressedRequestParameterName = slice
	}
	if internaltypes.IsDefined(plan.LogRedirectURI) {
		boolVal := plan.LogRedirectURI.ValueBool()
		addRequest.LogRedirectURI = &boolVal
	}
	if internaltypes.IsDefined(plan.Asynchronous) {
		boolVal := plan.Asynchronous.ValueBool()
		addRequest.Asynchronous = &boolVal
	}
	if internaltypes.IsDefined(plan.AutoFlush) {
		boolVal := plan.AutoFlush.ValueBool()
		addRequest.AutoFlush = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BufferSize) {
		stringVal := plan.BufferSize.ValueString()
		addRequest.BufferSize = &stringVal
	}
	if internaltypes.IsDefined(plan.MaxStringLength) {
		intVal := int32(plan.MaxStringLength.ValueInt64())
		addRequest.MaxStringLength = &intVal
	}
	if internaltypes.IsDefined(plan.QueueSize) {
		intVal := int32(plan.QueueSize.ValueInt64())
		addRequest.QueueSize = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeInterval) {
		stringVal := plan.TimeInterval.ValueString()
		addRequest.TimeInterval = &stringVal
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

// Read a DetailedHttpOperationLogPublisherResponse object into the model struct
func readDetailedHttpOperationLogPublisherResponse(ctx context.Context, r *client.DetailedHttpOperationLogPublisherResponse, state *detailedHttpOperationLogPublisherResourceModel, expectedValues *detailedHttpOperationLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createDetailedHttpOperationLogPublisherOperations(plan detailedHttpOperationLogPublisherResourceModel, state detailedHttpOperationLogPublisherResourceModel) []client.Operation {
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
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRequestProtocol, state.LogRequestProtocol, "log-request-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SuppressedRequestParameterName, state.SuppressedRequestParameterName, "suppressed-request-parameter-name")
	operations.AddBoolOperationIfNecessary(&ops, plan.LogRedirectURI, state.LogRedirectURI, "log-redirect-uri")
	operations.AddBoolOperationIfNecessary(&ops, plan.Asynchronous, state.Asynchronous, "asynchronous")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoFlush, state.AutoFlush, "auto-flush")
	operations.AddStringOperationIfNecessary(&ops, plan.BufferSize, state.BufferSize, "buffer-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxStringLength, state.MaxStringLength, "max-string-length")
	operations.AddInt64OperationIfNecessary(&ops, plan.QueueSize, state.QueueSize, "queue-size")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeInterval, state.TimeInterval, "time-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *detailedHttpOperationLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan detailedHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddDetailedHttpOperationLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumdetailedHttpOperationLogPublisherSchemaUrn{client.ENUMDETAILEDHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERDETAILED_HTTP_OPERATION},
		plan.LogFile.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalDetailedHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Detailed Http Operation Log Publisher", err.Error())
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
		client.AddDetailedHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Detailed Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state detailedHttpOperationLogPublisherResourceModel
	readDetailedHttpOperationLogPublisherResponse(ctx, addResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *detailedHttpOperationLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state detailedHttpOperationLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Detailed Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDetailedHttpOperationLogPublisherResponse(ctx, readResponse.DetailedHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *detailedHttpOperationLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan detailedHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state detailedHttpOperationLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDetailedHttpOperationLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Detailed Http Operation Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDetailedHttpOperationLogPublisherResponse(ctx, updateResponse.DetailedHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *detailedHttpOperationLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state detailedHttpOperationLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Detailed Http Operation Log Publisher", err, httpResp)
		return
	}
}

func (r *detailedHttpOperationLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
