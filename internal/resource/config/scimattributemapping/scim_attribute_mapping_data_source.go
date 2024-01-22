package scimattributemapping

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10000/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &scimAttributeMappingDataSource{}
	_ datasource.DataSourceWithConfigure = &scimAttributeMappingDataSource{}
)

// Create a Scim Attribute Mapping data source
func NewScimAttributeMappingDataSource() datasource.DataSource {
	return &scimAttributeMappingDataSource{}
}

// scimAttributeMappingDataSource is the datasource implementation.
type scimAttributeMappingDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *scimAttributeMappingDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_attribute_mapping"
}

// Configure adds the provider configured client to the data source.
func (r *scimAttributeMappingDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type scimAttributeMappingDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Type                      types.String `tfsdk:"type"`
	ScimResourceTypeName      types.String `tfsdk:"scim_resource_type_name"`
	CorrelatedLDAPDataView    types.String `tfsdk:"correlated_ldap_data_view"`
	ScimResourceTypeAttribute types.String `tfsdk:"scim_resource_type_attribute"`
	LdapAttribute             types.String `tfsdk:"ldap_attribute"`
	Readable                  types.Bool   `tfsdk:"readable"`
	Writable                  types.Bool   `tfsdk:"writable"`
	Searchable                types.Bool   `tfsdk:"searchable"`
	Authoritative             types.Bool   `tfsdk:"authoritative"`
}

// GetSchema defines the schema for the datasource.
func (r *scimAttributeMappingDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Scim Attribute Mapping.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of SCIM Attribute Mapping resource. Options are ['scim-attribute-mapping']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scim_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Resource Type",
				Required:    true,
			},
			"correlated_ldap_data_view": schema.StringAttribute{
				Description: "The Correlated LDAP Data View that persists the mapped SCIM Resource Type attribute(s).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scim_resource_type_attribute": schema.StringAttribute{
				Description: "The attribute path of SCIM Resource Type attributes to be mapped.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"ldap_attribute": schema.StringAttribute{
				Description: "The LDAP attribute to be mapped, or the path to a specific field of an LDAP attribute with the JSON object attribute syntax.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"readable": schema.BoolAttribute{
				Description: "Specifies whether the mapping is used to map from LDAP attribute to SCIM Resource Type attribute in a read operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"writable": schema.BoolAttribute{
				Description: "Specifies that the mapping is used to map from SCIM Resource Type attribute to LDAP attribute in a write operation.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"searchable": schema.BoolAttribute{
				Description: "Specifies that the mapping is used to map from SCIM Resource Type attribute to LDAP attribute in a search filter.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"authoritative": schema.BoolAttribute{
				Description: "Specifies that the mapping is authoritative over other mappings for the same SCIM Resource Type attribute (for read operations).",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ScimAttributeMappingResponse object into the model struct
func readScimAttributeMappingResponseDataSource(ctx context.Context, r *client.ScimAttributeMappingResponse, state *scimAttributeMappingDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("scim-attribute-mapping")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CorrelatedLDAPDataView = internaltypes.StringTypeOrNil(r.CorrelatedLDAPDataView, false)
	state.ScimResourceTypeAttribute = types.StringValue(r.ScimResourceTypeAttribute)
	state.LdapAttribute = types.StringValue(r.LdapAttribute)
	state.Readable = internaltypes.BoolTypeOrNil(r.Readable)
	state.Writable = internaltypes.BoolTypeOrNil(r.Writable)
	state.Searchable = internaltypes.BoolTypeOrNil(r.Searchable)
	state.Authoritative = internaltypes.BoolTypeOrNil(r.Authoritative)
}

// Read resource information
func (r *scimAttributeMappingDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state scimAttributeMappingDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimAttributeMappingAPI.GetScimAttributeMapping(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ScimResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Attribute Mapping", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScimAttributeMappingResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
