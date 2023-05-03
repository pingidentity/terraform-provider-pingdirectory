package logfilerotationlistener

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
	_ resource.Resource                = &summarizeLogFileRotationListenerResource{}
	_ resource.ResourceWithConfigure   = &summarizeLogFileRotationListenerResource{}
	_ resource.ResourceWithImportState = &summarizeLogFileRotationListenerResource{}
	_ resource.Resource                = &defaultSummarizeLogFileRotationListenerResource{}
	_ resource.ResourceWithConfigure   = &defaultSummarizeLogFileRotationListenerResource{}
	_ resource.ResourceWithImportState = &defaultSummarizeLogFileRotationListenerResource{}
)

// Create a Summarize Log File Rotation Listener resource
func NewSummarizeLogFileRotationListenerResource() resource.Resource {
	return &summarizeLogFileRotationListenerResource{}
}

func NewDefaultSummarizeLogFileRotationListenerResource() resource.Resource {
	return &defaultSummarizeLogFileRotationListenerResource{}
}

// summarizeLogFileRotationListenerResource is the resource implementation.
type summarizeLogFileRotationListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultSummarizeLogFileRotationListenerResource is the resource implementation.
type defaultSummarizeLogFileRotationListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *summarizeLogFileRotationListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_summarize_log_file_rotation_listener"
}

func (r *defaultSummarizeLogFileRotationListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_summarize_log_file_rotation_listener"
}

// Configure adds the provider configured client to the resource.
func (r *summarizeLogFileRotationListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultSummarizeLogFileRotationListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type summarizeLogFileRotationListenerResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	OutputDirectory types.String `tfsdk:"output_directory"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *summarizeLogFileRotationListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	summarizeLogFileRotationListenerSchema(ctx, req, resp, false)
}

func (r *defaultSummarizeLogFileRotationListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	summarizeLogFileRotationListenerSchema(ctx, req, resp, true)
}

func summarizeLogFileRotationListenerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Summarize Log File Rotation Listener.",
		Attributes: map[string]schema.Attribute{
			"output_directory": schema.StringAttribute{
				Description: "The path to the directory in which the summarize-access-log output should be written. If no value is provided, the output file will be written into the same directory as the rotated log file.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log File Rotation Listener",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Log File Rotation Listener is enabled for use.",
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
func addOptionalSummarizeLogFileRotationListenerFields(ctx context.Context, addRequest *client.AddSummarizeLogFileRotationListenerRequest, plan summarizeLogFileRotationListenerResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.OutputDirectory) {
		addRequest.OutputDirectory = plan.OutputDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a SummarizeLogFileRotationListenerResponse object into the model struct
func readSummarizeLogFileRotationListenerResponse(ctx context.Context, r *client.SummarizeLogFileRotationListenerResponse, state *summarizeLogFileRotationListenerResourceModel, expectedValues *summarizeLogFileRotationListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.OutputDirectory = internaltypes.StringTypeOrNil(r.OutputDirectory, internaltypes.IsEmptyString(expectedValues.OutputDirectory))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createSummarizeLogFileRotationListenerOperations(plan summarizeLogFileRotationListenerResourceModel, state summarizeLogFileRotationListenerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.OutputDirectory, state.OutputDirectory, "output-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *summarizeLogFileRotationListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan summarizeLogFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddSummarizeLogFileRotationListenerRequest(plan.Id.ValueString(),
		[]client.EnumsummarizeLogFileRotationListenerSchemaUrn{client.ENUMSUMMARIZELOGFILEROTATIONLISTENERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FILE_ROTATION_LISTENERSUMMARIZE},
		plan.Enabled.ValueBool())
	addOptionalSummarizeLogFileRotationListenerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFileRotationListenerApi.AddLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFileRotationListenerRequest(
		client.AddSummarizeLogFileRotationListenerRequestAsAddLogFileRotationListenerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.AddLogFileRotationListenerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Summarize Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state summarizeLogFileRotationListenerResourceModel
	readSummarizeLogFileRotationListenerResponse(ctx, addResponse.SummarizeLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultSummarizeLogFileRotationListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan summarizeLogFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Summarize Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state summarizeLogFileRotationListenerResourceModel
	readSummarizeLogFileRotationListenerResponse(ctx, readResponse.SummarizeLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListener(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createSummarizeLogFileRotationListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Summarize Log File Rotation Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSummarizeLogFileRotationListenerResponse(ctx, updateResponse.SummarizeLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *summarizeLogFileRotationListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSummarizeLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSummarizeLogFileRotationListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readSummarizeLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readSummarizeLogFileRotationListener(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state summarizeLogFileRotationListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogFileRotationListenerApi.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Summarize Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSummarizeLogFileRotationListenerResponse(ctx, readResponse.SummarizeLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *summarizeLogFileRotationListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSummarizeLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultSummarizeLogFileRotationListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateSummarizeLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateSummarizeLogFileRotationListener(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan summarizeLogFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state summarizeLogFileRotationListenerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createSummarizeLogFileRotationListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Summarize Log File Rotation Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readSummarizeLogFileRotationListenerResponse(ctx, updateResponse.SummarizeLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultSummarizeLogFileRotationListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *summarizeLogFileRotationListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state summarizeLogFileRotationListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogFileRotationListenerApi.DeleteLogFileRotationListenerExecute(r.apiClient.LogFileRotationListenerApi.DeleteLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Summarize Log File Rotation Listener", err, httpResp)
		return
	}
}

func (r *summarizeLogFileRotationListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSummarizeLogFileRotationListener(ctx, req, resp)
}

func (r *defaultSummarizeLogFileRotationListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importSummarizeLogFileRotationListener(ctx, req, resp)
}

func importSummarizeLogFileRotationListener(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
