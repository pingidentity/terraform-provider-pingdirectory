package correlatedldapdataview

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
	_ datasource.DataSource              = &correlatedLdapDataViewDataSource{}
	_ datasource.DataSourceWithConfigure = &correlatedLdapDataViewDataSource{}
)

// Create a Correlated Ldap Data View data source
func NewCorrelatedLdapDataViewDataSource() datasource.DataSource {
	return &correlatedLdapDataViewDataSource{}
}

// correlatedLdapDataViewDataSource is the datasource implementation.
type correlatedLdapDataViewDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *correlatedLdapDataViewDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_correlated_ldap_data_view"
}

// Configure adds the provider configured client to the data source.
func (r *correlatedLdapDataViewDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type correlatedLdapDataViewDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	ScimResourceTypeName          types.String `tfsdk:"scim_resource_type_name"`
	StructuralLDAPObjectclass     types.String `tfsdk:"structural_ldap_objectclass"`
	AuxiliaryLDAPObjectclass      types.Set    `tfsdk:"auxiliary_ldap_objectclass"`
	IncludeBaseDN                 types.String `tfsdk:"include_base_dn"`
	IncludeFilter                 types.Set    `tfsdk:"include_filter"`
	IncludeOperationalAttribute   types.Set    `tfsdk:"include_operational_attribute"`
	CreateDNPattern               types.String `tfsdk:"create_dn_pattern"`
	PrimaryCorrelationAttribute   types.String `tfsdk:"primary_correlation_attribute"`
	SecondaryCorrelationAttribute types.String `tfsdk:"secondary_correlation_attribute"`
}

// GetSchema defines the schema for the datasource.
func (r *correlatedLdapDataViewDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Describes a Correlated Ldap Data View.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Name of this object.",
				Required:    true,
			},
			"scim_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Resource Type",
				Required:    true,
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this Correlated LDAP Data View.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this Correlated LDAP Data View.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory that can be accessed by this Correlated LDAP Data View.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this Correlated LDAP Data View.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_operational_attribute": schema.SetAttribute{
				Description: "Specifies the set of operational LDAP attributes to be provided by this Correlated LDAP Data View.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"create_dn_pattern": schema.StringAttribute{
				Description: "Specifies the template to use for the DN when creating new entries.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"primary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the parent SCIM Resource Type whose value will be used to match objects in the Correlated LDAP Data View. If multiple correlation attributes are required they may be created using additional correlation-attribute-pairs.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secondary_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the Correlated LDAP Data View whose value will be matched with the primary-correlation-attribute. If multiple correlation attributes are required they may be specified by creating additional correlation-attribute-pairs.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
}

// Read a CorrelatedLdapDataViewResponse object into the model struct
func readCorrelatedLdapDataViewResponseDataSource(ctx context.Context, r *client.CorrelatedLdapDataViewResponse, state *correlatedLdapDataViewDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.StructuralLDAPObjectclass = types.StringValue(r.StructuralLDAPObjectclass)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.IncludeBaseDN = types.StringValue(r.IncludeBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.IncludeOperationalAttribute = internaltypes.GetStringSet(r.IncludeOperationalAttribute)
	state.CreateDNPattern = internaltypes.StringTypeOrNil(r.CreateDNPattern, false)
	state.PrimaryCorrelationAttribute = types.StringValue(r.PrimaryCorrelationAttribute)
	state.SecondaryCorrelationAttribute = types.StringValue(r.SecondaryCorrelationAttribute)
}

// Read resource information
func (r *correlatedLdapDataViewDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state correlatedLdapDataViewDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.CorrelatedLdapDataViewApi.GetCorrelatedLdapDataView(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Id.ValueString(), state.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Correlated Ldap Data View", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readCorrelatedLdapDataViewResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
