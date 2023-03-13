package httpservletextension

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
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &scim2HttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &scim2HttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &scim2HttpServletExtensionResource{}
)

// Create a Scim2 Http Servlet Extension resource
func NewScim2HttpServletExtensionResource() resource.Resource {
	return &scim2HttpServletExtensionResource{}
}

// scim2HttpServletExtensionResource is the resource implementation.
type scim2HttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *scim2HttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_scim2_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *scim2HttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type scim2HttpServletExtensionResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	BaseContextPath             types.String `tfsdk:"base_context_path"`
	AccessTokenValidator        types.Set    `tfsdk:"access_token_validator"`
	MapAccessTokensToLocalUsers types.String `tfsdk:"map_access_tokens_to_local_users"`
	DebugEnabled                types.Bool   `tfsdk:"debug_enabled"`
	DebugLevel                  types.String `tfsdk:"debug_level"`
	DebugType                   types.Set    `tfsdk:"debug_type"`
	IncludeStackTrace           types.Bool   `tfsdk:"include_stack_trace"`
	SwaggerEnabled              types.Bool   `tfsdk:"swagger_enabled"`
	Description                 types.String `tfsdk:"description"`
	CrossOriginPolicy           types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader              types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *scim2HttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Scim2 Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"base_context_path": schema.StringAttribute{
				Description: "The context path to use to access the SCIM 2.0 interface. The value must start with a forward slash and must represent a valid HTTP context path.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_validator": schema.SetAttribute{
				Description: "If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this SCIM2 HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"map_access_tokens_to_local_users": schema.StringAttribute{
				Description: "Indicates whether the SCIM2 servlet should attempt to map the presented access token to a local user.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"debug_enabled": schema.BoolAttribute{
				Description: "Enables debug logging of the SCIM 2.0 SDK. Debug messages will be forwarded to the Directory Server debug logger with the scope of com.unboundid.directory.broker.http.scim2.extension.SCIM2HTTPServletExtension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"debug_level": schema.StringAttribute{
				Description: "The minimum debug level that should be used for messages to be logged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"debug_type": schema.SetAttribute{
				Description: "The types of debug messages that should be logged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"include_stack_trace": schema.BoolAttribute{
				Description: "Indicates whether a stack trace of the thread which called the debug method should be included in debug log messages.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"swagger_enabled": schema.BoolAttribute{
				Description: "Indicates whether the SCIM2 HTTP Servlet Extension will generate a Swagger specification document.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cross_origin_policy": schema.StringAttribute{
				Description: "The cross-origin request policy to use for the HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"correlation_id_response_header": schema.StringAttribute{
				Description: "Specifies the name of the HTTP response header that will contain a correlation ID value. Example values are \"Correlation-Id\", \"X-Amzn-Trace-Id\", and \"X-Request-Id\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a Scim2HttpServletExtensionResponse object into the model struct
func readScim2HttpServletExtensionResponse(ctx context.Context, r *client.Scim2HttpServletExtensionResponse, state *scim2HttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.MapAccessTokensToLocalUsers = internaltypes.StringTypeOrNil(
		client.StringPointerEnumhttpServletExtensionMapAccessTokensToLocalUsersProp(r.MapAccessTokensToLocalUsers), true)
	state.DebugEnabled = internaltypes.BoolTypeOrNil(r.DebugEnabled)
	state.DebugLevel = types.StringValue(r.DebugLevel.String())
	state.DebugType = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionDebugTypeProp(r.DebugType))
	state.IncludeStackTrace = types.BoolValue(r.IncludeStackTrace)
	state.SwaggerEnabled = internaltypes.BoolTypeOrNil(r.SwaggerEnabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createScim2HttpServletExtensionOperations(plan scim2HttpServletExtensionResourceModel, state scim2HttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.MapAccessTokensToLocalUsers, state.MapAccessTokensToLocalUsers, "map-access-tokens-to-local-users")
	operations.AddBoolOperationIfNecessary(&ops, plan.DebugEnabled, state.DebugEnabled, "debug-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.DebugLevel, state.DebugLevel, "debug-level")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DebugType, state.DebugType, "debug-type")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeStackTrace, state.IncludeStackTrace, "include-stack-trace")
	operations.AddBoolOperationIfNecessary(&ops, plan.SwaggerEnabled, state.SwaggerEnabled, "swagger-enabled")
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
func (r *scim2HttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan scim2HttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim2 Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state scim2HttpServletExtensionResourceModel
	readScim2HttpServletExtensionResponse(ctx, readResponse.Scim2HttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createScim2HttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim2 Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScim2HttpServletExtensionResponse(ctx, updateResponse.Scim2HttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *scim2HttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state scim2HttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim2 Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScim2HttpServletExtensionResponse(ctx, readResponse.Scim2HttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *scim2HttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan scim2HttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state scim2HttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createScim2HttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Scim2 Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readScim2HttpServletExtensionResponse(ctx, updateResponse.Scim2HttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *scim2HttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *scim2HttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
