package sensitiveattribute

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
	_ datasource.DataSource              = &sensitiveAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &sensitiveAttributeDataSource{}
)

// Create a Sensitive Attribute data source
func NewSensitiveAttributeDataSource() datasource.DataSource {
	return &sensitiveAttributeDataSource{}
}

// sensitiveAttributeDataSource is the datasource implementation.
type sensitiveAttributeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *sensitiveAttributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_attribute"
}

// Configure adds the provider configured client to the data source.
func (r *sensitiveAttributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type sensitiveAttributeDataSourceModel struct {
	Id                                           types.String `tfsdk:"id"`
	Name                                         types.String `tfsdk:"name"`
	Type                                         types.String `tfsdk:"type"`
	Description                                  types.String `tfsdk:"description"`
	AttributeType                                types.Set    `tfsdk:"attribute_type"`
	IncludeDefaultSensitiveOperationalAttributes types.Bool   `tfsdk:"include_default_sensitive_operational_attributes"`
	AllowInReturnedEntries                       types.String `tfsdk:"allow_in_returned_entries"`
	AllowInFilter                                types.String `tfsdk:"allow_in_filter"`
	AllowInAdd                                   types.String `tfsdk:"allow_in_add"`
	AllowInCompare                               types.String `tfsdk:"allow_in_compare"`
	AllowInModify                                types.String `tfsdk:"allow_in_modify"`
}

// GetSchema defines the schema for the datasource.
func (r *sensitiveAttributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Sensitive Attribute.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Sensitive Attribute resource. Options are ['sensitive-attribute']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Sensitive Attribute",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_type": schema.SetAttribute{
				Description: "The name(s) or OID(s) of the attribute types for attributes whose values may be considered sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"include_default_sensitive_operational_attributes": schema.BoolAttribute{
				Description: "Indicates whether to automatically include any server-generated operational attributes that may contain sensitive data.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_in_returned_entries": schema.StringAttribute{
				Description: "Indicates whether sensitive attributes should be included in entries returned to the client. This includes not only search result entries, but also other forms including in the values of controls like the pre-read, post-read, get authorization entry, and LDAP join response controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_in_filter": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to include sensitive attributes in search filters. This also includes filters that may be used in other forms, including assertion and LDAP join request controls.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_in_add": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to include sensitive attributes in add requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_in_compare": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to target sensitive attributes with compare requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_in_modify": schema.StringAttribute{
				Description: "Indicates whether clients will be allowed to target sensitive attributes with modify requests.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SensitiveAttributeResponse object into the model struct
func readSensitiveAttributeResponseDataSource(ctx context.Context, r *client.SensitiveAttributeResponse, state *sensitiveAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("sensitive-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AttributeType = internaltypes.GetStringSet(r.AttributeType)
	state.IncludeDefaultSensitiveOperationalAttributes = internaltypes.BoolTypeOrNil(r.IncludeDefaultSensitiveOperationalAttributes)
	state.AllowInReturnedEntries = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInReturnedEntriesProp(r.AllowInReturnedEntries), false)
	state.AllowInFilter = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInFilterProp(r.AllowInFilter), false)
	state.AllowInAdd = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInAddProp(r.AllowInAdd), false)
	state.AllowInCompare = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInCompareProp(r.AllowInCompare), false)
	state.AllowInModify = internaltypes.StringTypeOrNil(
		client.StringPointerEnumsensitiveAttributeAllowInModifyProp(r.AllowInModify), false)
}

// Read resource information
func (r *sensitiveAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state sensitiveAttributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SensitiveAttributeAPI.GetSensitiveAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Sensitive Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSensitiveAttributeResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
