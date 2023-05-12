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
	_ resource.Resource                = &generatePasswordExtendedOperationHandlerResource{}
	_ resource.ResourceWithConfigure   = &generatePasswordExtendedOperationHandlerResource{}
	_ resource.ResourceWithImportState = &generatePasswordExtendedOperationHandlerResource{}
)

// Create a Generate Password Extended Operation Handler resource
func NewGeneratePasswordExtendedOperationHandlerResource() resource.Resource {
	return &generatePasswordExtendedOperationHandlerResource{}
}

// generatePasswordExtendedOperationHandlerResource is the resource implementation.
type generatePasswordExtendedOperationHandlerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *generatePasswordExtendedOperationHandlerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_generate_password_extended_operation_handler"
}

// Configure adds the provider configured client to the resource.
func (r *generatePasswordExtendedOperationHandlerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type generatePasswordExtendedOperationHandlerResourceModel struct {
	Id                                   types.String `tfsdk:"id"`
	LastUpdated                          types.String `tfsdk:"last_updated"`
	Notifications                        types.Set    `tfsdk:"notifications"`
	RequiredActions                      types.Set    `tfsdk:"required_actions"`
	DefaultPasswordPolicy                types.String `tfsdk:"default_password_policy"`
	DefaultPasswordGenerator             types.String `tfsdk:"default_password_generator"`
	MaximumPasswordsPerRequest           types.Int64  `tfsdk:"maximum_passwords_per_request"`
	MaximumValidationAttemptsPerPassword types.Int64  `tfsdk:"maximum_validation_attempts_per_password"`
	Description                          types.String `tfsdk:"description"`
	Enabled                              types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *generatePasswordExtendedOperationHandlerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schema := schema.Schema{
		Description: "Manages a Generate Password Extended Operation Handler.",
		Attributes: map[string]schema.Attribute{
			"default_password_policy": schema.StringAttribute{
				Description: "The default password policy that should be used when generating and validating passwords if the request does not specify an alternate policy. If this is not provided, then this Generate Password Extended Operation Handler will use the default password policy defined in the global configuration.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_password_generator": schema.StringAttribute{
				Description: "The default password generator that will be used if the selected password policy is not configured with a password generator.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"maximum_passwords_per_request": schema.Int64Attribute{
				Description: "The maximum number of passwords that may be generated and returned to the client for a single request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"maximum_validation_attempts_per_password": schema.Int64Attribute{
				Description: "The maximum number of attempts that the server may use to generate a password that passes validation.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Extended Operation Handler",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Extended Operation Handler is enabled (that is, whether the types of extended operations are allowed in the server).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonSchema(&schema, true)
	resp.Schema = schema
}

// Read a GeneratePasswordExtendedOperationHandlerResponse object into the model struct
func readGeneratePasswordExtendedOperationHandlerResponse(ctx context.Context, r *client.GeneratePasswordExtendedOperationHandlerResponse, state *generatePasswordExtendedOperationHandlerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.DefaultPasswordPolicy = internaltypes.StringTypeOrNil(r.DefaultPasswordPolicy, true)
	state.DefaultPasswordGenerator = types.StringValue(r.DefaultPasswordGenerator)
	state.MaximumPasswordsPerRequest = internaltypes.Int64TypeOrNil(r.MaximumPasswordsPerRequest)
	state.MaximumValidationAttemptsPerPassword = internaltypes.Int64TypeOrNil(r.MaximumValidationAttemptsPerPassword)
	state.Description = internaltypes.StringTypeOrNil(r.Description, true)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createGeneratePasswordExtendedOperationHandlerOperations(plan generatePasswordExtendedOperationHandlerResourceModel, state generatePasswordExtendedOperationHandlerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultPasswordPolicy, state.DefaultPasswordPolicy, "default-password-policy")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultPasswordGenerator, state.DefaultPasswordGenerator, "default-password-generator")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumPasswordsPerRequest, state.MaximumPasswordsPerRequest, "maximum-passwords-per-request")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaximumValidationAttemptsPerPassword, state.MaximumValidationAttemptsPerPassword, "maximum-validation-attempts-per-password")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *generatePasswordExtendedOperationHandlerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan generatePasswordExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generate Password Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state generatePasswordExtendedOperationHandlerResourceModel
	readGeneratePasswordExtendedOperationHandlerResponse(ctx, readResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createGeneratePasswordExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generate Password Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGeneratePasswordExtendedOperationHandlerResponse(ctx, updateResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
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
func (r *generatePasswordExtendedOperationHandlerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state generatePasswordExtendedOperationHandlerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.GetExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Generate Password Extended Operation Handler", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readGeneratePasswordExtendedOperationHandlerResponse(ctx, readResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *generatePasswordExtendedOperationHandlerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan generatePasswordExtendedOperationHandlerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state generatePasswordExtendedOperationHandlerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandler(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createGeneratePasswordExtendedOperationHandlerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.ExtendedOperationHandlerApi.UpdateExtendedOperationHandlerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Generate Password Extended Operation Handler", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readGeneratePasswordExtendedOperationHandlerResponse(ctx, updateResponse.GeneratePasswordExtendedOperationHandlerResponse, &state, &resp.Diagnostics)
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
func (r *generatePasswordExtendedOperationHandlerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *generatePasswordExtendedOperationHandlerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
