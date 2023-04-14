package externalserver

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
	_ resource.Resource                = &nokiaDsExternalServerResource{}
	_ resource.ResourceWithConfigure   = &nokiaDsExternalServerResource{}
	_ resource.ResourceWithImportState = &nokiaDsExternalServerResource{}
	_ resource.Resource                = &defaultNokiaDsExternalServerResource{}
	_ resource.ResourceWithConfigure   = &defaultNokiaDsExternalServerResource{}
	_ resource.ResourceWithImportState = &defaultNokiaDsExternalServerResource{}
)

// Create a Nokia Ds External Server resource
func NewNokiaDsExternalServerResource() resource.Resource {
	return &nokiaDsExternalServerResource{}
}

func NewDefaultNokiaDsExternalServerResource() resource.Resource {
	return &defaultNokiaDsExternalServerResource{}
}

// nokiaDsExternalServerResource is the resource implementation.
type nokiaDsExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultNokiaDsExternalServerResource is the resource implementation.
type defaultNokiaDsExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *nokiaDsExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_nokia_ds_external_server"
}

func (r *defaultNokiaDsExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_nokia_ds_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *nokiaDsExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultNokiaDsExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type nokiaDsExternalServerResourceModel struct {
	Id                                     types.String `tfsdk:"id"`
	LastUpdated                            types.String `tfsdk:"last_updated"`
	Notifications                          types.Set    `tfsdk:"notifications"`
	RequiredActions                        types.Set    `tfsdk:"required_actions"`
	VerifyCredentialsMethod                types.String `tfsdk:"verify_credentials_method"`
	UseAdministrativeOperationControl      types.Bool   `tfsdk:"use_administrative_operation_control"`
	ServerHostName                         types.String `tfsdk:"server_host_name"`
	ServerPort                             types.Int64  `tfsdk:"server_port"`
	Location                               types.String `tfsdk:"location"`
	BindDN                                 types.String `tfsdk:"bind_dn"`
	Password                               types.String `tfsdk:"password"`
	PassphraseProvider                     types.String `tfsdk:"passphrase_provider"`
	ConnectionSecurity                     types.String `tfsdk:"connection_security"`
	AuthenticationMethod                   types.String `tfsdk:"authentication_method"`
	HealthCheckConnectTimeout              types.String `tfsdk:"health_check_connect_timeout"`
	MaxConnectionAge                       types.String `tfsdk:"max_connection_age"`
	MinExpiredConnectionDisconnectInterval types.String `tfsdk:"min_expired_connection_disconnect_interval"`
	ConnectTimeout                         types.String `tfsdk:"connect_timeout"`
	MaxResponseSize                        types.String `tfsdk:"max_response_size"`
	KeyManagerProvider                     types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider                   types.String `tfsdk:"trust_manager_provider"`
	InitialConnections                     types.Int64  `tfsdk:"initial_connections"`
	MaxConnections                         types.Int64  `tfsdk:"max_connections"`
	DefunctConnectionResultCode            types.Set    `tfsdk:"defunct_connection_result_code"`
	AbandonOnTimeout                       types.Bool   `tfsdk:"abandon_on_timeout"`
	Description                            types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *nokiaDsExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	nokiaDsExternalServerSchema(ctx, req, resp, false)
}

func (r *defaultNokiaDsExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	nokiaDsExternalServerSchema(ctx, req, resp, true)
}

func nokiaDsExternalServerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Nokia Ds External Server.",
		Attributes: map[string]schema.Attribute{
			"verify_credentials_method": schema.StringAttribute{
				Description: "The mechanism to use to verify user credentials while ensuring that the ability to process other operations is not impacted by an alternate authorization identity.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_administrative_operation_control": schema.BoolAttribute{
				Description: "Indicates whether to include the administrative operation request control in requests sent to this server which are intended for administrative operations (e.g., health checking) rather than requests directly from clients.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"server_host_name": schema.StringAttribute{
				Description: "The host name or IP address of the target LDAP server.",
				Required:    true,
			},
			"server_port": schema.Int64Attribute{
				Description: "The port number on which the server listens for requests.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"location": schema.StringAttribute{
				Description: "Specifies the location for the LDAP External Server.",
				Optional:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "The DN to use to bind to the target LDAP server if simple authentication is required.",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The login password for the specified user.",
				Optional:    true,
				Sensitive:   true,
			},
			"passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the login password for the specified user.",
				Optional:    true,
			},
			"connection_security": schema.StringAttribute{
				Description: "The mechanism to use to secure communication with the directory server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"authentication_method": schema.StringAttribute{
				Description: "The mechanism to use to authenticate to the target server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"health_check_connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established for the purpose of performing a health check. If the connection cannot be established within this length of time, the server will be classified as unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_connection_age": schema.StringAttribute{
				Description: "Specifies the maximum length of time that connections to this server should be allowed to remain established before being closed and replaced with newly-established connections.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"min_expired_connection_disconnect_interval": schema.StringAttribute{
				Description: "Specifies the minimum length of time that should pass between connection closures as a result of the connections being established for longer than the maximum connection age. This may help avoid cases in which a large number of connections are closed and re-established in a short period of time because of the maximum connection age.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established before giving up and considering the server unavailable.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_response_size": schema.StringAttribute{
				Description: "Specifies the maximum response size that should be supported for messages received from the LDAP external server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "The key manager provider to use if SSL or StartTLS is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Optional:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use if SSL or StartTLS is to be used for connection-level security.",
				Optional:    true,
			},
			"initial_connections": schema.Int64Attribute{
				Description: "The number of connections to initially establish to the LDAP external server. A value of zero indicates that the number of connections should be dynamically based on the number of available worker threads. This will be ignored when using a thread-local connection pool.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"max_connections": schema.Int64Attribute{
				Description: "The maximum number of concurrent connections to maintain for the LDAP external server. A value of zero indicates that the number of connections should be dynamically based on the number of available worker threads. This will be ignored when using a thread-local connection pool.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"defunct_connection_result_code": schema.SetAttribute{
				Description: "Specifies the operation result code values that should cause the associated connection should be considered defunct. If an operation fails with one of these result codes, then it will be terminated and an attempt will be made to establish a new connection in its place.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"abandon_on_timeout": schema.BoolAttribute{
				Description: "Indicates whether to send an abandon request for an operation for which a response timeout is encountered. A request which has timed out on one server may be retried on another server regardless of whether an abandon request is sent, but if the initial attempt is not abandoned then a long-running operation may unnecessarily continue to consume processing resources on the initial server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
				Optional:    true,
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
func addOptionalNokiaDsExternalServerFields(ctx context.Context, addRequest *client.AddNokiaDsExternalServerRequest, plan nokiaDsExternalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VerifyCredentialsMethod) {
		verifyCredentialsMethod, err := client.NewEnumexternalServerVerifyCredentialsMethodPropFromValue(plan.VerifyCredentialsMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.VerifyCredentialsMethod = verifyCredentialsMethod
	}
	if internaltypes.IsDefined(plan.UseAdministrativeOperationControl) {
		addRequest.UseAdministrativeOperationControl = plan.UseAdministrativeOperationControl.ValueBoolPointer()
	}
	if internaltypes.IsDefined(plan.ServerPort) {
		addRequest.ServerPort = plan.ServerPort.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Location) {
		addRequest.Location = plan.Location.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.BindDN) {
		addRequest.BindDN = plan.BindDN.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PassphraseProvider) {
		addRequest.PassphraseProvider = plan.PassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectionSecurity) {
		connectionSecurity, err := client.NewEnumexternalServerConnectionSecurityPropFromValue(plan.ConnectionSecurity.ValueString())
		if err != nil {
			return err
		}
		addRequest.ConnectionSecurity = connectionSecurity
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthenticationMethod) {
		authenticationMethod, err := client.NewEnumexternalServerNokiaDsAuthenticationMethodPropFromValue(plan.AuthenticationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.AuthenticationMethod = authenticationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HealthCheckConnectTimeout) {
		addRequest.HealthCheckConnectTimeout = plan.HealthCheckConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxConnectionAge) {
		addRequest.MaxConnectionAge = plan.MaxConnectionAge.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MinExpiredConnectionDisconnectInterval) {
		addRequest.MinExpiredConnectionDisconnectInterval = plan.MinExpiredConnectionDisconnectInterval.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		addRequest.ConnectTimeout = plan.ConnectTimeout.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxResponseSize) {
		addRequest.MaxResponseSize = plan.MaxResponseSize.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyManagerProvider) {
		addRequest.KeyManagerProvider = plan.KeyManagerProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		addRequest.TrustManagerProvider = plan.TrustManagerProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitialConnections) {
		addRequest.InitialConnections = plan.InitialConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxConnections) {
		addRequest.MaxConnections = plan.MaxConnections.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DefunctConnectionResultCode) {
		var slice []string
		plan.DefunctConnectionResultCode.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumexternalServerDefunctConnectionResultCodeProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumexternalServerDefunctConnectionResultCodePropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.DefunctConnectionResultCode = enumSlice
	}
	if internaltypes.IsDefined(plan.AbandonOnTimeout) {
		addRequest.AbandonOnTimeout = plan.AbandonOnTimeout.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Read a NokiaDsExternalServerResponse object into the model struct
func readNokiaDsExternalServerResponse(ctx context.Context, r *client.NokiaDsExternalServerResponse, state *nokiaDsExternalServerResourceModel, expectedValues *nokiaDsExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.VerifyCredentialsMethod = types.StringValue(r.VerifyCredentialsMethod.String())
	state.UseAdministrativeOperationControl = internaltypes.BoolTypeOrNil(r.UseAdministrativeOperationControl)
	state.ServerHostName = types.StringValue(r.ServerHostName)
	state.ServerPort = types.Int64Value(r.ServerPort)
	state.Location = internaltypes.StringTypeOrNil(r.Location, internaltypes.IsEmptyString(expectedValues.Location))
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, internaltypes.IsEmptyString(expectedValues.BindDN))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.Password = expectedValues.Password
	state.PassphraseProvider = internaltypes.StringTypeOrNil(r.PassphraseProvider, internaltypes.IsEmptyString(expectedValues.PassphraseProvider))
	state.ConnectionSecurity = types.StringValue(r.ConnectionSecurity.String())
	state.AuthenticationMethod = types.StringValue(r.AuthenticationMethod.String())
	state.HealthCheckConnectTimeout = internaltypes.StringTypeOrNil(r.HealthCheckConnectTimeout, internaltypes.IsEmptyString(expectedValues.HealthCheckConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("health_check_connect_timeout",
		expectedValues.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, diagnostics)
	state.MaxConnectionAge = types.StringValue(r.MaxConnectionAge)
	config.CheckMismatchedPDFormattedAttributes("max_connection_age",
		expectedValues.MaxConnectionAge, state.MaxConnectionAge, diagnostics)
	state.MinExpiredConnectionDisconnectInterval = internaltypes.StringTypeOrNil(r.MinExpiredConnectionDisconnectInterval, internaltypes.IsEmptyString(expectedValues.MinExpiredConnectionDisconnectInterval))
	config.CheckMismatchedPDFormattedAttributes("min_expired_connection_disconnect_interval",
		expectedValues.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, diagnostics)
	state.ConnectTimeout = types.StringValue(r.ConnectTimeout)
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.MaxResponseSize = types.StringValue(r.MaxResponseSize)
	config.CheckMismatchedPDFormattedAttributes("max_response_size",
		expectedValues.MaxResponseSize, state.MaxResponseSize, diagnostics)
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.InitialConnections = internaltypes.Int64TypeOrNil(r.InitialConnections)
	state.MaxConnections = internaltypes.Int64TypeOrNil(r.MaxConnections)
	state.DefunctConnectionResultCode = internaltypes.GetStringSet(
		client.StringSliceEnumexternalServerDefunctConnectionResultCodeProp(r.DefunctConnectionResultCode))
	state.AbandonOnTimeout = internaltypes.BoolTypeOrNil(r.AbandonOnTimeout)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createNokiaDsExternalServerOperations(plan nokiaDsExternalServerResourceModel, state nokiaDsExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.VerifyCredentialsMethod, state.VerifyCredentialsMethod, "verify-credentials-method")
	operations.AddBoolOperationIfNecessary(&ops, plan.UseAdministrativeOperationControl, state.UseAdministrativeOperationControl, "use-administrative-operation-control")
	operations.AddStringOperationIfNecessary(&ops, plan.ServerHostName, state.ServerHostName, "server-host-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerPort, state.ServerPort, "server-port")
	operations.AddStringOperationIfNecessary(&ops, plan.Location, state.Location, "location")
	operations.AddStringOperationIfNecessary(&ops, plan.BindDN, state.BindDN, "bind-dn")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.PassphraseProvider, state.PassphraseProvider, "passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectionSecurity, state.ConnectionSecurity, "connection-security")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthenticationMethod, state.AuthenticationMethod, "authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HealthCheckConnectTimeout, state.HealthCheckConnectTimeout, "health-check-connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxConnectionAge, state.MaxConnectionAge, "max-connection-age")
	operations.AddStringOperationIfNecessary(&ops, plan.MinExpiredConnectionDisconnectInterval, state.MinExpiredConnectionDisconnectInterval, "min-expired-connection-disconnect-interval")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectTimeout, state.ConnectTimeout, "connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxResponseSize, state.MaxResponseSize, "max-response-size")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddInt64OperationIfNecessary(&ops, plan.InitialConnections, state.InitialConnections, "initial-connections")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxConnections, state.MaxConnections, "max-connections")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefunctConnectionResultCode, state.DefunctConnectionResultCode, "defunct-connection-result-code")
	operations.AddBoolOperationIfNecessary(&ops, plan.AbandonOnTimeout, state.AbandonOnTimeout, "abandon-on-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *nokiaDsExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan nokiaDsExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddNokiaDsExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumnokiaDsExternalServerSchemaUrn{client.ENUMNOKIADSEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERNOKIA_DS},
		plan.ServerHostName.ValueString())
	err := addOptionalNokiaDsExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Nokia Ds External Server", err.Error())
		return
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExternalServerApi.AddExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExternalServerRequest(
		client.AddNokiaDsExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Nokia Ds External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state nokiaDsExternalServerResourceModel
	readNokiaDsExternalServerResponse(ctx, addResponse.NokiaDsExternalServerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultNokiaDsExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan nokiaDsExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Nokia Ds External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state nokiaDsExternalServerResourceModel
	readNokiaDsExternalServerResponse(ctx, readResponse.NokiaDsExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createNokiaDsExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Nokia Ds External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readNokiaDsExternalServerResponse(ctx, updateResponse.NokiaDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *nokiaDsExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readNokiaDsExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultNokiaDsExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readNokiaDsExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readNokiaDsExternalServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state nokiaDsExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Nokia Ds External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readNokiaDsExternalServerResponse(ctx, readResponse.NokiaDsExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *nokiaDsExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateNokiaDsExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultNokiaDsExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateNokiaDsExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateNokiaDsExternalServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan nokiaDsExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state nokiaDsExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createNokiaDsExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Nokia Ds External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readNokiaDsExternalServerResponse(ctx, updateResponse.NokiaDsExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultNokiaDsExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *nokiaDsExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state nokiaDsExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Nokia Ds External Server", err, httpResp)
		return
	}
}

func (r *nokiaDsExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importNokiaDsExternalServer(ctx, req, resp)
}

func (r *defaultNokiaDsExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importNokiaDsExternalServer(ctx, req, resp)
}

func importNokiaDsExternalServer(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
