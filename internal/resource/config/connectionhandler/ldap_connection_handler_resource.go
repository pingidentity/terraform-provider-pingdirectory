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
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &ldapConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &ldapConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &ldapConnectionHandlerResource{}
	_ resource.Resource                = &defaultLdapConnectionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultLdapConnectionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultLdapConnectionHandlerResource{}
)

// Create a Ldap Connection Handler resource
func NewLdapConnectionHandlerResource() resource.Resource {
	return &ldapConnectionHandlerResource{}
}

func NewDefaultLdapConnectionHandlerResource() resource.Resource {
	return &defaultLdapConnectionHandlerResource{}
}

// ldapConnectionHandlerResource is the resource implementation.
type ldapConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultLdapConnectionHandlerResource is the resource implementation.
type defaultLdapConnectionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *ldapConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_connection_handler"
}

func (r *defaultLdapConnectionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_ldap_connection_handler"
}

// Configure adds the provider configured client to the resource.
func (r *ldapConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultLdapConnectionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type ldapConnectionHandlerResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	LastUpdated                            types.String `tfsdk:"last_updated"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	ListenAddress                          types.Set    `tfsdk:"listen_address"`
	ListenPort                             types.Int64  `tfsdk:"listen_port"`
	UseSSL                                 types.Bool   `tfsdk:"use_ssl"`
	AllowStartTLS                          types.Bool   `tfsdk:"allow_start_tls"`
	SslCertNickname                        types.String `tfsdk:"ssl_cert_nickname"`
	KeyManagerProvider                     types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider                   types.String `tfsdk:"trust_manager_provider"`
	AllowLDAPV2                            types.Bool   `tfsdk:"allow_ldap_v2"`
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
func (r *ldapConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapConnectionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultLdapConnectionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	ldapConnectionHandlerSchema(ctx, req, resp, true)
}

func ldapConnectionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Ldap Connection Handler.",
		Attributes: map[string]schema.Attribute{
			"listen_address": schema.SetAttribute{
				Description: "Specifies the address or set of addresses on which this LDAP Connection Handler should listen for connections from LDAP clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"listen_port": schema.Int64Attribute{
				Description: "Specifies the port number on which the LDAP Connection Handler will listen for connections from clients.",
				Required:    true,
			},
			"use_ssl": schema.BoolAttribute{
				Description: "Indicates whether the LDAP Connection Handler should use SSL.",
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
				Description: "Specifies the nickname (also called the alias) of the certificate that the LDAP Connection Handler should use when performing SSL communication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "Specifies the name of the key manager that should be used with this LDAP Connection Handler .",
				Optional:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "Specifies the name of the trust manager that should be used with the LDAP Connection Handler .",
				Optional:    true,
			},
			"allow_ldap_v2": schema.BoolAttribute{
				Description: "Indicates whether connections from LDAPv2 clients are allowed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the SSL cipher suites that are allowed for use in SSL or StartTLS communication. The set of supported cipher suites can be viewed via the ssl context monitor entry.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"denied_client": schema.SetAttribute{
				Description: "Specifies a set of address masks that determines the addresses of the clients that are not allowed to establish connections to this connection handler.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
func addOptionalLdapConnectionHandlerFields(ctx context.Context, addRequest *client.AddLdapConnectionHandlerRequest, plan ldapConnectionHandlerResourceModel) error {
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
		intVal := int32(plan.MaxCancelHandlers.ValueInt64())
		addRequest.MaxCancelHandlers = &intVal
	}
	if internaltypes.IsDefined(plan.NumAcceptHandlers) {
		intVal := int32(plan.NumAcceptHandlers.ValueInt64())
		addRequest.NumAcceptHandlers = &intVal
	}
	if internaltypes.IsDefined(plan.NumRequestHandlers) {
		intVal := int32(plan.NumRequestHandlers.ValueInt64())
		addRequest.NumRequestHandlers = &intVal
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
		intVal := int32(plan.AcceptBacklog.ValueInt64())
		addRequest.AcceptBacklog = &intVal
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

// Read a LdapConnectionHandlerResponse object into the model struct
func readLdapConnectionHandlerResponse(ctx context.Context, r *client.LdapConnectionHandlerResponse, state *ldapConnectionHandlerResourceModel, expectedValues *ldapConnectionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ListenAddress = internaltypes.GetStringSet(r.ListenAddress)
	state.ListenPort = types.Int64Value(int64(r.ListenPort))
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
}

// Create any update operations necessary to make the state match the plan
func createLdapConnectionHandlerOperations(plan ldapConnectionHandlerResourceModel, state ldapConnectionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ListenAddress, state.ListenAddress, "listen-address")
	operations.AddInt64OperationIfNecessary(&ops, plan.ListenPort, state.ListenPort, "listen-port")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseSSL, state.UseSSL, "use-ssl")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowStartTLS, state.AllowStartTLS, "allow-start-tls")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.AllowLDAPV2, state.AllowLDAPV2, "allow-ldap-v2")
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

// Create a new resource
func (r *ldapConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddLdapConnectionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumldapConnectionHandlerSchemaUrn{client.ENUMLDAPCONNECTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONNECTION_HANDLERLDAP},
		int32(plan.ListenPort.ValueInt64()),
		plan.Enabled.ValueBool())
	err := addOptionalLdapConnectionHandlerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Ldap Connection Handler", err.Error())
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
		client.AddLdapConnectionHandlerRequestAsAddConnectionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.AddConnectionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ldap Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state ldapConnectionHandlerResourceModel
	readLdapConnectionHandlerResponse(ctx, addResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultLdapConnectionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ldapConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state ldapConnectionHandlerResourceModel
	readLdapConnectionHandlerResponse(ctx, readResponse.LdapConnectionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createLdapConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapConnectionHandlerResponse(ctx, updateResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *ldapConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapConnectionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readLdapConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readLdapConnectionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state ldapConnectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConnectionHandlerApi.GetConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Connection Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapConnectionHandlerResponse(ctx, readResponse.LdapConnectionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *ldapConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultLdapConnectionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateLdapConnectionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateLdapConnectionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan ldapConnectionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state ldapConnectionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConnectionHandlerApi.UpdateConnectionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createLdapConnectionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConnectionHandlerApi.UpdateConnectionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ldap Connection Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readLdapConnectionHandlerResponse(ctx, updateResponse.LdapConnectionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultLdapConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *ldapConnectionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state ldapConnectionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConnectionHandlerApi.DeleteConnectionHandlerExecute(r.apiClient.ConnectionHandlerApi.DeleteConnectionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ldap Connection Handler", err, httpResp)
		return
	}
}

func (r *ldapConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapConnectionHandler(ctx, req, resp)
}

func (r *defaultLdapConnectionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importLdapConnectionHandler(ctx, req, resp)
}

func importLdapConnectionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
