package scimsubattribute

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
	_ datasource.DataSource              = &scimSubattributeDataSource{}
	_ datasource.DataSourceWithConfigure = &scimSubattributeDataSource{}
)

// Create a Scim Subattribute data source
func NewScimSubattributeDataSource() datasource.DataSource {
	return &scimSubattributeDataSource{}
}

// scimSubattributeDataSource is the datasource implementation.
type scimSubattributeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *scimSubattributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_scim_subattribute"
}

// Configure adds the provider configured client to the data source.
func (r *scimSubattributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type scimSubattributeDataSourceModel struct {
	Id                types.String `tfsdk:"id"`
	Name              types.String `tfsdk:"name"`
	ResourceType      types.String `tfsdk:"resource_type"`
	ScimAttributeName types.String `tfsdk:"scim_attribute_name"`
	ScimSchemaName    types.String `tfsdk:"scim_schema_name"`
	Description       types.String `tfsdk:"description"`
	Type              types.String `tfsdk:"type"`
	Required          types.Bool   `tfsdk:"required"`
	CaseExact         types.Bool   `tfsdk:"case_exact"`
	MultiValued       types.Bool   `tfsdk:"multi_valued"`
	CanonicalValue    types.Set    `tfsdk:"canonical_value"`
	Mutability        types.String `tfsdk:"mutability"`
	Returned          types.String `tfsdk:"returned"`
	ReferenceType     types.Set    `tfsdk:"reference_type"`
}

// GetSchema defines the schema for the datasource.
func (r *scimSubattributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Scim Subattribute.",
		Attributes: map[string]schema.Attribute{
			"resource_type": schema.StringAttribute{
				Description: "The type of SCIM Subattribute resource. Options are ['scim-subattribute']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"scim_attribute_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Attribute",
				Required:    true,
			},
			"scim_schema_name": schema.StringAttribute{
				Description: "Name of the parent SCIM Schema",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this SCIM Subattribute",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Specifies the data type for this sub-attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"required": schema.BoolAttribute{
				Description: "Specifies whether this sub-attribute is required.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"case_exact": schema.BoolAttribute{
				Description: "Specifies whether the sub-attribute values are case sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"multi_valued": schema.BoolAttribute{
				Description: "Specifies whether this attribute may have multiple values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"canonical_value": schema.SetAttribute{
				Description: "Specifies the suggested canonical type values for the sub-attribute.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"mutability": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the sub-attribute can be written.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"returned": schema.StringAttribute{
				Description: "Specifies the circumstances under which the values of the sub-attribute are returned in response to a request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"reference_type": schema.SetAttribute{
				Description: "Specifies the SCIM resource types that may be referenced. This property is only applicable for sub-attributes that are of type 'reference'. Valid values are: A SCIM resource type (e.g., 'User' or 'Group'), 'external' - indicating the resource is an external resource (e.g., such as a photo), or 'uri' - indicating that the reference is to a service endpoint or an identifier (such as a schema urn).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ScimSubattributeResponse object into the model struct
func readScimSubattributeResponseDataSource(ctx context.Context, r *client.ScimSubattributeResponse, state *scimSubattributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.ResourceType = types.StringValue("scim-subattribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Type = types.StringValue(r.Type.String())
	state.Required = types.BoolValue(r.Required)
	state.CaseExact = types.BoolValue(r.CaseExact)
	state.MultiValued = types.BoolValue(r.MultiValued)
	state.CanonicalValue = internaltypes.GetStringSet(r.CanonicalValue)
	state.Mutability = types.StringValue(r.Mutability.String())
	state.Returned = types.StringValue(r.Returned.String())
	state.ReferenceType = internaltypes.GetStringSet(r.ReferenceType)
}

// Read resource information
func (r *scimSubattributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state scimSubattributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ScimSubattributeAPI.GetScimSubattribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.ScimAttributeName.ValueString(), state.ScimSchemaName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Scim Subattribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readScimSubattributeResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
