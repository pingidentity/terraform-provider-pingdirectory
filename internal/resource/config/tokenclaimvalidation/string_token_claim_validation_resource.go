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
	_ resource.Resource                = &stringTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &stringTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &stringTokenClaimValidationResource{}
	_ resource.Resource                = &defaultStringTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &defaultStringTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &defaultStringTokenClaimValidationResource{}
)

// Create a String Token Claim Validation resource
func NewStringTokenClaimValidationResource() resource.Resource {
	return &stringTokenClaimValidationResource{}
}

func NewDefaultStringTokenClaimValidationResource() resource.Resource {
	return &defaultStringTokenClaimValidationResource{}
}

// stringTokenClaimValidationResource is the resource implementation.
type stringTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultStringTokenClaimValidationResource is the resource implementation.
type defaultStringTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *stringTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_string_token_claim_validation"
}

func (r *defaultStringTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_string_token_claim_validation"
}

// Configure adds the provider configured client to the resource.
func (r *stringTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultStringTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type stringTokenClaimValidationResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	IdTokenValidatorName types.String `tfsdk:"id_token_validator_name"`
	AnyRequiredValue     types.Set    `tfsdk:"any_required_value"`
	Description          types.String `tfsdk:"description"`
	ClaimName            types.String `tfsdk:"claim_name"`
}

// GetSchema defines the schema for the resource.
func (r *stringTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	stringTokenClaimValidationSchema(ctx, req, resp, false)
}

func (r *defaultStringTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	stringTokenClaimValidationSchema(ctx, req, resp, true)
}

func stringTokenClaimValidationSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a String Token Claim Validation.",
		Attributes: map[string]schema.Attribute{
			"id_token_validator_name": schema.StringAttribute{
				Description: "Name of the parent ID Token Validator",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"any_required_value": schema.SetAttribute{
				Description: "The set of values that the claim may have to be considered valid.",
				Required:    true,
				ElementType: types.StringType,
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
func addOptionalStringTokenClaimValidationFields(ctx context.Context, addRequest *client.AddStringTokenClaimValidationRequest, plan stringTokenClaimValidationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a StringTokenClaimValidationResponse object into the model struct
func readStringTokenClaimValidationResponse(ctx context.Context, r *client.StringTokenClaimValidationResponse, state *stringTokenClaimValidationResourceModel, expectedValues *stringTokenClaimValidationResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.IdTokenValidatorName = expectedValues.IdTokenValidatorName
	state.AnyRequiredValue = internaltypes.GetStringSet(r.AnyRequiredValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.ClaimName = types.StringValue(r.ClaimName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createStringTokenClaimValidationOperations(plan stringTokenClaimValidationResourceModel, state stringTokenClaimValidationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyRequiredValue, state.AnyRequiredValue, "any-required-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.ClaimName, state.ClaimName, "claim-name")
	return ops
}

// Create a new resource
func (r *stringTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stringTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var AnyRequiredValueSlice []string
	plan.AnyRequiredValue.ElementsAs(ctx, &AnyRequiredValueSlice, false)
	addRequest := client.NewAddStringTokenClaimValidationRequest(plan.Id.ValueString(),
		[]client.EnumstringTokenClaimValidationSchemaUrn{client.ENUMSTRINGTOKENCLAIMVALIDATIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TOKEN_CLAIM_VALIDATIONSTRING},
		AnyRequiredValueSlice,
		plan.ClaimName.ValueString())
	addOptionalStringTokenClaimValidationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TokenClaimValidationApi.AddTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IdTokenValidatorName.ValueString())
	apiAddRequest = apiAddRequest.AddTokenClaimValidationRequest(
		client.AddStringTokenClaimValidationRequestAsAddTokenClaimValidationRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.AddTokenClaimValidationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the String Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state stringTokenClaimValidationResourceModel
	readStringTokenClaimValidationResponse(ctx, addResponse.StringTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultStringTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stringTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the String Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state stringTokenClaimValidationResourceModel
	readStringTokenClaimValidationResponse(ctx, readResponse.StringTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TokenClaimValidationApi.UpdateTokenClaimValidation(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString())
	ops := createStringTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the String Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStringTokenClaimValidationResponse(ctx, updateResponse.StringTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
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
func (r *stringTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readStringTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultStringTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readStringTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readStringTokenClaimValidation(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state stringTokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TokenClaimValidationApi.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString(), state.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the String Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readStringTokenClaimValidationResponse(ctx, readResponse.StringTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *stringTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateStringTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultStringTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateStringTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateStringTokenClaimValidation(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan stringTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state stringTokenClaimValidationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TokenClaimValidationApi.UpdateTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString())

	// Determine what update operations are necessary
	ops := createStringTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TokenClaimValidationApi.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the String Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStringTokenClaimValidationResponse(ctx, updateResponse.StringTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultStringTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *stringTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state stringTokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TokenClaimValidationApi.DeleteTokenClaimValidationExecute(r.apiClient.TokenClaimValidationApi.DeleteTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.IdTokenValidatorName.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the String Token Claim Validation", err, httpResp)
		return
	}
}

func (r *stringTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStringTokenClaimValidation(ctx, req, resp)
}

func (r *defaultStringTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStringTokenClaimValidation(ctx, req, resp)
}

func importStringTokenClaimValidation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [id-token-validator-name]/[token-claim-validation-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id_token_validator_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}
