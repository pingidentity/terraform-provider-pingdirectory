package logfilerotationlistener

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &copyLogFileRotationListenerResource{}
	_ resource.ResourceWithConfigure   = &copyLogFileRotationListenerResource{}
	_ resource.ResourceWithImportState = &copyLogFileRotationListenerResource{}
	_ resource.Resource                = &defaultCopyLogFileRotationListenerResource{}
	_ resource.ResourceWithConfigure   = &defaultCopyLogFileRotationListenerResource{}
	_ resource.ResourceWithImportState = &defaultCopyLogFileRotationListenerResource{}
)

// Create a Copy Log File Rotation Listener resource
func NewCopyLogFileRotationListenerResource() resource.Resource {
	return &copyLogFileRotationListenerResource{}
}

func NewDefaultCopyLogFileRotationListenerResource() resource.Resource {
	return &defaultCopyLogFileRotationListenerResource{}
}

// copyLogFileRotationListenerResource is the resource implementation.
type copyLogFileRotationListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultCopyLogFileRotationListenerResource is the resource implementation.
type defaultCopyLogFileRotationListenerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *copyLogFileRotationListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_copy_log_file_rotation_listener"
}

func (r *defaultCopyLogFileRotationListenerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_copy_log_file_rotation_listener"
}

// Configure adds the provider configured client to the resource.
func (r *copyLogFileRotationListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultCopyLogFileRotationListenerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type copyLogFileRotationListenerResourceModel struct {
	Id              types.String `tfsdk:"id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	Notifications   types.Set    `tfsdk:"notifications"`
	RequiredActions types.Set    `tfsdk:"required_actions"`
	CopyToDirectory types.String `tfsdk:"copy_to_directory"`
	CompressOnCopy  types.Bool   `tfsdk:"compress_on_copy"`
	Description     types.String `tfsdk:"description"`
	Enabled         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *copyLogFileRotationListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	copyLogFileRotationListenerSchema(ctx, req, resp, false)
}

func (r *defaultCopyLogFileRotationListenerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	copyLogFileRotationListenerSchema(ctx, req, resp, true)
}

func copyLogFileRotationListenerSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Copy Log File Rotation Listener.",
		Attributes: map[string]schema.Attribute{
			"copy_to_directory": schema.StringAttribute{
				Description: "The path to the directory to which log files should be copied. It must be different from the directory to which the log file is originally written, and administrators should ensure that the filesystem has sufficient space to hold files as they are copied.",
				Required:    true,
			},
			"compress_on_copy": schema.BoolAttribute{
				Description: "Indicates whether the file should be gzip-compressed as it is copied into the destination directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
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
func addOptionalCopyLogFileRotationListenerFields(ctx context.Context, addRequest *client.AddCopyLogFileRotationListenerRequest, plan copyLogFileRotationListenerResourceModel) {
	if internaltypes.IsDefined(plan.CompressOnCopy) {
		addRequest.CompressOnCopy = plan.CompressOnCopy.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a CopyLogFileRotationListenerResponse object into the model struct
func readCopyLogFileRotationListenerResponse(ctx context.Context, r *client.CopyLogFileRotationListenerResponse, state *copyLogFileRotationListenerResourceModel, expectedValues *copyLogFileRotationListenerResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.CopyToDirectory = types.StringValue(r.CopyToDirectory)
	state.CompressOnCopy = internaltypes.BoolTypeOrNil(r.CompressOnCopy)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createCopyLogFileRotationListenerOperations(plan copyLogFileRotationListenerResourceModel, state copyLogFileRotationListenerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.CopyToDirectory, state.CopyToDirectory, "copy-to-directory")
	operations.AddBoolOperationIfNecessary(&ops, plan.CompressOnCopy, state.CompressOnCopy, "compress-on-copy")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *copyLogFileRotationListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan copyLogFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddCopyLogFileRotationListenerRequest(plan.Id.ValueString(),
		[]client.EnumcopyLogFileRotationListenerSchemaUrn{client.ENUMCOPYLOGFILEROTATIONLISTENERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0LOG_FILE_ROTATION_LISTENERCOPY},
		plan.CopyToDirectory.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalCopyLogFileRotationListenerFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.LogFileRotationListenerApi.AddLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddLogFileRotationListenerRequest(
		client.AddCopyLogFileRotationListenerRequestAsAddLogFileRotationListenerRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.AddLogFileRotationListenerExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Copy Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state copyLogFileRotationListenerResourceModel
	readCopyLogFileRotationListenerResponse(ctx, addResponse.CopyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultCopyLogFileRotationListenerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan copyLogFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Copy Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state copyLogFileRotationListenerResourceModel
	readCopyLogFileRotationListenerResponse(ctx, readResponse.CopyLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListener(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createCopyLogFileRotationListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Copy Log File Rotation Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCopyLogFileRotationListenerResponse(ctx, updateResponse.CopyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *copyLogFileRotationListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCopyLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCopyLogFileRotationListenerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readCopyLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readCopyLogFileRotationListener(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state copyLogFileRotationListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.LogFileRotationListenerApi.GetLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Copy Log File Rotation Listener", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCopyLogFileRotationListenerResponse(ctx, readResponse.CopyLogFileRotationListenerResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *copyLogFileRotationListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCopyLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultCopyLogFileRotationListenerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateCopyLogFileRotationListener(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateCopyLogFileRotationListener(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan copyLogFileRotationListenerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state copyLogFileRotationListenerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createCopyLogFileRotationListenerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.LogFileRotationListenerApi.UpdateLogFileRotationListenerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Copy Log File Rotation Listener", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCopyLogFileRotationListenerResponse(ctx, updateResponse.CopyLogFileRotationListenerResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultCopyLogFileRotationListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *copyLogFileRotationListenerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state copyLogFileRotationListenerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.LogFileRotationListenerApi.DeleteLogFileRotationListenerExecute(r.apiClient.LogFileRotationListenerApi.DeleteLogFileRotationListener(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Copy Log File Rotation Listener", err, httpResp)
		return
	}
}

func (r *copyLogFileRotationListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCopyLogFileRotationListener(ctx, req, resp)
}

func (r *defaultCopyLogFileRotationListenerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importCopyLogFileRotationListener(ctx, req, resp)
}

func importCopyLogFileRotationListener(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
