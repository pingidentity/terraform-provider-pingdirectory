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
	_ resource.Resource                = &pingOneHttpExternalServerResource{}
	_ resource.ResourceWithConfigure   = &pingOneHttpExternalServerResource{}
	_ resource.ResourceWithImportState = &pingOneHttpExternalServerResource{}
)

// Create a Ping One Http External Server resource
func NewPingOneHttpExternalServerResource() resource.Resource {
	return &pingOneHttpExternalServerResource{}
}

// pingOneHttpExternalServerResource is the resource implementation.
type pingOneHttpExternalServerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pingOneHttpExternalServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ping_one_http_external_server"
}

// Configure adds the provider configured client to the resource.
func (r *pingOneHttpExternalServerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type pingOneHttpExternalServerResourceModel struct {
	Id                         types.String `tfsdk:"id"`
	LastUpdated                types.String `tfsdk:"last_updated"`
	Notifications              types.Set    `tfsdk:"notifications"`
	RequiredActions            types.Set    `tfsdk:"required_actions"`
	HostnameVerificationMethod types.String `tfsdk:"hostname_verification_method"`
	TrustManagerProvider       types.String `tfsdk:"trust_manager_provider"`
	ConnectTimeout             types.String `tfsdk:"connect_timeout"`
	ResponseTimeout            types.String `tfsdk:"response_timeout"`
	Description                types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *pingOneHttpExternalServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Ping One Http External Server.",
		Attributes: map[string]schema.Attribute{
			"hostname_verification_method": schema.StringAttribute{
				Description: "The mechanism for checking if the hostname in the PingOne ID Token Validator's base-url value matches the name(s) stored inside the X.509 certificate presented by PingOne.",
				Optional:    true,
				Computed:    true,
			},
			"trust_manager_provider": schema.StringAttribute{
				Description: "The trust manager provider to use for HTTPS connection-level security.",
				Optional:    true,
				Computed:    true,
			},
			"connect_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for a connection to be established before aborting a request to PingOne.",
				Optional:    true,
				Computed:    true,
			},
			"response_timeout": schema.StringAttribute{
				Description: "Specifies the maximum length of time to wait for response data to be read from an established connection before aborting a request to PingOne.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this External Server",
				Optional:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalPingOneHttpExternalServerFields(ctx context.Context, addRequest *client.AddPingOneHttpExternalServerRequest, plan pingOneHttpExternalServerResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HostnameVerificationMethod) {
		hostnameVerificationMethod, err := client.NewEnumexternalServerHostnameVerificationMethodPropFromValue(plan.HostnameVerificationMethod.ValueString())
		if err != nil {
			return err
		}
		addRequest.HostnameVerificationMethod = hostnameVerificationMethod
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustManagerProvider) {
		stringVal := plan.TrustManagerProvider.ValueString()
		addRequest.TrustManagerProvider = &stringVal
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

// Read a PingOneHttpExternalServerResponse object into the model struct
func readPingOneHttpExternalServerResponse(ctx context.Context, r *client.PingOneHttpExternalServerResponse, state *pingOneHttpExternalServerResourceModel, expectedValues *pingOneHttpExternalServerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.HostnameVerificationMethod = internaltypes.StringTypeOrNil(
		client.StringPointerEnumexternalServerHostnameVerificationMethodProp(r.HostnameVerificationMethod), internaltypes.IsEmptyString(expectedValues.HostnameVerificationMethod))
	state.TrustManagerProvider = internaltypes.StringTypeOrNil(r.TrustManagerProvider, internaltypes.IsEmptyString(expectedValues.TrustManagerProvider))
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
func createPingOneHttpExternalServerOperations(plan pingOneHttpExternalServerResourceModel, state pingOneHttpExternalServerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.HostnameVerificationMethod, state.HostnameVerificationMethod, "hostname-verification-method")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustManagerProvider, state.TrustManagerProvider, "trust-manager-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.ConnectTimeout, state.ConnectTimeout, "connect-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.ResponseTimeout, state.ResponseTimeout, "response-timeout")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *pingOneHttpExternalServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pingOneHttpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPingOneHttpExternalServerRequest(plan.Id.ValueString(),
		[]client.EnumpingOneHttpExternalServerSchemaUrn{client.ENUMPINGONEHTTPEXTERNALSERVERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTERNAL_SERVERPING_ONE_HTTP})
	err := addOptionalPingOneHttpExternalServerFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Ping One Http External Server", err.Error())
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
		client.AddPingOneHttpExternalServerRequestAsAddExternalServerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExternalServerApi.AddExternalServerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Ping One Http External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pingOneHttpExternalServerResourceModel
	readPingOneHttpExternalServerResponse(ctx, addResponse.PingOneHttpExternalServerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *pingOneHttpExternalServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state pingOneHttpExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExternalServerApi.GetExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ping One Http External Server", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPingOneHttpExternalServerResponse(ctx, readResponse.PingOneHttpExternalServerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pingOneHttpExternalServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan pingOneHttpExternalServerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pingOneHttpExternalServerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ExternalServerApi.UpdateExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPingOneHttpExternalServerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExternalServerApi.UpdateExternalServerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Ping One Http External Server", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPingOneHttpExternalServerResponse(ctx, updateResponse.PingOneHttpExternalServerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *pingOneHttpExternalServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pingOneHttpExternalServerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExternalServerApi.DeleteExternalServerExecute(r.apiClient.ExternalServerApi.DeleteExternalServer(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Ping One Http External Server", err, httpResp)
		return
	}
}

func (r *pingOneHttpExternalServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
