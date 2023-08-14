package webapplicationextension

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
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
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &webApplicationExtensionResource{}
	_ resource.ResourceWithConfigure   = &webApplicationExtensionResource{}
	_ resource.ResourceWithImportState = &webApplicationExtensionResource{}
	_ resource.Resource                = &defaultWebApplicationExtensionResource{}
	_ resource.ResourceWithConfigure   = &defaultWebApplicationExtensionResource{}
	_ resource.ResourceWithImportState = &defaultWebApplicationExtensionResource{}
)

// Create a Web Application Extension resource
func NewWebApplicationExtensionResource() resource.Resource {
	return &webApplicationExtensionResource{}
}

func NewDefaultWebApplicationExtensionResource() resource.Resource {
	return &defaultWebApplicationExtensionResource{}
}

// webApplicationExtensionResource is the resource implementation.
type webApplicationExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultWebApplicationExtensionResource is the resource implementation.
type defaultWebApplicationExtensionResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *webApplicationExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_application_extension"
}

func (r *defaultWebApplicationExtensionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_web_application_extension"
}

// Configure adds the provider configured client to the resource.
func (r *webApplicationExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultWebApplicationExtensionResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type webApplicationExtensionResourceModel struct {
	Id                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	LastUpdated              types.String `tfsdk:"last_updated"`
	Notifications            types.Set    `tfsdk:"notifications"`
	RequiredActions          types.Set    `tfsdk:"required_actions"`
	Type                     types.String `tfsdk:"type"`
	Description              types.String `tfsdk:"description"`
	BaseContextPath          types.String `tfsdk:"base_context_path"`
	WarFile                  types.String `tfsdk:"war_file"`
	DocumentRootDirectory    types.String `tfsdk:"document_root_directory"`
	DeploymentDescriptorFile types.String `tfsdk:"deployment_descriptor_file"`
	TemporaryDirectory       types.String `tfsdk:"temporary_directory"`
	InitParameter            types.Set    `tfsdk:"init_parameter"`
}

type defaultWebApplicationExtensionResourceModel struct {
	Id                                  types.String `tfsdk:"id"`
	Name                                types.String `tfsdk:"name"`
	LastUpdated                         types.String `tfsdk:"last_updated"`
	Notifications                       types.Set    `tfsdk:"notifications"`
	RequiredActions                     types.Set    `tfsdk:"required_actions"`
	Type                                types.String `tfsdk:"type"`
	SsoEnabled                          types.Bool   `tfsdk:"sso_enabled"`
	OidcClientID                        types.String `tfsdk:"oidc_client_id"`
	OidcClientSecret                    types.String `tfsdk:"oidc_client_secret"`
	OidcClientSecretPassphraseProvider  types.String `tfsdk:"oidc_client_secret_passphrase_provider"`
	OidcIssuerURL                       types.String `tfsdk:"oidc_issuer_url"`
	OidcTrustStoreFile                  types.String `tfsdk:"oidc_trust_store_file"`
	OidcTrustStoreType                  types.String `tfsdk:"oidc_trust_store_type"`
	OidcTrustStorePinPassphraseProvider types.String `tfsdk:"oidc_trust_store_pin_passphrase_provider"`
	OidcStrictHostnameVerification      types.Bool   `tfsdk:"oidc_strict_hostname_verification"`
	OidcTrustAll                        types.Bool   `tfsdk:"oidc_trust_all"`
	LdapServer                          types.String `tfsdk:"ldap_server"`
	TrustStoreFile                      types.String `tfsdk:"trust_store_file"`
	TrustStoreType                      types.String `tfsdk:"trust_store_type"`
	TrustStorePinPassphraseProvider     types.String `tfsdk:"trust_store_pin_passphrase_provider"`
	LogFile                             types.String `tfsdk:"log_file"`
	Complexity                          types.String `tfsdk:"complexity"`
	Description                         types.String `tfsdk:"description"`
	BaseContextPath                     types.String `tfsdk:"base_context_path"`
	WarFile                             types.String `tfsdk:"war_file"`
	DocumentRootDirectory               types.String `tfsdk:"document_root_directory"`
	DeploymentDescriptorFile            types.String `tfsdk:"deployment_descriptor_file"`
	TemporaryDirectory                  types.String `tfsdk:"temporary_directory"`
	InitParameter                       types.Set    `tfsdk:"init_parameter"`
}

// GetSchema defines the schema for the resource.
func (r *webApplicationExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	webApplicationExtensionSchema(ctx, req, resp, false)
}

func (r *defaultWebApplicationExtensionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	webApplicationExtensionSchema(ctx, req, resp, true)
}

func webApplicationExtensionSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Web Application Extension.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Web Application Extension resource. Options are ['console', 'generic']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"generic"}...),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Web Application Extension",
				Optional:    true,
			},
			"base_context_path": schema.StringAttribute{
				Description: "Specifies the base context path that should be used by HTTP clients to reference content. The value must start with a forward slash and at least one additional character and must represent a valid HTTP context path.",
				Required:    true,
			},
			"war_file": schema.StringAttribute{
				Description: "Specifies the path to a standard web application archive (WAR) file.",
				Optional:    true,
			},
			"document_root_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory on the local filesystem containing the files to be served by this Web Application Extension. The path must exist, and it must be a directory.",
				Optional:    true,
			},
			"deployment_descriptor_file": schema.StringAttribute{
				Description: "Specifies the path to the deployment descriptor file when used with document-root-directory.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"temporary_directory": schema.StringAttribute{
				Description: "Specifies the path to the directory that may be used to store temporary files such as extracted WAR files and compiled JSP files.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"init_parameter": schema.SetAttribute{
				Description: "Specifies an initialization parameter to pass into the web application during startup.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		typeAttr.Validators = []validator.String{
			stringvalidator.OneOf([]string{"console", "generic"}...),
		}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		schemaDef.Attributes["sso_enabled"] = schema.BoolAttribute{
			Description: "Indicates that SSO login into the Administrative Console is enabled.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["oidc_client_id"] = schema.StringAttribute{
			Description: "The client ID to use when authenticating to the OpenID Connect provider.",
			Optional:    true,
		}
		schemaDef.Attributes["oidc_client_secret"] = schema.StringAttribute{
			Description: "The client secret to use when authenticating to the OpenID Connect provider.",
			Optional:    true,
			Sensitive:   true,
		}
		schemaDef.Attributes["oidc_client_secret_passphrase_provider"] = schema.StringAttribute{
			Description: "A passphrase provider that may be used to obtain the client secret to use when authenticating to the OpenID Connect provider.",
			Optional:    true,
		}
		schemaDef.Attributes["oidc_issuer_url"] = schema.StringAttribute{
			Description: "The issuer URL of the OpenID Connect provider.",
			Optional:    true,
		}
		schemaDef.Attributes["oidc_trust_store_file"] = schema.StringAttribute{
			Description: "Specifies the path to the truststore file used by this application to evaluate OIDC provider certificates. If this field is left blank, the default JVM trust store will be used.",
			Optional:    true,
		}
		schemaDef.Attributes["oidc_trust_store_type"] = schema.StringAttribute{
			Description: "Specifies the format for the data in the OIDC trust store file.",
			Optional:    true,
		}
		schemaDef.Attributes["oidc_trust_store_pin_passphrase_provider"] = schema.StringAttribute{
			Description: "The passphrase provider that may be used to obtain the PIN for the trust store used with OIDC providers. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
			Optional:    true,
		}
		schemaDef.Attributes["oidc_strict_hostname_verification"] = schema.BoolAttribute{
			Description: "Controls whether or not hostname verification is performed, which checks if the hostname of the OIDC provider matches the name(s) stored inside the certificate it provides. This property should only be set to false for testing purposes.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["oidc_trust_all"] = schema.BoolAttribute{
			Description: "Controls whether or not this application will always trust any certificate that is presented to it, regardless of its contents. This property should only be set to true for testing purposes.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["ldap_server"] = schema.StringAttribute{
			Description: "The LDAP URL used to connect to the managed server.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["trust_store_file"] = schema.StringAttribute{
			Description: "Specifies the path to the truststore file, which is used by this application to establish trust of managed servers.",
			Optional:    true,
		}
		schemaDef.Attributes["trust_store_type"] = schema.StringAttribute{
			Description: "Specifies the format for the data in the trust store file.",
			Optional:    true,
		}
		schemaDef.Attributes["trust_store_pin_passphrase_provider"] = schema.StringAttribute{
			Description: "The passphrase provider that may be used to obtain the PIN for the trust store used with managed LDAP servers. This is only required if a trust store file is required, and if that trust store requires a PIN to access its contents.",
			Optional:    true,
		}
		schemaDef.Attributes["log_file"] = schema.StringAttribute{
			Description: "The path to the log file for the web application.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		schemaDef.Attributes["complexity"] = schema.StringAttribute{
			Description: "Specifies the maximum complexity level for managed configuration elements.",
			Optional:    true,
			Computed:    true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		}
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsWebApplicationExtension() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			path.MatchRoot("war_file"),
			path.MatchRoot("document_root_directory"),
		),
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"console"},
			resourcevalidator.Conflicting(
				path.MatchRoot("oidc_client_secret"),
				path.MatchRoot("oidc_client_secret_passphrase_provider"),
			),
		),
	}
}

