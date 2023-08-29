package passphraseprovider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/configvalidators"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/operations"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/version"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &passphraseProviderResource{}
	_ resource.ResourceWithConfigure   = &passphraseProviderResource{}
	_ resource.ResourceWithImportState = &passphraseProviderResource{}
	_ resource.Resource                = &defaultPassphraseProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultPassphraseProviderResource{}
	_ resource.ResourceWithImportState = &defaultPassphraseProviderResource{}
)

// Create a Passphrase Provider resource
func NewPassphraseProviderResource() resource.Resource {
	return &passphraseProviderResource{}
}

func NewDefaultPassphraseProviderResource() resource.Resource {
	return &defaultPassphraseProviderResource{}
}

// passphraseProviderResource is the resource implementation.
type passphraseProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultPassphraseProviderResource is the resource implementation.
type defaultPassphraseProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *passphraseProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_passphrase_provider"
}

func (r *defaultPassphraseProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_passphrase_provider"
}

// Configure adds the provider configured client to the resource.
func (r *passphraseProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultPassphraseProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passphraseProviderResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	LastUpdated               types.String `tfsdk:"last_updated"`
	Notifications             types.Set    `tfsdk:"notifications"`
	RequiredActions           types.Set    `tfsdk:"required_actions"`
	Type                      types.String `tfsdk:"type"`
	ExtensionClass            types.String `tfsdk:"extension_class"`
	ExtensionArgument         types.Set    `tfsdk:"extension_argument"`
	VaultExternalServer       types.String `tfsdk:"vault_external_server"`
	VaultSecretPath           types.String `tfsdk:"vault_secret_path"`
	VaultSecretFieldName      types.String `tfsdk:"vault_secret_field_name"`
	ConjurExternalServer      types.String `tfsdk:"conjur_external_server"`
	ConjurSecretRelativePath  types.String `tfsdk:"conjur_secret_relative_path"`
	PasswordFile              types.String `tfsdk:"password_file"`
	KeyVaultURI               types.String `tfsdk:"key_vault_uri"`
	AzureAuthenticationMethod types.String `tfsdk:"azure_authentication_method"`
	HttpProxyExternalServer   types.String `tfsdk:"http_proxy_external_server"`
	SecretName                types.String `tfsdk:"secret_name"`
	ObscuredValue             types.String `tfsdk:"obscured_value"`
	AwsExternalServer         types.String `tfsdk:"aws_external_server"`
	SecretID                  types.String `tfsdk:"secret_id"`
	SecretFieldName           types.String `tfsdk:"secret_field_name"`
	SecretVersionID           types.String `tfsdk:"secret_version_id"`
	SecretVersionStage        types.String `tfsdk:"secret_version_stage"`
	MaxCacheDuration          types.String `tfsdk:"max_cache_duration"`
	EnvironmentVariable       types.String `tfsdk:"environment_variable"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *passphraseProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passphraseProviderSchema(ctx, req, resp, false)
}

func (r *defaultPassphraseProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	passphraseProviderSchema(ctx, req, resp, true)
}

func passphraseProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Passphrase Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Passphrase Provider resource. Options are ['environment-variable', 'amazon-secrets-manager', 'obscured-value', 'azure-key-vault', 'file-based', 'conjur', 'vault', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"environment-variable", "amazon-secrets-manager", "obscured-value", "azure-key-vault", "file-based", "conjur", "vault", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Passphrase Provider.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Passphrase Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"vault_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Vault instance containing the passphrase.",
				Optional:    true,
			},
			"vault_secret_path": schema.StringAttribute{
				Description: "The path to the desired secret in the Vault service. This will be appended to the value of the base-url property for the associated Vault external server.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vault_secret_field_name": schema.StringAttribute{
				Description: "The name of the field in the Vault secret record that contains the passphrase to use to generate the encryption key.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"conjur_external_server": schema.StringAttribute{
				Description: "An external server definition with information needed to connect and authenticate to the Conjur instance containing the passphrase.",
				Optional:    true,
			},
			"conjur_secret_relative_path": schema.StringAttribute{
				Description: "The portion of the path that follows the account name in the URI needed to obtain the desired secret. Any special characters in the path must be URL-encoded.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"password_file": schema.StringAttribute{
				Description: "The path to the file containing the passphrase.",
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
				Description: "Supported in PingDirectory product version 9.2.0.0+. A reference to an HTTP proxy server that should be used for requests sent to the Azure service.",
				Optional:    true,
			},
			"secret_name": schema.StringAttribute{
				Description: "The name of the secret to retrieve.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"obscured_value": schema.StringAttribute{
				Description: "The value to be stored in an obscured form.",
				Optional:    true,
				Sensitive:   true,
			},
			"aws_external_server": schema.StringAttribute{
				Description: "The external server with information to use when interacting with the AWS Secrets Manager.",
				Optional:    true,
			},
			"secret_id": schema.StringAttribute{
				Description: "The Amazon Resource Name (ARN) or the user-friendly name of the secret to be retrieved.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_field_name": schema.StringAttribute{
				Description: "The name of the JSON field whose value is the passphrase that will be retrieved.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_version_id": schema.StringAttribute{
				Description: "The unique identifier for the version of the secret to be retrieved.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"secret_version_stage": schema.StringAttribute{
				Description: "The staging label for the version of the secret to be retrieved.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"max_cache_duration": schema.StringAttribute{
				Description:         "When the `type` attribute is set to  one of [`amazon-secrets-manager`, `vault`]: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Vault. A value of zero seconds indicates that the provider should always attempt to read the passphrase from Vault. When the `type` attribute is set to `azure-key-vault`: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Azure Key Vault. A value of zero seconds indicates that the provider should always attempt to read the passphrase from the Azure service. When the `type` attribute is set to `file-based`: The maximum length of time that the passphrase provider may cache the passphrase that has been read from the target file. A value of zero seconds indicates that the provider should always attempt to read the passphrase from the file. When the `type` attribute is set to `conjur`: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Conjur. A value of zero seconds indicates that the provider should always attempt to read the passphrase from Conjur.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - One of [`amazon-secrets-manager`, `vault`]: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Vault. A value of zero seconds indicates that the provider should always attempt to read the passphrase from Vault.\n  - `azure-key-vault`: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Azure Key Vault. A value of zero seconds indicates that the provider should always attempt to read the passphrase from the Azure service.\n  - `file-based`: The maximum length of time that the passphrase provider may cache the passphrase that has been read from the target file. A value of zero seconds indicates that the provider should always attempt to read the passphrase from the file.\n  - `conjur`: The maximum length of time that the passphrase provider may cache the passphrase that has been read from Conjur. A value of zero seconds indicates that the provider should always attempt to read the passphrase from Conjur.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_variable": schema.StringAttribute{
				Description: "The name of the environment variable that is expected to hold the passphrase.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Passphrase Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Passphrase Provider is enabled for use in the server.",
				Required:    true,
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
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Validate that any restrictions are met in the plan
func (r *passphraseProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPassphraseProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPassphraseProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanPassphraseProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanPassphraseProvider(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9200)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model passphraseProviderResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.HttpProxyExternalServer) {
		resp.Diagnostics.AddError("Attribute 'http_proxy_external_server' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsPassphraseProvider() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("environment_variable"),
			path.MatchRoot("type"),
			[]string{"environment-variable"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("aws_external_server"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("secret_id"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("secret_field_name"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("secret_version_id"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("secret_version_stage"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("max_cache_duration"),
			path.MatchRoot("type"),
			[]string{"amazon-secrets-manager", "azure-key-vault", "file-based", "conjur", "vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("obscured_value"),
			path.MatchRoot("type"),
			[]string{"obscured-value"},
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
			path.MatchRoot("secret_name"),
			path.MatchRoot("type"),
			[]string{"azure-key-vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("password_file"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conjur_external_server"),
			path.MatchRoot("type"),
			[]string{"conjur"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("conjur_secret_relative_path"),
			path.MatchRoot("type"),
			[]string{"conjur"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_external_server"),
			path.MatchRoot("type"),
			[]string{"vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_secret_path"),
			path.MatchRoot("type"),
			[]string{"vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("vault_secret_field_name"),
			path.MatchRoot("type"),
			[]string{"vault"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"environment-variable",
			[]path.Expression{path.MatchRoot("environment_variable")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"amazon-secrets-manager",
			[]path.Expression{path.MatchRoot("aws_external_server"), path.MatchRoot("secret_id"), path.MatchRoot("secret_field_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"obscured-value",
			[]path.Expression{path.MatchRoot("obscured_value")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"azure-key-vault",
			[]path.Expression{path.MatchRoot("key_vault_uri"), path.MatchRoot("azure_authentication_method"), path.MatchRoot("secret_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"file-based",
			[]path.Expression{path.MatchRoot("password_file")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"conjur",
			[]path.Expression{path.MatchRoot("conjur_external_server"), path.MatchRoot("conjur_secret_relative_path")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"vault",
			[]path.Expression{path.MatchRoot("vault_external_server"), path.MatchRoot("vault_secret_path"), path.MatchRoot("vault_secret_field_name")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r passphraseProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPassphraseProvider()
}

// Add config validators
func (r defaultPassphraseProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsPassphraseProvider()
}

// Add optional fields to create request for environment-variable passphrase-provider
func addOptionalEnvironmentVariablePassphraseProviderFields(ctx context.Context, addRequest *client.AddEnvironmentVariablePassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for amazon-secrets-manager passphrase-provider
func addOptionalAmazonSecretsManagerPassphraseProviderFields(ctx context.Context, addRequest *client.AddAmazonSecretsManagerPassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecretVersionID) {
		addRequest.SecretVersionID = plan.SecretVersionID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SecretVersionStage) {
		addRequest.SecretVersionStage = plan.SecretVersionStage.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxCacheDuration) {
		addRequest.MaxCacheDuration = plan.MaxCacheDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for obscured-value passphrase-provider
func addOptionalObscuredValuePassphraseProviderFields(ctx context.Context, addRequest *client.AddObscuredValuePassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for azure-key-vault passphrase-provider
func addOptionalAzureKeyVaultPassphraseProviderFields(ctx context.Context, addRequest *client.AddAzureKeyVaultPassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.HttpProxyExternalServer) {
		addRequest.HttpProxyExternalServer = plan.HttpProxyExternalServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxCacheDuration) {
		addRequest.MaxCacheDuration = plan.MaxCacheDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for file-based passphrase-provider
func addOptionalFileBasedPassphraseProviderFields(ctx context.Context, addRequest *client.AddFileBasedPassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxCacheDuration) {
		addRequest.MaxCacheDuration = plan.MaxCacheDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for conjur passphrase-provider
func addOptionalConjurPassphraseProviderFields(ctx context.Context, addRequest *client.AddConjurPassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxCacheDuration) {
		addRequest.MaxCacheDuration = plan.MaxCacheDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for vault passphrase-provider
func addOptionalVaultPassphraseProviderFields(ctx context.Context, addRequest *client.AddVaultPassphraseProviderRequest, plan passphraseProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.MaxCacheDuration) {
		addRequest.MaxCacheDuration = plan.MaxCacheDuration.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party passphrase-provider
func addOptionalThirdPartyPassphraseProviderFields(ctx context.Context, addRequest *client.AddThirdPartyPassphraseProviderRequest, plan passphraseProviderResourceModel) {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populatePassphraseProviderUnknownValues(model *passphraseProviderResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.MaxCacheDuration.IsUnknown() || model.MaxCacheDuration.IsNull() {
		model.MaxCacheDuration = types.StringValue("")
	}
	if model.ObscuredValue.IsUnknown() {
		model.ObscuredValue = types.StringNull()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *passphraseProviderResourceModel) populateAllComputedStringAttributes() {
	if model.SecretVersionStage.IsUnknown() || model.SecretVersionStage.IsNull() {
		model.SecretVersionStage = types.StringValue("")
	}
	if model.VaultSecretFieldName.IsUnknown() || model.VaultSecretFieldName.IsNull() {
		model.VaultSecretFieldName = types.StringValue("")
	}
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.SecretVersionID.IsUnknown() || model.SecretVersionID.IsNull() {
		model.SecretVersionID = types.StringValue("")
	}
	if model.PasswordFile.IsUnknown() || model.PasswordFile.IsNull() {
		model.PasswordFile = types.StringValue("")
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
	if model.ConjurSecretRelativePath.IsUnknown() || model.ConjurSecretRelativePath.IsNull() {
		model.ConjurSecretRelativePath = types.StringValue("")
	}
	if model.AzureAuthenticationMethod.IsUnknown() || model.AzureAuthenticationMethod.IsNull() {
		model.AzureAuthenticationMethod = types.StringValue("")
	}
	if model.SecretFieldName.IsUnknown() || model.SecretFieldName.IsNull() {
		model.SecretFieldName = types.StringValue("")
	}
	if model.SecretName.IsUnknown() || model.SecretName.IsNull() {
		model.SecretName = types.StringValue("")
	}
	if model.ConjurExternalServer.IsUnknown() || model.ConjurExternalServer.IsNull() {
		model.ConjurExternalServer = types.StringValue("")
	}
	if model.SecretID.IsUnknown() || model.SecretID.IsNull() {
		model.SecretID = types.StringValue("")
	}
	if model.VaultSecretPath.IsUnknown() || model.VaultSecretPath.IsNull() {
		model.VaultSecretPath = types.StringValue("")
	}
	if model.EnvironmentVariable.IsUnknown() || model.EnvironmentVariable.IsNull() {
		model.EnvironmentVariable = types.StringValue("")
	}
	if model.KeyVaultURI.IsUnknown() || model.KeyVaultURI.IsNull() {
		model.KeyVaultURI = types.StringValue("")
	}
	if model.VaultExternalServer.IsUnknown() || model.VaultExternalServer.IsNull() {
		model.VaultExternalServer = types.StringValue("")
	}
}

// Read a EnvironmentVariablePassphraseProviderResponse object into the model struct
func readEnvironmentVariablePassphraseProviderResponse(ctx context.Context, r *client.EnvironmentVariablePassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("environment-variable")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.EnvironmentVariable = types.StringValue(r.EnvironmentVariable)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a AmazonSecretsManagerPassphraseProviderResponse object into the model struct
func readAmazonSecretsManagerPassphraseProviderResponse(ctx context.Context, r *client.AmazonSecretsManagerPassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("amazon-secrets-manager")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AwsExternalServer = types.StringValue(r.AwsExternalServer)
	state.SecretID = types.StringValue(r.SecretID)
	state.SecretFieldName = types.StringValue(r.SecretFieldName)
	state.SecretVersionID = internaltypes.StringTypeOrNil(r.SecretVersionID, internaltypes.IsEmptyString(expectedValues.SecretVersionID))
	state.SecretVersionStage = internaltypes.StringTypeOrNil(r.SecretVersionStage, internaltypes.IsEmptyString(expectedValues.SecretVersionStage))
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, true)
	config.CheckMismatchedPDFormattedAttributes("max_cache_duration",
		expectedValues.MaxCacheDuration, state.MaxCacheDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a ObscuredValuePassphraseProviderResponse object into the model struct
func readObscuredValuePassphraseProviderResponse(ctx context.Context, r *client.ObscuredValuePassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("obscured-value")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a AzureKeyVaultPassphraseProviderResponse object into the model struct
func readAzureKeyVaultPassphraseProviderResponse(ctx context.Context, r *client.AzureKeyVaultPassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("azure-key-vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyVaultURI = types.StringValue(r.KeyVaultURI)
	state.AzureAuthenticationMethod = types.StringValue(r.AzureAuthenticationMethod)
	state.HttpProxyExternalServer = internaltypes.StringTypeOrNil(r.HttpProxyExternalServer, internaltypes.IsEmptyString(expectedValues.HttpProxyExternalServer))
	state.SecretName = types.StringValue(r.SecretName)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, true)
	config.CheckMismatchedPDFormattedAttributes("max_cache_duration",
		expectedValues.MaxCacheDuration, state.MaxCacheDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a FileBasedPassphraseProviderResponse object into the model struct
func readFileBasedPassphraseProviderResponse(ctx context.Context, r *client.FileBasedPassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordFile = types.StringValue(r.PasswordFile)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, true)
	config.CheckMismatchedPDFormattedAttributes("max_cache_duration",
		expectedValues.MaxCacheDuration, state.MaxCacheDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a ConjurPassphraseProviderResponse object into the model struct
func readConjurPassphraseProviderResponse(ctx context.Context, r *client.ConjurPassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("conjur")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ConjurExternalServer = types.StringValue(r.ConjurExternalServer)
	state.ConjurSecretRelativePath = types.StringValue(r.ConjurSecretRelativePath)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, true)
	config.CheckMismatchedPDFormattedAttributes("max_cache_duration",
		expectedValues.MaxCacheDuration, state.MaxCacheDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a VaultPassphraseProviderResponse object into the model struct
func readVaultPassphraseProviderResponse(ctx context.Context, r *client.VaultPassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("vault")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.VaultExternalServer = types.StringValue(r.VaultExternalServer)
	state.VaultSecretPath = types.StringValue(r.VaultSecretPath)
	state.VaultSecretFieldName = types.StringValue(r.VaultSecretFieldName)
	state.MaxCacheDuration = internaltypes.StringTypeOrNil(r.MaxCacheDuration, true)
	config.CheckMismatchedPDFormattedAttributes("max_cache_duration",
		expectedValues.MaxCacheDuration, state.MaxCacheDuration, diagnostics)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Read a ThirdPartyPassphraseProviderResponse object into the model struct
func readThirdPartyPassphraseProviderResponse(ctx context.Context, r *client.ThirdPartyPassphraseProviderResponse, state *passphraseProviderResourceModel, expectedValues *passphraseProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populatePassphraseProviderUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *passphraseProviderResourceModel) setStateValuesNotReturnedByAPI(expectedValues *passphraseProviderResourceModel) {
	if !expectedValues.ObscuredValue.IsUnknown() {
		state.ObscuredValue = expectedValues.ObscuredValue
	}
}

// Create any update operations necessary to make the state match the plan
func createPassphraseProviderOperations(plan passphraseProviderResourceModel, state passphraseProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultExternalServer, state.VaultExternalServer, "vault-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretPath, state.VaultSecretPath, "vault-secret-path")
	operations.AddStringOperationIfNecessary(&ops, plan.VaultSecretFieldName, state.VaultSecretFieldName, "vault-secret-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurExternalServer, state.ConjurExternalServer, "conjur-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.ConjurSecretRelativePath, state.ConjurSecretRelativePath, "conjur-secret-relative-path")
	operations.AddStringOperationIfNecessary(&ops, plan.PasswordFile, state.PasswordFile, "password-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyVaultURI, state.KeyVaultURI, "key-vault-uri")
	operations.AddStringOperationIfNecessary(&ops, plan.AzureAuthenticationMethod, state.AzureAuthenticationMethod, "azure-authentication-method")
	operations.AddStringOperationIfNecessary(&ops, plan.HttpProxyExternalServer, state.HttpProxyExternalServer, "http-proxy-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretName, state.SecretName, "secret-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ObscuredValue, state.ObscuredValue, "obscured-value")
	operations.AddStringOperationIfNecessary(&ops, plan.AwsExternalServer, state.AwsExternalServer, "aws-external-server")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretID, state.SecretID, "secret-id")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretFieldName, state.SecretFieldName, "secret-field-name")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretVersionID, state.SecretVersionID, "secret-version-id")
	operations.AddStringOperationIfNecessary(&ops, plan.SecretVersionStage, state.SecretVersionStage, "secret-version-stage")
	operations.AddStringOperationIfNecessary(&ops, plan.MaxCacheDuration, state.MaxCacheDuration, "max-cache-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.EnvironmentVariable, state.EnvironmentVariable, "environment-variable")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a environment-variable passphrase-provider
func (r *passphraseProviderResource) CreateEnvironmentVariablePassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddEnvironmentVariablePassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumenvironmentVariablePassphraseProviderSchemaUrn{client.ENUMENVIRONMENTVARIABLEPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERENVIRONMENT_VARIABLE},
		plan.EnvironmentVariable.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalEnvironmentVariablePassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddEnvironmentVariablePassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readEnvironmentVariablePassphraseProviderResponse(ctx, addResponse.EnvironmentVariablePassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a amazon-secrets-manager passphrase-provider
func (r *passphraseProviderResource) CreateAmazonSecretsManagerPassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddAmazonSecretsManagerPassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumamazonSecretsManagerPassphraseProviderSchemaUrn{client.ENUMAMAZONSECRETSMANAGERPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERAMAZON_SECRETS_MANAGER},
		plan.AwsExternalServer.ValueString(),
		plan.SecretID.ValueString(),
		plan.SecretFieldName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAmazonSecretsManagerPassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddAmazonSecretsManagerPassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readAmazonSecretsManagerPassphraseProviderResponse(ctx, addResponse.AmazonSecretsManagerPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a obscured-value passphrase-provider
func (r *passphraseProviderResource) CreateObscuredValuePassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddObscuredValuePassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumobscuredValuePassphraseProviderSchemaUrn{client.ENUMOBSCUREDVALUEPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDEROBSCURED_VALUE},
		plan.ObscuredValue.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalObscuredValuePassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddObscuredValuePassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readObscuredValuePassphraseProviderResponse(ctx, addResponse.ObscuredValuePassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a azure-key-vault passphrase-provider
func (r *passphraseProviderResource) CreateAzureKeyVaultPassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddAzureKeyVaultPassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumazureKeyVaultPassphraseProviderSchemaUrn{client.ENUMAZUREKEYVAULTPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERAZURE_KEY_VAULT},
		plan.KeyVaultURI.ValueString(),
		plan.AzureAuthenticationMethod.ValueString(),
		plan.SecretName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalAzureKeyVaultPassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddAzureKeyVaultPassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readAzureKeyVaultPassphraseProviderResponse(ctx, addResponse.AzureKeyVaultPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a file-based passphrase-provider
func (r *passphraseProviderResource) CreateFileBasedPassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddFileBasedPassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedPassphraseProviderSchemaUrn{client.ENUMFILEBASEDPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERFILE_BASED},
		plan.PasswordFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedPassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddFileBasedPassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readFileBasedPassphraseProviderResponse(ctx, addResponse.FileBasedPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a conjur passphrase-provider
func (r *passphraseProviderResource) CreateConjurPassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddConjurPassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumconjurPassphraseProviderSchemaUrn{client.ENUMCONJURPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERCONJUR},
		plan.ConjurExternalServer.ValueString(),
		plan.ConjurSecretRelativePath.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalConjurPassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddConjurPassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readConjurPassphraseProviderResponse(ctx, addResponse.ConjurPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a vault passphrase-provider
func (r *passphraseProviderResource) CreateVaultPassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddVaultPassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumvaultPassphraseProviderSchemaUrn{client.ENUMVAULTPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERVAULT},
		plan.VaultExternalServer.ValueString(),
		plan.VaultSecretPath.ValueString(),
		plan.VaultSecretFieldName.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalVaultPassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddVaultPassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readVaultPassphraseProviderResponse(ctx, addResponse.VaultPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party passphrase-provider
func (r *passphraseProviderResource) CreateThirdPartyPassphraseProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan passphraseProviderResourceModel) (*passphraseProviderResourceModel, error) {
	addRequest := client.NewAddThirdPartyPassphraseProviderRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyPassphraseProviderSchemaUrn{client.ENUMTHIRDPARTYPASSPHRASEPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0PASSPHRASE_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyPassphraseProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.PassphraseProviderApi.AddPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddPassphraseProviderRequest(
		client.AddThirdPartyPassphraseProviderRequestAsAddPassphraseProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.PassphraseProviderApi.AddPassphraseProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Passphrase Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state passphraseProviderResourceModel
	readThirdPartyPassphraseProviderResponse(ctx, addResponse.ThirdPartyPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *passphraseProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passphraseProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *passphraseProviderResourceModel
	var err error
	if plan.Type.ValueString() == "environment-variable" {
		state, err = r.CreateEnvironmentVariablePassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "amazon-secrets-manager" {
		state, err = r.CreateAmazonSecretsManagerPassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "obscured-value" {
		state, err = r.CreateObscuredValuePassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "azure-key-vault" {
		state, err = r.CreateAzureKeyVaultPassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "file-based" {
		state, err = r.CreateFileBasedPassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "conjur" {
		state, err = r.CreateConjurPassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "vault" {
		state, err = r.CreateVaultPassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyPassphraseProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

	state.setStateValuesNotReturnedByAPI(&plan)
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
func (r *defaultPassphraseProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan passphraseProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PassphraseProviderApi.GetPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Passphrase Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state passphraseProviderResourceModel
	if readResponse.EnvironmentVariablePassphraseProviderResponse != nil {
		readEnvironmentVariablePassphraseProviderResponse(ctx, readResponse.EnvironmentVariablePassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerPassphraseProviderResponse != nil {
		readAmazonSecretsManagerPassphraseProviderResponse(ctx, readResponse.AmazonSecretsManagerPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ObscuredValuePassphraseProviderResponse != nil {
		readObscuredValuePassphraseProviderResponse(ctx, readResponse.ObscuredValuePassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultPassphraseProviderResponse != nil {
		readAzureKeyVaultPassphraseProviderResponse(ctx, readResponse.AzureKeyVaultPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedPassphraseProviderResponse != nil {
		readFileBasedPassphraseProviderResponse(ctx, readResponse.FileBasedPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurPassphraseProviderResponse != nil {
		readConjurPassphraseProviderResponse(ctx, readResponse.ConjurPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultPassphraseProviderResponse != nil {
		readVaultPassphraseProviderResponse(ctx, readResponse.VaultPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassphraseProviderResponse != nil {
		readThirdPartyPassphraseProviderResponse(ctx, readResponse.ThirdPartyPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.PassphraseProviderApi.UpdatePassphraseProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createPassphraseProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.PassphraseProviderApi.UpdatePassphraseProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Passphrase Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.EnvironmentVariablePassphraseProviderResponse != nil {
			readEnvironmentVariablePassphraseProviderResponse(ctx, updateResponse.EnvironmentVariablePassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AmazonSecretsManagerPassphraseProviderResponse != nil {
			readAmazonSecretsManagerPassphraseProviderResponse(ctx, updateResponse.AmazonSecretsManagerPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ObscuredValuePassphraseProviderResponse != nil {
			readObscuredValuePassphraseProviderResponse(ctx, updateResponse.ObscuredValuePassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AzureKeyVaultPassphraseProviderResponse != nil {
			readAzureKeyVaultPassphraseProviderResponse(ctx, updateResponse.AzureKeyVaultPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedPassphraseProviderResponse != nil {
			readFileBasedPassphraseProviderResponse(ctx, updateResponse.FileBasedPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConjurPassphraseProviderResponse != nil {
			readConjurPassphraseProviderResponse(ctx, updateResponse.ConjurPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VaultPassphraseProviderResponse != nil {
			readVaultPassphraseProviderResponse(ctx, updateResponse.VaultPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPassphraseProviderResponse != nil {
			readThirdPartyPassphraseProviderResponse(ctx, updateResponse.ThirdPartyPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	state.populateAllComputedStringAttributes()
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *passphraseProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPassphraseProvider(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultPassphraseProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readPassphraseProvider(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readPassphraseProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state passphraseProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.PassphraseProviderApi.GetPassphraseProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Passphrase Provider", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Passphrase Provider", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.EnvironmentVariablePassphraseProviderResponse != nil {
		readEnvironmentVariablePassphraseProviderResponse(ctx, readResponse.EnvironmentVariablePassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AmazonSecretsManagerPassphraseProviderResponse != nil {
		readAmazonSecretsManagerPassphraseProviderResponse(ctx, readResponse.AmazonSecretsManagerPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ObscuredValuePassphraseProviderResponse != nil {
		readObscuredValuePassphraseProviderResponse(ctx, readResponse.ObscuredValuePassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.AzureKeyVaultPassphraseProviderResponse != nil {
		readAzureKeyVaultPassphraseProviderResponse(ctx, readResponse.AzureKeyVaultPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.FileBasedPassphraseProviderResponse != nil {
		readFileBasedPassphraseProviderResponse(ctx, readResponse.FileBasedPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ConjurPassphraseProviderResponse != nil {
		readConjurPassphraseProviderResponse(ctx, readResponse.ConjurPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.VaultPassphraseProviderResponse != nil {
		readVaultPassphraseProviderResponse(ctx, readResponse.VaultPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPassphraseProviderResponse != nil {
		readThirdPartyPassphraseProviderResponse(ctx, readResponse.ThirdPartyPassphraseProviderResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *passphraseProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePassphraseProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultPassphraseProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updatePassphraseProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updatePassphraseProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan passphraseProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state passphraseProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.PassphraseProviderApi.UpdatePassphraseProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createPassphraseProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.PassphraseProviderApi.UpdatePassphraseProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Passphrase Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.EnvironmentVariablePassphraseProviderResponse != nil {
			readEnvironmentVariablePassphraseProviderResponse(ctx, updateResponse.EnvironmentVariablePassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AmazonSecretsManagerPassphraseProviderResponse != nil {
			readAmazonSecretsManagerPassphraseProviderResponse(ctx, updateResponse.AmazonSecretsManagerPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ObscuredValuePassphraseProviderResponse != nil {
			readObscuredValuePassphraseProviderResponse(ctx, updateResponse.ObscuredValuePassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.AzureKeyVaultPassphraseProviderResponse != nil {
			readAzureKeyVaultPassphraseProviderResponse(ctx, updateResponse.AzureKeyVaultPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.FileBasedPassphraseProviderResponse != nil {
			readFileBasedPassphraseProviderResponse(ctx, updateResponse.FileBasedPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ConjurPassphraseProviderResponse != nil {
			readConjurPassphraseProviderResponse(ctx, updateResponse.ConjurPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.VaultPassphraseProviderResponse != nil {
			readVaultPassphraseProviderResponse(ctx, updateResponse.VaultPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyPassphraseProviderResponse != nil {
			readThirdPartyPassphraseProviderResponse(ctx, updateResponse.ThirdPartyPassphraseProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	} else {
		tflog.Warn(ctx, "No configuration API operations created for update")
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
// This config object is edit-only, so Terraform can't delete it.
// After running a delete, Terraform will just "forget" about this object and it can be managed elsewhere.
func (r *defaultPassphraseProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *passphraseProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state passphraseProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.PassphraseProviderApi.DeletePassphraseProviderExecute(r.apiClient.PassphraseProviderApi.DeletePassphraseProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Passphrase Provider", err, httpResp)
		return
	}
}

func (r *passphraseProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPassphraseProvider(ctx, req, resp)
}

func (r *defaultPassphraseProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importPassphraseProvider(ctx, req, resp)
}

func importPassphraseProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
