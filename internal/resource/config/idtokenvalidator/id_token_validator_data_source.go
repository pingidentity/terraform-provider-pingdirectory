package idtokenvalidator

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
	_ datasource.DataSource              = &idTokenValidatorDataSource{}
	_ datasource.DataSourceWithConfigure = &idTokenValidatorDataSource{}
)

// Create a Id Token Validator data source
func NewIdTokenValidatorDataSource() datasource.DataSource {
	return &idTokenValidatorDataSource{}
}

// idTokenValidatorDataSource is the datasource implementation.
type idTokenValidatorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *idTokenValidatorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_id_token_validator"
}

// Configure adds the provider configured client to the data source.
func (r *idTokenValidatorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type idTokenValidatorDataSourceModel struct {
	Id                                 types.String `tfsdk:"id"`
	Name                               types.String `tfsdk:"name"`
	Type                               types.String `tfsdk:"type"`
	AllowedSigningAlgorithm            types.Set    `tfsdk:"allowed_signing_algorithm"`
	SigningCertificate                 types.Set    `tfsdk:"signing_certificate"`
	IssuerURL                          types.String `tfsdk:"issuer_url"`
	JwksEndpointPath                   types.String `tfsdk:"jwks_endpoint_path"`
	OpenIDConnectProvider              types.String `tfsdk:"openid_connect_provider"`
	OpenIDConnectMetadataCacheDuration types.String `tfsdk:"openid_connect_metadata_cache_duration"`
	Description                        types.String `tfsdk:"description"`
	Enabled                            types.Bool   `tfsdk:"enabled"`
	IdentityMapper                     types.String `tfsdk:"identity_mapper"`
	SubjectClaimName                   types.String `tfsdk:"subject_claim_name"`
	ClockSkewGracePeriod               types.String `tfsdk:"clock_skew_grace_period"`
	JwksCacheDuration                  types.String `tfsdk:"jwks_cache_duration"`
	EvaluationOrderIndex               types.Int64  `tfsdk:"evaluation_order_index"`
}

// GetSchema defines the schema for the datasource.
func (r *idTokenValidatorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Id Token Validator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of ID Token Validator resource. Options are ['ping-one', 'openid-connect']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_signing_algorithm": schema.SetAttribute{
				Description: "Specifies an allow list of JWT signing algorithms that will be accepted by the OpenID Connect ID Token Validator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"signing_certificate": schema.SetAttribute{
				Description: "Specifies the locally stored certificates that may be used to validate the signature of an incoming ID token. This property may be specified if a JWKS endpoint should not be used to retrieve public signing keys.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"issuer_url": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ping-one`: Specifies a PingOne base issuer URL. When the `type` attribute is set to `openid-connect`: Specifies the OpenID Connect provider's issuer URL.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ping-one`: Specifies a PingOne base issuer URL.\n  - `openid-connect`: Specifies the OpenID Connect provider's issuer URL.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"jwks_endpoint_path": schema.StringAttribute{
				Description: "The relative path to the JWKS endpoint from which to retrieve one or more public signing keys that may be used to validate the signature of an incoming ID token. This path is relative to the base_url property defined for the validator's OpenID Connect provider. If jwks-endpoint-path is specified, the OpenID Connect ID Token Validator will not consult locally stored certificates for validating token signatures.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"openid_connect_provider": schema.StringAttribute{
				Description:         "When the `type` attribute is set to `ping-one`: Specifies HTTPS connection settings for the PingOne OpenID Connect provider. When the `type` attribute is set to `openid-connect`: Specifies the OpenID Connect provider that issues ID tokens handled by this OpenID Connect ID Token Validator. This property is used in conjunction with the jwks-endpoint-path property.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `ping-one`: Specifies HTTPS connection settings for the PingOne OpenID Connect provider.\n  - `openid-connect`: Specifies the OpenID Connect provider that issues ID tokens handled by this OpenID Connect ID Token Validator. This property is used in conjunction with the jwks-endpoint-path property.",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"openid_connect_metadata_cache_duration": schema.StringAttribute{
				Description: "How often the PingOne ID Token Validator should refresh its stored cache of OpenID Connect-related metadata.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this ID Token Validator",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this ID Token Validator is enabled for use in the Directory Server.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"identity_mapper": schema.StringAttribute{
				Description: "Specifies the name of the Identity Mapper that should be used to correlate an ID token subject value to a user entry. The claim name from which to obtain the subject (i.e. the currently logged-in user) may be configured using the subject-claim-name property.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"subject_claim_name": schema.StringAttribute{
				Description: "The name of the token claim that contains the subject; i.e., the authenticated user.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"clock_skew_grace_period": schema.StringAttribute{
				Description: "Specifies the amount of clock skew that is tolerated by the ID Token Validator when evaluating whether a token is within its valid time interval. The duration specified by this parameter will be subtracted from the token's not-before (nbf) time and added to the token's expiration (exp) time, if present, to allow for any time difference between the local server's clock and the token issuer's clock.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"jwks_cache_duration": schema.StringAttribute{
				Description: "How often the ID Token Validator should refresh its cache of JWKS token signing keys.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"evaluation_order_index": schema.Int64Attribute{
				Description: "When multiple ID Token Validators are defined for a single Directory Server, this property determines the order in which the ID Token Validators are consulted. Values of this property must be unique among all ID Token Validators defined within Directory Server but not necessarily contiguous. ID Token Validators with lower values will be evaluated first to determine if they are able to validate the ID token.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a PingOneIdTokenValidatorResponse object into the model struct
func readPingOneIdTokenValidatorResponseDataSource(ctx context.Context, r *client.PingOneIdTokenValidatorResponse, state *idTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ping-one")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IssuerURL = types.StringValue(r.IssuerURL)
	state.OpenIDConnectProvider = types.StringValue(r.OpenIDConnectProvider)
	state.OpenIDConnectMetadataCacheDuration = internaltypes.StringTypeOrNil(r.OpenIDConnectMetadataCacheDuration, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, false)
	state.JwksCacheDuration = internaltypes.StringTypeOrNil(r.JwksCacheDuration, false)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
}

// Read a OpenidConnectIdTokenValidatorResponse object into the model struct
func readOpenidConnectIdTokenValidatorResponseDataSource(ctx context.Context, r *client.OpenidConnectIdTokenValidatorResponse, state *idTokenValidatorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("openid-connect")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllowedSigningAlgorithm = internaltypes.GetStringSet(
		client.StringSliceEnumidTokenValidatorAllowedSigningAlgorithmProp(r.AllowedSigningAlgorithm))
	state.SigningCertificate = internaltypes.GetStringSet(r.SigningCertificate)
	state.OpenIDConnectProvider = internaltypes.StringTypeOrNil(r.OpenIDConnectProvider, false)
	state.JwksEndpointPath = internaltypes.StringTypeOrNil(r.JwksEndpointPath, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.IdentityMapper = types.StringValue(r.IdentityMapper)
	state.SubjectClaimName = internaltypes.StringTypeOrNil(r.SubjectClaimName, false)
	state.IssuerURL = types.StringValue(r.IssuerURL)
	state.ClockSkewGracePeriod = internaltypes.StringTypeOrNil(r.ClockSkewGracePeriod, false)
	state.JwksCacheDuration = internaltypes.StringTypeOrNil(r.JwksCacheDuration, false)
	state.EvaluationOrderIndex = types.Int64Value(r.EvaluationOrderIndex)
}

// Read resource information
func (r *idTokenValidatorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state idTokenValidatorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdTokenValidatorApi.GetIdTokenValidator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Id Token Validator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.PingOneIdTokenValidatorResponse != nil {
		readPingOneIdTokenValidatorResponseDataSource(ctx, readResponse.PingOneIdTokenValidatorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.OpenidConnectIdTokenValidatorResponse != nil {
		readOpenidConnectIdTokenValidatorResponseDataSource(ctx, readResponse.OpenidConnectIdTokenValidatorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
