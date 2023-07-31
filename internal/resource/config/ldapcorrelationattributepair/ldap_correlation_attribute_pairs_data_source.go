package ldapcorrelationattributepair

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v9300/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &ldapCorrelationAttributePairsDataSource{}
	_ datasource.DataSourceWithConfigure = &ldapCorrelationAttributePairsDataSource{}
)

// Create a Ldap Correlation Attribute Pairs data source
func NewLdapCorrelationAttributePairsDataSource() datasource.DataSource {
	return &ldapCorrelationAttributePairsDataSource{}
}

// ldapCorrelationAttributePairsDataSource is the datasource implementation.
type ldapCorrelationAttributePairsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *ldapCorrelationAttributePairsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ldap_correlation_attribute_pairs"
}

// Configure adds the provider configured client to the data source.
func (r *ldapCorrelationAttributePairsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type ldapCorrelationAttributePairsDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Filter                     types.String `tfsdk:"filter"`
	Ids                        types.Set    `tfsdk:"ids"`
	CorrelatedLdapDataViewName types.String `tfsdk:"correlated_ldap_data_view_name"`
	ScimResourceTypeName       types.String `tfsdk:"scim_resource_type_name"`
}

// GetSchema defines the schema for the datasource.
func (r *ldapCorrelationAttributePairsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Lists Ldap Correlation Attribute Pair objects in the server configuration.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Placeholder name of this object required by Terraform.",
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
			"filter": schema.StringAttribute{
				Description: "SCIM filter used when searching the configuration.",
				Optional:    true,
			},
			"ids": schema.SetAttribute{
				Description: "Ldap Correlation Attribute Pair IDs found in the configuration",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Read resource information
func (r *ldapCorrelationAttributePairsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state ldapCorrelationAttributePairsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	listRequest := r.apiClient.LdapCorrelationAttributePairApi.ListLdapCorrelationAttributePairs(config.ProviderBasicAuthContext(ctx, r.providerConfig), state.CorrelatedLdapDataViewName.ValueString(), state.ScimResourceTypeName.ValueString())
	if internaltypes.IsDefined(state.Filter) {
		listRequest = listRequest.Filter(state.Filter.ValueString())
	}

	readResponse, httpResp, err := r.apiClient.LdapCorrelationAttributePairApi.ListLdapCorrelationAttributePairsExecute(listRequest)
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while listing the Ldap Correlation Attribute Pair objects", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	ids := []attr.Value{}
	for _, response := range readResponse.Resources {
		ids = append(ids, types.StringValue(response.Id))
	}

	state.Ids, diags = types.SetValue(types.StringType, ids)
	resp.Diagnostics.Append(diags...)
	state.Id = types.StringValue("id")

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
