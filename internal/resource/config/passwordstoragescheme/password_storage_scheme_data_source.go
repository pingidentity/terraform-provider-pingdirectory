package passwordstoragescheme

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passwordStorageSchemeDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordStorageSchemeDataSource{}
)

// Create a Password Storage Scheme data source
func NewPasswordStorageSchemeDataSource() datasource.DataSource {
	return &passwordStorageSchemeDataSource{}
}

// passwordStorageSchemeDataSource is the datasource implementation.
type passwordStorageSchemeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passwordStorageSchemeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_storage_scheme"
}

// Configure adds the provider configured client to the data source.
func (r *passwordStorageSchemeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordStorageSchemeDataSourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Type                              types.String `tfsdk:"type"`
	ScryptCpuMemoryCostFactorExponent types.Int64  `tfsdk:"scrypt_cpu_memory_cost_factor_exponent"`
	ScryptBlockSize                   types.Int64  `tfsdk:"scrypt_block_size"`
	ScryptParallelizationParameter    types.Int64  `tfsdk:"scrypt_parallelization_parameter"`
	ConjurExternalServer              types.String `tfsdk:"conjur_external_server"`
	KeyVaultURI                       types.String `tfsdk:"key_vault_uri"`
	AzureAuthenticationMethod         types.String `tfsdk:"azure_authentication_method"`
	HttpProxyExternalServer           types.String `tfsdk:"http_proxy_external_server"`
	AwsExternalServer                 types.String `tfsdk:"aws_external_server"`
	BcryptCostFactor                  types.Int64  `tfsdk:"bcrypt_cost_factor"`
	EncryptionSettingsDefinitionID    types.String `tfsdk:"encryption_settings_definition_id"`
	DigestAlgorithm                   types.String `tfsdk:"digest_algorithm"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	VaultExternalServer               types.String `tfsdk:"vault_external_server"`
	DefaultField                      types.String `tfsdk:"default_field"`
	PasswordEncodingMechanism         types.String `tfsdk:"password_encoding_mechanism"`
	NumDigestRounds                   types.Int64  `tfsdk:"num_digest_rounds"`
	MaxPasswordLength                 types.Int64  `tfsdk:"max_password_length"`
	IterationCount                    types.Int64  `tfsdk:"iteration_count"`
	ParallelismFactor                 types.Int64  `tfsdk:"parallelism_factor"`
	MemoryUsageKb                     types.Int64  `tfsdk:"memory_usage_kb"`
	SaltLengthBytes                   types.Int64  `tfsdk:"salt_length_bytes"`
	DerivedKeyLengthBytes             types.Int64  `tfsdk:"derived_key_length_bytes"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *passwordStorageSchemeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Password Storage Scheme.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Storage Scheme resource. Options are ['salted-sha256', 'argon2d', 'crypt', 'argon2i', 'base64', 'salted-md5', 'aes', 'argon2id', 'vault', 'third-party', 'argon2', 'third-party-enhanced', 'pbkdf2', 'rc4', 'salted-sha384', 'triple-des', 'clear', 'aes-256', 'bcrypt', 'blowfish', 'sha1', 'amazon-secrets-manager', 'azure-key-vault', 'conjur', 'salted-sha1', 'salted-sha512', 'scrypt', 'md5']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scrypt_cpu_memory_cost_factor_exponent": schema.Int64Attribute{
				Description: "Specifies the exponent that should be used for the CPU/memory cost factor. The cost factor must be a power of two, so the value of this property represents the power to which two is raised. The CPU/memory cost factor specifies the number of iterations required for encoding the password, and also affects the amount of memory required during processing. A higher cost factor requires more processing and more memory to generate a password, which makes attacks against the password more expensive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scrypt_block_size": schema.Int64Attribute{
				Description: "Specifies the block size for the digest that will be used in the course of encoding passwords. Increasing the block size while keeping the CPU/memory cost factor constant will increase the amount of memory required to encode a password, but it also increases the ratio of sequential memory access to random memory access (and sequential memory access is generally faster than random memory access).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scrypt_parallelization_parameter": schema.Int64Attribute{
				Description: "Specifies the number of times that scrypt has to perform the entire encoding process to produce the final result.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur instance containing user passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_vault_uri": schema.StringAttribute{
				Description: "The URI that identifies the Azure Key Vault from which the secret is to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"azure_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Azure service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Azure service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS Secrets Manager service.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bcrypt_cost_factor": schema.Int64Attribute{
				Description: "Specifies the cost factor to use when encoding passwords with Bcrypt. A higher cost factor requires more processing to generate a password, which makes attacks against the password more expensive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "The identifier for the encryption settings definition that should be used to derive the encryption key to use when encrypting new passwords. If this is not provided, the server's preferred encryption settings definition will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"digest_algorithm": schema.StringAttribute{
				Description: "Specifies the digest algorithm that will be used when encoding passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party`: The fully-qualified name of the Java class providing the logic for the Third Party Password Storage Scheme. When the `type` attribute is set to `third-party-enhanced`: The fully-qualified name of the Java class providing the logic for the Third Party Enhanced Password Storage Scheme.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party`: The fully-qualified name of the Java class providing the logic for the Third Party Password Storage Scheme.\n  - `third-party-enhanced`: The fully-qualified name of the Java class providing the logic for the Third Party Enhanced Password Storage Scheme.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"extension_argument": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party`: The set of arguments used to customize the behavior for the Third Party Password Storage Scheme. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-enhanced`: The set of arguments used to customize the behavior for the Third Party Enhanced Password Storage Scheme. Each configuration property should be given in the form 'name=value'.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party`: The set of arguments used to customize the behavior for the Third Party Password Storage Scheme. Each configuration property should be given in the form 'name=value'.\n  - `third-party-enhanced`: The set of arguments used to customize the behavior for the Third Party Enhanced Password Storage Scheme. Each configuration property should be given in the form 'name=value'.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault instance containing the passphrase.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_field": schema.StringAttribute{
				Description: "The default name of the field in JSON objects contained in the AWS Secrets Manager service that contains the password for the target user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_encoding_mechanism": schema.StringAttribute{
				Description: "Specifies the mechanism that should be used to encode clear-text passwords for use with this scheme.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"num_digest_rounds": schema.Int64Attribute{
				Description: "Specifies the number of digest rounds to use for the SHA-2 encodings. This will not be used for the legacy or MD5-based encodings.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"max_password_length": schema.Int64Attribute{
				Description: "Specifies the maximum allowed length, in bytes, for passwords encoded with this scheme, which can help mitigate denial of service attacks from clients that attempt to bind with very long passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"iteration_count": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of rounds of cryptographic processing required in the course of encoding each password. When the `type` attribute is set to `pbkdf2`: Specifies the number of iterations to use when encoding passwords. The value must be greater than or equal to 1000.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of rounds of cryptographic processing required in the course of encoding each password.\n  - `pbkdf2`: Specifies the number of iterations to use when encoding passwords. The value must be greater than or equal to 1000.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"parallelism_factor": schema.Int64Attribute{
				Description: "The number of concurrent threads that will be used in the course of encoding each password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"memory_usage_kb": schema.Int64Attribute{
				Description: "The number of kilobytes of memory that must be used in the course of encoding each password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"salt_length_bytes": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`salted-sha256`, `salted-md5`, `salted-sha384`, `salted-sha1`, `salted-sha512`]: Specifies the number of bytes to use for the generated salt. When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the generated salt. When the `type` attribute is set to `pbkdf2`: Specifies the number of bytes to use for the generated salt. The value must be greater than or equal to 8.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`salted-sha256`, `salted-md5`, `salted-sha384`, `salted-sha1`, `salted-sha512`]: Specifies the number of bytes to use for the generated salt.\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the generated salt.\n  - `pbkdf2`: Specifies the number of bytes to use for the generated salt. The value must be greater than or equal to 8.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"derived_key_length_bytes": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the derived key. The value must be greater than or equal to 8 and less than or equal to 512. When the `type` attribute is set to `pbkdf2`: Specifies the number of bytes to use for the derived key. The value must be greater than or equal to 8.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the derived key. The value must be greater than or equal to 8 and less than or equal to 512.\n  - `pbkdf2`: Specifies the number of bytes to use for the derived key. The value must be greater than or equal to 8.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Storage Scheme",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`salted-sha256`, `argon2d`, `crypt`, `argon2i`, `aes`, `argon2id`, `vault`, `third-party`, `argon2`, `third-party-enhanced`, `pbkdf2`, `salted-sha384`, `aes-256`, `bcrypt`, `blowfish`, `amazon-secrets-manager`, `azure-key-vault`, `conjur`, `salted-sha512`, `scrypt`]: Indicates whether the Password Storage Scheme is enabled for use. When the `type` attribute is set to `base64`: Indicates whether the Base64 Password Storage Scheme is enabled for use. When the `type` attribute is set to `salted-md5`: Indicates whether the Salted MD5 Password Storage Scheme is enabled for use. When the `type` attribute is set to `rc4`: Indicates whether the RC4 Password Storage Scheme is enabled for use. When the `type` attribute is set to `triple-des`: Indicates whether the Triple DES Password Storage Scheme is enabled for use. When the `type` attribute is set to `clear`: Indicates whether the Clear Password Storage Scheme is enabled for use. When the `type` attribute is set to `sha1`: Indicates whether the SHA1 Password Storage Scheme is enabled for use. When the `type` attribute is set to `salted-sha1`: Indicates whether the Salted SHA1 Password Storage Scheme is enabled for use. When the `type` attribute is set to `md5`: Indicates whether the MD5 Password Storage Scheme is enabled for use.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`salted-sha256`, `argon2d`, `crypt`, `argon2i`, `aes`, `argon2id`, `vault`, `third-party`, `argon2`, `third-party-enhanced`, `pbkdf2`, `salted-sha384`, `aes-256`, `bcrypt`, `blowfish`, `amazon-secrets-manager`, `azure-key-vault`, `conjur`, `salted-sha512`, `scrypt`]: Indicates whether the Password Storage Scheme is enabled for use.\n  - `base64`: Indicates whether the Base64 Password Storage Scheme is enabled for use.\n  - `salted-md5`: Indicates whether the Salted MD5 Password Storage Scheme is enabled for use.\n  - `rc4`: Indicates whether the RC4 Password Storage Scheme is enabled for use.\n  - `triple-des`: Indicates whether the Triple DES Password Storage Scheme is enabled for use.\n  - `clear`: Indicates whether the Clear Password Storage Scheme is enabled for use.\n  - `sha1`: Indicates whether the SHA1 Password Storage Scheme is enabled for use.\n  - `salted-sha1`: Indicates whether the Salted SHA1 Password Storage Scheme is enabled for use.\n  - `md5`: Indicates whether the MD5 Password Storage Scheme is enabled for use.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SaltedSha256PasswordStorageSchemeResponse object into the model struct
func readSaltedSha256PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.SaltedSha256PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha256")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Argon2dPasswordStorageSchemeResponse object into the model struct
func readArgon2dPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Argon2dPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2d")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a CryptPasswordStorageSchemeResponse object into the model struct
func readCryptPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.CryptPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("crypt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordEncodingMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordStorageSchemePasswordEncodingMechanismProp(r.PasswordEncodingMechanism), false)
	state.NumDigestRounds = internaltypes.Int64TypeOrNil(r.NumDigestRounds)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Argon2iPasswordStorageSchemeResponse object into the model struct
func readArgon2iPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Argon2iPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2i")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Base64PasswordStorageSchemeResponse object into the model struct
func readBase64PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Base64PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("base64")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SaltedMd5PasswordStorageSchemeResponse object into the model struct
func readSaltedMd5PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.SaltedMd5PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AesPasswordStorageSchemeResponse object into the model struct
func readAesPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.AesPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aes")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Argon2idPasswordStorageSchemeResponse object into the model struct
func readArgon2idPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Argon2idPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2id")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a VaultPasswordStorageSchemeResponse object into the model struct
func readVaultPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.VaultPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultExternalServer = types.StringValue(r.VaultExternalServer)
	state.DefaultField = internaltypes.StringTypeOrNil(r.DefaultField, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyPasswordStorageSchemeResponse object into the model struct
func readThirdPartyPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.ThirdPartyPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Argon2PasswordStorageSchemeResponse object into the model struct
func readArgon2PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Argon2PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyEnhancedPasswordStorageSchemeResponse object into the model struct
func readThirdPartyEnhancedPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.ThirdPartyEnhancedPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-enhanced")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Pbkdf2PasswordStorageSchemeResponse object into the model struct
func readPbkdf2PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Pbkdf2PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pbkdf2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DigestAlgorithm = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordStorageSchemeDigestAlgorithmProp(r.DigestAlgorithm), false)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Rc4PasswordStorageSchemeResponse object into the model struct
func readRc4PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Rc4PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("rc4")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SaltedSha384PasswordStorageSchemeResponse object into the model struct
func readSaltedSha384PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.SaltedSha384PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha384")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a TripleDesPasswordStorageSchemeResponse object into the model struct
func readTripleDesPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.TripleDesPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("triple-des")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a ClearPasswordStorageSchemeResponse object into the model struct
func readClearPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.ClearPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("clear")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a Aes256PasswordStorageSchemeResponse object into the model struct
func readAes256PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Aes256PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aes-256")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a BcryptPasswordStorageSchemeResponse object into the model struct
func readBcryptPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.BcryptPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("bcrypt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BcryptCostFactor = internaltypes.Int64TypeOrNil(r.BcryptCostFactor)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a BlowfishPasswordStorageSchemeResponse object into the model struct
func readBlowfishPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.BlowfishPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("blowfish")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Sha1PasswordStorageSchemeResponse object into the model struct
func readSha1PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Sha1PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("sha1")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a AmazonSecretsManagerPasswordStorageSchemeResponse object into the model struct
func readAmazonSecretsManagerPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.AmazonSecretsManagerPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-secrets-manager")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.DefaultField = internaltypes.StringTypeOrNil(r.DefaultField, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AzureKeyVaultPasswordStorageSchemeResponse object into the model struct
func readAzureKeyVaultPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.AzureKeyVaultPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("azure-key-vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ConjurPasswordStorageSchemeResponse object into the model struct
func readConjurPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.ConjurPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a SaltedSha1PasswordStorageSchemeResponse object into the model struct
func readSaltedSha1PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.SaltedSha1PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha1")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read a SaltedSha512PasswordStorageSchemeResponse object into the model struct
func readSaltedSha512PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.SaltedSha512PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha512")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ScryptPasswordStorageSchemeResponse object into the model struct
func readScryptPasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.ScryptPasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("scrypt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScryptCpuMemoryCostFactorExponent = internaltypes.Int64TypeOrNil(r.ScryptCpuMemoryCostFactorExponent)
	state.ScryptBlockSize = internaltypes.Int64TypeOrNil(r.ScryptBlockSize)
	state.ScryptParallelizationParameter = internaltypes.Int64TypeOrNil(r.ScryptParallelizationParameter)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Md5PasswordStorageSchemeResponse object into the model struct
func readMd5PasswordStorageSchemeResponseDataSource(ctx context.Context, r *client.Md5PasswordStorageSchemeResponse, state *passwordStorageSchemeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
}

// Read resource information
func (r *passwordStorageSchemeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passwordStorageSchemeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.GetPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Storage Scheme", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SaltedSha256PasswordStorageSchemeResponse != nil {
		readSaltedSha256PasswordStorageSchemeResponseDataSource(ctx, readResponse.SaltedSha256PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2dPasswordStorageSchemeResponse != nil {
		readArgon2dPasswordStorageSchemeResponseDataSource(ctx, readResponse.Argon2dPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.CryptPasswordStorageSchemeResponse != nil {
		readCryptPasswordStorageSchemeResponseDataSource(ctx, readResponse.CryptPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2iPasswordStorageSchemeResponse != nil {
		readArgon2iPasswordStorageSchemeResponseDataSource(ctx, readResponse.Argon2iPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Base64PasswordStorageSchemeResponse != nil {
		readBase64PasswordStorageSchemeResponseDataSource(ctx, readResponse.Base64PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedMd5PasswordStorageSchemeResponse != nil {
		readSaltedMd5PasswordStorageSchemeResponseDataSource(ctx, readResponse.SaltedMd5PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AesPasswordStorageSchemeResponse != nil {
		readAesPasswordStorageSchemeResponseDataSource(ctx, readResponse.AesPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2idPasswordStorageSchemeResponse != nil {
		readArgon2idPasswordStorageSchemeResponseDataSource(ctx, readResponse.Argon2idPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VaultPasswordStorageSchemeResponse != nil {
		readVaultPasswordStorageSchemeResponseDataSource(ctx, readResponse.VaultPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordStorageSchemeResponse != nil {
		readThirdPartyPasswordStorageSchemeResponseDataSource(ctx, readResponse.ThirdPartyPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2PasswordStorageSchemeResponse != nil {
		readArgon2PasswordStorageSchemeResponseDataSource(ctx, readResponse.Argon2PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse != nil {
		readThirdPartyEnhancedPasswordStorageSchemeResponseDataSource(ctx, readResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Pbkdf2PasswordStorageSchemeResponse != nil {
		readPbkdf2PasswordStorageSchemeResponseDataSource(ctx, readResponse.Pbkdf2PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Rc4PasswordStorageSchemeResponse != nil {
		readRc4PasswordStorageSchemeResponseDataSource(ctx, readResponse.Rc4PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha384PasswordStorageSchemeResponse != nil {
		readSaltedSha384PasswordStorageSchemeResponseDataSource(ctx, readResponse.SaltedSha384PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.TripleDesPasswordStorageSchemeResponse != nil {
		readTripleDesPasswordStorageSchemeResponseDataSource(ctx, readResponse.TripleDesPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ClearPasswordStorageSchemeResponse != nil {
		readClearPasswordStorageSchemeResponseDataSource(ctx, readResponse.ClearPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Aes256PasswordStorageSchemeResponse != nil {
		readAes256PasswordStorageSchemeResponseDataSource(ctx, readResponse.Aes256PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BcryptPasswordStorageSchemeResponse != nil {
		readBcryptPasswordStorageSchemeResponseDataSource(ctx, readResponse.BcryptPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.BlowfishPasswordStorageSchemeResponse != nil {
		readBlowfishPasswordStorageSchemeResponseDataSource(ctx, readResponse.BlowfishPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Sha1PasswordStorageSchemeResponse != nil {
		readSha1PasswordStorageSchemeResponseDataSource(ctx, readResponse.Sha1PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerPasswordStorageSchemeResponse != nil {
		readAmazonSecretsManagerPasswordStorageSchemeResponseDataSource(ctx, readResponse.AmazonSecretsManagerPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultPasswordStorageSchemeResponse != nil {
		readAzureKeyVaultPasswordStorageSchemeResponseDataSource(ctx, readResponse.AzureKeyVaultPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurPasswordStorageSchemeResponse != nil {
		readConjurPasswordStorageSchemeResponseDataSource(ctx, readResponse.ConjurPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha1PasswordStorageSchemeResponse != nil {
		readSaltedSha1PasswordStorageSchemeResponseDataSource(ctx, readResponse.SaltedSha1PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha512PasswordStorageSchemeResponse != nil {
		readSaltedSha512PasswordStorageSchemeResponseDataSource(ctx, readResponse.SaltedSha512PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ScryptPasswordStorageSchemeResponse != nil {
		readScryptPasswordStorageSchemeResponseDataSource(ctx, readResponse.ScryptPasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Md5PasswordStorageSchemeResponse != nil {
		readMd5PasswordStorageSchemeResponseDataSource(ctx, readResponse.Md5PasswordStorageSchemeResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
