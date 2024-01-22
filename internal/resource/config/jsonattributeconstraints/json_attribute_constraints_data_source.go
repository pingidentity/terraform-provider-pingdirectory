package jsonattributeconstraints

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
	_ datasource.DataSource              = &jsonAttributeConstraintsDataSource{}
	_ datasource.DataSourceWithConfigure = &jsonAttributeConstraintsDataSource{}
)

// Create a Json Attribute Constraints data source
func NewJsonAttributeConstraintsDataSource() datasource.DataSource {
	return &jsonAttributeConstraintsDataSource{}
}

// jsonAttributeConstraintsDataSource is the datasource implementation.
type jsonAttributeConstraintsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *jsonAttributeConstraintsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_json_attribute_constraints"
}

// Configure adds the provider configured client to the data source.
func (r *jsonAttributeConstraintsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type jsonAttributeConstraintsDataSourceModel struct {
	Id                 types.String `tfsdk:"id"`
	Type               types.String `tfsdk:"type"`
	Description        types.String `tfsdk:"description"`
	Enabled            types.Bool   `tfsdk:"enabled"`
	AttributeType      types.String `tfsdk:"attribute_type"`
	AllowUnnamedFields types.Bool   `tfsdk:"allow_unnamed_fields"`
}

// GetSchema defines the schema for the datasource.
func (r *jsonAttributeConstraintsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Json Attribute Constraints.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of JSON Attribute Constraints resource. Options are ['json-attribute-constraints']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this JSON Attribute Constraints",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether this JSON Attribute Constraints is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"attribute_type": schema.StringAttribute{
				Description: "The name or OID of the LDAP attribute type whose values will be subject to the associated field constraints. This attribute type must be defined in the server schema, and it must have a \"JSON object\" syntax.",
				Required:    true,
			},
			"allow_unnamed_fields": schema.BoolAttribute{
				Description: "Indicates whether JSON objects stored as values of attributes with the associated attribute-type will be permitted to include fields for which there is no subordinate json-field-constraints definition. If unnamed fields are allowed, then no constraints will be imposed on the values of those fields. However, if unnamed fields are not allowed, then the server will reject any attempt to store a JSON object with a field for which there is no corresponding json-fields-constraints definition.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a JsonAttributeConstraintsResponse object into the model struct
func readJsonAttributeConstraintsResponseDataSource(ctx context.Context, r *client.JsonAttributeConstraintsResponse, state *jsonAttributeConstraintsDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-attribute-constraints")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = internaltypes.BoolTypeOrNil(r.Enabled)
	state.AttributeType = types.StringValue(r.AttributeType)
	state.AllowUnnamedFields = internaltypes.BoolTypeOrNil(r.AllowUnnamedFields)
}

// Read resource information
func (r *jsonAttributeConstraintsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state jsonAttributeConstraintsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.JsonAttributeConstraintsAPI.GetJsonAttributeConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.AttributeType.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Attribute Constraints", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJsonAttributeConstraintsResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
