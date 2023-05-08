package cipherstreamprovider

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
	_ resource.Resource                = &fileBasedCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &fileBasedCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &fileBasedCipherStreamProviderResource{}
	_ resource.Resource                = &defaultFileBasedCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultFileBasedCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultFileBasedCipherStreamProviderResource{}
)

// Create a File Based Cipher Stream Provider resource
func NewFileBasedCipherStreamProviderResource() resource.Resource {
	return &fileBasedCipherStreamProviderResource{}
}

func NewDefaultFileBasedCipherStreamProviderResource() resource.Resource {
	return &defaultFileBasedCipherStreamProviderResource{}
}

// fileBasedCipherStreamProviderResource is the resource implementation.
type fileBasedCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFileBasedCipherStreamProviderResource is the resource implementation.
type defaultFileBasedCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fileBasedCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_based_cipher_stream_provider"
}

func (r *defaultFileBasedCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_file_based_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *fileBasedCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFileBasedCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type fileBasedCipherStreamProviderResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	LastUpdated         types.String `tfsdk:"last_updated"`
	Notifications       types.Set    `tfsdk:"notifications"`
	RequiredActions     types.Set    `tfsdk:"required_actions"`
	PasswordFile        types.String `tfsdk:"password_file"`
	WaitForPasswordFile types.Bool   `tfsdk:"wait_for_password_file"`
	Description         types.String `tfsdk:"description"`
	Enabled             types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *fileBasedCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileBasedCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultFileBasedCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileBasedCipherStreamProviderSchema(ctx, req, resp, true)
}

func fileBasedCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a File Based Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"password_file": schema.StringAttribute{
				Description: "The path to the file containing the password to use when generating ciphers.",
				Required:    true,
			},
			"wait_for_password_file": schema.BoolAttribute{
				Description: "Indicates whether the server should wait for the password file to become available if it does not exist.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Cipher Stream Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Cipher Stream Provider is enabled for use in the Directory Server.",
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
func addOptionalFileBasedCipherStreamProviderFields(ctx context.Context, addRequest *client.AddFileBasedCipherStreamProviderRequest, plan fileBasedCipherStreamProviderResourceModel) {
	if internaltypes.IsDefined(plan.WaitForPasswordFile) {
		addRequest.WaitForPasswordFile = plan.WaitForPasswordFile.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a FileBasedCipherStreamProviderResponse object into the model struct
func readFileBasedCipherStreamProviderResponse(ctx context.Context, r *client.FileBasedCipherStreamProviderResponse, state *fileBasedCipherStreamProviderResourceModel, expectedValues *fileBasedCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.PasswordFile = types.StringValue(r.PasswordFile)
	state.WaitForPasswordFile = internaltypes.BoolTypeOrNil(r.WaitForPasswordFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFileBasedCipherStreamProviderOperations(plan fileBasedCipherStreamProviderResourceModel, state fileBasedCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordFile, state.PasswordFile, "password-file")
	operations.AddBoolOperationIfNecessary(&ops, plan.WaitForPasswordFile, state.WaitForPasswordFile, "wait-for-password-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *fileBasedCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFileBasedCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumfileBasedCipherStreamProviderSchemaUrn{client.ENUMFILEBASEDCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERFILE_BASED},
		plan.PasswordFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddFileBasedCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the File Based Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fileBasedCipherStreamProviderResourceModel
	readFileBasedCipherStreamProviderResponse(ctx, addResponse.FileBasedCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFileBasedCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Based Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fileBasedCipherStreamProviderResourceModel
	readFileBasedCipherStreamProviderResponse(ctx, readResponse.FileBasedCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFileBasedCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Based Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileBasedCipherStreamProviderResponse(ctx, updateResponse.FileBasedCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fileBasedCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileBasedCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileBasedCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileBasedCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFileBasedCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fileBasedCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Based Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFileBasedCipherStreamProviderResponse(ctx, readResponse.FileBasedCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fileBasedCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileBasedCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileBasedCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileBasedCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFileBasedCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fileBasedCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileBasedCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFileBasedCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Based Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileBasedCipherStreamProviderResponse(ctx, updateResponse.FileBasedCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFileBasedCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fileBasedCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileBasedCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the File Based Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *fileBasedCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileBasedCipherStreamProvider(ctx, req, resp)
}

func (r *defaultFileBasedCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileBasedCipherStreamProvider(ctx, req, resp)
}

func importFileBasedCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
