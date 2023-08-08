package consentservice

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
	_ datasource.DataSource              = &consentServiceDataSource{}
	_ datasource.DataSourceWithConfigure = &consentServiceDataSource{}
)

// Create a Consent Service data source
func NewConsentServiceDataSource() datasource.DataSource {
	return &consentServiceDataSource{}
}

// consentServiceDataSource is the datasource implementation.
type consentServiceDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *consentServiceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_consent_service"
}

// Configure adds the provider configured client to the data source.
func (r *consentServiceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type consentServiceDataSourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Type                        types.String `tfsdk:"type"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	BaseDN                      types.String `tfsdk:"base_dn"`
	BindDN                      types.String `tfsdk:"bind_dn"`
	SearchSizeLimit             types.Int64  `tfsdk:"search_size_limit"`
	ConsentRecordIdentityMapper types.Set    `tfsdk:"consent_record_identity_mapper"`
	ServiceAccountDN            types.Set    `tfsdk:"service_account_dn"`
	UnprivilegedConsentScope    types.String `tfsdk:"unprivileged_consent_scope"`
	PrivilegedConsentScope      types.String `tfsdk:"privileged_consent_scope"`
	Audience                    types.String `tfsdk:"audience"`
}

// GetSchema defines the schema for the datasource.
func (r *consentServiceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Consent Service.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Consent Service resource. Options are ['consent-service']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Consent Service is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"base_dn": schema.StringAttribute{
				Description: "The base DN under which consent records are stored.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"bind_dn": schema.StringAttribute{
				Description: "The DN of an internal service account used by the Consent Service to make internal LDAP requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_size_limit": schema.Int64Attribute{
				Description: "The maximum number of consent resources that may be returned from a search request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"consent_record_identity_mapper": schema.SetAttribute{
				Description: "If specified, the Identity Mapper(s) that may be used to map consent record subject and actor values to DNs. This is typically only needed if privileged API clients will be used.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"service_account_dn": schema.SetAttribute{
				Description: "The set of account DNs that the Consent Service will consider to be privileged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"unprivileged_consent_scope": schema.StringAttribute{
				Description: "The name of a scope that must be present in an access token accepted by the Consent Service for unprivileged clients.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"privileged_consent_scope": schema.StringAttribute{
				Description: "The name of a scope that must be present in an access token accepted by the Consent Service if the client is to be considered privileged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"audience": schema.StringAttribute{
				Description: "A string or URI that identifies the Consent Service in the context of OAuth2 authorization.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a ConsentServiceResponse object into the model struct
func readConsentServiceResponseDataSource(ctx context.Context, r *client.ConsentServiceResponse, state *consentServiceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("consent-service")
	// Placeholder id value required by test framework
	state.Id = types.StringValue("id")
	state.Enabled = types.BoolValue(r.Enabled)
	state.BaseDN = internaltypes.StringTypeOrNil(r.BaseDN, false)
	state.BindDN = internaltypes.StringTypeOrNil(r.BindDN, false)
	state.SearchSizeLimit = internaltypes.Int64TypeOrNil(r.SearchSizeLimit)
	state.ConsentRecordIdentityMapper = internaltypes.GetStringSet(r.ConsentRecordIdentityMapper)
	state.ServiceAccountDN = internaltypes.GetStringSet(r.ServiceAccountDN)
	state.UnprivilegedConsentScope = internaltypes.StringTypeOrNil(r.UnprivilegedConsentScope, false)
	state.PrivilegedConsentScope = internaltypes.StringTypeOrNil(r.PrivilegedConsentScope, false)
	state.Audience = internaltypes.StringTypeOrNil(r.Audience, false)
}

// Read resource information
func (r *consentServiceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state consentServiceDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConsentServiceApi.GetConsentService(
		config.ProviderBasicAuthContext(ctx, r.providerConfig)).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Consent Service", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConsentServiceResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
