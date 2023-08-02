package delegatedadmincorrelatedrestresource

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
	_ datasource.DataSource              = &delegatedAdminCorrelatedRestResourceDataSource{}
	_ datasource.DataSourceWithConfigure = &delegatedAdminCorrelatedRestResourceDataSource{}
)

// Create a Delegated Admin Correlated Rest Resource data source
func NewDelegatedAdminCorrelatedRestResourceDataSource() datasource.DataSource {
	return &delegatedAdminCorrelatedRestResourceDataSource{}
}

// delegatedAdminCorrelatedRestResourceDataSource is the datasource implementation.
type delegatedAdminCorrelatedRestResourceDataSource struct {
	providerConfig internaltypes.ProviderConfiguration
	apiClient      *client.APIClient
}

// Metadata returns the data source type name.
func (r *delegatedAdminCorrelatedRestResourceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_delegated_admin_correlated_rest_resource"
}

// Configure adds the provider configured client to the data source.
func (r *delegatedAdminCorrelatedRestResourceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerCfg := req.ProviderData.(internaltypes.ResourceConfiguration)
	r.providerConfig = providerCfg.ProviderConfig
	r.apiClient = providerCfg.ApiClientV9300
}

type delegatedAdminCorrelatedRestResourceDataSourceModel struct {
	Id                                        types.String `tfsdk:"id"`
	Name                                      types.String `tfsdk:"name"`
	RestResourceTypeName                      types.String `tfsdk:"rest_resource_type_name"`
	DisplayName                               types.String `tfsdk:"display_name"`
	CorrelatedRESTResource                    types.String `tfsdk:"correlated_rest_resource"`
	PrimaryRESTResourceCorrelationAttribute   types.String `tfsdk:"primary_rest_resource_correlation_attribute"`
	SecondaryRESTResourceCorrelationAttribute types.String `tfsdk:"secondary_rest_resource_correlation_attribute"`
	UseSecondaryValueForLinking               types.Bool   `tfsdk:"use_secondary_value_for_linking"`
}

// GetSchema defines the schema for the datasource.
func (r *delegatedAdminCorrelatedRestResourceDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	schemaDef := schema.Schema{
		Description: "Describes a Delegated Admin Correlated Rest Resource.",
		Attributes: map[string]schema.Attribute{
			"rest_resource_type_name": schema.StringAttribute{
				Description: "Name of the parent REST Resource Type",
				Required:    true,
			},
			"display_name": schema.StringAttribute{
				Description: "A human readable display name for this Delegated Admin Correlated REST Resource.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"correlated_rest_resource": schema.StringAttribute{
				Description: "The REST Resource Type that will be linked to this REST Resource Type.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"primary_rest_resource_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the parent REST Resource Type whose value will be used to match objects in the Delegated Admin Correlated REST Resource. This attribute must be writeable when use-secondary-value-for-linking is enabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"secondary_rest_resource_correlation_attribute": schema.StringAttribute{
				Description: "The LDAP attribute from the Delegated Admin Correlated REST Resource whose value will be matched with the primary-rest-resource-correlation-attribute. This attribute must be writeable when use-secondary-value-for-linking is disabled.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
			"use_secondary_value_for_linking": schema.BoolAttribute{
				Description: "Indicates whether links should be created using the secondary correlation attribute value.",
				Required:    false,
				Optional:    false,
				Computed:    true,
			},
		},
	}
	config.AddCommonDataSourceSchema(&schemaDef, true)
	resp.Schema = schemaDef
}

// Read a DelegatedAdminCorrelatedRestResourceResponse object into the model struct
func readDelegatedAdminCorrelatedRestResourceResponseDataSource(ctx context.Context, r *client.DelegatedAdminCorrelatedRestResourceResponse, state *delegatedAdminCorrelatedRestResourceDataSourceModel, diagnostics *diag.Diagnostics) {
	state.Id = types.StringValue(r.Id)
	state.Name = types.StringValue(r.Id)
	state.DisplayName = types.StringValue(r.DisplayName)
	state.CorrelatedRESTResource = types.StringValue(r.CorrelatedRESTResource)
	state.PrimaryRESTResourceCorrelationAttribute = types.StringValue(r.PrimaryRESTResourceCorrelationAttribute)
	state.SecondaryRESTResourceCorrelationAttribute = types.StringValue(r.SecondaryRESTResourceCorrelationAttribute)
	state.UseSecondaryValueForLinking = internaltypes.BoolTypeOrNil(r.UseSecondaryValueForLinking)
}

// Read resource information
func (r *delegatedAdminCorrelatedRestResourceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state delegatedAdminCorrelatedRestResourceDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	readResponse, httpResp, err := r.apiClient.DelegatedAdminCorrelatedRestResourceApi.GetDelegatedAdminCorrelatedRestResource(
		config.ProviderBasicAuthContext(ctx, r.providerConfig), state.Name.ValueString(), state.RestResourceTypeName.ValueString()).Execute()
	if err != nil {
		config.ReportHttpError(ctx, &resp.Diagnostics, "An error occurred while getting the Delegated Admin Correlated Rest Resource", err, httpResp)
		return
	}

	// Log response JSON
	responseJson, err := readResponse.MarshalJSON()
	if err == nil {
		tflog.Debug(ctx, "Read response: "+string(responseJson))
	}

	// Read the response into the state
	readDelegatedAdminCorrelatedRestResourceResponseDataSource(ctx, readResponse, &state, &resp.Diagnostics)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
