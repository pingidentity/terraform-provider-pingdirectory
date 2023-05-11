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
	_ resource.Resource                = &singleUseTokensExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &singleUseTokensExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &singleUseTokensExtendedOperationHandlerResource{}
	_ resource.Resource                = &defaultSingleUseTokensExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultSingleUseTokensExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultSingleUseTokensExtendedOperationHandlerResource{}
)

// Create a Single Use Tokens Extended Operation Handler resource
func NewSingleUseTokensExtendedOperationHandlerResource() resource.Resource {
	return &singleUseTokensExtendedOperationHandlerResource{}
}

func NewDefaultSingleUseTokensExtendedOperationHandlerResource() resource.Resource {
	return &defaultSingleUseTokensExtendedOperationHandlerResource{}
}

// singleUseTokensExtendedOperationHandlerResource is the resource implementation.
type singleUseTokensExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSingleUseTokensExtendedOperationHandlerResource is the resource implementation.
type defaultSingleUseTokensExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *singleUseTokensExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_single_use_tokens_extended_operation_handler"
}

func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_single_use_tokens_extended_operation_handler"
}

// Configure adds the provider configured client to the resource.
func (r *singleUseTokensExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type singleUseTokensExtendedOperationHandlerResourceModel struct {
	Id                                    types.String `tfsdk:"id"`
	LastUpdated                           types.String `tfsdk:"last_updated"`
	Notifications                         types.Set    `tfsdk:"notifications"`
	RequiredActions                       types.Set    `tfsdk:"required_actions"`
	PasswordGenerator                     types.String `tfsdk:"password_generator"`
	DefaultOTPDeliveryMechanism           types.Set    `tfsdk:"default_otp_delivery_mechanism"`
	DefaultSingleUseTokenValidityDuration types.String `tfsdk:"default_single_use_token_validity_duration"`
	Description                           types.String `tfsdk:"description"`
	Enabled                               types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *singleUseTokensExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	singleUseTokensExtendedOperationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	singleUseTokensExtendedOperationHandlerSchema(ctx, req, resp, true)
}

func singleUseTokensExtendedOperationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Single Use Tokens Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"password_generator": schema.StringAttribute{
				Description: "The password generator that will be used to create the single-use token values to be delivered to the end user.",
				Required:    true,
			},
			"default_otp_delivery_mechanism": schema.SetAttribute{
				Description: "The set of delivery mechanisms that may be used to deliver single-use tokens to users in requests that do not specify one or more preferred delivery mechanisms.",
				Required:    true,
				ElementType: types.StringType,
			},
			"default_single_use_token_validity_duration": schema.StringAttribute{
				Description: "The default length of time that a single-use token will be considered valid by the server if the client doesn't specify a duration in the deliver single-use token request.",
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
func addOptionalSingleUseTokensExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddSingleUseTokensExtendedOperationHandlerRequest, plan singleUseTokensExtendedOperationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultSingleUseTokenValidityDuration) {
		addRequest.DefaultSingleUseTokenValidityDuration = plan.DefaultSingleUseTokenValidityDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a SingleUseTokensExtendedOperationHandlerResponse object into the model struct
func readSingleUseTokensExtendedOperationHandlerResponse(ctx context.Context, r *client.SingleUseTokensExtendedOperationHandlerResponse, state *singleUseTokensExtendedOperationHandlerResourceModel, expectedValues *singleUseTokensExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PasswordGenerator = types.StringValue(r.PasswordGenerator)
	state.DefaultOTPDeliveryMechanism = internaltypes.GetStringSet(r.DefaultOTPDeliveryMechanism)
	state.DefaultSingleUseTokenValidityDuration = internaltypes.StringTypeOrNil(r.DefaultSingleUseTokenValidityDuration, internaltypes.IsEmptyString(expectedValues.DefaultSingleUseTokenValidityDuration))
	config.CheckMismatchedPDFormattedAttributes("default_single_use_token_validity_duration",
		expectedValues.DefaultSingleUseTokenValidityDuration, state.DefaultSingleUseTokenValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSingleUseTokensExtendedOperationHandlerOperations(plan singleUseTokensExtendedOperationHandlerResourceModel, state singleUseTokensExtendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordGenerator, state.PasswordGenerator, "password-generator")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.DefaultOTPDeliveryMechanism, state.DefaultOTPDeliveryMechanism, "default-otp-delivery-mechanism")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultSingleUseTokenValidityDuration, state.DefaultSingleUseTokenValidityDuration, "default-single-use-token-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *singleUseTokensExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan singleUseTokensExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var DefaultOTPDeliveryMechanismSlice []string
	plan.DefaultOTPDeliveryMechanism.ElementsAs(ctx, &DefaultOTPDeliveryMechanismSlice, false)
	addRequest := client.NewAddSingleUseTokensExtendedOperationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumsingleUseTokensExtendedOperationHandlerSchemaUrn{client.ENUMSINGLEUSETOKENSEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERSINGLE_USE_TOKENS},
		plan.PasswordGenerator.ValueString(),
		DefaultOTPDeliveryMechanismSlice,
		plan.Enabled.ValueBool())
	addOptionalSingleUseTokensExtendedOperationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddSingleUseTokensExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Single Use Tokens Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state singleUseTokensExtendedOperationHandlerResourceModel
	readSingleUseTokensExtendedOperationHandlerResponse(ctx, addResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan singleUseTokensExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Single Use Tokens Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state singleUseTokensExtendedOperationHandlerResourceModel
	readSingleUseTokensExtendedOperationHandlerResponse(ctx, readResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSingleUseTokensExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Single Use Tokens Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSingleUseTokensExtendedOperationHandlerResponse(ctx, updateResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *singleUseTokensExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSingleUseTokensExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSingleUseTokensExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSingleUseTokensExtendedOperationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state singleUseTokensExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Single Use Tokens Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSingleUseTokensExtendedOperationHandlerResponse(ctx, readResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *singleUseTokensExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSingleUseTokensExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSingleUseTokensExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSingleUseTokensExtendedOperationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan singleUseTokensExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state singleUseTokensExtendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSingleUseTokensExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Single Use Tokens Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSingleUseTokensExtendedOperationHandlerResponse(ctx, updateResponse.SingleUseTokensExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSingleUseTokensExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *singleUseTokensExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state singleUseTokensExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandlerExecute(r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Single Use Tokens Extended Operation Handler", err, httpResp)
		return
	}
}

func (r *singleUseTokensExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSingleUseTokensExtendedOperationHandler(ctx, req, resp)
}

func (r *defaultSingleUseTokensExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSingleUseTokensExtendedOperationHandler(ctx, req, resp)
}

func importSingleUseTokensExtendedOperationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
