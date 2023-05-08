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
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &amazonKeyManagementServiceCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &amazonKeyManagementServiceCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &amazonKeyManagementServiceCipherStreamProviderResource{}
	_ resource.Resource                = &defaultAmazonKeyManagementServiceCipherStreamProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultAmazonKeyManagementServiceCipherStreamProviderResource{}
	_ resource.ResourceWithImportState = &defaultAmazonKeyManagementServiceCipherStreamProviderResource{}
)

// Create a Amazon Key Management Service Cipher Stream Provider resource
func NewAmazonKeyManagementServiceCipherStreamProviderResource() resource.Resource {
	return &amazonKeyManagementServiceCipherStreamProviderResource{}
}

func NewDefaultAmazonKeyManagementServiceCipherStreamProviderResource() resource.Resource {
	return &defaultAmazonKeyManagementServiceCipherStreamProviderResource{}
}

// amazonKeyManagementServiceCipherStreamProviderResource is the resource implementation.
type amazonKeyManagementServiceCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAmazonKeyManagementServiceCipherStreamProviderResource is the resource implementation.
type defaultAmazonKeyManagementServiceCipherStreamProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *amazonKeyManagementServiceCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_amazon_key_management_service_cipher_stream_provider"
}

func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_amazon_key_management_service_cipher_stream_provider"
}

// Configure adds the provider configured client to the resource.
func (r *amazonKeyManagementServiceCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9200
}

