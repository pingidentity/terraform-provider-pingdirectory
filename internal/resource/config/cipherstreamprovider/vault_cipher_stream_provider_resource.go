package cipherstreamprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
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
	_ resource.Resource                = &vaultCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &vaultCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &vaultCipherStreamProviderResource{}
	_ resource.Resource                = &defaultVaultCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultVaultCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultVaultCipherStreamProviderResource{}
)

// Create a Vault Cipher Stream Provider resource
func NewVaultCipherStreamProviderResource() resource.Resource {
	return &vaultCipherStreamProviderResource{}
}

func NewDefaultVaultCipherStreamProviderResource() resource.Resource {
	return &defaultVaultCipherStreamProviderResource{}
}

// vaultCipherStreamProviderResource is the resource implementation.
type vaultCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultVaultCipherStreamProviderResource is the resource implementation.
type defaultVaultCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *vaultCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_cipher_stream_provider"
}

func (r *defaultVaultCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_vault_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *vaultCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultVaultCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type vaultCipherStreamProviderResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	LastUpdated                 types.String `tfsdk:"last_updated"`
	Notifications               types.Set    `tfsdk:"notifications"`
	RequiredActions             types.Set    `tfsdk:"required_actions"`
	VaultExternalServer         types.String `tfsdk:"vault_external_server"`
	VaultServerBaseURI          types.Set    `tfsdk:"vault_server_base_uri"`
	VaultAuthenticationMethod   types.String `tfsdk:"vault_authentication_method"`
	VaultSecretPath             types.String `tfsdk:"vault_secret_path"`
	VaultSecretFieldName        types.String `tfsdk:"vault_secret_field_name"`
	VaultEncryptionMetadataFile types.String `tfsdk:"vault_encryption_metadata_file"`
	TrustStoreFile              types.String `tfsdk:"trust_store_file"`
	TrustStorePin               types.String `tfsdk:"trust_store_pin"`
	TrustStoreType              types.String `tfsdk:"trust_store_type"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *vaultCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	vaultCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultVaultCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	vaultCipherStreamProviderSchema(ctx, req, resp, true)
}

func vaultCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Vault Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault server.",
				Optional:    true,
			},
			"vault_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the Vault server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://vault.example.com:8200/\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"vault_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Vault server.",
				Optional:    true,
			},
			"vault_secret_path": schema.StringAttribute{
				Description: "The path to the desired secret in the Vault service. This will be appended to the value of the base-url property for the associated Vault external server.",
				Required:    true,
			},
			"vault_secret_field_name": schema.StringAttribute{
				Description: "The name of the field in the Vault secret record that contains the passphrase to use to generate the encryption key.",
				Required:    true,
			},
			"vault_encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Vault Cipher Stream Provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"trust_store_file": schema.StringAttribute{
				Description: "The path to a file containing the information needed to trust the certificate presented by the Vault servers.",
				Optional:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "The passphrase needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Optional:    true,
				Sensitive:   true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "The store type for the specified trust store file. The value should likely be one of \"JKS\" or \"PKCS12\".",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
func addOptionalVaultCipherStreamProviderFields(ctx context.Context, addRequest *client.AddVaultCipherStreamProviderRequest, plan vaultCipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VaultExternalServer) {
		addRequest.VaultExternalServer = plan.VaultExternalServer.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.VaultServerBaseURI) {
		var slice []string
		plan.VaultServerBaseURI.ElementsAs(ctx, &slice, false)
		addRequest.VaultServerBaseURI = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VaultAuthenticationMethod) {
		addRequest.VaultAuthenticationMethod = plan.VaultAuthenticationMethod.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.VaultEncryptionMetadataFile) {
		addRequest.VaultEncryptionMetadataFile = plan.VaultEncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreFile) {
		addRequest.TrustStoreFile = plan.TrustStoreFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStorePin) {
		addRequest.TrustStorePin = plan.TrustStorePin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TrustStoreType) {
		addRequest.TrustStoreType = plan.TrustStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a VaultCipherStreamProviderResponse object into the model struct
func readVaultCipherStreamProviderResponse(ctx context.Context, r *client.VaultCipherStreamProviderResponse, state *vaultCipherStreamProviderResourceModel, expectedValues *vaultCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.VaultExternalServer = internaltypes.StringTypeOrNil(r.VaultExternalServer, internaltypes.IsEmptyString(expectedValues.VaultExternalServer))
	state.VaultServerBaseURI = internaltypes.GetStringSet(r.VaultServerBaseURI)
	state.VaultAuthenticationMethod = internaltypes.StringTypeOrNil(r.VaultAuthenticationMethod, internaltypes.IsEmptyString(expectedValues.VaultAuthenticationMethod))
	state.VaultSecretPath = types.StringValue(r.VaultSecretPath)
	state.VaultSecretFieldName = types.StringValue(r.VaultSecretFieldName)
	state.VaultEncryptionMetadataFile = types.StringValue(r.VaultEncryptionMetadataFile)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, internaltypes.IsEmptyString(expectedValues.TrustStoreFile))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.TrustStorePin = expectedValues.TrustStorePin
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, internaltypes.IsEmptyString(expectedValues.TrustStoreType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createVaultCipherStreamProviderOperations(plan vaultCipherStreamProviderResourceModel, state vaultCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.VaultExternalServer, state.VaultExternalServer, "vault-external-server")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.VaultServerBaseURI, state.VaultServerBaseURI, "vault-server-base-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultAuthenticationMethod, state.VaultAuthenticationMethod, "vault-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretPath, state.VaultSecretPath, "vault-secret-path")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretFieldName, state.VaultSecretFieldName, "vault-secret-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultEncryptionMetadataFile, state.VaultEncryptionMetadataFile, "vault-encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePin, state.TrustStorePin, "trust-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *vaultCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vaultCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddVaultCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumvaultCipherStreamProviderSchemaUrn{client.ENUMVAULTCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERVAULT},
		plan.VaultSecretPath.ValueString(),
		plan.VaultSecretFieldName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalVaultCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddVaultCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Vault Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state vaultCipherStreamProviderResourceModel
	readVaultCipherStreamProviderResponse(ctx, addResponse.VaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultVaultCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan vaultCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state vaultCipherStreamProviderResourceModel
	readVaultCipherStreamProviderResponse(ctx, readResponse.VaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createVaultCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Vault Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readVaultCipherStreamProviderResponse(ctx, updateResponse.VaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *vaultCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVaultCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readVaultCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state vaultCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Vault Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readVaultCipherStreamProviderResponse(ctx, readResponse.VaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *vaultCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultVaultCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateVaultCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan vaultCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state vaultCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createVaultCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Vault Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readVaultCipherStreamProviderResponse(ctx, updateResponse.VaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultVaultCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *vaultCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state vaultCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Vault Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *vaultCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVaultCipherStreamProvider(ctx, req, resp)
}

func (r *defaultVaultCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importVaultCipherStreamProvider(ctx, req, resp)
}

func importVaultCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
