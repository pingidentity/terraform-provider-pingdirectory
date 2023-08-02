package logfieldbehavior

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
	_ datasource.DataSource              = &logFieldBehaviorDataSource{}
	_ datasource.DataSourceWithConfigure = &logFieldBehaviorDataSource{}
)

// Create a Log Field Behavior data source
func NewLogFieldBehaviorDataSource() datasource.DataSource {
	return &logFieldBehaviorDataSource{}
}

// logFieldBehaviorDataSource is the datasource implementation.
type logFieldBehaviorDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logFieldBehaviorDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_field_behavior"
}

// Configure adds the provider configured client to the data source.
func (r *logFieldBehaviorDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type logFieldBehaviorDataSourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Type                             types.String `tfsdk:"type"`
	PreserveField                    types.Set    `tfsdk:"preserve_field"`
	PreserveFieldName                types.Set    `tfsdk:"preserve_field_name"`
	OmitField                        types.Set    `tfsdk:"omit_field"`
	OmitFieldName                    types.Set    `tfsdk:"omit_field_name"`
	RedactEntireValueField           types.Set    `tfsdk:"redact_entire_value_field"`
	RedactEntireValueFieldName       types.Set    `tfsdk:"redact_entire_value_field_name"`
	RedactValueComponentsField       types.Set    `tfsdk:"redact_value_components_field"`
	RedactValueComponentsFieldName   types.Set    `tfsdk:"redact_value_components_field_name"`
	TokenizeEntireValueField         types.Set    `tfsdk:"tokenize_entire_value_field"`
	TokenizeEntireValueFieldName     types.Set    `tfsdk:"tokenize_entire_value_field_name"`
	TokenizeValueComponentsField     types.Set    `tfsdk:"tokenize_value_components_field"`
	TokenizeValueComponentsFieldName types.Set    `tfsdk:"tokenize_value_components_field_name"`
	Description                      types.String `tfsdk:"description"`
	DefaultBehavior                  types.String `tfsdk:"default_behavior"`
}