type amazonKeyManagementServiceCipherStreamProviderResourceModel struct {
	Id                      types.String `tfsdk:"id"`
	LastUpdated             types.String `tfsdk:"last_updated"`
	Notifications           types.Set    `tfsdk:"notifications"`
	RequiredActions         types.Set    `tfsdk:"required_actions"`
	EncryptedPassphraseFile types.String `tfsdk:"encrypted_passphrase_file"`
	AwsExternalServer       types.String `tfsdk:"aws_external_server"`
	AwsAccessKeyID          types.String `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKey      types.String `tfsdk:"aws_secret_access_key"`
	AwsRegionName           types.String `tfsdk:"aws_region_name"`
	KmsEncryptionKeyArn     types.String `tfsdk:"kms_encryption_key_arn"`
	Description             types.String `tfsdk:"description"`
	Enabled                 types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *amazonKeyManagementServiceCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	amazonKeyManagementServiceCipherStreamProviderSchema(ctx, req, resp, false)
}

func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	amazonKeyManagementServiceCipherStreamProviderSchema(ctx, req, resp, true)
}

func amazonKeyManagementServiceCipherStreamProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, setOptionalToComputed bool) {
	schema := schema.Schema{
		Description: "Manages a Amazon Key Management Service Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"encrypted_passphrase_file": schema.StringAttribute{
				Description: "The path to a file that will hold the encrypted passphrase used by this cipher stream provider.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the Amazon Key Management Service.",
				Optional:    true,
			},
			"aws_access_key_id": schema.StringAttribute{
				Description: "The access key ID that will be used if this cipher stream provider will authenticate to the Amazon Key Management Service using an access key rather than an IAM role associated with an EC2 instance.",
				Optional:    true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				Description: "The secret access key that will be used if this cipher stream provider will authenticate to the Amazon Key Management Service using an access key rather than an IAM role associated with an EC2 instance.",
				Optional:    true,
				Sensitive:   true,
			},
			"aws_region_name": schema.StringAttribute{
				Description: "The name of the Amazon Web Services region that holds the encryption key. This is optional, and if it is not provided, then the server will attempt to determine the region from the key ARN.",
				Optional:    true,
			},
			"kms_encryption_key_arn": schema.StringAttribute{
				Description: "The Amazon resource name (ARN) for the KMS key that will be used to encrypt the contents of the passphrase file. This key must exist, and the AWS client must have access to encrypt and decrypt data using this key.",
				Required:    true,
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
func (r amazonKeyManagementServiceCipherStreamProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.Implies(
			path.MatchRoot("aws_access_key_id"),
			path.MatchRoot("aws_secret_access_key"),
		),
		configvalidators.Implies(
			path.MatchRoot("aws_secret_access_key"),
			path.MatchRoot("aws_access_key_id"),
		),
	}
}

// Add optional fields to create request
func addOptionalAmazonKeyManagementServiceCipherStreamProviderFields(ctx context.Context, addRequest *client.AddAmazonKeyManagementServiceCipherStreamProviderRequest, plan amazonKeyManagementServiceCipherStreamProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptedPassphraseFile) {
		addRequest.EncryptedPassphraseFile = plan.EncryptedPassphraseFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsExternalServer) {
		addRequest.AwsExternalServer = plan.AwsExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsAccessKeyID) {
		addRequest.AwsAccessKeyID = plan.AwsAccessKeyID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsSecretAccessKey) {
		addRequest.AwsSecretAccessKey = plan.AwsSecretAccessKey.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AwsRegionName) {
		addRequest.AwsRegionName = plan.AwsRegionName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Read a AmazonKeyManagementServiceCipherStreamProviderResponse object into the model struct
func readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx context.Context, r *client.AmazonKeyManagementServiceCipherStreamProviderResponse, state *amazonKeyManagementServiceCipherStreamProviderResourceModel, expectedValues *amazonKeyManagementServiceCipherStreamProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.EncryptedPassphraseFile = types.StringValue(r.EncryptedPassphraseFile)
	state.AwsExternalServer = internaltypes.StringTypeOrNil(r.AwsExternalServer, internaltypes.IsEmptyString(expectedValues.AwsExternalServer))
	state.AwsAccessKeyID = internaltypes.StringTypeOrNil(r.AwsAccessKeyID, internaltypes.IsEmptyString(expectedValues.AwsAccessKeyID))
	// Obscured values aren't returned from the PD Configuration API - just use the expected value
	state.AwsSecretAccessKey = expectedValues.AwsSecretAccessKey
	state.AwsRegionName = internaltypes.StringTypeOrNil(r.AwsRegionName, internaltypes.IsEmptyString(expectedValues.AwsRegionName))
	state.KmsEncryptionKeyArn = types.StringValue(r.KmsEncryptionKeyArn)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createAmazonKeyManagementServiceCipherStreamProviderOperations(plan amazonKeyManagementServiceCipherStreamProviderResourceModel, state amazonKeyManagementServiceCipherStreamProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptedPassphraseFile, state.EncryptedPassphraseFile, "encrypted-passphrase-file")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsExternalServer, state.AwsExternalServer, "aws-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsAccessKeyID, state.AwsAccessKeyID, "aws-access-key-id")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsSecretAccessKey, state.AwsSecretAccessKey, "aws-secret-access-key")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsRegionName, state.AwsRegionName, "aws-region-name")
	operations.AddStringOperationIfNecessary(&ops, plan.KmsEncryptionKeyArn, state.KmsEncryptionKeyArn, "kms-encryption-key-arn")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a new resource
func (r *amazonKeyManagementServiceCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan amazonKeyManagementServiceCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	addRequest := client.NewAddAmazonKeyManagementServiceCipherStreamProviderRequest(plan.Id.ValueString(),
		[]client.EnumamazonKeyManagementServiceCipherStreamProviderSchemaUrn{client.ENUMAMAZONKEYMANAGEMENTSERVICECIPHERSTREAMPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0CIPHER_STREAM_PROVIDERAMAZON_KEY_MANAGEMENT_SERVICE},
		plan.KmsEncryptionKeyArn.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAmazonKeyManagementServiceCipherStreamProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.CipherStreamProviderApi.AddCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddCipherStreamProviderRequest(
		client.AddAmazonKeyManagementServiceCipherStreamProviderRequestAsAddCipherStreamProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.AddCipherStreamProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Amazon Key Management Service Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state amazonKeyManagementServiceCipherStreamProviderResourceModel
	readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, addResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)

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
func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan amazonKeyManagementServiceCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Amazon Key Management Service Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state amazonKeyManagementServiceCipherStreamProviderResourceModel
	readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Id.ValueString())
	ops := createAmazonKeyManagementServiceCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Amazon Key Management Service Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, updateResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *amazonKeyManagementServiceCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAmazonKeyManagementServiceCipherStreamProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state amazonKeyManagementServiceCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Id.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Amazon Key Management Service Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update a resource
func (r *amazonKeyManagementServiceCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAmazonKeyManagementServiceCipherStreamProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan amazonKeyManagementServiceCipherStreamProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state amazonKeyManagementServiceCipherStreamProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.CipherStreamProviderApi.UpdateCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Id.ValueString())

	// Determine what update operations are necessary
	ops := createAmazonKeyManagementServiceCipherStreamProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.CipherStreamProviderApi.UpdateCipherStreamProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Amazon Key Management Service Cipher Stream Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readAmazonKeyManagementServiceCipherStreamProviderResponse(ctx, updateResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *amazonKeyManagementServiceCipherStreamProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state amazonKeyManagementServiceCipherStreamProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProviderExecute(r.apiClient.CipherStreamProviderApi.DeleteCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Amazon Key Management Service Cipher Stream Provider", err, httpResp)
		return
	}
}

func (r *amazonKeyManagementServiceCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp)
}

func (r *defaultAmazonKeyManagementServiceCipherStreamProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAmazonKeyManagementServiceCipherStreamProvider(ctx, req, resp)
}

func importAmazonKeyManagementServiceCipherStreamProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
