package logpublisher

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &consoleJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &consoleJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithImportState = &consoleJsonHttpOperationLogPublisherResource{}
	_ resource.Resource                = &defaultConsoleJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithConfigure   = &defaultConsoleJsonHttpOperationLogPublisherResource{}
	_ resource.ResourceWithImportState = &defaultConsoleJsonHttpOperationLogPublisherResource{}
)

// Create a Console Json Http Operation Log Publisher resource
func NewConsoleJsonHttpOperationLogPublisherResource() resource.Resource {
	return &consoleJsonHttpOperationLogPublisherResource{}
}

func NewDefaultConsoleJsonHttpOperationLogPublisherResource() resource.Resource {
	return &defaultConsoleJsonHttpOperationLogPublisherResource{}
}

// consoleJsonHttpOperationLogPublisherResource is the resource implementation.
type consoleJsonHttpOperationLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConsoleJsonHttpOperationLogPublisherResource is the resource implementation.
type defaultConsoleJsonHttpOperationLogPublisherResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *consoleJsonHttpOperationLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_console_json_http_operation_log_publisher"
}

func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_console_json_http_operation_log_publisher"
}

// Configure adds the provider configured client to the resource.
func (r *consoleJsonHttpOperationLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type consoleJsonHttpOperationLogPublisherResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	Enabled                               types.Bool   `tfsdk:"enabled"`
	OutputLocation                        types.String `tfsdk:"output_location"`
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
	LoggingErrorBehavior                  types.String `tfsdk:"logging_error_behavior"`
}

// GetSchema defines the schema for the resource.
func (r *consoleJsonHttpOperationLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	consoleJsonHttpOperationLogPublisherSchema(ctx, req, resp, false)
}

func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	consoleJsonHttpOperationLogPublisherSchema(ctx, req, resp, true)
}

func consoleJsonHttpOperationLogPublisherSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Console Json Http Operation Log Publisher.",
		Attributes: map[string]schema.Attribute{
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Console JSON HTTP Operation Log Publisher is enabled for use.",
				Required:    true,
			},
			"output_location": schema.StringAttribute{
				Description: "Specifies the output stream to which JSON-formatted HTTP operation log messages should be written.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"log_requests": schema.BoolAttribute{
				Description: "Indicates whether to record a log message with information about requests received from the client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"log_results": schema.BoolAttribute{
				Description: "Indicates whether to record a log message with information about the result of processing a requested HTTP operation.",
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
			"include_request_details_in_result_messages": schema.BoolAttribute{
				Description: "Indicates whether result log messages should include all of the elements of request log messages. This may be used to record a single message per operation with details about both the request and response.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
			"suppressed_request_parameter_name": schema.SetAttribute{
				Description: "Specifies the case-insensitive names of request parameters that should be omitted from log messages (e.g., for the purpose of brevity or security). This will only be used if the log-request-parameters property has a value of parameter-names or parameter-names-and-values.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"log_request_protocol": schema.BoolAttribute{
				Description: "Indicates whether request log messages should include information about the HTTP version specified in the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
			"write_multi_line_messages": schema.BoolAttribute{
				Description: "Indicates whether the JSON objects should use a multi-line representation (with each object field and array value on its own line) that may be easier for administrators to read, but each message will be larger (because of additional spaces and end-of-line markers), and it may be more difficult to consume and parse through some text-oriented tools.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Publisher",
				Optional:    true,
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
func addOptionalConsoleJsonHttpOperationLogPublisherFields(ctx context.Context, addRequest *client.AddConsoleJsonHttpOperationLogPublisherRequest, plan consoleJsonHttpOperationLogPublisherResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OutputLocation) {
		outputLocation, err := client.NewEnumlogPublisherOutputLocationPropFromValue(plan.OutputLocation.ValueString())
		if err != nil {
			return err
		}
		addRequest.OutputLocation = outputLocation
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

// Read a ConsoleJsonHttpOperationLogPublisherResponse object into the model struct
func readConsoleJsonHttpOperationLogPublisherResponse(ctx context.Context, r *client.ConsoleJsonHttpOperationLogPublisherResponse, state *consoleJsonHttpOperationLogPublisherResourceModel, expectedValues *consoleJsonHttpOperationLogPublisherResourceModel, diagnostics *diag.Diagnostics) {
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
}

// Create any update operations necessary to make the state match the plan
func createConsoleJsonHttpOperationLogPublisherOperations(plan consoleJsonHttpOperationLogPublisherResourceModel, state consoleJsonHttpOperationLogPublisherResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.OutputLocation, state.OutputLocation, "output-location")
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
	operations.AddStringOperationIfNecessary(&ops, plan.LoggingErrorBehavior, state.LoggingErrorBehavior, "logging-error-behavior")
	return ops
}

// Create a new resource
func (r *consoleJsonHttpOperationLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consoleJsonHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddConsoleJsonHttpOperationLogPublisherRequest(plan.Id.ValueString(),
		[]client.EnumconsoleJsonHttpOperationLogPublisherSchemaUrn{client.ENUMCONSOLEJSONHTTPOPERATIONLOGPUBLISHERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_PUBLISHERCONSOLE_JSON_HTTP_OPERATION},
		plan.Enabled.ValueBool())
	err := addOptionalConsoleJsonHttpOperationLogPublisherFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Console Json Http Operation Log Publisher", err.Error())
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
		client.AddConsoleJsonHttpOperationLogPublisherRequestAsAddLogPublisherRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogPublisherApi.AddLogPublisherExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Console Json Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state consoleJsonHttpOperationLogPublisherResourceModel
	readConsoleJsonHttpOperationLogPublisherResponse(ctx, addResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan consoleJsonHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Console Json Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state consoleJsonHttpOperationLogPublisherResourceModel
	readConsoleJsonHttpOperationLogPublisherResponse(ctx, readResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogPublisherApi.UpdateLogPublisher(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConsoleJsonHttpOperationLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Console Json Http Operation Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsoleJsonHttpOperationLogPublisherResponse(ctx, updateResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *consoleJsonHttpOperationLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConsoleJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConsoleJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConsoleJsonHttpOperationLogPublisher(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state consoleJsonHttpOperationLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogPublisherApi.GetLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Console Json Http Operation Log Publisher", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsoleJsonHttpOperationLogPublisherResponse(ctx, readResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *consoleJsonHttpOperationLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConsoleJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConsoleJsonHttpOperationLogPublisher(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConsoleJsonHttpOperationLogPublisher(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan consoleJsonHttpOperationLogPublisherResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state consoleJsonHttpOperationLogPublisherResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogPublisherApi.UpdateLogPublisher(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConsoleJsonHttpOperationLogPublisherOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogPublisherApi.UpdateLogPublisherExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Console Json Http Operation Log Publisher", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConsoleJsonHttpOperationLogPublisherResponse(ctx, updateResponse.ConsoleJsonHttpOperationLogPublisherResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConsoleJsonHttpOperationLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *consoleJsonHttpOperationLogPublisherResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state consoleJsonHttpOperationLogPublisherResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogPublisherApi.DeleteLogPublisherExecute(r.apiClient.LogPublisherApi.DeleteLogPublisher(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Console Json Http Operation Log Publisher", err, httpResp)
		return
	}
}

func (r *consoleJsonHttpOperationLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConsoleJsonHttpOperationLogPublisher(ctx, req, resp)
}

func (r *defaultConsoleJsonHttpOperationLogPublisherResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConsoleJsonHttpOperationLogPublisher(ctx, req, resp)
}

func importConsoleJsonHttpOperationLogPublisher(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
