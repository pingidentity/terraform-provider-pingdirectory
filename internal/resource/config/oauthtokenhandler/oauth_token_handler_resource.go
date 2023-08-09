package oauthtokenhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &oauthTokenHandlerResource{}
	_ resource.ResourceWithConfigure   = &oauthTokenHandlerResource{}
	_ resource.ResourceWithImportState = &oauthTokenHandlerResource{}
	_ resource.Resource                = &defaultOauthTokenHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultOauthTokenHandlerResource{}
	_ resource.ResourceWithImportState = &defaultOauthTokenHandlerResource{}
)

// Create a Oauth Token Handler resource
func NewOauthTokenHandlerResource() resource.Resource {
	return &oauthTokenHandlerResource{}
}

func NewDefaultOauthTokenHandlerResource() resource.Resource {
	return &defaultOauthTokenHandlerResource{}
}

// oauthTokenHandlerResource is the resource implementation.
type oauthTokenHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultOauthTokenHandlerResource is the resource implementation.
type defaultOauthTokenHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *oauthTokenHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_oauth_token_handler"
}

func (r *defaultOauthTokenHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_oauth_token_handler"
}

// Configure adds the provider configured client to the resource.
func (r *oauthTokenHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultOauthTokenHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type oauthTokenHandlerResourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	LastUpdated       types.String `tfsdk:"last_updated"`
	Notifications     types.Set    `tfsdk:"notifications"`
	RequiredActions   types.Set    `tfsdk:"required_actions"`
	Type              types.String `tfsdk:"type"`
	ExtensionClass    types.String `tfsdk:"extension_class"`
	ExtensionArgument types.Set    `tfsdk:"extension_argument"`
	ScriptClass       types.String `tfsdk:"script_class"`
	ScriptArgument    types.Set    `tfsdk:"script_argument"`
	Description       types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *oauthTokenHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthTokenHandlerSchema(ctx, req, resp, false)
}

func (r *defaultOauthTokenHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	oauthTokenHandlerSchema(ctx, req, resp, true)
}

func oauthTokenHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Oauth Token Handler.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of OAuth Token Handler resource. Options are ['groovy-scripted', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"groovy-scripted", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party OAuth Token Handler.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party OAuth Token Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted OAuth Token Handler.",
				Optional:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted OAuth Token Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this OAuth Token Handler",
				Optional:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsOauthTokenHandler() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_argument"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("script_class"),
			path.MatchRoot("type"),
			[]string{"groovy-scripted"},
		),
	}
}

// Add config validators
func (r oauthTokenHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsOauthTokenHandler()
}

// Add config validators
func (r defaultOauthTokenHandlerResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsOauthTokenHandler()
}

