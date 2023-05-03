package logretentionpolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &fileCountLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &fileCountLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &fileCountLogRetentionPolicyResource{}
	_ resource.Resource                = &defaultFileCountLogRetentionPolicyResource{}
	_ resource.ResourceWithConfigure   = &defaultFileCountLogRetentionPolicyResource{}
	_ resource.ResourceWithImportState = &defaultFileCountLogRetentionPolicyResource{}
)

// Create a File Count Log Retention Policy resource
func NewFileCountLogRetentionPolicyResource() resource.Resource {
	return &fileCountLogRetentionPolicyResource{}
}

func NewDefaultFileCountLogRetentionPolicyResource() resource.Resource {
	return &defaultFileCountLogRetentionPolicyResource{}
}

// fileCountLogRetentionPolicyResource is the resource implementation.
type fileCountLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFileCountLogRetentionPolicyResource is the resource implementation.
type defaultFileCountLogRetentionPolicyResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fileCountLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_count_log_retention_policy"
}

func (r *defaultFileCountLogRetentionPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_file_count_log_retention_policy"
}

// Configure adds the provider configured client to the resource.
func (r *fileCountLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFileCountLogRetentionPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type fileCountLogRetentionPolicyResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	NumberOfFiles   types.Int64  `tfsdk:"number_of_files"`
	Description     types.String `tfsdk:"description"`
}

// GetSchema defines the schema for the resource.
func (r *fileCountLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileCountLogRetentionPolicySchema(ctx, req, resp, false)
}

func (r *defaultFileCountLogRetentionPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileCountLogRetentionPolicySchema(ctx, req, resp, true)
}

func fileCountLogRetentionPolicySchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a File Count Log Retention Policy.",
		Attributes: map[string]schema.Attribute{
			"number_of_files": schema.Int64Attribute{
				Description: "Specifies the number of archived log files to retain before the oldest ones are cleaned.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Retention Policy",
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
func addOptionalFileCountLogRetentionPolicyFields(ctx context.Context, addRequest *client.AddFileCountLogRetentionPolicyRequest, plan fileCountLogRetentionPolicyResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a FileCountLogRetentionPolicyResponse object into the model struct
func readFileCountLogRetentionPolicyResponse(ctx context.Context, r *client.FileCountLogRetentionPolicyResponse, state *fileCountLogRetentionPolicyResourceModel, expectedValues *fileCountLogRetentionPolicyResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.NumberOfFiles = types.Int64Value(r.NumberOfFiles)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFileCountLogRetentionPolicyOperations(plan fileCountLogRetentionPolicyResourceModel, state fileCountLogRetentionPolicyResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddInt64OperationIfNecessary(&ops, plan.NumberOfFiles, state.NumberOfFiles, "number-of-files")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	return ops
}

// Create a new resource
func (r *fileCountLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileCountLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFileCountLogRetentionPolicyRequest(plan.Id.ValueString(),
		[]client.EnumfileCountLogRetentionPolicySchemaUrn{client.ENUMFILECOUNTLOGRETENTIONPOLICYSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_RETENTION_POLICYFILE_COUNT},
		plan.NumberOfFiles.ValueInt64())
	addOptionalFileCountLogRetentionPolicyFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogRetentionPolicyApi.AddLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogRetentionPolicyRequest(
		client.AddFileCountLogRetentionPolicyRequestAsAddLogRetentionPolicyRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.AddLogRetentionPolicyExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the File Count Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fileCountLogRetentionPolicyResourceModel
	readFileCountLogRetentionPolicyResponse(ctx, addResponse.FileCountLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFileCountLogRetentionPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileCountLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Count Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fileCountLogRetentionPolicyResourceModel
	readFileCountLogRetentionPolicyResponse(ctx, readResponse.FileCountLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicy(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFileCountLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Count Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileCountLogRetentionPolicyResponse(ctx, updateResponse.FileCountLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fileCountLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileCountLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileCountLogRetentionPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileCountLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFileCountLogRetentionPolicy(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fileCountLogRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogRetentionPolicyApi.GetLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Count Log Retention Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFileCountLogRetentionPolicyResponse(ctx, readResponse.FileCountLogRetentionPolicyResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fileCountLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileCountLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileCountLogRetentionPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileCountLogRetentionPolicy(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFileCountLogRetentionPolicy(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fileCountLogRetentionPolicyResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileCountLogRetentionPolicyResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFileCountLogRetentionPolicyOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogRetentionPolicyApi.UpdateLogRetentionPolicyExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Count Log Retention Policy", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileCountLogRetentionPolicyResponse(ctx, updateResponse.FileCountLogRetentionPolicyResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFileCountLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fileCountLogRetentionPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileCountLogRetentionPolicyResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogRetentionPolicyApi.DeleteLogRetentionPolicyExecute(r.apiClient.LogRetentionPolicyApi.DeleteLogRetentionPolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the File Count Log Retention Policy", err, httpResp)
		return
	}
}

func (r *fileCountLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileCountLogRetentionPolicy(ctx, req, resp)
}

func (r *defaultFileCountLogRetentionPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileCountLogRetentionPolicy(ctx, req, resp)
}

func importFileCountLogRetentionPolicy(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
