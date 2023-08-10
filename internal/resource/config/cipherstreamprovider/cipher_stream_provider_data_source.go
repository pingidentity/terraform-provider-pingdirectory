package cipherstreamprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &cipherStreamProviderDataSource{}
	_ datasource.DataSourceWithConfigure = &cipherStreamProviderDataSource{}
)

// Create a Cipher Stream Provider data source
func NewCipherStreamProviderDataSource() datasource.DataSource {
	return &cipherStreamProviderDataSource{}
}

// cipherStreamProviderDataSource is the datasource implementation.
type cipherStreamProviderDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *cipherStreamProviderDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cipher_stream_provider"
}

// Configure adds the provider configured client to the data source.
func (r *cipherStreamProviderDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type cipherStreamProviderDataSourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	VaultExternalServer             types.String `tfsdk:"vault_external_server"`
	VaultServerBaseURI              types.Set    `tfsdk:"vault_server_base_uri"`
	VaultAuthenticationMethod       types.String `tfsdk:"vault_authentication_method"`
	VaultSecretPath                 types.String `tfsdk:"vault_secret_path"`
	VaultSecretFieldName            types.String `tfsdk:"vault_secret_field_name"`
	VaultEncryptionMetadataFile     types.String `tfsdk:"vault_encryption_metadata_file"`
	TrustStoreFile                  types.String `tfsdk:"trust_store_file"`
	TrustStorePin                   types.String `tfsdk:"trust_store_pin"`
	TrustStoreType                  types.String `tfsdk:"trust_store_type"`
	Pkcs11ProviderClass             types.String `tfsdk:"pkcs11_provider_class"`
	Pkcs11ProviderConfigurationFile types.String `tfsdk:"pkcs11_provider_configuration_file"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinEnvironmentVariable  types.String `tfsdk:"key_store_pin_environment_variable"`
	Pkcs11KeyStoreType              types.String `tfsdk:"pkcs11_key_store_type"`
	SslCertNickname                 types.String `tfsdk:"ssl_cert_nickname"`
	ConjurExternalServer            types.String `tfsdk:"conjur_external_server"`
	ConjurSecretRelativePath        types.String `tfsdk:"conjur_secret_relative_path"`
	PasswordFile                    types.String `tfsdk:"password_file"`
	WaitForPasswordFile             types.Bool   `tfsdk:"wait_for_password_file"`
	KeyVaultURI                     types.String `tfsdk:"key_vault_uri"`
	AzureAuthenticationMethod       types.String `tfsdk:"azure_authentication_method"`
	HttpProxyExternalServer         types.String `tfsdk:"http_proxy_external_server"`
	SecretName                      types.String `tfsdk:"secret_name"`
	EncryptedPassphraseFile         types.String `tfsdk:"encrypted_passphrase_file"`
	SecretID                        types.String `tfsdk:"secret_id"`
	SecretFieldName                 types.String `tfsdk:"secret_field_name"`
	SecretVersionID                 types.String `tfsdk:"secret_version_id"`
	SecretVersionStage              types.String `tfsdk:"secret_version_stage"`
	EncryptionMetadataFile          types.String `tfsdk:"encryption_metadata_file"`
	AwsExternalServer               types.String `tfsdk:"aws_external_server"`
	AwsAccessKeyID                  types.String `tfsdk:"aws_access_key_id"`
	AwsSecretAccessKey              types.String `tfsdk:"aws_secret_access_key"`
	AwsRegionName                   types.String `tfsdk:"aws_region_name"`
	KmsEncryptionKeyArn             types.String `tfsdk:"kms_encryption_key_arn"`
	IterationCount                  types.Int64  `tfsdk:"iteration_count"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *cipherStreamProviderDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Cipher Stream Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Cipher Stream Provider resource. Options are ['amazon-key-management-service', 'amazon-secrets-manager', 'azure-key-vault', 'file-based', 'wait-for-passphrase', 'conjur', 'pkcs11', 'vault', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Cipher Stream Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Cipher Stream Provider. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_server_base_uri": schema.SetAttribute{
				Description: "The base URL needed to access the Vault server. The base URL should consist of the protocol (\"http\" or \"https\"), the server address (resolvable name or IP address), and the port number. For example, \"https://vault.example.com:8200/\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"vault_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Vault server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_secret_path": schema.StringAttribute{
				Description: "The path to the desired secret in the Vault service. This will be appended to the value of the base-url property for the associated Vault external server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_secret_field_name": schema.StringAttribute{
				Description: "The name of the field in the Vault secret record that contains the passphrase to use to generate the encryption key.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"vault_encryption_metadata_file": schema.StringAttribute{
				Description: "The path to a file that will hold metadata about the encryption performed by this Vault Cipher Stream Provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_file": schema.StringAttribute{
				Description: "The path to a file containing the information needed to trust the certificate presented by the Vault servers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"trust_store_pin": schema.StringAttribute{
				Description: "The passphrase needed to access the contents of the trust store. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"trust_store_type": schema.StringAttribute{
				Description: "The store type for the specified trust store file. The value should likely be one of \"JKS\" or \"PKCS12\".",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pkcs11_provider_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java security provider class that implements support for interacting with PKCS #11 tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pkcs11_provider_configuration_file": schema.StringAttribute{
				Description: "The path to the file to use to configure the security provider that implements support for interacting with PKCS #11 tokens.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_store_pin": schema.StringAttribute{
				Description: "The clear-text user PIN needed to interact with the PKCS #11 token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description: "The path to a file containing the user PIN needed to interact with the PKCS #11 token. The file must exist and must contain exactly one line with a clear-text representation of the PIN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"key_store_pin_environment_variable": schema.StringAttribute{
				Description: "The name of an environment variable whose value is the user PIN needed to interact with the PKCS #11 token. The environment variable must be defined and must contain a clear-text representation of the PIN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"pkcs11_key_store_type": schema.StringAttribute{
				Description: "The key store type to use when obtaining an instance of a key store for interacting with a PKCS #11 token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ssl_cert_nickname": schema.StringAttribute{
				Description: "The alias for the certificate in the PKCS #11 token that will be used to wrap the encryption key. The target certificate must exist in the PKCS #11 token, and it must have an RSA key pair because the JVM does not currently provide adequate key wrapping support for elliptic curve key pairs.  If you have also configured the server to use a PKCS #11 token for accessing listener certificates, we strongly recommend that you use a different certificate to protect the contents of the encryption settings database than you use for negotiating TLS sessions with clients. It is imperative that the certificate used by this PKCS11 Cipher Stream Provider remain constant for the life of the provider because if the certificate were to be replaced, then the contents of the encryption settings database could become inaccessible. Unlike with listener certificates used for TLS negotiation that need to be replaced on a regular basis, this PKCS11 Cipher Stream Provider does not consider the validity period for the associated certificate, and it will continue to function even after the certificate has expired.  If you need to rotate the certificate used to protect the server's encryption settings database, you should first install the desired new certificate in the PKCS #11 token under a different alias. Then, you should create a new instance of this PKCS11 Cipher Stream Provider that is configured to use that certificate, and that also uses a different value for the encryption-metadata-file because the information in that file is tied to the certificate used to generate it. Finally, you will need to update the global configuration so that the encryption-settings-cipher-stream-provider property references the new cipher stream provider rather than this one. The update to the global configuration must be done with the server online so that it can properly re-encrypt the contents of the encryption settings database with the correct key tied to the new certificate.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"conjur_secret_relative_path": schema.StringAttribute{
				Description: "The portion of the path that follows the account name in the URI needed to obtain the secret passphrase to use to generate the encryption key. Any special characters in the path must be URL-encoded.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_file": schema.StringAttribute{
				Description: "The path to the file containing the password to use when generating ciphers.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"wait_for_password_file": schema.BoolAttribute{
				Description: "Indicates whether the server should wait for the password file to become available if it does not exist.",
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
			"secret_name": schema.StringAttribute{
				Description: "The name of the secret to retrieve.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encrypted_passphrase_file": schema.StringAttribute{
				Description: "The path to a file that will hold the encrypted passphrase used by this cipher stream provider.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_id": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_field_name": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the passphrase that will be used to generate the encryption key for protecting the contents of the encryption settings database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_version_id": schema.StringAttribute{
				Description: "The unique identifier for the version of the secret to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secret_version_stage": schema.StringAttribute{
				Description: "The staging label for the version of the secret to be retrieved.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"encryption_metadata_file": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `amazon-secrets-manager`: The path to a file that will hold metadata about the encryption performed by this Amazon Secrets Manager Cipher Stream Provider. When the `type` attribute is set to `azure-key-vault`: The path to a file that will hold metadata about the encryption performed by this Azure Key Vault Cipher Stream Provider. When the `type` attribute is set to `file-based`: The path to a file that will hold metadata about the encryption performed by this File Based Cipher Stream Provider. When the `type` attribute is set to `conjur`: The path to a file that will hold metadata about the encryption performed by this Conjur Cipher Stream Provider. When the `type` attribute is set to `pkcs11`: The path to a file that will hold metadata about the encryption performed by this PKCS11 Cipher Stream Provider.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `amazon-secrets-manager`: The path to a file that will hold metadata about the encryption performed by this Amazon Secrets Manager Cipher Stream Provider.\n  - `azure-key-vault`: The path to a file that will hold metadata about the encryption performed by this Azure Key Vault Cipher Stream Provider.\n  - `file-based`: The path to a file that will hold metadata about the encryption performed by this File Based Cipher Stream Provider.\n  - `conjur`: The path to a file that will hold metadata about the encryption performed by this Conjur Cipher Stream Provider.\n  - `pkcs11`: The path to a file that will hold metadata about the encryption performed by this PKCS11 Cipher Stream Provider.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"aws_external_server": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `amazon-key-management-service`: The external server with information to use when interacting with the Amazon Key Management Service. When the `type` attribute is set to `amazon-secrets-manager`: The external server with information to use when interacting with the AWS Secrets Manager.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `amazon-key-management-service`: The external server with information to use when interacting with the Amazon Key Management Service.\n  - `amazon-secrets-manager`: The external server with information to use when interacting with the AWS Secrets Manager.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"aws_access_key_id": schema.StringAttribute{
				Description: "The access key ID that will be used if this cipher stream provider will authenticate to the Amazon Key Management Service using an access key rather than an IAM role associated with an EC2 instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"aws_secret_access_key": schema.StringAttribute{
				Description: "The secret access key that will be used if this cipher stream provider will authenticate to the Amazon Key Management Service using an access key rather than an IAM role associated with an EC2 instance.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				Sensitive:   true,
			},
			"aws_region_name": schema.StringAttribute{
				Description: "The name of the Amazon Web Services region that holds the encryption key. This is optional, and if it is not provided, then the server will attempt to determine the region from the key ARN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"kms_encryption_key_arn": schema.StringAttribute{
				Description: "The Amazon resource name (ARN) for the KMS key that will be used to encrypt the contents of the passphrase file. This key must exist, and the AWS client must have access to encrypt and decrypt data using this key.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"iteration_count": schema.Int64Attribute{
				Description: "Supported in PingDirectory product version 9.3.0.0+. The PBKDF2 iteration count that will be used when deriving the encryption key used to protect the encryption settings database.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Cipher Stream Provider",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Cipher Stream Provider is enabled for use in the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a AmazonKeyManagementServiceCipherStreamProviderResponse object into the model struct
func readAmazonKeyManagementServiceCipherStreamProviderResponseDataSource(ctx context.Context, r *client.AmazonKeyManagementServiceCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-key-management-service")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EncryptedPassphraseFile = types.StringValue(r.EncryptedPassphraseFile)
	state.AwsExternalServer = internaltypes.StringTypeOrNil(r.AwsExternalServer, false)
	state.AwsAccessKeyID = internaltypes.StringTypeOrNil(r.AwsAccessKeyID, false)
	state.AwsRegionName = internaltypes.StringTypeOrNil(r.AwsRegionName, false)
	state.KmsEncryptionKeyArn = types.StringValue(r.KmsEncryptionKeyArn)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AmazonSecretsManagerCipherStreamProviderResponse object into the model struct
func readAmazonSecretsManagerCipherStreamProviderResponseDataSource(ctx context.Context, r *client.AmazonSecretsManagerCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-secrets-manager")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.SecretID = types.StringValue(r.SecretID)
	state.SecretFieldName = types.StringValue(r.SecretFieldName)
	state.SecretVersionID = internaltypes.StringTypeOrNil(r.SecretVersionID, false)
	state.SecretVersionStage = internaltypes.StringTypeOrNil(r.SecretVersionStage, false)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AzureKeyVaultCipherStreamProviderResponse object into the model struct
func readAzureKeyVaultCipherStreamProviderResponseDataSource(ctx context.Context, r *client.AzureKeyVaultCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("azure-key-vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, false)
	state.SecretName = types.StringValue(r.SecretName)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a FileBasedCipherStreamProviderResponse object into the model struct
func readFileBasedCipherStreamProviderResponseDataSource(ctx context.Context, r *client.FileBasedCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordFile = types.StringValue(r.PasswordFile)
	state.WaitForPasswordFile = internaltypes.BoolTypeOrNil(r.WaitForPasswordFile)
	state.EncryptionMetadataFile = internaltypes.StringTypeOrNil(r.EncryptionMetadataFile, false)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a WaitForPassphraseCipherStreamProviderResponse object into the model struct
func readWaitForPassphraseCipherStreamProviderResponseDataSource(ctx context.Context, r *client.WaitForPassphraseCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("wait-for-passphrase")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ConjurCipherStreamProviderResponse object into the model struct
func readConjurCipherStreamProviderResponseDataSource(ctx context.Context, r *client.ConjurCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.ConjurSecretRelativePath = types.StringValue(r.ConjurSecretRelativePath)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a Pkcs11CipherStreamProviderResponse object into the model struct
func readPkcs11CipherStreamProviderResponseDataSource(ctx context.Context, r *client.Pkcs11CipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pkcs11")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Pkcs11ProviderClass = internaltypes.StringTypeOrNil(r.Pkcs11ProviderClass, false)
	state.Pkcs11ProviderConfigurationFile = internaltypes.StringTypeOrNil(r.Pkcs11ProviderConfigurationFile, false)
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, false)
	state.KeyStorePinEnvironmentVariable = internaltypes.StringTypeOrNil(r.KeyStorePinEnvironmentVariable, false)
	state.Pkcs11KeyStoreType = internaltypes.StringTypeOrNil(r.Pkcs11KeyStoreType, false)
	state.SslCertNickname = types.StringValue(r.SslCertNickname)
	state.EncryptionMetadataFile = types.StringValue(r.EncryptionMetadataFile)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a VaultCipherStreamProviderResponse object into the model struct
func readVaultCipherStreamProviderResponseDataSource(ctx context.Context, r *client.VaultCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultExternalServer = internaltypes.StringTypeOrNil(r.VaultExternalServer, false)
	state.VaultServerBaseURI = internaltypes.GetStringSet(r.VaultServerBaseURI)
	state.VaultAuthenticationMethod = internaltypes.StringTypeOrNil(r.VaultAuthenticationMethod, false)
	state.VaultSecretPath = types.StringValue(r.VaultSecretPath)
	state.VaultSecretFieldName = types.StringValue(r.VaultSecretFieldName)
	state.VaultEncryptionMetadataFile = types.StringValue(r.VaultEncryptionMetadataFile)
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, false)
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, false)
	state.IterationCount = internaltypes.Int64TypeOrNil(r.IterationCount)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyCipherStreamProviderResponse object into the model struct
func readThirdPartyCipherStreamProviderResponseDataSource(ctx context.Context, r *client.ThirdPartyCipherStreamProviderResponse, state *cipherStreamProviderDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *cipherStreamProviderDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state cipherStreamProviderDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CipherStreamProviderApi.GetCipherStreamProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Cipher Stream Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse != nil {
		readAmazonKeyManagementServiceCipherStreamProviderResponseDataSource(ctx, readResponse.AmazonKeyManagementServiceCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerCipherStreamProviderResponse != nil {
		readAmazonSecretsManagerCipherStreamProviderResponseDataSource(ctx, readResponse.AmazonSecretsManagerCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultCipherStreamProviderResponse != nil {
		readAzureKeyVaultCipherStreamProviderResponseDataSource(ctx, readResponse.AzureKeyVaultCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedCipherStreamProviderResponse != nil {
		readFileBasedCipherStreamProviderResponseDataSource(ctx, readResponse.FileBasedCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.WaitForPassphraseCipherStreamProviderResponse != nil {
		readWaitForPassphraseCipherStreamProviderResponseDataSource(ctx, readResponse.WaitForPassphraseCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurCipherStreamProviderResponse != nil {
		readConjurCipherStreamProviderResponseDataSource(ctx, readResponse.ConjurCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.Pkcs11CipherStreamProviderResponse != nil {
		readPkcs11CipherStreamProviderResponseDataSource(ctx, readResponse.Pkcs11CipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.VaultCipherStreamProviderResponse != nil {
		readVaultCipherStreamProviderResponseDataSource(ctx, readResponse.VaultCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyCipherStreamProviderResponse != nil {
		readThirdPartyCipherStreamProviderResponseDataSource(ctx, readResponse.ThirdPartyCipherStreamProviderResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
