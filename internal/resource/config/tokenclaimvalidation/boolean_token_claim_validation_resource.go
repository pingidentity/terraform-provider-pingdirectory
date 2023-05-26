package tokenclaimvalidation

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	_ resource.Resource                = &booleanTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &booleanTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &booleanTokenClaimValidationResource{}
	_ resource.Resource                = &defaultBooleanTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &defaultBooleanTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &defaultBooleanTokenClaimValidationResource{}
)

// Create a Boolean Token Claim Validation resource
func NewBooleanTokenClaimValidationResource() resource.Resource {
	return &booleanTokenClaimValidationResource{}
}

func NewDefaultBooleanTokenClaimValidationResource() resource.Resource {
	return &defaultBooleanTokenClaimValidationResource{}
}

// booleanTokenClaimValidationResource is the resource implementation.
type booleanTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultBooleanTokenClaimValidationResource is the resource implementation.
type defaultBooleanTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *booleanTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_boolean_token_claim_validation"
}

func (r *defaultBooleanTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_boolean_token_claim_validation"
}

// Configure adds the provider configured client to the resource.
func (r *booleanTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultBooleanTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type booleanTokenClaimValidationResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	IdTokenValidatorName types.String `tfsdk:"id_token_validator_name"`
	RequiredValue        types.String `tfsdk:"required_value"`
	Description          types.String `tfsdk:"description"`
	ClaimName            types.String `tfsdk:"claim_name"`
}

// GetSchema defines the schema for the resource.
func (r *booleanTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	booleanTokenClaimValidationSchema(ctx, req, resp, false)
}

func (r *defaultBooleanTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	booleanTokenClaimValidationSchema(ctx, req, resp, true)
}

func booleanTokenClaimValidationSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Boolean Token Claim Validation.",
		Attributes: map[string]schema.Attribute{
			"id_token_validator_name": schema.StringAttribute{
				Description: "Name of the parent ID Token Validator",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"required_value": schema.StringAttribute{
				Description: "Specifies the boolean claim's required value.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Token Claim Validation",
				Optional:    true,
			},
			"claim_name": schema.StringAttribute{
				Description: "The name of the claim to be validated.",
				Required:    true,
			},
		},
	}
	if setOptionalToComputed {
		config.SetAllAttributesToOptionalAndComputed(&schema, []string{"id", "id_token_validator_name"})
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Add optional fields to create request
func addOptionalBooleanTokenClaimValidationFields(ctx context.Context, addRequest *client.AddBooleanTokenClaimValidationRequest, plan booleanTokenClaimValidationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a BooleanTokenClaimValidationResponse object into the model struct
func readBooleanTokenClaimValidationResponse(ctx context.Context, r *client.BooleanTokenClaimValidationResponse, state *booleanTokenClaimValidationResourceModel, expectedValues *booleanTokenClaimValidationResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.IdTokenValidatorName = expectedValues.IdTokenValidatorName
	state.RequiredValue = types.StringValue(r.RequiredValue.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.ClaimName = types.StringValue(r.ClaimName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createBooleanTokenClaimValidationOperations(plan booleanTokenClaimValidationResourceModel, state booleanTokenClaimValidationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.RequiredValue, state.RequiredValue, "required-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.ClaimName, state.ClaimName, "claim-name")
	return ops
}

// Create a new resource
func (r *booleanTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan booleanTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	requiredValue, err := client.NewEnumtokenClaimValidationRequiredValuePropFromValue(plan.RequiredValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for RequiredValue", err.Error())
		return
	}
	addRequest := client.NewAddBooleanTokenClaimValidationRequest(plan.Id.ValueString(),
		[]client.EnumbooleanTokenClaimValidationSchemaUrn{client.ENUMBOOLEANTOKENCLAIMVALIDATIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TOKEN_CLAIM_VALIDATIONBOOLEAN},
		*requiredValue,
		plan.ClaimName.ValueString())
	addOptionalBooleanTokenClaimValidationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TokenClaimValidationApi.AddTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IdTokenValidatorName.ValueString())
	apiAddRequest = apiAddRequest.AddTokenClaimValidationRequest(
		client.AddBooleanTokenClaimValidationRequestAsAddTokenClaimValidationRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.AddTokenClaimValidationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Boolean Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state booleanTokenClaimValidationResourceModel
	readBooleanTokenClaimValidationResponse(ctx, addResponse.BooleanTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultBooleanTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan booleanTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Boolean Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state booleanTokenClaimValidationResourceModel
	readBooleanTokenClaimValidationResponse(ctx, readResponse.BooleanTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TokenClaimValidationApi.UpdateTokenClaimValidation(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString())
	ops := createBooleanTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Boolean Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readBooleanTokenClaimValidationResponse(ctx, updateResponse.BooleanTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
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
func (r *booleanTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readBooleanTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultBooleanTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readBooleanTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readBooleanTokenClaimValidation(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state booleanTokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TokenClaimValidationApi.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString(), state.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Boolean Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readBooleanTokenClaimValidationResponse(ctx, readResponse.BooleanTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *booleanTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateBooleanTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultBooleanTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateBooleanTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateBooleanTokenClaimValidation(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan booleanTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state booleanTokenClaimValidationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TokenClaimValidationApi.UpdateTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString())

	// Determine what update operations are necessary
	ops := createBooleanTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TokenClaimValidationApi.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Boolean Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readBooleanTokenClaimValidationResponse(ctx, updateResponse.BooleanTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultBooleanTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *booleanTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state booleanTokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TokenClaimValidationApi.DeleteTokenClaimValidationExecute(r.apiClient.TokenClaimValidationApi.DeleteTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.IdTokenValidatorName.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Boolean Token Claim Validation", err, httpResp)
		return
	}
}

func (r *booleanTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importBooleanTokenClaimValidation(ctx, req, resp)
}

func (r *defaultBooleanTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importBooleanTokenClaimValidation(ctx, req, resp)
}

func importBooleanTokenClaimValidation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [id-token-validator-name]/[token-claim-validation-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id_token_validator_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}