// Add config validators
func (r webApplicationExtensionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsWebApplicationExtension()
}

// Add config validators
func (r defaultWebApplicationExtensionResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	validators := []resource.ConfigValidator{
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("complexity"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_trust_all"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("sso_enabled"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_issuer_url"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_client_id"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_pin_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_strict_hostname_verification"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_type"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_client_secret"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_trust_store_pin_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("trust_store_file"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_client_secret_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("log_file"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_trust_store_type"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("oidc_trust_store_file"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("ldap_server"),
			path.MatchRoot("type"),
			[]string{"console"},
		),
	}
	return append(configValidatorsWebApplicationExtension(), validators...)
}

// Add optional fields to create request for generic web-application-extension
func addOptionalGenericWebApplicationExtensionFields(ctx context.Context, addRequest *client.AddGenericWebApplicationExtensionRequest, plan webApplicationExtensionResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.WarFile) {
		addRequest.WarFile = plan.WarFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DocumentRootDirectory) {
		addRequest.DocumentRootDirectory = plan.DocumentRootDirectory.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.DeploymentDescriptorFile) {
		addRequest.DeploymentDescriptorFile = plan.DeploymentDescriptorFile.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.TemporaryDirectory) {
		addRequest.TemporaryDirectory = plan.TemporaryDirectory.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.InitParameter) {
		var slice []string
		plan.InitParameter.ElementsAs(ctx, &slice, false)
		addRequest.InitParameter = slice
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateWebApplicationExtensionUnknownValuesDefault(ctx context.Context, model *defaultWebApplicationExtensionResourceModel) {
	if model.OidcClientSecret.IsUnknown() {
		model.OidcClientSecret = types.StringNull()
	}
}

// Read a ConsoleWebApplicationExtensionResponse object into the model struct
func readConsoleWebApplicationExtensionResponseDefault(ctx context.Context, r *client.ConsoleWebApplicationExtensionResponse, state *defaultWebApplicationExtensionResourceModel, expectedValues *defaultWebApplicationExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("console")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.SsoEnabled = internaltypes.BoolTypeOrNil(r.SsoEnabled)
	state.OidcClientID = internaltypes.StringTypeOrNil(r.OidcClientID, internaltypes.IsEmptyString(expectedValues.OidcClientID))
	state.OidcClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.OidcClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.OidcClientSecretPassphraseProvider))
	state.OidcIssuerURL = internaltypes.StringTypeOrNil(r.OidcIssuerURL, internaltypes.IsEmptyString(expectedValues.OidcIssuerURL))
	state.OidcTrustStoreFile = internaltypes.StringTypeOrNil(r.OidcTrustStoreFile, internaltypes.IsEmptyString(expectedValues.OidcTrustStoreFile))
	state.OidcTrustStoreType = internaltypes.StringTypeOrNil(r.OidcTrustStoreType, internaltypes.IsEmptyString(expectedValues.OidcTrustStoreType))
	state.OidcTrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.OidcTrustStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.OidcTrustStorePinPassphraseProvider))
	state.OidcStrictHostnameVerification = internaltypes.BoolTypeOrNil(r.OidcStrictHostnameVerification)
	state.OidcTrustAll = internaltypes.BoolTypeOrNil(r.OidcTrustAll)
	state.LdapServer = internaltypes.StringTypeOrNil(r.LdapServer, internaltypes.IsEmptyString(expectedValues.LdapServer))
	state.TrustStoreFile = internaltypes.StringTypeOrNil(r.TrustStoreFile, internaltypes.IsEmptyString(expectedValues.TrustStoreFile))
	state.TrustStoreType = internaltypes.StringTypeOrNil(r.TrustStoreType, internaltypes.IsEmptyString(expectedValues.TrustStoreType))
	state.TrustStorePinPassphraseProvider = internaltypes.StringTypeOrNil(r.TrustStorePinPassphraseProvider, internaltypes.IsEmptyString(expectedValues.TrustStorePinPassphraseProvider))
	state.LogFile = internaltypes.StringTypeOrNil(r.LogFile, internaltypes.IsEmptyString(expectedValues.LogFile))
	state.Complexity = internaltypes.StringTypeOrNil(
		client.StringPointerEnumwebApplicationExtensionComplexityProp(r.Complexity), internaltypes.IsEmptyString(expectedValues.Complexity))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, internaltypes.IsEmptyString(expectedValues.WarFile))
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, internaltypes.IsEmptyString(expectedValues.DocumentRootDirectory))
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, internaltypes.IsEmptyString(expectedValues.DeploymentDescriptorFile))
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, internaltypes.IsEmptyString(expectedValues.TemporaryDirectory))
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateWebApplicationExtensionUnknownValuesDefault(ctx, state)
}

