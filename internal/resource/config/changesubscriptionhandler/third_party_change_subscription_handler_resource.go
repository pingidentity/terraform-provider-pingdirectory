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
	_ resource.Resource                = &thirdPartyChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithConfigure   = &thirdPartyChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithImportState = &thirdPartyChangeSubscriptionHandlerResource{}
	_ resource.Resource                = &defaultThirdPartyChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultThirdPartyChangeSubscriptionHandlerResource{}
	_ resource.ResourceWithImportState = &defaultThirdPartyChangeSubscriptionHandlerResource{}
)

// Create a Third Party Change Subscription Handler resource
func NewThirdPartyChangeSubscriptionHandlerResource() resource.Resource {
	return &thirdPartyChangeSubscriptionHandlerResource{}
}

func NewDefaultThirdPartyChangeSubscriptionHandlerResource() resource.Resource {
	return &defaultThirdPartyChangeSubscriptionHandlerResource{}
}

// thirdPartyChangeSubscriptionHandlerResource is the resource implementation.
type thirdPartyChangeSubscriptionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultThirdPartyChangeSubscriptionHandlerResource is the resource implementation.
type defaultThirdPartyChangeSubscriptionHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *thirdPartyChangeSubscriptionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_third_party_change_subscription_handler"
}

func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_third_party_change_subscription_handler"
}

// Configure adds the provider configured client to the resource.
func (r *thirdPartyChangeSubscriptionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type thirdPartyChangeSubscriptionHandlerResourceModel struct {
	Id                 types.String `tfsdk:"id"`
	LastUpdated        types.String `tfsdk:"last_updated"`
	Notifications      types.Set    `tfsdk:"notifications"`
	RequiredActions    types.Set    `tfsdk:"required_actions"`
	ExtensionClass     types.String `tfsdk:"extension_class"`
	ExtensionArgument  types.Set    `tfsdk:"extension_argument"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	ChangeSubscription types.Set    `tfsdk:"change_subscription"`
}

// GetSchema defines the schema for the resource.
func (r *thirdPartyChangeSubscriptionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyChangeSubscriptionHandlerSchema(ctx, req, resp, false)
}

func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	thirdPartyChangeSubscriptionHandlerSchema(ctx, req, resp, true)
}

func thirdPartyChangeSubscriptionHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Third Party Change Subscription Handler.",
		Attributes: map[string]schema.Attribute{
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Change Subscription Handler.",
				Required:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Change Subscription Handler. Each configuration property should be given in the form 'name=value'.",
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
func addOptionalThirdPartyChangeSubscriptionHandlerFields(ctx context.Context, addRequest *client.AddThirdPartyChangeSubscriptionHandlerRequest, plan thirdPartyChangeSubscriptionHandlerResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
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

// Read a ThirdPartyChangeSubscriptionHandlerResponse object into the model struct
func readThirdPartyChangeSubscriptionHandlerResponse(ctx context.Context, r *client.ThirdPartyChangeSubscriptionHandlerResponse, state *thirdPartyChangeSubscriptionHandlerResourceModel, expectedValues *thirdPartyChangeSubscriptionHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.ChangeSubscription = internaltypes.GetStringSet(r.ChangeSubscription)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createThirdPartyChangeSubscriptionHandlerOperations(plan thirdPartyChangeSubscriptionHandlerResourceModel, state thirdPartyChangeSubscriptionHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ChangeSubscription, state.ChangeSubscription, "change-subscription")
	return ops
}

// Create a new resource
func (r *thirdPartyChangeSubscriptionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyChangeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddThirdPartyChangeSubscriptionHandlerRequest(plan.Id.ValueString(),
		[]client.EnumthirdPartyChangeSubscriptionHandlerSchemaUrn{client.ENUMTHIRDPARTYCHANGESUBSCRIPTIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CHANGE_SUBSCRIPTION_HANDLERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyChangeSubscriptionHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ChangeSubscriptionHandlerApi.AddChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddChangeSubscriptionHandlerRequest(
		client.AddThirdPartyChangeSubscriptionHandlerRequestAsAddChangeSubscriptionHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.AddChangeSubscriptionHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Third Party Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state thirdPartyChangeSubscriptionHandlerResourceModel
	readThirdPartyChangeSubscriptionHandlerResponse(ctx, addResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan thirdPartyChangeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state thirdPartyChangeSubscriptionHandlerResourceModel
	readThirdPartyChangeSubscriptionHandlerResponse(ctx, readResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createThirdPartyChangeSubscriptionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Change Subscription Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyChangeSubscriptionHandlerResponse(ctx, updateResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *thirdPartyChangeSubscriptionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readThirdPartyChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readThirdPartyChangeSubscriptionHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state thirdPartyChangeSubscriptionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ChangeSubscriptionHandlerApi.GetChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Third Party Change Subscription Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readThirdPartyChangeSubscriptionHandlerResponse(ctx, readResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *thirdPartyChangeSubscriptionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateThirdPartyChangeSubscriptionHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateThirdPartyChangeSubscriptionHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan thirdPartyChangeSubscriptionHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state thirdPartyChangeSubscriptionHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createThirdPartyChangeSubscriptionHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ChangeSubscriptionHandlerApi.UpdateChangeSubscriptionHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Third Party Change Subscription Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readThirdPartyChangeSubscriptionHandlerResponse(ctx, updateResponse.ThirdPartyChangeSubscriptionHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultThirdPartyChangeSubscriptionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *thirdPartyChangeSubscriptionHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state thirdPartyChangeSubscriptionHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ChangeSubscriptionHandlerApi.DeleteChangeSubscriptionHandlerExecute(r.apiClient.ChangeSubscriptionHandlerApi.DeleteChangeSubscriptionHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Third Party Change Subscription Handler", err, httpResp)
		return
	}
}

func (r *thirdPartyChangeSubscriptionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyChangeSubscriptionHandler(ctx, req, resp)
}

func (r *defaultThirdPartyChangeSubscriptionHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importThirdPartyChangeSubscriptionHandler(ctx, req, resp)
}

func importThirdPartyChangeSubscriptionHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
