package extendedoperationhandler

import (
	"context"
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
	_ resource.Resource                = &deliverPasswordResetTokenExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &deliverPasswordResetTokenExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &deliverPasswordResetTokenExtendedOperationHandlerResource{}
	_ resource.Resource                = &defaultDeliverPasswordResetTokenExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultDeliverPasswordResetTokenExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultDeliverPasswordResetTokenExtendedOperationHandlerResource{}
)

// Create a Deliver Password Reset Token Extended Operation Handler resource
func NewDeliverPasswordResetTokenExtendedOperationHandlerResource() resource.Resource {
	return &deliverPasswordResetTokenExtendedOperationHandlerResource{}
}

func NewDefaultDeliverPasswordResetTokenExtendedOperationHandlerResource() resource.Resource {
	return &defaultDeliverPasswordResetTokenExtendedOperationHandlerResource{}
}

// deliverPasswordResetTokenExtendedOperationHandlerResource is the resource implementation.
type deliverPasswordResetTokenExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultDeliverPasswordResetTokenExtendedOperationHandlerResource is the resource implementation.
type defaultDeliverPasswordResetTokenExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_deliver_password_reset_token_extended_operation_handler"
}

func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_deliver_password_reset_token_extended_operation_handler"
}

// Configure adds the provider configured client to the resource.
func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type deliverPasswordResetTokenExtendedOperationHandlerResourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	LastUpdated                        types.String `tfsdk:"last_updated"`
	Notifications                      types.Set    `tfsdk:"notifications"`
	RequiredActions                    types.Set    `tfsdk:"required_actions"`
	PasswordGenerator                  types.String `tfsdk:"password_generator"`
	DefaultTokenDeliveryMechanism      types.Set    `tfsdk:"default_token_delivery_mechanism"`
	PasswordResetTokenValidityDuration types.String `tfsdk:"password_reset_token_validity_duration"`
	Description                        types.String `tfsdk:"description"`
	Enabled                            types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	deliverPasswordResetTokenExtendedOperationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	deliverPasswordResetTokenExtendedOperationHandlerSchema(ctx, req, resp, true)
}

func deliverPasswordResetTokenExtendedOperationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Deliver Password Reset Token Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"password_generator": schema.StringAttribute{
				Description: "The password generator that will be used to create the password reset token values to be delivered to the end user.",
				Required:    true,
			},
			"default_token_delivery_mechanism": schema.SetAttribute{
				Description: "The set of delivery mechanisms that may be used to deliver password reset tokens to users for requests that do not specify one or more preferred delivery mechanisms.",
				Required:    true,
				ElementType: types.StringType,
			},
			"password_reset_token_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a password reset token should be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Extended Operation Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Extended Operation Handler is enabled (that is, whether the types of extended operations are allowed in the server).",
				Required:    true,
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
func addOptionalDeliverPasswordResetTokenExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddDeliverPasswordResetTokenExtendedOperationHandlerRequest, plan deliverPasswordResetTokenExtendedOperationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordResetTokenValidityDuration) {
		addRequest.PasswordResetTokenValidityDuration = plan.PasswordResetTokenValidityDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a DeliverPasswordResetTokenExtendedOperationHandlerResponse object into the model struct
func readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx context.Context, r *client.DeliverPasswordResetTokenExtendedOperationHandlerResponse, state *deliverPasswordResetTokenExtendedOperationHandlerResourceModel, expectedValues *deliverPasswordResetTokenExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultTokenDeliveryMechanism = internaltypes.GetStringSet(r.DefaultTokenDeliveryMechanism)
	state.PasswordResetTokenValidityDuration = types.StringValue(r.PasswordResetTokenValidityDuration)
	config.CheckMismatchedPDFormattedAttributes("password_reset_token_validity_duration",
		expectedValues.PasswordResetTokenValidityDuration, state.PasswordResetTokenValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createDeliverPasswordResetTokenExtendedOperationHandlerOperations(plan deliverPasswordResetTokenExtendedOperationHandlerResourceModel, state deliverPasswordResetTokenExtendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordGenerator, state.PasswordGenerator, "password-generator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultTokenDeliveryMechanism, state.DefaultTokenDeliveryMechanism, "default-token-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordResetTokenValidityDuration, state.PasswordResetTokenValidityDuration, "password-reset-token-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var DefaultTokenDeliveryMechanismSlice []string
	plan.DefaultTokenDeliveryMechanism.ElementsAs(ctx, &DefaultTokenDeliveryMechanismSlice, false)
	addRequest := client.NewAddDeliverPasswordResetTokenExtendedOperationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumdeliverPasswordResetTokenExtendedOperationHandlerSchemaUrn{client.ENUMDELIVERPASSWORDRESETTOKENEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERDELIVER_PASSWORD_RESET_TOKEN},
		plan.PasswordGenerator.ValueString(),
		DefaultTokenDeliveryMechanismSlice,
		plan.Enabled.ValueBool())
	addOptionalDeliverPasswordResetTokenExtendedOperationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddDeliverPasswordResetTokenExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Deliver Password Reset Token Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, addResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Deliver Password Reset Token Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createDeliverPasswordResetTokenExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Deliver Password Reset Token Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readDeliverPasswordResetTokenExtendedOperationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Deliver Password Reset Token Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, readResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateDeliverPasswordResetTokenExtendedOperationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createDeliverPasswordResetTokenExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Deliver Password Reset Token Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readDeliverPasswordResetTokenExtendedOperationHandlerResponse(ctx, updateResponse.DeliverPasswordResetTokenExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state deliverPasswordResetTokenExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandlerExecute(r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Deliver Password Reset Token Extended Operation Handler", err, httpResp)
		return
	}
}

func (r *deliverPasswordResetTokenExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp)
}

func (r *defaultDeliverPasswordResetTokenExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importDeliverPasswordResetTokenExtendedOperationHandler(ctx, req, resp)
}

func importDeliverPasswordResetTokenExtendedOperationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
