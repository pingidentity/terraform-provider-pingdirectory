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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &azureKeyVaultCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &azureKeyVaultCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &azureKeyVaultCipherStreamProviderResource{}
	_ resource.Resource                = &defaultAzureKeyVaultCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultAzureKeyVaultCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultAzureKeyVaultCipherStreamProviderResource{}
)

// Create a Azure Key Vault Cipher Stream Provider resource
func NewAzureKeyVaultCipherStreamProviderResource() resource.Resource {
	return &azureKeyVaultCipherStreamProviderResource{}
}

func NewDefaultAzureKeyVaultCipherStreamProviderResource() resource.Resource {
	return &defaultAzureKeyVaultCipherStreamProviderResource{}
}

// azureKeyVaultCipherStreamProviderResource is the resource implementation.
type azureKeyVaultCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAzureKeyVaultCipherStreamProviderResource is the resource implementation.
type defaultAzureKeyVaultCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *azureKeyVaultCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_azure_key_vault_cipher_stream_provider"
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_azure_key_vault_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *azureKeyVaultCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type azureKeyVaultCipherStreamProviderResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	KeyVaultURI               types.String `tfsdk:"key_vault_uri"`
	AzureAuthenticationMethod types.String `tfsdk:"azure_authentication_method"`
	HttpProxyExternalServer   types.String `tfsdk:"http_proxy_external_server"`
	SecretName                types.String `tfsdk:"secret_name"`
	EncryptionMetadataFile    types.String `tfsdk:"encryption_metadata_file"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *azureKeyVaultCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	azureKeyVaultCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	azureKeyVaultCipherStreamProviderSchema(ctx, req, resp, true)
}

func azureKeyVaultCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Azure Key Vault Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"key_vault_uri": schema.StringAttribute{
				Description: "The URI that identifies the Azure Key Vault from which the secret is to be retrieved.",
				Required:    true,
			},
			"azure_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Azure service.",
				Required:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the Azure service. Supported in PingDirectory product version 9.2.0.0+.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secret_name": schema.StringAttribute{
				Description: "The name of the secret to retrieve.",
				Required:    true,
			},
			"encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Azure Key Vault Cipher Stream Provider.",
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

// Validate that any version restrictions are met in the plan
func (r *azureKeyVaultCipherStreamProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanAzureKeyVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanAzureKeyVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanAzureKeyVaultCipherStreamProvider(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model azureKeyVaultCipherStreamProviderResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add optional fields to create request
func addOptionalAzureKeyVaultCipherStreamProviderFields(ctx context.Context, addRequest *client.AddAzureKeyVaultCipherStreamProviderRequest, plan azureKeyVaultCipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
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

// Read a AzureKeyVaultCipherStreamProviderResponse object into the model struct
func readAzureKeyVaultCipherStreamProviderResponse(ctx context.Context, r *client.AzureKeyVaultCipherStreamProviderResponse, state *azureKeyVaultCipherStreamProviderResourceModel, expectedValues *azureKeyVaultCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.SecretName = types.StringValue(r.SecretName)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAzureKeyVaultCipherStreamProviderOperations(plan azureKeyVaultCipherStreamProviderResourceModel, state azureKeyVaultCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.KeyVaultURI, state.KeyVaultURI, "key-vault-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.AzureAuthenticationMethod, state.AzureAuthenticationMethod, "azure-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretName, state.SecretName, "secret-name")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionMetadataFile, state.EncryptionMetadataFile, "encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *azureKeyVaultCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan azureKeyVaultCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAzureKeyVaultCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumazureKeyVaultCipherStreamProviderSchemaUrn{client.ENUMAZUREKEYVAULTCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERAZURE_KEY_VAULT},
		plan.KeyVaultURI.ValueString(),
		plan.AzureAuthenticationMethod.ValueString(),
		plan.SecretName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAzureKeyVaultCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddAzureKeyVaultCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Azure Key Vault Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state azureKeyVaultCipherStreamProviderResourceModel
	readAzureKeyVaultCipherStreamProviderResponse(ctx, addResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAzureKeyVaultCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan azureKeyVaultCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Azure Key Vault Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state azureKeyVaultCipherStreamProviderResourceModel
	readAzureKeyVaultCipherStreamProviderResponse(ctx, readResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAzureKeyVaultCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Azure Key Vault Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAzureKeyVaultCipherStreamProviderResponse(ctx, updateResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *azureKeyVaultCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAzureKeyVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAzureKeyVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAzureKeyVaultCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state azureKeyVaultCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Azure Key Vault Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAzureKeyVaultCipherStreamProviderResponse(ctx, readResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *azureKeyVaultCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAzureKeyVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAzureKeyVaultCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAzureKeyVaultCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan azureKeyVaultCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state azureKeyVaultCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAzureKeyVaultCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Azure Key Vault Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAzureKeyVaultCipherStreamProviderResponse(ctx, updateResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAzureKeyVaultCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *azureKeyVaultCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state azureKeyVaultCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Azure Key Vault Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *azureKeyVaultCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAzureKeyVaultCipherStreamProvider(ctx, req, resp)
}

func (r *defaultAzureKeyVaultCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAzureKeyVaultCipherStreamProvider(ctx, req, resp)
}

func importAzureKeyVaultCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
