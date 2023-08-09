package accesstokenvalidator

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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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
	_ resource.Resource                = &accessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &accessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &accessTokenValidatorResource{}
	_ resource.Resource                = &defaultAccessTokenValidatorResource{}
	_ resource.ResourceWithConfigure   = &defaultAccessTokenValidatorResource{}
	_ resource.ResourceWithImportState = &defaultAccessTokenValidatorResource{}
)

// Create a Access Token Validator resource
func NewAccessTokenValidatorResource() resource.Resource {
	return &accessTokenValidatorResource{}
}

func NewDefaultAccessTokenValidatorResource() resource.Resource {
	return &defaultAccessTokenValidatorResource{}
}

// accessTokenValidatorResource is the resource implementation.
type accessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// defaultAccessTokenValidatorResource is the resource implementation.
type defaultAccessTokenValidatorResource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the resource type name.
func (r *accessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_access_token_validator"
}

func (r *defaultAccessTokenValidatorResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_default_access_token_validator"
}

// Configure adds the provider configured client to the resource.
func (r *accessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

func (r *defaultAccessTokenValidatorResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type accessTokenValidatorResourceModel struct {
	Id                                types.String `tfsdk:"id"`
	Name                              types.String `tfsdk:"name"`
	LastUpdated                       types.String `tfsdk:"last_updated"`
	Notifications                     types.Set    `tfsdk:"notifications"`
	RequiredActions                   types.Set    `tfsdk:"required_actions"`
	Type                              types.String `tfsdk:"type"`
	ExtensionClass                    types.String `tfsdk:"extension_class"`
	ExtensionArgument                 types.Set    `tfsdk:"extension_argument"`
	AllowedSigningAlgorithm           types.Set    `tfsdk:"allowed_signing_algorithm"`
	SigningCertificate                types.Set    `tfsdk:"signing_certificate"`
	JwksEndpointPath                  types.String `tfsdk:"jwks_endpoint_path"`
	EncryptionKeyPair                 types.String `tfsdk:"encryption_key_pair"`
	AllowedKeyEncryptionAlgorithm     types.Set    `tfsdk:"allowed_key_encryption_algorithm"`
	AllowedContentEncryptionAlgorithm types.Set    `tfsdk:"allowed_content_encryption_algorithm"`
	ClockSkewGracePeriod              types.String `tfsdk:"clock_skew_grace_period"`
	ClientIDClaimName                 types.String `tfsdk:"client_id_claim_name"`
	ScopeClaimName                    types.String `tfsdk:"scope_claim_name"`
	ClientID                          types.String `tfsdk:"client_id"`
	ClientSecret                      types.String `tfsdk:"client_secret"`
	ClientSecretPassphraseProvider    types.String `tfsdk:"client_secret_passphrase_provider"`
	IncludeAudParameter               types.Bool   `tfsdk:"include_aud_parameter"`
	AccessTokenManagerID              types.String `tfsdk:"access_token_manager_id"`
	EndpointCacheRefresh              types.String `tfsdk:"endpoint_cache_refresh"`
	EvaluationOrderIndex              types.Int64  `tfsdk:"evaluation_order_index"`
	AuthorizationServer               types.String `tfsdk:"authorization_server"`
	IdentityMapper                    types.String `tfsdk:"identity_mapper"`
	SubjectClaimName                  types.String `tfsdk:"subject_claim_name"`
	Description                       types.String `tfsdk:"description"`
	Enabled                           types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the resource.
func (r *accessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	accessTokenValidatorSchema(ctx, req, resp, false)
}

func (r *defaultAccessTokenValidatorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	accessTokenValidatorSchema(ctx, req, resp, true)
}

func accessTokenValidatorSchema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse, isDefault bool) {
	schemaDef := schema.Schema{
		Description: "Manages a Access Token Validator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Access Token Validator resource. Options are ['ping-federate', 'jwt', 'mock', 'third-party']",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"ping-federate", "jwt", "mock", "third-party"}...),
				},
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Access Token Validator.",
				Optional:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Access Token Validator. Each configuration property should be given in the form 'name=value'.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_signing_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT signing algorithms that will be accepted by the JWT Access Token Validator.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"signing_certificate": schema.SetAttribute{
				Description: "Specifies the locally stored certificates that may be used to validate the signature of an incoming JWT access token. If this property is specified, the JWT Access Token Validator will not use a JWKS endpoint to retrieve public keys.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"jwks_endpoint_path": schema.StringAttribute{
				Description: "The relative path to JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming JWT access token. This path is relative to the base_url property defined for the validator's external authorization server. If jwks-endpoint-path is specified, the JWT Access Token Validator will not consult locally stored certificates for validating token signatures.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_key_pair": schema.StringAttribute{
				Description: "The public-private key pair that is used to encrypt the JWT payload. If specified, the JWT Access Token Validator will use the private key to decrypt the JWT payload, and the public key must be exported to the Authorization Server that is issuing access tokens.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_key_encryption_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT key encryption algorithms that will be accepted by the JWT Access Token Validator. This setting is only used if encryption-key-pair is set.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"allowed_content_encryption_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT content encryption algorithms that will be accepted by the JWT Access Token Validator.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"clock_skew_grace_period": schema.StringAttribute{
				Description: "Specifies the amount of clock skew that is tolerated by the JWT Access Token Validator when evaluating whether a token is within its valid time interval. The duration specified by this parameter will be subtracted from the token's not-before (nbf) time and added to the token's expiration (exp) time, if present, to allow for any time difference between the local server's clock and the token issuer's clock.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the OAuth2 client Id.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"scope_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the scopes granted by the token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"client_id": schema.StringAttribute{
				Description: "The client identifier to use when authenticating to the PingFederate authorization server.",
				Optional:    true,
			},
			"client_secret": schema.StringAttribute{
				Description: "The client secret to use when authenticating to the PingFederate authorization server.",
				Optional:    true,
				Sensitive:   true,
			},
			"client_secret_passphrase_provider": schema.StringAttribute{
				Description: "The passphrase provider for obtaining the client secret to use when authenticating to the PingFederate authorization server.",
				Optional:    true,
			},
			"include_aud_parameter": schema.BoolAttribute{
				Description: "Whether to include the incoming request URL as the \"aud\" parameter when calling the PingFederate introspection endpoint. This property is ignored if the access-token-manager-id property is set.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"access_token_manager_id": schema.StringAttribute{
				Description: "The Access Token Manager instance ID to specify when calling the PingFederate introspection endpoint. If this property is set the include-aud-parameter property is ignored.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"endpoint_cache_refresh": schema.StringAttribute{
				Description: "How often the Access Token Validator should refresh its stored value of the PingFederate server's token introspection endpoint.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple Ping Federate Access Token Validators are defined for a single Directory Server, this property determines the evaluation order for determining the correct validator class for an access token received by the Directory Server. Values of this property must be unique among all Ping Federate Access Token Validators defined within Directory Server but not necessarily contiguous. Ping Federate Access Token Validators with a smaller value will be evaluated first to determine if they are able to validate the access token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"authorization_server": schema.StringAttribute{
				Description: "Specifies the external server that will be used to aid in validating access tokens. In most cases this will be the Authorization Server that minted the token.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used for associating user entries with Bearer token subject names. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject, i.e. the logged-in user in an access token. This property goes hand-in-hand with the identity-mapper property and tells the Identity Mapper which field to use to look up the user entry on the server.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("sub"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Description: "A description for this Access Token Validator",
				Optional:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this Access Token Validator is enabled for use in Directory Server.",
				Required:    true,
			},
		},
	}
	if isDefault {
		typeAttr := schemaDef.Attributes["type"].(schema.StringAttribute)
		typeAttr.Optional = false
		typeAttr.Required = false
		typeAttr.Computed = true
		typeAttr.PlanModifiers = []planmodifier.String{}
		schemaDef.Attributes["type"] = typeAttr
		// Add any default properties and set optional properties to computed where necessary
		config.SetAttributesToOptionalAndComputedAndRemoveDefaults(&schemaDef, []string{"type"})
	}
	config.AddCommonResourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Add config validators that apply to both default_ and non-default_
func configValidatorsAccessTokenValidator() []resource.ConfigValidator {
	return []resource.ConfigValidator{
		configvalidators.ImpliesOtherValidator(
			path.MatchRoot("type"),
			[]string{"ping-federate"},
			resourcevalidator.ExactlyOneOf(
				path.MatchRoot("client_secret"),
				path.MatchRoot("client_secret_passphrase_provider"),
			),
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("client_secret_passphrase_provider"),
			path.MatchRoot("type"),
			[]string{"ping-federate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("access_token_manager_id"),
			path.MatchRoot("type"),
			[]string{"ping-federate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_content_encryption_algorithm"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_signing_algorithm"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("signing_certificate"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("scope_claim_name"),
			path.MatchRoot("type"),
			[]string{"jwt", "mock"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("encryption_key_pair"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("endpoint_cache_refresh"),
			path.MatchRoot("type"),
			[]string{"ping-federate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("jwks_endpoint_path"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("include_aud_parameter"),
			path.MatchRoot("type"),
			[]string{"ping-federate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("client_secret"),
			path.MatchRoot("type"),
			[]string{"ping-federate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("authorization_server"),
			path.MatchRoot("type"),
			[]string{"ping-federate", "jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_argument"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("clock_skew_grace_period"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("extension_class"),
			path.MatchRoot("type"),
			[]string{"third-party"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("allowed_key_encryption_algorithm"),
			path.MatchRoot("type"),
			[]string{"jwt"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("client_id"),
			path.MatchRoot("type"),
			[]string{"ping-federate"},
		),
		configvalidators.ImpliesOtherAttributeOneOfString(
			path.MatchRoot("client_id_claim_name"),
			path.MatchRoot("type"),
			[]string{"jwt", "mock"},
		),
	}
}

// Add config validators
func (r accessTokenValidatorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsAccessTokenValidator()
}

// Add config validators
func (r defaultAccessTokenValidatorResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return configValidatorsAccessTokenValidator()
}

// Add optional fields to create request for ping-federate access-token-validator
func addOptionalPingFederateAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddPingFederateAccessTokenValidatorRequest, plan accessTokenValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientSecret) {
		addRequest.ClientSecret = plan.ClientSecret.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientSecretPassphraseProvider) {
		addRequest.ClientSecretPassphraseProvider = plan.ClientSecretPassphraseProvider.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.IncludeAudParameter) {
		addRequest.IncludeAudParameter = plan.IncludeAudParameter.ValueBoolPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AccessTokenManagerID) {
		addRequest.AccessTokenManagerID = plan.AccessTokenManagerID.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EndpointCacheRefresh) {
		addRequest.EndpointCacheRefresh = plan.EndpointCacheRefresh.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		addRequest.EvaluationOrderIndex = plan.EvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthorizationServer) {
		addRequest.AuthorizationServer = plan.AuthorizationServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		addRequest.SubjectClaimName = plan.SubjectClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for jwt access-token-validator
func addOptionalJwtAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddJwtAccessTokenValidatorRequest, plan accessTokenValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.AllowedSigningAlgorithm) {
		var slice []string
		plan.AllowedSigningAlgorithm.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumaccessTokenValidatorAllowedSigningAlgorithmProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumaccessTokenValidatorAllowedSigningAlgorithmPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedSigningAlgorithm = enumSlice
	}
	if internaltypes.IsDefined(plan.SigningCertificate) {
		var slice []string
		plan.SigningCertificate.ElementsAs(ctx, &slice, false)
		addRequest.SigningCertificate = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.JwksEndpointPath) {
		addRequest.JwksEndpointPath = plan.JwksEndpointPath.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.EncryptionKeyPair) {
		addRequest.EncryptionKeyPair = plan.EncryptionKeyPair.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.AllowedKeyEncryptionAlgorithm) {
		var slice []string
		plan.AllowedKeyEncryptionAlgorithm.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedKeyEncryptionAlgorithm = enumSlice
	}
	if internaltypes.IsDefined(plan.AllowedContentEncryptionAlgorithm) {
		var slice []string
		plan.AllowedContentEncryptionAlgorithm.ElementsAs(ctx, &slice, false)
		enumSlice := make([]client.EnumaccessTokenValidatorAllowedContentEncryptionAlgorithmProp, len(slice))
		for i := 0; i < len(slice); i++ {
			enumVal, err := client.NewEnumaccessTokenValidatorAllowedContentEncryptionAlgorithmPropFromValue(slice[i])
			if err != nil {
				return err
			}
			enumSlice[i] = *enumVal
		}
		addRequest.AllowedContentEncryptionAlgorithm = enumSlice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClockSkewGracePeriod) {
		addRequest.ClockSkewGracePeriod = plan.ClockSkewGracePeriod.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientIDClaimName) {
		addRequest.ClientIDClaimName = plan.ClientIDClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ScopeClaimName) {
		addRequest.ScopeClaimName = plan.ScopeClaimName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		addRequest.EvaluationOrderIndex = plan.EvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.AuthorizationServer) {
		addRequest.AuthorizationServer = plan.AuthorizationServer.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		addRequest.SubjectClaimName = plan.SubjectClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for mock access-token-validator
func addOptionalMockAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddMockAccessTokenValidatorRequest, plan accessTokenValidatorResourceModel) error {
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ClientIDClaimName) {
		addRequest.ClientIDClaimName = plan.ClientIDClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.ScopeClaimName) {
		addRequest.ScopeClaimName = plan.ScopeClaimName.ValueStringPointer()
	}
	if internaltypes.IsDefined(plan.EvaluationOrderIndex) {
		addRequest.EvaluationOrderIndex = plan.EvaluationOrderIndex.ValueInt64Pointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		addRequest.SubjectClaimName = plan.SubjectClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Add optional fields to create request for third-party access-token-validator
func addOptionalThirdPartyAccessTokenValidatorFields(ctx context.Context, addRequest *client.AddThirdPartyAccessTokenValidatorRequest, plan accessTokenValidatorResourceModel) error {
	if internaltypes.IsDefined(plan.ExtensionArgument) {
		var slice []string
		plan.ExtensionArgument.ElementsAs(ctx, &slice, false)
		addRequest.ExtensionArgument = slice
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.IdentityMapper) {
		addRequest.IdentityMapper = plan.IdentityMapper.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.SubjectClaimName) {
		addRequest.SubjectClaimName = plan.SubjectClaimName.ValueStringPointer()
	}
	// Empty strings are treated as equivalent to null
	if internaltypes.IsNonEmptyString(plan.Description) {
		addRequest.Description = plan.Description.ValueStringPointer()
	}
	return nil
}

// Populate any unknown values or sets that have a nil ElementType, to avoid errors when setting the state
func populateAccessTokenValidatorUnknownValues(ctx context.Context, model *accessTokenValidatorResourceModel) {
	if model.AllowedKeyEncryptionAlgorithm.ElementType(ctx) == nil {
		model.AllowedKeyEncryptionAlgorithm = types.SetNull(types.StringType)
	}
	if model.SigningCertificate.ElementType(ctx) == nil {
		model.SigningCertificate = types.SetNull(types.StringType)
	}
	if model.AllowedSigningAlgorithm.ElementType(ctx) == nil {
		model.AllowedSigningAlgorithm = types.SetNull(types.StringType)
	}
	if model.ExtensionArgument.ElementType(ctx) == nil {
		model.ExtensionArgument = types.SetNull(types.StringType)
	}
	if model.AllowedContentEncryptionAlgorithm.ElementType(ctx) == nil {
		model.AllowedContentEncryptionAlgorithm = types.SetNull(types.StringType)
	}
	if model.ClientSecret.IsUnknown() {
		model.ClientSecret = types.StringNull()
	}
}

// Read a PingFederateAccessTokenValidatorResponse object into the model struct
func readPingFederateAccessTokenValidatorResponse(ctx context.Context, r *client.PingFederateAccessTokenValidatorResponse, state *accessTokenValidatorResourceModel, expectedValues *accessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-federate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ClientID = types.StringValue(r.ClientID)
	state.ClientSecretPassphraseProvider = internaltypes.StringTypeOrNil(r.ClientSecretPassphraseProvider, internaltypes.IsEmptyString(expectedValues.ClientSecretPassphraseProvider))
	state.IncludeAudParameter = internaltypes.BoolTypeOrNil(r.IncludeAudParameter)
	state.AccessTokenManagerID = internaltypes.StringTypeOrNil(r.AccessTokenManagerID, internaltypes.IsEmptyString(expectedValues.AccessTokenManagerID))
	state.EndpointCacheRefresh = internaltypes.StringTypeOrNil(r.EndpointCacheRefresh, internaltypes.IsEmptyString(expectedValues.EndpointCacheRefresh))
	config.CheckMismatchedPDFormattedAttributes("endpoint_cache_refresh",
		expectedValues.EndpointCacheRefresh, state.EndpointCacheRefresh, diagnostics)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.AuthorizationServer = internaltypes.StringTypeOrNil(r.AuthorizationServer, internaltypes.IsEmptyString(expectedValues.AuthorizationServer))
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccessTokenValidatorUnknownValues(ctx, state)
}

// Read a JwtAccessTokenValidatorResponse object into the model struct
func readJwtAccessTokenValidatorResponse(ctx context.Context, r *client.JwtAccessTokenValidatorResponse, state *accessTokenValidatorResourceModel, expectedValues *accessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("jwt")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowedSigningAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedSigningAlgorithmProp(r.AllowedSigningAlgorithm))
	state.SigningCertificate = internaltypes.GetStringSet(r.SigningCertificate)
	state.JwksEndpointPath = internaltypes.StringTypeOrNil(r.JwksEndpointPath, internaltypes.IsEmptyString(expectedValues.JwksEndpointPath))
	state.EncryptionKeyPair = internaltypes.StringTypeOrNil(r.EncryptionKeyPair, internaltypes.IsEmptyString(expectedValues.EncryptionKeyPair))
	state.AllowedKeyEncryptionAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedKeyEncryptionAlgorithmProp(r.AllowedKeyEncryptionAlgorithm))
	state.AllowedContentEncryptionAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumaccessTokenValidatorAllowedContentEncryptionAlgorithmProp(r.AllowedContentEncryptionAlgorithm))
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, internaltypes.IsEmptyString(expectedValues.ClockSkewGracePeriod))
	config.CheckMismatchedPDFormattedAttributes("clock_skew_grace_period",
		expectedValues.ClockSkewGracePeriod, state.ClockSkewGracePeriod, diagnostics)
	state.ClientIDClaimName = internaltypes.StringTypeOrNil(r.ClientIDClaimName, internaltypes.IsEmptyString(expectedValues.ClientIDClaimName))
	state.ScopeClaimName = internaltypes.StringTypeOrNil(r.ScopeClaimName, internaltypes.IsEmptyString(expectedValues.ScopeClaimName))
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.AuthorizationServer = internaltypes.StringTypeOrNil(r.AuthorizationServer, internaltypes.IsEmptyString(expectedValues.AuthorizationServer))
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccessTokenValidatorUnknownValues(ctx, state)
}

// Read a MockAccessTokenValidatorResponse object into the model struct
func readMockAccessTokenValidatorResponse(ctx context.Context, r *client.MockAccessTokenValidatorResponse, state *accessTokenValidatorResourceModel, expectedValues *accessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("mock")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ClientIDClaimName = internaltypes.StringTypeOrNil(r.ClientIDClaimName, internaltypes.IsEmptyString(expectedValues.ClientIDClaimName))
	state.ScopeClaimName = internaltypes.StringTypeOrNil(r.ScopeClaimName, internaltypes.IsEmptyString(expectedValues.ScopeClaimName))
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccessTokenValidatorUnknownValues(ctx, state)
}

// Read a ThirdPartyAccessTokenValidatorResponse object into the model struct
func readThirdPartyAccessTokenValidatorResponse(ctx context.Context, r *client.ThirdPartyAccessTokenValidatorResponse, state *accessTokenValidatorResourceModel, expectedValues *accessTokenValidatorResourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.IdentityMapper = internaltypes.StringTypeOrNil(r.IdentityMapper, internaltypes.IsEmptyString(expectedValues.IdentityMapper))
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, internaltypes.IsEmptyString(expectedValues.SubjectClaimName))
	state.Description = internaltypes.StringTypeOrNil(r.Description, internaltypes.IsEmptyString(expectedValues.Description))
	state.Enabled = types.BoolValue(r.Enabled)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
	state.Notifications, state.RequiredActions = config.ReadMessages(ctx, r.Urnpingidentityschemasconfigurationmessages20, diagnostics)
	populateAccessTokenValidatorUnknownValues(ctx, state)
}

// Set any properties that aren't returned by the API in the state, based on some expected value (usually the plan value)
// This will include any parent endpoint names and any obscured (sensitive) attributes
func (state *accessTokenValidatorResourceModel) setStateValuesNotReturnedByAPI(expectedValues *accessTokenValidatorResourceModel) {
	if !expectedValues.ClientSecret.IsUnknown() {
		state.ClientSecret = expectedValues.ClientSecret
	}
}

// Create any update operations necessary to make the state match the plan
func createAccessTokenValidatorOperations(plan accessTokenValidatorResourceModel, state accessTokenValidatorResourceModel) []client.Operation {
	var ops []client.Operation
	operations.AddStringOperationIfNecessary(&ops, plan.ExtensionClass, state.ExtensionClass, "extension-class")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.ExtensionArgument, state.ExtensionArgument, "extension-argument")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedSigningAlgorithm, state.AllowedSigningAlgorithm, "allowed-signing-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.SigningCertificate, state.SigningCertificate, "signing-certificate")
	operations.AddStringOperationIfNecessary(&ops, plan.JwksEndpointPath, state.JwksEndpointPath, "jwks-endpoint-path")
	operations.AddStringOperationIfNecessary(&ops, plan.EncryptionKeyPair, state.EncryptionKeyPair, "encryption-key-pair")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedKeyEncryptionAlgorithm, state.AllowedKeyEncryptionAlgorithm, "allowed-key-encryption-algorithm")
	operations.AddStringSetOperationsIfNecessary(&ops, plan.AllowedContentEncryptionAlgorithm, state.AllowedContentEncryptionAlgorithm, "allowed-content-encryption-algorithm")
	operations.AddStringOperationIfNecessary(&ops, plan.ClockSkewGracePeriod, state.ClockSkewGracePeriod, "clock-skew-grace-period")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientIDClaimName, state.ClientIDClaimName, "client-id-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ScopeClaimName, state.ScopeClaimName, "scope-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientID, state.ClientID, "client-id")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientSecret, state.ClientSecret, "client-secret")
	operations.AddStringOperationIfNecessary(&ops, plan.ClientSecretPassphraseProvider, state.ClientSecretPassphraseProvider, "client-secret-passphrase-provider")
	operations.AddBoolOperationIfNecessary(&ops, plan.IncludeAudParameter, state.IncludeAudParameter, "include-aud-parameter")
	operations.AddStringOperationIfNecessary(&ops, plan.AccessTokenManagerID, state.AccessTokenManagerID, "access-token-manager-id")
	operations.AddStringOperationIfNecessary(&ops, plan.EndpointCacheRefresh, state.EndpointCacheRefresh, "endpoint-cache-refresh")
	operations.AddInt64OperationIfNecessary(&ops, plan.EvaluationOrderIndex, state.EvaluationOrderIndex, "evaluation-order-index")
	operations.AddStringOperationIfNecessary(&ops, plan.AuthorizationServer, state.AuthorizationServer, "authorization-server")
	operations.AddStringOperationIfNecessary(&ops, plan.IdentityMapper, state.IdentityMapper, "identity-mapper")
	operations.AddStringOperationIfNecessary(&ops, plan.SubjectClaimName, state.SubjectClaimName, "subject-claim-name")
	operations.AddStringOperationIfNecessary(&ops, plan.Description, state.Description, "description")
	operations.AddBoolOperationIfNecessary(&ops, plan.Enabled, state.Enabled, "enabled")
	return ops
}

// Create a ping-federate access-token-validator
func (r *accessTokenValidatorResource) CreatePingFederateAccessTokenValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accessTokenValidatorResourceModel) (*accessTokenValidatorResourceModel, error) {
	addRequest := client.NewAddPingFederateAccessTokenValidatorRequest(plan.Name.ValueString(),
		[]client.EnumpingFederateAccessTokenValidatorSchemaUrn{client.ENUMPINGFEDERATEACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORPING_FEDERATE},
		plan.ClientID.ValueString(),
		plan.Enabled.ValueBool())
	err := addOptionalPingFederateAccessTokenValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Access Token Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddPingFederateAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Access Token Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accessTokenValidatorResourceModel
	readPingFederateAccessTokenValidatorResponse(ctx, addResponse.PingFederateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a jwt access-token-validator
func (r *accessTokenValidatorResource) CreateJwtAccessTokenValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accessTokenValidatorResourceModel) (*accessTokenValidatorResourceModel, error) {
	addRequest := client.NewAddJwtAccessTokenValidatorRequest(plan.Name.ValueString(),
		[]client.EnumjwtAccessTokenValidatorSchemaUrn{client.ENUMJWTACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORJWT},
		plan.Enabled.ValueBool())
	err := addOptionalJwtAccessTokenValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Access Token Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddJwtAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Access Token Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accessTokenValidatorResourceModel
	readJwtAccessTokenValidatorResponse(ctx, addResponse.JwtAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a mock access-token-validator
func (r *accessTokenValidatorResource) CreateMockAccessTokenValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accessTokenValidatorResourceModel) (*accessTokenValidatorResourceModel, error) {
	addRequest := client.NewAddMockAccessTokenValidatorRequest(plan.Name.ValueString(),
		[]client.EnummockAccessTokenValidatorSchemaUrn{client.ENUMMOCKACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORMOCK},
		plan.Enabled.ValueBool())
	err := addOptionalMockAccessTokenValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Access Token Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddMockAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Access Token Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accessTokenValidatorResourceModel
	readMockAccessTokenValidatorResponse(ctx, addResponse.MockAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a third-party access-token-validator
func (r *accessTokenValidatorResource) CreateThirdPartyAccessTokenValidator(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan accessTokenValidatorResourceModel) (*accessTokenValidatorResourceModel, error) {
	addRequest := client.NewAddThirdPartyAccessTokenValidatorRequest(plan.Name.ValueString(),
		[]client.EnumthirdPartyAccessTokenValidatorSchemaUrn{client.ENUMTHIRDPARTYACCESSTOKENVALIDATORSCHEMAURN_URNPINGIDENTITYSCHEMASCONFIGURATION2_0ACCESS_TOKEN_VALIDATORTHIRD_PARTY},
		plan.ExtensionClass.ValueString(),
		plan.Enabled.ValueBool(),
		plan.EvaluationOrderIndex.ValueInt64())
	err := addOptionalThirdPartyAccessTokenValidatorFields(ctx, addRequest, plan)
	if err != nil {
		resp.Diagnostics.AddError("Failed to add optional properties to add request for Access Token Validator", err.Error())
		return nil, err
	}
	// Log request JSON
	requestJson, err := addRequest.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add request: "+string(requestJson))
	}
	apiAddRequest := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig))
	apiAddRequest = apiAddRequest.AddAccessTokenValidatorRequest(
		client.AddThirdPartyAccessTokenValidatorRequestAsAddAccessTokenValidatorRequest(addRequest))

	addResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.AddAccessTokenValidatorExecute(apiAddRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while creating the Access Token Validator", err, httpResp)
		return nil, err
	}

	// Log response JSON
	responseJson, err := addResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Add response: "+string(responseJson))
	}

	// Read the response into the state
	var state accessTokenValidatorResourceModel
	readThirdPartyAccessTokenValidatorResponse(ctx, addResponse.ThirdPartyAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
	return &state, nil
}

// Create a new resource
func (r *accessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state *accessTokenValidatorResourceModel
	var err error
	if plan.Type.ValueString() == "ping-federate" {
		state, err = r.CreatePingFederateAccessTokenValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "jwt" {
		state, err = r.CreateJwtAccessTokenValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "mock" {
		state, err = r.CreateMockAccessTokenValidator(ctx, req, resp, plan)
		if err != nil {
			return
		}
	}
	if plan.Type.ValueString() == "third-party" {
		state, err = r.CreateThirdPartyAccessTokenValidator(ctx, req, resp, plan)
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
func (r *defaultAccessTokenValidatorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan accessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the existing configuration
	var state accessTokenValidatorResourceModel
	if readResponse.PingFederateAccessTokenValidatorResponse != nil {
		readPingFederateAccessTokenValidatorResponse(ctx, readResponse.PingFederateAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JwtAccessTokenValidatorResponse != nil {
		readJwtAccessTokenValidatorResponse(ctx, readResponse.JwtAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MockAccessTokenValidatorResponse != nil {
		readMockAccessTokenValidatorResponse(ctx, readResponse.MockAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccessTokenValidatorResponse != nil {
		readThirdPartyAccessTokenValidatorResponse(ctx, readResponse.ThirdPartyAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}

	// Determine what changes are needed to match the plan
	updateRequest := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(config.ProviderBasicAuthContext(ctx, r.providerConfig), plan.Name.ValueString())
	ops := createAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := r.apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.PingFederateAccessTokenValidatorResponse != nil {
			readPingFederateAccessTokenValidatorResponse(ctx, updateResponse.PingFederateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JwtAccessTokenValidatorResponse != nil {
			readJwtAccessTokenValidatorResponse(ctx, updateResponse.JwtAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MockAccessTokenValidatorResponse != nil {
			readMockAccessTokenValidatorResponse(ctx, updateResponse.MockAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAccessTokenValidatorResponse != nil {
			readThirdPartyAccessTokenValidatorResponse(ctx, updateResponse.ThirdPartyAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *accessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAccessTokenValidatorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	readAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func readAccessTokenValidator(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Get current state
	var state accessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := apiClient.AccessTokenValidatorApi.GetAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Access Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PingFederateAccessTokenValidatorResponse != nil {
		readPingFederateAccessTokenValidatorResponse(ctx, readResponse.PingFederateAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.JwtAccessTokenValidatorResponse != nil {
		readJwtAccessTokenValidatorResponse(ctx, readResponse.JwtAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.MockAccessTokenValidatorResponse != nil {
		readMockAccessTokenValidatorResponse(ctx, readResponse.MockAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyAccessTokenValidatorResponse != nil {
		readThirdPartyAccessTokenValidatorResponse(ctx, readResponse.ThirdPartyAccessTokenValidatorResponse, &state, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update a resource
func (r *accessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func (r *defaultAccessTokenValidatorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	updateAccessTokenValidator(ctx, req, resp, r.apiClient, r.providerConfig)
}

func updateAccessTokenValidator(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, apiClient *client.APIClient, providerConfig internaltypes.ProviderConfiguration) {
	// Retrieve values from plan
	var plan accessTokenValidatorResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the current state to see how any attributes are changing
	var state accessTokenValidatorResourceModel
	req.State.Get(ctx, &state)
	updateRequest := apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, providerConfig), plan.Name.ValueString())

	// Determine what update operations are necessary
	ops := createAccessTokenValidatorOperations(plan, state)
	if len(ops) > 0 {
		updateRequest = updateRequest.UpdateRequest(*client.NewUpdateRequest(ops))
		// Log operations
		operations.LogUpdateOperations(ctx, ops)

		updateResponse, httpResp, err := apiClient.AccessTokenValidatorApi.UpdateAccessTokenValidatorExecute(updateRequest)
		if err != nil {
			config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while updating the Access Token Validator", err, httpResp)
			return
		}

		// Log response JSON
		responseJson, err := updateResponse.MarshalJSON()
		if err == nil {
			tflog.Debug(ctx, "Update response: "+string(responseJson))
		}

		// Read the response
		if updateResponse.PingFederateAccessTokenValidatorResponse != nil {
			readPingFederateAccessTokenValidatorResponse(ctx, updateResponse.PingFederateAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.JwtAccessTokenValidatorResponse != nil {
			readJwtAccessTokenValidatorResponse(ctx, updateResponse.JwtAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.MockAccessTokenValidatorResponse != nil {
			readMockAccessTokenValidatorResponse(ctx, updateResponse.MockAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
		}
		if updateResponse.ThirdPartyAccessTokenValidatorResponse != nil {
			readThirdPartyAccessTokenValidatorResponse(ctx, updateResponse.ThirdPartyAccessTokenValidatorResponse, &state, &plan, &resp.Diagnostics)
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
func (r *defaultAccessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No implementation necessary
}

func (r *accessTokenValidatorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state accessTokenValidatorResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidatorExecute(r.apiClient.AccessTokenValidatorApi.DeleteAccessTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()))
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while deleting the Access Token Validator", err, httpResp)
		return
	}
}

func (r *accessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAccessTokenValidator(ctx, req, resp)
}

func (r *defaultAccessTokenValidatorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	importAccessTokenValidator(ctx, req, resp)
}

func importAccessTokenValidator(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to name attribute
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
