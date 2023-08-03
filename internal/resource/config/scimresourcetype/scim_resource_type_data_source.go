package scimresourcetype

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
	_ datasource.DataSource              = &scimResourceTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &scimResourceTypeDataSource{}
)

// Create a Scim Resource Type data source
func NewScimResourceTypeDataSource() datasource.DataSource {
	return &scimResourceTypeDataSource{}
}

// scimResourceTypeDataSource is the datasource implementation.
type scimResourceTypeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *scimResourceTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_resource_type"
}

// Configure adds the provider configured client to the data source.
func (r *scimResourceTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type scimResourceTypeDataSourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Name                        types.String `tfsdk:"name"`
	Type                        types.String `tfsdk:"type"`
	CoreSchema                  types.String `tfsdk:"core_schema"`
	RequiredSchemaExtension     types.Set    `tfsdk:"required_schema_extension"`
	OptionalSchemaExtension     types.Set    `tfsdk:"optional_schema_extension"`
	Description                 types.String `tfsdk:"description"`
	Enabled                     types.Bool   `tfsdk:"enabled"`
	Endpoint                    types.String `tfsdk:"endpoint"`
	LookthroughLimit            types.Int64  `tfsdk:"lookthrough_limit"`
	SchemaCheckingOption        types.Set    `tfsdk:"schema_checking_option"`
	StructuralLDAPObjectclass   types.String `tfsdk:"structural_ldap_objectclass"`
	AuxiliaryLDAPObjectclass    types.Set    `tfsdk:"auxiliary_ldap_objectclass"`
	IncludeBaseDN               types.String `tfsdk:"include_base_dn"`
	IncludeFilter               types.Set    `tfsdk:"include_filter"`
	IncludeOperationalAttribute types.Set    `tfsdk:"include_operational_attribute"`
	CreateDNPattern             types.String `tfsdk:"create_dn_pattern"`
}

// GetSchema defines the schema for the datasource.
func (r *scimResourceTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Scim Resource Type.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of SCIM Resource Type resource. Options are ['ldap-pass-through', 'ldap-mapping']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"core_schema": schema.StringAttribute{
				Description: "The core schema enforced on core attributes at the top level of a SCIM resource representation exposed by thisMapping SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"required_schema_extension": schema.SetAttribute{
				Description: "Required additive schemas that are enforced on extension attributes in a SCIM resource representation for this Mapping SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"optional_schema_extension": schema.SetAttribute{
				Description: "Optional additive schemas that are enforced on extension attributes in a SCIM resource representation for this Mapping SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this SCIM Resource Type",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the SCIM Resource Type is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"endpoint": schema.StringAttribute{
				Description: "The HTTP addressable endpoint of this SCIM Resource Type relative to the '/scim/v2' base URL. Do not include a leading '/'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"lookthrough_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that the SCIM Resource Type should \"look through\" in the course of processing a search request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"schema_checking_option": schema.SetAttribute{
				Description: "Options to alter the way schema checking is performed during create or modify requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory that can be accessed by this SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this SCIM Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_operational_attribute": schema.SetAttribute{
				Description: "Specifies the set of operational LDAP attributes to be provided by this SCIM Resource Type.",
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
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a LdapPassThroughScimResourceTypeResponse object into the model struct
func readLdapPassThroughScimResourceTypeResponseDataSource(ctx context.Context, r *client.LdapPassThroughScimResourceTypeResponse, state *scimResourceTypeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-pass-through")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Endpoint = types.StringValue(r.Endpoint)
	state.LookthroughLimit = internaltypes.Int64TypeOrNil(r.LookthroughLimit)
	state.SchemaCheckingOption = internaltypes.GetStringSet(
		client.StringSliceEnumscimResourceTypeSchemaCheckingOptionProp(r.SchemaCheckingOption))
	state.StructuralLDAPObjectclass = internaltypes.StringTypeOrNil(r.StructuralLDAPObjectclass, false)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.IncludeBaseDN = internaltypes.StringTypeOrNil(r.IncludeBaseDN, false)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.IncludeOperationalAttribute = internaltypes.GetStringSet(r.IncludeOperationalAttribute)
	state.CreateDNPattern = internaltypes.StringTypeOrNil(r.CreateDNPattern, false)
}

// Read a LdapMappingScimResourceTypeResponse object into the model struct
func readLdapMappingScimResourceTypeResponseDataSource(ctx context.Context, r *client.LdapMappingScimResourceTypeResponse, state *scimResourceTypeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("ldap-mapping")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.CoreSchema = types.StringValue(r.CoreSchema)
	state.RequiredSchemaExtension = internaltypes.GetStringSet(r.RequiredSchemaExtension)
	state.OptionalSchemaExtension = internaltypes.GetStringSet(r.OptionalSchemaExtension)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.Endpoint = types.StringValue(r.Endpoint)
	state.LookthroughLimit = internaltypes.Int64TypeOrNil(r.LookthroughLimit)
	state.SchemaCheckingOption = internaltypes.GetStringSet(
		client.StringSliceEnumscimResourceTypeSchemaCheckingOptionProp(r.SchemaCheckingOption))
	state.StructuralLDAPObjectclass = internaltypes.StringTypeOrNil(r.StructuralLDAPObjectclass, false)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.IncludeBaseDN = internaltypes.StringTypeOrNil(r.IncludeBaseDN, false)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.IncludeOperationalAttribute = internaltypes.GetStringSet(r.IncludeOperationalAttribute)
	state.CreateDNPattern = internaltypes.StringTypeOrNil(r.CreateDNPattern, false)
}

// Read resource information
func (r *scimResourceTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state scimResourceTypeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimResourceTypeApi.GetScimResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.LdapPassThroughScimResourceTypeResponse != nil {
		readLdapPassThroughScimResourceTypeResponseDataSource(ctx, readResponse.LdapPassThroughScimResourceTypeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.LdapMappingScimResourceTypeResponse != nil {
		readLdapMappingScimResourceTypeResponseDataSource(ctx, readResponse.LdapMappingScimResourceTypeResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