// GetSchema defines the schema for the datasource.
func (r *logFieldBehaviorDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log Field Behavior.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Field Behavior resource. Options are ['text-access', 'json-formatted-access']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"preserve_field": schema.SetAttribute{
				Description: "The log fields whose values should be logged with the intended value. The values for these fields will be preserved, although they may be sanitized for parsability or safety purposes (for example, to escape special characters in the value), and values that are too long may be truncated.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"preserve_field_name": schema.SetAttribute{
				Description: "The names of any custom fields whose values should be preserved. This should generally only be used for fields that are not available through the preserve-field property (for example, custom log fields defined in Server SDK extensions).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"omit_field": schema.SetAttribute{
				Description: "The log fields that should be omitted entirely from log messages. Neither the field name nor value will be included.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"omit_field_name": schema.SetAttribute{
				Description: "The names of any custom fields that should be omitted from log messages. This should generally only be used for fields that are not available through the omit-field property (for example, custom log fields defined in Server SDK extensions).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"redact_entire_value_field": schema.SetAttribute{
				Description: "The log fields whose values should be completely redacted in log messages. The field name will be included, but with a fixed value that does not reflect the actual value for the field.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"redact_entire_value_field_name": schema.SetAttribute{
				Description: "The names of any custom fields whose values should be completely redacted. This should generally only be used for fields that are not available through the redact-entire-value-field property (for example, custom log fields defined in Server SDK extensions).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"redact_value_components_field": schema.SetAttribute{
				Description: "The log fields whose values will include redacted components.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"redact_value_components_field_name": schema.SetAttribute{
				Description: "The names of any custom fields for which to redact components within the value. This should generally only be used for fields that are not available through the redact-value-components-field property (for example, custom log fields defined in Server SDK extensions).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tokenize_entire_value_field": schema.SetAttribute{
				Description: "The log fields whose values should be completely tokenized in log messages. The field name will be included, but the value will be replaced with a token that does not reveal the actual value, but that is generated from the value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tokenize_entire_value_field_name": schema.SetAttribute{
				Description: "The names of any custom fields whose values should be completely tokenized. This should generally only be used for fields that are not available through the tokenize-entire-value-field property (for example, custom log fields defined in Server SDK extensions).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tokenize_value_components_field": schema.SetAttribute{
				Description: "The log fields whose values will include tokenized components.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"tokenize_value_components_field_name": schema.SetAttribute{
				Description: "The names of any custom fields for which to tokenize components within the value. This should generally only be used for fields that are not available through the tokenize-value-components-field property (for example, custom log fields defined in Server SDK extensions).",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Behavior",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_behavior": schema.StringAttribute{
				Description: "The default behavior that the server should exhibit for fields for which no explicit behavior is defined. If no default behavior is defined, the server will fall back to using the default behavior configured for the syntax used for each log field.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a TextAccessLogFieldBehaviorResponse object into the model struct
func readTextAccessLogFieldBehaviorResponseDataSource(ctx context.Context, r *client.TextAccessLogFieldBehaviorResponse, state *logFieldBehaviorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("text-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PreserveField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessPreserveFieldProp(r.PreserveField))
	state.PreserveFieldName = internaltypes.GetStringSet(r.PreserveFieldName)
	state.OmitField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessOmitFieldProp(r.OmitField))
	state.OmitFieldName = internaltypes.GetStringSet(r.OmitFieldName)
	state.RedactEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessRedactEntireValueFieldProp(r.RedactEntireValueField))
	state.RedactEntireValueFieldName = internaltypes.GetStringSet(r.RedactEntireValueFieldName)
	state.RedactValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessRedactValueComponentsFieldProp(r.RedactValueComponentsField))
	state.RedactValueComponentsFieldName = internaltypes.GetStringSet(r.RedactValueComponentsFieldName)
	state.TokenizeEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessTokenizeEntireValueFieldProp(r.TokenizeEntireValueField))
	state.TokenizeEntireValueFieldName = internaltypes.GetStringSet(r.TokenizeEntireValueFieldName)
	state.TokenizeValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorTextAccessTokenizeValueComponentsFieldProp(r.TokenizeValueComponentsField))
	state.TokenizeValueComponentsFieldName = internaltypes.GetStringSet(r.TokenizeValueComponentsFieldName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldBehaviorDefaultBehaviorProp(r.DefaultBehavior), false)
}

// Read a JsonFormattedAccessLogFieldBehaviorResponse object into the model struct
func readJsonFormattedAccessLogFieldBehaviorResponseDataSource(ctx context.Context, r *client.JsonFormattedAccessLogFieldBehaviorResponse, state *logFieldBehaviorDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json-formatted-access")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.PreserveField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessPreserveFieldProp(r.PreserveField))
	state.PreserveFieldName = internaltypes.GetStringSet(r.PreserveFieldName)
	state.OmitField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessOmitFieldProp(r.OmitField))
	state.OmitFieldName = internaltypes.GetStringSet(r.OmitFieldName)
	state.RedactEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessRedactEntireValueFieldProp(r.RedactEntireValueField))
	state.RedactEntireValueFieldName = internaltypes.GetStringSet(r.RedactEntireValueFieldName)
	state.RedactValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessRedactValueComponentsFieldProp(r.RedactValueComponentsField))
	state.RedactValueComponentsFieldName = internaltypes.GetStringSet(r.RedactValueComponentsFieldName)
	state.TokenizeEntireValueField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessTokenizeEntireValueFieldProp(r.TokenizeEntireValueField))
	state.TokenizeEntireValueFieldName = internaltypes.GetStringSet(r.TokenizeEntireValueFieldName)
	state.TokenizeValueComponentsField = internaltypes.GetStringSet(
		client.StringSliceEnumlogFieldBehaviorJsonFormattedAccessTokenizeValueComponentsFieldProp(r.TokenizeValueComponentsField))
	state.TokenizeValueComponentsFieldName = internaltypes.GetStringSet(r.TokenizeValueComponentsFieldName)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldBehaviorDefaultBehaviorProp(r.DefaultBehavior), false)
}

// Read resource information
func (r *logFieldBehaviorDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logFieldBehaviorDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldBehaviorApi.GetLogFieldBehavior(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Behavior", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.TextAccessLogFieldBehaviorResponse != nil {
		readTextAccessLogFieldBehaviorResponseDataSource(ctx, readResponse.TextAccessLogFieldBehaviorResponse, &state, &resp.Diagnostics)
	}
	if readResponse.JsonFormattedAccessLogFieldBehaviorResponse != nil {
		readJsonFormattedAccessLogFieldBehaviorResponseDataSource(ctx, readResponse.JsonFormattedAccessLogFieldBehaviorResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
