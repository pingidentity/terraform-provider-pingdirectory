package extendedoperationhandler

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &validateTotpPasswordExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &validateTotpPasswordExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &validateTotpPasswordExtendedOperationHandlerResource{}
	_ resource.Resource                = &defaultValidateTotpPasswordExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &defaultValidateTotpPasswordExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &defaultValidateTotpPasswordExtendedOperationHandlerResource{}
)

// Create a Validate Totp Password Extended Operation Handler resource
func NewValidateTotpPasswordExtendedOperationHandlerResource() resource.Resource {
	return &validateTotpPasswordExtendedOperationHandlerResource{}
}

func NewDefaultValidateTotpPasswordExtendedOperationHandlerResource() resource.Resource {
	return &defaultValidateTotpPasswordExtendedOperationHandlerResource{}
}

// validateTotpPasswordExtendedOperationHandlerResource is the resource implementation.
type validateTotpPasswordExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultValidateTotpPasswordExtendedOperationHandlerResource is the resource implementation.
type defaultValidateTotpPasswordExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *validateTotpPasswordExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_validate_totp_password_extended_operation_handler"
}

func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_validate_totp_password_extended_operation_handler"
}

// Configure adds the provider configured client to the resource.
func (r *validateTotpPasswordExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type validateTotpPasswordExtendedOperationHandlerResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	SharedSecretAttributeType types.String `tfsdk:"shared_secret_attribute_type"`
	TimeIntervalDuration      types.String `tfsdk:"time_interval_duration"`
	AdjacentIntervalsToCheck  types.Int64  `tfsdk:"adjacent_intervals_to_check"`
	PreventTOTPReuse          types.Bool   `tfsdk:"prevent_totp_reuse"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *validateTotpPasswordExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	validateTotpPasswordExtendedOperationHandlerSchema(ctx, req, resp, false)
}

func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	validateTotpPasswordExtendedOperationHandlerSchema(ctx, req, resp, true)
}

func validateTotpPasswordExtendedOperationHandlerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Validate Totp Password Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"shared_secret_attribute_type": schema.StringAttribute{
				Description: "The name or OID of the attribute that will be used to hold the shared secret key used during TOTP processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"time_interval_duration": schema.StringAttribute{
				Description: "The duration of the time interval used for TOTP processing.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"adjacent_intervals_to_check": schema.Int64Attribute{
				Description: "The number of adjacent time intervals (both before and after the current time) that should be checked when performing authentication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"prevent_totp_reuse": schema.BoolAttribute{
				Description: "Indicates whether to prevent clients from re-using TOTP passwords.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
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
func addOptionalValidateTotpPasswordExtendedOperationHandlerFields(ctx context.Context, addRequest *client.AddValidateTotpPasswordExtendedOperationHandlerRequest, plan validateTotpPasswordExtendedOperationHandlerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SharedSecretAttributeType) {
		addRequest.SharedSecretAttributeType = plan.SharedSecretAttributeType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TimeIntervalDuration) {
		addRequest.TimeIntervalDuration = plan.TimeIntervalDuration.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AdjacentIntervalsToCheck) {
		addRequest.AdjacentIntervalsToCheck = plan.AdjacentIntervalsToCheck.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.PreventTOTPReuse) {
		addRequest.PreventTOTPReuse = plan.PreventTOTPReuse.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a ValidateTotpPasswordExtendedOperationHandlerResponse object into the model struct
func readValidateTotpPasswordExtendedOperationHandlerResponse(ctx context.Context, r *client.ValidateTotpPasswordExtendedOperationHandlerResponse, state *validateTotpPasswordExtendedOperationHandlerResourceModel, expectedValues *validateTotpPasswordExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.SharedSecretAttributeType = internaltypes.StringTypeOrNil(r.SharedSecretAttributeType, internaltypes.IsEmptyString(expectedValues.SharedSecretAttributeType))
	state.TimeIntervalDuration = internaltypes.StringTypeOrNil(r.TimeIntervalDuration, internaltypes.IsEmptyString(expectedValues.TimeIntervalDuration))
	config.CheckMismatchedPDFormattedAttributes("time_interval_duration",
		expectedValues.TimeIntervalDuration, state.TimeIntervalDuration, diagnostics)
	state.AdjacentIntervalsToCheck = internaltypes.Int64TypeOrNil(r.AdjacentIntervalsToCheck)
	state.PreventTOTPReuse = internaltypes.BoolTypeOrNil(r.PreventTOTPReuse)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createValidateTotpPasswordExtendedOperationHandlerOperations(plan validateTotpPasswordExtendedOperationHandlerResourceModel, state validateTotpPasswordExtendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.SharedSecretAttributeType, state.SharedSecretAttributeType, "shared-secret-attribute-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TimeIntervalDuration, state.TimeIntervalDuration, "time-interval-duration")
	operations.AddInt64OperationIfNecessary(&ops, plan.AdjacentIntervalsToCheck, state.AdjacentIntervalsToCheck, "adjacent-intervals-to-check")
	operations.AddBoolOperationIfNecessary(&ops, plan.PreventTOTPReuse, state.PreventTOTPReuse, "prevent-totp-reuse")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *validateTotpPasswordExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan validateTotpPasswordExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddValidateTotpPasswordExtendedOperationHandlerRequest(plan.Id.ValueString(),
		[]client.EnumvalidateTotpPasswordExtendedOperationHandlerSchemaUrn{client.ENUMVALIDATETOTPPASSWORDEXTENDEDOPERATIONHANDLERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0EXTENDED_OPERATION_HANDLERVALIDATE_TOTP_PASSWORD},
		plan.Enabled.ValueBool())
	addOptionalValidateTotpPasswordExtendedOperationHandlerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddExtendedOperationHandlerRequest(
		client.AddValidateTotpPasswordExtendedOperationHandlerRequestAsAddExtendedOperationHandlerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.AddExtendedOperationHandlerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Validate Totp Password Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state validateTotpPasswordExtendedOperationHandlerResourceModel
	readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, addResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan validateTotpPasswordExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Validate Totp Password Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state validateTotpPasswordExtendedOperationHandlerResourceModel
	readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createValidateTotpPasswordExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Validate Totp Password Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *validateTotpPasswordExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readValidateTotpPasswordExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readValidateTotpPasswordExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readValidateTotpPasswordExtendedOperationHandler(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state validateTotpPasswordExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Validate Totp Password Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, readResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *validateTotpPasswordExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateValidateTotpPasswordExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateValidateTotpPasswordExtendedOperationHandler(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateValidateTotpPasswordExtendedOperationHandler(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan validateTotpPasswordExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state validateTotpPasswordExtendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createValidateTotpPasswordExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Validate Totp Password Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readValidateTotpPasswordExtendedOperationHandlerResponse(ctx, updateResponse.ValidateTotpPasswordExtendedOperationHandlerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *validateTotpPasswordExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state validateTotpPasswordExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandlerExecute(r.apiClient.ExtendedOperationHandlerApi.DeleteExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Validate Totp Password Extended Operation Handler", err, httpResp)
		return
	}
}

func (r *validateTotpPasswordExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importValidateTotpPasswordExtendedOperationHandler(ctx, req, resp)
}

func (r *defaultValidateTotpPasswordExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importValidateTotpPasswordExtendedOperationHandler(ctx, req, resp)
}

func importValidateTotpPasswordExtendedOperationHandler(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
