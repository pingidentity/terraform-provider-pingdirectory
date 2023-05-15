package keymanagerprovider

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
	_ resource.Resource                = &fileBasedKeyManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &fileBasedKeyManagerProviderResource{}
	_ resource.ResourceWithImportState = &fileBasedKeyManagerProviderResource{}
	_ resource.Resource                = &defaultFileBasedKeyManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultFileBasedKeyManagerProviderResource{}
	_ resource.ResourceWithImportState = &defaultFileBasedKeyManagerProviderResource{}
)

// Create a File Based Key Manager Provider resource
func NewFileBasedKeyManagerProviderResource() resource.Resource {
	return &fileBasedKeyManagerProviderResource{}
}

func NewDefaultFileBasedKeyManagerProviderResource() resource.Resource {
	return &defaultFileBasedKeyManagerProviderResource{}
}

// fileBasedKeyManagerProviderResource is the resource implementation.
type fileBasedKeyManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultFileBasedKeyManagerProviderResource is the resource implementation.
type defaultFileBasedKeyManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *fileBasedKeyManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file_based_key_manager_provider"
}

func (r *defaultFileBasedKeyManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_file_based_key_manager_provider"
}

// Configure adds the provider configured client to the resource.
func (r *fileBasedKeyManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultFileBasedKeyManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type fileBasedKeyManagerProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	KeyStoreFile                    types.String `tfsdk:"key_store_file"`
	KeyStoreType                    types.String `tfsdk:"key_store_type"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinPassphraseProvider   types.String `tfsdk:"key_store_pin_passphrase_provider"`
	PrivateKeyPin                   types.String `tfsdk:"private_key_pin"`
	PrivateKeyPinFile               types.String `tfsdk:"private_key_pin_file"`
	PrivateKeyPinPassphraseProvider types.String `tfsdk:"private_key_pin_passphrase_provider"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *fileBasedKeyManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileBasedKeyManagerProviderSchema(ctx, req, resp, false)
}

func (r *defaultFileBasedKeyManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	fileBasedKeyManagerProviderSchema(ctx, req, resp, true)
}

func fileBasedKeyManagerProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a File Based Key Manager Provider.",
		Attributes: map[string]schema.Attribute{
			"key_store_file": schema.StringAttribute{
				Description: "Specifies the path to the file that contains the private key information. This may be an absolute path, or a path that is relative to the Directory Server instance root.",
				Required:    true,
			},
			"key_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the key store file.",
				Optional:    true,
			},
			"key_store_pin": schema.StringAttribute{
				Description: "Specifies the PIN needed to access the File Based Key Manager Provider.",
				Optional:    true,
				Sensitive:   true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider.",
				Optional:    true,
			},
			"key_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider.",
				Optional:    true,
			},
			"private_key_pin": schema.StringAttribute{
				Description: "Specifies the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Optional:    true,
				Sensitive:   true,
			},
			"private_key_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Optional:    true,
			},
			"private_key_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Key Manager Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Key Manager Provider is enabled for use.",
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
func addOptionalFileBasedKeyManagerProviderFields(ctx context.Context, addRequest *client.AddFileBasedKeyManagerProviderRequest, plan fileBasedKeyManagerProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStoreType) {
		addRequest.KeyStoreType = plan.KeyStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStorePin) {
		addRequest.KeyStorePin = plan.KeyStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStorePinFile) {
		addRequest.KeyStorePinFile = plan.KeyStorePinFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStorePinPassphraseProvider) {
		addRequest.KeyStorePinPassphraseProvider = plan.KeyStorePinPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKeyPin) {
		addRequest.PrivateKeyPin = plan.PrivateKeyPin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKeyPinFile) {
		addRequest.PrivateKeyPinFile = plan.PrivateKeyPinFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKeyPinPassphraseProvider) {
		addRequest.PrivateKeyPinPassphraseProvider = plan.PrivateKeyPinPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a FileBasedKeyManagerProviderResponse object into the model struct
func readFileBasedKeyManagerProviderResponse(ctx context.Context, r *client.FileBasedKeyManagerProviderResponse, state *fileBasedKeyManagerProviderResourceModel, expectedValues *fileBasedKeyManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.KeyStoreFile = types.StringValue(r.KeyStoreFile)
	state.KeyStoreType = internaltypes.StringTypeOrNil(r.KeyStoreType, internaltypes.IsEmptyString(expectedValues.KeyStoreType))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.KeyStorePin = expectedValues.KeyStorePin
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, internaltypes.IsEmptyString(expectedValues.KeyStorePinFile))
	state.KeyStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.KeyStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.KeyStorePinPassphraseProvider))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.PrivateKeyPin = expectedValues.PrivateKeyPin
	state.PrivateKeyPinFile = internaltypes.StringTypeOrNil(r.PrivateKeyPinFile, internaltypes.IsEmptyString(expectedValues.PrivateKeyPinFile))
	state.PrivateKeyPinPassphraseProvider = internaltypes.StringTypeOrNil(r.PrivateKeyPinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.PrivateKeyPinPassphraseProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createFileBasedKeyManagerProviderOperations(plan fileBasedKeyManagerProviderResourceModel, state fileBasedKeyManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStoreFile, state.KeyStoreFile, "key-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStoreType, state.KeyStoreType, "key-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePin, state.KeyStorePin, "key-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinFile, state.KeyStorePinFile, "key-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinPassphraseProvider, state.KeyStorePinPassphraseProvider, "key-store-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKeyPin, state.PrivateKeyPin, "private-key-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKeyPinFile, state.PrivateKeyPinFile, "private-key-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKeyPinPassphraseProvider, state.PrivateKeyPinPassphraseProvider, "private-key-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *fileBasedKeyManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedKeyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddFileBasedKeyManagerProviderRequest(plan.Id.ValueString(),
		[]client.EnumfileBasedKeyManagerProviderSchemaUrn{client.ENUMFILEBASEDKEYMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0KEY_MANAGER_PROVIDERFILE_BASED},
		plan.KeyStoreFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedKeyManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.KeyManagerProviderApi.AddKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyManagerProviderRequest(
		client.AddFileBasedKeyManagerProviderRequestAsAddKeyManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.AddKeyManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the File Based Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state fileBasedKeyManagerProviderResourceModel
	readFileBasedKeyManagerProviderResponse(ctx, addResponse.FileBasedKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultFileBasedKeyManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan fileBasedKeyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Based Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state fileBasedKeyManagerProviderResourceModel
	readFileBasedKeyManagerProviderResponse(ctx, readResponse.FileBasedKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.KeyManagerProviderApi.UpdateKeyManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createFileBasedKeyManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.UpdateKeyManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Based Key Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileBasedKeyManagerProviderResponse(ctx, updateResponse.FileBasedKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *fileBasedKeyManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileBasedKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileBasedKeyManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readFileBasedKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readFileBasedKeyManagerProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state fileBasedKeyManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the File Based Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readFileBasedKeyManagerProviderResponse(ctx, readResponse.FileBasedKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *fileBasedKeyManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileBasedKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultFileBasedKeyManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateFileBasedKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateFileBasedKeyManagerProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan fileBasedKeyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state fileBasedKeyManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.KeyManagerProviderApi.UpdateKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createFileBasedKeyManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.KeyManagerProviderApi.UpdateKeyManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the File Based Key Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readFileBasedKeyManagerProviderResponse(ctx, updateResponse.FileBasedKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultFileBasedKeyManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *fileBasedKeyManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state fileBasedKeyManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.KeyManagerProviderApi.DeleteKeyManagerProviderExecute(r.apiClient.KeyManagerProviderApi.DeleteKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the File Based Key Manager Provider", err, httpResp)
		return
	}
}

func (r *fileBasedKeyManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileBasedKeyManagerProvider(ctx, req, resp)
}

func (r *defaultFileBasedKeyManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importFileBasedKeyManagerProvider(ctx, req, resp)
}

func importFileBasedKeyManagerProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
