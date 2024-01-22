package delegatedadminattribute

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
	_ datasource.DataSource              = &delegatedAdminAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminAttributeDataSource{}
)

// Create a Delegated Admin Attribute data source
func NewDelegatedAdminAttributeDataSource() datasource.DataSource {
	return &delegatedAdminAttributeDataSource{}
}

// delegatedAdminAttributeDataSource is the datasource implementation.
type delegatedAdminAttributeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminAttributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_attribute"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminAttributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type delegatedAdminAttributeDataSourceModel struct {
	Id                    types.String `tfsdk:"id"`
	Type                  types.String `tfsdk:"type"`
	RestResourceTypeName  types.String `tfsdk:"rest_resource_type_name"`
	AllowedMIMEType       types.Set    `tfsdk:"allowed_mime_type"`
	Description           types.String `tfsdk:"description"`
	AttributeType         types.String `tfsdk:"attribute_type"`
	DisplayName           types.String `tfsdk:"display_name"`
	Mutability            types.String `tfsdk:"mutability"`
	IncludeInSummary      types.Bool   `tfsdk:"include_in_summary"`
	MultiValued           types.Bool   `tfsdk:"multi_valued"`
	AttributeCategory     types.String `tfsdk:"attribute_category"`
	DisplayOrderIndex     types.Int64  `tfsdk:"display_order_index"`
	ReferenceResourceType types.String `tfsdk:"reference_resource_type"`
	AttributePresentation types.String `tfsdk:"attribute_presentation"`
	DateTimeFormat        types.String `tfsdk:"date_time_format"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminAttributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Delegated Admin Attribute.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Delegated Admin Attribute resource. Options are ['certificate', 'photo', 'generic']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"rest_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent REST Resource Type",
				Required:    true,
			},
			"allowed_mime_type": schema.SetAttribute{
				Description: "The list of file types allowed to be uploaded. If no types are specified, then all types will be allowed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Delegated Admin Attribute",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the name or OID of the LDAP attribute type.",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"mutability": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the attribute can be written.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"include_in_summary": schema.BoolAttribute{
				Description: "Indicates whether this Delegated Admin Attribute is to be included in the summary display for a resource.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"multi_valued": schema.BoolAttribute{
				Description: "Indicates whether this Delegated Admin Attribute may have multiple values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_category": schema.StringAttribute{
				Description: "Specifies which attribute category this attribute belongs to.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"display_order_index": schema.Int64Attribute{
				Description: "This property determines a display order for attributes within a given attribute category. Attributes are ordered within their category based on this index from least to greatest.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"reference_resource_type": schema.StringAttribute{
				Description: "For LDAP attributes with DN syntax, specifies what kind of resource is referenced.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_presentation": schema.StringAttribute{
				Description: "Indicates how the attribute is presented to the user of the app.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"date_time_format": schema.StringAttribute{
				Description: "Specifies the format string that is used to present a date and/or time value to the user of the app. This property only applies to LDAP attribute types whose LDAP syntax is GeneralizedTime and is ignored if the attribute type has any other syntax.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a CertificateDelegatedAdminAttributeResponse object into the model struct
func readCertificateDelegatedAdminAttributeResponseDataSource(ctx context.Context, r *client.CertificateDelegatedAdminAttributeResponse, state *delegatedAdminAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("certificate")
	state.Id = types.StringValue(r.Id)
	state.AllowedMIMEType = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminAttributeCertificateAllowedMIMETypeProp(r.AllowedMIMEType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, false)
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, false)
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), false)
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, false)
}

// Read a PhotoDelegatedAdminAttributeResponse object into the model struct
func readPhotoDelegatedAdminAttributeResponseDataSource(ctx context.Context, r *client.PhotoDelegatedAdminAttributeResponse, state *delegatedAdminAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("photo")
	state.Id = types.StringValue(r.Id)
	state.AllowedMIMEType = internaltypes.GetStringSet(
		client.StringSliceEnumdelegatedAdminAttributePhotoAllowedMIMETypeProp(r.AllowedMIMEType))
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, false)
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, false)
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), false)
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, false)
}

// Read a GenericDelegatedAdminAttributeResponse object into the model struct
func readGenericDelegatedAdminAttributeResponseDataSource(ctx context.Context, r *client.GenericDelegatedAdminAttributeResponse, state *delegatedAdminAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.IncludeInSummary = types.BoolValue(r.IncludeInSummary)
	state.AttributeCategory = internaltypes.StringTypeOrNil(r.AttributeCategory, false)
	state.DisplayOrderIndex = types.Int64Value(r.DisplayOrderIndex)
	state.ReferenceResourceType = internaltypes.StringTypeOrNil(r.ReferenceResourceType, false)
	state.AttributePresentation = internaltypes.StringTypeOrNil(
		client.StringPointerEnumdelegatedAdminAttributeAttributePresentationProp(r.AttributePresentation), false)
	state.DateTimeFormat = internaltypes.StringTypeOrNil(r.DateTimeFormat, false)
}

// Read resource information
func (r *delegatedAdminAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminAttributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminAttributeAPI.GetDelegatedAdminAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AttributeType.ValueString(), state.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.CertificateDelegatedAdminAttributeResponse != nil {
		readCertificateDelegatedAdminAttributeResponseDataSource(ctx, readResponse.CertificateDelegatedAdminAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PhotoDelegatedAdminAttributeResponse != nil {
		readPhotoDelegatedAdminAttributeResponseDataSource(ctx, readResponse.PhotoDelegatedAdminAttributeResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericDelegatedAdminAttributeResponse != nil {
		readGenericDelegatedAdminAttributeResponseDataSource(ctx, readResponse.GenericDelegatedAdminAttributeResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
