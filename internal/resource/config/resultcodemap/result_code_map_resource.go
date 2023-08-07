package resultcodemap

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &resultCodeMapResource{}
	_ resource.ResourceWithConfigure   = &resultCodeMapResource{}
	_ resource.ResourceWithImportState = &resultCodeMapResource{}
	_ resource.Resource                = &defaultResultCodeMapResource{}
	_ resource.ResourceWithConfigure   = &defaultResultCodeMapResource{}
	_ resource.ResourceWithImportState = &defaultResultCodeMapResource{}
)

// Create a Result Code Map resource
func NewResultCodeMapResource() resource.Resource {
	return &resultCodeMapResource{}
}

func NewDefaultResultCodeMapResource() resource.Resource {
	return &defaultResultCodeMapResource{}
}

// resultCodeMapResource is the resource implementation.
type resultCodeMapResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultResultCodeMapResource is the resource implementation.
type defaultResultCodeMapResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *resultCodeMapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_result_code_map"
}

func (r *defaultResultCodeMapResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_result_code_map"
}

// Configure adds the provider configured client to the resource.
func (r *resultCodeMapResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultResultCodeMapResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type resultCodeMapResourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	LastUpdated                   types.String `tfsdk:"last_updated"`
	Notifications                 types.Set    `tfsdk:"notifications"`
	RequiredActions               types.Set    `tfsdk:"required_actions"`
	Type                          types.String `tfsdk:"type"`
	Description                   types.String `tfsdk:"description"`
	BindAccountLockedResultCode   types.Int64  `tfsdk:"bind_account_locked_result_code"`
	BindMissingUserResultCode     types.Int64  `tfsdk:"bind_missing_user_result_code"`
	BindMissingPasswordResultCode types.Int64  `tfsdk:"bind_missing_password_result_code"`
	ServerErrorResultCode         types.Int64  `tfsdk:"server_error_result_code"`
}

// GetSchema defines the schema for the resource.
func (r *resultCodeMapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resultCodeMapSchema(ctx, req, resp, false)
}

func (r *defaultResultCodeMapResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resultCodeMapSchema(ctx, req, resp, true)
}

func resultCodeMapSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Result Code Map.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Result Code Map resource. Options are ['result-code-map']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("result-code-map"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"result-code-map"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Result Code Map",
				Optional:    true,
			},
			"bind_account_locked_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a bind attempt fails because the user's account is locked as a result of too many failed authentication attempts.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"bind_missing_user_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a bind attempt fails because the target user entry does not exist in the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"bind_missing_password_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a password-based bind attempt fails because the target user entry does not have a password.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"server_error_result_code": schema.Int64Attribute{
				Description: "Specifies the result code that should be returned if a generic error occurs within the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputed(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for result-code-map result-code-map
func addOptionalResultCodeMapFields(ctx context.Context, addRequest *client.AddResultCodeMapRequest, plan resultCodeMapResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.BindAccountLockedResultCode) {
		addRequest.BindAccountLockedResultCode = plan.BindAccountLockedResultCode.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.BindMissingUserResultCode) {
		addRequest.BindMissingUserResultCode = plan.BindMissingUserResultCode.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.BindMissingPasswordResultCode) {
		addRequest.BindMissingPasswordResultCode = plan.BindMissingPasswordResultCode.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ServerErrorResultCode) {
		addRequest.ServerErrorResultCode = plan.ServerErrorResultCode.ValueInt64Pointer()
	}
}

// Read a ResultCodeMapResponse object into the model struct
func readResultCodeMapResponse(ctx context.Context, r *client.ResultCodeMapResponse, state *resultCodeMapResourceModel, expectedValues *resultCodeMapResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("result-code-map")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.BindAccountLockedResultCode = internaltypes.Int64TypeOrNil(r.BindAccountLockedResultCode)
	state.BindMissingUserResultCode = internaltypes.Int64TypeOrNil(r.BindMissingUserResultCode)
	state.BindMissingPasswordResultCode = internaltypes.Int64TypeOrNil(r.BindMissingPasswordResultCode)
	state.ServerErrorResultCode = internaltypes.Int64TypeOrNil(r.ServerErrorResultCode)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createResultCodeMapOperations(plan resultCodeMapResourceModel, state resultCodeMapResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddInt64OperationIfNecessary(&ops, plan.BindAccountLockedResultCode, state.BindAccountLockedResultCode, "bind-account-locked-result-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.BindMissingUserResultCode, state.BindMissingUserResultCode, "bind-missing-user-result-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.BindMissingPasswordResultCode, state.BindMissingPasswordResultCode, "bind-missing-password-result-code")
	operations.AddInt64OperationIfNecessary(&ops, plan.ServerErrorResultCode, state.ServerErrorResultCode, "server-error-result-code")
	return ops
}

// Create a result-code-map result-code-map
func (r *resultCodeMapResource) CreateResultCodeMap(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan resultCodeMapResourceModel) (*resultCodeMapResourceModel, error) {
	addRequest := client.NewAddResultCodeMapRequest(plan.Name.ValueString())
	addOptionalResultCodeMapFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ResultCodeMapApi.AddResultCodeMap(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddResultCodeMapRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ResultCodeMapApi.AddResultCodeMapExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Result Code Map", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state resultCodeMapResourceModel
	readResultCodeMapResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *resultCodeMapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resultCodeMapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateResultCodeMap(ctx, req, resp, plan)
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
func (r *defaultResultCodeMapResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan resultCodeMapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ResultCodeMapApi.GetResultCodeMap(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Result Code Map", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state resultCodeMapResourceModel
	readResultCodeMapResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ResultCodeMapApi.UpdateResultCodeMap(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createResultCodeMapOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ResultCodeMapApi.UpdateResultCodeMapExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Result Code Map", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readResultCodeMapResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *resultCodeMapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readResultCodeMap(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultResultCodeMapResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readResultCodeMap(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readResultCodeMap(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state resultCodeMapResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ResultCodeMapApi.GetResultCodeMap(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Result Code Map", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readResultCodeMapResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *resultCodeMapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateResultCodeMap(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultResultCodeMapResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateResultCodeMap(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateResultCodeMap(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan resultCodeMapResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state resultCodeMapResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ResultCodeMapApi.UpdateResultCodeMap(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createResultCodeMapOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ResultCodeMapApi.UpdateResultCodeMapExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Result Code Map", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readResultCodeMapResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultResultCodeMapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *resultCodeMapResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state resultCodeMapResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ResultCodeMapApi.DeleteResultCodeMapExecute(r.apiClient.ResultCodeMapApi.DeleteResultCodeMap(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Result Code Map", err, httpResp)
		return
	}
}

func (r *resultCodeMapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importResultCodeMap(ctx, req, resp)
}

func (r *defaultResultCodeMapResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importResultCodeMap(ctx, req, resp)
}

func importResultCodeMap(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
