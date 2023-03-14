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
	_ resource.Resource                = &availabilityStateHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &availabilityStateHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &availabilityStateHttpServletExtensionResource{}
	_ resource.Resource                = &defaultAvailabilityStateHttpServletExtensionResource{}
	_ resource.ResourceWithConfigure   = &defaultAvailabilityStateHttpServletExtensionResource{}
	_ resource.ResourceWithImportState = &defaultAvailabilityStateHttpServletExtensionResource{}
)

// Create a Availability State Http Servlet Extension resource
func NewAvailabilityStateHttpServletExtensionResource() resource.Resource {
	return &availabilityStateHttpServletExtensionResource{}
}

func NewDefaultAvailabilityStateHttpServletExtensionResource() resource.Resource {
	return &defaultAvailabilityStateHttpServletExtensionResource{}
}

// availabilityStateHttpServletExtensionResource is the resource implementation.
type availabilityStateHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAvailabilityStateHttpServletExtensionResource is the resource implementation.
type defaultAvailabilityStateHttpServletExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *availabilityStateHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_availability_state_http_servlet_extension"
}

func (r *defaultAvailabilityStateHttpServletExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_availability_state_http_servlet_extension"
}

// Configure adds the provider configured client to the resource.
func (r *availabilityStateHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAvailabilityStateHttpServletExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type availabilityStateHttpServletExtensionResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	BaseContextPath             types.String `tfsdk:"base_context_path"`
	AvailableStatusCode         types.Int64  `tfsdk:"available_status_code"`
	DegradedStatusCode          types.Int64  `tfsdk:"degraded_status_code"`
	UnavailableStatusCode       types.Int64  `tfsdk:"unavailable_status_code"`
	OverrideStatusCode          types.Int64  `tfsdk:"override_status_code"`
	IncludeResponseBody         types.Bool   `tfsdk:"include_response_body"`
	AdditionalResponseContents  types.String `tfsdk:"additional_response_contents"`
	Description                 types.String `tfsdk:"description"`
	CrossOriginPolicy           types.String `tfsdk:"cross_origin_policy"`
	ResponseHeader              types.Set    `tfsdk:"response_header"`
	CorrelationIDResponseHeader types.String `tfsdk:"correlation_id_response_header"`
}

// GetSchema defines the schema for the resource.
func (r *availabilityStateHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	availabilityStateHttpServletExtensionSchema(ctx, req, resp, false)
}

func (r *defaultAvailabilityStateHttpServletExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	availabilityStateHttpServletExtensionSchema(ctx, req, resp, true)
}

func availabilityStateHttpServletExtensionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Availability State Http Servlet Extension.",
		Attributes: map[string]schema.Attribute{
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that HTTP clients should use to access this servlet. The value must start with a forward slash and must represent a valid HTTP context path.",
				Required:    true,
			},
			"available_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be available.",
				Required:    true,
			},
			"degraded_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be degraded.",
				Required:    true,
			},
			"unavailable_status_code": schema.Int64Attribute{
				Description: "Specifies the HTTP status code that the servlet should return if the server considers itself to be unavailable.",
				Required:    true,
			},
			"override_status_code": schema.Int64Attribute{
				Description: "Specifies a HTTP status code that the servlet should always return, regardless of the server's availability. If this value is defined, it will override the availability-based return codes.",
				Optional:    true,
			},
			"include_response_body": schema.BoolAttribute{
				Description: "Indicates whether the response should include a body that is a JSON object.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"additional_response_contents": schema.StringAttribute{
				Description: "A JSON-formatted string containing additional fields to be returned in the response body. For example, an additional-response-contents value of '{ \"key\": \"value\" }' would add the key and value to the root of the JSON response body.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this HTTP Servlet Extension",
				Optional:    true,
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
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalAvailabilityStateHttpServletExtensionFields(ctx context.Context, addRequest *client.AddAvailabilityStateHttpServletExtensionRequest, plan availabilityStateHttpServletExtensionResourceModel) {
	if internaltypes.IsDefined(plan.OverrideStatusCode) {
		intVal := int32(plan.OverrideStatusCode.ValueInt64())
		addRequest.OverrideStatusCode = &intVal
	}
	if internaltypes.IsDefined(plan.IncludeResponseBody) {
		boolVal := plan.IncludeResponseBody.ValueBool()
		addRequest.IncludeResponseBody = &boolVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AdditionalResponseContents) {
		stringVal := plan.AdditionalResponseContents.ValueString()
		addRequest.AdditionalResponseContents = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CrossOriginPolicy) {
		stringVal := plan.CrossOriginPolicy.ValueString()
		addRequest.CrossOriginPolicy = &stringVal
	}
	if internaltypes.IsDefined(plan.ResponseHeader) {
		var slice []string
		plan.ResponseHeader.ElementsAs(ctx, &slice, false)
		addRequest.ResponseHeader = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.CorrelationIDResponseHeader) {
		stringVal := plan.CorrelationIDResponseHeader.ValueString()
		addRequest.CorrelationIDResponseHeader = &stringVal
	}
}

// Read a AvailabilityStateHttpServletExtensionResponse object into the model struct
func readAvailabilityStateHttpServletExtensionResponse(ctx context.Context, r *client.AvailabilityStateHttpServletExtensionResponse, state *availabilityStateHttpServletExtensionResourceModel, expectedValues *availabilityStateHttpServletExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.AvailableStatusCode = types.Int64Value(int64(r.AvailableStatusCode))
	state.DegradedStatusCode = types.Int64Value(int64(r.DegradedStatusCode))
	state.UnavailableStatusCode = types.Int64Value(int64(r.UnavailableStatusCode))
	state.OverrideStatusCode = internaltypes.Int64TypeOrNil(r.OverrideStatusCode)
	state.IncludeResponseBody = internaltypes.BoolTypeOrNil(r.IncludeResponseBody)
	state.AdditionalResponseContents = internaltypes.StringTypeOrNil(r.AdditionalResponseContents, internaltypes.IsEmptyString(expectedValues.AdditionalResponseContents))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.CrossOriginPolicy = internaltypes.StringTypeOrNil(r.CrossOriginPolicy, internaltypes.IsEmptyString(expectedValues.CrossOriginPolicy))
	state.ResponseHeader = internaltypes.GetStringSet(r.ResponseHeader)
	state.CorrelationIDResponseHeader = internaltypes.StringTypeOrNil(r.CorrelationIDResponseHeader, internaltypes.IsEmptyString(expectedValues.CorrelationIDResponseHeader))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAvailabilityStateHttpServletExtensionOperations(plan availabilityStateHttpServletExtensionResourceModel, state availabilityStateHttpServletExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddInt64OperationIfNecessary(&ops, plan.AvailableStatusCode, state.AvailableStatusCode, "available-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.DegradedStatusCode, state.DegradedStatusCode, "degraded-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.UnavailableStatusCode, state.UnavailableStatusCode, "unavailable-status-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.OverrideStatusCode, state.OverrideStatusCode, "override-status-code")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeResponseBody, state.IncludeResponseBody, "include-response-body")
	operations.AddStringOperationIfNecessary(&ops, plan.AdditionalResponseContents, state.AdditionalResponseContents, "additional-response-contents")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.CrossOriginPolicy, state.CrossOriginPolicy, "cross-origin-policy")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ResponseHeader, state.ResponseHeader, "response-header")
	operations.AddStringOperationIfNecessary(&ops, plan.CorrelationIDResponseHeader, state.CorrelationIDResponseHeader, "correlation-id-response-header")
	return ops
}

// Create a new resource
func (r *availabilityStateHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan availabilityStateHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAvailabilityStateHttpServletExtensionRequest(plan.Id.ValueString(),
		[]client.EnumavailabilityStateHttpServletExtensionSchemaUrn{client.ENUMAVAILABILITYSTATEHTTPSERVLETEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0HTTP_SERVLET_EXTENSIONAVAILABILITY_STATE},
		plan.BaseContextPath.ValueString(),
		int32(plan.AvailableStatusCode.ValueInt64()),
		int32(plan.DegradedStatusCode.ValueInt64()),
		int32(plan.UnavailableStatusCode.ValueInt64()))
	addOptionalAvailabilityStateHttpServletExtensionFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.HttpServletExtensionApi.AddHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddHttpServletExtensionRequest(
		client.AddAvailabilityStateHttpServletExtensionRequestAsAddHttpServletExtensionRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.AddHttpServletExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Availability State Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state availabilityStateHttpServletExtensionResourceModel
	readAvailabilityStateHttpServletExtensionResponse(ctx, addResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAvailabilityStateHttpServletExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan availabilityStateHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Availability State Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state availabilityStateHttpServletExtensionResourceModel
	readAvailabilityStateHttpServletExtensionResponse(ctx, readResponse.AvailabilityStateHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAvailabilityStateHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Availability State Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAvailabilityStateHttpServletExtensionResponse(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *availabilityStateHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAvailabilityStateHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAvailabilityStateHttpServletExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAvailabilityStateHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAvailabilityStateHttpServletExtension(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state availabilityStateHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.HttpServletExtensionApi.GetHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Availability State Http Servlet Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAvailabilityStateHttpServletExtensionResponse(ctx, readResponse.AvailabilityStateHttpServletExtensionResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *availabilityStateHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAvailabilityStateHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAvailabilityStateHttpServletExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAvailabilityStateHttpServletExtension(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAvailabilityStateHttpServletExtension(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan availabilityStateHttpServletExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state availabilityStateHttpServletExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.HttpServletExtensionApi.UpdateHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAvailabilityStateHttpServletExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.HttpServletExtensionApi.UpdateHttpServletExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Availability State Http Servlet Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAvailabilityStateHttpServletExtensionResponse(ctx, updateResponse.AvailabilityStateHttpServletExtensionResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAvailabilityStateHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *availabilityStateHttpServletExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state availabilityStateHttpServletExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtensionExecute(r.apiClient.HttpServletExtensionApi.DeleteHttpServletExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Availability State Http Servlet Extension", err, httpResp)
		return
	}
}

func (r *availabilityStateHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAvailabilityStateHttpServletExtension(ctx, req, resp)
}

func (r *defaultAvailabilityStateHttpServletExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAvailabilityStateHttpServletExtension(ctx, req, resp)
}

func importAvailabilityStateHttpServletExtension(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
