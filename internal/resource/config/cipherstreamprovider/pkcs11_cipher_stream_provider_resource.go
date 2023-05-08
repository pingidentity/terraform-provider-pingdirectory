package cipherstreamprovider

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
	_ resource.Resource                = &pkcs11CipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &pkcs11CipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &pkcs11CipherStreamProviderResource{}
	_ resource.Resource                = &defaultPkcs11CipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultPkcs11CipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultPkcs11CipherStreamProviderResource{}
)

// Create a Pkcs11 Cipher Stream Provider resource
func NewPkcs11CipherStreamProviderResource() resource.Resource {
	return &pkcs11CipherStreamProviderResource{}
}

func NewDefaultPkcs11CipherStreamProviderResource() resource.Resource {
	return &defaultPkcs11CipherStreamProviderResource{}
}

// pkcs11CipherStreamProviderResource is the resource implementation.
type pkcs11CipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPkcs11CipherStreamProviderResource is the resource implementation.
type defaultPkcs11CipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *pkcs11CipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pkcs11_cipher_stream_provider"
}

func (r *defaultPkcs11CipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_pkcs11_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *pkcs11CipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultPkcs11CipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type pkcs11CipherStreamProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Pkcs11ProviderClass             types.String `tfsdk:"pkcs11_provider_class"`
	Pkcs11ProviderConfigurationFile types.String `tfsdk:"pkcs11_provider_configuration_file"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinEnvironmentVariable  types.String `tfsdk:"key_store_pin_environment_variable"`
	Pkcs11KeyStoreType              types.String `tfsdk:"pkcs11_key_store_type"`
	SslCertNickname                 types.String `tfsdk:"ssl_cert_nickname"`
	EncryptionMetadataFile          types.String `tfsdk:"encryption_metadata_file"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *pkcs11CipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pkcs11CipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultPkcs11CipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	pkcs11CipherStreamProviderSchema(ctx, req, resp, true)
}

