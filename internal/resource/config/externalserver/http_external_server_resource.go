package externalserver

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
	_ resource.Resource                = &httpExternalServerResource{}
	_ resource.ResourceWithConfigure   = &httpExternalServerResource{}
	_ resource.ResourceWithImportState = &httpExternalServerResource{}
	_ resource.Resource                = &defaultHttpExternalServerResource{}
	_ resource.ResourceWithConfigure   = &defaultHttpExternalServerResource{}
	_ resource.ResourceWithImportState = &defaultHttpExternalServerResource{}
)

// Create a Http External Server resource
func NewHttpExternalServerResource() resource.Resource {
	return &httpExternalServerResource{}
}

func NewDefaultHttpExternalServerResource() resource.Resource {
	return &defaultHttpExternalServerResource{}
}

// httpExternalServerResource is the resource implementation.
type httpExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultHttpExternalServerResource is the resource implementation.
type defaultHttpExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *httpExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_external_server"
}

func (r *defaultHttpExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_http_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *httpExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultHttpExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type httpExternalServerResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
	BaseURL                    types.String `tfsdk:"base_url"`
	HostnameVerificationMethod types.String `tfsdk:"hostname_verification_method"`
	KeyManagerProvider         types.String `tfsdk:"key_manager_provider"`
	TrustManagerProvider       types.String `tfsdk:"trust_manager_provider"`
	SslCertNickname            types.String `tfsdk:"ssl_cert_nickname"`
	ConnectTimeout             types.String `tfsdk:"connect_timeout"`
	ResponseTimeout            types.String `tfsdk:"response_timeout"`
	Description                types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *httpExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpExternalServerSchema(ctx, req, resp, false)
}

func (r *defaultHttpExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpExternalServerSchema(ctx, req, resp, true)
}

func httpExternalServerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Http External Server.",
		Attributes: map[string]schema.Attribute{
			"base_url": schema.StringAttribute{
				Description: "The base URL of the external server, optionally including port number, for example \"https://externalService:9031\".",
				Required:    true,
			},
			"hostname_verification_method": schema.StringAttribute{
				Description: "The mechanism for checking if the hostname of the HTTP External Server matches the name(s) stored inside the server's X.509 certificate. This is only applicable if SSL is being used for connection security.",
				Optional:    true,
				Computed:    true,
			},
			"key_manager_provider": schema.StringAttribute{
				Description: "The key manager provider to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property (except when using the Null key manager provider) you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Optional:    true,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use if SSL (HTTPS) is to be used for connection-level security.",
				Optional:    true,
				Computed:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "The certificate alias within the keystore to use if SSL (HTTPS) is to be used for connection-level security. When specifying a value for this property you must ensure that the external server trusts this server's public certificate by adding this server's public certificate to the external server's trust store.",
				Optional:    true,
				Computed:    true,
			},
			"connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established before aborting a request to the server.",
				Optional:    true,
				Computed:    true,
			},
			"response_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to the server.",
				Optional:    true,
				Computed:    true,
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
func addOptionalHttpExternalServerFields(ctx context.Context, addRequest *client.AddHttpExternalServerRequest, plan httpExternalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HostnameVerificationMethod) {
		hostnameVerificationMethod, err := client.NewEnumexternalServerHostnameVerificationMethodPropFromValue(plan.HostnameVerificationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.HostnameVerificationMethod = hostnameVerificationMethod
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
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SslCertNickname) {
		stringVal := plan.SslCertNickname.ValueString()
		addRequest.SslCertNickname = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ConnectTimeout) {
		stringVal := plan.ConnectTimeout.ValueString()
		addRequest.ConnectTimeout = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ResponseTimeout) {
		stringVal := plan.ResponseTimeout.ValueString()
		addRequest.ResponseTimeout = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	return nil
}

// Read a HttpExternalServerResponse object into the model struct
func readHttpExternalServerResponse(ctx context.Context, r *client.HttpExternalServerResponse, state *httpExternalServerResourceModel, expectedValues *httpExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseURL = types.StringValue(r.BaseURL)
	state.HostnameVerificationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerHostnameVerificationMethodProp(r.HostnameVerificationMethod), internaltypes.IsEmptyString(expectedValues.HostnameVerificationMethod))
	state.KeyManagerProvider = internaltypes.StringTypeOrNil(r.KeyManagerProvider, internaltypes.IsEmptyString(expectedValues.KeyManagerProvider))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, internaltypes.IsEmptyString(expectedValues.SslCertNickname))
	state.ConnectTimeout = internaltypes.StringTypeOrNil(r.ConnectTimeout, internaltypes.IsEmptyString(expectedValues.ConnectTimeout))
	config.CheckMismatchedPDFormattedAttributes("connect_timeout",
		expectedValues.ConnectTimeout, state.ConnectTimeout, diagnostics)
	state.ResponseTimeout = internaltypes.StringTypeOrNil(r.ResponseTimeout, internaltypes.IsEmptyString(expectedValues.ResponseTimeout))
	config.CheckMismatchedPDFormattedAttributes("response_timeout",
		expectedValues.ResponseTimeout, state.ResponseTimeout, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createHttpExternalServerOperations(plan httpExternalServerResourceModel, state httpExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseURL, state.BaseURL, "base-url")
	operations.AddStringOperationIfNecessary(&ops, plan.HostnameVerificationMethod, state.HostnameVerificationMethod, "hostname-verification-method")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyManagerProvider, state.KeyManagerProvider, "key-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectTimeout, state.ConnectTimeout, "connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.ResponseTimeout, state.ResponseTimeout, "response-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *httpExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddHttpExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumhttpExternalServerSchemaUrn{client.ENUMHTTPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERHTTP},
		plan.BaseURL.ValueString())
	err := addOptionalHttpExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Http External Server", err.Error())
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
		client.AddHttpExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpExternalServerResourceModel
	readHttpExternalServerResponse(ctx, addResponse.HttpExternalServerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultHttpExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state httpExternalServerResourceModel
	readHttpExternalServerResponse(ctx, readResponse.HttpExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createHttpExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpExternalServerResponse(ctx, updateResponse.HttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *httpExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultHttpExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readHttpExternalServer(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state httpExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpExternalServerResponse(ctx, readResponse.HttpExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *httpExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultHttpExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpExternalServer(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateHttpExternalServer(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan httpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state httpExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createHttpExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpExternalServerResponse(ctx, updateResponse.HttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultHttpExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *httpExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state httpExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Http External Server", err, httpResp)
		return
	}
}

func (r *httpExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpExternalServer(ctx, req, resp)
}

func (r *defaultHttpExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpExternalServer(ctx, req, resp)
}

func importHttpExternalServer(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