// Read a GenericWebApplicationExtensionResponse object into the model struct
func readGenericWebApplicationExtensionResponse(ctx context.Context, r *client.GenericWebApplicationExtensionResponse, state *webApplicationExtensionResourceModel, expectedValues *webApplicationExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, internaltypes.IsEmptyString(expectedValues.WarFile))
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, internaltypes.IsEmptyString(expectedValues.DocumentRootDirectory))
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, internaltypes.IsEmptyString(expectedValues.DeploymentDescriptorFile))
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, internaltypes.IsEmptyString(expectedValues.TemporaryDirectory))
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
}

// Read a GenericWebApplicationExtensionResponse object into the model struct
func readGenericWebApplicationExtensionResponseDefault(ctx context.Context, r *client.GenericWebApplicationExtensionResponse, state *defaultWebApplicationExtensionResourceModel, expectedValues *defaultWebApplicationExtensionResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.BaseContextPath = types.StringValue(r.BaseContextPath)
	state.WarFile = internaltypes.StringTypeOrNil(r.WarFile, internaltypes.IsEmptyString(expectedValues.WarFile))
	state.DocumentRootDirectory = internaltypes.StringTypeOrNil(r.DocumentRootDirectory, internaltypes.IsEmptyString(expectedValues.DocumentRootDirectory))
	state.DeploymentDescriptorFile = internaltypes.StringTypeOrNil(r.DeploymentDescriptorFile, internaltypes.IsEmptyString(expectedValues.DeploymentDescriptorFile))
	state.TemporaryDirectory = internaltypes.StringTypeOrNil(r.TemporaryDirectory, internaltypes.IsEmptyString(expectedValues.TemporaryDirectory))
	state.InitParameter = internaltypes.GetStringSet(r.InitParameter)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateWebApplicationExtensionUnknownValuesDefault(ctx, state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *defaultWebApplicationExtensionResourceModel) setStateValuesNotReturnedByAPI(expectedValues *defaultWebApplicationExtensionResourceModel) {
	if !expectedValues.OidcClientSecret.IsUnknown() {
		state.OidcClientSecret = expectedValues.OidcClientSecret
	}
}

// Create any update operations necessary to make the state match the plan
func createWebApplicationExtensionOperations(plan webApplicationExtensionResourceModel, state webApplicationExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.WarFile, state.WarFile, "war-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DeploymentDescriptorFile, state.DeploymentDescriptorFile, "deployment-descriptor-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TemporaryDirectory, state.TemporaryDirectory, "temporary-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InitParameter, state.InitParameter, "init-parameter")
	return ops
}

// Create any update operations necessary to make the state match the plan
func createWebApplicationExtensionOperationsDefault(plan defaultWebApplicationExtensionResourceModel, state defaultWebApplicationExtensionResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddBoolOperationIfNecessary(&ops, plan.SsoEnabled, state.SsoEnabled, "sso-enabled")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcClientID, state.OidcClientID, "oidc-client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcClientSecret, state.OidcClientSecret, "oidc-client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcClientSecretPassphraseProvider, state.OidcClientSecretPassphraseProvider, "oidc-client-secret-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcIssuerURL, state.OidcIssuerURL, "oidc-issuer-url")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcTrustStoreFile, state.OidcTrustStoreFile, "oidc-trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcTrustStoreType, state.OidcTrustStoreType, "oidc-trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.OidcTrustStorePinPassphraseProvider, state.OidcTrustStorePinPassphraseProvider, "oidc-trust-store-pin-passphrase-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.OidcStrictHostnameVerification, state.OidcStrictHostnameVerification, "oidc-strict-hostname-verification")
	operations.AddBoolOperationIfNecessary(&ops, plan.OidcTrustAll, state.OidcTrustAll, "oidc-trust-all")
	operations.AddStringOperationIfNecessary(&ops, plan.LdapServer, state.LdapServer, "ldap-server")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreFile, state.TrustStoreFile, "trust-store-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStoreType, state.TrustStoreType, "trust-store-type")
	operations.AddStringOperationIfNecessary(&ops, plan.TrustStorePinPassphraseProvider, state.TrustStorePinPassphraseProvider, "trust-store-pin-passphrase-provider")
	operations.AddStringOperationIfNecessary(&ops, plan.LogFile, state.LogFile, "log-file")
	operations.AddStringOperationIfNecessary(&ops, plan.Complexity, state.Complexity, "complexity")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddStringOperationIfNecessary(&ops, plan.BaseContextPath, state.BaseContextPath, "base-context-path")
	operations.AddStringOperationIfNecessary(&ops, plan.WarFile, state.WarFile, "war-file")
	operations.AddStringOperationIfNecessary(&ops, plan.DocumentRootDirectory, state.DocumentRootDirectory, "document-root-directory")
	operations.AddStringOperationIfNecessary(&ops, plan.DeploymentDescriptorFile, state.DeploymentDescriptorFile, "deployment-descriptor-file")
	operations.AddStringOperationIfNecessary(&ops, plan.TemporaryDirectory, state.TemporaryDirectory, "temporary-directory")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.InitParameter, state.InitParameter, "init-parameter")
	return ops
}

