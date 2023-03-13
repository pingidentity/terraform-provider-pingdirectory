package httpservletextension

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
	_ resource.Resource                = &directoryRestApiHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &directoryRestApiHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &directoryRestApiHttpServletExtensionResource{}
)

// Create a Directory Rest Api Http Servlet Extension resource
func NewDirectoryRestApiHttpServletExtensionResource() resource.Resource {
	return &directoryRestApiHttpServletExtensionResource{}
}

// directoryRestApiHttpServletExtensionResource is the resource implementation.
type directoryRestApiHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *directoryRestApiHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_directory_rest_api_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *directoryRestApiHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type directoryRestApiHttpServletExtensionResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	BasicAuthEnabled            types.Bool   `tfsdk:"basic_auth_enabled"`
	IdentityMapper              types.String `tfsdk:"identity_mapper"`
	AccessTokenValidator        types.Set    `tfsdk:"access_token_validator"`
	AccessTokenScope            types.String `tfsdk:"access_token_scope"`
	Audience                    types.String `tfsdk:"audience"`
	MaxPageSize                 types.Int64  `tfsdk:"max_page_size"`
	SchemasEndpointObjectclass  types.Set    `tfsdk:"schemas_endpoint_objectclass"`
	DefaultOperationalAttribute types.Set    `tfsdk:"default_operational_attribute"`
	RejectExpansionAttribute    types.Set    `tfsdk:"reject_expansion_attribute"`
	AllowedControl              types.Set    `tfsdk:"allowed_control"`
	Description                 types.String `tfsdk:"description"`
	CrossOriginPolicy           types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader              types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *directoryRestApiHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Directory Rest Api Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"basic_auth_enabled": schema.BoolAttribute{
				Description: "Enables HTTP Basic authentication, using a username and password. The Identity Mapper specified by the identity-mapper property will be used to map the username to a DN.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the Identity Mapper that is to be used for associating user entries with basic authentication usernames.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_validator": schema.SetAttribute{
				Description: "If specified, the Access Token Validator(s) that may be used to validate access tokens for requests submitted to this Directory REST API HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"access_token_scope": schema.StringAttribute{
				Description: "The name of a scope that must be present in an access token accepted by the Directory REST API HTTP Servlet Extension.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"audience": schema.StringAttribute{
				Description: "A string or URI that identifies the Directory REST API HTTP Servlet Extension in the context of OAuth2 authorization.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"max_page_size": schema.Int64Attribute{
				Description: "The maximum number of entries to be returned in one page of search results.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"schemas_endpoint_objectclass": schema.SetAttribute{
				Description: "The list of object classes which will be returned by the schemas endpoint.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"default_operational_attribute": schema.SetAttribute{
				Description: "A set of operational attributes that will be returned with entries by default.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"reject_expansion_attribute": schema.SetAttribute{
				Description: "A set of attributes which the client is not allowed to provide for the expand query parameters. This should be used for attributes that could either have a large number of values or that reference entries that are very large like groups.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"allowed_control": schema.SetAttribute{
				Description: "Specifies the names of any request controls that should be allowed by the Directory REST API. Any request that contains a critical control not in this list will be rejected. Any non-critical request control which is not supported by the Directory REST API will be removed from the request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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

// Read a DirectoryRestApiHttpServletExtensionResponse object into the model struct
func readDirectoryRestApiHttpServletExtensionResponse(ctx context.Context, r *client.DirectoryRestApiHttpServletExtensionResponse, state *directoryRestApiHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BasicAuthEnabled = internaltypes.BoolTypeOrNil(r.BasicAuthEnabled)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, true)
	state.AccessTokenValidator = internaltypes.GetStringSet(r.AccessTokenValidator)
	state.AccessTokenScope = internaltypes.StringTypeOrNil(r.AccessTokenScope, true)
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, true)
	state.MaxPageSize = internaltypes.Int64TypeOrNil(r.MaxPageSize)
	state.SchemasEndpointObjectclass = internaltypes.GetStringSet(r.SchemasEndpointObjectclass)
	state.DefaultOperationalAttribute = internaltypes.GetStringSet(r.DefaultOperationalAttribute)
	state.RejectExpansionAttribute = internaltypes.GetStringSet(r.RejectExpansionAttribute)
	state.AllowedControl = internaltypes.GetStringSet(
		client.StringSliceEnumhttpServletExtensionAllowedControlProp(r.AllowedControl))
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, true)
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDirectoryRestApiHttpServletExtensionOperations(plan directoryRestApiHttpServletExtensionResourceModel, state directoryRestApiHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.BasicAuthEnabled, state.BasicAuthEnabled, "basic-auth-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AccessTokenValidator, state.AccessTokenValidator, "access-token-validator")
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTokenScope, state.AccessTokenScope, "access-token-scope")
	operations.AddStringOperationIfNecessary(&ops, plan.Audience, state.Audience, "audience")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxPageSize, state.MaxPageSize, "max-page-size")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SchemasEndpointObjectclass, state.SchemasEndpointObjectclass, "schemas-endpoint-objectclass")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultOperationalAttribute, state.DefaultOperationalAttribute, "default-operational-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.RejectExpansionAttribute, state.RejectExpansionAttribute, "reject-expansion-attribute")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedControl, state.AllowedControl, "allowed-control")
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
func (r *directoryRestApiHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan directoryRestApiHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Directory Rest Api Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state directoryRestApiHttpServletExtensionResourceModel
	readDirectoryRestApiHttpServletExtensionResponse(ctx, readResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDirectoryRestApiHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Directory Rest Api Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDirectoryRestApiHttpServletExtensionResponse(ctx, updateResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *directoryRestApiHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state directoryRestApiHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Directory Rest Api Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDirectoryRestApiHttpServletExtensionResponse(ctx, readResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *directoryRestApiHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan directoryRestApiHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state directoryRestApiHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDirectoryRestApiHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Directory Rest Api Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDirectoryRestApiHttpServletExtensionResponse(ctx, updateResponse.DirectoryRestApiHttpServletExtensionResponse, &state, &resp.Diagnostics)
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
func (r *directoryRestApiHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *directoryRestApiHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
