package constructedattribute

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10200/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &constructedAttributeDataSource{}
	_ datasource.DataSourceWithConfigure = &constructedAttributeDataSource{}
)

// Create a Constructed Attribute data source
func NewConstructedAttributeDataSource() datasource.DataSource {
	return &constructedAttributeDataSource{}
}

// constructedAttributeDataSource is the datasource implementation.
type constructedAttributeDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *constructedAttributeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_constructed_attribute"
}

// Configure adds the provider configured client to the data source.
func (r *constructedAttributeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type constructedAttributeDataSourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Type          types.String `tfsdk:"type"`
	Description   types.String `tfsdk:"description"`
	AttributeType types.String `tfsdk:"attribute_type"`
	ValuePattern  types.Set    `tfsdk:"value_pattern"`
}

// GetSchema defines the schema for the datasource.
func (r *constructedAttributeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Constructed Attribute.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Constructed Attribute resource. Options are ['constructed-attribute']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Constructed Attribute",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "Specifies the attribute type for the attribute whose values are to be constructed.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"value_pattern": schema.SetAttribute{
				Description: "Specifies a pattern for constructing the attribute value using fixed text and attribute values from the entry.",
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

// Read a ConstructedAttributeResponse object into the model struct
func readConstructedAttributeResponseDataSource(ctx context.Context, r *client.ConstructedAttributeResponse, state *constructedAttributeDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("constructed-attribute")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.ValuePattern = internaltypes.GetStringSet(r.ValuePattern)
}

// Read resource information
func (r *constructedAttributeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state constructedAttributeDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.ConstructedAttributeAPI.GetConstructedAttribute(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Constructed Attribute", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readConstructedAttributeResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
