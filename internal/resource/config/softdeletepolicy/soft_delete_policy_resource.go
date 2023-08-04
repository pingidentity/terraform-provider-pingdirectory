package softdeletepolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	_ resource.Resource                = &softDeletePolicyResource{}
	_ resource.ResourceWithConfigure   = &softDeletePolicyResource{}
	_ resource.ResourceWithImportState = &softDeletePolicyResource{}
	_ resource.Resource                = &defaultSoftDeletePolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultSoftDeletePolicyResource{}
	_ resource.ResourceWithImportState = &defaultSoftDeletePolicyResource{}
)

// Create a Soft Delete Policy resource
func NewSoftDeletePolicyResource() resource.Resource {
	return &softDeletePolicyResource{}
}

func NewDefaultSoftDeletePolicyResource() resource.Resource {
	return &defaultSoftDeletePolicyResource{}
}

// softDeletePolicyResource is the resource implementation.
type softDeletePolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSoftDeletePolicyResource is the resource implementation.
type defaultSoftDeletePolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *softDeletePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_soft_delete_policy"
}

func (r *defaultSoftDeletePolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_soft_delete_policy"
}

// Configure adds the provider configured client to the resource.
func (r *softDeletePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultSoftDeletePolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type softDeletePolicyResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	LastUpdated                      types.String `tfsdk:"last_updated"`
	Notifications                    types.Set    `tfsdk:"notifications"`
	RequiredActions                  types.Set    `tfsdk:"required_actions"`
	Type                             types.String `tfsdk:"type"`
	Description                      types.String `tfsdk:"description"`
	AutoSoftDeleteConnectionCriteria types.String `tfsdk:"auto_soft_delete_connection_criteria"`
	AutoSoftDeleteRequestCriteria    types.String `tfsdk:"auto_soft_delete_request_criteria"`
	SoftDeleteRetentionTime          types.String `tfsdk:"soft_delete_retention_time"`
	SoftDeleteRetainNumberOfEntries  types.Int64  `tfsdk:"soft_delete_retain_number_of_entries"`
}

// GetSchema defines the schema for the resource.
func (r *softDeletePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	softDeletePolicySchema(ctx, req, resp, false)
}

func (r *defaultSoftDeletePolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	softDeletePolicySchema(ctx, req, resp, true)
}

func softDeletePolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Soft Delete Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Soft Delete Policy resource. Options are ['soft-delete-policy']",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("soft-delete-policy"),
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"soft-delete-policy"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Soft Delete Policy",
				Optional:    true,
			},
			"auto_soft_delete_connection_criteria": schema.StringAttribute{
				Description: "Connection criteria used to automatically identify a delete operation for processing as a soft delete request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"auto_soft_delete_request_criteria": schema.StringAttribute{
				Description: "Request criteria used to automatically identify a delete operation for processing as a soft delete request.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_delete_retention_time": schema.StringAttribute{
				Description: "Specifies the maximum length of time that soft delete entries are retained before they are eligible to purged automatically.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"soft_delete_retain_number_of_entries": schema.Int64Attribute{
				Description: "Specifies the number of soft deleted entries to retain before the oldest entries are purged.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		// Add any default properties and set optional properties to computed where necessary
		config.SetAllAttributesToOptionalAndComputed(&schemaDef)
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add optional fields to create request for soft-delete-policy soft-delete-policy
func addOptionalSoftDeletePolicyFields(ctx context.Context, addRequest *client.AddSoftDeletePolicyRequest, plan softDeletePolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AutoSoftDeleteConnectionCriteria) {
		addRequest.AutoSoftDeleteConnectionCriteria = plan.AutoSoftDeleteConnectionCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AutoSoftDeleteRequestCriteria) {
		addRequest.AutoSoftDeleteRequestCriteria = plan.AutoSoftDeleteRequestCriteria.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SoftDeleteRetentionTime) {
		addRequest.SoftDeleteRetentionTime = plan.SoftDeleteRetentionTime.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.SoftDeleteRetainNumberOfEntries) {
		addRequest.SoftDeleteRetainNumberOfEntries = plan.SoftDeleteRetainNumberOfEntries.ValueInt64Pointer()
	}
}

// Read a SoftDeletePolicyResponse object into the model struct
func readSoftDeletePolicyResponse(ctx context.Context, r *client.SoftDeletePolicyResponse, state *softDeletePolicyResourceModel, expectedValues *softDeletePolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("soft-delete-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.AutoSoftDeleteConnectionCriteria = internaltypes.StringTypeOrNil(r.AutoSoftDeleteConnectionCriteria, internaltypes.IsEmptyString(expectedValues.AutoSoftDeleteConnectionCriteria))
	state.AutoSoftDeleteRequestCriteria = internaltypes.StringTypeOrNil(r.AutoSoftDeleteRequestCriteria, internaltypes.IsEmptyString(expectedValues.AutoSoftDeleteRequestCriteria))
	state.SoftDeleteRetentionTime = internaltypes.StringTypeOrNil(r.SoftDeleteRetentionTime, internaltypes.IsEmptyString(expectedValues.SoftDeleteRetentionTime))
	config.CheckMismatchedPDFormattedAttributes("soft_delete_retention_time",
		expectedValues.SoftDeleteRetentionTime, state.SoftDeleteRetentionTime, diagnostics)
	state.SoftDeleteRetainNumberOfEntries = internaltypes.Int64TypeOrNil(r.SoftDeleteRetainNumberOfEntries)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSoftDeletePolicyOperations(plan softDeletePolicyResourceModel, state softDeletePolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.AutoSoftDeleteConnectionCriteria, state.AutoSoftDeleteConnectionCriteria, "auto-soft-delete-connection-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.AutoSoftDeleteRequestCriteria, state.AutoSoftDeleteRequestCriteria, "auto-soft-delete-request-criteria")
	operations.AddStringOperationIfNecessary(&ops, plan.SoftDeleteRetentionTime, state.SoftDeleteRetentionTime, "soft-delete-retention-time")
	operations.AddInt64OperationIfNecessary(&ops, plan.SoftDeleteRetainNumberOfEntries, state.SoftDeleteRetainNumberOfEntries, "soft-delete-retain-number-of-entries")
	return ops
}

// Create a soft-delete-policy soft-delete-policy
func (r *softDeletePolicyResource) CreateSoftDeletePolicy(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan softDeletePolicyResourceModel) (*softDeletePolicyResourceModel, error) {
	addRequest := client.NewAddSoftDeletePolicyRequest(plan.Name.ValueString())
	addOptionalSoftDeletePolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.SoftDeletePolicyApi.AddSoftDeletePolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddSoftDeletePolicyRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.SoftDeletePolicyApi.AddSoftDeletePolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Soft Delete Policy", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state softDeletePolicyResourceModel
	readSoftDeletePolicyResponse(ctx, addResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *softDeletePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan softDeletePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateSoftDeletePolicy(ctx, req, resp, plan)
	if err != nil {
		return
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
func (r *defaultSoftDeletePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan softDeletePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SoftDeletePolicyApi.GetSoftDeletePolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Soft Delete Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state softDeletePolicyResourceModel
	readSoftDeletePolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.SoftDeletePolicyApi.UpdateSoftDeletePolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createSoftDeletePolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.SoftDeletePolicyApi.UpdateSoftDeletePolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Soft Delete Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSoftDeletePolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *softDeletePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSoftDeletePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSoftDeletePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSoftDeletePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSoftDeletePolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state softDeletePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.SoftDeletePolicyApi.GetSoftDeletePolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Soft Delete Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSoftDeletePolicyResponse(ctx, readResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *softDeletePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSoftDeletePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSoftDeletePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSoftDeletePolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSoftDeletePolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan softDeletePolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state softDeletePolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.SoftDeletePolicyApi.UpdateSoftDeletePolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createSoftDeletePolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.SoftDeletePolicyApi.UpdateSoftDeletePolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Soft Delete Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSoftDeletePolicyResponse(ctx, updateResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSoftDeletePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *softDeletePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state softDeletePolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.SoftDeletePolicyApi.DeleteSoftDeletePolicyExecute(r.apiClient.SoftDeletePolicyApi.DeleteSoftDeletePolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Soft Delete Policy", err, httpResp)
		return
	}
}

func (r *softDeletePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSoftDeletePolicy(ctx, req, resp)
}

func (r *defaultSoftDeletePolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSoftDeletePolicy(ctx, req, resp)
}

func importSoftDeletePolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
