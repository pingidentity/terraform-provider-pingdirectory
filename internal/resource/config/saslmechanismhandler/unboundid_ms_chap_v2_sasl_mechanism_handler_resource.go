package saslmechanismhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &unboundidMsChapV2SaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &unboundidMsChapV2SaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &unboundidMsChapV2SaslMechanismHandlerResource{}
	_ resource.Resource                = &defaultUnboundidMsChapV2SaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultUnboundidMsChapV2SaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &defaultUnboundidMsChapV2SaslMechanismHandlerResource{}
)

// Create a Unboundid Ms Chap V2 Sasl Mechanism Handler resource
func NewUnboundidMsChapV2SaslMechanismHandlerResource() resource.Resource {
	return &unboundidMsChapV2SaslMechanismHandlerResource{}
}

func NewDefaultUnboundidMsChapV2SaslMechanismHandlerResource() resource.Resource {
	return &defaultUnboundidMsChapV2SaslMechanismHandlerResource{}
}

// unboundidMsChapV2SaslMechanismHandlerResource is the resource implementation.
type unboundidMsChapV2SaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultUnboundidMsChapV2SaslMechanismHandlerResource is the resource implementation.
type defaultUnboundidMsChapV2SaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *unboundidMsChapV2SaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unboundid_ms_chap_v2_sasl_mechanism_handler"
}

func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_unboundid_ms_chap_v2_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the resource.
func (r *unboundidMsChapV2SaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type unboundidMsChapV2SaslMechanismHandlerResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	IdentityMapper  types.String `tfsdk:"identity_mapper"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *unboundidMsChapV2SaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	unboundidMsChapV2SaslMechanismHandlerSchema(ctx, req, resp, false)
}

func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	unboundidMsChapV2SaslMechanismHandlerSchema(ctx, req, resp, true)
}

func unboundidMsChapV2SaslMechanismHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Unboundid Ms Chap V2 Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that should be used to identify the entry associated with the username provided in the bind request.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this SASL Mechanism Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SASL mechanism handler is enabled for use.",
				Required:    true,
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
func addOptionalUnboundidMsChapV2SaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddUnboundidMsChapV2SaslMechanismHandlerRequest, plan unboundidMsChapV2SaslMechanismHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a UnboundidMsChapV2SaslMechanismHandlerResponse object into the model struct
func readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx context.Context, r *client.UnboundidMsChapV2SaslMechanismHandlerResponse, state *unboundidMsChapV2SaslMechanismHandlerResourceModel, expectedValues *unboundidMsChapV2SaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createUnboundidMsChapV2SaslMechanismHandlerOperations(plan unboundidMsChapV2SaslMechanismHandlerResourceModel, state unboundidMsChapV2SaslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *unboundidMsChapV2SaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan unboundidMsChapV2SaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddUnboundidMsChapV2SaslMechanismHandlerRequest(plan.Id.ValueString(),
		[]client.EnumunboundidMsChapV2SaslMechanismHandlerSchemaUrn{client.ENUMUNBOUNDIDMSCHAPV2SASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLERUNBOUNDID_MS_CHAP_V2},
		plan.IdentityMapper.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalUnboundidMsChapV2SaslMechanismHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerApi.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddUnboundidMsChapV2SaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Unboundid Ms Chap V2 Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state unboundidMsChapV2SaslMechanismHandlerResourceModel
	readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, addResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan unboundidMsChapV2SaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unboundid Ms Chap V2 Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state unboundidMsChapV2SaslMechanismHandlerResourceModel
	readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createUnboundidMsChapV2SaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unboundid Ms Chap V2 Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *unboundidMsChapV2SaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readUnboundidMsChapV2SaslMechanismHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state unboundidMsChapV2SaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unboundid Ms Chap V2 Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, readResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *unboundidMsChapV2SaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateUnboundidMsChapV2SaslMechanismHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan unboundidMsChapV2SaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state unboundidMsChapV2SaslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createUnboundidMsChapV2SaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unboundid Ms Chap V2 Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUnboundidMsChapV2SaslMechanismHandlerResponse(ctx, updateResponse.UnboundidMsChapV2SaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *unboundidMsChapV2SaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state unboundidMsChapV2SaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SaslMechanismHandlerApi.DeleteSaslMechanismHandlerExecute(r.apiClient.SaslMechanismHandlerApi.DeleteSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Unboundid Ms Chap V2 Sasl Mechanism Handler", err, httpResp)
		return
	}
}

func (r *unboundidMsChapV2SaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp)
}

func (r *defaultUnboundidMsChapV2SaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUnboundidMsChapV2SaslMechanismHandler(ctx, req, resp)
}

func importUnboundidMsChapV2SaslMechanismHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