// Add optional fields to create request for groovy-scripted oauth-token-handler
func addOptionalGroovyScriptedOauthTokenHandlerFields(ctx context.Context, addRequest *client.AddGroovyScriptedOauthTokenHandlerRequest, plan oauthTokenHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party oauth-token-handler
func addOptionalThirdPartyOauthTokenHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyOauthTokenHandlerRequest, plan oauthTokenHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateOauthTokenHandlerUnknownValues(ctx context.Context, model *oauthTokenHandlerResourceModel) {
	if model.ScriptArgument.ElementType(ctx) == nil {
		model.ScriptArgument = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
}

// Read a GroovyScriptedOauthTokenHandlerResponse object into the model struct
func readGroovyScriptedOauthTokenHandlerResponse(ctx context.Context, r *client.GroovyScriptedOauthTokenHandlerResponse, state *oauthTokenHandlerResourceModel, expectedValues *oauthTokenHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateOauthTokenHandlerUnknownValues(ctx, state)
}

// Read a ThirdPartyOauthTokenHandlerResponse object into the model struct
func readThirdPartyOauthTokenHandlerResponse(ctx context.Context, r *client.ThirdPartyOauthTokenHandlerResponse, state *oauthTokenHandlerResourceModel, expectedValues *oauthTokenHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateOauthTokenHandlerUnknownValues(ctx, state)
}

// Create any update operations necessary to make the state match the plan
func createOauthTokenHandlerOperations(plan oauthTokenHandlerResourceModel, state oauthTokenHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a groovy-scripted oauth-token-handler
func (r *oauthTokenHandlerResource) CreateGroovyScriptedOauthTokenHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan oauthTokenHandlerResourceModel) (*oauthTokenHandlerResourceModel, error) {
	addRequest := client.NewAddGroovyScriptedOauthTokenHandlerRequest(plan.Name.ValueString(),
		[]client.EnumgroovyScriptedOauthTokenHandlerSchemaUrn{client.ENUMGROOVYSCRIPTEDOAUTHTOKENHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0OAUTH_TOKEN_HANDLERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString())
	addOptionalGroovyScriptedOauthTokenHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.OauthTokenHandlerApi.AddOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddOauthTokenHandlerRequest(
		client.AddGroovyScriptedOauthTokenHandlerRequestAsAddOauthTokenHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.OauthTokenHandlerApi.AddOauthTokenHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Oauth Token Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthTokenHandlerResourceModel
	readGroovyScriptedOauthTokenHandlerResponse(ctx, addResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party oauth-token-handler
func (r *oauthTokenHandlerResource) CreateThirdPartyOauthTokenHandler(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan oauthTokenHandlerResourceModel) (*oauthTokenHandlerResourceModel, error) {
	addRequest := client.NewAddThirdPartyOauthTokenHandlerRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyOauthTokenHandlerSchemaUrn{client.ENUMTHIRDPARTYOAUTHTOKENHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0OAUTH_TOKEN_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString())
	addOptionalThirdPartyOauthTokenHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.OauthTokenHandlerApi.AddOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddOauthTokenHandlerRequest(
		client.AddThirdPartyOauthTokenHandlerRequestAsAddOauthTokenHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.OauthTokenHandlerApi.AddOauthTokenHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Oauth Token Handler", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state oauthTokenHandlerResourceModel
	readThirdPartyOauthTokenHandlerResponse(ctx, addResponse.ThirdPartyOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *oauthTokenHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan oauthTokenHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *oauthTokenHandlerResourceModel
	var err error
	if plan.Type.ValueString() == "groovy-scripted" {
		state, err = r.CreateGroovyScriptedOauthTokenHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyOauthTokenHandler(ctx, req, resp, plan)
		if err != nil {
			return
		}
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
func (r *defaultOauthTokenHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan oauthTokenHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.OauthTokenHandlerApi.GetOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Oauth Token Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state oauthTokenHandlerResourceModel
	if readResponse.GroovyScriptedOauthTokenHandlerResponse != nil {
		readGroovyScriptedOauthTokenHandlerResponse(ctx, readResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyOauthTokenHandlerResponse != nil {
		readThirdPartyOauthTokenHandlerResponse(ctx, readResponse.ThirdPartyOauthTokenHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createOauthTokenHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Oauth Token Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GroovyScriptedOauthTokenHandlerResponse != nil {
			readGroovyScriptedOauthTokenHandlerResponse(ctx, updateResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyOauthTokenHandlerResponse != nil {
			readThirdPartyOauthTokenHandlerResponse(ctx, updateResponse.ThirdPartyOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *oauthTokenHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOauthTokenHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readOauthTokenHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state oauthTokenHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.OauthTokenHandlerApi.GetOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Oauth Token Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GroovyScriptedOauthTokenHandlerResponse != nil {
		readGroovyScriptedOauthTokenHandlerResponse(ctx, readResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyOauthTokenHandlerResponse != nil {
		readThirdPartyOauthTokenHandlerResponse(ctx, readResponse.ThirdPartyOauthTokenHandlerResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *oauthTokenHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultOauthTokenHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateOauthTokenHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan oauthTokenHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state oauthTokenHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createOauthTokenHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Oauth Token Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GroovyScriptedOauthTokenHandlerResponse != nil {
			readGroovyScriptedOauthTokenHandlerResponse(ctx, updateResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyOauthTokenHandlerResponse != nil {
			readThirdPartyOauthTokenHandlerResponse(ctx, updateResponse.ThirdPartyOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultOauthTokenHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *oauthTokenHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state oauthTokenHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.OauthTokenHandlerApi.DeleteOauthTokenHandlerExecute(r.apiClient.OauthTokenHandlerApi.DeleteOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Oauth Token Handler", err, httpResp)
		return
	}
}

func (r *oauthTokenHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOauthTokenHandler(ctx, req, resp)
}

func (r *defaultOauthTokenHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importOauthTokenHandler(ctx, req, resp)
}

func importOauthTokenHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
