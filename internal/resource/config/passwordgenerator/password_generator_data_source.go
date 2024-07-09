package passwordgenerator

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	client "github.com/pingidentity/pingdirectory-go-client/v10100/configurationapi"
	"github.com/pingidentity/terraform-provider-pingdirectory/internal/resource/config"
	internaltypes "github.com/pingidentity/terraform-provider-pingdirectory/internal/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &passwordGeneratorDataSource{}
	_ datasource.DataSourceWithConfigure = &passwordGeneratorDataSource{}
)

// Create a Password Generator data source
func NewPasswordGeneratorDataSource() datasource.DataSource {
	return &passwordGeneratorDataSource{}
}

// passwordGeneratorDataSource is the datasource implementation.
type passwordGeneratorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *passwordGeneratorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password_generator"
}

// Configure adds the provider configured client to the data source.
func (r *passwordGeneratorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type passwordGeneratorDataSourceModel struct {
	Id                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name"`
	Type                      types.String `tfsdk:"type"`
	ExtensionClass            types.String `tfsdk:"extension_class"`
	ExtensionArgument         types.Set    `tfsdk:"extension_argument"`
	DictionaryFile            types.String `tfsdk:"dictionary_file"`
	MinimumPasswordCharacters types.Int64  `tfsdk:"minimum_password_characters"`
	MinimumPasswordWords      types.Int64  `tfsdk:"minimum_password_words"`
	CapitalizeWords           types.Bool   `tfsdk:"capitalize_words"`
	ScriptClass               types.String `tfsdk:"script_class"`
	ScriptArgument            types.Set    `tfsdk:"script_argument"`
	PasswordCharacterSet      types.Set    `tfsdk:"password_character_set"`
	PasswordFormat            types.String `tfsdk:"password_format"`
	Description               types.String `tfsdk:"description"`
	Enabled                   types.Bool   `tfsdk:"enabled"`
}

// GetSchema defines the schema for the datasource.
func (r *passwordGeneratorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Password Generator.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Password Generator resource. Options are ['random', 'groovy-scripted', 'passphrase', 'third-party']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Java class providing the logic for the Third Party Password Generator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"extension_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Third Party Password Generator. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"dictionary_file": schema.StringAttribute{
				Description: "The path to the dictionary file that will be used to obtain the words for use in generated passwords.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_password_characters": schema.Int64Attribute{
				Description: "The minimum number of characters that generated passwords will be required to have.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"minimum_password_words": schema.Int64Attribute{
				Description: "The minimum number of words that must be concatenated in the course of generating a password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"capitalize_words": schema.BoolAttribute{
				Description: "Indicates whether to capitalize each word used in the generated password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_class": schema.StringAttribute{
				Description: "The fully-qualified name of the Groovy class providing the logic for the Groovy Scripted Password Generator.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"script_argument": schema.SetAttribute{
				Description: "The set of arguments used to customize the behavior for the Scripted Password Generator. Each configuration property should be given in the form 'name=value'.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"password_character_set": schema.SetAttribute{
				Description: "Specifies one or more named character sets.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"password_format": schema.StringAttribute{
				Description: "Specifies the format to use for the generated password.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Password Generator",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Indicates whether the Password Generator is enabled for use.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a RandomPasswordGeneratorResponse object into the model struct
func readRandomPasswordGeneratorResponseDataSource(ctx context.Context, r *client.RandomPasswordGeneratorResponse, state *passwordGeneratorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("random")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PasswordCharacterSet = internaltypes.GetStringSet(r.PasswordCharacterSet)
	state.PasswordFormat = types.StringValue(r.PasswordFormat)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a GroovyScriptedPasswordGeneratorResponse object into the model struct
func readGroovyScriptedPasswordGeneratorResponseDataSource(ctx context.Context, r *client.GroovyScriptedPasswordGeneratorResponse, state *passwordGeneratorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("groovy-scripted")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ScriptClass = types.StringValue(r.ScriptClass)
	state.ScriptArgument = internaltypes.GetStringSet(r.ScriptArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a PassphrasePasswordGeneratorResponse object into the model struct
func readPassphrasePasswordGeneratorResponseDataSource(ctx context.Context, r *client.PassphrasePasswordGeneratorResponse, state *passwordGeneratorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("passphrase")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DictionaryFile = types.StringValue(r.DictionaryFile)
	state.MinimumPasswordCharacters = internaltypes.Int64TypeOrNil(r.MinimumPasswordCharacters)
	state.MinimumPasswordWords = internaltypes.Int64TypeOrNil(r.MinimumPasswordWords)
	state.CapitalizeWords = internaltypes.BoolTypeOrNil(r.CapitalizeWords)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read a ThirdPartyPasswordGeneratorResponse object into the model struct
func readThirdPartyPasswordGeneratorResponseDataSource(ctx context.Context, r *client.ThirdPartyPasswordGeneratorResponse, state *passwordGeneratorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("third-party")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.ExtensionClass = types.StringValue(r.ExtensionClass)
	state.ExtensionArgument = internaltypes.GetStringSet(r.ExtensionArgument)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.Enabled = types.BoolValue(r.Enabled)
}

// Read resource information
func (r *passwordGeneratorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state passwordGeneratorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.PasswordGeneratorAPI.GetPasswordGenerator(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Password Generator", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.RandomPasswordGeneratorResponse != nil {
		readRandomPasswordGeneratorResponseDataSource(ctx, readResponse.RandomPasswordGeneratorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GroovyScriptedPasswordGeneratorResponse != nil {
		readGroovyScriptedPasswordGeneratorResponseDataSource(ctx, readResponse.GroovyScriptedPasswordGeneratorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.PassphrasePasswordGeneratorResponse != nil {
		readPassphrasePasswordGeneratorResponseDataSource(ctx, readResponse.PassphrasePasswordGeneratorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.ThirdPartyPasswordGeneratorResponse != nil {
		readThirdPartyPasswordGeneratorResponseDataSource(ctx, readResponse.ThirdPartyPasswordGeneratorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
