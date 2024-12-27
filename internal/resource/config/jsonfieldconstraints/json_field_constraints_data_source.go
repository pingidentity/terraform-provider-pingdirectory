package jsonfieldconstraints

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
	_ datasource.DataSource              = &jsonFieldConstraintsDataSource{}
	_ datasource.DataSourceWithConfigure = &jsonFieldConstraintsDataSource{}
)

// Create a Json Field Constraints data source
func NewJsonFieldConstraintsDataSource() datasource.DataSource {
	return &jsonFieldConstraintsDataSource{}
}

// jsonFieldConstraintsDataSource is the datasource implementation.
type jsonFieldConstraintsDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *jsonFieldConstraintsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_json_field_constraints"
}

// Configure adds the provider configured client to the data source.
func (r *jsonFieldConstraintsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type jsonFieldConstraintsDataSourceModel struct {
	Id                            types.String `tfsdk:"id"`
	Type                          types.String `tfsdk:"type"`
	JsonAttributeConstraintsName  types.String `tfsdk:"json_attribute_constraints_name"`
	Description                   types.String `tfsdk:"description"`
	JsonField                     types.String `tfsdk:"json_field"`
	ValueType                     types.String `tfsdk:"value_type"`
	IsRequired                    types.Bool   `tfsdk:"is_required"`
	IsArray                       types.String `tfsdk:"is_array"`
	AllowNullValue                types.Bool   `tfsdk:"allow_null_value"`
	AllowEmptyObject              types.Bool   `tfsdk:"allow_empty_object"`
	IndexValues                   types.Bool   `tfsdk:"index_values"`
	IndexEntryLimit               types.Int64  `tfsdk:"index_entry_limit"`
	PrimeIndex                    types.Bool   `tfsdk:"prime_index"`
	CacheMode                     types.String `tfsdk:"cache_mode"`
	TokenizeValues                types.Bool   `tfsdk:"tokenize_values"`
	AllowedValue                  types.Set    `tfsdk:"allowed_value"`
	AllowedValueRegularExpression types.Set    `tfsdk:"allowed_value_regular_expression"`
	MinimumNumericValue           types.String `tfsdk:"minimum_numeric_value"`
	MaximumNumericValue           types.String `tfsdk:"maximum_numeric_value"`
	MinimumValueLength            types.Int64  `tfsdk:"minimum_value_length"`
	MaximumValueLength            types.Int64  `tfsdk:"maximum_value_length"`
	MinimumValueCount             types.Int64  `tfsdk:"minimum_value_count"`
	MaximumValueCount             types.Int64  `tfsdk:"maximum_value_count"`
}

// GetSchema defines the schema for the datasource.
func (r *jsonFieldConstraintsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Json Field Constraints.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of JSON Field Constraints resource. Options are ['json-field-constraints']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"json_attribute_constraints_name": schema.StringAttribute{
				Description: "Name of the parent JSON Attribute Constraints",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this JSON Field Constraints",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"json_field": schema.StringAttribute{
				Description: "The full name of the JSON field to which these constraints apply.",
				Required:    true,
			},
			"value_type": schema.StringAttribute{
				Description: "The data type that will be required for values of the target field.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"is_required": schema.BoolAttribute{
				Description: "Indicates whether the target field must be present in JSON objects stored as values of the associated attribute type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"is_array": schema.StringAttribute{
				Description: "Indicates whether the value of the target field may be an array of values rather than a single value. If this property is set to \"required\" or \"optional\", then the constraints defined for this field will be applied to each element of the array.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_null_value": schema.BoolAttribute{
				Description: "Indicates whether the target field may have a value that is the JSON null value as an alternative to a value (or array of values) of the specified value-type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allow_empty_object": schema.BoolAttribute{
				Description: "Indicates whether the target field may have a value that is an empty JSON object (i.e., a JSON object with zero fields). This may only be set to true if value-type property is set to object.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_values": schema.BoolAttribute{
				Description: "Indicates whether backends that support JSON indexing should maintain an index for values of the target field.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"index_entry_limit": schema.Int64Attribute{
				Description: "The maximum number of entries that may contain a particular value for the target field before the server will stop maintaining the index for that value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"prime_index": schema.BoolAttribute{
				Description: "Indicates whether backends that support database priming should load the contents of the associated JSON index into memory whenever the backend is opened.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"cache_mode": schema.StringAttribute{
				Description: "Specifies the behavior that the server should exhibit when caching data for the associated JSON index. This can be useful in environments in which the system does not have enough memory to fully cache the entire data set, as it makes it possible to prioritize which data is the most important to keep in memory.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"tokenize_values": schema.BoolAttribute{
				Description: "Indicates whether the backend should attempt to assign a compact token for each distinct value for the target field in an attempt to reduce the encoded size of the field in JSON objects. These tokens would be assigned prior to using any from the token set used for automatic compaction of some JSON string values.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"allowed_value": schema.SetAttribute{
				Description: "Specifies an explicit set of string values that will be the only values permitted for the target field. If a set of allowed values is defined, then the server will reject any attempt to store a JSON object with a value for the target field that is not included in that set.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"allowed_value_regular_expression": schema.SetAttribute{
				Description: "Specifies an explicit set of regular expressions that may be used to restrict the set of values that may be used for the target field. If a set of allowed value regular expressions is defined, then the server will reject any attempt to store a JSON object with a value for the target field that does not match at least one of those regular expressions.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"minimum_numeric_value": schema.StringAttribute{
				Description: "Specifies the smallest numeric value that may be used as the value for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is less than that minimum numeric value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_numeric_value": schema.StringAttribute{
				Description: "Specifies the largest numeric value that may be used as the value for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is greater than that maximum numeric value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_value_length": schema.Int64Attribute{
				Description: "Specifies the smallest number of characters that may be present in string values of the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is shorter than that minimum value length.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_value_length": schema.Int64Attribute{
				Description: "Specifies the largest number of characters that may be present in string values of the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is longer than that maximum value length.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_value_count": schema.Int64Attribute{
				Description: "Specifies the smallest number of elements that may be present in an array of values for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is an array with fewer than this number of elements.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"maximum_value_count": schema.Int64Attribute{
				Description: "Specifies the largest number of elements that may be present in an array of values for the target field. If configured, then the server will reject any attempt to store a JSON object with a value for the target field that is an array with more than this number of elements.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, false)
	resp.Schema = schemaDef
}

// Read a JsonFieldConstraintsResponse object into the model struct
func readJsonFieldConstraintsResponseDataSource(ctx context.Context, r *client.JsonFieldConstraintsResponse, state *jsonFieldConstraintsDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-field-constraints")
	state.Id = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.JsonField = types.StringValue(r.JsonField)
	state.ValueType = types.StringValue(r.ValueType.String())
	state.IsRequired = internaltypes.BoolTypeOrNil(r.IsRequired)
	state.IsArray = internaltypes.StringTypeOrNil(
		client.StringPointerEnumjsonFieldConstraintsIsArrayProp(r.IsArray), false)
	state.AllowNullValue = internaltypes.BoolTypeOrNil(r.AllowNullValue)
	state.AllowEmptyObject = internaltypes.BoolTypeOrNil(r.AllowEmptyObject)
	state.IndexValues = internaltypes.BoolTypeOrNil(r.IndexValues)
	state.IndexEntryLimit = internaltypes.Int64TypeOrNil(r.IndexEntryLimit)
	state.PrimeIndex = internaltypes.BoolTypeOrNil(r.PrimeIndex)
	state.CacheMode = internaltypes.StringTypeOrNil(
		client.StringPointerEnumjsonFieldConstraintsCacheModeProp(r.CacheMode), false)
	state.TokenizeValues = internaltypes.BoolTypeOrNil(r.TokenizeValues)
	state.AllowedValue = internaltypes.GetStringSet(r.AllowedValue)
	state.AllowedValueRegularExpression = internaltypes.GetStringSet(r.AllowedValueRegularExpression)
	state.MinimumNumericValue = internaltypes.StringTypeOrNil(r.MinimumNumericValue, false)
	state.MaximumNumericValue = internaltypes.StringTypeOrNil(r.MaximumNumericValue, false)
	state.MinimumValueLength = internaltypes.Int64TypeOrNil(r.MinimumValueLength)
	state.MaximumValueLength = internaltypes.Int64TypeOrNil(r.MaximumValueLength)
	state.MinimumValueCount = internaltypes.Int64TypeOrNil(r.MinimumValueCount)
	state.MaximumValueCount = internaltypes.Int64TypeOrNil(r.MaximumValueCount)
}

// Read resource information
func (r *jsonFieldConstraintsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state jsonFieldConstraintsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.JsonFieldConstraintsAPI.GetJsonFieldConstraints(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.JsonField.ValueString(), state.JsonAttributeConstraintsName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Json Field Constraints", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readJsonFieldConstraintsResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
