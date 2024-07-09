package ldapcorrelationattributepair

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
	_ datasource.DataSource              = &ldapCorrelationAttributePairDataSource{}
	_ datasource.DataSourceWithConfigure = &ldapCorrelationAttributePairDataSource{}
)

// Create a Ldap Correlation Attribute Pair data source
func NewLdapCorrelationAttributePairDataSource() datasource.DataSource {
	return &ldapCorrelationAttributePairDataSource{}
}

// ldapCorrelationAttributePairDataSource is the datasource implementation.
type ldapCorrelationAttributePairDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *ldapCorrelationAttributePairDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_correlation_attribute_pair"
}

// Configure adds the provider configured client to the data source.
func (r *ldapCorrelationAttributePairDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type ldapCorrelationAttributePairDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Name                          types.String `tfsdk:"name"`
	Type                          types.String `tfsdk:"type"`
	CorrelatedLdapDataViewName    types.String `tfsdk:"correlated_ldap_data_view_name"`
	ScimResourceTypeName          types.String `tfsdk:"scim_resource_type_name"`
	PrimaryCorrelationAttribute   types.String `tfsdk:"primary_correlation_attribute"`
	SecondaryCorrelationAttribute types.String `tfsdk:"secondary_correlation_attribute"`
}

// GetSchema defines the schema for the datasource.
func (r *ldapCorrelationAttributePairDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Ldap Correlation Attribute Pair.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of LDAP Correlation Attribute Pair resource. Options are ['ldap-correlation-attribute-pair']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"correlated_ldap_data_view_name": schema.StringAttribute{
				Description: "Name of the parent Correlated LDAP Data View",
				Required:    true,
			},
			"scim_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Resource Type",
				Required:    true,
			},
			"primary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the base SCIM Resource Type whose value will be used to match objects in the Correlated LDAP Data View.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secondary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the Correlated LDAP Data View whose value will be matched.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a LdapCorrelationAttributePairResponse object into the model struct
func readLdapCorrelationAttributePairResponseDataSource(ctx context.Context, r *client.LdapCorrelationAttributePairResponse, state *ldapCorrelationAttributePairDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-correlation-attribute-pair")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PrimaryCorrelationAttribute = types.StringValue(r.PrimaryCorrelationAttribute)
	state.SecondaryCorrelationAttribute = types.StringValue(r.SecondaryCorrelationAttribute)
}

// Read resource information
func (r *ldapCorrelationAttributePairDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state ldapCorrelationAttributePairDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LdapCorrelationAttributePairAPI.GetLdapCorrelationAttributePair(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.CorrelatedLdapDataViewName.ValueString(), state.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Ldap Correlation Attribute Pair", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readLdapCorrelationAttributePairResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
