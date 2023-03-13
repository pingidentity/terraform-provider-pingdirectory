package connectionhandler

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
	client "github.com/pingidentity/pingdirectory-go-client/v9100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &httpConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &httpConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &httpConnectionHandlerResource{}
	_ resource.Resource                = &defaultHttpConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultHttpConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultHttpConnectionHandlerResource{}
)

// Create a Http Connection Handler resource
func NewHttpConnectionHandlerResource() resource.Resource {
	return &httpConnectionHandlerResource{}
}

func NewDefaultHttpConnectionHandlerResource() resource.Resource {
	return &defaultHttpConnectionHandlerResource{}
}

// httpConnectionHandlerResource is the resource implementation.
type httpConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultHttpConnectionHandlerResource is the resource implementation.
type defaultHttpConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *httpConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_connection_handler"
}

func (r *defaultHttpConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_http_connection_handler"
}

// Configure adds the provider configured client to the resource.
func (r *httpConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultHttpConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type httpConnectionHandlerResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	ListenAddress                   types.String `tfsdk:"listen_address"`
	ListenPort                      types.Int64  `tfsdk:"listen_port"`
	UseSSL                          types.Bool   `tfsdk:"use_ssl"`
	SslCertNickname                 types.String `tfsdk:"ssl_cert_nickname"`
	HttpServletExtension            types.Set    `tfsdk:"http_servlet_extension"`
	WebApplicationExtension         types.Set    `tfsdk:"web_application_extension"`
	HttpOperationLogPublisher       types.Set    `tfsdk:"http_operation_log_publisher"`
	SslProtocol                     types.Set    `tfsdk:"ssl_protocol"`
	SslCipherSuite                  types.Set    `tfsdk:"ssl_cipher_suite"`
	KeyManagerProvider              types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider            types.String `tfsdk:"trust_manager_provider"`
	NumRequestHandlers              types.Int64  `tfsdk:"num_request_handlers"`
	KeepStats                       types.Bool   `tfsdk:"keep_stats"`
	AcceptBacklog                   types.Int64  `tfsdk:"accept_backlog"`
	AllowTCPReuseAddress            types.Bool   `tfsdk:"allow_tcp_reuse_address"`
	IdleTimeLimit                   types.String `tfsdk:"idle_time_limit"`
	LowResourcesConnectionThreshold types.Int64  `tfsdk:"low_resources_connection_threshold"`
	LowResourcesIdleTimeLimit       types.String `tfsdk:"low_resources_idle_time_limit"`
	EnableMultipartMIMEParameters   types.Bool   `tfsdk:"enable_multipart_mime_parameters"`
	UseForwardedHeaders             types.Bool   `tfsdk:"use_forwarded_headers"`
	HttpRequestHeaderSize           types.Int64  `tfsdk:"http_request_header_size"`
	ResponseHeader                  types.Set    `tfsdk:"response_header"`
	UseCorrelationIDHeader          types.Bool   `tfsdk:"use_correlation_id_header"`
	CorrelationIDResponseHeader     types.String `tfsdk:"correlation_id_response_header"`
	CorrelationIDRequestHeader      types.Set    `tfsdk:"correlation_id_request_header"`
	SslClientAuthPolicy             types.String `tfsdk:"ssl_client_auth_policy"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *httpConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpConnectionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultHttpConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpConnectionHandlerSchema(ctx, req, resp, true)
}

func httpConnectionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Http Connection Handler.",
		Attributes: map[string]schema.Attribute{
			"listen_address": schema.StringAttribute{
				Description: "Specifies the address on which to listen for connections from HTTP clients. If no value is defined, the server will listen on all addresses on all interfaces.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"listen_port": schema.Int64Attribute{
				Description: "Specifies the port number on which the HTTP Connection Handler will listen for connections from clients.",
				Required:    true,
			},
			"use_ssl": schema.BoolAttribute{
				Description: "Indicates whether the HTTP Connection Handler should use SSL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "Specifies the nickname (also called the alias) of the certificate that the HTTP Connection Handler should use when performing SSL communication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"http_servlet_extension": schema.SetAttribute{
				Description: "Specifies information about servlets that will be provided via this connection handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"web_application_extension": schema.SetAttribute{
				Description: "Specifies information about web applications that will be provided via this connection handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"http_operation_log_publisher": schema.SetAttribute{
				Description: "Specifies the set of HTTP operation loggers that should be used to log information about requests and responses for operations processed through this HTTP Connection Handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"ssl_protocol": schema.SetAttribute{
				Description: "Specifies the names of the SSL protocols that are allowed for use in SSL communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the SSL cipher suites that are allowed for use in SSL communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies the key manager provider that will be used to obtain the certificate to present to HTTPS clients.",
				Optional:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies the trust manager provider that will be used to validate any certificates presented by HTTPS clients.",
				Optional:    true,
			},
			"num_request_handlers": schema.Int64Attribute{
				Description: "Specifies the number of threads that will be used for accepting connections and reading requests from clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"keep_stats": schema.BoolAttribute{
				Description: "Indicates whether to enable statistics collection for this connection handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"accept_backlog": schema.Int64Attribute{
				Description: "Specifies the number of concurrent outstanding connection attempts that the connection handler should allow. The default value should be acceptable in most cases, but it may need to be increased in environments that may attempt to establish large numbers of connections simultaneously.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"allow_tcp_reuse_address": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to reuse socket descriptors. This may be useful in environments with a high rate of connection establishment and termination.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"idle_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum idle time for a connection. The max idle time is applied when waiting for a new request to be received on a connection, when reading the headers and content of a request, or when writing the headers and content of a response.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"low_resources_connection_threshold": schema.Int64Attribute{
				Description: "Specifies the number of connections, which if exceeded, places this handler in a low resource state where a different idle time limit is applied on the connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"low_resources_idle_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum idle time for a connection when this handler is in a low resource state as defined by low-resource-connections. The max idle time is applied when waiting for a new request to be received on a connection, when reading the headers and content of a request, or when writing the headers and content of a response.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_multipart_mime_parameters": schema.BoolAttribute{
				Description: "Determines whether request form parameters submitted in multipart/ form-data (RFC 2388) format should be processed as request parameters.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"use_forwarded_headers": schema.BoolAttribute{
				Description: "Indicates whether to use \"Forwarded\" and \"X-Forwarded-*\" request headers to override corresponding HTTP request information available during request processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"http_request_header_size": schema.Int64Attribute{
				Description: "Specifies the maximum buffer size of an http request including the request uri and all of the request headers.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"use_correlation_id_header": schema.BoolAttribute{
				Description: "If enabled, a correlation ID header will be added to outgoing HTTP responses.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"correlation_id_request_header": schema.SetAttribute{
				Description: "Specifies the set of HTTP request headers that may contain a value to be used as the correlation ID. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"ssl_client_auth_policy": schema.StringAttribute{
				Description: "Specifies the policy that the HTTP Connection Handler should use regarding client SSL certificates. In order for a client certificate to be accepted it must be known to the trust-manager-provider associated with this HTTP Connection Handler. Client certificates received by the HTTP Connection Handler are by default used for TLS mutual authentication only, as there is no support for user authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Connection Handler is enabled.",
				Required:    true,
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
func addOptionalHttpConnectionHandlerFields(ctx context.Context, addRequest *client.AddHttpConnectionHandlerRequest, plan httpConnectionHandlerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ListenAddress) {
		stringVal := plan.ListenAddress.ValueString()
		addRequest.ListenAddress = &stringVal
	}
	if internaltypes.IsDefined(plan.UseSSL) {
		boolVal := plan.UseSSL.ValueBool()
		addRequest.UseSSL = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslCertNickname) {
		stringVal := plan.SslCertNickname.ValueString()
		addRequest.SslCertNickname = &stringVal
	}
	if internaltypes.IsDefined(plan.HttpServletExtension) {
		var slice []string
		plan.HttpServletExtension.ElementsAs(ctx, &slice, false)
		addRequest.HttpServletExtension = slice
	}
	if internaltypes.IsDefined(plan.WebApplicationExtension) {
		var slice []string
		plan.WebApplicationExtension.ElementsAs(ctx, &slice, false)
		addRequest.WebApplicationExtension = slice
	}
	if internaltypes.IsDefined(plan.HttpOperationLogPublisher) {
		var slice []string
		plan.HttpOperationLogPublisher.ElementsAs(ctx, &slice, false)
		addRequest.HttpOperationLogPublisher = slice
	}
	if internaltypes.IsDefined(plan.SslProtocol) {
		var slice []string
		plan.SslProtocol.ElementsAs(ctx, &slice, false)
		addRequest.SslProtocol = slice
	}
	if internaltypes.IsDefined(plan.SslCipherSuite) {
		var slice []string
		plan.SslCipherSuite.ElementsAs(ctx, &slice, false)
		addRequest.SslCipherSuite = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		stringVal := plan.KeyManagerProvider.ValueString()
		addRequest.KeyManagerProvider = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		stringVal := plan.TrustManagerProvider.ValueString()
		addRequest.TrustManagerProvider = &stringVal
	}
	if internaltypes.IsDefined(plan.NumRequestHandlers) {
		intVal := int32(plan.NumRequestHandlers.ValueInt64())
		addRequest.NumRequestHandlers = &intVal
	}
	if internaltypes.IsDefined(plan.KeepStats) {
		boolVal := plan.KeepStats.ValueBool()
		addRequest.KeepStats = &boolVal
	}
	if internaltypes.IsDefined(plan.AcceptBacklog) {
		intVal := int32(plan.AcceptBacklog.ValueInt64())
		addRequest.AcceptBacklog = &intVal
	}
	if internaltypes.IsDefined(plan.AllowTCPReuseAddress) {
		boolVal := plan.AllowTCPReuseAddress.ValueBool()
		addRequest.AllowTCPReuseAddress = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdleTimeLimit) {
		stringVal := plan.IdleTimeLimit.ValueString()
		addRequest.IdleTimeLimit = &stringVal
	}
	if internaltypes.IsDefined(plan.LowResourcesConnectionThreshold) {
		intVal := int32(plan.LowResourcesConnectionThreshold.ValueInt64())
		addRequest.LowResourcesConnectionThreshold = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LowResourcesIdleTimeLimit) {
		stringVal := plan.LowResourcesIdleTimeLimit.ValueString()
		addRequest.LowResourcesIdleTimeLimit = &stringVal
	}
	if internaltypes.IsDefined(plan.EnableMultipartMIMEParameters) {
		boolVal := plan.EnableMultipartMIMEParameters.ValueBool()
		addRequest.EnableMultipartMIMEParameters = &boolVal
	}
	if internaltypes.IsDefined(plan.UseForwardedHeaders) {
		boolVal := plan.UseForwardedHeaders.ValueBool()
		addRequest.UseForwardedHeaders = &boolVal
	}
	if internaltypes.IsDefined(plan.HttpRequestHeaderSize) {
		intVal := int32(plan.HttpRequestHeaderSize.ValueInt64())
		addRequest.HttpRequestHeaderSize = &intVal
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	if internaltypes.IsDefined(plan.UseCorrelationIDHeader) {
		boolVal := plan.UseCorrelationIDHeader.ValueBool()
		addRequest.UseCorrelationIDHeader = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		stringVal := plan.CorrelationIDResponseHeader.ValueString()
		addRequest.CorrelationIDResponseHeader = &stringVal
	}
	if internaltypes.IsDefined(plan.CorrelationIDRequestHeader) {
		var slice []string
		plan.CorrelationIDRequestHeader.ElementsAs(ctx, &slice, false)
		addRequest.CorrelationIDRequestHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslClientAuthPolicy) {
		sslClientAuthPolicy, err := client.NewEnumconnectionHandlerSslClientAuthPolicyPropFromValue(plan.SslClientAuthPolicy.ValueString())
		if err != nil {
			return err
		}
		addRequest.SslClientAuthPolicy = sslClientAuthPolicy
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a HttpConnectionHandlerResponse object into the model struct
func readHttpConnectionHandlerResponse(ctx context.Context, r *client.HttpConnectionHandlerResponse, state *httpConnectionHandlerResourceModel, expectedValues *httpConnectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ListenAddress = internaltypes.StringTypeOrNil(r.ListenAddress, internaltypes.IsEmptyString(expectedValues.ListenAddress))
	state.ListenPort = types.Int64Value(int64(r.ListenPort))
	state.UseSSL = internaltypes.BoolTypeOrNil(r.UseSSL)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, internaltypes.IsEmptyString(expectedValues.SslCertNickname))
	state.HttpServletExtension = internaltypes.GetStringSet(r.HttpServletExtension)
	state.WebApplicationExtension = internaltypes.GetStringSet(r.WebApplicationExtension)
	state.HttpOperationLogPublisher = internaltypes.GetStringSet(r.HttpOperationLogPublisher)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.NumRequestHandlers = internaltypes.Int64TypeOrNil(r.NumRequestHandlers)
	state.KeepStats = internaltypes.BoolTypeOrNil(r.KeepStats)
	state.AcceptBacklog = internaltypes.Int64TypeOrNil(r.AcceptBacklog)
	state.AllowTCPReuseAddress = internaltypes.BoolTypeOrNil(r.AllowTCPReuseAddress)
	state.IdleTimeLimit = internaltypes.StringTypeOrNil(r.IdleTimeLimit, internaltypes.IsEmptyString(expectedValues.IdleTimeLimit))
	config.CheckMismatchedPDFormattedAttributes("idle_time_limit",
		expectedValues.IdleTimeLimit, state.IdleTimeLimit, diagnostics)
	state.LowResourcesConnectionThreshold = internaltypes.Int64TypeOrNil(r.LowResourcesConnectionThreshold)
	state.LowResourcesIdleTimeLimit = internaltypes.StringTypeOrNil(r.LowResourcesIdleTimeLimit, internaltypes.IsEmptyString(expectedValues.LowResourcesIdleTimeLimit))
	config.CheckMismatchedPDFormattedAttributes("low_resources_idle_time_limit",
		expectedValues.LowResourcesIdleTimeLimit, state.LowResourcesIdleTimeLimit, diagnostics)
	state.EnableMultipartMIMEParameters = internaltypes.BoolTypeOrNil(r.EnableMultipartMIMEParameters)
	state.UseForwardedHeaders = internaltypes.BoolTypeOrNil(r.UseForwardedHeaders)
	state.HttpRequestHeaderSize = internaltypes.Int64TypeOrNil(r.HttpRequestHeaderSize)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.UseCorrelationIDHeader = internaltypes.BoolTypeOrNil(r.UseCorrelationIDHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.CorrelationIDRequestHeader = internaltypes.GetStringSet(r.CorrelationIDRequestHeader)
	state.SslClientAuthPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionHandlerSslClientAuthPolicyProp(r.SslClientAuthPolicy), internaltypes.IsEmptyString(expectedValues.SslClientAuthPolicy))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createHttpConnectionHandlerOperations(plan httpConnectionHandlerResourceModel, state httpConnectionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ListenAddress, state.ListenAddress, "listen-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.ListenPort, state.ListenPort, "listen-port")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSSL, state.UseSSL, "use-ssl")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HttpServletExtension, state.HttpServletExtension, "http-servlet-extension")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WebApplicationExtension, state.WebApplicationExtension, "web-application-extension")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HttpOperationLogPublisher, state.HttpOperationLogPublisher, "http-operation-log-publisher")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SslProtocol, state.SslProtocol, "ssl-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SslCipherSuite, state.SslCipherSuite, "ssl-cipher-suite")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumRequestHandlers, state.NumRequestHandlers, "num-request-handlers")
	operations.AddBoolOperationIfNecessary(&ops, plan.KeepStats, state.KeepStats, "keep-stats")
	operations.AddInt64OperationIfNecessary(&ops, plan.AcceptBacklog, state.AcceptBacklog, "accept-backlog")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowTCPReuseAddress, state.AllowTCPReuseAddress, "allow-tcp-reuse-address")
	operations.AddStringOperationIfNecessary(&ops, plan.IdleTimeLimit, state.IdleTimeLimit, "idle-time-limit")
	operations.AddInt64OperationIfNecessary(&ops, plan.LowResourcesConnectionThreshold, state.LowResourcesConnectionThreshold, "low-resources-connection-threshold")
	operations.AddStringOperationIfNecessary(&ops, plan.LowResourcesIdleTimeLimit, state.LowResourcesIdleTimeLimit, "low-resources-idle-time-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableMultipartMIMEParameters, state.EnableMultipartMIMEParameters, "enable-multipart-mime-parameters")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseForwardedHeaders, state.UseForwardedHeaders, "use-forwarded-headers")
	operations.AddInt64OperationIfNecessary(&ops, plan.HttpRequestHeaderSize, state.HttpRequestHeaderSize, "http-request-header-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseCorrelationIDHeader, state.UseCorrelationIDHeader, "use-correlation-id-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CorrelationIDRequestHeader, state.CorrelationIDRequestHeader, "correlation-id-request-header")
	operations.AddStringOperationIfNecessary(&ops, plan.SslClientAuthPolicy, state.SslClientAuthPolicy, "ssl-client-auth-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *httpConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddHttpConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumhttpConnectionHandlerSchemaUrn{client.ENUMHTTPCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERHTTP},
		int32(plan.ListenPort.ValueInt64()),
		plan.Enabled.ValueBool())
	err := addOptionalHttpConnectionHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http Connection Handler", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionHandlerApi.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddHttpConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.AddConnectionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpConnectionHandlerResourceModel
	readHttpConnectionHandlerResponse(ctx, addResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultHttpConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state httpConnectionHandlerResourceModel
	readHttpConnectionHandlerResponse(ctx, readResponse.HttpConnectionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createHttpConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpConnectionHandlerResponse(ctx, updateResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *httpConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultHttpConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readHttpConnectionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state httpConnectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpConnectionHandlerResponse(ctx, readResponse.HttpConnectionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *httpConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultHttpConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateHttpConnectionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan httpConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state httpConnectionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionHandlerApi.UpdateConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createHttpConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpConnectionHandlerResponse(ctx, updateResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultHttpConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *httpConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state httpConnectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionHandlerApi.DeleteConnectionHandlerExecute(r.apiClient.ConnectionHandlerApi.DeleteConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Http Connection Handler", err, httpResp)
		return
	}
}

func (r *httpConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpConnectionHandler(ctx, req, resp)
}

func (r *defaultHttpConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpConnectionHandler(ctx, req, resp)
}

func importHttpConnectionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
