package conjurauthenticationmethod

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
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
	_ resource.Resource                = &conjurAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &conjurAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &conjurAuthenticationMethodResource{}
	_ resource.Resource                = &defaultConjurAuthenticationMethodResource{}
	_ resource.ResourceWithConfigure   = &defaultConjurAuthenticationMethodResource{}
	_ resource.ResourceWithImportState = &defaultConjurAuthenticationMethodResource{}
)

// Create a Conjur Authentication Method resource
func NewConjurAuthenticationMethodResource() resource.Resource {
	return &conjurAuthenticationMethodResource{}
}

func NewDefaultConjurAuthenticationMethodResource() resource.Resource {
	return &defaultConjurAuthenticationMethodResource{}
}

// conjurAuthenticationMethodResource is the resource implementation.
type conjurAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConjurAuthenticationMethodResource is the resource implementation.
type defaultConjurAuthenticationMethodResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *conjurAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_conjur_authentication_method"
}

func (r *defaultConjurAuthenticationMethodResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_conjur_authentication_method"
}

// Configure adds the provider configured client to the resource.
func (r *conjurAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultConjurAuthenticationMethodResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type conjurAuthenticationMethodResourceModel struct {
	Id              types.String `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	Type            types.String `tfsdk:"type"`
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	ApiKey          types.String `tfsdk:"api_key"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *conjurAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	conjurAuthenticationMethodSchema(ctx, req, resp, false)
}

func (r *defaultConjurAuthenticationMethodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	conjurAuthenticationMethodSchema(ctx, req, resp, true)
}

func conjurAuthenticationMethodSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Conjur Authentication Method.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Conjur Authentication Method resource. Options are ['api-key']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("api-key"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"api-key"}...),
				},
			},
			"username": schema.StringAttribute{
				Description: "The username for the user to authenticate.",
				Required:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password for the user to authenticate. This will be used to obtain an API key for the target user.",
				Optional:    true,
				Sensitive:   true,
			},
			"api_key": schema.StringAttribute{
				Description: "The API key for the user to authenticate.",
				Optional:    true,
				Sensitive:   true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Conjur Authentication Method",
				Optional:    true,
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
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for api-key conjur-authentication-method
func addOptionalApiKeyConjurAuthenticationMethodFields(ctx context.Context, addRequest *client.AddApiKeyConjurAuthenticationMethodRequest, plan conjurAuthenticationMethodResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Password) {
		addRequest.Password = plan.Password.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ApiKey) {
		addRequest.ApiKey = plan.ApiKey.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateConjurAuthenticationMethodUnknownValues(ctx context.Context, model *conjurAuthenticationMethodResourceModel) {
	if model.Password.IsUnknown() {
		model.Password = types.StringNull()
	}
	if model.ApiKey.IsUnknown() {
		model.ApiKey = types.StringNull()
	}
}

// Read a ApiKeyConjurAuthenticationMethodResponse object into the model struct
func readApiKeyConjurAuthenticationMethodResponse(ctx context.Context, r *client.ApiKeyConjurAuthenticationMethodResponse, state *conjurAuthenticationMethodResourceModel, expectedValues *conjurAuthenticationMethodResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("api-key")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Username = types.StringValue(r.Username)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateConjurAuthenticationMethodUnknownValues(ctx, state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *conjurAuthenticationMethodResourceModel) setStateValuesNotReturnedByAPI(expectedValues *conjurAuthenticationMethodResourceModel) {
	if !expectedValues.Password.IsUnknown() {
		state.Password = expectedValues.Password
	}
	if !expectedValues.ApiKey.IsUnknown() {
		state.ApiKey = expectedValues.ApiKey
	}
}

// Create any update operations necessary to make the state match the plan
func createConjurAuthenticationMethodOperations(plan conjurAuthenticationMethodResourceModel, state conjurAuthenticationMethodResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Username, state.Username, "username")
	operations.AddStringOperationIfNecessary(&ops, plan.Password, state.Password, "password")
	operations.AddStringOperationIfNecessary(&ops, plan.ApiKey, state.ApiKey, "api-key")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a api-key conjur-authentication-method
func (r *conjurAuthenticationMethodResource) CreateApiKeyConjurAuthenticationMethod(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan conjurAuthenticationMethodResourceModel) (*conjurAuthenticationMethodResourceModel, error) {
	addRequest := client.NewAddApiKeyConjurAuthenticationMethodRequest(plan.Name.ValueString(),
		[]client.EnumapiKeyConjurAuthenticationMethodSchemaUrn{client.ENUMAPIKEYCONJURAUTHENTICATIONMETHODSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CONJUR_AUTHENTICATION_METHODAPI_KEY},
		plan.Username.ValueString())
	addOptionalApiKeyConjurAuthenticationMethodFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ConjurAuthenticationMethodApi.AddConjurAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddApiKeyConjurAuthenticationMethodRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.ConjurAuthenticationMethodApi.AddConjurAuthenticationMethodExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Conjur Authentication Method", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state conjurAuthenticationMethodResourceModel
	readApiKeyConjurAuthenticationMethodResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *conjurAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan conjurAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateApiKeyConjurAuthenticationMethod(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultConjurAuthenticationMethodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan conjurAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConjurAuthenticationMethodApi.GetConjurAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur Authentication Method", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state conjurAuthenticationMethodResourceModel
	readApiKeyConjurAuthenticationMethodResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ConjurAuthenticationMethodApi.UpdateConjurAuthenticationMethod(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createConjurAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ConjurAuthenticationMethodApi.UpdateConjurAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Conjur Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readApiKeyConjurAuthenticationMethodResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *conjurAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConjurAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultConjurAuthenticationMethodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConjurAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readConjurAuthenticationMethod(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state conjurAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ConjurAuthenticationMethodApi.GetConjurAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur Authentication Method", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur Authentication Method", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readApiKeyConjurAuthenticationMethodResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *conjurAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConjurAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConjurAuthenticationMethodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConjurAuthenticationMethod(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConjurAuthenticationMethod(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan conjurAuthenticationMethodResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state conjurAuthenticationMethodResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ConjurAuthenticationMethodApi.UpdateConjurAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createConjurAuthenticationMethodOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ConjurAuthenticationMethodApi.UpdateConjurAuthenticationMethodExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Conjur Authentication Method", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readApiKeyConjurAuthenticationMethodResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
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
func (r *defaultConjurAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *conjurAuthenticationMethodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state conjurAuthenticationMethodResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ConjurAuthenticationMethodApi.DeleteConjurAuthenticationMethodExecute(r.apiClient.ConjurAuthenticationMethodApi.DeleteConjurAuthenticationMethod(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Conjur Authentication Method", err, httpResp)
		return
	}
}

func (r *conjurAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConjurAuthenticationMethod(ctx, req, resp)
}

func (r *defaultConjurAuthenticationMethodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConjurAuthenticationMethod(ctx, req, resp)
}

func importConjurAuthenticationMethod(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
