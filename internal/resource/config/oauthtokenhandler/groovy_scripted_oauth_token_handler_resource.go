package oauthtokenhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &groovyScriptedOauthTokenHandlerResource{}
	_ resource.ResourceWithConfigure   = &groovyScriptedOauthTokenHandlerResource{}
	_ resource.ResourceWithImportState = &groovyScriptedOauthTokenHandlerResource{}
	_ resource.Resource                = &defaultGroovyScriptedOauthTokenHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultGroovyScriptedOauthTokenHandlerResource{}
	_ resource.ResourceWithImportState = &defaultGroovyScriptedOauthTokenHandlerResource{}
)

// Create a Groovy Scripted Oauth Token Handler resource
func NewGroovyScriptedOauthTokenHandlerResource() resource.Resource {
	return &groovyScriptedOauthTokenHandlerResource{}
}

func NewDefaultGroovyScriptedOauthTokenHandlerResource() resource.Resource {
	return &defaultGroovyScriptedOauthTokenHandlerResource{}
}

// groovyScriptedOauthTokenHandlerResource is the resource implementation.
type groovyScriptedOauthTokenHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGroovyScriptedOauthTokenHandlerResource is the resource implementation.
type defaultGroovyScriptedOauthTokenHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *groovyScriptedOauthTokenHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groovy_scripted_oauth_token_handler"
}

func (r *defaultGroovyScriptedOauthTokenHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_groovy_scripted_oauth_token_handler"
}

// Configure adds the provider configured client to the resource.
func (r *groovyScriptedOauthTokenHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultGroovyScriptedOauthTokenHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type groovyScriptedOauthTokenHandlerResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	ScriptClass     types.String `tfsdk:"script_class"`
	ScriptArgument  types.Set    `tfsdk:"script_argument"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *groovyScriptedOauthTokenHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedOauthTokenHandlerSchema(ctx, req, resp, false)
}

func (r *defaultGroovyScriptedOauthTokenHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedOauthTokenHandlerSchema(ctx, req, resp, true)
}

func groovyScriptedOauthTokenHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Groovy Scripted Oauth Token Handler.",
		Attributes: map[string]schema.Attribute{
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted OAuth Token Handler.",
				Required:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted OAuth Token Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this OAuth Token Handler",
				Optional:    true,
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
func addOptionalGroovyScriptedOauthTokenHandlerFields(ctx context.Context, addRequest *client.AddGroovyScriptedOauthTokenHandlerRequest, plan groovyScriptedOauthTokenHandlerResourceModel) {
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

// Read a GroovyScriptedOauthTokenHandlerResponse object into the model struct
func readGroovyScriptedOauthTokenHandlerResponse(ctx context.Context, r *client.GroovyScriptedOauthTokenHandlerResponse, state *groovyScriptedOauthTokenHandlerResourceModel, expectedValues *groovyScriptedOauthTokenHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGroovyScriptedOauthTokenHandlerOperations(plan groovyScriptedOauthTokenHandlerResourceModel, state groovyScriptedOauthTokenHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *groovyScriptedOauthTokenHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedOauthTokenHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGroovyScriptedOauthTokenHandlerRequest(plan.Id.ValueString(),
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
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Groovy Scripted Oauth Token Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state groovyScriptedOauthTokenHandlerResourceModel
	readGroovyScriptedOauthTokenHandlerResponse(ctx, addResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultGroovyScriptedOauthTokenHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedOauthTokenHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.OauthTokenHandlerApi.GetOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Oauth Token Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state groovyScriptedOauthTokenHandlerResourceModel
	readGroovyScriptedOauthTokenHandlerResponse(ctx, readResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGroovyScriptedOauthTokenHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Oauth Token Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedOauthTokenHandlerResponse(ctx, updateResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *groovyScriptedOauthTokenHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedOauthTokenHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGroovyScriptedOauthTokenHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state groovyScriptedOauthTokenHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.OauthTokenHandlerApi.GetOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Oauth Token Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGroovyScriptedOauthTokenHandlerResponse(ctx, readResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *groovyScriptedOauthTokenHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedOauthTokenHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedOauthTokenHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGroovyScriptedOauthTokenHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan groovyScriptedOauthTokenHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state groovyScriptedOauthTokenHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGroovyScriptedOauthTokenHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.OauthTokenHandlerApi.UpdateOauthTokenHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Oauth Token Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedOauthTokenHandlerResponse(ctx, updateResponse.GroovyScriptedOauthTokenHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGroovyScriptedOauthTokenHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *groovyScriptedOauthTokenHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groovyScriptedOauthTokenHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.OauthTokenHandlerApi.DeleteOauthTokenHandlerExecute(r.apiClient.OauthTokenHandlerApi.DeleteOauthTokenHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Groovy Scripted Oauth Token Handler", err, httpResp)
		return
	}
}

func (r *groovyScriptedOauthTokenHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedOauthTokenHandler(ctx, req, resp)
}

func (r *defaultGroovyScriptedOauthTokenHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedOauthTokenHandler(ctx, req, resp)
}

func importGroovyScriptedOauthTokenHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
