package cryptomanager

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &cryptoManagerResource{}
	_ resource.ResourceWithConfigure   = &cryptoManagerResource{}
	_ resource.ResourceWithImportState = &cryptoManagerResource{}
)

// Create a Crypto Manager resource
func NewCryptoManagerResource() resource.Resource {
	return &cryptoManagerResource{}
}

// cryptoManagerResource is the resource implementation.
type cryptoManagerResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *cryptoManagerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_crypto_manager"
}

// Configure adds the provider configured client to the resource.
func (r *cryptoManagerResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type cryptoManagerResourceModel struct {
	Id                               types.String `tfsdk:"id"`
	LastUpdated                      types.String `tfsdk:"last_updated"`
	Notifications                    types.Set    `tfsdk:"notifications"`
	RequiredActions                  types.Set    `tfsdk:"required_actions"`
	Type                             types.String `tfsdk:"type"`
	DigestAlgorithm                  types.String `tfsdk:"digest_algorithm"`
	MacAlgorithm                     types.String `tfsdk:"mac_algorithm"`
	MacKeyLength                     types.Int64  `tfsdk:"mac_key_length"`
	SigningEncryptionSettingsID      types.String `tfsdk:"signing_encryption_settings_id"`
	CipherTransformation             types.String `tfsdk:"cipher_transformation"`
	CipherKeyLength                  types.Int64  `tfsdk:"cipher_key_length"`
	KeyWrappingTransformation        types.String `tfsdk:"key_wrapping_transformation"`
	SslProtocol                      types.Set    `tfsdk:"ssl_protocol"`
	SslCipherSuite                   types.Set    `tfsdk:"ssl_cipher_suite"`
	OutboundSSLProtocol              types.Set    `tfsdk:"outbound_ssl_protocol"`
	OutboundSSLCipherSuite           types.Set    `tfsdk:"outbound_ssl_cipher_suite"`
	EnableSha1CipherSuites           types.Bool   `tfsdk:"enable_sha_1_cipher_suites"`
	EnableRsaKeyExchangeCipherSuites types.Bool   `tfsdk:"enable_rsa_key_exchange_cipher_suites"`
	SslCertNickname                  types.String `tfsdk:"ssl_cert_nickname"`
}

// GetSchema defines the schema for the resource.
func (r *cryptoManagerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Manages a Crypto Manager.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Crypto Manager resource. Options are ['crypto-manager']",
				Optional:    false,
				Required:    false,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"crypto-manager"}...),
				},
			},
			"digest_algorithm": schema.StringAttribute{
				Description: "Specifies the preferred message digest algorithm for the Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_algorithm": schema.StringAttribute{
				Description: "Specifies the preferred MAC algorithm for the Directory Server.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mac_key_length": schema.Int64Attribute{
				Description: "Specifies the key length in bits for the preferred MAC algorithm.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"signing_encryption_settings_id": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. The ID of the encryption settings definition to use for generating digital signatures. If this is not specified, then the server's preferred encryption settings definition will be used.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cipher_transformation": schema.StringAttribute{
				Description: "Specifies the cipher for the Directory Server using the syntax algorithm/mode/padding.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cipher_key_length": schema.Int64Attribute{
				Description: "Specifies the key length in bits for the preferred cipher.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"key_wrapping_transformation": schema.StringAttribute{
				Description: "The preferred key wrapping transformation for the Directory Server. This value must be the same for all server instances in a replication topology.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_protocol": schema.SetAttribute{
				Description: "Specifies the names of TLS protocols that are allowed for use in secure communication.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the TLS cipher suites that are allowed for use in secure communication.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"outbound_ssl_protocol": schema.SetAttribute{
				Description: "Specifies the names of the TLS protocols that will be enabled for outbound connections initiated by the Directory Server.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"outbound_ssl_cipher_suite": schema.SetAttribute{
				Description: "Specifies the names of the TLS cipher suites that will be enabled for outbound connections initiated by the Directory Server.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_sha_1_cipher_suites": schema.BoolAttribute{
				Description: "Indicates whether to enable support for TLS cipher suites that use the SHA-1 digest algorithm. The SHA-1 digest algorithm is no longer considered secure and is not recommended for use.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"enable_rsa_key_exchange_cipher_suites": schema.BoolAttribute{
				Description: "Indicates whether to enable support for TLS cipher suites that use the RSA key exchange algorithm. Cipher suites that rely on RSA key exchange are not recommended because they do not support forward secrecy, which means that if the private key is compromised, then any communication negotiated using that private key should also be considered compromised.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "Specifies the nickname (also called the alias) of the certificate that the Crypto Manager should use when performing SSL communication.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	config.AddCommonResourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *cryptoManagerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	compare, err := version.Compare(r.providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model cryptoManagerResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.SigningEncryptionSettingsID) {
		resp.Diagnostics.AddError("Attribute 'signing_encryption_settings_id' not supported by PingDirectory version "+r.providerConfig.ProductVersion, "")
	}
}

// Read a CryptoManagerResponse object into the model struct
func readCryptoManagerResponse(ctx context.Context, r *client.CryptoManagerResponse, state *cryptoManagerResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("crypto-manager")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.DigestAlgorithm = internaltypes.StringTypeOrNil(r.DigestAlgorithm, true)
	state.MacAlgorithm = internaltypes.StringTypeOrNil(r.MacAlgorithm, true)
	state.MacKeyLength = internaltypes.Int64TypeOrNil(r.MacKeyLength)
	state.SigningEncryptionSettingsID = internaltypes.StringTypeOrNil(r.SigningEncryptionSettingsID, true)
	state.CipherTransformation = internaltypes.StringTypeOrNil(r.CipherTransformation, true)
	state.CipherKeyLength = internaltypes.Int64TypeOrNil(r.CipherKeyLength)
	state.KeyWrappingTransformation = internaltypes.StringTypeOrNil(r.KeyWrappingTransformation, true)
	state.SslProtocol = internaltypes.GetStringSet(r.SslProtocol)
	state.SslCipherSuite = internaltypes.GetStringSet(r.SslCipherSuite)
	state.OutboundSSLProtocol = internaltypes.GetStringSet(r.OutboundSSLProtocol)
	state.OutboundSSLCipherSuite = internaltypes.GetStringSet(r.OutboundSSLCipherSuite)
	state.EnableSha1CipherSuites = internaltypes.BoolTypeOrNil(r.EnableSha1CipherSuites)
	state.EnableRsaKeyExchangeCipherSuites = internaltypes.BoolTypeOrNil(r.EnableRsaKeyExchangeCipherSuites)
	state.SslCertNickname = internaltypes.StringTypeOrNil(r.SslCertNickname, true)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Create any update operations necessary to make the state match the plan
func createCryptoManagerOperations(plan cryptoManagerResourceModel, state cryptoManagerResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.DigestAlgorithm, state.DigestAlgorithm, "digest-algorithm")
	operations.AddStringOperationIfNecessary(&ops, plan.MacAlgorithm, state.MacAlgorithm, "mac-algorithm")
	operations.AddInt64OperationIfNecessary(&ops, plan.MacKeyLength, state.MacKeyLength, "mac-key-length")
	operations.AddStringOperationIfNecessary(&ops, plan.SigningEncryptionSettingsID, state.SigningEncryptionSettingsID, "signing-encryption-settings-id")
	operations.AddStringOperationIfNecessary(&ops, plan.CipherTransformation, state.CipherTransformation, "cipher-transformation")
	operations.AddInt64OperationIfNecessary(&ops, plan.CipherKeyLength, state.CipherKeyLength, "cipher-key-length")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyWrappingTransformation, state.KeyWrappingTransformation, "key-wrapping-transformation")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SslProtocol, state.SslProtocol, "ssl-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SslCipherSuite, state.SslCipherSuite, "ssl-cipher-suite")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OutboundSSLProtocol, state.OutboundSSLProtocol, "outbound-ssl-protocol")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.OutboundSSLCipherSuite, state.OutboundSSLCipherSuite, "outbound-ssl-cipher-suite")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableSha1CipherSuites, state.EnableSha1CipherSuites, "enable-sha-1-cipher-suites")
	operations.AddBoolOperationIfNecessary(&ops, plan.EnableRsaKeyExchangeCipherSuites, state.EnableRsaKeyExchangeCipherSuites, "enable-rsa-key-exchange-cipher-suites")
	operations.AddStringOperationIfNecessary(&ops, plan.SslCertNickname, state.SslCertNickname, "ssl-cert-nickname")
	return ops
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *cryptoManagerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan cryptoManagerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CryptoManagerApi.GetCryptoManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Crypto Manager", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state cryptoManagerResourceModel
	readCryptoManagerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.CryptoManagerApi.UpdateCryptoManager(config.ProviderBasicAuthContext(ctx, r.providerConfig))
	ops := createCryptoManagerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CryptoManagerApi.UpdateCryptoManagerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Crypto Manager", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCryptoManagerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *cryptoManagerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state cryptoManagerResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CryptoManagerApi.GetCryptoManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Crypto Manager", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCryptoManagerResponse(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *cryptoManagerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan cryptoManagerResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state cryptoManagerResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.CryptoManagerApi.UpdateCryptoManager(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))

	// Determine what update operations are necessary
	ops := createCryptoManagerOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.CryptoManagerApi.UpdateCryptoManagerExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Crypto Manager", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		readCryptoManagerResponse(ctx, updateResponse, &state, &resp.Diagnostics)
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
func (r *cryptoManagerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *cryptoManagerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Set a placeholder id value to appease terraform.
	// The real attributes will be imported when terraform performs a read after the import.
	// If no value is set here, Terraform will error out when importing.
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), "id")...)
}
