package config

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
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &httpServletCrossOriginPolicyResource{}
	_ resource.ResourceWithConfigure   = &httpServletCrossOriginPolicyResource{}
	_ resource.ResourceWithImportState = &httpServletCrossOriginPolicyResource{}
	_ resource.Resource                = &defaultHttpServletCrossOriginPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultHttpServletCrossOriginPolicyResource{}
	_ resource.ResourceWithImportState = &defaultHttpServletCrossOriginPolicyResource{}
)

// Create a Http Servlet Cross Origin Policy resource
func NewHttpServletCrossOriginPolicyResource() resource.Resource {
	return &httpServletCrossOriginPolicyResource{}
}

func NewDefaultHttpServletCrossOriginPolicyResource() resource.Resource {
	return &defaultHttpServletCrossOriginPolicyResource{}
}

// httpServletCrossOriginPolicyResource is the resource implementation.
type httpServletCrossOriginPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultHttpServletCrossOriginPolicyResource is the resource implementation.
type defaultHttpServletCrossOriginPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *httpServletCrossOriginPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_http_servlet_cross_origin_policy"
}

func (r *defaultHttpServletCrossOriginPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_http_servlet_cross_origin_policy"
}

// Configure adds the provider configured client to the resource.
func (r *httpServletCrossOriginPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultHttpServletCrossOriginPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type httpServletCrossOriginPolicyResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Description          types.String `tfsdk:"description"`
	CorsAllowedMethods   types.Set    `tfsdk:"cors_allowed_methods"`
	CorsAllowedOrigins   types.Set    `tfsdk:"cors_allowed_origins"`
	CorsExposedHeaders   types.Set    `tfsdk:"cors_exposed_headers"`
	CorsAllowedHeaders   types.Set    `tfsdk:"cors_allowed_headers"`
	CorsPreflightMaxAge  types.String `tfsdk:"cors_preflight_max_age"`
	CorsAllowCredentials types.Bool   `tfsdk:"cors_allow_credentials"`
}

// GetSchema defines the schema for the resource.
func (r *httpServletCrossOriginPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpServletCrossOriginPolicySchema(ctx, req, resp, false)
}

func (r *defaultHttpServletCrossOriginPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	httpServletCrossOriginPolicySchema(ctx, req, resp, true)
}

func httpServletCrossOriginPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Http Servlet Cross Origin Policy.",
		Attributes: map[string]schema.Attribute{
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Cross Origin Policy",
				Optional:    true,
			},
			"cors_allowed_methods": schema.SetAttribute{
				Description: "A list of HTTP methods allowed for cross-origin access to resources. i.e. one or more of GET, POST, PUT, DELETE, etc.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"cors_allowed_origins": schema.SetAttribute{
				Description: "A list of origins that are allowed to execute cross-origin requests.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"cors_exposed_headers": schema.SetAttribute{
				Description: "A list of HTTP headers other than the simple response headers that browsers are allowed to access.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"cors_allowed_headers": schema.SetAttribute{
				Description: "A list of HTTP headers that are supported by the resource and can be specified in a cross-origin request.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"cors_preflight_max_age": schema.StringAttribute{
				Description: "The maximum amount of time that a preflight request can be cached by a client.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cors_allow_credentials": schema.BoolAttribute{
				Description: "Indicates whether the servlet extension allows CORS requests with username/password credentials.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		SetAllAttributesToOptionalAndComputed(&schemaDef, []string{"id"})
	}
	AddCommonSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for http-servlet-cross-origin-policy http-servlet-cross-origin-policy
func addOptionalHttpServletCrossOriginPolicyFields(ctx context.Context, addRequest *client.AddHttpServletCrossOriginPolicyRequest, plan httpServletCrossOriginPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CorsAllowedMethods) {
		var slice []string
		plan.CorsAllowedMethods.ElementsAs(ctx, &slice, false)
		addRequest.CorsAllowedMethods = slice
	}
	if internaltypes.IsDefined(plan.CorsAllowedOrigins) {
		var slice []string
		plan.CorsAllowedOrigins.ElementsAs(ctx, &slice, false)
		addRequest.CorsAllowedOrigins = slice
	}
	if internaltypes.IsDefined(plan.CorsExposedHeaders) {
		var slice []string
		plan.CorsExposedHeaders.ElementsAs(ctx, &slice, false)
		addRequest.CorsExposedHeaders = slice
	}
	if internaltypes.IsDefined(plan.CorsAllowedHeaders) {
		var slice []string
		plan.CorsAllowedHeaders.ElementsAs(ctx, &slice, false)
		addRequest.CorsAllowedHeaders = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorsPreflightMaxAge) {
		addRequest.CorsPreflightMaxAge = plan.CorsPreflightMaxAge.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.CorsAllowCredentials) {
		addRequest.CorsAllowCredentials = plan.CorsAllowCredentials.ValueBoolPointer()
	}
}

