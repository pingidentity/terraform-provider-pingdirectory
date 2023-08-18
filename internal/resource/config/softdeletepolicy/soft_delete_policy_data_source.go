package softdeletepolicy

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
	_ datasource.DataSource              = &softDeletePolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &softDeletePolicyDataSource{}
)

// Create a Soft Delete Policy data source
func NewSoftDeletePolicyDataSource() datasource.DataSource {
	return &softDeletePolicyDataSource{}
}

// softDeletePolicyDataSource is the datasource implementation.
type softDeletePolicyDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *softDeletePolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_soft_delete_policy"
}

// Configure adds the provider configured client to the data source.
func (r *softDeletePolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClient
}

type softDeletePolicyDataSourceModel struct {
	Id                               types.String `tfsdk:"id"`
	Name                             types.String `tfsdk:"name"`
	Type                             types.String `tfsdk:"type"`
	Description                      types.String `tfsdk:"description"`
	AutoSoftDeleteConnectionCriteria types.String `tfsdk:"auto_soft_delete_connection_criteria"`
	AutoSoftDeleteRequestCriteria    types.String `tfsdk:"auto_soft_delete_request_criteria"`
	SoftDeleteRetentionTime          types.String `tfsdk:"soft_delete_retention_time"`
	SoftDeleteRetainNumberOfEntries  types.Int64  `tfsdk:"soft_delete_retain_number_of_entries"`
}

// GetSchema defines the schema for the datasource.
func (r *softDeletePolicyDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Soft Delete Policy.",
		Attributes: map[string]schema.Attribute{
			"type": schema.StringAttribute{
				Description: "The type of Soft Delete Policy resource. Options are ['soft-delete-policy']",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				Description: "A description for this Soft Delete Policy",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auto_soft_delete_connection_criteria": schema.StringAttribute{
				Description: "Connection criteria used to automatically identify a delete operation for processing as a soft delete request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"auto_soft_delete_request_criteria": schema.StringAttribute{
				Description: "Request criteria used to automatically identify a delete operation for processing as a soft delete request.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"soft_delete_retention_time": schema.StringAttribute{
				Description: "Specifies the maximum length of time that soft delete entries are retained before they are eligible to purged automatically.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"soft_delete_retain_number_of_entries": schema.Int64Attribute{
				Description: "Specifies the number of soft deleted entries to retain before the oldest entries are purged.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a SoftDeletePolicyResponse object into the model struct
func readSoftDeletePolicyResponseDataSource(ctx context.Context, r *client.SoftDeletePolicyResponse, state *softDeletePolicyDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Type = types.StringValue("soft-delete-policy")
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.Description = internaltypes.StringTypeOrNil(r.Description, false)
	state.AutoSoftDeleteConnectionCriteria = internaltypes.StringTypeOrNil(r.AutoSoftDeleteConnectionCriteria, false)
	state.AutoSoftDeleteRequestCriteria = internaltypes.StringTypeOrNil(r.AutoSoftDeleteRequestCriteria, false)
	state.SoftDeleteRetentionTime = internaltypes.StringTypeOrNil(r.SoftDeleteRetentionTime, false)
	state.SoftDeleteRetainNumberOfEntries = internaltypes.Int64TypeOrNil(r.SoftDeleteRetainNumberOfEntries)
}

// Read resource information
func (r *softDeletePolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state softDeletePolicyDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.SoftDeletePolicyApi.GetSoftDeletePolicy(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Soft Delete Policy", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readSoftDeletePolicyResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
