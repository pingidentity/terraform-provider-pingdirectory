package identitymapper

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
	_ datasource.DataSource              = &identityMapperDataSource{}
	_ datasource.DataSourceWithConfigure = &identityMapperDataSource{}
)

// Create a Identity Mapper data source
func NewIdentityMapperDataSource() datasource.DataSource {
	return &identityMapperDataSource{}
}

// identityMapperDataSource is the datasource implementation.
type identityMapperDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *identityMapperDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_identity_mapper"
}

// Configure adds the provider configured client to the data source.
func (r *identityMapperDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type identityMapperDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Type                      types.String `tfsdk:"type"`
	ExtensionClass            types.String `tfsdk:"extension_class"`
	ExtensionArgument         types.Set    `tfsdk:"extension_argument"`
	AllIncludedIdentityMapper types.Set    `tfsdk:"all_included_identity_mapper"`
	AnyIncludedIdentityMapper types.Set    `tfsdk:"any_included_identity_mapper"`
	ScriptClass               types.String `tfsdk:"script_class"`
	ScriptArgument            types.Set    `tfsdk:"script_argument"`
	MatchAttribute            types.Set    `tfsdk:"match_attribute"`
	MatchPattern              types.String `tfsdk:"match_pattern"`
	ReplacePattern            types.String `tfsdk:"replace_pattern"`
	MatchBaseDN               types.Set    `tfsdk:"match_base_dn"`
	MatchFilter               types.String `tfsdk:"match_filter"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *identityMapperDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Identity Mapper.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Identity Mapper resource. Options are ['exact-match', 'groovy-scripted', 'regular-expression', 'aggregate', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Identity Mapper.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Identity Mapper. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"all_included_identity_mapper": schema.SetAttribute{
				Description: "The set of identity mappers that must all match the target entry. Each identity mapper must uniquely match the same target entry. If any of the identity mappers match multiple entries, if any of them match zero entries, or if any of them match different entries, then the mapping will fail.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"any_included_identity_mapper": schema.SetAttribute{
				Description: "The set of identity mappers that will be used to identify the target entry. At least one identity mapper must uniquely match an entry. If multiple identity mappers match entries, then they must all uniquely match the same entry. If none of the identity mappers match any entries, if any of them match multiple entries, or if any of them match different entries, then the mapping will fail.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Identity Mapper.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Identity Mapper. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"match_attribute": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `exact-match`: Specifies the attribute whose value should exactly match the ID string provided to this identity mapper. When the `type` attribute is set to `regular-expression`: Specifies the name or OID of the attribute whose value should match the provided identifier string after it has been processed by the associated regular expression.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `exact-match`: Specifies the attribute whose value should exactly match the ID string provided to this identity mapper.\n  - `regular-expression`: Specifies the name or OID of the attribute whose value should match the provided identifier string after it has been processed by the associated regular expression.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"match_pattern": schema.StringAttribute{
				Description: "Specifies the regular expression pattern that is used to identify portions of the ID string that will be replaced.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"replace_pattern": schema.StringAttribute{
				Description: "Specifies the replacement pattern that should be used for substrings in the ID string that match the provided regular expression pattern.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"match_base_dn": schema.SetAttribute{
				Description:         "When the `type` attribute is set to `exact-match`: Specifies the set of base DNs below which to search for users. When the `type` attribute is set to `regular-expression`: Specifies the base DN(s) that should be used when performing searches to map the provided ID string to a user entry. If multiple values are given, searches are performed below all the specified base DNs.",
				MarkdownDescription: "When the `type` attribute is set to:\n  - `exact-match`: Specifies the set of base DNs below which to search for users.\n  - `regular-expression`: Specifies the base DN(s) that should be used when performing searches to map the provided ID string to a user entry. If multiple values are given, searches are performed below all the specified base DNs.",
				Required:            false,
				Optional:            false,
				Computed:            true,
				ElementType:         types.StringType,
			},
			"match_filter": schema.StringAttribute{
				Description: "An optional filter that mapped users must match.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Identity Mapper",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Identity Mapper is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a ExactMatchIdentityMapperResponse object into the model struct
func readExactMatchIdentityMapperResponseDataSource(ctx context.Context, r *client.ExactMatchIdentityMapperResponse, state *identityMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("exact-match")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.MatchBaseDN = internaltypes.GetStringSet(r.MatchBaseDN)
	state.MatchFilter = internaltypes.StringTypeOrNil(r.MatchFilter, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GroovyScriptedIdentityMapperResponse object into the model struct
func readGroovyScriptedIdentityMapperResponseDataSource(ctx context.Context, r *client.GroovyScriptedIdentityMapperResponse, state *identityMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a RegularExpressionIdentityMapperResponse object into the model struct
func readRegularExpressionIdentityMapperResponseDataSource(ctx context.Context, r *client.RegularExpressionIdentityMapperResponse, state *identityMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("regular-expression")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.MatchAttribute = internaltypes.GetStringSet(r.MatchAttribute)
	state.MatchBaseDN = internaltypes.GetStringSet(r.MatchBaseDN)
	state.MatchFilter = internaltypes.StringTypeOrNil(r.MatchFilter, false)
	state.MatchPattern = types.StringValue(r.MatchPattern)
	state.ReplacePattern = internaltypes.StringTypeOrNil(r.ReplacePattern, false)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a AggregateIdentityMapperResponse object into the model struct
func readAggregateIdentityMapperResponseDataSource(ctx context.Context, r *client.AggregateIdentityMapperResponse, state *identityMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("aggregate")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.AllIncludedIdentityMapper = internaltypes.GetStringSet(r.AllIncludedIdentityMapper)
	state.AnyIncludedIdentityMapper = internaltypes.GetStringSet(r.AnyIncludedIdentityMapper)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyIdentityMapperResponse object into the model struct
func readThirdPartyIdentityMapperResponseDataSource(ctx context.Context, r *client.ThirdPartyIdentityMapperResponse, state *identityMapperDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *identityMapperDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state identityMapperDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.IdentityMapperApi.GetIdentityMapper(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Identity Mapper", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.ExactMatchIdentityMapperResponse != nil {
		readExactMatchIdentityMapperResponseDataSource(ctx, readResponse.ExactMatchIdentityMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedIdentityMapperResponse != nil {
		readGroovyScriptedIdentityMapperResponseDataSource(ctx, readResponse.GroovyScriptedIdentityMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.RegularExpressionIdentityMapperResponse != nil {
		readRegularExpressionIdentityMapperResponseDataSource(ctx, readResponse.RegularExpressionIdentityMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AggregateIdentityMapperResponse != nil {
		readAggregateIdentityMapperResponseDataSource(ctx, readResponse.AggregateIdentityMapperResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyIdentityMapperResponse != nil {
		readThirdPartyIdentityMapperResponseDataSource(ctx, readResponse.ThirdPartyIdentityMapperResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
