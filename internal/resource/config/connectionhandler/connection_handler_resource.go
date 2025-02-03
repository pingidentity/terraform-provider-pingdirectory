// Copyright © 2025 Ping Identity Corporation

package connectionhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
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
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type connectionHandlerResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	Name                                   types.String `tfsdk:"name"`
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
	EnableSniHostnameChecks                types.Bool   `tfsdk:"enable_sni_hostname_checks"`
	SendRejectionNotice                    types.Bool   `tfsdk:"send_rejection_notice"`
	FailedBindResponseDelay                types.String `tfsdk:"failed_bind_response_delay"`
	MaxRequestSize                         types.String `tfsdk:"max_request_size"`
	MaxCancelHandlers                      types.Int64  `tfsdk:"max_cancel_handlers"`
	NumAcceptHandlers                      types.Int64  `tfsdk:"num_accept_handlers"`
	NumRequestHandlers                     types.Int64  `tfsdk:"num_request_handlers"`
	RequestHandlerPerConnection            types.Bool   `tfsdk:"request_handler_per_connection"`
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
				Description:         "When the `type` attribute is set to `ldap`: Specifies the address or set of addresses on which this LDAP Connection Handler should listen for connections from LDAP clients. When the `type` attribute is set to `http`: Specifies the address on which to listen for connections from HTTP clients. If no value is defined, the server will listen on all addresses on all interfaces.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the address or set of addresses on which this LDAP Connection Handler should listen for connections from LDAP clients.\n  - `http`: Specifies the address on which to listen for connections from HTTP clients. If no value is defined, the server will listen on all addresses on all interfaces.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"listen_port": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `jmx`: Specifies the port number on which the JMX Connection Handler will listen for connections from clients. When the `type` attribute is set to `ldap`: Specifies the port number on which the LDAP Connection Handler will listen for connections from clients. When the `type` attribute is set to `http`: Specifies the port number on which the HTTP Connection Handler will listen for connections from clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Specifies the port number on which the JMX Connection Handler will listen for connections from clients.\n  - `ldap`: Specifies the port number on which the LDAP Connection Handler will listen for connections from clients.\n  - `http`: Specifies the port number on which the HTTP Connection Handler will listen for connections from clients.",
				Optional:            true,
			},
			"ldif_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory in which the LDIF files should be placed.",
				Optional:    true,
				Computed:    true,
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
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"web_application_extension": schema.SetAttribute{
				Description: "Specifies information about web applications that will be provided via this connection handler.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"http_operation_log_publisher": schema.SetAttribute{
				Description: "Specifies the set of HTTP operation loggers that should be used to log information about requests and responses for operations processed through this HTTP Connection Handler.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"use_ssl": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `jmx`: Indicates whether the JMX Connection Handler should use SSL. When the `type` attribute is set to `ldap`: Indicates whether the LDAP Connection Handler should use SSL. When the `type` attribute is set to `http`: Indicates whether the HTTP Connection Handler should use SSL.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Indicates whether the JMX Connection Handler should use SSL.\n  - `ldap`: Indicates whether the LDAP Connection Handler should use SSL.\n  - `http`: Indicates whether the HTTP Connection Handler should use SSL.",
				Optional:            true,
				Computed:            true,
			},
			"allow_start_tls": schema.BoolAttribute{
				Description: "Indicates whether clients are allowed to use StartTLS.",
				Optional:    true,
				Computed:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `jmx`: Specifies the nickname (also called the alias) of the certificate that the JMX Connection Handler should use when performing SSL communication. When the `type` attribute is set to `ldap`: Specifies the nickname (also called the alias) of the certificate that the LDAP Connection Handler should use when performing SSL or StartTLS communication. When the `type` attribute is set to `http`: Specifies the nickname (also called the alias) of the certificate that the HTTP Connection Handler should use when performing SSL communication.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Specifies the nickname (also called the alias) of the certificate that the JMX Connection Handler should use when performing SSL communication.\n  - `ldap`: Specifies the nickname (also called the alias) of the certificate that the LDAP Connection Handler should use when performing SSL or StartTLS communication.\n  - `http`: Specifies the nickname (also called the alias) of the certificate that the HTTP Connection Handler should use when performing SSL communication.",
				Optional:            true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `jmx`: Specifies the name of the key manager that should be used with this JMX Connection Handler . When the `type` attribute is set to `ldap`: Specifies the name of the key manager that should be used with this LDAP Connection Handler . When the `type` attribute is set to `http`: Specifies the key manager provider that will be used to obtain the certificate to present to HTTPS clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Specifies the name of the key manager that should be used with this JMX Connection Handler .\n  - `ldap`: Specifies the name of the key manager that should be used with this LDAP Connection Handler .\n  - `http`: Specifies the key manager provider that will be used to obtain the certificate to present to HTTPS clients.",
				Optional:            true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the name of the trust manager that should be used with the LDAP Connection Handler . When the `type` attribute is set to `http`: Specifies the trust manager provider that will be used to validate any certificates presented by HTTPS clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the name of the trust manager that should be used with the LDAP Connection Handler .\n  - `http`: Specifies the trust manager provider that will be used to validate any certificates presented by HTTPS clients.",
				Optional:            true,
			},
			"keep_stats": schema.BoolAttribute{
				Description: "Indicates whether to enable statistics collection for this connection handler.",
				Optional:    true,
				Computed:    true,
			},
			"allow_ldap_v2": schema.BoolAttribute{
				Description: "Indicates whether connections from LDAPv2 clients are allowed.",
				Optional:    true,
				Computed:    true,
			},
			"allow_tcp_reuse_address": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to reuse socket descriptors. This may be useful in environments with a high rate of connection establishment and termination.",
				Optional:    true,
				Computed:    true,
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
			},
			"use_forwarded_headers": schema.BoolAttribute{
				Description: "Indicates whether to use \"Forwarded\" and \"X-Forwarded-*\" request headers to override corresponding HTTP request information available during request processing.",
				Optional:    true,
				Computed:    true,
			},
			"http_request_header_size": schema.Int64Attribute{
				Description: "Specifies the maximum buffer size of an http request including the request uri and all of the request headers.",
				Optional:    true,
				Computed:    true,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"use_correlation_id_header": schema.BoolAttribute{
				Description: "If enabled, a correlation ID header will be added to outgoing HTTP responses.",
				Optional:    true,
				Computed:    true,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
			},
			"correlation_id_request_header": schema.SetAttribute{
				Description: "Specifies the set of HTTP request headers that may contain a value to be used as the correlation ID. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"use_tcp_keep_alive": schema.BoolAttribute{
				Description: "Indicates whether the LDAP Connection Handler should use TCP keep-alive.",
				Optional:    true,
				Computed:    true,
			},
			"enable_sni_hostname_checks": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Requires SNI hostnames to match or else throw an Invalid SNI error.",
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
			},
			"max_cancel_handlers": schema.Int64Attribute{
				Description: "Specifies the maximum number of threads that are used to process cancel and abandon requests from clients.",
				Optional:    true,
				Computed:    true,
			},
			"num_accept_handlers": schema.Int64Attribute{
				Description: "Specifies the number of threads that are used to accept new client connections, and to perform any initial preparation on those connections that may be needed before the connection can be used to read requests and send responses.",
				Optional:    true,
				Computed:    true,
			},
			"num_request_handlers": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the number of request handlers that are used to read requests from clients. When the `type` attribute is set to `http`: Specifies the number of threads that will be used for accepting connections and reading requests from clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the number of request handlers that are used to read requests from clients.\n  - `http`: Specifies the number of threads that will be used for accepting connections and reading requests from clients.",
				Optional:            true,
				Computed:            true,
			},
			"request_handler_per_connection": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Indicates whether a separate request handler thread should be created for each client connection, which can help avoid starvation of client connections for cases in which one or more clients send large numbers of concurrent asynchronous requests. This should only be used for cases in which a relatively small number of connections will be established at any given time, the connections established will generally be long-lived, and at least one client may send high volumes of asynchronous requests. This property can be used to alleviate possible blocking during long-running TLS negotiation on a single request handler which can result in it being unable to acknowledge further client requests until the TLS negotation completes or times out.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_client_auth_policy": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the policy that the LDAP Connection Handler should use regarding client SSL certificates. When the `type` attribute is set to `http`: Specifies the policy that the HTTP Connection Handler should use regarding client SSL certificates. In order for a client certificate to be accepted it must be known to the trust-manager-provider associated with this HTTP Connection Handler. Client certificates received by the HTTP Connection Handler are by default used for TLS mutual authentication only, as there is no support for user authentication.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the policy that the LDAP Connection Handler should use regarding client SSL certificates.\n  - `http`: Specifies the policy that the HTTP Connection Handler should use regarding client SSL certificates. In order for a client certificate to be accepted it must be known to the trust-manager-provider associated with this HTTP Connection Handler. Client certificates received by the HTTP Connection Handler are by default used for TLS mutual authentication only, as there is no support for user authentication.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"disabled", "optional", "required"}...),
				},
			},
			"accept_backlog": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the maximum number of pending connection attempts that are allowed to queue up in the accept backlog before the server starts rejecting new connection attempts. When the `type` attribute is set to `http`: Specifies the number of concurrent outstanding connection attempts that the connection handler should allow. The default value should be acceptable in most cases, but it may need to be increased in environments that may attempt to establish large numbers of connections simultaneously.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the maximum number of pending connection attempts that are allowed to queue up in the accept backlog before the server starts rejecting new connection attempts.\n  - `http`: Specifies the number of concurrent outstanding connection attempts that the connection handler should allow. The default value should be acceptable in most cases, but it may need to be increased in environments that may attempt to establish large numbers of connections simultaneously.",
				Optional:            true,
				Computed:            true,
			},
			"ssl_protocol": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the names of the TLS protocols that are allowed for use in SSL or StartTLS communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry. When the `type` attribute is set to `http`: Specifies the names of the SSL protocols that are allowed for use in SSL communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the names of the TLS protocols that are allowed for use in SSL or StartTLS communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.\n  - `http`: Specifies the names of the SSL protocols that are allowed for use in SSL communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the names of the TLS cipher suites that are allowed for use in SSL or StartTLS communication. The set of supported cipher suites can be viewed via the ssl context monitor entry. When the `type` attribute is set to `http`: Specifies the names of the SSL cipher suites that are allowed for use in SSL communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the names of the TLS cipher suites that are allowed for use in SSL or StartTLS communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.\n  - `http`: Specifies the names of the SSL cipher suites that are allowed for use in SSL communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
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
			},
			"close_connections_when_unavailable": schema.BoolAttribute{
				Description: "Indicates whether all connections associated with this LDAP Connection Handler should be closed and no new connections accepted when the server has determined that it is \"unavailable.\" This allows clients (or a network load balancer) to route requests to another server.",
				Optional:    true,
				Computed:    true,
			},
			"close_connections_on_explicit_gc": schema.BoolAttribute{
				Description: "Indicates whether all connections associated with this LDAP Connection Handler should be closed before an explicit garbage collection is performed to allow clients to route requests to another server.",
				Optional:    true,
				Computed:    true,
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
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"denied_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are not allowed to establish connections to this connection handler.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
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
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *connectionHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
	var planModel, configModel connectionHandlerResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for jmx type
	if resourceType == "jmx" {
		if !internaltypes.IsDefined(configModel.UseSSL) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseSSL.Equal(defaultVal) {
				planModel.UseSSL = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ldap type
	if resourceType == "ldap" {
		if !internaltypes.IsDefined(configModel.UseSSL) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseSSL.Equal(defaultVal) {
				planModel.UseSSL = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AllowStartTLS) {
			defaultVal := types.BoolValue(false)
			if !planModel.AllowStartTLS.Equal(defaultVal) {
				planModel.AllowStartTLS = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AllowLDAPV2) {
			defaultVal := types.BoolValue(true)
			if !planModel.AllowLDAPV2.Equal(defaultVal) {
				planModel.AllowLDAPV2 = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseTCPKeepAlive) {
			defaultVal := types.BoolValue(true)
			if !planModel.UseTCPKeepAlive.Equal(defaultVal) {
				planModel.UseTCPKeepAlive = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SendRejectionNotice) {
			defaultVal := types.BoolValue(true)
			if !planModel.SendRejectionNotice.Equal(defaultVal) {
				planModel.SendRejectionNotice = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxRequestSize) {
			defaultVal := types.StringValue("5 megabytes")
			if !planModel.MaxRequestSize.Equal(defaultVal) {
				planModel.MaxRequestSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxCancelHandlers) {
			defaultVal := types.Int64Value(16)
			if !planModel.MaxCancelHandlers.Equal(defaultVal) {
				planModel.MaxCancelHandlers = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.NumAcceptHandlers) {
			defaultVal := types.Int64Value(0)
			if !planModel.NumAcceptHandlers.Equal(defaultVal) {
				planModel.NumAcceptHandlers = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.NumRequestHandlers) {
			defaultVal := types.Int64Value(0)
			if !planModel.NumRequestHandlers.Equal(defaultVal) {
				planModel.NumRequestHandlers = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AcceptBacklog) {
			defaultVal := types.Int64Value(128)
			if !planModel.AcceptBacklog.Equal(defaultVal) {
				planModel.AcceptBacklog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AutoAuthenticateUsingClientCertificate) {
			defaultVal := types.BoolValue(false)
			if !planModel.AutoAuthenticateUsingClientCertificate.Equal(defaultVal) {
				planModel.AutoAuthenticateUsingClientCertificate = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CloseConnectionsWhenUnavailable) {
			defaultVal := types.BoolValue(false)
			if !planModel.CloseConnectionsWhenUnavailable.Equal(defaultVal) {
				planModel.CloseConnectionsWhenUnavailable = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CloseConnectionsOnExplicitGC) {
			defaultVal := types.BoolValue(false)
			if !planModel.CloseConnectionsOnExplicitGC.Equal(defaultVal) {
				planModel.CloseConnectionsOnExplicitGC = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for ldif type
	if resourceType == "ldif" {
		if !internaltypes.IsDefined(configModel.LdifDirectory) {
			defaultVal := types.StringValue("config/auto-process-ldif")
			if !planModel.LdifDirectory.Equal(defaultVal) {
				planModel.LdifDirectory = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for http type
	if resourceType == "http" {
		if !internaltypes.IsDefined(configModel.UseSSL) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseSSL.Equal(defaultVal) {
				planModel.UseSSL = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.NumRequestHandlers) {
			defaultVal := types.Int64Value(0)
			if !planModel.NumRequestHandlers.Equal(defaultVal) {
				planModel.NumRequestHandlers = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.KeepStats) {
			defaultVal := types.BoolValue(true)
			if !planModel.KeepStats.Equal(defaultVal) {
				planModel.KeepStats = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AcceptBacklog) {
			defaultVal := types.Int64Value(128)
			if !planModel.AcceptBacklog.Equal(defaultVal) {
				planModel.AcceptBacklog = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.AllowTCPReuseAddress) {
			defaultVal := types.BoolValue(true)
			if !planModel.AllowTCPReuseAddress.Equal(defaultVal) {
				planModel.AllowTCPReuseAddress = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.LowResourcesConnectionThreshold) {
			defaultVal := types.Int64Value(0)
			if !planModel.LowResourcesConnectionThreshold.Equal(defaultVal) {
				planModel.LowResourcesConnectionThreshold = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.EnableMultipartMIMEParameters) {
			defaultVal := types.BoolValue(false)
			if !planModel.EnableMultipartMIMEParameters.Equal(defaultVal) {
				planModel.EnableMultipartMIMEParameters = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseForwardedHeaders) {
			defaultVal := types.BoolValue(false)
			if !planModel.UseForwardedHeaders.Equal(defaultVal) {
				planModel.UseForwardedHeaders = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.HttpRequestHeaderSize) {
			defaultVal := types.Int64Value(8192)
			if !planModel.HttpRequestHeaderSize.Equal(defaultVal) {
				planModel.HttpRequestHeaderSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.UseCorrelationIDHeader) {
			defaultVal := types.BoolValue(true)
			if !planModel.UseCorrelationIDHeader.Equal(defaultVal) {
				planModel.UseCorrelationIDHeader = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.CorrelationIDResponseHeader) {
			defaultVal := types.StringValue("Correlation-Id")
			if !planModel.CorrelationIDResponseHeader.Equal(defaultVal) {
				planModel.CorrelationIDResponseHeader = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SslClientAuthPolicy) {
			defaultVal := types.StringValue("disabled")
			if !planModel.SslClientAuthPolicy.Equal(defaultVal) {
				planModel.SslClientAuthPolicy = defaultVal
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

func (r *defaultConnectionHandlerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanConnectionHandler(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory10000)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model connectionHandlerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.RequestHandlerPerConnection) {
		resp.Diagnostics.AddError("Attribute 'request_handler_per_connection' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
	if internaltypes.IsDefined(model.EnableSniHostnameChecks) {
		resp.Diagnostics.AddError("Attribute 'enable_sni_hostname_checks' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

func (model *connectionHandlerResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "jmx" {
		model.KeepStats = types.BoolNull()
		model.EnableMultipartMIMEParameters = types.BoolNull()
		model.AllowLDAPV2 = types.BoolNull()
		model.EnableSniHostnameChecks = types.BoolNull()
		model.NumAcceptHandlers = types.Int64Null()
		model.SendRejectionNotice = types.BoolNull()
		model.PollInterval = types.StringNull()
		model.AcceptBacklog = types.Int64Null()
		model.AllowTCPReuseAddress = types.BoolNull()
		model.LowResourcesIdleTimeLimit = types.StringNull()
		model.CorrelationIDResponseHeader = types.StringNull()
		model.MaxRequestSize = types.StringNull()
		model.MaxBlockedWriteTimeLimit = types.StringNull()
		model.CloseConnectionsWhenUnavailable = types.BoolNull()
		model.ListenAddress, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxCancelHandlers = types.Int64Null()
		model.AutoAuthenticateUsingClientCertificate = types.BoolNull()
		model.LowResourcesConnectionThreshold = types.Int64Null()
		model.SslClientAuthPolicy = types.StringNull()
		model.UseForwardedHeaders = types.BoolNull()
		model.AllowStartTLS = types.BoolNull()
		model.UseTCPKeepAlive = types.BoolNull()
		model.LdifDirectory = types.StringNull()
		model.RequestHandlerPerConnection = types.BoolNull()
		model.UseCorrelationIDHeader = types.BoolNull()
		model.IdleTimeLimit = types.StringNull()
		model.HttpRequestHeaderSize = types.Int64Null()
		model.FailedBindResponseDelay = types.StringNull()
		model.CloseConnectionsOnExplicitGC = types.BoolNull()
		model.NumRequestHandlers = types.Int64Null()
	}
	if resourceType == "ldap" {
		model.KeepStats = types.BoolNull()
		model.EnableMultipartMIMEParameters = types.BoolNull()
		model.EnableSniHostnameChecks = types.BoolNull()
		model.PollInterval = types.StringNull()
		model.AllowTCPReuseAddress = types.BoolNull()
		model.LowResourcesIdleTimeLimit = types.StringNull()
		model.CorrelationIDResponseHeader = types.StringNull()
		model.LowResourcesConnectionThreshold = types.Int64Null()
		model.UseForwardedHeaders = types.BoolNull()
		model.LdifDirectory = types.StringNull()
		model.UseCorrelationIDHeader = types.BoolNull()
		model.IdleTimeLimit = types.StringNull()
		model.HttpRequestHeaderSize = types.Int64Null()
	}
	if resourceType == "ldif" {
		model.KeepStats = types.BoolNull()
		model.EnableMultipartMIMEParameters = types.BoolNull()
		model.AllowLDAPV2 = types.BoolNull()
		model.EnableSniHostnameChecks = types.BoolNull()
		model.NumAcceptHandlers = types.Int64Null()
		model.SendRejectionNotice = types.BoolNull()
		model.AcceptBacklog = types.Int64Null()
		model.AllowTCPReuseAddress = types.BoolNull()
		model.LowResourcesIdleTimeLimit = types.StringNull()
		model.CorrelationIDResponseHeader = types.StringNull()
		model.UseSSL = types.BoolNull()
		model.MaxRequestSize = types.StringNull()
		model.MaxBlockedWriteTimeLimit = types.StringNull()
		model.CloseConnectionsWhenUnavailable = types.BoolNull()
		model.ListenAddress, _ = types.SetValue(types.StringType, []attr.Value{})
		model.MaxCancelHandlers = types.Int64Null()
		model.AutoAuthenticateUsingClientCertificate = types.BoolNull()
		model.LowResourcesConnectionThreshold = types.Int64Null()
		model.SslClientAuthPolicy = types.StringNull()
		model.UseForwardedHeaders = types.BoolNull()
		model.AllowStartTLS = types.BoolNull()
		model.UseTCPKeepAlive = types.BoolNull()
		model.RequestHandlerPerConnection = types.BoolNull()
		model.UseCorrelationIDHeader = types.BoolNull()
		model.IdleTimeLimit = types.StringNull()
		model.HttpRequestHeaderSize = types.Int64Null()
		model.FailedBindResponseDelay = types.StringNull()
		model.CloseConnectionsOnExplicitGC = types.BoolNull()
		model.NumRequestHandlers = types.Int64Null()
	}
	if resourceType == "http" {
		model.AllowLDAPV2 = types.BoolNull()
		model.NumAcceptHandlers = types.Int64Null()
		model.SendRejectionNotice = types.BoolNull()
		model.PollInterval = types.StringNull()
		model.MaxRequestSize = types.StringNull()
		model.MaxBlockedWriteTimeLimit = types.StringNull()
		model.CloseConnectionsWhenUnavailable = types.BoolNull()
		model.MaxCancelHandlers = types.Int64Null()
		model.AutoAuthenticateUsingClientCertificate = types.BoolNull()
		model.AllowStartTLS = types.BoolNull()
		model.UseTCPKeepAlive = types.BoolNull()
		model.LdifDirectory = types.StringNull()
		model.RequestHandlerPerConnection = types.BoolNull()
		model.FailedBindResponseDelay = types.StringNull()
		model.CloseConnectionsOnExplicitGC = types.BoolNull()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsConnectionHandler() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"http"},
			resourcevalidator.AtLeastOneOf(
				path.MatchRoot("http_servlet_extension"),
				path.MatchRoot("web_application_extension"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("listen_port"),
			path.MatchRoot("type"),
			[]string{"jmx", "ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_ssl"),
			path.MatchRoot("type"),
			[]string{"jmx", "ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ssl_cert_nickname"),
			path.MatchRoot("type"),
			[]string{"jmx", "ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_manager_provider"),
			path.MatchRoot("type"),
			[]string{"jmx", "ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_client"),
			path.MatchRoot("type"),
			[]string{"jmx", "ldap", "ldif"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("denied_client"),
			path.MatchRoot("type"),
			[]string{"jmx", "ldap", "ldif"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("listen_address"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_start_tls"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_manager_provider"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_ldap_v2"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_tcp_keep_alive"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("send_rejection_notice"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("failed_bind_response_delay"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_request_size"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_cancel_handlers"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("num_accept_handlers"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("num_request_handlers"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("request_handler_per_connection"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ssl_client_auth_policy"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("accept_backlog"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ssl_protocol"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ssl_cipher_suite"),
			path.MatchRoot("type"),
			[]string{"ldap", "http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_blocked_write_time_limit"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("auto_authenticate_using_client_certificate"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("close_connections_when_unavailable"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("close_connections_on_explicit_gc"),
			path.MatchRoot("type"),
			[]string{"ldap"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ldif_directory"),
			path.MatchRoot("type"),
			[]string{"ldif"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("poll_interval"),
			path.MatchRoot("type"),
			[]string{"ldif"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_servlet_extension"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("web_application_extension"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_operation_log_publisher"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("keep_stats"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allow_tcp_reuse_address"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("idle_time_limit"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("low_resources_connection_threshold"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("low_resources_idle_time_limit"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("enable_multipart_mime_parameters"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_forwarded_headers"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_request_header_size"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("response_header"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("use_correlation_id_header"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("correlation_id_response_header"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("correlation_id_request_header"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("enable_sni_hostname_checks"),
			path.MatchRoot("type"),
			[]string{"http"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"jmx",
			[]path.Expression{path.MatchRoot("listen_port")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"ldap",
			[]path.Expression{path.MatchRoot("listen_port")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"http",
			[]path.Expression{path.MatchRoot("listen_port")},
		),
	}
}

// Add config validators
func (r connectionHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsConnectionHandler()
}

// Add config validators
func (r defaultConnectionHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsConnectionHandler()
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
	if internaltypes.IsDefined(plan.RequestHandlerPerConnection) {
		addRequest.RequestHandlerPerConnection = plan.RequestHandlerPerConnection.ValueBoolPointer()
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
	if internaltypes.IsDefined(plan.EnableSniHostnameChecks) {
		addRequest.EnableSniHostnameChecks = plan.EnableSniHostnameChecks.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateConnectionHandlerUnknownValues(model *connectionHandlerResourceModel) {
	if model.SslCipherSuite.IsUnknown() || model.SslCipherSuite.IsNull() {
		model.SslCipherSuite, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.SslProtocol.IsUnknown() || model.SslProtocol.IsNull() {
		model.SslProtocol, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ResponseHeader.IsUnknown() || model.ResponseHeader.IsNull() {
		model.ResponseHeader, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllowedClient.IsUnknown() || model.AllowedClient.IsNull() {
		model.AllowedClient, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.WebApplicationExtension.IsUnknown() || model.WebApplicationExtension.IsNull() {
		model.WebApplicationExtension, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.HttpServletExtension.IsUnknown() || model.HttpServletExtension.IsNull() {
		model.HttpServletExtension, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.ListenAddress.IsUnknown() || model.ListenAddress.IsNull() {
		model.ListenAddress, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.CorrelationIDRequestHeader.IsUnknown() || model.CorrelationIDRequestHeader.IsNull() {
		model.CorrelationIDRequestHeader, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.HttpOperationLogPublisher.IsUnknown() || model.HttpOperationLogPublisher.IsNull() {
		model.HttpOperationLogPublisher, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.DeniedClient.IsUnknown() || model.DeniedClient.IsNull() {
		model.DeniedClient, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *connectionHandlerResourceModel) populateAllComputedStringAttributes() {
	if model.TrustManagerProvider.IsUnknown() || model.TrustManagerProvider.IsNull() {
		model.TrustManagerProvider = types.StringValue("")
	}
	if model.LowResourcesIdleTimeLimit.IsUnknown() || model.LowResourcesIdleTimeLimit.IsNull() {
		model.LowResourcesIdleTimeLimit = types.StringValue("")
	}
	if model.FailedBindResponseDelay.IsUnknown() || model.FailedBindResponseDelay.IsNull() {
		model.FailedBindResponseDelay = types.StringValue("")
	}
	if model.MaxBlockedWriteTimeLimit.IsUnknown() || model.MaxBlockedWriteTimeLimit.IsNull() {
		model.MaxBlockedWriteTimeLimit = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.PollInterval.IsUnknown() || model.PollInterval.IsNull() {
		model.PollInterval = types.StringValue("")
	}
	if model.SslClientAuthPolicy.IsUnknown() || model.SslClientAuthPolicy.IsNull() {
		model.SslClientAuthPolicy = types.StringValue("")
	}
	if model.CorrelationIDResponseHeader.IsUnknown() || model.CorrelationIDResponseHeader.IsNull() {
		model.CorrelationIDResponseHeader = types.StringValue("")
	}
	if model.SslCertNickname.IsUnknown() || model.SslCertNickname.IsNull() {
		model.SslCertNickname = types.StringValue("")
	}
	if model.KeyManagerProvider.IsUnknown() || model.KeyManagerProvider.IsNull() {
		model.KeyManagerProvider = types.StringValue("")
	}
	if model.LdifDirectory.IsUnknown() || model.LdifDirectory.IsNull() {
		model.LdifDirectory = types.StringValue("")
	}
	if model.MaxRequestSize.IsUnknown() || model.MaxRequestSize.IsNull() {
		model.MaxRequestSize = types.StringValue("")
	}
	if model.IdleTimeLimit.IsUnknown() || model.IdleTimeLimit.IsNull() {
		model.IdleTimeLimit = types.StringValue("")
	}
}

// Read a JmxConnectionHandlerResponse object into the model struct
func readJmxConnectionHandlerResponse(ctx context.Context, r *client.JmxConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jmx")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ListenPort = types.Int64Value(r.ListenPort)
	state.UseSSL = internaltypes.BoolTypeOrNil(r.UseSSL)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, internaltypes.IsEmptyString(expectedValues.SslCertNickname))
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionHandlerUnknownValues(state)
}

// Read a LdapConnectionHandlerResponse object into the model struct
func readLdapConnectionHandlerResponse(ctx context.Context, r *client.LdapConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	state.FailedBindResponseDelay = internaltypes.StringTypeOrNil(r.FailedBindResponseDelay, true)
	config.CheckMismatchedPDFormattedAttributes("failed_bind_response_delay",
		expectedValues.FailedBindResponseDelay, state.FailedBindResponseDelay, diagnostics)
	state.MaxRequestSize = internaltypes.StringTypeOrNil(r.MaxRequestSize, true)
	config.CheckMismatchedPDFormattedAttributes("max_request_size",
		expectedValues.MaxRequestSize, state.MaxRequestSize, diagnostics)
	state.MaxCancelHandlers = internaltypes.Int64TypeOrNil(r.MaxCancelHandlers)
	state.NumAcceptHandlers = internaltypes.Int64TypeOrNil(r.NumAcceptHandlers)
	state.NumRequestHandlers = internaltypes.Int64TypeOrNil(r.NumRequestHandlers)
	state.RequestHandlerPerConnection = internaltypes.BoolTypeOrNil(r.RequestHandlerPerConnection)
	state.SslClientAuthPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionHandlerSslClientAuthPolicyProp(r.SslClientAuthPolicy), true)
	state.AcceptBacklog = internaltypes.Int64TypeOrNil(r.AcceptBacklog)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.MaxBlockedWriteTimeLimit = internaltypes.StringTypeOrNil(r.MaxBlockedWriteTimeLimit, true)
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
	populateConnectionHandlerUnknownValues(state)
}

// Read a LdifConnectionHandlerResponse object into the model struct
func readLdifConnectionHandlerResponse(ctx context.Context, r *client.LdifConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.PollInterval = types.StringValue(r.PollInterval)
	config.CheckMismatchedPDFormattedAttributes("poll_interval",
		expectedValues.PollInterval, state.PollInterval, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionHandlerUnknownValues(state)
}

// Read a HttpConnectionHandlerResponse object into the model struct
func readHttpConnectionHandlerResponse(ctx context.Context, r *client.HttpConnectionHandlerResponse, state *connectionHandlerResourceModel, expectedValues *connectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("http")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
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
	state.IdleTimeLimit = internaltypes.StringTypeOrNil(r.IdleTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("idle_time_limit",
		expectedValues.IdleTimeLimit, state.IdleTimeLimit, diagnostics)
	state.LowResourcesConnectionThreshold = internaltypes.Int64TypeOrNil(r.LowResourcesConnectionThreshold)
	state.LowResourcesIdleTimeLimit = internaltypes.StringTypeOrNil(r.LowResourcesIdleTimeLimit, true)
	config.CheckMismatchedPDFormattedAttributes("low_resources_idle_time_limit",
		expectedValues.LowResourcesIdleTimeLimit, state.LowResourcesIdleTimeLimit, diagnostics)
	state.EnableMultipartMIMEParameters = internaltypes.BoolTypeOrNil(r.EnableMultipartMIMEParameters)
	state.UseForwardedHeaders = internaltypes.BoolTypeOrNil(r.UseForwardedHeaders)
	state.HttpRequestHeaderSize = internaltypes.Int64TypeOrNil(r.HttpRequestHeaderSize)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.UseCorrelationIDHeader = internaltypes.BoolTypeOrNil(r.UseCorrelationIDHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.CorrelationIDRequestHeader = internaltypes.GetStringSet(r.CorrelationIDRequestHeader)
	state.SslClientAuthPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionHandlerSslClientAuthPolicyProp(r.SslClientAuthPolicy), true)
	state.EnableSniHostnameChecks = internaltypes.BoolTypeOrNil(r.EnableSniHostnameChecks)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConnectionHandlerUnknownValues(state)
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
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableSniHostnameChecks, state.EnableSniHostnameChecks, "enable-sni-hostname-checks")
	operations.AddBoolOperationIfNecessary(&ops, plan.SendRejectionNotice, state.SendRejectionNotice, "send-rejection-notice")
	operations.AddStringOperationIfNecessary(&ops, plan.FailedBindResponseDelay, state.FailedBindResponseDelay, "failed-bind-response-delay")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxRequestSize, state.MaxRequestSize, "max-request-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxCancelHandlers, state.MaxCancelHandlers, "max-cancel-handlers")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumAcceptHandlers, state.NumAcceptHandlers, "num-accept-handlers")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumRequestHandlers, state.NumRequestHandlers, "num-request-handlers")
	operations.AddBoolOperationIfNecessary(&ops, plan.RequestHandlerPerConnection, state.RequestHandlerPerConnection, "request-handler-per-connection")
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
	addRequest := client.NewAddJmxConnectionHandlerRequest([]client.EnumjmxConnectionHandlerSchemaUrn{client.ENUMJMXCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERJMX},
		plan.ListenPort.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
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
	apiAddRequest := r.apiClient.ConnectionHandlerAPI.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddJmxConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.AddConnectionHandlerExecute(apiAddRequest)
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
	addRequest := client.NewAddLdapConnectionHandlerRequest([]client.EnumldapConnectionHandlerSchemaUrn{client.ENUMLDAPCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERLDAP},
		plan.ListenPort.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
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
	apiAddRequest := r.apiClient.ConnectionHandlerAPI.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddLdapConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.AddConnectionHandlerExecute(apiAddRequest)
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
	addRequest := client.NewAddLdifConnectionHandlerRequest([]client.EnumldifConnectionHandlerSchemaUrn{client.ENUMLDIFCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERLDIF},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
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
	apiAddRequest := r.apiClient.ConnectionHandlerAPI.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddLdifConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.AddConnectionHandlerExecute(apiAddRequest)
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
	addRequest := client.NewAddHttpConnectionHandlerRequest([]client.EnumhttpConnectionHandlerSchemaUrn{client.ENUMHTTPCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERHTTP},
		plan.ListenPort.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
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
	apiAddRequest := r.apiClient.ConnectionHandlerAPI.AddConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddConnectionHandlerRequest(
		client.AddHttpConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.AddConnectionHandlerExecute(apiAddRequest)
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

	readResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
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

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionHandlerAPI.UpdateConnectionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.UpdateConnectionHandlerExecute(updateRequest)
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
		if updateResponse.JmxConnectionHandlerResponse != nil {
			readJmxConnectionHandlerResponse(ctx, updateResponse.JmxConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapConnectionHandlerResponse != nil {
			readLdapConnectionHandlerResponse(ctx, updateResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdifConnectionHandlerResponse != nil {
			readLdifConnectionHandlerResponse(ctx, updateResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HttpConnectionHandlerResponse != nil {
			readHttpConnectionHandlerResponse(ctx, updateResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *connectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readConnectionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state connectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionHandlerAPI.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Connection Handler", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Connection Handler", err, httpResp)
		}
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

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
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
	updateRequest := apiClient.ConnectionHandlerAPI.UpdateConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionHandlerAPI.UpdateConnectionHandlerExecute(updateRequest)
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
		if updateResponse.JmxConnectionHandlerResponse != nil {
			readJmxConnectionHandlerResponse(ctx, updateResponse.JmxConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdapConnectionHandlerResponse != nil {
			readLdapConnectionHandlerResponse(ctx, updateResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.LdifConnectionHandlerResponse != nil {
			readLdifConnectionHandlerResponse(ctx, updateResponse.LdifConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.HttpConnectionHandlerResponse != nil {
			readHttpConnectionHandlerResponse(ctx, updateResponse.HttpConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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

	httpResp, err := r.apiClient.ConnectionHandlerAPI.DeleteConnectionHandlerExecute(r.apiClient.ConnectionHandlerAPI.DeleteConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
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
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