func pkcs11CipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Pkcs11 Cipher Stream Provider.",
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
			"key_store_pin": schema.StringAttribute{
				Description: "The clear-text user PIN needed to interact with the PKCS #11 token.",
				Optional:    true,
				Sensitive:   true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description: "The path to a file containing the user PIN needed to interact with the PKCS #11 token. The file must exist and must contain exactly one line with a clear-text representation of the PIN.",
				Optional:    true,
			},
			"key_store_pin_environment_variable": schema.StringAttribute{
				Description: "The name of an environment variable whose value is the user PIN needed to interact with the PKCS #11 token. The environment variable must be defined and must contain a clear-text representation of the PIN.",
				Optional:    true,
			},
			"pkcs11_key_store_type": schema.StringAttribute{
				Description: "The key store type to use when obtaining an instance of a key store for interacting with a PKCS #11 token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "The alias for the certificate in the PKCS #11 token that will be used to wrap the encryption key. The target certificate must exist in the PKCS #11 token, and it must have an RSA key pair because the JVM does not currently provide adequate key wrapping support for elliptic curve key pairs.  If you have also configured the server to use a PKCS #11 token for accessing listener certificates, we strongly recommend that you use a different certificate to protect the contents of the encryption settings database than you use for negotiating TLS sessions with clients. It is imperative that the certificate used by this PKCS11 Cipher Stream Provider remain constant for the life of the provider because if the certificate were to be replaced, then the contents of the encryption settings database could become inaccessible. Unlike with listener certificates used for TLS negotiation that need to be replaced on a regular basis, this PKCS11 Cipher Stream Provider does not consider the validity period for the associated certificate, and it will continue to function even after the certificate has expired.  If you need to rotate the certificate used to protect the server's encryption settings database, you should first install the desired new certificate in the PKCS #11 token under a different alias. Then, you should create a new instance of this PKCS11 Cipher Stream Provider that is configured to use that certificate, and that also uses a different value for the encryption-metadata-file because the information in that file is tied to the certificate used to generate it. Finally, you will need to update the global configuration so that the encryption-settings-cipher-stream-provider property references the new cipher stream provider rather than this one. The update to the global configuration must be done with the server online so that it can properly re-encrypt the contents of the encryption settings database with the correct key tied to the new certificate.",
				Required:    true,
			},
			"encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this PKCS11 Cipher Stream Provider.",
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
func addOptionalPkcs11CipherStreamProviderFields(ctx context.Context, addRequest *client.AddPkcs11CipherStreamProviderRequest, plan pkcs11CipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderClass) {
		addRequest.Pkcs11ProviderClass = plan.Pkcs11ProviderClass.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderConfigurationFile) {
		addRequest.Pkcs11ProviderConfigurationFile = plan.Pkcs11ProviderConfigurationFile.ValueStringPointer()
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
	if internaltypes.IsNonEmptyString(plan.KeyStorePinEnvironmentVariable) {
		addRequest.KeyStorePinEnvironmentVariable = plan.KeyStorePinEnvironmentVariable.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11KeyStoreType) {
		addRequest.Pkcs11KeyStoreType = plan.Pkcs11KeyStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionMetadataFile) {
		addRequest.EncryptionMetadataFile = plan.EncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a Pkcs11CipherStreamProviderResponse object into the model struct
func readPkcs11CipherStreamProviderResponse(ctx context.Context, r *client.Pkcs11CipherStreamProviderResponse, state *pkcs11CipherStreamProviderResourceModel, expectedValues *pkcs11CipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Pkcs11ProviderClass = internaltypes.StringTypeOrNil(r.Pkcs11ProviderClass, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderClass))
	state.Pkcs11ProviderConfigurationFile = internaltypes.StringTypeOrNil(r.Pkcs11ProviderConfigurationFile, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderConfigurationFile))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.KeyStorePin = expectedValues.KeyStorePin
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, internaltypes.IsEmptyString(expectedValues.KeyStorePinFile))
	state.KeyStorePinEnvironmentVariable = internaltypes.StringTypeOrNil(r.KeyStorePinEnvironmentVariable, internaltypes.IsEmptyString(expectedValues.KeyStorePinEnvironmentVariable))
	state.Pkcs11KeyStoreType = internaltypes.StringTypeOrNil(r.Pkcs11KeyStoreType, internaltypes.IsEmptyString(expectedValues.Pkcs11KeyStoreType))
	state.SslCertNickname = types.StringValue(r.SslCertNickname)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createPkcs11CipherStreamProviderOperations(plan pkcs11CipherStreamProviderResourceModel, state pkcs11CipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderClass, state.Pkcs11ProviderClass, "pkcs11-provider-class")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderConfigurationFile, state.Pkcs11ProviderConfigurationFile, "pkcs11-provider-configuration-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePin, state.KeyStorePin, "key-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinFile, state.KeyStorePinFile, "key-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinEnvironmentVariable, state.KeyStorePinEnvironmentVariable, "key-store-pin-environment-variable")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11KeyStoreType, state.Pkcs11KeyStoreType, "pkcs11-key-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionMetadataFile, state.EncryptionMetadataFile, "encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *pkcs11CipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pkcs11CipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddPkcs11CipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.Enumpkcs11CipherStreamProviderSchemaUrn{client.ENUMPKCS11CIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERPKCS11},
		plan.SslCertNickname.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalPkcs11CipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddPkcs11CipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Pkcs11 Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state pkcs11CipherStreamProviderResourceModel
	readPkcs11CipherStreamProviderResponse(ctx, addResponse.Pkcs11CipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultPkcs11CipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan pkcs11CipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pkcs11 Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state pkcs11CipherStreamProviderResourceModel
	readPkcs11CipherStreamProviderResponse(ctx, readResponse.Pkcs11CipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createPkcs11CipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pkcs11 Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPkcs11CipherStreamProviderResponse(ctx, updateResponse.Pkcs11CipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *pkcs11CipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPkcs11CipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPkcs11CipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPkcs11CipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readPkcs11CipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state pkcs11CipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Pkcs11 Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readPkcs11CipherStreamProviderResponse(ctx, readResponse.Pkcs11CipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *pkcs11CipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePkcs11CipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPkcs11CipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePkcs11CipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePkcs11CipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan pkcs11CipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state pkcs11CipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createPkcs11CipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Pkcs11 Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readPkcs11CipherStreamProviderResponse(ctx, updateResponse.Pkcs11CipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultPkcs11CipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *pkcs11CipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state pkcs11CipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Pkcs11 Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *pkcs11CipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPkcs11CipherStreamProvider(ctx, req, resp)
}

func (r *defaultPkcs11CipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPkcs11CipherStreamProvider(ctx, req, resp)
}

func importPkcs11CipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
