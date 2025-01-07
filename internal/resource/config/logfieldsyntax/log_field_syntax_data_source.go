package logfieldsyntax

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
	_ datasource.DataSource              = &logFieldSyntaxDataSource{}
	_ datasource.DataSourceWithConfigure = &logFieldSyntaxDataSource{}
)

// Create a Log Field Syntax data source
func NewLogFieldSyntaxDataSource() datasource.DataSource {
	return &logFieldSyntaxDataSource{}
}

// logFieldSyntaxDataSource is the datasource implementation.
type logFieldSyntaxDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *logFieldSyntaxDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_log_field_syntax"
}

// Configure adds the provider configured client to the data source.
func (r *logFieldSyntaxDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type logFieldSyntaxDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Name                       types.String `tfsdk:"name"`
	Type                       types.String `tfsdk:"type"`
	IncludedSensitiveAttribute types.Set    `tfsdk:"included_sensitive_attribute"`
	ExcludedSensitiveAttribute types.Set    `tfsdk:"excluded_sensitive_attribute"`
	IncludedSensitiveField     types.Set    `tfsdk:"included_sensitive_field"`
	ExcludedSensitiveField     types.Set    `tfsdk:"excluded_sensitive_field"`
	Description                types.String `tfsdk:"description"`
	DefaultBehavior            types.String `tfsdk:"default_behavior"`
}

// GetSchema defines the schema for the datasource.
func (r *logFieldSyntaxDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Log Field Syntax.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Log Field Syntax resource. Options are ['json', 'attribute-based', 'generic']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"included_sensitive_attribute": schema.SetAttribute{
				Description: "The set of attribute types that will be considered sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_sensitive_attribute": schema.SetAttribute{
				Description: "The set of attribute types that will not be considered sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"included_sensitive_field": schema.SetAttribute{
				Description: "The names of the JSON fields that will be considered sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"excluded_sensitive_field": schema.SetAttribute{
				Description: "The names of the JSON fields that will not be considered sensitive.",
				Required:    false,
				Optional:    false,
				Computed:    true,
				ElementType: types.StringType,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Log Field Syntax",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"default_behavior": schema.StringAttribute{
				Description: "The default behavior that the server should exhibit when logging fields with this syntax. This may be overridden on a per-field basis.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a JsonLogFieldSyntaxResponse object into the model struct
func readJsonLogFieldSyntaxResponseDataSource(ctx context.Context, r *client.JsonLogFieldSyntaxResponse, state *logFieldSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("json")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IncludedSensitiveField = internaltypes.GetStringSet(r.IncludedSensitiveField)
	state.ExcludedSensitiveField = internaltypes.GetStringSet(r.ExcludedSensitiveField)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldSyntaxDefaultBehaviorProp(r.DefaultBehavior), false)
}

// Read a AttributeBasedLogFieldSyntaxResponse object into the model struct
func readAttributeBasedLogFieldSyntaxResponseDataSource(ctx context.Context, r *client.AttributeBasedLogFieldSyntaxResponse, state *logFieldSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("attribute-based")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.IncludedSensitiveAttribute = internaltypes.GetStringSet(r.IncludedSensitiveAttribute)
	state.ExcludedSensitiveAttribute = internaltypes.GetStringSet(r.ExcludedSensitiveAttribute)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldSyntaxDefaultBehaviorProp(r.DefaultBehavior), false)
}

// Read a GenericLogFieldSyntaxResponse object into the model struct
func readGenericLogFieldSyntaxResponseDataSource(ctx context.Context, r *client.GenericLogFieldSyntaxResponse, state *logFieldSyntaxDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("generic")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.DefaultBehavior = internaltypes.StringTypeOrNil(
		client.StringPointerEnumlogFieldSyntaxDefaultBehaviorProp(r.DefaultBehavior), false)
}

// Read resource information
func (r *logFieldSyntaxDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state logFieldSyntaxDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.LogFieldSyntaxAPI.GetLogFieldSyntax(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Log Field Syntax", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	if readResponse.JsonLogFieldSyntaxResponse != nil {
		readJsonLogFieldSyntaxResponseDataSource(ctx, readResponse.JsonLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.AttributeBasedLogFieldSyntaxResponse != nil {
		readAttributeBasedLogFieldSyntaxResponseDataSource(ctx, readResponse.AttributeBasedLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}
	if readResponse.GenericLogFieldSyntaxResponse != nil {
		readGenericLogFieldSyntaxResponseDataSource(ctx, readResponse.GenericLogFieldSyntaxResponse, &state, &resp.Diagnostics)
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