// Create a generic web-application-extension
func (r *webApplicationExtensionResource) CreateGenericWebApplicationExtension(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan webApplicationExtensionResourceModel) (*webApplicationExtensionResourceModel, error) {
	addRequest := client.NewAddGenericWebApplicationExtensionRequest(plan.Name.ValueString(),
		[]client.EnumgenericWebApplicationExtensionSchemaUrn{client.ENUMGENERICWEBAPPLICATIONEXTENSIONSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0WEB_APPLICATION_EXTENSIONGENERIC},
		plan.BaseContextPath.ValueString())
	err := addOptionalGenericWebApplicationExtensionFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Web Application Extension", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.WebApplicationExtensionApi.AddWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddGenericWebApplicationExtensionRequest(*addRequest)

	addResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.AddWebApplicationExtensionExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Web Application Extension", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state webApplicationExtensionResourceModel
	readGenericWebApplicationExtensionResponse(ctx, addResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *webApplicationExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan webApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.CreateGenericWebApplicationExtension(ctx, req, resp, plan)
	if err != nil {
		return
	}

	// Populate Computed attribute values
	state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))

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
func (r *defaultWebApplicationExtensionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan defaultWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state defaultWebApplicationExtensionResourceModel
	if readResponse.ConsoleWebApplicationExtensionResponse != nil {
		readConsoleWebApplicationExtensionResponseDefault(ctx, readResponse.ConsoleWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.GenericWebApplicationExtensionResponse != nil {
		readGenericWebApplicationExtensionResponseDefault(ctx, readResponse.GenericWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createWebApplicationExtensionOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ConsoleWebApplicationExtensionResponse != nil {
			readConsoleWebApplicationExtensionResponseDefault(ctx, updateResponse.ConsoleWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GenericWebApplicationExtensionResponse != nil {
			readGenericWebApplicationExtensionResponseDefault(ctx, updateResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		// Update computed values
		state.LastUpdated = types.StringValue(string(time.Now().Format(time.RFC850)))
	}

	state.setStateValuesNotReturnedByAPI(&plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r *webApplicationExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state webApplicationExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		if httpResp.StatusCode == 404 {
			config.ReportHttpErrorAsWarning(ctx, &resp.Diagnostics, "An error occurred while getting the Web Application Extension", err, httpResp)
			resp.State.RemoveResource(ctx)
		} else {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Web Application Extension", err, httpResp)
		}
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.GenericWebApplicationExtensionResponse != nil {
		readGenericWebApplicationExtensionResponse(ctx, readResponse.GenericWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *defaultWebApplicationExtensionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state defaultWebApplicationExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.GetWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Web Application Extension", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ConsoleWebApplicationExtensionResponse != nil {
		readConsoleWebApplicationExtensionResponseDefault(ctx, readResponse.ConsoleWebApplicationExtensionResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *webApplicationExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan webApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state webApplicationExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createWebApplicationExtensionOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.GenericWebApplicationExtensionResponse != nil {
			readGenericWebApplicationExtensionResponse(ctx, updateResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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

func (r *defaultWebApplicationExtensionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan defaultWebApplicationExtensionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state defaultWebApplicationExtensionResourceModel
	req.State.Get(ctx, &state)
	updateRequest := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createWebApplicationExtensionOperationsDefault(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.WebApplicationExtensionApi.UpdateWebApplicationExtensionExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Web Application Extension", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.ConsoleWebApplicationExtensionResponse != nil {
			readConsoleWebApplicationExtensionResponseDefault(ctx, updateResponse.ConsoleWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.GenericWebApplicationExtensionResponse != nil {
			readGenericWebApplicationExtensionResponseDefault(ctx, updateResponse.GenericWebApplicationExtensionResponse, &state, &plan, &resp.Diagnostics)
		}
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
func (r *defaultWebApplicationExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *webApplicationExtensionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state webApplicationExtensionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.WebApplicationExtensionApi.DeleteWebApplicationExtensionExecute(r.apiClient.WebApplicationExtensionApi.DeleteWebApplicationExtension(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil && httpResp.StatusCode != 404 {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Web Application Extension", err, httpResp)
		return
	}
}

func (r *webApplicationExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importWebApplicationExtension(ctx, req, resp)
}

func (r *defaultWebApplicationExtensionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importWebApplicationExtension(ctx, req, resp)
}

func importWebApplicationExtension(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
