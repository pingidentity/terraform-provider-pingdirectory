package keymanagerprovider

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
	_ resource.Resource                = &pkcs11KeyManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &pkcs11KeyManagerProviderResource{}
	_ resource.ResourceWithImportState = &pkcs11KeyManagerProviderResource{}
	_ resource.Resource                = &defaultPkcs11KeyManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultPkcs11KeyManagerProviderResource{}
	_ resource.ResourceWithImportState = &defaultPkcs11KeyManagerProviderResource{}
)

// Create a Pkcs11 Key Manager Provider resource
func NewPkcs11KeyManagerProviderResource() resource.Resource {
	return &pkcs11KeyManagerProviderResource{}
}

func NewDefaultPkcs11KeyManagerProviderResource() resource.Resource {
	return &defaultPkcs11KeyManagerProviderResource{}
}

// pkcs11KeyManagerProviderResource is the resource implementation.
type pkcs11KeyManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPkcs11KeyManagerProviderResource is the resource implementation.
type defaultPkcs11KeyManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pkcs11KeyManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pkcs11_key_manager_provider"
}

func (r *defaultPkcs11KeyManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_pkcs11_key_manager_provider"
}

// Configure adds the provider configured client to the resource.
func (r *pkcs11KeyManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultPkcs11KeyManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type pkcs11KeyManagerProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Pkcs11ProviderClass             types.String `tfsdk:"pkcs11_provider_class"`
	Pkcs11ProviderConfigurationFile types.String `tfsdk:"pkcs11_provider_configuration_file"`
	Pkcs11KeyStoreType              types.String `tfsdk:"pkcs11_key_store_type"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinPassphraseProvider   types.String `tfsdk:"key_store_pin_passphrase_provider"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *pkcs11KeyManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pkcs11KeyManagerProviderSchema(ctx, req, resp, false)
}

func (r *defaultPkcs11KeyManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pkcs11KeyManagerProviderSchema(ctx, req, resp, true)
}

func pkcs11KeyManagerProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Pkcs11 Key Manager Provider.",
		Attributes: map[string]schema.Attribute{
			"pkcs11_provider_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java security provider class that implements support for interacting with PKCS #11 tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pkcs11_provider_configuration_file": schema.StringAttribute{
				Description: "The path to the file to use to configure the security provider that implements support for interacting with PKCS #11 tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"pkcs11_key_store_type": schema.StringAttribute{
				Description: "The key store type to use when obtaining an instance of a key store for interacting with a PKCS #11 token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_store_pin": schema.StringAttribute{
				Description: "Specifies the PIN needed to access the PKCS11 Key Manager Provider.",
				Optional:    true,
				Sensitive:   true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				Optional:    true,
			},
			"key_store_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
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
func addOptionalPkcs11KeyManagerProviderFields(ctx context.Context, addRequest *client.AddPkcs11KeyManagerProviderRequest, plan pkcs11KeyManagerProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderClass) {
		addRequest.Pkcs11ProviderClass = plan.Pkcs11ProviderClass.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderConfigurationFile) {
		addRequest.Pkcs11ProviderConfigurationFile = plan.Pkcs11ProviderConfigurationFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11KeyStoreType) {
		addRequest.Pkcs11KeyStoreType = plan.Pkcs11KeyStoreType.ValueStringPointer()
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
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a Pkcs11KeyManagerProviderResponse object into the model struct
func readPkcs11KeyManagerProviderResponse(ctx context.Context, r *client.Pkcs11KeyManagerProviderResponse, state *pkcs11KeyManagerProviderResourceModel, expectedValues *pkcs11KeyManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Pkcs11ProviderClass = internaltypes.StringTypeOrNil(r.Pkcs11ProviderClass, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderClass))
	state.Pkcs11ProviderConfigurationFile = internaltypes.StringTypeOrNil(r.Pkcs11ProviderConfigurationFile, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderConfigurationFile))
	state.Pkcs11KeyStoreType = internaltypes.StringTypeOrNil(r.Pkcs11KeyStoreType, internaltypes.IsEmptyString(expectedValues.Pkcs11KeyStoreType))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.KeyStorePin = expectedValues.KeyStorePin
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, internaltypes.IsEmptyString(expectedValues.KeyStorePinFile))
	state.KeyStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.KeyStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.KeyStorePinPassphraseProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPkcs11KeyManagerProviderOperations(plan pkcs11KeyManagerProviderResourceModel, state pkcs11KeyManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderClass, state.Pkcs11ProviderClass, "pkcs11-provider-class")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderConfigurationFile, state.Pkcs11ProviderConfigurationFile, "pkcs11-provider-configuration-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11KeyStoreType, state.Pkcs11KeyStoreType, "pkcs11-key-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePin, state.KeyStorePin, "key-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinFile, state.KeyStorePinFile, "key-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinPassphraseProvider, state.KeyStorePinPassphraseProvider, "key-store-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *pkcs11KeyManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pkcs11KeyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPkcs11KeyManagerProviderRequest(plan.Id.ValueString(),
		[]client.Enumpkcs11KeyManagerProviderSchemaUrn{client.ENUMPKCS11KEYMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0KEY_MANAGER_PROVIDERPKCS11},
		plan.Enabled.ValueBool())
	addOptionalPkcs11KeyManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.KeyManagerProviderApi.AddKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyManagerProviderRequest(
		client.AddPkcs11KeyManagerProviderRequestAsAddKeyManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.AddKeyManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pkcs11 Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pkcs11KeyManagerProviderResourceModel
	readPkcs11KeyManagerProviderResponse(ctx, addResponse.Pkcs11KeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultPkcs11KeyManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pkcs11KeyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pkcs11 Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state pkcs11KeyManagerProviderResourceModel
	readPkcs11KeyManagerProviderResponse(ctx, readResponse.Pkcs11KeyManagerProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.KeyManagerProviderApi.UpdateKeyManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPkcs11KeyManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.UpdateKeyManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pkcs11 Key Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPkcs11KeyManagerProviderResponse(ctx, updateResponse.Pkcs11KeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *pkcs11KeyManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPkcs11KeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPkcs11KeyManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPkcs11KeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPkcs11KeyManagerProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state pkcs11KeyManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pkcs11 Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPkcs11KeyManagerProviderResponse(ctx, readResponse.Pkcs11KeyManagerProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pkcs11KeyManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePkcs11KeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPkcs11KeyManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePkcs11KeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePkcs11KeyManagerProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan pkcs11KeyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pkcs11KeyManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.KeyManagerProviderApi.UpdateKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPkcs11KeyManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.KeyManagerProviderApi.UpdateKeyManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pkcs11 Key Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPkcs11KeyManagerProviderResponse(ctx, updateResponse.Pkcs11KeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPkcs11KeyManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pkcs11KeyManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pkcs11KeyManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.KeyManagerProviderApi.DeleteKeyManagerProviderExecute(r.apiClient.KeyManagerProviderApi.DeleteKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Pkcs11 Key Manager Provider", err, httpResp)
		return
	}
}

func (r *pkcs11KeyManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPkcs11KeyManagerProvider(ctx, req, resp)
}

func (r *defaultPkcs11KeyManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPkcs11KeyManagerProvider(ctx, req, resp)
}

func importPkcs11KeyManagerProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
