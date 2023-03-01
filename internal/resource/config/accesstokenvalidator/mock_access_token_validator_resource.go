package accesstokenvalidator

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
	_ resource.Resource                = &mockAccessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &mockAccessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &mockAccessTokenValidatorResource{}
)

// Create a Mock Access Token Validator resource
func NewMockAccessTokenValidatorResource() resource.Resource {
	return &mockAccessTokenValidatorResource{}
}

// mockAccessTokenValidatorResource is the resource implementation.
type mockAccessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *mockAccessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mock_access_token_validator"
}

// Configure adds the provider configured client to the resource.
func (r *mockAccessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type mockAccessTokenValidatorResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	ClientIDClaimName    types.String `tfsdk:"client_id_claim_name"`
	ScopeClaimName       types.String `tfsdk:"scope_claim_name"`
	EvaluationOrderIndex types.Int64  `tfsdk:"evaluation_order_index"`
	IdentityMapper       types.String `tfsdk:"identity_mapper"`
	SubjectClaimName     types.String `tfsdk:"subject_claim_name"`
	Description          types.String `tfsdk:"description"`
	Enabled              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *mockAccessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Mock Access Token Validator.",
		Attributes: map[string]schema.Attribute{
			"client_id_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the OAuth2 client ID.",
				Optional:    true,
				Computed:    true,
			},
			"scope_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the scopes granted by the token.",
				Optional:    true,
				Computed:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple Mock Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Mock Access Token Validators defined within Directory Server but not necessarily contiguous. Mock Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.",
				Optional:    true,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Optional:    true,
				Computed:    true,
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Access Token Validator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Access Token Validator is enabled for use in Directory Server.",
				Required:    true,
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalMockAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddMockAccessTokenValidatorRequest, plan mockAccessTokenValidatorResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientIDClaimName) {
		stringVal := plan.ClientIDClaimName.ValueString()
		addRequest.ClientIDClaimName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ScopeClaimName) {
		stringVal := plan.ScopeClaimName.ValueString()
		addRequest.ScopeClaimName = &stringVal
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		intVal := int32(plan.EvaluationOrderIndex.ValueInt64())
		addRequest.EvaluationOrderIndex = &intVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		stringVal := plan.IdentityMapper.ValueString()
		addRequest.IdentityMapper = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		stringVal := plan.SubjectClaimName.ValueString()
		addRequest.SubjectClaimName = &stringVal
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		stringVal := plan.Description.ValueString()
		addRequest.Description = &stringVal
	}
}

// Read a MockAccessTokenValidatorResponse object into the model struct
func readMockAccessTokenValidatorResponse(ctx context.Context, r *client.MockAccessTokenValidatorResponse, state *mockAccessTokenValidatorResourceModel, expectedValues *mockAccessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ClientIDClaimName = internaltypes.StringTypeOrNil(r.ClientIDClaimName, internaltypes.IsEmptyString(expectedValues.ClientIDClaimName))
	state.ScopeClaimName = internaltypes.StringTypeOrNil(r.ScopeClaimName, internaltypes.IsEmptyString(expectedValues.ScopeClaimName))
	state.EvaluationOrderIndex = types.Int64Value(int64(r.EvaluationOrderIndex))
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createMockAccessTokenValidatorOperations(plan mockAccessTokenValidatorResourceModel, state mockAccessTokenValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ClientIDClaimName, state.ClientIDClaimName, "client-id-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ScopeClaimName, state.ScopeClaimName, "scope-claim-name")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectClaimName, state.SubjectClaimName, "subject-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *mockAccessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan mockAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddMockAccessTokenValidatorRequest(plan.Id.ValueString(),
		[]client.EnummockAccessTokenValidatorSchemaUrn{client.ENUMMOCKACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORMOCK},
		plan.Enabled.ValueBool())
	addOptionalMockAccessTokenValidatorFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddMockAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Mock Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state mockAccessTokenValidatorResourceModel
	readMockAccessTokenValidatorResponse(ctx, addResponse.MockAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)

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
func (r *mockAccessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state mockAccessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Mock Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readMockAccessTokenValidatorResponse(ctx, readResponse.MockAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *mockAccessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan mockAccessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state mockAccessTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createMockAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Mock Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readMockAccessTokenValidatorResponse(ctx, updateResponse.MockAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *mockAccessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state mockAccessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidatorExecute(r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Mock Access Token Validator", err, httpResp)
		return
	}
}

func (r *mockAccessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
