package changesubscriptionhandler

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
	_ resource.Resource                = &groovyScriptedChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithConfigure   = &groovyScriptedChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithImportState = &groovyScriptedChangeSubscriptionHandlerResource{}
	_ resource.Resource                = &defaultGroovyScriptedChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultGroovyScriptedChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultGroovyScriptedChangeSubscriptionHandlerResource{}
)

// Create a Groovy Scripted Change Subscription Handler resource
func NewGroovyScriptedChangeSubscriptionHandlerResource() resource.Resource {
	return &groovyScriptedChangeSubscriptionHandlerResource{}
}

func NewDefaultGroovyScriptedChangeSubscriptionHandlerResource() resource.Resource {
	return &defaultGroovyScriptedChangeSubscriptionHandlerResource{}
}

// groovyScriptedChangeSubscriptionHandlerResource is the resource implementation.
type groovyScriptedChangeSubscriptionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultGroovyScriptedChangeSubscriptionHandlerResource is the resource implementation.
type defaultGroovyScriptedChangeSubscriptionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *groovyScriptedChangeSubscriptionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_groovy_scripted_change_subscription_handler"
}

func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_groovy_scripted_change_subscription_handler"
}

// Configure adds the provider configured client to the resource.
func (r *groovyScriptedChangeSubscriptionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type groovyScriptedChangeSubscriptionHandlerResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	ScriptClass        types.String `tfsdk:"script_class"`
	ScriptArgument     types.Set    `tfsdk:"script_argument"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ChangeSubscription types.Set    `tfsdk:"change_subscription"`
}

// GetSchema defines the schema for the resource.
func (r *groovyScriptedChangeSubscriptionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedChangeSubscriptionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	groovyScriptedChangeSubscriptionHandlerSchema(ctx, req, resp, true)
}

func groovyScriptedChangeSubscriptionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Groovy Scripted Change Subscription Handler.",
		Attributes: map[string]schema.Attribute{
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Change Subscription Handler.",
				Required:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Change Subscription Handler. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Change Subscription Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this change subscription handler is enabled within the server.",
				Required:    true,
			},
			"change_subscription": schema.SetAttribute{
				Description: "The set of change subscriptions for which this change subscription handler should be notified. If no values are provided then it will be notified for all change subscriptions defined in the server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
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
func addOptionalGroovyScriptedChangeSubscriptionHandlerFields(ctx context.Context, addRequest *client.AddGroovyScriptedChangeSubscriptionHandlerRequest, plan groovyScriptedChangeSubscriptionHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ScriptArgument) {
		var slice []string
		plan.ScriptArgument.ElementsAs(ctx, &slice, false)
		addRequest.ScriptArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.ChangeSubscription) {
		var slice []string
		plan.ChangeSubscription.ElementsAs(ctx, &slice, false)
		addRequest.ChangeSubscription = slice
	}
}

// Read a GroovyScriptedChangeSubscriptionHandlerResponse object into the model struct
func readGroovyScriptedChangeSubscriptionHandlerResponse(ctx context.Context, r *client.GroovyScriptedChangeSubscriptionHandlerResponse, state *groovyScriptedChangeSubscriptionHandlerResourceModel, expectedValues *groovyScriptedChangeSubscriptionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGroovyScriptedChangeSubscriptionHandlerOperations(plan groovyScriptedChangeSubscriptionHandlerResourceModel, state groovyScriptedChangeSubscriptionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ScriptClass, state.ScriptClass, "script-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ScriptArgument, state.ScriptArgument, "script-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangeSubscription, state.ChangeSubscription, "change-subscription")
	return ops
}

// Create a new resource
func (r *groovyScriptedChangeSubscriptionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedChangeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddGroovyScriptedChangeSubscriptionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumgroovyScriptedChangeSubscriptionHandlerSchemaUrn{client.ENUMGROOVYSCRIPTEDCHANGESUBSCRIPTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CHANGE_SUBSCRIPTION_HANDLERGROOVY_SCRIPTED},
		plan.ScriptClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalGroovyScriptedChangeSubscriptionHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ChangeSubscriptionHandlerApi.AddChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddChangeSubscriptionHandlerRequest(
		client.AddGroovyScriptedChangeSubscriptionHandlerRequestAsAddChangeSubscriptionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.AddChangeSubscriptionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Groovy Scripted Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state groovyScriptedChangeSubscriptionHandlerResourceModel
	readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, addResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan groovyScriptedChangeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state groovyScriptedChangeSubscriptionHandlerResourceModel
	readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, readResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGroovyScriptedChangeSubscriptionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Change Subscription Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, updateResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *groovyScriptedChangeSubscriptionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readGroovyScriptedChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readGroovyScriptedChangeSubscriptionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state groovyScriptedChangeSubscriptionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Groovy Scripted Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, readResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *groovyScriptedChangeSubscriptionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateGroovyScriptedChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateGroovyScriptedChangeSubscriptionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan groovyScriptedChangeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state groovyScriptedChangeSubscriptionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGroovyScriptedChangeSubscriptionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Groovy Scripted Change Subscription Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGroovyScriptedChangeSubscriptionHandlerResponse(ctx, updateResponse.GroovyScriptedChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *groovyScriptedChangeSubscriptionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state groovyScriptedChangeSubscriptionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.DeleteChangeSubscriptionHandlerExecute(r.apiClient.ChangeSubscriptionHandlerApi.DeleteChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Groovy Scripted Change Subscription Handler", err, httpResp)
		return
	}
}

func (r *groovyScriptedChangeSubscriptionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedChangeSubscriptionHandler(ctx, req, resp)
}

func (r *defaultGroovyScriptedChangeSubscriptionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importGroovyScriptedChangeSubscriptionHandler(ctx, req, resp)
}

func importGroovyScriptedChangeSubscriptionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
