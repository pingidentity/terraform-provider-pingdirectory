package connectionhandler

import (
	"context"
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
	_ resource.Resource                = &connectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &connectionHandlerResource{}
	_ resource.ResourceWithImportState = &connectionHandlerResource{}
	_ resource.Resource                = &defaultConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultConnectionHandlerResource{}
)

// Create a Connection Handler resource
func NewConnectionHandlerResource() resource.Resource {
	return &connectionHandlerResource{}
}

func NewDefaultConnectionHandlerResource() resource.Resource {
	return &defaultConnectionHandlerResource{}
}

// connectionHandlerResource is the resource implementation.
type connectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConnectionHandlerResource is the resource implementation.
type defaultConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *connectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection_handler"
}

func (r *defaultConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_connection_handler"
}

// Configure adds the provider configured client to the resource.
func (r *connectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type connectionHandlerResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	LastUpdated                            types.String `tfsdk:"last_updated"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	Type                                   types.String `tfsdk:"type"`
	ListenAddress                          types.Set    `tfsdk:"listen_address"`
	ListenPort                             types.Int64  `tfsdk:"listen_port"`
	LdifDirectory                          types.String `tfsdk:"ldif_directory"`
	PollInterval                           types.String `tfsdk:"poll_interval"`
	HttpServletExtension                   types.Set    `tfsdk:"http_servlet_extension"`
	WebApplicationExtension                types.Set    `tfsdk:"web_application_extension"`
	HttpOperationLogPublisher              types.Set    `tfsdk:"http_operation_log_publisher"`
	UseSSL                                 types.Bool   `tfsdk:"use_ssl"`
	AllowStartTLS                          types.Bool   `tfsdk:"allow_start_tls"`
	SslCertNickname                        types.String `tfsdk:"ssl_cert_nickname"`
	KeyManagerProvider                     types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider                   types.String `tfsdk:"trust_manager_provider"`
	KeepStats                              types.Bool   `tfsdk:"keep_stats"`
	AllowLDAPV2                            types.Bool   `tfsdk:"allow_ldap_v2"`
	AllowTCPReuseAddress                   types.Bool   `tfsdk:"allow_tcp_reuse_address"`
	IdleTimeLimit                          types.String `tfsdk:"idle_time_limit"`
	LowResourcesConnectionThreshold        types.Int64  `tfsdk:"low_resources_connection_threshold"`
	LowResourcesIdleTimeLimit              types.String `tfsdk:"low_resources_idle_time_limit"`
	EnableMultipartMIMEParameters          types.Bool   `tfsdk:"enable_multipart_mime_parameters"`
	UseForwardedHeaders                    types.Bool   `tfsdk:"use_forwarded_headers"`
	HttpRequestHeaderSize                  types.Int64  `tfsdk:"http_request_header_size"`
	ResponseHeader                         types.Set    `tfsdk:"response_header"`
	UseCorrelationIDHeader                 types.Bool   `tfsdk:"use_correlation_id_header"`
	CorrelationIDResponseHeader            types.String `tfsdk:"correlation_id_response_header"`
	CorrelationIDRequestHeader             types.Set    `tfsdk:"correlation_id_request_header"`
	UseTCPKeepAlive                        types.Bool   `tfsdk:"use_tcp_keep_alive"`
	SendRejectionNotice                    types.Bool   `tfsdk:"send_rejection_notice"`
	FailedBindResponseDelay                types.String `tfsdk:"failed_bind_response_delay"`
	MaxRequestSize                         types.String `tfsdk:"max_request_size"`
	MaxCancelHandlers                      types.Int64  `tfsdk:"max_cancel_handlers"`
	NumAcceptHandlers                      types.Int64  `tfsdk:"num_accept_handlers"`
	NumRequestHandlers                     types.Int64  `tfsdk:"num_request_handlers"`
	SslClientAuthPolicy                    types.String `tfsdk:"ssl_client_auth_policy"`
	AcceptBacklog                          types.Int64  `tfsdk:"accept_backlog"`
	SslProtocol                            types.Set    `tfsdk:"ssl_protocol"`
	SslCipherSuite                         types.Set    `tfsdk:"ssl_cipher_suite"`
	MaxBlockedWriteTimeLimit               types.String `tfsdk:"max_blocked_write_time_limit"`
	AutoAuthenticateUsingClientCertificate types.Bool   `tfsdk:"auto_authenticate_using_client_certificate"`
	CloseConnectionsWhenUnavailable        types.Bool   `tfsdk:"close_connections_when_unavailable"`
	CloseConnectionsOnExplicitGC           types.Bool   `tfsdk:"close_connections_on_explicit_gc"`
	Description                            types.String `tfsdk:"description"`
	Enabled                                types.Bool   `tfsdk:"enabled"`
	AllowedClient                          types.Set    `tfsdk:"allowed_client"`
	DeniedClient                           types.Set    `tfsdk:"denied_client"`
}

// GetSchema defines the schema for the resource.
func (r *connectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	connectionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	connectionHandlerSchema(ctx, req, resp, true)
}

func connectionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Connection Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Connection Handler resource. Options are ['jmx', 'ldap', 'ldif', 'http']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"jmx", "ldap", "ldif", "http"}...),
				},
			},
			"listen_address": schema.SetAttribute{
				Description: "Specifies the address or set of addresses on which this LDAP Connection Handler should listen for connections from LDAP clients.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"listen_port": schema.Int64Attribute{
				Description: "Specifies the port number on which the JMX Connection Handler will listen for connections from clients.",
				Optional:    true,
			},
			"ldif_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory in which the LDIF files should be placed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"poll_interval": schema.StringAttribute{
				Description: "Specifies how frequently the LDIF connection handler should check the LDIF directory to determine whether a new LDIF file has been added.",
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
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"web_application_extension": schema.SetAttribute{
				Description: "Specifies information about web applications that will be provided via this connection handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"http_operation_log_publisher": schema.SetAttribute{
				Description: "Specifies the set of HTTP operation loggers that should be used to log information about requests and responses for operations processed through this HTTP Connection Handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"use_ssl": schema.BoolAttribute{
				Description: "Indicates whether the JMX Connection Handler should use SSL.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_start_tls": schema.BoolAttribute{
				Description: "Indicates whether clients are allowed to use StartTLS.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "Specifies the nickname (also called the alias) of the certificate that the JMX Connection Handler should use when performing SSL communication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies the name of the key manager that should be used with this JMX Connection Handler .",
				Optional:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies the name of the trust manager that should be used with the LDAP Connection Handler .",
				Optional:    true,
			},
			"keep_stats": schema.BoolAttribute{
				Description: "Indicates whether to enable statistics collection for this connection handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"allow_ldap_v2": schema.BoolAttribute{
				Description: "Indicates whether connections from LDAPv2 clients are allowed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"use_tcp_keep_alive": schema.BoolAttribute{
				Description: "Indicates whether the LDAP Connection Handler should use TCP keep-alive.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"send_rejection_notice": schema.BoolAttribute{
				Description: "Indicates whether the LDAP Connection Handler should send a notice of disconnection extended response message to the client if a new connection is rejected for some reason.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"failed_bind_response_delay": schema.StringAttribute{
				Description: "Specifies the length of time that the server should delay the response to non-successful bind operations. A value of zero milliseconds indicates that non-successful bind operations should not be delayed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_request_size": schema.StringAttribute{
				Description: "Specifies the size of the largest LDAP request message that will be allowed by this LDAP Connection handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_cancel_handlers": schema.Int64Attribute{
				Description: "Specifies the maximum number of threads that are used to process cancel and abandon requests from clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_accept_handlers": schema.Int64Attribute{
				Description: "Specifies the number of threads that are used to accept new client connections, and to perform any initial preparation on those connections that may be needed before the connection can be used to read requests and send responses.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"num_request_handlers": schema.Int64Attribute{
				Description: "Specifies the number of request handlers that are used to read requests from clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ssl_client_auth_policy": schema.StringAttribute{
				Description: "Specifies the policy that the LDAP Connection Handler should use regarding client SSL certificates.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"accept_backlog": schema.Int64Attribute{
				Description: "Specifies the maximum number of pending connection attempts that are allowed to queue up in the accept backlog before the server starts rejecting new connection attempts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ssl_protocol": schema.SetAttribute{
				Description: "Specifies the names of the SSL protocols that are allowed for use in SSL or StartTLS communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the SSL cipher suites that are allowed for use in SSL or StartTLS communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"max_blocked_write_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that attempts to write data to LDAP clients should be allowed to block.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_authenticate_using_client_certificate": schema.BoolAttribute{
				Description: "Indicates whether to attempt to automatically authenticate a client connection that has established a secure communication channel (using either SSL or StartTLS) and presented its own client certificate. Generally, clients should use the SASL EXTERNAL mechanism to authenticate using a client certificate, but some clients may not support that capability and/or may expect automatic authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"close_connections_when_unavailable": schema.BoolAttribute{
				Description: "Indicates whether all connections associated with this LDAP Connection Handler should be closed and no new connections accepted when the server has determined that it is \"unavailable.\" This allows clients (or a network load balancer) to route requests to another server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"close_connections_on_explicit_gc": schema.BoolAttribute{
				Description: "Indicates whether all connections associated with this LDAP Connection Handler should be closed before an explicit garbage collection is performed to allow clients to route requests to another server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
			"allowed_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are allowed to establish connections to this connection handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"denied_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are not allowed to establish connections to this connection handler.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"jmx", "ldap", "ldif", "http"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	config.AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *connectionHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConnectionHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanConnectionHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	var model connectionHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.KeepStats) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'keep_stats' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'keep_stats', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.EnableMultipartMIMEParameters) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'enable_multipart_mime_parameters' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'enable_multipart_mime_parameters', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.AllowLDAPV2) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'allow_ldap_v2' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allow_ldap_v2', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.SslCipherSuite) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'ssl_cipher_suite' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ssl_cipher_suite', the 'type' attribute must be one of ['ldap', 'http']")
	}
	if internaltypes.IsDefined(model.WebApplicationExtension) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'web_application_extension' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'web_application_extension', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.HttpOperationLogPublisher) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'http_operation_log_publisher' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'http_operation_log_publisher', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.NumAcceptHandlers) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'num_accept_handlers' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'num_accept_handlers', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.SendRejectionNotice) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'send_rejection_notice' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'send_rejection_notice', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.PollInterval) && model.Type.ValueString() != "ldif" {
		resp.Diagnostics.AddError("Attribute 'poll_interval' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'poll_interval', the 'type' attribute must be one of ['ldif']")
	}
	if internaltypes.IsDefined(model.KeyManagerProvider) && model.Type.ValueString() != "jmx" && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'key_manager_provider' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'key_manager_provider', the 'type' attribute must be one of ['jmx', 'ldap', 'http']")
	}
	if internaltypes.IsDefined(model.AcceptBacklog) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'accept_backlog' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'accept_backlog', the 'type' attribute must be one of ['ldap', 'http']")
	}
	if internaltypes.IsDefined(model.DeniedClient) && model.Type.ValueString() != "jmx" && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "ldif" {
		resp.Diagnostics.AddError("Attribute 'denied_client' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'denied_client', the 'type' attribute must be one of ['jmx', 'ldap', 'ldif']")
	}
	if internaltypes.IsDefined(model.TrustManagerProvider) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'trust_manager_provider' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'trust_manager_provider', the 'type' attribute must be one of ['ldap', 'http']")
	}
	if internaltypes.IsDefined(model.AllowTCPReuseAddress) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'allow_tcp_reuse_address' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allow_tcp_reuse_address', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.LowResourcesIdleTimeLimit) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'low_resources_idle_time_limit' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'low_resources_idle_time_limit', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.CorrelationIDResponseHeader) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'correlation_id_response_header' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'correlation_id_response_header', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.ListenPort) && model.Type.ValueString() != "jmx" && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'listen_port' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'listen_port', the 'type' attribute must be one of ['jmx', 'ldap', 'http']")
	}
	if internaltypes.IsDefined(model.UseSSL) && model.Type.ValueString() != "jmx" && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'use_ssl' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_ssl', the 'type' attribute must be one of ['jmx', 'ldap', 'http']")
	}
	if internaltypes.IsDefined(model.ResponseHeader) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'response_header' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'response_header', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.AllowedClient) && model.Type.ValueString() != "jmx" && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "ldif" {
		resp.Diagnostics.AddError("Attribute 'allowed_client' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allowed_client', the 'type' attribute must be one of ['jmx', 'ldap', 'ldif']")
	}
	if internaltypes.IsDefined(model.MaxRequestSize) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'max_request_size' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_request_size', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.MaxBlockedWriteTimeLimit) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'max_blocked_write_time_limit' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_blocked_write_time_limit', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.CorrelationIDRequestHeader) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'correlation_id_request_header' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'correlation_id_request_header', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.CloseConnectionsWhenUnavailable) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'close_connections_when_unavailable' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'close_connections_when_unavailable', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.ListenAddress) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'listen_address' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'listen_address', the 'type' attribute must be one of ['ldap', 'http']")
	}
	if internaltypes.IsDefined(model.MaxCancelHandlers) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'max_cancel_handlers' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'max_cancel_handlers', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.HttpServletExtension) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'http_servlet_extension' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'http_servlet_extension', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.AutoAuthenticateUsingClientCertificate) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'auto_authenticate_using_client_certificate' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'auto_authenticate_using_client_certificate', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.LowResourcesConnectionThreshold) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'low_resources_connection_threshold' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'low_resources_connection_threshold', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.SslClientAuthPolicy) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'ssl_client_auth_policy' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ssl_client_auth_policy', the 'type' attribute must be one of ['ldap', 'http']")
	}
	if internaltypes.IsDefined(model.UseForwardedHeaders) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'use_forwarded_headers' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_forwarded_headers', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.SslProtocol) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'ssl_protocol' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ssl_protocol', the 'type' attribute must be one of ['ldap', 'http']")
	}
	if internaltypes.IsDefined(model.AllowStartTLS) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'allow_start_tls' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'allow_start_tls', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.UseTCPKeepAlive) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'use_tcp_keep_alive' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_tcp_keep_alive', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.LdifDirectory) && model.Type.ValueString() != "ldif" {
		resp.Diagnostics.AddError("Attribute 'ldif_directory' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ldif_directory', the 'type' attribute must be one of ['ldif']")
	}
	if internaltypes.IsDefined(model.SslCertNickname) && model.Type.ValueString() != "jmx" && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'ssl_cert_nickname' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'ssl_cert_nickname', the 'type' attribute must be one of ['jmx', 'ldap', 'http']")
	}
	if internaltypes.IsDefined(model.UseCorrelationIDHeader) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'use_correlation_id_header' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'use_correlation_id_header', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.IdleTimeLimit) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'idle_time_limit' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'idle_time_limit', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.HttpRequestHeaderSize) && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'http_request_header_size' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'http_request_header_size', the 'type' attribute must be one of ['http']")
	}
	if internaltypes.IsDefined(model.FailedBindResponseDelay) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'failed_bind_response_delay' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'failed_bind_response_delay', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.CloseConnectionsOnExplicitGC) && model.Type.ValueString() != "ldap" {
		resp.Diagnostics.AddError("Attribute 'close_connections_on_explicit_gc' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'close_connections_on_explicit_gc', the 'type' attribute must be one of ['ldap']")
	}
	if internaltypes.IsDefined(model.NumRequestHandlers) && model.Type.ValueString() != "ldap" && model.Type.ValueString() != "http" {
		resp.Diagnostics.AddError("Attribute 'num_request_handlers' not supported by pingdirectory_connection_handler resources with 'type' '"+model.Type.ValueString()+"'",
			"When using attribute 'num_request_handlers', the 'type' attribute must be one of ['ldap', 'http']")
	}
}

// Add optional fields to create request for jmx connection-handler
func addOptionalJmxConnectionHandlerFields(ctx context.Context, addRequest *client.AddJmxConnectionHandlerRequest, plan connectionHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.UseSSL) {
		addRequest.UseSSL = plan.UseSSL.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslCertNickname) {
		addRequest.SslCertNickname = plan.SslCertNickname.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowedClient) {
		var slice []string
		plan.AllowedClient.ElementsAs(ctx, &slice, false)
		addRequest.AllowedClient = slice
	}
	if internaltypes.IsDefined(plan.DeniedClient) {
		var slice []string
		plan.DeniedClient.ElementsAs(ctx, &slice, false)
		addRequest.DeniedClient = slice
	}
	return nil
}

// Add optional fields to create request for ldap connection-handler
func addOptionalLdapConnectionHandlerFields(ctx context.Context, addRequest *client.AddLdapConnectionHandlerRequest, plan connectionHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.ListenAddress) {
		var slice []string
		plan.ListenAddress.ElementsAs(ctx, &slice, false)
		addRequest.ListenAddress = slice
	}
	if internaltypes.IsDefined(plan.UseSSL) {
		addRequest.UseSSL = plan.UseSSL.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AllowStartTLS) {
		addRequest.AllowStartTLS = plan.AllowStartTLS.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslCertNickname) {
		addRequest.SslCertNickname = plan.SslCertNickname.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowLDAPV2) {
		addRequest.AllowLDAPV2 = plan.AllowLDAPV2.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseTCPKeepAlive) {
		addRequest.UseTCPKeepAlive = plan.UseTCPKeepAlive.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.SendRejectionNotice) {
		addRequest.SendRejectionNotice = plan.SendRejectionNotice.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.FailedBindResponseDelay) {
		addRequest.FailedBindResponseDelay = plan.FailedBindResponseDelay.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxRequestSize) {
		addRequest.MaxRequestSize = plan.MaxRequestSize.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.MaxCancelHandlers) {
		addRequest.MaxCancelHandlers = plan.MaxCancelHandlers.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumAcceptHandlers) {
		addRequest.NumAcceptHandlers = plan.NumAcceptHandlers.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.NumRequestHandlers) {
		addRequest.NumRequestHandlers = plan.NumRequestHandlers.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslClientAuthPolicy) {
		sslClientAuthPolicy, err := client.NewEnumconnectionHandlerSslClientAuthPolicyPropFromValue(plan.SslClientAuthPolicy.ValueString())
		if err != nil {
			return err
		}
		addRequest.SslClientAuthPolicy = sslClientAuthPolicy
	}
	if internaltypes.IsDefined(plan.AcceptBacklog) {
		addRequest.AcceptBacklog = plan.AcceptBacklog.ValueInt64Pointer()
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
	if internaltypes.IsNonEmptyString(plan.MaxBlockedWriteTimeLimit) {
		addRequest.MaxBlockedWriteTimeLimit = plan.MaxBlockedWriteTimeLimit.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AutoAuthenticateUsingClientCertificate) {
		addRequest.AutoAuthenticateUsingClientCertificate = plan.AutoAuthenticateUsingClientCertificate.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CloseConnectionsWhenUnavailable) {
		addRequest.CloseConnectionsWhenUnavailable = plan.CloseConnectionsWhenUnavailable.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.CloseConnectionsOnExplicitGC) {
		addRequest.CloseConnectionsOnExplicitGC = plan.CloseConnectionsOnExplicitGC.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowedClient) {
		var slice []string
		plan.AllowedClient.ElementsAs(ctx, &slice, false)
		addRequest.AllowedClient = slice
	}
	if internaltypes.IsDefined(plan.DeniedClient) {
		var slice []string
		plan.DeniedClient.ElementsAs(ctx, &slice, false)
		addRequest.DeniedClient = slice
	}
	return nil
}

// Add optional fields to create request for ldif connection-handler
func addOptionalLdifConnectionHandlerFields(ctx context.Context, addRequest *client.AddLdifConnectionHandlerRequest, plan connectionHandlerResourceModel) error {
	if internaltypes.IsDefined(plan.AllowedClient) {
		var slice []string
		plan.AllowedClient.ElementsAs(ctx, &slice, false)
		addRequest.AllowedClient = slice
	}
	if internaltypes.IsDefined(plan.DeniedClient) {
		var slice []string
		plan.DeniedClient.ElementsAs(ctx, &slice, false)
		addRequest.DeniedClient = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LdifDirectory) {
		addRequest.LdifDirectory = plan.LdifDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PollInterval) {
		addRequest.PollInterval = plan.PollInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for http connection-handler
func addOptionalHttpConnectionHandlerFields(ctx context.Context, addRequest *client.AddHttpConnectionHandlerRequest, plan connectionHandlerResourceModel) error {
	// Treat this set as a single string
	if internaltypes.IsDefined(plan.ListenAddress) && len(plan.ListenAddress.Elements()) > 0 {
		addRequest.ListenAddress = plan.ListenAddress.Elements()[0].(types.String).ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.UseSSL) {
		addRequest.UseSSL = plan.UseSSL.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslCertNickname) {
		addRequest.SslCertNickname = plan.SslCertNickname.ValueStringPointer()
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
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.NumRequestHandlers) {
		addRequest.NumRequestHandlers = plan.NumRequestHandlers.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.KeepStats) {
		addRequest.KeepStats = plan.KeepStats.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.AcceptBacklog) {
		addRequest.AcceptBacklog = plan.AcceptBacklog.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.AllowTCPReuseAddress) {
		addRequest.AllowTCPReuseAddress = plan.AllowTCPReuseAddress.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdleTimeLimit) {
		addRequest.IdleTimeLimit = plan.IdleTimeLimit.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.LowResourcesConnectionThreshold) {
		addRequest.LowResourcesConnectionThreshold = plan.LowResourcesConnectionThreshold.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.LowResourcesIdleTimeLimit) {
		addRequest.LowResourcesIdleTimeLimit = plan.LowResourcesIdleTimeLimit.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EnableMultipartMIMEParameters) {
		addRequest.EnableMultipartMIMEParameters = plan.EnableMultipartMIMEParameters.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.UseForwardedHeaders) {
		addRequest.UseForwardedHeaders = plan.UseForwardedHeaders.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.HttpRequestHeaderSize) {
		addRequest.HttpRequestHeaderSize = plan.HttpRequestHeaderSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	if internaltypes.IsDefined(plan.UseCorrelationIDHeader) {
		addRequest.UseCorrelationIDHeader = plan.UseCorrelationIDHeader.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		addRequest.CorrelationIDResponseHeader = plan.CorrelationIDResponseHeader.ValueStringPointer()
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
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any sets that have a nil ElementType, to avoid a nil pointer when setting the state
func populateConnectionHandlerNilSets(ctx context.Context, model *connectionHandlerResourceModel) {
	if model.SslCipherSuite.ElementType(ctx) == nil {
		model.SslCipherSuite = types.SetNull(types.StringType)
	}
	if model.SslProtocol.ElementType(ctx) == nil {
		model.SslProtocol = types.SetNull(types.StringType)
	}
	if model.ResponseHeader.ElementType(ctx) == nil {
		model.ResponseHeader = types.SetNull(types.StringType)
	}
	if model.AllowedClient.ElementType(ctx) == nil {
		model.AllowedClient = types.SetNull(types.StringType)
	}
	if model.WebApplicationExtension.ElementType(ctx) == nil {
		model.WebApplicationExtension = types.SetNull(types.StringType)
	}
	if model.HttpServletExtension.ElementType(ctx) == nil {
		model.HttpServletExtension = types.SetNull(types.StringType)
	}
	if model.ListenAddress.ElementType(ctx) == nil {
		model.ListenAddress = types.SetNull(types.StringType)
	}
	if model.CorrelationIDRequestHeader.ElementType(ctx) == nil {
		model.CorrelationIDRequestHeader = types.SetNull(types.StringType)
	}
	if model.HttpOperationLogPublisher.ElementType(ctx) == nil {
		model.HttpOperationLogPublisher = types.SetNull(types.StringType)
	}
	if model.DeniedClient.ElementType(ctx) == nil {
		model.DeniedClient = types.SetNull(types.StringType)
	}
}

// Read a JmxConnectionHandlerResponse object into the model struct
func readJmxConnectionHandlerResponse(ctx context.Context, r *client.JmxConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jmx")
	state.Id = types.StringValue(r.Id)
	state.ListenPort = types.Int64Value(r.ListenPort)
	state.UseSSL = internaltypes.BoolTypeOrNil(r.UseSSL)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, internaltypes.IsEmptyString(expectedValues.SslCertNickname))
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionHandlerNilSets(ctx, state)
}

// Read a LdapConnectionHandlerResponse object into the model struct
func readLdapConnectionHandlerResponse(ctx context.Context, r *client.LdapConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.ListenAddress = internaltypes.GetStringSet(r.ListenAddress)
	state.ListenPort = types.Int64Value(r.ListenPort)
	state.UseSSL = internaltypes.BoolTypeOrNil(r.UseSSL)
	state.AllowStartTLS = internaltypes.BoolTypeOrNil(r.AllowStartTLS)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, internaltypes.IsEmptyString(expectedValues.SslCertNickname))
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.AllowLDAPV2 = internaltypes.BoolTypeOrNil(r.AllowLDAPV2)
	state.UseTCPKeepAlive = internaltypes.BoolTypeOrNil(r.UseTCPKeepAlive)
	state.SendRejectionNotice = internaltypes.BoolTypeOrNil(r.SendRejectionNotice)
	state.FailedBindResponseDelay = internaltypes.StringTypeOrNil(r.FailedBindResponseDelay, internaltypes.IsEmptyString(expectedValues.FailedBindResponseDelay))
	config.CheckMismatchedPDFormattedAttributes("failed_bind_response_delay",
		expectedValues.FailedBindResponseDelay, state.FailedBindResponseDelay, diagnostics)
	state.MaxRequestSize = internaltypes.StringTypeOrNil(r.MaxRequestSize, internaltypes.IsEmptyString(expectedValues.MaxRequestSize))
	config.CheckMismatchedPDFormattedAttributes("max_request_size",
		expectedValues.MaxRequestSize, state.MaxRequestSize, diagnostics)
	state.MaxCancelHandlers = internaltypes.Int64TypeOrNil(r.MaxCancelHandlers)
	state.NumAcceptHandlers = internaltypes.Int64TypeOrNil(r.NumAcceptHandlers)
	state.NumRequestHandlers = internaltypes.Int64TypeOrNil(r.NumRequestHandlers)
	state.SslClientAuthPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionHandlerSslClientAuthPolicyProp(r.SslClientAuthPolicy), internaltypes.IsEmptyString(expectedValues.SslClientAuthPolicy))
	state.AcceptBacklog = internaltypes.Int64TypeOrNil(r.AcceptBacklog)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.MaxBlockedWriteTimeLimit = internaltypes.StringTypeOrNil(r.MaxBlockedWriteTimeLimit, internaltypes.IsEmptyString(expectedValues.MaxBlockedWriteTimeLimit))
	config.CheckMismatchedPDFormattedAttributes("max_blocked_write_time_limit",
		expectedValues.MaxBlockedWriteTimeLimit, state.MaxBlockedWriteTimeLimit, diagnostics)
	state.AutoAuthenticateUsingClientCertificate = internaltypes.BoolTypeOrNil(r.AutoAuthenticateUsingClientCertificate)
	state.CloseConnectionsWhenUnavailable = internaltypes.BoolTypeOrNil(r.CloseConnectionsWhenUnavailable)
	state.CloseConnectionsOnExplicitGC = internaltypes.BoolTypeOrNil(r.CloseConnectionsOnExplicitGC)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionHandlerNilSets(ctx, state)
}

// Read a LdifConnectionHandlerResponse object into the model struct
func readLdifConnectionHandlerResponse(ctx context.Context, r *client.LdifConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif")
	state.Id = types.StringValue(r.Id)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.PollInterval = types.StringValue(r.PollInterval)
	config.CheckMismatchedPDFormattedAttributes("poll_interval",
		expectedValues.PollInterval, state.PollInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionHandlerNilSets(ctx, state)
}

// Read a HttpConnectionHandlerResponse object into the model struct
func readHttpConnectionHandlerResponse(ctx context.Context, r *client.HttpConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http")
	state.Id = types.StringValue(r.Id)
	listenAddressValues := []string{}
	listenAddressType := internaltypes.StringTypeOrNil(r.ListenAddress, false)
	if !listenAddressType.IsNull() {
		listenAddressValues = append(listenAddressValues, listenAddressType.ValueString())
	}
	state.ListenAddress = internaltypes.GetStringSet(listenAddressValues)
	state.ListenPort = types.Int64Value(r.ListenPort)
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
	populateConnectionHandlerNilSets(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createConnectionHandlerOperations(plan connectionHandlerResourceModel, state connectionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ListenAddress, state.ListenAddress, "listen-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.ListenPort, state.ListenPort, "listen-port")
	operations.AddStringOperationIfNecessary(&ops, plan.LdifDirectory, state.LdifDirectory, "ldif-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.PollInterval, state.PollInterval, "poll-interval")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HttpServletExtension, state.HttpServletExtension, "http-servlet-extension")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.WebApplicationExtension, state.WebApplicationExtension, "web-application-extension")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.HttpOperationLogPublisher, state.HttpOperationLogPublisher, "http-operation-log-publisher")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSSL, state.UseSSL, "use-ssl")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowStartTLS, state.AllowStartTLS, "allow-start-tls")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.KeepStats, state.KeepStats, "keep-stats")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowLDAPV2, state.AllowLDAPV2, "allow-ldap-v2")
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
	operations.AddBoolOperationIfNecessary(&ops, plan.UseTCPKeepAlive, state.UseTCPKeepAlive, "use-tcp-keep-alive")
	operations.AddBoolOperationIfNecessary(&ops, plan.SendRejectionNotice, state.SendRejectionNotice, "send-rejection-notice")
	operations.AddStringOperationIfNecessary(&ops, plan.FailedBindResponseDelay, state.FailedBindResponseDelay, "failed-bind-response-delay")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxRequestSize, state.MaxRequestSize, "max-request-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxCancelHandlers, state.MaxCancelHandlers, "max-cancel-handlers")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumAcceptHandlers, state.NumAcceptHandlers, "num-accept-handlers")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumRequestHandlers, state.NumRequestHandlers, "num-request-handlers")
	operations.AddStringOperationIfNecessary(&ops, plan.SslClientAuthPolicy, state.SslClientAuthPolicy, "ssl-client-auth-policy")
	operations.AddInt64OperationIfNecessary(&ops, plan.AcceptBacklog, state.AcceptBacklog, "accept-backlog")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SslProtocol, state.SslProtocol, "ssl-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SslCipherSuite, state.SslCipherSuite, "ssl-cipher-suite")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxBlockedWriteTimeLimit, state.MaxBlockedWriteTimeLimit, "max-blocked-write-time-limit")
	operations.AddBoolOperationIfNecessary(&ops, plan.AutoAuthenticateUsingClientCertificate, state.AutoAuthenticateUsingClientCertificate, "auto-authenticate-using-client-certificate")
	operations.AddBoolOperationIfNecessary(&ops, plan.CloseConnectionsWhenUnavailable, state.CloseConnectionsWhenUnavailable, "close-connections-when-unavailable")
	operations.AddBoolOperationIfNecessary(&ops, plan.CloseConnectionsOnExplicitGC, state.CloseConnectionsOnExplicitGC, "close-connections-on-explicit-gc")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedClient, state.AllowedClient, "allowed-client")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DeniedClient, state.DeniedClient, "denied-client")
	return ops
}

// Create a jmx connection-handler
func (r *connectionHandlerResource) CreateJmxConnectionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionHandlerResourceModel) (*connectionHandlerResourceModel, error) {
	addRequest := client.NewAddJmxConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumjmxConnectionHandlerSchemaUrn{client.ENUMJMXCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERJMX},
		plan.ListenPort.ValueInt64(),
		plan.Enabled.ValueBool())
	err := addOptionalJmxConnectionHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionHandlerApi.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddJmxConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.AddConnectionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionHandlerResourceModel
	readJmxConnectionHandlerResponse(ctx, addResponse.JmxConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldap connection-handler
func (r *connectionHandlerResource) CreateLdapConnectionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionHandlerResourceModel) (*connectionHandlerResourceModel, error) {
	addRequest := client.NewAddLdapConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumldapConnectionHandlerSchemaUrn{client.ENUMLDAPCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERLDAP},
		plan.ListenPort.ValueInt64(),
		plan.Enabled.ValueBool())
	err := addOptionalLdapConnectionHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionHandlerApi.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddLdapConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.AddConnectionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionHandlerResourceModel
	readLdapConnectionHandlerResponse(ctx, addResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a ldif connection-handler
func (r *connectionHandlerResource) CreateLdifConnectionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionHandlerResourceModel) (*connectionHandlerResourceModel, error) {
	addRequest := client.NewAddLdifConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumldifConnectionHandlerSchemaUrn{client.ENUMLDIFCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERLDIF},
		plan.Enabled.ValueBool())
	err := addOptionalLdifConnectionHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Handler", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConnectionHandlerApi.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddLdifConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.AddConnectionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionHandlerResourceModel
	readLdifConnectionHandlerResponse(ctx, addResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a http connection-handler
func (r *connectionHandlerResource) CreateHttpConnectionHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan connectionHandlerResourceModel) (*connectionHandlerResourceModel, error) {
	addRequest := client.NewAddHttpConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumhttpConnectionHandlerSchemaUrn{client.ENUMHTTPCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERHTTP},
		plan.ListenPort.ValueInt64(),
		plan.Enabled.ValueBool())
	err := addOptionalHttpConnectionHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Connection Handler", err.Error())
		return nil, err
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Connection Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state connectionHandlerResourceModel
	readHttpConnectionHandlerResponse(ctx, addResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *connectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan connectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *connectionHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "jmx" {
		state, err = r.CreateJmxConnectionHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldap" {
		state, err = r.CreateLdapConnectionHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "ldif" {
		state, err = r.CreateLdifConnectionHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "http" {
		state, err = r.CreateHttpConnectionHandler(ctx, req, resp, plan)
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
func (r *defaultConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan connectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state connectionHandlerResourceModel
	if plan.Type.ValueString() == "jmx" {
		readJmxConnectionHandlerResponse(ctx, readResponse.JmxConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "ldap" {
		readLdapConnectionHandlerResponse(ctx, readResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "ldif" {
		readLdifConnectionHandlerResponse(ctx, readResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	}
	if plan.Type.ValueString() == "http" {
		readHttpConnectionHandlerResponse(ctx, readResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "jmx" {
			readJmxConnectionHandlerResponse(ctx, updateResponse.JmxConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldap" {
			readLdapConnectionHandlerResponse(ctx, updateResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldif" {
			readLdifConnectionHandlerResponse(ctx, updateResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "http" {
			readHttpConnectionHandlerResponse(ctx, updateResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *connectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConnectionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state connectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.JmxConnectionHandlerResponse != nil {
		readJmxConnectionHandlerResponse(ctx, readResponse.JmxConnectionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdapConnectionHandlerResponse != nil {
		readLdapConnectionHandlerResponse(ctx, readResponse.LdapConnectionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.LdifConnectionHandlerResponse != nil {
		readLdifConnectionHandlerResponse(ctx, readResponse.LdifConnectionHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.HttpConnectionHandlerResponse != nil {
		readHttpConnectionHandlerResponse(ctx, readResponse.HttpConnectionHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *connectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConnectionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan connectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state connectionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionHandlerApi.UpdateConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if plan.Type.ValueString() == "jmx" {
			readJmxConnectionHandlerResponse(ctx, updateResponse.JmxConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldap" {
			readLdapConnectionHandlerResponse(ctx, updateResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "ldif" {
			readLdifConnectionHandlerResponse(ctx, updateResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if plan.Type.ValueString() == "http" {
			readHttpConnectionHandlerResponse(ctx, updateResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *connectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state connectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionHandlerApi.DeleteConnectionHandlerExecute(r.apiClient.ConnectionHandlerApi.DeleteConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Connection Handler", err, httpResp)
		return
	}
}

func (r *connectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConnectionHandler(ctx, req, resp)
}

func (r *defaultConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConnectionHandler(ctx, req, resp)
}

func importConnectionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
