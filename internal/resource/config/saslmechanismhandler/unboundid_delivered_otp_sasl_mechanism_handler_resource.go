package saslmechanismhandler

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
	_ resource.Resource                = &unboundidDeliveredOtpSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &unboundidDeliveredOtpSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &unboundidDeliveredOtpSaslMechanismHandlerResource{}
	_ resource.Resource                = &defaultUnboundidDeliveredOtpSaslMechanismHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultUnboundidDeliveredOtpSaslMechanismHandlerResource{}
	_ resource.ResourceWithImportState = &defaultUnboundidDeliveredOtpSaslMechanismHandlerResource{}
)

// Create a Unboundid Delivered Otp Sasl Mechanism Handler resource
func NewUnboundidDeliveredOtpSaslMechanismHandlerResource() resource.Resource {
	return &unboundidDeliveredOtpSaslMechanismHandlerResource{}
}

func NewDefaultUnboundidDeliveredOtpSaslMechanismHandlerResource() resource.Resource {
	return &defaultUnboundidDeliveredOtpSaslMechanismHandlerResource{}
}

// unboundidDeliveredOtpSaslMechanismHandlerResource is the resource implementation.
type unboundidDeliveredOtpSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultUnboundidDeliveredOtpSaslMechanismHandlerResource is the resource implementation.
type defaultUnboundidDeliveredOtpSaslMechanismHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unboundid_delivered_otp_sasl_mechanism_handler"
}

func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_unboundid_delivered_otp_sasl_mechanism_handler"
}

// Configure adds the provider configured client to the resource.
func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type unboundidDeliveredOtpSaslMechanismHandlerResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	Notifications       types.Set    `tfsdk:"notifications"`
	RequiredActions     types.Set    `tfsdk:"required_actions"`
	IdentityMapper      types.String `tfsdk:"identity_mapper"`
	OtpValidityDuration types.String `tfsdk:"otp_validity_duration"`
	Description         types.String `tfsdk:"description"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	unboundidDeliveredOtpSaslMechanismHandlerSchema(ctx, req, resp, false)
}

func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	unboundidDeliveredOtpSaslMechanismHandlerSchema(ctx, req, resp, true)
}

func unboundidDeliveredOtpSaslMechanismHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Unboundid Delivered Otp Sasl Mechanism Handler.",
		Attributes: map[string]schema.Attribute{
			"identity_mapper": schema.StringAttribute{
				Description: "The identity mapper that should be used to identify the user(s) targeted in the authentication and/or authorization identities contained in the bind request. This will only be used for \"u:\"-style identities.",
				Required:    true,
			},
			"otp_validity_duration": schema.StringAttribute{
				Description: "The maximum length of time that a one-time password value should be considered valid.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this SASL Mechanism Handler",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SASL mechanism handler is enabled for use.",
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
func addOptionalUnboundidDeliveredOtpSaslMechanismHandlerFields(ctx context.Context, addRequest *client.AddUnboundidDeliveredOtpSaslMechanismHandlerRequest, plan unboundidDeliveredOtpSaslMechanismHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OtpValidityDuration) {
		addRequest.OtpValidityDuration = plan.OtpValidityDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a UnboundidDeliveredOtpSaslMechanismHandlerResponse object into the model struct
func readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx context.Context, r *client.UnboundidDeliveredOtpSaslMechanismHandlerResponse, state *unboundidDeliveredOtpSaslMechanismHandlerResourceModel, expectedValues *unboundidDeliveredOtpSaslMechanismHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.OtpValidityDuration = types.StringValue(r.OtpValidityDuration)
	config.CheckMismatchedPDFormattedAttributes("otp_validity_duration",
		expectedValues.OtpValidityDuration, state.OtpValidityDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createUnboundidDeliveredOtpSaslMechanismHandlerOperations(plan unboundidDeliveredOtpSaslMechanismHandlerResourceModel, state unboundidDeliveredOtpSaslMechanismHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.OtpValidityDuration, state.OtpValidityDuration, "otp-validity-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddUnboundidDeliveredOtpSaslMechanismHandlerRequest(plan.Id.ValueString(),
		[]client.EnumunboundidDeliveredOtpSaslMechanismHandlerSchemaUrn{client.ENUMUNBOUNDIDDELIVEREDOTPSASLMECHANISMHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0SASL_MECHANISM_HANDLERUNBOUNDID_DELIVERED_OTP},
		plan.IdentityMapper.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalUnboundidDeliveredOtpSaslMechanismHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SaslMechanismHandlerApi.AddSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSaslMechanismHandlerRequest(
		client.AddUnboundidDeliveredOtpSaslMechanismHandlerRequestAsAddSaslMechanismHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.AddSaslMechanismHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Unboundid Delivered Otp Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, addResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unboundid Delivered Otp Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createUnboundidDeliveredOtpSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unboundid Delivered Otp Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readUnboundidDeliveredOtpSaslMechanismHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SaslMechanismHandlerApi.GetSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Unboundid Delivered Otp Sasl Mechanism Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, readResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateUnboundidDeliveredOtpSaslMechanismHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createUnboundidDeliveredOtpSaslMechanismHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SaslMechanismHandlerApi.UpdateSaslMechanismHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Unboundid Delivered Otp Sasl Mechanism Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readUnboundidDeliveredOtpSaslMechanismHandlerResponse(ctx, updateResponse.UnboundidDeliveredOtpSaslMechanismHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state unboundidDeliveredOtpSaslMechanismHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SaslMechanismHandlerApi.DeleteSaslMechanismHandlerExecute(r.apiClient.SaslMechanismHandlerApi.DeleteSaslMechanismHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Unboundid Delivered Otp Sasl Mechanism Handler", err, httpResp)
		return
	}
}

func (r *unboundidDeliveredOtpSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp)
}

func (r *defaultUnboundidDeliveredOtpSaslMechanismHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importUnboundidDeliveredOtpSaslMechanismHandler(ctx, req, resp)
}

func importUnboundidDeliveredOtpSaslMechanismHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
