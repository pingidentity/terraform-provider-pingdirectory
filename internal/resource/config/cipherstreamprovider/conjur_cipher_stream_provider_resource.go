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
	_ resource.Resource                = &conjurCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &conjurCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &conjurCipherStreamProviderResource{}
	_ resource.Resource                = &defaultConjurCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultConjurCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultConjurCipherStreamProviderResource{}
)

// Create a Conjur Cipher Stream Provider resource
func NewConjurCipherStreamProviderResource() resource.Resource {
	return &conjurCipherStreamProviderResource{}
}

func NewDefaultConjurCipherStreamProviderResource() resource.Resource {
	return &defaultConjurCipherStreamProviderResource{}
}

// conjurCipherStreamProviderResource is the resource implementation.
type conjurCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultConjurCipherStreamProviderResource is the resource implementation.
type defaultConjurCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *conjurCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_conjur_cipher_stream_provider"
}

func (r *defaultConjurCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_conjur_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *conjurCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultConjurCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type conjurCipherStreamProviderResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	ConjurExternalServer     types.String `tfsdk:"conjur_external_server"`
	ConjurSecretRelativePath types.String `tfsdk:"conjur_secret_relative_path"`
	EncryptionMetadataFile   types.String `tfsdk:"encryption_metadata_file"`
	Description              types.String `tfsdk:"description"`
	Enabled                  types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *conjurCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	conjurCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultConjurCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	conjurCipherStreamProviderSchema(ctx, req, resp, true)
}

func conjurCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Conjur Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur server.",
				Required:    true,
			},
			"conjur_secret_relative_path": schema.StringAttribute{
				Description: "The portion of the path that follows the account name in the URI needed to obtain the secret passphrase to use to generate the encryption key. Any special characters in the path must be URL-encoded.",
				Required:    true,
			},
			"encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Conjur Cipher Stream Provider.",
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
func addOptionalConjurCipherStreamProviderFields(ctx context.Context, addRequest *client.AddConjurCipherStreamProviderRequest, plan conjurCipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionMetadataFile) {
		addRequest.EncryptionMetadataFile = plan.EncryptionMetadataFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a ConjurCipherStreamProviderResponse object into the model struct
func readConjurCipherStreamProviderResponse(ctx context.Context, r *client.ConjurCipherStreamProviderResponse, state *conjurCipherStreamProviderResourceModel, expectedValues *conjurCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.ConjurSecretRelativePath = types.StringValue(r.ConjurSecretRelativePath)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createConjurCipherStreamProviderOperations(plan conjurCipherStreamProviderResourceModel, state conjurCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurExternalServer, state.ConjurExternalServer, "conjur-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurSecretRelativePath, state.ConjurSecretRelativePath, "conjur-secret-relative-path")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionMetadataFile, state.EncryptionMetadataFile, "encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *conjurCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan conjurCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddConjurCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumconjurCipherStreamProviderSchemaUrn{client.ENUMCONJURCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERCONJUR},
		plan.ConjurExternalServer.ValueString(),
		plan.ConjurSecretRelativePath.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalConjurCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddConjurCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Conjur Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state conjurCipherStreamProviderResourceModel
	readConjurCipherStreamProviderResponse(ctx, addResponse.ConjurCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultConjurCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan conjurCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state conjurCipherStreamProviderResourceModel
	readConjurCipherStreamProviderResponse(ctx, readResponse.ConjurCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createConjurCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Conjur Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConjurCipherStreamProviderResponse(ctx, updateResponse.ConjurCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *conjurCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConjurCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConjurCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readConjurCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readConjurCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state conjurCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Conjur Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConjurCipherStreamProviderResponse(ctx, readResponse.ConjurCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *conjurCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConjurCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultConjurCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateConjurCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateConjurCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan conjurCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state conjurCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createConjurCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Conjur Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readConjurCipherStreamProviderResponse(ctx, updateResponse.ConjurCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultConjurCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *conjurCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state conjurCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Conjur Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *conjurCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConjurCipherStreamProvider(ctx, req, resp)
}

func (r *defaultConjurCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importConjurCipherStreamProvider(ctx, req, resp)
}

func importConjurCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
