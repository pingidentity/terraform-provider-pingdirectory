package restresourcetype

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
	_ datasource.DataSource              = &restResourceTypeDataSource{}
	_ datasource.DataSourceWithConfigure = &restResourceTypeDataSource{}
)

// Create a Rest Resource Type data source
func NewRestResourceTypeDataSource() datasource.DataSource {
	return &restResourceTypeDataSource{}
}

// restResourceTypeDataSource is the datasource implementation.
type restResourceTypeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *restResourceTypeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rest_resource_type"
}

// Configure adds the provider configured client to the data source.
func (r *restResourceTypeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type restResourceTypeDataSourceModel struct {
	Id                             types.String `tfsdk:"id"`
	Name                           types.String `tfsdk:"name"`
	Type                           types.String `tfsdk:"type"`
	PasswordAttributeCategory      types.String `tfsdk:"password_attribute_category"`
	PasswordDisplayOrderIndex      types.Int64  `tfsdk:"password_display_order_index"`
	Description                    types.String `tfsdk:"description"`
	Enabled                        types.Bool   `tfsdk:"enabled"`
	ResourceEndpoint               types.String `tfsdk:"resource_endpoint"`
	StructuralLDAPObjectclass      types.String `tfsdk:"structural_ldap_objectclass"`
	AuxiliaryLDAPObjectclass       types.Set    `tfsdk:"auxiliary_ldap_objectclass"`
	SearchBaseDN                   types.String `tfsdk:"search_base_dn"`
	IncludeFilter                  types.Set    `tfsdk:"include_filter"`
	ParentDN                       types.String `tfsdk:"parent_dn"`
	ParentResourceType             types.String `tfsdk:"parent_resource_type"`
	RelativeDNFromParentResource   types.String `tfsdk:"relative_dn_from_parent_resource"`
	CreateRDNAttributeType         types.String `tfsdk:"create_rdn_attribute_type"`
	PostCreateConstructedAttribute types.Set    `tfsdk:"post_create_constructed_attribute"`
	UpdateConstructedAttribute     types.Set    `tfsdk:"update_constructed_attribute"`
	DisplayName                    types.String `tfsdk:"display_name"`
	SearchFilterPattern            types.String `tfsdk:"search_filter_pattern"`
	PrimaryDisplayAttributeType    types.String `tfsdk:"primary_display_attribute_type"`
	DelegatedAdminSearchSizeLimit  types.Int64  `tfsdk:"delegated_admin_search_size_limit"`
	DelegatedAdminReportSizeLimit  types.Int64  `tfsdk:"delegated_admin_report_size_limit"`
	MembersColumnName              types.String `tfsdk:"members_column_name"`
	NonmembersColumnName           types.String `tfsdk:"nonmembers_column_name"`
}

// GetSchema defines the schema for the datasource.
func (r *restResourceTypeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Rest Resource Type.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of REST Resource Type resource. Options are ['user', 'generic', 'group']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_attribute_category": schema.StringAttribute{
				Description: "Specifies which attribute category the password belongs to.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"password_display_order_index": schema.Int64Attribute{
				Description: "This property determines the display order for the password within its attribute category. Attributes are ordered within their category based on this index from least to greatest.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this REST Resource Type",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the REST Resource Type is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"resource_endpoint": schema.StringAttribute{
				Description: "The HTTP addressable endpoint of this REST Resource Type relative to a REST API base URL. Do not include a leading '/'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"structural_ldap_objectclass": schema.StringAttribute{
				Description: "Specifies the LDAP structural object class that should be exposed by this REST Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auxiliary_ldap_objectclass": schema.SetAttribute{
				Description: "Specifies an auxiliary LDAP object class that should be exposed by this REST Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"search_base_dn": schema.StringAttribute{
				Description: "Specifies the base DN of the branch of the LDAP directory where resources of this type are located.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_filter": schema.SetAttribute{
				Description: "The set of LDAP filters that define the LDAP entries that should be included in this REST Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"parent_dn": schema.StringAttribute{
				Description: "Specifies the DN of the parent entry for new resources of this type, when a parent resource is not provided by the app. The parent DN must be at or below the search base of this resource type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"parent_resource_type": schema.StringAttribute{
				Description: "Specifies the name of another resource type which may be a parent of new resources of this type. The search base DN of the parent resource type must be at or above the search base DN of this resource type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"relative_dn_from_parent_resource": schema.StringAttribute{
				Description: "Specifies a template for a relative DN from the parent resource which identifies the parent entry for a new resource of this type. If this property is not specified then new resources are created immediately below the parent resource or parent DN.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"create_rdn_attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type to be used as the RDN of new resources.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"post_create_constructed_attribute": schema.SetAttribute{
				Description: "Specifies an attribute whose values are to be constructed when a new resource is created. The values are only set at creation time. Subsequent modifications to attributes in the constructed attribute value-pattern are not propagated here.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"update_constructed_attribute": schema.SetAttribute{
				Description: "Specifies an attribute whose values are to be constructed when a resource is updated. The constructed values replace any existing values of the attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this REST Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"search_filter_pattern": schema.StringAttribute{
				Description: "Specifies the LDAP filter that should be used when searching for resources matching provided search text. All attribute types in the filter pattern referencing the search text must have a Delegated Admin Attribute definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"primary_display_attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type which is the primary display attribute. This attribute type must be in the search filter pattern and must have a Delegated Admin Attribute definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"delegated_admin_search_size_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that may be returned from a search request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"delegated_admin_report_size_limit": schema.Int64Attribute{
				Description: "The maximum number of resources that may be included in a report.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"members_column_name": schema.StringAttribute{
				Description: "Specifies the name of the group member column that will be displayed in the Delegated Admin UI",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"nonmembers_column_name": schema.StringAttribute{
				Description: "Specifies the name of the group nonmember column that will be displayed in the Delegated Admin UI",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a UserRestResourceTypeResponse object into the model struct
func readUserRestResourceTypeResponseDataSource(ctx context.Context, r *client.UserRestResourceTypeResponse, state *restResourceTypeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("user")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordAttributeCategory = internaltypes.StringTypeOrNil(r.PasswordAttributeCategory, false)
	state.PasswordDisplayOrderIndex = internaltypes.Int64TypeOrNil(r.PasswordDisplayOrderIndex)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ResourceEndpoint = types.StringValue(r.ResourceEndpoint)
	state.StructuralLDAPObjectclass = types.StringValue(r.StructuralLDAPObjectclass)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.SearchBaseDN = types.StringValue(r.SearchBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ParentDN = internaltypes.StringTypeOrNil(r.ParentDN, false)
	state.ParentResourceType = internaltypes.StringTypeOrNil(r.ParentResourceType, false)
	state.RelativeDNFromParentResource = internaltypes.StringTypeOrNil(r.RelativeDNFromParentResource, false)
	state.CreateRDNAttributeType = internaltypes.StringTypeOrNil(r.CreateRDNAttributeType, false)
	state.PostCreateConstructedAttribute = internaltypes.GetStringSet(r.PostCreateConstructedAttribute)
	state.UpdateConstructedAttribute = internaltypes.GetStringSet(r.UpdateConstructedAttribute)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, false)
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, false)
	state.PrimaryDisplayAttributeType = internaltypes.StringTypeOrNil(r.PrimaryDisplayAttributeType, false)
	state.DelegatedAdminSearchSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminSearchSizeLimit)
	state.DelegatedAdminReportSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminReportSizeLimit)
	state.MembersColumnName = internaltypes.StringTypeOrNil(r.MembersColumnName, false)
	state.NonmembersColumnName = internaltypes.StringTypeOrNil(r.NonmembersColumnName, false)
}

// Read a GenericRestResourceTypeResponse object into the model struct
func readGenericRestResourceTypeResponseDataSource(ctx context.Context, r *client.GenericRestResourceTypeResponse, state *restResourceTypeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ResourceEndpoint = types.StringValue(r.ResourceEndpoint)
	state.StructuralLDAPObjectclass = types.StringValue(r.StructuralLDAPObjectclass)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.SearchBaseDN = types.StringValue(r.SearchBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ParentDN = internaltypes.StringTypeOrNil(r.ParentDN, false)
	state.ParentResourceType = internaltypes.StringTypeOrNil(r.ParentResourceType, false)
	state.RelativeDNFromParentResource = internaltypes.StringTypeOrNil(r.RelativeDNFromParentResource, false)
	state.CreateRDNAttributeType = internaltypes.StringTypeOrNil(r.CreateRDNAttributeType, false)
	state.PostCreateConstructedAttribute = internaltypes.GetStringSet(r.PostCreateConstructedAttribute)
	state.UpdateConstructedAttribute = internaltypes.GetStringSet(r.UpdateConstructedAttribute)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, false)
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, false)
	state.PrimaryDisplayAttributeType = internaltypes.StringTypeOrNil(r.PrimaryDisplayAttributeType, false)
	state.DelegatedAdminSearchSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminSearchSizeLimit)
	state.DelegatedAdminReportSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminReportSizeLimit)
	state.MembersColumnName = internaltypes.StringTypeOrNil(r.MembersColumnName, false)
	state.NonmembersColumnName = internaltypes.StringTypeOrNil(r.NonmembersColumnName, false)
}

// Read a GroupRestResourceTypeResponse object into the model struct
func readGroupRestResourceTypeResponseDataSource(ctx context.Context, r *client.GroupRestResourceTypeResponse, state *restResourceTypeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("group")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MembersColumnName = internaltypes.StringTypeOrNil(r.MembersColumnName, false)
	state.NonmembersColumnName = internaltypes.StringTypeOrNil(r.NonmembersColumnName, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
	state.ResourceEndpoint = types.StringValue(r.ResourceEndpoint)
	state.StructuralLDAPObjectclass = types.StringValue(r.StructuralLDAPObjectclass)
	state.AuxiliaryLDAPObjectclass = internaltypes.GetStringSet(r.AuxiliaryLDAPObjectclass)
	state.SearchBaseDN = types.StringValue(r.SearchBaseDN)
	state.IncludeFilter = internaltypes.GetStringSet(r.IncludeFilter)
	state.ParentDN = internaltypes.StringTypeOrNil(r.ParentDN, false)
	state.ParentResourceType = internaltypes.StringTypeOrNil(r.ParentResourceType, false)
	state.RelativeDNFromParentResource = internaltypes.StringTypeOrNil(r.RelativeDNFromParentResource, false)
	state.CreateRDNAttributeType = internaltypes.StringTypeOrNil(r.CreateRDNAttributeType, false)
	state.PostCreateConstructedAttribute = internaltypes.GetStringSet(r.PostCreateConstructedAttribute)
	state.UpdateConstructedAttribute = internaltypes.GetStringSet(r.UpdateConstructedAttribute)
	state.DisplayName = internaltypes.StringTypeOrNil(r.DisplayName, false)
	state.SearchFilterPattern = internaltypes.StringTypeOrNil(r.SearchFilterPattern, false)
	state.PrimaryDisplayAttributeType = internaltypes.StringTypeOrNil(r.PrimaryDisplayAttributeType, false)
	state.DelegatedAdminSearchSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminSearchSizeLimit)
	state.DelegatedAdminReportSizeLimit = internaltypes.Int64TypeOrNil(r.DelegatedAdminReportSizeLimit)
}

// Read resource information
func (r *restResourceTypeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state restResourceTypeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.RestResourceTypeApi.GetRestResourceType(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Rest Resource Type", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.UserRestResourceTypeResponse != nil {
		readUserRestResourceTypeResponseDataSource(ctx, readResponse.UserRestResourceTypeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericRestResourceTypeResponse != nil {
		readGenericRestResourceTypeResponseDataSource(ctx, readResponse.GenericRestResourceTypeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroupRestResourceTypeResponse != nil {
		readGroupRestResourceTypeResponseDataSource(ctx, readResponse.GroupRestResourceTypeResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
