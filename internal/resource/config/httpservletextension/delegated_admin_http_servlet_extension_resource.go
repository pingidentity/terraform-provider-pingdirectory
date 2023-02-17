package httpservletextension

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
	client "github.com/pingidentity/pingdirectory-go-client/v9100"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &delegatedAdminHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &delegatedAdminHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &delegatedAdminHttpServletExtensionResource{}
)

// Create a Delegated Admin Http Servlet Extension resource
func NewDelegatedAdminHttpServletExtensionResource() resource.Resource {
	return &delegatedAdminHttpServletExtensionResource{}
}

// delegatedAdminHttpServletExtensionResource is the resource implementation.
type delegatedAdminHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *delegatedAdminHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *delegatedAdminHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminHttpServletExtensionResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	BasicAuthEnabled            types.Bool   `tfsdk:"basic_auth_enabled"`
	IdentityMapper              types.String `tfsdk:"identity_mapper"`
	AccessTokenValidator        types.Set    `tfsdk:"access_token_validator"`
	AccessTokenScope            types.String `tfsdk:"access_token_scope"`
	Audience                    types.String `tfsdk:"audience"`
	Description                 types.String `tfsdk:"description"`
	CrossOriginPolicy           types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader              types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *delegatedAdminHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Delegated Admin Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"basic_auth_enabled": schema.BoolAttribute{
				Description: "Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN.",
				Optional:    true,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the Identity Mapper that is to be used for associating user entries with basic authentication user names.",
				Optional:    true,
				Computed:    true,
			},
			"access_token_validator": schema.SetAttribute{
				Description: "If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Delegated Admin HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"access_token_scope": schema.StringAttribute{
				Description: "The name of a scope that must be present in an access token accepted by the Delegated Admin HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
			},
			"audience": schema.StringAttribute{
				Description: "A string or URI that identifies the Delegated Admin HTTP Servlet Extension in the context of OAuth2 authorization.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
				Computed:    true,
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
			},
			"response_header": schema.SetAttribute{
				Description: "Specifies HTTP header fields and values added to response headers for all requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a DelegatedAdminHttpServletExtensionResponse object into the model struct
func readDelegatedAdminHttpServletExtensionResponse(ctx context.Context, r *client.DelegatedAdminHttpServletExtensionResponse, state *delegatedAdminHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.AccessTokenScope = internaltypes.StringTypeOrNil(r.AccessTokenScope, true)
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, true)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDelegatedAdminHttpServletExtensionOperations(plan delegatedAdminHttpServletExtensionResourceModel, state delegatedAdminHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.BasicAuthEnabled, state.BasicAuthEnabled, "basic-auth-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTokenScope, state.AccessTokenScope, "access-token-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Audience, state.Audience, "audience")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *delegatedAdminHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state delegatedAdminHttpServletExtensionResourceModel
	readDelegatedAdminHttpServletExtensionResponse(ctx, readResponse.DelegatedAdminHttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDelegatedAdminHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminHttpServletExtensionResponse(ctx, updateResponse.DelegatedAdminHttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *delegatedAdminHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state delegatedAdminHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminHttpServletExtensionResponse(ctx, readResponse.DelegatedAdminHttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *delegatedAdminHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan delegatedAdminHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state delegatedAdminHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDelegatedAdminHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Delegated Admin Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDelegatedAdminHttpServletExtensionResponse(ctx, updateResponse.DelegatedAdminHttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *delegatedAdminHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *delegatedAdminHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
