// Copyright © 2025 Ping Identity Corporation

package passwordstoragescheme

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passwordStorageSchemeResource{}
	_ resource.ResourceWithConfigure   = &passwordStorageSchemeResource{}
	_ resource.ResourceWithImportState = &passwordStorageSchemeResource{}
	_ resource.Resource                = &defaultPasswordStorageSchemeResource{}
	_ resource.ResourceWithConfigure   = &defaultPasswordStorageSchemeResource{}
	_ resource.ResourceWithImportState = &defaultPasswordStorageSchemeResource{}
)

// Create a Password Storage Scheme resource
func NewPasswordStorageSchemeResource() resource.Resource {
	return &passwordStorageSchemeResource{}
}

func NewDefaultPasswordStorageSchemeResource() resource.Resource {
	return &defaultPasswordStorageSchemeResource{}
}

// passwordStorageSchemeResource is the resource implementation.
type passwordStorageSchemeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPasswordStorageSchemeResource is the resource implementation.
type defaultPasswordStorageSchemeResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passwordStorageSchemeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_storage_scheme"
}

func (r *defaultPasswordStorageSchemeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_password_storage_scheme"
}

// Configure adds the provider configured client to the resource.
func (r *passwordStorageSchemeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPasswordStorageSchemeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordStorageSchemeResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
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
	EncodedPasswordCacheSize          types.Int64  `tfsdk:"encoded_password_cache_size"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *passwordStorageSchemeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordStorageSchemeSchema(ctx, req, resp, false)
}

func (r *defaultPasswordStorageSchemeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passwordStorageSchemeSchema(ctx, req, resp, true)
}

func passwordStorageSchemeSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Password Storage Scheme.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Storage Scheme resource. Options are ['salted-sha256', 'argon2d', 'crypt', 'argon2i', 'base64', 'salted-md5', 'aes', 'argon2id', 'vault', 'third-party', 'argon2', 'third-party-enhanced', 'pbkdf2', 'rc4', 'salted-sha384', 'triple-des', 'clear', 'aes-256', 'bcrypt', 'blowfish', 'sha1', 'amazon-secrets-manager', 'azure-key-vault', 'conjur', 'salted-sha1', 'salted-sha512', 'scrypt', 'md5']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"salted-sha256", "argon2d", "crypt", "argon2i", "base64", "salted-md5", "aes", "argon2id", "vault", "third-party", "argon2", "third-party-enhanced", "pbkdf2", "rc4", "salted-sha384", "triple-des", "clear", "aes-256", "bcrypt", "blowfish", "sha1", "amazon-secrets-manager", "azure-key-vault", "conjur", "salted-sha1", "salted-sha512", "scrypt", "md5"}...),
				},
			},
			"scrypt_cpu_memory_cost_factor_exponent": schema.Int64Attribute{
				Description: "Specifies the exponent that should be used for the CPU/memory cost factor. The cost factor must be a power of two, so the value of this property represents the power to which two is raised. The CPU/memory cost factor specifies the number of iterations required for encoding the password, and also affects the amount of memory required during processing. A higher cost factor requires more processing and more memory to generate a password, which makes attacks against the password more expensive.",
				Optional:    true,
				Computed:    true,
			},
			"scrypt_block_size": schema.Int64Attribute{
				Description: "Specifies the block size for the digest that will be used in the course of encoding passwords. Increasing the block size while keeping the CPU/memory cost factor constant will increase the amount of memory required to encode a password, but it also increases the ratio of sequential memory access to random memory access (and sequential memory access is generally faster than random memory access).",
				Optional:    true,
				Computed:    true,
			},
			"scrypt_parallelization_parameter": schema.Int64Attribute{
				Description: "Specifies the number of times that scrypt has to perform the entire encoding process to produce the final result.",
				Optional:    true,
				Computed:    true,
			},
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur instance containing user passwords.",
				Optional:    true,
			},
			"key_vault_uri": schema.StringAttribute{
				Description: "The URI that identifies the Azure Key Vault from which the secret is to be retrieved.",
				Optional:    true,
			},
			"azure_authentication_method": schema.StringAttribute{
				Description: "The mechanism used to authenticate to the Azure service.",
				Optional:    true,
			},
			"http_proxy_external_server": schema.StringAttribute{
				Description: "A reference to an HTTP proxy server that should be used for requests sent to the Azure service.",
				Optional:    true,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS Secrets Manager service.",
				Optional:    true,
			},
			"bcrypt_cost_factor": schema.Int64Attribute{
				Description: "Specifies the cost factor to use when encoding passwords with Bcrypt. A higher cost factor requires more processing to generate a password, which makes attacks against the password more expensive.",
				Optional:    true,
				Computed:    true,
			},
			"encryption_settings_definition_id": schema.StringAttribute{
				Description: "The identifier for the encryption settings definition that should be used to derive the encryption key to use when encrypting new passwords. If this is not provided, the server's preferred encryption settings definition will be used.",
				Optional:    true,
			},
			"digest_algorithm": schema.StringAttribute{
				Description: "Specifies the digest algorithm that will be used when encoding passwords.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"sha-1", "sha-256", "sha-384", "sha-512"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `third-party`: The fully-qualified name of the Java class providing the logic for the Third Party Password Storage Scheme. When the `type` attribute is set to `third-party-enhanced`: The fully-qualified name of the Java class providing the logic for the Third Party Enhanced Password Storage Scheme.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party`: The fully-qualified name of the Java class providing the logic for the Third Party Password Storage Scheme.\n  - `third-party-enhanced`: The fully-qualified name of the Java class providing the logic for the Third Party Enhanced Password Storage Scheme.",
				Optional:            true,
			},
			"extension_argument": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `third-party`: The set of arguments used to customize the behavior for the Third Party Password Storage Scheme. Each configuration property should be given in the form 'name=value'. When the `type` attribute is set to `third-party-enhanced`: The set of arguments used to customize the behavior for the Third Party Enhanced Password Storage Scheme. Each configuration property should be given in the form 'name=value'.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `third-party`: The set of arguments used to customize the behavior for the Third Party Password Storage Scheme. Each configuration property should be given in the form 'name=value'.\n  - `third-party-enhanced`: The set of arguments used to customize the behavior for the Third Party Enhanced Password Storage Scheme. Each configuration property should be given in the form 'name=value'.",
				Optional:            true,
				Computed:            true,
				Default:             internaltypes.EmptySetDefault(types.StringType),
				ElementType:         types.StringType,
			},
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault instance containing the passphrase.",
				Optional:    true,
			},
			"default_field": schema.StringAttribute{
				Description: "The default name of the field in JSON objects contained in the AWS Secrets Manager service that contains the password for the target user.",
				Optional:    true,
			},
			"password_encoding_mechanism": schema.StringAttribute{
				Description: "Specifies the mechanism that should be used to encode clear-text passwords for use with this scheme.",
				Optional:    true,
				Computed:    true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"crypt", "md5", "sha-2-256", "sha-2-512"}...),
				},
			},
			"num_digest_rounds": schema.Int64Attribute{
				Description: "Specifies the number of digest rounds to use for the SHA-2 encodings. This will not be used for the legacy or MD5-based encodings.",
				Optional:    true,
				Computed:    true,
			},
			"max_password_length": schema.Int64Attribute{
				Description: "Specifies the maximum allowed length, in bytes, for passwords encoded with this scheme, which can help mitigate denial of service attacks from clients that attempt to bind with very long passwords.",
				Optional:    true,
				Computed:    true,
			},
			"iteration_count": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of rounds of cryptographic processing required in the course of encoding each password. When the `type` attribute is set to `pbkdf2`: Specifies the number of iterations to use when encoding passwords. The value must be greater than or equal to 1000.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of rounds of cryptographic processing required in the course of encoding each password.\n  - `pbkdf2`: Specifies the number of iterations to use when encoding passwords. The value must be greater than or equal to 1000.",
				Optional:            true,
				Computed:            true,
			},
			"parallelism_factor": schema.Int64Attribute{
				Description: "The number of concurrent threads that will be used in the course of encoding each password.",
				Optional:    true,
			},
			"memory_usage_kb": schema.Int64Attribute{
				Description: "The number of kilobytes of memory that must be used in the course of encoding each password.",
				Optional:    true,
			},
			"salt_length_bytes": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`salted-sha256`, `salted-md5`, `salted-sha384`, `salted-sha1`, `salted-sha512`]: Specifies the number of bytes to use for the generated salt. When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the generated salt. When the `type` attribute is set to `pbkdf2`: Specifies the number of bytes to use for the generated salt. The value must be greater than or equal to 8.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`salted-sha256`, `salted-md5`, `salted-sha384`, `salted-sha1`, `salted-sha512`]: Specifies the number of bytes to use for the generated salt.\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the generated salt.\n  - `pbkdf2`: Specifies the number of bytes to use for the generated salt. The value must be greater than or equal to 8.",
				Optional:            true,
				Computed:            true,
			},
			"derived_key_length_bytes": schema.Int64Attribute{
				Description:         "When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the derived key. The value must be greater than or equal to 8 and less than or equal to 512. When the `type` attribute is set to `pbkdf2`: Specifies the number of bytes to use for the derived key. The value must be greater than or equal to 8.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The number of bytes to use for the derived key. The value must be greater than or equal to 8 and less than or equal to 512.\n  - `pbkdf2`: Specifies the number of bytes to use for the derived key. The value must be greater than or equal to 8.",
				Optional:            true,
				Computed:            true,
			},
			"encoded_password_cache_size": schema.Int64Attribute{
				Description:         "Supported in PingDirectory product version 10.2.0.0+. When the `type` attribute is set to  one of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The maximum number of Argon2-encoded passwords to cache for faster verification. When the `type` attribute is set to `pbkdf2`: The maximum number of PBKDF2-encoded passwords to cache for faster verification. When the `type` attribute is set to `bcrypt`: The maximum number of Bcrypt-encoded passwords to cache for faster verification. When the `type` attribute is set to `scrypt`: The maximum number of scrypt-encoded passwords to cache for faster verification.",
				MarkdownDescription: "Supported in PingDirectory product version 10.2.0.0+. When the `type` attribute is set to:\n  - One of [`argon2d`, `argon2i`, `argon2id`, `argon2`]: The maximum number of Argon2-encoded passwords to cache for faster verification.\n  - `pbkdf2`: The maximum number of PBKDF2-encoded passwords to cache for faster verification.\n  - `bcrypt`: The maximum number of Bcrypt-encoded passwords to cache for faster verification.\n  - `scrypt`: The maximum number of scrypt-encoded passwords to cache for faster verification.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Storage Scheme",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description:         "When the `type` attribute is set to  one of [`salted-sha256`, `argon2d`, `crypt`, `argon2i`, `aes`, `argon2id`, `vault`, `third-party`, `argon2`, `third-party-enhanced`, `pbkdf2`, `salted-sha384`, `aes-256`, `bcrypt`, `blowfish`, `amazon-secrets-manager`, `azure-key-vault`, `conjur`, `salted-sha512`, `scrypt`]: Indicates whether the Password Storage Scheme is enabled for use. When the `type` attribute is set to `base64`: Indicates whether the Base64 Password Storage Scheme is enabled for use. When the `type` attribute is set to `salted-md5`: Indicates whether the Salted MD5 Password Storage Scheme is enabled for use. When the `type` attribute is set to `rc4`: Indicates whether the RC4 Password Storage Scheme is enabled for use. When the `type` attribute is set to `triple-des`: Indicates whether the Triple DES Password Storage Scheme is enabled for use. When the `type` attribute is set to `clear`: Indicates whether the Clear Password Storage Scheme is enabled for use. When the `type` attribute is set to `sha1`: Indicates whether the SHA1 Password Storage Scheme is enabled for use. When the `type` attribute is set to `salted-sha1`: Indicates whether the Salted SHA1 Password Storage Scheme is enabled for use. When the `type` attribute is set to `md5`: Indicates whether the MD5 Password Storage Scheme is enabled for use.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`salted-sha256`, `argon2d`, `crypt`, `argon2i`, `aes`, `argon2id`, `vault`, `third-party`, `argon2`, `third-party-enhanced`, `pbkdf2`, `salted-sha384`, `aes-256`, `bcrypt`, `blowfish`, `amazon-secrets-manager`, `azure-key-vault`, `conjur`, `salted-sha512`, `scrypt`]: Indicates whether the Password Storage Scheme is enabled for use.\n  - `base64`: Indicates whether the Base64 Password Storage Scheme is enabled for use.\n  - `salted-md5`: Indicates whether the Salted MD5 Password Storage Scheme is enabled for use.\n  - `rc4`: Indicates whether the RC4 Password Storage Scheme is enabled for use.\n  - `triple-des`: Indicates whether the Triple DES Password Storage Scheme is enabled for use.\n  - `clear`: Indicates whether the Clear Password Storage Scheme is enabled for use.\n  - `sha1`: Indicates whether the SHA1 Password Storage Scheme is enabled for use.\n  - `salted-sha1`: Indicates whether the Salted SHA1 Password Storage Scheme is enabled for use.\n  - `md5`: Indicates whether the MD5 Password Storage Scheme is enabled for use.",
				Required:            true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{
			stringplanmodifier.UseStateForUnknown(),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	} else {
		// Add RequiresReplace modifier for read-only attributes
		extensionClassAttr := schemaDef.Attributes["extension_class"].(schema.StringAttribute)
		extensionClassAttr.PlanModifiers = append(extensionClassAttr.PlanModifiers, stringplanmodifier.RequiresReplace())
		schemaDef.Attributes["extension_class"] = extensionClassAttr
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan and set any type-specific defaults
func (r *passwordStorageSchemeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPasswordStorageScheme(ctx, req, resp, r.apiClient, r.providerConfig)
	var planModel, configModel passwordStorageSchemeResourceModel
	req.Config.Get(ctx, &configModel)
	req.Plan.Get(ctx, &planModel)
	resourceType := planModel.Type.ValueString()
	anyDefaultsSet := false
	// Set defaults for crypt type
	if resourceType == "crypt" {
		if !internaltypes.IsDefined(configModel.PasswordEncodingMechanism) {
			defaultVal := types.StringValue("sha-2-256")
			if !planModel.PasswordEncodingMechanism.Equal(defaultVal) {
				planModel.PasswordEncodingMechanism = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.NumDigestRounds) {
			defaultVal := types.Int64Value(5000)
			if !planModel.NumDigestRounds.Equal(defaultVal) {
				planModel.NumDigestRounds = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxPasswordLength) {
			defaultVal := types.Int64Value(200)
			if !planModel.MaxPasswordLength.Equal(defaultVal) {
				planModel.MaxPasswordLength = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for pbkdf2 type
	if resourceType == "pbkdf2" {
		if !internaltypes.IsDefined(configModel.DigestAlgorithm) {
			defaultVal := types.StringValue("sha-1")
			if !planModel.DigestAlgorithm.Equal(defaultVal) {
				planModel.DigestAlgorithm = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.IterationCount) {
			defaultVal := types.Int64Value(10000)
			if !planModel.IterationCount.Equal(defaultVal) {
				planModel.IterationCount = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.SaltLengthBytes) {
			defaultVal := types.Int64Value(16)
			if !planModel.SaltLengthBytes.Equal(defaultVal) {
				planModel.SaltLengthBytes = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.DerivedKeyLengthBytes) {
			defaultVal := types.Int64Value(32)
			if !planModel.DerivedKeyLengthBytes.Equal(defaultVal) {
				planModel.DerivedKeyLengthBytes = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxPasswordLength) {
			defaultVal := types.Int64Value(200)
			if !planModel.MaxPasswordLength.Equal(defaultVal) {
				planModel.MaxPasswordLength = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for bcrypt type
	if resourceType == "bcrypt" {
		if !internaltypes.IsDefined(configModel.BcryptCostFactor) {
			defaultVal := types.Int64Value(10)
			if !planModel.BcryptCostFactor.Equal(defaultVal) {
				planModel.BcryptCostFactor = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	// Set defaults for scrypt type
	if resourceType == "scrypt" {
		if !internaltypes.IsDefined(configModel.ScryptCpuMemoryCostFactorExponent) {
			defaultVal := types.Int64Value(14)
			if !planModel.ScryptCpuMemoryCostFactorExponent.Equal(defaultVal) {
				planModel.ScryptCpuMemoryCostFactorExponent = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ScryptBlockSize) {
			defaultVal := types.Int64Value(8)
			if !planModel.ScryptBlockSize.Equal(defaultVal) {
				planModel.ScryptBlockSize = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.ScryptParallelizationParameter) {
			defaultVal := types.Int64Value(1)
			if !planModel.ScryptParallelizationParameter.Equal(defaultVal) {
				planModel.ScryptParallelizationParameter = defaultVal
				anyDefaultsSet = true
			}
		}
		if !internaltypes.IsDefined(configModel.MaxPasswordLength) {
			defaultVal := types.Int64Value(200)
			if !planModel.MaxPasswordLength.Equal(defaultVal) {
				planModel.MaxPasswordLength = defaultVal
				anyDefaultsSet = true
			}
		}
	}
	if anyDefaultsSet {
		planModel.Notifications = types.SetUnknown(types.StringType)
		planModel.RequiredActions = types.SetUnknown(config.GetRequiredActionsObjectType())
	}
	planModel.setNotApplicableAttrsNull()
	resp.Plan.Set(ctx, &planModel)
}

func (r *defaultPasswordStorageSchemeResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPasswordStorageScheme(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanPasswordStorageScheme(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory10200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model passwordStorageSchemeResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsDefined(model.EncodedPasswordCacheSize) {
		resp.Diagnostics.AddError("Attribute 'encoded_password_cache_size' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

func (model *passwordStorageSchemeResourceModel) setNotApplicableAttrsNull() {
	resourceType := model.Type.ValueString()
	// Set any not applicable computed attributes to null for each type
	if resourceType == "argon2d" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "crypt" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.IterationCount = types.Int64Null()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "argon2i" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "argon2id" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "vault" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "third-party" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "argon2" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "third-party-enhanced" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "pbkdf2" {
		model.ScryptBlockSize = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "aes-256" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "bcrypt" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
	}
	if resourceType == "amazon-secrets-manager" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "azure-key-vault" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "conjur" {
		model.ScryptBlockSize = types.Int64Null()
		model.DigestAlgorithm = types.StringNull()
		model.MaxPasswordLength = types.Int64Null()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.EncodedPasswordCacheSize = types.Int64Null()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.ScryptParallelizationParameter = types.Int64Null()
		model.ScryptCpuMemoryCostFactorExponent = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
	if resourceType == "scrypt" {
		model.DigestAlgorithm = types.StringNull()
		model.IterationCount = types.Int64Null()
		model.PasswordEncodingMechanism = types.StringNull()
		model.DerivedKeyLengthBytes = types.Int64Null()
		model.NumDigestRounds = types.Int64Null()
		model.SaltLengthBytes = types.Int64Null()
		model.BcryptCostFactor = types.Int64Null()
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsPasswordStorageScheme() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("salt_length_bytes"),
			path.MatchRoot("type"),
			[]string{"salted-sha256", "argon2d", "argon2i", "salted-md5", "argon2id", "argon2", "pbkdf2", "salted-sha384", "salted-sha1", "salted-sha512"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("iteration_count"),
			path.MatchRoot("type"),
			[]string{"argon2d", "argon2i", "argon2id", "argon2", "pbkdf2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("parallelism_factor"),
			path.MatchRoot("type"),
			[]string{"argon2d", "argon2i", "argon2id", "argon2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("memory_usage_kb"),
			path.MatchRoot("type"),
			[]string{"argon2d", "argon2i", "argon2id", "argon2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("derived_key_length_bytes"),
			path.MatchRoot("type"),
			[]string{"argon2d", "argon2i", "argon2id", "argon2", "pbkdf2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encoded_password_cache_size"),
			path.MatchRoot("type"),
			[]string{"argon2d", "argon2i", "argon2id", "argon2", "pbkdf2", "bcrypt", "scrypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password_encoding_mechanism"),
			path.MatchRoot("type"),
			[]string{"crypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("num_digest_rounds"),
			path.MatchRoot("type"),
			[]string{"crypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_password_length"),
			path.MatchRoot("type"),
			[]string{"crypt", "pbkdf2", "scrypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_external_server"),
			path.MatchRoot("type"),
			[]string{"vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("default_field"),
			path.MatchRoot("type"),
			[]string{"vault", "amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party", "third-party-enhanced"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party", "third-party-enhanced"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("digest_algorithm"),
			path.MatchRoot("type"),
			[]string{"pbkdf2"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encryption_settings_definition_id"),
			path.MatchRoot("type"),
			[]string{"aes-256"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("bcrypt_cost_factor"),
			path.MatchRoot("type"),
			[]string{"bcrypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("aws_external_server"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_vault_uri"),
			path.MatchRoot("type"),
			[]string{"azure-key-vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("azure_authentication_method"),
			path.MatchRoot("type"),
			[]string{"azure-key-vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("http_proxy_external_server"),
			path.MatchRoot("type"),
			[]string{"azure-key-vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conjur_external_server"),
			path.MatchRoot("type"),
			[]string{"conjur"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("scrypt_cpu_memory_cost_factor_exponent"),
			path.MatchRoot("type"),
			[]string{"scrypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("scrypt_block_size"),
			path.MatchRoot("type"),
			[]string{"scrypt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("scrypt_parallelization_parameter"),
			path.MatchRoot("type"),
			[]string{"scrypt"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"argon2d",
			[]path.Expression{path.MatchRoot("salt_length_bytes"), path.MatchRoot("enabled"), path.MatchRoot("iteration_count"), path.MatchRoot("parallelism_factor"), path.MatchRoot("memory_usage_kb"), path.MatchRoot("derived_key_length_bytes")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"argon2i",
			[]path.Expression{path.MatchRoot("salt_length_bytes"), path.MatchRoot("enabled"), path.MatchRoot("iteration_count"), path.MatchRoot("parallelism_factor"), path.MatchRoot("memory_usage_kb"), path.MatchRoot("derived_key_length_bytes")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"argon2id",
			[]path.Expression{path.MatchRoot("salt_length_bytes"), path.MatchRoot("enabled"), path.MatchRoot("iteration_count"), path.MatchRoot("parallelism_factor"), path.MatchRoot("memory_usage_kb"), path.MatchRoot("derived_key_length_bytes")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"argon2",
			[]path.Expression{path.MatchRoot("salt_length_bytes"), path.MatchRoot("enabled"), path.MatchRoot("iteration_count"), path.MatchRoot("parallelism_factor"), path.MatchRoot("memory_usage_kb"), path.MatchRoot("derived_key_length_bytes")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"crypt",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"vault",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("vault_external_server")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party-enhanced",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"pbkdf2",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"aes-256",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"bcrypt",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"amazon-secrets-manager",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("aws_external_server")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"azure-key-vault",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("key_vault_uri"), path.MatchRoot("azure_authentication_method")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"conjur",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("conjur_external_server")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"scrypt",
			[]path.Expression{path.MatchRoot("enabled")},
		),
	}
}

// Add config validators
func (r passwordStorageSchemeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPasswordStorageScheme()
}

// Add config validators
func (r defaultPasswordStorageSchemeResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPasswordStorageScheme()
}

// Add optional fields to create request for argon2d password-storage-scheme
func addOptionalArgon2dPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddArgon2dPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for crypt password-storage-scheme
func addOptionalCryptPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddCryptPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PasswordEncodingMechanism) {
		passwordEncodingMechanism, err := client.NewEnumpasswordStorageSchemePasswordEncodingMechanismPropFromValue(plan.PasswordEncodingMechanism.ValueString())
		if err != nil {
			return err
		}
		addRequest.PasswordEncodingMechanism = passwordEncodingMechanism
	}
	if internaltypes.IsDefined(plan.NumDigestRounds) {
		addRequest.NumDigestRounds = plan.NumDigestRounds.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxPasswordLength) {
		addRequest.MaxPasswordLength = plan.MaxPasswordLength.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for argon2i password-storage-scheme
func addOptionalArgon2iPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddArgon2iPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for argon2id password-storage-scheme
func addOptionalArgon2idPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddArgon2idPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for vault password-storage-scheme
func addOptionalVaultPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddVaultPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultField) {
		addRequest.DefaultField = plan.DefaultField.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party password-storage-scheme
func addOptionalThirdPartyPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddThirdPartyPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for argon2 password-storage-scheme
func addOptionalArgon2PasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddArgon2PasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party-enhanced password-storage-scheme
func addOptionalThirdPartyEnhancedPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddThirdPartyEnhancedPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for pbkdf2 password-storage-scheme
func addOptionalPbkdf2PasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddPbkdf2PasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DigestAlgorithm) {
		digestAlgorithm, err := client.NewEnumpasswordStorageSchemeDigestAlgorithmPropFromValue(plan.DigestAlgorithm.ValueString())
		if err != nil {
			return err
		}
		addRequest.DigestAlgorithm = digestAlgorithm
	}
	if internaltypes.IsDefined(plan.IterationCount) {
		addRequest.IterationCount = plan.IterationCount.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.SaltLengthBytes) {
		addRequest.SaltLengthBytes = plan.SaltLengthBytes.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.DerivedKeyLengthBytes) {
		addRequest.DerivedKeyLengthBytes = plan.DerivedKeyLengthBytes.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxPasswordLength) {
		addRequest.MaxPasswordLength = plan.MaxPasswordLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for aes-256 password-storage-scheme
func addOptionalAes256PasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddAes256PasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionSettingsDefinitionID) {
		addRequest.EncryptionSettingsDefinitionID = plan.EncryptionSettingsDefinitionID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for bcrypt password-storage-scheme
func addOptionalBcryptPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddBcryptPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.BcryptCostFactor) {
		addRequest.BcryptCostFactor = plan.BcryptCostFactor.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for amazon-secrets-manager password-storage-scheme
func addOptionalAmazonSecretsManagerPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddAmazonSecretsManagerPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DefaultField) {
		addRequest.DefaultField = plan.DefaultField.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for azure-key-vault password-storage-scheme
func addOptionalAzureKeyVaultPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddAzureKeyVaultPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for conjur password-storage-scheme
func addOptionalConjurPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddConjurPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for scrypt password-storage-scheme
func addOptionalScryptPasswordStorageSchemeFields(ctx context.Context, addRequest *client.AddScryptPasswordStorageSchemeRequest, plan passwordStorageSchemeResourceModel) error {
	if internaltypes.IsDefined(plan.ScryptCpuMemoryCostFactorExponent) {
		addRequest.ScryptCpuMemoryCostFactorExponent = plan.ScryptCpuMemoryCostFactorExponent.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ScryptBlockSize) {
		addRequest.ScryptBlockSize = plan.ScryptBlockSize.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.ScryptParallelizationParameter) {
		addRequest.ScryptParallelizationParameter = plan.ScryptParallelizationParameter.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.MaxPasswordLength) {
		addRequest.MaxPasswordLength = plan.MaxPasswordLength.ValueInt64Pointer()
	}
	if internaltypes.IsDefined(plan.EncodedPasswordCacheSize) {
		addRequest.EncodedPasswordCacheSize = plan.EncodedPasswordCacheSize.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populatePasswordStorageSchemeUnknownValues(model *passwordStorageSchemeResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *passwordStorageSchemeResourceModel) populateAllComputedStringAttributes() {
	if model.PasswordEncodingMechanism.IsUnknown() || model.PasswordEncodingMechanism.IsNull() {
		model.PasswordEncodingMechanism = types.StringValue("")
	}
	if model.ConjurExternalServer.IsUnknown() || model.ConjurExternalServer.IsNull() {
		model.ConjurExternalServer = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.HttpProxyExternalServer.IsUnknown() || model.HttpProxyExternalServer.IsNull() {
		model.HttpProxyExternalServer = types.StringValue("")
	}
	if model.AwsExternalServer.IsUnknown() || model.AwsExternalServer.IsNull() {
		model.AwsExternalServer = types.StringValue("")
	}
	if model.DefaultField.IsUnknown() || model.DefaultField.IsNull() {
		model.DefaultField = types.StringValue("")
	}
	if model.EncryptionSettingsDefinitionID.IsUnknown() || model.EncryptionSettingsDefinitionID.IsNull() {
		model.EncryptionSettingsDefinitionID = types.StringValue("")
	}
	if model.KeyVaultURI.IsUnknown() || model.KeyVaultURI.IsNull() {
		model.KeyVaultURI = types.StringValue("")
	}
	if model.AzureAuthenticationMethod.IsUnknown() || model.AzureAuthenticationMethod.IsNull() {
		model.AzureAuthenticationMethod = types.StringValue("")
	}
	if model.VaultExternalServer.IsUnknown() || model.VaultExternalServer.IsNull() {
		model.VaultExternalServer = types.StringValue("")
	}
	if model.DigestAlgorithm.IsUnknown() || model.DigestAlgorithm.IsNull() {
		model.DigestAlgorithm = types.StringValue("")
	}
}

// Read a SaltedSha256PasswordStorageSchemeResponse object into the model struct
func readSaltedSha256PasswordStorageSchemeResponse(ctx context.Context, r *client.SaltedSha256PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha256")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Argon2dPasswordStorageSchemeResponse object into the model struct
func readArgon2dPasswordStorageSchemeResponse(ctx context.Context, r *client.Argon2dPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2d")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a CryptPasswordStorageSchemeResponse object into the model struct
func readCryptPasswordStorageSchemeResponse(ctx context.Context, r *client.CryptPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("crypt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordEncodingMechanism = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordStorageSchemePasswordEncodingMechanismProp(r.PasswordEncodingMechanism), true)
	state.NumDigestRounds = internaltypes.Int64TypeOrNil(r.NumDigestRounds)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Argon2iPasswordStorageSchemeResponse object into the model struct
func readArgon2iPasswordStorageSchemeResponse(ctx context.Context, r *client.Argon2iPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2i")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Base64PasswordStorageSchemeResponse object into the model struct
func readBase64PasswordStorageSchemeResponse(ctx context.Context, r *client.Base64PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("base64")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a SaltedMd5PasswordStorageSchemeResponse object into the model struct
func readSaltedMd5PasswordStorageSchemeResponse(ctx context.Context, r *client.SaltedMd5PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a AesPasswordStorageSchemeResponse object into the model struct
func readAesPasswordStorageSchemeResponse(ctx context.Context, r *client.AesPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aes")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Argon2idPasswordStorageSchemeResponse object into the model struct
func readArgon2idPasswordStorageSchemeResponse(ctx context.Context, r *client.Argon2idPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2id")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a VaultPasswordStorageSchemeResponse object into the model struct
func readVaultPasswordStorageSchemeResponse(ctx context.Context, r *client.VaultPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultExternalServer = types.StringValue(r.VaultExternalServer)
	state.DefaultField = internaltypes.StringTypeOrNil(r.DefaultField, internaltypes.IsEmptyString(expectedValues.DefaultField))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a ThirdPartyPasswordStorageSchemeResponse object into the model struct
func readThirdPartyPasswordStorageSchemeResponse(ctx context.Context, r *client.ThirdPartyPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Argon2PasswordStorageSchemeResponse object into the model struct
func readArgon2PasswordStorageSchemeResponse(ctx context.Context, r *client.Argon2PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("argon2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.ParallelismFactor = types.Int64Value(r.ParallelismFactor)
	state.MemoryUsageKb = types.Int64Value(r.MemoryUsageKb)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a ThirdPartyEnhancedPasswordStorageSchemeResponse object into the model struct
func readThirdPartyEnhancedPasswordStorageSchemeResponse(ctx context.Context, r *client.ThirdPartyEnhancedPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party-enhanced")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Pbkdf2PasswordStorageSchemeResponse object into the model struct
func readPbkdf2PasswordStorageSchemeResponse(ctx context.Context, r *client.Pbkdf2PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pbkdf2")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DigestAlgorithm = internaltypes.StringTypeOrNil(
		client.StringPointerEnumpasswordStorageSchemeDigestAlgorithmProp(r.DigestAlgorithm), true)
	state.IterationCount = types.Int64Value(r.IterationCount)
	state.SaltLengthBytes = types.Int64Value(r.SaltLengthBytes)
	state.DerivedKeyLengthBytes = types.Int64Value(r.DerivedKeyLengthBytes)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Rc4PasswordStorageSchemeResponse object into the model struct
func readRc4PasswordStorageSchemeResponse(ctx context.Context, r *client.Rc4PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("rc4")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a SaltedSha384PasswordStorageSchemeResponse object into the model struct
func readSaltedSha384PasswordStorageSchemeResponse(ctx context.Context, r *client.SaltedSha384PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha384")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a TripleDesPasswordStorageSchemeResponse object into the model struct
func readTripleDesPasswordStorageSchemeResponse(ctx context.Context, r *client.TripleDesPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("triple-des")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a ClearPasswordStorageSchemeResponse object into the model struct
func readClearPasswordStorageSchemeResponse(ctx context.Context, r *client.ClearPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("clear")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Aes256PasswordStorageSchemeResponse object into the model struct
func readAes256PasswordStorageSchemeResponse(ctx context.Context, r *client.Aes256PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aes-256")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EncryptionSettingsDefinitionID = internaltypes.StringTypeOrNil(r.EncryptionSettingsDefinitionID, internaltypes.IsEmptyString(expectedValues.EncryptionSettingsDefinitionID))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a BcryptPasswordStorageSchemeResponse object into the model struct
func readBcryptPasswordStorageSchemeResponse(ctx context.Context, r *client.BcryptPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("bcrypt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.BcryptCostFactor = internaltypes.Int64TypeOrNil(r.BcryptCostFactor)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a BlowfishPasswordStorageSchemeResponse object into the model struct
func readBlowfishPasswordStorageSchemeResponse(ctx context.Context, r *client.BlowfishPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("blowfish")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Sha1PasswordStorageSchemeResponse object into the model struct
func readSha1PasswordStorageSchemeResponse(ctx context.Context, r *client.Sha1PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("sha1")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a AmazonSecretsManagerPasswordStorageSchemeResponse object into the model struct
func readAmazonSecretsManagerPasswordStorageSchemeResponse(ctx context.Context, r *client.AmazonSecretsManagerPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-secrets-manager")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.DefaultField = internaltypes.StringTypeOrNil(r.DefaultField, internaltypes.IsEmptyString(expectedValues.DefaultField))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a AzureKeyVaultPasswordStorageSchemeResponse object into the model struct
func readAzureKeyVaultPasswordStorageSchemeResponse(ctx context.Context, r *client.AzureKeyVaultPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("azure-key-vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a ConjurPasswordStorageSchemeResponse object into the model struct
func readConjurPasswordStorageSchemeResponse(ctx context.Context, r *client.ConjurPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a SaltedSha1PasswordStorageSchemeResponse object into the model struct
func readSaltedSha1PasswordStorageSchemeResponse(ctx context.Context, r *client.SaltedSha1PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha1")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a SaltedSha512PasswordStorageSchemeResponse object into the model struct
func readSaltedSha512PasswordStorageSchemeResponse(ctx context.Context, r *client.SaltedSha512PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("salted-sha512")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SaltLengthBytes = internaltypes.Int64TypeOrNil(r.SaltLengthBytes)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a ScryptPasswordStorageSchemeResponse object into the model struct
func readScryptPasswordStorageSchemeResponse(ctx context.Context, r *client.ScryptPasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("scrypt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScryptCpuMemoryCostFactorExponent = internaltypes.Int64TypeOrNil(r.ScryptCpuMemoryCostFactorExponent)
	state.ScryptBlockSize = internaltypes.Int64TypeOrNil(r.ScryptBlockSize)
	state.ScryptParallelizationParameter = internaltypes.Int64TypeOrNil(r.ScryptParallelizationParameter)
	state.MaxPasswordLength = internaltypes.Int64TypeOrNil(r.MaxPasswordLength)
	state.EncodedPasswordCacheSize = internaltypes.Int64TypeOrNil(r.EncodedPasswordCacheSize)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Read a Md5PasswordStorageSchemeResponse object into the model struct
func readMd5PasswordStorageSchemeResponse(ctx context.Context, r *client.Md5PasswordStorageSchemeResponse, state *passwordStorageSchemeResourceModel, expectedValues *passwordStorageSchemeResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("md5")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePasswordStorageSchemeUnknownValues(state)
}

// Create any update operations necessary to make the state match the plan
func createPasswordStorageSchemeOperations(plan passwordStorageSchemeResourceModel, state passwordStorageSchemeResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddInt64OperationIfNecessary(&ops, plan.ScryptCpuMemoryCostFactorExponent, state.ScryptCpuMemoryCostFactorExponent, "scrypt-cpu-memory-cost-factor-exponent")
	operations.AddInt64OperationIfNecessary(&ops, plan.ScryptBlockSize, state.ScryptBlockSize, "scrypt-block-size")
	operations.AddInt64OperationIfNecessary(&ops, plan.ScryptParallelizationParameter, state.ScryptParallelizationParameter, "scrypt-parallelization-parameter")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurExternalServer, state.ConjurExternalServer, "conjur-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyVaultURI, state.KeyVaultURI, "key-vault-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.AzureAuthenticationMethod, state.AzureAuthenticationMethod, "azure-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsExternalServer, state.AwsExternalServer, "aws-external-server")
	operations.AddInt64OperationIfNecessary(&ops, plan.BcryptCostFactor, state.BcryptCostFactor, "bcrypt-cost-factor")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionSettingsDefinitionID, state.EncryptionSettingsDefinitionID, "encryption-settings-definition-id")
	operations.AddStringOperationIfNecessary(&ops, plan.DigestAlgorithm, state.DigestAlgorithm, "digest-algorithm")
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultExternalServer, state.VaultExternalServer, "vault-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.DefaultField, state.DefaultField, "default-field")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordEncodingMechanism, state.PasswordEncodingMechanism, "password-encoding-mechanism")
	operations.AddInt64OperationIfNecessary(&ops, plan.NumDigestRounds, state.NumDigestRounds, "num-digest-rounds")
	operations.AddInt64OperationIfNecessary(&ops, plan.MaxPasswordLength, state.MaxPasswordLength, "max-password-length")
	operations.AddInt64OperationIfNecessary(&ops, plan.IterationCount, state.IterationCount, "iteration-count")
	operations.AddInt64OperationIfNecessary(&ops, plan.ParallelismFactor, state.ParallelismFactor, "parallelism-factor")
	operations.AddInt64OperationIfNecessary(&ops, plan.MemoryUsageKb, state.MemoryUsageKb, "memory-usage-kb")
	operations.AddInt64OperationIfNecessary(&ops, plan.SaltLengthBytes, state.SaltLengthBytes, "salt-length-bytes")
	operations.AddInt64OperationIfNecessary(&ops, plan.DerivedKeyLengthBytes, state.DerivedKeyLengthBytes, "derived-key-length-bytes")
	operations.AddInt64OperationIfNecessary(&ops, plan.EncodedPasswordCacheSize, state.EncodedPasswordCacheSize, "encoded-password-cache-size")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a argon2 password-storage-scheme
func (r *passwordStorageSchemeResource) CreateArgon2PasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddArgon2PasswordStorageSchemeRequest([]client.Enumargon2PasswordStorageSchemeSchemaUrn{client.ENUMARGON2PASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEARGON2},
		plan.IterationCount.ValueInt64(),
		plan.ParallelismFactor.ValueInt64(),
		plan.MemoryUsageKb.ValueInt64(),
		plan.SaltLengthBytes.ValueInt64(),
		plan.DerivedKeyLengthBytes.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalArgon2PasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddArgon2PasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readArgon2PasswordStorageSchemeResponse(ctx, addResponse.Argon2PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party-enhanced password-storage-scheme
func (r *passwordStorageSchemeResource) CreateThirdPartyEnhancedPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddThirdPartyEnhancedPasswordStorageSchemeRequest([]client.EnumthirdPartyEnhancedPasswordStorageSchemeSchemaUrn{client.ENUMTHIRDPARTYENHANCEDPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMETHIRD_PARTY_ENHANCED},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalThirdPartyEnhancedPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddThirdPartyEnhancedPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readThirdPartyEnhancedPasswordStorageSchemeResponse(ctx, addResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a pbkdf2 password-storage-scheme
func (r *passwordStorageSchemeResource) CreatePbkdf2PasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddPbkdf2PasswordStorageSchemeRequest([]client.Enumpbkdf2PasswordStorageSchemeSchemaUrn{client.ENUMPBKDF2PASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEPBKDF2},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalPbkdf2PasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddPbkdf2PasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readPbkdf2PasswordStorageSchemeResponse(ctx, addResponse.Pbkdf2PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a argon2d password-storage-scheme
func (r *passwordStorageSchemeResource) CreateArgon2dPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddArgon2dPasswordStorageSchemeRequest([]client.Enumargon2dPasswordStorageSchemeSchemaUrn{client.ENUMARGON2DPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEARGON2D},
		plan.IterationCount.ValueInt64(),
		plan.ParallelismFactor.ValueInt64(),
		plan.MemoryUsageKb.ValueInt64(),
		plan.SaltLengthBytes.ValueInt64(),
		plan.DerivedKeyLengthBytes.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalArgon2dPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddArgon2dPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readArgon2dPasswordStorageSchemeResponse(ctx, addResponse.Argon2dPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a crypt password-storage-scheme
func (r *passwordStorageSchemeResource) CreateCryptPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddCryptPasswordStorageSchemeRequest([]client.EnumcryptPasswordStorageSchemeSchemaUrn{client.ENUMCRYPTPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMECRYPT},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalCryptPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddCryptPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readCryptPasswordStorageSchemeResponse(ctx, addResponse.CryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a argon2i password-storage-scheme
func (r *passwordStorageSchemeResource) CreateArgon2iPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddArgon2iPasswordStorageSchemeRequest([]client.Enumargon2iPasswordStorageSchemeSchemaUrn{client.ENUMARGON2IPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEARGON2I},
		plan.IterationCount.ValueInt64(),
		plan.ParallelismFactor.ValueInt64(),
		plan.MemoryUsageKb.ValueInt64(),
		plan.SaltLengthBytes.ValueInt64(),
		plan.DerivedKeyLengthBytes.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalArgon2iPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddArgon2iPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readArgon2iPasswordStorageSchemeResponse(ctx, addResponse.Argon2iPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a aes-256 password-storage-scheme
func (r *passwordStorageSchemeResource) CreateAes256PasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddAes256PasswordStorageSchemeRequest([]client.Enumaes256PasswordStorageSchemeSchemaUrn{client.ENUMAES256PASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEAES_256},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalAes256PasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddAes256PasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readAes256PasswordStorageSchemeResponse(ctx, addResponse.Aes256PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a bcrypt password-storage-scheme
func (r *passwordStorageSchemeResource) CreateBcryptPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddBcryptPasswordStorageSchemeRequest([]client.EnumbcryptPasswordStorageSchemeSchemaUrn{client.ENUMBCRYPTPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEBCRYPT},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalBcryptPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddBcryptPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readBcryptPasswordStorageSchemeResponse(ctx, addResponse.BcryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a argon2id password-storage-scheme
func (r *passwordStorageSchemeResource) CreateArgon2idPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddArgon2idPasswordStorageSchemeRequest([]client.Enumargon2idPasswordStorageSchemeSchemaUrn{client.ENUMARGON2IDPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEARGON2ID},
		plan.IterationCount.ValueInt64(),
		plan.ParallelismFactor.ValueInt64(),
		plan.MemoryUsageKb.ValueInt64(),
		plan.SaltLengthBytes.ValueInt64(),
		plan.DerivedKeyLengthBytes.ValueInt64(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalArgon2idPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddArgon2idPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readArgon2idPasswordStorageSchemeResponse(ctx, addResponse.Argon2idPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a amazon-secrets-manager password-storage-scheme
func (r *passwordStorageSchemeResource) CreateAmazonSecretsManagerPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddAmazonSecretsManagerPasswordStorageSchemeRequest([]client.EnumamazonSecretsManagerPasswordStorageSchemeSchemaUrn{client.ENUMAMAZONSECRETSMANAGERPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEAMAZON_SECRETS_MANAGER},
		plan.AwsExternalServer.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalAmazonSecretsManagerPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddAmazonSecretsManagerPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readAmazonSecretsManagerPasswordStorageSchemeResponse(ctx, addResponse.AmazonSecretsManagerPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a azure-key-vault password-storage-scheme
func (r *passwordStorageSchemeResource) CreateAzureKeyVaultPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddAzureKeyVaultPasswordStorageSchemeRequest([]client.EnumazureKeyVaultPasswordStorageSchemeSchemaUrn{client.ENUMAZUREKEYVAULTPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEAZURE_KEY_VAULT},
		plan.KeyVaultURI.ValueString(),
		plan.AzureAuthenticationMethod.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalAzureKeyVaultPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddAzureKeyVaultPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readAzureKeyVaultPasswordStorageSchemeResponse(ctx, addResponse.AzureKeyVaultPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a conjur password-storage-scheme
func (r *passwordStorageSchemeResource) CreateConjurPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddConjurPasswordStorageSchemeRequest([]client.EnumconjurPasswordStorageSchemeSchemaUrn{client.ENUMCONJURPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMECONJUR},
		plan.ConjurExternalServer.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalConjurPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddConjurPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readConjurPasswordStorageSchemeResponse(ctx, addResponse.ConjurPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a scrypt password-storage-scheme
func (r *passwordStorageSchemeResource) CreateScryptPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddScryptPasswordStorageSchemeRequest([]client.EnumscryptPasswordStorageSchemeSchemaUrn{client.ENUMSCRYPTPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMESCRYPT},
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalScryptPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddScryptPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readScryptPasswordStorageSchemeResponse(ctx, addResponse.ScryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a vault password-storage-scheme
func (r *passwordStorageSchemeResource) CreateVaultPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddVaultPasswordStorageSchemeRequest([]client.EnumvaultPasswordStorageSchemeSchemaUrn{client.ENUMVAULTPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMEVAULT},
		plan.VaultExternalServer.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalVaultPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddVaultPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readVaultPasswordStorageSchemeResponse(ctx, addResponse.VaultPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party password-storage-scheme
func (r *passwordStorageSchemeResource) CreateThirdPartyPasswordStorageScheme(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passwordStorageSchemeResourceModel) (*passwordStorageSchemeResourceModel, error) {
	addRequest := client.NewAddThirdPartyPasswordStorageSchemeRequest([]client.EnumthirdPartyPasswordStorageSchemeSchemaUrn{client.ENUMTHIRDPARTYPASSWORDSTORAGESCHEMESCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSWORD_STORAGE_SCHEMETHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.Name.ValueString())
	err := addOptionalThirdPartyPasswordStorageSchemeFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Password Storage Scheme", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPasswordStorageSchemeRequest(
		client.AddThirdPartyPasswordStorageSchemeRequestAsAddPasswordStorageSchemeRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.AddPasswordStorageSchemeExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Password Storage Scheme", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passwordStorageSchemeResourceModel
	readThirdPartyPasswordStorageSchemeResponse(ctx, addResponse.ThirdPartyPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passwordStorageSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordStorageSchemeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *passwordStorageSchemeResourceModel
	var err error
	if plan.Type.ValueString() == "argon2" {
		state, err = r.CreateArgon2PasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party-enhanced" {
		state, err = r.CreateThirdPartyEnhancedPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "pbkdf2" {
		state, err = r.CreatePbkdf2PasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "argon2d" {
		state, err = r.CreateArgon2dPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "crypt" {
		state, err = r.CreateCryptPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "argon2i" {
		state, err = r.CreateArgon2iPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "aes-256" {
		state, err = r.CreateAes256PasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "bcrypt" {
		state, err = r.CreateBcryptPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "argon2id" {
		state, err = r.CreateArgon2idPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "amazon-secrets-manager" {
		state, err = r.CreateAmazonSecretsManagerPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "azure-key-vault" {
		state, err = r.CreateAzureKeyVaultPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "conjur" {
		state, err = r.CreateConjurPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "scrypt" {
		state, err = r.CreateScryptPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "vault" {
		state, err = r.CreateVaultPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPasswordStorageScheme(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, *state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Create a new resource
// For edit only resources like this, create doesn't actually "create" anything - it "adopts" the existing
// config object into management by terraform. This method reads the existing config object
// and makes any changes needed to make it match the plan - similar to the Update method.
func (r *defaultPasswordStorageSchemeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passwordStorageSchemeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.GetPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Storage Scheme", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passwordStorageSchemeResourceModel
	if readResponse.SaltedSha256PasswordStorageSchemeResponse != nil {
		readSaltedSha256PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha256PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2dPasswordStorageSchemeResponse != nil {
		readArgon2dPasswordStorageSchemeResponse(ctx, readResponse.Argon2dPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CryptPasswordStorageSchemeResponse != nil {
		readCryptPasswordStorageSchemeResponse(ctx, readResponse.CryptPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2iPasswordStorageSchemeResponse != nil {
		readArgon2iPasswordStorageSchemeResponse(ctx, readResponse.Argon2iPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Base64PasswordStorageSchemeResponse != nil {
		readBase64PasswordStorageSchemeResponse(ctx, readResponse.Base64PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedMd5PasswordStorageSchemeResponse != nil {
		readSaltedMd5PasswordStorageSchemeResponse(ctx, readResponse.SaltedMd5PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AesPasswordStorageSchemeResponse != nil {
		readAesPasswordStorageSchemeResponse(ctx, readResponse.AesPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2idPasswordStorageSchemeResponse != nil {
		readArgon2idPasswordStorageSchemeResponse(ctx, readResponse.Argon2idPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultPasswordStorageSchemeResponse != nil {
		readVaultPasswordStorageSchemeResponse(ctx, readResponse.VaultPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordStorageSchemeResponse != nil {
		readThirdPartyPasswordStorageSchemeResponse(ctx, readResponse.ThirdPartyPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2PasswordStorageSchemeResponse != nil {
		readArgon2PasswordStorageSchemeResponse(ctx, readResponse.Argon2PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse != nil {
		readThirdPartyEnhancedPasswordStorageSchemeResponse(ctx, readResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Pbkdf2PasswordStorageSchemeResponse != nil {
		readPbkdf2PasswordStorageSchemeResponse(ctx, readResponse.Pbkdf2PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Rc4PasswordStorageSchemeResponse != nil {
		readRc4PasswordStorageSchemeResponse(ctx, readResponse.Rc4PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha384PasswordStorageSchemeResponse != nil {
		readSaltedSha384PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha384PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TripleDesPasswordStorageSchemeResponse != nil {
		readTripleDesPasswordStorageSchemeResponse(ctx, readResponse.TripleDesPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ClearPasswordStorageSchemeResponse != nil {
		readClearPasswordStorageSchemeResponse(ctx, readResponse.ClearPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Aes256PasswordStorageSchemeResponse != nil {
		readAes256PasswordStorageSchemeResponse(ctx, readResponse.Aes256PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BcryptPasswordStorageSchemeResponse != nil {
		readBcryptPasswordStorageSchemeResponse(ctx, readResponse.BcryptPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BlowfishPasswordStorageSchemeResponse != nil {
		readBlowfishPasswordStorageSchemeResponse(ctx, readResponse.BlowfishPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Sha1PasswordStorageSchemeResponse != nil {
		readSha1PasswordStorageSchemeResponse(ctx, readResponse.Sha1PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerPasswordStorageSchemeResponse != nil {
		readAmazonSecretsManagerPasswordStorageSchemeResponse(ctx, readResponse.AmazonSecretsManagerPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultPasswordStorageSchemeResponse != nil {
		readAzureKeyVaultPasswordStorageSchemeResponse(ctx, readResponse.AzureKeyVaultPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurPasswordStorageSchemeResponse != nil {
		readConjurPasswordStorageSchemeResponse(ctx, readResponse.ConjurPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha1PasswordStorageSchemeResponse != nil {
		readSaltedSha1PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha1PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha512PasswordStorageSchemeResponse != nil {
		readSaltedSha512PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha512PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ScryptPasswordStorageSchemeResponse != nil {
		readScryptPasswordStorageSchemeResponse(ctx, readResponse.ScryptPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Md5PasswordStorageSchemeResponse != nil {
		readMd5PasswordStorageSchemeResponse(ctx, readResponse.Md5PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PasswordStorageSchemeAPI.UpdatePasswordStorageScheme(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPasswordStorageSchemeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PasswordStorageSchemeAPI.UpdatePasswordStorageSchemeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Storage Scheme", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SaltedSha256PasswordStorageSchemeResponse != nil {
			readSaltedSha256PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha256PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2dPasswordStorageSchemeResponse != nil {
			readArgon2dPasswordStorageSchemeResponse(ctx, updateResponse.Argon2dPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CryptPasswordStorageSchemeResponse != nil {
			readCryptPasswordStorageSchemeResponse(ctx, updateResponse.CryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2iPasswordStorageSchemeResponse != nil {
			readArgon2iPasswordStorageSchemeResponse(ctx, updateResponse.Argon2iPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Base64PasswordStorageSchemeResponse != nil {
			readBase64PasswordStorageSchemeResponse(ctx, updateResponse.Base64PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedMd5PasswordStorageSchemeResponse != nil {
			readSaltedMd5PasswordStorageSchemeResponse(ctx, updateResponse.SaltedMd5PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AesPasswordStorageSchemeResponse != nil {
			readAesPasswordStorageSchemeResponse(ctx, updateResponse.AesPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2idPasswordStorageSchemeResponse != nil {
			readArgon2idPasswordStorageSchemeResponse(ctx, updateResponse.Argon2idPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VaultPasswordStorageSchemeResponse != nil {
			readVaultPasswordStorageSchemeResponse(ctx, updateResponse.VaultPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPasswordStorageSchemeResponse != nil {
			readThirdPartyPasswordStorageSchemeResponse(ctx, updateResponse.ThirdPartyPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2PasswordStorageSchemeResponse != nil {
			readArgon2PasswordStorageSchemeResponse(ctx, updateResponse.Argon2PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse != nil {
			readThirdPartyEnhancedPasswordStorageSchemeResponse(ctx, updateResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Pbkdf2PasswordStorageSchemeResponse != nil {
			readPbkdf2PasswordStorageSchemeResponse(ctx, updateResponse.Pbkdf2PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Rc4PasswordStorageSchemeResponse != nil {
			readRc4PasswordStorageSchemeResponse(ctx, updateResponse.Rc4PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedSha384PasswordStorageSchemeResponse != nil {
			readSaltedSha384PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha384PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TripleDesPasswordStorageSchemeResponse != nil {
			readTripleDesPasswordStorageSchemeResponse(ctx, updateResponse.TripleDesPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ClearPasswordStorageSchemeResponse != nil {
			readClearPasswordStorageSchemeResponse(ctx, updateResponse.ClearPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Aes256PasswordStorageSchemeResponse != nil {
			readAes256PasswordStorageSchemeResponse(ctx, updateResponse.Aes256PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BcryptPasswordStorageSchemeResponse != nil {
			readBcryptPasswordStorageSchemeResponse(ctx, updateResponse.BcryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BlowfishPasswordStorageSchemeResponse != nil {
			readBlowfishPasswordStorageSchemeResponse(ctx, updateResponse.BlowfishPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Sha1PasswordStorageSchemeResponse != nil {
			readSha1PasswordStorageSchemeResponse(ctx, updateResponse.Sha1PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AmazonSecretsManagerPasswordStorageSchemeResponse != nil {
			readAmazonSecretsManagerPasswordStorageSchemeResponse(ctx, updateResponse.AmazonSecretsManagerPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AzureKeyVaultPasswordStorageSchemeResponse != nil {
			readAzureKeyVaultPasswordStorageSchemeResponse(ctx, updateResponse.AzureKeyVaultPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConjurPasswordStorageSchemeResponse != nil {
			readConjurPasswordStorageSchemeResponse(ctx, updateResponse.ConjurPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedSha1PasswordStorageSchemeResponse != nil {
			readSaltedSha1PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha1PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedSha512PasswordStorageSchemeResponse != nil {
			readSaltedSha512PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha512PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ScryptPasswordStorageSchemeResponse != nil {
			readScryptPasswordStorageSchemeResponse(ctx, updateResponse.ScryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Md5PasswordStorageSchemeResponse != nil {
			readMd5PasswordStorageSchemeResponse(ctx, updateResponse.Md5PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
	}

	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *passwordStorageSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordStorageScheme(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPasswordStorageSchemeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPasswordStorageScheme(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPasswordStorageScheme(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state passwordStorageSchemeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PasswordStorageSchemeAPI.GetPasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp != nil && httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Password Storage Scheme", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Storage Scheme", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.SaltedSha256PasswordStorageSchemeResponse != nil {
		readSaltedSha256PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha256PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2dPasswordStorageSchemeResponse != nil {
		readArgon2dPasswordStorageSchemeResponse(ctx, readResponse.Argon2dPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CryptPasswordStorageSchemeResponse != nil {
		readCryptPasswordStorageSchemeResponse(ctx, readResponse.CryptPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2iPasswordStorageSchemeResponse != nil {
		readArgon2iPasswordStorageSchemeResponse(ctx, readResponse.Argon2iPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Base64PasswordStorageSchemeResponse != nil {
		readBase64PasswordStorageSchemeResponse(ctx, readResponse.Base64PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedMd5PasswordStorageSchemeResponse != nil {
		readSaltedMd5PasswordStorageSchemeResponse(ctx, readResponse.SaltedMd5PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AesPasswordStorageSchemeResponse != nil {
		readAesPasswordStorageSchemeResponse(ctx, readResponse.AesPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2idPasswordStorageSchemeResponse != nil {
		readArgon2idPasswordStorageSchemeResponse(ctx, readResponse.Argon2idPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultPasswordStorageSchemeResponse != nil {
		readVaultPasswordStorageSchemeResponse(ctx, readResponse.VaultPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordStorageSchemeResponse != nil {
		readThirdPartyPasswordStorageSchemeResponse(ctx, readResponse.ThirdPartyPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Argon2PasswordStorageSchemeResponse != nil {
		readArgon2PasswordStorageSchemeResponse(ctx, readResponse.Argon2PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse != nil {
		readThirdPartyEnhancedPasswordStorageSchemeResponse(ctx, readResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Pbkdf2PasswordStorageSchemeResponse != nil {
		readPbkdf2PasswordStorageSchemeResponse(ctx, readResponse.Pbkdf2PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Rc4PasswordStorageSchemeResponse != nil {
		readRc4PasswordStorageSchemeResponse(ctx, readResponse.Rc4PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha384PasswordStorageSchemeResponse != nil {
		readSaltedSha384PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha384PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.TripleDesPasswordStorageSchemeResponse != nil {
		readTripleDesPasswordStorageSchemeResponse(ctx, readResponse.TripleDesPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ClearPasswordStorageSchemeResponse != nil {
		readClearPasswordStorageSchemeResponse(ctx, readResponse.ClearPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Aes256PasswordStorageSchemeResponse != nil {
		readAes256PasswordStorageSchemeResponse(ctx, readResponse.Aes256PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BcryptPasswordStorageSchemeResponse != nil {
		readBcryptPasswordStorageSchemeResponse(ctx, readResponse.BcryptPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.BlowfishPasswordStorageSchemeResponse != nil {
		readBlowfishPasswordStorageSchemeResponse(ctx, readResponse.BlowfishPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Sha1PasswordStorageSchemeResponse != nil {
		readSha1PasswordStorageSchemeResponse(ctx, readResponse.Sha1PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerPasswordStorageSchemeResponse != nil {
		readAmazonSecretsManagerPasswordStorageSchemeResponse(ctx, readResponse.AmazonSecretsManagerPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultPasswordStorageSchemeResponse != nil {
		readAzureKeyVaultPasswordStorageSchemeResponse(ctx, readResponse.AzureKeyVaultPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurPasswordStorageSchemeResponse != nil {
		readConjurPasswordStorageSchemeResponse(ctx, readResponse.ConjurPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha1PasswordStorageSchemeResponse != nil {
		readSaltedSha1PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha1PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.SaltedSha512PasswordStorageSchemeResponse != nil {
		readSaltedSha512PasswordStorageSchemeResponse(ctx, readResponse.SaltedSha512PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ScryptPasswordStorageSchemeResponse != nil {
		readScryptPasswordStorageSchemeResponse(ctx, readResponse.ScryptPasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Md5PasswordStorageSchemeResponse != nil {
		readMd5PasswordStorageSchemeResponse(ctx, readResponse.Md5PasswordStorageSchemeResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *passwordStorageSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordStorageScheme(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPasswordStorageSchemeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePasswordStorageScheme(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePasswordStorageScheme(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passwordStorageSchemeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passwordStorageSchemeResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PasswordStorageSchemeAPI.UpdatePasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPasswordStorageSchemeOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PasswordStorageSchemeAPI.UpdatePasswordStorageSchemeExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Password Storage Scheme", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.SaltedSha256PasswordStorageSchemeResponse != nil {
			readSaltedSha256PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha256PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2dPasswordStorageSchemeResponse != nil {
			readArgon2dPasswordStorageSchemeResponse(ctx, updateResponse.Argon2dPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CryptPasswordStorageSchemeResponse != nil {
			readCryptPasswordStorageSchemeResponse(ctx, updateResponse.CryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2iPasswordStorageSchemeResponse != nil {
			readArgon2iPasswordStorageSchemeResponse(ctx, updateResponse.Argon2iPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Base64PasswordStorageSchemeResponse != nil {
			readBase64PasswordStorageSchemeResponse(ctx, updateResponse.Base64PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedMd5PasswordStorageSchemeResponse != nil {
			readSaltedMd5PasswordStorageSchemeResponse(ctx, updateResponse.SaltedMd5PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AesPasswordStorageSchemeResponse != nil {
			readAesPasswordStorageSchemeResponse(ctx, updateResponse.AesPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2idPasswordStorageSchemeResponse != nil {
			readArgon2idPasswordStorageSchemeResponse(ctx, updateResponse.Argon2idPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VaultPasswordStorageSchemeResponse != nil {
			readVaultPasswordStorageSchemeResponse(ctx, updateResponse.VaultPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPasswordStorageSchemeResponse != nil {
			readThirdPartyPasswordStorageSchemeResponse(ctx, updateResponse.ThirdPartyPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Argon2PasswordStorageSchemeResponse != nil {
			readArgon2PasswordStorageSchemeResponse(ctx, updateResponse.Argon2PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse != nil {
			readThirdPartyEnhancedPasswordStorageSchemeResponse(ctx, updateResponse.ThirdPartyEnhancedPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Pbkdf2PasswordStorageSchemeResponse != nil {
			readPbkdf2PasswordStorageSchemeResponse(ctx, updateResponse.Pbkdf2PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Rc4PasswordStorageSchemeResponse != nil {
			readRc4PasswordStorageSchemeResponse(ctx, updateResponse.Rc4PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedSha384PasswordStorageSchemeResponse != nil {
			readSaltedSha384PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha384PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.TripleDesPasswordStorageSchemeResponse != nil {
			readTripleDesPasswordStorageSchemeResponse(ctx, updateResponse.TripleDesPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ClearPasswordStorageSchemeResponse != nil {
			readClearPasswordStorageSchemeResponse(ctx, updateResponse.ClearPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Aes256PasswordStorageSchemeResponse != nil {
			readAes256PasswordStorageSchemeResponse(ctx, updateResponse.Aes256PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BcryptPasswordStorageSchemeResponse != nil {
			readBcryptPasswordStorageSchemeResponse(ctx, updateResponse.BcryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.BlowfishPasswordStorageSchemeResponse != nil {
			readBlowfishPasswordStorageSchemeResponse(ctx, updateResponse.BlowfishPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Sha1PasswordStorageSchemeResponse != nil {
			readSha1PasswordStorageSchemeResponse(ctx, updateResponse.Sha1PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AmazonSecretsManagerPasswordStorageSchemeResponse != nil {
			readAmazonSecretsManagerPasswordStorageSchemeResponse(ctx, updateResponse.AmazonSecretsManagerPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AzureKeyVaultPasswordStorageSchemeResponse != nil {
			readAzureKeyVaultPasswordStorageSchemeResponse(ctx, updateResponse.AzureKeyVaultPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConjurPasswordStorageSchemeResponse != nil {
			readConjurPasswordStorageSchemeResponse(ctx, updateResponse.ConjurPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedSha1PasswordStorageSchemeResponse != nil {
			readSaltedSha1PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha1PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.SaltedSha512PasswordStorageSchemeResponse != nil {
			readSaltedSha512PasswordStorageSchemeResponse(ctx, updateResponse.SaltedSha512PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ScryptPasswordStorageSchemeResponse != nil {
			readScryptPasswordStorageSchemeResponse(ctx, updateResponse.ScryptPasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Md5PasswordStorageSchemeResponse != nil {
			readMd5PasswordStorageSchemeResponse(ctx, updateResponse.Md5PasswordStorageSchemeResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultPasswordStorageSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passwordStorageSchemeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passwordStorageSchemeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PasswordStorageSchemeAPI.DeletePasswordStorageSchemeExecute(r.apiClient.PasswordStorageSchemeAPI.DeletePasswordStorageScheme(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && (httpResp == nil || httpResp.StatusCode != 404) {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Password Storage Scheme", err, httpResp)
		return
	}
}

func (r *passwordStorageSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordStorageScheme(ctx, req, resp)
}

func (r *defaultPasswordStorageSchemeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPasswordStorageScheme(ctx, req, resp)
}

func importPasswordStorageScheme(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