// Read a HttpServletCrossOriginPolicyResponse object into the model struct
func readHttpServletCrossOriginPolicyResponse(ctx context.Context, r *client.HttpServletCrossOriginPolicyResponse, state *httpServletCrossOriginPolicyResourceModel, expectedValues *httpServletCrossOriginPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CorsAllowedMethods = internaltypes.GetStringSet(r.CorsAllowedMethods)
	state.CorsAllowedOrigins = internaltypes.GetStringSet(r.CorsAllowedOrigins)
	state.CorsExposedHeaders = internaltypes.GetStringSet(r.CorsExposedHeaders)
	state.CorsAllowedHeaders = internaltypes.GetStringSet(r.CorsAllowedHeaders)
	state.CorsPreflightMaxAge = internaltypes.StringTypeOrNil(r.CorsPreflightMaxAge, internaltypes.IsEmptyString(expectedValues.CorsPreflightMaxAge))
	CheckMismatchedPDFormattedAttributes("cors_preflight_max_age",
		expectedValues.CorsPreflightMaxAge, state.CorsPreflightMaxAge, diagnostics)
	state.CorsAllowCredentials = internaltypes.BoolTypeOrNil(r.CorsAllowCredentials)
	state.Notifications, state.RequiredActions = ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createHttpServletCrossOriginPolicyOperations(plan httpServletCrossOriginPolicyResourceModel, state httpServletCrossOriginPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CorsAllowedMethods, state.CorsAllowedMethods, "cors-allowed-methods")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CorsAllowedOrigins, state.CorsAllowedOrigins, "cors-allowed-origins")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CorsExposedHeaders, state.CorsExposedHeaders, "cors-exposed-headers")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.CorsAllowedHeaders, state.CorsAllowedHeaders, "cors-allowed-headers")
	operations.AddStringOperationIfNecessary(&ops, plan.CorsPreflightMaxAge, state.CorsPreflightMaxAge, "cors-preflight-max-age")
	operations.AddBoolOperationIfNecessary(&ops, plan.CorsAllowCredentials, state.CorsAllowCredentials, "cors-allow-credentials")
	return ops
}

// Create a http-servlet-cross-origin-policy http-servlet-cross-origin-policy
func (r *httpServletCrossOriginPolicyResource) CreateHttpServletCrossOriginPolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan httpServletCrossOriginPolicyResourceModel) (*httpServletCrossOriginPolicyResourceModel, error) {
	addRequest := client.NewAddHttpServletCrossOriginPolicyRequest(plan.Id.ValueString())
	addOptionalHttpServletCrossOriginPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletCrossOriginPolicyApi.AddHttpServletCrossOriginPolicy(
		ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletCrossOriginPolicyRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.HttpServletCrossOriginPolicyApi.AddHttpServletCrossOriginPolicyExecute(apiAddRequest)
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Http Servlet Cross Origin Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state httpServletCrossOriginPolicyResourceModel
	readHttpServletCrossOriginPolicyResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *httpServletCrossOriginPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpServletCrossOriginPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateHttpServletCrossOriginPolicy(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultHttpServletCrossOriginPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan httpServletCrossOriginPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletCrossOriginPolicyApi.GetHttpServletCrossOriginPolicy(
		ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Cross Origin Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state httpServletCrossOriginPolicyResourceModel
	readHttpServletCrossOriginPolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletCrossOriginPolicyApi.UpdateHttpServletCrossOriginPolicy(ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createHttpServletCrossOriginPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletCrossOriginPolicyApi.UpdateHttpServletCrossOriginPolicyExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Servlet Cross Origin Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpServletCrossOriginPolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *httpServletCrossOriginPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpServletCrossOriginPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultHttpServletCrossOriginPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readHttpServletCrossOriginPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readHttpServletCrossOriginPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state httpServletCrossOriginPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.HttpServletCrossOriginPolicyApi.GetHttpServletCrossOriginPolicy(
		ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Http Servlet Cross Origin Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readHttpServletCrossOriginPolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *httpServletCrossOriginPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpServletCrossOriginPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultHttpServletCrossOriginPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateHttpServletCrossOriginPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateHttpServletCrossOriginPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan httpServletCrossOriginPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state httpServletCrossOriginPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.HttpServletCrossOriginPolicyApi.UpdateHttpServletCrossOriginPolicy(
		ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createHttpServletCrossOriginPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.HttpServletCrossOriginPolicyApi.UpdateHttpServletCrossOriginPolicyExecute(updateRequest)
		if err != nil {
			ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Http Servlet Cross Origin Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readHttpServletCrossOriginPolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultHttpServletCrossOriginPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *httpServletCrossOriginPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state httpServletCrossOriginPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.HttpServletCrossOriginPolicyApi.DeleteHttpServletCrossOriginPolicyExecute(r.apiClient.HttpServletCrossOriginPolicyApi.DeleteHttpServletCrossOriginPolicy(
		ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Http Servlet Cross Origin Policy", err, httpResp)
		return
	}
}

func (r *httpServletCrossOriginPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpServletCrossOriginPolicy(ctx, req, resp)
}

func (r *defaultHttpServletCrossOriginPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importHttpServletCrossOriginPolicy(ctx, req, resp)
}

func importHttpServletCrossOriginPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
