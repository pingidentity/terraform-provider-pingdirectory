package cipherstreamprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
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
	_ resource.Resource                = &amazonSecretsManagerCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &amazonSecretsManagerCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &amazonSecretsManagerCipherStreamProviderResource{}
	_ resource.Resource                = &defaultAmazonSecretsManagerCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultAmazonSecretsManagerCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultAmazonSecretsManagerCipherStreamProviderResource{}
)

// Create a Amazon Secrets Manager Cipher Stream Provider resource
func NewAmazonSecretsManagerCipherStreamProviderResource() resource.Resource {
	return &amazonSecretsManagerCipherStreamProviderResource{}
}

func NewDefaultAmazonSecretsManagerCipherStreamProviderResource() resource.Resource {
	return &defaultAmazonSecretsManagerCipherStreamProviderResource{}
}

// amazonSecretsManagerCipherStreamProviderResource is the resource implementation.
type amazonSecretsManagerCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAmazonSecretsManagerCipherStreamProviderResource is the resource implementation.
type defaultAmazonSecretsManagerCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *amazonSecretsManagerCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_amazon_secrets_manager_cipher_stream_provider"
}

func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_amazon_secrets_manager_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *amazonSecretsManagerCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type amazonSecretsManagerCipherStreamProviderResourceModel struct {
	Id                     types.String `tfsdk:"id"`
	LastUpdated            types.String `tfsdk:"last_updated"`
	Notifications          types.Set    `tfsdk:"notifications"`
	RequiredActions        types.Set    `tfsdk:"required_actions"`
	AwsExternalServer      types.String `tfsdk:"aws_external_server"`
	SecretID               types.String `tfsdk:"secret_id"`
	SecretFieldName        types.String `tfsdk:"secret_field_name"`
	SecretVersionID        types.String `tfsdk:"secret_version_id"`
	SecretVersionStage     types.String `tfsdk:"secret_version_stage"`
	EncryptionMetadataFile types.String `tfsdk:"encryption_metadata_file"`
	Description            types.String `tfsdk:"description"`
	Enabled                types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *amazonSecretsManagerCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	amazonSecretsManagerCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	amazonSecretsManagerCipherStreamProviderSchema(ctx, req, resp, true)
}

func amazonSecretsManagerCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Amazon Secrets Manager Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS Secrets Manager.",
				Required:    true,
			},
			"secret_id": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.",
				Required:    true,
			},
			"secret_field_name": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the passphrase that will be used to generate the encryption key for protecting the contents of the encryption settings database.",
				Required:    true,
			},
			"secret_version_id": schema.StringAttribute{
				Description: "The unique identifier for the version of the secret to be retrieved.",
				Optional:    true,
			},
			"secret_version_stage": schema.StringAttribute{
				Description: "The staging label for the version of the secret to be retrieved.",
				Optional:    true,
			},
			"encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Amazon Secrets Manager Cipher Stream Provider.",
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

// Add config validators
func (r amazonSecretsManagerCipherStreamProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.Conflicting(
			path.MatchRoot("secret_version_id"),
			path.MatchRoot("secret_version_stage"),
		),
	}
}

// Add optional fields to create request
func addOptionalAmazonSecretsManagerCipherStreamProviderFields(ctx context.Context, addRequest *client.AddAmazonSecretsManagerCipherStreamProviderRequest, plan amazonSecretsManagerCipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecretVersionID) {
		addRequest.SecretVersionID = plan.SecretVersionID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecretVersionStage) {
		addRequest.SecretVersionStage = plan.SecretVersionStage.ValueStringPointer()
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

// Read a AmazonSecretsManagerCipherStreamProviderResponse object into the model struct
func readAmazonSecretsManagerCipherStreamProviderResponse(ctx context.Context, r *client.AmazonSecretsManagerCipherStreamProviderResponse, state *amazonSecretsManagerCipherStreamProviderResourceModel, expectedValues *amazonSecretsManagerCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.SecretID = types.StringValue(r.SecretID)
	state.SecretFieldName = types.StringValue(r.SecretFieldName)
	state.SecretVersionID = internaltypes.StringTypeOrNil(r.SecretVersionID, internaltypes.IsEmptyString(expectedValues.SecretVersionID))
	state.SecretVersionStage = internaltypes.StringTypeOrNil(r.SecretVersionStage, internaltypes.IsEmptyString(expectedValues.SecretVersionStage))
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAmazonSecretsManagerCipherStreamProviderOperations(plan amazonSecretsManagerCipherStreamProviderResourceModel, state amazonSecretsManagerCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.AwsExternalServer, state.AwsExternalServer, "aws-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretID, state.SecretID, "secret-id")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretFieldName, state.SecretFieldName, "secret-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretVersionID, state.SecretVersionID, "secret-version-id")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretVersionStage, state.SecretVersionStage, "secret-version-stage")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionMetadataFile, state.EncryptionMetadataFile, "encryption-metadata-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *amazonSecretsManagerCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan amazonSecretsManagerCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAmazonSecretsManagerCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumamazonSecretsManagerCipherStreamProviderSchemaUrn{client.ENUMAMAZONSECRETSMANAGERCIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERAMAZON_SECRETS_MANAGER},
		plan.AwsExternalServer.ValueString(),
		plan.SecretID.ValueString(),
		plan.SecretFieldName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAmazonSecretsManagerCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddAmazonSecretsManagerCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Amazon Secrets Manager Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state amazonSecretsManagerCipherStreamProviderResourceModel
	readAmazonSecretsManagerCipherStreamProviderResponse(ctx, addResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan amazonSecretsManagerCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Amazon Secrets Manager Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state amazonSecretsManagerCipherStreamProviderResourceModel
	readAmazonSecretsManagerCipherStreamProviderResponse(ctx, readResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAmazonSecretsManagerCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Amazon Secrets Manager Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAmazonSecretsManagerCipherStreamProviderResponse(ctx, updateResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *amazonSecretsManagerCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAmazonSecretsManagerCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAmazonSecretsManagerCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAmazonSecretsManagerCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state amazonSecretsManagerCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Amazon Secrets Manager Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAmazonSecretsManagerCipherStreamProviderResponse(ctx, readResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *amazonSecretsManagerCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAmazonSecretsManagerCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAmazonSecretsManagerCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAmazonSecretsManagerCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan amazonSecretsManagerCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state amazonSecretsManagerCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAmazonSecretsManagerCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Amazon Secrets Manager Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAmazonSecretsManagerCipherStreamProviderResponse(ctx, updateResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *amazonSecretsManagerCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state amazonSecretsManagerCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Amazon Secrets Manager Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *amazonSecretsManagerCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAmazonSecretsManagerCipherStreamProvider(ctx, req, resp)
}

func (r *defaultAmazonSecretsManagerCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAmazonSecretsManagerCipherStreamProvider(ctx, req, resp)
}

func importAmazonSecretsManagerCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
