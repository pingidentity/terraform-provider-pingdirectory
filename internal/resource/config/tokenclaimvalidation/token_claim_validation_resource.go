// Copyright © 2025 Ping Identity Corporation

package tokenclaimvalidation

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &tokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &tokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &tokenClaimValidationResource{}
	_ resource.Resource                = &defaultTokenClaimValidationResource{}
	_ resource.ResourceWithConfigure   = &defaultTokenClaimValidationResource{}
	_ resource.ResourceWithImportState = &defaultTokenClaimValidationResource{}
)

// Create a Token Claim Validation resource
func NewTokenClaimValidationResource() resource.Resource {
	return &tokenClaimValidationResource{}
}

func NewDefaultTokenClaimValidationResource() resource.Resource {
	return &defaultTokenClaimValidationResource{}
}

// tokenClaimValidationResource is the resource implementation.
type tokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultTokenClaimValidationResource is the resource implementation.
type defaultTokenClaimValidationResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *tokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_token_claim_validation"
}

func (r *defaultTokenClaimValidationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_token_claim_validation"
}

// Configure adds the provider configured client to the resource.
func (r *tokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultTokenClaimValidationResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type tokenClaimValidationResourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Notifications        types.Set    `tfsdk:"notifications"`
	RequiredActions      types.Set    `tfsdk:"required_actions"`
	Type                 types.String `tfsdk:"type"`
	IdTokenValidatorName types.String `tfsdk:"id_token_validator_name"`
	RequiredValue        types.String `tfsdk:"required_value"`
	AllRequiredValue     types.Set    `tfsdk:"all_required_value"`
	AnyRequiredValue     types.Set    `tfsdk:"any_required_value"`
	Description          types.String `tfsdk:"description"`
	ClaimName            types.String `tfsdk:"claim_name"`
}

// GetSchema defines the schema for the resource.
func (r *tokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tokenClaimValidationSchema(ctx, req, resp, false)
}

func (r *defaultTokenClaimValidationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	tokenClaimValidationSchema(ctx, req, resp, true)
}

func tokenClaimValidationSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Token Claim Validation.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Token Claim Validation resource. Options are ['string-array', 'boolean', 'string']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"string-array", "boolean", "string"}...),
				},
			},
			"id_token_validator_name": schema.StringAttribute{
				Description: "Name of the parent ID Token Validator",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"required_value": schema.StringAttribute{
				Description: "Specifies the boolean claim's required value.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"true", "false"}...),
				},
			},
			"all_required_value": schema.SetAttribute{
				Description: "The set of all values that the claim must have to be considered valid.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"any_required_value": schema.SetAttribute{
				Description: "The set of values that the claim may have to be considered valid.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
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
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type", "id_token_validator_name"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsTokenClaimValidation() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"string-array"},
			resourcevalidator.AtLeastOneOf(
				path.MatchRoot("all_required_value"),
				path.MatchRoot("any_required_value"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("all_required_value"),
			path.MatchRoot("type"),
			[]string{"string-array"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("any_required_value"),
			path.MatchRoot("type"),
			[]string{"string-array", "string"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("required_value"),
			path.MatchRoot("type"),
			[]string{"boolean"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"string",
			[]path.Expression{path.MatchRoot("any_required_value")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"boolean",
			[]path.Expression{path.MatchRoot("required_value")},
		),
	}
}

// Add config validators
func (r tokenClaimValidationResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsTokenClaimValidation()
}

// Add config validators
func (r defaultTokenClaimValidationResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsTokenClaimValidation()
}

// Add optional fields to create request for string-array token-claim-validation
func addOptionalStringArrayTokenClaimValidationFields(ctx context.Context, addRequest *client.AddStringArrayTokenClaimValidationRequest, plan tokenClaimValidationResourceModel) {
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

// Add optional fields to create request for boolean token-claim-validation
func addOptionalBooleanTokenClaimValidationFields(ctx context.Context, addRequest *client.AddBooleanTokenClaimValidationRequest, plan tokenClaimValidationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for string token-claim-validation
func addOptionalStringTokenClaimValidationFields(ctx context.Context, addRequest *client.AddStringTokenClaimValidationRequest, plan tokenClaimValidationResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateTokenClaimValidationUnknownValues(model *tokenClaimValidationResourceModel) {
	if model.AnyRequiredValue.IsUnknown() || model.AnyRequiredValue.IsNull() {
		model.AnyRequiredValue, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.AllRequiredValue.IsUnknown() || model.AllRequiredValue.IsNull() {
		model.AllRequiredValue, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *tokenClaimValidationResourceModel) populateAllComputedStringAttributes() {
	if model.ClaimName.IsUnknown() || model.ClaimName.IsNull() {
		model.ClaimName = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.RequiredValue.IsUnknown() || model.RequiredValue.IsNull() {
		model.RequiredValue = types.StringValue("")
	}
}

// Read a StringArrayTokenClaimValidationResponse object into the model struct
func readStringArrayTokenClaimValidationResponse(ctx context.Context, r *client.StringArrayTokenClaimValidationResponse, state *tokenClaimValidationResourceModel, expectedValues *tokenClaimValidationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("string-array")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllRequiredValue = internaltypes.GetStringSet(r.AllRequiredValue)
	state.AnyRequiredValue = internaltypes.GetStringSet(r.AnyRequiredValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.ClaimName = types.StringValue(r.ClaimName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTokenClaimValidationUnknownValues(state)
}

// Read a BooleanTokenClaimValidationResponse object into the model struct
func readBooleanTokenClaimValidationResponse(ctx context.Context, r *client.BooleanTokenClaimValidationResponse, state *tokenClaimValidationResourceModel, expectedValues *tokenClaimValidationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("boolean")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.RequiredValue = types.StringValue(r.RequiredValue.String())
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.ClaimName = types.StringValue(r.ClaimName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTokenClaimValidationUnknownValues(state)
}

// Read a StringTokenClaimValidationResponse object into the model struct
func readStringTokenClaimValidationResponse(ctx context.Context, r *client.StringTokenClaimValidationResponse, state *tokenClaimValidationResourceModel, expectedValues *tokenClaimValidationResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("string")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AnyRequiredValue = internaltypes.GetStringSet(r.AnyRequiredValue)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.ClaimName = types.StringValue(r.ClaimName)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateTokenClaimValidationUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *tokenClaimValidationResourceModel) setStateValuesNotReturnedByAPI(expectedValues *tokenClaimValidationResourceModel) {
	if !expectedValues.IdTokenValidatorName.IsUnknown() {
		state.IdTokenValidatorName = expectedValues.IdTokenValidatorName
	}
}

// Create any update operations necessary to make the state match the plan
func createTokenClaimValidationOperations(plan tokenClaimValidationResourceModel, state tokenClaimValidationResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.RequiredValue, state.RequiredValue, "required-value")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllRequiredValue, state.AllRequiredValue, "all-required-value")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AnyRequiredValue, state.AnyRequiredValue, "any-required-value")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.ClaimName, state.ClaimName, "claim-name")
	return ops
}

// Create a string-array token-claim-validation
func (r *tokenClaimValidationResource) CreateStringArrayTokenClaimValidation(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan tokenClaimValidationResourceModel) (*tokenClaimValidationResourceModel, error) {
	addRequest := client.NewAddStringArrayTokenClaimValidationRequest([]client.EnumstringArrayTokenClaimValidationSchemaUrn{client.ENUMSTRINGARRAYTOKENCLAIMVALIDATIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TOKEN_CLAIM_VALIDATIONSTRING_ARRAY},
		plan.ClaimName.ValueString(),
		plan.Name.ValueString())
	addOptionalStringArrayTokenClaimValidationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TokenClaimValidationAPI.AddTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IdTokenValidatorName.ValueString())
	apiAddRequest = apiAddRequest.AddTokenClaimValidationRequest(
		client.AddStringArrayTokenClaimValidationRequestAsAddTokenClaimValidationRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.AddTokenClaimValidationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Token Claim Validation", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state tokenClaimValidationResourceModel
	readStringArrayTokenClaimValidationResponse(ctx, addResponse.StringArrayTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a boolean token-claim-validation
func (r *tokenClaimValidationResource) CreateBooleanTokenClaimValidation(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan tokenClaimValidationResourceModel) (*tokenClaimValidationResourceModel, error) {
	requiredValue, err := client.NewEnumtokenClaimValidationRequiredValuePropFromValue(plan.RequiredValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse enum value for RequiredValue", err.Error())
		return nil, err
	}
	addRequest := client.NewAddBooleanTokenClaimValidationRequest([]client.EnumbooleanTokenClaimValidationSchemaUrn{client.ENUMBOOLEANTOKENCLAIMVALIDATIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TOKEN_CLAIM_VALIDATIONBOOLEAN},
		*requiredValue,
		plan.ClaimName.ValueString(),
		plan.Name.ValueString())
	addOptionalBooleanTokenClaimValidationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TokenClaimValidationAPI.AddTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IdTokenValidatorName.ValueString())
	apiAddRequest = apiAddRequest.AddTokenClaimValidationRequest(
		client.AddBooleanTokenClaimValidationRequestAsAddTokenClaimValidationRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.AddTokenClaimValidationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Token Claim Validation", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state tokenClaimValidationResourceModel
	readBooleanTokenClaimValidationResponse(ctx, addResponse.BooleanTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a string token-claim-validation
func (r *tokenClaimValidationResource) CreateStringTokenClaimValidation(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan tokenClaimValidationResourceModel) (*tokenClaimValidationResourceModel, error) {
	var AnyRequiredValueSlice []string
	plan.AnyRequiredValue.ElementsAs(ctx, &AnyRequiredValueSlice, false)
	addRequest := client.NewAddStringTokenClaimValidationRequest([]client.EnumstringTokenClaimValidationSchemaUrn{client.ENUMSTRINGTOKENCLAIMVALIDATIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0TOKEN_CLAIM_VALIDATIONSTRING},
		AnyRequiredValueSlice,
		plan.ClaimName.ValueString(),
		plan.Name.ValueString())
	addOptionalStringTokenClaimValidationFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.TokenClaimValidationAPI.AddTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.IdTokenValidatorName.ValueString())
	apiAddRequest = apiAddRequest.AddTokenClaimValidationRequest(
		client.AddStringTokenClaimValidationRequestAsAddTokenClaimValidationRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.AddTokenClaimValidationExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Token Claim Validation", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state tokenClaimValidationResourceModel
	readStringTokenClaimValidationResponse(ctx, addResponse.StringTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *tokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan tokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *tokenClaimValidationResourceModel
	var err error
	if plan.Type.ValueString() == "string-array" {
		state, err = r.CreateStringArrayTokenClaimValidation(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "boolean" {
		state, err = r.CreateBooleanTokenClaimValidation(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "string" {
		state, err = r.CreateStringTokenClaimValidation(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.setStateValuesNotReturnedByAPI(&plan)
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
func (r *defaultTokenClaimValidationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan tokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Token Claim Validation", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state tokenClaimValidationResourceModel
	if readResponse.StringArrayTokenClaimValidationResponse != nil {
		readStringArrayTokenClaimValidationResponse(ctx, readResponse.StringArrayTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BooleanTokenClaimValidationResponse != nil {
		readBooleanTokenClaimValidationResponse(ctx, readResponse.BooleanTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StringTokenClaimValidationResponse != nil {
		readStringTokenClaimValidationResponse(ctx, readResponse.StringTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.TokenClaimValidationAPI.UpdateTokenClaimValidation(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString(), plan.IdTokenValidatorName.ValueString())
	ops := createTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.TokenClaimValidationAPI.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.StringArrayTokenClaimValidationResponse != nil {
			readStringArrayTokenClaimValidationResponse(ctx, updateResponse.StringArrayTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BooleanTokenClaimValidationResponse != nil {
			readBooleanTokenClaimValidationResponse(ctx, updateResponse.BooleanTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StringTokenClaimValidationResponse != nil {
			readStringTokenClaimValidationResponse(ctx, updateResponse.StringTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *tokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultTokenClaimValidationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readTokenClaimValidation(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state tokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.TokenClaimValidationAPI.GetTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString(), state.IdTokenValidatorName.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Token Claim Validation", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Token Claim Validation", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.StringArrayTokenClaimValidationResponse != nil {
		readStringArrayTokenClaimValidationResponse(ctx, readResponse.StringArrayTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BooleanTokenClaimValidationResponse != nil {
		readBooleanTokenClaimValidationResponse(ctx, readResponse.BooleanTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.StringTokenClaimValidationResponse != nil {
		readStringTokenClaimValidationResponse(ctx, readResponse.StringTokenClaimValidationResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *tokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultTokenClaimValidationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateTokenClaimValidation(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateTokenClaimValidation(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan tokenClaimValidationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state tokenClaimValidationResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.TokenClaimValidationAPI.UpdateTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString(), plan.IdTokenValidatorName.ValueString())

	// Determine what update operations are necessary
	ops := createTokenClaimValidationOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.TokenClaimValidationAPI.UpdateTokenClaimValidationExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Token Claim Validation", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.StringArrayTokenClaimValidationResponse != nil {
			readStringArrayTokenClaimValidationResponse(ctx, updateResponse.StringArrayTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BooleanTokenClaimValidationResponse != nil {
			readBooleanTokenClaimValidationResponse(ctx, updateResponse.BooleanTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.StringTokenClaimValidationResponse != nil {
			readStringTokenClaimValidationResponse(ctx, updateResponse.StringTokenClaimValidationResponse, &state, &plan, &resp.Diagnostics)
		}
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultTokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *tokenClaimValidationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state tokenClaimValidationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.TokenClaimValidationAPI.DeleteTokenClaimValidationExecute(r.apiClient.TokenClaimValidationAPI.DeleteTokenClaimValidation(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.IdTokenValidatorName.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Token Claim Validation", err, httpResp)
		return
	}
}

func (r *tokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTokenClaimValidation(ctx, req, resp)
}

func (r *defaultTokenClaimValidationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importTokenClaimValidation(ctx, req, resp)
}

func importTokenClaimValidation(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	split := strings.Split(req.ID, "/")
	if len(split) != 2 {
		resp.Diagnostics.AddError("Invalid import id for resource", "Expected [id-token-validator-name]/[token-claim-validation-name]. Got: "+req.ID)
		return
	}
	// Set the required attributes to read the resource
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id_token_validator_name"), split[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), split[1])...)
}
