package connectionhandler

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &connectionHandlerDataSource{}
	_ datasource.DataSourceWithConfigure = &connectionHandlerDataSource{}
)

// Create a Connection Handler data source
func NewConnectionHandlerDataSource() datasource.DataSource {
	return &connectionHandlerDataSource{}
}

// connectionHandlerDataSource is the datasource implementation.
type connectionHandlerDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *connectionHandlerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_connection_handler"
}

// Configure adds the provider configured client to the data source.
func (r *connectionHandlerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type connectionHandlerDataSourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	Name                                   types.String `tfsdk:"name"`
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

// GetSchema defines the schema for the datasource.
func (r *connectionHandlerDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Connection Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Connection Handler resource. Options are ['jmx', 'ldap', 'ldif', 'http']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"listen_address": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the address or set of addresses on which this LDAP Connection Handler should listen for connections from LDAP clients. When the `type` attribute is set to `http`: Specifies the address on which to listen for connections from HTTP clients. If no value is defined, the server will listen on all addresses on all interfaces.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the address or set of addresses on which this LDAP Connection Handler should listen for connections from LDAP clients.\n  - `http`: Specifies the address on which to listen for connections from HTTP clients. If no value is defined, the server will listen on all addresses on all interfaces.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"listen_port": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `jmx`: Specifies the port number on which the JMX Connection Handler will listen for connections from clients. When the `type` attribute is set to `ldap`: Specifies the port number on which the LDAP Connection Handler will listen for connections from clients. When the `type` attribute is set to `http`: Specifies the port number on which the HTTP Connection Handler will listen for connections from clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Specifies the port number on which the JMX Connection Handler will listen for connections from clients.\n  - `ldap`: Specifies the port number on which the LDAP Connection Handler will listen for connections from clients.\n  - `http`: Specifies the port number on which the HTTP Connection Handler will listen for connections from clients.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"ldif_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory in which the LDIF files should be placed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"poll_interval": schema.StringAttribute{
				Description: "Specifies how frequently the LDIF connection handler should check the LDIF directory to determine whether a new LDIF file has been added.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_servlet_extension": schema.SetAttribute{
				Description: "Specifies information about servlets that will be provided via this connection handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"web_application_extension": schema.SetAttribute{
				Description: "Specifies information about web applications that will be provided via this connection handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"http_operation_log_publisher": schema.SetAttribute{
				Description: "Specifies the set of HTTP operation loggers that should be used to log information about requests and responses for operations processed through this HTTP Connection Handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"use_ssl": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to `jmx`: Indicates whether the JMX Connection Handler should use SSL. When the `type` attribute is set to `ldap`: Indicates whether the LDAP Connection Handler should use SSL. When the `type` attribute is set to `http`: Indicates whether the HTTP Connection Handler should use SSL.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Indicates whether the JMX Connection Handler should use SSL.\n  - `ldap`: Indicates whether the LDAP Connection Handler should use SSL.\n  - `http`: Indicates whether the HTTP Connection Handler should use SSL.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"allow_start_tls": schema.BoolAttribute{
				Description: "Indicates whether clients are allowed to use StartTLS.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `jmx`: Specifies the nickname (also called the alias) of the certificate that the JMX Connection Handler should use when performing SSL communication. When the `type` attribute is set to `ldap`: Specifies the nickname (also called the alias) of the certificate that the LDAP Connection Handler should use when performing SSL communication. When the `type` attribute is set to `http`: Specifies the nickname (also called the alias) of the certificate that the HTTP Connection Handler should use when performing SSL communication.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Specifies the nickname (also called the alias) of the certificate that the JMX Connection Handler should use when performing SSL communication.\n  - `ldap`: Specifies the nickname (also called the alias) of the certificate that the LDAP Connection Handler should use when performing SSL communication.\n  - `http`: Specifies the nickname (also called the alias) of the certificate that the HTTP Connection Handler should use when performing SSL communication.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `jmx`: Specifies the name of the key manager that should be used with this JMX Connection Handler . When the `type` attribute is set to `ldap`: Specifies the name of the key manager that should be used with this LDAP Connection Handler . When the `type` attribute is set to `http`: Specifies the key manager provider that will be used to obtain the certificate to present to HTTPS clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `jmx`: Specifies the name of the key manager that should be used with this JMX Connection Handler .\n  - `ldap`: Specifies the name of the key manager that should be used with this LDAP Connection Handler .\n  - `http`: Specifies the key manager provider that will be used to obtain the certificate to present to HTTPS clients.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the name of the trust manager that should be used with the LDAP Connection Handler . When the `type` attribute is set to `http`: Specifies the trust manager provider that will be used to validate any certificates presented by HTTPS clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the name of the trust manager that should be used with the LDAP Connection Handler .\n  - `http`: Specifies the trust manager provider that will be used to validate any certificates presented by HTTPS clients.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"keep_stats": schema.BoolAttribute{
				Description: "Indicates whether to enable statistics collection for this connection handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_ldap_v2": schema.BoolAttribute{
				Description: "Indicates whether connections from LDAPv2 clients are allowed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_tcp_reuse_address": schema.BoolAttribute{
				Description: "Indicates whether the server should attempt to reuse socket descriptors. This may be useful in environments with a high rate of connection establishment and termination.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"idle_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum idle time for a connection. The max idle time is applied when waiting for a new request to be received on a connection, when reading the headers and content of a request, or when writing the headers and content of a response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"low_resources_connection_threshold": schema.Int64Attribute{
				Description: "Specifies the number of connections, which if exceeded, places this handler in a low resource state where a different idle time limit is applied on the connections.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"low_resources_idle_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum idle time for a connection when this handler is in a low resource state as defined by low-resource-connections. The max idle time is applied when waiting for a new request to be received on a connection, when reading the headers and content of a request, or when writing the headers and content of a response.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_multipart_mime_parameters": schema.BoolAttribute{
				Description: "Determines whether request form parameters submitted in multipart/ form-data (RFC 2388) format should be processed as request parameters.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_forwarded_headers": schema.BoolAttribute{
				Description: "Indicates whether to use \"Forwarded\" and \"X-Forwarded-*\" request headers to override corresponding HTTP request information available during request processing.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_request_header_size": schema.Int64Attribute{
				Description: "Specifies the maximum buffer size of an http request including the request uri and all of the request headers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"use_correlation_id_header": schema.BoolAttribute{
				Description: "If enabled, a correlation ID header will be added to outgoing HTTP responses.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"correlation_id_request_header": schema.SetAttribute{
				Description: "Specifies the set of HTTP request headers that may contain a value to be used as the correlation ID. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"use_tcp_keep_alive": schema.BoolAttribute{
				Description: "Indicates whether the LDAP Connection Handler should use TCP keep-alive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enable_sni_hostname_checks": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Requires SNI hostnames to match or else throw an Invalid SNI error.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"send_rejection_notice": schema.BoolAttribute{
				Description: "Indicates whether the LDAP Connection Handler should send a notice of disconnection extended response message to the client if a new connection is rejected for some reason.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"failed_bind_response_delay": schema.StringAttribute{
				Description: "Specifies the length of time that the server should delay the response to non-successful bind operations. A value of zero milliseconds indicates that non-successful bind operations should not be delayed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_request_size": schema.StringAttribute{
				Description: "Specifies the size of the largest LDAP request message that will be allowed by this LDAP Connection handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_cancel_handlers": schema.Int64Attribute{
				Description: "Specifies the maximum number of threads that are used to process cancel and abandon requests from clients.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_accept_handlers": schema.Int64Attribute{
				Description: "Specifies the number of threads that are used to accept new client connections, and to perform any initial preparation on those connections that may be needed before the connection can be used to read requests and send responses.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_request_handlers": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the number of request handlers that are used to read requests from clients. When the `type` attribute is set to `http`: Specifies the number of threads that will be used for accepting connections and reading requests from clients.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the number of request handlers that are used to read requests from clients.\n  - `http`: Specifies the number of threads that will be used for accepting connections and reading requests from clients.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"request_handler_per_connection": schema.BoolAttribute{
				Description: "Supported in PingDirectory product version 10.0.0.0+. Indicates whether a separate request handler thread should be created for each client connection, which can help avoid starvation of client connections for cases in which one or more clients send large numbers of concurrent asynchronous requests. This should only be used for cases in which a relatively small number of connections will be established at any given time, the connections established will generally be long-lived, and at least one client may send high volumes of asynchronous requests. This property can be used to alleviate possible blocking during long-running TLS negotiation on a single request handler which can result in it being unable to acknowledge further client requests until the TLS negotation completes or times out.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ssl_client_auth_policy": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the policy that the LDAP Connection Handler should use regarding client SSL certificates. When the `type` attribute is set to `http`: Specifies the policy that the HTTP Connection Handler should use regarding client SSL certificates. In order for a client certificate to be accepted it must be known to the trust-manager-provider associated with this HTTP Connection Handler. Client certificates received by the HTTP Connection Handler are by default used for TLS mutual authentication only, as there is no support for user authentication.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the policy that the LDAP Connection Handler should use regarding client SSL certificates.\n  - `http`: Specifies the policy that the HTTP Connection Handler should use regarding client SSL certificates. In order for a client certificate to be accepted it must be known to the trust-manager-provider associated with this HTTP Connection Handler. Client certificates received by the HTTP Connection Handler are by default used for TLS mutual authentication only, as there is no support for user authentication.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"accept_backlog": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the maximum number of pending connection attempts that are allowed to queue up in the accept backlog before the server starts rejecting new connection attempts. When the `type` attribute is set to `http`: Specifies the number of concurrent outstanding connection attempts that the connection handler should allow. The default value should be acceptable in most cases, but it may need to be increased in environments that may attempt to establish large numbers of connections simultaneously.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the maximum number of pending connection attempts that are allowed to queue up in the accept backlog before the server starts rejecting new connection attempts.\n  - `http`: Specifies the number of concurrent outstanding connection attempts that the connection handler should allow. The default value should be acceptable in most cases, but it may need to be increased in environments that may attempt to establish large numbers of connections simultaneously.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"ssl_protocol": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the names of the TLS protocols that are allowed for use in SSL or StartTLS communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry. When the `type` attribute is set to `http`: Specifies the names of the SSL protocols that are allowed for use in SSL communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the names of the TLS protocols that are allowed for use in SSL or StartTLS communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.\n  - `http`: Specifies the names of the SSL protocols that are allowed for use in SSL communication. The set of supported ssl protocols can be viewed via the ssl context monitor entry.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `ldap`: Specifies the names of the TLS cipher suites that are allowed for use in SSL or StartTLS communication. The set of supported cipher suites can be viewed via the ssl context monitor entry. When the `type` attribute is set to `http`: Specifies the names of the SSL cipher suites that are allowed for use in SSL communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ldap`: Specifies the names of the TLS cipher suites that are allowed for use in SSL or StartTLS communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.\n  - `http`: Specifies the names of the SSL cipher suites that are allowed for use in SSL communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"max_blocked_write_time_limit": schema.StringAttribute{
				Description: "Specifies the maximum length of time that attempts to write data to LDAP clients should be allowed to block.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auto_authenticate_using_client_certificate": schema.BoolAttribute{
				Description: "Indicates whether to attempt to automatically authenticate a client connection that has established a secure communication channel (using either SSL or StartTLS) and presented its own client certificate. Generally, clients should use the SASL EXTERNAL mechanism to authenticate using a client certificate, but some clients may not support that capability and/or may expect automatic authentication.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"close_connections_when_unavailable": schema.BoolAttribute{
				Description: "Indicates whether all connections associated with this LDAP Connection Handler should be closed and no new connections accepted when the server has determined that it is \"unavailable.\" This allows clients (or a network load balancer) to route requests to another server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"close_connections_on_explicit_gc": schema.BoolAttribute{
				Description: "Indicates whether all connections associated with this LDAP Connection Handler should be closed before an explicit garbage collection is performed to allow clients to route requests to another server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Connection Handler",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Connection Handler is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are allowed to establish connections to this connection handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"denied_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are not allowed to establish connections to this connection handler.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a JmxConnectionHandlerResponse object into the model struct
func readJmxConnectionHandlerResponseDataSource(ctx context.Context, r *client.JmxConnectionHandlerResponse, state *connectionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jmx")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ListenPort = types.Int64Value(r.ListenPort)
	state.UseSSL = internaltypes.BoolTypeOrNil(r.UseSSL)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, false)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
}

// Read a LdapConnectionHandlerResponse object into the model struct
func readLdapConnectionHandlerResponseDataSource(ctx context.Context, r *client.LdapConnectionHandlerResponse, state *connectionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ListenAddress = internaltypes.GetStringSet(r.ListenAddress)
	state.ListenPort = types.Int64Value(r.ListenPort)
	state.UseSSL = internaltypes.BoolTypeOrNil(r.UseSSL)
	state.AllowStartTLS = internaltypes.BoolTypeOrNil(r.AllowStartTLS)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, false)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.AllowLDAPV2 = internaltypes.BoolTypeOrNil(r.AllowLDAPV2)
	state.UseTCPKeepAlive = internaltypes.BoolTypeOrNil(r.UseTCPKeepAlive)
	state.SendRejectionNotice = internaltypes.BoolTypeOrNil(r.SendRejectionNotice)
	state.FailedBindResponseDelay = internaltypes.StringTypeOrNil(r.FailedBindResponseDelay, false)
	state.MaxRequestSize = internaltypes.StringTypeOrNil(r.MaxRequestSize, false)
	state.MaxCancelHandlers = internaltypes.Int64TypeOrNil(r.MaxCancelHandlers)
	state.NumAcceptHandlers = internaltypes.Int64TypeOrNil(r.NumAcceptHandlers)
	state.NumRequestHandlers = internaltypes.Int64TypeOrNil(r.NumRequestHandlers)
	state.RequestHandlerPerConnection = internaltypes.BoolTypeOrNil(r.RequestHandlerPerConnection)
	state.SslClientAuthPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionHandlerSslClientAuthPolicyProp(r.SslClientAuthPolicy), false)
	state.AcceptBacklog = internaltypes.Int64TypeOrNil(r.AcceptBacklog)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.MaxBlockedWriteTimeLimit = internaltypes.StringTypeOrNil(r.MaxBlockedWriteTimeLimit, false)
	state.AutoAuthenticateUsingClientCertificate = internaltypes.BoolTypeOrNil(r.AutoAuthenticateUsingClientCertificate)
	state.CloseConnectionsWhenUnavailable = internaltypes.BoolTypeOrNil(r.CloseConnectionsWhenUnavailable)
	state.CloseConnectionsOnExplicitGC = internaltypes.BoolTypeOrNil(r.CloseConnectionsOnExplicitGC)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
}

// Read a LdifConnectionHandlerResponse object into the model struct
func readLdifConnectionHandlerResponseDataSource(ctx context.Context, r *client.LdifConnectionHandlerResponse, state *connectionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldif")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowedClient = internaltypes.GetStringSet(r.AllowedClient)
	state.DeniedClient = internaltypes.GetStringSet(r.DeniedClient)
	state.LdifDirectory = types.StringValue(r.LdifDirectory)
	state.PollInterval = types.StringValue(r.PollInterval)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a HttpConnectionHandlerResponse object into the model struct
func readHttpConnectionHandlerResponseDataSource(ctx context.Context, r *client.HttpConnectionHandlerResponse, state *connectionHandlerDataSourceModel, diagnostics *diag.Diagnostics) {
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
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, false)
	state.HttpServletExtension = internaltypes.GetStringSet(r.HttpServletExtension)
	state.WebApplicationExtension = internaltypes.GetStringSet(r.WebApplicationExtension)
	state.HttpOperationLogPublisher = internaltypes.GetStringSet(r.HttpOperationLogPublisher)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, false)
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, false)
	state.NumRequestHandlers = internaltypes.Int64TypeOrNil(r.NumRequestHandlers)
	state.KeepStats = internaltypes.BoolTypeOrNil(r.KeepStats)
	state.AcceptBacklog = internaltypes.Int64TypeOrNil(r.AcceptBacklog)
	state.AllowTCPReuseAddress = internaltypes.BoolTypeOrNil(r.AllowTCPReuseAddress)
	state.IdleTimeLimit = internaltypes.StringTypeOrNil(r.IdleTimeLimit, false)
	state.LowResourcesConnectionThreshold = internaltypes.Int64TypeOrNil(r.LowResourcesConnectionThreshold)
	state.LowResourcesIdleTimeLimit = internaltypes.StringTypeOrNil(r.LowResourcesIdleTimeLimit, false)
	state.EnableMultipartMIMEParameters = internaltypes.BoolTypeOrNil(r.EnableMultipartMIMEParameters)
	state.UseForwardedHeaders = internaltypes.BoolTypeOrNil(r.UseForwardedHeaders)
	state.HttpRequestHeaderSize = internaltypes.Int64TypeOrNil(r.HttpRequestHeaderSize)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.UseCorrelationIDHeader = internaltypes.BoolTypeOrNil(r.UseCorrelationIDHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, false)
	state.CorrelationIDRequestHeader = internaltypes.GetStringSet(r.CorrelationIDRequestHeader)
	state.SslClientAuthPolicy = internaltypes.StringTypeOrNil(
		client.StringPointerEnumconnectionHandlerSslClientAuthPolicyProp(r.SslClientAuthPolicy), false)
	state.EnableSniHostnameChecks = internaltypes.BoolTypeOrNil(r.EnableSniHostnameChecks)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *connectionHandlerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state connectionHandlerDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionHandlerAPI.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
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
		readJmxConnectionHandlerResponseDataSource(ctx, readResponse.JmxConnectionHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapConnectionHandlerResponse != nil {
		readLdapConnectionHandlerResponseDataSource(ctx, readResponse.LdapConnectionHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdifConnectionHandlerResponse != nil {
		readLdifConnectionHandlerResponseDataSource(ctx, readResponse.LdifConnectionHandlerResponse, &state, &resp.Diagnostics)
	}
	if readResponse.HttpConnectionHandlerResponse != nil {
		readHttpConnectionHandlerResponseDataSource(ctx, readResponse.HttpConnectionHandlerResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
