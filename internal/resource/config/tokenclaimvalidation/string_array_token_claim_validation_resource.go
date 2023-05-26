package tokenclaimvalidation

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
	_ resource.Resource                = &stringArrayTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &stringArrayTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &stringArrayTokenClaimValidationResource{}
	_ resource.Resource                = &defaultStringArrayTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &defaultStringArrayTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &defaultStringArrayTokenClaimValidationResource{}
)

// Create a String Array Token Claim Validation resource
func NewStringArrayTokenClaimValidationResource() resource.Resource {
	return &stringArrayTokenClaimValidationResource{}
}

func NewDefaultStringArrayTokenClaimValidationResource() resource.Resource {
	return &defaultStringArrayTokenClaimValidationResource{}
}

// stringArrayTokenClaimValidationResource is the resource implementation.
type stringArrayTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultStringArrayTokenClaimValidationResource is the resource implementation.
type defaultStringArrayTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *stringArrayTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_string_array_token_claim_validation"
}

func (r *defaultStringArrayTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_string_array_token_claim_validation"
}

// Configure adds the provider configured client to the resource.
func (r *stringArrayTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultStringArrayTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type stringArrayTokenClaimValidationResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	LastUpdated          types.String `tfsdk:"last_updated"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	IdTokenValidatorName types.String `tfsdk:"id_token_validator_name"`
	AllRequiredValue     types.Set    `tfsdk:"all_required_value"`
	AnyRequiredValue     types.Set    `tfsdk:"any_required_value"`
	Description          types.String `tfsdk:"description"`
	ClaimName            types.String `tfsdk:"claim_name"`
}

// GetSchema defines the schema for the resource.
func (r *stringArrayTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	stringArrayTokenClaimValidationSchema(ctx, req, resp, false)
}

func (r *defaultStringArrayTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	stringArrayTokenClaimValidationSchema(ctx, req, resp, true)
}

func stringArrayTokenClaimValidationSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a String Array Token Claim Validation.",
		Attributes: map[string]schema.Attribute{
			"id_token_validator_name": schema.StringAttribute{
				Description: "Name of the parent ID Token Validator",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"all_required_value": schema.SetAttribute{
				Description: "The set of all values that the claim must have to be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"any_required_value": schema.SetAttribute{
				Description: "The set of values that the claim may have to be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
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

// Add config validators
func (r stringArrayTokenClaimValidationResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("any_required_value"),
			path.MatchRoot("all_required_value"),
		),
	}
}

// Add optional fields to create request
func addOptionalStringArrayTokenClaimValidationFields(ctx context.Context, addRequest *client.AddStringArrayTokenClaimValidationRequest, plan stringArrayTokenClaimValidationResourceModel) {
	if internaltypes.IsDefined(plan.AllRequiredValue) {
		var slice []string
		plan.AllRequiredValue.ElementsAs(ctx, &slice, false)
		addRequest.AllRequiredValue = slice
	}
	if internaltypes.IsDefined(plan.AnyRequiredValue) {
		var slice []string
		plan.AnyRequiredValue.ElementsAs(ctx, &slice, false)
		addRequest.AnyRequiredValue = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a StringArrayTokenClaimValidationResponse object into the model struct
func readStringArrayTokenClaimValidationResponse(ctx context.Context, r *client.StringArrayTokenClaimValidationResponse, state *stringArrayTokenClaimValidationResourceModel, expectedValues *stringArrayTokenClaimValidationResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.IdTokenValidatorName = expectedValues.IdTokenValidatorName
	state.AllRequiredValue = internaltypes.GetStringSet(r.AllRequiredValue)
	state.AnyRequiredValue = internaltypes.GetStringSet(r.AnyRequiredValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.ClaimName = types.StringValue(r.ClaimName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createStringArrayTokenClaimValidationOperations(plan stringArrayTokenClaimValidationResourceModel, state stringArrayTokenClaimValidationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllRequiredValue, state.AllRequiredValue, "all-required-value")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyRequiredValue, state.AnyRequiredValue, "any-required-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.ClaimName, state.ClaimName, "claim-name")
	return ops
}

// Create a new resource
func (r *stringArrayTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stringArrayTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddStringArrayTokenClaimValidationRequest(plan.Id.ValueString(),
		[]client.EnumstringArrayTokenClaimValidationSchemaUrn{client.ENUMSTRINGARRAYTOKENCLAIMVALIDATIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TOKEN_CLAIM_VALIDATIONSTRING_ARRAY},
		plan.ClaimName.ValueString())
	addOptionalStringArrayTokenClaimValidationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TokenClaimValidationApi.AddTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IdTokenValidatorName.ValueString())
	apiAddRequest = apiAddRequest.AddTokenClaimValidationRequest(
		client.AddStringArrayTokenClaimValidationRequestAsAddTokenClaimValidationRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.AddTokenClaimValidationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the String Array Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state stringArrayTokenClaimValidationResourceModel
	readStringArrayTokenClaimValidationResponse(ctx, addResponse.StringArrayTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultStringArrayTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan stringArrayTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the String Array Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state stringArrayTokenClaimValidationResourceModel
	readStringArrayTokenClaimValidationResponse(ctx, readResponse.StringArrayTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TokenClaimValidationApi.UpdateTokenClaimValidation(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString())
	ops := createStringArrayTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TokenClaimValidationApi.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the String Array Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStringArrayTokenClaimValidationResponse(ctx, updateResponse.StringArrayTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
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
func (r *stringArrayTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readStringArrayTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultStringArrayTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readStringArrayTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readStringArrayTokenClaimValidation(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state stringArrayTokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TokenClaimValidationApi.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString(), state.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the String Array Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readStringArrayTokenClaimValidationResponse(ctx, readResponse.StringArrayTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *stringArrayTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateStringArrayTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultStringArrayTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateStringArrayTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateStringArrayTokenClaimValidation(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan stringArrayTokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state stringArrayTokenClaimValidationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TokenClaimValidationApi.UpdateTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString(), plan.IdTokenValidatorName.ValueString())

	// Determine what update operations are necessary
	ops := createStringArrayTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TokenClaimValidationApi.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the String Array Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readStringArrayTokenClaimValidationResponse(ctx, updateResponse.StringArrayTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultStringArrayTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *stringArrayTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state stringArrayTokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TokenClaimValidationApi.DeleteTokenClaimValidationExecute(r.apiClient.TokenClaimValidationApi.DeleteTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.IdTokenValidatorName.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the String Array Token Claim Validation", err, httpResp)
		return
	}
}

func (r *stringArrayTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStringArrayTokenClaimValidation(ctx, req, resp)
}

func (r *defaultStringArrayTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importStringArrayTokenClaimValidation(ctx, req, resp)
}

func importStringArrayTokenClaimValidation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [id-token-validator-name]/[token-claim-validation-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id_token_validator_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), split[1])...)
}
