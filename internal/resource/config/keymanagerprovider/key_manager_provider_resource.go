package keymanagerprovider

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
	_ resource.Resource                = &keyManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &keyManagerProviderResource{}
	_ resource.ResourceWithImportState = &keyManagerProviderResource{}
	_ resource.Resource                = &defaultKeyManagerProviderResource{}
	_ resource.ResourceWithConfigure   = &defaultKeyManagerProviderResource{}
	_ resource.ResourceWithImportState = &defaultKeyManagerProviderResource{}
)

// Create a Key Manager Provider resource
func NewKeyManagerProviderResource() resource.Resource {
	return &keyManagerProviderResource{}
}

func NewDefaultKeyManagerProviderResource() resource.Resource {
	return &defaultKeyManagerProviderResource{}
}

// keyManagerProviderResource is the resource implementation.
type keyManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultKeyManagerProviderResource is the resource implementation.
type defaultKeyManagerProviderResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *keyManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_key_manager_provider"
}

func (r *defaultKeyManagerProviderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_key_manager_provider"
}

// Configure adds the provider configured client to the resource.
func (r *keyManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

func (r *defaultKeyManagerProviderResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type keyManagerProviderResourceModel struct {
	Id                              types.String `tfsdk:"id"`
	Name                            types.String `tfsdk:"name"`
	LastUpdated                     types.String `tfsdk:"last_updated"`
	Notifications                   types.Set    `tfsdk:"notifications"`
	RequiredActions                 types.Set    `tfsdk:"required_actions"`
	Type                            types.String `tfsdk:"type"`
	ExtensionClass                  types.String `tfsdk:"extension_class"`
	ExtensionArgument               types.Set    `tfsdk:"extension_argument"`
	Pkcs11ProviderClass             types.String `tfsdk:"pkcs11_provider_class"`
	Pkcs11ProviderConfigurationFile types.String `tfsdk:"pkcs11_provider_configuration_file"`
	Pkcs11KeyStoreType              types.String `tfsdk:"pkcs11_key_store_type"`
	Pkcs11MaxCacheDuration          types.String `tfsdk:"pkcs11_max_cache_duration"`
	KeyStoreFile                    types.String `tfsdk:"key_store_file"`
	KeyStoreType                    types.String `tfsdk:"key_store_type"`
	KeyStorePin                     types.String `tfsdk:"key_store_pin"`
	KeyStorePinFile                 types.String `tfsdk:"key_store_pin_file"`
	KeyStorePinPassphraseProvider   types.String `tfsdk:"key_store_pin_passphrase_provider"`
	PrivateKeyPin                   types.String `tfsdk:"private_key_pin"`
	PrivateKeyPinFile               types.String `tfsdk:"private_key_pin_file"`
	PrivateKeyPinPassphraseProvider types.String `tfsdk:"private_key_pin_passphrase_provider"`
	Description                     types.String `tfsdk:"description"`
	Enabled                         types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *keyManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	keyManagerProviderSchema(ctx, req, resp, false)
}

func (r *defaultKeyManagerProviderResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	keyManagerProviderSchema(ctx, req, resp, true)
}

func keyManagerProviderSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Key Manager Provider.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Key Manager Provider resource. Options are ['file-based', 'custom', 'pkcs11', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"file-based", "custom", "pkcs11", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Key Manager Provider.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Key Manager Provider. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				Default:     internaltypes.EmptySetDefault(types.StringType),
				ElementType: types.StringType,
			},
			"pkcs11_provider_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java security provider class that implements support for interacting with PKCS #11 tokens.",
				Optional:    true,
			},
			"pkcs11_provider_configuration_file": schema.StringAttribute{
				Description: "The path to the file to use to configure the security provider that implements support for interacting with PKCS #11 tokens.",
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
			"pkcs11_max_cache_duration": schema.StringAttribute{
				Description: "Supported in PingDirectory product version 9.2.0.1+. The maximum length of time that data retrieved from PKCS #11 tokens may be cached for reuse. Caching might be necessary if there is noticable latency when accessing the token, for example if the token uses a remote key store. A value of zero milliseconds indicates that no caching should be performed.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"key_store_file": schema.StringAttribute{
				Description: "Specifies the path to the file that contains the private key information. This may be an absolute path, or a path that is relative to the Directory Server instance root.",
				Optional:    true,
			},
			"key_store_type": schema.StringAttribute{
				Description: "Specifies the format for the data in the key store file.",
				Optional:    true,
			},
			"key_store_pin": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `file-based`: Specifies the PIN needed to access the File Based Key Manager Provider. When the `type` attribute is set to `pkcs11`: Specifies the PIN needed to access the PKCS11 Key Manager Provider.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `file-based`: Specifies the PIN needed to access the File Based Key Manager Provider.\n  - `pkcs11`: Specifies the PIN needed to access the PKCS11 Key Manager Provider.",
				Optional:            true,
				Sensitive:           true,
			},
			"key_store_pin_file": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `file-based`: Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider. When the `type` attribute is set to `pkcs11`: Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `file-based`: Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider.\n  - `pkcs11`: Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				Optional:            true,
			},
			"key_store_pin_passphrase_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `file-based`: The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider. When the `type` attribute is set to `pkcs11`: The passphrase provider to use to obtain the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `file-based`: The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider.\n  - `pkcs11`: The passphrase provider to use to obtain the clear-text PIN needed to access the PKCS11 Key Manager Provider.",
				Optional:            true,
			},
			"private_key_pin": schema.StringAttribute{
				Description: "Specifies the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Optional:    true,
				Sensitive:   true,
			},
			"private_key_pin_file": schema.StringAttribute{
				Description: "Specifies the path to the text file whose only contents should be a single line containing the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Optional:    true,
			},
			"private_key_pin_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider to use to obtain the clear-text PIN needed to access the File Based Key Manager Provider private key. If no private key PIN is specified the PIN defaults to the key store PIN.",
				Optional:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Key Manager Provider",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Key Manager Provider is enabled for use.",
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
func (r *keyManagerProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultKeyManagerProviderResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	modifyPlanKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func modifyPlanKeyManagerProvider(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	compare, err := version.Compare(providerConfig.ProductVersion, version.PingDirectory9201)
	if err != nil {
		resp.Diagnostics.AddError("Failed to compare PingDirectory versions", err.Error())
		return
	}
	if compare >= 0 {
		// Every remaining property is supported
		return
	}
	var model keyManagerProviderResourceModel
	req.Plan.Get(ctx, &model)
	if internaltypes.IsNonEmptyString(model.Pkcs11MaxCacheDuration) {
		resp.Diagnostics.AddError("Attribute 'pkcs11_max_cache_duration' not supported by PingDirectory version "+providerConfig.ProductVersion, "")
	}
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsKeyManagerProvider() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_store_file"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_store_type"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_store_pin"),
			path.MatchRoot("type"),
			[]string{"file-based", "pkcs11"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_store_pin_file"),
			path.MatchRoot("type"),
			[]string{"file-based", "pkcs11"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("key_store_pin_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"file-based", "pkcs11"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("private_key_pin"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("private_key_pin_file"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("private_key_pin_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"file-based"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("pkcs11_provider_class"),
			path.MatchRoot("type"),
			[]string{"pkcs11"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("pkcs11_provider_configuration_file"),
			path.MatchRoot("type"),
			[]string{"pkcs11"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("pkcs11_key_store_type"),
			path.MatchRoot("type"),
			[]string{"pkcs11"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("pkcs11_max_cache_duration"),
			path.MatchRoot("type"),
			[]string{"pkcs11"},
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
			"file-based",
			[]path.Expression{path.MatchRoot("key_store_file"), path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"pkcs11",
			[]path.Expression{path.MatchRoot("enabled")},
		),
		configvalidators.ValueImpliesAttributeRequired(
			path.MatchRoot("type"),
			"third-party",
			[]path.Expression{path.MatchRoot("enabled"), path.MatchRoot("extension_class")},
		),
	}
}

// Add config validators
func (r keyManagerProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsKeyManagerProvider()
}

// Add config validators
func (r defaultKeyManagerProviderResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsKeyManagerProvider()
}

// Add optional fields to create request for file-based key-manager-provider
func addOptionalFileBasedKeyManagerProviderFields(ctx context.Context, addRequest *client.AddFileBasedKeyManagerProviderRequest, plan keyManagerProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.KeyStoreType) {
		addRequest.KeyStoreType = plan.KeyStoreType.ValueStringPointer()
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
	if internaltypes.IsNonEmptyString(plan.KeyStorePinPassphraseProvider) {
		addRequest.KeyStorePinPassphraseProvider = plan.KeyStorePinPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKeyPin) {
		addRequest.PrivateKeyPin = plan.PrivateKeyPin.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKeyPinFile) {
		addRequest.PrivateKeyPinFile = plan.PrivateKeyPinFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.PrivateKeyPinPassphraseProvider) {
		addRequest.PrivateKeyPinPassphraseProvider = plan.PrivateKeyPinPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for pkcs11 key-manager-provider
func addOptionalPkcs11KeyManagerProviderFields(ctx context.Context, addRequest *client.AddPkcs11KeyManagerProviderRequest, plan keyManagerProviderResourceModel) {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderClass) {
		addRequest.Pkcs11ProviderClass = plan.Pkcs11ProviderClass.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11ProviderConfigurationFile) {
		addRequest.Pkcs11ProviderConfigurationFile = plan.Pkcs11ProviderConfigurationFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11KeyStoreType) {
		addRequest.Pkcs11KeyStoreType = plan.Pkcs11KeyStoreType.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Pkcs11MaxCacheDuration) {
		addRequest.Pkcs11MaxCacheDuration = plan.Pkcs11MaxCacheDuration.ValueStringPointer()
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
	if internaltypes.IsNonEmptyString(plan.KeyStorePinPassphraseProvider) {
		addRequest.KeyStorePinPassphraseProvider = plan.KeyStorePinPassphraseProvider.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
}

// Add optional fields to create request for third-party key-manager-provider
func addOptionalThirdPartyKeyManagerProviderFields(ctx context.Context, addRequest *client.AddThirdPartyKeyManagerProviderRequest, plan keyManagerProviderResourceModel) {
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
func populateKeyManagerProviderUnknownValues(model *keyManagerProviderResourceModel) {
	if model.ExtensionArgument.IsUnknown() || model.ExtensionArgument.IsNull() {
		model.ExtensionArgument, _ = types.SetValue(types.StringType, []attr.Value{})
	}
	if model.Pkcs11KeyStoreType.IsUnknown() || model.Pkcs11KeyStoreType.IsNull() {
		model.Pkcs11KeyStoreType = types.StringValue("")
	}
	if model.Pkcs11MaxCacheDuration.IsUnknown() || model.Pkcs11MaxCacheDuration.IsNull() {
		model.Pkcs11MaxCacheDuration = types.StringValue("")
	}
	if model.KeyStorePin.IsUnknown() {
		model.KeyStorePin = types.StringNull()
	}
	if model.PrivateKeyPin.IsUnknown() {
		model.PrivateKeyPin = types.StringNull()
	}
}

// Populate any computed string values with empty strings, since that is equivalent to null to PD. This will reduce noise in plan output
func (model *keyManagerProviderResourceModel) populateAllComputedStringAttributes() {
	if model.Description.IsUnknown() || model.Description.IsNull() {
		model.Description = types.StringValue("")
	}
	if model.KeyStoreType.IsUnknown() || model.KeyStoreType.IsNull() {
		model.KeyStoreType = types.StringValue("")
	}
	if model.PrivateKeyPinFile.IsUnknown() || model.PrivateKeyPinFile.IsNull() {
		model.PrivateKeyPinFile = types.StringValue("")
	}
	if model.ExtensionClass.IsUnknown() || model.ExtensionClass.IsNull() {
		model.ExtensionClass = types.StringValue("")
	}
	if model.Pkcs11ProviderClass.IsUnknown() || model.Pkcs11ProviderClass.IsNull() {
		model.Pkcs11ProviderClass = types.StringValue("")
	}
	if model.KeyStorePinFile.IsUnknown() || model.KeyStorePinFile.IsNull() {
		model.KeyStorePinFile = types.StringValue("")
	}
	if model.PrivateKeyPinPassphraseProvider.IsUnknown() || model.PrivateKeyPinPassphraseProvider.IsNull() {
		model.PrivateKeyPinPassphraseProvider = types.StringValue("")
	}
	if model.Pkcs11ProviderConfigurationFile.IsUnknown() || model.Pkcs11ProviderConfigurationFile.IsNull() {
		model.Pkcs11ProviderConfigurationFile = types.StringValue("")
	}
	if model.KeyStorePinPassphraseProvider.IsUnknown() || model.KeyStorePinPassphraseProvider.IsNull() {
		model.KeyStorePinPassphraseProvider = types.StringValue("")
	}
	if model.KeyStoreFile.IsUnknown() || model.KeyStoreFile.IsNull() {
		model.KeyStoreFile = types.StringValue("")
	}
}

// Read a FileBasedKeyManagerProviderResponse object into the model struct
func readFileBasedKeyManagerProviderResponse(ctx context.Context, r *client.FileBasedKeyManagerProviderResponse, state *keyManagerProviderResourceModel, expectedValues *keyManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("file-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.KeyStoreFile = types.StringValue(r.KeyStoreFile)
	state.KeyStoreType = internaltypes.StringTypeOrNil(r.KeyStoreType, internaltypes.IsEmptyString(expectedValues.KeyStoreType))
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, internaltypes.IsEmptyString(expectedValues.KeyStorePinFile))
	state.KeyStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.KeyStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.KeyStorePinPassphraseProvider))
	state.PrivateKeyPinFile = internaltypes.StringTypeOrNil(r.PrivateKeyPinFile, internaltypes.IsEmptyString(expectedValues.PrivateKeyPinFile))
	state.PrivateKeyPinPassphraseProvider = internaltypes.StringTypeOrNil(r.PrivateKeyPinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.PrivateKeyPinPassphraseProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateKeyManagerProviderUnknownValues(state)
}

// Read a CustomKeyManagerProviderResponse object into the model struct
func readCustomKeyManagerProviderResponse(ctx context.Context, r *client.CustomKeyManagerProviderResponse, state *keyManagerProviderResourceModel, expectedValues *keyManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("custom")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateKeyManagerProviderUnknownValues(state)
}

// Read a Pkcs11KeyManagerProviderResponse object into the model struct
func readPkcs11KeyManagerProviderResponse(ctx context.Context, r *client.Pkcs11KeyManagerProviderResponse, state *keyManagerProviderResourceModel, expectedValues *keyManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("pkcs11")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Pkcs11ProviderClass = internaltypes.StringTypeOrNil(r.Pkcs11ProviderClass, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderClass))
	state.Pkcs11ProviderConfigurationFile = internaltypes.StringTypeOrNil(r.Pkcs11ProviderConfigurationFile, internaltypes.IsEmptyString(expectedValues.Pkcs11ProviderConfigurationFile))
	state.Pkcs11KeyStoreType = internaltypes.StringTypeOrNil(r.Pkcs11KeyStoreType, true)
	state.Pkcs11MaxCacheDuration = internaltypes.StringTypeOrNil(r.Pkcs11MaxCacheDuration, true)
	config.CheckMismatchedPDFormattedAttributes("pkcs11_max_cache_duration",
		expectedValues.Pkcs11MaxCacheDuration, state.Pkcs11MaxCacheDuration, diagnostics)
	state.KeyStorePinFile = internaltypes.StringTypeOrNil(r.KeyStorePinFile, internaltypes.IsEmptyString(expectedValues.KeyStorePinFile))
	state.KeyStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.KeyStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.KeyStorePinPassphraseProvider))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateKeyManagerProviderUnknownValues(state)
}

// Read a ThirdPartyKeyManagerProviderResponse object into the model struct
func readThirdPartyKeyManagerProviderResponse(ctx context.Context, r *client.ThirdPartyKeyManagerProviderResponse, state *keyManagerProviderResourceModel, expectedValues *keyManagerProviderResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateKeyManagerProviderUnknownValues(state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *keyManagerProviderResourceModel) setStateValuesNotReturnedByAPI(expectedValues *keyManagerProviderResourceModel) {
	if !expectedValues.KeyStorePin.IsUnknown() {
		state.KeyStorePin = expectedValues.KeyStorePin
	}
	if !expectedValues.PrivateKeyPin.IsUnknown() {
		state.PrivateKeyPin = expectedValues.PrivateKeyPin
	}
}

// Create any update operations necessary to make the state match the plan
func createKeyManagerProviderOperations(plan keyManagerProviderResourceModel, state keyManagerProviderResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderClass, state.Pkcs11ProviderClass, "pkcs11-provider-class")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11ProviderConfigurationFile, state.Pkcs11ProviderConfigurationFile, "pkcs11-provider-configuration-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11KeyStoreType, state.Pkcs11KeyStoreType, "pkcs11-key-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.Pkcs11MaxCacheDuration, state.Pkcs11MaxCacheDuration, "pkcs11-max-cache-duration")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStoreFile, state.KeyStoreFile, "key-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStoreType, state.KeyStoreType, "key-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePin, state.KeyStorePin, "key-store-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinFile, state.KeyStorePinFile, "key-store-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.KeyStorePinPassphraseProvider, state.KeyStorePinPassphraseProvider, "key-store-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKeyPin, state.PrivateKeyPin, "private-key-pin")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKeyPinFile, state.PrivateKeyPinFile, "private-key-pin-file")
	operations.AddStringOperationIfNecessary(&ops, plan.PrivateKeyPinPassphraseProvider, state.PrivateKeyPinPassphraseProvider, "private-key-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a file-based key-manager-provider
func (r *keyManagerProviderResource) CreateFileBasedKeyManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan keyManagerProviderResourceModel) (*keyManagerProviderResourceModel, error) {
	addRequest := client.NewAddFileBasedKeyManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumfileBasedKeyManagerProviderSchemaUrn{client.ENUMFILEBASEDKEYMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0KEY_MANAGER_PROVIDERFILE_BASED},
		plan.KeyStoreFile.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalFileBasedKeyManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.KeyManagerProviderApi.AddKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyManagerProviderRequest(
		client.AddFileBasedKeyManagerProviderRequestAsAddKeyManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.AddKeyManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Key Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state keyManagerProviderResourceModel
	readFileBasedKeyManagerProviderResponse(ctx, addResponse.FileBasedKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a pkcs11 key-manager-provider
func (r *keyManagerProviderResource) CreatePkcs11KeyManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan keyManagerProviderResourceModel) (*keyManagerProviderResourceModel, error) {
	addRequest := client.NewAddPkcs11KeyManagerProviderRequest(plan.Name.ValueString(),
		[]client.Enumpkcs11KeyManagerProviderSchemaUrn{client.ENUMPKCS11KEYMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0KEY_MANAGER_PROVIDERPKCS11},
		plan.Enabled.ValueBool())
	addOptionalPkcs11KeyManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.KeyManagerProviderApi.AddKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyManagerProviderRequest(
		client.AddPkcs11KeyManagerProviderRequestAsAddKeyManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.AddKeyManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Key Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state keyManagerProviderResourceModel
	readPkcs11KeyManagerProviderResponse(ctx, addResponse.Pkcs11KeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party key-manager-provider
func (r *keyManagerProviderResource) CreateThirdPartyKeyManagerProvider(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan keyManagerProviderResourceModel) (*keyManagerProviderResourceModel, error) {
	addRequest := client.NewAddThirdPartyKeyManagerProviderRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyKeyManagerProviderSchemaUrn{client.ENUMTHIRDPARTYKEYMANAGERPROVIDERSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0KEY_MANAGER_PROVIDERTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool())
	addOptionalThirdPartyKeyManagerProviderFields(ctx, addRequest, plan)
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.KeyManagerProviderApi.AddKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddKeyManagerProviderRequest(
		client.AddThirdPartyKeyManagerProviderRequestAsAddKeyManagerProviderRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.AddKeyManagerProviderExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Key Manager Provider", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state keyManagerProviderResourceModel
	readThirdPartyKeyManagerProviderResponse(ctx, addResponse.ThirdPartyKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *keyManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan keyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *keyManagerProviderResourceModel
	var err error
	if plan.Type.ValueString() == "file-based" {
		state, err = r.CreateFileBasedKeyManagerProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "pkcs11" {
		state, err = r.CreatePkcs11KeyManagerProvider(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyKeyManagerProvider(ctx, req, resp, plan)
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
func (r *defaultKeyManagerProviderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan keyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Manager Provider", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state keyManagerProviderResourceModel
	if readResponse.FileBasedKeyManagerProviderResponse != nil {
		readFileBasedKeyManagerProviderResponse(ctx, readResponse.FileBasedKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomKeyManagerProviderResponse != nil {
		readCustomKeyManagerProviderResponse(ctx, readResponse.CustomKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Pkcs11KeyManagerProviderResponse != nil {
		readPkcs11KeyManagerProviderResponse(ctx, readResponse.Pkcs11KeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyKeyManagerProviderResponse != nil {
		readThirdPartyKeyManagerProviderResponse(ctx, readResponse.ThirdPartyKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.KeyManagerProviderApi.UpdateKeyManagerProvider(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createKeyManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.KeyManagerProviderApi.UpdateKeyManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.FileBasedKeyManagerProviderResponse != nil {
			readFileBasedKeyManagerProviderResponse(ctx, updateResponse.FileBasedKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomKeyManagerProviderResponse != nil {
			readCustomKeyManagerProviderResponse(ctx, updateResponse.CustomKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Pkcs11KeyManagerProviderResponse != nil {
			readPkcs11KeyManagerProviderResponse(ctx, updateResponse.Pkcs11KeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyKeyManagerProviderResponse != nil {
			readThirdPartyKeyManagerProviderResponse(ctx, updateResponse.ThirdPartyKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *keyManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig, false)
}

func (r *defaultKeyManagerProviderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig, true)
}

func readKeyManagerProvider(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration, isDefault bool) {
	// Get current state
	var state keyManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.KeyManagerProviderApi.GetKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 && !isDefault {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Key Manager Provider", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Key Manager Provider", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.FileBasedKeyManagerProviderResponse != nil {
		readFileBasedKeyManagerProviderResponse(ctx, readResponse.FileBasedKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.CustomKeyManagerProviderResponse != nil {
		readCustomKeyManagerProviderResponse(ctx, readResponse.CustomKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.Pkcs11KeyManagerProviderResponse != nil {
		readPkcs11KeyManagerProviderResponse(ctx, readResponse.Pkcs11KeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyKeyManagerProviderResponse != nil {
		readThirdPartyKeyManagerProviderResponse(ctx, readResponse.ThirdPartyKeyManagerProviderResponse, &state, &state, &resp.Diagnostics)
	}

	if isDefault {
		state.populateAllComputedStringAttributes()
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *keyManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultKeyManagerProviderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateKeyManagerProvider(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateKeyManagerProvider(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan keyManagerProviderResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state keyManagerProviderResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.KeyManagerProviderApi.UpdateKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createKeyManagerProviderOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.KeyManagerProviderApi.UpdateKeyManagerProviderExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Key Manager Provider", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.FileBasedKeyManagerProviderResponse != nil {
			readFileBasedKeyManagerProviderResponse(ctx, updateResponse.FileBasedKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.CustomKeyManagerProviderResponse != nil {
			readCustomKeyManagerProviderResponse(ctx, updateResponse.CustomKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.Pkcs11KeyManagerProviderResponse != nil {
			readPkcs11KeyManagerProviderResponse(ctx, updateResponse.Pkcs11KeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyKeyManagerProviderResponse != nil {
			readThirdPartyKeyManagerProviderResponse(ctx, updateResponse.ThirdPartyKeyManagerProviderResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultKeyManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *keyManagerProviderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state keyManagerProviderResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.KeyManagerProviderApi.DeleteKeyManagerProviderExecute(r.apiClient.KeyManagerProviderApi.DeleteKeyManagerProvider(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Key Manager Provider", err, httpResp)
		return
	}
}

func (r *keyManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKeyManagerProvider(ctx, req, resp)
}

func (r *defaultKeyManagerProviderResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importKeyManagerProvider(ctx, req, resp)
}

func importKeyManagerProvider(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
